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
