package gate

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
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
	require.True(t, res.Proved)
	require.False(t, res.Unobservable)
	require.Empty(t, res.Reason)
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
	require.True(t, res.Proved, "a different rule's broken polls must not affect this rule's verdict")
}

// --- Check 1: sentinel ---

func TestProveCoverage_NoSentinelIsUnobservable(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(time.Minute)
	rt := newRuleTimings(30*time.Second, 60)
	def := Definition{UID: "r1", Title: "R1"}

	res := proveCoverage(Header{StartedAt: from.Add(-time.Hour)}, nil, nil, rt, def, from, to, 0)
	require.False(t, res.Proved)
	require.Equal(t, ReasonNoSentinel, res.Reason, "an absent sentinel must never be a pass")
}

func TestProveCoverage_SentinelBeforeGraceIsUnobservable(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(10 * time.Minute)
	grace := 2 * time.Minute
	rt := newRuleTimings(30*time.Second, 60)
	def := Definition{UID: "r1", Title: "R1"}

	sentinel := to.Add(grace).Add(-time.Second) // one second short of to+grace
	res := proveCoverage(Header{StartedAt: from.Add(-time.Hour)}, nil, &sentinel, rt, def, from, to, grace)
	require.Equal(t, ReasonSentinelEarly, res.Reason)
	require.True(t, res.Unobservable)
	require.False(t, res.Proved, "a reason string with no consequence is not a coverage failure")

	// The consequence: decide() must turn this into exit 2, never a pass.
	defs := []Definition{def}
	drt := map[string]ruleTimings{def.UID: rt}
	gt := globalTimings{transitionGrace: grace}
	pol := Policy{From: from, To: to}
	dres, err := decide(Header{StartedAt: from.Add(-time.Hour)}, nil, &sentinel, defs, drt, gt, pol)
	require.Error(t, err, "a sentinel short of to+grace must fail the run")
	require.Len(t, dres.Verdicts, 1)
	require.Equal(t, OutcomeUnobservable, dres.Verdicts[0].Outcome)
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
	require.True(t, res.Proved, "sentinel exactly at to+grace must satisfy check 1")
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
	require.Equal(t, ReasonFromBeforeRecord, res.Reason)
	require.True(t, res.Unobservable)
	require.False(t, res.Proved)

	// The consequence: decide() must turn this into exit 2, never a pass.
	defs := []Definition{def}
	drt := map[string]ruleTimings{def.UID: rt}
	gt := globalTimings{}
	pol := Policy{From: from, To: to}
	dres, err := decide(Header{StartedAt: started}, nil, &sentinel, defs, drt, gt, pol)
	require.Error(t, err, "`from` before the recording started must fail the run")
	require.Len(t, dres.Verdicts, 1)
	require.Equal(t, OutcomeUnobservable, dres.Verdicts[0].Outcome)
}

// The from-bounds check compares at whole-second granularity: a whole-second
// `from` may precede the recorder's sub-second StartedAt INSIDE the same second
// without being judged early. That one sliver is the --from truncation, not a
// blind interval, so the window is still proved.
func TestProveCoverage_FromSameSecondAsStartedAtIsProved(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(10 * time.Minute)
	started := from.Add(500 * time.Millisecond)
	rt := newRuleTimings(30*time.Second, 60)
	def := Definition{UID: "r1", Title: "R1"}

	var polls []Poll
	for ts := from; !ts.After(to); ts = ts.Add(30 * time.Second) {
		polls = append(polls, Poll{RuleUID: "r1", GrafanaNow: ts, Found: true, Health: "ok", State: "inactive", LastEvaluation: ts})
	}
	sentinel := to

	res := proveCoverage(Header{StartedAt: started}, polls, &sentinel, rt, def, from, to, 0)
	require.True(t, res.Proved)
	require.False(t, res.Unobservable)
	require.Empty(t, res.Reason)
}

// Exactly one whole second later is a different second: even at the boundary,
// the whole-second comparison reads it as before, however healthy the polls.
func TestProveCoverage_FromExactlyOneSecondBeforeStartedAtIsUnobservable(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(10 * time.Minute)
	started := from.Add(time.Second)
	rt := newRuleTimings(30*time.Second, 60)
	def := Definition{UID: "r1", Title: "R1"}

	sentinel := to
	res := proveCoverage(Header{StartedAt: started}, nil, &sentinel, rt, def, from, to, 0)
	require.Equal(t, ReasonFromBeforeRecord, res.Reason)
	require.True(t, res.Unobservable)
	require.False(t, res.Proved)
}

// A sub-second sliver that straddles the second boundary is still "before":
// 900ms into one second vs 100ms into the next are distinct seconds, so the
// 200ms gap is a from_before_record, not rounding noise.
func TestProveCoverage_FromSubSecondEarlierAcrossSecondBoundaryIsUnobservable(t *testing.T) {
	base := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	from := base.Add(900 * time.Millisecond)
	started := base.Add(time.Second + 100*time.Millisecond)
	to := from.Add(10 * time.Minute)
	rt := newRuleTimings(30*time.Second, 60)
	def := Definition{UID: "r1", Title: "R1"}

	sentinel := to
	res := proveCoverage(Header{StartedAt: started}, nil, &sentinel, rt, def, from, to, 0)
	require.Equal(t, ReasonFromBeforeRecord, res.Reason)
	require.True(t, res.Unobservable)
	require.False(t, res.Proved)
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
	require.Equal(t, ReasonHeartbeatGap, res.Reason, "healthy edges with a hole in the middle must still fail")
	// The gap is the SPACING between the two polls (598s), not either
	// boundary segment (1s each) — pin the actual values, not just the verdict.
	require.Equal(t, 598*time.Second, res.LargestGap)
	require.True(t, res.LargestGapAt.Equal(from.Add(time.Second)))
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
	require.True(t, res.Proved, "one failed evaluation must not fail an otherwise clean window")
	require.True(t, anyContains(res.Notes, "health=error"), "want a health=error note even though it did not fail the window")
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
	require.Equal(t, ReasonHealthError, res.Reason, "a run that outlasts healthGrace")
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
	require.True(t, res.Proved, "health=nodata for the WHOLE window must still not be fatal by itself")
	require.True(t, anyContains(res.Notes, "health=nodata"))
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
	require.NotEqual(t, ReasonStaleEvaluation, res.Reason, "liveness must be absolute, never a delta against a previous poll")
	require.Zero(t, res.BlindFor)
	require.True(t, res.Proved)
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
	require.Equal(t, ReasonStaleEvaluation, res.Reason)
	require.Equal(t, 3*time.Minute, res.BlindFor)
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
	require.NotEqual(t, ReasonStaleEvaluation, res.Reason, "a zero lastEvaluation on a paused poll must not trigger check 6")
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
	require.Equal(t, ReasonPausedInWindow, res.Reason)
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
	require.NotEqual(t, ReasonPausedInWindow, res.Reason, "a paused poll after windowEnd tripped check 7")
	require.True(t, res.Proved)
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
	require.Equal(t, ReasonRuleAbsent, res.Reason)
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
	require.True(t, res.Proved, "KeepLast is a note, never fatal")
	require.True(t, anyContains(res.Notes, "KeepLast"), "comma-joined membership, not a literal-key match")
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
			require.True(t, res.Proved, "a declared KeepLast is a note, never fatal")
			require.True(t, anyContains(res.Notes, "KeepLast"),
				"from the definition alone, with zero KeepLast reasons observed")
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
	require.True(t, res.Proved, "a constant clock skew must not itself read as a coverage gap")
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
		require.NoError(t, err)

		var polls []Poll
		for ts := from; !ts.After(windowEnd); ts = ts.Add(120 * time.Second) {
			polls = append(polls, Poll{RuleUID: "r1", GrafanaNow: ts, Found: true, Health: "ok", LastEvaluation: ts})
		}
		sentinel := windowEnd

		res := proveCoverage(h, polls, &sentinel, rt["r1"], defs[0], from, windowEnd, 0)
		require.True(t, res.Proved,
			"maxGap must come from the recorded 120s cadence, not the 30s default")
	})

	t.Run("faster override still catches a real recorder gap", func(t *testing.T) {
		windowEnd := from.Add(20 * time.Minute)
		h := Header{
			StartedAt: from.Add(-time.Hour),
			Rules:     []LoggedRule{{UID: "r1", Title: "R1", IntervalSeconds: 300, PollEverySeconds: 5}},
		}
		defs := []Definition{{UID: "r1", Title: "R1", IntervalSeconds: 300}}
		rt, _, err := DeriveTimingsFromLog(h, defs)
		require.NoError(t, err)

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
		require.Equal(t, ReasonHeartbeatGap, res.Reason,
			"if maxGap had been re-derived from the 300s definition, this 250s gap would pass silently")
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
	require.Equal(t, ReasonStaleEvaluation, res.Reason,
		"a zero lastEvaluation on a found, non-paused poll must fail closed")
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
	require.Equal(t, ReasonHeartbeatGap, res.Reason,
		"the poll's own %s skew bound must push it past the threshold", bound)
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
	require.Equal(t, ReasonPausedInWindow, res.Reason, "the FIRST check to fail names the reason")
	require.True(t, anyContains(res.Notes, "paused"))
	require.True(t, anyContains(res.Notes, "no rule"),
		"a later failure must still be recorded, not swallowed once Reason is already set")
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
	require.Equal(t, ReasonHeartbeatGap, res.Reason,
		"proveCoverage has no 'skipped' concept, so decide must handle a skipped rule's classification itself")
}
