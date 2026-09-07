package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunVersion(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"version"}, &stdout, &stderr)
	require.Equal(t, 0, code)
	out := stdout.String()
	for _, want := range []string{"version:", "commit:", "date:", "builtBy:"} {
		require.Contains(t, out, want)
	}
	require.Empty(t, stderr.String())
}

func TestRunVersion_RejectsArgs(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"version", "extra"}, &stdout, &stderr)
	require.Equal(t, 2, code)
}
