package gate

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRuleKey_GrafanaKeepsUID(t *testing.T) {
	require.Equal(t, "rule1", ruleKey("", "", "title", "rule1"))
}

// A ds key must not collide with a Grafana uid and must not let group/name
// separators collide: a name containing ":" or "/" is still a distinct tuple.
func TestRuleKey_DatasourceInjectivity(t *testing.T) {
	keys := []string{
		ruleKey("dsA", "g", "n", ""),
		ruleKey("dsA", "g/n", "", ""),
		ruleKey("dsA", "g", "/n", ""),
		ruleKey("dsB", "g", "n", ""),
		ruleKey("", "g", "n", ""),
	}
	seen := map[string]bool{}
	for _, k := range keys {
		require.True(t, len(k) > len(dsKeyPrefix) && k[:len(dsKeyPrefix)] == dsKeyPrefix)
		require.False(t, seen[k], "key %q collided", k)
		seen[k] = true
	}
	require.Equal(t, "u1", ruleKey("dsA", "g", "n", "u1"), "a uid wins over the ds tuple")
}

func TestDefKey_FallsBackToUID(t *testing.T) {
	require.Equal(t, "u1", defKey(Definition{UID: "u1"}))
	require.Equal(t, "k1", defKey(Definition{Key: "k1", UID: "u1"}))
}

func TestPollKey_FallsBackToUID(t *testing.T) {
	require.Equal(t, "u1", pollKey(Poll{RuleUID: "u1"}))
	require.Equal(t, "k1", pollKey(Poll{RuleKey: "k1", RuleUID: "u1"}))
}
