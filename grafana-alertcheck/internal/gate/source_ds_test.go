package gate

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// The []-suffixed keys are an upstream vmalert quirk: plain rule_name= is
// ignored, so the raw query must carry the bracketed forms.
func TestDatasourceQuery_BracketedKeys(t *testing.T) {
	require.Empty(t, datasourceQuery(nil, "", ""), "no filters means the bulk request")
	got := datasourceQuery([]string{"A"}, "G", "F")
	require.Equal(t, "?file%5B%5D=F&rule_group%5B%5D=G&rule_name%5B%5D=A", got)
}

func TestDiscoverRuleSources_StrictFilterAndProbe(t *testing.T) {
	datasources := `[
		{"uid":"vm","name":"VictoriaMetrics - Prod","type":"prometheus","jsonData":{"manageAlerts":true}},
		{"uid":"ash","name":"AlertStateHistoryBackend","type":"prometheus","jsonData":{"manageAlerts":false}},
		{"uid":"loki","name":"Loki","type":"loki","jsonData":{"manageAlerts":true}}
	]`
	var probed []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/datasources":
			_, _ = w.Write([]byte(datasources))
		case "/api/prometheus/vm/api/v1/rules":
			probed = append(probed, r.URL.RawQuery)
			_, _ = w.Write([]byte(`{"status":"success","data":{"groups":[]}}`))
		default:
			t.Errorf("unexpected path %q", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()

	src := NewHTTPSource(srv.URL, "", newFakeClock(time.Now()))
	got, err := src.DiscoverRuleSources(context.Background())
	require.NoError(t, err)
	require.Equal(t, []RuleSource{{UID: "vm", Name: "VictoriaMetrics - Prod"}}, got)
	require.Equal(t, []string{"rule_name%5B%5D=__probe__"}, probed)
}

func TestDiscoverRuleSources_ProbeFailureIsHardError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/datasources":
			_, _ = w.Write([]byte(`[{"uid":"vm","name":"VictoriaMetrics - Prod","type":"prometheus","jsonData":{"manageAlerts":true}}]`))
		default:
			w.WriteHeader(http.StatusForbidden)
		}
	}))
	defer srv.Close()

	src := NewHTTPSource(srv.URL, "", newFakeClock(time.Now()))
	_, err := src.DiscoverRuleSources(context.Background())
	require.Error(t, err)
	require.Contains(t, err.Error(), "VictoriaMetrics - Prod")
	require.Contains(t, err.Error(), "vm")
}

func TestDatasourceDefinitions_FilteredQuery(t *testing.T) {
	var gotQuery string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.RawQuery
		require.Equal(t, "/api/prometheus/vm/api/v1/rules", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(readFixture(t, "ds_rules.json"))
	}))
	defer srv.Close()

	src := NewHTTPSource(srv.URL, "", newFakeClock(time.Now()))
	defs, err := src.DatasourceDefinitions(context.Background(), RuleSource{UID: "vm", Name: "VM"}, []string{"ExampleTargetDown"})
	require.NoError(t, err)
	require.Equal(t, "rule_name%5B%5D=ExampleTargetDown", gotQuery)
	require.Len(t, defs, 1)
	require.Equal(t, KindDatasourceManaged, defs[0].Kind)
	require.Equal(t, "VM", defs[0].DatasourceName)
}

func TestDSFilterNames(t *testing.T) {
	filters, fetchAll := dsFilterNames([]string{"devex-cicd/prod/griddle-github: ContainersNotReady"})
	require.False(t, fetchAll)
	require.ElementsMatch(t, []string{
		"devex-cicd/prod/griddle-github: ContainersNotReady",
		"griddle-github: ContainersNotReady",
	}, filters)

	_, fetchAll = dsFilterNames([]string{"key:ds:[\"a\"]"})
	require.True(t, fetchAll)
}

// A datasource rule whose name contains "/" is fetched by its full name, so
// loadDefinitions must request the whole input, not just the last segment.
func TestLoadDefinitions_SlashyDatasourceName(t *testing.T) {
	name := "devex-cicd/prod/griddle-github: ContainersNotReady"
	f := newFakeSource()
	f.ruleSources = []RuleSource{{UID: "vm", Name: "VM"}}
	f.dsDefs = map[string][]Definition{"vm": {{
		Key: ruleKey("vm", "G", name, "f", ""), Title: name, Group: "G",
		Kind: KindDatasourceManaged, DatasourceUID: "vm", DatasourceName: "VM",
	}}}
	defs, err := loadDefinitions(context.Background(), f, []string{name}, false)
	require.NoError(t, err)
	require.Len(t, defs, 1)
	require.Equal(t, name, defs[0].Title)
}

func TestFakeSource_DatasourceScriptedByKey(t *testing.T) {
	f := newFakeSource()
	key := ruleKey("vm", "G", "A", "f", "")
	f.scriptKey(key, Observation{Rules: []StateRule{{Key: key, DatasourceUID: "vm", Title: "A"}}}, nil)
	obs, err := f.RuleState(context.Background(), RuleRef{Key: key, Kind: KindDatasourceManaged, DatasourceUID: "vm"})
	require.NoError(t, err)
	require.Len(t, obs.Rules, 1)
}

func TestLoadDefinitions_DiscoversAndDropsRulerDatasourceRules(t *testing.T) {
	f := newFakeSource()
	f.defs = []Definition{
		{Key: "g1", UID: "g1", Title: "Grafana Rule", Kind: KindGrafanaManaged},
		{Title: "RulerDsRule", Kind: KindDatasourceManaged}, // no datasource UID: dropped
	}
	f.ruleSources = []RuleSource{{UID: "vm", Name: "VM"}}
	f.dsDefs = map[string][]Definition{"vm": {dsDef("A"), dsDef("B")}}

	all, err := loadDefinitions(context.Background(), f, nil, true)
	require.NoError(t, err)
	require.Len(t, all, 3, "Grafana + two ds, ruler ds rule dropped")

	filtered, err := loadDefinitions(context.Background(), f, []string{"A"}, false)
	require.NoError(t, err)
	require.Len(t, filtered, 2)
	require.Equal(t, "A", filtered[1].Title)
}

func TestRuleState_DatasourceAssertsAllFilters(t *testing.T) {
	var gotPath, gotQuery string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotQuery = r.URL.Path, r.URL.RawQuery
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(readFixture(t, "ds_rules.json"))
	}))
	defer srv.Close()

	src := NewHTTPSource(srv.URL, "", newFakeClock(time.Now()))
	ref := RuleRef{
		Key:  ruleKey("vm", "ExampleMetrics", "ExampleTargetDown", "/etc/vm/rules/example.yml", ""),
		Kind: KindDatasourceManaged, DatasourceUID: "vm",
		Group: "ExampleMetrics", Name: "ExampleTargetDown", File: "/etc/vm/rules/example.yml",
	}
	obs, err := src.RuleState(context.Background(), ref)
	require.NoError(t, err)
	require.Equal(t, "/api/prometheus/vm/api/v1/rules", gotPath)
	require.Equal(t,
		"file%5B%5D=%2Fetc%2Fvm%2Frules%2Fexample.yml&rule_group%5B%5D=ExampleMetrics&rule_name%5B%5D=ExampleTargetDown",
		gotQuery)
	require.Len(t, obs.Rules, 2)
	require.Equal(t, ref.Key, obs.Rules[0].Key)
}
