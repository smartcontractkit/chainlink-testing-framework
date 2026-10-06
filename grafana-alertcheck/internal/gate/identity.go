package gate

import (
	"encoding/json"
	"fmt"
)

// dsKeyPrefix marks a datasource-managed key, so a caller without the
// Definition can still tell the two source kinds apart.
const dsKeyPrefix = "ds:"

// ruleKey is the one map key across both source kinds: a Grafana rule keeps its
// uid; a datasource rule has none, so it gets a JSON tuple (file included — a
// Prometheus group name is only unique within a file).
func ruleKey(dsUID, group, name, file, uid string) string {
	if uid != "" {
		return uid
	}
	b, _ := json.Marshal([4]string{dsUID, group, name, file})
	return dsKeyPrefix + string(b)
}

// defKey is a Definition's map key: Key when set, UID otherwise (a Definition
// built directly by a test may carry only UID).
func defKey(d Definition) string {
	if d.Key != "" {
		return d.Key
	}
	return d.UID
}

// loggedKey is a LoggedRule's map key, mirroring defKey.
func loggedKey(lr LoggedRule) string {
	if lr.Key != "" {
		return lr.Key
	}
	return lr.UID
}

// pollKey is a Poll's map key: rule_key when written, rule_uid otherwise (a v1
// log written before rule_key existed).
func pollKey(p Poll) string {
	if p.RuleKey != "" {
		return p.RuleKey
	}
	return p.RuleUID
}

// stateRuleKey is a StateRule's map key, mirroring defKey.
func stateRuleKey(r StateRule) string {
	if r.Key != "" {
		return r.Key
	}
	return r.UID
}

// verdictKey is a RuleVerdict's map key, mirroring defKey.
func verdictKey(v RuleVerdict) string {
	if v.RuleKey != "" {
		return v.RuleKey
	}
	return v.RuleUID
}

// rejectDuplicateKeys fails closed on two definitions sharing a key. A Grafana
// uid is unique by construction; a backend may serve two distinct
// datasource-managed rules under one (datasource, group, name, file), and those
// cannot be told apart, so a selection that matches both must be narrowed
// rather than silently observing one of them.
func rejectDuplicateKeys(defs []Definition) error {
	seen := make(map[string]Definition, len(defs))
	for _, d := range defs {
		key := defKey(d)
		prev, ok := seen[key]
		if !ok {
			seen[key] = d
			continue
		}
		if d.Kind == KindDatasourceManaged {
			return fmt.Errorf(
				"two datasource-managed rules share the identity %s: %s and %s; narrow the selection (e.g. a distinguishing label) so only one matches",
				key, describeRule(prev), describeRule(d))
		}
		return fmt.Errorf("two rules share the key %s: %s and %s", key, describeRule(prev), describeRule(d))
	}
	return nil
}

// describeRule names a definition for a duplicate-key error.
func describeRule(d Definition) string {
	if d.Kind == KindDatasourceManaged {
		src := d.DatasourceName
		if src == "" {
			src = d.DatasourceUID
		}
		return fmt.Sprintf("datasource %q, group %q, file %q, name %q", src, d.Group, d.File, d.Title)
	}
	return fmt.Sprintf("uid %s, title %q", d.UID, d.Title)
}
