package cli

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/dexiotropic/kubenv/internal/render"
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

// Entry point for the kubenv CLI
func Run(args []string, stdin io.Reader, stdout, stderr io.Writer, environ []string) error {
	if len(args) == 0 {
		return usage(stderr)
	}

	switch args[0] {
	case "render":
		return runRender(args[1:], stdin, stdout, environ)
	case "apply":
		return runApply(args[1:], stdin, stdout, stderr, environ)
	case "version":
		_, err := fmt.Fprintln(stdout, "kubenv dev")
		return err
	default:
		return usage(stderr)
	}
}

// Entry point for the kubectl plugin
func RunKubectlPlugin(args []string, stdin io.Reader, stdout, stderr io.Writer, environ []string) error {
	if len(args) == 0 {
		return kubectlPluginUsage(stderr)
	}

	// If the first argument is a known command, run it directly. Otherwise, look for "apply" to determine if we're in plugin mode.
	// For example, "kubectl kubenv render -f manifest.yaml" should run the render command
	if isCommand(args[0]) {
		return Run(args, stdin, stdout, stderr, environ)
	}

	applyIndex := indexOf(args, "apply")
	if applyIndex == -1 {
		return kubectlPluginUsage(stderr)
	}

	output, _, err := renderCommand(args[:applyIndex], stdin, environ)
	if err != nil {
		return err
	}

	return kubectlApply(output, stdout, stderr, args[applyIndex+1:])
}

func runRender(args []string, stdin io.Reader, stdout io.Writer, environ []string) error {
	output, _, err := renderCommand(args, stdin, environ)
	if err != nil {
		return err
	}

	_, err = stdout.Write(output)
	return err
}

func runApply(args []string, stdin io.Reader, stdout, stderr io.Writer, environ []string) error {
	output, kubectlArgs, err := renderCommand(args, stdin, environ)
	if err != nil {
		return err
	}

	return kubectlApply(output, stdout, stderr, kubectlArgs)
}

func renderCommand(args []string, stdin io.Reader, environ []string) ([]byte, []string, error) {
	options, err := parseRenderOptions(args)
	if err != nil {
		return nil, nil, err
	}

	input, err := readInputs(stdin, options.filePaths)
	if err != nil {
		return nil, nil, err
	}

	vars, err := loadVariables(environ, options.useDotenv, options.envFile, options.ignoreProcessEnv, options.setValues)
	if err != nil {
		return nil, nil, err
	}

	output, err := render.Strict(input, vars)
	if err != nil {
		return nil, nil, err
	}

	return output, options.extraArgs, nil
}

func usage(w io.Writer) error {
	_, _ = fmt.Fprintln(w, "usage: kubenv <render|apply|version> [flags]")
	return errors.New("invalid command")
}

func kubectlPluginUsage(w io.Writer) error {
	_, _ = fmt.Fprintln(w, "usage: kubectl kubenv [kubenv flags] apply [kubectl apply flags]")
	_, _ = fmt.Fprintln(w, "   or: kubectl kubenv <render|apply|version> [flags]")
	return errors.New("invalid command")
}

type renderOptions struct {
	filePaths        []string
	useDotenv        bool
	envFile          string
	ignoreProcessEnv bool
	setValues        []string
	extraArgs        []string
}

func parseRenderOptions(args []string) (renderOptions, error) {
	fs := flag.NewFlagSet("render", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	var options renderOptions
	fs.Var((*stringSliceFlag)(&options.filePaths), "f", "manifest file to render; may be repeated")
	fs.BoolVar(&options.useDotenv, "dotenv", false, "load variables from .env")
	fs.StringVar(&options.envFile, "dotenv-file", "", "load variables from a specific dotenv file")
	fs.BoolVar(&options.ignoreProcessEnv, "ignore-process-env", false, "skip reading variables from the process environment")
	fs.Var((*stringSliceFlag)(&options.setValues), "set", "override a variable with KEY=VALUE")
	if err := fs.Parse(args); err != nil {
		return renderOptions{}, err
	}
	if options.useDotenv && options.envFile != "" {
		return renderOptions{}, errors.New("--env and --env-file cannot be used together")
	}

	options.extraArgs = fs.Args()
	return options, nil
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

		for key, value := range dotenvVars {
			vars[key] = value
		}
	}

	if !ignoreProcessEnv {
		for key, value := range render.FromEnviron(environ) {
			vars[key] = value
		}
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

type stringSliceFlag []string

func (f *stringSliceFlag) String() string {
	return strings.Join(*f, ",")
}

func (f *stringSliceFlag) Set(value string) error {
	*f = append(*f, value)
	return nil
}

func isCommand(arg string) bool {
	switch arg {
	case "render", "apply", "version":
		return true
	default:
		return false
	}
}

func indexOf(items []string, target string) int {
	for i, item := range items {
		if item == target {
			return i
		}
	}
	return -1
}
