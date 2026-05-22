package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

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

	filePath := fs.String("f", "", "manifest file to render")
	if err := fs.Parse(args); err != nil {
		return err
	}

	var input []byte
	var err error
	if *filePath == "" {
		input, err = io.ReadAll(stdin)
	} else {
		input, err = os.ReadFile(*filePath)
	}
	if err != nil {
		return err
	}

	output, err := render.Strict(input, render.FromEnviron(environ))
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
