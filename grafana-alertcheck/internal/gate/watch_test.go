package gate

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// The two fixture rules every prepareWatch test below uses: one live, one
// paused in its definition. Both are addressed by uid:, because the ruler
// fixture deliberately contains a 2-way title collision and a title would make
// the tests depend on which side of it they hit.
const (
	watchActiveUID   = "rule0000009"
	watchActiveTitle = "Example Failure Ratio Above 10 Percent"
	watchPausedUID   = "rule0000007"
	watchPausedTitle = "example_workflow_paused_rule"
)

// loopSource answers every RuleState call from a responder that also sees the
// call count, so a recorder-loop test can make the answer depend on virtual
// time or fail on the Nth poll. The loop never reads Version or Definitions —
// the parent did that before detaching — so both fail loudly here.
type loopSource struct {
	mu      sync.Mutex
	calls   map[string]int
	respond func(title string, call int) (Observation, error)
}

func newLoopSource(respond func(title string, call int) (Observation, error)) *loopSource {
	return &loopSource{calls: map[string]int{}, respond: respond}
}

func (s *loopSource) Version(context.Context) (string, error) {
	return "", errors.New("loopSource: the recorder loop must not read the version")
}

func (s *loopSource) Definitions(context.Context) ([]Definition, error) {
	return nil, errors.New("loopSource: the recorder loop must not read the definitions")
}

func (s *loopSource) RuleState(_ context.Context, title string) (Observation, error) {
	s.mu.Lock()
	s.calls[title]++
	call := s.calls[title]
	s.mu.Unlock()
	return s.respond(title, call)
}

var _ Source = (*loopSource)(nil)

// testStateRule is one rule as the state endpoint would return it, healthy and
// evaluated at grafanaNow.
func testStateRule(uid, title string, interval time.Duration, grafanaNow time.Time, instances ...Instance) StateRule {
	totals := map[string]int{"normal": len(instances)}
	return StateRule{
		UID: uid, Title: title, Folder: "F", Group: "G",
		Interval: interval, State: "inactive", Health: "ok",
		LastEvaluation: grafanaNow, Totals: totals, Instances: instances,
	}
}

// newLoopWriter opens a log with a header already written, exactly as the
// parent hands it to the child.
func newLoopWriter(t *testing.T, path string, clock Clock) *Writer {
	t.Helper()
	w, err := NewWriter(path, clock)
	require.NoError(t, err)
	require.NoError(t, w.WriteHeader(testHeader()))
	return w
}

func countPolls(polls []Poll, uid string) int {
	n := 0
	for _, p := range polls {
		if p.RuleUID == uid {
			n++
		}
	}
	return n
}

// The per-rule schedule seen from the recorder: a 10s rule beside a 300s one
// keeps its own 5s cadence instead of dragging the slack rule along with it or
// being slowed to its pace.
func TestWatchLoopPollsEachRuleAtItsOwnCadence(t *testing.T) {
	const tightUID, slackUID = "tight", "slack"
	path := filepath.Join(t.TempDir(), "log.jsonl")
	clock := newVirtualClock(testNow)
	w := newLoopWriter(t, path, clock)

	src := newLoopSource(func(title string, _ int) (Observation, error) {
		uid := tightUID
		if title == "Slack Rule" {
			uid = slackUID
		}
		now := clock.Now()
		return observation(now, testStateRule(uid, title, time.Minute, now)), nil
	})

	err := watchLoop(context.Background(), watchLoopConfig{
		Src:     src,
		Writer:  w,
		Reducer: NewReducer(),
		Titles:  map[string]string{tightUID: "Tight Rule", slackUID: "Slack Rule"},
		Cadence: map[string]time.Duration{
			tightUID: 5 * time.Second,
			slackUID: 150 * time.Second,
		},
		Until:       testNow.Add(300 * time.Second),
		Concurrency: 2,
		Clock:       clock,
	})
	require.NoError(t, err)

	_, polls, sentinel, readErr := ReadLog(path)
	require.NoError(t, readErr)
	require.NotNil(t, sentinel, "no stopped sentinel after a clean stop")
	require.False(t, sentinel.Before(testNow.Add(300*time.Second)))
	// 300s of window at 5s and 150s, minus the initial stagger offset of up to
	// one cadence: 59-60 and 1-2. The assertion is the ratio, not the exact
	// count — a single global cycle would give both rules the same number.
	got := countPolls(polls, tightUID)
	require.GreaterOrEqual(t, got, 59)
	require.LessOrEqual(t, got, 61)
	got = countPolls(polls, slackUID)
	require.GreaterOrEqual(t, got, 1)
	require.LessOrEqual(t, got, 3)
}

// Fail-closed from the recorder's side: a recorder that dies must look exactly
// like a coverage gap, so it must not sign off the log on its way out.
func TestWatchLoopHardErrorLeavesNoSentinel(t *testing.T) {
	path := filepath.Join(t.TempDir(), "log.jsonl")
	clock := newVirtualClock(testNow)
	w := newLoopWriter(t, path, clock)

	boom := errors.New("grafana went away for good")
	src := newLoopSource(func(title string, call int) (Observation, error) {
		if call >= 2 {
			return Observation{}, boom
		}
		now := clock.Now()
		return observation(now, testStateRule("r1", title, time.Minute, now)), nil
	})

	err := watchLoop(context.Background(), watchLoopConfig{
		Src:         src,
		Writer:      w,
		Reducer:     NewReducer(),
		Titles:      map[string]string{"r1": "Example"},
		Cadence:     map[string]time.Duration{"r1": 30 * time.Second},
		Until:       testNow.Add(time.Hour),
		Concurrency: 1,
		Clock:       clock,
	})
	require.ErrorIs(t, err, boom)

	_, polls, sentinel, readErr := ReadLog(path)
	require.NoError(t, readErr)
	require.Nil(t, sentinel, "check would read that as a finished window")
	require.Len(t, polls, 1, "want the 1 that succeeded before the failure")
}

// SIGTERM arriving while a poll is in flight is a clean stop, so the aborted
// poll's error must not suppress the sentinel — otherwise every normal check
// run, which stops the recorder exactly this way, would end unobservable.
func TestWatchLoopSignalDuringPollIsACleanStop(t *testing.T) {
	path := filepath.Join(t.TempDir(), "log.jsonl")
	clock := newVirtualClock(testNow)
	w := newLoopWriter(t, path, clock)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	src := newLoopSource(func(title string, call int) (Observation, error) {
		if call >= 2 {
			// The signal lands while this request is out.
			cancel()
			return Observation{}, ctx.Err()
		}
		now := clock.Now()
		return observation(now, testStateRule("r1", title, time.Minute, now)), nil
	})

	require.NoError(t, watchLoop(ctx, watchLoopConfig{
		Src:         src,
		Writer:      w,
		Reducer:     NewReducer(),
		Titles:      map[string]string{"r1": "Example"},
		Cadence:     map[string]time.Duration{"r1": 30 * time.Second},
		Concurrency: 1,
		Clock:       clock,
	}))

	_, _, sentinel, err := ReadLog(path)
	require.NoError(t, err)
	require.NotNil(t, sentinel, "no sentinel after a signalled stop; check would call a fully observed window unobservable")
}

// TestWatchLoopWithNothingToPollStillFinishesTheLog covers the every-rule-is-
// paused case: there is nothing to record, but "the recorder ran and finished"
// is still what check has to prove about the window.
func TestWatchLoopWithNothingToPollStillFinishesTheLog(t *testing.T) {
	path := filepath.Join(t.TempDir(), "log.jsonl")
	clock := newVirtualClock(testNow)
	w := newLoopWriter(t, path, clock)

	src := newLoopSource(func(title string, _ int) (Observation, error) {
		return Observation{}, fmt.Errorf("nothing should be polled, got %q", title)
	})

	require.NoError(t, watchLoop(context.Background(), watchLoopConfig{
		Src:         src,
		Writer:      w,
		Reducer:     NewReducer(),
		Titles:      map[string]string{},
		Cadence:     map[string]time.Duration{},
		Until:       testNow.Add(time.Minute),
		Concurrency: 1,
		Clock:       clock,
	}))

	_, polls, sentinel, err := ReadLog(path)
	require.NoError(t, err)
	require.Empty(t, polls)
	require.NotNil(t, sentinel, "no sentinel: check cannot tell this recording from one that died")
}

// TestWatchLoopPollBatchKeepsTheHeartbeatsItGot: one rule's failure must not
// discard another rule's observed heartbeat, or a single transport failure
// turns into a coverage gap for every rule that answered.
func TestWatchLoopPollBatchKeepsTheHeartbeatsItGot(t *testing.T) {
	path := filepath.Join(t.TempDir(), "log.jsonl")
	clock := newVirtualClock(testNow)
	w := newLoopWriter(t, path, clock)

	boom := errors.New("one rule is unreachable")
	src := newLoopSource(func(title string, _ int) (Observation, error) {
		if title == "Broken" {
			return Observation{}, boom
		}
		now := clock.Now()
		return observation(now, testStateRule("ok", title, time.Minute, now)), nil
	})

	cfg := watchLoopConfig{
		Src:         src,
		Writer:      w,
		Reducer:     NewReducer(),
		Titles:      map[string]string{"ok": "Healthy", "bad": "Broken"},
		Cadence:     map[string]time.Duration{"ok": 30 * time.Second, "bad": 30 * time.Second},
		Concurrency: 2,
		Clock:       clock,
	}
	require.ErrorIs(t, cfg.pollBatch(context.Background(), []string{"ok", "bad"}), boom)
	require.NoError(t, w.Close())

	_, polls, _, err := ReadLog(path)
	require.NoError(t, err)
	require.Len(t, polls, 1)
	require.Equal(t, "ok", polls[0].RuleUID)
}

// The vanish-versus-clear distinction at the one seam the parent/child handoff
// introduces. The parent observes a firing instance; the child starts with a
// fresh Reducer and sees the instance gone. Seeded, that is a vanish — a
// discontinuity. Unseeded, it is nothing at all, and the instance silently
// leaves the record as if it had never been bad.
func TestReducerSeedFromKeepsMarkersAcrossTheHandoff(t *testing.T) {
	firing := testInstance(StateFiring, "", "b")
	key := instanceKey(firing.Labels)
	parentPoll := Poll{RuleUID: "r1", Found: true, Abnormal: []Instance{firing}}
	// The child's first response: the instance is gone from the response
	// entirely, which is a vanish and never a clear.
	childObs := observation(testNow, testStateRule("r1", "Example", time.Minute, testNow))

	t.Run("seeded", func(t *testing.T) {
		r := NewReducer()
		r.seedFrom([]Poll{parentPoll})
		p := r.Reduce("r1", childObs)
		require.Contains(t, p.Vanished, key)
		require.Empty(t, p.Cleared, "a vanish is not a recovery")
	})

	t.Run("unseeded loses the transition", func(t *testing.T) {
		p := NewReducer().Reduce("r1", childObs)
		require.Empty(t, p.Vanished, "this subtest exists to show the seed is what produces the marker")
	})

	t.Run("a not-found poll does not clear the seed", func(t *testing.T) {
		r := NewReducer()
		r.seedFrom([]Poll{parentPoll, {RuleUID: "r1", Found: false}})
		p := r.Reduce("r1", childObs)
		require.Contains(t, p.Vanished, key, "an absent rule leaves the abnormal set untouched")
	})
}

// watchTestConfig is a prepareWatch config over a temp log, with the notes
// captured so the tests can assert on what an operator is told.
func watchTestConfig(t *testing.T, notes *strings.Builder, alerts ...string) WatchConfig {
	t.Helper()
	return WatchConfig{
		URL:         "https://grafana.example.com",
		Token:       "secret-token",
		Alerts:      alerts,
		Out:         filepath.Join(t.TempDir(), "log.jsonl"),
		Concurrency: 2,
		Clock:       newFakeClock(testNow),
		Notes:       notes,
	}.withDefaults()
}

// watchTestSource is a fakeSource with the real ruler fixture and a scripted
// state response for the live rule only. The paused rule is deliberately
// unscripted: fakeSource errors on an unscripted title, so any attempt to poll
// it fails the test rather than passing silently.
func watchTestSource(t *testing.T, obs Observation) *fakeSource {
	t.Helper()
	src := newFakeSource()
	src.version = "13.1.0"
	src.defs = rulerDefs(t)
	src.script(watchActiveTitle, obs, nil)
	return src
}

func liveObservation(grafanaNow time.Time) Observation {
	return observation(grafanaNow, testStateRule(watchActiveUID, watchActiveTitle, time.Minute, grafanaNow,
		testInstance(StateNormal, "", "a")))
}

// A rule paused in its definition is skipped, never waited for. Waiting for one
// either hangs forever or errors before the deploy — and the header must still
// name it, so check can report it as skipped rather than lose it.
func TestPrepareWatchDoesNotWaitForPausedRules(t *testing.T) {
	var notes strings.Builder
	cfg := watchTestConfig(t, &notes, "uid:"+watchActiveUID, "uid:"+watchPausedUID)
	src := watchTestSource(t, liveObservation(testNow))

	prep, err := prepareWatch(context.Background(), cfg, src)
	require.NoError(t, err)
	require.NoError(t, prep.writer.Close())

	header, polls, sentinel, err := ReadLog(cfg.Out)
	require.NoError(t, err)
	require.Nil(t, sentinel, "the parent wrote a sentinel; that would tell check the recording ended before the child started")

	require.Len(t, header.Rules, 2)
	for _, lr := range header.Rules {
		require.Positive(t, lr.PollEverySeconds,
			"check needs a positive cadence to derive maxGap from")
		if lr.UID == watchPausedUID {
			require.True(t, lr.IsPaused, "want the resolve-time snapshot to say true")
		}
	}

	// One poll, for the live rule only — and it is already in the log before
	// prepareWatch returned, which is the whole point of the record step.
	require.Len(t, polls, 1)
	require.Equal(t, watchActiveUID, polls[0].RuleUID)
	require.True(t, polls[0].Found)
	require.True(t, polls[0].GrafanaNow.Equal(testNow))
	// The poll record holds the state histogram, asserted through a real
	// prepareWatch()/Reducer call rather than log_test.go's hand-built
	// Writer/ReadLog round trip.
	require.Equal(t, map[string]int{"normal": 1}, polls[0].Histogram)
	require.Contains(t, notes.String(), watchPausedTitle)
	require.Contains(t, notes.String(), "paused")
}

// One authority for the cadence, from the writing side: whatever
// --poll-interval resolves to is what the header records, because that is the
// only value check may derive maxGap from.
func TestPrepareWatchHeaderRecordsTheOverriddenCadence(t *testing.T) {
	var notes strings.Builder
	cfg := watchTestConfig(t, &notes, "uid:"+watchActiveUID)
	cfg.PollEvery = 120 * time.Second // the rule evaluates every 60s
	src := watchTestSource(t, liveObservation(testNow))

	prep, err := prepareWatch(context.Background(), cfg, src)
	require.NoError(t, err)
	defer prep.writer.Close()

	require.Equal(t, float64(120), prep.header.Rules[0].PollEverySeconds,
		"the override, used verbatim and never clamped")
	require.Equal(t, 240*time.Second, prep.timings[watchActiveUID].maxGap)
	require.Contains(t, notes.String(), "--poll-interval")
}

// The budget check runs on the latencies the parent just measured, before the
// deploy runs.
func TestPrepareWatchFailsWhenTheScheduleDoesNotFit(t *testing.T) {
	var notes strings.Builder
	cfg := watchTestConfig(t, &notes, "uid:"+watchActiveUID)

	obs := liveObservation(testNow)
	obs.Latency = 60 * time.Second // against a 30s cadence
	src := watchTestSource(t, obs)

	_, err := prepareWatch(context.Background(), cfg, src)
	require.Error(t, err, "a schedule cannot hold its own cadence")
	assertBudgetMessage(t, err.Error())
}

// Normal instances are verified visible at the one place it is still cheap:
// the first observation. If the state endpoint stops returning them, the
// reduction's predicate quietly inverts.
func TestPrepareWatchVerifiesNormalInstancesAreVisible(t *testing.T) {
	var notes strings.Builder
	cfg := watchTestConfig(t, &notes, "uid:"+watchActiveUID)

	rule := testStateRule(watchActiveUID, watchActiveTitle, time.Minute, testNow, testInstance(StateFiring, "", "b"))
	rule.Totals = map[string]int{"alerting": 1, "normal": 4} // claims normals it did not return
	src := watchTestSource(t, observation(testNow, rule))

	_, err := prepareWatch(context.Background(), cfg, src)
	require.Error(t, err, "totals claim normal instances the response omitted")
	require.Contains(t, err.Error(), "no longer returns normal instances")

	// The failure happens before any poll is appended, so the log holds a
	// header and nothing else.
	_, polls, _, readErr := ReadLog(cfg.Out)
	require.NoError(t, readErr)
	require.Empty(t, polls)
}

func TestPrepareWatchRejectsAnUnsupportedGrafana(t *testing.T) {
	var notes strings.Builder
	cfg := watchTestConfig(t, &notes, "uid:"+watchActiveUID)
	src := watchTestSource(t, liveObservation(testNow))
	src.version = "12.4.0"

	_, err := prepareWatch(context.Background(), cfg, src)
	require.Error(t, err)
	require.Contains(t, err.Error(), "12.4.0")
	require.Contains(t, err.Error(), "13.0.0")
}

// A rule that resolved in the ruler API but is absent from the state endpoint
// is recorded as Found=false — authoritative evidence the coverage proof turns
// into unobservable — not silently dropped.
func TestPrepareWatchNotesAnAbsentRule(t *testing.T) {
	var notes strings.Builder
	cfg := watchTestConfig(t, &notes, "uid:"+watchActiveUID)
	src := watchTestSource(t, observation(testNow)) // an authoritative, empty 2xx

	prep, err := prepareWatch(context.Background(), cfg, src)
	require.NoError(t, err)
	require.NoError(t, prep.writer.Close())

	_, polls, _, err := ReadLog(cfg.Out)
	require.NoError(t, err)
	require.Len(t, polls, 1)
	require.False(t, polls[0].Found, "want one poll recorded as not found")
	require.Contains(t, notes.String(), "absent from the state endpoint")
}

func TestWatchConfigValidation(t *testing.T) {
	base := func() WatchConfig {
		return WatchConfig{
			URL:    "https://grafana.example.com",
			Alerts: []string{"Example"},
			Out:    filepath.Join(t.TempDir(), "log.jsonl"),
			Clock:  newFakeClock(testNow),
		}
	}

	for _, tc := range []struct {
		name   string
		mutate func(*WatchConfig)
		want   string
	}{
		{"no url", func(c *WatchConfig) { c.URL = "" }, "url"},
		{"no log path", func(c *WatchConfig) { c.Out = "" }, "no log path"},
		{"no alerts", func(c *WatchConfig) { c.Alerts = nil }, "no alert names"},
		{"blank alerts only", func(c *WatchConfig) { c.Alerts = []string{"", "  "} }, "no alert names"},
		{"until in the past", func(c *WatchConfig) { c.Until = testNow.Add(-time.Second) }, "not in the future"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cfg := base()
			tc.mutate(&cfg)
			err := cfg.withDefaults().validate()
			require.Errorf(t, err, "validate: no error, want one naming %q", tc.want)
			require.Containsf(t, err.Error(), tc.want, "validate error")
		})
	}

	t.Run("defaults derive the pidfile and daemon log from the log path", func(t *testing.T) {
		cfg := base().withDefaults()
		require.Equal(t, cfg.Out+".pid", cfg.PidFile)
		require.NotEmpty(t, cfg.DaemonLog, "a detached child would have nowhere to explain a failure")
		require.NoError(t, cfg.validate())
	})
}

// The fail-open direction, checked on the child's side: a log recorded at 5s on
// a 300s rule must schedule at 5s.
// Re-deriving from the interval would give 150s — and every real 250s hole in
// that recording would pass.
func TestChildScheduleUsesTheRecordedCadence(t *testing.T) {
	h := Header{Rules: []LoggedRule{
		{UID: "fast", Title: "Fast", IntervalSeconds: 300, PollEverySeconds: 5},
		{UID: "paused", Title: "Paused", IntervalSeconds: 60, PollEverySeconds: 30, IsPaused: true},
	}}

	titles, cadence, err := childSchedule(h)
	require.NoError(t, err)
	_, ok := titles["paused"]
	require.False(t, ok, "the child scheduled a rule that was paused when the window opened")
	require.Equal(t, 5*time.Second, cadence["fast"])
}

func TestChildScheduleRejectsAnUnusableHeader(t *testing.T) {
	for _, tc := range []struct {
		name string
		h    Header
		want string
	}{
		{
			"no recorded cadence",
			Header{Rules: []LoggedRule{{UID: "r1", Title: "Example", IntervalSeconds: 60}}},
			"no cadence",
		},
		{
			"the same rule twice",
			Header{Rules: []LoggedRule{
				{UID: "r1", Title: "Example", IntervalSeconds: 60, PollEverySeconds: 30},
				{UID: "r1", Title: "Example", IntervalSeconds: 60, PollEverySeconds: 300},
			}},
			"twice",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, _, err := childSchedule(tc.h)
			require.Errorf(t, err, "childSchedule: no error, want one naming %q", tc.want)
			require.Contains(t, err.Error(), tc.want)
		})
	}
}

// TestChildArgsCarryNoSecretsAndNoRuleSet: the child's command line lands in
// the process table and in CI logs. Everything it needs about the rules comes
// from the header, and everything about the connection comes from the
// environment — so argv holds the log path and the two run facts only.
func TestChildArgsCarryNoSecretsAndNoRuleSet(t *testing.T) {
	cfg := WatchConfig{
		URL:         "https://grafana.example.com",
		Token:       "secret-token",
		Alerts:      []string{"Example"},
		Folder:      "F",
		Out:         "/tmp/log.jsonl",
		PidFile:     "/tmp/log.jsonl.pid",
		Until:       testNow.Add(time.Hour),
		PollEvery:   17 * time.Second,
		Concurrency: 3,
	}
	args := childArgs(cfg)
	joined := strings.Join(args, " ")

	for _, want := range []string{DaemonChildFlag, "--out /tmp/log.jsonl", "--concurrency 3", "--until ", ReadyFDFlag + " 3"} {
		require.Contains(t, joined, want)
	}
	for _, forbidden := range []string{"secret-token", "Example", "--folder", "--poll-interval", "--pidfile"} {
		require.NotContains(t, joined, forbidden)
	}
}

func TestPidFileRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "log.jsonl.pid")
	require.NoError(t, writePidFile(path, 4242))
	pid, err := ReadPidFile(path)
	require.NoError(t, err)
	require.Equal(t, 4242, pid)

	t.Run("garbage is an error, never a pid", func(t *testing.T) {
		bad := filepath.Join(t.TempDir(), "bad.pid")
		require.NoError(t, os.WriteFile(bad, []byte("not-a-pid\n"), 0o644))
		_, err := ReadPidFile(bad)
		require.Error(t, err, "no error on an unparseable pidfile")
	})
}
