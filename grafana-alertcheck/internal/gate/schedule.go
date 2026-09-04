package gate

import (
	"fmt"
	"math/rand/v2"
	"sort"
	"strings"
	"time"
)

// SkewHardLimit is the largest runner↔Grafana clock skew a run tolerates.
// Exported so the CLI reports it verbatim next to a measured skew.
const SkewHardLimit = 60 * time.Second

// fromFutureTolerance is how far ahead of the runner's clock a `from` may sit
// before check refuses it — the same 60s as SkewHardLimit, since a future `from`
// can only be clock disagreement. Once-per-run input validation, not a coverage
// check, so Check applies it and proveCoverage does not.
const fromFutureTolerance = 60 * time.Second

// minDrainTimeout floors drainTimeout (otherwise 2 × max intervalSeconds) so a
// fleet of tight rules still lets a healthy in-flight poll land.
const minDrainTimeout = 2 * time.Minute

// graceWarnFraction is the share of the window above which transitionGrace is
// worth warning about.
const graceWarnFraction = 0.25

// ruleTimings groups the per-rule thresholds derived from a rule's poll
// cadence and its own evaluation interval.
type ruleTimings struct {
	pollEvery      time.Duration
	maxGap         time.Duration
	healthGrace    time.Duration
	evalStaleAfter time.Duration
}

// globalTimings groups the values that apply to the whole run rather than to
// one rule: transitionGrace and drainTimeout are each derived once, across
// every non-skipped watched rule, not per rule.
type globalTimings struct {
	transitionGrace time.Duration
	// graceSource names, and already carries the `for` value of, the rule that
	// set transitionGrace — one string field rather than a second
	// (rule, duration) pair, matching this struct's fixed shape. "none" when no
	// rule contributed (transitionGrace is then 0).
	graceSource  string
	drainTimeout time.Duration
}

// newRuleTimings derives one rule's thresholds from its resolved cadence and
// evaluation interval. pollEvery is an input — already resolved for the
// caller's mode — so no caller can compute maxGap against the wrong authority.
func newRuleTimings(pollEvery time.Duration, intervalSeconds int) ruleTimings {
	interval := time.Duration(intervalSeconds) * time.Second
	maxGap := 2 * pollEvery
	healthGrace := max(maxGap, interval)
	return ruleTimings{
		pollEvery:      pollEvery,
		maxGap:         maxGap,
		healthGrace:    healthGrace,
		evalStaleAfter: 2 * interval,
	}
}

// defaultPollEvery is the default per-rule cadence: half the rule's own
// evaluation interval.
func defaultPollEvery(intervalSeconds int) time.Duration {
	return time.Duration(intervalSeconds) * time.Second / 2
}

// DeriveTimings computes every resolved rule's ruleTimings (keyed by UID) plus
// the shared globalTimings. A non-zero override is used verbatim for every rule
// and never clamped to the default — an override above intervalSeconds/2 widens
// maxGap and is reported as a note, not corrected.
func DeriveTimings(defs []Definition, override time.Duration) (rules map[string]ruleTimings, global globalTimings, notes []string) {
	rules = make(map[string]ruleTimings, len(defs))
	for _, d := range defs {
		def := defaultPollEvery(d.IntervalSeconds)
		pollEvery := def
		if override > 0 {
			pollEvery = override
			if override > def {
				notes = append(notes, fmt.Sprintf(
					"rule %s: --poll-interval %s exceeds half its %ds evaluation interval (%s); maxGap widens accordingly",
					d.Title, override, d.IntervalSeconds, def))
			}
		}
		rules[d.UID] = newRuleTimings(pollEvery, d.IntervalSeconds)
	}
	// In this mode defs ARE the start-of-step snapshot, so they answer what was
	// paused at the window open; only the log-mode counterpart uses the header.
	return rules, deriveGlobalTimings(defs, pausedSet(defs)), notes
}

// pausedSet is Header.pausedAtStart's counterpart for a set of definitions
// resolved at the start of the step, which is the one moment a definition can
// answer "was this paused when the window opened".
func pausedSet(defs []Definition) map[string]bool {
	paused := make(map[string]bool, len(defs))
	for _, d := range defs {
		paused[d.UID] = d.IsPaused
	}
	return paused
}

// DeriveTimingsFromLog is DeriveTimings' log-mode counterpart: pollEvery comes
// from the header (the cadence actually used), not the definitions — re-deriving
// it here would compare recorded gaps against default-cadence thresholds, an
// exit 2 on a clean window (slower override) or a silently passing recorder gap
// (faster override). evalStaleAfter still comes from defs (2 × intervalSeconds).
//
// Three header shapes are hard errors rather than a best-effort derivation,
// because each would silently widen a threshold: a rule with no matching
// definition, a non-positive recorded cadence, and a UID listed twice
// (last-one-wins would widen maxGap on a corrupt log).
//
// It checks only the header-to-defs direction. A definition absent from the
// header is Check's log-identity validation to judge, not this function's.
func DeriveTimingsFromLog(h Header, defs []Definition) (rules map[string]ruleTimings, global globalTimings, err error) {
	byUID := make(map[string]Definition, len(defs))
	for _, d := range defs {
		byUID[d.UID] = d
	}

	rules = make(map[string]ruleTimings, len(h.Rules))
	for _, lr := range h.Rules {
		def, ok := byUID[lr.UID]
		if !ok {
			return nil, globalTimings{}, fmt.Errorf(
				"log header names rule %s (%q), which no current definition matches", lr.UID, lr.Title)
		}
		if _, duplicate := rules[lr.UID]; duplicate {
			return nil, globalTimings{}, fmt.Errorf(
				"log header names rule %s (%q) twice; its recorded cadence is ambiguous", lr.UID, lr.Title)
		}
		if lr.PollEverySeconds <= 0 {
			return nil, globalTimings{}, fmt.Errorf(
				"log header records poll_every_seconds=%v for rule %s (%q); the recorded cadence is required to derive maxGap",
				lr.PollEverySeconds, lr.UID, lr.Title)
		}
		pollEvery := time.Duration(lr.PollEverySeconds * float64(time.Second))
		rules[lr.UID] = newRuleTimings(pollEvery, def.IntervalSeconds)
	}
	// The header, not defs, decides which rules are excluded from the grace:
	// defs were resolved after the window closed. See deriveGlobalTimings.
	return rules, deriveGlobalTimings(defs, h.pausedAtStart()), nil
}

// deriveGlobalTimings computes transitionGrace and drainTimeout over defs. A
// rule skipped at the window open is excluded from the transitionGrace max (its
// `for` can never fire in-window); drainTimeout runs over every resolved rule.
//
// The exclusion authority is pausedAtStart, never Definition.IsPaused: log-mode
// defs are re-resolved after the window closed. Reading the late definitions
// was a quiet fail-open — a rule paused after `to` would drop out of the max,
// collapse the grace past windowEnd (the classification bound AND collection
// deadline), and pass a window the surfacing poll was never recorded for.
func deriveGlobalTimings(defs []Definition, pausedAtStart map[string]bool) globalTimings {
	var g globalTimings
	var maxInterval time.Duration
	for _, d := range defs {
		interval := time.Duration(d.IntervalSeconds) * time.Second
		if interval > maxInterval {
			maxInterval = interval
		}
		if pausedAtStart[d.UID] {
			continue
		}
		if candidate := d.For + interval; candidate > g.transitionGrace {
			g.transitionGrace = candidate
			g.graceSource = fmt.Sprintf("%s (for=%s, interval=%s)", d.Title, d.For, interval)
		}
	}
	g.drainTimeout = max(2*maxInterval, minDrainTimeout)
	return g
}

// Scheduler drives one per-rule schedule, never a global cycle: a rule at
// intervalSeconds=10 alongside twenty at 300 keeps its own 5s cadence without
// forcing the same cadence onto the other twenty.
type Scheduler struct {
	next  map[string]time.Time
	every map[string]time.Duration
}

// NewScheduler builds a Scheduler over per-rule cadences, staggering each
// rule's initial next-due time across [0, pollEvery) so a phase-aligned fleet
// (which would void CheckBudget's burst bound) never arises by construction.
// It takes cadences, not ruleTimings: a scheduler only decides when to poll and
// must not be handed coverage thresholds it never applies.
func NewScheduler(every map[string]time.Duration, now time.Time) *Scheduler {
	s := &Scheduler{
		next:  make(map[string]time.Time, len(every)),
		every: make(map[string]time.Duration, len(every)),
	}
	for uid, pollEvery := range every {
		s.every[uid] = pollEvery
		var offset time.Duration
		if pollEvery > 0 {
			offset = rand.N(pollEvery)
		}
		s.next[uid] = now.Add(offset)
	}
	return s
}

// Due returns the due UIDs, earliest-due-first. Ties break by tightest cadence
// first: the burst bound assumes a newly-due tight rule waits at most one
// in-flight request, which only holds if a simultaneous batch serves the
// tightest rule first. A map-order tie-break would silently void that.
func (s *Scheduler) Due(now time.Time) []string {
	var due []string
	for uid, t := range s.next {
		if !t.After(now) {
			due = append(due, uid)
		}
	}
	sort.Slice(due, func(i, j int) bool {
		a, b := due[i], due[j]
		if !s.next[a].Equal(s.next[b]) {
			return s.next[a].Before(s.next[b])
		}
		if s.every[a] != s.every[b] {
			return s.every[a] < s.every[b]
		}
		return a < b // stable, deterministic fallback for an exact tie
	})
	return due
}

// Mark records that uid was just polled at now, scheduling its next poll one
// cadence later.
func (s *Scheduler) Mark(uid string, now time.Time) {
	s.next[uid] = now.Add(s.every[uid])
}

// earliestDue returns the earliest next-due time (false when empty). The loop
// waits exactly that long instead of a fixed tick, which would poll slack rules
// early (wasting budget) or wake late for the tightest rule (opening a gap).
func (s *Scheduler) earliestDue() (time.Time, bool) {
	var earliest time.Time
	for _, t := range s.next {
		if earliest.IsZero() || t.Before(earliest) {
			earliest = t
		}
	}
	return earliest, !earliest.IsZero()
}

// CheckBudget proves at start that a fully resolved schedule can be served. t
// and measured are keyed by UID, and measured must carry every UID in t (a rule
// never measured cannot have its budget proved). It fails on any of three
// conditions:
//
//   - utilization — the long-run request rate exceeds the concurrency;
//   - a single rule's request cannot fit inside its own cadence;
//   - the burst bound — the slowest measured request is slower than the fleet's
//     tightest cadence, which can open a mid-run gap beyond that rule's maxGap.
//
// The message names only the three operator controls: concurrency,
// poll-interval, and the alert list.
func CheckBudget(t map[string]ruleTimings, measured map[string]time.Duration, concurrency int) error {
	if len(t) == 0 {
		return nil
	}

	uids := make([]string, 0, len(t))
	for uid := range t {
		uids = append(uids, uid)
	}
	sort.Strings(uids) // deterministic message order

	for _, uid := range uids {
		if _, ok := measured[uid]; !ok {
			return fmt.Errorf("schedule budget: rule %s was never measured", uid)
		}
	}

	var utilization float64
	tightestUID := uids[0] // the rule with the smallest pollEvery seen so far — a UID, not a duration
	var maxMeasuredUID string
	var maxMeasured time.Duration
	var overCadence []string
	for _, uid := range uids {
		rt, m := t[uid], measured[uid]
		utilization += float64(m) / float64(rt.pollEvery)
		if t[tightestUID].pollEvery > rt.pollEvery {
			tightestUID = uid
		}
		if m > maxMeasured {
			maxMeasured, maxMeasuredUID = m, uid
		}
		if m > rt.pollEvery {
			overCadence = append(overCadence, uid)
		}
	}

	var problems []string
	if utilization > float64(concurrency) {
		problems = append(problems, fmt.Sprintf("utilization %.2f exceeds concurrency %d", utilization, concurrency))
	}
	for _, uid := range overCadence {
		problems = append(problems, fmt.Sprintf(
			"rule %s: measured %s exceeds its own poll-interval %s", uid, measured[uid], t[uid].pollEvery))
	}
	if maxMeasured > t[tightestUID].pollEvery {
		problems = append(problems, fmt.Sprintf(
			"burst bound: rule %s's measured %s exceeds the fleet's tightest poll-interval %s (rule %s)",
			maxMeasuredUID, maxMeasured, t[tightestUID].pollEvery, tightestUID))
	}

	if len(problems) == 0 {
		return nil
	}

	var b strings.Builder
	fmt.Fprintf(&b, "schedule does not fit at concurrency %d:\n", concurrency)
	for _, uid := range uids {
		fmt.Fprintf(&b, "  rule %s: measured %s, poll-interval %s\n", uid, measured[uid], t[uid].pollEvery)
	}
	for _, p := range problems {
		fmt.Fprintf(&b, "  - %s\n", p)
	}
	b.WriteString("fix by: raising concurrency, raising poll-interval, or watching fewer alerts")
	return fmt.Errorf("%s", b.String())
}

// StartupSummary formats the pre-run print an operator sees before the wait:
// the total planned run time and the rule (with its `for` value) that set
// transitionGrace, plus a warning when the grace eats more than
// graceWarnFraction of the requested window. from/to are the requested
// classification window.
func StartupSummary(from, to time.Time, global globalTimings) (summary, warning string) {
	window := to.Sub(from)
	total := window + global.transitionGrace + global.drainTimeout
	source := global.graceSource
	if source == "" {
		source = "none"
	}
	summary = fmt.Sprintf(
		"planned run time: %s (window %s + transitionGrace %s [source: %s] + drainTimeout %s)",
		total, window, global.transitionGrace, source, global.drainTimeout)

	if window > 0 && float64(global.transitionGrace) > float64(window)*graceWarnFraction {
		warning = fmt.Sprintf(
			"transitionGrace %s is more than %.0f%% of the window %s (source: %s) — the window may be too short for this alert's `for`",
			global.transitionGrace, graceWarnFraction*100, window, source)
	}
	return summary, warning
}
