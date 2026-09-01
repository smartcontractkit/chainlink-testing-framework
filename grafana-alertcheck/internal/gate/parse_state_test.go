package gate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func readFixture(t *testing.T, name string) []byte {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatalf("reading fixture %s: %v", name, err)
	}
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
				if r.UID != "rule0000001" {
					t.Errorf("UID = %q, want rule0000001", r.UID)
				}
				if r.Folder != "ExampleTeam" || r.Group != "Example Service - Prod" {
					t.Errorf("Folder/Group = %q/%q, want ExampleTeam/Example Service - Prod", r.Folder, r.Group)
				}
				if r.Health != "ok" || r.State != "inactive" {
					t.Errorf("Health/State = %q/%q, want ok/inactive", r.Health, r.State)
				}
				if r.Interval.Seconds() != 60 {
					t.Errorf("Interval = %v, want 60s", r.Interval)
				}
				if r.IsPaused {
					t.Errorf("IsPaused = true, want false")
				}
				if len(r.Instances) != 1 || r.Instances[0].State != StateNormal {
					t.Fatalf("Instances = %+v, want one normal instance", r.Instances)
				}

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
				if !maps.Equal(inst.Labels, wantLabels) {
					t.Errorf("Labels = %+v, want %+v", inst.Labels, wantLabels)
				}
				wantActiveAt, err := time.Parse(time.RFC3339, "2026-08-31T08:02:50Z")
				if err != nil {
					t.Fatalf("test setup: %v", err)
				}
				if !inst.ActiveAt.Equal(wantActiveAt) {
					t.Errorf("ActiveAt = %v, want %v", inst.ActiveAt, wantActiveAt)
				}
				if inst.Value != "" {
					t.Errorf("Value = %q, want empty string", inst.Value)
				}
			},
		},
		{
			fixture:       "state_paused.json",
			wantRules:     1,
			wantInstances: 0,
			checkFirst: func(t *testing.T, r StateRule) {
				if !r.IsPaused {
					t.Errorf("IsPaused = false, want true")
				}
				if !r.LastEvaluation.IsZero() {
					t.Errorf("LastEvaluation = %v, want zero time", r.LastEvaluation)
				}
				if r.Health != "ok" || r.State != "inactive" {
					t.Errorf("Health/State = %q/%q, want ok/inactive", r.Health, r.State)
				}
			},
		},
		{
			fixture:       "state_health_error.json",
			wantRules:     1,
			wantInstances: 1,
			checkFirst: func(t *testing.T, r StateRule) {
				if r.Health != "error" {
					t.Errorf("Health = %q, want error", r.Health)
				}
				if r.LastError == "" {
					t.Errorf("LastError is empty, want a message")
				}
				if len(r.Instances) != 1 || r.Instances[0].State != StateError {
					t.Fatalf("Instances = %+v, want one error instance", r.Instances)
				}
			},
		},
		{
			fixture:       "state_health_nodata.json",
			wantRules:     1,
			wantInstances: 1,
			checkFirst: func(t *testing.T, r StateRule) {
				if r.Health != "nodata" {
					t.Errorf("Health = %q, want nodata", r.Health)
				}
				if len(r.Instances) != 1 || r.Instances[0].State != StateNodata {
					t.Fatalf("Instances = %+v, want one nodata instance", r.Instances)
				}
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
				if !ok || errInst.State != StateNormal {
					t.Errorf(`want an instance with State=normal Reason="Error", got %+v`, byReason["Error"])
				}
				nodataInst, ok := byReason["NoData"]
				if !ok || nodataInst.State != StateNormal {
					t.Errorf(`want an instance with State=normal Reason="NoData", got %+v`, byReason["NoData"])
				}
				plain, ok := byReason[""]
				if !ok || plain.State != StateNormal {
					t.Errorf(`want a plain State=normal Reason="" instance, got %+v`, byReason[""])
				}
			},
		},
		{
			fixture:       "state_missing_optional.json",
			wantRules:     1,
			wantInstances: 0,
			checkFirst: func(t *testing.T, r StateRule) {
				if r.Instances != nil {
					t.Errorf("Instances = %+v, want nil", r.Instances)
				}
				if r.Totals != nil {
					t.Errorf("Totals = %+v, want nil", r.Totals)
				}
			},
		},
		{
			fixture:       "state_only_active_instances.json",
			wantRules:     1,
			wantInstances: 1,
			checkFirst: func(t *testing.T, r StateRule) {
				if len(r.Instances) != 1 || r.Instances[0].State != StateFiring {
					t.Fatalf("Instances = %+v, want one firing instance", r.Instances)
				}
				if r.Totals["normal"] == 0 {
					t.Errorf(`Totals["normal"] = 0, want >0 (the totals/instances mismatch this fixture exists to capture)`)
				}
			},
		},
	}

	for _, c := range cases {
		t.Run(c.fixture, func(t *testing.T) {
			rules, err := ParseState(readFixture(t, c.fixture))
			if err != nil {
				t.Fatalf("ParseState(%s): unexpected error: %v", c.fixture, err)
			}
			if len(rules) != c.wantRules {
				t.Fatalf("ParseState(%s): got %d rules, want %d", c.fixture, len(rules), c.wantRules)
			}
			if got := len(rules[0].Instances); got != c.wantInstances {
				t.Fatalf("ParseState(%s): got %d instances, want %d", c.fixture, got, c.wantInstances)
			}
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
			if err == nil {
				t.Fatalf("ParseState(%s): expected an error, got none", c.fixture)
			}
			for _, want := range c.wantContains {
				if !strings.Contains(err.Error(), want) {
					t.Errorf("ParseState(%s): error %q does not mention %q", c.fixture, err.Error(), want)
				}
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
			if err == nil {
				t.Errorf("normalizeInstanceState(%q): expected an error, got none", c.in)
			}
			continue
		}
		if err != nil {
			t.Errorf("normalizeInstanceState(%q): unexpected error: %v", c.in, err)
			continue
		}
		if state != c.wantState || reason != c.wantReason {
			t.Errorf("normalizeInstanceState(%q) = (%q, %q), want (%q, %q)", c.in, state, reason, c.wantState, c.wantReason)
		}
	}
}

func TestInstanceKey(t *testing.T) {
	a := instanceKey(map[string]string{"b": "2", "a": "1"})
	b := instanceKey(map[string]string{"a": "1", "b": "2"})
	if a != b {
		t.Errorf("instanceKey order-independence: %q != %q", a, b)
	}
	if a != "a=1\nb=2\n" {
		t.Errorf("instanceKey = %q, want a=1\\nb=2\\n", a)
	}

	diff := instanceKey(map[string]string{"a": "1", "b": "3"})
	if a == diff {
		t.Errorf("instanceKey should differ when a label value differs")
	}

	if instanceKey(nil) != "" {
		t.Errorf("instanceKey(nil) = %q, want empty string", instanceKey(nil))
	}
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
			if err != nil {
				t.Fatalf("ParseState: %v", err)
			}
			if len(rules) != 1 {
				t.Fatalf("rules = %+v, want one", rules)
			}
		})
	}
}

// labels is optional at the INSTANCE level (opt(m, "labels", ...) in
// parseInstance), distinct from the rule-level labels state_missing_optional.json
// already covers — an instance can exist with no labels of its own.
func TestParseState_InstanceWithoutLabelsParses(t *testing.T) {
	body := minimalStateBody(`,"alerts":[{"state":"Normal","activeAt":"2026-01-01T00:00:00Z"}]`)
	rules, err := ParseState(body)
	if err != nil {
		t.Fatalf("ParseState: %v", err)
	}
	if len(rules) != 1 || len(rules[0].Instances) != 1 {
		t.Fatalf("rules = %+v, want one rule with one instance", rules)
	}
	if got := rules[0].Instances[0].Labels; len(got) != 0 {
		t.Errorf("Instance.Labels = %v, want empty/nil", got)
	}
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
	if err := json.Unmarshal(base, &top); err != nil {
		t.Fatalf("synthesize: %v", err)
	}
	var data map[string]json.RawMessage
	if err := json.Unmarshal(top["data"], &data); err != nil {
		t.Fatalf("synthesize: %v", err)
	}
	var groups []map[string]json.RawMessage
	if err := json.Unmarshal(data["groups"], &groups); err != nil {
		t.Fatalf("synthesize: %v", err)
	}
	var rules []map[string]json.RawMessage
	if err := json.Unmarshal(groups[0]["rules"], &rules); err != nil {
		t.Fatalf("synthesize: %v", err)
	}
	var alerts []map[string]json.RawMessage
	if err := json.Unmarshal(rules[0]["alerts"], &alerts); err != nil {
		t.Fatalf("synthesize: %v", err)
	}
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
	if err != nil {
		t.Fatalf("synthesize: %v", err)
	}
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
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return json.RawMessage(b)
}

func TestParseState_HighCardinality(t *testing.T) {
	body := synthesizeHighCardinalityState(t, 445, 2004)
	if !bytes.Contains(body, []byte("alerting-0")) {
		t.Fatalf("synthesized body missing expected content")
	}

	rules, err := ParseState(body)
	if err != nil {
		t.Fatalf("ParseState: unexpected error: %v", err)
	}
	if len(rules) != 1 {
		t.Fatalf("got %d rules, want 1", len(rules))
	}
	r := rules[0]
	if len(r.Instances) != 445+2004 {
		t.Fatalf("got %d instances, want %d", len(r.Instances), 445+2004)
	}

	var firing, normal int
	for _, inst := range r.Instances {
		switch inst.State {
		case StateFiring:
			firing++
		case StateNormal:
			normal++
		default:
			t.Fatalf("unexpected instance state %q", inst.State)
		}
	}
	if firing != 445 || normal != 2004 {
		t.Fatalf("got firing=%d normal=%d, want firing=445 normal=2004", firing, normal)
	}

	// Each synthesized instance carries a distinct "instance" label; confirm
	// Labels actually made it through parsing (not just State) by checking
	// instanceKey produces one unique key per instance, with no collisions.
	seen := make(map[string]bool, len(r.Instances))
	for _, inst := range r.Instances {
		if inst.Labels == nil {
			t.Fatalf("instance has nil Labels")
		}
		k := instanceKey(inst.Labels)
		if seen[k] {
			t.Fatalf("duplicate instance key %q", k)
		}
		seen[k] = true
	}
	if len(seen) != 445+2004 {
		t.Fatalf("got %d unique instance keys, want %d", len(seen), 445+2004)
	}
}
