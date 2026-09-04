package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRun_NoArgs(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run(nil, &stdout, &stderr)
	require.Equal(t, 2, code)
	require.Contains(t, stderr.String(), "usage")
}

func TestRun_Help(t *testing.T) {
	for _, flag := range []string{"-h", "-help", "--help"} {
		t.Run(flag, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := run([]string{flag}, &stdout, &stderr)
			require.Equal(t, 0, code, "requested help is not a could-not-check condition")
			require.Contains(t, stdout.String(), "usage")
			require.Empty(t, stderr.String(), "help goes to stdout")
		})
	}
}

func TestRun_UnknownSubcommand(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"bogus"}, &stdout, &stderr)
	require.Equal(t, 2, code)
	require.Contains(t, stderr.String(), `"bogus"`)
}

func TestRun_List_MissingEnv(t *testing.T) {
	t.Setenv("GRAFANA_URL", "")
	t.Setenv("GRAFANA_TOKEN", "")
	var stdout, stderr bytes.Buffer
	code := run([]string{"list"}, &stdout, &stderr)
	require.Equal(t, 2, code)
	require.Contains(t, stderr.String(), "GRAFANA_URL")
}
