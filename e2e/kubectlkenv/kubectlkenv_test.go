package kubectlkenv_test

import (
	"strings"
	"testing"

	"github.com/dexiotropic/kubenv/e2e/internal/e2etest"
)

func TestRenderSubcommand(t *testing.T) {
	result, err := e2etest.RunGo(
		t,
		"kind: ConfigMap\nmetadata:\n  name: {{ env.NAME }}\n",
		nil,
		"./cmd/kubectl-kenv",
		"render",
		"--set", "NAME=demo",
	)
	if err != nil {
		t.Fatalf("go run kubectl-kenv render: %v\nstderr:\n%s", err, result.Stderr)
	}

	if got := result.Stdout; got != "kind: ConfigMap\nmetadata:\n  name: demo\n" {
		t.Fatalf("unexpected stdout: %q", got)
	}
}

func TestRenderSubcommandSupportsShellStyle(t *testing.T) {
	result, err := e2etest.RunGo(
		t,
		"msg: $GREETING ${TARGET}\n",
		nil,
		"./cmd/kubectl-kenv",
		"render",
		"--shell-style",
		"--set", "GREETING=hello",
		"--set", "TARGET=world",
	)
	if err != nil {
		t.Fatalf("go run kubectl-kenv render --shell-style: %v\nstderr:\n%s", err, result.Stderr)
	}

	if got := result.Stdout; got != "msg: hello world\n" {
		t.Fatalf("unexpected stdout: %q", got)
	}
}

func TestHelpMentionsImplicitApplySyntax(t *testing.T) {
	result, err := e2etest.RunGo(t, "", nil, "./cmd/kubectl-kenv", "--help")
	if err != nil {
		t.Fatalf("go run kubectl-kenv --help: %v\nstderr:\n%s", err, result.Stderr)
	}

	if !strings.Contains(result.Stdout, "kubectl kenv [kubenv flags] apply [kubectl apply flags]") {
		t.Fatalf("unexpected help output: %q", result.Stdout)
	}
}
