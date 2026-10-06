package gate

import (
	"fmt"
	"time"
)

// TerminationKind is the kind of terminal verdict. It reaches the JSON output.
type TerminationKind string

const (
	// TerminationViolation is a post-`from` bad onset (new_failure or unstable).
	TerminationViolation TerminationKind = "violation"
	// TerminationNotVerified is an inability that has already happened.
	TerminationNotVerified TerminationKind = "not_verified"
)

// Termination is why a fail-fast run stopped before the window closed. At is
// the runner-domain detection time.
type Termination struct {
	Kind    TerminationKind    `json:"kind"`
	Alert   string             `json:"alert,omitempty"`
	RuleKey string             `json:"rule_key,omitempty"`
	RuleUID string             `json:"rule_uid,omitempty"`
	Outcome Outcome            `json:"outcome,omitempty"`
	Reason  UnobservableReason `json:"reason,omitempty"`
	At      time.Time          `json:"at"`
}

// terminalVerdict reports the first monotone terminal condition among the polls
// observed so far, treating at as the provisional end of the window. PURE.
//
// Only two conditions qualify, because only they can never become a pass: an
// inability that already happened, and a post-`from` bad onset (new_failure or
// unstable). `recovered` forgives an observed bad state and is reserved for
// bad-at-`from`, so a preexisting condition is deliberately not terminal.
// not_verified beats violation, as it does at the end of a full run.
func terminalVerdict(h Header, polls []Poll, defs []Definition, rt map[string]RuleTimings,
	pol Policy, from, at time.Time) (Termination, bool) {

	badStates := badStateSet(pol.States)
	pausedAtStart := h.pausedAtStart()

	var violation *Termination
	for _, def := range defs {
		key := defKey(def)
		if pausedAtStart[key] {
			continue
		}
		// The synthetic sentinel at `at` satisfies check 1, leaving only the
		// checks decidable from the polls so far. The policy-specific nodata
		// escalation is applied here too, or a configured terminal inability
		// would never fail fast.
		cov := proveCoverage(h, polls, &at, rt[key], def, from, at, 0)
		if pol.NodataIsUnobservable {
			applyNodataPolicy(def, polls, &cov, rt[key], from, at)
		}
		if cov.Unobservable {
			return Termination{
				Kind:    TerminationNotVerified,
				Alert:   def.Title,
				RuleKey: key,
				RuleUID: def.UID,
				Outcome: OutcomeNotVerified,
				Reason:  cov.Reason,
				At:      at,
			}, true
		}

		outcome, _, _ := classifyRule(def, rt[key], polls, from, at, badStates, pol.Preexisting)
		if outcome == OutcomeNewFailure || outcome == OutcomeUnstable {
			if violation == nil {
				v := Termination{
					Kind:    TerminationViolation,
					Alert:   def.Title,
					RuleKey: key,
					RuleUID: def.UID,
					Outcome: outcome,
					At:      at,
				}
				violation = &v
			}
		}
	}
	if violation != nil {
		return *violation, true
	}
	return Termination{}, false
}

// earlyResult classifies the evidence observed when a terminal condition ended
// the run, and marks the Result as stopped early. It reuses decide unchanged by
// clamping the window to term.At and zeroing the grace (a synthetic sentinel
// at term.At satisfies check 1), then restores the requested window and real
// thresholds for reporting.
func earlyResult(h Header, polls []Poll, defs []Definition, rt map[string]RuleTimings,
	gt GlobalTimings, pol Policy, term Termination) (Result, error) {

	earlyPol := pol
	earlyPol.To = term.At
	earlyGT := gt
	earlyGT.transitionGrace = 0

	sentinel := term.At
	res, err := decide(h, polls, &sentinel, defs, rt, earlyGT, earlyPol)

	// An early Result can never be a pass; fail closed if classification lost
	// the terminal verdict.
	if err == nil && len(res.Violations) == 0 {
		return res, fmt.Errorf("gate: fail-fast stopped on %s (%s) but classification found no violation; refusing to report a pass",
			term.Kind, term.Alert)
	}

	res.From = pol.From
	res.To = pol.To
	res.Global = GlobalThresholds{
		TransitionGrace: gt.transitionGrace,
		GraceSource:     graceSourceOrNone(gt.graceSource),
		DrainTimeout:    gt.drainTimeout,
	}
	t := term
	res.TerminatedEarly = &t
	return res, err
}
