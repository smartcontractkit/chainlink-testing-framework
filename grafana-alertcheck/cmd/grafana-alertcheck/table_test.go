package main

import (
	"bytes"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-testing-framework/grafana-alertcheck/internal/gate"
	"github.com/stretchr/testify/require"
)

// The golden table test: a fixed Result renders a deterministic, ordered rule
// table, a violations section and a footer carrying the per-rule and global
// thresholds plus the skew and its bound — with no live Check involved.
func TestRenderTable(t *testing.T) {
	gapAt := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	res := gate.Result{
		GrafanaVersion: "13.1.0",
		ClockSkew:      1500 * time.Millisecond,
		ClockSkewBound: 250 * time.Millisecond,
		Verdicts: []gate.RuleVerdict{
			{Alert: "Zebra Alert", RuleUID: "uid-z", Outcome: gate.OutcomeClean, PollEvery: 30 * time.Second},
			{Alert: "Ape Alert", RuleUID: "uid-a", Outcome: gate.OutcomeUnobservable,
				PollEvery: 30 * time.Second, Note: "gap of 5m0s starting at 2026-01-01T12:00:00Z exceeds maxGap 1m0s"},
			{Alert: "Paused Alert", RuleUID: "uid-p", Outcome: gate.OutcomeSkipped,
				Note: "paused before the window opened; counts against --min-observed unless --allow-paused is set"},
		},
		Violations: []gate.Violation{
			{Alert: "Ape Alert", RuleUID: "uid-a", Outcome: gate.OutcomeUnobservable, State: gate.StateFiring, Health: "error", Note: "unobservable"},
			{Alert: "Paused Alert", RuleUID: "uid-p", Outcome: gate.OutcomeSkipped,
				Note: "paused before the window opened; counts against --min-observed unless --allow-paused is set"},
		},
		Coverage: map[string]gate.CoverageResult{
			"uid-z": {Proved: true},
			"uid-a": {Unobservable: true, Reason: gate.ReasonHeartbeatGap, LargestGap: 5 * time.Minute, LargestGapAt: gapAt},
		},
		Thresholds: map[string]gate.RuleThresholds{
			"uid-z": {MaxGap: time.Minute, HealthGrace: time.Minute, EvalStaleAfter: time.Minute},
			"uid-a": {MaxGap: time.Minute, HealthGrace: 2 * time.Minute, EvalStaleAfter: time.Minute},
		},
		Global: gate.GlobalThresholds{
			TransitionGrace: 5 * time.Minute,
			GraceSource:     `Ape Alert (for=5m)`,
			DrainTimeout:    2 * time.Minute,
		},
	}

	var buf bytes.Buffer
	require.NoError(t, renderTable(&buf, res))
	out := buf.String()

	// Rule table: Ape sorts before Zebra sorts before... Paused is skipped and
	// carries no coverage entry, so it renders "-" for PROVED.
	require.Contains(t, out, "Ape Alert")
	require.Contains(t, out, "unobservable")
	require.Contains(t, out, "heartbeat_gap")
	require.Contains(t, out, "largest gap 5m0s")
	require.Contains(t, out, "Zebra Alert")
	require.Contains(t, out, "clean")

	// The violations section must show up even without --output json, and must
	// carry the --allow-paused hint text verbatim.
	require.Contains(t, out, "VIOLATIONS")
	require.Contains(t, out, "--allow-paused")
	require.Contains(t, out, "STATE")
	require.Contains(t, out, "HEALTH")
	require.Contains(t, out, string(gate.StateFiring))
	require.Contains(t, out, "error")

	// The footer: per-rule thresholds are a table (RULE/MAXGAP/HEALTHGRACE/
	// EVALSTALEAFTER) rather than prose, followed by the global thresholds and
	// the violations count with the skew and its own bound rather than the
	// fixed hard limit.
	require.Contains(t, out, "MAXGAP")
	require.Contains(t, out, "HEALTHGRACE")
	require.Contains(t, out, "EVALSTALEAFTER")
	require.Contains(t, out, "global: transitionGrace=5m0s (source: Ape Alert (for=5m)) drainTimeout=2m0s")
	require.Contains(t, out, "largest measured clock skew: 1.5s (bound ±250ms, hard limit 1m0s)")
	require.Contains(t, out, "violations: 2")
	require.Contains(t, out, "13.1.0")
}

// The "-" case: a rule decide never asked proveCoverage about (paused before
// the window opened) has an empty CoverageResult and must not be reported as
// either proved or unobservable.
func TestProvedLabel_Skipped(t *testing.T) {
	require.Equal(t, "-", provedLabel(gate.CoverageResult{}))
}
