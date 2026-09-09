package gate

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// Every time literal in this file is UTC and built with time.Date, so it
// carries no monotonic reading and survives a JSON round trip byte-identical —
// which is what lets the round-trip tests below compare whole Poll values
// instead of comparing field by field.
var testNow = time.Date(2026, 8, 31, 9, 0, 0, 0, time.UTC)

func testInstance(state State, reason, instanceLabel string) Instance {
	return Instance{
		Labels:   map[string]string{"alertname": "Example", "instance": instanceLabel},
		State:    state,
		Reason:   reason,
		ActiveAt: testNow,
	}
}

// observation wraps rules into an Observation with plausible timing numbers,
// deliberately not round millisecond values so a lost conversion at the ms
// boundary shows up as a wrong number rather than a coincidentally equal one.
func observation(grafanaNow time.Time, rules ...StateRule) Observation {
	return Observation{
		Rules:      rules,
		GrafanaNow: grafanaNow,
		Skew:       1500 * time.Millisecond,
		SkewBound:  40 * time.Millisecond,
		Latency:    1800 * time.Millisecond,
	}
}

func TestLogReduceKeepsOnlyAbnormalInstances(t *testing.T) {
	rule := StateRule{
		UID: "rule1", Title: "Example", Folder: "F", Group: "G",
		Interval: time.Minute, State: "firing", Health: "ok",
		LastEvaluation: testNow, Totals: map[string]int{"alerting": 1, "normal": 2},
		Instances: []Instance{
			testInstance(StateNormal, "", "a"),
			testInstance(StateFiring, "", "b"),
			// Both composites are canonical normal: they must NOT be
			// retained as abnormal, and their reasons must still be counted.
			testInstance(StateNormal, "NoData", "c"),
			testInstance(StateNormal, "Error", "d"),
		},
	}

	p := NewReducer().Reduce("rule1", observation(testNow, rule))

	require.True(t, p.Found)
	require.Len(t, p.Abnormal, 1)
	require.Equal(t, "b", p.Abnormal[0].Labels["instance"])
	require.Equal(t, map[string]int{"NoData": 1, "Error": 1}, p.Reasons)
	// The histogram is a verbatim copy of the response totals — raw keys, no
	// normalization.
	require.Equal(t, map[string]int{"alerting": 1, "normal": 2}, p.Histogram)
	// Rule-level state and health stay raw and unnormalized.
	require.Equal(t, "firing", p.State)
	require.Equal(t, "ok", p.Health)
	require.Equal(t, 1500*time.Millisecond, p.Skew())
	require.Equal(t, 40*time.Millisecond, p.SkewBound())
	require.Equal(t, 1800*time.Millisecond, p.Latency())
	require.Zero(t, p.Reasons["MissingSeries"])
}

// A filtered response can hold several rules sharing one title, so the reducer
// must select by UID.
func TestLogReduceSelectsRuleByUID(t *testing.T) {
	first := StateRule{UID: "ruleA", Title: "Same Title", Health: "ok", State: "inactive", LastEvaluation: testNow}
	second := StateRule{
		UID: "ruleB", Title: "Same Title", Health: "error", State: "firing", LastEvaluation: testNow,
		Instances: []Instance{testInstance(StateFiring, "", "x")},
	}

	p := NewReducer().Reduce("ruleB", observation(testNow, first, second))

	require.Equal(t, "error", p.Health)
	require.Len(t, p.Abnormal, 1)
}

func TestLogReduceRuleAbsentIsAuthoritative(t *testing.T) {
	other := StateRule{UID: "other", Title: "Other", Health: "ok", State: "inactive", LastEvaluation: testNow}

	p := NewReducer().Reduce("rule1", observation(testNow, other))

	require.False(t, p.Found, "a rule absent from an authoritative 2xx")
	require.Equal(t, "rule1", p.RuleUID, "an absent rule is still attributed")
	// The heartbeat still exists: a not-found poll is evidence that Grafana
	// answered at this time, which the coverage proof reads.
	require.True(t, p.GrafanaNow.Equal(testNow))
	require.NotZero(t, p.Latency())
	require.Empty(t, p.Health)
	require.Nil(t, p.Abnormal)
}

// An instance that leaves the abnormal set is resolved against the SAME
// response, and MissingSeries is a vanish, never a recovery.
func TestTransitionMarkersClearedVersusVanished(t *testing.T) {
	badKey := instanceKey(testInstance(StateFiring, "", "b").Labels)

	cases := []struct {
		name         string
		second       []Instance
		wantCleared  []string
		wantVanished []string
	}{
		{
			name:        "present as canonical normal is a clear",
			second:      []Instance{testInstance(StateNormal, "", "b")},
			wantCleared: []string{badKey},
		},
		{
			name:         "fully absent is a discontinuity",
			second:       nil,
			wantVanished: []string{badKey},
		},
		{
			name:         "Normal (MissingSeries) is the vanish in disguise",
			second:       []Instance{testInstance(StateNormal, "MissingSeries", "b")},
			wantVanished: []string{badKey},
		},
		{
			name:         "a comma-joined reason naming MissingSeries still vanishes",
			second:       []Instance{testInstance(StateNormal, "KeepLast, MissingSeries", "b")},
			wantVanished: []string{badKey},
		},
		{
			name:        "an unrelated comma-joined reason still clears",
			second:      []Instance{testInstance(StateNormal, "KeepLast, Updated", "b")},
			wantCleared: []string{badKey},
		},
		{
			name:   "still abnormal is neither",
			second: []Instance{testInstance(StateFiring, "", "b")},
		},
		{
			name:         "abnormal under a different state, then gone, still vanishes",
			second:       []Instance{testInstance(StateNormal, "", "unrelated")},
			wantVanished: []string{badKey},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := NewReducer()
			firing := StateRule{
				UID: "rule1", Health: "ok", State: "firing", LastEvaluation: testNow,
				Instances: []Instance{testInstance(StateFiring, "", "b")},
			}
			first := r.Reduce("rule1", observation(testNow, firing))
			require.Nil(t, first.Cleared, "first poll produced cleared markers with no previous poll")
			require.Nil(t, first.Vanished, "first poll produced vanished markers with no previous poll")

			next := StateRule{UID: "rule1", Health: "ok", State: "inactive", LastEvaluation: testNow, Instances: c.second}
			p := r.Reduce("rule1", observation(testNow.Add(30*time.Second), next))

			require.Equal(t, c.wantCleared, p.Cleared)
			require.Equal(t, c.wantVanished, p.Vanished)
		})
	}
}

// A departed key must be resolved against the response the reducer is holding,
// not against a later one — so a rule that goes absent and comes back with the
// instance missing still reports the vanish rather than losing it.
func TestTransitionMarkersSurviveAnAbsentPoll(t *testing.T) {
	r := NewReducer()
	firing := StateRule{
		UID: "rule1", Health: "ok", State: "firing", LastEvaluation: testNow,
		Instances: []Instance{testInstance(StateFiring, "", "b")},
	}
	r.Reduce("rule1", observation(testNow, firing))

	absent := r.Reduce("rule1", observation(testNow.Add(30*time.Second)))
	require.Nil(t, absent.Vanished, "an absent rule produced vanished markers")
	require.Nil(t, absent.Cleared, "an absent rule produced cleared markers")

	back := StateRule{UID: "rule1", Health: "ok", State: "inactive", LastEvaluation: testNow}
	p := r.Reduce("rule1", observation(testNow.Add(60*time.Second), back))
	require.Len(t, p.Vanished, 1, "want the instance that disappeared across the absent poll")
}

func TestTransitionMarkersAreSortedAndPerRule(t *testing.T) {
	r := NewReducer()
	ruleOne := StateRule{
		UID: "rule1", Health: "ok", State: "firing", LastEvaluation: testNow,
		Instances: []Instance{
			testInstance(StateFiring, "", "z"),
			testInstance(StateFiring, "", "a"),
			testInstance(StateFiring, "", "m"),
		},
	}
	ruleTwo := StateRule{
		UID: "rule2", Health: "ok", State: "firing", LastEvaluation: testNow,
		Instances: []Instance{testInstance(StateFiring, "", "q")},
	}
	r.Reduce("rule1", observation(testNow, ruleOne, ruleTwo))
	r.Reduce("rule2", observation(testNow, ruleOne, ruleTwo))

	clearedOne := StateRule{UID: "rule1", Health: "ok", State: "inactive", LastEvaluation: testNow}
	p := r.Reduce("rule1", observation(testNow.Add(time.Minute), clearedOne, ruleTwo))
	require.Len(t, p.Vanished, 3)
	for i := 1; i < len(p.Vanished); i++ {
		require.Less(t, p.Vanished[i-1], p.Vanished[i], "Vanished is not sorted: %q", p.Vanished)
	}
	// rule2's own abnormal set is untouched by rule1's transitions.
	q := r.Reduce("rule2", observation(testNow.Add(time.Minute), clearedOne, ruleTwo))
	require.Nil(t, q.Cleared, "rule2 picked up rule1's transitions")
	require.Nil(t, q.Vanished, "rule2 picked up rule1's transitions")
}

// The reduction depends on the state endpoint returning normal instances. If it
// ever stops, that must fail loudly at start, never be assumed.
func TestLogVerifyNormalInstancesVisible(t *testing.T) {
	cases := []struct {
		fixture   string
		wantError bool
	}{
		{"state_one_instance.json", false},
		{"state_reason_composite.json", false}, // composites plus one plain Normal
		{"state_paused.json", false},           // no totals, no instances
		{"state_missing_optional.json", false}, // no totals key at all
		{"state_health_error.json", false},     // totals {"error":1} claims no normal
		{"state_only_active_instances.json", true},
	}

	for _, c := range cases {
		t.Run(c.fixture, func(t *testing.T) {
			rules, err := ParseState(readFixture(t, c.fixture))
			require.NoError(t, err)
			err = VerifyNormalInstancesVisible(rules)
			if c.wantError {
				require.Error(t, err)
				require.Contains(t, err.Error(), "no longer returns normal instances")
				return
			}
			require.NoError(t, err)
		})
	}
}

// The totals vocabulary is mixed and its case has drifted, so the check sums
// every key that lowercases to normal or inactive rather than indexing one
// literal key.
func TestLogVerifyNormalInstancesVisibleVocabularies(t *testing.T) {
	cases := []struct {
		name      string
		totals    map[string]int
		wantError bool
	}{
		{"lowercase normal", map[string]int{"alerting": 1, "normal": 4}, true},
		{"capitalized Normal", map[string]int{"Alerting": 1, "Normal": 4}, true},
		{"rule vocabulary inactive", map[string]int{"firing": 2, "inactive": 363}, true},
		{"no normal claimed", map[string]int{"alerting": 1}, false},
		{"nil totals", nil, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			rules := []StateRule{{
				UID: "rule1", Title: "Example", Totals: c.totals,
				Instances: []Instance{testInstance(StateFiring, "", "b")},
			}}
			err := VerifyNormalInstancesVisible(rules)
			require.Equal(t, c.wantError, err != nil)
		})
	}
}

// The two authorities: the header owns the recording facts (the cadence
// actually used), the ruler API owns the rule facts. Mixing them up is
// fail-open in the faster-override direction, so this pins both.
func TestLogModeCadenceComesFromTheHeader(t *testing.T) {
	defs := []Definition{{UID: "rule1", Title: "Example", IntervalSeconds: 300, For: time.Minute}}

	h := testHeader()
	h.Rules[0].IntervalSeconds = 300
	h.Rules[0].PollEverySeconds = 5 // an operator override far tighter than the default 150s

	rt, _, err := DeriveTimingsFromLog(h, defs)
	require.NoError(t, err)
	got := rt["rule1"]
	require.Equal(t, 5*time.Second, got.pollEvery, "the header's 5s, not the default 150s")
	// maxGap and healthGrace follow the recorded cadence; without this a 250s
	// hole in a log recorded at 5s would pass silently.
	require.Equal(t, 10*time.Second, got.maxGap)
	require.Equal(t, 300*time.Second, got.healthGrace)
	// evalStaleAfter is a rule fact, so it stays 2 x intervalSeconds from the
	// definitions regardless of how often the gate polled.
	require.Equal(t, 600*time.Second, got.evalStaleAfter)

	// A log that cannot say how often it was written cannot have its coverage
	// proved, and neither can one naming a rule that no longer resolves.
	missingCadence := testHeader()
	missingCadence.Rules[0].PollEverySeconds = 0
	require.Error(t, func() error { _, _, err := DeriveTimingsFromLog(missingCadence, defs); return err }(),
		"a header with no recorded cadence was accepted")
	require.Error(t, func() error { _, _, err := DeriveTimingsFromLog(testHeader(), nil); return err }(),
		"a header naming an unresolvable rule was accepted")

	// A duplicated UID must not resolve last-one-wins: the slower duplicate
	// would widen maxGap, which is fail-open through log corruption alone.
	duplicated := testHeader()
	slower := duplicated.Rules[0]
	slower.PollEverySeconds = 600
	duplicated.Rules = append(duplicated.Rules, slower)
	require.Error(t, func() error { _, _, err := DeriveTimingsFromLog(duplicated, defs); return err }(),
		"a header naming one rule twice was accepted")
}

// watch polls a fleet concurrently through one Reducer, so the marker
// state it holds per rule must be safe under -race — a latent data race here
// surfaces as a wrong transition, which is the one thing markers exist to get
// right.
func TestLogReduceIsSafeForConcurrentUse(t *testing.T) {
	r := NewReducer()
	rules := make([]StateRule, 0, 8)
	for i := range 8 {
		rules = append(rules, StateRule{
			UID: fmt.Sprintf("rule%d", i), Health: "ok", State: "firing", LastEvaluation: testNow,
			Instances: []Instance{testInstance(StateFiring, "", "b")},
		})
	}
	obs := observation(testNow, rules...)

	var wg sync.WaitGroup
	for range 3 {
		for _, rule := range rules {
			wg.Go(func() { r.Reduce(rule.UID, obs) })
		}
	}
	wg.Wait()

	// Each rule's abnormal instance never left, so no round may invent a
	// transition — the concurrency must not corrupt the per-rule state either.
	for _, rule := range rules {
		p := r.Reduce(rule.UID, obs)
		require.Nilf(t, p.Cleared, "rule %s: cleared markers after concurrent reduction", rule.UID)
		require.Nilf(t, p.Vanished, "rule %s: vanished markers after concurrent reduction", rule.UID)
	}
}

// A not-found poll has no evaluation time, and the artifact is read by humans
// and jq — the zero time must not appear as though it were real.
func TestLogPollOmitsTheZeroEvaluationTime(t *testing.T) {
	absent := NewReducer().Reduce("rule1", observation(testNow))
	b, err := json.Marshal(pollRecord{Type: RecordPoll, Poll: absent})
	require.NoError(t, err)
	require.NotContains(t, string(b), "0001-01-01")
	require.NotContains(t, string(b), "last_evaluation")

	// A real evaluation time still round-trips.
	found := StateRule{UID: "rule1", Health: "ok", State: "inactive", LastEvaluation: testNow}
	p := NewReducer().Reduce("rule1", observation(testNow, found))
	b, err = json.Marshal(pollRecord{Type: RecordPoll, Poll: p})
	require.NoError(t, err)
	var back pollRecord
	require.NoError(t, json.Unmarshal(b, &back))
	require.True(t, back.LastEvaluation.Equal(testNow))
}

func testHeader() Header {
	return Header{
		URL:            "https://grafana.example.com",
		GrafanaVersion: "13.1.0",
		StartedAt:      testNow,
		Rules: []LoggedRule{{
			UID: "rule1", Title: "Example", Folder: "F", Group: "G",
			ForSeconds: 300, IntervalSeconds: 60, NoDataState: "OK", ExecErrState: "OK",
			PollEverySeconds: 30,
		}},
	}
}

func newTestWriter(t *testing.T, path string) (*Writer, *fakeClock) {
	t.Helper()
	clock := newFakeClock(testNow)
	w, err := NewWriter(path, clock)
	require.NoError(t, err)
	return w, clock
}

func TestWriterReadLogRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "log.jsonl")
	w, clock := newTestWriter(t, path)

	h := testHeader()
	require.NoError(t, w.WriteHeader(h))

	r := NewReducer()
	firing := StateRule{
		UID: "rule1", Health: "ok", State: "firing", LastError: "", LastEvaluation: testNow,
		Totals:    map[string]int{"alerting": 1},
		Instances: []Instance{testInstance(StateFiring, "", "b")},
	}
	cleared := StateRule{UID: "rule1", Health: "ok", State: "inactive", LastEvaluation: testNow.Add(time.Minute)}
	want := []Poll{
		r.Reduce("rule1", observation(testNow, firing)),
		r.Reduce("rule1", observation(testNow.Add(time.Minute), cleared)),
	}
	for _, p := range want {
		require.NoError(t, w.WritePoll(p))
	}

	clock.Advance(2 * time.Minute)
	require.NoError(t, w.Stop())

	gotHeader, gotPolls, sentinel, err := ReadLog(path)
	require.NoError(t, err)
	h.SchemaVersion = LogSchemaVersion // WriteHeader stamps it
	require.Equal(t, h, gotHeader, "header round trip")
	require.Equal(t, want, gotPolls, "poll round trip")
	require.NotNil(t, sentinel, "sentinel is nil after Stop")
	// Stop stamps the recorder's own stop time and makes no comparison
	// against `to` — watch never knows it.
	require.True(t, sentinel.Equal(testNow.Add(2*time.Minute)))
}

// The log is append-only. A second run against the same path must never
// destroy the evidence the first one recorded.
func TestWriterAppendsAndNeverTruncates(t *testing.T) {
	path := filepath.Join(t.TempDir(), "log.jsonl")
	w, _ := newTestWriter(t, path)
	require.NoError(t, w.WriteHeader(testHeader()))
	require.NoError(t, w.WritePoll(Poll{RuleUID: "rule1", Found: true, GrafanaNow: testNow}))
	require.NoError(t, w.Close())
	before, err := os.ReadFile(path)
	require.NoError(t, err)

	// The handoff: the parent wrote the header and closed; the child
	// reopens the same path and appends without a second header.
	child, _ := newTestWriter(t, path)
	require.NoError(t, child.WritePoll(Poll{RuleUID: "rule1", Found: true, GrafanaNow: testNow.Add(time.Minute)}))
	require.NoError(t, child.Stop())

	after, err := os.ReadFile(path)
	require.NoError(t, err)
	require.True(t, strings.HasPrefix(string(after), string(before)), "reopening the log rewrote earlier records")
	_, polls, sentinel, err := ReadLog(path)
	require.NoError(t, err)
	require.Len(t, polls, 2)
	require.NotNil(t, sentinel)
}

// Two recorders on one log means one of them is recording a window nobody
// will classify, so the second writer fails immediately — it never blocks.
func TestWriterSecondWriterFails(t *testing.T) {
	path := filepath.Join(t.TempDir(), "log.jsonl")
	first, _ := newTestWriter(t, path)
	defer first.Close()

	done := make(chan error, 1)
	go func() {
		_, err := NewWriter(path, newFakeClock(testNow))
		done <- err
	}()

	select {
	case err := <-done:
		require.Error(t, err, "a second writer took the lock")
		require.Contains(t, err.Error(), "another writer")
	case <-time.After(5 * time.Second):
		require.Fail(t, "the second NewWriter blocked instead of failing immediately")
	}
}

func TestWriterHeaderRefusesANonEmptyLog(t *testing.T) {
	path := filepath.Join(t.TempDir(), "log.jsonl")
	w, _ := newTestWriter(t, path)
	require.NoError(t, w.WriteHeader(testHeader()))
	require.Error(t, w.WriteHeader(testHeader()), "a second header was accepted")
	require.NoError(t, w.Close())

	reopened, _ := newTestWriter(t, path)
	defer reopened.Close()
	require.Error(t, reopened.WriteHeader(testHeader()), "a header was accepted on a non-empty log")
}

func TestSentinelStopIsIdempotentAndLast(t *testing.T) {
	path := filepath.Join(t.TempDir(), "log.jsonl")
	w, _ := newTestWriter(t, path)
	require.NoError(t, w.WriteHeader(testHeader()))
	require.NoError(t, w.Stop())
	// watch reaches Stop from both a signal handler and a defer; a second
	// sentinel would be indistinguishable from a second writer.
	require.NoError(t, w.Stop())
	// Nothing may be appended after the sentinel — not even by the same writer.
	require.Error(t, w.WritePoll(Poll{RuleUID: "rule1"}), "WritePoll after Stop was accepted")

	b, err := os.ReadFile(path)
	require.NoError(t, err)
	lines := strings.Split(strings.TrimSuffix(string(b), "\n"), "\n")
	require.Len(t, lines, 2, "header + one sentinel")
	require.Contains(t, lines[1], `"type":"stopped"`)
}

// Close is the parent's handoff path: a sentinel there would tell check the
// recording ended before the child had even started.
func TestSentinelCloseWritesNone(t *testing.T) {
	path := filepath.Join(t.TempDir(), "log.jsonl")
	w, _ := newTestWriter(t, path)
	require.NoError(t, w.WriteHeader(testHeader()))
	require.NoError(t, w.Close())

	_, polls, sentinel, err := ReadLog(path)
	require.NoError(t, err)
	require.Nil(t, sentinel, "Close wrote a sentinel")
	require.Nil(t, polls)
}

// An unfinished recording reads cleanly with a nil sentinel — ReadLog reports
// the absence and the coverage proof turns it into unobservable. It is never
// ReadLog's job to call that a failure, and never anyone's job to call it a
// pass.
func TestReadLogWithoutASentinel(t *testing.T) {
	path := filepath.Join(t.TempDir(), "log.jsonl")
	w, _ := newTestWriter(t, path)
	require.NoError(t, w.WriteHeader(testHeader()))
	require.NoError(t, w.WritePoll(Poll{RuleUID: "rule1", Found: true, GrafanaNow: testNow}))
	require.NoError(t, w.Close())

	_, polls, sentinel, err := ReadLog(path)
	require.NoError(t, err)
	require.Nil(t, sentinel)
	require.Len(t, polls, 1)
}

// The read rules are deliberately the crudest possible: any unparseable
// line is an error, full stop — including the last one, and including a last
// line that follows a sentinel.
func TestReadLogRejectsBadLogs(t *testing.T) {
	header := func(version int) string {
		h := testHeader()
		h.SchemaVersion = version
		b, err := json.Marshal(headerRecord{Type: RecordHeader, Header: h})
		require.NoError(t, err, "marshal header")
		return string(b)
	}
	poll := func() string {
		b, err := json.Marshal(pollRecord{Type: RecordPoll, Poll: Poll{RuleUID: "rule1", Found: true, GrafanaNow: testNow}})
		require.NoError(t, err, "marshal poll")
		return string(b)
	}
	sentinel := func() string {
		b, err := json.Marshal(stoppedRecord{Type: RecordStopped, At: testNow})
		require.NoError(t, err, "marshal sentinel")
		return string(b)
	}

	cases := []struct {
		name    string
		content string
		wantIn  string
	}{
		{"empty file", "", "empty"},
		{"no header", poll() + "\n", "header must be line 1"},
		{"header not first", poll() + "\n" + header(LogSchemaVersion) + "\n", "header must be line 1"},
		{"second header", header(LogSchemaVersion) + "\n" + header(LogSchemaVersion) + "\n", "appear only once"},
		{"wrong schema version", header(2) + "\n" + poll() + "\n", "schema version"},
		{
			name:    "unparseable last line",
			content: header(LogSchemaVersion) + "\n" + poll() + "\n" + `{"type":"poll","rule_ui`,
			wantIn:  "unparseable",
		},
		{
			// A preceding sentinel makes no difference: a truncated tail is
			// evidence that something killed the recorder.
			name:    "unparseable line after the sentinel",
			content: header(LogSchemaVersion) + "\n" + sentinel() + "\n" + `{"type":"pol`,
			wantIn:  "unparseable",
		},
		{
			name:    "unparseable middle line",
			content: header(LogSchemaVersion) + "\n" + `{"type":` + "\n" + poll() + "\n",
			wantIn:  "unparseable",
		},
		{
			name:    "empty middle line",
			content: header(LogSchemaVersion) + "\n\n" + poll() + "\n",
			wantIn:  "unparseable",
		},
		{
			name:    "a record after the sentinel",
			content: header(LogSchemaVersion) + "\n" + sentinel() + "\n" + poll() + "\n",
			wantIn:  "second writer",
		},
		{
			name:    "unknown record type",
			content: header(LogSchemaVersion) + "\n" + `{"type":"heartbeat"}` + "\n",
			wantIn:  "unknown record type",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "log.jsonl")
			require.NoError(t, os.WriteFile(path, []byte(c.content), 0o600))
			_, _, _, err := ReadLog(path)
			require.Error(t, err)
			require.Contains(t, err.Error(), c.wantIn)
		})
	}
}

func TestReadLogMissingFile(t *testing.T) {
	_, _, _, err := ReadLog(filepath.Join(t.TempDir(), "absent.jsonl"))
	require.Error(t, err)
}

// Per-poll log size must not grow across polls on a high-cardinality
// rule, and the one firing instance among 2446 must still be attributed by its
// labels. The reduction makes size independent of NORMAL cardinality — the
// firing instances are still stored, which is why a clear shrinks the record.
func TestLogSizeIsFlatAcrossPollsOnAHighCardinalityRule(t *testing.T) {
	body := synthesizeHighCardinalityState(t, 1, 2445)
	rules, err := ParseState(body)
	require.NoError(t, err)
	require.Len(t, rules[0].Instances, 2446)
	uid := rules[0].UID

	path := filepath.Join(t.TempDir(), "log.jsonl")
	w, _ := newTestWriter(t, path)
	require.NoError(t, w.WriteHeader(testHeader()))

	r := NewReducer()
	var sizes []int64
	// Measure from the end of the header line, so sizes[0] is the first poll
	// record alone rather than the header plus it.
	info, err := os.Stat(path)
	require.NoError(t, err)
	previous := info.Size()
	for i := range 5 {
		p := r.Reduce(uid, observation(testNow.Add(time.Duration(i)*30*time.Second), rules[0]))
		require.Lenf(t, p.Abnormal, 1, "poll %d", i)
		require.Equalf(t, "alerting-0", p.Abnormal[0].Labels["instance"], "poll %d: the firing instance lost its identity", i)
		require.NoError(t, w.WritePoll(p))
		info, err := os.Stat(path)
		require.NoError(t, err)
		sizes = append(sizes, info.Size()-previous)
		previous = info.Size()
	}

	for i := 1; i < len(sizes); i++ {
		require.Equal(t, sizes[0], sizes[i], "per-poll size grew across polls: %v", sizes)
	}
	// One firing instance among 2446 costs a few hundred bytes, against the
	// ~600 KB the unreduced response carries.
	require.LessOrEqual(t, sizes[0], int64(2048))

	// When the firing instance clears, the record collapses further and the
	// transition is still attributed.
	rules[0].Instances[0].State = StateNormal
	p := r.Reduce(uid, observation(testNow.Add(5*30*time.Second), rules[0]))
	require.Len(t, p.Cleared, 1)
	require.Empty(t, p.Abnormal)
	require.NoError(t, w.Stop())
}

// The log must stay readable by anything that reads JSONL, one flat object per
// line with its type tag — an uploaded artifact is read by humans and by jq,
// not only by ReadLog.
func TestLogRecordsAreFlatOneLineObjects(t *testing.T) {
	path := filepath.Join(t.TempDir(), "log.jsonl")
	w, _ := newTestWriter(t, path)
	require.NoError(t, w.WriteHeader(testHeader()))
	require.NoError(t, w.WritePoll(Poll{RuleUID: "rule1", Found: true, GrafanaNow: testNow}))
	require.NoError(t, w.Stop())

	b, err := os.ReadFile(path)
	require.NoError(t, err)
	lines := strings.Split(strings.TrimSuffix(string(b), "\n"), "\n")
	wantTypes := []RecordType{RecordHeader, RecordPoll, RecordStopped}
	require.Len(t, lines, len(wantTypes))
	for i, line := range lines {
		var m map[string]json.RawMessage
		require.NoErrorf(t, json.Unmarshal([]byte(line), &m), "line %d is not one JSON object", i+1)
		var gotType RecordType
		require.NoErrorf(t, json.Unmarshal(m["type"], &gotType), "line %d has no type tag", i+1)
		require.Equalf(t, wantTypes[i], gotType, "line %d type", i+1)
		_, nested := m["header"]
		require.Falsef(t, nested, "line %d wraps its payload instead of being flat", i+1)
	}
}
