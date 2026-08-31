package gate

import (
	"bufio"
	"context"
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
)

// The one rule every test in this file watches, unless it says otherwise: a
// 60s evaluation interval, no `for`, not paused. Every derived value follows
// from those three numbers, and the tests assert against them by name rather
// than by magic constant:
//
//	pollEvery       30s   (§5: intervalSeconds/2)
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
// normal instances that the instance list does not contain fails §3.2's
// verification (VerifyNormalInstancesVisible), which is a different failure
// from the one most of these tests are about.
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
// §19.1 step 1 — configuration validation
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
			// §19.1 step 1: an empty Alerts is an error — but only without a log.
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
			// §19.1 step 3, the other direction: the log names the alert set.
			name:    "log mode with alerts",
			mutate:  func(c *Config) { c.Log = "log.jsonl"; c.Alerts = []string{"A"} },
			wantErr: "--alerts is refused with a recorded log",
		},
		{
			// §7 — never a warning-and-continue.
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
			if err == nil {
				t.Fatalf("validate() = nil, want an error containing %q", tc.wantErr)
			}
			if !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("validate() = %q, want it to contain %q", err, tc.wantErr)
			}
		})
	}
}

// A past `to` WITH a log is explicitly not a special mode (§7, §24.3): the
// collection loop's condition is already true and the evidence classifies
// immediately. No branch, and no refusal.
func TestCheckValidateAcceptsAPastToWithALog(t *testing.T) {
	cfg := Config{
		URL:   "https://grafana.example.com",
		Log:   "log.jsonl",
		From:  testNow.Add(-10 * time.Minute),
		To:    testNow.Add(-time.Minute),
		Clock: newFakeClock(testNow),
	}.withDefaults()

	if err := cfg.validate(); err != nil {
		t.Fatalf("validate() = %v, want nil", err)
	}
	if cfg.PidFile != "log.jsonl.pid" {
		t.Errorf("PidFile = %q, want the <in>.pid default", cfg.PidFile)
	}
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
	if err != nil {
		t.Fatalf("check() = %v, want nil\nnotes:\n%s", err, notesOf(cfg))
	}
	// H7: a pass is exactly this shape.
	if len(res.Violations) != 0 {
		t.Fatalf("Violations = %+v, want none", res.Violations)
	}
	if len(res.Verdicts) != 1 || res.Verdicts[0].Outcome != OutcomeClean {
		t.Fatalf("Verdicts = %+v, want one clean verdict", res.Verdicts)
	}
	if cov := res.Coverage[checkUID]; !cov.Proved || cov.Unobservable {
		t.Fatalf("Coverage = %+v, want proved", cov)
	}

	// The collection loop ran to to+transitionGrace and no further (H5).
	windowEnd := cfg.To.Add(checkGrace)
	if clock.Now().Before(windowEnd) {
		t.Errorf("stopped collecting at %s, before to+grace %s", clock.Now(), windowEnd)
	}
	// One measurement-pass poll plus one every 30s across the 6-minute
	// collection, plus the drain wait's own polls. The exact count depends on
	// the scheduler's random stagger, so assert the order of magnitude a full
	// window implies rather than an exact number.
	if got := src.callCount(checkTitle); got < 12 {
		t.Errorf("polled %d times, want at least the ~13 a full 6-minute window at 30s implies", got)
	}
	if notes := notesOf(cfg); !strings.Contains(notes, "planned run time") {
		t.Errorf("§13.2 requires the planned run time at start; notes were:\n%s", notes)
	}
}

// H5: a certain violation does not release the runner early, and it does not
// stop the gate reporting exit-1 shape — violations with a nil error.
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
	if err != nil {
		t.Fatalf("check() = %v, want nil (a violation is exit 1, not an error)", err)
	}
	if len(res.Violations) != 1 {
		t.Fatalf("Violations = %+v, want exactly one", res.Violations)
	}
	if got := res.Violations[0].Outcome; got != OutcomePersistentlyBad {
		t.Errorf("Outcome = %q, want %q", got, OutcomePersistentlyBad)
	}
	if windowEnd := cfg.To.Add(checkGrace); clock.Now().Before(windowEnd) {
		t.Errorf("exited early at %s; H5 requires collecting to %s", clock.Now(), windowEnd)
	}
}

// §4.2/§22.4: in single-step mode an explicit `from` earlier than the first
// observation is a DECLARED blind interval — a warning and a pass, naming the
// exact interval it cannot see. Recorder mode keeps P7 check 2 strict.
func TestCheckSingleStepFromBeforeFirstObservationWarnsAndPasses(t *testing.T) {
	clock := newVirtualClock(testNow)
	cfg := baseConfig(t, clock)
	cfg.From = testNow.Add(-2 * time.Minute)
	src := newCheckSource(func(_ string, _ int) (Observation, error) {
		return healthyObservation(clock.Now()), nil
	})

	res, err := check(context.Background(), cfg, src)
	if err != nil {
		t.Fatalf("check() = %v, want a pass with a warning\nnotes:\n%s", err, notesOf(cfg))
	}
	notes := notesOf(cfg)
	if !strings.Contains(notes, "cannot see [") || !strings.Contains(notes, testNow.Format(time.RFC3339)) {
		t.Errorf("want a warning naming the unseen interval; notes were:\n%s", notes)
	}
	// The classified window is the clamped one, and Result says so rather than
	// reporting a window the run never proved.
	if !res.From.Equal(testNow) {
		t.Errorf("Result.From = %s, want the clamped %s", res.From, testNow)
	}
}

// §19.3 case 1: the failure limit was exceeded. The measurement pass succeeds
// and the collection loop then hits a terminal failure, so this exercises the
// path a live run really takes.
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
	if err == nil {
		t.Fatalf("check() = nil, want the collection failure to fail closed")
	}
	if !strings.Contains(err.Error(), "collect evidence") {
		t.Errorf("err = %q, want it to name the collection step", err)
	}
	if len(res.Violations) != 0 {
		t.Errorf("Violations = %+v; an error must never be reported as a verdict", res.Violations)
	}
}

// §19.3 case 2: the resolution of the definitions failed. Both shapes — the
// ruler read itself failing, and a name that resolves to nothing.
func TestCheckFailClosedOnDefinitionResolution(t *testing.T) {
	t.Run("ruler read fails", func(t *testing.T) {
		clock := newVirtualClock(testNow)
		cfg := baseConfig(t, clock)
		src := newCheckSource(nil)
		src.defsErr = errors.New("502 bad gateway")

		if _, err := check(context.Background(), cfg, src); err == nil ||
			!strings.Contains(err.Error(), "read rule definitions") {
			t.Fatalf("check() = %v, want a definitions-read failure", err)
		}
	})

	t.Run("unknown alert name", func(t *testing.T) {
		clock := newVirtualClock(testNow)
		cfg := baseConfig(t, clock)
		cfg.Alerts = []string{"No Such Rule"}
		src := newCheckSource(nil)

		if _, err := check(context.Background(), cfg, src); err == nil ||
			!strings.Contains(err.Error(), "no rule matched") {
			t.Fatalf("check() = %v, want a no-match failure", err)
		}
	})
}

// The version gate (§2.7 control 2): an unsupported Grafana is exit 2 before
// anything else is attempted.
func TestCheckRefusesUnsupportedGrafanaVersion(t *testing.T) {
	clock := newVirtualClock(testNow)
	cfg := baseConfig(t, clock)
	src := newCheckSource(nil)
	src.version = "12.4.0"

	if _, err := check(context.Background(), cfg, src); err == nil ||
		!strings.Contains(err.Error(), "unsupported grafana version") {
		t.Fatalf("check() = %v, want the version gate to refuse 12.4.0", err)
	}
}

// §5.2: the budget is checked against the latencies the measurement pass
// actually measured, and a schedule that cannot fit errors at START rather
// than producing a gap-riddled recording nobody can classify.
func TestCheckSingleStepRefusesAScheduleThatDoesNotFit(t *testing.T) {
	clock := newVirtualClock(testNow)
	cfg := baseConfig(t, clock)
	src := newCheckSource(func(_ string, _ int) (Observation, error) {
		obs := healthyObservation(clock.Now())
		obs.Latency = 45 * time.Second // longer than the rule's own 30s cadence
		return obs, nil
	})

	_, err := check(context.Background(), cfg, src)
	if err == nil {
		t.Fatalf("check() = nil, want the budget check to refuse the schedule")
	}
	for _, want := range []string{"raising concurrency", "raising poll-interval", "watching fewer alerts"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("err = %q, want it to name the control %q (§5.1)", err, want)
		}
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
	if err != nil {
		t.Fatalf("NewWriter: %v", err)
	}
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
	if err := w.WriteHeader(header); err != nil {
		t.Fatalf("WriteHeader: %v", err)
	}
	for at := start; !at.After(end); at = at.Add(checkPollEvery) {
		if err := w.WritePoll(Poll{
			RuleUID: checkUID, GrafanaNow: at, Found: true,
			State: "inactive", Health: "ok", LastEvaluation: at.Add(-lastEvalLag),
		}); err != nil {
			t.Fatalf("WritePoll: %v", err)
		}
	}
	if err := w.Stop(); err != nil {
		t.Fatalf("Stop: %v", err)
	}
	return path
}

// deadPid returns a pid that is guaranteed to have exited — the normal state
// of a recorder by the time check signals it, since a recorder given --until
// (or one that finished cleanly) is already gone.
func deadPid(t *testing.T) int {
	t.Helper()
	cmd := exec.Command("/bin/sh", "-c", "exit 0")
	if err := cmd.Start(); err != nil {
		t.Fatalf("start a throwaway process: %v", err)
	}
	pid := cmd.Process.Pid
	if err := cmd.Wait(); err != nil {
		t.Fatalf("wait for the throwaway process: %v", err)
	}
	return pid
}

func writePid(t *testing.T, path, contents string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatalf("write pidfile: %v", err)
	}
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
		t.Errorf("the drain wait polled %q although the log already proves the evaluations", title)
		return Observation{}, errors.New("unexpected poll")
	})

	res, err := check(context.Background(), cfg, src)
	if err != nil {
		t.Fatalf("check() = %v, want nil\nnotes:\n%s", err, notesOf(cfg))
	}
	if len(res.Violations) != 0 {
		t.Fatalf("Violations = %+v, want none", res.Violations)
	}
	if len(res.Verdicts) != 1 || res.Verdicts[0].Outcome != OutcomeClean {
		t.Fatalf("Verdicts = %+v, want one clean verdict", res.Verdicts)
	}
	if res.GrafanaVersion != "13.1.0" {
		t.Errorf("GrafanaVersion = %q, want the recorded one", res.GrafanaVersion)
	}
	// The collection loop still waited out to+transitionGrace (H5) even though
	// the recorder had already finished.
	if clock.Now().Before(windowEnd) {
		t.Errorf("returned at %s, before to+grace %s", clock.Now(), windowEnd)
	}
}

// §19.3 case 3: the identity of the log is not correct. The check runs against
// the header read EARLY, so it fails before the window's wait rather than
// after it.
func TestCheckFailClosedOnWrongLogIdentity(t *testing.T) {
	t.Run("different url", func(t *testing.T) {
		dir := t.TempDir()
		windowEnd := testNow.Add(5*time.Minute + checkGrace)
		logPath := recordedLog(t, dir, "https://other.example.com",
			testNow.Add(-time.Minute), testNow.Add(-time.Minute), windowEnd, windowEnd.Add(30*time.Second), 0)

		clock := newVirtualClock(testNow)
		cfg := recorderConfig(t, clock, logPath)
		_, err := check(context.Background(), cfg, newCheckSource(nil))
		if err == nil || !strings.Contains(err.Error(), "log identity") {
			t.Fatalf("check() = %v, want a log-identity failure", err)
		}
		if !clock.Now().Equal(testNow) {
			t.Errorf("the identity check waited out the window (now %s); it must fail before the wait", clock.Now())
		}
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
		if err == nil || !strings.Contains(err.Error(), "log identity") {
			t.Fatalf("check() = %v, want a log-identity failure", err)
		}
	})
}

// §19.3 case 4: the coverage proof failed. A hole in the middle of the
// recording is not saved by healthy data at both ends (§22.4).
func TestCheckFailClosedOnCoverageGap(t *testing.T) {
	dir := t.TempDir()
	windowEnd := testNow.Add(5*time.Minute + checkGrace)
	path := filepath.Join(dir, "log.jsonl")
	clock := newFakeClock(windowEnd.Add(30 * time.Second))
	w, err := NewWriter(path, clock)
	if err != nil {
		t.Fatalf("NewWriter: %v", err)
	}
	if err := w.WriteHeader(Header{
		URL: "https://grafana.example.com", GrafanaVersion: "13.1.0", StartedAt: testNow.Add(-time.Minute),
		Rules: []LoggedRule{{UID: checkUID, Title: checkTitle, IntervalSeconds: 60, PollEverySeconds: checkPollEvery.Seconds()}},
	}); err != nil {
		t.Fatalf("WriteHeader: %v", err)
	}
	for at := testNow.Add(-time.Minute); !at.After(windowEnd.Add(30 * time.Second)); at = at.Add(checkPollEvery) {
		// A three-minute hole in the middle of the window.
		if at.After(testNow.Add(time.Minute)) && at.Before(testNow.Add(4*time.Minute)) {
			continue
		}
		if err := w.WritePoll(Poll{
			RuleUID: checkUID, GrafanaNow: at, Found: true, State: "inactive", Health: "ok", LastEvaluation: at,
		}); err != nil {
			t.Fatalf("WritePoll: %v", err)
		}
	}
	if err := w.Stop(); err != nil {
		t.Fatalf("Stop: %v", err)
	}
	writePid(t, path+".pid", fmt.Sprintf("%d\n", deadPid(t)))

	cfg := recorderConfig(t, newVirtualClock(testNow), path)
	res, err := check(context.Background(), cfg, newCheckSource(nil))
	if err == nil {
		t.Fatalf("check() = nil, want the coverage gap to fail closed")
	}
	if got := res.Coverage[checkUID].Reason; got != ReasonHeartbeatGap {
		t.Errorf("Reason = %q, want %q", got, ReasonHeartbeatGap)
	}
	if got := res.Verdicts[0].Outcome; got != OutcomeUnobservable {
		t.Errorf("Outcome = %q, want %q", got, OutcomeUnobservable)
	}
}

// §19.3 case 5: the drain limit passed. The recording itself is clean, so this
// isolates the drain wait — the rule simply never evaluates through the end of
// the window, and a rule that cannot answer that question is unobservable.
func TestCheckFailClosedOnDrainTimeout(t *testing.T) {
	dir := t.TempDir()
	windowEnd := testNow.Add(5*time.Minute + checkGrace)
	// A 45s lag keeps every poll inside evalStaleAfter (120s), so P7 check 6
	// is silent and only the drain wait can fail.
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
	if err == nil {
		t.Fatalf("check() = nil, want the drain limit to fail closed")
	}
	if got := res.Coverage[checkUID].Reason; got != ReasonDrainTimeout {
		t.Errorf("Reason = %q, want %q", got, ReasonDrainTimeout)
	}
	if got := res.Verdicts[0].Outcome; got != OutcomeUnobservable {
		t.Errorf("Outcome = %q, want %q", got, OutcomeUnobservable)
	}
	if !strings.Contains(res.Verdicts[0].Note, "drain limit") {
		t.Errorf("Note = %q, want it to explain the drain limit", res.Verdicts[0].Note)
	}
	if waited := clock.Now().Sub(windowEnd); waited < checkDrainLimit {
		t.Errorf("gave up after %s of drain wait, want the full %s", waited, checkDrainLimit)
	}
}

// §14.5: a rule the state endpoint no longer serves is knowable on the FIRST
// drain poll, and the answer is rule_absent — the fault — rather than
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
	// An authoritative 2xx that parsed and carries no matching rule. P2
	// retried every transport failure long before an Observation exists, so
	// this is a deletion, not a hiccup.
	src := newCheckSource(func(_ string, _ int) (Observation, error) {
		return Observation{GrafanaNow: clock.Now()}, nil
	})

	res, err := check(context.Background(), cfg, src)
	if err == nil {
		t.Fatalf("check() = nil, want a deleted rule to fail closed")
	}
	if got := res.Coverage[checkUID].Reason; got != ReasonRuleAbsent {
		t.Errorf("Reason = %q, want %q — the fault, not the wait", got, ReasonRuleAbsent)
	}
	if got := src.callCount(checkTitle); got != 1 {
		t.Errorf("polled %d times, want exactly 1: the absence is knowable on the first poll", got)
	}
	if waited := clock.Now().Sub(windowEnd); waited >= checkDrainLimit {
		t.Errorf("spent %s in the drain wait, want it to conclude at once", waited)
	}
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
	if err != nil {
		t.Fatalf("NewWriter: %v", err)
	}
	if err := w.WriteHeader(Header{
		URL: "https://grafana.example.com", GrafanaVersion: "13.1.0", StartedAt: testNow.Add(-time.Minute),
		Rules: []LoggedRule{{
			UID: checkUID, Title: checkTitle, IntervalSeconds: 60,
			IsPaused: false, PollEverySeconds: checkPollEvery.Seconds(),
		}},
	}); err != nil {
		t.Fatalf("WriteHeader: %v", err)
	}
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
		if err := w.WritePoll(p); err != nil {
			t.Fatalf("WritePoll: %v", err)
		}
	}
	if err := w.Stop(); err != nil {
		t.Fatalf("Stop: %v", err)
	}
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
	res, err, cfg := pausedAfterWindowCheck(t, false)
	if err != nil {
		t.Fatalf("check() = %v, want a classified verdict\nnotes:\n%s", err, notesOf(cfg))
	}
	if got := res.Verdicts[0].Outcome; got != OutcomeNewlyBad {
		t.Fatalf("Outcome = %q, want %q: the rule was active for the whole window and fired inside it", got, OutcomeNewlyBad)
	}
	if len(res.Violations) != 1 || res.Violations[0].Outcome != OutcomeNewlyBad {
		t.Fatalf("Violations = %+v, want the firing reported", res.Violations)
	}
	if strings.Contains(res.Verdicts[0].Note, "paused before the window opened") {
		t.Errorf("Note = %q, which the log's own polls contradict", res.Verdicts[0].Note)
	}
}

// The regression pin for the loophole this fix closed. Reading skipped from
// the post-window definition made the rule skipped; --allow-paused then made
// skipped free; and a window in which the alert fired reported exit 0. The
// default message names --allow-paused, so an operator was led straight to it.
func TestCheckAllowPausedCannotExcuseARulePausedAfterItFired(t *testing.T) {
	res, err, cfg := pausedAfterWindowCheck(t, true)
	if err != nil {
		t.Fatalf("check() = %v, want a classified verdict\nnotes:\n%s", err, notesOf(cfg))
	}
	if len(res.Violations) == 0 {
		t.Fatalf("Violations = none with --allow-paused: the run passed over a window in which the alert fired")
	}
}

// The other direction, unchanged: a rule the HEADER says was paused when the
// recording opened is genuinely skipped. It has no polls, so no coverage is
// attempted for it, and --allow-paused behaves as it always did.
func TestCheckHeaderPausedRuleStaysSkipped(t *testing.T) {
	dir := t.TempDir()
	windowEnd := testNow.Add(5*time.Minute + checkGrace)
	path := filepath.Join(dir, "log.jsonl")
	w, err := NewWriter(path, newFakeClock(windowEnd.Add(30*time.Second)))
	if err != nil {
		t.Fatalf("NewWriter: %v", err)
	}
	// Named in the header, is_paused true, and no poll records at all — the
	// shape watch writes for a rule paused before the window opened (P6
	// deviation 4).
	if err := w.WriteHeader(Header{
		URL: "https://grafana.example.com", GrafanaVersion: "13.1.0", StartedAt: testNow.Add(-time.Minute),
		Rules: []LoggedRule{{
			UID: checkUID, Title: checkTitle, IntervalSeconds: 60,
			IsPaused: true, PollEverySeconds: checkPollEvery.Seconds(),
		}},
	}); err != nil {
		t.Fatalf("WriteHeader: %v", err)
	}
	if err := w.Stop(); err != nil {
		t.Fatalf("Stop: %v", err)
	}
	writePid(t, path+".pid", fmt.Sprintf("%d\n", deadPid(t)))

	run := func(allowPaused bool) (Result, error) {
		cfg := recorderConfig(t, newVirtualClock(testNow), path)
		cfg.AllowPaused = allowPaused
		// The definition is unpaused now; the header still decides.
		src := newCheckSource(func(title string, _ int) (Observation, error) {
			t.Errorf("the drain wait polled skipped rule %q", title)
			return Observation{}, errors.New("unexpected poll")
		})
		return check(context.Background(), cfg, src)
	}

	res, err := run(false)
	if err != nil {
		t.Fatalf("check() = %v, want exit-1 shape: a skipped rule is a known condition, not an inability", err)
	}
	if got := res.Verdicts[0].Outcome; got != OutcomeSkipped {
		t.Fatalf("Outcome = %q, want %q", got, OutcomeSkipped)
	}
	if _, ok := res.Coverage[checkUID]; ok {
		t.Errorf("Coverage[%s] present, want absent: a skipped rule has no coverage to prove", checkUID)
	}
	if len(res.Violations) != 1 {
		t.Errorf("Violations = %+v, want the MinObserved shortfall (§12.1)", res.Violations)
	}

	res, err = run(true)
	if err != nil || len(res.Violations) != 0 {
		t.Errorf("with --allow-paused: err = %v, Violations = %+v, want a pass", err, res.Violations)
	}
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
	if err == nil {
		t.Fatalf("check() = nil, want a rule that stopped evaluating to fail closed")
	}
	if got := src.callCount(checkTitle); got != 1 {
		t.Errorf("polled %d times, want exactly 1: a paused rule can never catch up", got)
	}
	if got := res.Coverage[checkUID].Reason; got != ReasonDrainTimeout {
		t.Errorf("Reason = %q, want %q — the vocabulary is published (§19.0), so the detail goes in the note", got, ReasonDrainTimeout)
	}
	if !strings.Contains(res.Verdicts[0].Note, "paused before it evaluated through") {
		t.Errorf("Note = %q, want it to say the rule was paused", res.Verdicts[0].Note)
	}
	if waited := clock.Now().Sub(windowEnd); waited >= checkDrainLimit {
		t.Errorf("spent %s in the drain wait, want it to conclude at once", waited)
	}
}

// P6's obligation on this phase: an absent or unparseable pidfile is never
// "there was nothing to stop". The parent writes the pidfile only once the
// child reports that it is recording, so a missing one means the recording
// never started — and the log must not be read at all.
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
			if err == nil || !strings.Contains(err.Error(), "cannot stop the recorder") {
				t.Fatalf("check() = %v, want a refusal to stop the recorder", err)
			}
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
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	if err := cmd.Start(); err != nil {
		t.Fatalf("start the lock holder: %v", err)
	}
	t.Cleanup(func() {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
	})
	if _, err := bufio.NewReader(stdout).ReadString('\n'); err != nil {
		t.Fatalf("the lock holder never reported holding the lock: %v", err)
	}
	return cmd.Process.Pid
}

// §4.4 step 4: a recorder that will not let go of the log means the log may
// still be appended to, and a log a writer can change cannot be read at all.
func TestCheckFailsWhenTheRecorderWillNotExit(t *testing.T) {
	dir := t.TempDir()
	windowEnd := testNow.Add(5*time.Minute + checkGrace)
	logPath := recordedLog(t, dir, "https://grafana.example.com",
		testNow.Add(-time.Minute), testNow.Add(-time.Minute), windowEnd, windowEnd.Add(30*time.Second), 0)

	writePid(t, logPath+".pid", fmt.Sprintf("%d\n", startLockHolder(t, logPath)))

	cfg := recorderConfig(t, newVirtualClock(testNow), logPath)
	_, err := check(context.Background(), cfg, newCheckSource(nil))
	if err == nil || !strings.Contains(err.Error(), "still holds") {
		t.Fatalf("check() = %v, want the stop wait to time out on the lock", err)
	}
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
	if err := bystander.Start(); err != nil {
		t.Fatalf("start the bystander: %v", err)
	}
	t.Cleanup(func() {
		_ = bystander.Process.Kill()
		_ = bystander.Wait()
	})
	writePid(t, logPath+".pid", fmt.Sprintf("%d\n", bystander.Process.Pid))

	cfg := recorderConfig(t, newVirtualClock(testNow), logPath)
	if _, err := check(context.Background(), cfg, newCheckSource(nil)); err != nil {
		t.Fatalf("check() = %v, want nil\nnotes:\n%s", err, notesOf(cfg))
	}
	if err := syscall.Kill(bystander.Process.Pid, 0); err != nil {
		t.Fatalf("the bystander is gone (%v): check signalled a process that was not the recorder", err)
	}
}

// P5's "two authorities", from check's side: maxGap comes from the cadence the
// header records, never from a re-derivation off intervalSeconds. The
// fail-open direction is the one asserted — a log recorded at 5s on a 60s rule
// must still fail on a hole a re-derived 30s maxGap would have forgiven.
func TestCheckDerivesMaxGapFromTheRecordedCadence(t *testing.T) {
	dir := t.TempDir()
	windowEnd := testNow.Add(5*time.Minute + checkGrace)
	path := filepath.Join(dir, "log.jsonl")
	w, err := NewWriter(path, newFakeClock(windowEnd.Add(30*time.Second)))
	if err != nil {
		t.Fatalf("NewWriter: %v", err)
	}
	if err := w.WriteHeader(Header{
		URL: "https://grafana.example.com", GrafanaVersion: "13.1.0", StartedAt: testNow.Add(-time.Minute),
		Rules: []LoggedRule{{UID: checkUID, Title: checkTitle, IntervalSeconds: 60, PollEverySeconds: 5}},
	}); err != nil {
		t.Fatalf("WriteHeader: %v", err)
	}
	for at := testNow.Add(-time.Minute); !at.After(windowEnd.Add(30 * time.Second)); at = at.Add(5 * time.Second) {
		// A 20s hole: under the recorded 5s cadence maxGap is 10s and this
		// fails; under a cadence re-derived from intervalSeconds it would be
		// 60s and the hole would pass unseen.
		if at.After(testNow.Add(time.Minute)) && at.Before(testNow.Add(80*time.Second)) {
			continue
		}
		if err := w.WritePoll(Poll{
			RuleUID: checkUID, GrafanaNow: at, Found: true, State: "inactive", Health: "ok", LastEvaluation: at,
		}); err != nil {
			t.Fatalf("WritePoll: %v", err)
		}
	}
	if err := w.Stop(); err != nil {
		t.Fatalf("Stop: %v", err)
	}
	writePid(t, path+".pid", fmt.Sprintf("%d\n", deadPid(t)))

	cfg := recorderConfig(t, newVirtualClock(testNow), path)
	res, err := check(context.Background(), cfg, newCheckSource(nil))
	if err == nil {
		t.Fatalf("check() = nil; a 20s hole exceeds the 10s maxGap the recorded 5s cadence implies")
	}
	if got := res.Coverage[checkUID].Reason; got != ReasonHeartbeatGap {
		t.Errorf("Reason = %q, want %q", got, ReasonHeartbeatGap)
	}
}

// ---------------------------------------------------------------------------
// The pieces, in isolation
// ---------------------------------------------------------------------------

// The drain wait's one comparison is cross-domain (§16), and its uncertainty
// is spent in the fail-closed direction: an evaluation that only MIGHT have
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
			if got := evaluatedThrough(tc.lastEval, tc.skew, tc.bound, end); got != tc.wantSatisfied {
				t.Errorf("evaluatedThrough() = %v, want %v", got, tc.wantSatisfied)
			}
		})
	}
}

// H6 through the merge: a drain timeout on one rule and a coverage failure on
// another must both reach the message. Neither error may shadow the other.
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
	if err == nil {
		t.Fatal("mergeDrainTimeouts() = nil, want an error naming the newly unobservable rule")
	}
	// Its own shape: joined with decide's, two counts under one identical
	// phrase would read as a contradiction rather than as two findings.
	if !strings.Contains(err.Error(), "unobservable at the drain wait") {
		t.Errorf("err = %q, want the drain wait's own error shape", err)
	}
	// Only A is newly unobservable; B was already, so naming it twice would
	// only lengthen the message.
	if !strings.Contains(err.Error(), "A ("+string(ReasonDrainTimeout)+")") {
		t.Errorf("err = %q, want it to name A's drain timeout", err)
	}
	if strings.Contains(err.Error(), "B (") {
		t.Errorf("err = %q, want it not to re-report B, which decide already reported", err)
	}
	if got := merged.Coverage["a"].Reason; got != ReasonDrainTimeout {
		t.Errorf("Coverage[a].Reason = %q, want %q", got, ReasonDrainTimeout)
	}
	// B keeps the reason the coverage proof gave it — the FIRST reason wins,
	// as it does inside proveCoverage.
	if got := merged.Coverage["b"].Reason; got != ReasonHeartbeatGap {
		t.Errorf("Coverage[b].Reason = %q, want the earlier %q", got, ReasonHeartbeatGap)
	}
	if merged.Verdicts[0].Outcome != OutcomeUnobservable {
		t.Errorf("Verdicts[0].Outcome = %q, want %q", merged.Verdicts[0].Outcome, OutcomeUnobservable)
	}
}

// ReadLogHeader is the one read of a log a writer may still hold, so its
// refusals matter as much as its successes.
func TestReadLogHeader(t *testing.T) {
	dir := t.TempDir()

	t.Run("reads line 1 while the log keeps growing", func(t *testing.T) {
		path := filepath.Join(dir, "growing.jsonl")
		w, err := NewWriter(path, newFakeClock(testNow))
		if err != nil {
			t.Fatalf("NewWriter: %v", err)
		}
		defer w.Close()
		if err := w.WriteHeader(testHeader()); err != nil {
			t.Fatalf("WriteHeader: %v", err)
		}
		if err := w.WritePoll(Poll{RuleUID: "rule1", GrafanaNow: testNow, Found: true}); err != nil {
			t.Fatalf("WritePoll: %v", err)
		}

		h, err := ReadLogHeader(path)
		if err != nil {
			t.Fatalf("ReadLogHeader: %v", err)
		}
		if h.URL != testHeader().URL || len(h.Rules) != 1 {
			t.Errorf("header = %+v, want the written one", h)
		}
	})

	t.Run("a half-written header is not a header", func(t *testing.T) {
		path := filepath.Join(dir, "torn.jsonl")
		if err := os.WriteFile(path, []byte(`{"type":"header","url":"htt`), 0o644); err != nil {
			t.Fatalf("write: %v", err)
		}
		if _, err := ReadLogHeader(path); err == nil ||
			!strings.Contains(err.Error(), "no complete header") {
			t.Fatalf("ReadLogHeader() = %v, want a refusal", err)
		}
	})

	t.Run("a wrong schema version is refused", func(t *testing.T) {
		path := filepath.Join(dir, "old.jsonl")
		if err := os.WriteFile(path, []byte(`{"type":"header","schema_version":99,"url":"u"}`+"\n"), 0o644); err != nil {
			t.Fatalf("write: %v", err)
		}
		if _, err := ReadLogHeader(path); err == nil ||
			!strings.Contains(err.Error(), "schema version 99") {
			t.Fatalf("ReadLogHeader() = %v, want a schema refusal", err)
		}
	})
}
