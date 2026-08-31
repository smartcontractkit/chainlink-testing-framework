package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRun_NoArgs(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run(nil, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("code = %d, want 2", code)
	}
	if !strings.Contains(stderr.String(), "usage") {
		t.Fatalf("stderr = %q, want a usage message", stderr.String())
	}
}

func TestRun_Help(t *testing.T) {
	for _, flag := range []string{"-h", "-help", "--help"} {
		t.Run(flag, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := run([]string{flag}, &stdout, &stderr)
			if code != 0 {
				t.Fatalf("code = %d, want 0 (requested help is not a could-not-check condition)", code)
			}
			if !strings.Contains(stdout.String(), "usage") {
				t.Fatalf("stdout = %q, want a usage message", stdout.String())
			}
			if stderr.String() != "" {
				t.Fatalf("stderr = %q, want empty — help goes to stdout", stderr.String())
			}
		})
	}
}

func TestRun_UnknownSubcommand(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"bogus"}, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("code = %d, want 2", code)
	}
	if !strings.Contains(stderr.String(), `"bogus"`) {
		t.Fatalf("stderr = %q, want it to name the unknown subcommand", stderr.String())
	}
}

func TestRun_List_MissingEnv(t *testing.T) {
	t.Setenv("GRAFANA_URL", "")
	t.Setenv("GRAFANA_TOKEN", "")
	var stdout, stderr bytes.Buffer
	code := run([]string{"list"}, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("code = %d, want 2", code)
	}
	if !strings.Contains(stderr.String(), "GRAFANA_URL") {
		t.Fatalf("stderr = %q, want it to name the missing env var", stderr.String())
	}
}
