package gate

import (
	"fmt"
	"slices"
	"strings"
	"time"
)

// ReasonNodata is decide's own unobservable reason (§10.1/§10.2): proveCoverage
// (P7) deliberately never sets it — health=nodata is a note there, never fatal,
// because escalating it needs Policy.NodataIsUnobservable, and the pure
// coverage layer has no Policy to consult (coverage.go, check 5). decide is
// the seam that DOES have a Policy, so the escalation lives here.
const ReasonNodata UnobservableReason = "nodata"

// Outcome is the verdict of one instance's timeline, and — after decide takes
// the worst across a rule's instances — of the rule itself (§9). It is a
// published JSON output (§19.0): the three fail values stay distinct even
// though v1 maps all three to exit 1, because a later reason string cannot
// recover the information a single "fail" value would have thrown away, and
// because splitting them later would break a published interface for no gain.
type Outcome string

const (
	OutcomeClean           Outcome = "clean"
	OutcomeNewlyBad        Outcome = "newly_bad"
	OutcomeRecovered       Outcome = "recovered"
	OutcomePersistentlyBad Outcome = "persistently_bad"
	OutcomeFlapping        Outcome = "flapping"
	OutcomeSkipped         Outcome = "skipped"
	OutcomeUnobservable    Outcome = "unobservable"
)

// PreexistingPolicy governs only the ONE ambiguous case in the outcome table:
// an instance that was already bad when the window opened. A newly_bad or
// flapping instance is a fail under every policy (§11.3) — the plan lists
// them among the outcomes that "do not change" — so this type only ever
// changes how `recovered` and `persistently_bad` are judged (isViolation
// below).
type PreexistingPolicy string

const (
	// PreexistingFailUnlessRecovered is the default (§11.7): a preexisting
	// instance that clears and stays clear is a pass (`recovered`); one that
	// never clears is still a fail (`persistently_bad`).
	PreexistingFailUnlessRecovered PreexistingPolicy = "fail-unless-recovered"
	// PreexistingFail makes ANY preexisting instance a fail, even one that
	// recovers — for a user who wants no benefit of the doubt for a
	// condition this release did not cause.
	PreexistingFail PreexistingPolicy = "fail"
	// PreexistingIgnore disregards a preexisting instance entirely, whether
	// it recovers or stays bad for the whole window: only a genuinely NEW
	// bad episode (newly_bad or flapping) can fail the rule.
	PreexistingIgnore PreexistingPolicy = "ignore"
)

// Violation is one instance whose timeline outcome counts against the run,
// after the preexisting policy has been applied (isViolation below).
type Violation struct {
	Alert, RuleUID string
	Outcome        Outcome
	State          State
	Health         string // raw, reporting-only, like Poll.Health (P1.2a)
	LastError      string
	// FirstSeen is the episode's onset, in the runner domain (§16): activeAt
	// translated by its poll's own skew when the episode opened strictly
	// inside the window, or `from` itself when the instance was already bad
	// at window-open (preexisting) — never a raw, untranslated Grafana
	// timestamp.
	FirstSeen time.Time
	// ClearedAt is zero unless the episode closed via a genuine Cleared
	// event, also translated to the runner domain.
	ClearedAt      time.Time
	InstanceLabels map[string]string
	// Note carries an explanation for a Violation that has no instance
	// behind it — the synthetic MinObserved shortfall entry decide emits
	// when the deficit exceeds what any named paused rule explains (§12).
	// LastError is reporting-only rule state from a real poll and must not
	// double as a message field for a Violation that never touched one.
	Note string
}

// RuleVerdict is one rule's worst-of outcome (§9), always present for every
// resolved rule — Verdicts includes the passes, not only the failures — so a
// human reading the table sees every alert that was asked for, not only the
// ones that misbehaved.
type RuleVerdict struct {
	Alert, RuleUID string
	Outcome        Outcome
	BadFor         time.Duration // total wall-clock time any instance was bad inside the window, overlaps merged
	PollEvery      time.Duration
	Note           string
}

// Policy is decide's narrowed, pure-layer view of Config/Cfg (§9's P9
// comment): the classification knobs and the window, nothing else. No URL,
// no token, no I/O handles — those never reach the pure layer.
type Policy struct {
	States                            []State
	Preexisting                       PreexistingPolicy
	MinObserved                       int
	AllowPaused, NodataIsUnobservable bool
	From, To                          time.Time
}

// Result is decide's whole answer: everything §20.2's table and the action's
// JSON outputs need. Coverage carries one CoverageResult per non-skipped
// rule — no separate Interval type anywhere in the project (§2's
// simplification table).
type Result struct {
	From, To       time.Time
	GrafanaVersion string
	ClockSkew      time.Duration // the largest |skew| across every poll decide was given, not only the ones a rule's window actually used
	Coverage       map[string]CoverageResult
	Verdicts       []RuleVerdict
	Violations     []Violation
}

// episode is one contiguous, policy-bad span of one instance's timeline,
// already resolved to the runner domain and clamped to [from, windowEnd]. It
// never crosses a genuine Cleared event (H2): a Vanished marker freezes the
// state instead of closing the episode, which is what keeps a vanish from
// ever reading as a recovery.
type episode struct {
	start, end        time.Time
	closedByRealClear bool
}

// instanceTimeline accumulates one instance's walk across a rule's in-window
// polls. preexisting is decided once, the first time this key is seen bad:
// by the translated ActiveAt against `from` (§16), never by which poll
// happened to report it first — a poll's own cadence is not evidence of when
// the condition actually began (F1/F2).
type instanceTimeline struct {
	labels       map[string]string
	preexisting  bool
	seen         bool
	badOpen      bool
	episodeStart time.Time
	lastState    State
	lastHealth   string
	lastError    string
	episodes     []episode
}

// runnerTime translates a Grafana-domain timestamp recorded on poll p into
// the runner domain, undoing that poll's own measured skew (§16). GrafanaNow
// and ActiveAt come from the same response, so the same poll's skew applies
// to both. This is the single implementation of that translation for the
// package (same drift argument as pollsForRule, F5): coverage.go's window
// membership test and heartbeat boundary segments call it too, rather than
// each keeping its own copy of `p.GrafanaNow.Add(-p.Skew())` that could
// silently diverge from this one.
func runnerTime(p Poll, grafanaDomain time.Time) time.Time {
	return grafanaDomain.Add(-p.Skew())
}

// classifyRule builds every instance timeline for one rule across
// [from, windowEnd] and reduces them to the rule's worst outcome (§9), its
// merged BadFor, and the Violations the preexisting policy actually charges
// against the run. It is PURE: no I/O, no clock reads (§2) — decide supplies
// windowEnd (to + transitionGrace) rather than this function deriving it, so
// a test can pin the boundary directly.
//
// polls need not be pre-filtered to this rule, matching proveCoverage's own
// contract (§14.5): selection is by def.UID.
func classifyRule(def Definition, polls []Poll, from, windowEnd time.Time, badStates map[State]bool, pol PreexistingPolicy) (Outcome, time.Duration, []Violation) {
	rulePolls := pollsForRule(polls, def.UID)
	inWindow := inWindowPolls(rulePolls, from, windowEnd)

	timelines := make(map[string]*instanceTimeline)
	order := make([]string, 0)

	// get backfills labels the first time a real Instance is seen (F4): a key
	// can be created earlier by a bare Cleared/Vanished marker, which carries
	// no labels of its own, and the instance later re-firing must not report
	// an empty InstanceLabels just because of which event happened to create
	// the timeline first.
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
	closeEpisode := func(tl *instanceTimeline, end time.Time, real bool) {
		// inWindowPolls admits a poll whose translated time is up to its own
		// skew bound PAST windowEnd (the membership test widens the boundary
		// outward, §16). Without this clamp a genuine Cleared event on such a
		// poll would produce an episode.end slightly beyond windowEnd,
		// contradicting the episode type's own "clamped to
		// [from, windowEnd]" contract.
		if end.After(windowEnd) {
			end = windowEnd
		}
		// Different polls can carry different measured skews. In theory a
		// closing poll's translated time could land before the opening
		// poll's — skew is capped at skewHardLimit (60s), so this is remote,
		// not impossible — and a negative span would feed mergeDurations a
		// duration that subtracts instead of adds. Clamp rather than trust
		// the arithmetic never to invert.
		if end.Before(tl.episodeStart) {
			end = tl.episodeStart
		}
		tl.episodes = append(tl.episodes, episode{start: tl.episodeStart, end: end, closedByRealClear: real})
		tl.badOpen = false
	}
	// onsetOf resolves a fresh episode's start: the instance's own ActiveAt,
	// translated to the runner domain by this poll's skew, clamped so it
	// never reads as starting before the window opened.
	onsetOf := func(p Poll, inst Instance) time.Time {
		start := runnerTime(p, inst.ActiveAt)
		if start.Before(from) {
			start = from
		}
		return start
	}

	for _, p := range inWindow {
		byKey := make(map[string]Instance, len(p.Abnormal))
		for _, inst := range p.Abnormal {
			byKey[instanceKey(inst.Labels)] = inst
		}

		for key, inst := range byKey {
			tl := get(key, inst.Labels)
			bad := badStates[inst.State]
			switch {
			case !tl.seen:
				tl.seen = true
				if bad {
					// Fail-closed (§16): only call an onset "preexisting"
					// when even the worst-case skew error still puts it at
					// or before `from`. An onset that might really have
					// landed just inside the window must classify as a new
					// episode, never earn the `recovered` benefit of the
					// doubt it would get if it later clears (F1/F2).
					activeAtRunner := runnerTime(p, inst.ActiveAt)
					tl.preexisting = !activeAtRunner.Add(p.SkewBound()).After(from)
					if tl.preexisting {
						openEpisode(tl, from)
					} else {
						openEpisode(tl, onsetOf(p, inst))
					}
				}
			case bad && !tl.badOpen:
				openEpisode(tl, onsetOf(p, inst))
			case !bad && tl.badOpen:
				closeEpisode(tl, runnerTime(p, p.GrafanaNow), true)
			}
			tl.lastState, tl.lastHealth, tl.lastError = inst.State, p.Health, p.LastError
		}

		for _, key := range p.Cleared {
			tl := get(key, nil)
			if !tl.seen {
				// Cleared on the very first mention means the transition
				// happened between the poll just before this one (possibly
				// pre-window) and this one: there is no window-internal
				// evidence that it was ever bad, so it is neither
				// preexisting nor a new episode.
				tl.seen = true
				continue
			}
			if tl.badOpen {
				closeEpisode(tl, runnerTime(p, p.GrafanaNow), true)
			}
			tl.lastHealth, tl.lastError = p.Health, p.LastError
		}

		// Vanished is a deliberate no-op (H2): freeze whatever badOpen/preexisting
		// already holds. An instance that vanishes while bad must stay bad, and
		// one that vanishes while never having been bad must stay uninteresting.
		for _, key := range p.Vanished {
			tl := get(key, nil)
			tl.seen = true
			tl.lastHealth = p.Health
		}
	}

	// Multiple instances can appear for the first time within the same poll,
	// and map iteration order is nondeterministic; sort so this pure
	// function's Violations/BadFor output is stable across runs given the
	// same input, like log.go sorts Cleared/Vanished for the same reason.
	slices.Sort(order)

	var (
		outcome Outcome = OutcomeClean
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
			instOutcome = OutcomeFlapping
		case tl.preexisting:
			if tl.episodes[0].closedByRealClear {
				instOutcome = OutcomeRecovered
			} else {
				instOutcome = OutcomePersistentlyBad
			}
		default:
			// A genuinely new onset always fails, whether or not it later
			// clears within the window (§11.4 point 3): only a PREEXISTING
			// condition earns the benefit of `recovered`.
			instOutcome = OutcomeNewlyBad
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
// once the preexisting policy is applied. newly_bad and flapping always do
// (§11.3): both contain a genuinely new bad episode, so no policy forgives
// them. recovered and persistently_bad are, by classifyRule's construction,
// ALWAYS preexisting (a non-preexisting single episode is newly_bad instead,
// regardless of whether it clears) — so these are the only two policy can
// change, and isViolation needs no separate preexisting flag to know that.
func isViolation(o Outcome, pol PreexistingPolicy) bool {
	switch o {
	case OutcomeNewlyBad, OutcomeFlapping:
		return true
	case OutcomePersistentlyBad:
		return pol != PreexistingIgnore
	case OutcomeRecovered:
		return pol == PreexistingFail
	default:
		return false
	}
}

// outcomeRank orders outcomes for classifyRule's worst-of reduction across a
// rule's instances (§9). The three fail values, and recovered above clean,
// give it exactly the ordering the table requires —
// "unobservable > {flapping, persistently_bad, newly_bad} > recovered >
// skipped > clean" — with unobservable and skipped applied outside this
// function (decide owns both: unobservable from CoverageResult, skipped from
// Definition.IsPaused). The table does not distinguish among the three fail
// values, so their relative order here (flapping above persistently_bad
// above newly_bad) is an arbitrary but fixed and documented tie-break, not a
// claim that one is worse than another.
func outcomeRank(o Outcome) int {
	switch o {
	case OutcomeFlapping:
		return 4
	case OutcomePersistentlyBad:
		return 3
	case OutcomeNewlyBad:
		return 2
	case OutcomeRecovered:
		return 1
	default: // OutcomeClean
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
// same selection proveCoverage uses (§14.5: selection is by UID, never by
// title) — stable, because two polls sharing a coarse Date header must not
// reorder nondeterministically in a pure function. This is the single
// filter+sort implementation for the package (F5): proveCoverage calls it
// too, rather than keeping its own copy that could silently drift from this
// one's membership test.
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

// badStateSet turns Policy.States into a lookup set, defaulting to {firing}
// (§13) when the caller leaves States empty — decide applies the default
// itself so a test can pass a zero-value Policy and get v1's real default,
// rather than relying on a CLI layer that does not exist yet.
func badStateSet(states []State) map[State]bool {
	if len(states) == 0 {
		states = []State{StateFiring}
	}
	set := make(map[State]bool, len(states))
	for _, s := range states {
		set[s] = true
	}
	return set
}

// decide is the pure seam between the collected evidence and the CLI's exit
// code: nearly every §22 test targets this function, not Check (P9). It
// combines proveCoverage's nine checks with classifyRule's timelines under
// one Policy, and OWNS the H6 mapping: any unobservable rule makes decide
// return a non-nil error, which P10's CLI maps to exit 2 unconditionally
// (H7) — never to 0 or 1, and never suppressed by a real violation found
// alongside it.
//
// Result is fully populated even when the returned error is non-nil: H7's
// "err != nil, the violation list is irrelevant" means the CALLER must not
// use Violations to second-guess the error, not that Result stops being
// useful for the human table on exit 2.
func decide(h Header, polls []Poll, sentinel *time.Time, defs []Definition,
	rt map[string]ruleTimings, gt globalTimings, pol Policy) (Result, error) {

	badStates := badStateSet(pol.States)

	result := Result{
		From:           pol.From,
		To:             pol.To,
		GrafanaVersion: h.GrafanaVersion,
		Coverage:       make(map[string]CoverageResult),
	}
	for _, p := range polls {
		s := p.Skew()
		if s < 0 {
			s = -s
		}
		if s > result.ClockSkew {
			result.ClockSkew = s
		}
	}

	minObserved := pol.MinObserved
	if minObserved == 0 {
		minObserved = len(defs)
	}

	windowEnd := pol.To.Add(gt.transitionGrace)

	var (
		skippedRules      []Definition
		watchedCount      int
		anyUnobservable   bool
		unobservableNames []string
	)

	for _, def := range defs {
		if def.IsPaused {
			skippedRules = append(skippedRules, def)
			result.Verdicts = append(result.Verdicts, RuleVerdict{
				Alert: def.Title, RuleUID: def.UID, Outcome: OutcomeSkipped,
				PollEvery: rt[def.UID].pollEvery,
				Note:      "paused before the window opened",
			})
			continue
		}
		watchedCount++

		t := rt[def.UID]
		cov := proveCoverage(h, polls, sentinel, t, def, pol.From, pol.To, gt.transitionGrace)

		if pol.NodataIsUnobservable && !cov.Unobservable {
			inWindow := inWindowPolls(pollsForRule(polls, def.UID), pol.From, windowEnd)
			if runLen, sawAny := longestHealthRun(inWindow, "nodata"); sawAny && runLen > t.healthGrace {
				cov.Unobservable = true
				cov.Proved = false
				if cov.Reason == "" {
					cov.Reason = ReasonNodata
				}
				cov.Notes = append(cov.Notes, fmt.Sprintf(
					"rule %q: health=nodata for %s exceeds healthGrace %s and --nodata-is-unobservable is set",
					def.Title, runLen, t.healthGrace))
			}
		}
		result.Coverage[def.UID] = cov

		outcome, badFor, viols := classifyRule(def, polls, pol.From, windowEnd, badStates, pol.Preexisting)
		if cov.Unobservable {
			outcome = OutcomeUnobservable
			anyUnobservable = true
			unobservableNames = append(unobservableNames, fmt.Sprintf("%s (%s)", def.Title, cov.Reason))
		}
		result.Violations = append(result.Violations, viols...)
		result.Verdicts = append(result.Verdicts, RuleVerdict{
			Alert: def.Title, RuleUID: def.UID, Outcome: outcome, BadFor: badFor,
			PollEvery: t.pollEvery, Note: strings.Join(cov.Notes, "; "),
		})
	}

	// MinObserved (§12): default len(defs) after the collapse (already done
	// by Resolve before decide ever sees defs). skipped rules count against
	// it unless AllowPaused says otherwise. A shortfall counts toward exit 1
	// (§9.1), never exit 2 — decide never returns an error for this — and H7
	// requires it to surface through Violations like any other fail reason,
	// so a shortfall always produces at least one, even when no rule is
	// paused at all (an operator-supplied MinObserved that simply exceeds
	// what could ever be resolved).
	counted := watchedCount
	var attributable []Definition
	if pol.AllowPaused {
		counted += len(skippedRules)
	} else {
		attributable = skippedRules
	}
	if shortfall := minObserved - counted; shortfall > 0 {
		attributed := 0
		for _, def := range attributable {
			if attributed >= shortfall {
				break
			}
			// §12.1 requires the paused rule and --allow-paused both be
			// named to the user; naming the rule is this Violation's job,
			// the --allow-paused hint is the CLI table/renderer's (P10) —
			// tracked here so it is not dropped when that phase is built.
			result.Violations = append(result.Violations, Violation{
				Alert: def.Title, RuleUID: def.UID, Outcome: OutcomeSkipped,
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
				Outcome: OutcomeSkipped,
				Note:    fmt.Sprintf("min-observed %d exceeds the %d rule(s) counted as observed", minObserved, counted),
			})
		}
	}

	if anyUnobservable {
		return result, fmt.Errorf("gate: %d rule(s) unobservable: %s", len(unobservableNames), strings.Join(unobservableNames, "; "))
	}
	return result, nil
}
