package gate

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

// LogSchemaVersion is the version stamped into every log header. A log with
// any other value is a read error, never a best-effort read: the log is the
// gate's only evidence, and misreading a stale shape is a fail-open.
const LogSchemaVersion = 1

// RecordType tags each JSONL line. There are exactly three, and a poll record
// IS the heartbeat — there is deliberately no separate heartbeat type.
type RecordType string

const (
	RecordHeader  RecordType = "header"
	RecordPoll    RecordType = "poll"
	RecordStopped RecordType = "stopped"
)

// missingSeriesReason is the reason Grafana parks a disappearing series at
// ("Normal (MissingSeries)") for a couple of evaluations before deleting the
// instance. Reading that as a recovery would turn a disappearing series into a
// fake recovery, so the markers below route it to Vanished.
const missingSeriesReason = "MissingSeries"

// LoggedRule is the per-rule identity written into the header. Together with
// the header URL it IS the log's identity, which check validates, and it
// supplies the alert set in check mode.
type LoggedRule struct {
	UID    string `json:"uid"`
	Title  string `json:"title"`
	Folder string `json:"folder"`
	Group  string `json:"group"`
	// ForSeconds, IntervalSeconds, NoDataState and ExecErrState are purely
	// forensic: a resolve-time snapshot that makes the uploaded artifact
	// self-describing to a human reading it after the runner is gone. check
	// never converts them back into a Definition — it always re-resolves
	// definitions from the ruler API.
	ForSeconds      float64 `json:"for_seconds"`
	IntervalSeconds int     `json:"interval_seconds"`
	// IsPaused is load-bearing (beside PollEverySeconds): the pause state at
	// record start, the only moment `skipped` can honestly mean. decide reads
	// it via Header.pausedAtStart, never a ruler read taken after the window.
	IsPaused     bool   `json:"is_paused"`
	NoDataState  string `json:"no_data_state"`
	ExecErrState string `json:"exec_err_state"`
	// PollEverySeconds is the cadence this recording ACTUALLY used. Load-bearing:
	// check derives maxGap from it, never from the definitions — getting that
	// wrong is fail-open in the faster-override direction.
	PollEverySeconds float64 `json:"poll_every_seconds"`
}

// Header is the log's first line: what was recorded, from where, and when the
// recording started. It carries no States field — recording is deliberately
// unfiltered, so the same log can be re-classified under different --states
// without re-recording.
type Header struct {
	SchemaVersion  int          `json:"schema_version"`
	URL            string       `json:"url"` // the log's identity
	GrafanaVersion string       `json:"grafana_version"`
	StartedAt      time.Time    `json:"started_at"` // the record start
	Rules          []LoggedRule `json:"rules"`      // THE alert set
}

// pausedAtStart reports, per rule UID, whether the rule was paused when the
// recording opened. That instant — and no other — is what `skipped` means: a
// rule nobody was watching on purpose.
//
// It is the authority for `skipped` in BOTH modes, and the reason is that no
// other source knows the right moment. `check` re-resolves the definitions
// AFTER the window closed, so Definition.IsPaused there describes the present,
// not the window: a rule that fired and was then paused would read as skipped,
// its firing would never be classified, and under --allow-paused the run would
// pass. The header cannot drift that way, because watch stamps it before the
// deploy step runs and single-step check stamps it from definitions resolved at
// the start of its own step.
//
// A UID the header does not name is reported NOT paused, which is the safe
// direction: it then reaches proveCoverage with no polls and fails closed as
// a heartbeat gap, rather than being waved through as legitimately unwatched.
func (h Header) pausedAtStart() map[string]bool {
	paused := make(map[string]bool, len(h.Rules))
	for _, lr := range h.Rules {
		paused[lr.UID] = lr.IsPaused
	}
	return paused
}

// Poll is one reduced observation of one rule — the log's heartbeat and the
// only input the pure coverage and classification layers ever see.
type Poll struct {
	RuleUID    string    `json:"rule_uid"`
	GrafanaNow time.Time `json:"grafana_now"` // the response's Date header
	// SkewMS, SkewBoundMS and LatencyMS are milliseconds for JSONL
	// compactness ONLY. The pure layer never touches raw ms: it reads
	// Skew(), SkewBound() and Latency() below, which convert at the
	// (de)serialization boundary.
	SkewMS      int64 `json:"skew_ms"`
	SkewBoundMS int64 `json:"skew_bound_ms"`
	LatencyMS   int64 `json:"latency_ms"`
	// Found false means an authoritative 2xx in which this rule was absent —
	// never a transport failure, which the transport retries and never turns
	// into a Poll. The coverage proof turns it into unobservable.
	Found bool `json:"found"`
	// State, Health and LastError are the raw rule-level strings, reporting
	// only and never classified.
	State     string `json:"state,omitempty"`
	Health    string `json:"health,omitempty"`
	LastError string `json:"last_error,omitempty"`
	// omitzero, not omitempty: a not-found poll (and a paused rule) has no
	// evaluation time, and writing "0001-01-01T00:00:00Z" into an artifact
	// humans and jq read invites reading it as a real timestamp.
	LastEvaluation time.Time      `json:"last_evaluation,omitzero"`
	IsPaused       bool           `json:"is_paused"`
	Histogram      map[string]int `json:"histogram,omitempty"` // written, never analysed
	// Reasons counts this poll's non-empty instance reasons, e.g.
	// {"NoData":1091,"Error":14}; nil when none. Reporting-only, and the only
	// place composite states stay visible (they are canonical normal, dropped
	// from Abnormal). Keys are raw reason strings and can be comma-joined
	// composites ("KeepLast, MissingSeries"), so consumers must test membership
	// via reasonNames and never index a literal key.
	Reasons map[string]int `json:"reasons,omitempty"`
	// Abnormal holds the instances whose CANONICAL state is not normal.
	// "Normal (NoData)" and "Normal (Error)" are canonical normal and are
	// deliberately not retained here.
	Abnormal []Instance `json:"abnormal,omitempty"`
	// Cleared and Vanished are instance keys that left the abnormal set,
	// resolved against the SAME response — a clear and a discontinuity are not
	// the same fact.
	Cleared  []string `json:"cleared,omitempty"`
	Vanished []string `json:"vanished,omitempty"`
}

// Skew is the signed clock skew of this poll.
func (p Poll) Skew() time.Duration { return time.Duration(p.SkewMS) * time.Millisecond }

// SkewBound is the uncertainty on Skew — the tolerance every cross-domain
// comparison applies alongside it.
func (p Poll) SkewBound() time.Duration { return time.Duration(p.SkewBoundMS) * time.Millisecond }

// Latency is the wall time this poll's request took, feeding the budget check.
func (p Poll) Latency() time.Duration { return time.Duration(p.LatencyMS) * time.Millisecond }

// Reducer turns each Observation into the single Poll record that goes into
// the log. It holds the previous poll's abnormal instance keys per rule, which
// is all the state the transition markers need.
//
// A Reducer is safe for concurrent use: watch polls a fleet of rules
// concurrently and every one of those goroutines reduces through the same
// instance, because the per-rule marker state has to live in one place. The
// lock is per-Reducer rather than per-rule — Reduce only touches maps and
// slices, so it never blocks on I/O while holding it.
type Reducer struct {
	mu           sync.Mutex
	prevAbnormal map[string]map[string]struct{}
}

func NewReducer() *Reducer {
	return &Reducer{prevAbnormal: make(map[string]map[string]struct{})}
}

// Reduce selects the rule identified by uid out of obs and reduces it to a
// Poll. Selection is BY UID, never by title: a filtered response can carry
// several rules sharing one title, and picking the first would silently watch
// the wrong rule.
//
// The reduction keeps the rule-level fields, the raw totals histogram, the
// reason counts, and only the instances whose canonical state is not normal.
// That makes per-poll size independent of NORMAL cardinality — not of
// cardinality outright: a rule with 449 firing instances still stores all 449.
func (r *Reducer) Reduce(uid string, obs Observation) Poll {
	r.mu.Lock()
	defer r.mu.Unlock()

	p := Poll{
		RuleUID:     uid,
		GrafanaNow:  obs.GrafanaNow,
		SkewMS:      obs.Skew.Milliseconds(),
		SkewBoundMS: obs.SkewBound.Milliseconds(),
		LatencyMS:   obs.Latency.Milliseconds(),
	}

	rule := stateRuleByUID(obs.Rules, uid)
	if rule == nil {
		// An authoritative "the rule is absent". No markers are computed and
		// the previous abnormal set is kept untouched: if the rule comes back
		// with an instance missing, the next poll still reports that instance
		// as vanished rather than losing the transition entirely.
		return p
	}

	p.Found = true
	p.State = rule.State
	p.Health = rule.Health
	p.LastError = rule.LastError
	p.LastEvaluation = rule.LastEvaluation
	p.IsPaused = rule.IsPaused
	p.Histogram = rule.Totals

	// present indexes every instance in THIS response, normal ones included —
	// the markers below must resolve a departed key against the same response,
	// which is impossible from the abnormal subset alone.
	present := make(map[string]Instance, len(rule.Instances))
	curAbnormal := make(map[string]struct{})
	for _, inst := range rule.Instances {
		key := instanceKey(inst.Labels)
		present[key] = inst
		if inst.Reason != "" {
			if p.Reasons == nil {
				p.Reasons = make(map[string]int)
			}
			p.Reasons[inst.Reason]++
		}
		if inst.State != StateNormal {
			p.Abnormal = append(p.Abnormal, inst)
			curAbnormal[key] = struct{}{}
		}
	}

	for key := range r.prevAbnormal[uid] {
		if _, still := curAbnormal[key]; still {
			continue
		}
		inst, found := present[key]
		switch {
		case !found:
			// Fully absent from the response: a discontinuity, not a recovery.
			p.Vanished = append(p.Vanished, key)
		case reasonNames(inst.Reason, missingSeriesReason):
			// The vanish in disguise, caught one poll earlier than the fully
			// absent case.
			p.Vanished = append(p.Vanished, key)
		default:
			// Present as canonical normal without a MissingSeries reason.
			p.Cleared = append(p.Cleared, key)
		}
	}
	// Map iteration is unordered; sort so a log line is byte-stable for a
	// given poll and a golden fixture stays meaningful.
	sort.Strings(p.Cleared)
	sort.Strings(p.Vanished)

	r.prevAbnormal[uid] = curAbnormal
	return p
}

// seedFrom restores the marker state above from polls that are already in the
// log, so the first poll a NEW Reducer produces compares against the last poll
// the previous one wrote instead of against an empty set.
//
// It exists for the one place a recording changes hands: watch's parent takes
// the first observation of every rule and its detached child continues from
// there. Without the seed, an instance that is abnormal in the parent's
// observation and gone by the child's first poll produces no marker at all —
// it leaves the record as though it had never been bad, the same fail-open a
// misread MissingSeries causes, reached through the handoff instead.
//
// Not-found polls are skipped, mirroring Reduce: an absent rule leaves the
// previous abnormal set untouched rather than emptying it.
func (r *Reducer) seedFrom(polls []Poll) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, p := range polls {
		if !p.Found {
			continue
		}
		keys := make(map[string]struct{}, len(p.Abnormal))
		for _, inst := range p.Abnormal {
			keys[instanceKey(inst.Labels)] = struct{}{}
		}
		r.prevAbnormal[p.RuleUID] = keys
	}
}

// stateRuleByUID picks one rule out of a state response BY UID (nil = the
// authoritative "rule absent"). Never by title: the ?rule_name= filter is a
// title filter and can return several rules sharing a title. The single
// selection for the package — Reduce and the drain wait both use it.
func stateRuleByUID(rules []StateRule, uid string) *StateRule {
	for i := range rules {
		if rules[i].UID == uid {
			return &rules[i]
		}
	}
	return nil
}

// reasonNames reports whether reason names want. Newer Grafana versions
// comma-join several reasons into one string, so this tests membership rather
// than equality.
func reasonNames(reason, want string) bool {
	for part := range strings.SplitSeq(reason, ",") {
		if strings.TrimSpace(part) == want {
			return true
		}
	}
	return false
}

// VerifyNormalInstancesVisible checks, on a first observation, that the state
// endpoint really returns normal instances: if it ever stops, the reduction's
// "keep the non-normal" becomes "keep everything the API sent", a silent
// fail-open in the transition markers. Counts are summed over every totals key
// whose lowercased name is "normal" or "inactive" — never a literal key, since
// the vocabulary is mixed and its case has drifted.
func VerifyNormalInstancesVisible(rules []StateRule) error {
	for _, r := range rules {
		var claimed int
		for k, v := range r.Totals {
			switch strings.ToLower(k) {
			case "normal", "inactive":
				claimed += v
			}
		}
		if claimed == 0 {
			continue
		}
		if hasNormalInstance(r.Instances) {
			continue
		}
		return fmt.Errorf(
			"rule %q (%s): totals claim %d normal instances but the response returned none — "+
				"the state endpoint no longer returns normal instances, which the reduction depends on",
			r.Title, r.UID, claimed)
	}
	return nil
}

func hasNormalInstance(instances []Instance) bool {
	for _, inst := range instances {
		if inst.State == StateNormal {
			return true
		}
	}
	return false
}

// headerRecord, pollRecord and stoppedRecord are the three wire shapes. The
// type tag is a real field on each line rather than an envelope, so a human
// (or jq) reading an uploaded log sees flat records.
type headerRecord struct {
	Type RecordType `json:"type"`
	Header
}

type pollRecord struct {
	Type RecordType `json:"type"`
	Poll
}

type stoppedRecord struct {
	Type RecordType `json:"type"`
	At   time.Time  `json:"at"`
}

// Writer appends records to the JSONL log. It is append-only by construction —
// O_APPEND|O_CREATE|O_WRONLY, never O_TRUNC — so no writer can ever destroy
// evidence a previous one recorded. An exclusive non-blocking flock makes a
// second writer fail immediately rather than interleave.
type Writer struct {
	mu      sync.Mutex
	f       *os.File
	enc     *json.Encoder
	clock   Clock
	stopped bool
}

// NewWriter opens path for appending and takes the exclusive lock. A second
// writer on the same path fails here, immediately — it never blocks and never
// waits, because two recorders on one log means one of them is recording a
// window nobody will classify.
func NewWriter(path string, clock Clock) (*Writer, error) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, fmt.Errorf("open log %s: %w", path, err)
	}
	if err := lockExclusive(f); err != nil {
		f.Close()
		return nil, fmt.Errorf("lock log %s: %w (another writer holds it)", path, err)
	}
	return &Writer{f: f, enc: json.NewEncoder(f), clock: clock}, nil
}

// WriteHeader writes line 1 and stamps the current schema version, so no
// caller can leave it at zero. It refuses a non-empty file: the log already
// has a header, and a second one would make ReadLog's "header is line 1"
// contract a lie. In watch's handoff the parent writes the header and the
// detached child only appends polls.
func (w *Writer) WriteHeader(h Header) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.stopped {
		return fmt.Errorf("log writer already stopped")
	}
	info, err := w.f.Stat()
	if err != nil {
		return fmt.Errorf("stat log: %w", err)
	}
	if info.Size() != 0 {
		return fmt.Errorf("log %s is not empty: it already has a header", w.f.Name())
	}
	h.SchemaVersion = LogSchemaVersion
	if err := w.enc.Encode(headerRecord{Type: RecordHeader, Header: h}); err != nil {
		return fmt.Errorf("write log header: %w", err)
	}
	return nil
}

// WritePoll appends one poll record — the heartbeat.
func (w *Writer) WritePoll(p Poll) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.stopped {
		return fmt.Errorf("log writer already stopped")
	}
	if err := w.enc.Encode(pollRecord{Type: RecordPoll, Poll: p}); err != nil {
		return fmt.Errorf("write poll for rule %s: %w", p.RuleUID, err)
	}
	return nil
}

// Stop finishes recording in a fixed order that must not be rearranged: let the
// in-flight write finish, append the sentinel, fsync, release — any other order
// can leave a sentinel that was never preceded by the polls it vouches for.
// The sentinel uses the writer's OWN stop time; check does the `to` comparison
// after this has exited. Calling Stop twice is a no-op (watch reaches it from a
// signal handler and a defer).
func (w *Writer) Stop() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.stopped {
		return nil
	}
	w.stopped = true

	encErr := w.enc.Encode(stoppedRecord{Type: RecordStopped, At: w.clock.Now()})
	syncErr := w.f.Sync()
	closeErr := w.f.Close()

	switch {
	case encErr != nil:
		return fmt.Errorf("write stopped sentinel: %w", encErr)
	case syncErr != nil:
		return fmt.Errorf("fsync log: %w", syncErr)
	case closeErr != nil:
		return fmt.Errorf("close log: %w", closeErr)
	}
	return nil
}

// Close releases the file and the lock WITHOUT writing a sentinel. It exists
// for exactly one caller: watch's parent, which writes the header and then
// hands the log to the detached child that will finish it. A sentinel here
// would tell check the recording ended before the child had even started.
func (w *Writer) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.stopped {
		return nil
	}
	w.stopped = true
	if err := w.f.Close(); err != nil {
		return fmt.Errorf("close log: %w", err)
	}
	return nil
}

// ReadLogHeader reads ONLY line 1 — the one read safe while a writer may still
// hold the log. The header is written once by watch's parent before any child
// appends a byte, so line 1 is immutable. It lets check fail closed EARLY on a
// wrong URL or unresolvable rule; it is advisory only, and the authoritative
// identity read is still ReadLog after the writer exits.
func ReadLogHeader(path string) (Header, error) {
	f, err := os.Open(path)
	if err != nil {
		return Header{}, fmt.Errorf("read log header %s: %w", path, err)
	}
	defer f.Close()

	line, err := bufio.NewReader(f).ReadString('\n')
	if err != nil {
		// io.EOF included: a log whose first line has no terminating newline is
		// a log whose header was never fully written, which is not a header.
		return Header{}, fmt.Errorf("log %s: no complete header on line 1: %w", path, err)
	}

	var rec headerRecord
	if err := json.Unmarshal([]byte(line), &rec); err != nil {
		return Header{}, fmt.Errorf("log %s line 1: unparseable header: %w", path, err)
	}
	if rec.Type != RecordHeader {
		return Header{}, fmt.Errorf("log %s line 1: got record type %q; the header must be line 1", path, rec.Type)
	}
	if rec.SchemaVersion != LogSchemaVersion {
		return Header{}, fmt.Errorf(
			"log %s: schema version %d is not %d — this log was written by a different version of the gate",
			path, rec.SchemaVersion, LogSchemaVersion)
	}
	return rec.Header, nil
}

// ReadLog reads the whole log once and returns its header, its polls in
// recorded order, and the sentinel time when one is present (nil when the
// recording never finished — check turns that into unobservable, never a
// pass).
//
// Call this only after the writer has exited. Reading a log a writer can still
// append to can only produce a shorter window than the one that was recorded.
//
// The parse rules are deliberately the crudest possible: the header must be
// line 1 with a matching schema version, and ANY unparseable line — including
// the last, or one after a sentinel — is an error, full stop. A truncated log
// is evidence something killed the recorder, which must not pass.
func ReadLog(path string) (Header, []Poll, *time.Time, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return Header{}, nil, nil, fmt.Errorf("read log %s: %w", path, err)
	}

	lines := strings.Split(string(b), "\n")
	// A complete record always ends with the encoder's newline, so the split
	// leaves one trailing empty element. Drop exactly that one; any other
	// empty line stays and fails below as unparseable.
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	if len(lines) == 0 {
		return Header{}, nil, nil, fmt.Errorf("log %s is empty: the header must be line 1", path)
	}

	var (
		header   Header
		polls    []Poll
		sentinel *time.Time
	)
	for i, line := range lines {
		var probe struct {
			Type RecordType `json:"type"`
		}
		if err := json.Unmarshal([]byte(line), &probe); err != nil {
			return Header{}, nil, nil, fmt.Errorf("log %s line %d: unparseable record: %w", path, i+1, err)
		}
		if sentinel != nil {
			return Header{}, nil, nil, fmt.Errorf(
				"log %s line %d: a %q record follows the stopped sentinel — the log had a second writer",
				path, i+1, probe.Type)
		}
		if (i == 0) != (probe.Type == RecordHeader) {
			return Header{}, nil, nil, fmt.Errorf(
				"log %s line %d: got record type %q; the header must be line 1 and appear only once",
				path, i+1, probe.Type)
		}

		switch probe.Type {
		case RecordHeader:
			var rec headerRecord
			if err := json.Unmarshal([]byte(line), &rec); err != nil {
				return Header{}, nil, nil, fmt.Errorf("log %s line 1: unparseable header: %w", path, err)
			}
			if rec.SchemaVersion != LogSchemaVersion {
				return Header{}, nil, nil, fmt.Errorf(
					"log %s: schema version %d is not %d — this log was written by a different version of the gate",
					path, rec.SchemaVersion, LogSchemaVersion)
			}
			header = rec.Header
		case RecordPoll:
			var rec pollRecord
			if err := json.Unmarshal([]byte(line), &rec); err != nil {
				return Header{}, nil, nil, fmt.Errorf("log %s line %d: unparseable poll: %w", path, i+1, err)
			}
			polls = append(polls, rec.Poll)
		case RecordStopped:
			var rec stoppedRecord
			if err := json.Unmarshal([]byte(line), &rec); err != nil {
				return Header{}, nil, nil, fmt.Errorf("log %s line %d: unparseable sentinel: %w", path, i+1, err)
			}
			at := rec.At
			sentinel = &at
		default:
			return Header{}, nil, nil, fmt.Errorf("log %s line %d: unknown record type %q", path, i+1, probe.Type)
		}
	}

	return header, polls, sentinel, nil
}
