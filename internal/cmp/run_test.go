package cmp

import (
	"bytes"
	"testing"
)

func TestRunUsesCMPStringParameters(t *testing.T) {
	var stdout bytes.Buffer
	err := Run(
		nil,
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
		nil,
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
		nil,
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
		nil,
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
		nil,
		bytes.NewBufferString("msg: {{ env.GREETING }}\n"),
		bytes.NewBuffer(nil),
		bytes.NewBuffer(nil),
		[]string{`ARGOCD_APP_PARAMETERS=not-json`},
	)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRunShowsHelp(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := Run([]string{"--help"}, bytes.NewBuffer(nil), &stdout, &stderr, nil)
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if got := stdout.String(); !bytes.Contains([]byte(got), []byte("ARGOCD_APP_PARAMETERS")) {
		t.Fatalf("unexpected help output: %q", got)
	}
	if got := stderr.String(); got != "" {
		t.Fatalf("unexpected stderr: %q", got)
	}
}
