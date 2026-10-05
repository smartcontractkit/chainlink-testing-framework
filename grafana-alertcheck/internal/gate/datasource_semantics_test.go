package gate

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func dsDef(name string) Definition {
	return Definition{
		Key: ruleKey("vm", "G", name, ""), Title: name, Group: "G",
		Kind: KindDatasourceManaged, DatasourceUID: "vm", DatasourceName: "VM",
		IntervalSeconds: 60,
	}
}

// dsPoll is a datasource poll with the fields the pure layer reads.
func dsPoll(key string, at time.Time, health string) Poll {
	return Poll{RuleKey: key, GrafanaNow: at, Found: true, Health: health, LastEvaluation: at}
}

// A datasource instance that is bad at `from` and then leaves the active set is
// a recovery: Reduce turns the departure into Cleared, so classifyRule sees a
// real clear and the run passes.
func TestDecide_DatasourceDepartureIsRecovered(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(10 * time.Minute)
	def := dsDef("A")
	key := defKey(def)
	rt := map[string]RuleTimings{key: newRuleTimings(30*time.Second, 60)}
	pol := Policy{From: from, To: to}

	var polls []Poll
	for ts := from; !ts.After(to); ts = ts.Add(30 * time.Second) {
		switch {
		case ts.Equal(from):
			polls = append(polls, Poll{
				RuleKey: key, GrafanaNow: ts, Found: true, Health: "ok", LastEvaluation: ts,
				Abnormal: []Instance{{Labels: map[string]string{"x": "y"}, State: StateFiring, ActiveAt: from.Add(-time.Hour)}},
			})
		case ts.Equal(from.Add(5 * time.Minute)):
			p := dsPoll(key, ts, "ok")
			p.Cleared = []string{instanceKey(map[string]string{"x": "y"})}
			polls = append(polls, p)
		default:
			polls = append(polls, dsPoll(key, ts, "ok"))
		}
	}
	sentinel := to
	res, err := decide(Header{StartedAt: from.Add(-time.Hour)}, polls, &sentinel, []Definition{def}, rt, GlobalTimings{}, pol)
	require.NoError(t, err)
	require.Empty(t, res.Violations)
	require.Equal(t, OutcomeRecovered, res.Verdicts[0].Outcome)
	require.Contains(t, res.Verdicts[0].Note, "treated as a recovery")
}

// The same shape for a Grafana rule, but a VANISH rather than a clear, stays
// still_failing: a disappearing series must not read as a recovery.
func TestDecide_GrafanaVanishStaysFailing(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(10 * time.Minute)
	def := Definition{Key: "r1", UID: "r1", Title: "R1"}
	rt := map[string]RuleTimings{"r1": newRuleTimings(30*time.Second, 60)}
	pol := Policy{From: from, To: to}

	var polls []Poll
	for ts := from; !ts.After(to); ts = ts.Add(30 * time.Second) {
		if ts.Equal(from) {
			polls = append(polls, Poll{
				RuleUID: "r1", RuleKey: "r1", GrafanaNow: ts, Found: true, Health: "ok", LastEvaluation: ts,
				Abnormal: []Instance{{Labels: map[string]string{"x": "y"}, State: StateFiring, ActiveAt: from.Add(-time.Hour)}},
			})
			continue
		}
		p := quietPoll("r1", ts)
		if ts.Equal(from.Add(5 * time.Minute)) {
			p.Vanished = []string{instanceKey(map[string]string{"x": "y"})}
		}
		polls = append(polls, p)
	}
	sentinel := to
	res, err := decide(Header{StartedAt: from.Add(-time.Hour)}, polls, &sentinel, []Definition{def}, rt, GlobalTimings{}, pol)
	require.NoError(t, err)
	require.NotEmpty(t, res.Violations, "a vanish must stay a failure")
	require.Equal(t, OutcomeStillFailing, res.Verdicts[0].Outcome)
}

// A datasource health=err normalizes to "error" and, sustained past
// healthGrace, makes the rule unobservable through the ordinary check 4.
func TestDecide_DatasourceHealthErrIsUnobservable(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(10 * time.Minute)
	def := dsDef("A")
	key := defKey(def)
	rt := map[string]RuleTimings{key: newRuleTimings(30*time.Second, 60)}
	pol := Policy{From: from, To: to}

	var polls []Poll
	for ts := from; !ts.After(to); ts = ts.Add(30 * time.Second) {
		polls = append(polls, dsPoll(key, ts, "error"))
	}
	sentinel := to
	res, err := decide(Header{StartedAt: from.Add(-time.Hour)}, polls, &sentinel, []Definition{def}, rt, GlobalTimings{}, pol)
	require.Error(t, err)
	require.Equal(t, OutcomeNotVerified, res.Verdicts[0].Outcome)
	require.Equal(t, ReasonHealthError, res.Coverage[key].Reason)
}
