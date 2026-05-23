package cli

import (
	"errors"
	"fmt"
	"io"
)

// Run is the entrypoint for the kubenv CLI.
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

func usage(w io.Writer) error {
	_, _ = fmt.Fprintln(w, "usage: kubenv <render|apply|version> [flags]")
	return errors.New("invalid command")
}
