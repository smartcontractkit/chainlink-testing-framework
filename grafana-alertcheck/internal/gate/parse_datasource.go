package gate

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// ParseDatasourceRules parses a Prometheus/vmalert rules response
// (/api/prometheus/{uid}/api/v1/rules) into StateRules. It shares parseInstance
// with the Grafana parser but uses the datasource vocabulary: lowercase instance
// states, health "err" (not "error"), and a zero lastEvaluation is allowed
// (liveness treats zero as maximally stale). A missing or unparseable required
// field is an error, never a zero value.
func ParseDatasourceRules(body []byte, dsUID, dsName string) ([]StateRule, error) {
	var top map[string]json.RawMessage
	if err := json.Unmarshal(body, &top); err != nil {
		return nil, fmt.Errorf("datasource rules response: %w", err)
	}
	var dataRaw json.RawMessage
	if err := req(top, "data", &dataRaw); err != nil {
		return nil, fmt.Errorf("datasource rules response: %w", err)
	}
	var data map[string]json.RawMessage
	if err := json.Unmarshal(dataRaw, &data); err != nil {
		return nil, fmt.Errorf("datasource rules response: data: %w", err)
	}
	var groupsRaw []json.RawMessage
	if err := req(data, "groups", &groupsRaw); err != nil {
		return nil, fmt.Errorf("datasource rules response: %w", err)
	}

	var rules []StateRule
	for gi, groupRaw := range groupsRaw {
		var group map[string]json.RawMessage
		if err := json.Unmarshal(groupRaw, &group); err != nil {
			return nil, fmt.Errorf("datasource rules response: group %d: %w", gi, err)
		}
		var groupName, file string
		if err := req(group, "name", &groupName); err != nil {
			return nil, fmt.Errorf("datasource rules response: group %d: %w", gi, err)
		}
		if err := opt(group, "file", &file); err != nil {
			return nil, fmt.Errorf("datasource rules response: group %q: %w", groupName, err)
		}
		var intervalSeconds float64
		if err := opt(group, "interval", &intervalSeconds); err != nil {
			return nil, fmt.Errorf("datasource rules response: group %q: %w", groupName, err)
		}
		interval := time.Duration(intervalSeconds * float64(time.Second))

		var rulesRaw []json.RawMessage
		if err := req(group, "rules", &rulesRaw); err != nil {
			return nil, fmt.Errorf("datasource rules response: group %q: %w", groupName, err)
		}
		for ri, ruleRaw := range rulesRaw {
			rule, err := parseDatasourceRule(ruleRaw, dsUID, groupName, file, interval)
			if err != nil {
				return nil, fmt.Errorf("datasource rules response: group %q: rule %d: %w", groupName, ri, err)
			}
			rules = append(rules, rule)
		}
	}
	return rules, nil
}

func parseDatasourceRule(raw json.RawMessage, dsUID, group, file string, interval time.Duration) (StateRule, error) {
	var m map[string]json.RawMessage
	if err := json.Unmarshal(raw, &m); err != nil {
		return StateRule{}, err
	}
	var name, ruleType string
	if err := req(m, "name", &name); err != nil {
		return StateRule{}, err
	}
	if err := req(m, "type", &ruleType); err != nil {
		return StateRule{}, fmt.Errorf("rule %q: %w", name, err)
	}

	r := StateRule{
		Key:           ruleKey(dsUID, group, name, ""),
		Title:         name,
		Group:         group,
		File:          file,
		Interval:      interval,
		DatasourceUID: dsUID,
		Type:          ruleType,
	}
	if err := opt(m, "query", &r.Query); err != nil {
		return StateRule{}, fmt.Errorf("rule %q: %w", name, err)
	}
	if err := opt(m, "labels", &r.Labels); err != nil {
		return StateRule{}, fmt.Errorf("rule %q: %w", name, err)
	}
	var durationSeconds float64
	if err := opt(m, "duration", &durationSeconds); err != nil {
		return StateRule{}, fmt.Errorf("rule %q: %w", name, err)
	}
	r.For = time.Duration(durationSeconds * float64(time.Second))

	// Recording rules carry no state, health or alerts; DefinitionsFromDatasource
	// drops them, so parsing only the shared fields keeps a recording rule from
	// failing on fields it was never going to have.
	if ruleType == "recording" {
		return r, nil
	}

	if err := req(m, "health", &r.Health); err != nil {
		return StateRule{}, fmt.Errorf("rule %q: %w", name, err)
	}
	if r.Health == "err" {
		r.Health = "error"
	}
	var lastEvalStr string
	if err := opt(m, "lastEvaluation", &lastEvalStr); err != nil {
		return StateRule{}, fmt.Errorf("rule %q: %w", name, err)
	}
	if lastEvalStr != "" {
		lastEval, err := time.Parse(time.RFC3339, lastEvalStr)
		if err != nil {
			return StateRule{}, fmt.Errorf("rule %q: lastEvaluation: %w", name, err)
		}
		r.LastEvaluation = lastEval
	}

	if err := opt(m, "state", &r.State); err != nil {
		return StateRule{}, fmt.Errorf("rule %q: %w", name, err)
	}
	var alertsRaw []json.RawMessage
	if err := opt(m, "alerts", &alertsRaw); err != nil {
		return StateRule{}, fmt.Errorf("rule %q: %w", name, err)
	}
	instances := make([]Instance, 0, len(alertsRaw))
	for ii, ar := range alertsRaw {
		inst, err := parseInstanceWith(ar, normalizeDatasourceInstanceState)
		if err != nil {
			return StateRule{}, fmt.Errorf("rule %q: instance %d: %w", name, ii, err)
		}
		instances = append(instances, inst)
	}
	r.Instances = instances
	return r, nil
}

// datasourceInstanceStates is the strict datasource instance-state vocabulary.
// Only the two active states exist — a resolved instance is absent from the
// response, not reported as normal.
var datasourceInstanceStates = map[string]State{
	"firing":  StateFiring,
	"pending": StatePending,
}

func normalizeDatasourceInstanceState(s string) (State, string, error) {
	base, reason := s, ""
	if i := strings.Index(s, " ("); i != -1 && strings.HasSuffix(s, ")") {
		base, reason = s[:i], s[i+2:len(s)-1]
	}
	state, ok := datasourceInstanceStates[base]
	if !ok {
		return "", "", fmt.Errorf("unrecognized datasource instance state %q", s)
	}
	return state, reason, nil
}

// DefinitionsFromDatasource converts datasource rule states into Definitions,
// keeping only alerting rules. A datasource-managed rule has no pause signal and
// no uid, so PauseObservable is false and UID stays empty.
func DefinitionsFromDatasource(rules []StateRule, dsUID, dsName string) []Definition {
	var defs []Definition
	for _, r := range rules {
		if r.Type != "alerting" {
			continue
		}
		defs = append(defs, Definition{
			Key:             r.Key,
			Title:           r.Title,
			Group:           r.Group,
			File:            r.File,
			For:             r.For,
			Labels:          r.Labels,
			Kind:            KindDatasourceManaged,
			DatasourceUID:   dsUID,
			DatasourceName:  dsName,
			PauseObservable: false,
			IntervalSeconds: int(r.Interval / time.Second),
		})
	}
	return defs
}
