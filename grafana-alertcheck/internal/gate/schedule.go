package gate

import (
	"fmt"
	"math/rand/v2"
	"sort"
	"strings"
	"time"
)

// SkewHardLimit is one of §5's filled-in values (basis: §16; §22.11 asserts
// 120s errors, 30s does not). Defined here, in schedule.go's named-constants
// block, per §5's instruction — it moved out of source.go now that P4 exists;
// P2 needed it before this file did, so it started there. Exported (P10) so
// the CLI can report it verbatim next to a measured skew instead of keeping
// its own mirrored copy.
const SkewHardLimit = 60 * time.Second

// fromFutureTolerance is how far ahead of the runner's own clock a supplied
// `from` may sit before check refuses it (§7: "from in the future, more than
// the skew tolerance — error"). §7 names no number, so this is the judgment
// call §5's table records: the same 60s as SkewHardLimit, because the only
// legitimate reason for a `from` in the future is clock disagreement between
// the deploy step and the check step, and that is bounded by the same figure.
// It is once-per-run input validation, not a per-rule coverage check, so
// Check applies it (P9) and proveCoverage does not.
const fromFutureTolerance = 60 * time.Second

// minDrainTimeout is §5's floor on drainTimeout: max(2 x max(intervalSeconds),
// 2m). Without the floor, a fleet of very tight rules would derive a
// drainTimeout too short to let a healthy in-flight poll land.
const minDrainTimeout = 2 * time.Minute

// graceWarnFraction is §13.2's threshold for warning that transitionGrace eats
// too much of the requested window: "approximately one quarter of the window".
const graceWarnFraction = 0.25

// ruleTimings groups the per-rule threshold values §5/§10.1/§14.1 derive from
// a rule's poll cadence and its own evaluation interval.
type ruleTimings struct {
	pollEvery      time.Duration
	maxGap         time.Duration
	healthGrace    time.Duration
	evalStaleAfter time.Duration
}

// globalTimings groups the values that apply to the whole run rather than to
// one rule: §13.1's transitionGrace and §19's drainTimeout are each derived
// once, across every non-skipped watched rule, not per rule.
type globalTimings struct {
	transitionGrace time.Duration
	// graceSource names, and already carries the `for` value of, the rule
	// that set transitionGrace (§13.2 requires printing both) — one string
	// field rather than a second (rule, duration) pair, matching this
	// struct's fixed shape. "none" when no rule contributed (transitionGrace
	// is then 0).
	graceSource  string
	drainTimeout time.Duration
}

// newRuleTimings derives one rule's thresholds from its fully-resolved poll
// cadence and its evaluation interval (§5, §10.1, §14.1). pollEvery arrives
// already resolved for the caller's mode — the §5 default, the operator's
// --poll-interval override, or (in log mode, a later phase) the cadence
// recorded in the log header. Deriving pollEvery inline here, instead of
// accepting it as an input, would let a caller in the wrong mode compute
// maxGap against the wrong authority — see the "Two authorities" note in P5.
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

// defaultPollEvery is §5's default per-rule cadence: half the rule's own
// evaluation interval.
func defaultPollEvery(intervalSeconds int) time.Duration {
	return time.Duration(intervalSeconds) * time.Second / 2
}

// DeriveTimings computes every resolved rule's ruleTimings, keyed by UID,
// plus the shared globalTimings, from resolved definitions and watch's
// optional --poll-interval override (0 = no override: use each rule's §5
// default of half its own interval). Per §5.1, a supplied override is used
// verbatim for every rule and is never clamped down to the default even when
// it exceeds intervalSeconds/2 — that case is reported back as a note, not
// silently corrected or refused, because clamping would defeat the one knob
// §5.1 gives an operator for making a tight schedule fit.
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
	// In this mode the defs ARE the start-of-step snapshot — watch resolves
	// them before it detaches, and single-step check before its first
	// observation — so they can answer what was paused when the window opened.
	// Only the log-mode counterpart below has to look elsewhere.
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

// DeriveTimingsFromLog is DeriveTimings' log-mode counterpart, and the two
// authorities of P5 are the whole reason it exists as a separate function.
// pollEvery comes from the header — the cadence the recording ACTUALLY used,
// after any --poll-interval override — and maxGap and healthGrace follow from
// it. Re-deriving pollEvery from defs here would compare gaps recorded at the
// override cadence against thresholds computed from the default: exit 2 on a
// clean window when the override was slower, and, worse, a real recorder gap
// passing silently when it was faster.
//
// evalStaleAfter still comes from defs (2 x intervalSeconds): it is a property
// of the rule's own evaluation cadence and is unaffected by how often the gate
// polled.
//
// Three shapes of header are errors rather than a best-effort derivation,
// because each one would otherwise widen a threshold silently:
//
//   - a rule with no matching definition — a log that names a rule nobody can
//     resolve cannot have that rule's coverage proved;
//   - a non-positive recorded cadence — a log that cannot say how often it was
//     written cannot have maxGap derived, and defaulting the cadence would
//     prove a window that was never observed;
//   - the same UID twice — last-one-wins would take whichever cadence happened
//     to be written last, and a slower duplicate widens maxGap. That is a
//     fail-open reachable through nothing but log corruption.
//
// It checks only the header-to-defs direction. The opposite direction — a
// resolved definition absent from the header — is NOT this function's to
// judge: it is §19.1 step 3's log-identity validation, and it belongs to P9's
// Check, which is the only caller that knows both sets and can name the
// mismatch. Without that check a definition simply gets no timings entry, and
// a downstream lookup would read a zero maxGap: fail-closed (every gap
// exceeds it) but silent, so P9 must reject the set mismatch by name rather
// than let a rule fail for an unexplained reason.
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

// deriveGlobalTimings computes transitionGrace and drainTimeout over defs
// (§5, §13.1, §19).
//
// A rule paused before the window opened — skipped, §12 — is excluded from the
// transitionGrace max: its `for` value can never fire during the window, so
// counting it would only inflate the wait past what any watched rule actually
// needs (a judgment call the v2 plan makes explicitly for this formula; §19's
// drainTimeout carries no such exclusion, so it still runs over every resolved
// rule).
//
// "Before the window opened" is the whole content of that exclusion, so the
// authority is pausedAtStart and NEVER Definition.IsPaused: in log mode the
// definitions are re-resolved after the window closed. Reading them instead
// was a fail-open, and a quiet one. transitionGrace is what lets a condition
// arising just before `to` be seen when it surfaces at to + `for`, and
// windowEnd is BOTH the classification bound and the collection deadline — so
// a rule somebody paused after `to` dropped out of the max, the grace
// collapsed, the surfacing poll was never even recorded, and the run reported
// clean. With one watched rule the shrink is total.
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

// Scheduler drives one per-rule schedule, never a global cycle (§5): a rule
// at intervalSeconds=10 alongside twenty at 300 keeps its own 5s cadence
// without forcing the same cadence onto the other twenty.
type Scheduler struct {
	next  map[string]time.Time
	every map[string]time.Duration
}

// NewScheduler builds a Scheduler over per-rule cadences (keyed by UID),
// staggering each rule's initial next-due time across [0, pollEvery) so the
// fleet does not start phase-aligned (§5's burst-bound proof depends on this:
// an already-staggered fleet only re-aligns by chance, briefly, not by
// construction).
//
// It takes cadences rather than whole ruleTimings on purpose: a scheduler
// decides when to poll and nothing else, so it must not be handed maxGap,
// healthGrace or evalStaleAfter. Those are coverage thresholds, they are
// applied by the pure layer at classification time, and the recorder that
// drives this scheduler never applies them at all.
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

// Due returns the UIDs whose next-due time has arrived, earliest-due-first.
// Ties (equal next-due time) break by tightest cadence first: the burst-bound
// proof in §5 assumes a newly-due tight rule waits at most for one in-flight
// request, which only holds if a simultaneous batch serves the tightest rule
// ahead of slacker ones. A tie-break that instead followed map iteration
// order would silently void that proof — nothing else would fail until a
// phase-aligned fleet opened a mid-run gap in production.
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

// earliestDue returns the earliest scheduled next-due time, and false when the
// scheduler holds no rules at all. The recorder's loop (P6) waits exactly that
// long instead of waking on a fixed tick: a fixed tick either polls a slack
// rule early — spending request budget the §5 formulas already accounted for —
// or wakes too late for the tightest rule and opens a gap inside its own
// maxGap.
func (s *Scheduler) earliestDue() (time.Time, bool) {
	var earliest time.Time
	for _, t := range s.next {
		if earliest.IsZero() || t.Before(earliest) {
			earliest = t
		}
	}
	return earliest, !earliest.IsZero()
}

// CheckBudget applies §5's error-at-start check to a fully resolved schedule.
// t and measured are both keyed by rule UID; measured must carry every UID in
// t; a rule this run never measured can't have its budget proved, and a
// silent zero-duration default would be exactly the kind of pass-on-an-
// unproven-window bug §5 exists to catch. CheckBudget fails when any of three
// conditions holds (sanity-checked against §22.3's mixed-interval regression
// in this phase's tests):
//
//   - utilization: the long-run request rate exceeds what concurrency serves;
//   - a single rule's own request cannot fit inside its own cadence;
//   - the burst bound: the slowest measured request is slower than the
//     fleet's tightest cadence, which — even under earliest-due-first
//     ordering — can open a mid-run gap bigger than that rule's maxGap.
//
// The message never suggests a single interval (§5.1) — only the three
// controls an operator actually has: concurrency, poll-interval, and the
// alert list.
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

// StartupSummary formats §13.2's required pre-run print: the total planned
// run time and the rule (with its `for` value) that set transitionGrace, plus
// a warning when the grace eats more than graceWarnFraction of the requested
// window. from/to are the requested classification window.
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
