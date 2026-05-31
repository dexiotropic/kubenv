package version

import (
	"strings"
	"testing"
)

func TestStringIncludesVersionCommitAndDate(t *testing.T) {
	got := String()

	for _, want := range []string{version, commit, date} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q to contain %q", got, want)
		}
	}
}

func TestCommitReturnsBuildCommit(t *testing.T) {
	if got := Commit(); got != commit {
		t.Fatalf("unexpected commit: %q", got)
	}
}

func TestDateReturnsBuildDate(t *testing.T) {
	if got := Date(); got != date {
		t.Fatalf("unexpected date: %q", got)
	}
}
