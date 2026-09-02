package gate

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDeriveTimings_Default(t *testing.T) {
	defs := []Definition{
		{UID: "r1", Title: "R1", IntervalSeconds: 60},
	}
	rules, _, notes := DeriveTimings(defs, 0)
	require.Empty(t, notes)
	rt := rules["r1"]
	require.Equal(t, 30*time.Second, rt.pollEvery)
	require.Equal(t, 60*time.Second, rt.maxGap)
	require.Equal(t, 60*time.Second, rt.healthGrace)
	require.Equal(t, 120*time.Second, rt.evalStaleAfter)
}

func TestDeriveTimings_OverrideVerbatimNoClamp(t *testing.T) {
	defs := []Definition{
		{UID: "r1", Title: "R1", IntervalSeconds: 10}, // default pollEvery = 5s
	}
	rules, _, notes := DeriveTimings(defs, 20*time.Second)
	rt := rules["r1"]
	require.Equal(t, 20*time.Second, rt.pollEvery, "the override verbatim (20s), never clamped down to the 5s default")
	require.Equal(t, 40*time.Second, rt.maxGap)
	require.Len(t, notes, 1)
	require.Contains(t, notes[0], "R1")
}

func TestDeriveTimings_OverrideBelowDefaultNoNote(t *testing.T) {
	defs := []Definition{{UID: "r1", Title: "R1", IntervalSeconds: 60}} // default pollEvery = 30s
	_, _, notes := DeriveTimings(defs, 5*time.Second)
	require.Empty(t, notes)
}

func TestDeriveTimings_TransitionGraceExcludesSkippedRule(t *testing.T) {
	defs := []Definition{
		{UID: "r1", Title: "Tight", IntervalSeconds: 60, For: time.Minute},
		{UID: "r2", Title: "PausedLongFor", IntervalSeconds: 60, For: time.Hour, IsPaused: true},
	}
	_, global, _ := DeriveTimings(defs, 0)
	want := time.Minute + 60*time.Second // r1's for+interval; r2 (skipped) must not win despite its huge `for`
	require.Equal(t, want, global.transitionGrace)
	require.Contains(t, global.graceSource, "Tight")
}

// `for: 1d` and `for: 1w` parse correctly (parse_ruler_test.go),
// but that alone never proves they flow into transitionGrace — a Prometheus
// duration parser that silently truncated to time.Duration's other units, or
// a transitionGrace derivation that only ever saw hand-built values, could
// each pass every existing test and still be wrong together. This drives the
// real ruler_rules.json fixture (rule0000010, for:1w, DERIVED to exercise the
// w unit — testdata/README.md) through ParseDefinitions and DeriveTimings.
func TestDeriveTimings_RealForOneWeekRuleSetsTransitionGrace(t *testing.T) {
	defs := rulerDefs(t)
	_, global, notes := DeriveTimings(defs, 0)
	require.Empty(t, notes, "no --poll-interval override is given, so no override note should fire")

	want := 7*24*time.Hour + 60*time.Second // rule0000010: for=1w, intervalSeconds=60
	require.Equal(t, want, global.transitionGrace)
	require.Contains(t, global.graceSource, "Example Failure Ratio Above 10 Percent Weekly")
}

func TestDeriveTimings_TransitionGraceZeroWhenAllSkipped(t *testing.T) {
	defs := []Definition{{UID: "r1", Title: "R1", IntervalSeconds: 60, For: time.Hour, IsPaused: true}}
	_, global, _ := DeriveTimings(defs, 0)
	require.Zero(t, global.transitionGrace)
}

// TestDeriveTimingsFromLog_TransitionGraceFollowsTheHeaderNotTheDefinition
// pins the log-mode authority for the grace exclusion. Definitions are
// re-resolved AFTER the window closed, so "paused" in a definition says
// nothing about whether the rule was watched during it.
//
// The fail-open direction is the first case. transitionGrace exists so a
// condition arising just before `to` is still seen when it surfaces at
// to + `for`, and windowEnd is both the classification bound and the
// collection deadline — so a rule somebody paused after `to` dropping out of
// the max collapses the grace, the surfacing poll is never recorded, and the
// run reports clean.
func TestDeriveTimingsFromLog_TransitionGraceFollowsTheHeaderNotTheDefinition(t *testing.T) {
	loggedRule := func(uid string, pausedAtStart bool) LoggedRule {
		return LoggedRule{UID: uid, Title: uid, IntervalSeconds: 60, PollEverySeconds: 30, IsPaused: pausedAtStart}
	}
	// The definition says paused in BOTH cases: it is the post-window reading,
	// and it must change nothing.
	defs := []Definition{{UID: "r1", Title: "R1", IntervalSeconds: 60, For: 5 * time.Minute, IsPaused: true}}
	want := 5*time.Minute + 60*time.Second

	t.Run("header says active: the rule stays in the max", func(t *testing.T) {
		h := Header{Rules: []LoggedRule{loggedRule("r1", false)}}
		_, global, err := DeriveTimingsFromLog(h, defs)
		require.NoError(t, err)
		require.Equal(t, want, global.transitionGrace,
			"the rule was active when the recording opened, so a pause applied afterwards must not shrink the window")
		require.Contains(t, global.graceSource, "R1")
	})

	t.Run("header says paused: the rule stays out", func(t *testing.T) {
		h := Header{Rules: []LoggedRule{loggedRule("r1", true)}}
		_, global, err := DeriveTimingsFromLog(h, defs)
		require.NoError(t, err)
		require.Zero(t, global.transitionGrace,
			"a rule paused before the window opened can never fire during it")
	})

	t.Run("drainTimeout counts every rule either way", func(t *testing.T) {
		// drainTimeout carries no pause exclusion, so both headers give the
		// same floor-bound value.
		for _, pausedAtStart := range []bool{false, true} {
			h := Header{Rules: []LoggedRule{loggedRule("r1", pausedAtStart)}}
			_, global, err := DeriveTimingsFromLog(h, defs)
			require.NoError(t, err)
			require.Equalf(t, minDrainTimeout, global.drainTimeout, "pausedAtStart=%v", pausedAtStart)
		}
	})
}

func TestDeriveTimings_DrainTimeoutIncludesPaused(t *testing.T) {
	defs := []Definition{{UID: "r1", Title: "R1", IntervalSeconds: 10}}
	_, global, _ := DeriveTimings(defs, 0)
	require.Equal(t, minDrainTimeout, global.drainTimeout)
}

func TestDeriveTimings_DrainTimeoutFloor(t *testing.T) {
	defs := []Definition{
		{UID: "r1", Title: "Tight", IntervalSeconds: 60, For: time.Minute},
		{UID: "r2", Title: "PausedLongFor", IntervalSeconds: 180, For: time.Hour, IsPaused: true},
	}
	_, global, _ := DeriveTimings(defs, 0)
	// double the longest interval (2 * 180s) should be the drain timeout
	require.Equal(t, 2*180*time.Second, global.drainTimeout)
}

func TestDeriveTimings_DrainTimeoutAboveFloor(t *testing.T) {
	defs := []Definition{{UID: "r1", Title: "R1", IntervalSeconds: 300}} // 2x300s = 600s > 2m floor
	_, global, _ := DeriveTimings(defs, 0)
	require.Equal(t, 600*time.Second, global.drainTimeout)
}

// The ordering invariant the burst bound depends on: when several rules become
// due at the exact same instant, Due must serve the tightest cadence first, not
// whatever order the underlying map happens to iterate in. A refactor that
// loses this ordering must fail here, not in a production phase-aligned gap.
func TestScheduler_DueOrderingTiesBreakByTightestCadence(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	s := &Scheduler{
		next: map[string]time.Time{
			"slack1": now, "slack2": now, "tight": now, "slack3": now,
		},
		every: map[string]time.Duration{
			"slack1": 300 * time.Second,
			"slack2": 300 * time.Second,
			"tight":  10 * time.Second,
			"slack3": 300 * time.Second,
		},
	}
	due := s.Due(now)
	require.Len(t, due, 4)
	require.Equal(t, "tight", due[0], "the tightest-cadence rule must be first when all are simultaneously due")
}

func TestScheduler_DueExcludesNotYetDue(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	s := &Scheduler{
		next:  map[string]time.Time{"soon": now.Add(-time.Second), "later": now.Add(time.Minute)},
		every: map[string]time.Duration{"soon": 10 * time.Second, "later": 10 * time.Second},
	}
	due := s.Due(now)
	require.Len(t, due, 1)
	require.Equal(t, "soon", due[0])
}

func TestScheduler_MarkAdvancesNextDue(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	s := &Scheduler{
		next:  map[string]time.Time{"r1": now},
		every: map[string]time.Duration{"r1": 30 * time.Second},
	}
	s.Mark("r1", now)
	require.Empty(t, s.Due(now), "next due is 30s out")
	require.Len(t, s.Due(now.Add(30*time.Second)), 1)
}

// TestScheduler_PerRuleCadenceOverTime simulates a run and counts how often
// each rule comes due: schedules are per rule, never a global cycle. A tight
// rule must be polled at its own cadence regardless
// of what slower rules in the same fleet need, and a slack rule must never be
// forced onto the tight rule's cadence.
func TestScheduler_PerRuleCadenceOverTime(t *testing.T) {
	start := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	rules := map[string]time.Duration{
		"tight": 10 * time.Second,
		"slack": 300 * time.Second,
	}
	s := NewScheduler(rules, start)

	const runFor = 900 * time.Second
	const step = time.Second
	counts := map[string]int{}
	for elapsed := time.Duration(0); elapsed <= runFor; elapsed += step {
		now := start.Add(elapsed)
		for _, uid := range s.Due(now) {
			counts[uid]++
			s.Mark(uid, now)
		}
	}

	// 900s of runtime: "tight" (10s cadence) polls ~90 times, "slack" (300s
	// cadence) ~3 times. Assert the ratio holds rather than an exact count,
	// since the staggered initial offset shifts each by up to one cadence.
	require.GreaterOrEqual(t, counts["tight"], 85)
	require.LessOrEqual(t, counts["tight"], 91)
	require.GreaterOrEqual(t, counts["slack"], 2)
	require.LessOrEqual(t, counts["slack"], 4)
	require.Less(t, counts["slack"], counts["tight"],
		"schedules must be per rule, not a shared global cycle")
}

func TestNewScheduler_StaggersWithinPollEvery(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	rules := map[string]time.Duration{"r1": 100 * time.Second}
	s := NewScheduler(rules, now)
	offset := s.next["r1"].Sub(now)
	require.GreaterOrEqual(t, offset, time.Duration(0))
	require.Less(t, offset, 100*time.Second)
}

// One rule at 10s beside twenty at 300s, all measured ~1.8s, must not error at
// any reasonable concurrency — the exact case a naive worst-case-slot
// simulation would wrongly fail.
func TestCheckBudget_MixedIntervalRegression(t *testing.T) {
	timings := map[string]ruleTimings{"tight": {pollEvery: 5 * time.Second}}
	measured := map[string]time.Duration{"tight": 1800 * time.Millisecond}
	for i := range 20 {
		uid := uidN(i)
		timings[uid] = ruleTimings{pollEvery: 150 * time.Second}
		measured[uid] = 1800 * time.Millisecond
	}
	require.NoError(t, CheckBudget(timings, measured, 1))
}

func TestCheckBudget_UtilizationExceeded(t *testing.T) {
	timings := map[string]ruleTimings{
		"a": {pollEvery: 10 * time.Second},
		"b": {pollEvery: 10 * time.Second},
	}
	measured := map[string]time.Duration{"a": 9 * time.Second, "b": 9 * time.Second}
	err := CheckBudget(timings, measured, 1)
	require.Error(t, err)
	assertBudgetMessage(t, err.Error())
}

func TestCheckBudget_SingleRuleExceedsOwnCadence(t *testing.T) {
	timings := map[string]ruleTimings{"slow": {pollEvery: 5 * time.Second}}
	measured := map[string]time.Duration{"slow": 6 * time.Second}
	err := CheckBudget(timings, measured, 10)
	require.Error(t, err, "measured 6s exceeds its own 5s poll-interval")
	assertBudgetMessage(t, err.Error())
}

func TestCheckBudget_BurstBoundViolation(t *testing.T) {
	// Utilization is trivially fine, but the slower rule's request time (3s)
	// exceeds the tighter rule's cadence (2s) — a mid-run gap risk even
	// though no single rule breaches its own cadence and utilization is low.
	timings := map[string]ruleTimings{
		"tight": {pollEvery: 2 * time.Second},
		"slow":  {pollEvery: 100 * time.Second},
	}
	measured := map[string]time.Duration{"tight": 100 * time.Millisecond, "slow": 3 * time.Second}
	err := CheckBudget(timings, measured, 10)
	require.Error(t, err, "slow's 3s measured exceeds tight's 2s cadence")
	require.Contains(t, err.Error(), "burst bound")
	assertBudgetMessage(t, err.Error())
}

func TestCheckBudget_BurstBoundOKWhenNotExceeded(t *testing.T) {
	timings := map[string]ruleTimings{
		"tight": {pollEvery: 5 * time.Second},
		"slow":  {pollEvery: 100 * time.Second},
	}
	measured := map[string]time.Duration{"tight": 100 * time.Millisecond, "slow": 1800 * time.Millisecond}
	require.NoError(t, CheckBudget(timings, measured, 10))
}

func TestCheckBudget_MissingMeasurementIsAnError(t *testing.T) {
	timings := map[string]ruleTimings{"r1": {pollEvery: 30 * time.Second}}
	err := CheckBudget(timings, map[string]time.Duration{}, 10)
	require.Error(t, err, "r1 was never measured (fail closed, not a silent zero)")
}

func TestCheckBudget_MissingMixedMeasurementIsAnError(t *testing.T) {
	timings := map[string]ruleTimings{
		"tight": {pollEvery: 5 * time.Second},
		"slow":  {pollEvery: 100 * time.Second},
	}
	measured := map[string]time.Duration{"tight": 100 * time.Millisecond}
	require.Error(t, CheckBudget(timings, measured, 10),
		"slow was never measured (fail closed, not a silent zero)")
}

func TestCheckBudget_EmptyScheduleIsFine(t *testing.T) {
	require.NoError(t, CheckBudget(nil, nil, 1))
}

// assertBudgetMessage checks the message contents: a measured duration is
// present, and all three controls are named — never a single suggested
// interval.
func assertBudgetMessage(t *testing.T, msg string) {
	t.Helper()
	for _, want := range []string{"measured", "concurrency", "poll-interval", "fewer"} {
		require.Contains(t, msg, want)
	}
}

func uidN(i int) string {
	return "slack" + string(rune('a'+i))
}

func TestStartupSummary_WarningWhenGraceTooLarge(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(10 * time.Minute)
	global := globalTimings{transitionGrace: 5 * time.Minute, graceSource: "R (for=4m30s, interval=30s)", drainTimeout: time.Minute}
	summary, warning := StartupSummary(from, to, global)
	require.Contains(t, summary, "planned run time")
	require.NotEmpty(t, warning, "transitionGrace (5m) > 1/4 of the 10m window")
	require.Contains(t, warning, "R (for=4m30s, interval=30s)")
}

// The test above pins the warning formula with a hand-built globalTimings.
// This drives the same warning off the real ruler_rules.json fixture's for:1w
// rule instead, tying ParseDefinitions and DeriveTimings into the warning end
// to end.
func TestStartupSummary_RealForOneWeekRuleTriggersWarning(t *testing.T) {
	defs := rulerDefs(t)
	_, global, notes := DeriveTimings(defs, 0)
	require.Empty(t, notes)

	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(10 * time.Minute) // transitionGrace (>1w) dwarfs 1/4 of this window
	summary, warning := StartupSummary(from, to, global)
	require.Contains(t, summary, "planned run time")
	require.NotEmpty(t, warning)
	require.Contains(t, warning, "Example Failure Ratio Above 10 Percent Weekly")
}

func TestStartupSummary_NoWarningWhenGraceSmall(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(time.Hour)
	global := globalTimings{transitionGrace: time.Minute, graceSource: "R (for=30s, interval=30s)", drainTimeout: time.Minute}
	_, warning := StartupSummary(from, to, global)
	require.Empty(t, warning)
}

func TestStartupSummary_NoGraceSourceReadsNone(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(time.Hour)
	summary, _ := StartupSummary(from, to, globalTimings{})
	require.Contains(t, summary, "none")
}
