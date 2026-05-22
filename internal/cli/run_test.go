package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunRenderLoadsDefaultDotenv(t *testing.T) {
	t.Setenv("GREETING", "")
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, ".env"), "GREETING=hello\n")

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("Chdir: %v", err)
	}
	defer func() {
		_ = os.Chdir(cwd)
	}()

	var stdout bytes.Buffer
	err = Run([]string{"render", "--env"}, strings.NewReader("msg: {{ env.GREETING }}\n"), &stdout, ioDiscard{}, []string{})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if got := stdout.String(); got != "msg: hello\n" {
		t.Fatalf("unexpected output: %q", got)
	}
}

func TestRunRenderLoadsCustomEnvFile(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env.dev")
	writeFile(t, envPath, "GREETING=hello\n")

	var stdout bytes.Buffer
	err := Run([]string{"render", "--env-file", envPath}, strings.NewReader("msg: {{ env.GREETING }}\n"), &stdout, ioDiscard{}, []string{})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if got := stdout.String(); got != "msg: hello\n" {
		t.Fatalf("unexpected output: %q", got)
	}
}

func TestRunRenderPrecedenceCLIThenProcessThenFile(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env.dev")
	writeFile(t, envPath, "GREETING=file\n")

	var stdout bytes.Buffer
	err := Run(
		[]string{"render", "--env-file", envPath, "--set", "GREETING=cli"},
		strings.NewReader("msg: {{ env.GREETING }}\n"),
		&stdout,
		ioDiscard{},
		[]string{"GREETING=process"},
	)
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if got := stdout.String(); got != "msg: cli\n" {
		t.Fatalf("unexpected output: %q", got)
	}
}

func TestRunRenderCanReadMultipleFiles(t *testing.T) {
	dir := t.TempDir()
	firstPath := filepath.Join(dir, "first.yaml")
	secondPath := filepath.Join(dir, "second.yaml")
	writeFile(t, firstPath, "first: {{ env.GREETING }}\n")
	writeFile(t, secondPath, "second: {{ env.NAME }}\n")

	var stdout bytes.Buffer
	err := Run(
		[]string{"render", "-f", firstPath, "-f", secondPath},
		strings.NewReader(""),
		&stdout,
		ioDiscard{},
		[]string{"GREETING=hello", "NAME=world"},
	)
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if got := stdout.String(); got != "first: hello\n\n---\nsecond: world\n" {
		t.Fatalf("unexpected output: %q", got)
	}
}

func TestRunRenderCanIgnoreProcessEnv(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env.dev")
	writeFile(t, envPath, "GREETING=file\n")

	var stdout bytes.Buffer
	err := Run(
		[]string{"render", "--env-file", envPath, "--ignore-process-env"},
		strings.NewReader("msg: {{ env.GREETING }}\n"),
		&stdout,
		ioDiscard{},
		[]string{"GREETING=process"},
	)
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if got := stdout.String(); got != "msg: file\n" {
		t.Fatalf("unexpected output: %q", got)
	}
}

func TestRunRenderFailsWhenEnvFlagsConflict(t *testing.T) {
	var stdout bytes.Buffer
	err := Run(
		[]string{"render", "--env", "--env-file", ".env.dev"},
		strings.NewReader("msg: {{ env.GREETING }}\n"),
		&stdout,
		ioDiscard{},
		nil,
	)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRunRenderFailsOnInvalidSetSyntax(t *testing.T) {
	var stdout bytes.Buffer
	err := Run(
		[]string{"render", "--set", "GREETING"},
		strings.NewReader("msg: {{ env.GREETING }}\n"),
		&stdout,
		ioDiscard{},
		nil,
	)
	if err == nil {
		t.Fatal("expected error")
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile(%q): %v", path, err)
	}
}

type ioDiscard struct{}

func (ioDiscard) Write(p []byte) (int, error) {
	return len(p), nil
}
