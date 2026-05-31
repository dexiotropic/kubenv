package argocdcmp_test

import (
	"testing"

	"github.com/dexiotropic/kubenv/e2e/internal/e2etest"
)

func TestCMPUsesExplicitParameters(t *testing.T) {
	result, err := e2etest.RunGo(
		t,
		"msg: {{ env.GREETING }}\n",
		[]string{`ARGOCD_APP_PARAMETERS=[{"name":"GREETING","string":"hello"}]`},
		"./cmd/kubenv-argocd-cmp",
	)
	if err != nil {
		t.Fatalf("go run kubenv-argocd-cmp: %v\nstderr:\n%s", err, result.Stderr)
	}

	if got := result.Stdout; got != "msg: hello\n" {
		t.Fatalf("unexpected stdout: %q", got)
	}
}

func TestCMPSupportsShellStyleOption(t *testing.T) {
	result, err := e2etest.RunGo(
		t,
		"msg: $GREETING ${TARGET}\n",
		[]string{`ARGOCD_APP_PARAMETERS=[{"name":"kubenv","map":{"shell-style":"true"}},{"name":"vars","map":{"GREETING":"hello","TARGET":"world"}}]`},
		"./cmd/kubenv-argocd-cmp",
	)
	if err != nil {
		t.Fatalf("go run kubenv-argocd-cmp shell-style: %v\nstderr:\n%s", err, result.Stderr)
	}

	if got := result.Stdout; got != "msg: hello world\n" {
		t.Fatalf("unexpected stdout: %q", got)
	}
}
