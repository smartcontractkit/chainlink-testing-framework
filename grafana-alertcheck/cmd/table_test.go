package main

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/grafana-alertcheck/internal/gate"
)

// The golden table test: a fixed Result renders a deterministic, ordered rule
// table, a violations section and a footer carrying the per-rule limits and
// globals in plain words — with no live Check involved.
func TestRenderTable(t *testing.T) {
	gapAt := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	res := gate.Result{
		GrafanaVersion: "13.1.0",
		ClockSkew:      1500 * time.Millisecond,
		ClockSkewBound: 250 * time.Millisecond,
		Verdicts: []gate.RuleVerdict{
			{Alert: "Zebra Alert", RuleUID: "uid-z", Outcome: gate.OutcomeHealthy, PollEvery: 30 * time.Second},
			{Alert: "Ape Alert", RuleUID: "uid-a", Outcome: gate.OutcomeNotVerified,
				PollEvery: 30 * time.Second, Note: "gap of 5m0s starting at 2026-01-01T12:00:00Z exceeds maxGap 1m0s"},
			{Alert: "Paused Alert", RuleUID: "uid-p", Outcome: gate.OutcomePaused,
				Note: "paused before the window opened; counts against --min-observed unless --allow-paused is set"},
		},
		Violations: []gate.Violation{
			{Alert: "Ape Alert", RuleUID: "uid-a", Outcome: gate.OutcomeNotVerified, State: gate.StateFiring, Health: "error", Note: "not verified"},
			{Alert: "Paused Alert", RuleUID: "uid-p", Outcome: gate.OutcomePaused,
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

	// The rule table: Ape sorts before Zebra sorts before... Paused carries no
	// coverage entry, so it renders "-" for WINDOW COVERED.
	require.Contains(t, out, "ALERT")
	require.Contains(t, out, "Ape Alert")
	require.Contains(t, out, "not_verified")
	require.Contains(t, out, "heartbeat_gap")
	require.Contains(t, out, "largest gap 5m0s")
	require.Contains(t, out, "Zebra Alert")
	require.Contains(t, out, "healthy")
	require.Contains(t, out, "WINDOW COVERED")

	// The violations section must show up even without --output json, and must
	// carry the --allow-paused hint text verbatim. INSTANCES is a single word
	// so the count under it cannot read as a second, empty column.
	require.Contains(t, out, "VIOLATIONS")
	require.Contains(t, out, "--allow-paused")
	require.Contains(t, out, "GRAFANA STATE")
	require.Contains(t, out, "GRAFANA HEALTH")
	require.Contains(t, out, "INSTANCES")
	require.Contains(t, out, string(gate.StateFiring))
	require.Contains(t, out, "error")

	// The limits table names each threshold in plain words and explains it
	// right below, so an operator does not have to consult the docs.
	require.Contains(t, out, "LIMITS USED")
	require.Contains(t, out, "MAX GAP WITHOUT CHECK")
	require.Contains(t, out, "QUERY FAILING FOR")
	require.Contains(t, out, "NO EVALUATION FOR")
	require.Contains(t, out, "the longest gap between two checks")
	require.Contains(t, out, "without evaluating the alert")

	// The global footer in plain words.
	require.Contains(t, out, "extra watching after your window: +5m0s")
	require.Contains(t, out, "slowest: Ape Alert (for=5m)")
	require.Contains(t, out, "max wait for all alerts to finish evaluating: 2m0s")
	require.Contains(t, out, "clock difference from Grafana: 1.5s, accurate to ±250ms (checks fail above 1m0s); Grafana 13.1.0")
	require.Contains(t, out, "violations: 2")
}

// The "-" case: a rule decide never asked proveCoverage about (paused before
// the window opened) has an empty CoverageResult and must not be reported as
// either covered or not verified.
func TestProvedLabel_Paused(t *testing.T) {
	require.Equal(t, "-", provedLabel(gate.CoverageResult{}))
}

// groupedViolations collapses a rule's many firing instances into one row per
// rendered signature, each with a count.
func TestGroupedViolations(t *testing.T) {
	in := []gate.Violation{
		{Alert: "OCR2 Consensus failure", RuleUID: "uid-o", Outcome: gate.OutcomeStillFailing, State: gate.StateFiring, Health: "ok"},
		{Alert: "OCR2 Consensus failure", RuleUID: "uid-o", Outcome: gate.OutcomeStillFailing, State: gate.StateFiring, Health: "ok"},
		{Alert: "OCR2 Consensus failure", RuleUID: "uid-o", Outcome: gate.OutcomeStillFailing, State: gate.StateFiring, Health: "ok"},
		{Alert: "OCR2 Consensus failure", RuleUID: "uid-o", Outcome: gate.OutcomeStillFailing, State: gate.StateFiring, Health: "error"},
		{Alert: "Other Alert", RuleUID: "uid-p", Outcome: gate.OutcomeNewFailure, State: gate.StateFiring, Health: "ok", Note: "x"},
		{Alert: "Other Alert", RuleUID: "uid-p", Outcome: gate.OutcomeNewFailure, State: gate.StateFiring, Health: "ok", Note: "x"},
	}

	got := groupedViolations(in)

	require.Len(t, got, 3)
	counts := map[string]int{}
	for _, g := range got {
		counts[g.v.Health+"|"+string(g.v.Outcome)] = g.n
	}
	require.Equal(t, 3, counts["ok|"+string(gate.OutcomeStillFailing)])
	require.Equal(t, 1, counts["error|"+string(gate.OutcomeStillFailing)])
	require.Equal(t, 2, counts["ok|"+string(gate.OutcomeNewFailure)])
}
