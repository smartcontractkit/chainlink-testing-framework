package gate

import (
	"strings"
	"testing"
	"time"
)

func TestProveCoverage_CleanWindowIsProved(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(10 * time.Minute)
	rt := newRuleTimings(30*time.Second, 60)
	def := Definition{UID: "r1", Title: "R1"}

	var polls []Poll
	for ts := from; !ts.After(to); ts = ts.Add(30 * time.Second) {
		polls = append(polls, Poll{RuleUID: "r1", GrafanaNow: ts, Found: true, Health: "ok", State: "inactive", LastEvaluation: ts})
	}
	sentinel := to

	res := proveCoverage(Header{StartedAt: from.Add(-time.Hour)}, polls, &sentinel, rt, def, from, to, 0)
	if !res.Proved || res.Unobservable || res.Reason != "" {
		t.Fatalf("res = %+v, want a clean proved window", res)
	}
}

func TestProveCoverage_FiltersPollsByUID(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(10 * time.Minute)
	rt := newRuleTimings(30*time.Second, 60)
	def := Definition{UID: "r1", Title: "R1"}

	var polls []Poll
	for ts := from; !ts.After(to); ts = ts.Add(30 * time.Second) {
		polls = append(polls, Poll{RuleUID: "r1", GrafanaNow: ts, Found: true, Health: "ok", LastEvaluation: ts})
		// A different rule's polls, deliberately broken, must never
		// contaminate r1's verdict: proveCoverage selects by UID itself.
		polls = append(polls, Poll{RuleUID: "other", GrafanaNow: ts, Found: false})
	}
	sentinel := to

	res := proveCoverage(Header{StartedAt: from.Add(-time.Hour)}, polls, &sentinel, rt, def, from, to, 0)
	if !res.Proved {
		t.Fatalf("res = %+v, want proved: a different rule's broken polls must not affect this rule's verdict", res)
	}
}

// --- Check 1: sentinel ---

func TestProveCoverage_NoSentinelIsUnobservable(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(time.Minute)
	rt := newRuleTimings(30*time.Second, 60)
	def := Definition{UID: "r1", Title: "R1"}

	res := proveCoverage(Header{StartedAt: from.Add(-time.Hour)}, nil, nil, rt, def, from, to, 0)
	if res.Proved || res.Reason != ReasonNoSentinel {
		t.Fatalf("res = %+v, want unobservable/no_sentinel: an absent sentinel must never be a pass", res)
	}
}

func TestProveCoverage_SentinelBeforeGraceIsUnobservable(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(10 * time.Minute)
	grace := 2 * time.Minute
	rt := newRuleTimings(30*time.Second, 60)
	def := Definition{UID: "r1", Title: "R1"}

	sentinel := to.Add(grace).Add(-time.Second) // one second short of to+grace
	res := proveCoverage(Header{StartedAt: from.Add(-time.Hour)}, nil, &sentinel, rt, def, from, to, grace)
	if res.Reason != ReasonSentinelEarly {
		t.Fatalf("Reason = %q, want sentinel_early", res.Reason)
	}
	if !res.Unobservable || res.Proved {
		t.Fatalf("res = %+v, want Unobservable and not Proved — a reason string with no consequence is not a coverage failure", res)
	}

	// The consequence: decide() must turn this into exit 2, never a pass.
	defs := []Definition{def}
	drt := map[string]ruleTimings{def.UID: rt}
	gt := globalTimings{transitionGrace: grace}
	pol := Policy{From: from, To: to}
	dres, err := decide(Header{StartedAt: from.Add(-time.Hour)}, nil, &sentinel, defs, drt, gt, pol)
	if err == nil {
		t.Fatalf("decide() err = nil, want non-nil: a sentinel short of to+grace must fail the run")
	}
	if len(dres.Verdicts) != 1 || dres.Verdicts[0].Outcome != OutcomeUnobservable {
		t.Fatalf("Verdicts = %+v, want one unobservable verdict", dres.Verdicts)
	}
}

func TestProveCoverage_SentinelExactlyAtGraceIsFine(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(10 * time.Minute)
	grace := 2 * time.Minute
	windowEnd := to.Add(grace)
	rt := newRuleTimings(30*time.Second, 60)
	def := Definition{UID: "r1", Title: "R1"}

	var polls []Poll
	for ts := from; !ts.After(windowEnd); ts = ts.Add(30 * time.Second) {
		polls = append(polls, Poll{RuleUID: "r1", GrafanaNow: ts, Found: true, Health: "ok", LastEvaluation: ts})
	}
	sentinel := windowEnd

	res := proveCoverage(Header{StartedAt: from.Add(-time.Hour)}, polls, &sentinel, rt, def, from, to, grace)
	if !res.Proved {
		t.Fatalf("Proved = false, want true: sentinel exactly at to+grace must satisfy check 1: %+v", res)
	}
}

// --- Check 2: from bounds ---

func TestProveCoverage_FromBeforeRecordIsUnobservable(t *testing.T) {
	started := time.Date(2026, 1, 1, 1, 0, 0, 0, time.UTC)
	from := started.Add(-time.Minute) // the requested window opens before recording started
	to := from.Add(10 * time.Minute)
	rt := newRuleTimings(30*time.Second, 60)
	def := Definition{UID: "r1", Title: "R1"}

	sentinel := to
	res := proveCoverage(Header{StartedAt: started}, nil, &sentinel, rt, def, from, to, 0)
	if res.Reason != ReasonFromBeforeRecord {
		t.Fatalf("Reason = %q, want from_before_record", res.Reason)
	}
	if !res.Unobservable || res.Proved {
		t.Fatalf("res = %+v, want Unobservable and not Proved — a reason string with no consequence is not a coverage failure", res)
	}

	// The consequence: decide() must turn this into exit 2, never a pass.
	defs := []Definition{def}
	drt := map[string]ruleTimings{def.UID: rt}
	gt := globalTimings{}
	pol := Policy{From: from, To: to}
	dres, err := decide(Header{StartedAt: started}, nil, &sentinel, defs, drt, gt, pol)
	if err == nil {
		t.Fatalf("decide() err = nil, want non-nil: `from` before the recording started must fail the run")
	}
	if len(dres.Verdicts) != 1 || dres.Verdicts[0].Outcome != OutcomeUnobservable {
		t.Fatalf("Verdicts = %+v, want one unobservable verdict", dres.Verdicts)
	}
}

// --- Check 3: heartbeat continuity ---

// The core heartbeat regression: data at both ends with a hole between is not
// enough.
func TestProveCoverage_HeartbeatGapBetweenBoundariesIsUnobservable(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(10 * time.Minute)
	rt := newRuleTimings(30*time.Second, 60) // maxGap = 60s
	def := Definition{UID: "r1", Title: "R1"}

	polls := []Poll{
		{RuleUID: "r1", GrafanaNow: from.Add(time.Second), Found: true, Health: "ok", LastEvaluation: from},
		{RuleUID: "r1", GrafanaNow: to.Add(-time.Second), Found: true, Health: "ok", LastEvaluation: to},
	}
	sentinel := to

	res := proveCoverage(Header{StartedAt: from.Add(-time.Hour)}, polls, &sentinel, rt, def, from, to, 0)
	if res.Reason != ReasonHeartbeatGap {
		t.Fatalf("Reason = %q, want heartbeat_gap: healthy edges with a hole in the middle must still fail", res.Reason)
	}
	// The gap is the SPACING between the two polls (598s), not either
	// boundary segment (1s each) — pin the actual values, not just the verdict.
	if res.LargestGap != 598*time.Second {
		t.Fatalf("LargestGap = %s, want 598s (the spacing between the two polls, not a boundary segment)", res.LargestGap)
	}
	wantAt := from.Add(time.Second)
	if !res.LargestGapAt.Equal(wantAt) {
		t.Fatalf("LargestGapAt = %s, want %s (where the gap starts, at the first poll)", res.LargestGapAt, wantAt)
	}
}

// --- Check 4/5: health ---

func TestProveCoverage_HealthErrorShortBlipPassesWithNote(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(10 * time.Minute)
	rt := newRuleTimings(30*time.Second, 60) // healthGrace = 60s
	def := Definition{UID: "r1", Title: "R1"}

	var polls []Poll
	blip := from.Add(2 * time.Minute)
	for ts := from; !ts.After(to); ts = ts.Add(30 * time.Second) {
		health := "ok"
		if ts.Equal(blip) {
			health = "error" // one isolated failed evaluation
		}
		polls = append(polls, Poll{RuleUID: "r1", GrafanaNow: ts, Found: true, Health: health, LastEvaluation: ts})
	}
	sentinel := to

	res := proveCoverage(Header{StartedAt: from.Add(-time.Hour)}, polls, &sentinel, rt, def, from, to, 0)
	if !res.Proved {
		t.Fatalf("Proved = false, want true: one failed evaluation must not fail an otherwise clean window: %+v", res)
	}
	if !anyContains(res.Notes, "health=error") {
		t.Fatalf("Notes = %v, want a health=error note even though it did not fail the window", res.Notes)
	}
}

func TestProveCoverage_HealthErrorSustainedIsUnobservable(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(10 * time.Minute)
	rt := newRuleTimings(30*time.Second, 60) // healthGrace = 60s
	def := Definition{UID: "r1", Title: "R1"}

	runStart, runEnd := from.Add(2*time.Minute), from.Add(5*time.Minute) // a 3-minute run, well past healthGrace
	var polls []Poll
	for ts := from; !ts.After(to); ts = ts.Add(30 * time.Second) {
		health := "ok"
		if !ts.Before(runStart) && !ts.After(runEnd) {
			health = "error"
		}
		polls = append(polls, Poll{RuleUID: "r1", GrafanaNow: ts, Found: true, Health: health, LastEvaluation: ts})
	}
	sentinel := to

	res := proveCoverage(Header{StartedAt: from.Add(-time.Hour)}, polls, &sentinel, rt, def, from, to, 0)
	if res.Reason != ReasonHealthError {
		t.Fatalf("Reason = %q, want health_error for a run that outlasts healthGrace", res.Reason)
	}
}

func TestProveCoverage_HealthNodataNeverFatalHere(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(10 * time.Minute)
	rt := newRuleTimings(30*time.Second, 60)
	def := Definition{UID: "r1", Title: "R1"}

	var polls []Poll
	for ts := from; !ts.After(to); ts = ts.Add(30 * time.Second) {
		polls = append(polls, Poll{RuleUID: "r1", GrafanaNow: ts, Found: true, Health: "nodata", LastEvaluation: ts})
	}
	sentinel := to

	res := proveCoverage(Header{StartedAt: from.Add(-time.Hour)}, polls, &sentinel, rt, def, from, to, 0)
	if !res.Proved {
		t.Fatalf("Proved = false, want true: health=nodata for the WHOLE window must still not be fatal by itself "+
			"(escalating it is Policy.NodataIsUnobservable's job, applied by decide): %+v", res)
	}
	if !anyContains(res.Notes, "health=nodata") {
		t.Fatalf("Notes = %v, want a health=nodata note", res.Notes)
	}
}

// --- Check 6: liveness ---

// A healthy rule polled at intervalSeconds/2, across the full window, must
// show zero staleness violations. lastEvaluation only advances once per full
// evaluation interval here — the realistic shape a delta check misreads as
// stale on roughly half of all polls.
func TestProveCoverage_LivenessAbsoluteNeverFalseStale(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	pollEvery := 30 * time.Second
	intervalSeconds := 60
	windowEnd := from.Add(10 * time.Minute)
	rt := newRuleTimings(pollEvery, intervalSeconds)
	def := Definition{UID: "r1", Title: "R1"}

	var polls []Poll
	lastEval := from
	for ts := from; !ts.After(windowEnd); ts = ts.Add(pollEvery) {
		if ts.Sub(lastEval) >= time.Duration(intervalSeconds)*time.Second {
			lastEval = ts
		}
		polls = append(polls, Poll{RuleUID: "r1", GrafanaNow: ts, Found: true, Health: "ok", LastEvaluation: lastEval})
	}
	sentinel := windowEnd

	res := proveCoverage(Header{StartedAt: from.Add(-time.Hour)}, polls, &sentinel, rt, def, from, windowEnd, 0)
	if res.Reason == ReasonStaleEvaluation || res.BlindFor != 0 {
		t.Fatalf("proveCoverage flagged staleness on a healthy rule polled at intervalSeconds/2 — liveness must be absolute, "+
			"never a delta against a previous poll: %+v", res)
	}
	if !res.Proved {
		t.Fatalf("Proved = false, want true: %+v (notes: %v)", res, res.Notes)
	}
}

func TestProveCoverage_StaleEvaluationIsUnobservable(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(10 * time.Minute)
	rt := newRuleTimings(30*time.Second, 60) // evalStaleAfter = 120s
	def := Definition{UID: "r1", Title: "R1"}

	// Dense, otherwise-healthy polling so heartbeat continuity (check 3)
	// stays intact — only check 6 should be able to fire.
	polls := denseHealthyPolls("r1", from, to, 30*time.Second)
	staleAt := from.Add(5 * time.Minute)
	for i := range polls {
		if polls[i].GrafanaNow.Equal(staleAt) {
			polls[i].LastEvaluation = staleAt.Add(-3 * time.Minute)
		}
	}
	sentinel := to

	res := proveCoverage(Header{StartedAt: from.Add(-time.Hour)}, polls, &sentinel, rt, def, from, to, 0)
	if res.Reason != ReasonStaleEvaluation {
		t.Fatalf("Reason = %q, want stale_evaluation", res.Reason)
	}
	if res.BlindFor != 3*time.Minute {
		t.Fatalf("BlindFor = %s, want 3m", res.BlindFor)
	}
}

func TestProveCoverage_ZeroLastEvaluationNeverFalseStale(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(10 * time.Minute)
	rt := newRuleTimings(30*time.Second, 60)
	def := Definition{UID: "r1", Title: "R1"}

	// A paused rule legitimately reports the zero time; check 6 must not read
	// that as an enormous staleness violation. Check 7 is its detector.
	polls := []Poll{
		{RuleUID: "r1", GrafanaNow: from.Add(time.Minute), Found: true, IsPaused: true},
	}
	sentinel := to

	res := proveCoverage(Header{StartedAt: from.Add(-time.Hour)}, polls, &sentinel, rt, def, from, to, 0)
	if res.Reason == ReasonStaleEvaluation {
		t.Fatalf("a zero lastEvaluation on a paused poll must not trigger check 6: %+v", res)
	}
}

// --- Check 7: isPaused in-window ---

func TestProveCoverage_PausedInWindowIsUnobservable(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(10 * time.Minute)
	rt := newRuleTimings(30*time.Second, 60)
	def := Definition{UID: "r1", Title: "R1"}

	// Dense, otherwise-healthy polling so heartbeat continuity (check 3)
	// stays intact — only check 7 should be able to fire.
	polls := denseHealthyPolls("r1", from, to, 30*time.Second)
	pausedAt := from.Add(5 * time.Minute)
	for i := range polls {
		if polls[i].GrafanaNow.Equal(pausedAt) {
			polls[i].IsPaused = true
			polls[i].LastEvaluation = time.Time{} // legal only while paused
		}
	}
	sentinel := to

	res := proveCoverage(Header{StartedAt: from.Add(-time.Hour)}, polls, &sentinel, rt, def, from, to, 0)
	if res.Reason != ReasonPausedInWindow {
		t.Fatalf("Reason = %q, want paused_in_window", res.Reason)
	}
}

// TestProveCoverage_PausedAfterWindowIsFine pins check 7's respect for the
// window boundary: a poll that reports paused but lands beyond windowEnd (a
// rule paused only after THIS release window closed) is filtered out by
// inWindowPolls and must not fail the window. Without that filter, a pause in
// the next release's window would wrongly fail this one.
func TestProveCoverage_PausedAfterWindowIsFine(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(10 * time.Minute)
	rt := newRuleTimings(30*time.Second, 60)
	def := Definition{UID: "r1", Title: "R1"}

	polls := denseHealthyPolls("r1", from, to, 30*time.Second)
	polls = append(polls, Poll{
		RuleUID: "r1", GrafanaNow: to.Add(2 * time.Minute),
		Found: true, Health: "ok", IsPaused: true,
	})
	sentinel := to

	res := proveCoverage(Header{StartedAt: from.Add(-time.Hour)}, polls, &sentinel, rt, def, from, to, 0)
	if res.Reason == ReasonPausedInWindow {
		t.Fatalf("a paused poll after windowEnd tripped check 7: %+v", res.Notes)
	}
	if !res.Proved {
		t.Fatalf("Proved = false, want a clean window: %+v", res.Notes)
	}
}

// --- Check 8: rule absent ---

func TestProveCoverage_RuleAbsentIsUnobservable(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(10 * time.Minute)
	rt := newRuleTimings(30*time.Second, 60)
	def := Definition{UID: "r1", Title: "R1"}

	// Dense, otherwise-healthy polling so heartbeat continuity (check 3)
	// stays intact — only check 8 should be able to fire.
	polls := denseHealthyPolls("r1", from, to, 30*time.Second)
	absentAt := from.Add(5 * time.Minute)
	for i := range polls {
		if polls[i].GrafanaNow.Equal(absentAt) {
			polls[i].Found = false
			polls[i].Health = ""
			polls[i].LastEvaluation = time.Time{}
		}
	}
	sentinel := to

	res := proveCoverage(Header{StartedAt: from.Add(-time.Hour)}, polls, &sentinel, rt, def, from, to, 0)
	if res.Reason != ReasonRuleAbsent {
		t.Fatalf("Reason = %q, want rule_absent", res.Reason)
	}
}

// denseHealthyPolls builds a clean poll sequence at a fixed cadence, with
// zero staleness and nothing abnormal — the baseline the single-check tests
// mutate exactly one poll of, so heartbeat continuity (check 3) never
// confounds the check under test.
func denseHealthyPolls(uid string, from, to time.Time, every time.Duration) []Poll {
	var out []Poll
	for ts := from; !ts.After(to); ts = ts.Add(every) {
		out = append(out, Poll{RuleUID: uid, GrafanaNow: ts, Found: true, Health: "ok", LastEvaluation: ts})
	}
	return out
}

// --- Check 9: KeepLast ---

func TestProveCoverage_KeepLastObservedIsNoteOnly(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(10 * time.Minute)
	rt := newRuleTimings(30*time.Second, 60)
	def := Definition{UID: "r1", Title: "R1"}

	var polls []Poll
	for ts := from; !ts.After(to); ts = ts.Add(30 * time.Second) {
		polls = append(polls, Poll{
			RuleUID: "r1", GrafanaNow: ts, Found: true, Health: "ok", LastEvaluation: ts,
			// A comma-joined composite — reasonsContain must match by
			// membership, never by an exact key.
			Reasons: map[string]int{"KeepLast, MissingSeries": 1},
		})
	}
	sentinel := to

	res := proveCoverage(Header{StartedAt: from.Add(-time.Hour)}, polls, &sentinel, rt, def, from, to, 0)
	if !res.Proved {
		t.Fatalf("Proved = false, want true: KeepLast is a note, never fatal: %+v", res)
	}
	if !anyContains(res.Notes, "KeepLast") {
		t.Fatalf("Notes = %v, want a KeepLast note (comma-joined membership, not a literal-key match)", res.Notes)
	}
}

// KeepLast in the CONFIGURATION gives a note — a different claim from the
// observed-reason test above. A rule DECLARED with
// no_data_state or exec_err_state = KeepLast is a standing blind spot
// whether or not any poll ever actually reports the reason, so the note
// must fire off the definition alone, over an otherwise perfectly healthy
// window with zero KeepLast reasons anywhere in it.
func TestProveCoverage_KeepLastConfiguredIsNoteOnly(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(10 * time.Minute)
	rt := newRuleTimings(30*time.Second, 60)

	tests := []struct {
		name string
		def  Definition
	}{
		{"no_data_state", Definition{UID: "r1", Title: "R1", NoDataState: "KeepLast"}},
		{"exec_err_state", Definition{UID: "r1", Title: "R1", ExecErrState: "KeepLast"}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			polls := denseHealthyPolls("r1", from, to, 30*time.Second)
			sentinel := to

			res := proveCoverage(Header{StartedAt: from.Add(-time.Hour)}, polls, &sentinel, rt, tc.def, from, to, 0)
			if !res.Proved {
				t.Fatalf("Proved = false, want true: a declared KeepLast is a note, never fatal: %+v", res)
			}
			if !anyContains(res.Notes, "KeepLast") {
				t.Fatalf("Notes = %v, want a KeepLast note from the definition alone, with zero KeepLast reasons observed", res.Notes)
			}
		})
	}
}

// --- Clock domains ---

// A constant clock skew on every poll must not itself read as a coverage gap
// or a from-before-record violation, because every cross-domain comparison
// translates by that poll's own skew first.
func TestProveCoverage_SkewTranslationAtWindowBoundary(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(10 * time.Minute)
	rt := newRuleTimings(30*time.Second, 60)
	def := Definition{UID: "r1", Title: "R1"}

	const skew = 45 * time.Second // Grafana's clock reads 45s ahead of the runner's
	const bound = 5 * time.Second

	var polls []Poll
	for ts := from; !ts.After(to); ts = ts.Add(30 * time.Second) {
		grafanaTime := ts.Add(skew)
		polls = append(polls, Poll{
			RuleUID: "r1", GrafanaNow: grafanaTime, SkewMS: skew.Milliseconds(), SkewBoundMS: bound.Milliseconds(),
			Found: true, Health: "ok", LastEvaluation: grafanaTime,
		})
	}
	sentinel := to

	res := proveCoverage(Header{StartedAt: from.Add(-time.Hour)}, polls, &sentinel, rt, def, from, to, 0)
	if !res.Proved {
		t.Fatalf("res = %+v, want proved: a constant clock skew must not itself read as a coverage gap", res)
	}
}

// --- Override round-trip: one authority for the cadence ---

// This exercises DeriveTimingsFromLog and proveCoverage together, exactly as
// check does, to prove maxGap tracks the RECORDED cadence and never a
// re-derivation from the rule's own evaluation interval.
func TestProveCoverage_OverrideRoundTrip(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	t.Run("slower override on a tighter rule classifies clean", func(t *testing.T) {
		windowEnd := from.Add(10 * time.Minute)
		h := Header{
			StartedAt: from.Add(-time.Hour),
			Rules:     []LoggedRule{{UID: "r1", Title: "R1", IntervalSeconds: 60, PollEverySeconds: 120}},
		}
		defs := []Definition{{UID: "r1", Title: "R1", IntervalSeconds: 60}}
		rt, _, err := DeriveTimingsFromLog(h, defs)
		if err != nil {
			t.Fatalf("DeriveTimingsFromLog: %v", err)
		}

		var polls []Poll
		for ts := from; !ts.After(windowEnd); ts = ts.Add(120 * time.Second) {
			polls = append(polls, Poll{RuleUID: "r1", GrafanaNow: ts, Found: true, Health: "ok", LastEvaluation: ts})
		}
		sentinel := windowEnd

		res := proveCoverage(h, polls, &sentinel, rt["r1"], defs[0], from, windowEnd, 0)
		if !res.Proved {
			t.Fatalf("Proved = false, want true (maxGap must come from the recorded 120s cadence, not the 30s default): %+v", res)
		}
	})

	t.Run("faster override still catches a real recorder gap", func(t *testing.T) {
		windowEnd := from.Add(20 * time.Minute)
		h := Header{
			StartedAt: from.Add(-time.Hour),
			Rules:     []LoggedRule{{UID: "r1", Title: "R1", IntervalSeconds: 300, PollEverySeconds: 5}},
		}
		defs := []Definition{{UID: "r1", Title: "R1", IntervalSeconds: 300}}
		rt, _, err := DeriveTimingsFromLog(h, defs)
		if err != nil {
			t.Fatalf("DeriveTimingsFromLog: %v", err)
		}

		var polls []Poll
		ts := from
		for range 20 {
			polls = append(polls, Poll{RuleUID: "r1", GrafanaNow: ts, Found: true, Health: "ok", LastEvaluation: ts})
			ts = ts.Add(5 * time.Second)
		}
		// The one real gap: 250s, nowhere near this recording's actual 5s
		// cadence. Resume 5s polling afterward all the way to windowEnd, so
		// this hole is the ONLY gap in the window — otherwise an uncovered
		// tail would exceed even the WRONG (definition-derived) 300s maxGap
		// on its own, and the test could not tell the two derivations apart.
		ts = ts.Add(250 * time.Second)
		for !ts.After(windowEnd) {
			polls = append(polls, Poll{RuleUID: "r1", GrafanaNow: ts, Found: true, Health: "ok", LastEvaluation: ts})
			ts = ts.Add(5 * time.Second)
		}
		sentinel := windowEnd

		res := proveCoverage(h, polls, &sentinel, rt["r1"], defs[0], from, windowEnd, 0)
		if res.Reason != ReasonHeartbeatGap {
			t.Fatalf("Reason = %q, want heartbeat_gap: if maxGap had been re-derived from the 300s definition instead of "+
				"the recorded 5s cadence, this 250s gap would pass silently — the fail-open direction", res.Reason)
		}
	})
}

func anyContains(notes []string, substr string) bool {
	for _, n := range notes {
		if strings.Contains(n, substr) {
			return true
		}
	}
	return false
}

// --- Check 6, tightened: a corrupted log must not silently disable liveness ---

// TestProveCoverage_ZeroLastEvaluationWithoutPauseIsStale guards check 6's
// skip condition. ReadLog does no field validation, so a log line can claim
// found:true, is_paused:false and still carry a zero LastEvaluation (a
// corrupted write, a hand-edited fixture, a future log format bug). That
// combination must read as maximally stale, not be waved through the way a
// legitimately paused poll's zero time is — the skip must key off
// IsPaused/Found, never off LastEvaluation being zero.
func TestProveCoverage_ZeroLastEvaluationWithoutPauseIsStale(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(10 * time.Minute)
	rt := newRuleTimings(30*time.Second, 60)
	def := Definition{UID: "r1", Title: "R1"}

	polls := denseHealthyPolls("r1", from, to, 30*time.Second)
	corruptAt := from.Add(5 * time.Minute)
	for i := range polls {
		if polls[i].GrafanaNow.Equal(corruptAt) {
			polls[i].LastEvaluation = time.Time{} // found:true, is_paused:false, yet zero
		}
	}
	sentinel := to

	res := proveCoverage(Header{StartedAt: from.Add(-time.Hour)}, polls, &sentinel, rt, def, from, to, 0)
	if res.Reason != ReasonStaleEvaluation {
		t.Fatalf("Reason = %q, want stale_evaluation: a zero lastEvaluation on a found, non-paused poll must fail "+
			"closed, not be silently skipped as if it were a legitimately paused observation", res.Reason)
	}
}

// --- Check 3, tightened: the boundary segments must widen by the skew bound ---

// The two boundary segments take their own poll's bound as the tolerance: a
// boundary gap that lands EXACTLY at maxGap must still fail once the poll's
// own skew bound is added, because the translation is only a best
// estimate and understating the gap by up to the bound would be fail-open.
func TestProveCoverage_BoundaryGapWidensBySkewBound(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(10 * time.Minute)
	rt := newRuleTimings(30*time.Second, 60) // maxGap = 60s
	def := Definition{UID: "r1", Title: "R1"}

	const bound = 5 * time.Second
	first := Poll{
		RuleUID: "r1", GrafanaNow: from.Add(rt.maxGap), Found: true, Health: "ok",
		LastEvaluation: from.Add(rt.maxGap), SkewBoundMS: bound.Milliseconds(),
	}
	rest := denseHealthyPolls("r1", from.Add(rt.maxGap+30*time.Second), to, 30*time.Second)
	polls := append([]Poll{first}, rest...)
	sentinel := to

	res := proveCoverage(Header{StartedAt: from.Add(-time.Hour)}, polls, &sentinel, rt, def, from, to, 0)
	if res.Reason != ReasonHeartbeatGap {
		t.Fatalf("Reason = %q, want heartbeat_gap: the leading boundary segment sits at EXACTLY maxGap (60s) before "+
			"widening; the poll's own %s skew bound must push it past the threshold, not just the skew translation", res.Reason, bound)
	}
}

// --- Multi-failure contract ---

// Two checks failing in the same rule: check 7 (paused in-window) runs before
// check 8 (rule absent), so Reason must name the pause even though the rule
// also goes absent later — and the later failure must still add its own Note
// rather than being swallowed once Reason is set.
func TestProveCoverage_MultipleFailuresReasonIsFirstButAllNoted(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(10 * time.Minute)
	rt := newRuleTimings(30*time.Second, 60)
	def := Definition{UID: "r1", Title: "R1"}

	polls := denseHealthyPolls("r1", from, to, 30*time.Second)
	pausedAt := from.Add(3 * time.Minute)
	absentAt := from.Add(6 * time.Minute)
	for i := range polls {
		switch {
		case polls[i].GrafanaNow.Equal(pausedAt):
			polls[i].IsPaused = true
			polls[i].LastEvaluation = time.Time{}
		case polls[i].GrafanaNow.Equal(absentAt):
			polls[i].Found = false
			polls[i].Health = ""
			polls[i].LastEvaluation = time.Time{}
		}
	}
	sentinel := to

	res := proveCoverage(Header{StartedAt: from.Add(-time.Hour)}, polls, &sentinel, rt, def, from, to, 0)
	if res.Reason != ReasonPausedInWindow {
		t.Fatalf("Reason = %q, want paused_in_window: the FIRST check to fail names the reason", res.Reason)
	}
	if !anyContains(res.Notes, "paused") {
		t.Fatalf("Notes = %v, want a note about the pause", res.Notes)
	}
	if !anyContains(res.Notes, "no rule") {
		t.Fatalf("Notes = %v, want a note about the absence too — a later failure must still be recorded, "+
			"not swallowed once Reason is already set", res.Notes)
	}
}

// --- Skipped rules ---

// A known limit of this function's contract, not a bug in it: a rule paused
// BEFORE the window opened is never scheduled or polled (watch.go), so it
// reaches proveCoverage with zero polls at all. proveCoverage has no notion of
// "skipped" — that classification belongs to the definitions
// (LoggedRule.IsPaused / Definition.IsPaused), never to the polls — so it
// reports the whole window as one big heartbeat_gap instead. decide is what
// reads skipped status from the header and never calls this function for such
// a rule; this pins the behavior it relies on not reaching.
func TestProveCoverage_SkippedRuleWithZeroPollsPinnedAsHeartbeatGap(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(10 * time.Minute)
	rt := newRuleTimings(30*time.Second, 60)
	def := Definition{UID: "r1", Title: "R1", IsPaused: true}

	sentinel := to
	res := proveCoverage(Header{StartedAt: from.Add(-time.Hour)}, nil, &sentinel, rt, def, from, to, 0)
	if res.Reason != ReasonHeartbeatGap {
		t.Fatalf("Reason = %q, want heartbeat_gap (pinned, not the desired end state): proveCoverage has no "+
			"'skipped' concept, so decide must handle a skipped rule's classification itself, before or "+
			"instead of calling this function", res.Reason)
	}
}
