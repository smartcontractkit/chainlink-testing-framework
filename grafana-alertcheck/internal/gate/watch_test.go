package gate

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"
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
	if err != nil {
		t.Fatalf("NewWriter: %v", err)
	}
	if err := w.WriteHeader(testHeader()); err != nil {
		t.Fatalf("WriteHeader: %v", err)
	}
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
	if err != nil {
		t.Fatalf("watchLoop: %v", err)
	}

	_, polls, sentinel, readErr := ReadLog(path)
	if readErr != nil {
		t.Fatalf("ReadLog: %v", readErr)
	}
	if sentinel == nil {
		t.Fatal("no stopped sentinel after a clean stop")
	}
	if sentinel.Before(testNow.Add(300 * time.Second)) {
		t.Errorf("sentinel at %s, want >= the stop time %s", sentinel, testNow.Add(300*time.Second))
	}
	// 300s of window at 5s and 150s, minus the initial stagger offset of up to
	// one cadence: 59-60 and 1-2. The assertion is the ratio, not the exact
	// count — a single global cycle would give both rules the same number.
	if got := countPolls(polls, tightUID); got < 59 || got > 61 {
		t.Errorf("tight rule polled %d times, want ~60 (300s at 5s)", got)
	}
	if got := countPolls(polls, slackUID); got < 1 || got > 3 {
		t.Errorf("slack rule polled %d times, want ~2 (300s at 150s)", got)
	}
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
	if !errors.Is(err, boom) {
		t.Fatalf("watchLoop error = %v, want %v", err, boom)
	}

	_, polls, sentinel, readErr := ReadLog(path)
	if readErr != nil {
		t.Fatalf("ReadLog: %v", readErr)
	}
	if sentinel != nil {
		t.Errorf("sentinel at %s after a failed recording; check would read that as a finished window", sentinel)
	}
	if len(polls) != 1 {
		t.Errorf("kept %d polls, want the 1 that succeeded before the failure", len(polls))
	}
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

	if err := watchLoop(ctx, watchLoopConfig{
		Src:         src,
		Writer:      w,
		Reducer:     NewReducer(),
		Titles:      map[string]string{"r1": "Example"},
		Cadence:     map[string]time.Duration{"r1": 30 * time.Second},
		Concurrency: 1,
		Clock:       clock,
	}); err != nil {
		t.Fatalf("watchLoop: %v", err)
	}

	if _, _, sentinel, err := ReadLog(path); err != nil {
		t.Fatalf("ReadLog: %v", err)
	} else if sentinel == nil {
		t.Error("no sentinel after a signalled stop; check would call a fully observed window unobservable")
	}
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

	if err := watchLoop(context.Background(), watchLoopConfig{
		Src:         src,
		Writer:      w,
		Reducer:     NewReducer(),
		Titles:      map[string]string{},
		Cadence:     map[string]time.Duration{},
		Until:       testNow.Add(time.Minute),
		Concurrency: 1,
		Clock:       clock,
	}); err != nil {
		t.Fatalf("watchLoop: %v", err)
	}

	_, polls, sentinel, err := ReadLog(path)
	if err != nil {
		t.Fatalf("ReadLog: %v", err)
	}
	if len(polls) != 0 {
		t.Errorf("wrote %d polls with nothing to poll", len(polls))
	}
	if sentinel == nil {
		t.Error("no sentinel: check cannot tell this recording from one that died")
	}
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
	if err := cfg.pollBatch(context.Background(), []string{"ok", "bad"}); !errors.Is(err, boom) {
		t.Fatalf("pollBatch error = %v, want %v", err, boom)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	_, polls, _, err := ReadLog(path)
	if err != nil {
		t.Fatalf("ReadLog: %v", err)
	}
	if len(polls) != 1 || polls[0].RuleUID != "ok" {
		t.Errorf("polls = %+v, want the one heartbeat that was actually observed", polls)
	}
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
		if !slices.Contains(p.Vanished, key) {
			t.Errorf("vanished = %v, want it to contain %q", p.Vanished, key)
		}
		if len(p.Cleared) != 0 {
			t.Errorf("cleared = %v, want none: a vanish is not a recovery", p.Cleared)
		}
	})

	t.Run("unseeded loses the transition", func(t *testing.T) {
		p := NewReducer().Reduce("r1", childObs)
		if len(p.Vanished) != 0 {
			t.Fatalf("vanished = %v; this subtest exists to show the seed is what produces the marker", p.Vanished)
		}
	})

	t.Run("a not-found poll does not clear the seed", func(t *testing.T) {
		r := NewReducer()
		r.seedFrom([]Poll{parentPoll, {RuleUID: "r1", Found: false}})
		if p := r.Reduce("r1", childObs); !slices.Contains(p.Vanished, key) {
			t.Errorf("vanished = %v, want it to contain %q: an absent rule leaves the abnormal set untouched", p.Vanished, key)
		}
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
	if err != nil {
		t.Fatalf("prepareWatch: %v", err)
	}
	if err := prep.writer.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	header, polls, sentinel, err := ReadLog(cfg.Out)
	if err != nil {
		t.Fatalf("ReadLog: %v", err)
	}
	if sentinel != nil {
		t.Error("the parent wrote a sentinel; that would tell check the recording ended before the child started")
	}

	if len(header.Rules) != 2 {
		t.Fatalf("header names %d rules, want both the live and the paused one", len(header.Rules))
	}
	for _, lr := range header.Rules {
		if lr.PollEverySeconds <= 0 {
			t.Errorf("header rule %s records poll_every_seconds=%v; check needs a positive cadence to derive maxGap from", lr.UID, lr.PollEverySeconds)
		}
		if lr.UID == watchPausedUID && !lr.IsPaused {
			t.Errorf("header rule %s: is_paused = false, want the resolve-time snapshot to say true", lr.UID)
		}
	}

	// One poll, for the live rule only — and it is already in the log before
	// prepareWatch returned, which is the whole point of the record step.
	if len(polls) != 1 || polls[0].RuleUID != watchActiveUID {
		t.Fatalf("polls = %+v, want exactly one first observation of %s", polls, watchActiveUID)
	}
	if !polls[0].Found || !polls[0].GrafanaNow.Equal(testNow) {
		t.Errorf("first poll = %+v, want a found observation at %s", polls[0], testNow)
	}
	// The poll record holds the state histogram, asserted through a real
	// prepareWatch()/Reducer call rather than log_test.go's hand-built
	// Writer/ReadLog round trip.
	if want := map[string]int{"normal": 1}; !maps.Equal(polls[0].Histogram, want) {
		t.Errorf("Histogram = %v, want %v: watch must record the state histogram on every poll it writes", polls[0].Histogram, want)
	}
	if !strings.Contains(notes.String(), watchPausedTitle) || !strings.Contains(notes.String(), "paused") {
		t.Errorf("notes do not mention the paused rule:\n%s", notes.String())
	}
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
	if err != nil {
		t.Fatalf("prepareWatch: %v", err)
	}
	defer prep.writer.Close()

	if got := prep.header.Rules[0].PollEverySeconds; got != 120 {
		t.Errorf("header poll_every_seconds = %v, want 120 (the override, used verbatim and never clamped)", got)
	}
	if got := prep.timings[watchActiveUID].maxGap; got != 240*time.Second {
		t.Errorf("maxGap = %s, want 240s (2 x the recorded cadence)", got)
	}
	if !strings.Contains(notes.String(), "--poll-interval") {
		t.Errorf("notes do not report that the override exceeds half the evaluation interval:\n%s", notes.String())
	}
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
	if err == nil {
		t.Fatal("prepareWatch: no error on a schedule that cannot hold its own cadence")
	}
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
	if err == nil {
		t.Fatal("prepareWatch: no error when totals claim normal instances the response omitted")
	}
	if !strings.Contains(err.Error(), "no longer returns normal instances") {
		t.Errorf("error does not say the endpoint stopped returning normal instances: %v", err)
	}

	// The failure happens before any poll is appended, so the log holds a
	// header and nothing else.
	if _, polls, _, readErr := ReadLog(cfg.Out); readErr != nil {
		t.Fatalf("ReadLog: %v", readErr)
	} else if len(polls) != 0 {
		t.Errorf("wrote %d polls from an observation it refused to trust", len(polls))
	}
}

func TestPrepareWatchRejectsAnUnsupportedGrafana(t *testing.T) {
	var notes strings.Builder
	cfg := watchTestConfig(t, &notes, "uid:"+watchActiveUID)
	src := watchTestSource(t, liveObservation(testNow))
	src.version = "12.4.0"

	if _, err := prepareWatch(context.Background(), cfg, src); err == nil {
		t.Fatal("prepareWatch: no error on an unsupported grafana version")
	} else if !strings.Contains(err.Error(), "12.4.0") || !strings.Contains(err.Error(), "13.0.0") {
		t.Errorf("error names neither what was found nor what is supported: %v", err)
	}
}

// A rule that resolved in the ruler API but is absent from the state endpoint
// is recorded as Found=false — authoritative evidence the coverage proof turns
// into unobservable — not silently dropped.
func TestPrepareWatchNotesAnAbsentRule(t *testing.T) {
	var notes strings.Builder
	cfg := watchTestConfig(t, &notes, "uid:"+watchActiveUID)
	src := watchTestSource(t, observation(testNow)) // an authoritative, empty 2xx

	prep, err := prepareWatch(context.Background(), cfg, src)
	if err != nil {
		t.Fatalf("prepareWatch: %v", err)
	}
	if err := prep.writer.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	_, polls, _, err := ReadLog(cfg.Out)
	if err != nil {
		t.Fatalf("ReadLog: %v", err)
	}
	if len(polls) != 1 || polls[0].Found {
		t.Fatalf("polls = %+v, want one poll recorded as not found", polls)
	}
	if !strings.Contains(notes.String(), "absent from the state endpoint") {
		t.Errorf("notes do not warn about the absent rule:\n%s", notes.String())
	}
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
			if err == nil {
				t.Fatalf("validate: no error, want one naming %q", tc.want)
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Errorf("validate error = %v, want it to name %q", err, tc.want)
			}
		})
	}

	t.Run("defaults derive the pidfile and daemon log from the log path", func(t *testing.T) {
		cfg := base().withDefaults()
		if cfg.PidFile != cfg.Out+".pid" {
			t.Errorf("PidFile = %q, want %q — check finds the recorder by this convention", cfg.PidFile, cfg.Out+".pid")
		}
		if cfg.DaemonLog == "" {
			t.Error("DaemonLog is empty: a detached child would have nowhere to explain a failure")
		}
		if err := cfg.validate(); err != nil {
			t.Errorf("validate: %v", err)
		}
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
	if err != nil {
		t.Fatalf("childSchedule: %v", err)
	}
	if _, ok := titles["paused"]; ok {
		t.Error("the child scheduled a rule that was paused when the window opened")
	}
	if got := cadence["fast"]; got != 5*time.Second {
		t.Errorf("pollEvery = %s, want 5s from the header, not %s from the interval", got, defaultPollEvery(300))
	}
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
			if _, _, err := childSchedule(tc.h); err == nil {
				t.Fatalf("childSchedule: no error, want one naming %q", tc.want)
			} else if !strings.Contains(err.Error(), tc.want) {
				t.Errorf("error = %v, want it to name %q", err, tc.want)
			}
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
		if !strings.Contains(joined, want) {
			t.Errorf("child args %q do not contain %q", joined, want)
		}
	}
	for _, forbidden := range []string{"secret-token", "Example", "--folder", "--poll-interval", "--pidfile"} {
		if strings.Contains(joined, forbidden) {
			t.Errorf("child args %q contain %q, which must not reach argv", joined, forbidden)
		}
	}
}

func TestPidFileRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "log.jsonl.pid")
	if err := writePidFile(path, 4242); err != nil {
		t.Fatalf("writePidFile: %v", err)
	}
	pid, err := ReadPidFile(path)
	if err != nil {
		t.Fatalf("ReadPidFile: %v", err)
	}
	if pid != 4242 {
		t.Errorf("pid = %d, want 4242", pid)
	}

	t.Run("garbage is an error, never a pid", func(t *testing.T) {
		bad := filepath.Join(t.TempDir(), "bad.pid")
		if err := os.WriteFile(bad, []byte("not-a-pid\n"), 0o644); err != nil {
			t.Fatalf("write: %v", err)
		}
		if _, err := ReadPidFile(bad); err == nil {
			t.Error("ReadPidFile: no error on an unparseable pidfile")
		}
	})
}
