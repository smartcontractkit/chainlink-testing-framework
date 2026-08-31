package gate

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"
)

// Every time literal in this file is UTC and built with time.Date, so it
// carries no monotonic reading and survives a JSON round trip byte-identical —
// which is what lets the round-trip tests below use reflect.DeepEqual on whole
// Poll values instead of comparing field by field.
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
			// Both composites are canonical normal (P1.2a): they must NOT be
			// retained as abnormal, and their reasons must still be counted.
			testInstance(StateNormal, "NoData", "c"),
			testInstance(StateNormal, "Error", "d"),
		},
	}

	p := NewReducer().Reduce("rule1", observation(testNow, rule))

	if !p.Found {
		t.Fatalf("Found = false, want true")
	}
	if len(p.Abnormal) != 1 || p.Abnormal[0].Labels["instance"] != "b" {
		t.Errorf("Abnormal = %+v, want only the firing instance b", p.Abnormal)
	}
	if want := map[string]int{"NoData": 1, "Error": 1}; !reflect.DeepEqual(p.Reasons, want) {
		t.Errorf("Reasons = %v, want %v", p.Reasons, want)
	}
	// The histogram is a verbatim copy of the response totals — raw keys, no
	// normalization (§4.9).
	if want := map[string]int{"alerting": 1, "normal": 2}; !reflect.DeepEqual(p.Histogram, want) {
		t.Errorf("Histogram = %v, want %v", p.Histogram, want)
	}
	// Rule-level state and health stay raw and unnormalized (P1.2a).
	if p.State != "firing" || p.Health != "ok" {
		t.Errorf("State/Health = %q/%q, want firing/ok", p.State, p.Health)
	}
	if p.Skew() != 1500*time.Millisecond || p.SkewBound() != 40*time.Millisecond || p.Latency() != 1800*time.Millisecond {
		t.Errorf("durations = %s/%s/%s, want 1.5s/40ms/1.8s", p.Skew(), p.SkewBound(), p.Latency())
	}
	if p.Reasons["MissingSeries"] != 0 {
		t.Errorf("unexpected MissingSeries count")
	}
}

// A filtered response can hold several rules sharing one title (the known
// 2-way collision, §14.5), so the reducer must select by UID.
func TestLogReduceSelectsRuleByUID(t *testing.T) {
	first := StateRule{UID: "ruleA", Title: "Same Title", Health: "ok", State: "inactive", LastEvaluation: testNow}
	second := StateRule{
		UID: "ruleB", Title: "Same Title", Health: "error", State: "firing", LastEvaluation: testNow,
		Instances: []Instance{testInstance(StateFiring, "", "x")},
	}

	p := NewReducer().Reduce("ruleB", observation(testNow, first, second))

	if p.Health != "error" || len(p.Abnormal) != 1 {
		t.Errorf("reduced the wrong rule: %+v", p)
	}
}

func TestLogReduceRuleAbsentIsAuthoritative(t *testing.T) {
	other := StateRule{UID: "other", Title: "Other", Health: "ok", State: "inactive", LastEvaluation: testNow}

	p := NewReducer().Reduce("rule1", observation(testNow, other))

	if p.Found {
		t.Errorf("Found = true, want false for a rule absent from an authoritative 2xx")
	}
	if p.RuleUID != "rule1" {
		t.Errorf("RuleUID = %q, want rule1 — an absent rule is still attributed", p.RuleUID)
	}
	// The heartbeat still exists: a not-found poll is evidence that Grafana
	// answered at this time, which the coverage proof reads.
	if !p.GrafanaNow.Equal(testNow) || p.Latency() == 0 {
		t.Errorf("absent-rule poll lost its timing evidence: %+v", p)
	}
	if p.Health != "" || p.Abnormal != nil {
		t.Errorf("absent-rule poll carries rule fields: %+v", p)
	}
}

// H2: an instance that leaves the abnormal set is resolved against the SAME
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
			if first.Cleared != nil || first.Vanished != nil {
				t.Fatalf("first poll produced markers with no previous poll: %+v", first)
			}

			next := StateRule{UID: "rule1", Health: "ok", State: "inactive", LastEvaluation: testNow, Instances: c.second}
			p := r.Reduce("rule1", observation(testNow.Add(30*time.Second), next))

			if !reflect.DeepEqual(p.Cleared, c.wantCleared) {
				t.Errorf("Cleared = %q, want %q", p.Cleared, c.wantCleared)
			}
			if !reflect.DeepEqual(p.Vanished, c.wantVanished) {
				t.Errorf("Vanished = %q, want %q", p.Vanished, c.wantVanished)
			}
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
	if absent.Vanished != nil || absent.Cleared != nil {
		t.Fatalf("an absent rule produced markers: %+v", absent)
	}

	back := StateRule{UID: "rule1", Health: "ok", State: "inactive", LastEvaluation: testNow}
	p := r.Reduce("rule1", observation(testNow.Add(60*time.Second), back))
	if len(p.Vanished) != 1 {
		t.Errorf("Vanished = %q, want the instance that disappeared across the absent poll", p.Vanished)
	}
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
	if len(p.Vanished) != 3 {
		t.Fatalf("Vanished = %q, want 3 keys", p.Vanished)
	}
	for i := 1; i < len(p.Vanished); i++ {
		if p.Vanished[i-1] > p.Vanished[i] {
			t.Errorf("Vanished is not sorted: %q", p.Vanished)
		}
	}
	// rule2's own abnormal set is untouched by rule1's transitions.
	q := r.Reduce("rule2", observation(testNow.Add(time.Minute), clearedOne, ruleTwo))
	if q.Cleared != nil || q.Vanished != nil {
		t.Errorf("rule2 picked up rule1's transitions: %+v", q)
	}
}

// §3.2: the reduction depends on the state endpoint returning normal instances.
// If it ever stops, that must fail loudly at start, never be assumed.
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
			if err != nil {
				t.Fatalf("ParseState: %v", err)
			}
			err = VerifyNormalInstancesVisible(rules)
			if c.wantError {
				if err == nil {
					t.Fatalf("VerifyNormalInstancesVisible: want an error, got nil")
				}
				if !strings.Contains(err.Error(), "§3.2") {
					t.Errorf("error does not name §3.2: %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("VerifyNormalInstancesVisible: unexpected error: %v", err)
			}
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
			if (err != nil) != c.wantError {
				t.Errorf("VerifyNormalInstancesVisible: error = %v, want error = %v", err, c.wantError)
			}
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
	if err != nil {
		t.Fatalf("DeriveTimingsFromLog: %v", err)
	}
	got := rt["rule1"]
	if got.pollEvery != 5*time.Second {
		t.Errorf("pollEvery = %s, want the header's 5s, not the default 150s", got.pollEvery)
	}
	// maxGap and healthGrace follow the recorded cadence; without this a 250s
	// hole in a log recorded at 5s would pass silently.
	if got.maxGap != 10*time.Second {
		t.Errorf("maxGap = %s, want 10s (2 x the recorded cadence)", got.maxGap)
	}
	if got.healthGrace != 300*time.Second {
		t.Errorf("healthGrace = %s, want 300s (max(maxGap, interval))", got.healthGrace)
	}
	// evalStaleAfter is a rule fact, so it stays 2 x intervalSeconds from the
	// definitions regardless of how often the gate polled.
	if got.evalStaleAfter != 600*time.Second {
		t.Errorf("evalStaleAfter = %s, want 600s from the definition's interval", got.evalStaleAfter)
	}

	// A log that cannot say how often it was written cannot have its coverage
	// proved, and neither can one naming a rule that no longer resolves.
	missingCadence := testHeader()
	missingCadence.Rules[0].PollEverySeconds = 0
	if _, _, err := DeriveTimingsFromLog(missingCadence, defs); err == nil {
		t.Errorf("a header with no recorded cadence was accepted")
	}
	if _, _, err := DeriveTimingsFromLog(testHeader(), nil); err == nil {
		t.Errorf("a header naming an unresolvable rule was accepted")
	}

	// A duplicated UID must not resolve last-one-wins: the slower duplicate
	// would widen maxGap, which is fail-open through log corruption alone.
	duplicated := testHeader()
	slower := duplicated.Rules[0]
	slower.PollEverySeconds = 600
	duplicated.Rules = append(duplicated.Rules, slower)
	if _, _, err := DeriveTimingsFromLog(duplicated, defs); err == nil {
		t.Errorf("a header naming one rule twice was accepted")
	}
}

// watch polls a fleet concurrently through one Reducer (P6), so the marker
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
		if p.Cleared != nil || p.Vanished != nil {
			t.Errorf("rule %s: markers after concurrent reduction: %+v", rule.UID, p)
		}
	}
}

// A not-found poll has no evaluation time, and the artifact is read by humans
// and jq (§21.3) — the zero time must not appear as though it were real.
func TestLogPollOmitsTheZeroEvaluationTime(t *testing.T) {
	absent := NewReducer().Reduce("rule1", observation(testNow))
	b, err := json.Marshal(pollRecord{Type: RecordPoll, Poll: absent})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if strings.Contains(string(b), "0001-01-01") {
		t.Errorf("a not-found poll wrote the zero time: %s", b)
	}
	if strings.Contains(string(b), "last_evaluation") {
		t.Errorf("a not-found poll wrote last_evaluation at all: %s", b)
	}

	// A real evaluation time still round-trips.
	found := StateRule{UID: "rule1", Health: "ok", State: "inactive", LastEvaluation: testNow}
	p := NewReducer().Reduce("rule1", observation(testNow, found))
	b, err = json.Marshal(pollRecord{Type: RecordPoll, Poll: p})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var back pollRecord
	if err := json.Unmarshal(b, &back); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !back.LastEvaluation.Equal(testNow) {
		t.Errorf("last_evaluation = %s, want %s", back.LastEvaluation, testNow)
	}
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
	if err != nil {
		t.Fatalf("NewWriter: %v", err)
	}
	return w, clock
}

func TestWriterReadLogRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "log.jsonl")
	w, clock := newTestWriter(t, path)

	h := testHeader()
	if err := w.WriteHeader(h); err != nil {
		t.Fatalf("WriteHeader: %v", err)
	}

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
		if err := w.WritePoll(p); err != nil {
			t.Fatalf("WritePoll: %v", err)
		}
	}

	clock.Advance(2 * time.Minute)
	if err := w.Stop(); err != nil {
		t.Fatalf("Stop: %v", err)
	}

	gotHeader, gotPolls, sentinel, err := ReadLog(path)
	if err != nil {
		t.Fatalf("ReadLog: %v", err)
	}
	h.SchemaVersion = LogSchemaVersion // WriteHeader stamps it
	if !reflect.DeepEqual(gotHeader, h) {
		t.Errorf("header round trip:\n got %+v\nwant %+v", gotHeader, h)
	}
	if !reflect.DeepEqual(gotPolls, want) {
		t.Errorf("poll round trip:\n got %+v\nwant %+v", gotPolls, want)
	}
	if sentinel == nil {
		t.Fatalf("sentinel is nil after Stop")
	}
	// Stop stamps the recorder's own stop time and makes no comparison
	// against `to` — watch never knows it (§4.5).
	if !sentinel.Equal(testNow.Add(2 * time.Minute)) {
		t.Errorf("sentinel = %s, want the writer's stop time %s", sentinel, testNow.Add(2*time.Minute))
	}
}

// §8: the log is append-only. A second run against the same path must never
// destroy the evidence the first one recorded.
func TestWriterAppendsAndNeverTruncates(t *testing.T) {
	path := filepath.Join(t.TempDir(), "log.jsonl")
	w, _ := newTestWriter(t, path)
	if err := w.WriteHeader(testHeader()); err != nil {
		t.Fatalf("WriteHeader: %v", err)
	}
	if err := w.WritePoll(Poll{RuleUID: "rule1", Found: true, GrafanaNow: testNow}); err != nil {
		t.Fatalf("WritePoll: %v", err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read: %v", err)
	}

	// The P6 handoff: the parent wrote the header and closed; the child
	// reopens the same path and appends without a second header.
	child, _ := newTestWriter(t, path)
	if err := child.WritePoll(Poll{RuleUID: "rule1", Found: true, GrafanaNow: testNow.Add(time.Minute)}); err != nil {
		t.Fatalf("child WritePoll: %v", err)
	}
	if err := child.Stop(); err != nil {
		t.Fatalf("child Stop: %v", err)
	}

	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if !strings.HasPrefix(string(after), string(before)) {
		t.Fatalf("reopening the log rewrote earlier records:\n%s", after)
	}
	_, polls, sentinel, err := ReadLog(path)
	if err != nil {
		t.Fatalf("ReadLog: %v", err)
	}
	if len(polls) != 2 || sentinel == nil {
		t.Errorf("got %d polls, sentinel %v; want 2 polls and a sentinel", len(polls), sentinel)
	}
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
		if err == nil {
			t.Fatalf("a second writer took the lock")
		}
		if !strings.Contains(err.Error(), "another writer") {
			t.Errorf("error does not name the conflict: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatalf("the second NewWriter blocked instead of failing immediately")
	}
}

func TestWriterHeaderRefusesANonEmptyLog(t *testing.T) {
	path := filepath.Join(t.TempDir(), "log.jsonl")
	w, _ := newTestWriter(t, path)
	if err := w.WriteHeader(testHeader()); err != nil {
		t.Fatalf("WriteHeader: %v", err)
	}
	if err := w.WriteHeader(testHeader()); err == nil {
		t.Fatalf("a second header was accepted")
	}
	if err := w.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	reopened, _ := newTestWriter(t, path)
	defer reopened.Close()
	if err := reopened.WriteHeader(testHeader()); err == nil {
		t.Fatalf("a header was accepted on a non-empty log")
	}
}

func TestSentinelStopIsIdempotentAndLast(t *testing.T) {
	path := filepath.Join(t.TempDir(), "log.jsonl")
	w, _ := newTestWriter(t, path)
	if err := w.WriteHeader(testHeader()); err != nil {
		t.Fatalf("WriteHeader: %v", err)
	}
	if err := w.Stop(); err != nil {
		t.Fatalf("Stop: %v", err)
	}
	// watch reaches Stop from both a signal handler and a defer; a second
	// sentinel would be indistinguishable from a second writer.
	if err := w.Stop(); err != nil {
		t.Errorf("second Stop: %v", err)
	}
	// Nothing may be appended after the sentinel — not even by the same writer.
	if err := w.WritePoll(Poll{RuleUID: "rule1"}); err == nil {
		t.Errorf("WritePoll after Stop was accepted")
	}

	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	lines := strings.Split(strings.TrimSuffix(string(b), "\n"), "\n")
	if len(lines) != 2 {
		t.Fatalf("got %d lines, want header + one sentinel:\n%s", len(lines), b)
	}
	if !strings.Contains(lines[1], `"type":"stopped"`) {
		t.Errorf("last line is not the sentinel: %s", lines[1])
	}
}

// Close is the parent's handoff path in P6: a sentinel there would tell check
// the recording ended before the child had even started.
func TestSentinelCloseWritesNone(t *testing.T) {
	path := filepath.Join(t.TempDir(), "log.jsonl")
	w, _ := newTestWriter(t, path)
	if err := w.WriteHeader(testHeader()); err != nil {
		t.Fatalf("WriteHeader: %v", err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	_, polls, sentinel, err := ReadLog(path)
	if err != nil {
		t.Fatalf("ReadLog: %v", err)
	}
	if sentinel != nil {
		t.Errorf("Close wrote a sentinel: %s", sentinel)
	}
	if polls != nil {
		t.Errorf("polls = %+v, want none", polls)
	}
}

// An unfinished recording reads cleanly with a nil sentinel — ReadLog reports
// the absence and P7 turns it into unobservable. It is never ReadLog's job to
// call that a failure, and never anyone's job to call it a pass.
func TestReadLogWithoutASentinel(t *testing.T) {
	path := filepath.Join(t.TempDir(), "log.jsonl")
	w, _ := newTestWriter(t, path)
	if err := w.WriteHeader(testHeader()); err != nil {
		t.Fatalf("WriteHeader: %v", err)
	}
	if err := w.WritePoll(Poll{RuleUID: "rule1", Found: true, GrafanaNow: testNow}); err != nil {
		t.Fatalf("WritePoll: %v", err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	_, polls, sentinel, err := ReadLog(path)
	if err != nil {
		t.Fatalf("ReadLog: %v", err)
	}
	if sentinel != nil || len(polls) != 1 {
		t.Errorf("got %d polls, sentinel %v; want 1 poll and no sentinel", len(polls), sentinel)
	}
}

// The read rules are deliberately the crudest possible (§24.2): any unparseable
// line is an error, full stop — including the last one, and including a last
// line that follows a sentinel.
func TestReadLogRejectsBadLogs(t *testing.T) {
	header := func(version int) string {
		h := testHeader()
		h.SchemaVersion = version
		b, err := json.Marshal(headerRecord{Type: RecordHeader, Header: h})
		if err != nil {
			t.Fatalf("marshal header: %v", err)
		}
		return string(b)
	}
	poll := func() string {
		b, err := json.Marshal(pollRecord{Type: RecordPoll, Poll: Poll{RuleUID: "rule1", Found: true, GrafanaNow: testNow}})
		if err != nil {
			t.Fatalf("marshal poll: %v", err)
		}
		return string(b)
	}
	sentinel := func() string {
		b, err := json.Marshal(stoppedRecord{Type: RecordStopped, At: testNow})
		if err != nil {
			t.Fatalf("marshal sentinel: %v", err)
		}
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
			if err := os.WriteFile(path, []byte(c.content), 0o600); err != nil {
				t.Fatalf("write: %v", err)
			}
			_, _, _, err := ReadLog(path)
			if err == nil {
				t.Fatalf("ReadLog: want an error, got nil")
			}
			if !strings.Contains(err.Error(), c.wantIn) {
				t.Errorf("error %q does not contain %q", err, c.wantIn)
			}
		})
	}
}

func TestReadLogMissingFile(t *testing.T) {
	_, _, _, err := ReadLog(filepath.Join(t.TempDir(), "absent.jsonl"))
	if err == nil {
		t.Fatalf("ReadLog on a missing log: want an error, got nil")
	}
}

// §22.3: per-poll log size must not grow across polls on a high-cardinality
// rule, and the one firing instance among 2446 must still be attributed by its
// labels. The reduction makes size independent of NORMAL cardinality — the
// firing instances are still stored, which is why a clear shrinks the record.
func TestLogSizeIsFlatAcrossPollsOnAHighCardinalityRule(t *testing.T) {
	body := synthesizeHighCardinalityState(t, 1, 2445)
	rules, err := ParseState(body)
	if err != nil {
		t.Fatalf("ParseState: %v", err)
	}
	if len(rules[0].Instances) != 2446 {
		t.Fatalf("got %d instances, want 2446", len(rules[0].Instances))
	}
	uid := rules[0].UID

	path := filepath.Join(t.TempDir(), "log.jsonl")
	w, _ := newTestWriter(t, path)
	if err := w.WriteHeader(testHeader()); err != nil {
		t.Fatalf("WriteHeader: %v", err)
	}

	r := NewReducer()
	var sizes []int64
	// Measure from the end of the header line, so sizes[0] is the first poll
	// record alone rather than the header plus it.
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	previous := info.Size()
	for i := range 5 {
		p := r.Reduce(uid, observation(testNow.Add(time.Duration(i)*30*time.Second), rules[0]))
		if len(p.Abnormal) != 1 {
			t.Fatalf("poll %d: Abnormal = %d instances, want the single firing one", i, len(p.Abnormal))
		}
		if got := p.Abnormal[0].Labels["instance"]; got != "alerting-0" {
			t.Fatalf("poll %d: the firing instance lost its identity: %q", i, got)
		}
		if err := w.WritePoll(p); err != nil {
			t.Fatalf("WritePoll: %v", err)
		}
		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("stat: %v", err)
		}
		sizes = append(sizes, info.Size()-previous)
		previous = info.Size()
	}

	for i := 1; i < len(sizes); i++ {
		if sizes[i] != sizes[0] {
			t.Errorf("per-poll size grew across polls: %v", sizes)
		}
	}
	// One firing instance among 2446 costs a few hundred bytes, against the
	// ~600 KB the unreduced response carries.
	if sizes[0] > 2048 {
		t.Errorf("per-poll size %d bytes is not a reduction of a %d-byte response", sizes[0], len(body))
	}

	// When the firing instance clears, the record collapses further and the
	// transition is still attributed.
	rules[0].Instances[0].State = StateNormal
	p := r.Reduce(uid, observation(testNow.Add(5*30*time.Second), rules[0]))
	if len(p.Cleared) != 1 || len(p.Abnormal) != 0 {
		t.Errorf("cleared poll = %+v, want exactly one cleared key and no abnormal instances", p)
	}
	if err := w.Stop(); err != nil {
		t.Fatalf("Stop: %v", err)
	}
}

// The log must stay readable by anything that reads JSONL, one flat object per
// line with its type tag — an uploaded artifact (§21.3) is read by humans and
// by jq, not only by ReadLog.
func TestLogRecordsAreFlatOneLineObjects(t *testing.T) {
	path := filepath.Join(t.TempDir(), "log.jsonl")
	w, _ := newTestWriter(t, path)
	if err := w.WriteHeader(testHeader()); err != nil {
		t.Fatalf("WriteHeader: %v", err)
	}
	if err := w.WritePoll(Poll{RuleUID: "rule1", Found: true, GrafanaNow: testNow}); err != nil {
		t.Fatalf("WritePoll: %v", err)
	}
	if err := w.Stop(); err != nil {
		t.Fatalf("Stop: %v", err)
	}

	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	lines := strings.Split(strings.TrimSuffix(string(b), "\n"), "\n")
	wantTypes := []RecordType{RecordHeader, RecordPoll, RecordStopped}
	if len(lines) != len(wantTypes) {
		t.Fatalf("got %d lines, want %d:\n%s", len(lines), len(wantTypes), b)
	}
	for i, line := range lines {
		var m map[string]json.RawMessage
		if err := json.Unmarshal([]byte(line), &m); err != nil {
			t.Fatalf("line %d is not one JSON object: %v", i+1, err)
		}
		var gotType RecordType
		if err := json.Unmarshal(m["type"], &gotType); err != nil {
			t.Fatalf("line %d has no type tag: %v", i+1, err)
		}
		if gotType != wantTypes[i] {
			t.Errorf("line %d type = %q, want %q", i+1, gotType, wantTypes[i])
		}
		if _, nested := m["header"]; nested {
			t.Errorf("line %d wraps its payload instead of being flat: %s", i+1, line)
		}
	}
}
