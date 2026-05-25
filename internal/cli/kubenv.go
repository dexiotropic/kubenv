package cli

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/dexiotropic/kubenv/internal/version"
	ucli "github.com/urfave/cli/v3"
)

// Run is the entrypoint for the kubenv CLI.
func Run(args []string, stdin io.Reader, stdout, stderr io.Writer, environ []string) error {
	command := newKubenvCommand(stdin, stdout, stderr, environ)
	return command.Run(context.Background(), append([]string{"kubenv"}, args...))
}

func newKubenvCommand(stdin io.Reader, stdout, stderr io.Writer, environ []string) *ucli.Command {
	return &ucli.Command{
		Name:        "kubenv",
		Usage:       "Render Kubernetes manifests with strict variable substitution",
		Description: "kubenv renders manifests with strict {{ env.NAME }} substitution.",
		Version:     version.String(),
		Reader:      stdin,
		Writer:      stdout,
		ErrWriter:   stderr,
		OnUsageError: func(_ context.Context, _ *ucli.Command, err error, _ bool) error {
			return err
		},
		Action: func(_ context.Context, cmd *ucli.Command) error {
			return ucli.ShowRootCommandHelp(cmd)
		},
		Commands: []*ucli.Command{
			{
				Name:      "render",
				Usage:     "Render manifests to stdout",
				UsageText: "kubenv render [flags]",
				Flags:     renderFlags(),
				Action: func(_ context.Context, cmd *ucli.Command) error {
					options, err := renderOptionsFromCommand(cmd)
					if err != nil {
						return err
					}
					if len(options.extraArgs) != 0 {
						return fmt.Errorf("render does not accept positional arguments: %s", stringsJoin(options.extraArgs))
					}

					output, err := renderWithOptions(options, stdin, environ)
					if err != nil {
						return err
					}

					_, err = stdout.Write(output)
					return err
				},
			},
			{
				Name:      "apply",
				Usage:     "Render manifests and run kubectl apply -f -",
				UsageText: "kubenv apply [flags] -- [kubectl apply flags]",
				Flags:     renderFlags(),
				Action: func(_ context.Context, cmd *ucli.Command) error {
					options, err := renderOptionsFromCommand(cmd)
					if err != nil {
						return err
					}

					output, err := renderWithOptions(options, stdin, environ)
					if err != nil {
						return err
					}

					return kubectlApply(output, stdout, stderr, options.extraArgs)
				},
			},
			{
				Name:  "version",
				Usage: "Print version information",
				Action: func(_ context.Context, _ *ucli.Command) error {
					_, err := fmt.Fprintln(stdout, version.String())
					return err
				},
			},
		},
	}
}

func stringsJoin(values []string) string {
	return strings.Join(values, " ")
}
