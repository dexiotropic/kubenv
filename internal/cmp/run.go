package cmp

import (
	"io"

	"github.com/dexiotropic/kubenv/internal/render"
)

func Run(stdin io.Reader, stdout, _ io.Writer, environ []string) error {
	input, err := io.ReadAll(stdin)
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
