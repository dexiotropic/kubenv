package cli

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"maps"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/dexiotropic/kubenv/internal/render"
	ucli "github.com/urfave/cli/v3"
)

var kubectlApply = func(input []byte, stdout, stderr io.Writer, args []string) error {
	commandArgs := append([]string{"apply"}, args...)
	commandArgs = append(commandArgs, "-f", "-")

	cmd := exec.Command("kubectl", commandArgs...)
	cmd.Stdin = bytes.NewReader(input)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	return cmd.Run()
}

type renderOptions struct {
	filePaths        []string
	recursive        bool
	style            render.Style
	useDotenv        bool
	envFile          string
	ignoreProcessEnv bool
	setValues        []string
	extraArgs        []string
}

func renderFlags() []ucli.Flag {
	return []ucli.Flag{
		&ucli.StringSliceFlag{
			Name:    "file",
			Aliases: []string{"f"},
			Usage:   "manifest source to render; file, directory, glob, or URL; may be repeated",
		},
		&ucli.BoolFlag{
			Name:    "recursive",
			Aliases: []string{"R"},
			Usage:   "walk directories passed with -f recursively",
		},
		&ucli.BoolFlag{
			Name:  "shell-style",
			Usage: "render shell-style placeholders like $VAR and ${VAR} instead of {{ env.NAME }}",
		},
		&ucli.BoolFlag{
			Name:  "dotenv",
			Usage: "load variables from .env",
		},
		&ucli.StringFlag{
			Name:  "dotenv-file",
			Usage: "load variables from a specific dotenv file",
		},
		&ucli.BoolFlag{
			Name:  "ignore-process-env",
			Usage: "skip reading variables from the process environment",
		},
		&ucli.StringSliceFlag{
			Name:  "set",
			Usage: "override a variable with KEY=VALUE",
		},
	}
}

func renderOptionsFromCommand(cmd *ucli.Command) (renderOptions, error) {
	options := renderOptions{
		filePaths:        cmd.StringSlice("file"),
		recursive:        cmd.Bool("recursive"),
		style:            render.StyleExplicit,
		useDotenv:        cmd.Bool("dotenv"),
		envFile:          cmd.String("dotenv-file"),
		ignoreProcessEnv: cmd.Bool("ignore-process-env"),
		setValues:        cmd.StringSlice("set"),
		extraArgs:        cmd.Args().Slice(),
	}
	if cmd.Bool("shell-style") {
		options.style = render.StyleShell
	}

	if options.useDotenv && options.envFile != "" {
		return renderOptions{}, errors.New("--dotenv and --dotenv-file cannot be used together")
	}

	return options, nil
}

func renderWithOptions(options renderOptions, stdin io.Reader, environ []string) ([]byte, error) {
	input, err := readInputs(stdin, options.filePaths, options.recursive)
	if err != nil {
		return nil, err
	}

	vars, err := loadVariables(environ, options.useDotenv, options.envFile, options.ignoreProcessEnv, options.setValues)
	if err != nil {
		return nil, err
	}

	return render.StrictWithStyle(input, vars, options.style)
}

func loadVariables(environ []string, useDotenv bool, envFile string, ignoreProcessEnv bool, setValues []string) (map[string]string, error) {
	vars := map[string]string{}

	if useDotenv || envFile != "" {
		path := envFile
		if path == "" {
			path = ".env"
		}

		dotenvVars, err := parseDotenvFile(path)
		if err != nil {
			return nil, err
		}

		maps.Copy(vars, dotenvVars)
	}

	if !ignoreProcessEnv {
		maps.Copy(vars, render.FromEnviron(environ))
	}

	for _, item := range setValues {
		key, value, ok := strings.Cut(item, "=")
		if !ok || key == "" {
			return nil, fmt.Errorf("invalid --set value %q: expected KEY=VALUE", item)
		}
		vars[key] = value
	}

	return vars, nil
}

var httpClient = http.DefaultClient

func readInputs(stdin io.Reader, filePaths []string, recursive bool) ([]byte, error) {
	if len(filePaths) == 0 {
		return io.ReadAll(stdin)
	}

	resolvedPaths, err := resolveInputPaths(filePaths, recursive)
	if err != nil {
		return nil, err
	}

	var documents [][]byte
	for _, filePath := range resolvedPaths {
		input, err := readInputPath(filePath)
		if err != nil {
			return nil, err
		}
		documents = append(documents, input)
	}

	return bytes.Join(documents, []byte("\n---\n")), nil
}

func resolveInputPaths(filePaths []string, recursive bool) ([]string, error) {
	var resolved []string
	for _, filePath := range filePaths {
		paths, err := resolveInputPath(filePath, recursive)
		if err != nil {
			return nil, err
		}
		resolved = append(resolved, paths...)
	}

	return resolved, nil
}

func resolveInputPath(filePath string, recursive bool) ([]string, error) {
	if isRemoteURL(filePath) {
		return []string{filePath}, nil
	}

	if strings.ContainsAny(filePath, "*?[") {
		matches, err := filepath.Glob(filePath)
		if err != nil {
			return nil, fmt.Errorf("expand glob %q: %w", filePath, err)
		}
		sort.Strings(matches)

		var resolved []string
		for _, match := range matches {
			paths, err := resolveInputPath(match, recursive)
			if err != nil {
				return nil, err
			}
			resolved = append(resolved, paths...)
		}
		return resolved, nil
	}

	info, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return []string{filePath}, nil
	}

	return resolveDirectoryInputs(filePath, recursive)
}

func resolveDirectoryInputs(root string, recursive bool) ([]string, error) {
	var matches []string

	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			if !recursive && path != root {
				return filepath.SkipDir
			}
			return nil
		}
		if isManifestFile(path) {
			matches = append(matches, filepath.Clean(path))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Strings(matches)
	return matches, nil
}

func isManifestFile(path string) bool {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".yaml", ".yml", ".json":
		return true
	default:
		return false
	}
}

func readInputPath(path string) ([]byte, error) {
	if isRemoteURL(path) {
		return readRemoteInput(path)
	}

	return os.ReadFile(path)
}

func readRemoteInput(rawURL string) ([]byte, error) {
	response, err := httpClient.Get(rawURL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("fetch %q: unexpected status %s", rawURL, response.Status)
	}

	return io.ReadAll(response.Body)
}

func isRemoteURL(value string) bool {
	parsed, err := url.Parse(value)
	if err != nil {
		return false
	}

	return parsed.Scheme == "http" || parsed.Scheme == "https"
}
