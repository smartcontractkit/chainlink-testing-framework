package gate

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func terminalHeader(from time.Time) Header {
	return Header{
		SchemaVersion: LogSchemaVersion,
		StartedAt:     from.Add(-time.Minute),
		Rules: []LoggedRule{{
			UID: checkUID, Title: checkTitle, IntervalSeconds: 60,
			NoDataState: "OK", ExecErrState: "OK",
			PollEverySeconds: checkPollEvery.Seconds(),
		}},
	}
}

// badFromPolls is healthy until onset, then carries the same firing instance.
func badFromPolls(uid string, from, at, onset time.Time) []Poll {
	var out []Poll
	for ts := from; !ts.After(at); ts = ts.Add(checkPollEvery) {
		if ts.Before(onset) {
			out = append(out, quietPoll(uid, ts))
		} else {
			out = append(out, abnormalPoll(uid, ts, StateFiring, lbl("a"), onset))
		}
	}
	return out
}

func TestTerminalVerdict(t *testing.T) {
	from := testNow
	at := from.Add(5 * time.Minute)
	def := checkDef()
	rt := map[string]RuleTimings{checkUID: newRuleTimings(checkPollEvery, 60)}
	pol := Policy{From: from, To: at}

	tests := []struct {
		name       string
		polls      func() []Poll
		want       bool
		wantKind   TerminationKind
		wantReason UnobservableReason
		wantOut    Outcome
	}{
		{
			name:  "clean window is not terminal",
			polls: func() []Poll { return denseHealthyPolls(checkUID, from, at, checkPollEvery) },
			want:  false,
		},
		{
			name:     "post-from onset is a terminal violation",
			polls:    func() []Poll { return badFromPolls(checkUID, from, at, from.Add(time.Minute)) },
			want:     true,
			wantKind: TerminationViolation,
			wantOut:  OutcomeNewlyBad,
		},
		{
			// Bad before `from` can still become `recovered` (a pass).
			name: "preexisting bad is not terminal",
			polls: func() []Poll {
				var out []Poll
				for ts := from; !ts.After(at); ts = ts.Add(checkPollEvery) {
					out = append(out, abnormalPoll(checkUID, ts, StateFiring, lbl("a"), from.Add(-10*time.Minute)))
				}
				return out
			},
			want: false,
		},
		{
			name: "heartbeat gap is terminal unobservable",
			polls: func() []Poll {
				return append(denseHealthyPolls(checkUID, from, from.Add(time.Minute), checkPollEvery),
					quietPoll(checkUID, at))
			},
			want:       true,
			wantKind:   TerminationUnobservable,
			wantReason: ReasonHeartbeatGap,
		},
		{
			name: "sustained health=error is terminal unobservable",
			polls: func() []Poll {
				var out []Poll
				for ts := from; !ts.After(at); ts = ts.Add(checkPollEvery) {
					p := quietPoll(checkUID, ts)
					p.Health = "error"
					out = append(out, p)
				}
				return out
			},
			want:       true,
			wantKind:   TerminationUnobservable,
			wantReason: ReasonHealthError,
		},
		{
			name: "in-window pause is terminal unobservable",
			polls: func() []Poll {
				out := denseHealthyPolls(checkUID, from, at, checkPollEvery)
				out = append(out, Poll{RuleUID: checkUID, GrafanaNow: from.Add(time.Minute), Found: true, Health: "ok", IsPaused: true})
				return out
			},
			want:       true,
			wantKind:   TerminationUnobservable,
			wantReason: ReasonPausedInWindow,
		},
		{
			name: "absent rule is terminal unobservable",
			polls: func() []Poll {
				out := denseHealthyPolls(checkUID, from, at, checkPollEvery)
				out = append(out, Poll{RuleUID: checkUID, GrafanaNow: from.Add(time.Minute)})
				return out
			},
			want:       true,
			wantKind:   TerminationUnobservable,
			wantReason: ReasonRuleAbsent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			term, ok := terminalVerdict(terminalHeader(from), tt.polls(), []Definition{def}, rt, pol, from, at)
			require.Equal(t, tt.want, ok)
			if !tt.want {
				return
			}
			require.Equal(t, tt.wantKind, term.Kind)
			require.Equal(t, tt.wantReason, term.Reason)
			if tt.wantOut != "" {
				require.Equal(t, tt.wantOut, term.Outcome)
			}
			require.True(t, term.At.Equal(at))
		})
	}
}

// A sustained health=nodata run is terminal only when the policy escalates it,
// matching decide's end-of-run behavior.
func TestTerminalVerdictNodataIsTerminalOnlyWhenConfigured(t *testing.T) {
	from := testNow
	at := from.Add(5 * time.Minute)
	def := checkDef()
	rt := map[string]RuleTimings{checkUID: newRuleTimings(checkPollEvery, 60)}

	var polls []Poll
	for ts := from; !ts.After(at); ts = ts.Add(checkPollEvery) {
		p := quietPoll(checkUID, ts)
		p.Health = "nodata"
		polls = append(polls, p)
	}

	_, ok := terminalVerdict(terminalHeader(from), polls, []Definition{def}, rt, Policy{From: from, To: at}, from, at)
	require.False(t, ok, "health=nodata is only a note by default")

	term, ok := terminalVerdict(terminalHeader(from), polls, []Definition{def}, rt,
		Policy{From: from, To: at, NodataIsUnobservable: true}, from, at)
	require.True(t, ok)
	require.Equal(t, TerminationUnobservable, term.Kind)
	require.Equal(t, ReasonNodata, term.Reason)
}

// Inability beats violation mid-window, as at the end of a full run.
func TestTerminalVerdictUnobservableBeatsViolation(t *testing.T) {
	from := testNow
	at := from.Add(5 * time.Minute)

	violating := checkDef() // checkUID/checkTitle, fires after from
	gapped := Definition{UID: "rule-two", Title: "Rule Two", IntervalSeconds: 60}
	rt := map[string]RuleTimings{
		checkUID:   newRuleTimings(checkPollEvery, 60),
		"rule-two": newRuleTimings(checkPollEvery, 60),
	}

	polls := badFromPolls(checkUID, from, at, from.Add(time.Minute))
	// rule-two has a hole: nothing in the second half of the window.
	polls = append(polls, denseHealthyPolls("rule-two", from, from.Add(time.Minute), checkPollEvery)...)

	term, ok := terminalVerdict(terminalHeader(from), polls, []Definition{violating, gapped}, rt, Policy{From: from, To: at}, from, at)
	require.True(t, ok)
	require.Equal(t, TerminationUnobservable, term.Kind)
	require.Equal(t, "rule-two", term.RuleUID)
	require.Equal(t, ReasonHeartbeatGap, term.Reason)
}
