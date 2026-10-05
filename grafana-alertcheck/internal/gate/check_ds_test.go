package gate

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// End-to-end through check: a datasource-managed rule that is bad at `from` and
// leaves the active set mid-window classifies as recovered, exit 0.
func TestCheck_DatasourceFireAndResolveIsRecovered(t *testing.T) {
	clock := newVirtualClock(testNow)
	from := testNow
	to := testNow.Add(5 * time.Minute)
	def := dsDef("A")
	key := defKey(def)

	src := newCheckSource(nil)
	src.defs = nil
	src.ruleSources = []RuleSource{{UID: "vm", Name: "VM"}}
	src.dsDefs = map[string][]Definition{"vm": {def}}
	src.dsRespond = func(_ string, _ int) (Observation, error) {
		now := clock.Now()
		rule := StateRule{
			Key: key, DatasourceUID: "vm", Title: "A", Group: "G", Type: "alerting",
			Health: "ok", LastEvaluation: now,
		}
		if now.Before(from.Add(2 * time.Minute)) {
			rule.Instances = []Instance{{
				Labels: map[string]string{"x": "y"}, State: StateFiring, ActiveAt: from.Add(-time.Hour),
			}}
		}
		return Observation{Rules: []StateRule{rule}, GrafanaNow: now, Latency: 100 * time.Millisecond}, nil
	}

	cfg := Config{
		URL:    "https://grafana.example.com",
		Alerts: []string{"A"},
		From:   from,
		To:     to,
		Clock:  clock,
		Notes:  &strings.Builder{},
	}.withDefaults()

	res, err := check(context.Background(), cfg, src)
	require.NoError(t, err)
	require.Empty(t, res.Violations)
	require.Equal(t, OutcomeRecovered, res.Verdicts[0].Outcome)
	require.Contains(t, res.Verdicts[0].Note, "check 7 skipped")
}
