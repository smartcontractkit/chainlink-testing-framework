package gate

import (
	"fmt"
	"time"
)

// keepLastReason is the instance Reason that check 9 watches for.
const keepLastReason = "KeepLast"

// Two things this file leaves to its callers: the "from too far ahead" bound is
// Config.validate's once-per-run input validation (check 2 owns only the
// "from < StartedAt" half), and a rule paused at the window open never reaches
// proveCoverage — decide reads `skipped` from Header.pausedAtStart first, so a
// paused rule's zero polls read as skipped, not as one large heartbeat gap.

// UnobservableReason names why proveCoverage could not prove a rule's window.
// It reaches the JSON output, so it is a published vocabulary like Outcome;
// prose belongs in Notes.
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
	// ReasonDrainTimeout is set by check.go's drain wait, never by
	// proveCoverage: the wait is I/O and must not be added to this pure
	// function — that would put HTTP inside the pure layer and destroy the
	// seam this design depends on.
	ReasonDrainTimeout UnobservableReason = "drain_timeout"
)

// CoverageResult is proveCoverage's whole answer for one rule. No interval
// list: proved-or-not plus the largest gap and where is everything a human
// reads on exit 2, and everything the rendered table needs.
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

// proveCoverage applies the nine coverage checks to one rule's polls. PURE: no
// HTTP, no files, no clock reads — everything arrives as an argument. polls need
// not be pre-filtered to this rule (selection is by def.UID). Every check runs
// even after Unobservable is set, so LargestGap and the notes are complete on
// exit 2; Reason names only the FIRST check that failed.
func proveCoverage(h Header, polls []Poll, sentinel *time.Time, t ruleTimings, def Definition,
	from, to time.Time, grace time.Duration) CoverageResult {

	windowEnd := to.Add(grace)

	// pollsForRule (classify.go) is the single filter+sort implementation; this
	// and classifyRule must not carry two independent copies.
	rulePolls := pollsForRule(polls, def.UID)

	var res CoverageResult
	fail := func(reason UnobservableReason, note string) {
		res.Unobservable = true
		if res.Reason == "" {
			res.Reason = reason
		}
		res.Notes = append(res.Notes, fmt.Sprintf("rule %q: %s", def.Title, note))
	}

	// Check 1 — sentinel. Present and At >= to+grace -> coverage provable;
	// absent, or short of it, is never a pass. A recorder that died early must
	// look exactly like a coverage gap, because it is one.
	switch {
	case sentinel == nil:
		fail(ReasonNoSentinel, "no stopped sentinel: the recorder never reported finishing")
	case sentinel.Before(windowEnd):
		fail(ReasonSentinelEarly, fmt.Sprintf("stopped sentinel at %s is before the required %s (to+grace)",
			sentinel.Format(time.RFC3339), windowEnd.Format(time.RFC3339)))
	}

	// Check 2 — from bounds: from < StartedAt makes coverage unprovable, no
	// matter how healthy the polls that DO exist look. Both are runner-domain
	// clock reads (the recorder's own Clock.Now()), so no cross-domain
	// translation applies here. The other half of the bound — from too far
	// ahead of the runner's clock — is Check's input validation, once per run
	// rather than per rule.
	if from.Before(h.StartedAt) {
		fail(ReasonFromBeforeRecord, fmt.Sprintf(
			"requested from %s is before recording started at %s", from.Format(time.RFC3339), h.StartedAt.Format(time.RFC3339)))
	}

	// Filtered once and threaded through every remaining check.
	inWindow := inWindowPolls(rulePolls, from, windowEnd)

	// Check 3 — heartbeat continuity. Data at both ends with a hole in between
	// is not enough: this scans every gap inside the window, not just its
	// edges.
	res.LargestGap, res.LargestGapAt = ruleHeartbeatGap(inWindow, from, windowEnd)
	if res.LargestGap > t.maxGap {
		fail(ReasonHeartbeatGap, fmt.Sprintf(
			"gap of %s starting at %s exceeds maxGap %s", res.LargestGap, res.LargestGapAt.Format(time.RFC3339), t.maxGap))
	}

	// Check 4 — health=="error". A short blip is a note only — one failed
	// evaluation must not exit 2 over an otherwise clean window; only a run
	// longer than healthGrace consumes coverage.
	if runLen, sawAny := longestHealthRun(inWindow, "error"); sawAny {
		res.Notes = append(res.Notes, fmt.Sprintf("rule %q: health=error observed (longest run %s)", def.Title, runLen))
		if runLen > t.healthGrace {
			fail(ReasonHealthError, fmt.Sprintf("health=error for %s exceeds healthGrace %s", runLen, t.healthGrace))
		}
	}

	// Check 5 — health=="nodata". Never fatal here (most of the fleet runs
	// no_data_state:OK, so it would block healthy idle deploys). Escalating
	// under Policy.NodataIsUnobservable is decide's job, since this pure
	// function has no Policy to consult.
	if _, sawAny := longestHealthRun(inWindow, "nodata"); sawAny {
		res.Notes = append(res.Notes, fmt.Sprintf("rule %q: health=nodata observed (not fatal; see --nodata-is-unobservable)", def.Title))
	}

	// Check 6 — liveness. Same-domain (GrafanaNow and LastEvaluation are from
	// the SAME response), so raw values — never a delta against a previous
	// poll, which reports stale ~half the polls of a healthy rule. Skipped only
	// for IsPaused (check 7) or !Found (check 8); a zero LastEvaluation on a
	// found, unpaused poll is treated as maximally stale, not waved through.
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

	// Check 7 — isPaused in-window. The PRIMARY pause detector: liveness
	// (check 6) is only the backup for what IsPaused cannot show (a deleted
	// rule, a stopped scheduler, a blocked evaluation). This is what catches
	// pause-then-unpause, which the drain wait alone passes.
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

	// Check 8 — rule absent. Found==false is authoritative (the transport
	// already retried every transient failure before a Poll record ever
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

	// Check 9 — KeepLast. Two non-fatal notes: DECLARED (the rule is configured
	// with no_data_state/exec_err_state=KeepLast, read from def so it fires once),
	// and OBSERVED (an instance reported KeepLast in-window; Reasons keys can be
	// comma-joined, so membership via reasonsContain, never a literal index).
	if def.NoDataState == keepLastReason || def.ExecErrState == keepLastReason {
		res.Notes = append(res.Notes, fmt.Sprintf(
			"rule %q: configured with no_data_state/exec_err_state=KeepLast — a stale state can continue past a real fault", def.Title))
	}
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

// inWindowPolls filters to polls inside [from, windowEnd] via the cross-domain
// membership test: each GrafanaNow is translated to the runner domain by its
// own skew, widened by its skew bound, so clock imprecision never excludes a
// genuinely in-window poll. Everything downstream reads same-domain raw fields;
// only this filter and check 3's boundary segments cross domains.
func inWindowPolls(polls []Poll, from, windowEnd time.Time) []Poll {
	var out []Poll
	for _, p := range polls {
		bound := p.SkewBound()
		runner := runnerTime(p, p.GrafanaNow)
		if runner.Before(from.Add(-bound)) || runner.After(windowEnd.Add(bound)) {
			continue
		}
		out = append(out, p)
	}
	return out
}

// ruleHeartbeatGap finds the largest unobserved span inside [from, windowEnd],
// including the two boundary segments — which is why data at both ends with a
// hole in the middle still fails. in must be filtered (inWindowPolls) and
// sorted by GrafanaNow.
//
// Boundary segments are cross-domain, so each translated poll time is widened
// by its skew bound on the side that makes the gap LARGER (never smaller — an
// uncertain boundary must read as at least as big a gap as it might be).
// Consecutive-poll spacing is same-domain and uses the raw GrafanaNow diff.
func ruleHeartbeatGap(in []Poll, from, windowEnd time.Time) (largestGap time.Duration, largestGapAt time.Time) {
	if len(in) == 0 {
		return windowEnd.Sub(from), from
	}

	runnerOf := func(p Poll) time.Time { return runnerTime(p, p.GrafanaNow) }

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

// longestHealthRun returns the longest contiguous span of polls reading the
// given rule-level Health, measured incrementally so a run still failing at the
// last in-window poll is measured correctly without assuming anything past the
// window.
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
