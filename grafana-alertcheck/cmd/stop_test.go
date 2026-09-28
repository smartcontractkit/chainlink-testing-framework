package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunStop_RequiresOut(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"stop"}, &stdout, &stderr)
	require.Equal(t, 2, code)
	require.Contains(t, stderr.String(), "--out is required")
}

// Nothing recorded: an `if: always()` stop still exits 0.
func TestRunStop_NoPidfileIsANoOp(t *testing.T) {
	out := filepath.Join(t.TempDir(), "log.jsonl")
	require.NoError(t, os.WriteFile(out, nil, 0o644)) // nolint:gosec // test-only temp file

	var stdout, stderr bytes.Buffer
	code := run([]string{"stop", "--out", out}, &stdout, &stderr)
	require.Equal(t, 0, code)
	require.Contains(t, stderr.String(), "nothing to stop")
}
