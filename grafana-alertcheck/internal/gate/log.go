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
// gate's only evidence, and misreading a stale shape is a fail-open (§5).
const LogSchemaVersion = 1

// RecordType tags each JSONL line. There are exactly three, and a poll record
// IS the heartbeat — there is deliberately no separate heartbeat type (§4.6).
type RecordType string

const (
	RecordHeader  RecordType = "header"
	RecordPoll    RecordType = "poll"
	RecordStopped RecordType = "stopped"
)

// missingSeriesReason is the reason Grafana parks a disappearing series at
// ("Normal (MissingSeries)") for a couple of evaluations before deleting the
// instance. Reading that as a recovery is H2's named bug, so the markers below
// route it to Vanished (P1.2a).
const missingSeriesReason = "MissingSeries"

// LoggedRule is the per-rule identity written into the header. Together with
// the header URL it IS the log's identity, which check validates (§19.1 step
// 3), and it supplies the alert set in check mode.
type LoggedRule struct {
	UID    string `json:"uid"`
	Title  string `json:"title"`
	Folder string `json:"folder"`
	Group  string `json:"group"`
	// ForSeconds, IntervalSeconds, NoDataState and ExecErrState are purely
	// forensic: a resolve-time snapshot that makes the uploaded artifact
	// self-describing to a human reading it after the runner is gone (§21.3).
	// check never converts them back into a Definition — it always re-resolves
	// definitions from the ruler API (§19.1 step 2).
	ForSeconds      float64 `json:"for_seconds"`
	IntervalSeconds int     `json:"interval_seconds"`
	// IsPaused is NOT forensic, and is the second load-bearing field here
	// beside PollEverySeconds. It is the pause state at record start, which is
	// the only moment `skipped` can honestly mean (§12), and decide reads it
	// through Header.pausedAtStart rather than reading Definition.IsPaused off
	// a ruler read taken after the window had already closed. See that method
	// for what goes wrong the other way.
	IsPaused     bool   `json:"is_paused"`
	NoDataState  string `json:"no_data_state"`
	ExecErrState string `json:"exec_err_state"`
	// PollEverySeconds is the cadence this recording ACTUALLY used, after any
	// --poll-interval override. Load-bearing, not forensic: check derives
	// maxGap from it and never re-derives it from the definitions. Getting
	// that wrong is fail-open in the faster-override direction — a real
	// recorder gap would pass silently (see "Two authorities", P5).
	PollEverySeconds float64 `json:"poll_every_seconds"`
}

// Header is the log's first line: what was recorded, from where, and when the
// recording started. It carries no States field — recording is deliberately
// unfiltered, so the same log can be re-classified under different --states
// without re-recording (P6).
type Header struct {
	SchemaVersion  int          `json:"schema_version"`
	URL            string       `json:"url"` // the log's identity (§19.1 step 3)
	GrafanaVersion string       `json:"grafana_version"`
	StartedAt      time.Time    `json:"started_at"` // the record start (§7 validation)
	Rules          []LoggedRule `json:"rules"`      // THE alert set (§19.1 step 3)
}

// pausedAtStart reports, per rule UID, whether the rule was paused when the
// recording opened. That instant — and no other — is what `skipped` means
// (§12): a rule nobody was watching on purpose.
//
// It is the authority for `skipped` in BOTH modes, and the reason is that no
// other source knows the right moment. `check` re-resolves the definitions
// AFTER the window closed (§19.1 step 2), so Definition.IsPaused there
// describes the present, not the window: a rule that fired and was then
// paused would read as skipped, its firing would never be classified, and
// under --allow-paused the run would pass. The header cannot drift that way,
// because watch stamps it before the deploy step runs and single-step check
// stamps it from definitions resolved at the start of its own step.
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
	GrafanaNow time.Time `json:"grafana_now"` // the Date header — H4
	// SkewMS, SkewBoundMS and LatencyMS are milliseconds for JSONL
	// compactness ONLY. The pure layer never touches raw ms: it reads
	// Skew(), SkewBound() and Latency() below, which convert at the
	// (de)serialization boundary.
	SkewMS      int64 `json:"skew_ms"`
	SkewBoundMS int64 `json:"skew_bound_ms"`
	LatencyMS   int64 `json:"latency_ms"`
	// Found false means an authoritative 2xx in which this rule was absent
	// (§14.5) — never a transport failure, which P2 retried and never turns
	// into a Poll. P7 check 8 turns it into unobservable.
	Found bool `json:"found"`
	// State, Health and LastError are the raw rule-level strings, reporting
	// only and never classified (P1.2a).
	State     string `json:"state,omitempty"`
	Health    string `json:"health,omitempty"`
	LastError string `json:"last_error,omitempty"`
	// omitzero, not omitempty: a not-found poll (and a paused rule, §2.3) has
	// no evaluation time, and writing "0001-01-01T00:00:00Z" into an artifact
	// humans and jq read (§21.3) invites reading it as a real timestamp.
	LastEvaluation time.Time      `json:"last_evaluation,omitzero"`
	IsPaused       bool           `json:"is_paused"`
	Histogram      map[string]int `json:"histogram,omitempty"` // §4.9 — written, never analysed
	// Reasons counts this poll's non-empty instance reasons, e.g.
	// {"NoData":1091,"Error":14}; nil when none. Reporting-only, and the ONLY
	// place composite states stay visible: they are canonical normal (so they
	// are dropped from Abnormal) and `totals` never carries composite keys.
	//
	// The KEYS are raw reason strings and can be comma-joined composites
	// ("KeepLast, MissingSeries") — newer Grafana versions join several
	// reasons into one. So any consumer, P7 check 9's KeepLast note included,
	// must test membership across the keys with reasonNames and must NEVER
	// index a literal key: reasons["KeepLast"] misses every composite.
	Reasons map[string]int `json:"reasons,omitempty"`
	// Abnormal holds the instances whose CANONICAL state is not normal
	// (§4.6). "Normal (NoData)" and "Normal (Error)" are canonical normal and
	// are deliberately not retained here (P1.2a).
	Abnormal []Instance `json:"abnormal,omitempty"`
	// Cleared and Vanished are instance keys (§4.7): keys that left the
	// abnormal set, resolved against the SAME response — a clear and a
	// discontinuity are not the same fact (H2).
	Cleared  []string `json:"cleared,omitempty"`
	Vanished []string `json:"vanished,omitempty"`
}

// Skew is the signed clock skew of this poll (§16).
func (p Poll) Skew() time.Duration { return time.Duration(p.SkewMS) * time.Millisecond }

// SkewBound is the uncertainty on Skew — the tolerance every cross-domain
// comparison in P7 applies alongside it.
func (p Poll) SkewBound() time.Duration { return time.Duration(p.SkewBoundMS) * time.Millisecond }

// Latency is the wall time this poll's request took, feeding §5.2's budget check.
func (p Poll) Latency() time.Duration { return time.Duration(p.LatencyMS) * time.Millisecond }

// Reducer turns each Observation into the single Poll record that goes into
// the log. It holds the previous poll's abnormal instance keys per rule, which
// is all the state the transition markers need (§4.7).
//
// A Reducer is safe for concurrent use: watch polls a fleet of rules
// concurrently (P6) and every one of those goroutines reduces through the same
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
// several rules sharing one title (the known 2-way collision, §14.5), and
// picking the first would silently watch the wrong rule.
//
// The reduction (§4.6) keeps the rule-level fields, the raw totals histogram,
// the reason counts, and only the instances whose canonical state is not
// normal. That makes per-poll size independent of NORMAL cardinality — not of
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
	// the markers below must resolve a departed key against the same response
	// (H2), which is impossible from the abnormal subset alone.
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
			// absent case — H2's named bug.
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
// there (P6). Without the seed, an instance that is abnormal in the parent's
// observation and gone by the child's first poll produces no marker at all —
// it leaves the record as though it had never been bad, which is H2's
// fail-open reached through the handoff rather than through a reason string.
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

// stateRuleByUID picks one rule out of a state-endpoint response BY UID, and
// nil means the response is an authoritative "the rule is absent" (§14.5).
//
// Never by title: the ?rule_name= filter is a title filter, and a filtered
// response can carry several rules sharing one title (the known 2-way
// collision), so picking the first would silently watch the wrong rule. This
// is the single implementation of that selection for the package — Reduce
// above and the drain wait (check.go) both call it, for the same reason
// pollsForRule (classify.go) is shared between proveCoverage and classifyRule:
// two copies of a membership test are two chances for one to drift.
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
// than equality (P7 check 9 needs the same test for KeepLast).
func reasonNames(reason, want string) bool {
	for part := range strings.SplitSeq(reason, ",") {
		if strings.TrimSpace(part) == want {
			return true
		}
	}
	return false
}

// VerifyNormalInstancesVisible checks §3.2's assumption on a first
// observation: that the state endpoint really does return normal instances,
// not only the abnormal ones. If it ever stops doing so, the reduction's
// "keep the non-normal instances" becomes "keep everything the API happened to
// send" and the transition markers lose their ground truth — a silent
// fail-open. So this is verified at start, never assumed.
//
// The counts are summed over every totals key whose LOWERCASED name is
// "normal" or "inactive". Never index one literal key: the captured
// vocabulary is mixed across rules ({"alerting":445,"normal":2004} on one,
// {"firing":2,"inactive":363} on another) and its case has already drifted
// from the original recon. Composite states never appear in totals — Grafana
// counts a "Normal (NoData)" instance under normal.
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
				"the state endpoint no longer returns normal instances, which the §3.2 reduction depends on",
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

// Writer appends records to the JSONL log. It is append-only by construction
// (§8): O_APPEND|O_CREATE|O_WRONLY, never O_TRUNC, so no writer can ever
// destroy evidence a previous one recorded. An exclusive non-blocking flock
// makes a second writer fail immediately rather than interleave.
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
// contract a lie. In the P6 handoff the parent writes the header and the child
// only appends polls.
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

// Stop finishes recording in the fixed §4.4 order, which must not be
// reordered: let the in-flight write finish (the mutex), append the stopped
// sentinel, fsync, then release. Any other order can leave a log whose last
// durable byte is a sentinel that was never actually preceded by the polls it
// vouches for.
//
// Stop writes the sentinel with the recorder's OWN stop time and makes no
// comparison against `to` — watch never knows `to` or the transition grace.
// check does that comparison, after this writer has exited (§4.5).
//
// Calling Stop twice is a no-op: watch reaches it from both a signal handler
// and a defer, and a second sentinel would be indistinguishable from a second
// writer.
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
// hands the log to the detached child that will finish it (P6). A sentinel
// here would tell check the recording ended before the child had even
// started.
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

// ReadLogHeader reads ONLY line 1 and is the one read of a log that a writer
// may still hold. That is safe for exactly one line and for no other: the
// header is written once, by watch's parent, before any child appends a byte,
// the file is opened O_APPEND and never O_TRUNC (§8), so line 1 is complete
// and immutable for the whole life of the recording.
//
// It exists so check can fail closed EARLY (§19.1 steps 3-4): the log's
// identity, the rule set and the cadences are all knowable at the start, and
// discovering a wrong URL or an unresolvable rule after a ten-minute wait
// helps nobody. It is advisory only — the authoritative read is still ReadLog,
// once, after the writer has exited (§4.4 step 4), and check re-validates the
// identity against that header rather than trusting this one.
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
// Call this only after the writer has exited (§4.4 step 4). Reading a log a
// writer can still append to can only produce a shorter window than the one
// that was recorded.
//
// The parse rules are deliberately the crudest possible (§24.2): the header
// must be line 1 with a matching schema version, and ANY unparseable line —
// including the last one, and including a last line that follows a sentinel —
// is an error, full stop. No heuristics, no discarding an untidy tail: a
// truncated log is evidence that something killed the recorder, which is
// exactly what must not pass.
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
