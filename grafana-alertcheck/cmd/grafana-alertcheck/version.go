package main

import (
	"fmt"
	"io"
)

// Build metadata. These are package-level variables so goreleaser's ldflags
// (-X) can stamp them at build time; left unstamped they fall back to the
// "dev" defaults below, which is what a plain `go build` produces.
var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
	builtBy = "unknown"
)

// runVersion prints the build metadata to stdout. Unlike list/watch/check it
// needs no Grafana connection, so it never touches the environment or the
// network; it exists purely so operators can answer "what am I running?"
// against a deployed binary.
func runVersion(args []string, stdout, stderr io.Writer) int {
	if len(args) != 0 {
		fmt.Fprintf(stderr, "version takes no arguments, got %v\n", args)
		return 2
	}
	fmt.Fprintf(stdout, "version: %s\ncommit:  %s\ndate:    %s\nbuiltBy: %s\n", version, commit, date, builtBy)
	return 0
}
