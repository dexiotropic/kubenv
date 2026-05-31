package e2etest

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

type Result struct {
	Stdout string
	Stderr string
}

func RepoRoot(t *testing.T) string {
	t.Helper()

	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}

	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", "..", ".."))
}

func RunGo(t *testing.T, stdin string, env []string, args ...string) (Result, error) {
	t.Helper()

	cmd := exec.Command("go", append([]string{"run"}, args...)...)
	cmd.Dir = RepoRoot(t)
	cmd.Env = append(os.Environ(), env...)
	cmd.Stdin = bytes.NewBufferString(stdin)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	return Result{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}, err
}
