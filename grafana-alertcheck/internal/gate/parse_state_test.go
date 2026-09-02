package gate

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func readFixture(t *testing.T, name string) []byte {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("testdata", name))
	require.NoErrorf(t, err, "reading fixture %s", name)
	return b
}

func TestParseState_HappyPaths(t *testing.T) {
	cases := []struct {
		fixture       string
		wantRules     int
		wantInstances int
		checkFirst    func(t *testing.T, r StateRule)
	}{
		{
			fixture:       "state_one_instance.json",
			wantRules:     1,
			wantInstances: 1,
			checkFirst: func(t *testing.T, r StateRule) {
				require.Equal(t, "rule0000001", r.UID)
				require.Equal(t, "ExampleTeam", r.Folder)
				require.Equal(t, "Example Service - Prod", r.Group)
				require.Equal(t, "ok", r.Health)
				require.Equal(t, "inactive", r.State)
				require.Equal(t, float64(60), r.Interval.Seconds())
				require.False(t, r.IsPaused)
				require.Len(t, r.Instances, 1)
				require.Equal(t, StateNormal, r.Instances[0].State)

				inst := r.Instances[0]
				wantLabels := map[string]string{
					"alertname":      "Example High CPU Usage",
					"env":            "prod",
					"grafana_folder": "ExampleTeam",
					"app_instance":   "example-app-instance",
					"service":        "example-svc",
					"severity":       "critical",
					"team":           "example-team",
				}
				require.Equal(t, wantLabels, inst.Labels)
				wantActiveAt, err := time.Parse(time.RFC3339, "2026-08-31T08:02:50Z")
				require.NoError(t, err, "test setup")
				require.True(t, inst.ActiveAt.Equal(wantActiveAt))
				require.Empty(t, inst.Value)
			},
		},
		{
			fixture:       "state_paused.json",
			wantRules:     1,
			wantInstances: 0,
			checkFirst: func(t *testing.T, r StateRule) {
				require.True(t, r.IsPaused)
				require.True(t, r.LastEvaluation.IsZero())
				require.Equal(t, "ok", r.Health)
				require.Equal(t, "inactive", r.State)
			},
		},
		{
			fixture:       "state_health_error.json",
			wantRules:     1,
			wantInstances: 1,
			checkFirst: func(t *testing.T, r StateRule) {
				require.Equal(t, "error", r.Health)
				require.NotEmpty(t, r.LastError)
				require.Len(t, r.Instances, 1)
				require.Equal(t, StateError, r.Instances[0].State)
			},
		},
		{
			fixture:       "state_health_nodata.json",
			wantRules:     1,
			wantInstances: 1,
			checkFirst: func(t *testing.T, r StateRule) {
				require.Equal(t, "nodata", r.Health)
				require.Len(t, r.Instances, 1)
				require.Equal(t, StateNodata, r.Instances[0].State)
			},
		},
		{
			fixture:       "state_reason_composite.json",
			wantRules:     1,
			wantInstances: 3,
			checkFirst: func(t *testing.T, r StateRule) {
				byReason := map[string]Instance{}
				for _, inst := range r.Instances {
					byReason[inst.Reason] = inst
				}
				errInst, ok := byReason["Error"]
				require.True(t, ok, `want an instance with Reason="Error"`)
				require.Equal(t, StateNormal, errInst.State)
				nodataInst, ok := byReason["NoData"]
				require.True(t, ok, `want an instance with Reason="NoData"`)
				require.Equal(t, StateNormal, nodataInst.State)
				plain, ok := byReason[""]
				require.True(t, ok, `want a plain Reason="" instance`)
				require.Equal(t, StateNormal, plain.State)
			},
		},
		{
			fixture:       "state_missing_optional.json",
			wantRules:     1,
			wantInstances: 0,
			checkFirst: func(t *testing.T, r StateRule) {
				require.Nil(t, r.Instances)
				require.Nil(t, r.Totals)
			},
		},
		{
			fixture:       "state_only_active_instances.json",
			wantRules:     1,
			wantInstances: 1,
			checkFirst: func(t *testing.T, r StateRule) {
				require.Len(t, r.Instances, 1)
				require.Equal(t, StateFiring, r.Instances[0].State)
				require.NotZero(t, r.Totals["normal"],
					"the totals/instances mismatch this fixture exists to capture")
			},
		},
	}

	for _, c := range cases {
		t.Run(c.fixture, func(t *testing.T) {
			rules, err := ParseState(readFixture(t, c.fixture))
			require.NoErrorf(t, err, "ParseState(%s)", c.fixture)
			require.Lenf(t, rules, c.wantRules, "ParseState(%s)", c.fixture)
			require.Lenf(t, rules[0].Instances, c.wantInstances, "ParseState(%s)", c.fixture)
			if c.checkFirst != nil {
				c.checkFirst(t, rules[0])
			}
		})
	}
}

// The strict-parsing regression suite: it doesn't just check err != nil (a
// stray comma in a fixture would keep that green forever while the actual
// check regressed) — it asserts the error names the specific offending field
// or value, so a check going missing fails loudly here instead of surviving
// unnoticed.
func TestParseState_MustError(t *testing.T) {
	cases := []struct {
		fixture      string
		wantContains []string
	}{
		{"state_missing_health.json", []string{`"health"`}},
		{"state_missing_lasteval.json", []string{`"lastEvaluation"`}},
		{"state_missing_state.json", []string{`"state"`}},
		{"state_missing_interval.json", []string{`"interval"`}},
		{"state_missing_file.json", []string{`"file"`}},
		{"state_missing_name.json", []string{`"name"`}},
		{"state_zerotime_unpaused.json", []string{"zero time", "isPaused"}},
		{"state_unknown_state.json", []string{`"Weird (NoData)"`}},
	}
	for _, c := range cases {
		t.Run(c.fixture, func(t *testing.T) {
			_, err := ParseState(readFixture(t, c.fixture))
			require.Errorf(t, err, "ParseState(%s): expected an error, got none", c.fixture)
			for _, want := range c.wantContains {
				require.Containsf(t, err.Error(), want, "ParseState(%s): error", c.fixture)
			}
		})
	}
}

func TestParseNormalizeInstanceState(t *testing.T) {
	cases := []struct {
		in         string
		wantState  State
		wantReason string
		wantErr    bool
	}{
		{"Normal", StateNormal, "", false},
		{"Alerting", StateFiring, "", false},
		{"Pending", StatePending, "", false},
		{"NoData", StateNodata, "", false},
		{"Error", StateError, "", false},
		{"Normal (NoData)", StateNormal, "NoData", false},
		{"Normal (Error)", StateNormal, "Error", false},
		{"Normal (MissingSeries)", StateNormal, "MissingSeries", false},
		{"Weird (NoData)", "", "", true},
		{"Weird", "", "", true},
		{"", "", "", true},
	}
	for _, c := range cases {
		state, reason, err := normalizeInstanceState(c.in)
		if c.wantErr {
			require.Errorf(t, err, "normalizeInstanceState(%q)", c.in)
			continue
		}
		require.NoErrorf(t, err, "normalizeInstanceState(%q)", c.in)
		require.Equalf(t, c.wantState, state, "normalizeInstanceState(%q)", c.in)
		require.Equalf(t, c.wantReason, reason, "normalizeInstanceState(%q)", c.in)
	}
}

func TestInstanceKey(t *testing.T) {
	a := instanceKey(map[string]string{"b": "2", "a": "1"})
	b := instanceKey(map[string]string{"a": "1", "b": "2"})
	require.Equal(t, b, a, "instanceKey order-independence")
	require.Equal(t, "a=1\nb=2\n", a)

	diff := instanceKey(map[string]string{"a": "1", "b": "3"})
	require.NotEqual(t, a, diff, "instanceKey should differ when a label value differs")

	require.Empty(t, instanceKey(nil))
}

// minimalStateBody is the smallest legal state response: one group, one
// rule, no optional keys at all, plus whatever extra is spliced in verbatim
// before the rule's closing brace — for isolating one optional key at a time
// rather than relying on a fixture that removes several together.
func minimalStateBody(extraRuleJSON string) []byte {
	return fmt.Appendf(nil,
		`{"status":"success","data":{"groups":[{"file":"F","name":"G","interval":60,`+
			`"rules":[{"uid":"r1","name":"R1","state":"inactive","health":"ok","isPaused":false,`+
			`"lastEvaluation":"2026-01-01T00:00:00Z"%s}]}]}}`, extraRuleJSON)
}

// keepFiringFor is optional alongside alerts/totals/labels, but
// state_missing_optional.json removes it together with everything else — never
// in isolation, so a regression that made it required specifically would not
// be caught by that fixture alone.
func TestParseState_KeepFiringForIsOptional(t *testing.T) {
	tests := []struct {
		name  string
		extra string
	}{
		{"present", `,"keepFiringFor":300`},
		{"absent", ""},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rules, err := ParseState(minimalStateBody(tc.extra))
			require.NoError(t, err)
			require.Len(t, rules, 1)
		})
	}
}

// labels is optional at the INSTANCE level (opt(m, "labels", ...) in
// parseInstance), distinct from the rule-level labels state_missing_optional.json
// already covers — an instance can exist with no labels of its own.
func TestParseState_InstanceWithoutLabelsParses(t *testing.T) {
	body := minimalStateBody(`,"alerts":[{"state":"Normal","activeAt":"2026-01-01T00:00:00Z"}]`)
	rules, err := ParseState(body)
	require.NoError(t, err)
	require.Len(t, rules, 1)
	require.Len(t, rules[0].Instances, 1)
	require.Empty(t, rules[0].Instances[0].Labels)
}

// synthesizeHighCardinalityState builds a state response with a single rule
// holding `alerting` Alerting instances and `normal` Normal instances, by
// cloning the one real instance in state_one_instance.json. It is never
// committed — the 2446-instance rule this stands in for is ~600 KB and exists
// only to prove the parser and the reducer don't choke on real fleet
// cardinality.
func synthesizeHighCardinalityState(t *testing.T, alerting, normal int) []byte {
	t.Helper()
	base := readFixture(t, "state_one_instance.json")

	var top map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(base, &top))
	var data map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(top["data"], &data))
	var groups []map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(data["groups"], &groups))
	var rules []map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(groups[0]["rules"], &rules))
	var alerts []map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(rules[0]["alerts"], &alerts))
	template := alerts[0]

	newAlerts := make([]map[string]json.RawMessage, 0, alerting+normal)
	for i := range alerting {
		inst := cloneRawMap(template)
		inst["state"] = mustRaw(t, "Alerting")
		inst["labels"] = mustRaw(t, map[string]string{"instance": fmt.Sprintf("alerting-%d", i)})
		newAlerts = append(newAlerts, inst)
	}
	for i := range normal {
		inst := cloneRawMap(template)
		inst["state"] = mustRaw(t, "Normal")
		inst["labels"] = mustRaw(t, map[string]string{"instance": fmt.Sprintf("normal-%d", i)})
		newAlerts = append(newAlerts, inst)
	}

	rules[0]["alerts"] = mustRaw(t, newAlerts)
	rules[0]["totals"] = mustRaw(t, map[string]int{"alerting": alerting, "normal": normal})
	rules[0]["totalsFiltered"] = rules[0]["totals"]
	groups[0]["rules"] = mustRaw(t, rules)
	data["groups"] = mustRaw(t, groups)
	top["data"] = mustRaw(t, data)

	out, err := json.Marshal(top)
	require.NoError(t, err)
	return out
}

func cloneRawMap(m map[string]json.RawMessage) map[string]json.RawMessage {
	out := make(map[string]json.RawMessage, len(m))
	for k, v := range m {
		cp := make(json.RawMessage, len(v))
		copy(cp, v)
		out[k] = cp
	}
	return out
}

func mustRaw(t *testing.T, v any) json.RawMessage {
	t.Helper()
	b, err := json.Marshal(v)
	require.NoError(t, err)
	return json.RawMessage(b)
}

func TestParseState_HighCardinality(t *testing.T) {
	body := synthesizeHighCardinalityState(t, 445, 2004)
	require.Contains(t, string(body), "alerting-0")

	rules, err := ParseState(body)
	require.NoError(t, err)
	require.Len(t, rules, 1)
	r := rules[0]
	require.Len(t, r.Instances, 445+2004)

	var firing, normal int
	for _, inst := range r.Instances {
		switch inst.State {
		case StateFiring:
			firing++
		case StateNormal:
			normal++
		default:
			require.Fail(t, fmt.Sprintf("unexpected instance state %q", inst.State))
		}
	}
	require.Equal(t, 445, firing)
	require.Equal(t, 2004, normal)

	// Each synthesized instance carries a distinct "instance" label; confirm
	// Labels actually made it through parsing (not just State) by checking
	// instanceKey produces one unique key per instance, with no collisions.
	seen := make(map[string]bool, len(r.Instances))
	for _, inst := range r.Instances {
		require.NotNil(t, inst.Labels)
		k := instanceKey(inst.Labels)
		require.Falsef(t, seen[k], "duplicate instance key %q", k)
		seen[k] = true
	}
	require.Len(t, seen, 445+2004)
}
