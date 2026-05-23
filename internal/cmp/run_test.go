package cmp

import (
	"bytes"
	"testing"
)

func TestRunUsesCMPStringParameters(t *testing.T) {
	var stdout bytes.Buffer
	err := Run(
		bytes.NewBufferString("msg: {{ env.GREETING }}\n"),
		&stdout,
		bytes.NewBuffer(nil),
		[]string{`ARGOCD_APP_PARAMETERS=[{"name":"GREETING","string":"hello"}]`},
	)
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if got := stdout.String(); got != "msg: hello\n" {
		t.Fatalf("unexpected output: %q", got)
	}
}

func TestRunUsesCMPMapParameters(t *testing.T) {
	var stdout bytes.Buffer
	err := Run(
		bytes.NewBufferString("msg: {{ env.GREETING }} {{ env.NAME }}\n"),
		&stdout,
		bytes.NewBuffer(nil),
		[]string{`ARGOCD_APP_PARAMETERS=[{"name":"vars","map":{"GREETING":"hello","NAME":"world"}}]`},
	)
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if got := stdout.String(); got != "msg: hello world\n" {
		t.Fatalf("unexpected output: %q", got)
	}
}

func TestRunCMPParametersOverrideProcessEnv(t *testing.T) {
	var stdout bytes.Buffer
	err := Run(
		bytes.NewBufferString("msg: {{ env.GREETING }}\n"),
		&stdout,
		bytes.NewBuffer(nil),
		[]string{
			"GREETING=process",
			`ARGOCD_APP_PARAMETERS=[{"name":"GREETING","string":"parameter"}]`,
		},
	)
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if got := stdout.String(); got != "msg: parameter\n" {
		t.Fatalf("unexpected output: %q", got)
	}
}

func TestRunUsesPrefixedPluginEnv(t *testing.T) {
	var stdout bytes.Buffer
	err := Run(
		bytes.NewBufferString("msg: {{ env.GREETING }}\n"),
		&stdout,
		bytes.NewBuffer(nil),
		[]string{"ARGOCD_ENV_GREETING=hello"},
	)
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if got := stdout.String(); got != "msg: hello\n" {
		t.Fatalf("unexpected output: %q", got)
	}
}

func TestRunFailsOnInvalidCMPParameters(t *testing.T) {
	err := Run(
		bytes.NewBufferString("msg: {{ env.GREETING }}\n"),
		bytes.NewBuffer(nil),
		bytes.NewBuffer(nil),
		[]string{`ARGOCD_APP_PARAMETERS=not-json`},
	)
	if err == nil {
		t.Fatal("expected error")
	}
}
