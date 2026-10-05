package gate

import (
	"context"
	"fmt"
	"slices"
	"sync"
	"time"
)

// fakeClock is a manually-advanced Clock — no test in this package sleeps on
// real time. It is goroutine-safe (a concurrent fleet under -race must not trip
// on the double itself), but After always fires immediately, regardless of the
// requested duration or whether Advance was ever called. That is enough for the
// retry/backoff tests, which only need to avoid a real sleep. It is NOT enough
// for a test that must prove a wait did not fire early — e.g. asserting Due()
// does not return a rule before its next-due time. Use virtualClock below for
// that: it makes a wait and the passage of time the same event.
type fakeClock struct {
	mu  sync.Mutex
	now time.Time
}

func newFakeClock(now time.Time) *fakeClock { return &fakeClock{now: now} }

func (c *fakeClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.now
}

func (c *fakeClock) Advance(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.now = c.now.Add(d)
}

func (c *fakeClock) After(d time.Duration) <-chan time.Time {
	c.mu.Lock()
	fireAt := c.now.Add(d)
	c.mu.Unlock()
	ch := make(chan time.Time, 1)
	ch <- fireAt
	return ch
}

// virtualClock is a Clock in which time moves only when something waits for
// it: After(d) jumps Now() forward by d and fires at once. That makes a
// recorder-loop test both instant and exact — a loop that waits for its next
// scheduled poll gets that poll's time, never an early or a late wake — and it
// terminates, which a clock whose After fires without advancing Now does not
// (the loop would spin forever on a rule that never comes due).
//
// It is goroutine-safe, but a test that advances time from two goroutines gets
// what it deserves: use it from the loop under test only.
type virtualClock struct {
	mu  sync.Mutex
	now time.Time
}

func newVirtualClock(now time.Time) *virtualClock { return &virtualClock{now: now} }

func (c *virtualClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.now
}

// Advance moves the clock without a wait. It exists for a test that must model
// time passing inside a non-waiting section — the measurement pass of a
// single-step check — while collectUntil below still advances through After.
func (c *virtualClock) Advance(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.now = c.now.Add(d)
}

func (c *virtualClock) After(d time.Duration) <-chan time.Time {
	c.mu.Lock()
	if d > 0 {
		c.now = c.now.Add(d)
	}
	fireAt := c.now
	c.mu.Unlock()
	ch := make(chan time.Time, 1)
	ch <- fireAt
	return ch
}

// steppingClock advances by a fixed step on every Now() call, so a test can
// assert exact latency/skew-bound arithmetic (doRequest's three clock reads
// per attempt) without depending on real wall-clock timing.
type steppingClock struct {
	mu   sync.Mutex
	now  time.Time
	step time.Duration
}

func (c *steppingClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	t := c.now
	c.now = c.now.Add(c.step)
	return t
}

func (c *steppingClock) After(d time.Duration) <-chan time.Time {
	c.mu.Lock()
	fireAt := c.now.Add(d)
	c.mu.Unlock()
	ch := make(chan time.Time, 1)
	ch <- fireAt
	return ch
}

// scriptedObservation is one canned (Observation, error) pair a fakeSource
// returns from RuleState, in the order scripted.
type scriptedObservation struct {
	obs Observation
	err error
}

// fakeSource is a scripted Source with no HTTP, goroutine-safe so a test that
// polls several rules concurrently can share one instance across goroutines
// without tripping -race. A test that needs it to behave like a live server
// under concurrent load beyond simple locking should verify that assumption
// rather than take this comment's word for it.
type fakeSource struct {
	mu sync.Mutex

	version    string
	versionErr error

	defs    []Definition
	defsErr error

	// ruleSources is what DiscoverRuleSources returns; dsDefs maps a datasource
	// UID to the definitions DatasourceDefinitions returns.
	ruleSources []RuleSource
	sourcesErr  error
	dsDefs      map[string][]Definition

	// states maps a rule title (Grafana) or key (datasource) to a queue of
	// scripted results, popped one per call to RuleState. Once the queue is down
	// to its last entry, that entry repeats — so a test can script the
	// interesting transitions and let a long collection loop settle into steady
	// state without scripting every single poll.
	states map[string][]scriptedObservation
}

func newFakeSource() *fakeSource {
	return &fakeSource{states: make(map[string][]scriptedObservation)}
}

func (f *fakeSource) Version(_ context.Context) (string, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.version, f.versionErr
}

func (f *fakeSource) GrafanaDefinitions(_ context.Context) ([]Definition, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	// Defensive copy: Definition is a value type, so Clone copies the
	// full slice contents, not just the header — callers are free to mutate
	// what they got back without racing or corrupting later reads.
	return slices.Clone(f.defs), f.defsErr
}

func (f *fakeSource) DiscoverRuleSources(_ context.Context) ([]RuleSource, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return slices.Clone(f.ruleSources), f.sourcesErr
}

func (f *fakeSource) DatasourceDefinitions(_ context.Context, src RuleSource, names []string) ([]Definition, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	defs := slices.Clone(f.dsDefs[src.UID])
	if len(names) == 0 {
		return defs, nil
	}
	want := make(map[string]bool, len(names))
	for _, n := range names {
		want[n] = true
	}
	var out []Definition
	for _, d := range defs {
		if want[d.Title] {
			out = append(out, d)
		}
	}
	return out, nil
}

func (f *fakeSource) RuleState(_ context.Context, ref RuleRef) (Observation, error) {
	key := ref.Title
	if ref.Kind == KindDatasourceManaged {
		key = ref.Key
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	q := f.states[key]
	if len(q) == 0 {
		return Observation{}, fmt.Errorf("fakeSource: no scripted response for %q", key)
	}
	next := q[0]
	if len(q) > 1 {
		f.states[key] = q[1:]
	}
	// Defensive copy of the shared Rules slice so a caller mutating the
	// returned Observation can't corrupt the scripted state other calls read.
	next.obs.Rules = slices.Clone(next.obs.Rules)
	return next.obs, next.err
}

// script appends one scripted (Observation, error) pair to be returned, in
// order, by RuleState(ctx, ref) for the Grafana rule with this title.
func (f *fakeSource) script(title string, obs Observation, err error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.states[title] = append(f.states[title], scriptedObservation{obs: obs, err: err})
}

// scriptKey is script's datasource counterpart, keyed by the rule key.
func (f *fakeSource) scriptKey(key string, obs Observation, err error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.states[key] = append(f.states[key], scriptedObservation{obs: obs, err: err})
}

var _ Source = (*fakeSource)(nil)
