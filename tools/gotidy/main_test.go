package main

import (
	"testing"

	"github.com/test-go/testify/require"
)

func TestTidy(t *testing.T) {
	project := "."

	Main(project, false, false, false)
}

func TestCompareFiles(t *testing.T) {
	before := "abc"
	after := "def"

	// should have a diff
	found := CompareFiles(before, after)
	require.Equal(t, "\x1b[31mabc\x1b[0m\x1b[32mdef\x1b[0m", found, "Did not get expected diff")

	// should not have a diff
	found = CompareFiles(after, after)
	require.Equal(t, "", found, "Found a diff when there should not have been one")
}
