package e2etest

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestRepoRootPointsAtRepository(t *testing.T) {
	root := RepoRoot(t)

	if !filepath.IsAbs(root) {
		t.Fatalf("expected absolute repo root, got %q", root)
	}
	if filepath.Base(root) != "kubenv" {
		t.Fatalf("unexpected repo root: %q", root)
	}
}

func TestRunGoExecutesCommandFromRepoRoot(t *testing.T) {
	result, err := RunGo(t, "", nil, "./cmd/kubenv", "--help")
	if err != nil {
		t.Fatalf("RunGo returned error: %v\nstderr:\n%s", err, result.Stderr)
	}

	if !strings.Contains(result.Stdout, "Render manifests to stdout") {
		t.Fatalf("unexpected stdout: %q", result.Stdout)
	}
}
