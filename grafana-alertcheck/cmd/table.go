package main

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"
	"time"

	"github.com/smartcontractkit/chainlink-testing-framework/grafana-alertcheck/internal/gate"
)

// renderTable is the human table. It always writes to the writer it is given,
// which the caller (runCheck) always points at stderr — stdout is reserved for
// the machine-readable --output json.
//
// Three titled tables, in order (the name column is RULE in all of them — one
// row is one resolved alert rule, never a firing instance):
//
//  1. RESULTS, one line per rule: outcome, BadFor, pollEvery, proved-or-not
//     with the largest gap;
//  2. VIOLATIONS, one line per Violation (only when any): a rule's worst-of
//     outcome does not carry the State/Health of the instance that actually
//     caused it — Violation does — so this is also where those two columns
//     appear, sorted after the result table rather than folded into it, and it
//     is the only place an operator running WITHOUT --output json sees the
//     --allow-paused hint that Violation.Note already carries (classify.go);
//  3. THRESHOLDS, the numbers that answer "why" on exit 2: each non-skipped
//     rule's maxGap/healthGrace/evalStaleAfter, followed by the global
//     transitionGrace and drainTimeout, and the largest measured clock skew
//     alongside its own error bound (RTT/2) — SkewHardLimit is a separate,
//     fixed input threshold and is reported next to it, never as if it were
//     that bound.
func renderTable(w io.Writer, res gate.Result) error {
	alertOf := make(map[string]string, len(res.Verdicts))
	for _, v := range res.Verdicts {
		alertOf[v.RuleUID] = v.Alert
	}

	enabled := colorEnabled(w)
	// A blank line separates the result table from the notes the gate streamed
	// before it (planned run time, warning, min-observed, collecting, drain
	// wait), so the verdict reads as its own section rather than the tail of a
	// wall of progress text.
	fmt.Fprintln(w)

	fmt.Fprintln(w, "RESULTS")
	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)
	fmt.Fprintln(tw, "RULE\tOUTCOME\tBADFOR\tPOLLEVERY\tPROVED\tNOTE")
	for _, v := range sortedVerdicts(res.Verdicts) {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t%s\n",
			v.Alert, v.Outcome, v.BadFor.Round(time.Second), v.PollEvery.Round(time.Second),
			provedLabel(res.Coverage[v.RuleUID]), v.Note)
	}
	if err := tw.Flush(); err != nil {
		return fmt.Errorf("render table: %w", err)
	}

	if len(res.Violations) > 0 {
		fmt.Fprintln(w, "\nVIOLATIONS")
		vtw := tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)
		fmt.Fprintln(vtw, "RULE\tOUTCOME\tSTATE\tHEALTH\tNOTE")
		for _, v := range sortedViolations(res.Violations) {
			fmt.Fprintf(vtw, "%s\t%s\t%s\t%s\t%s\n", alertLabel(v, alertOf), v.Outcome, v.State, v.Health, v.Note)
		}
		if err := vtw.Flush(); err != nil {
			return fmt.Errorf("render table: %w", err)
		}
	}

	// The per-rule thresholds answer "why" on exit 2: a table, not the prose
	// "rule NAME: maxGap=... healthGrace=... evalStaleAfter=..." that repeated
	// the rule name a fourth time. It is separated from the result above by a
	// blank line.
	fmt.Fprintln(w)
	fmt.Fprintln(w, "THRESHOLDS")
	ttw := tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)
	fmt.Fprintln(ttw, "RULE\tMAXGAP\tHEALTHGRACE\tEVALSTALEAFTER")
	for _, uid := range sortedThresholdUIDs(res.Thresholds, alertOf) {
		t := res.Thresholds[uid]
		fmt.Fprintf(ttw, "%s\t%s\t%s\t%s\n",
			alertOr(uid, alertOf), t.MaxGap, t.HealthGrace, t.EvalStaleAfter)
	}
	if err := ttw.Flush(); err != nil {
		return fmt.Errorf("render table: %w", err)
	}

	fmt.Fprintln(w)
	fmt.Fprintf(w, "global: transitionGrace=%s (source: %s) drainTimeout=%s\n",
		res.Global.TransitionGrace, res.Global.GraceSource, res.Global.DrainTimeout)
	fmt.Fprintf(w, "largest measured clock skew: %s (bound ±%s, hard limit %s), grafana %s\n",
		res.ClockSkew.Round(time.Millisecond), res.ClockSkewBound.Round(time.Millisecond),
		gate.SkewHardLimit, res.GrafanaVersion)
	// The verdict — the single number a terminal operator reads last — sits on
	// its own line at the very bottom, separated from the diagnostics above and
	// from the shell prompt below.
	fmt.Fprintf(w, "\n%s\n\n", violationsLabel(len(res.Violations), enabled))
	return nil
}

// violationsLabel colours the "violations: N" prefix of the footer: green for a
// clean run, red otherwise. The rest of the line is written uncoloured.
func violationsLabel(n int, enabled bool) string {
	s := fmt.Sprintf("violations: %d", n)
	if !enabled {
		return s
	}
	if n == 0 {
		return ansiGreen + s + ansiReset
	}
	return ansiRed + s + ansiReset
}

// provedLabel is the table's PROVED column: "yes" for a clean coverage
// proof, "no" with the reason and largest gap for an unobservable rule, and
// "-" for a rule decide never asked proveCoverage about at all (skipped —
// paused before the window opened).
func provedLabel(cov gate.CoverageResult) string {
	if cov.Reason == "" && !cov.Unobservable && !cov.Proved {
		return "-"
	}
	if cov.Unobservable {
		if cov.LargestGap > 0 {
			return fmt.Sprintf("no (%s; largest gap %s at %s)", cov.Reason,
				cov.LargestGap.Round(time.Second), cov.LargestGapAt.Format(time.RFC3339))
		}
		return fmt.Sprintf("no (%s)", cov.Reason)
	}
	return "yes"
}

// alertLabel resolves a Violation's alert name. Most violations already
// carry it directly; the synthetic MinObserved-shortfall entry with no named
// rule (classify.go) has an empty Alert and an empty RuleUID, so alertOf
// cannot resolve it either — "-" says plainly that this row is not about a
// specific rule.
func alertLabel(v gate.Violation, alertOf map[string]string) string {
	if v.Alert != "" {
		return v.Alert
	}
	if a, ok := alertOf[v.RuleUID]; ok {
		return a
	}
	return "-"
}

func alertOr(uid string, alertOf map[string]string) string {
	if a, ok := alertOf[uid]; ok {
		return a
	}
	return uid
}

func sortedVerdicts(in []gate.RuleVerdict) []gate.RuleVerdict {
	out := append([]gate.RuleVerdict(nil), in...)
	sort.Slice(out, func(i, j int) bool { return out[i].Alert < out[j].Alert })
	return out
}

func sortedViolations(in []gate.Violation) []gate.Violation {
	out := append([]gate.Violation(nil), in...)
	sort.SliceStable(out, func(i, j int) bool { return out[i].Alert < out[j].Alert })
	return out
}

func sortedThresholdUIDs(thresholds map[string]gate.RuleThresholds, alertOf map[string]string) []string {
	uids := make([]string, 0, len(thresholds))
	for uid := range thresholds {
		uids = append(uids, uid)
	}
	sort.Slice(uids, func(i, j int) bool { return alertOr(uids[i], alertOf) < alertOr(uids[j], alertOf) })
	return uids
}
