package gate

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// One tight rule beside twenty slack ones stays safe even though the pass
// (37.8s) dwarfs the tight rule's 10s maxGap: the slack rules are not due yet.
func TestCheckStartupHandoff_MixedIntervalFleetFits(t *testing.T) {
	base := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	const latency = 1800 * time.Millisecond
	timings := map[string]RuleTimings{"tight": {pollEvery: 5 * time.Second, maxGap: 10 * time.Second}}
	measured := map[string]time.Duration{"tight": latency}
	first := []Poll{{RuleUID: "tight", GrafanaNow: base, Found: true}}
	for i := range 20 {
		uid := uidN(i)
		timings[uid] = RuleTimings{pollEvery: 150 * time.Second, maxGap: 300 * time.Second}
		measured[uid] = latency
		first = append(first, Poll{RuleUID: uid, GrafanaNow: base.Add(time.Duration(i+1) * latency), Found: true})
	}
	readyAt := base.Add(21 * latency)

	require.NoError(t, CheckStartupHandoff(timings, measured, first, readyAt, readyAt, 1))
}

// A tight rule observed late queues behind the rules due before it: 50 rules at
// 200ms make a 10s pass, so the last one's first poll is ~10s out against 8s.
func TestCheckStartupHandoff_RefusesATightRuleBehindALongBacklog(t *testing.T) {
	base := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	const (
		n       = 50
		latency = 200 * time.Millisecond
	)
	timings := make(map[string]RuleTimings, n)
	measured := make(map[string]time.Duration, n)
	first := make([]Poll, 0, n)
	for i := range n {
		uid := fmt.Sprintf("r%03d", i)
		timings[uid] = RuleTimings{title: fmt.Sprintf("Rule %03d", i), pollEvery: 4 * time.Second, maxGap: 8 * time.Second}
		measured[uid] = latency
		first = append(first, Poll{RuleUID: uid, GrafanaNow: base.Add(time.Duration(i) * latency), Found: true})
	}
	readyAt := base.Add(n * latency)

	err := CheckStartupHandoff(timings, measured, first, readyAt, readyAt, 1)
	require.Error(t, err, "a 10s backlog cannot fit an 8s maxGap")
	assertBudgetMessage(t, err.Error())
	require.Contains(t, err.Error(), "raising --concurrency to at least 2",
		"the error must name the concurrency that would fit, not just the lever")
	require.Regexp(t, `rule "Rule \d+" \(r\d+\)`, err.Error(),
		"the offending rules must be named by title, not only by UID")

	require.NoError(t, CheckStartupHandoff(timings, measured, first, readyAt, readyAt, 2),
		"doubling concurrency halves the backlog drain")
}

func TestCheckStartupHandoff_MissingInputsFailClosed(t *testing.T) {
	timings := map[string]RuleTimings{"r1": {pollEvery: 30 * time.Second, maxGap: time.Minute}}
	measured := map[string]time.Duration{"r1": time.Second}

	t.Run("no first observation", func(t *testing.T) {
		err := CheckStartupHandoff(timings, measured, nil, testNow, testNow, 1)
		require.Error(t, err, "a rule never observed cannot be proved")
		require.Contains(t, err.Error(), "first observation")
	})

	t.Run("no measurement", func(t *testing.T) {
		first := []Poll{{RuleUID: "r1", GrafanaNow: testNow, Found: true}}
		err := CheckStartupHandoff(timings, nil, first, testNow, testNow, 1)
		require.Error(t, err, "a rule never measured cannot have its backlog bounded")
		require.Contains(t, err.Error(), "never measured")
	})

	t.Run("empty schedule", func(t *testing.T) {
		require.NoError(t, CheckStartupHandoff(nil, nil, nil, testNow, testNow, 1))
	})
}

// Release delay and queueing ADD, they do not `max`: a rule due at +5s behind a
// batch released at +4.9s finishes at +9.9s, not 5.5s.
func TestCheckStartupHandoff_ReleaseDelayAndQueueAdd(t *testing.T) {
	base := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	timings := make(map[string]RuleTimings, 10)
	measured := make(map[string]time.Duration, 10)
	first := make([]Poll, 0, 10)

	for i := range 9 {
		uid := fmt.Sprintf("backlog%d", i)
		timings[uid] = RuleTimings{pollEvery: 4900 * time.Millisecond, maxGap: 9800 * time.Millisecond}
		measured[uid] = 500 * time.Millisecond
		first = append(first, Poll{RuleUID: uid, GrafanaNow: base, Found: true}) // due +4.9s
	}
	timings["tight"] = RuleTimings{pollEvery: 2750 * time.Millisecond, maxGap: 5500 * time.Millisecond}
	measured["tight"] = 500 * time.Millisecond
	first = append(first, Poll{RuleUID: "tight", GrafanaNow: base.Add(2250 * time.Millisecond), Found: true}) // due +5s

	err := CheckStartupHandoff(timings, measured, first, base, base, 1)
	require.Error(t, err, "the tight rule is polled after the 4.5s backlog it is queued behind")
	require.Contains(t, err.Error(), "tight")
}

// The recorder's whole-second `from` tolerance can open the window 999ms early,
// so the gap must be budgeted from that earlier instant.
func TestCheckStartupHandoff_BudgetsTheWholeSecondFromTolerance(t *testing.T) {
	base := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	const latency = 1200 * time.Millisecond
	timings := map[string]RuleTimings{
		"a":     {pollEvery: 2 * time.Second, maxGap: 4 * time.Second},
		"b":     {pollEvery: 2 * time.Second, maxGap: 4 * time.Second},
		"tight": {pollEvery: 2 * time.Second, maxGap: 4 * time.Second},
	}
	measured := map[string]time.Duration{"a": latency, "b": latency, "tight": latency}
	obs := base.Add(-1500 * time.Millisecond) // due base+0.5s, inside the first batch
	first := []Poll{
		{RuleUID: "a", GrafanaNow: obs, Found: true},
		{RuleUID: "b", GrafanaNow: obs, Found: true},
		{RuleUID: "tight", GrafanaNow: obs, Found: true},
	}
	readyAt := base.Add(900 * time.Millisecond) // truncates to base

	// Single-step: the clamp is exact, and the 3.6s gap fits the 4s maxGap.
	require.NoError(t, CheckStartupHandoff(timings, measured, first, readyAt, readyAt, 1))
	// Recorder mode: the same run can open at base, making the gap 4.5s.
	require.Error(t, CheckStartupHandoff(timings, measured, first, readyAt, readyAt.Truncate(time.Second), 1))
}

// A handoff that fails even at one worker per rule: no concurrency can fix it,
// so the message must not suggest raising it.
func TestCheckStartupHandoff_NoConcurrencyCanFixIt(t *testing.T) {
	base := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	timings := map[string]RuleTimings{"r": {pollEvery: time.Second, maxGap: 2 * time.Second}}
	measured := map[string]time.Duration{"r": 100 * time.Millisecond}
	first := []Poll{{RuleUID: "r", GrafanaNow: base, Found: true}}

	err := CheckStartupHandoff(timings, measured, first, base, base.Add(-time.Second), 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "fix by: raising poll-interval")
	require.NotContains(t, err.Error(), "raising concurrency")
}
