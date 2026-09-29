package main

import (
	"fmt"
	"io"
	"sort"
	"strconv"
	"text/tabwriter"
	"time"

	"github.com/smartcontractkit/chainlink-testing-framework/grafana-alertcheck/internal/gate"
)

// limitsLegend explains each LIMITS USED column in one plain sentence, so the
// table needs no documentation lookup.
const limitsLegend = `  max gap without check — the longest gap between two checks we accept before we say the alert was not watched.
  query failing for — how long Grafana may keep failing to run the alert's query before we stop trusting its state.
  no evaluation for — how long Grafana may go without evaluating the alert before we stop trusting its state.`

// renderTable is the human table. It always writes to the writer it is given,
// which the caller (runCheck) always points at stderr — stdout is reserved for
// the machine-readable --output json.
//
// Three titled tables, in order (the ALERT column is one resolved alert rule,
// never a firing instance):
//
//  1. RESULTS, one line per rule: verdict, time broken, check cadence,
//     whether the window was observed, and any notes;
//  2. VIOLATIONS, one line per distinct violation, with the raw Grafana state
//     and health and the number of instances it stands for;
//  3. LIMITS USED, the coverage thresholds that answer "why" on exit 2. Each
//     column is named in plain words and explained by the legend below it, so
//     the table needs no documentation lookup.
//
// The global footer then reports the extra observation time, the drain limit
// and the largest measured clock difference, also in plain words.
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

	if te := res.TerminatedEarly; te != nil {
		detail := string(te.Reason)
		if detail == "" {
			detail = string(te.Outcome)
		}
		fmt.Fprintf(w, "EARLY EXIT: %s %q at %s (%s); the window [%s, %s] was not fully observed\n\n",
			te.Kind, te.Alert, te.At.Format(time.RFC3339), detail,
			res.From.Format(time.RFC3339), res.To.Format(time.RFC3339))
	}

	fmt.Fprintln(w, "RESULTS")
	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)
	fmt.Fprintln(tw, "ALERT\tVERDICT\tBROKEN FOR\tCHECKED EVERY\tWINDOW COVERED\tDETAILS")
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
		fmt.Fprintln(vtw, "ALERT\tVERDICT\tGRAFANA STATE\tGRAFANA HEALTH\tINSTANCES\tDETAILS")
		for _, g := range groupedViolations(res.Violations) {
			fmt.Fprintf(vtw, "%s\t%s\t%s\t%s\t%s\t%s\n", alertLabel(g.v, alertOf), g.v.Outcome, g.v.State, g.v.Health, instanceCount(g), g.v.Note)
		}
		if err := vtw.Flush(); err != nil {
			return fmt.Errorf("render table: %w", err)
		}
	}

	// The per-rule limits answer "why" on exit 2. The columns are spelled out
	// and explained by limitsLegend right below, so an operator does not have
	// to look anything up.
	fmt.Fprintln(w)
	fmt.Fprintln(w, "LIMITS USED")
	ttw := tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)
	fmt.Fprintln(ttw, "ALERT\tMAX GAP WITHOUT CHECK\tQUERY FAILING FOR\tNO EVALUATION FOR")
	for _, uid := range sortedThresholdUIDs(res.Thresholds, alertOf) {
		t := res.Thresholds[uid]
		fmt.Fprintf(ttw, "%s\t%s\t%s\t%s\n",
			alertOr(uid, alertOf), t.MaxGap, t.HealthGrace, t.EvalStaleAfter)
	}
	if err := ttw.Flush(); err != nil {
		return fmt.Errorf("render table: %w", err)
	}
	fmt.Fprintln(w, limitsLegend)

	fmt.Fprintln(w)
	if res.Global.TransitionGrace > 0 {
		fmt.Fprintf(w, "extra watching after your window: +%s — so an alert that only starts firing at the end is still caught (slowest: %s)\n",
			res.Global.TransitionGrace, res.Global.GraceSource)
	} else {
		fmt.Fprintln(w, "extra watching after your window: none")
	}
	fmt.Fprintf(w, "max wait for all alerts to finish evaluating: %s\n", res.Global.DrainTimeout)
	fmt.Fprintf(w, "clock difference from Grafana: %s, accurate to ±%s (checks fail above %s); Grafana %s\n",
		res.ClockSkew.Round(time.Millisecond), res.ClockSkewBound.Round(time.Millisecond),
		gate.SkewHardLimit, res.GrafanaVersion)
	// The verdict — the single number a terminal operator reads last — sits on
	// its own line at the very bottom, separated from the diagnostics above and
	// from the shell prompt below.
	fmt.Fprintf(w, "\n%s\n\n", violationsLabel(len(res.Violations), enabled))
	return nil
}

// violationsLabel colours the "violations: N" prefix of the footer: green when
// there are none, red otherwise. The rest of the line is written uncoloured.
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

// provedLabel is the table's WINDOW COVERED column: "yes" for a fully
// observed window, "no" with the reason and largest gap for a not-verified
// rule, and "-" for a rule decide never asked proveCoverage about at all
// (paused before the window opened).
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

type violationGroup struct {
	v gate.Violation
	n int
}

func groupedViolations(in []gate.Violation) []violationGroup {
	sorted := append([]gate.Violation(nil), in...)
	sort.SliceStable(sorted, func(i, j int) bool {
		return violationSignature(sorted[i]) < violationSignature(sorted[j])
	})
	var out []violationGroup
	for _, v := range sorted {
		if n := len(out); n > 0 && sameRendered(out[n-1].v, v) {
			out[n-1].n++
		} else {
			out = append(out, violationGroup{v: v, n: 1})
		}
	}
	return out
}

func violationSignature(v gate.Violation) string {
	return v.Alert + "\x00" + v.RuleUID + "\x00" + string(v.Outcome) + "\x00" + string(v.State) + "\x00" + v.Health + "\x00" + v.Note
}

func sameRendered(a, b gate.Violation) bool {
	return violationSignature(a) == violationSignature(b)
}

// instanceCount is the INSTANCES column: how many alert instances one grouped
// violation row stands for. Paused and not-counted rows stand for no instance
// at all, so they render "-".
func instanceCount(g violationGroup) string {
	if g.v.Outcome == gate.OutcomePaused || g.v.Outcome == gate.OutcomeNotCounted {
		return "-"
	}
	return strconv.Itoa(g.n)
}

func sortedThresholdUIDs(thresholds map[string]gate.RuleThresholds, alertOf map[string]string) []string {
	uids := make([]string, 0, len(thresholds))
	for uid := range thresholds {
		uids = append(uids, uid)
	}
	sort.Slice(uids, func(i, j int) bool { return alertOr(uids[i], alertOf) < alertOr(uids[j], alertOf) })
	return uids
}
