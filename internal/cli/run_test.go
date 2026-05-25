package cli

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunShowsHelp(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := Run([]string{"--help"}, strings.NewReader(""), &stdout, &stderr, nil)
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	assertContains(t, stdout.String(), "render")
	assertContains(t, stdout.String(), "apply")
	assertContains(t, stdout.String(), "version")
	if got := stderr.String(); got != "" {
		t.Fatalf("unexpected stderr: %q", got)
	}
}

func TestRunRenderShowsHelp(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := Run([]string{"render", "--help"}, strings.NewReader(""), &stdout, &stderr, nil)
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	assertContains(t, stdout.String(), "--dotenv")
	assertContains(t, stdout.String(), "--dotenv-file")
	assertContains(t, stdout.String(), "--shell-style")
	assertContains(t, stdout.String(), "--set")
	if got := stderr.String(); got != "" {
		t.Fatalf("unexpected stderr: %q", got)
	}
}

func TestRunKubectlPluginShowsHelp(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := RunKubectlPlugin([]string{"--help"}, strings.NewReader(""), &stdout, &stderr, nil)
	if err != nil {
		t.Fatalf("RunKubectlPlugin returned error: %v", err)
	}

	assertContains(t, stdout.String(), "kubectl env [kubenv flags] apply [kubectl apply flags]")
	assertContains(t, stdout.String(), "render")
	assertContains(t, stdout.String(), "apply")
	if got := stderr.String(); got != "" {
		t.Fatalf("unexpected stderr: %q", got)
	}
}

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

func TestRunRenderSupportsShellStylePlaceholders(t *testing.T) {
	var stdout bytes.Buffer
	err := Run(
		[]string{"render", "--shell-style", "--set", "GREETING=hello", "--set", "TARGET=world"},
		strings.NewReader("msg: $GREETING ${TARGET}\n"),
		&stdout,
		ioDiscard{},
		nil,
	)
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if got := stdout.String(); got != "msg: hello world\n" {
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

func TestRunRenderCanReadDirectory(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "first.yaml"), "first: {{ env.GREETING }}\n")
	writeFile(t, filepath.Join(dir, "second.json"), "{\"second\":\"{{ env.NAME }}\"}\n")
	writeFile(t, filepath.Join(dir, "ignore.txt"), "ignored")
	subdir := filepath.Join(dir, "nested")
	if err := os.Mkdir(subdir, 0o755); err != nil {
		t.Fatalf("Mkdir: %v", err)
	}
	writeFile(t, filepath.Join(subdir, "third.yaml"), "third: {{ env.EXTRA }}\n")

	var stdout bytes.Buffer
	err := Run(
		[]string{"render", "-f", dir},
		strings.NewReader(""),
		&stdout,
		ioDiscard{},
		[]string{"GREETING=hello", "NAME=world", "EXTRA=ignored"},
	)
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if got := stdout.String(); got != "first: hello\n\n---\n{\"second\":\"world\"}\n" {
		t.Fatalf("unexpected output: %q", got)
	}
}

func TestRunRenderCanReadDirectoryRecursively(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "first.yaml"), "first: {{ env.GREETING }}\n")
	subdir := filepath.Join(dir, "nested")
	if err := os.Mkdir(subdir, 0o755); err != nil {
		t.Fatalf("Mkdir: %v", err)
	}
	writeFile(t, filepath.Join(subdir, "second.yaml"), "second: {{ env.NAME }}\n")

	var stdout bytes.Buffer
	err := Run(
		[]string{"render", "-f", dir, "--recursive"},
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

func TestRunRenderCanExpandGlobPatterns(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "a.yaml"), "a: {{ env.FIRST }}\n")
	writeFile(t, filepath.Join(dir, "b.yaml"), "b: {{ env.SECOND }}\n")
	writeFile(t, filepath.Join(dir, "c.txt"), "ignored")

	var stdout bytes.Buffer
	err := Run(
		[]string{"render", "-f", filepath.Join(dir, "*.yaml")},
		strings.NewReader(""),
		&stdout,
		ioDiscard{},
		[]string{"FIRST=hello", "SECOND=world"},
	)
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if got := stdout.String(); got != "a: hello\n\n---\nb: world\n" {
		t.Fatalf("unexpected output: %q", got)
	}
}

func TestRunRenderCanReadRemoteURL(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/manifest.yaml" {
			http.NotFound(w, r)
			return
		}
		_, _ = io.WriteString(w, "msg: {{ env.GREETING }}\n")
	}))
	t.Cleanup(server.Close)

	var stdout bytes.Buffer
	err := Run(
		[]string{"render", "-f", server.URL + "/manifest.yaml"},
		strings.NewReader(""),
		&stdout,
		ioDiscard{},
		[]string{"GREETING=hello"},
	)
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if got := stdout.String(); got != "msg: hello\n" {
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

func TestRunKubectlPluginSupportsRecursiveDirectoriesBeforeApply(t *testing.T) {
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

	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "first.yaml"), "first: {{ env.GREETING }}\n")
	subdir := filepath.Join(dir, "nested")
	if err := os.Mkdir(subdir, 0o755); err != nil {
		t.Fatalf("Mkdir: %v", err)
	}
	writeFile(t, filepath.Join(subdir, "second.yaml"), "second: {{ env.NAME }}\n")

	err := RunKubectlPlugin(
		[]string{"-f", dir, "--recursive", "apply", "--namespace", "demo"},
		strings.NewReader(""),
		ioDiscard{},
		ioDiscard{},
		[]string{"GREETING=hello", "NAME=world"},
	)
	if err != nil {
		t.Fatalf("RunKubectlPlugin returned error: %v", err)
	}

	if got := string(gotInput); got != "first: hello\n\n---\nsecond: world\n" {
		t.Fatalf("unexpected kubectl stdin: %q", got)
	}
	expected := []string{"--namespace", "demo"}
	if strings.Join(gotArgs, ",") != strings.Join(expected, ",") {
		t.Fatalf("unexpected kubectl args: %v", gotArgs)
	}
}

func TestRunKubectlPluginSupportsShellStyleBeforeApply(t *testing.T) {
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
		[]string{"--shell-style", "--set", "GREETING=hello", "apply"},
		strings.NewReader("msg: $GREETING\n"),
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

func assertContains(t *testing.T, got, want string) {
	t.Helper()
	if !strings.Contains(got, want) {
		t.Fatalf("expected %q to contain %q", got, want)
	}
}
