package gate

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// The one rule every test in this file watches, unless it says otherwise: a
// 60s evaluation interval, no `for`, not paused. Every derived value follows
// from those three numbers, and the tests assert against them by name rather
// than by magic constant:
//
//	pollEvery       30s   (intervalSeconds/2)
//	maxGap          60s   (2 x pollEvery)
//	healthGrace     60s   (max(maxGap, interval))
//	evalStaleAfter  120s  (2 x interval)
//	transitionGrace 60s   (for + interval)
//	drainTimeout    2m    (max(2 x interval, 2m))
const (
	checkUID   = "rule-one"
	checkTitle = "Rule One"

	checkPollEvery  = 30 * time.Second
	checkGrace      = 60 * time.Second
	checkDrainLimit = 2 * time.Minute
)

func checkDef() Definition {
	return Definition{
		UID: checkUID, Title: checkTitle, Folder: "F", Group: "G",
		IntervalSeconds: 60, NoDataState: "OK", ExecErrState: "OK",
		Kind: KindGrafanaManaged,
	}
}

// checkSource is a Source whose state answers depend on virtual time and on
// the call count, which is what a collection-loop test needs: fakeSource's
// static script cannot express "healthy for the whole window" without
// scripting every poll, and loopSource (watch_test.go) deliberately refuses
// Version and Definitions because the recorder's child never reads them.
type checkSource struct {
	mu sync.Mutex

	version    string
	versionErr error
	defs       []Definition
	defsErr    error

	calls   map[string]int
	respond func(title string, call int) (Observation, error)
}

func newCheckSource(respond func(title string, call int) (Observation, error)) *checkSource {
	return &checkSource{
		version: "13.1.0",
		defs:    []Definition{checkDef()},
		calls:   map[string]int{},
		respond: respond,
	}
}

func (s *checkSource) Version(context.Context) (string, error) { return s.version, s.versionErr }

func (s *checkSource) Definitions(context.Context) ([]Definition, error) { return s.defs, s.defsErr }

// RuleState answers from the responder. A nil responder means the test
// expects no state read at all — it fails with a message rather than a nil
// dereference, because "this path must not poll" is an assertion several tests
// here make on purpose.
func (s *checkSource) RuleState(_ context.Context, title string) (Observation, error) {
	s.mu.Lock()
	s.calls[title]++
	call := s.calls[title]
	s.mu.Unlock()
	if s.respond == nil {
		return Observation{}, fmt.Errorf("checkSource: this test expects no state read, but %q was polled", title)
	}
	return s.respond(title, call)
}

func (s *checkSource) callCount(title string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.calls[title]
}

var _ Source = (*checkSource)(nil)

// checkStateRule builds one state-endpoint rule whose totals agree with the
// instances it carries. That agreement is load-bearing: a totals map claiming
// normal instances that the instance list does not contain fails
// VerifyNormalInstancesVisible, which is a different failure from the one most
// of these tests are about.
func checkStateRule(lastEval time.Time, insts ...Instance) StateRule {
	totals := map[string]int{}
	for _, i := range insts {
		totals[string(i.State)]++
	}
	return StateRule{
		UID: checkUID, Title: checkTitle, Folder: "F", Group: "G",
		Interval: time.Minute, State: "inactive", Health: "ok",
		LastEvaluation: lastEval, Totals: totals, Instances: insts,
	}
}

// healthyObservation is a poll of a rule that evaluated at this instant, with
// no skew at all — every test that is not ABOUT skew uses zero so its
// arithmetic reads directly off the timestamps.
func healthyObservation(now time.Time, insts ...Instance) Observation {
	return Observation{
		Rules:      []StateRule{checkStateRule(now, insts...)},
		GrafanaNow: now,
		Latency:    200 * time.Millisecond,
	}
}

// baseConfig is a single-step run over [now, now+5m]: window 5m, grace 60s, so
// the collection loop ends at now+6m.
func baseConfig(t *testing.T, clock Clock) Config {
	t.Helper()
	now := clock.Now()
	return Config{
		URL:    "https://grafana.example.com",
		Alerts: []string{"uid:" + checkUID},
		From:   now,
		To:     now.Add(5 * time.Minute),
		Clock:  clock,
		Notes:  &strings.Builder{},
	}.withDefaults()
}

func notesOf(cfg Config) string { return cfg.Notes.(*strings.Builder).String() }

// ---------------------------------------------------------------------------
// Configuration validation
// ---------------------------------------------------------------------------

func TestCheckValidateRejectsBadConfigurations(t *testing.T) {
	clock := newFakeClock(testNow)
	base := func() Config {
		return Config{
			URL:   "https://grafana.example.com",
			From:  testNow,
			To:    testNow.Add(5 * time.Minute),
			Clock: clock,
		}
	}

	tests := []struct {
		name    string
		mutate  func(*Config)
		wantErr string
	}{
		{
			name:    "no url",
			mutate:  func(c *Config) { c.URL = ""; c.Alerts = []string{"A"} },
			wantErr: "no grafana url",
		},
		{
			name:    "no to",
			mutate:  func(c *Config) { c.To = time.Time{}; c.Alerts = []string{"A"} },
			wantErr: "no `to`",
		},
		{
			// An empty Alerts is an error — but only without a log.
			name:    "single-step without alerts",
			mutate:  func(c *Config) {},
			wantErr: "no alert names given",
		},
		{
			// An alerts file ending in a newline must not read as a named alert.
			name:    "single-step with only blank alert lines",
			mutate:  func(c *Config) { c.Alerts = []string{"", "   "} },
			wantErr: "no alert names given",
		},
		{
			// The other direction: with a log, the log names the alert set.
			name:    "log mode with alerts",
			mutate:  func(c *Config) { c.Log = "log.jsonl"; c.Alerts = []string{"A"} },
			wantErr: "--alerts is refused with a recorded log",
		},
		{
			// Never a warning-and-continue.
			name:    "log mode without from",
			mutate:  func(c *Config) { c.Log = "log.jsonl"; c.From = time.Time{} },
			wantErr: "the deploy step must emit a completion timestamp",
		},
		{
			name: "from beyond the future tolerance",
			mutate: func(c *Config) {
				c.Alerts = []string{"A"}
				c.From = testNow.Add(2 * time.Minute)
				c.To = testNow.Add(10 * time.Minute)
			},
			wantErr: "ahead of this runner's clock",
		},
		{
			name:    "to before from",
			mutate:  func(c *Config) { c.Alerts = []string{"A"}; c.To = testNow.Add(-time.Minute) },
			wantErr: "is before `from`",
		},
		{
			// A past `to` is only "not a special mode" WITH a log: without one
			// the coverage window ends before the first observation exists.
			name: "single-step with a to already past",
			mutate: func(c *Config) {
				c.Alerts = []string{"A"}
				c.From = testNow.Add(-10 * time.Minute)
				c.To = testNow.Add(-time.Minute)
			},
			wantErr: "can only be classified from a recording",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cfg := base()
			tc.mutate(&cfg)
			err := cfg.withDefaults().validate()
			require.Errorf(t, err, "validate()")
			require.Contains(t, err.Error(), tc.wantErr)
		})
	}
}

// A past `to` WITH a log is not a special mode: the collection loop's condition
// is already true and the evidence classifies immediately. No branch, and no
// refusal.
func TestCheckValidateAcceptsAPastToWithALog(t *testing.T) {
	cfg := Config{
		URL:   "https://grafana.example.com",
		Log:   "log.jsonl",
		From:  testNow.Add(-10 * time.Minute),
		To:    testNow.Add(-time.Minute),
		Clock: newFakeClock(testNow),
	}.withDefaults()

	require.NoError(t, cfg.validate())
	require.Equal(t, "log.jsonl.pid", cfg.PidFile)
}

// ---------------------------------------------------------------------------
// Single-step mode
// ---------------------------------------------------------------------------

func TestCheckSingleStepCleanWindowPasses(t *testing.T) {
	clock := newVirtualClock(testNow)
	cfg := baseConfig(t, clock)
	src := newCheckSource(func(_ string, _ int) (Observation, error) {
		return healthyObservation(clock.Now()), nil
	})

	res, err := check(context.Background(), cfg, src)
	require.NoError(t, err)
	// A pass is exactly this shape.
	require.Empty(t, res.Violations)
	require.Len(t, res.Verdicts, 1)
	require.Equal(t, OutcomeClean, res.Verdicts[0].Outcome)
	cov := res.Coverage[checkUID]
	require.True(t, cov.Proved)
	require.False(t, cov.Unobservable)

	// The collection loop ran to to+transitionGrace and no further.
	windowEnd := cfg.To.Add(checkGrace)
	require.False(t, clock.Now().Before(windowEnd))
	// One measurement-pass poll plus one every 30s across the 6-minute
	// collection, plus the drain wait's own polls. The exact count depends on
	// the scheduler's random stagger, so assert the order of magnitude a full
	// window implies rather than an exact number.
	require.GreaterOrEqual(t, src.callCount(checkTitle), 12)
	require.Contains(t, notesOf(cfg), "planned run time")
}

// resolve_test.go proves the collapse-note-plus-satisfied-MinObserved path at
// Resolve() directly; this drives the same shape through check() end to end —
// the two input names must collapse to one verdict, the run must pass, and the
// collapse note must reach the run's own notes, not just Resolve()'s return
// value.
func TestCheckSingleStepDuplicateAlertNamesCollapseWithNote(t *testing.T) {
	clock := newVirtualClock(testNow)
	cfg := baseConfig(t, clock)
	cfg.Alerts = []string{"uid:" + checkUID, checkTitle} // the same rule, named two different ways
	src := newCheckSource(func(_ string, _ int) (Observation, error) {
		return healthyObservation(clock.Now()), nil
	})

	res, err := check(context.Background(), cfg, src)
	require.NoError(t, err)
	require.Len(t, res.Verdicts, 1, "the duplicate must collapse to a single rule")
	require.Empty(t, res.Violations, "MinObserved must be satisfied by the post-collapse count of 1")
	require.Contains(t, notesOf(cfg), "counted once")
}

// A rule with health=error for the whole window is unobservable, exit 2 —
// driven from the real "[JD] No Job Proposals" capture (testdata/README.md),
// not a synthetic Poll table, so a change in how the real payload shapes
// health/lastError cannot slip past a hand-built fixture that happens to still
// look right.
func TestCheckSingleStepContinuousHealthErrorIsUnobservable(t *testing.T) {
	body := readFixture(t, "state_health_error.json")
	rules, err := ParseState(body)
	require.NoError(t, err)
	base := rules[0]
	def := Definition{
		UID: base.UID, Title: base.Title, Folder: base.Folder, Group: base.Group,
		IntervalSeconds: int(base.Interval / time.Second), NoDataState: "OK", ExecErrState: "OK",
		Kind: KindGrafanaManaged,
	}

	clock := newVirtualClock(testNow)
	cfg := Config{
		URL: "https://grafana.example.com", Alerts: []string{"uid:" + def.UID},
		From: testNow, To: testNow.Add(5 * time.Minute), Clock: clock, Notes: &strings.Builder{},
	}.withDefaults()

	src := newCheckSource(func(_ string, _ int) (Observation, error) {
		// Every field but LastEvaluation stays exactly as the real capture
		// shaped it (health=error, the real lastError text, the real Error
		// instance); LastEvaluation tracks the poll so staleness — a
		// different coverage check — never becomes the actual cause.
		r := base
		r.LastEvaluation = clock.Now()
		return Observation{Rules: []StateRule{r}, GrafanaNow: clock.Now(), Latency: 200 * time.Millisecond}, nil
	})
	src.defs = []Definition{def}

	res, err := check(context.Background(), cfg, src)
	require.Error(t, err, "continuous health=error must be unobservable")
	require.Len(t, res.Verdicts, 1)
	require.Equal(t, OutcomeUnobservable, res.Verdicts[0].Outcome)
	require.Equal(t, ReasonHealthError, res.Coverage[def.UID].Reason)
}

// A certain violation does not release the runner early, and it does not stop
// the gate reporting exit-1 shape — violations with a nil error.
func TestCheckSingleStepFiringInstanceReportsWithoutExitingEarly(t *testing.T) {
	clock := newVirtualClock(testNow)
	cfg := baseConfig(t, clock)
	firing := Instance{
		Labels:   map[string]string{"alertname": "Rule One", "instance": "a"},
		State:    StateFiring,
		ActiveAt: testNow.Add(-10 * time.Minute), // bad before the window opened
	}
	src := newCheckSource(func(_ string, _ int) (Observation, error) {
		return healthyObservation(clock.Now(), firing), nil
	})

	res, err := check(context.Background(), cfg, src)
	require.NoError(t, err, "a violation is exit 1, not an error")
	require.Len(t, res.Violations, 1)
	require.Equal(t, OutcomePersistentlyBad, res.Violations[0].Outcome)
	require.False(t, clock.Now().Before(cfg.To.Add(checkGrace)), "exited early; collection must run to to+grace")
}

// A newly_bad instance at from+30s gives exit 1, but ONLY after
// to+transitionGrace. The test above covers a rule already bad before the
// window opened (persistently_bad); this covers a fresh onset just inside the
// window, which must not release the runner the instant it is first observed.
func TestCheckSingleStepNewOnsetDoesNotExitEarly(t *testing.T) {
	clock := newVirtualClock(testNow)
	cfg := baseConfig(t, clock)
	onset := testNow.Add(30 * time.Second)
	src := newCheckSource(func(_ string, _ int) (Observation, error) {
		now := clock.Now()
		if now.Before(onset) {
			return healthyObservation(now), nil
		}
		firing := Instance{
			Labels:   map[string]string{"alertname": checkTitle, "instance": "a"},
			State:    StateFiring,
			ActiveAt: onset,
		}
		return healthyObservation(now, firing), nil
	})

	res, err := check(context.Background(), cfg, src)
	require.NoError(t, err)
	require.Len(t, res.Violations, 1)
	require.Equal(t, OutcomeNewlyBad, res.Violations[0].Outcome)
	require.False(t, clock.Now().Before(cfg.To.Add(checkGrace)),
		"exited early; collection must run to to+grace even for a fresh onset at from+30s")
}

// An ABSENT `from` in single-step mode (as opposed to recorder mode, which
// hard-errors — TestCheckValidateRejectsBadConfigurations's "log mode without
// from") falls back to the start of this check step, with the same
// declared-blind-interval warning as an explicit early `from`.
func TestCheckSingleStepAbsentFromFallsBackToStepStart(t *testing.T) {
	clock := newVirtualClock(testNow)
	cfg := baseConfig(t, clock)
	cfg.From = time.Time{}
	src := newCheckSource(func(_ string, _ int) (Observation, error) {
		return healthyObservation(clock.Now()), nil
	})

	res, err := check(context.Background(), cfg, src)
	require.NoError(t, err, "an absent `from` in single-step mode is a fallback, not an error")
	require.Contains(t, notesOf(cfg), "no `from` given")
	require.True(t, res.From.Equal(testNow))
}

// In single-step mode an explicit `from` earlier than the first observation is
// a DECLARED blind interval — a warning and a pass, naming the exact interval
// it cannot see. Recorder mode keeps the from-bounds coverage check strict.
func TestCheckSingleStepFromBeforeFirstObservationWarnsAndPasses(t *testing.T) {
	clock := newVirtualClock(testNow)
	cfg := baseConfig(t, clock)
	cfg.From = testNow.Add(-2 * time.Minute)
	src := newCheckSource(func(_ string, _ int) (Observation, error) {
		return healthyObservation(clock.Now()), nil
	})

	res, err := check(context.Background(), cfg, src)
	require.NoError(t, err)
	notes := notesOf(cfg)
	require.Contains(t, notes, "cannot see [")
	require.Contains(t, notes, testNow.Format(time.RFC3339))
	// The classified window is the clamped one, and Result says so rather than
	// reporting a window the run never proved.
	require.True(t, res.From.Equal(testNow))
}

// The failure limit was exceeded. The measurement pass succeeds and the
// collection loop then hits a terminal failure, so this exercises the path a
// live run really takes.
func TestCheckFailClosedOnExhaustedRetries(t *testing.T) {
	clock := newVirtualClock(testNow)
	cfg := baseConfig(t, clock)
	src := newCheckSource(func(_ string, call int) (Observation, error) {
		if call > 1 {
			return Observation{}, &RetryExhaustedError{Failures: 6, Cause: errors.New("connection refused")}
		}
		return healthyObservation(clock.Now()), nil
	})

	res, err := check(context.Background(), cfg, src)
	require.Error(t, err, "the collection failure to fail closed")
	require.Contains(t, err.Error(), "collect evidence")
	require.Empty(t, res.Violations, "an error must never be reported as a verdict")
}

// The resolution of the definitions failed. Both shapes — the ruler read
// itself failing, and a name that resolves to nothing.
func TestCheckFailClosedOnDefinitionResolution(t *testing.T) {
	t.Run("ruler read fails", func(t *testing.T) {
		clock := newVirtualClock(testNow)
		cfg := baseConfig(t, clock)
		src := newCheckSource(nil)
		src.defsErr = errors.New("502 bad gateway")

		_, err := check(context.Background(), cfg, src)
		require.Error(t, err)
		require.Contains(t, err.Error(), "read rule definitions")
	})

	t.Run("unknown alert name", func(t *testing.T) {
		clock := newVirtualClock(testNow)
		cfg := baseConfig(t, clock)
		cfg.Alerts = []string{"No Such Rule"}
		src := newCheckSource(nil)

		_, err := check(context.Background(), cfg, src)
		require.Error(t, err)
		require.Contains(t, err.Error(), "no rule matched")
	})
}

// The version gate: an unsupported Grafana is exit 2 before anything else is
// attempted.
func TestCheckRefusesUnsupportedGrafanaVersion(t *testing.T) {
	clock := newVirtualClock(testNow)
	cfg := baseConfig(t, clock)
	src := newCheckSource(nil)
	src.version = "12.4.0"

	_, err := check(context.Background(), cfg, src)
	require.Error(t, err)
	require.Contains(t, err.Error(), "unsupported grafana version")
}

// The budget is checked against the latencies the measurement pass actually
// measured, and a schedule that cannot fit errors at START rather than
// producing a gap-riddled recording nobody can classify.
func TestCheckSingleStepRefusesAScheduleThatDoesNotFit(t *testing.T) {
	clock := newVirtualClock(testNow)
	cfg := baseConfig(t, clock)
	src := newCheckSource(func(_ string, _ int) (Observation, error) {
		obs := healthyObservation(clock.Now())
		obs.Latency = 45 * time.Second // longer than the rule's own 30s cadence
		return obs, nil
	})

	_, err := check(context.Background(), cfg, src)
	require.Error(t, err, "the budget check to refuse the schedule")
	for _, want := range []string{"raising concurrency", "raising poll-interval", "watching fewer alerts"} {
		require.Contains(t, err.Error(), want)
	}
}

// ---------------------------------------------------------------------------
// Recorder mode
// ---------------------------------------------------------------------------

// recordedLog writes a log the way watch would have: a header, one poll every
// 30s over [start, end], and a stopped sentinel at sentinelAt. lastEvalLag is
// how far behind each poll's own GrafanaNow its lastEvaluation sits, which is
// what the drain-wait tests vary.
func recordedLog(t *testing.T, dir string, url string, startedAt, start, end, sentinelAt time.Time, lastEvalLag time.Duration) string {
	t.Helper()
	path := filepath.Join(dir, "log.jsonl")
	clock := newFakeClock(sentinelAt)
	w, err := NewWriter(path, clock)
	require.NoError(t, err)
	header := Header{
		URL:            url,
		GrafanaVersion: "13.1.0",
		StartedAt:      startedAt,
		Rules: []LoggedRule{{
			UID: checkUID, Title: checkTitle, Folder: "F", Group: "G",
			IntervalSeconds: 60, NoDataState: "OK", ExecErrState: "OK",
			PollEverySeconds: checkPollEvery.Seconds(),
		}},
	}
	require.NoError(t, w.WriteHeader(header))
	for at := start; !at.After(end); at = at.Add(checkPollEvery) {
		require.NoError(t, w.WritePoll(Poll{
			RuleUID: checkUID, GrafanaNow: at, Found: true,
			State: "inactive", Health: "ok", LastEvaluation: at.Add(-lastEvalLag),
		}))
	}
	require.NoError(t, w.Stop())
	return path
}

// deadPid returns a pid that is guaranteed to have exited — the normal state
// of a recorder by the time check signals it, since a recorder given --until
// (or one that finished cleanly) is already gone.
func deadPid(t *testing.T) int {
	t.Helper()
	cmd := exec.Command("/bin/sh", "-c", "exit 0")
	require.NoError(t, cmd.Start())
	pid := cmd.Process.Pid
	require.NoError(t, cmd.Wait())
	return pid
}

func writePid(t *testing.T, path, contents string) {
	t.Helper()
	require.NoError(t, os.WriteFile(path, []byte(contents), 0o644))
}

// recorderConfig points check at a recording of [testNow-1m, windowEnd+30s]
// over the window [testNow, testNow+5m].
func recorderConfig(t *testing.T, clock Clock, logPath string) Config {
	t.Helper()
	return Config{
		URL:   "https://grafana.example.com",
		Log:   logPath,
		From:  testNow,
		To:    testNow.Add(5 * time.Minute),
		Clock: clock,
		Notes: &strings.Builder{},
	}.withDefaults()
}

func TestCheckRecorderModeCleanWindowPasses(t *testing.T) {
	dir := t.TempDir()
	windowEnd := testNow.Add(5*time.Minute + checkGrace)
	logPath := recordedLog(t, dir, "https://grafana.example.com",
		testNow.Add(-time.Minute), testNow.Add(-time.Minute), windowEnd.Add(30*time.Second), windowEnd.Add(30*time.Second), 0)
	writePid(t, logPath+".pid", fmt.Sprintf("%d\n", deadPid(t)))

	clock := newVirtualClock(testNow.Add(time.Minute))
	cfg := recorderConfig(t, clock, logPath)
	// The drain wait is satisfied from the log's own evidence, so the source
	// must never be asked for a state — asserted by the nil responder.
	src := newCheckSource(func(title string, _ int) (Observation, error) {
		require.Fail(t, fmt.Sprintf("the drain wait polled %q although the log already proves the evaluations", title))
		return Observation{}, errors.New("unexpected poll")
	})

	res, err := check(context.Background(), cfg, src)
	require.NoError(t, err)
	require.Empty(t, res.Violations)
	require.Len(t, res.Verdicts, 1)
	require.Equal(t, OutcomeClean, res.Verdicts[0].Outcome)
	require.Equal(t, "13.1.0", res.GrafanaVersion)
	// The collection loop still waited out to+transitionGrace even though the
	// recorder had already finished.
	require.False(t, clock.Now().Before(windowEnd))
}

// The identity of the log is not correct. The check runs against the header
// read EARLY, so it fails before the window's wait rather than after it.
func TestCheckFailClosedOnWrongLogIdentity(t *testing.T) {
	t.Run("different url", func(t *testing.T) {
		dir := t.TempDir()
		windowEnd := testNow.Add(5*time.Minute + checkGrace)
		logPath := recordedLog(t, dir, "https://other.example.com",
			testNow.Add(-time.Minute), testNow.Add(-time.Minute), windowEnd, windowEnd.Add(30*time.Second), 0)

		clock := newVirtualClock(testNow)
		cfg := recorderConfig(t, clock, logPath)
		_, err := check(context.Background(), cfg, newCheckSource(nil))
		require.Error(t, err)
		require.Contains(t, err.Error(), "log identity")
		require.True(t, clock.Now().Equal(testNow), "it must fail before the wait")
	})

	t.Run("rule no longer resolves", func(t *testing.T) {
		dir := t.TempDir()
		windowEnd := testNow.Add(5*time.Minute + checkGrace)
		logPath := recordedLog(t, dir, "https://grafana.example.com",
			testNow.Add(-time.Minute), testNow.Add(-time.Minute), windowEnd, windowEnd.Add(30*time.Second), 0)

		cfg := recorderConfig(t, newVirtualClock(testNow), logPath)
		src := newCheckSource(nil)
		src.defs = []Definition{{UID: "somebody-else", Title: "Other", Kind: KindGrafanaManaged, IntervalSeconds: 60}}

		_, err := check(context.Background(), cfg, src)
		require.Error(t, err)
		require.Contains(t, err.Error(), "log identity")
	})
}

// `from` before the recording's StartedAt is statically knowable from the
// header (immutable line 1), so check fails closed on it BEFORE the window's
// wait — exactly like the identity check above — rather than surfacing a
// from_before_record verdict only after the drain.
func TestCheckFailFastWhenFromPrecedesRecordStart(t *testing.T) {
	dir := t.TempDir()
	startedAt := testNow.Add(time.Minute) // the recording opened a minute AFTER `from`
	windowEnd := testNow.Add(5*time.Minute + checkGrace)
	// The poll range is irrelevant to the assertion: the fail-fast reads
	// StartedAt from the header alone, before any polling would matter.
	logPath := recordedLog(t, dir, "https://grafana.example.com",
		startedAt, startedAt, windowEnd.Add(30*time.Second), windowEnd.Add(30*time.Second), 0)

	clock := newVirtualClock(testNow)
	cfg := recorderConfig(t, clock, logPath) // From = testNow, before StartedAt
	_, err := check(context.Background(), cfg, newCheckSource(nil))
	require.Error(t, err)
	require.Contains(t, err.Error(), "before recording started")
	require.True(t, clock.Now().Equal(testNow), "it must fail before the wait")
}

// A whole-second `from` in the same second as the recording's sub-second
// StartedAt is not a blind interval: the whole-second comparison lets the run
// proceed to a clean pass instead of the fail-fast above.
func TestCheckRecorderModeFromSameSecondAsStartedAtPasses(t *testing.T) {
	dir := t.TempDir()
	windowEnd := testNow.Add(5*time.Minute + checkGrace)
	// StartedAt is 500ms after `from` (testNow via recorderConfig) — the same
	// whole second. Polls still cover the whole window.
	logPath := recordedLog(t, dir, "https://grafana.example.com",
		testNow.Add(500*time.Millisecond), testNow.Add(-time.Minute), windowEnd.Add(30*time.Second), windowEnd.Add(30*time.Second), 0)
	writePid(t, logPath+".pid", fmt.Sprintf("%d\n", deadPid(t)))

	clock := newVirtualClock(testNow.Add(time.Minute))
	cfg := recorderConfig(t, clock, logPath)
	src := newCheckSource(func(title string, _ int) (Observation, error) {
		require.Fail(t, fmt.Sprintf("the drain wait polled %q although the log already proves the evaluations", title))
		return Observation{}, errors.New("unexpected poll")
	})

	res, err := check(context.Background(), cfg, src)
	require.NoError(t, err)
	require.Empty(t, res.Violations)
	require.Len(t, res.Verdicts, 1)
	require.Equal(t, OutcomeClean, res.Verdicts[0].Outcome)
}

// The coverage proof failed: a hole in the middle of the recording is not
// saved by healthy data at both ends.
func TestCheckFailClosedOnCoverageGap(t *testing.T) {
	dir := t.TempDir()
	windowEnd := testNow.Add(5*time.Minute + checkGrace)
	path := filepath.Join(dir, "log.jsonl")
	clock := newFakeClock(windowEnd.Add(30 * time.Second))
	w, err := NewWriter(path, clock)
	require.NoError(t, err)
	require.NoError(t, w.WriteHeader(Header{
		URL: "https://grafana.example.com", GrafanaVersion: "13.1.0", StartedAt: testNow.Add(-time.Minute),
		Rules: []LoggedRule{{UID: checkUID, Title: checkTitle, IntervalSeconds: 60, PollEverySeconds: checkPollEvery.Seconds()}},
	}))
	for at := testNow.Add(-time.Minute); !at.After(windowEnd.Add(30 * time.Second)); at = at.Add(checkPollEvery) {
		// A three-minute hole in the middle of the window.
		if at.After(testNow.Add(time.Minute)) && at.Before(testNow.Add(4*time.Minute)) {
			continue
		}
		require.NoError(t, w.WritePoll(Poll{
			RuleUID: checkUID, GrafanaNow: at, Found: true, State: "inactive", Health: "ok", LastEvaluation: at,
		}))
	}
	require.NoError(t, w.Stop())
	writePid(t, path+".pid", fmt.Sprintf("%d\n", deadPid(t)))

	cfg := recorderConfig(t, newVirtualClock(testNow), path)
	res, err := check(context.Background(), cfg, newCheckSource(nil))
	require.Error(t, err, "the coverage gap to fail closed")
	require.Equal(t, ReasonHeartbeatGap, res.Coverage[checkUID].Reason)
	require.Equal(t, OutcomeUnobservable, res.Verdicts[0].Outcome)
}

// An episode fully between the deploy and the start of the check: recorder
// mode must find this at the LEADING edge of the window too, right after `from`
// (the deploy's completion), not only in the middle
// (TestCheckFailClosedOnCoverageGap above). No poll exists for [from, from+3m),
// so whatever happened there is invisible to every per-poll check and only the
// coverage gap itself can catch it — the reason the recorder exists at all.
func TestCheckRecorderModeFindsAGapImmediatelyAfterTheDeploy(t *testing.T) {
	dir := t.TempDir()
	windowEnd := testNow.Add(5*time.Minute + checkGrace)
	path := filepath.Join(dir, "log.jsonl")
	clock := newFakeClock(windowEnd.Add(30 * time.Second))
	w, err := NewWriter(path, clock)
	require.NoError(t, err)
	require.NoError(t, w.WriteHeader(Header{
		URL: "https://grafana.example.com", GrafanaVersion: "13.1.0", StartedAt: testNow.Add(-time.Minute),
		Rules: []LoggedRule{{UID: checkUID, Title: checkTitle, IntervalSeconds: 60, PollEverySeconds: checkPollEvery.Seconds()}},
	}))
	gapEnd := testNow.Add(3 * time.Minute) // nothing recorded from `from` (testNow) to here
	for at := gapEnd; !at.After(windowEnd.Add(30 * time.Second)); at = at.Add(checkPollEvery) {
		require.NoError(t, w.WritePoll(Poll{
			RuleUID: checkUID, GrafanaNow: at, Found: true, State: "inactive", Health: "ok", LastEvaluation: at,
		}))
	}
	require.NoError(t, w.Stop())
	writePid(t, path+".pid", fmt.Sprintf("%d\n", deadPid(t)))

	cfg := recorderConfig(t, newVirtualClock(testNow), path)
	res, err := check(context.Background(), cfg, newCheckSource(nil))
	require.Error(t, err, "a hole right after the deploy hides whatever happened there")
	require.Len(t, res.Verdicts, 1)
	require.Equal(t, OutcomeUnobservable, res.Verdicts[0].Outcome, "never clean")
}

// The drain limit passed. The recording itself is clean, so this isolates the
// drain wait — the rule simply never evaluates through the end of the window,
// and a rule that cannot answer that question is unobservable.
func TestCheckFailClosedOnDrainTimeout(t *testing.T) {
	dir := t.TempDir()
	windowEnd := testNow.Add(5*time.Minute + checkGrace)
	// A 45s lag keeps every poll inside evalStaleAfter (120s), so the liveness
	// coverage check is silent and only the drain wait can fail.
	logPath := recordedLog(t, dir, "https://grafana.example.com",
		testNow.Add(-time.Minute), testNow.Add(-time.Minute), windowEnd, windowEnd.Add(30*time.Second), 45*time.Second)
	writePid(t, logPath+".pid", fmt.Sprintf("%d\n", deadPid(t)))

	clock := newVirtualClock(testNow)
	cfg := recorderConfig(t, clock, logPath)
	frozen := windowEnd.Add(-45 * time.Second)
	src := newCheckSource(func(_ string, _ int) (Observation, error) {
		now := clock.Now()
		return Observation{
			Rules:      []StateRule{checkStateRule(frozen)},
			GrafanaNow: now,
		}, nil
	})

	res, err := check(context.Background(), cfg, src)
	require.Error(t, err, "the drain limit to fail closed")
	require.Equal(t, ReasonDrainTimeout, res.Coverage[checkUID].Reason)
	require.Equal(t, OutcomeUnobservable, res.Verdicts[0].Outcome)
	require.Contains(t, res.Verdicts[0].Note, "drain limit")
	require.GreaterOrEqual(t, clock.Now().Sub(windowEnd), checkDrainLimit,
		"the rule never evaluates through the window, so the drain wait must run its full limit")
}

// A rule the state endpoint no longer serves is knowable on the FIRST drain
// poll, and the answer is rule_absent — the fault — rather than
// drain_timeout, which would only name the wait. It must not spend the whole
// drain limit to reach it.
func TestCheckDrainWaitNamesADeletedRuleAtOnce(t *testing.T) {
	dir := t.TempDir()
	windowEnd := testNow.Add(5*time.Minute + checkGrace)
	logPath := recordedLog(t, dir, "https://grafana.example.com",
		testNow.Add(-time.Minute), testNow.Add(-time.Minute), windowEnd, windowEnd.Add(30*time.Second), 45*time.Second)
	writePid(t, logPath+".pid", fmt.Sprintf("%d\n", deadPid(t)))

	clock := newVirtualClock(testNow)
	cfg := recorderConfig(t, clock, logPath)
	// An authoritative 2xx that parsed and carries no matching rule. The
	// transport retried every transient failure long before an Observation
	// exists, so this is a deletion, not a hiccup.
	src := newCheckSource(func(_ string, _ int) (Observation, error) {
		return Observation{GrafanaNow: clock.Now()}, nil
	})

	res, err := check(context.Background(), cfg, src)
	require.Error(t, err, "a deleted rule to fail closed")
	require.Equal(t, ReasonRuleAbsent, res.Coverage[checkUID].Reason, "the fault, not the wait")
	require.Equal(t, 1, src.callCount(checkTitle), "the absence is knowable on the first poll")
	require.Less(t, clock.Now().Sub(windowEnd), checkDrainLimit)
}

// ---------------------------------------------------------------------------
// `skipped` comes from the header, not from a definition read after the window
// ---------------------------------------------------------------------------

// pausedAfterWindowLog records a rule that was ACTIVE at record start and that
// fired inside the window. The caller then tells check that the rule's current
// definition says paused — the state somebody set after the fact.
func pausedAfterWindowLog(t *testing.T, dir string, firesAt time.Time, end, sentinelAt time.Time) string {
	t.Helper()
	path := filepath.Join(dir, "log.jsonl")
	w, err := NewWriter(path, newFakeClock(sentinelAt))
	require.NoError(t, err)
	require.NoError(t, w.WriteHeader(Header{
		URL: "https://grafana.example.com", GrafanaVersion: "13.1.0", StartedAt: testNow.Add(-time.Minute),
		Rules: []LoggedRule{{
			UID: checkUID, Title: checkTitle, IntervalSeconds: 60,
			IsPaused: false, PollEverySeconds: checkPollEvery.Seconds(),
		}},
	}))
	firing := Instance{
		Labels:   map[string]string{"alertname": checkTitle, "instance": "a"},
		State:    StateFiring,
		ActiveAt: firesAt,
	}
	for at := testNow.Add(-time.Minute); !at.After(end); at = at.Add(checkPollEvery) {
		p := Poll{RuleUID: checkUID, GrafanaNow: at, Found: true, State: "inactive", Health: "ok", LastEvaluation: at}
		if !at.Before(firesAt) {
			p.State = "firing"
			p.Abnormal = []Instance{firing}
		}
		require.NoError(t, w.WritePoll(p))
	}
	require.NoError(t, w.Stop())
	return path
}

// pausedAfterWindowCheck runs the timeline above. The recording reaches past
// to + transitionGrace, which for this 60s rule is to + 60s: the fresh
// definition says paused, but that no longer shrinks the grace — the header
// does, and the header says the rule was active (deriveGlobalTimings).
func pausedAfterWindowCheck(t *testing.T, allowPaused bool) (Result, error, Config) {
	t.Helper()
	dir := t.TempDir()
	to := testNow.Add(5 * time.Minute)
	end := to.Add(checkGrace + 30*time.Second)
	logPath := pausedAfterWindowLog(t, dir, testNow.Add(2*time.Minute), end, end)
	writePid(t, logPath+".pid", fmt.Sprintf("%d\n", deadPid(t)))

	cfg := recorderConfig(t, newVirtualClock(testNow), logPath)
	cfg.AllowPaused = allowPaused

	src := newCheckSource(nil)
	paused := checkDef()
	paused.IsPaused = true // somebody paused it after the alert started paging
	src.defs = []Definition{paused}

	res, err := check(context.Background(), cfg, src)
	return res, err, cfg
}

// The window's own evidence outranks a definition read after it closed: a rule
// that was active at record start is classified, whatever its pause state is
// by the time check resolves the definitions.
func TestCheckPausingARuleAfterTheWindowDoesNotMakeItSkipped(t *testing.T) {
	res, err, _ := pausedAfterWindowCheck(t, false)
	require.NoError(t, err)
	require.Equal(t, OutcomeNewlyBad, res.Verdicts[0].Outcome,
		"the rule was active for the whole window and fired inside it")
	require.Len(t, res.Violations, 1)
	require.Equal(t, OutcomeNewlyBad, res.Violations[0].Outcome)
	require.NotContains(t, res.Verdicts[0].Note, "paused before the window opened")
}

// The regression pin for the loophole this fix closed. Reading skipped from
// the post-window definition made the rule skipped; --allow-paused then made
// skipped free; and a window in which the alert fired reported exit 0. The
// default message names --allow-paused, so an operator was led straight to it.
func TestCheckAllowPausedCannotExcuseARulePausedAfterItFired(t *testing.T) {
	res, err, _ := pausedAfterWindowCheck(t, true)
	require.NoError(t, err)
	require.NotEmpty(t, res.Violations, "the run passed over a window in which the alert fired")
}

// The other direction, unchanged: a rule the HEADER says was paused when the
// recording opened is genuinely skipped. It has no polls, so no coverage is
// attempted for it, and --allow-paused behaves as it always did.
func TestCheckHeaderPausedRuleStaysSkipped(t *testing.T) {
	dir := t.TempDir()
	windowEnd := testNow.Add(5*time.Minute + checkGrace)
	path := filepath.Join(dir, "log.jsonl")
	w, err := NewWriter(path, newFakeClock(windowEnd.Add(30*time.Second)))
	require.NoError(t, err)
	// Named in the header, is_paused true, and no poll records at all — the
	// shape watch writes for a rule paused before the window opened.
	require.NoError(t, w.WriteHeader(Header{
		URL: "https://grafana.example.com", GrafanaVersion: "13.1.0", StartedAt: testNow.Add(-time.Minute),
		Rules: []LoggedRule{{
			UID: checkUID, Title: checkTitle, IntervalSeconds: 60,
			IsPaused: true, PollEverySeconds: checkPollEvery.Seconds(),
		}},
	}))
	require.NoError(t, w.Stop())
	writePid(t, path+".pid", fmt.Sprintf("%d\n", deadPid(t)))

	run := func(allowPaused bool) (Result, error) {
		cfg := recorderConfig(t, newVirtualClock(testNow), path)
		cfg.AllowPaused = allowPaused
		// The definition is unpaused now; the header still decides.
		src := newCheckSource(func(title string, _ int) (Observation, error) {
			require.Fail(t, fmt.Sprintf("the drain wait polled skipped rule %q", title))
			return Observation{}, errors.New("unexpected poll")
		})
		return check(context.Background(), cfg, src)
	}

	res, err := run(false)
	require.NoError(t, err, "a skipped rule is a known condition, not an inability")
	require.Equal(t, OutcomeSkipped, res.Verdicts[0].Outcome)
	_, ok := res.Coverage[checkUID]
	require.False(t, ok, "a skipped rule has no coverage to prove")
	require.Len(t, res.Violations, 1, "the MinObserved shortfall")

	res, err = run(true)
	require.NoError(t, err)
	require.Empty(t, res.Violations, "with --allow-paused: want a pass")
}

// A paused rule does not evaluate, so it can never catch up: the drain wait
// must conclude on the first poll instead of spending the whole limit.
func TestCheckDrainWaitConcludesAtOnceOnAPausedRule(t *testing.T) {
	dir := t.TempDir()
	windowEnd := testNow.Add(5*time.Minute + checkGrace)
	logPath := recordedLog(t, dir, "https://grafana.example.com",
		testNow.Add(-time.Minute), testNow.Add(-time.Minute), windowEnd, windowEnd.Add(30*time.Second), 45*time.Second)
	writePid(t, logPath+".pid", fmt.Sprintf("%d\n", deadPid(t)))

	clock := newVirtualClock(testNow)
	cfg := recorderConfig(t, clock, logPath)
	src := newCheckSource(func(_ string, _ int) (Observation, error) {
		now := clock.Now()
		rule := checkStateRule(windowEnd.Add(-45 * time.Second))
		rule.IsPaused = true
		return Observation{Rules: []StateRule{rule}, GrafanaNow: now}, nil
	})

	res, err := check(context.Background(), cfg, src)
	require.Error(t, err, "a rule that stopped evaluating must fail closed")
	require.Equal(t, 1, src.callCount(checkTitle), "a paused rule can never catch up")
	require.Equal(t, ReasonDrainTimeout, res.Coverage[checkUID].Reason,
		"the vocabulary is published, so the detail goes in the note")
	require.Contains(t, res.Verdicts[0].Note, "paused before it evaluated through")
	require.Less(t, clock.Now().Sub(windowEnd), checkDrainLimit)
}

// An absent or unparseable pidfile is never "there was nothing to stop". The
// parent writes the pidfile only once the child reports that it is recording,
// so a missing one means the recording never started — and the log must not be
// read at all.
func TestCheckRefusesToReadALogItCannotStop(t *testing.T) {
	windowEnd := testNow.Add(5*time.Minute + checkGrace)

	tests := []struct {
		name    string
		pidfile string // "" = do not create one
	}{
		{name: "missing pidfile"},
		{name: "unparseable pidfile", pidfile: "not-a-pid\n"},
		{name: "empty pidfile", pidfile: ""},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			logPath := recordedLog(t, dir, "https://grafana.example.com",
				testNow.Add(-time.Minute), testNow.Add(-time.Minute), windowEnd, windowEnd.Add(30*time.Second), 0)
			if i != 0 {
				writePid(t, logPath+".pid", tc.pidfile)
			}

			cfg := recorderConfig(t, newVirtualClock(testNow), logPath)
			_, err := check(context.Background(), cfg, newCheckSource(nil))
			require.Error(t, err)
			require.Contains(t, err.Error(), "cannot stop the recorder")
		})
	}
}

// startLockHolder re-execs this test binary as a process that holds the log's
// flock and ignores SIGTERM, and returns its pid once the lock is genuinely
// held. See lockHolderEnv (watch_daemon_test.go) for why it must be a separate
// real process rather than a shell one-liner.
func startLockHolder(t *testing.T, logPath string) int {
	t.Helper()
	cmd := exec.Command(os.Args[0])
	cmd.Env = append(os.Environ(), lockHolderEnv+"="+logPath)
	cmd.Stderr = os.Stderr
	stdout, err := cmd.StdoutPipe()
	require.NoError(t, err)
	require.NoError(t, cmd.Start())
	t.Cleanup(func() {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
	})
	_, err = bufio.NewReader(stdout).ReadString('\n')
	require.NoError(t, err, "the lock holder never reported holding the lock")
	return cmd.Process.Pid
}

// A recorder that will not let go of the log means the log may still be
// appended to, and a log a writer can change cannot be read at all.
func TestCheckFailsWhenTheRecorderWillNotExit(t *testing.T) {
	dir := t.TempDir()
	windowEnd := testNow.Add(5*time.Minute + checkGrace)
	logPath := recordedLog(t, dir, "https://grafana.example.com",
		testNow.Add(-time.Minute), testNow.Add(-time.Minute), windowEnd, windowEnd.Add(30*time.Second), 0)

	writePid(t, logPath+".pid", fmt.Sprintf("%d\n", startLockHolder(t, logPath)))

	cfg := recorderConfig(t, newVirtualClock(testNow), logPath)
	_, err := check(context.Background(), cfg, newCheckSource(nil))
	require.Error(t, err)
	require.Contains(t, err.Error(), "still holds")
}

// The regression pin for a stray SIGTERM. Nothing removes the pidfile when a
// recorder exits cleanly — the parent has returned and the child never learns
// the path — so after a --until run, a supported flow, the pidfile names a pid
// the operating system is free to hand to somebody else. The flock, not the
// pid, is what says whether a writer exists.
func TestCheckDoesNotSignalABystanderHoldingAReusedPid(t *testing.T) {
	dir := t.TempDir()
	windowEnd := testNow.Add(5*time.Minute + checkGrace)
	logPath := recordedLog(t, dir, "https://grafana.example.com",
		testNow.Add(-time.Minute), testNow.Add(-time.Minute), windowEnd.Add(30*time.Second), windowEnd.Add(30*time.Second), 0)

	// An innocent process that happens to hold the pid the finished recorder
	// left behind. It does not hold the log's lock, because it is not a
	// recorder.
	bystander := exec.Command("sleep", "30")
	require.NoError(t, bystander.Start())
	t.Cleanup(func() {
		_ = bystander.Process.Kill()
		_ = bystander.Wait()
	})
	writePid(t, logPath+".pid", fmt.Sprintf("%d\n", bystander.Process.Pid))

	cfg := recorderConfig(t, newVirtualClock(testNow), logPath)
	_, err := check(context.Background(), cfg, newCheckSource(nil))
	require.NoError(t, err)
	require.NoError(t, syscall.Kill(bystander.Process.Pid, 0),
		"check signalled a process that was not the recorder")
}

// A dead pidfile (the recorder process has already exited, holding no flock)
// with NO sentinel in the log — the shape a killed `watch` leaves behind —
// must not hang the stop wait: the flock is free immediately, so
// check reads the log at once, finds no sentinel, and fails closed.
func TestCheckDeadPidWithNoSentinelIsUnobservable(t *testing.T) {
	dir := t.TempDir()
	windowEnd := testNow.Add(5*time.Minute + checkGrace)
	logPath := filepath.Join(dir, "log.jsonl")
	w, err := NewWriter(logPath, newFakeClock(testNow))
	require.NoError(t, err)
	require.NoError(t, w.WriteHeader(Header{
		URL: "https://grafana.example.com", GrafanaVersion: "13.1.0", StartedAt: testNow.Add(-time.Minute),
		Rules: []LoggedRule{{
			UID: checkUID, Title: checkTitle, Folder: "F", Group: "G",
			IntervalSeconds: 60, NoDataState: "OK", ExecErrState: "OK",
			PollEverySeconds: checkPollEvery.Seconds(),
		}},
	}))
	// Healthy heartbeats all the way past windowEnd — evaluatedThrough is
	// satisfied, so the drain wait needs no live re-poll — but no sentinel is
	// ever written: the recorder died before it could call Stop.
	for at := testNow.Add(-time.Minute); !at.After(windowEnd.Add(30 * time.Second)); at = at.Add(checkPollEvery) {
		require.NoError(t, w.WritePoll(Poll{RuleUID: checkUID, GrafanaNow: at, Found: true, State: "inactive", Health: "ok", LastEvaluation: at}))
	}
	require.NoError(t, w.Close()) // no sentinel — a clean exit would call Stop
	writePid(t, logPath+".pid", fmt.Sprintf("%d\n", deadPid(t)))

	cfg := recorderConfig(t, newVirtualClock(testNow), logPath)
	res, err := check(context.Background(), cfg, newCheckSource(nil))
	require.Error(t, err, "no sentinel means the recorder never proved it ran to the end")
	require.Len(t, res.Verdicts, 1)
	require.Equal(t, OutcomeUnobservable, res.Verdicts[0].Outcome)
}

// An incomplete last line gives exit 2. log_test.go's TestReadLogRejectsBadLogs
// pins ReadLog's own error and TestExitCode pins that any non-nil error maps to
// exit 2, but only this feeds a genuinely truncated log through check() itself:
// a raw file with a valid header and poll, then a torn JSON tail, exactly what
// a recorder killed mid-write leaves behind.
func TestCheckRecorderModeTruncatedLogFailsClosed(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "log.jsonl")

	h := Header{
		SchemaVersion: LogSchemaVersion,
		URL:           "https://grafana.example.com", GrafanaVersion: "13.1.0", StartedAt: testNow.Add(-time.Minute),
		Rules: []LoggedRule{{UID: checkUID, Title: checkTitle, IntervalSeconds: 60, PollEverySeconds: checkPollEvery.Seconds()}},
	}
	hb, err := json.Marshal(headerRecord{Type: RecordHeader, Header: h})
	require.NoError(t, err)
	pb, err := json.Marshal(pollRecord{Type: RecordPoll, Poll: Poll{
		RuleUID: checkUID, GrafanaNow: testNow, Found: true, State: "inactive", Health: "ok", LastEvaluation: testNow,
	}})
	require.NoError(t, err)
	content := string(hb) + "\n" + string(pb) + "\n" + `{"type":"poll","rule_ui` // torn mid-write
	require.NoError(t, os.WriteFile(path, []byte(content), 0o644))
	writePid(t, path+".pid", fmt.Sprintf("%d\n", deadPid(t)))

	cfg := recorderConfig(t, newVirtualClock(testNow), path)
	_, err = check(context.Background(), cfg, newCheckSource(nil))
	require.Error(t, err)
	require.Contains(t, err.Error(), "unparseable")
}

// One authority for the cadence, from check's side: maxGap comes from the
// cadence the header records, never from a re-derivation off intervalSeconds.
// The fail-open direction is the one asserted — a log recorded at 5s on a 60s
// rule must still fail on a hole a re-derived 30s maxGap would have forgiven.
func TestCheckDerivesMaxGapFromTheRecordedCadence(t *testing.T) {
	dir := t.TempDir()
	windowEnd := testNow.Add(5*time.Minute + checkGrace)
	path := filepath.Join(dir, "log.jsonl")
	w, err := NewWriter(path, newFakeClock(windowEnd.Add(30*time.Second)))
	require.NoError(t, err)
	require.NoError(t, w.WriteHeader(Header{
		URL: "https://grafana.example.com", GrafanaVersion: "13.1.0", StartedAt: testNow.Add(-time.Minute),
		Rules: []LoggedRule{{UID: checkUID, Title: checkTitle, IntervalSeconds: 60, PollEverySeconds: 5}},
	}))
	for at := testNow.Add(-time.Minute); !at.After(windowEnd.Add(30 * time.Second)); at = at.Add(5 * time.Second) {
		// A 20s hole: under the recorded 5s cadence maxGap is 10s and this
		// fails; under a cadence re-derived from intervalSeconds it would be
		// 60s and the hole would pass unseen.
		if at.After(testNow.Add(time.Minute)) && at.Before(testNow.Add(80*time.Second)) {
			continue
		}
		require.NoError(t, w.WritePoll(Poll{
			RuleUID: checkUID, GrafanaNow: at, Found: true, State: "inactive", Health: "ok", LastEvaluation: at,
		}))
	}
	require.NoError(t, w.Stop())
	writePid(t, path+".pid", fmt.Sprintf("%d\n", deadPid(t)))

	cfg := recorderConfig(t, newVirtualClock(testNow), path)
	res, err := check(context.Background(), cfg, newCheckSource(nil))
	require.Error(t, err, "a 20s hole exceeds the 10s maxGap the recorded 5s cadence implies")
	require.Equal(t, ReasonHeartbeatGap, res.Coverage[checkUID].Reason)
}

// ---------------------------------------------------------------------------
// The pieces, in isolation
// ---------------------------------------------------------------------------

// The drain wait's one comparison is cross-domain, and its uncertainty is
// spent in the fail-closed direction: an evaluation that only MIGHT have
// reached the end of the window does not count as one that did.
func TestEvaluatedThroughSpendsItsUncertaintyFailingClosed(t *testing.T) {
	end := testNow

	tests := []struct {
		name          string
		lastEval      time.Time
		skew, bound   time.Duration
		wantSatisfied bool
	}{
		{name: "zero lastEvaluation never satisfies", lastEval: time.Time{}, wantSatisfied: false},
		{name: "exactly at the end, no skew", lastEval: end, wantSatisfied: true},
		{name: "one second short", lastEval: end.Add(-time.Second), wantSatisfied: false},
		{
			name: "far enough past the end to absorb the bound",
			// Grafana runs 10s fast; the reading translates back to end+5s and
			// the 1s bound still leaves it past the end.
			lastEval: end.Add(16 * time.Second), skew: 10 * time.Second, bound: time.Second,
			wantSatisfied: true,
		},
		{
			name: "inside the bound is not proof",
			// Translated it lands exactly on the end, so the bound can put it
			// either side — which is not an evaluation THROUGH the end.
			lastEval: end.Add(10 * time.Second), skew: 10 * time.Second, bound: time.Second,
			wantSatisfied: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.wantSatisfied, evaluatedThrough(tc.lastEval, tc.skew, tc.bound, end))
		})
	}
}

// A drain timeout on one rule and a coverage failure on another must both
// reach the message. Neither error may shadow the other.
func TestMergeDrainTimeoutsNamesEveryUnobservableRule(t *testing.T) {
	res := Result{
		Coverage: map[string]CoverageResult{
			"a": {Proved: true},
			"b": {Unobservable: true, Reason: ReasonHeartbeatGap, Notes: []string{"rule \"B\": gap"}},
		},
		Verdicts: []RuleVerdict{
			{Alert: "A", RuleUID: "a", Outcome: OutcomeClean},
			{Alert: "B", RuleUID: "b", Outcome: OutcomeUnobservable},
		},
	}

	merged, err := mergeDrainTimeouts(res, map[string]drainVerdict{
		"a": {reason: ReasonDrainTimeout, note: "rule \"A\": did not evaluate through the end within the drain limit"},
		"b": {reason: ReasonDrainTimeout, note: "rule \"B\": did not evaluate through the end within the drain limit"},
	})
	require.Error(t, err, "naming the newly unobservable rule")
	require.Contains(t, err.Error(), "unobservable at the drain wait")
	// Only A is newly unobservable; B was already, so naming it twice would
	// only lengthen the message.
	require.Contains(t, err.Error(), "A ("+string(ReasonDrainTimeout)+")")
	require.NotContains(t, err.Error(), "B (")
	require.Equal(t, ReasonDrainTimeout, merged.Coverage["a"].Reason)
	// B keeps the reason the coverage proof gave it — the FIRST reason wins,
	// as it does inside proveCoverage.
	require.Equal(t, ReasonHeartbeatGap, merged.Coverage["b"].Reason)
	require.Equal(t, OutcomeUnobservable, merged.Verdicts[0].Outcome)
}

// ReadLogHeader is the one read of a log a writer may still hold, so its
// refusals matter as much as its successes.
func TestReadLogHeader(t *testing.T) {
	dir := t.TempDir()

	t.Run("reads line 1 while the log keeps growing", func(t *testing.T) {
		path := filepath.Join(dir, "growing.jsonl")
		w, err := NewWriter(path, newFakeClock(testNow))
		require.NoError(t, err)
		defer w.Close()
		require.NoError(t, w.WriteHeader(testHeader()))
		require.NoError(t, w.WritePoll(Poll{RuleUID: "rule1", GrafanaNow: testNow, Found: true}))

		h, err := ReadLogHeader(path)
		require.NoError(t, err)
		require.Equal(t, testHeader().URL, h.URL)
		require.Len(t, h.Rules, 1)
	})

	t.Run("a half-written header is not a header", func(t *testing.T) {
		path := filepath.Join(dir, "torn.jsonl")
		require.NoError(t, os.WriteFile(path, []byte(`{"type":"header","url":"htt`), 0o644))
		_, err := ReadLogHeader(path)
		require.Error(t, err)
		require.Contains(t, err.Error(), "no complete header")
	})

	t.Run("a wrong schema version is refused", func(t *testing.T) {
		path := filepath.Join(dir, "old.jsonl")
		require.NoError(t, os.WriteFile(path, []byte(`{"type":"header","schema_version":99,"url":"u"}`+"\n"), 0o644))
		_, err := ReadLogHeader(path)
		require.Error(t, err)
		require.Contains(t, err.Error(), "schema version 99")
	})
}
