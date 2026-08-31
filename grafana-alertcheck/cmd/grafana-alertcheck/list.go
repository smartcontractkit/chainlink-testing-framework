package main

import (
	"context"
	"fmt"
	"io"
	"sort"
	"text/tabwriter"

	"github.com/smartcontractkit/chainlink-testing-framework/grafana-alertcheck/internal/gate"
)

// runList reads every rule definition from the ruler endpoint and prints one
// line per rule: its kind, its Folder/Group/Title, and its uid. This is what
// makes the gate runnable end to end before any coverage logic exists (§9
// rule 4) — it validates auth, the ruler parse, and the shapes Resolve
// matches against, all against a real Grafana. It is also the "did you mean"
// surface §17.2's no-match error points operators at.
func runList(args []string, stdout, stderr io.Writer) int {
	if len(args) != 0 {
		fmt.Fprintf(stderr, "list takes no arguments, got %v\n", args)
		return 2
	}

	url, token, err := grafanaEnv()
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}

	// context.Background(), no outer deadline: httpSource bounds every single
	// attempt with its http.Client's 30s Timeout (source.go) and gives up
	// after maxSequentialFailures consecutive transport errors, so this call
	// always terminates. It can still take minutes end-to-end under repeated
	// transient failures (5 retries * up to 30s backoff each, per call) — an
	// acceptable wait for an interactive `list`, not for `watch`/`check`,
	// which get their own deadlines from `--until`/`to` in P10.
	src := gate.NewHTTPSource(url, token, gate.SystemClock{})
	version, err := src.Version(context.Background())
	if err != nil {
		fmt.Fprintf(stderr, "checking grafana version: %v\n", err)
		return 2
	}
	if err := gate.CheckGrafanaVersion(version); err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}

	defs, err := src.Definitions(context.Background())
	if err != nil {
		fmt.Fprintf(stderr, "reading rule definitions: %v\n", err)
		return 2
	}

	sort.Slice(defs, func(i, j int) bool {
		if defs[i].Folder != defs[j].Folder {
			return defs[i].Folder < defs[j].Folder
		}
		if defs[i].Group != defs[j].Group {
			return defs[i].Group < defs[j].Group
		}
		return defs[i].Title < defs[j].Title
	})

	tw := tabwriter.NewWriter(stdout, 0, 4, 2, ' ', 0)
	fmt.Fprintln(tw, "KIND\tFOLDER\tGROUP\tTITLE\tUID")
	for _, d := range defs {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\n", kindLabel(d.Kind), d.Folder, d.Group, d.Title, uidOrDash(d.UID))
	}
	if err := tw.Flush(); err != nil {
		fmt.Fprintf(stderr, "writing output: %v\n", err)
		return 2
	}
	return 0
}

func kindLabel(k gate.RuleKind) string {
	switch k {
	case gate.KindDatasourceManaged:
		return "datasource-managed"
	case gate.KindRecording:
		return "recording"
	default:
		return "grafana-managed"
	}
}

func uidOrDash(uid string) string {
	if uid == "" {
		return "-"
	}
	return uid
}
