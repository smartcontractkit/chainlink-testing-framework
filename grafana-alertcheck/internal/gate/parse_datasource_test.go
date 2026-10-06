package gate

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestParseDatasourceRules_Fixture(t *testing.T) {
	rules, err := ParseDatasourceRules(readFixture(t, "ds_rules.json"), "ds-uid")
	require.NoError(t, err)
	require.Len(t, rules, 2, "recording rules are parsed but filtered later")

	alert := rules[0]
	require.Equal(t, "ExampleTargetDown", alert.Title)
	require.Equal(t, "ExampleMetrics", alert.Group)
	require.Equal(t, "/etc/vm/rules/example.yml", alert.File)
	require.Equal(t, "ds-uid", alert.DatasourceUID)
	require.Equal(t, "alerting", alert.Type)
	require.Equal(t, "up == 0", alert.Query)
	require.Equal(t, 5*time.Minute, alert.For)
	require.Equal(t, "firing", alert.State)
	require.Equal(t, "ok", alert.Health)
	require.Empty(t, alert.UID, "a datasource rule has no uid")
	require.Equal(t, ruleKey("ds-uid", "ExampleMetrics", "ExampleTargetDown", "/etc/vm/rules/example.yml", ""), alert.Key)
	require.Len(t, alert.Instances, 1)
	require.Equal(t, StateFiring, alert.Instances[0].State)
	require.Nil(t, alert.Totals)
}

func TestDefinitionsFromDatasource_FiltersRecording(t *testing.T) {
	rules, err := ParseDatasourceRules(readFixture(t, "ds_rules.json"), "ds-uid")
	require.NoError(t, err)

	defs := DefinitionsFromDatasource(rules, "ds-uid", "ExampleMetrics")
	require.Len(t, defs, 1, "only the alerting rule becomes a Definition")
	require.Equal(t, KindDatasourceManaged, defs[0].Kind)
	require.Equal(t, "ExampleMetrics", defs[0].DatasourceName)
	require.False(t, defs[0].PauseObservable)
	require.Equal(t, 60, defs[0].IntervalSeconds)
}

func TestParseDatasourceRules_HealthErrNormalizes(t *testing.T) {
	body := []byte(`{"status":"success","data":{"groups":[{"name":"g","file":"f","interval":60,"rules":[
		{"name":"A","type":"alerting","health":"err","lastEvaluation":"2026-08-01T00:00:00Z","state":"firing"}]}]}}`)
	rules, err := ParseDatasourceRules(body, "d")
	require.NoError(t, err)
	require.Equal(t, "error", rules[0].Health)
}

func TestParseDatasourceRules_LastError(t *testing.T) {
	body := []byte(`{"status":"success","data":{"groups":[{"name":"g","rules":[
		{"name":"A","type":"alerting","health":"err","lastError":"query failed: bad","state":"firing"}]}]}}`)
	rules, err := ParseDatasourceRules(body, "d")
	require.NoError(t, err)
	require.Equal(t, "error", rules[0].Health)
	require.Equal(t, "query failed: bad", rules[0].LastError)
}

func TestParseDatasourceRules_ZeroLastEvaluationAllowed(t *testing.T) {
	body := []byte(`{"status":"success","data":{"groups":[{"name":"g","rules":[
		{"name":"A","type":"alerting","health":"ok","state":"pending"}]}]}}`)
	rules, err := ParseDatasourceRules(body, "d")
	require.NoError(t, err)
	require.True(t, rules[0].LastEvaluation.IsZero())
}

func TestParseDatasourceRules_UnknownInstanceStateIsError(t *testing.T) {
	body := []byte(`{"status":"success","data":{"groups":[{"name":"g","rules":[
		{"name":"A","type":"alerting","health":"ok","state":"firing","alerts":[
			{"labels":{},"state":"inactive","activeAt":"2026-08-01T00:00:00Z"}]}]}]}}`)
	_, err := ParseDatasourceRules(body, "d")
	require.Error(t, err)
	require.Contains(t, err.Error(), "unrecognized datasource instance state")
}

func TestParseDatasourceRules_PendingInstance(t *testing.T) {
	body := []byte(`{"status":"success","data":{"groups":[{"name":"g","rules":[
		{"name":"A","type":"alerting","health":"ok","state":"pending","alerts":[
			{"labels":{"x":"y"},"state":"pending","activeAt":"2026-08-01T00:00:00Z","value":"1"}]}]}]}}`)
	rules, err := ParseDatasourceRules(body, "d")
	require.NoError(t, err)
	require.Equal(t, StatePending, rules[0].Instances[0].State)
}
