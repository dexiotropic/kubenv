package cli

import (
	"bytes"
	"io"
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
	err = Run([]string{"render", "--dotenv"}, strings.NewReader("msg: {{ env.GREETING }}\n"), &stdout, ioDiscard{}, []string{})
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
	err := Run([]string{"render", "--dotenv-file", envPath}, strings.NewReader("msg: {{ env.GREETING }}\n"), &stdout, ioDiscard{}, []string{})
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
		[]string{"render", "--dotenv-file", envPath, "--set", "GREETING=cli"},
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
		[]string{"render", "--dotenv-file", envPath, "--ignore-process-env"},
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

func TestRunApplyInvokesKubectlWithRenderedManifest(t *testing.T) {
	original := kubectlApply
	t.Cleanup(func() {
		kubectlApply = original
	})

	var gotInput []byte
	var gotArgs []string
	kubectlApply = func(input []byte, _, _ io.Writer, args []string) error {
		gotInput = append([]byte(nil), input...)
		gotArgs = append([]string(nil), args...)
		return nil
	}

	var stdout bytes.Buffer
	err := Run(
		[]string{"apply", "--set", "GREETING=hello"},
		strings.NewReader("msg: {{ env.GREETING }}\n"),
		&stdout,
		ioDiscard{},
		nil,
	)
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if got := string(gotInput); got != "msg: hello\n" {
		t.Fatalf("unexpected kubectl stdin: %q", got)
	}
	if len(gotArgs) != 0 {
		t.Fatalf("unexpected kubectl args: %v", gotArgs)
	}
}

func TestRunApplyPassesArgsAfterSeparator(t *testing.T) {
	original := kubectlApply
	t.Cleanup(func() {
		kubectlApply = original
	})

	var gotArgs []string
	kubectlApply = func(_ []byte, _, _ io.Writer, args []string) error {
		gotArgs = append([]string(nil), args...)
		return nil
	}

	var stdout bytes.Buffer
	err := Run(
		[]string{"apply", "--set", "GREETING=hello", "--", "--namespace", "demo", "--server-side"},
		strings.NewReader("msg: {{ env.GREETING }}\n"),
		&stdout,
		ioDiscard{},
		nil,
	)
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	expected := []string{"--namespace", "demo", "--server-side"}
	if strings.Join(gotArgs, ",") != strings.Join(expected, ",") {
		t.Fatalf("unexpected kubectl args: %v", gotArgs)
	}
}

func TestRunKubectlPluginSupportsFlagsBeforeApply(t *testing.T) {
	original := kubectlApply
	t.Cleanup(func() {
		kubectlApply = original
	})

	var gotInput []byte
	var gotArgs []string
	kubectlApply = func(input []byte, _, _ io.Writer, args []string) error {
		gotInput = append([]byte(nil), input...)
		gotArgs = append([]string(nil), args...)
		return nil
	}

	var stdout bytes.Buffer
	err := RunKubectlPlugin(
		[]string{"--dotenv-file", "testdata.env", "--set", "GREETING=hello", "apply", "--namespace", "demo"},
		strings.NewReader("msg: {{ env.GREETING }}\n"),
		&stdout,
		ioDiscard{},
		nil,
	)
	if err == nil {
		t.Fatal("expected error for missing env file in plugin mode")
	}

	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env.dev")
	writeFile(t, envPath, "NAME=file\n")

	err = RunKubectlPlugin(
		[]string{"--dotenv-file", envPath, "--set", "GREETING=hello", "apply", "--namespace", "demo"},
		strings.NewReader("msg: {{ env.GREETING }} {{ env.NAME }}\n"),
		&stdout,
		ioDiscard{},
		nil,
	)
	if err != nil {
		t.Fatalf("RunKubectlPlugin returned error: %v", err)
	}

	if got := string(gotInput); got != "msg: hello file\n" {
		t.Fatalf("unexpected kubectl stdin: %q", got)
	}

	expected := []string{"--namespace", "demo"}
	if strings.Join(gotArgs, ",") != strings.Join(expected, ",") {
		t.Fatalf("unexpected kubectl args: %v", gotArgs)
	}
}

func TestRunKubectlPluginPassesThroughDirectSyntax(t *testing.T) {
	original := kubectlApply
	t.Cleanup(func() {
		kubectlApply = original
	})

	var gotInput []byte
	kubectlApply = func(input []byte, _, _ io.Writer, _ []string) error {
		gotInput = append([]byte(nil), input...)
		return nil
	}

	err := RunKubectlPlugin(
		[]string{"apply", "--set", "GREETING=hello"},
		strings.NewReader("msg: {{ env.GREETING }}\n"),
		ioDiscard{},
		ioDiscard{},
		nil,
	)
	if err != nil {
		t.Fatalf("RunKubectlPlugin returned error: %v", err)
	}

	if got := string(gotInput); got != "msg: hello\n" {
		t.Fatalf("unexpected kubectl stdin: %q", got)
	}
}

func TestRunRenderFailsWhenEnvFlagsConflict(t *testing.T) {
	var stdout bytes.Buffer
	err := Run(
		[]string{"render", "--dotenv", "--dotenv-file", ".env.dev"},
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
