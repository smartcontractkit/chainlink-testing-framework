package gate

import (
	"fmt"
	"slices"
	"strings"
	"time"
)

// ReasonNodata is decide's own unobservable reason: proveCoverage never sets it
// — health=nodata is a note there, never fatal — because escalating it needs
// Policy.NodataIsUnobservable, which only decide (the Policy-holding seam) has.
const ReasonNodata UnobservableReason = "nodata"

// Outcome is the verdict of one instance's timeline, and (after decide takes
// the worst across instances) of the rule. It is a published JSON output: the
// fail values stay distinct even though v1 maps them all to exit 1, so a later
// version can split them without breaking the interface.
type Outcome string

const (
	OutcomeHealthy      Outcome = "healthy"
	OutcomeNewFailure   Outcome = "new_failure"
	OutcomeRecovered    Outcome = "recovered"
	OutcomeStillFailing Outcome = "still_failing"
	OutcomeUnstable     Outcome = "unstable"
	OutcomePaused       Outcome = "paused"
	OutcomeNotVerified  Outcome = "not_verified"
	// OutcomeNotCounted is synthetic: decide uses it for a --min-observed
	// deficit row that no resolved rule explains. It is not a verdict on an
	// alert, so it is deliberately not named after one.
	OutcomeNotCounted Outcome = "not_counted"
)

// PreexistingPolicy governs only the ONE ambiguous case in the outcome table:
// an instance that was already bad when the window opened. A new_failure or
// unstable instance is a fail under every policy, so this type only ever
// changes how `recovered` and `still_failing` are judged (isViolation
// below).
type PreexistingPolicy string

const (
	// PreexistingFailUnlessRecovered is the default: a preexisting instance
	// that clears and stays clear is a pass (`recovered`); one that never
	// clears is still a fail (`still_failing`).
	PreexistingFailUnlessRecovered PreexistingPolicy = "fail-unless-recovered"
	// PreexistingFail makes ANY preexisting instance a fail, even one that
	// recovers — for a user who wants no benefit of the doubt for a
	// condition this release did not cause.
	PreexistingFail PreexistingPolicy = "fail"
	// PreexistingIgnore disregards a preexisting instance entirely, whether
	// it recovers or stays bad for the whole window: only a genuinely NEW
	// bad episode (new_failure or unstable) can fail the rule.
	PreexistingIgnore PreexistingPolicy = "ignore"
)

// Violation is one instance whose timeline outcome counts against the run,
// after the preexisting policy has been applied (isViolation below).
type Violation struct {
	Alert, RuleUID string
	Outcome        Outcome
	State          State
	Health         string // raw, reporting-only, like Poll.Health
	LastError      string
	// FirstSeen is the episode's onset in the runner domain (translated by the
	// poll's own skew), or `from` when preexisting — never a raw Grafana time.
	FirstSeen time.Time
	// ClearedAt is zero unless the episode closed via a genuine Cleared event.
	ClearedAt      time.Time
	InstanceLabels map[string]string
	// Note explains a Violation with no instance behind it — decide's synthetic
	// MinObserved shortfall — and must not double as LastError (reporting-only
	// rule state from a real poll).
	Note string
}

// RuleVerdict is one rule's worst-of outcome, present for every resolved rule
// (passes included) so the table shows every alert asked for.
type RuleVerdict struct {
	Alert, RuleUID string
	Outcome        Outcome
	BadFor         time.Duration // total wall-clock time any instance was bad inside the window, overlaps merged
	PollEvery      time.Duration
	Note           string
}

// Policy is decide's narrowed, pure-layer view of a Config: the classification
// knobs and the window, nothing else. No URL, no token, no I/O handles — those
// never reach the pure layer.
type Policy struct {
	States                            []State
	Preexisting                       PreexistingPolicy
	MinObserved                       int
	AllowPaused, NodataIsUnobservable bool
	From, To                          time.Time
}

// RuleThresholds is one non-skipped rule's resolved coverage thresholds,
// carried on Result so the CLI's table can print the numbers that answer "why"
// on exit 2 without decide exposing the unexported ruleTimings type itself.
type RuleThresholds struct {
	MaxGap         time.Duration
	HealthGrace    time.Duration
	EvalStaleAfter time.Duration
}

// GlobalThresholds is the run-wide half of the same information:
// transitionGrace and drainTimeout apply once, across every non-skipped
// watched rule, not per rule (globalTimings).
type GlobalThresholds struct {
	TransitionGrace time.Duration
	// GraceSource names, and already carries the `for` value of, the rule that
	// set TransitionGrace — an operator has to see both. "none" when no rule
	// contributed (TransitionGrace is then 0).
	GraceSource  string
	DrainTimeout time.Duration
}

// Result is decide's whole answer: everything the human table and the JSON
// output need. Coverage carries one CoverageResult per non-skipped rule —
// there is deliberately no separate Interval type anywhere in the project.
type Result struct {
	From, To       time.Time
	GrafanaVersion string
	ClockSkew      time.Duration // the largest |skew| across every poll decide was given, not only the ones a rule's window actually used
	// ClockSkewBound is the skew BOUND (RTT/2) of that SAME poll — not
	// the largest bound seen overall, which would pair a wide bound from an
	// unrelated slow request with the worst skew and misstate how tightly
	// that skew is actually known. SkewHardLimit is a separate, fixed input
	// validation threshold (source.go) and is not an error bound on this
	// value; the CLI prints both, but must not conflate them.
	ClockSkewBound time.Duration
	Coverage       map[string]CoverageResult
	// Thresholds carries one RuleThresholds per rule Coverage also covers —
	// every non-skipped rule, keyed by UID. A skipped rule has neither: it was
	// never scheduled, so it has no maxGap/healthGrace/evalStaleAfter to
	// report.
	Thresholds map[string]RuleThresholds
	Global     GlobalThresholds
	Verdicts   []RuleVerdict
	Violations []Violation
	// TerminatedEarly is set only when fail-fast stopped before the window
	// closed; the coverage proof is then over [from, At]. To and Global below
	// still report the requested values. Published JSON output.
	TerminatedEarly *Termination `json:"terminated_early,omitempty"`
}

// episode is one contiguous, policy-bad span of one instance's timeline,
// already resolved to the runner domain and clamped to [from, windowEnd]. It
// never crosses a genuine Cleared event: a Vanished marker freezes the state
// instead of closing the episode, which is what keeps a vanish from ever
// reading as a recovery.
type episode struct {
	start, end        time.Time
	closedByRealClear bool
}

// instanceTimeline accumulates one instance's walk across a rule's in-window
// polls. preexisting is decided once, the first time this key is seen bad: by
// the translated ActiveAt against `from`, never by which poll happened to
// report it first — a poll's own cadence is not evidence of when the condition
// actually began.
//
// recoveryOpen marks an episode last seen Recovering: an extension poll may
// close it only until recoveryDeadline.
type instanceTimeline struct {
	labels           map[string]string
	preexisting      bool
	seen             bool
	badOpen          bool
	recoveryOpen     bool
	recoveryDeadline time.Time
	episodeStart     time.Time
	lastState        State
	lastHealth       string
	lastError        string
	episodes         []episode
}

// runnerTime translates a Grafana-domain timestamp into the runner domain by
// undoing poll p's measured skew. The single implementation for the package —
// coverage.go's window-membership and heartbeat-boundary checks use it too.
func runnerTime(p Poll, grafanaDomain time.Time) time.Time {
	return grafanaDomain.Add(-p.Skew())
}

// classifyRule builds every instance timeline for one rule across
// [from, windowEnd] and reduces them to the rule's worst outcome, merged
// BadFor, and the Violations the preexisting policy charges against the run.
// PURE: no I/O, no clock reads; polls need not be pre-filtered to this rule.
// t bounds a recovering episode, identically to the recovery wait.
func classifyRule(def Definition, t RuleTimings, polls []Poll, from, windowEnd time.Time, badStates map[State]bool, pol PreexistingPolicy) (Outcome, time.Duration, []Violation) {
	rulePolls := pollsForRule(polls, def.UID)

	timelines := make(map[string]*instanceTimeline)
	order := make([]string, 0)

	// get backfills labels on the first real Instance: a bare Cleared/Vanished
	// marker can create the timeline first (with no labels), and a later re-fire
	// must not report an empty InstanceLabels.
	get := func(key string, labels map[string]string) *instanceTimeline {
		tl, ok := timelines[key]
		if !ok {
			tl = &instanceTimeline{labels: labels}
			timelines[key] = tl
			order = append(order, key)
			return tl
		}
		if tl.labels == nil && labels != nil {
			tl.labels = labels
		}
		return tl
	}

	openEpisode := func(tl *instanceTimeline, start time.Time) {
		tl.badOpen = true
		tl.episodeStart = start
	}
	closeEpisode := func(tl *instanceTimeline, end time.Time, isReal bool) {
		// inWindowPolls widens its boundary outward by the skew bound, so a
		// translated end can land past windowEnd or before episodeStart; clamp
		// both, otherwise mergeDurations gets an inverted span.
		if end.After(windowEnd) {
			end = windowEnd
		}
		if end.Before(tl.episodeStart) {
			end = tl.episodeStart
		}
		tl.episodes = append(tl.episodes, episode{start: tl.episodeStart, end: end, closedByRealClear: isReal})
		tl.badOpen = false
		tl.recoveryOpen = false // resolved; a later re-fire is out of scope
		tl.recoveryDeadline = time.Time{}
	}
	// onsetOf resolves a fresh episode's start: the instance's own ActiveAt,
	// translated to the runner domain by this poll's skew, clamped to
	// [from, windowEnd].
	onsetOf := func(p Poll, inst Instance) time.Time {
		start := runnerTime(p, inst.ActiveAt)
		if start.Before(from) {
			start = from
		}
		if start.After(windowEnd) {
			start = windowEnd
		}
		return start
	}
	// openRecovering: Recovering is only reachable from Alerting, so the fire
	// predates this observation; ActiveAt is the recovery onset, not the fire.
	openRecovering := func(tl *instanceTimeline) {
		tl.preexisting = true
		openEpisode(tl, from)
	}

	// recoveryExpired: an extension clear past the deadline does not resolve.
	recoveryExpired := func(tl *instanceTimeline, p Poll) bool {
		return tl.recoveryDeadline.IsZero() || runnerTime(p, p.GrafanaNow).After(tl.recoveryDeadline)
	}

	for _, p := range rulePolls {
		if pollBeforeWindow(p, from) {
			continue
		}
		extension := pollAfterWindow(p, windowEnd)
		if extension && (p.IsPaused || !p.Found) {
			// Not a resolution: retire eligibility, leave the episode open.
			for _, tl := range timelines {
				tl.recoveryOpen = false
			}
			continue
		}

		byKey := make(map[string]Instance, len(p.Abnormal))
		for _, inst := range p.Abnormal {
			byKey[instanceKey(inst.Labels)] = inst
		}

		for key, inst := range byKey {
			tl, known := timelines[key]
			if extension && (!known || !tl.recoveryOpen) {
				continue // post-window onset: not ours to classify
			}
			tl = get(key, inst.Labels) // backfills labels on an existing bare timeline
			bad := badStates[inst.State]
			switch {
			case !tl.seen:
				tl.seen = true
				if bad {
					if inst.State == StateRecovering {
						openRecovering(tl)
					} else {
						// Fail-closed: "preexisting" only when even the
						// worst-case skew error places the onset at or before
						// `from`; an onset that might be in-window must
						// classify as a new episode.
						activeAtRunner := runnerTime(p, inst.ActiveAt)
						tl.preexisting = !activeAtRunner.Add(p.SkewBound()).After(from)
						if tl.preexisting {
							openEpisode(tl, from)
						} else {
							openEpisode(tl, onsetOf(p, inst))
						}
					}
				}
			case bad && !tl.badOpen:
				// Already observed in-window and not bad: a Recovering here
				// means the fire happened in-window (a missed Alerting poll
				// hides it), so it must never be preexisting.
				openEpisode(tl, onsetOf(p, inst))
			case !bad && tl.badOpen:
				if extension && recoveryExpired(tl, p) {
					continue // a clear past the deadline does not resolve
				}
				closeEpisode(tl, runnerTime(p, p.GrafanaNow), true)
			}
			if bad && inst.State == StateRecovering {
				tl.recoveryOpen = true
				tl.recoveryDeadline = runnerTime(p, inst.ActiveAt).Add(p.KeepFiringFor() + t.evalStaleAfter + t.pollEvery)
			}
			tl.lastState, tl.lastHealth, tl.lastError = inst.State, p.Health, p.LastError
		}

		for _, key := range p.Cleared {
			tl, known := timelines[key]
			if extension && (!known || !tl.recoveryOpen) {
				continue
			}
			if !known {
				tl = get(key, nil)
			}
			if !tl.seen {
				// Cleared on first mention: the transition happened pre-window,
				// with no in-window evidence it was ever bad.
				tl.seen = true
				continue
			}
			if tl.badOpen {
				if !extension || !recoveryExpired(tl, p) {
					closeEpisode(tl, runnerTime(p, p.GrafanaNow), true)
				}
			}
			tl.lastHealth, tl.lastError = p.Health, p.LastError
		}

		// Vanished is a deliberate no-op: freeze badOpen/preexisting as-is, so a
		// vanish while bad stays bad (never reading as a recovery).
		for _, key := range p.Vanished {
			tl, known := timelines[key]
			if extension && (!known || !tl.recoveryOpen) {
				continue
			}
			if !known {
				tl = get(key, nil)
			}
			tl.seen = true
			tl.lastHealth = p.Health
		}
	}

	// Map iteration order is nondeterministic; sort so Violations/BadFor output
	// is stable for a given input (like log.go sorts Cleared/Vanished).
	slices.Sort(order)

	var (
		outcome = OutcomeHealthy
		badFor  []episode
		viols   []Violation
	)

	for _, key := range order {
		tl := timelines[key]
		if tl.badOpen {
			closeEpisode(tl, windowEnd, false)
		}
		if len(tl.episodes) == 0 {
			continue
		}

		var instOutcome Outcome
		switch {
		case len(tl.episodes) > 1:
			instOutcome = OutcomeUnstable
		case tl.preexisting:
			if tl.episodes[0].closedByRealClear {
				instOutcome = OutcomeRecovered
			} else {
				instOutcome = OutcomeStillFailing
			}
		default:
			// A genuinely new onset fails whether or not it clears in-window;
			// only a preexisting condition earns `recovered`.
			instOutcome = OutcomeNewFailure
		}

		if outcomeRank(instOutcome) > outcomeRank(outcome) {
			outcome = instOutcome
		}
		badFor = append(badFor, tl.episodes...)

		if isViolation(instOutcome, pol) {
			var clearedAt time.Time
			last := tl.episodes[len(tl.episodes)-1]
			if last.closedByRealClear {
				clearedAt = last.end
			}
			viols = append(viols, Violation{
				Alert:          def.Title,
				RuleUID:        def.UID,
				Outcome:        instOutcome,
				State:          tl.lastState,
				Health:         tl.lastHealth,
				LastError:      tl.lastError,
				FirstSeen:      tl.episodes[0].start,
				ClearedAt:      clearedAt,
				InstanceLabels: tl.labels,
			})
		}
	}

	return outcome, mergeDurations(badFor), viols
}

// isViolation decides whether one instance's outcome counts against the run,
// once the preexisting policy is applied. new_failure and unstable always do:
// both contain a genuinely new bad episode, so no policy forgives them.
// recovered and still_failing are, by classifyRule's construction,
// ALWAYS preexisting (a non-preexisting single episode is new_failure instead,
// regardless of whether it clears) — so these are the only two policy can
// change, and isViolation needs no separate preexisting flag to know that.
func isViolation(o Outcome, pol PreexistingPolicy) bool {
	switch o {
	case OutcomeNewFailure, OutcomeUnstable:
		return true
	case OutcomeStillFailing:
		return pol != PreexistingIgnore
	case OutcomeRecovered:
		return pol == PreexistingFail
	default:
		return false
	}
}

// outcomeRank orders outcomes for classifyRule's worst-of reduction across a
// rule's instances:
//
//	not_verified > {unstable, still_failing, new_failure} > recovered >
//	paused > healthy
//
// with not_verified and paused applied outside this function (decide owns
// both: not_verified from CoverageResult, paused from the log header). The
// three fail values are not ranked against each other by anything that reads
// this, so their relative order here is an arbitrary but fixed tie-break, not
// a claim that one is worse than another.
func outcomeRank(o Outcome) int {
	switch o {
	case OutcomeUnstable:
		return 4
	case OutcomeStillFailing:
		return 3
	case OutcomeNewFailure:
		return 2
	case OutcomeRecovered:
		return 1
	default: // OutcomeHealthy
		return 0
	}
}

// mergeDurations sums the wall-clock time covered by a set of episodes,
// merging overlaps so a rule with several simultaneously-bad instances is
// not reported as bad for longer than it actually was.
func mergeDurations(eps []episode) time.Duration {
	if len(eps) == 0 {
		return 0
	}
	sorted := slices.Clone(eps)
	slices.SortStableFunc(sorted, func(a, b episode) int { return a.start.Compare(b.start) })

	var total time.Duration
	cur := sorted[0]
	for _, e := range sorted[1:] {
		if e.start.After(cur.end) {
			total += cur.end.Sub(cur.start)
			cur = e
			continue
		}
		if e.end.After(cur.end) {
			cur.end = e.end
		}
	}
	total += cur.end.Sub(cur.start)
	return total
}

// pollsForRule filters polls to one rule and sorts them by GrafanaNow, the
// same selection proveCoverage uses (by UID, never by title) — stable, because
// two polls sharing a coarse Date header must not reorder nondeterministically
// in a pure function. This is the single filter+sort implementation for the
// package: proveCoverage calls it too, rather than keeping its own copy that
// could silently drift from this one's membership test.
func pollsForRule(polls []Poll, uid string) []Poll {
	var out []Poll
	for _, p := range polls {
		if p.RuleUID == uid {
			out = append(out, p)
		}
	}
	slices.SortStableFunc(out, func(a, b Poll) int { return a.GrafanaNow.Compare(b.GrafanaNow) })
	return out
}

// badStateSet turns Policy.States into a lookup set, defaulting to
// {firing, recovering} when the caller leaves States empty — decide applies
// the default itself so a test can pass a zero-value Policy and get the real
// default, rather than depending on the CLI to have filled it in. Recovering
// is bad by default because the instance is still firing until its
// KeepFiringFor elapses; an explicit --states firing opts out.
func badStateSet(states []State) map[State]bool {
	if len(states) == 0 {
		states = []State{StateFiring, StateRecovering}
	}
	set := make(map[State]bool, len(states))
	for _, s := range states {
		set[s] = true
	}
	return set
}

// applyNodataPolicy promotes a sustained health=nodata run to unobservable
// when the caller has checked Policy.NodataIsUnobservable. proveCoverage never
// escalates nodata itself (most fleets map no-data to OK), so both the
// end-of-run decide and the fail-fast guard must apply this identically — one
// helper, so the two can never drift.
func applyNodataPolicy(def Definition, polls []Poll, cov *CoverageResult, t RuleTimings, from, windowEnd time.Time) {
	if cov.Unobservable {
		return
	}
	inWindow := inWindowPolls(pollsForRule(polls, def.UID), from, windowEnd)
	runLen, sawAny := longestHealthRun(inWindow, "nodata")
	if !sawAny || runLen <= t.healthGrace {
		return
	}
	cov.Unobservable = true
	cov.Proved = false
	if cov.Reason == "" {
		cov.Reason = ReasonNodata
	}
	cov.Notes = append(cov.Notes, fmt.Sprintf(
		"rule %q: health=nodata for %s exceeds healthGrace %s and --nodata-is-unobservable is set",
		def.Title, runLen, t.healthGrace))
}

// decide is the pure seam between the collected evidence and the CLI's exit
// code, and carries nearly the whole test suite because of it. It combines
// proveCoverage's nine checks with classifyRule's timelines under one Policy,
// and owns the inability-beats-violation rule: any not-verified rule makes
// decide return a non-nil error, which the CLI maps to exit 2 unconditionally
// — never to 0 or 1, and never suppressed by a real violation found alongside
// it.
//
// Result is fully populated even when the returned error is non-nil. A caller
// must not use Violations to second-guess the error, but Result stays useful
// for the human table on exit 2.
func decide(h Header, polls []Poll, sentinel *time.Time, defs []Definition,
	rt map[string]RuleTimings, gt GlobalTimings, pol Policy) (Result, error) {

	badStates := badStateSet(pol.States)

	result := Result{
		From:           pol.From,
		To:             pol.To,
		GrafanaVersion: h.GrafanaVersion,
		Coverage:       make(map[string]CoverageResult),
		Thresholds:     make(map[string]RuleThresholds),
		Global: GlobalThresholds{
			TransitionGrace: gt.transitionGrace,
			GraceSource:     graceSourceOrNone(gt.graceSource),
			DrainTimeout:    gt.drainTimeout,
		},
	}
	skewSeen := false
	for _, p := range polls {
		s := p.Skew()
		if s < 0 {
			s = -s
		}
		// The bound travels with its own poll's skew (see Result.ClockSkewBound),
		// overwritten in lockstep. >= rather than > so a bound is still assigned
		// when every poll's skew is exactly 0.
		if !skewSeen || s > result.ClockSkew {
			result.ClockSkew = s
			result.ClockSkewBound = p.SkewBound()
			skewSeen = true
		}
	}

	minObserved := pol.MinObserved
	if minObserved == 0 {
		minObserved = len(defs)
	}

	windowEnd := pol.To.Add(gt.transitionGrace)

	var (
		pausedRules       []Definition
		watchedCount      int
		anyUnobservable   bool
		unobservableNames []string
	)

	// `paused` is decided from the header, never from defs: defs are resolved
	// after the window closed, so Definition.IsPaused describes the present,
	// while Header.pausedAtStart describes the window open — the only moment
	// "paused before the window opened" can mean.
	pausedAtStart := h.pausedAtStart()

	for _, def := range defs {
		if pausedAtStart[def.UID] {
			pausedRules = append(pausedRules, def)
			result.Verdicts = append(result.Verdicts, RuleVerdict{
				Alert: def.Title, RuleUID: def.UID, Outcome: OutcomePaused,
				PollEvery: rt[def.UID].pollEvery,
				Note:      "paused before the window opened",
			})
			continue
		}
		watchedCount++

		t := rt[def.UID]
		cov := proveCoverage(h, polls, sentinel, t, def, pol.From, pol.To, gt.transitionGrace)

		if pol.NodataIsUnobservable {
			applyNodataPolicy(def, polls, &cov, t, pol.From, windowEnd)
		}
		result.Coverage[def.UID] = cov
		result.Thresholds[def.UID] = RuleThresholds{
			MaxGap:         t.maxGap,
			HealthGrace:    t.healthGrace,
			EvalStaleAfter: t.evalStaleAfter,
		}

		outcome, badFor, viols := classifyRule(def, t, polls, pol.From, windowEnd, badStates, pol.Preexisting)
		if cov.Unobservable {
			outcome = OutcomeNotVerified
			anyUnobservable = true
			unobservableNames = append(unobservableNames, fmt.Sprintf("%s (%s)", def.Title, cov.Reason))
		}
		result.Violations = append(result.Violations, viols...)
		result.Verdicts = append(result.Verdicts, RuleVerdict{
			Alert: def.Title, RuleUID: def.UID, Outcome: outcome, BadFor: badFor,
			PollEvery: t.pollEvery, Note: strings.Join(cov.Notes, "; "),
		})
	}

	// MinObserved defaults to len(defs) (post-collapse). A shortfall counts
	// toward exit 1, never exit 2, and surfaces through Violations — so it
	// always produces at least one, even when no rule is paused.
	counted := watchedCount
	var attributable []Definition
	if pol.AllowPaused {
		counted += len(pausedRules)
	} else {
		attributable = pausedRules
	}
	if shortfall := minObserved - counted; shortfall > 0 {
		attributed := 0
		for _, def := range attributable {
			if attributed >= shortfall {
				break
			}
			// The paused rule and --allow-paused must both be named to the
			// user; both live in this one Violation, in Note — the renderer
			// prints Note verbatim rather than re-deriving the hint, so the
			// exact wording here is what an operator reads.
			result.Violations = append(result.Violations, Violation{
				Alert: def.Title, RuleUID: def.UID, Outcome: OutcomePaused,
				Note: "paused before the window opened; counts against --min-observed unless --allow-paused is set",
			})
			attributed++
		}
		for ; attributed < shortfall; attributed++ {
			// No named rule explains this part of the deficit — e.g. an
			// operator-supplied --min-observed above what could ever be
			// resolved. Note, not LastError: LastError is reporting-only
			// rule state read from a real poll, and this Violation never
			// touched one.
			result.Violations = append(result.Violations, Violation{
				Outcome: OutcomeNotCounted,
				Note:    fmt.Sprintf("min-observed %d exceeds the %d rule(s) counted as observed", minObserved, counted),
			})
		}
	}

	if anyUnobservable {
		return result, fmt.Errorf("gate: %d rule(s) not verified: %s", len(unobservableNames), strings.Join(unobservableNames, "; "))
	}
	return result, nil
}
