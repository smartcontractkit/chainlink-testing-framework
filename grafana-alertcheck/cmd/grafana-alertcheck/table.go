package main

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"
	"time"

	"github.com/smartcontractkit/chainlink-testing-framework/grafana-alertcheck/internal/gate"
)

// renderTable is §20.2's required human table. It always writes to the
// writer it is given, which the caller (runCheck) always points at
// stderr — the human table is not the machine output §20.2 reserves stdout
// for.
//
// Three sections, in order:
//
//  1. one line per rule: outcome, BadFor, pollEvery, proved-or-not with the
//     largest gap;
//  2. one line per Violation (R2): a rule's worst-of outcome does not carry
//     the State/Health of the instance that actually caused it — Violation
//     does — so this is also where those two columns appear, sorted after
//     the rule table rather than folded into it, and it is the only place an
//     operator running WITHOUT --output json sees the §12.1 --allow-paused
//     hint that Violation.Note already carries (classify.go);
//  3. a footer with the per-rule thresholds and the run-wide numbers §20.2
//     says are the answer to "why" on exit 2: each non-skipped rule's
//     maxGap/healthGrace/evalStaleAfter, the global transitionGrace and
//     drainTimeout, and the largest measured clock skew alongside its own
//     error bound (RTT/2) — SkewHardLimit is a separate, fixed input
//     threshold and is reported next to it, never as if it were that bound.
func renderTable(w io.Writer, res gate.Result) error {
	alertOf := make(map[string]string, len(res.Verdicts))
	for _, v := range res.Verdicts {
		alertOf[v.RuleUID] = v.Alert
	}

	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)
	fmt.Fprintln(tw, "ALERT\tOUTCOME\tBADFOR\tPOLLEVERY\tPROVED\tNOTE")
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
		fmt.Fprintln(vtw, "ALERT\tOUTCOME\tSTATE\tHEALTH\tNOTE")
		for _, v := range sortedViolations(res.Violations) {
			fmt.Fprintf(vtw, "%s\t%s\t%s\t%s\t%s\n", alertLabel(v, alertOf), v.Outcome, v.State, v.Health, v.Note)
		}
		if err := vtw.Flush(); err != nil {
			return fmt.Errorf("render table: %w", err)
		}
	}

	fmt.Fprintln(w)
	for _, uid := range sortedThresholdUIDs(res.Thresholds, alertOf) {
		t := res.Thresholds[uid]
		fmt.Fprintf(w, "rule %s: maxGap=%s healthGrace=%s evalStaleAfter=%s\n",
			alertOr(uid, alertOf), t.MaxGap, t.HealthGrace, t.EvalStaleAfter)
	}
	fmt.Fprintf(w, "global: transitionGrace=%s (source: %s) drainTimeout=%s\n",
		res.Global.TransitionGrace, res.Global.GraceSource, res.Global.DrainTimeout)
	fmt.Fprintf(w, "violations: %d, largest measured clock skew: %s (bound ±%s, hard limit %s), grafana %s\n",
		len(res.Violations), res.ClockSkew.Round(time.Millisecond), res.ClockSkewBound.Round(time.Millisecond),
		gate.SkewHardLimit, res.GrafanaVersion)
	return nil
}

// provedLabel is the table's PROVED column: "yes" for a clean coverage
// proof, "no" with the reason and largest gap for an unobservable rule, and
// "-" for a rule decide never asked proveCoverage about at all (skipped —
// paused before the window opened, §12).
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
