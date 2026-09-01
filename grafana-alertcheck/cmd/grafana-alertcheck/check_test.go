package main

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/smartcontractkit/chainlink-testing-framework/grafana-alertcheck/internal/gate"
)

// The exit-code mapping, pinned directly against exitCode with no network
// involved: err != nil is exit 2 even alongside violations (an inability to
// check beats a violation), violations alone are exit 1, and neither is 0.
func TestExitCode(t *testing.T) {
	tests := []struct {
		name string
		res  gate.Result
		err  error
		want int
	}{
		{"pass", gate.Result{}, nil, 0},
		{"violation", gate.Result{Violations: []gate.Violation{{}}}, nil, 1},
		{"error alone", gate.Result{}, errors.New("boom"), 2},
		{"error beats violation", gate.Result{Violations: []gate.Violation{{}}}, errors.New("boom"), 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := exitCode(tt.res, tt.err); got != tt.want {
				t.Fatalf("exitCode(...) = %d, want %d", got, tt.want)
			}
		})
	}
}

func writeTempAlerts(t *testing.T) string {
	t.Helper()
	path := t.TempDir() + "/alerts.txt"
	if err := os.WriteFile(path, []byte("Some Alert\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

// The flag-validation matrix: every one of these must fail before any network
// call, because gate.Config.validate() runs first — an unreachable GRAFANA_URL
// succeeding or timing out is a different test than these, which check pure
// input validation.
func TestRunCheck_FlagValidation(t *testing.T) {
	tests := []struct {
		name    string
		env     bool
		args    func(t *testing.T) []string
		wantErr string
	}{
		{"missing env", false, func(t *testing.T) []string {
			return []string{"--to", "2026-01-01T00:00:00Z", "--alerts", writeTempAlerts(t)}
		}, "GRAFANA_URL"},
		{"missing to", true, func(t *testing.T) []string {
			return []string{"--alerts", writeTempAlerts(t)}
		}, "--to"},
		{"bad to", true, func(t *testing.T) []string {
			return []string{"--to", "not-a-time", "--alerts", writeTempAlerts(t)}
		}, "--to"},
		{"bad from", true, func(t *testing.T) []string {
			return []string{"--to", "2026-01-01T00:00:00Z", "--from", "not-a-time", "--alerts", writeTempAlerts(t)}
		}, "--from"},
		{"bad output", true, func(t *testing.T) []string {
			return []string{"--to", "2026-01-01T00:00:00Z", "--output", "xml", "--alerts", writeTempAlerts(t)}
		}, "--output"},
		{"bad states", true, func(t *testing.T) []string {
			return []string{"--to", "2026-01-01T00:00:00Z", "--states", "bogus", "--alerts", writeTempAlerts(t)}
		}, "--states"},
		{"states normal is rejected", true, func(t *testing.T) []string {
			// normal is the good state, never a state to classify AS bad:
			// accepting it would make --states normal fail every healthy
			// instance.
			return []string{"--to", "2026-01-01T00:00:00Z", "--states", "normal", "--alerts", writeTempAlerts(t)}
		}, "--states"},
		{"bad preexisting", true, func(t *testing.T) []string {
			return []string{"--to", "2026-01-01T00:00:00Z", "--preexisting", "bogus", "--alerts", writeTempAlerts(t)}
		}, "--preexisting"},
		{"alerts with in", true, func(t *testing.T) []string {
			return []string{"--to", "2026-01-01T00:00:00Z", "--in", "some.jsonl", "--alerts", writeTempAlerts(t)}
		}, "refused"},
		{"no alerts no in", true, func(t *testing.T) []string {
			return []string{"--to", "2026-01-01T00:00:00Z"}
		}, "no alert names"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.env {
				t.Setenv("GRAFANA_URL", "http://example.invalid")
				t.Setenv("GRAFANA_TOKEN", "test-token")
			} else {
				t.Setenv("GRAFANA_URL", "")
				t.Setenv("GRAFANA_TOKEN", "")
			}
			var stdout, stderr bytes.Buffer
			args := append([]string{"check"}, tt.args(t)...)
			code := run(args, &stdout, &stderr)
			if code != 2 {
				t.Fatalf("code = %d, want 2; stderr = %q", code, stderr.String())
			}
			if !strings.Contains(stderr.String(), tt.wantErr) {
				t.Fatalf("stderr = %q, want it to contain %q", stderr.String(), tt.wantErr)
			}
		})
	}
}

// A `to` already in the past with no recorded log cannot be classified from
// anything, because nothing ever observed the window.
func TestRunCheck_ToInPastNoLog(t *testing.T) {
	t.Setenv("GRAFANA_URL", "http://example.invalid")
	t.Setenv("GRAFANA_TOKEN", "test-token")

	var stdout, stderr bytes.Buffer
	code := run([]string{"check",
		"--from", "1999-01-01T00:00:00Z", "--to", "2000-01-01T00:00:00Z",
		"--alerts", writeTempAlerts(t),
	}, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("code = %d, want 2; stderr = %q", code, stderr.String())
	}
	if !strings.Contains(stderr.String(), "already passed") {
		t.Fatalf("stderr = %q, want the past-`to` refusal", stderr.String())
	}
}

// --output json never writes to stdout when Check was never reached, because
// there is no Result to encode — only the table (on stderr) can report a
// configuration failure.
func TestRunCheck_NoResultOnConfigError(t *testing.T) {
	t.Setenv("GRAFANA_URL", "")
	t.Setenv("GRAFANA_TOKEN", "")
	var stdout, stderr bytes.Buffer
	code := run([]string{"check", "--to", "2026-01-01T00:00:00Z", "--output", "json"}, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("code = %d, want 2", code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout = %q, want empty", stdout.String())
	}
}
