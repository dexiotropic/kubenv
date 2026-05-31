package cli

import (
	"context"
	"fmt"
	"io"

	"github.com/dexiotropic/kubenv/internal/version"
	ucli "github.com/urfave/cli/v3"
)

// RunKubectlPlugin is the entrypoint for the kubectl kenv plugin.
func RunKubectlPlugin(args []string, stdin io.Reader, stdout, stderr io.Writer, environ []string) error {
	if usesImplicitApplySyntax(args) {
		command := newKubectlImplicitApplyCommand(stdin, stdout, stderr, environ)
		return command.Run(context.Background(), append([]string{"kubectl kenv"}, args...))
	}

	command := newKubectlPluginCommand(stdin, stdout, stderr, environ)
	return command.Run(context.Background(), append([]string{"kubectl kenv"}, args...))
}

func newKubectlPluginCommand(stdin io.Reader, stdout, stderr io.Writer, environ []string) *ucli.Command {
	return &ucli.Command{
		Name:  "kubectl kenv",
		Usage: "Render manifests through kubectl",
		UsageText: `kubectl kenv [kubenv flags] apply [kubectl apply flags]
kubectl kenv <render|apply|version> [flags]`,
		Description: "Use the first form to put kubenv flags before apply and pass the remaining arguments directly to kubectl apply.",
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
				UsageText: "kubectl kenv render [flags]",
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
				UsageText: "kubectl kenv apply [flags] -- [kubectl apply flags]",
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

func newKubectlImplicitApplyCommand(stdin io.Reader, stdout, stderr io.Writer, environ []string) *ucli.Command {
	stopOnFirstArg := 1

	return &ucli.Command{
		Name:         "kubectl kenv",
		Usage:        "Render manifests and run kubectl apply -f -",
		UsageText:    "kubectl kenv [kubenv flags] apply [kubectl apply flags]",
		Description:  "Arguments before apply are handled by kubenv. Arguments after apply are forwarded to kubectl apply.",
		Version:      version.String(),
		Reader:       stdin,
		Writer:       stdout,
		ErrWriter:    stderr,
		Flags:        renderFlags(),
		StopOnNthArg: &stopOnFirstArg,
		OnUsageError: func(_ context.Context, _ *ucli.Command, err error, _ bool) error {
			return err
		},
		Action: func(_ context.Context, cmd *ucli.Command) error {
			options, err := renderOptionsFromCommand(cmd)
			if err != nil {
				return err
			}

			args := options.extraArgs
			if len(args) == 0 || args[0] != "apply" {
				return ucli.ShowRootCommandHelp(cmd)
			}
			options.extraArgs = args[1:]

			output, err := renderWithOptions(options, stdin, environ)
			if err != nil {
				return err
			}

			return kubectlApply(output, stdout, stderr, options.extraArgs)
		},
	}
}

func isCommand(arg string) bool {
	switch arg {
	case "render", "apply", "version", "help":
		return true
	default:
		return false
	}
}

func usesImplicitApplySyntax(args []string) bool {
	if len(args) == 0 || isCommand(args[0]) || isHelpFlag(args[0]) {
		return false
	}

	return indexOf(args, "apply") != -1
}

func isHelpFlag(arg string) bool {
	return arg == "-h" || arg == "--help"
}

func indexOf(items []string, target string) int {
	for i, item := range items {
		if item == target {
			return i
		}
	}
	return -1
}
