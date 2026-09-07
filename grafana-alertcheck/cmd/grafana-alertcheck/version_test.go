package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunVersion(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"version"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("code = %d, want 0", code)
	}
	out := stdout.String()
	for _, want := range []string{"version:", "commit:", "date:", "builtBy:"} {
		if !strings.Contains(out, want) {
			t.Errorf("stdout = %q, want it to contain %q", out, want)
		}
	}
	if stderr.String() != "" {
		t.Errorf("stderr = %q, want empty", stderr.String())
	}
}

func TestRunVersion_RejectsArgs(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"version", "extra"}, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("code = %d, want 2", code)
	}
}
