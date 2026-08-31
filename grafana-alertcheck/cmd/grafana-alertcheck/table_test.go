package main

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-testing-framework/grafana-alertcheck/internal/gate"
)

// TestRenderTable is the golden table test: a fixed Result renders a
// deterministic, ordered rule table, a violations section (R2) and a footer
// carrying the per-rule and global thresholds plus the skew and its bound
// (R3) — with no live Check involved.
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
	if err := renderTable(&buf, res); err != nil {
		t.Fatalf("renderTable: %v", err)
	}
	out := buf.String()

	// Rule table: Ape sorts before Zebra sorts before... Paused is skipped and
	// carries no coverage entry, so it renders "-" for PROVED.
	if !strings.Contains(out, "Ape Alert") || !strings.Contains(out, "unobservable") {
		t.Fatalf("out = %q, want Ape's unobservable row", out)
	}
	if !strings.Contains(out, "heartbeat_gap") || !strings.Contains(out, "largest gap 5m0s") {
		t.Fatalf("out = %q, want the coverage reason and largest gap", out)
	}
	if !strings.Contains(out, "Zebra Alert") || !strings.Contains(out, "clean") {
		t.Fatalf("out = %q, want Zebra's clean row", out)
	}

	// Violations section (R2): must show up even without --output json, and
	// must carry the §12.1 --allow-paused hint text verbatim.
	if !strings.Contains(out, "VIOLATIONS") {
		t.Fatalf("out = %q, want a VIOLATIONS section", out)
	}
	if !strings.Contains(out, "--allow-paused") {
		t.Fatalf("out = %q, want the §12.1 --allow-paused hint in the human table", out)
	}
	if !strings.Contains(out, "STATE") || !strings.Contains(out, "HEALTH") {
		t.Fatalf("out = %q, want the violations table to have STATE and HEALTH columns", out)
	}
	if !strings.Contains(out, string(gate.StateFiring)) || !strings.Contains(out, "error") {
		t.Fatalf("out = %q, want Ape's violation State/Health", out)
	}

	// Footer (R3): per-rule thresholds, global thresholds, and skew with its
	// own bound rather than the fixed hard limit.
	if !strings.Contains(out, "Ape Alert: maxGap=1m0s healthGrace=2m0s evalStaleAfter=1m0s") {
		t.Fatalf("out = %q, want Ape's per-rule thresholds", out)
	}
	if !strings.Contains(out, "Zebra Alert: maxGap=1m0s healthGrace=1m0s evalStaleAfter=1m0s") {
		t.Fatalf("out = %q, want Zebra's per-rule thresholds", out)
	}
	if strings.Contains(out, "Paused Alert: maxGap") {
		t.Fatalf("out = %q, a skipped rule must not report thresholds it never had (§12)", out)
	}
	if !strings.Contains(out, "global: transitionGrace=5m0s (source: Ape Alert (for=5m)) drainTimeout=2m0s") {
		t.Fatalf("out = %q, want the global thresholds line", out)
	}
	if !strings.Contains(out, "largest measured clock skew: 1.5s (bound ±250ms, hard limit 1m0s)") {
		t.Fatalf("out = %q, want the skew and its own bound, not the hard limit misused as one", out)
	}
	if !strings.Contains(out, "violations: 2") {
		t.Fatalf("out = %q, want the violation count", out)
	}
	if !strings.Contains(out, "13.1.0") {
		t.Fatalf("out = %q, want the grafana version", out)
	}
}

// TestProvedLabel_Skipped pins the "-" case: a rule decide never asked
// proveCoverage about (paused before the window opened, §12) has an empty
// CoverageResult and must not be reported as either proved or unobservable.
func TestProvedLabel_Skipped(t *testing.T) {
	if got := provedLabel(gate.CoverageResult{}); got != "-" {
		t.Fatalf("provedLabel(zero value) = %q, want \"-\"", got)
	}
}
