package gate

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"
)

// State is the canonical instance state (P1.2a). It is distinct from the raw,
// unnormalized vocabularies the API uses at the rule level and at the instance
// level — see normalizeInstanceState.
type State string

const (
	StateNormal  State = "normal"
	StateFiring  State = "firing"
	StatePending State = "pending"
	StateNodata  State = "nodata"
	StateError   State = "error"
)

// Instance is one entry of a rule's alerts[]. State is always canonical; Reason
// is the opaque suffix of a "State (Reason)" composite ("" when the API gave a
// bare state). Reason is reporting-only except for the H2 MissingSeries routing
// done downstream in the log markers.
//
// The json tags are for the JSONL log's abnormal-instance list (P5) only —
// parsing an API response never goes through them, because parseInstance
// decodes field by field through req/opt to keep H1's presence checks explicit.
type Instance struct {
	Labels   map[string]string `json:"labels"`
	State    State             `json:"state"`
	Reason   string            `json:"reason,omitempty"`
	ActiveAt time.Time         `json:"active_at"`
	Value    string            `json:"value,omitempty"`
}

// StateRule is one rule from the state endpoint
// (/api/prometheus/grafana/api/v1/rules), fully and strictly parsed (H1).
type StateRule struct {
	UID, Title, Folder, Group string
	Interval                  time.Duration
	// State and Health are raw, lowercase, and reporting-only — never
	// classified (P1.2a). State in particular is never normalized.
	State, Health  string
	LastError      string
	LastEvaluation time.Time
	IsPaused       bool
	Instances      []Instance
	Totals         map[string]int
}

// ParseState strictly parses a state-endpoint response body into its rules.
// A missing or unparseable required field (health, state, lastEvaluation on
// each rule; interval on each group) is an error, never a zero value (H1).
func ParseState(body []byte) ([]StateRule, error) {
	var top map[string]json.RawMessage
	if err := json.Unmarshal(body, &top); err != nil {
		return nil, fmt.Errorf("state response: %w", err)
	}

	var dataRaw json.RawMessage
	if err := req(top, "data", &dataRaw); err != nil {
		return nil, fmt.Errorf("state response: %w", err)
	}
	var data map[string]json.RawMessage
	if err := json.Unmarshal(dataRaw, &data); err != nil {
		return nil, fmt.Errorf("state response: data: %w", err)
	}

	var groupsRaw []json.RawMessage
	if err := req(data, "groups", &groupsRaw); err != nil {
		return nil, fmt.Errorf("state response: %w", err)
	}

	var rules []StateRule
	for gi, groupRaw := range groupsRaw {
		var group map[string]json.RawMessage
		if err := json.Unmarshal(groupRaw, &group); err != nil {
			return nil, fmt.Errorf("state response: group %d: %w", gi, err)
		}

		var folder, groupName string
		if err := req(group, "file", &folder); err != nil {
			return nil, fmt.Errorf("state response: group %d: %w", gi, err)
		}
		if err := req(group, "name", &groupName); err != nil {
			return nil, fmt.Errorf("state response: group %d (folder %q): %w", gi, folder, err)
		}

		var intervalSeconds float64
		if err := req(group, "interval", &intervalSeconds); err != nil {
			return nil, fmt.Errorf("state response: group %q: %w", groupName, err)
		}
		interval := time.Duration(intervalSeconds * float64(time.Second))

		var rulesRaw []json.RawMessage
		if err := req(group, "rules", &rulesRaw); err != nil {
			return nil, fmt.Errorf("state response: group %q: %w", groupName, err)
		}

		for ri, ruleRaw := range rulesRaw {
			rule, err := parseStateRule(ruleRaw, folder, groupName, interval)
			if err != nil {
				return nil, fmt.Errorf("state response: group %q (folder %q): rule %d: %w", groupName, folder, ri, err)
			}
			rules = append(rules, rule)
		}
	}
	return rules, nil
}

func parseStateRule(raw json.RawMessage, folder, group string, interval time.Duration) (StateRule, error) {
	var m map[string]json.RawMessage
	if err := json.Unmarshal(raw, &m); err != nil {
		return StateRule{}, fmt.Errorf("rule: %w", err)
	}

	var uid, name string
	if err := req(m, "uid", &uid); err != nil {
		return StateRule{}, fmt.Errorf("rule: %w", err)
	}
	if err := req(m, "name", &name); err != nil {
		return StateRule{}, fmt.Errorf("rule %q: %w", uid, err)
	}

	r := StateRule{UID: uid, Title: name, Folder: folder, Group: group, Interval: interval}

	if err := req(m, "state", &r.State); err != nil {
		return StateRule{}, fmt.Errorf("rule %q: %w", uid, err)
	}
	if err := req(m, "health", &r.Health); err != nil {
		return StateRule{}, fmt.Errorf("rule %q: %w", uid, err)
	}
	// isPaused is not one of H1's four named required fields, but this parser
	// extends that contract to it: the zero-time rule below can't tell a
	// paused rule from a broken one without it, and it's the primary
	// in-window pause detector (H2/§12.2) — a silent false default would be
	// exactly the fail-open bug H1 exists to kill.
	if err := req(m, "isPaused", &r.IsPaused); err != nil {
		return StateRule{}, fmt.Errorf("rule %q: %w", uid, err)
	}

	var lastEvalStr string
	if err := req(m, "lastEvaluation", &lastEvalStr); err != nil {
		return StateRule{}, fmt.Errorf("rule %q: %w", uid, err)
	}
	lastEval, err := time.Parse(time.RFC3339, lastEvalStr)
	if err != nil {
		return StateRule{}, fmt.Errorf("rule %q: lastEvaluation: %w", uid, err)
	}
	// The zero-time rule (§2.3): only a paused rule may report the zero time.
	if lastEval.IsZero() && !r.IsPaused {
		return StateRule{}, fmt.Errorf("rule %q: lastEvaluation is the zero time but isPaused is false", uid)
	}
	r.LastEvaluation = lastEval

	if err := opt(m, "lastError", &r.LastError); err != nil {
		return StateRule{}, fmt.Errorf("rule %q: %w", uid, err)
	}

	if err := opt(m, "totals", &r.Totals); err != nil {
		return StateRule{}, fmt.Errorf("rule %q: %w", uid, err)
	}

	var alertsRaw []json.RawMessage
	if err := opt(m, "alerts", &alertsRaw); err != nil {
		return StateRule{}, fmt.Errorf("rule %q: %w", uid, err)
	}
	if alertsRaw != nil {
		instances := make([]Instance, 0, len(alertsRaw))
		for ii, ar := range alertsRaw {
			inst, err := parseInstance(ar)
			if err != nil {
				return StateRule{}, fmt.Errorf("rule %q: instance %d: %w", uid, ii, err)
			}
			instances = append(instances, inst)
		}
		r.Instances = instances
	}

	return r, nil
}

func parseInstance(raw json.RawMessage) (Instance, error) {
	var m map[string]json.RawMessage
	if err := json.Unmarshal(raw, &m); err != nil {
		return Instance{}, fmt.Errorf("%w", err)
	}

	var rawState string
	if err := req(m, "state", &rawState); err != nil {
		return Instance{}, err
	}
	state, reason, err := normalizeInstanceState(rawState)
	if err != nil {
		return Instance{}, err
	}

	// activeAt is also not in H1's named list, extended here for the same
	// reason as StateRule.IsPaused: it's the onset time BadFor (P8) measures
	// from, so a silently zeroed one would misclassify how long an instance
	// has been bad rather than failing loudly.
	var activeAtStr string
	if err := req(m, "activeAt", &activeAtStr); err != nil {
		return Instance{}, err
	}
	activeAt, err := time.Parse(time.RFC3339, activeAtStr)
	if err != nil {
		return Instance{}, fmt.Errorf("activeAt: %w", err)
	}

	inst := Instance{State: state, Reason: reason, ActiveAt: activeAt}

	if err := opt(m, "value", &inst.Value); err != nil {
		return Instance{}, err
	}
	if err := opt(m, "labels", &inst.Labels); err != nil {
		return Instance{}, err
	}

	return inst, nil
}

// baseInstanceStates is the strict 5-value allowlist for the base of an
// instance state (P1.2a). Anything else — including an unrecognized base
// inside a "Base (Reason)" composite — is a parse error (H1, §2.7 control 3).
var baseInstanceStates = map[string]State{
	"Normal":   StateNormal,
	"Alerting": StateFiring,
	"Pending":  StatePending,
	"NoData":   StateNodata,
	"Error":    StateError,
}

// normalizeInstanceState normalizes an instance-level state string into its
// canonical base and an opaque reason. The composite "State (Reason)" form is
// parsed structurally — split on the first " (" with a trailing ")" — never by
// enumerating composites, because Grafana's reason vocabulary grows across
// versions and an unknown reason must not break parsing.
func normalizeInstanceState(s string) (State, string, error) {
	base, reason := s, ""
	if i := strings.Index(s, " ("); i != -1 && strings.HasSuffix(s, ")") {
		base = s[:i]
		reason = s[i+2 : len(s)-1]
	}
	state, ok := baseInstanceStates[base]
	if !ok {
		return "", "", fmt.Errorf("unrecognized instance state %q", s)
	}
	return state, reason, nil
}

// instanceKey is a stable identity for an instance's label set: a sorted
// "k=v\n" join. Used to correlate an instance across polls without hashing.
func instanceKey(labels map[string]string) string {
	keys := make([]string, 0, len(labels))
	for k := range labels {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var b strings.Builder
	for _, k := range keys {
		b.WriteString(k)
		b.WriteByte('=')
		b.WriteString(labels[k])
		b.WriteByte('\n')
	}
	return b.String()
}
