package gate

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"time"
)

// Check returns (Result, error) and no exit code: the code is a presentation
// decision the CLI makes. err != nil is exit 2 unconditionally, even alongside
// real violations; violations with err == nil is exit 1; neither is exit 0.
// Check never reads the environment — the CLI reads URL/token and passes them
// in, and the token must never reach a *flag.FlagSet.

// countdownEvery is how often the collection loop reports what it is waiting
// for; a silent wait is indistinguishable from a hung process.
const countdownEvery = 30 * time.Second

// recorderStopTimeout bounds the wait for the recorder's exit after SIGTERM.
// Everything after the signal is local (finish the in-flight write, sentinel,
// fsync), so this is loose; it stays a hard error because a log a writer still
// holds cannot be read.
const recorderStopTimeout = 30 * time.Second

// recorderStopPoll is how often the wait re-checks the lock. With no wait(2)
// on a detached session leader, its exit is observable only by polling.
const recorderStopPoll = 100 * time.Millisecond

// Config is check's whole input. It is the CLI's view of a run, and it is
// deliberately wider than Policy: Policy is the narrowed, pure-layer subset
// that reaches decide (classify.go), and the token is the field that must
// never cross that line.
type Config struct {
	// URL and Token are the connection details, read from the environment by
	// the CLI and never registered as flags. Token never enters the pure layer,
	// an error string, or a Result.
	URL, Token string

	// Alerts is REQUIRED in single-step mode and must be EMPTY in log mode:
	// with a log, the header IS the alert set, and there is nothing to compare
	// a second list against.
	Alerts []string
	Folder string

	States               []State
	Preexisting          PreexistingPolicy
	MinObserved          int
	AllowPaused          bool
	NodataIsUnobservable bool

	// From is the moment the deploy finished and To is the end of the work.
	// They are different moments and both come from the work. In recorder mode
	// an absent From is a hard error; in single-step mode it falls back to the
	// start of this step, with a blind-interval warning.
	From, To time.Time

	// Log is the path of a recording made by watch; "" selects single-step
	// mode. PidFile defaults to <Log>.pid, the convention watch's parent
	// writes and the only way check can reach the recorder it must stop before
	// it may read the log.
	Log     string
	PidFile string

	// There is deliberately NO PollEvery here, and `check` has no
	// --poll-interval flag. In log mode the cadence comes from the header —
	// the cadence the recording actually used — and a second authority would
	// let an operator silently widen maxGap over evidence that was recorded at
	// a different rate; in single-step mode the same process records and
	// classifies, so the default cadence is the only cadence there is.
	Concurrency int
	Clock       Clock

	// Notes is where the shell prints what an operator has to see while the
	// run is in progress: the planned run time, the grace and its source, the
	// countdown, the blind-interval warning. nil discards them. The library
	// renders no table — the CLI owns presentation.
	Notes io.Writer
}

func (cfg Config) withDefaults() Config {
	if cfg.Clock == nil {
		cfg.Clock = SystemClock{}
	}
	if cfg.Notes == nil {
		cfg.Notes = io.Discard
	}
	if cfg.Concurrency < 1 {
		cfg.Concurrency = 1
	}
	if cfg.PidFile == "" && cfg.Log != "" {
		cfg.PidFile = cfg.Log + ".pid"
	}
	return cfg
}

// namedAlerts returns the alert names that survive Resolve's trim-and-discard,
// so validation counts what Resolve will actually see rather than what the
// caller happened to pass (a file ending in a newline yields an empty line).
func (cfg Config) namedAlerts() []string {
	out := make([]string, 0, len(cfg.Alerts))
	for _, a := range cfg.Alerts {
		if strings.TrimSpace(a) != "" {
			out = append(out, a)
		}
	}
	return out
}

// Check is the I/O shell: HTTP, signals, the pidfile, file reads, the
// countdown print. Every correctness question it touches is answered elsewhere
// (proveCoverage, decide — both pure), which is the most important seam in the
// project. A pass is exactly len(Violations) == 0 && err == nil; every error
// path leaves err non-nil.
func Check(ctx context.Context, cfg Config) (Result, error) {
	cfg = cfg.withDefaults()
	if err := cfg.validate(); err != nil {
		return Result{}, err
	}
	// The Source is built here and injected into check() so every behaviour
	// below is testable against a scripted fake — the same seam prepareWatch
	// uses, and the reason this file needs no test-only setter.
	return check(ctx, cfg, NewHTTPSource(cfg.URL, cfg.Token, cfg.Clock))
}

// validate runs before any network call, so a configuration mistake costs
// nothing and, more importantly, is never discovered after a ten-minute wait.
func (cfg Config) validate() error {
	if cfg.URL == "" {
		return errors.New("check: no grafana url")
	}
	if cfg.To.IsZero() {
		return errors.New("check: no `to`: the end of the window is required")
	}

	named := cfg.namedAlerts()
	if cfg.Log == "" {
		// An empty Alerts is an error — but only without a log.
		if len(named) == 0 {
			return errors.New("check: no alert names given and no recorded log to take them from")
		}
	} else if len(named) > 0 {
		// The other direction: with a log, the alert set comes from the log.
		// Accepting both would mean reconciling two sets, which the log being
		// the one source removes entirely.
		return fmt.Errorf("check: --alerts is refused with a recorded log: %s already names the alert set it recorded", cfg.Log)
	}

	now := cfg.Clock.Now()
	// from mirrors what check() will use, so the window checks below judge the
	// window that will really be classified.
	from := cfg.From
	switch {
	case from.IsZero() && cfg.Log != "":
		// Never a warning-and-continue: falling back to the start of the check
		// step reinstates exactly the blind interval the recorder exists to
		// remove, which is the fail-open shape this design refuses.
		return errors.New("check: no `from` in recorder mode: the deploy step must emit a completion timestamp")
	case from.IsZero():
		// Single-step only. The caller sees the resulting blind interval named
		// exactly, once the first observation has fixed its end.
		from = now
	}

	if cfg.To.Before(from) {
		return fmt.Errorf("check: `to` %s is before `from` %s", cfg.To.Format(time.RFC3339), from.Format(time.RFC3339))
	}
	if from.After(now.Add(fromFutureTolerance)) {
		return fmt.Errorf("check: `from` %s is more than %s ahead of this runner's clock %s",
			from.Format(time.RFC3339), fromFutureTolerance, now.Format(time.RFC3339))
	}

	// A `to` in the past is fine WITH a log (the collection loop is already
	// done). Without one it is a request to prove a window nothing observed:
	// every heartbeat gap would measure negative, and the run would report a
	// proved window it never saw.
	if cfg.Log == "" && !cfg.To.After(now) {
		return fmt.Errorf("check: `to` %s has already passed and there is no recorded log: a window that ended before check started can only be classified from a recording",
			cfg.To.Format(time.RFC3339))
	}
	return nil
}

// check is Check with the Source injected, and its body is one commented block
// per stage of a run, in the order a run performs them.
func check(ctx context.Context, cfg Config, src Source) (Result, error) {
	// ---- Validate the configuration. --------------------------------------
	// Done by Check before this function is reached, except for the one part
	// that needs a clock reading kept for later: the single-step fallback for
	// an absent `from`.
	from := cfg.From
	if from.IsZero() {
		from = cfg.Clock.Now()
		fmt.Fprintf(cfg.Notes, "note: no `from` given; the window starts at the start of this step, %s\n",
			from.Format(time.RFC3339))
	}

	// ---- Resolve the definitions from the ruler API. ----------------------
	// Unconditional, in BOTH modes. A log's header supplies the alert set as
	// UIDs and the recording facts, never the rule facts: `for`,
	// intervalSeconds and Kind always come from a fresh ruler read, which is
	// why LoggedRule.ForSeconds is never converted back into a Definition.
	version, err := src.Version(ctx)
	if err != nil {
		return Result{}, fmt.Errorf("read grafana version: %w", err)
	}
	if err := CheckGrafanaVersion(version); err != nil {
		return Result{}, err
	}
	allDefs, err := src.Definitions(ctx)
	if err != nil {
		return Result{}, fmt.Errorf("read rule definitions: %w", err)
	}

	// ---- With a log, validate its identity. -------------------------------
	// The header is read early — line 1 only, the one line a writer can never
	// change — so a wrong URL or an unresolvable rule fails closed NOW. It is
	// advisory: the authoritative header is re-read once collection ends and
	// the writer has exited.
	var (
		resolved   []Definition
		notes      []string
		earlyHdr   Header
		logHasHdr  bool
		rt         map[string]ruleTimings
		gt         globalTimings
		timingNote []string
	)
	if cfg.Log != "" {
		earlyHdr, err = ReadLogHeader(cfg.Log)
		if err != nil {
			return Result{}, fmt.Errorf("log identity: %w", err)
		}
		logHasHdr = true
		// Fail fast on a statically-knowable bound violation. `from <
		// StartedAt` makes coverage unprovable no matter how healthy the polls
		// that DO exist look, and StartedAt is immutable (line 1, written
		// first), so this cannot disagree with the authoritative header read
		// later. proveCoverage's check 2 remains the backstop against the
		// authoritative header, so a bad advisory read can only ever fail
		// closed, never produce a false pass. This is recorder mode only: the
		// single-step branch has no header, and its own `from < startedAt` is
		// a warning-and-pass (see below), not an error.
		if from.Before(earlyHdr.StartedAt) {
			return Result{}, fmt.Errorf("check: `from` %s is before recording started at %s",
				from.Format(time.RFC3339), earlyHdr.StartedAt.Format(time.RFC3339))
		}
		resolved, notes, err = resolveFromLog(allDefs, earlyHdr, cfg)
	} else {
		resolved, notes, err = Resolve(allDefs, cfg.namedAlerts(), cfg.Folder)
	}
	if err != nil {
		return Result{}, err
	}
	if len(resolved) == 0 {
		// Reachable only from a header with an empty rule list. Left to run,
		// MinObserved would default to zero, no rule would be judged, and the
		// gate would return a pass over nothing at all.
		return Result{}, fmt.Errorf("check: no rules to classify")
	}
	for _, n := range notes {
		fmt.Fprintf(cfg.Notes, "note: %s\n", n)
	}

	// ---- Derive the timings, print them, fit the request budget. ----------
	if logHasHdr {
		// The header is the authority for the cadence actually recorded at;
		// re-deriving it from defs would compare gaps recorded at an override
		// cadence against thresholds computed from the default — fail-open in
		// the faster-override direction.
		rt, gt, err = DeriveTimingsFromLog(earlyHdr, resolved)
		if err != nil {
			return Result{}, fmt.Errorf("log identity: %w", err)
		}
	} else {
		rt, gt, timingNote = DeriveTimings(resolved, 0)
		for _, n := range timingNote {
			fmt.Fprintf(cfg.Notes, "note: %s\n", n)
		}
	}
	summary, warning := StartupSummary(from, cfg.To, gt)
	fmt.Fprintln(cfg.Notes, summary)
	if warning != "" {
		fmt.Fprintf(cfg.Notes, "warning: %s\n", warning)
	}

	// The measurement pass and budget check are single-step only: in recorder
	// mode watch already measured and checked the budget before detaching.
	var (
		header  Header
		initial []Poll
		reducer = NewReducer()
	)
	if !logHasHdr {
		// StartedAt is fixed before the pass rather than after it, so the
		// interval it claims to have observed can only be wider than the one
		// it really saw — and the first heartbeat's own boundary gap is what
		// proves that interval, not this timestamp.
		startedAt := cfg.Clock.Now()
		active := activeRules(resolved)
		var measured map[string]time.Duration
		initial, measured, err = firstObservations(ctx, src, active, reducer, cfg.Concurrency, cfg.Notes)
		if err != nil {
			return Result{}, err
		}
		if err := CheckBudget(activeTimingsOf(active, rt), measured, cfg.Concurrency); err != nil {
			return Result{}, err
		}

		// Single-step synthesis — how the pure layer stays unconditional. The
		// shell builds the Header and later stamps the sentinel itself, so the
		// sentinel and from-bounds coverage checks run exactly as they do over
		// a recording and no mode flag ever reaches proveCoverage or decide.
		header = Header{
			SchemaVersion:  LogSchemaVersion,
			URL:            cfg.URL,
			GrafanaVersion: version,
			StartedAt:      startedAt,
			Rules:          loggedRules(resolved, rt),
		}
		if from.Before(startedAt) {
			// The declared blind interval: in single-step mode this is a
			// warning and a pass, and ONLY here. Recorder mode keeps the
			// from-bounds coverage check strict, because there the recorder
			// was supposed to be watching and the gap means it was not.
			fmt.Fprintf(cfg.Notes, "warning: cannot see [%s, %s) — %s before the first observation; the window is classified from %s\n",
				from.Format(time.RFC3339), startedAt.Format(time.RFC3339),
				startedAt.Sub(from).Round(time.Second), startedAt.Format(time.RFC3339))
			from = startedAt
		}
	}

	// ---- Apply MinObserved. -----------------------------------------------
	// Its default is the resolved rule count AFTER duplicate names collapse,
	// which is len(resolved) by construction. decide defaults it identically;
	// it is resolved here as well so the value the run will judge against is
	// printed before the wait rather than inferred from the verdict afterwards.
	minObserved := cfg.MinObserved
	if minObserved == 0 {
		minObserved = len(resolved)
	}
	fmt.Fprintf(cfg.Notes, "min-observed: %d of %d resolved rule(s)\n", minObserved, len(resolved))

	// ---- Collect the evidence. --------------------------------------------
	// Collect ONLY. No classification happens here and there is no early exit,
	// even once a violation is certain: the loop always runs to
	// to + transitionGrace, which is what makes "did the early exit lose the
	// coverage proof?" a question that cannot be asked.
	windowEnd := cfg.To.Add(gt.transitionGrace)

	var poller *livePoller
	if !logHasHdr {
		poller = newLivePoller(src, reducer, activeRules(resolved), rt, cfg.Concurrency, cfg.Clock.Now())
	}
	collected, err := collectUntil(ctx, cfg, windowEnd, poller)
	if err != nil {
		// Nothing collected is classified; the count lets an operator tell a
		// run that failed at once from one that failed at minute nine.
		return Result{}, fmt.Errorf("collect evidence after %d poll(s): %w", len(collected), err)
	}

	var (
		polls    []Poll
		sentinel *time.Time
	)
	if logHasHdr {
		// In this order and no other: signal the writer, wait for its exit,
		// and only THEN read the log once. A log read while a writer can still
		// append can only yield a shorter window than the one that was
		// actually recorded.
		heldLog, err := stopRecorder(ctx, cfg)
		if err != nil {
			return Result{}, err
		}
		header, polls, sentinel, err = ReadLog(cfg.Log)
		// Held across the read so no writer can appear mid-read, then released
		// (everything past here works from memory, and the drain wait is minutes).
		_ = heldLog.Close()
		if err != nil {
			return Result{}, err
		}
		// The authoritative header wins: the advisory read was only a fail-fast.
		resolved, _, err = resolveFromLog(allDefs, header, cfg)
		if err != nil {
			return Result{}, err
		}
		// rt is re-derived from the authoritative header. windowEnd is NOT
		// recomputed: the loop already stopped at the earlier value, and moving
		// it afterwards would prove a window this run did not collect.
		if rt, gt, err = DeriveTimingsFromLog(header, resolved); err != nil {
			return Result{}, fmt.Errorf("log identity: %w", err)
		}
	} else {
		polls = make([]Poll, 0, len(initial)+len(collected))
		polls = append(polls, initial...)
		polls = append(polls, collected...)
		// The shell stamps the sentinel itself, when the collection loop
		// exits: by construction that is at or after to + transitionGrace, so
		// the sentinel check passes for the same reason a clean recorder stop
		// does, and for no other.
		stoppedAt := cfg.Clock.Now()
		sentinel = &stoppedAt
	}

	// ---- The drain wait. --------------------------------------------------
	// The last instance of the liveness check: did this rule evaluate through
	// the end of the window? It is I/O and it is deliberately NOT part of
	// proveCoverage — adding it there would put HTTP inside the pure layer and
	// destroy the seam this design depends on.
	drained, err := drainWait(ctx, cfg, src, resolved, header.pausedAtStart(), rt, polls, windowEnd, gt.drainTimeout)
	if err != nil {
		return Result{}, err
	}

	// ---- Classify. --------------------------------------------------------
	pol := Policy{
		States:               cfg.States,
		Preexisting:          cfg.Preexisting,
		MinObserved:          minObserved,
		AllowPaused:          cfg.AllowPaused,
		NodataIsUnobservable: cfg.NodataIsUnobservable,
		From:                 from,
		To:                   cfg.To,
	}
	result, decideErr := decide(header, polls, sentinel, resolved, rt, gt, pol)
	result, drainErr := mergeDrainTimeouts(result, drained)

	// ---- Return the Result. -----------------------------------------------
	// Both errors are joined rather than one shadowing the other: each names
	// rules the other does not, and on exit 2 that list IS the answer to
	// "why".
	return result, errors.Join(decideErr, drainErr)
}

// resolveFromLog is the log's identity check in practice: the URL must match
// and every header UID must still resolve against a fresh ruler read. The
// alert set is TAKEN from the log, never compared against --alerts (which the
// validator requires empty in log mode). Resolving through Resolve by uid:
// keeps one implementation of the resolution rules.
//
// Only the header-to-defs direction can fail: resolved is BUILT from the
// header, so no resolved definition can be absent from it.
func resolveFromLog(allDefs []Definition, h Header, cfg Config) ([]Definition, []string, error) {
	if h.URL != cfg.URL {
		return nil, nil, fmt.Errorf("log identity: %s recorded url %q but this run is configured for %q",
			cfg.Log, h.URL, cfg.URL)
	}
	names := make([]string, 0, len(h.Rules))
	for _, lr := range h.Rules {
		names = append(names, "uid:"+lr.UID)
	}
	resolved, notes, err := Resolve(allDefs, names, "")
	if err != nil {
		return nil, nil, fmt.Errorf("log identity: %s names a rule that no longer resolves: %w", cfg.Log, err)
	}
	return resolved, notes, nil
}

// activeRules drops the rules whose DEFINITION says paused. They are skipped:
// never polled, never waited for, and reported from the definitions alone — a
// skipped rule has no poll records at all, so it has no heartbeats to prove
// and no IsPaused poll to detect.
func activeRules(defs []Definition) []Definition {
	out := make([]Definition, 0, len(defs))
	for _, d := range defs {
		if !d.IsPaused {
			out = append(out, d)
		}
	}
	return out
}

// activeTimingsOf narrows the timings map to the rules that will actually be
// polled, which is what the request budget is spent on: a skipped rule
// consumes none of the capacity, so counting it would refuse schedules that
// fit.
func activeTimingsOf(active []Definition, rt map[string]ruleTimings) map[string]ruleTimings {
	out := make(map[string]ruleTimings, len(active))
	for _, d := range active {
		out[d.UID] = rt[d.UID]
	}
	return out
}

// livePoller is single-step mode's collection engine: the same per-rule
// scheduler and the same Reducer the recorder uses, writing into memory
// instead of a log. Log mode has none — the recorder is doing this work in
// another process — and collectUntil takes a nil poller for it.
type livePoller struct {
	src         Source
	reducer     *Reducer
	sched       *Scheduler
	titles      map[string]string // uid -> title: poll by title, select by UID
	concurrency int
}

func newLivePoller(src Source, reducer *Reducer, active []Definition, rt map[string]ruleTimings,
	concurrency int, now time.Time) *livePoller {

	titles := make(map[string]string, len(active))
	cadence := make(map[string]time.Duration, len(active))
	for _, d := range active {
		titles[d.UID] = d.Title
		cadence[d.UID] = rt[d.UID].pollEvery
	}
	return &livePoller{
		src:         src,
		reducer:     reducer,
		sched:       NewScheduler(cadence, now),
		titles:      titles,
		concurrency: concurrency,
	}
}

// poll runs one round of due rules and returns every poll that succeeded,
// alongside the first failure.
//
// The successes are NOT kept for the reason watchLoopConfig.pollBatch keeps
// its own: those go into a durable log that a later check will read, so
// dropping one would turn a single rule's transport failure into a coverage
// gap for the others. Here there is no later reader. A terminal failure during
// collection is exit 2 and check discards the whole collection, so these come
// back only to let the error say how far the run got before it stopped —
// which is the one part of it an operator can act on.
func (p *livePoller) poll(ctx context.Context, uids []string) ([]Poll, error) {
	observed, obsErr := observeAll(ctx, p.src, p.titles, uids, p.concurrency)
	out := make([]Poll, 0, len(uids))
	for _, uid := range uids {
		obs, ok := observed[uid]
		if !ok {
			continue
		}
		out = append(out, p.reducer.Reduce(uid, obs))
	}
	return out, obsErr
}

// collectUntil is the collection loop, shared by both modes. With a poller it
// polls each rule on its own cadence; with nil it only waits, because in
// recorder mode the evidence is being written by another process. Both print
// the same countdown, because both are the same silence to an operator
// watching a job.
//
// It never classifies and never exits early.
func collectUntil(ctx context.Context, cfg Config, deadline time.Time, p *livePoller) ([]Poll, error) {
	var (
		polls     []Poll
		lastPrint time.Time
	)
	for {
		now := cfg.Clock.Now()
		if !now.Before(deadline) {
			return polls, nil
		}
		if lastPrint.IsZero() || now.Sub(lastPrint) >= countdownEvery {
			fmt.Fprintf(cfg.Notes, "collecting: %s until the window closes at %s\n",
				deadline.Sub(now).Round(time.Second), deadline.Format(time.RFC3339))
			lastPrint = now
		}

		wait := min(deadline.Sub(now), countdownEvery)
		if p != nil {
			// Mark before polling, against the batch's own `now`: the next poll
			// is one cadence after this one was DUE, not after it returned, so
			// request latency cannot make the heartbeat spacing drift towards
			// maxGap (the same rule watchLoop follows).
			due := p.sched.Due(now)
			for _, uid := range due {
				p.sched.Mark(uid, now)
			}
			if len(due) > 0 {
				batch, err := p.poll(ctx, due)
				polls = append(polls, batch...)
				if err != nil {
					return polls, err
				}
			}
			if next, ok := p.sched.earliestDue(); ok {
				wait = min(wait, next.Sub(cfg.Clock.Now()))
			}
		}

		select {
		case <-ctx.Done():
			return polls, ctx.Err()
		case <-cfg.Clock.After(max(wait, 0)):
		}
	}
}

// stopRecorder signals the recorder and waits for it to go; the log may not be
// read until the writer has provably gone, so every failure is a hard error.
// It returns the log held under an exclusive flock, which the caller must keep
// open across ReadLog — the lock is the proof that no writer exists.
//
// Two authorities, only one of which is evidence:
//
//   - the PIDFILE says whether a recording ever started (it is written only
//     after the child reports ready, and removed on failure).
//   - the FLOCK says whether a writer exists right now. A pidfile can go stale
//     — nothing removes it on a clean --until stop, so it may name a pid
//     somebody else now owns — but the kernel drops a flock when the holder
//     exits, so the lock is always authoritative.
//
// So: read the pidfile to learn a recording happened, ask the lock whether it
// is still running, and signal only if it is.
func stopRecorder(ctx context.Context, cfg Config) (*os.File, error) {
	pid, err := ReadPidFile(cfg.PidFile)
	if err != nil {
		return nil, fmt.Errorf("cannot stop the recorder: %w; a pidfile is written only once a recorder reports that it is running, so an unreadable one means the recording never started", err)
	}

	log, err := os.Open(cfg.Log)
	if err != nil {
		return nil, fmt.Errorf("open %s to check for a writer: %w", cfg.Log, err)
	}

	held, err := tryLockExclusive(log)
	if err != nil {
		log.Close()
		return nil, err
	}
	if held {
		// No writer. Send no signal, whatever the pidfile says — the pid may
		// belong to somebody else entirely by now. Which of --until, a clean
		// stop and a death ended the recording is the sentinel's question,
		// answered by the coverage proof over the log this unblocks.
		fmt.Fprintf(cfg.Notes, "note: no writer holds %s; the recorder has already finished\n", cfg.Log)
		return log, nil
	}

	// The lock is held, so a writer is alive and the pidfile's pid cannot be
	// stale — the recorder that took the lock is the one the parent recorded.
	gone, err := signalRecorder(pid)
	if err != nil {
		log.Close()
		return nil, err
	}
	if gone {
		// A live writer holds the log and the pidfile names a process that
		// does not exist. That is a broken contract, not a case to reason
		// around: signalling the real holder would mean guessing who it is.
		log.Close()
		return nil, fmt.Errorf("a writer holds %s but pidfile %s names pid %d, which does not exist: the pidfile does not name the process that holds the log",
			cfg.Log, cfg.PidFile, pid)
	}

	// Wait on the LOCK, not on the pid: its release is the kernel-guaranteed
	// writer-is-gone event, and it carries no pid-reuse hazard.
	deadline := cfg.Clock.Now().Add(recorderStopTimeout)
	for {
		select {
		case <-ctx.Done():
			log.Close()
			return nil, ctx.Err()
		case <-cfg.Clock.After(recorderStopPoll):
		}

		held, err := tryLockExclusive(log)
		if err != nil {
			log.Close()
			return nil, err
		}
		if held {
			return log, nil
		}
		if !cfg.Clock.Now().Before(deadline) {
			log.Close()
			return nil, fmt.Errorf("recorder pid %d still holds %s %s after SIGTERM; refusing to read a log a writer can still append to",
				pid, cfg.Log, recorderStopTimeout)
		}
	}
}

// drainVerdict is what the drain wait concluded about one rule it could not
// clear. It carries the reason as well as the prose because the two outcomes
// are genuinely different faults: drain_timeout means the rule is still there
// and still behind, rule_absent means it is gone. Collapsing both into
// drain_timeout would name the wait instead of the fault, and Reason is a
// published vocabulary that reaches the JSON output.
type drainVerdict struct {
	reason UnobservableReason
	note   string
}

// drainWait is the final liveness check: did each rule evaluate through the
// end of the window? A rule that cannot answer within drainTimeout is
// unobservable, never a pass. It returns one verdict per rule it could not
// clear (keyed by UID); an error only for a hard failure of the wait itself.
//
// Two kinds of rule are excluded up front because draining them could not
// change a verdict: a rule the HEADER says was paused at the window open (the
// header, not the late-resolved definitions — see Header.pausedAtStart), and a
// rule whose last poll says Found == false (already unobservable via rule_absent).
func drainWait(ctx context.Context, cfg Config, src Source, defs []Definition, pausedAtStart map[string]bool,
	rt map[string]ruleTimings, polls []Poll, windowEnd time.Time, timeout time.Duration) (map[string]drainVerdict, error) {

	pending := make(map[string]string) // uid -> title, the shape observeAll wants
	for _, d := range defs {
		if pausedAtStart[d.UID] {
			continue
		}
		rulePolls := pollsForRule(polls, d.UID)
		if n := len(rulePolls); n > 0 && !rulePolls[n-1].Found {
			continue
		}
		// Evidence already in hand can satisfy the wait outright: a rule whose
		// recorded evaluations already reach past the end of the window has
		// answered the question, and polling it again asks nothing new.
		if !anyPollEvaluatedThrough(rulePolls, windowEnd) {
			pending[d.UID] = d.Title
		}
	}
	if len(pending) == 0 {
		return nil, nil
	}

	fmt.Fprintf(cfg.Notes, "drain wait: %d rule(s) have not yet evaluated through %s (limit %s)\n",
		len(pending), windowEnd.Format(time.RFC3339), timeout)

	deadline := cfg.Clock.Now().Add(timeout)
	verdicts := make(map[string]drainVerdict)
	for {
		uids := make([]string, 0, len(pending))
		for uid := range pending {
			uids = append(uids, uid)
		}
		sort.Strings(uids) // deterministic request order and message order

		observed, err := observeAll(ctx, src, pending, uids, cfg.Concurrency)
		if err != nil {
			return nil, fmt.Errorf("drain wait: %w", err)
		}
		for _, uid := range uids {
			obs, ok := observed[uid]
			if !ok {
				continue
			}
			rule := stateRuleByUID(obs.Rules, uid)
			if rule == nil {
				// A 2xx that parsed and carries no matching rule is an
				// authoritative "the rule is gone" — the transport retried
				// every transient failure long before this Observation
				// existed. It is knowable on the FIRST poll, so waiting the
				// rest of drainTimeout would spend two minutes to reach the
				// same verdict under a name that describes the wait rather
				// than the fault.
				verdicts[uid] = drainVerdict{
					reason: ReasonRuleAbsent,
					note: fmt.Sprintf("rule %q: absent from the state endpoint during the drain wait; there is no evaluation to wait for",
						pending[uid]),
				}
				delete(pending, uid)
				continue
			}
			if rule.IsPaused {
				// A paused rule does not evaluate, so this one can never catch
				// up and the rest of drainTimeout would buy nothing. The reason
				// stays drain_timeout: UnobservableReason is a published
				// vocabulary that reaches the JSON output, and the prose below
				// is where the detail belongs.
				verdicts[uid] = drainVerdict{
					reason: ReasonDrainTimeout,
					note: fmt.Sprintf("rule %q: paused before it evaluated through %s, so it never will",
						pending[uid], windowEnd.Format(time.RFC3339)),
				}
				delete(pending, uid)
				continue
			}
			if evaluatedThrough(rule.LastEvaluation, obs.Skew, obs.SkewBound, windowEnd) {
				delete(pending, uid)
			}
		}
		if len(pending) == 0 {
			return verdicts, nil
		}

		now := cfg.Clock.Now()
		if !now.Before(deadline) {
			for uid, title := range pending {
				verdicts[uid] = drainVerdict{
					reason: ReasonDrainTimeout,
					note: fmt.Sprintf("rule %q: did not evaluate through %s within the %s drain limit",
						title, windowEnd.Format(time.RFC3339), timeout),
				}
			}
			return verdicts, nil
		}

		// Re-ask no faster than the tightest cadence among the rules still
		// pending: a rule evaluating every 60s cannot answer differently 200ms
		// later, and hammering it would spend the run's request budget on
		// nothing.
		wait := deadline.Sub(now)
		for uid := range pending {
			if every := rt[uid].pollEvery; every > 0 {
				wait = min(wait, every)
			}
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-cfg.Clock.After(max(wait, 0)):
		}
	}
}

// anyPollEvaluatedThrough reports whether any recorded poll of a rule already
// proves it evaluated through windowEnd.
func anyPollEvaluatedThrough(polls []Poll, windowEnd time.Time) bool {
	for _, p := range polls {
		if !p.Found {
			continue
		}
		if evaluatedThrough(p.LastEvaluation, p.Skew(), p.SkewBound(), windowEnd) {
			return true
		}
	}
	return false
}

// evaluatedThrough is the drain wait's cross-domain comparison: a Grafana
// lastEvaluation is translated by its poll's skew, and the bound is SUBTRACTED
// (the pessimistic end) so an evaluation that only *might* have reached the
// window end is not counted as having reached it. A zero lastEvaluation never
// satisfies the wait.
func evaluatedThrough(lastEval time.Time, skew, bound time.Duration, windowEnd time.Time) bool {
	if lastEval.IsZero() {
		return false
	}
	return !lastEval.Add(-skew).Add(-bound).Before(windowEnd)
}

// mergeDrainTimeouts folds the drain wait's I/O verdicts into the pure Result,
// running after decide so that function stays pure of its arguments. It returns
// its own error rather than mutating decide's so neither hides the other: a run
// faulted by both the coverage proof and the drain wait must name both.
func mergeDrainTimeouts(res Result, drained map[string]drainVerdict) (Result, error) {
	if len(drained) == 0 {
		return res, nil
	}
	if res.Coverage == nil {
		res.Coverage = make(map[string]CoverageResult, len(drained))
	}

	var names []string
	for i := range res.Verdicts {
		uid := res.Verdicts[i].RuleUID
		verdict, ok := drained[uid]
		if !ok {
			continue
		}
		cov := res.Coverage[uid]
		cov.Unobservable = true
		cov.Proved = false
		if cov.Reason == "" {
			// The FIRST reason wins, as it does inside proveCoverage: a rule
			// the coverage proof already faulted keeps the fault it was
			// actually caught by.
			cov.Reason = verdict.reason
		}
		cov.Notes = append(cov.Notes, verdict.note)
		res.Coverage[uid] = cov

		if res.Verdicts[i].Outcome != OutcomeUnobservable {
			names = append(names, fmt.Sprintf("%s (%s)", res.Verdicts[i].Alert, verdict.reason))
		}
		res.Verdicts[i].Outcome = OutcomeUnobservable
		res.Verdicts[i].Note = strings.Join(cov.Notes, "; ")
	}
	if len(names) == 0 {
		// Every drained rule was already unobservable for an earlier reason,
		// so decide's own error already stops the run. Adding a second error
		// saying the same thing would only make the message longer.
		return res, nil
	}
	return res, fmt.Errorf("gate: %d rule(s) unobservable at the drain wait: %s", len(names), strings.Join(names, "; "))
}
