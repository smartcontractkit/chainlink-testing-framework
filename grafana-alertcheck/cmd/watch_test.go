package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// The record step's flag-validation matrix. Every case fails inside
// gate.WatchConfig.validate() or before it, so none needs a reachable Grafana.
func TestRunWatch_FlagValidation(t *testing.T) {
	tests := []struct {
		name    string
		env     bool
		args    func(t *testing.T) []string
		wantErr string
	}{
		{"missing env", false, func(t *testing.T) []string {
			return []string{"--out", t.TempDir() + "/log.jsonl", "--alerts", writeTempAlerts(t)}
		}, "GRAFANA_URL"},
		{"missing out", true, func(t *testing.T) []string {
			return []string{"--alerts", writeTempAlerts(t)}
		}, "no log path"},
		{"missing alerts", true, func(t *testing.T) []string {
			return []string{"--out", t.TempDir() + "/log.jsonl"}
		}, "no alert names"},
		{"bad until format", true, func(t *testing.T) []string {
			return []string{"--out", t.TempDir() + "/log.jsonl", "--alerts", writeTempAlerts(t), "--until", "not-a-time"}
		}, "--until"},
		{"until in the past", true, func(t *testing.T) []string {
			return []string{"--out", t.TempDir() + "/log.jsonl", "--alerts", writeTempAlerts(t), "--until", "2000-01-01T00:00:00Z"}
		}, "not in the future"},
		{"bad poll-interval", true, func(t *testing.T) []string {
			return []string{"--out", t.TempDir() + "/log.jsonl", "--alerts", writeTempAlerts(t), "--poll-interval", "not-a-duration"}
		}, "--poll-interval"},
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
			args := append([]string{"watch"}, tt.args(t)...)
			code := run(args, &stdout, &stderr)
			require.Equal(t, 2, code)
			require.Contains(t, stderr.String(), tt.wantErr)
		})
	}
}

// Seeing gate.DaemonChildFlag must dispatch to gate.RunDaemonChild, and the
// flag must never appear in watchUsage (an operator never types it).
func TestRunWatch_DaemonChildDispatch(t *testing.T) {
	t.Setenv("GRAFANA_URL", "http://example.invalid")
	t.Setenv("GRAFANA_TOKEN", "test-token")

	var stdout, stderr bytes.Buffer
	// No log at this path: RunDaemonChild fails trying to read it, which is
	// enough to prove dispatch happened without needing a real recording.
	missing := os.DevNull + ".missing"
	code := run([]string{"watch", "--daemon-child", "--out", missing, "--ready-fd", "0"}, &stdout, &stderr)
	require.Equal(t, 2, code)
	require.Contains(t, stderr.String(), missing)
	require.NotContains(t, watchUsage, "daemon-child")
	require.NotContains(t, watchUsage, "ready-fd")
}

func TestRunWatch_DaemonChild_MissingEnv(t *testing.T) {
	t.Setenv("GRAFANA_URL", "")
	t.Setenv("GRAFANA_TOKEN", "")

	var stdout, stderr bytes.Buffer
	code := run([]string{"watch", "--daemon-child", "--out", "log.jsonl"}, &stdout, &stderr)
	require.Equal(t, 2, code)
	require.Contains(t, stderr.String(), "GRAFANA_URL")
}
