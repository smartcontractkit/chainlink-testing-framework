package gate

import (
	"strings"
	"testing"
	"time"
)

func TestDeriveTimings_Default(t *testing.T) {
	defs := []Definition{
		{UID: "r1", Title: "R1", IntervalSeconds: 60},
	}
	rules, _, notes := DeriveTimings(defs, 0)
	if len(notes) != 0 {
		t.Fatalf("notes = %v, want none", notes)
	}
	rt := rules["r1"]
	if rt.pollEvery != 30*time.Second {
		t.Errorf("pollEvery = %s, want 30s", rt.pollEvery)
	}
	if rt.maxGap != 60*time.Second {
		t.Errorf("maxGap = %s, want 60s", rt.maxGap)
	}
	if rt.healthGrace != 60*time.Second {
		t.Errorf("healthGrace = %s, want 60s", rt.healthGrace)
	}
	if rt.evalStaleAfter != 120*time.Second {
		t.Errorf("evalStaleAfter = %s, want 120s", rt.evalStaleAfter)
	}
}

func TestDeriveTimings_OverrideVerbatimNoClamp(t *testing.T) {
	defs := []Definition{
		{UID: "r1", Title: "R1", IntervalSeconds: 10}, // default pollEvery = 5s
	}
	rules, _, notes := DeriveTimings(defs, 20*time.Second)
	rt := rules["r1"]
	if rt.pollEvery != 20*time.Second {
		t.Fatalf("pollEvery = %s, want the override verbatim (20s), never clamped down to the 5s default", rt.pollEvery)
	}
	if rt.maxGap != 40*time.Second {
		t.Errorf("maxGap = %s, want 2x the override (40s)", rt.maxGap)
	}
	if len(notes) != 1 || !strings.Contains(notes[0], "R1") {
		t.Fatalf("notes = %v, want one note naming R1's exceeded default", notes)
	}
}

func TestDeriveTimings_OverrideBelowDefaultNoNote(t *testing.T) {
	defs := []Definition{{UID: "r1", Title: "R1", IntervalSeconds: 60}} // default pollEvery = 30s
	_, _, notes := DeriveTimings(defs, 5*time.Second)
	if len(notes) != 0 {
		t.Fatalf("notes = %v, want none when the override tightens rather than exceeds the default", notes)
	}
}

func TestDeriveTimings_TransitionGraceExcludesSkippedRule(t *testing.T) {
	defs := []Definition{
		{UID: "r1", Title: "Tight", IntervalSeconds: 60, For: time.Minute},
		{UID: "r2", Title: "PausedLongFor", IntervalSeconds: 60, For: time.Hour, IsPaused: true},
	}
	_, global, _ := DeriveTimings(defs, 0)
	want := time.Minute + 60*time.Second // r1's for+interval; r2 (skipped) must not win despite its huge `for`
	if global.transitionGrace != want {
		t.Fatalf("transitionGrace = %s, want %s (paused rule r2 must be excluded from the max)", global.transitionGrace, want)
	}
	if !strings.Contains(global.graceSource, "Tight") {
		t.Errorf("graceSource = %q, want it to name the contributing rule Tight", global.graceSource)
	}
}

// §22.2/§22.3: `for: 1d` and `for: 1w` parse correctly (parse_ruler_test.go),
// but that alone never proves they flow into transitionGrace — a Prometheus
// duration parser that silently truncated to time.Duration's other units, or
// a transitionGrace derivation that only ever saw hand-built values, could
// each pass every existing test and still be wrong together. This drives the
// real ruler_rules.json fixture (rule0000010, for:1w, DERIVED to exercise the
// w unit — testdata/README.md) through ParseDefinitions and DeriveTimings.
func TestDeriveTimings_RealForOneWeekRuleSetsTransitionGrace(t *testing.T) {
	defs := rulerDefs(t)
	_, global, notes := DeriveTimings(defs, 0)
	if len(notes) != 0 {
		t.Fatalf("notes = %v, want none: no --poll-interval override is given, so no override note should fire", notes)
	}

	want := 7*24*time.Hour + 60*time.Second // rule0000010: for=1w, intervalSeconds=60
	if global.transitionGrace != want {
		t.Fatalf("transitionGrace = %s, want %s (rule0000010's for:1w plus its interval)", global.transitionGrace, want)
	}
	if !strings.Contains(global.graceSource, "Example Failure Ratio Above 10 Percent Weekly") {
		t.Errorf("graceSource = %q, want it to name rule0000010", global.graceSource)
	}
}

func TestDeriveTimings_TransitionGraceZeroWhenAllSkipped(t *testing.T) {
	defs := []Definition{{UID: "r1", Title: "R1", IntervalSeconds: 60, For: time.Hour, IsPaused: true}}
	_, global, _ := DeriveTimings(defs, 0)
	if global.transitionGrace != 0 {
		t.Fatalf("transitionGrace = %s, want 0 when every rule is skipped", global.transitionGrace)
	}
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
		if err != nil {
			t.Fatalf("DeriveTimingsFromLog: %v", err)
		}
		if global.transitionGrace != want {
			t.Fatalf("transitionGrace = %s, want %s: the rule was active when the recording opened, "+
				"so a pause applied afterwards must not shrink the window", global.transitionGrace, want)
		}
		if !strings.Contains(global.graceSource, "R1") {
			t.Errorf("graceSource = %q, want it to name R1", global.graceSource)
		}
	})

	t.Run("header says paused: the rule stays out", func(t *testing.T) {
		h := Header{Rules: []LoggedRule{loggedRule("r1", true)}}
		_, global, err := DeriveTimingsFromLog(h, defs)
		if err != nil {
			t.Fatalf("DeriveTimingsFromLog: %v", err)
		}
		if global.transitionGrace != 0 {
			t.Fatalf("transitionGrace = %s, want 0: a rule paused before the window opened can never fire during it",
				global.transitionGrace)
		}
	})

	t.Run("drainTimeout counts every rule either way", func(t *testing.T) {
		// §19 puts no pause exclusion on drainTimeout, so both headers give the
		// same floor-bound value.
		for _, pausedAtStart := range []bool{false, true} {
			h := Header{Rules: []LoggedRule{loggedRule("r1", pausedAtStart)}}
			_, global, err := DeriveTimingsFromLog(h, defs)
			if err != nil {
				t.Fatalf("DeriveTimingsFromLog: %v", err)
			}
			if global.drainTimeout != minDrainTimeout {
				t.Fatalf("drainTimeout = %s with pausedAtStart=%v, want the %s floor",
					global.drainTimeout, pausedAtStart, minDrainTimeout)
			}
		}
	})
}

func TestDeriveTimings_DrainTimeoutIncludesPaused(t *testing.T) {
	defs := []Definition{{UID: "r1", Title: "R1", IntervalSeconds: 10}}
	_, global, _ := DeriveTimings(defs, 0)
	if global.drainTimeout != minDrainTimeout {
		t.Fatalf("drainTimeout = %s, want the %s floor", global.drainTimeout, minDrainTimeout)
	}
}

func TestDeriveTimings_DrainTimeoutFloor(t *testing.T) {
	defs := []Definition{
		{UID: "r1", Title: "Tight", IntervalSeconds: 60, For: time.Minute},
		{UID: "r2", Title: "PausedLongFor", IntervalSeconds: 180, For: time.Hour, IsPaused: true},
	}
	_, global, _ := DeriveTimings(defs, 0)
	// double the longest interval (2 * 180s) should be the drain timeout
	if global.drainTimeout != 2*180*time.Second {
		t.Fatalf("drainTimeout = %s, want %s", global.drainTimeout, 180*time.Second)
	}
}

func TestDeriveTimings_DrainTimeoutAboveFloor(t *testing.T) {
	defs := []Definition{{UID: "r1", Title: "R1", IntervalSeconds: 300}} // 2x300s = 600s > 2m floor
	_, global, _ := DeriveTimings(defs, 0)
	if global.drainTimeout != 600*time.Second {
		t.Fatalf("drainTimeout = %s, want 600s", global.drainTimeout)
	}
}

// TestScheduler_DueOrderingTiesBreakByTightestCadence pins the ordering
// invariant the burst bound depends on (§5): when several rules become due at
// the exact same instant, Due must serve the tightest cadence first, not
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
	if len(due) != 4 || due[0] != "tight" {
		t.Fatalf("Due = %v, want the tightest-cadence rule (tight) first when all are simultaneously due", due)
	}
}

func TestScheduler_DueExcludesNotYetDue(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	s := &Scheduler{
		next:  map[string]time.Time{"soon": now.Add(-time.Second), "later": now.Add(time.Minute)},
		every: map[string]time.Duration{"soon": 10 * time.Second, "later": 10 * time.Second},
	}
	due := s.Due(now)
	if len(due) != 1 || due[0] != "soon" {
		t.Fatalf("Due = %v, want only [soon]", due)
	}
}

func TestScheduler_MarkAdvancesNextDue(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	s := &Scheduler{
		next:  map[string]time.Time{"r1": now},
		every: map[string]time.Duration{"r1": 30 * time.Second},
	}
	s.Mark("r1", now)
	if got := s.Due(now); len(got) != 0 {
		t.Fatalf("Due right after Mark = %v, want none (next due is 30s out)", got)
	}
	if got := s.Due(now.Add(30 * time.Second)); len(got) != 1 {
		t.Fatalf("Due at next-due time = %v, want [r1]", got)
	}
}

// TestScheduler_PerRuleCadenceOverTime simulates a run and counts how often
// each rule comes due, pinning §5's core claim: schedules are per rule, never
// a global cycle. A tight rule must be polled at its own cadence regardless
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
	if counts["tight"] < 85 || counts["tight"] > 91 {
		t.Errorf("tight polled %d times over 900s, want ~90 (its own 10s cadence)", counts["tight"])
	}
	if counts["slack"] < 2 || counts["slack"] > 4 {
		t.Errorf("slack polled %d times over 900s, want ~3 (its own 300s cadence, not tight's)", counts["slack"])
	}
	if counts["slack"] >= counts["tight"] {
		t.Fatalf("slack polled as often as tight (%d vs %d) — schedules must be per rule, not a shared global cycle", counts["slack"], counts["tight"])
	}
}

func TestNewScheduler_StaggersWithinPollEvery(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	rules := map[string]time.Duration{"r1": 100 * time.Second}
	s := NewScheduler(rules, now)
	offset := s.next["r1"].Sub(now)
	if offset < 0 || offset >= 100*time.Second {
		t.Fatalf("initial offset = %s, want within [0, 100s)", offset)
	}
}

// TestCheckBudget_MixedIntervalRegression is §22.3's sanity check from the
// plan: one rule at 10s beside twenty at 300s, all measured ~1.8s, must not
// error at any reasonable concurrency — the exact case a naive worst-case-slot
// simulation would wrongly fail.
func TestCheckBudget_MixedIntervalRegression(t *testing.T) {
	timings := map[string]ruleTimings{"tight": {pollEvery: 5 * time.Second}}
	measured := map[string]time.Duration{"tight": 1800 * time.Millisecond}
	for i := range 20 {
		uid := uidN(i)
		timings[uid] = ruleTimings{pollEvery: 150 * time.Second}
		measured[uid] = 1800 * time.Millisecond
	}
	if err := CheckBudget(timings, measured, 1); err != nil {
		t.Fatalf("CheckBudget = %v, want nil (utilization 0.6, burst bound 1.8s <= 5s)", err)
	}
}

func TestCheckBudget_UtilizationExceeded(t *testing.T) {
	timings := map[string]ruleTimings{
		"a": {pollEvery: 10 * time.Second},
		"b": {pollEvery: 10 * time.Second},
	}
	measured := map[string]time.Duration{"a": 9 * time.Second, "b": 9 * time.Second}
	err := CheckBudget(timings, measured, 1)
	if err == nil {
		t.Fatal("CheckBudget = nil, want an error: utilization 1.8 > concurrency 1")
	}
	assertBudgetMessage(t, err.Error())
}

func TestCheckBudget_SingleRuleExceedsOwnCadence(t *testing.T) {
	timings := map[string]ruleTimings{"slow": {pollEvery: 5 * time.Second}}
	measured := map[string]time.Duration{"slow": 6 * time.Second}
	err := CheckBudget(timings, measured, 10)
	if err == nil {
		t.Fatal("CheckBudget = nil, want an error: measured 6s exceeds its own 5s poll-interval")
	}
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
	if err == nil {
		t.Fatal("CheckBudget = nil, want a burst-bound error: slow's 3s measured exceeds tight's 2s cadence")
	}
	if !strings.Contains(err.Error(), "burst bound") {
		t.Errorf("error = %q, want it to name the burst bound", err.Error())
	}
	assertBudgetMessage(t, err.Error())
}

func TestCheckBudget_BurstBoundOKWhenNotExceeded(t *testing.T) {
	timings := map[string]ruleTimings{
		"tight": {pollEvery: 5 * time.Second},
		"slow":  {pollEvery: 100 * time.Second},
	}
	measured := map[string]time.Duration{"tight": 100 * time.Millisecond, "slow": 1800 * time.Millisecond}
	if err := CheckBudget(timings, measured, 10); err != nil {
		t.Fatalf("CheckBudget = %v, want nil (1.8s <= 5s tightest cadence)", err)
	}
}

func TestCheckBudget_MissingMeasurementIsAnError(t *testing.T) {
	timings := map[string]ruleTimings{"r1": {pollEvery: 30 * time.Second}}
	err := CheckBudget(timings, map[string]time.Duration{}, 10)
	if err == nil {
		t.Fatal("CheckBudget = nil, want an error: r1 was never measured (fail closed, not a silent zero)")
	}
}

func TestCheckBudget_MissingMixedMeasurementIsAnError(t *testing.T) {
	timings := map[string]ruleTimings{
		"tight": {pollEvery: 5 * time.Second},
		"slow":  {pollEvery: 100 * time.Second},
	}
	measured := map[string]time.Duration{"tight": 100 * time.Millisecond}
	if err := CheckBudget(timings, measured, 10); err == nil {
		t.Fatal("CheckBudget = nil, want an error: slow was never measured (fail closed, not a silent zero)")
	}
}

func TestCheckBudget_EmptyScheduleIsFine(t *testing.T) {
	if err := CheckBudget(nil, nil, 1); err != nil {
		t.Fatalf("CheckBudget = %v, want nil for an empty schedule", err)
	}
}

// assertBudgetMessage checks §5.1's required message contents: a measured
// duration is present, and all three controls are named — never a single
// suggested interval.
func assertBudgetMessage(t *testing.T, msg string) {
	t.Helper()
	for _, want := range []string{"measured", "concurrency", "poll-interval", "fewer"} {
		if !strings.Contains(msg, want) {
			t.Errorf("message %q missing %q", msg, want)
		}
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
	if !strings.Contains(summary, "planned run time") {
		t.Errorf("summary = %q, want it to name the planned run time", summary)
	}
	if warning == "" {
		t.Fatal("warning = \"\", want one: transitionGrace (5m) > 1/4 of the 10m window")
	}
	if !strings.Contains(warning, "R (for=4m30s, interval=30s)") {
		t.Errorf("warning = %q, want it to name the grace source", warning)
	}
}

// §22.3: "a rule with for: 15m in a 10-minute window gives the warning about
// a large grace period" — no such rule exists in the real capture
// (testdata/README.md), so the test above pins the mechanism with a
// hand-built globalTimings. This drives the same warning off the real
// ruler_rules.json fixture's for:1w rule instead, tying ParseDefinitions and
// DeriveTimings into the warning end to end, not just the warning formula in
// isolation.
func TestStartupSummary_RealForOneWeekRuleTriggersWarning(t *testing.T) {
	defs := rulerDefs(t)
	_, global, notes := DeriveTimings(defs, 0)
	if len(notes) != 0 {
		t.Fatalf("notes = %v, want none: no --poll-interval override is given, so no override note should fire", notes)
	}

	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(10 * time.Minute) // transitionGrace (>1w) dwarfs 1/4 of this window
	summary, warning := StartupSummary(from, to, global)
	if !strings.Contains(summary, "planned run time") {
		t.Errorf("summary = %q, want it to name the planned run time", summary)
	}
	if warning == "" {
		t.Fatal("warning = \"\", want one: a real for:1w rule's transitionGrace vastly exceeds 1/4 of a 10m window")
	}
	if !strings.Contains(warning, "Example Failure Ratio Above 10 Percent Weekly") {
		t.Errorf("warning = %q, want it to name rule0000010", warning)
	}
}

func TestStartupSummary_NoWarningWhenGraceSmall(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(time.Hour)
	global := globalTimings{transitionGrace: time.Minute, graceSource: "R (for=30s, interval=30s)", drainTimeout: time.Minute}
	_, warning := StartupSummary(from, to, global)
	if warning != "" {
		t.Fatalf("warning = %q, want none: 1m grace is well under 1/4 of a 1h window", warning)
	}
}

func TestStartupSummary_NoGraceSourceReadsNone(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(time.Hour)
	summary, _ := StartupSummary(from, to, globalTimings{})
	if !strings.Contains(summary, "none") {
		t.Fatalf("summary = %q, want it to read \"none\" when no rule set the grace", summary)
	}
}
