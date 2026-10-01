package gate

import (
	"container/heap"
	"fmt"
	"sort"
	"strings"
	"time"
)

// CheckStartupHandoff proves the first poll after the measurement pass arrives
// before each rule's maxGap — a transient the steady-state budget cannot see.
// It simulates the poller's first cycles and fails on any first poll past it.
// windowOpen is the earliest instant the classification window can open.
func CheckStartupHandoff(t map[string]RuleTimings, measured map[string]time.Duration, first []Poll, readyAt, windowOpen time.Time, concurrency int) error {
	if len(t) == 0 {
		return nil
	}
	if concurrency < 1 {
		concurrency = 1
	}
	problems, err := handoffProblems(t, measured, first, readyAt, windowOpen, concurrency)
	if err != nil {
		return err
	}
	if len(problems) == 0 {
		return nil
	}

	var b strings.Builder
	fmt.Fprintf(&b, "startup handoff does not fit at concurrency %d:\n", concurrency)
	const maxListed = 5
	for i, p := range problems {
		if i == maxListed {
			fmt.Fprintf(&b, "  - ... and %d more\n", len(problems)-maxListed)
			break
		}
		fmt.Fprintf(&b, "  - %s\n", p)
	}
	if minC := minHandoffConcurrency(t, measured, first, readyAt, windowOpen, concurrency, len(t)); minC > concurrency {
		fmt.Fprintf(&b, "fix by: raising --concurrency to at least %d (currently %d), raising poll-interval, or watching fewer alerts",
			minC, concurrency)
	} else {
		b.WriteString("fix by: raising poll-interval or watching fewer alerts")
	}
	return fmt.Errorf("%s", b.String())
}

// handoffProblems returns one message per rule whose simulated first poll is late.
func handoffProblems(t map[string]RuleTimings, measured map[string]time.Duration, first []Poll, readyAt, windowOpen time.Time, concurrency int) ([]string, error) {
	type job struct {
		uid       string
		title     string
		due       time.Time
		latency   time.Duration
		bound     time.Duration
		maxGap    time.Duration
		pollEvery time.Duration
	}
	jobs := make([]job, 0, len(t))
	for _, p := range first {
		rt, ok := t[p.RuleUID]
		if !ok {
			continue
		}
		m, ok := measured[p.RuleUID]
		if !ok {
			return nil, fmt.Errorf("startup handoff: rule %s was never measured", ruleLabel(rt.title, p.RuleUID))
		}
		jobs = append(jobs, job{
			uid:       p.RuleUID,
			title:     rt.title,
			due:       runnerTime(p, p.GrafanaNow).Add(rt.pollEvery),
			latency:   m,
			bound:     p.SkewBound(),
			maxGap:    rt.maxGap,
			pollEvery: rt.pollEvery,
		})
	}
	if len(jobs) != len(t) {
		return nil, fmt.Errorf("startup handoff: %d of %d rule(s) have a first observation", len(jobs), len(t))
	}
	// Matches Scheduler.Due: earliest-due first, ties by tightest cadence, then uid.
	sort.Slice(jobs, func(i, j int) bool {
		if !jobs[i].due.Equal(jobs[j].due) {
			return jobs[i].due.Before(jobs[j].due)
		}
		if jobs[i].pollEvery != jobs[j].pollEvery {
			return jobs[i].pollEvery < jobs[j].pollEvery
		}
		return jobs[i].uid < jobs[j].uid
	})

	var problems []string
	now := readyAt
	for i := 0; i < len(jobs); {
		if jobs[i].due.After(now) {
			now = jobs[i].due
		}
		free := newTimeHeap(concurrency, now)
		batchEnd := now
		j := i
		for ; j < len(jobs) && !jobs[j].due.After(now); j++ {
			start := free.pop()
			finish := start.Add(jobs[j].latency)
			// `bound` keeps a boundary-adjacent window open on the fail-closed side.
			if gap := finish.Sub(windowOpen) + jobs[j].bound; gap > jobs[j].maxGap {
				problems = append(problems, fmt.Sprintf(
					"rule %s: first poll after the window can open ~%s (measured latency %s) exceeds its maxGap %s",
					ruleLabel(jobs[j].title, jobs[j].uid), gap.Round(time.Millisecond), jobs[j].latency.Round(time.Millisecond), jobs[j].maxGap))
			}
			free.push(finish)
			if finish.After(batchEnd) {
				batchEnd = finish
			}
		}
		now = batchEnd
		i = j
	}
	return problems, nil
}

// minHandoffConcurrency is the smallest concurrency in from..maxConcurrency
// whose simulated handoff fits, or 0 when even one worker per rule does not.
func minHandoffConcurrency(t map[string]RuleTimings, measured map[string]time.Duration, first []Poll, readyAt, windowOpen time.Time, from, maxConcurrency int) int {
	if from < 1 {
		from = 1
	}
	if maxConcurrency < from {
		return 0
	}
	if problems, err := handoffProblems(t, measured, first, readyAt, windowOpen, maxConcurrency); err != nil || len(problems) > 0 {
		return 0
	}
	lo, hi := from, maxConcurrency // lo fails (the caller just checked), hi passes
	for lo+1 < hi {
		mid := lo + (hi-lo)/2
		problems, err := handoffProblems(t, measured, first, readyAt, windowOpen, mid)
		if err == nil && len(problems) == 0 {
			hi = mid
		} else {
			lo = mid
		}
	}
	return hi
}

// timeHeap is the min-heap of C worker free times within one handoff batch.
type timeHeap []time.Time

func newTimeHeap(n int, at time.Time) *timeHeap {
	h := make(timeHeap, n)
	for i := range h {
		h[i] = at
	}
	heap.Init(&h)
	return &h
}

func (h *timeHeap) pop() time.Time {
	return heap.Pop(h).(time.Time)
}

func (h *timeHeap) push(t time.Time) {
	heap.Push(h, t)
}

func (h timeHeap) Len() int           { return len(h) }
func (h timeHeap) Less(i, j int) bool { return h[i].Before(h[j]) }
func (h timeHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *timeHeap) Push(x any)        { *h = append(*h, x.(time.Time)) }
func (h *timeHeap) Pop() any          { old := *h; n := len(old); x := old[n-1]; *h = old[:n-1]; return x }
