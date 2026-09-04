package gate

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"
)

// RuleKind classifies a ruler-endpoint rule by shape, not by name. Resolve
// rejects KindDatasourceManaged and KindRecording, but only for rules a user
// actually named — ParseDefinitions itself never rejects.
type RuleKind int

const (
	KindGrafanaManaged RuleKind = iota
	KindDatasourceManaged
	KindRecording
)

// Definition is one rule from the ruler endpoint
// (/api/ruler/grafana/api/v1/rules). IntervalSeconds, NoDataState and
// ExecErrState live inside the grafana_alert block and are only populated for
// KindGrafanaManaged — a datasource-managed rule has no such block by
// definition. relativeTimeRange and keep_firing_for are deliberately not
// parsed: nothing in the gate reads them.
type Definition struct {
	UID, Title, Folder, FolderUID, Group string
	For                                  time.Duration
	IntervalSeconds                      int
	NoDataState                          string
	ExecErrState                         string
	IsPaused                             bool
	Kind                                 RuleKind
}

// ParseDefinitions strictly parses a ruler-endpoint response body
// (map[namespace][]group) into its rule definitions.
func ParseDefinitions(body []byte) ([]Definition, error) {
	var namespaces map[string][]json.RawMessage
	if err := json.Unmarshal(body, &namespaces); err != nil {
		return nil, fmt.Errorf("ruler response: %w", err)
	}

	// Map iteration order is nondeterministic; sort namespace names so
	// ParseDefinitions' output order is stable across calls — Resolve's
	// candidate listings and the golden tests depend on that.
	names := make([]string, 0, len(namespaces))
	for name := range namespaces {
		names = append(names, name)
	}
	sort.Strings(names)

	var defs []Definition
	for _, folder := range names {
		for gi, groupRaw := range namespaces[folder] {
			var group map[string]json.RawMessage
			if err := json.Unmarshal(groupRaw, &group); err != nil {
				return nil, fmt.Errorf("ruler response: namespace %q: group %d: %w", folder, gi, err)
			}

			var groupName string
			if err := req(group, "name", &groupName); err != nil {
				return nil, fmt.Errorf("ruler response: namespace %q: group %d: %w", folder, gi, err)
			}

			var rulesRaw []json.RawMessage
			if err := req(group, "rules", &rulesRaw); err != nil {
				return nil, fmt.Errorf("ruler response: namespace %q: group %q: %w", folder, groupName, err)
			}

			for ri, ruleRaw := range rulesRaw {
				def, err := parseDefinition(ruleRaw, folder, groupName)
				if err != nil {
					return nil, fmt.Errorf("ruler response: namespace %q: group %q: rule %d: %w", folder, groupName, ri, err)
				}
				defs = append(defs, def)
			}
		}
	}
	return defs, nil
}

func parseDefinition(raw json.RawMessage, folder, group string) (Definition, error) {
	var m map[string]json.RawMessage
	if err := json.Unmarshal(raw, &m); err != nil {
		return Definition{}, fmt.Errorf("%w", err)
	}

	var forStr string
	if err := opt(m, "for", &forStr); err != nil {
		return Definition{}, err
	}
	forDur, err := ParsePromDuration(forStr)
	if err != nil {
		return Definition{}, fmt.Errorf("for: %w", err)
	}

	var gaRaw json.RawMessage
	if err := opt(m, "grafana_alert", &gaRaw); err != nil {
		return Definition{}, err
	}
	if gaRaw == nil {
		// No grafana_alert block at all: a datasource-managed (native
		// Prometheus-format) rule. Its identity is the Prometheus rule
		// name — "alert" for an alerting rule, "record" for a recording
		// one — never a synthetic UID (Grafana's ruler API gives this
		// shape no uid at all; inventing one would be inventing shape).
		def := Definition{Folder: folder, Group: group, For: forDur, Kind: KindDatasourceManaged}
		if err := opt(m, "alert", &def.Title); err != nil {
			return Definition{}, err
		}
		if def.Title == "" {
			if err := opt(m, "record", &def.Title); err != nil {
				return Definition{}, err
			}
		}
		return def, nil
	}

	var ga map[string]json.RawMessage
	if err := json.Unmarshal(gaRaw, &ga); err != nil {
		return Definition{}, fmt.Errorf("grafana_alert: %w", err)
	}

	// uid identifies the rule regardless of kind — every grafana_alert
	// object Grafana emits, alerting or recording, carries one.
	var uid string
	if err := req(ga, "uid", &uid); err != nil {
		return Definition{}, fmt.Errorf("grafana_alert: %w", err)
	}
	def := Definition{Folder: folder, Group: group, For: forDur, UID: uid}

	// Classify by the presence of "record" before requiring anything else.
	// no_data_state/exec_err_state/is_paused/intervalSeconds are alerting-only
	// concepts a recording rule may not carry at all — its real shape is
	// unverified, none exist in the fleet capture — and Resolve refuses this
	// Kind categorically before any of this would gate a release.
	// Strict-parsing a recording rule into a hard error over fields it was
	// never going to use would brick `list` and every resolve for rules nobody
	// named.
	var record json.RawMessage
	if err := opt(ga, "record", &record); err != nil {
		return Definition{}, fmt.Errorf("rule %q: grafana_alert: %w", uid, err)
	}
	if record != nil {
		def.Kind = KindRecording
		if err := opt(ga, "title", &def.Title); err != nil {
			return Definition{}, fmt.Errorf("rule %q: grafana_alert: %w", uid, err)
		}
		if err := opt(ga, "namespace_uid", &def.FolderUID); err != nil {
			return Definition{}, fmt.Errorf("rule %q: grafana_alert: %w", uid, err)
		}
		if err := opt(ga, "intervalSeconds", &def.IntervalSeconds); err != nil {
			return Definition{}, fmt.Errorf("rule %q: grafana_alert: %w", uid, err)
		}
		if err := opt(ga, "is_paused", &def.IsPaused); err != nil {
			return Definition{}, fmt.Errorf("rule %q: grafana_alert: %w", uid, err)
		}
		return def, nil
	}

	def.Kind = KindGrafanaManaged
	if err := req(ga, "title", &def.Title); err != nil {
		return Definition{}, fmt.Errorf("rule %q: grafana_alert: %w", uid, err)
	}
	if err := req(ga, "namespace_uid", &def.FolderUID); err != nil {
		return Definition{}, fmt.Errorf("rule %q: grafana_alert: %w", uid, err)
	}
	if err := req(ga, "intervalSeconds", &def.IntervalSeconds); err != nil {
		return Definition{}, fmt.Errorf("rule %q: grafana_alert: %w", uid, err)
	}
	if err := req(ga, "no_data_state", &def.NoDataState); err != nil {
		return Definition{}, fmt.Errorf("rule %q: grafana_alert: %w", uid, err)
	}
	if err := req(ga, "exec_err_state", &def.ExecErrState); err != nil {
		return Definition{}, fmt.Errorf("rule %q: grafana_alert: %w", uid, err)
	}
	if err := req(ga, "is_paused", &def.IsPaused); err != nil {
		return Definition{}, fmt.Errorf("rule %q: grafana_alert: %w", uid, err)
	}

	return def, nil
}
