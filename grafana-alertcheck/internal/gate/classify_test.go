package gate

import (
	"testing"
	"time"
)

func lbl(name string) map[string]string { return map[string]string{"instance": name} }

// abnormalPoll builds one Poll carrying a single abnormal instance, with the
// bookkeeping classifyRule needs (RuleUID, GrafanaNow, Health, Abnormal).
func abnormalPoll(uid string, at time.Time, state State, labels map[string]string, activeAt time.Time) Poll {
	return Poll{
		RuleUID:        uid,
		GrafanaNow:     at,
		Found:          true,
		Health:         "ok",
		LastEvaluation: at,
		Abnormal:       []Instance{{Labels: labels, State: state, ActiveAt: activeAt}},
	}
}

func clearedPoll(uid string, at time.Time, cleared ...string) Poll {
	return Poll{RuleUID: uid, GrafanaNow: at, Found: true, Health: "ok", LastEvaluation: at, Cleared: cleared}
}

func vanishedPoll(uid string, at time.Time, vanished ...string) Poll {
	return Poll{RuleUID: uid, GrafanaNow: at, Found: true, Health: "ok", LastEvaluation: at, Vanished: vanished}
}

func quietPoll(uid string, at time.Time) Poll {
	return Poll{RuleUID: uid, GrafanaNow: at, Found: true, Health: "ok", LastEvaluation: at}
}

var defaultBad = badStateSet(nil) // {firing}

// pausedHeader builds the header decide reads `skipped` from: the pause state
// as of record start. Definition.IsPaused is deliberately NOT that authority
// — it comes from a ruler read taken after the window closed — so a test that
// wants a rule treated as skipped must say so HERE (Header.pausedAtStart).
func pausedHeader(startedAt time.Time, pausedUIDs ...string) Header {
	h := Header{SchemaVersion: LogSchemaVersion, StartedAt: startedAt}
	for _, uid := range pausedUIDs {
		h.Rules = append(h.Rules, LoggedRule{UID: uid, IsPaused: true})
	}
	return h
}

// --- clean / newly_bad ---

func TestClassifyRule_NoEvidenceIsClean(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(10 * time.Minute)
	def := Definition{UID: "r1", Title: "R1"}

	polls := []Poll{quietPoll("r1", from), quietPoll("r1", to)}
	outcome, badFor, viols := classifyRule(def, polls, from, to, defaultBad, PreexistingFailUnlessRecovered)
	if outcome != OutcomeClean || badFor != 0 || len(viols) != 0 {
		t.Fatalf("outcome=%v badFor=%v viols=%v, want clean/0/none", outcome, badFor, viols)
	}
}

func TestClassifyRule_NewOnsetInsideWindowIsNewlyBad(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(10 * time.Minute)
	onset := from.Add(5 * time.Minute)
	def := Definition{UID: "r1", Title: "R1"}

	polls := []Poll{
		quietPoll("r1", from),
		abnormalPoll("r1", onset, StateFiring, lbl("a"), onset),
		abnormalPoll("r1", to, StateFiring, lbl("a"), onset),
	}
	outcome, badFor, viols := classifyRule(def, polls, from, to, defaultBad, PreexistingFailUnlessRecovered)
	if outcome != OutcomeNewlyBad {
		t.Fatalf("outcome = %v, want newly_bad", outcome)
	}
	if want := to.Sub(onset); badFor != want {
		t.Fatalf("badFor = %v, want %v", badFor, want)
	}
	if len(viols) != 1 || viols[0].Outcome != OutcomeNewlyBad {
		t.Fatalf("viols = %+v, want exactly one newly_bad violation", viols)
	}
}

// TestClassifyRule_NewOnsetThatClearsStillFails pins §11.4 point 3: a
// genuinely new bad episode fails even if it clears again before the window
// ends — only a PREEXISTING condition earns the benefit of `recovered`.
func TestClassifyRule_NewOnsetThatClearsStillFails(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(10 * time.Minute)
	onset := from.Add(2 * time.Minute)
	clearAt := from.Add(3 * time.Minute)
	def := Definition{UID: "r1", Title: "R1"}

	polls := []Poll{
		quietPoll("r1", from),
		abnormalPoll("r1", onset, StateFiring, lbl("a"), onset),
		clearedPoll("r1", clearAt, instanceKey(lbl("a"))),
		quietPoll("r1", to),
	}
	outcome, _, viols := classifyRule(def, polls, from, to, defaultBad, PreexistingFailUnlessRecovered)
	if outcome != OutcomeNewlyBad {
		t.Fatalf("outcome = %v, want newly_bad even though it cleared", outcome)
	}
	if len(viols) != 1 {
		t.Fatalf("viols = %+v, want one violation", viols)
	}
}

// --- recovered / persistently_bad (preexisting) ---

func TestClassifyRule_PreexistingThatRecoversIsRecoveredAndNotAViolation(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(10 * time.Minute)
	clearAt := from.Add(8 * time.Minute)
	def := Definition{UID: "r1", Title: "R1"}
	key := instanceKey(lbl("a"))

	polls := []Poll{
		abnormalPoll("r1", from, StateFiring, lbl("a"), from.Add(-time.Hour)),
		clearedPoll("r1", clearAt, key),
		quietPoll("r1", to),
	}
	outcome, badFor, viols := classifyRule(def, polls, from, to, defaultBad, PreexistingFailUnlessRecovered)
	if outcome != OutcomeRecovered {
		t.Fatalf("outcome = %v, want recovered", outcome)
	}
	if want := clearAt.Sub(from); badFor != want {
		t.Fatalf("badFor = %v, want %v", badFor, want)
	}
	if len(viols) != 0 {
		t.Fatalf("viols = %+v, want none: default policy passes a recovered preexisting instance", viols)
	}
}

// §22.2's "late condition": bad for 58 of a 60-minute window, clear at
// minute 58, still passes with a large BadFor — never a fail against some
// derived deadline (e.g. "must clear before 90% of the window").
func TestClassifyRule_LateRecoveryPassesRegardlessOfHowLateItIs(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(60 * time.Minute)
	clearAt := from.Add(58 * time.Minute)
	def := Definition{UID: "r1", Title: "R1"}
	key := instanceKey(lbl("a"))

	polls := []Poll{
		abnormalPoll("r1", from, StateFiring, lbl("a"), from.Add(-time.Hour)),
		clearedPoll("r1", clearAt, key),
		quietPoll("r1", to),
	}
	outcome, badFor, viols := classifyRule(def, polls, from, to, defaultBad, PreexistingFailUnlessRecovered)
	if outcome != OutcomeRecovered {
		t.Fatalf("outcome = %v, want recovered even 58 minutes into a 60-minute window", outcome)
	}
	if want := clearAt.Sub(from); badFor != want {
		t.Fatalf("badFor = %v, want the full %v bad duration, not a value clamped against a deadline", badFor, want)
	}
	if len(viols) != 0 {
		t.Fatalf("viols = %+v, want none: there is no deadline a preexisting recovery must beat", viols)
	}
}

func TestClassifyRule_PreexistingStillBadAtWindowEndIsPersistentlyBad(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(10 * time.Minute)
	def := Definition{UID: "r1", Title: "R1"}

	polls := []Poll{
		abnormalPoll("r1", from, StateFiring, lbl("a"), from.Add(-time.Hour)),
		abnormalPoll("r1", to, StateFiring, lbl("a"), from.Add(-time.Hour)),
	}
	outcome, badFor, viols := classifyRule(def, polls, from, to, defaultBad, PreexistingFailUnlessRecovered)
	if outcome != OutcomePersistentlyBad {
		t.Fatalf("outcome = %v, want persistently_bad", outcome)
	}
	if badFor != to.Sub(from) {
		t.Fatalf("badFor = %v, want the full window %v", badFor, to.Sub(from))
	}
	if len(viols) != 1 || viols[0].Outcome != OutcomePersistentlyBad {
		t.Fatalf("viols = %+v, want one persistently_bad violation", viols)
	}
}

// --- flapping ---

func TestClassifyRule_ClearThenBadAgainIsFlapping(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(10 * time.Minute)
	def := Definition{UID: "r1", Title: "R1"}
	key := instanceKey(lbl("a"))

	polls := []Poll{
		abnormalPoll("r1", from, StateFiring, lbl("a"), from),
		clearedPoll("r1", from.Add(2*time.Minute), key),
		abnormalPoll("r1", from.Add(5*time.Minute), StateFiring, lbl("a"), from.Add(5*time.Minute)),
		quietPoll("r1", to),
	}
	outcome, _, viols := classifyRule(def, polls, from, to, defaultBad, PreexistingFailUnlessRecovered)
	if outcome != OutcomeFlapping {
		t.Fatalf("outcome = %v, want flapping", outcome)
	}
	if len(viols) != 1 || viols[0].Outcome != OutcomeFlapping {
		t.Fatalf("viols = %+v, want one flapping violation, always a fail regardless of policy", viols)
	}
}

// §22.2: "a clear and then a second bad state gives flapping, at each
// possible time of the second bad state." A table over where the second
// onset lands — immediately after the clear, mid-window, and right at the
// last instant before windowEnd — closes the boundary this single fixed
// timing above cannot.
func TestClassifyRule_FlappingAtEveryTimingOfTheSecondOnset(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(10 * time.Minute)
	clearAt := from.Add(2 * time.Minute)
	def := Definition{UID: "r1", Title: "R1"}
	key := instanceKey(lbl("a"))

	tests := []struct {
		name        string
		secondOnset time.Time
	}{
		{"immediately after the clear", clearAt.Add(time.Second)},
		{"mid-window", from.Add(5 * time.Minute)},
		{"the last instant before windowEnd", to.Add(-time.Second)},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			polls := []Poll{
				abnormalPoll("r1", from, StateFiring, lbl("a"), from),
				clearedPoll("r1", clearAt, key),
				abnormalPoll("r1", tc.secondOnset, StateFiring, lbl("a"), tc.secondOnset),
				quietPoll("r1", to),
			}
			outcome, _, viols := classifyRule(def, polls, from, to, defaultBad, PreexistingFailUnlessRecovered)
			if outcome != OutcomeFlapping {
				t.Fatalf("outcome = %v, want flapping for a second onset at %s", outcome, tc.secondOnset)
			}
			if len(viols) != 1 || viols[0].Outcome != OutcomeFlapping {
				t.Fatalf("viols = %+v, want one flapping violation", viols)
			}
		})
	}
}

// --- H2: vanished is a discontinuity, never a clear ---

func TestClassifyRule_VanishedWhileBadStaysPersistentlyBad(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(10 * time.Minute)
	def := Definition{UID: "r1", Title: "R1"}
	key := instanceKey(lbl("a"))

	polls := []Poll{
		abnormalPoll("r1", from, StateFiring, lbl("a"), from.Add(-time.Minute)),
		vanishedPoll("r1", from.Add(5*time.Minute), key),
		quietPoll("r1", to),
	}
	outcome, badFor, viols := classifyRule(def, polls, from, to, defaultBad, PreexistingFailUnlessRecovered)
	if outcome != OutcomePersistentlyBad {
		t.Fatalf("outcome = %v, want persistently_bad: a vanish must never read as a recovery (H2)", outcome)
	}
	if badFor != to.Sub(from) {
		t.Fatalf("badFor = %v, want the full window %v: the freeze must hold the episode open to windowEnd", badFor, to.Sub(from))
	}
	if len(viols) != 1 {
		t.Fatalf("viols = %+v, want one violation", viols)
	}
}

func TestClassifyRule_VanishedWhileNeverBadIsUninteresting(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(10 * time.Minute)
	def := Definition{UID: "r1", Title: "R1"}
	key := instanceKey(lbl("a"))

	// Pending is abnormal (non-normal) but not in the default {firing} bad
	// set, so its vanish must stay uninteresting too.
	polls := []Poll{
		abnormalPoll("r1", from, StatePending, lbl("a"), from),
		vanishedPoll("r1", from.Add(5*time.Minute), key),
		quietPoll("r1", to),
	}
	outcome, badFor, viols := classifyRule(def, polls, from, to, defaultBad, PreexistingFailUnlessRecovered)
	if outcome != OutcomeClean || badFor != 0 || len(viols) != 0 {
		t.Fatalf("outcome=%v badFor=%v viols=%v, want clean/0/none", outcome, badFor, viols)
	}
}

// --- preexisting policy ---

func TestClassifyRule_PreexistingPolicyFailFailsARecoveredInstance(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(10 * time.Minute)
	def := Definition{UID: "r1", Title: "R1"}
	key := instanceKey(lbl("a"))

	polls := []Poll{
		abnormalPoll("r1", from, StateFiring, lbl("a"), from.Add(-time.Hour)),
		clearedPoll("r1", from.Add(2*time.Minute), key),
		quietPoll("r1", to),
	}
	outcome, _, viols := classifyRule(def, polls, from, to, defaultBad, PreexistingFail)
	if outcome != OutcomeRecovered {
		t.Fatalf("outcome = %v, want recovered — the descriptive outcome does not change under policy=fail", outcome)
	}
	if len(viols) != 1 || viols[0].Outcome != OutcomeRecovered {
		t.Fatalf("viols = %+v, want one violation: policy=fail gives no benefit of the doubt to a preexisting instance", viols)
	}
}

func TestClassifyRule_PreexistingPolicyIgnoreForgivesPersistentlyBad(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(10 * time.Minute)
	def := Definition{UID: "r1", Title: "R1"}

	polls := []Poll{
		abnormalPoll("r1", from, StateFiring, lbl("a"), from.Add(-time.Hour)),
		abnormalPoll("r1", to, StateFiring, lbl("a"), from.Add(-time.Hour)),
	}
	outcome, _, viols := classifyRule(def, polls, from, to, defaultBad, PreexistingIgnore)
	if outcome != OutcomePersistentlyBad {
		t.Fatalf("outcome = %v, want persistently_bad — the descriptive outcome does not change under policy=ignore", outcome)
	}
	if len(viols) != 0 {
		t.Fatalf("viols = %+v, want none: policy=ignore disregards a preexisting instance even if it never recovers", viols)
	}
}

func TestClassifyRule_PreexistingPolicyIgnoreStillFailsANewOnset(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(10 * time.Minute)
	onset := from.Add(5 * time.Minute)
	def := Definition{UID: "r1", Title: "R1"}

	polls := []Poll{
		quietPoll("r1", from),
		abnormalPoll("r1", onset, StateFiring, lbl("a"), onset),
		abnormalPoll("r1", to, StateFiring, lbl("a"), onset),
	}
	outcome, _, viols := classifyRule(def, polls, from, to, defaultBad, PreexistingIgnore)
	if outcome != OutcomeNewlyBad || len(viols) != 1 {
		t.Fatalf("outcome=%v viols=%v, want newly_bad/1: ignore only forgives PREEXISTING badness", outcome, viols)
	}
}

// --- worst-of across instances ---

func TestClassifyRule_WorstOfMultipleInstancesWins(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(10 * time.Minute)
	def := Definition{UID: "r1", Title: "R1"}

	polls := []Poll{
		{
			RuleUID: "r1", GrafanaNow: from, Found: true, Health: "ok",
			Abnormal: []Instance{
				{Labels: lbl("a"), State: StateFiring, ActiveAt: from.Add(-time.Hour)}, // preexisting, will recover
				{Labels: lbl("b"), State: StateFiring, ActiveAt: from},                 // preexisting, will stay bad
			},
		},
		clearedPoll("r1", from.Add(2*time.Minute), instanceKey(lbl("a"))),
		{
			RuleUID: "r1", GrafanaNow: to, Found: true, Health: "ok",
			Abnormal: []Instance{{Labels: lbl("b"), State: StateFiring, ActiveAt: from}},
		},
	}
	outcome, _, viols := classifyRule(def, polls, from, to, defaultBad, PreexistingFailUnlessRecovered)
	if outcome != OutcomePersistentlyBad {
		t.Fatalf("outcome = %v, want persistently_bad: the worse of {recovered, persistently_bad}", outcome)
	}
	if len(viols) != 1 || viols[0].Outcome != OutcomePersistentlyBad {
		t.Fatalf("viols = %+v, want exactly the persistently_bad instance's violation", viols)
	}
}

// --- decide(): skipped rules, unobservable (H6), MinObserved, exit mapping (H7) ---

func TestDecide_SkippedRuleNeverReachesProveCoverage(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(10 * time.Minute)
	def := Definition{UID: "r1", Title: "R1", IsPaused: true}
	defs := []Definition{def}
	rt := map[string]ruleTimings{"r1": newRuleTimings(30*time.Second, 60)}
	gt := globalTimings{}
	pol := Policy{From: from, To: to, AllowPaused: true}

	// No polls, no sentinel at all: a heartbeat_gap/no_sentinel misclassification
	// here would mean proveCoverage ran for a skipped rule (§4.3's obligation).
	// The HEADER is what says paused — decide reads skipped from there, not
	// from def.IsPaused, which is a post-window reading (Header.pausedAtStart).
	res, err := decide(pausedHeader(from.Add(-time.Hour), "r1"), nil, nil, defs, rt, gt, pol)
	if err != nil {
		t.Fatalf("err = %v, want nil: a rule paused before the window is skipped, not unobservable", err)
	}
	if len(res.Verdicts) != 1 || res.Verdicts[0].Outcome != OutcomeSkipped {
		t.Fatalf("Verdicts = %+v, want exactly one skipped verdict", res.Verdicts)
	}
	if _, ok := res.Coverage["r1"]; ok {
		t.Fatalf("Coverage[r1] present, want absent: a skipped rule has no coverage to prove")
	}
}

func TestDecide_UnobservableRuleAlwaysReturnsAnError(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(10 * time.Minute)
	def := Definition{UID: "r1", Title: "R1"}
	defs := []Definition{def}
	rt := map[string]ruleTimings{"r1": newRuleTimings(30*time.Second, 60)}
	gt := globalTimings{}
	pol := Policy{From: from, To: to}

	// No sentinel at all: check 1 fails, so the rule is unobservable
	// regardless of anything else.
	res, err := decide(Header{StartedAt: from.Add(-time.Hour)}, nil, nil, defs, rt, gt, pol)
	if err == nil {
		t.Fatalf("err = nil, want non-nil: H6/H7 require an unobservable rule to always fail the run")
	}
	if len(res.Verdicts) != 1 || res.Verdicts[0].Outcome != OutcomeUnobservable {
		t.Fatalf("Verdicts = %+v, want exactly one unobservable verdict", res.Verdicts)
	}
}

// TestDecide_UnobservableWinsEvenAlongsideARealViolation pins H6 exactly:
// "Any unobservable rule -> exit 2, no exception, even alongside a real
// newly_bad."
func TestDecide_UnobservableWinsEvenAlongsideARealViolation(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(10 * time.Minute)
	onset := from.Add(5 * time.Minute)

	defBroken := Definition{UID: "broken", Title: "Broken"}
	defBad := Definition{UID: "bad", Title: "Bad"}
	defs := []Definition{defBroken, defBad}
	rt := map[string]ruleTimings{
		"broken": newRuleTimings(30*time.Second, 60),
		"bad":    newRuleTimings(30*time.Second, 60),
	}
	gt := globalTimings{}
	pol := Policy{From: from, To: to}

	var polls []Poll
	for ts := from; !ts.After(to); ts = ts.Add(30 * time.Second) {
		if ts.Equal(onset) || ts.After(onset) {
			polls = append(polls, abnormalPoll("bad", ts, StateFiring, lbl("a"), onset))
		} else {
			polls = append(polls, quietPoll("bad", ts))
		}
	}
	// "broken" gets no polls at all: no sentinel, no heartbeats -> unobservable.
	sentinel := to

	res, err := decide(Header{StartedAt: from.Add(-time.Hour)}, polls, &sentinel, defs, rt, gt, pol)
	if err == nil {
		t.Fatalf("err = nil, want non-nil: one rule is unobservable")
	}
	var gotBroken, gotBad Outcome
	for _, v := range res.Verdicts {
		switch v.RuleUID {
		case "broken":
			gotBroken = v.Outcome
		case "bad":
			gotBad = v.Outcome
		}
	}
	if gotBroken != OutcomeUnobservable {
		t.Fatalf("broken.Outcome = %v, want unobservable", gotBroken)
	}
	if gotBad != OutcomeNewlyBad {
		t.Fatalf("bad.Outcome = %v, want newly_bad: classification still runs and is still visible in Verdicts (H5)", gotBad)
	}
	if len(res.Violations) == 0 {
		t.Fatalf("Violations empty, want the newly_bad instance still reported even though the run fails on the unobservable rule")
	}
}

// §22.10: "a clean verdict with a coverage gap ... must never give exit 0",
// and "a recovered verdict and a skipped verdict also need proved coverage
// of the full window." One genuinely unobservable rule ("broken", zero
// polls) alongside a rule with each of the three favorable outcomes — none
// of them may waive the run.
func TestDecide_UnobservableRuleWinsOverEveryFavorableOutcome(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(10 * time.Minute)

	tests := []struct {
		name          string
		goodPolls     []Poll
		pausedAtStart bool
		wantOutcome   Outcome
	}{
		{
			name: "clean",
			goodPolls: func() []Poll {
				var polls []Poll
				for ts := from; !ts.After(to); ts = ts.Add(30 * time.Second) {
					polls = append(polls, quietPoll("good", ts))
				}
				return polls
			}(),
			wantOutcome: OutcomeClean,
		},
		{
			// Dense 30s-spaced polls throughout, so "good"'s own coverage
			// proves clean on its own — a sparse abnormal/cleared/quiet
			// triple (enough for classifyRule alone) would leave a
			// heartbeat gap that muddies which rule made the run fail.
			name: "recovered",
			goodPolls: func() []Poll {
				var polls []Poll
				clearAt := from.Add(3 * time.Minute)
				key := instanceKey(lbl("a"))
				cleared := false
				for ts := from; !ts.After(to); ts = ts.Add(30 * time.Second) {
					switch {
					case ts.Equal(clearAt):
						polls = append(polls, clearedPoll("good", ts, key))
						cleared = true
					case !cleared:
						polls = append(polls, abnormalPoll("good", ts, StateFiring, lbl("a"), from.Add(-time.Hour)))
					default:
						polls = append(polls, quietPoll("good", ts))
					}
				}
				return polls
			}(),
			wantOutcome: OutcomeRecovered,
		},
		{
			name:          "skipped",
			pausedAtStart: true,
			wantOutcome:   OutcomeSkipped,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			defs := []Definition{{UID: "good", Title: "Good"}, {UID: "broken", Title: "Broken"}}
			rt := map[string]ruleTimings{
				"good":   newRuleTimings(30*time.Second, 60),
				"broken": newRuleTimings(30*time.Second, 60),
			}
			gt := globalTimings{}
			pol := Policy{From: from, To: to}

			h := Header{StartedAt: from.Add(-time.Hour)}
			if tc.pausedAtStart {
				h.Rules = []LoggedRule{{UID: "good", IsPaused: true}}
			}

			// "broken" gets no polls at all: no sentinel-worthy heartbeats,
			// so it is unobservable regardless of "good".
			sentinel := to
			res, err := decide(h, tc.goodPolls, &sentinel, defs, rt, gt, pol)
			if err == nil {
				t.Fatalf("err = nil, want non-nil: 'broken' is unobservable regardless of 'good' being %s", tc.name)
			}
			var gotGood, gotBroken Outcome
			for _, v := range res.Verdicts {
				switch v.RuleUID {
				case "good":
					gotGood = v.Outcome
				case "broken":
					gotBroken = v.Outcome
				}
			}
			if gotGood != tc.wantOutcome {
				t.Errorf("good.Outcome = %v, want %v", gotGood, tc.wantOutcome)
			}
			if gotBroken != OutcomeUnobservable {
				t.Errorf("broken.Outcome = %v, want unobservable", gotBroken)
			}
		})
	}
}

// §22.10: the table above puts the coverage gap on a DIFFERENT rule from the
// one with the favorable outcome. This pins the tighter claim: a rule that
// itself recovers, but ALSO itself has a coverage gap, is still overridden to
// unobservable — the favorable classification of a rule is never a reason to
// skip that same rule's own coverage check.
func TestDecide_RecoveredOutcomeOverriddenByItsOwnCoverageGap(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(10 * time.Minute)
	def := Definition{UID: "r1", Title: "R1"}
	defs := []Definition{def}
	rt := map[string]ruleTimings{"r1": newRuleTimings(30*time.Second, 60)} // maxGap = 60s
	gt := globalTimings{}
	pol := Policy{From: from, To: to}
	key := instanceKey(lbl("a"))

	polls := []Poll{
		abnormalPoll("r1", from, StateFiring, lbl("a"), from.Add(-time.Hour)),
		clearedPoll("r1", from.Add(30*time.Second), key),
	}
	for ts := from.Add(time.Minute); !ts.After(to); ts = ts.Add(30 * time.Second) {
		// A gap from from+1.5m to from+4m — well past the 60s maxGap —
		// sitting entirely AFTER the clear, so classifyRule alone would
		// still call this rule `recovered`.
		if ts.After(from.Add(90*time.Second)) && ts.Before(from.Add(4*time.Minute)) {
			continue
		}
		polls = append(polls, quietPoll("r1", ts))
	}
	sentinel := to

	res, err := decide(Header{StartedAt: from.Add(-time.Hour)}, polls, &sentinel, defs, rt, gt, pol)
	if err == nil {
		t.Fatalf("err = nil, want non-nil: r1's own coverage gap must fail the run even though it recovered")
	}
	if len(res.Verdicts) != 1 || res.Verdicts[0].Outcome != OutcomeUnobservable {
		t.Fatalf("Verdicts = %+v, want unobservable, never recovered", res.Verdicts)
	}
	if cov := res.Coverage["r1"]; cov.Proved {
		t.Fatalf("Coverage = %+v, want not proved", cov)
	}
}

func TestDecide_CleanWindowIsAPass(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(10 * time.Minute)
	def := Definition{UID: "r1", Title: "R1"}
	defs := []Definition{def}
	rt := map[string]ruleTimings{"r1": newRuleTimings(30*time.Second, 60)}
	gt := globalTimings{}
	pol := Policy{From: from, To: to}

	var polls []Poll
	for ts := from; !ts.After(to); ts = ts.Add(30 * time.Second) {
		polls = append(polls, quietPoll("r1", ts))
	}
	sentinel := to

	res, err := decide(Header{StartedAt: from.Add(-time.Hour)}, polls, &sentinel, defs, rt, gt, pol)
	if err != nil {
		t.Fatalf("err = %v, want nil", err)
	}
	if len(res.Violations) != 0 {
		t.Fatalf("Violations = %+v, want none: H7 says a pass is exactly len(Violations)==0 && err==nil", res.Violations)
	}
	if res.Verdicts[0].Outcome != OutcomeClean {
		t.Fatalf("Outcome = %v, want clean", res.Verdicts[0].Outcome)
	}
}

// §22.7's second of the plan's "if only three tests could exist" cases: a
// pause and then an unpause inside the window, with an episode that would
// fire and resolve entirely inside the blind interval. A drain wait alone —
// "did the rule eventually evaluate through windowEnd?" — would see
// lastEvaluation catch up after the unpause and answer yes, a pass. decide()
// never runs a drain wait (that is check.go's I/O concern, §14.6); this pins
// that proveCoverage's own per-poll checks already refuse the window without
// one, so a live drain wait is not what is saving this case.
func TestDecide_PauseThenUnpauseWithHiddenEpisodeGivesUnobservableNotClean(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(20 * time.Minute)
	def := Definition{UID: "r1", Title: "R1"}
	defs := []Definition{def}
	rt := map[string]ruleTimings{"r1": newRuleTimings(30*time.Second, 60)}
	gt := globalTimings{}
	pol := Policy{From: from, To: to}

	pauseStart := from.Add(5 * time.Minute)
	pauseEnd := from.Add(10 * time.Minute)

	var polls []Poll
	for ts := from; !ts.After(pauseStart.Add(-30 * time.Second)); ts = ts.Add(30 * time.Second) {
		polls = append(polls, quietPoll("r1", ts))
	}
	for ts := pauseStart; !ts.After(pauseEnd); ts = ts.Add(30 * time.Second) {
		// No fire/resolve is ever observed here: the rule was not
		// evaluating, so any real episode inside this stretch is invisible
		// to every poll (§14.7).
		polls = append(polls, Poll{RuleUID: "r1", GrafanaNow: ts, Found: true, Health: "ok", IsPaused: true, LastEvaluation: pauseStart})
	}
	for ts := pauseEnd.Add(30 * time.Second); !ts.After(to); ts = ts.Add(30 * time.Second) {
		// Evaluations resume and catch straight up — a drain wait's final
		// "did it reach windowEnd" question would answer yes.
		polls = append(polls, quietPoll("r1", ts))
	}
	sentinel := to

	res, err := decide(Header{StartedAt: from.Add(-time.Hour)}, polls, &sentinel, defs, rt, gt, pol)
	if err == nil {
		t.Fatalf("err = nil, want the pause-then-unpause blind interval to fail closed")
	}
	if len(res.Verdicts) != 1 || res.Verdicts[0].Outcome != OutcomeUnobservable {
		t.Fatalf("Verdicts = %+v, want unobservable, never clean", res.Verdicts)
	}
	if cov := res.Coverage["r1"]; cov.Proved {
		t.Fatalf("Coverage = %+v, want not proved", cov)
	}
}

// --- MinObserved shortfall (§12) ---

func TestDecide_SkippedOnlyShortfallProducesAViolationWithoutAnError(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(10 * time.Minute)

	watched := Definition{UID: "watched", Title: "Watched"}
	paused := Definition{UID: "paused", Title: "Paused", IsPaused: true}
	defs := []Definition{watched, paused}
	rt := map[string]ruleTimings{
		"watched": newRuleTimings(30*time.Second, 60),
		"paused":  newRuleTimings(30*time.Second, 60),
	}
	gt := globalTimings{}
	// MinObserved defaults to len(defs) = 2, but only "watched" is observable.
	pol := Policy{From: from, To: to}

	var polls []Poll
	for ts := from; !ts.After(to); ts = ts.Add(30 * time.Second) {
		polls = append(polls, quietPoll("watched", ts))
	}
	sentinel := to

	res, err := decide(pausedHeader(from.Add(-time.Hour), "paused"), polls, &sentinel, defs, rt, gt, pol)
	if err != nil {
		t.Fatalf("err = %v, want nil: a shortfall caused only by a skipped rule is exit 1, not exit 2 (§9.1)", err)
	}
	if len(res.Violations) != 1 {
		t.Fatalf("Violations = %+v, want exactly one: H7 needs the shortfall visible through Violations to keep its equivalence", res.Violations)
	}
	if v := res.Violations[0]; v.Outcome != OutcomeSkipped || v.RuleUID != "paused" || v.Alert != "Paused" {
		t.Fatalf("Violations[0] = %+v, want Outcome=skipped naming the paused rule (§12.1: the message names the paused rule)", v)
	}
	if res.Violations[0].Note == "" {
		t.Fatalf("Violations[0].Note is empty, want an explanation: the shortfall reason must not be smuggled into LastError, " +
			"which is reporting-only rule state from a real poll this synthetic Violation never touched")
	}
}

// TestDecide_ExplicitMinObservedShortfallWithNoPausedRuleStillProducesAViolation
// pins F3: an operator-supplied MinObserved that exceeds what could ever be
// resolved is still a shortfall, even with zero paused rules to blame it on
// — H7 must not let this silently read as a pass.
func TestDecide_ExplicitMinObservedShortfallWithNoPausedRuleStillProducesAViolation(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(10 * time.Minute)

	def := Definition{UID: "r1", Title: "R1"}
	defs := []Definition{def}
	rt := map[string]ruleTimings{"r1": newRuleTimings(30*time.Second, 60)}
	gt := globalTimings{}
	pol := Policy{From: from, To: to, MinObserved: 3} // only one rule will ever be resolved

	var polls []Poll
	for ts := from; !ts.After(to); ts = ts.Add(30 * time.Second) {
		polls = append(polls, quietPoll("r1", ts))
	}
	sentinel := to

	res, err := decide(Header{StartedAt: from.Add(-time.Hour)}, polls, &sentinel, defs, rt, gt, pol)
	if err != nil {
		t.Fatalf("err = %v, want nil: an unmet MinObserved is exit 1, never exit 2", err)
	}
	if len(res.Violations) != 2 {
		t.Fatalf("Violations = %+v, want two: the shortfall (3-1=2) is not explained by any paused rule, "+
			"so H7 requires it to surface directly rather than pass silently", res.Violations)
	}
	for _, v := range res.Violations {
		if v.Outcome != OutcomeSkipped {
			t.Fatalf("Violations = %+v, want Outcome=skipped on the synthetic shortfall entries", res.Violations)
		}
	}
}

func TestDecide_AllowPausedSuppressesTheShortfall(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(10 * time.Minute)

	watched := Definition{UID: "watched", Title: "Watched"}
	paused := Definition{UID: "paused", Title: "Paused", IsPaused: true}
	defs := []Definition{watched, paused}
	rt := map[string]ruleTimings{
		"watched": newRuleTimings(30*time.Second, 60),
		"paused":  newRuleTimings(30*time.Second, 60),
	}
	gt := globalTimings{}
	pol := Policy{From: from, To: to, AllowPaused: true}

	var polls []Poll
	for ts := from; !ts.After(to); ts = ts.Add(30 * time.Second) {
		polls = append(polls, quietPoll("watched", ts))
	}
	sentinel := to

	res, err := decide(pausedHeader(from.Add(-time.Hour), "paused"), polls, &sentinel, defs, rt, gt, pol)
	if err != nil {
		t.Fatalf("err = %v, want nil", err)
	}
	if len(res.Violations) != 0 {
		t.Fatalf("Violations = %+v, want none: --allow-paused must suppress the shortfall entirely", res.Violations)
	}
}

// --- nodata escalation (decide's own Policy-driven check) ---

func TestDecide_NodataIsUnobservableEscalatesASustainedRun(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(10 * time.Minute)
	def := Definition{UID: "r1", Title: "R1"}
	defs := []Definition{def}
	rt := map[string]ruleTimings{"r1": newRuleTimings(30*time.Second, 60)} // healthGrace = max(60s,60s) = 60s
	gt := globalTimings{}
	pol := Policy{From: from, To: to, NodataIsUnobservable: true}

	var polls []Poll
	for ts := from; !ts.After(to); ts = ts.Add(30 * time.Second) {
		polls = append(polls, Poll{RuleUID: "r1", GrafanaNow: ts, Found: true, Health: "nodata", LastEvaluation: ts})
	}
	sentinel := to

	res, err := decide(Header{StartedAt: from.Add(-time.Hour)}, polls, &sentinel, defs, rt, gt, pol)
	if err == nil {
		t.Fatalf("err = nil, want non-nil: a sustained nodata run must be unobservable under --nodata-is-unobservable")
	}
	if res.Coverage["r1"].Reason != ReasonNodata {
		t.Fatalf("Reason = %q, want %q", res.Coverage["r1"].Reason, ReasonNodata)
	}
}

func TestDecide_NodataIsANoteByDefault(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(10 * time.Minute)
	def := Definition{UID: "r1", Title: "R1"}
	defs := []Definition{def}
	rt := map[string]ruleTimings{"r1": newRuleTimings(30*time.Second, 60)}
	gt := globalTimings{}
	pol := Policy{From: from, To: to} // NodataIsUnobservable defaults to false

	var polls []Poll
	for ts := from; !ts.After(to); ts = ts.Add(30 * time.Second) {
		polls = append(polls, Poll{RuleUID: "r1", GrafanaNow: ts, Found: true, Health: "nodata", LastEvaluation: ts})
	}
	sentinel := to

	res, err := decide(Header{StartedAt: from.Add(-time.Hour)}, polls, &sentinel, defs, rt, gt, pol)
	if err != nil {
		t.Fatalf("err = %v, want nil: 96%% of the fleet runs no_data_state:OK and must not fail by default", err)
	}
	if res.Coverage["r1"].Unobservable {
		t.Fatalf("Coverage[r1].Unobservable = true, want false by default")
	}
}

// --- F1/F2 regressions: preexisting is decided by ActiveAt, not poll timing ---

// TestClassifyRule_OnsetBetweenFromAndFirstPollIsNewlyBadNotRecovered pins
// F1: an instance whose true onset (ActiveAt) falls strictly inside the
// window — even though the first poll that happens to observe it already
// shows it bad — must never be treated as preexisting. If it then clears,
// the plan requires newly_bad (exit 1), not recovered (exit 0).
func TestClassifyRule_OnsetBetweenFromAndFirstPollIsNewlyBadNotRecovered(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(10 * time.Minute)
	onset := from.Add(1 * time.Minute)     // the true onset, strictly after `from`
	firstPoll := from.Add(2 * time.Minute) // the first poll that happens to observe it
	clearAt := from.Add(5 * time.Minute)
	def := Definition{UID: "r1", Title: "R1"}
	key := instanceKey(lbl("a"))

	polls := []Poll{
		abnormalPoll("r1", firstPoll, StateFiring, lbl("a"), onset),
		clearedPoll("r1", clearAt, key),
		quietPoll("r1", to),
	}
	outcome, badFor, viols := classifyRule(def, polls, from, to, defaultBad, PreexistingFailUnlessRecovered)
	if outcome != OutcomeNewlyBad {
		t.Fatalf("outcome = %v, want newly_bad: the onset is after `from`, so it is not preexisting even though "+
			"the FIRST in-window poll already observes it bad (F1)", outcome)
	}
	if len(viols) != 1 || viols[0].Outcome != OutcomeNewlyBad {
		t.Fatalf("viols = %+v, want one newly_bad violation: a policy=fail-unless-recovered default must still fail this", viols)
	}
	if want := clearAt.Sub(onset); badFor != want {
		t.Fatalf("badFor = %v, want %v: BadFor must count from the true onset, not from `from` (F1's overcount bug)", badFor, want)
	}
}

// TestClassifyRule_OnsetJustBeforeFromIsPreexisting is the mirror check: an
// onset at or before `from` (even if the first poll is later) is genuinely
// preexisting and, if it clears, is `recovered`.
func TestClassifyRule_OnsetJustBeforeFromIsPreexisting(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(10 * time.Minute)
	onset := from.Add(-time.Minute)
	firstPoll := from.Add(2 * time.Minute)
	clearAt := from.Add(5 * time.Minute)
	def := Definition{UID: "r1", Title: "R1"}
	key := instanceKey(lbl("a"))

	polls := []Poll{
		abnormalPoll("r1", firstPoll, StateFiring, lbl("a"), onset),
		clearedPoll("r1", clearAt, key),
		quietPoll("r1", to),
	}
	outcome, badFor, viols := classifyRule(def, polls, from, to, defaultBad, PreexistingFailUnlessRecovered)
	if outcome != OutcomeRecovered {
		t.Fatalf("outcome = %v, want recovered: the onset is at/before `from`, genuinely preexisting", outcome)
	}
	if len(viols) != 0 {
		t.Fatalf("viols = %+v, want none: default policy passes a recovered preexisting instance", viols)
	}
	if want := clearAt.Sub(from); badFor != want {
		t.Fatalf("badFor = %v, want %v: a preexisting episode's BadFor is clamped to window-open, not backdated past it", badFor, want)
	}
}

// TestClassifyRule_SkewTranslatesActiveAtAcrossTheWindowBoundary pins F2: a
// poll carrying a nonzero skew must have its ActiveAt (and GrafanaNow)
// translated to the runner domain before comparing against `from` — a raw,
// untranslated comparison would land on the wrong side of the F1 boundary
// check.
func TestClassifyRule_SkewTranslatesActiveAtAcrossTheWindowBoundary(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(10 * time.Minute)
	def := Definition{UID: "r1", Title: "R1"}

	// Grafana's clock reads 90s ahead of the runner's (skew = +90s). The
	// poll's raw GrafanaNow/ActiveAt both sit 90s past `from` in Grafana's
	// domain, but translate to exactly `from` in the runner domain — genuinely
	// preexisting once translated, and wrongly "newly_bad" if the skew is
	// ignored.
	skew := 90 * time.Second
	rawActiveAt := from.Add(skew)
	poll := Poll{
		RuleUID: "r1", GrafanaNow: from.Add(skew), Found: true, Health: "ok",
		LastEvaluation: from.Add(skew), SkewMS: skew.Milliseconds(),
		Abnormal: []Instance{{Labels: lbl("a"), State: StateFiring, ActiveAt: rawActiveAt}},
	}
	stillBad := poll
	stillBad.GrafanaNow = to.Add(skew)
	stillBad.LastEvaluation = to.Add(skew)

	outcome, badFor, _ := classifyRule(def, []Poll{poll, stillBad}, from, to, defaultBad, PreexistingFailUnlessRecovered)
	if outcome != OutcomePersistentlyBad {
		t.Fatalf("outcome = %v, want persistently_bad: a +90s skew must translate ActiveAt back to exactly `from` (F2)", outcome)
	}
	if badFor != to.Sub(from) {
		t.Fatalf("badFor = %v, want the full window %v", badFor, to.Sub(from))
	}
}

// --- F4: InstanceLabels must survive a timeline first created by a bare marker ---

func TestClassifyRule_LabelsSurviveWhenTimelineStartsFromAClearedMarker(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(10 * time.Minute)
	def := Definition{UID: "r1", Title: "R1"}
	key := instanceKey(lbl("a"))
	newOnset := from.Add(5 * time.Minute)

	polls := []Poll{
		// The very first mention of this key is a bare Cleared marker (its
		// prior bad episode, if any, started before the window) — no labels
		// travel with a Cleared/Vanished event.
		clearedPoll("r1", from.Add(1*time.Minute), key),
		abnormalPoll("r1", newOnset, StateFiring, lbl("a"), newOnset),
		quietPoll("r1", to),
	}
	_, _, viols := classifyRule(def, polls, from, to, defaultBad, PreexistingFailUnlessRecovered)
	if len(viols) != 1 {
		t.Fatalf("viols = %+v, want exactly one newly_bad violation", viols)
	}
	if viols[0].InstanceLabels == nil || viols[0].InstanceLabels["instance"] != "a" {
		t.Fatalf("InstanceLabels = %+v, want {instance: a}: labels must backfill even though the "+
			"timeline was first created by a label-less Cleared marker (F4)", viols[0].InstanceLabels)
	}
}

// TestClassifyRule_ViolationFieldsArePrecise pins FirstSeen/ClearedAt exactly,
// not just that a violation exists (F7).
func TestClassifyRule_ViolationFieldsArePrecise(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(10 * time.Minute)
	onset := from.Add(2 * time.Minute)
	clearAt := from.Add(3 * time.Minute)
	def := Definition{UID: "r1", Title: "R1"}

	polls := []Poll{
		quietPoll("r1", from),
		abnormalPoll("r1", onset, StateFiring, lbl("a"), onset),
		clearedPoll("r1", clearAt, instanceKey(lbl("a"))),
		quietPoll("r1", to),
	}
	_, _, viols := classifyRule(def, polls, from, to, defaultBad, PreexistingFailUnlessRecovered)
	if len(viols) != 1 {
		t.Fatalf("viols = %+v, want exactly one violation", viols)
	}
	v := viols[0]
	if !v.FirstSeen.Equal(onset) {
		t.Fatalf("FirstSeen = %v, want %v", v.FirstSeen, onset)
	}
	if !v.ClearedAt.Equal(clearAt) {
		t.Fatalf("ClearedAt = %v, want %v", v.ClearedAt, clearAt)
	}
	if v.InstanceLabels["instance"] != "a" {
		t.Fatalf("InstanceLabels = %+v, want {instance: a}", v.InstanceLabels)
	}
}

// TestClassifyRule_ClearedEventPastWindowEndClampsToWindowEnd pins the
// episode.end clamp: inWindowPolls admits a poll up to its own skew bound
// past windowEnd (§16's widened membership test), so a genuine Cleared event
// on such a poll must not leave the episode extending beyond windowEnd.
func TestClassifyRule_ClearedEventPastWindowEndClampsToWindowEnd(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(10 * time.Minute)
	def := Definition{UID: "r1", Title: "R1"}
	key := instanceKey(lbl("a"))

	bound := 30 * time.Second
	clearedAt := to.Add(20 * time.Second) // past windowEnd, but within the skew bound

	polls := []Poll{
		abnormalPoll("r1", from, StateFiring, lbl("a"), from.Add(-time.Hour)),
		{
			RuleUID: "r1", GrafanaNow: clearedAt, Found: true, Health: "ok",
			SkewBoundMS: bound.Milliseconds(), Cleared: []string{key},
		},
	}
	outcome, badFor, _ := classifyRule(def, polls, from, to, defaultBad, PreexistingFailUnlessRecovered)
	if outcome != OutcomeRecovered {
		t.Fatalf("outcome = %v, want recovered", outcome)
	}
	if badFor != to.Sub(from) {
		t.Fatalf("badFor = %v, want the window %v exactly: the episode end must clamp to windowEnd, "+
			"not extend to the late Cleared event's raw time", badFor, to.Sub(from))
	}
}

// TestClassifyRule_OnsetJustPastWindowEndIsNewlyBadNotClean pins the fail-closed
// reading of the upper boundary: an instance whose runner-domain onset lands
// only slightly past windowEnd (to + transitionGrace) is reachable at all only
// because inWindowPolls widens the boundary outward by the skew bound, so the
// gate cannot PROVE it belongs to the next window. It is charged as newly_bad —
// with BadFor truncated to zero — rather than silently forgiven as clean.
func TestClassifyRule_OnsetJustPastWindowEndIsNewlyBadNotClean(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(10 * time.Minute)
	grace := time.Minute
	windowEnd := to.Add(grace)
	def := Definition{UID: "r1", Title: "R1"}

	// A poll admitted only by its own skew bound: its GrafanaNow sits 20s past
	// windowEnd, inside the 30s tolerance. It carries an instance whose onset
	// is 10s past windowEnd — still "after the grace", but only by less than
	// the measurement's own uncertainty.
	bound := 30 * time.Second
	poll := Poll{
		RuleUID: "r1", GrafanaNow: windowEnd.Add(20 * time.Second), Found: true, Health: "ok",
		LastEvaluation: windowEnd.Add(20 * time.Second), SkewBoundMS: bound.Milliseconds(),
		Abnormal: []Instance{{Labels: lbl("a"), State: StateFiring, ActiveAt: windowEnd.Add(10 * time.Second)}},
	}

	outcome, badFor, viols := classifyRule(def, []Poll{poll}, from, windowEnd, defaultBad, PreexistingFailUnlessRecovered)
	if outcome != OutcomeNewlyBad {
		t.Fatalf("outcome = %v, want newly_bad: an onset past windowEnd seen only via the skew bound must fail closed", outcome)
	}
	if badFor != 0 {
		t.Fatalf("badFor = %v, want 0: the zero-length episode must truncate to the window end", badFor)
	}
	if len(viols) != 1 {
		t.Fatalf("viols = %+v, want exactly one newly_bad violation", viols)

	}
}

// §22.2: "a clear after `to` gives persistently_bad." classifyRule filters
// its input to [from, windowEnd] itself (inWindowPolls), so a Cleared event
// GENUINELY past windowEnd — well beyond any skew bound, unlike the clamp
// case above — never reaches the timeline at all: the instance is still bad
// at windowEnd as far as this window is concerned.
func TestClassifyRule_ClearAfterWindowEndIsPersistentlyBad(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(10 * time.Minute)
	def := Definition{UID: "r1", Title: "R1"}
	key := instanceKey(lbl("a"))

	polls := []Poll{
		abnormalPoll("r1", from, StateFiring, lbl("a"), from.Add(-time.Hour)),
		clearedPoll("r1", to.Add(time.Hour), key), // far past `to`, not a boundary case
	}
	outcome, badFor, viols := classifyRule(def, polls, from, to, defaultBad, PreexistingFailUnlessRecovered)
	if outcome != OutcomePersistentlyBad {
		t.Fatalf("outcome = %v, want persistently_bad: a clear outside the window must not read as a recovery", outcome)
	}
	if badFor != to.Sub(from) {
		t.Fatalf("badFor = %v, want the full window %v", badFor, to.Sub(from))
	}
	if len(viols) != 1 || viols[0].Outcome != OutcomePersistentlyBad {
		t.Fatalf("viols = %+v, want one persistently_bad violation", viols)
	}
}

// TestClassifyRule_CloseBeforeOpenClampsToZeroNotNegative pins the
// end-before-start clamp: two polls with different measured skews can
// translate so that a closing poll's runner-domain time lands before the
// opening poll's, which — unclamped — would feed mergeDurations a negative
// span.
func TestClassifyRule_CloseBeforeOpenClampsToZeroNotNegative(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(10 * time.Minute)
	onsetPoll := from.Add(5 * time.Minute)
	closePoll := from.Add(6 * time.Minute)
	closeSkew := 2 * time.Minute // translates closePoll back to from+4min, before onsetPoll's from+5min
	def := Definition{UID: "r1", Title: "R1"}
	key := instanceKey(lbl("a"))

	polls := []Poll{
		quietPoll("r1", from),
		abnormalPoll("r1", onsetPoll, StateFiring, lbl("a"), onsetPoll), // skew 0
		{
			RuleUID: "r1", GrafanaNow: closePoll, Found: true, Health: "ok",
			SkewMS: closeSkew.Milliseconds(), Cleared: []string{key},
		},
		quietPoll("r1", to),
	}
	outcome, badFor, viols := classifyRule(def, polls, from, to, defaultBad, PreexistingFailUnlessRecovered)
	if outcome != OutcomeNewlyBad {
		t.Fatalf("outcome = %v, want newly_bad", outcome)
	}
	if badFor < 0 {
		t.Fatalf("badFor = %v, want a non-negative duration even though the closing poll's translated "+
			"time landed before the opening poll's", badFor)
	}
	if badFor != 0 {
		t.Fatalf("badFor = %v, want 0: the clamp collapses the inverted span to a zero-length episode", badFor)
	}
	if len(viols) != 1 {
		t.Fatalf("viols = %+v, want one violation", viols)
	}
}

// --- mergeDurations ---

func TestMergeDurations_OverlappingEpisodesCountOnce(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	eps := []episode{
		{start: from, end: from.Add(5 * time.Minute)},
		{start: from.Add(2 * time.Minute), end: from.Add(8 * time.Minute)},   // overlaps the first
		{start: from.Add(20 * time.Minute), end: from.Add(21 * time.Minute)}, // disjoint
	}
	got := mergeDurations(eps)
	want := 8*time.Minute + 1*time.Minute // [0,8) merged = 8m, plus the disjoint 1m
	if got != want {
		t.Fatalf("mergeDurations = %v, want %v: two simultaneously-bad instances must not double-count their overlap", got, want)
	}
}

func TestMergeDurations_Empty(t *testing.T) {
	if got := mergeDurations(nil); got != 0 {
		t.Fatalf("mergeDurations(nil) = %v, want 0", got)
	}
}
