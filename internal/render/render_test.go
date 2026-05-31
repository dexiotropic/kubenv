package render

import (
	"bytes"
	"testing"
)

func TestStrictReplacesVariables(t *testing.T) {
	input := []byte("message: {{ env.GREETING }}\n")
	output, err := Strict(input, map[string]string{"GREETING": "hello"})
	if err != nil {
		t.Fatalf("Strict returned error: %v", err)
	}

	if !bytes.Equal(output, []byte("message: hello\n")) {
		t.Fatalf("unexpected output: %q", output)
	}
}

func TestStrictFailsOnMissingVariables(t *testing.T) {
	_, err := Strict([]byte("message: {{ env.GREETING }}\n"), map[string]string{})
	if err == nil {
		t.Fatal("expected missing variable error")
	}
}

func TestStrictWithShellStyleReplacesVariables(t *testing.T) {
	input := []byte("message: $GREETING ${TARGET}\n")
	output, err := StrictWithStyle(input, map[string]string{
		"GREETING": "hello",
		"TARGET":   "world",
	}, StyleShell)
	if err != nil {
		t.Fatalf("StrictWithStyle returned error: %v", err)
	}

	if !bytes.Equal(output, []byte("message: hello world\n")) {
		t.Fatalf("unexpected output: %q", output)
	}
}

func TestStrictWithShellStyleFailsOnMissingVariables(t *testing.T) {
	_, err := StrictWithStyle([]byte("message: $GREETING\n"), map[string]string{}, StyleShell)
	if err == nil {
		t.Fatal("expected missing variable error")
	}
}

func TestFromEnvironParsesEntriesWithEquals(t *testing.T) {
	vars := FromEnviron([]string{
		"GREETING=hello",
		"TARGET=world=wide",
		"INVALID",
	})

	if got := vars["GREETING"]; got != "hello" {
		t.Fatalf("unexpected GREETING: %q", got)
	}
	if got := vars["TARGET"]; got != "world=wide" {
		t.Fatalf("unexpected TARGET: %q", got)
	}
	if _, ok := vars["INVALID"]; ok {
		t.Fatal("expected invalid entry to be ignored")
	}
}

func TestStrictWithUnknownStyleFails(t *testing.T) {
	_, err := StrictWithStyle([]byte("message: {{ env.GREETING }}\n"), map[string]string{}, Style("unknown"))
	if err == nil {
		t.Fatal("expected unknown style error")
	}
}
