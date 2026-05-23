package cli

import (
	"errors"
	"fmt"
	"io"
)

// RunKubectlPlugin is the entrypoint for the kubectl env plugin.
func RunKubectlPlugin(args []string, stdin io.Reader, stdout, stderr io.Writer, environ []string) error {
	if len(args) == 0 {
		return kubectlPluginUsage(stderr)
	}

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

func kubectlPluginUsage(w io.Writer) error {
	_, _ = fmt.Fprintln(w, "usage: kubectl env [kubenv flags] apply [kubectl apply flags]")
	_, _ = fmt.Fprintln(w, "   or: kubectl env <render|apply|version> [flags]")
	return errors.New("invalid command")
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
