package gate

import (
	"fmt"
	"slices"
	"time"
)

// keepLastReason is the instance Reason that check 9 watches for (§10.2).
const keepLastReason = "KeepLast"

// Obligations this phase leaves for later ones — carried forward the same
// way P6's own deviations list did, so a later review has something concrete
// to check against:
//
//   - fromFutureTolerance (§5: 60s) has no constant and no hard-error check
//     anywhere yet. Check 2 below implements only "from < StartedAt"; the
//     second clause — from more than fromFutureTolerance ahead is a hard
//     error — is once-per-run input validation, not a per-rule coverage
//     check, and belongs to Check's construction in a later phase (P9).
//   - decide (P8) must read a rule's skipped status from the definitions
//     (LoggedRule.IsPaused / Definition.IsPaused), never from the polls, and
//     must do so BEFORE calling proveCoverage for that rule: a rule paused
//     before the window opened is never scheduled or polled (§4.3), so it
//     reaches this function with zero polls and today reads as one large
//     heartbeat_gap, not skipped (pinned by
//     TestProveCoverage_SkippedRuleWithZeroPollsPinnedAsHeartbeatGap).

// UnobservableReason names why proveCoverage could not prove a rule's window.
// It is machine-readable — this reaches the action's JSON outputs, so it is a
// published vocabulary like Outcome (§19.0); prose belongs in Notes.
type UnobservableReason string

const (
	ReasonNoSentinel       UnobservableReason = "no_sentinel"
	ReasonSentinelEarly    UnobservableReason = "sentinel_early"
	ReasonFromBeforeRecord UnobservableReason = "from_before_record"
	ReasonHeartbeatGap     UnobservableReason = "heartbeat_gap"
	ReasonHealthError      UnobservableReason = "health_error"
	ReasonStaleEvaluation  UnobservableReason = "stale_evaluation"
	ReasonPausedInWindow   UnobservableReason = "paused_in_window"
	ReasonRuleAbsent       UnobservableReason = "rule_absent"
	// ReasonDrainTimeout is set by check.go's drain wait (a later phase),
	// never by proveCoverage: the wait is I/O and must not be added to this
	// pure function — that would put HTTP inside the pure layer and destroy
	// the seam §2's architecture depends on.
	ReasonDrainTimeout UnobservableReason = "drain_timeout"
)

// CoverageResult is proveCoverage's whole answer for one rule. No interval
// list: proved-or-not plus the largest gap and where is everything a human
// reads on exit 2, and everything §20.2's table needs.
type CoverageResult struct {
	Proved       bool
	LargestGap   time.Duration
	LargestGapAt time.Time
	Unobservable bool
	Reason       UnobservableReason
	Notes        []string
	// BlindFor is the worst staleness (GrafanaNow - LastEvaluation) that
	// tripped check 6; zero when check 6 never fired.
	BlindFor time.Duration
}

// proveCoverage applies the nine coverage checks (§6, §10, §14) to one rule's
// polls and is PURE: no HTTP, no files, no clock reads — everything it needs
// arrives as an argument, which is what lets §22's tests build []Poll literals
// instead of a fixture server (§2).
//
// polls need not be pre-filtered to this rule: proveCoverage selects by
// def.UID itself, exactly as Reduce selects by UID rather than by title
// (§14.5) — a caller handing it a whole log's polls must not have to
// pre-filter to get a correct answer.
//
// Every check always runs, even once an earlier one has already set
// Unobservable: LargestGap and the notes are diagnostics an operator reads on
// exit 2 regardless of which check actually failed (§20.2). Reason names the
// FIRST check, in the order below, that failed; a later failure still adds
// its own Note.
func proveCoverage(h Header, polls []Poll, sentinel *time.Time, t ruleTimings, def Definition,
	from, to time.Time, grace time.Duration) CoverageResult {

	windowEnd := to.Add(grace)

	var rulePolls []Poll
	for _, p := range polls {
		if p.RuleUID == def.UID {
			rulePolls = append(rulePolls, p)
		}
	}
	// Stable, not sort.Slice: two polls sharing a GrafanaNow (a coarse Date
	// header, or a corrupted/replayed log) must not reorder nondeterministically
	// in a function that promises to be pure.
	slices.SortStableFunc(rulePolls, func(a, b Poll) int { return a.GrafanaNow.Compare(b.GrafanaNow) })

	var res CoverageResult
	fail := func(reason UnobservableReason, note string) {
		res.Unobservable = true
		if res.Reason == "" {
			res.Reason = reason
		}
		res.Notes = append(res.Notes, fmt.Sprintf("rule %q: %s", def.Title, note))
	}

	// Check 1 — sentinel (§4.5). Present and At >= to+grace -> coverage
	// provable; absent, or short of it, is never a pass. A recorder that died
	// early must look exactly like a coverage gap, because it is one.
	switch {
	case sentinel == nil:
		fail(ReasonNoSentinel, "no stopped sentinel: the recorder never reported finishing")
	case sentinel.Before(windowEnd):
		fail(ReasonSentinelEarly, fmt.Sprintf("stopped sentinel at %s is before the required %s (to+grace)",
			sentinel.Format(time.RFC3339), windowEnd.Format(time.RFC3339)))
	}

	// Check 2 — from bounds (§7), first sentence only: from < StartedAt makes
	// coverage unprovable, no matter how healthy the polls that DO exist look.
	// Both are runner-domain clock reads (the recorder's own Clock.Now()), so
	// no cross-domain translation applies here. The second sentence — from
	// more than fromFutureTolerance ahead is a hard error — is Check's input
	// validation, once per run rather than per rule, and belongs to a later
	// phase: this function has no error return, only a per-rule verdict.
	if from.Before(h.StartedAt) {
		fail(ReasonFromBeforeRecord, fmt.Sprintf(
			"requested from %s is before recording started at %s", from.Format(time.RFC3339), h.StartedAt.Format(time.RFC3339)))
	}

	// Filtered once, here, and threaded through every remaining check —
	// ruleHeartbeatGap included — rather than re-filtered per check: two
	// independent filters over the same polls would only invite one of them
	// drifting from the other's membership test.
	inWindow := inWindowPolls(rulePolls, from, windowEnd)

	// Check 3 — heartbeat continuity (§6). Data at both ends with a hole in
	// between is not enough (§22.4): this scans every gap inside the window,
	// not just its edges.
	res.LargestGap, res.LargestGapAt = ruleHeartbeatGap(inWindow, from, windowEnd)
	if res.LargestGap > t.maxGap {
		fail(ReasonHeartbeatGap, fmt.Sprintf(
			"gap of %s starting at %s exceeds maxGap %s", res.LargestGap, res.LargestGapAt.Format(time.RFC3339), t.maxGap))
	}

	// Check 4 — health=="error" (§10.1). A short blip is a note only (§22.1:
	// one failed evaluation must not exit 2 over an otherwise clean window);
	// only a run longer than healthGrace consumes coverage.
	if runLen, sawAny := longestHealthRun(inWindow, "error"); sawAny {
		res.Notes = append(res.Notes, fmt.Sprintf("rule %q: health=error observed (longest run %s)", def.Title, runLen))
		if runLen > t.healthGrace {
			fail(ReasonHealthError, fmt.Sprintf("health=error for %s exceeds healthGrace %s", runLen, t.healthGrace))
		}
	}

	// Check 5 — health=="nodata" (§10.1/§10.2). Never fatal here: 96% of the
	// fleet runs no_data_state:OK, so treating this as fatal by default would
	// block nearly every healthy deploy in an idle environment. Escalating it
	// under Policy.NodataIsUnobservable is decide's job (a later phase),
	// applied directly against the raw polls — this pure function has no
	// Policy to consult and must not invent one.
	if _, sawAny := longestHealthRun(inWindow, "nodata"); sawAny {
		res.Notes = append(res.Notes, fmt.Sprintf("rule %q: health=nodata observed (not fatal; see --nodata-is-unobservable)", def.Title))
	}

	// Check 6 — liveness (H3). Absolute only, per poll: GrafanaNow and
	// LastEvaluation are both Grafana-domain reads off the SAME response, so
	// this is a same-domain comparison and uses raw values — never a delta
	// against a previous poll, which reports stale on ~half the polls of a
	// perfectly healthy rule (polling runs at intervalSeconds/2).
	//
	// Skipped only for a poll whose own flags SAY there is nothing to check:
	// IsPaused (a zero LastEvaluation is legal only while paused, §2.3; check
	// 7 is its detector) or !Found (no rule, no evaluation; check 8 is its
	// detector). Deliberately NOT skipped merely because LastEvaluation is
	// zero: ReadLog does no field validation, so a corrupted or hand-edited
	// log line can claim found:true, is_paused:false and still carry a zero
	// LastEvaluation, and that combination must read as maximally stale
	// rather than being silently waved through.
	var staleCount int
	var worstStale time.Duration
	var worstStaleAt time.Time
	for _, p := range inWindow {
		if p.IsPaused || !p.Found {
			continue
		}
		if stale := p.GrafanaNow.Sub(p.LastEvaluation); stale > t.evalStaleAfter {
			staleCount++
			if stale > worstStale {
				worstStale, worstStaleAt = stale, p.GrafanaNow
			}
		}
	}
	if staleCount > 0 {
		res.BlindFor = worstStale
		fail(ReasonStaleEvaluation, fmt.Sprintf(
			"lastEvaluation stale on %d poll(s); worst %s (> evalStaleAfter %s) as of %s",
			staleCount, worstStale, t.evalStaleAfter, worstStaleAt.Format(time.RFC3339)))
	}

	// Check 7 — isPaused in-window (§12.2, §14.8). The PRIMARY pause
	// detector: liveness (check 6) is only the backup for what IsPaused
	// cannot show (a deleted rule, a stopped scheduler, a blocked
	// evaluation). This is what catches pause-then-unpause, which the drain
	// wait alone passes (§14.7).
	var pausedCount int
	var pausedAt time.Time
	for _, p := range inWindow {
		if p.IsPaused {
			pausedCount++
			if pausedAt.IsZero() {
				pausedAt = p.GrafanaNow
			}
		}
	}
	if pausedCount > 0 {
		fail(ReasonPausedInWindow, fmt.Sprintf("observed paused on %d poll(s), first at %s", pausedCount, pausedAt.Format(time.RFC3339)))
	}

	// Check 8 — rule absent (§14.5). Found==false is authoritative (P2
	// already retried every transport failure before a Poll record ever
	// exists): the rule resolved at resolve time but the state endpoint
	// stopped serving it. Never drop a watched rule from the verdict set
	// silently.
	var absentCount int
	var absentAt time.Time
	for _, p := range inWindow {
		if !p.Found {
			absentCount++
			if absentAt.IsZero() {
				absentAt = p.GrafanaNow
			}
		}
	}
	if absentCount > 0 {
		fail(ReasonRuleAbsent, fmt.Sprintf("state endpoint returned no rule on %d poll(s), first at %s", absentCount, absentAt.Format(time.RFC3339)))
	}

	// Check 9 — KeepLast (§10.2). A note, never fatal. It surfaces only as an
	// instance Reason after P1.2a's parsing, and Reasons keys can be
	// comma-joined composites, so membership (reasonsContain) is required —
	// indexing "KeepLast" directly would miss "KeepLast, MissingSeries".
	for _, p := range inWindow {
		if reasonsContain(p.Reasons, keepLastReason) {
			res.Notes = append(res.Notes, fmt.Sprintf(
				"rule %q: KeepLast observed at %s: a held-over state may hide a real blind spot", def.Title, p.GrafanaNow.Format(time.RFC3339)))
			break
		}
	}

	res.Proved = !res.Unobservable
	return res
}

// inWindowPolls filters polls to those inside [from, windowEnd] using the
// CROSS-DOMAIN membership test (§16): each poll's Grafana-domain GrafanaNow
// is translated to the runner domain by its OWN skew, and its own skew bound
// is the membership tolerance, so a poll that is genuinely inside the window
// is never excluded by ordinary clock imprecision.
//
// Everything downstream of this filter (health runs, liveness, pause,
// absence) reads the poll's raw fields: GrafanaNow paired with
// LastEvaluation on the SAME response, or one poll's GrafanaNow against the
// next's, are same-domain comparisons and need no translation (§16, "Clock
// domains" — only window membership and check 3's two boundary segments do).
func inWindowPolls(polls []Poll, from, windowEnd time.Time) []Poll {
	var out []Poll
	for _, p := range polls {
		bound := p.SkewBound()
		runner := p.GrafanaNow.Add(-p.Skew())
		if runner.Before(from.Add(-bound)) || runner.After(windowEnd.Add(bound)) {
			continue
		}
		out = append(out, p)
	}
	return out
}

// ruleHeartbeatGap finds the largest unobserved span inside [from, windowEnd]
// (§6), including the two boundary segments — which is why "data at both
// ends with a hole in the middle" still fails (§22.4): the segment between
// the polls just inside each edge is exactly what this measures. in must
// already be filtered to this window (inWindowPolls) and sorted by
// GrafanaNow — proveCoverage computes that filter once and threads it through
// every check, this one included, rather than each check re-filtering.
//
// The two boundary segments compare a Grafana-domain poll time against the
// runner-domain from/windowEnd, so each is translated by its own poll's skew
// AND widened by that same poll's skew bound (§16: "with that poll's bound as
// the tolerance") — on the side that makes the segment larger, never smaller,
// so an uncertain boundary reads as at least as big a gap as it might really
// be. Understating it by up to the bound would be fail-open. The spacing
// BETWEEN consecutive polls compares two Grafana-domain reads to each other —
// same domain — and uses the raw GrafanaNow difference, no bound needed.
func ruleHeartbeatGap(in []Poll, from, windowEnd time.Time) (largestGap time.Duration, largestGapAt time.Time) {
	if len(in) == 0 {
		return windowEnd.Sub(from), from
	}

	runnerOf := func(p Poll) time.Time { return p.GrafanaNow.Add(-p.Skew()) }

	first := in[0]
	if gap := runnerOf(first).Sub(from) + first.SkewBound(); gap > largestGap {
		largestGap, largestGapAt = gap, from
	}
	for i := 1; i < len(in); i++ {
		if gap := in[i].GrafanaNow.Sub(in[i-1].GrafanaNow); gap > largestGap {
			largestGap, largestGapAt = gap, runnerOf(in[i-1])
		}
	}
	last := in[len(in)-1]
	if gap := windowEnd.Sub(runnerOf(last)) + last.SkewBound(); gap > largestGap {
		largestGap, largestGapAt = gap, runnerOf(last)
	}
	return largestGap, largestGapAt
}

// longestHealthRun returns the longest contiguous wall-clock span (§10.1)
// during which polls — already sorted by GrafanaNow, same-domain spacing
// (§16) — read the given rule-level Health, and whether any poll matched it
// at all.
//
// It detects the span as it accumulates rather than waiting for the run to
// end, so an open-ended run that is still failing at the last poll in the
// window is measured correctly without needing data past the window: waiting
// for the run to "end" would have to assume the best case about what happens
// next, which is exactly what this gate must not do (§1).
func longestHealthRun(polls []Poll, health string) (longest time.Duration, sawAny bool) {
	var runStart time.Time
	for _, p := range polls {
		if p.Health != health {
			runStart = time.Time{}
			continue
		}
		sawAny = true
		if runStart.IsZero() {
			runStart = p.GrafanaNow
		}
		if span := p.GrafanaNow.Sub(runStart); span > longest {
			longest = span
		}
	}
	return longest, sawAny
}

// reasonsContain reports whether any key of reasons names want, honoring
// Grafana's comma-joined composite reason strings via reasonNames (log.go).
func reasonsContain(reasons map[string]int, want string) bool {
	for reason := range reasons {
		if reasonNames(reason, want) {
			return true
		}
	}
	return false
}
