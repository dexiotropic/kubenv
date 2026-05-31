package kubenv_test

import (
	"strings"
	"testing"

	"github.com/dexiotropic/kubenv/e2e/internal/e2etest"
)

func TestRenderExplicitPlaceholders(t *testing.T) {
	result, err := e2etest.RunGo(
		t,
		"msg: {{ env.GREETING }} {{ env.TARGET }}\n",
		nil,
		"./cmd/kubenv",
		"render",
		"--set", "GREETING=hello",
		"--set", "TARGET=world",
	)
	if err != nil {
		t.Fatalf("go run kubenv render: %v\nstderr:\n%s", err, result.Stderr)
	}

	if got := result.Stdout; got != "msg: hello world\n" {
		t.Fatalf("unexpected stdout: %q", got)
	}
}

func TestRenderShellStylePlaceholders(t *testing.T) {
	result, err := e2etest.RunGo(
		t,
		"msg: $GREETING ${TARGET}\n",
		nil,
		"./cmd/kubenv",
		"render",
		"--shell-style",
		"--set", "GREETING=hello",
		"--set", "TARGET=world",
	)
	if err != nil {
		t.Fatalf("go run kubenv render --shell-style: %v\nstderr:\n%s", err, result.Stderr)
	}

	if got := result.Stdout; got != "msg: hello world\n" {
		t.Fatalf("unexpected stdout: %q", got)
	}
}

func TestHelpIncludesShellStyle(t *testing.T) {
	result, err := e2etest.RunGo(t, "", nil, "./cmd/kubenv", "render", "--help")
	if err != nil {
		t.Fatalf("go run kubenv render --help: %v\nstderr:\n%s", err, result.Stderr)
	}

	if !strings.Contains(result.Stdout, "--shell-style") {
		t.Fatalf("help output missing --shell-style: %q", result.Stdout)
	}
}
