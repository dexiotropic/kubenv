package cli

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"maps"
	"os"
	"os/exec"
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
			Usage:   "manifest file to render; may be repeated",
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
		useDotenv:        cmd.Bool("dotenv"),
		envFile:          cmd.String("dotenv-file"),
		ignoreProcessEnv: cmd.Bool("ignore-process-env"),
		setValues:        cmd.StringSlice("set"),
		extraArgs:        cmd.Args().Slice(),
	}

	if options.useDotenv && options.envFile != "" {
		return renderOptions{}, errors.New("--dotenv and --dotenv-file cannot be used together")
	}

	return options, nil
}

func renderWithOptions(options renderOptions, stdin io.Reader, environ []string) ([]byte, error) {
	input, err := readInputs(stdin, options.filePaths)
	if err != nil {
		return nil, err
	}

	vars, err := loadVariables(environ, options.useDotenv, options.envFile, options.ignoreProcessEnv, options.setValues)
	if err != nil {
		return nil, err
	}

	return render.Strict(input, vars)
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

func readInputs(stdin io.Reader, filePaths []string) ([]byte, error) {
	if len(filePaths) == 0 {
		return io.ReadAll(stdin)
	}

	var documents [][]byte
	for _, filePath := range filePaths {
		input, err := os.ReadFile(filePath)
		if err != nil {
			return nil, err
		}
		documents = append(documents, input)
	}

	return bytes.Join(documents, []byte("\n---\n")), nil
}
