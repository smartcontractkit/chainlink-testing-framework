package gate

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

// DaemonChildFlag is the hidden flag the parent passes when it re-execs itself
// as the detached recorder. It is deliberately absent from the CLI's
// usage text: an operator never types it, and a child started by hand against
// a log no parent prepared fails immediately on the header read.
const DaemonChildFlag = "--daemon-child"

// ReadyFDFlag names the inherited descriptor the child reports readiness on.
// The parent passes the write end of a pipe as descriptor 3 and waits for one
// byte, so "the recorder is running" is a POSITIVE signal from the child
// itself — it has read the header, taken the log's flock and entered its poll
// loop — and not an assumption drawn from surviving a timer. A timer cannot
// tell a healthy child from one that is about to die on a slow runner, and
// getting that wrong means watch returns success over a recording that never
// happened.
const ReadyFDFlag = "--ready-fd"

// childReadyTimeout bounds that wait. Everything before the signal is local —
// fork, exec, one read of a log holding a header and a handful of polls — so
// the real figure is milliseconds; this is loose enough for a badly overloaded
// runner and still fails closed rather than hanging the pipeline.
const childReadyTimeout = 30 * time.Second

// daemonLogTailBytes bounds how much of a dead child's output the parent
// quotes back. A child dies in its first few lines or not at all.
const daemonLogTailBytes = 4096

// WatchConfig is the record step's whole input.
//
// It has no To field and must never gain one: watch writes the stopped
// sentinel with its OWN stop time and makes no comparison against `to`, which
// only check knows. Passing `to` here would give two components an
// opinion about the same comparison, and the recorder's opinion is the one
// that cannot be trusted — it exits before the grace it would have to wait for.
//
// It has no States field either, and watch has no --states flag: recording is
// deliberately unfiltered. The reduction keeps every non-normal instance and
// the transition markers key off the same predicate, so neither consults
// States. The payoff is real — because the log is raw evidence, one recording
// can be re-classified under different --states without re-recording — and the
// Header carries no States field for the same reason.
type WatchConfig struct {
	// URL and Token are the connection details. The CLI reads both from the
	// environment and never from a flag; Token is never logged and never
	// enters an error string.
	URL, Token string

	// Alerts are the operator-supplied names, one per line, in any of the forms
	// Resolve accepts. Empty lines are discarded by Resolve.
	Alerts []string
	Folder string

	// Out is the JSONL log path. PidFile and DaemonLog default to
	// <Out>.pid and <Out>.daemon.log — the same convention check uses to find
	// the recorder it must stop, so nothing has to be wired by hand.
	Out       string
	PidFile   string
	DaemonLog string

	// Until is an optional hard stop for the child. Zero means "record until
	// signalled", which is the normal case: check sends SIGTERM when its
	// collection loop ends.
	Until time.Time

	// PollEvery is the --poll-interval override, used verbatim for every rule
	// and never clamped. Zero means each rule polls at half its own evaluation
	// interval. Whatever this resolves to is written into the header as the
	// cadence actually used, and that header value — never a re-derivation from
	// the definitions — is what check derives maxGap from.
	PollEvery time.Duration

	Concurrency int
	Clock       Clock

	// Notes is where the parent prints what an operator has to see before the
	// deploy step runs: resolve notes, the cadence per rule, the rules it will
	// not wait for. nil discards them. The library prints nothing else — the
	// CLI owns presentation.
	Notes io.Writer
}

func (cfg WatchConfig) withDefaults() WatchConfig {
	if cfg.Clock == nil {
		cfg.Clock = SystemClock{}
	}
	if cfg.Notes == nil {
		cfg.Notes = io.Discard
	}
	if cfg.Concurrency < 1 {
		cfg.Concurrency = 1
	}
	if cfg.PidFile == "" && cfg.Out != "" {
		cfg.PidFile = cfg.Out + ".pid"
	}
	if cfg.DaemonLog == "" && cfg.Out != "" {
		cfg.DaemonLog = cfg.Out + ".daemon.log"
	}
	return cfg
}

func (cfg WatchConfig) validate() error {
	if cfg.URL == "" {
		return fmt.Errorf("watch: no grafana url (it is the log's identity, which check validates)")
	}
	if cfg.Out == "" {
		return fmt.Errorf("watch: no log path")
	}
	named := 0
	for _, a := range cfg.Alerts {
		if strings.TrimSpace(a) != "" {
			named++
		}
	}
	if named == 0 {
		return fmt.Errorf("watch: no alert names given; there is nothing to record")
	}
	// An --until already in the past would make the child stop before it ever
	// polled, and the parent would then report a child that never reported
	// ready — a true statement about a config mistake, but a confusing one.
	if !cfg.Until.IsZero() && !cfg.Until.After(cfg.Clock.Now()) {
		return fmt.Errorf("watch: --until %s is not in the future", cfg.Until.Format(time.RFC3339))
	}
	return nil
}

// Watch is the record step's parent process. It returns only once the window
// is genuinely being recorded:
//
//	version gate -> resolve definitions and names -> derive timings ->
//	open the log and write the header -> ONE observation of every non-skipped
//	rule -> verify normal instances are visible -> check the schedule budget ->
//	detach the child -> wait for the child to report that it is recording ->
//	write the pidfile -> return.
//
// The first-observation wait is not a convenience. Returning before it would
// leave the deploy inside [from, first_poll] with no evidence — the exact
// blind interval the record-then-check split exists to remove — and it is
// also what surfaces auth, name-resolution and parse failures BEFORE
// deploy.sh runs rather than ten minutes later.
func Watch(ctx context.Context, cfg WatchConfig) error {
	cfg = cfg.withDefaults()
	if err := cfg.validate(); err != nil {
		return err
	}

	src := NewHTTPSource(cfg.URL, cfg.Token, cfg.Clock)
	prep, err := prepareWatch(ctx, cfg, src)
	if err != nil {
		return err
	}

	// Hand the log over with Close, never Stop: a sentinel here would tell
	// check the recording ended before the child had even started. Closing
	// also releases the flock the child is about to take.
	if err := prep.writer.Close(); err != nil {
		return err
	}

	child, err := spawnChild(cfg)
	if err != nil {
		return fmt.Errorf("detach recorder: %w", err)
	}
	if err := waitForChildReady(cfg, child); err != nil {
		return err
	}

	// The PARENT writes the pidfile, not the child: check must find the pid the
	// instant Watch returns, and a child writing its own would race the very
	// next step of the pipeline. That is why the child is never given
	// --pidfile at all.
	//
	// It is written only once the child has reported ready, so no path through
	// this function leaves a pidfile naming a process that is not recording.
	// Pids are reused: a stale pidfile is a live process somewhere else, and a
	// pipeline that ignored this function's error would SIGTERM it.
	if err := writePidFile(cfg.PidFile, child.cmd.Process.Pid); err != nil {
		_ = child.cmd.Process.Kill()
		_ = os.Remove(cfg.PidFile)
		return err
	}
	fmt.Fprintf(cfg.Notes, "recording to %s (pid %d, pidfile %s, output %s)\n",
		cfg.Out, child.cmd.Process.Pid, cfg.PidFile, cfg.DaemonLog)
	return nil
}

// waitForChildReady waits for the child's own readiness byte, and treats every
// other outcome as a failure to record: the pipe closing without a byte (the
// child died on its way to the loop), the process exiting, or the timeout.
// Each of the three quotes the daemon log, which is the only place a detached
// process can explain itself.
func waitForChildReady(cfg WatchConfig, child detachedChild) error {
	defer child.ready.Close()

	// This goroutine outlives the function when the child keeps running. It
	// stays behind only to reap a child that dies while this short-lived parent
	// is still alive, and costs nothing.
	exited := make(chan error, 1)
	go func() { exited <- child.cmd.Wait() }()

	signalled := make(chan error, 1)
	go func() {
		_, err := child.ready.Read(make([]byte, 1))
		signalled <- err
	}()

	fail := func(format string, args ...any) error {
		_ = child.cmd.Process.Kill()
		return fmt.Errorf("%s; its output was:\n%s",
			fmt.Sprintf(format, args...), daemonLogTail(cfg.DaemonLog, child.logOffset))
	}

	select {
	case err := <-signalled:
		if err == nil {
			return nil
		}
		// The child closed the pipe — by exiting — without ever reporting that
		// it had the log and was polling.
		return fail("the detached recorder never reported ready (%v)", err)
	case waitErr := <-exited:
		status := "exit status 0"
		if waitErr != nil {
			status = waitErr.Error()
		}
		return fail("the detached recorder exited before it started recording (%s)", status)
	case <-cfg.Clock.After(childReadyTimeout):
		return fail("the detached recorder did not report ready within %s", childReadyTimeout)
	}
}

// daemonLogTail quotes the end of the daemon log, starting at from — the size
// the file had when THIS run opened it. The offset is what keeps the quote
// honest when several runs share one --daemon-log path: without it the tail can
// name a previous run's failure as the current one's cause.
func daemonLogTail(path string, from int64) string {
	b, err := os.ReadFile(path)
	if err != nil {
		return fmt.Sprintf("(daemon log %s is unreadable: %v)", path, err)
	}
	if from > 0 && from <= int64(len(b)) {
		b = b[from:]
	}
	if len(b) > daemonLogTailBytes {
		b = b[len(b)-daemonLogTailBytes:]
	}
	if len(b) == 0 {
		return fmt.Sprintf("(this run wrote nothing to the daemon log %s)", path)
	}
	return strings.TrimRight(string(b), "\n")
}

// preparedWatch is what the parent has established by the time it is willing
// to detach: an open log with a header and one poll per non-skipped rule, the
// timings that produced them, and the latencies it measured doing so.
type preparedWatch struct {
	writer   *Writer
	header   Header
	timings  map[string]ruleTimings
	measured map[string]time.Duration
}

// prepareWatch is everything the parent does before it detaches. It takes a
// Source rather than building one so the paused-rule, first-observation,
// instance-visibility and budget behaviours are all testable with a scripted
// fake — only the process spawning needs a real binary.
func prepareWatch(ctx context.Context, cfg WatchConfig, src Source) (*preparedWatch, error) {
	version, err := src.Version(ctx)
	if err != nil {
		return nil, fmt.Errorf("read grafana version: %w", err)
	}
	if err := CheckGrafanaVersion(version); err != nil {
		return nil, err
	}

	defs, err := src.Definitions(ctx)
	if err != nil {
		return nil, fmt.Errorf("read rule definitions: %w", err)
	}
	resolved, notes, err := Resolve(defs, cfg.Alerts, cfg.Folder)
	if err != nil {
		return nil, err
	}
	for _, n := range notes {
		fmt.Fprintf(cfg.Notes, "note: %s\n", n)
	}

	rt, _, timingNotes := DeriveTimings(resolved, cfg.PollEvery)
	for _, n := range timingNotes {
		fmt.Fprintf(cfg.Notes, "note: %s\n", n)
	}
	for _, d := range resolved {
		// A cadence of zero would make the child spin: every rule is due the
		// instant it was marked. It also cannot be written into the header,
		// where check requires a positive value to derive maxGap from.
		if rt[d.UID].pollEvery <= 0 {
			return nil, fmt.Errorf("rule %q (%s) reports intervalSeconds=%d: there is no poll cadence to record at",
				d.Title, d.UID, d.IntervalSeconds)
		}
	}

	writer, err := NewWriter(cfg.Out, cfg.Clock)
	if err != nil {
		return nil, err
	}
	prep, err := openRecording(ctx, cfg, src, writer, version, resolved, rt)
	if err != nil {
		// Close, never Stop. The log keeps whatever was written and gets no
		// sentinel, so nothing can later mistake it for a finished recording.
		_ = writer.Close()
		return nil, err
	}
	return prep, nil
}

// openRecording writes the header, takes the first observation of every rule
// the recorder will actually watch, appends those observations as the log's
// first heartbeats, and only then decides whether the schedule is feasible.
func openRecording(ctx context.Context, cfg WatchConfig, src Source, writer *Writer,
	version string, resolved []Definition, rt map[string]ruleTimings) (*preparedWatch, error) {

	header := Header{
		SchemaVersion:  LogSchemaVersion,
		URL:            cfg.URL,
		GrafanaVersion: version,
		StartedAt:      cfg.Clock.Now(),
		Rules:          loggedRules(resolved, rt),
	}
	if err := writer.WriteHeader(header); err != nil {
		return nil, err
	}

	// A rule whose DEFINITION says is_paused is skipped: it is not waited for,
	// not scheduled and never polled. Waiting for one either hangs forever or
	// errors before the deploy, and recording polls for it would report an
	// in-window pause (coverage check 7) for a rule that was already paused
	// when the window opened — turning a skipped rule's exit 1 into an exit 2.
	// The header still names it, with is_paused true, so check reports it as
	// skipped from the definitions.
	var active []Definition
	activeTimings := make(map[string]ruleTimings, len(resolved))
	for _, d := range resolved {
		if d.IsPaused {
			fmt.Fprintf(cfg.Notes, "note: rule %q (%s) is paused: recorded as skipped, not waited for\n", d.Title, d.UID)
			continue
		}
		active = append(active, d)
		activeTimings[d.UID] = rt[d.UID]
		fmt.Fprintf(cfg.Notes, "recording %q (%s) every %s (maxGap %s)\n", d.Title, d.UID, rt[d.UID].pollEvery, rt[d.UID].maxGap)
	}

	polls, measured, err := firstObservations(ctx, src, active, NewReducer(), cfg.Concurrency, cfg.Notes)
	if err != nil {
		return nil, err
	}
	for _, p := range polls {
		if err := writer.WritePoll(p); err != nil {
			return nil, err
		}
	}

	// Budget last, on the latencies just measured — never on a fixed estimate.
	// Only the active rules count: a skipped rule is never polled and consumes
	// none of the capacity.
	if err := CheckBudget(activeTimings, measured, cfg.Concurrency); err != nil {
		return nil, err
	}

	return &preparedWatch{writer: writer, header: header, timings: rt, measured: measured}, nil
}

// loggedRules snapshots the resolved definitions into the header's rule list.
// Every field but PollEverySeconds is forensic — a resolve-time snapshot that
// makes an uploaded log self-describing — while PollEverySeconds is
// load-bearing: it is the cadence this recording actually used, and check
// derives maxGap from it rather than from the definitions.
func loggedRules(defs []Definition, rt map[string]ruleTimings) []LoggedRule {
	out := make([]LoggedRule, 0, len(defs))
	for _, d := range defs {
		out = append(out, LoggedRule{
			UID:              d.UID,
			Title:            d.Title,
			Folder:           d.Folder,
			Group:            d.Group,
			ForSeconds:       d.For.Seconds(),
			IntervalSeconds:  d.IntervalSeconds,
			IsPaused:         d.IsPaused,
			NoDataState:      d.NoDataState,
			ExecErrState:     d.ExecErrState,
			PollEverySeconds: rt[d.UID].pollEvery.Seconds(),
		})
	}
	return out
}

// firstObservations takes one observation of every rule in active, verifies
// that normal instances are visible in those very responses, and reduces each
// into the poll record that IS the window's first heartbeat — plus the measured
// latency of each, which is the only honest input to the budget check (a fixed
// estimate is worthless when one rule's payload is ~230x another's).
//
// Both entry paths share it: watch's parent, before it detaches, and
// single-step check's measurement pass, which keeps the polls as evidence
// rather than writing them to a log. Keeping one implementation is the point —
// the instance-visibility verification and the "absent is a warning, not an
// error" rule are exactly the places where two copies would silently drift, and
// a drift in either direction is fail-open.
//
// polls come back in `active` order, so a log written from them is byte-stable
// for a given set of observations.
func firstObservations(ctx context.Context, src Source, active []Definition, reducer *Reducer,
	concurrency int, notes io.Writer) ([]Poll, map[string]time.Duration, error) {

	titles := make(map[string]string, len(active))
	uids := make([]string, 0, len(active))
	for _, d := range active {
		titles[d.UID] = d.Title
		uids = append(uids, d.UID)
	}

	observed, err := observeAll(ctx, src, titles, uids, concurrency)
	if err != nil {
		return nil, nil, err
	}

	// Verify this before anything downstream relies on it: if the state
	// endpoint ever stops returning normal instances, the reduction's "keep
	// the non-normal ones" silently becomes "keep everything it happened to
	// send" and the transition markers lose their ground truth.
	for _, d := range active {
		if err := VerifyNormalInstancesVisible(observed[d.UID].Rules); err != nil {
			return nil, nil, err
		}
	}

	polls := make([]Poll, 0, len(active))
	measured := make(map[string]time.Duration, len(active))
	for _, d := range active {
		obs := observed[d.UID]
		measured[d.UID] = obs.Latency
		poll := reducer.Reduce(d.UID, obs)
		if !poll.Found {
			// Authoritative, not transient (the transport already retried
			// every transient failure): the rule resolved in the ruler API but
			// the state endpoint does not serve it. Recorded as Found=false,
			// which the coverage proof turns into unobservable — a note rather
			// than an error here, because the state endpoint can lag a freshly
			// created rule and the coverage proof fails closed either way.
			fmt.Fprintf(notes, "warning: rule %q (%s) is absent from the state endpoint; recorded as not found\n", d.Title, d.UID)
		}
		polls = append(polls, poll)
	}
	return polls, measured, nil
}

// observeAll polls every rule in uids concurrently, bounded by concurrency,
// and returns one Observation per rule that answered. Every rule is polled by
// TITLE (the ?rule_name= filter is a title filter) and selected out of the
// response by UID — a filtered response can carry several rules sharing one
// title.
//
// It returns the successful observations alongside the first error in UID
// order, so a caller that wants to keep the good heartbeats can, and the error
// message is the same on every run.
func observeAll(ctx context.Context, src Source, titles map[string]string, uids []string, concurrency int) (map[string]Observation, error) {
	if concurrency < 1 {
		concurrency = 1
	}
	var (
		mu          sync.Mutex
		out         = make(map[string]Observation, len(uids))
		firstErr    error
		firstErrUID string
	)
	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup
	for _, uid := range uids {
		wg.Go(func() {
			sem <- struct{}{}
			defer func() { <-sem }()

			obs, err := src.RuleState(ctx, titles[uid])

			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				if firstErr == nil || uid < firstErrUID {
					firstErr, firstErrUID = err, uid
				}
				return
			}
			out[uid] = obs
		})
	}
	wg.Wait()

	if firstErr != nil {
		return out, fmt.Errorf("poll rule %q (%s): %w", titles[firstErrUID], firstErrUID, firstErr)
	}
	return out, nil
}

// DaemonChildConfig is the detached recorder's whole input, and it is
// deliberately tiny. The rule set and every cadence come from the header the
// parent already wrote — one source of truth, no parent/child drift, and it
// exercises ReadLog's header path — and the connection details come from the
// inherited environment. Only the run facts the header does not carry travel
// in argv.
type DaemonChildConfig struct {
	URL, Token  string // from the inherited environment, never from argv
	Out         string
	Until       time.Time
	Concurrency int
	Clock       Clock
	// ReadyFD is the inherited descriptor to report readiness on (ReadyFDFlag).
	// Zero means nobody is waiting — a hand-started child — and the report is
	// then skipped rather than written to stdin.
	ReadyFD int
}

// RunDaemonChild is the detached recorder. The CLI dispatches to it when it
// sees DaemonChildFlag; nothing else ever calls it.
//
// It re-reads the log the parent wrote, restores the transition-marker state
// from the polls already in it, reopens the log for appending, takes the flock
// the parent released, and then polls until it is signalled or reaches Until.
func RunDaemonChild(ctx context.Context, cfg DaemonChildConfig) error {
	if cfg.Clock == nil {
		cfg.Clock = SystemClock{}
	}
	if cfg.Concurrency < 1 {
		cfg.Concurrency = 1
	}
	if cfg.Out == "" {
		return fmt.Errorf("recorder: no log path")
	}

	// Safe to read: the parent closed its writer before spawning this process,
	// and no other writer can hold the log's flock.
	header, polls, sentinel, err := ReadLog(cfg.Out)
	if err != nil {
		return err
	}
	if sentinel != nil {
		return fmt.Errorf("log %s already carries a stopped sentinel: another recorder finished it", cfg.Out)
	}
	// The header's URL is the log's identity. Checking it here
	// catches a child that inherited an environment pointing somewhere else,
	// before it appends a single poll from the wrong Grafana.
	if header.URL != cfg.URL {
		return fmt.Errorf("log %s records url %q but this recorder is configured for %q", cfg.Out, header.URL, cfg.URL)
	}

	titles, cadence, err := childSchedule(header)
	if err != nil {
		return err
	}

	writer, err := NewWriter(cfg.Out, cfg.Clock)
	if err != nil {
		return err
	}

	reducer := NewReducer()
	reducer.seedFrom(polls)

	// SIGTERM is how check stops the recorder; SIGINT is the
	// same request from a human at a terminal. Both are clean stops, so both
	// end with a sentinel. Registered before the readiness report, so a signal
	// arriving the moment the parent unblocks is already handled.
	sigCtx, stop := signal.NotifyContext(ctx, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	// Everything that can fail before a single poll has now succeeded: the
	// header parsed, the identity matched, the flock is held. That — and not
	// the mere fact of having been started — is what the parent waits for.
	if err := reportReady(cfg.ReadyFD); err != nil {
		return err
	}

	return watchLoop(sigCtx, watchLoopConfig{
		Src:         NewHTTPSource(cfg.URL, cfg.Token, cfg.Clock),
		Writer:      writer,
		Reducer:     reducer,
		Titles:      titles,
		Cadence:     cadence,
		Until:       cfg.Until,
		Concurrency: cfg.Concurrency,
		Clock:       cfg.Clock,
	})
}

// reportReady writes one byte to the inherited readiness descriptor and closes
// it. fd 0 means no parent is waiting: descriptor 0 is stdin, so it can never
// be a readiness pipe, which makes the zero value safe as "absent".
func reportReady(fd int) error {
	if fd == 0 {
		return nil
	}
	pipe := os.NewFile(uintptr(fd), "ready")
	if pipe == nil {
		return fmt.Errorf("readiness descriptor %d is not open", fd)
	}
	defer pipe.Close()
	if _, err := pipe.Write([]byte{'1'}); err != nil {
		return fmt.Errorf("report ready on descriptor %d: %w", fd, err)
	}
	return nil
}

// childSchedule derives what the child polls, and how often, from the header
// alone. The cadence comes from PollEverySeconds — the cadence the recording
// actually uses — and is never re-derived from the rule's evaluation interval,
// which would be a second authority for the same value and is fail-open in the
// faster-override direction. Paused rules are excluded here for the same reason
// the parent never polls them.
//
// It returns cadences and nothing else. maxGap, healthGrace and evalStaleAfter
// are coverage thresholds applied by the pure layer at classification time, so
// the recorder must not carry them: it would only be able to misuse them.
func childSchedule(h Header) (titles map[string]string, cadence map[string]time.Duration, err error) {
	titles = make(map[string]string, len(h.Rules))
	cadence = make(map[string]time.Duration, len(h.Rules))
	for _, lr := range h.Rules {
		if lr.IsPaused {
			continue
		}
		if lr.PollEverySeconds <= 0 {
			return nil, nil, fmt.Errorf("log header records poll_every_seconds=%v for rule %s (%q): there is no cadence to record at",
				lr.PollEverySeconds, lr.UID, lr.Title)
		}
		if _, duplicate := titles[lr.UID]; duplicate {
			return nil, nil, fmt.Errorf("log header names rule %s (%q) twice; its recorded cadence is ambiguous", lr.UID, lr.Title)
		}
		titles[lr.UID] = lr.Title
		cadence[lr.UID] = time.Duration(lr.PollEverySeconds * float64(time.Second))
	}
	return titles, cadence, nil
}

// watchLoopConfig is the child's working state: what to poll, how often, and
// where to append it. There is no threshold in here and no policy — the child
// records and classifies nothing.
type watchLoopConfig struct {
	Src         Source
	Writer      *Writer
	Reducer     *Reducer
	Titles      map[string]string        // uid -> title: poll by title, select by UID
	Cadence     map[string]time.Duration // uid -> pollEvery, as recorded in the header
	Until       time.Time
	Concurrency int
	Clock       Clock
}

// watchLoop is the child's whole working life: poll the rules that are due,
// reduce each observation to one poll record, append it, and — on a clean stop
// only — finish the log with the stopped sentinel.
//
// The sentinel policy is the load-bearing part. A clean stop (a signal, or
// Until) writes it; a hard error does NOT. A recorder that died must look
// exactly like a coverage gap to check, because it is one — writing a
// sentinel on the way out of a failure would hand check a "recording finished"
// claim about a window that stopped being observed.
func watchLoop(ctx context.Context, cfg watchLoopConfig) error {
	sched := NewScheduler(cfg.Cadence, cfg.Clock.Now())

	for {
		if ctx.Err() != nil {
			return cfg.Writer.Stop()
		}
		now := cfg.Clock.Now()
		if !cfg.Until.IsZero() && !now.Before(cfg.Until) {
			return cfg.Writer.Stop()
		}

		due := sched.Due(now)
		if len(due) == 0 {
			wait, ok := untilNextPoll(sched, cfg.Until, now)
			if !ok {
				// Nothing will ever come due: every watched rule is paused and
				// there is no hard stop. Wait for the signal — and still write
				// a sentinel, because "the recorder ran and finished" is
				// exactly what check needs to prove about the window.
				<-ctx.Done()
				return cfg.Writer.Stop()
			}
			select {
			case <-ctx.Done():
				return cfg.Writer.Stop()
			case <-cfg.Clock.After(wait):
			}
			continue
		}

		// Mark before polling, against the batch's own now: the next poll is
		// one cadence after this one was DUE, not one cadence after it
		// returned, so request latency cannot make the heartbeat spacing drift
		// towards maxGap.
		for _, uid := range due {
			sched.Mark(uid, now)
		}

		pollErr := cfg.pollBatch(ctx, due)
		if ctx.Err() != nil {
			// Signalled while a poll was in flight. The aborted poll's error is
			// not a recorder failure, and a clean stop wins over it: finish the
			// in-flight write, then the sentinel.
			return cfg.Writer.Stop()
		}
		if pollErr != nil {
			return pollErr
		}
	}
}

// pollBatch polls one round of due rules and appends every poll that
// succeeded, in due order, before returning the first failure. Writing the
// successes first is deliberate: a heartbeat that was genuinely observed is
// evidence, and dropping it because a different rule failed would turn one
// rule's transport failure into a coverage gap for the others.
func (cfg watchLoopConfig) pollBatch(ctx context.Context, uids []string) error {
	observed, obsErr := observeAll(ctx, cfg.Src, cfg.Titles, uids, cfg.Concurrency)
	for _, uid := range uids {
		obs, ok := observed[uid]
		if !ok {
			continue
		}
		if err := cfg.Writer.WritePoll(cfg.Reducer.Reduce(uid, obs)); err != nil {
			return err
		}
	}
	return obsErr
}

// untilNextPoll returns how long to wait for the next scheduled poll, cut
// short by Until when that comes first. ok is false when nothing will ever
// come due: no rules to poll and no hard stop.
func untilNextPoll(sched *Scheduler, until, now time.Time) (time.Duration, bool) {
	next, hasNext := sched.earliestDue()
	switch {
	case hasNext && (until.IsZero() || next.Before(until)):
		// keep next
	case !until.IsZero():
		next = until
	default:
		return 0, false
	}
	return max(next.Sub(now), 0), true
}

// writePidFile records the child's pid where check looks for it (its
// --pidfile, default <in>.pid). The format is the decimal pid and a newline,
// so `kill $(cat log.jsonl.pid)` works and ReadPidFile stays trivial.
func writePidFile(path string, pid int) error {
	if err := os.WriteFile(path, []byte(strconv.Itoa(pid)+"\n"), 0o644); err != nil {
		return fmt.Errorf("write pidfile %s: %w", path, err)
	}
	return nil
}

// ReadPidFile is the other side of that contract: the pid of the recorder
// check must stop before it may read the log.
func ReadPidFile(path string) (int, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return 0, fmt.Errorf("read pidfile %s: %w", path, err)
	}
	text := strings.TrimSpace(string(b))
	pid, err := strconv.Atoi(text)
	if err != nil {
		return 0, fmt.Errorf("pidfile %s: unparseable pid %q", path, text)
	}
	if pid <= 0 {
		return 0, fmt.Errorf("pidfile %s: %d is not a pid", path, pid)
	}
	return pid, nil
}
