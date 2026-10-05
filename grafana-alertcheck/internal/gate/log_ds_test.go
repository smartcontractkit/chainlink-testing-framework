package gate

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func dsRule(uid, name string, insts ...Instance) StateRule {
	return StateRule{
		Key: ruleKey("vm", "G", name, ""), DatasourceUID: "vm",
		Title: name, Group: "G", Type: "alerting", Health: "ok",
		LastEvaluation: testNow, Instances: insts,
	}
}

// A datasource instance leaving the active set IS a resolution: vmalert only
// returns active instances, so a departure is Cleared, never Vanished.
func TestReduce_DatasourceDepartureIsCleared(t *testing.T) {
	firing := Instance{Labels: map[string]string{"x": "y"}, State: StateFiring, ActiveAt: testNow}
	r := NewReducer()
	key := ruleKey("vm", "G", "A", "")
	r.Reduce(key, observation(testNow, dsRule("", "A", firing)))
	p := r.Reduce(key, observation(testNow.Add(time.Minute), dsRule("", "A")))
	require.Equal(t, []string{instanceKey(firing.Labels)}, p.Cleared)
	require.Empty(t, p.Vanished)
}

// The same shape for a Grafana rule stays Vanished: a disappearing series must
// never read as a recovery.
func TestReduce_GrafanaDepartureIsVanished(t *testing.T) {
	firing := Instance{Labels: map[string]string{"x": "y"}, State: StateFiring, ActiveAt: testNow}
	r := NewReducer()
	r.Reduce("r1", observation(testNow, StateRule{UID: "r1", Title: "A", Health: "ok", LastEvaluation: testNow, Instances: []Instance{firing}}))
	p := r.Reduce("r1", observation(testNow.Add(time.Minute), StateRule{UID: "r1", Title: "A", Health: "ok", LastEvaluation: testNow.Add(time.Minute)}))
	require.Empty(t, p.Cleared)
	require.Equal(t, []string{instanceKey(firing.Labels)}, p.Vanished)
}

// A datasource poll records the key but no uid.
func TestReduce_DatasourcePollCarriesKeyNotUID(t *testing.T) {
	r := NewReducer()
	key := ruleKey("vm", "G", "A", "")
	p := r.Reduce(key, observation(testNow, dsRule("", "A")))
	require.Equal(t, key, p.RuleKey)
	require.Empty(t, p.RuleUID)
}

// A v1 log written before rule_key existed still reads: pollKey falls back to
// rule_uid.
func TestReadLog_OldPollWithoutRuleKey(t *testing.T) {
	p := Poll{RuleUID: "rule1", GrafanaNow: testNow, Found: true}
	require.Equal(t, "rule1", pollKey(p))
}

// The recorder's child rebuilds a datasource poll ref from the header alone.
func TestChildSchedule_DatasourceRef(t *testing.T) {
	def := dsDef("A")
	rt := map[string]RuleTimings{defKey(def): newRuleTimings(30*time.Second, 60)}
	h := Header{StartedAt: testNow, Rules: loggedRules([]Definition{def}, rt)}
	require.Equal(t, "datasource", h.Rules[0].SourceKind)
	require.Equal(t, "vm", h.Rules[0].DatasourceUID)

	refs, cadence, err := childSchedule(h)
	require.NoError(t, err)
	ref := refs[defKey(def)]
	require.Equal(t, KindDatasourceManaged, ref.Kind)
	require.Equal(t, "vm", ref.DatasourceUID)
	require.Equal(t, "G", ref.Group)
	require.Equal(t, "A", ref.Name)
	require.Equal(t, 30*time.Second, cadence[defKey(def)])

	if _, _, err := DeriveTimingsFromLog(h, []Definition{def}); err != nil {
		t.Fatalf("DeriveTimingsFromLog with a datasource rule: %v", err)
	}
}
