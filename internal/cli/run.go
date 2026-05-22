package cli

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/dexiotropic/kubenv/internal/render"
)

func Run(args []string, stdin io.Reader, stdout, stderr io.Writer, environ []string) error {
	if len(args) == 0 {
		return usage(stderr)
	}

	switch args[0] {
	case "render":
		return runRender(args[1:], stdin, stdout, environ)
	case "version":
		_, err := fmt.Fprintln(stdout, "kubenv dev")
		return err
	default:
		return usage(stderr)
	}
}

func runRender(args []string, stdin io.Reader, stdout io.Writer, environ []string) error {
	fs := flag.NewFlagSet("render", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	var filePaths stringSliceFlag
	fs.Var(&filePaths, "f", "manifest file to render; may be repeated")
	useDotenv := fs.Bool("env", false, "load variables from .env")
	envFile := fs.String("env-file", "", "load variables from a specific dotenv file")
	ignoreProcessEnv := fs.Bool("ignore-process-env", false, "skip reading variables from the process environment")
	var setValues stringSliceFlag
	fs.Var(&setValues, "set", "override a variable with KEY=VALUE")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *useDotenv && *envFile != "" {
		return errors.New("--env and --env-file cannot be used together")
	}

	input, err := readInputs(stdin, filePaths)
	if err != nil {
		return err
	}

	vars, err := loadVariables(environ, *useDotenv, *envFile, *ignoreProcessEnv, setValues)
	if err != nil {
		return err
	}

	output, err := render.Strict(input, vars)
	if err != nil {
		return err
	}

	_, err = stdout.Write(output)
	return err
}

func usage(w io.Writer) error {
	_, _ = fmt.Fprintln(w, "usage: kubenv <render|version> [flags]")
	return errors.New("invalid command")
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
