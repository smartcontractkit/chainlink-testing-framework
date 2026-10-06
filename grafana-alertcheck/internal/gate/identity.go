package gate

import (
	"encoding/json"
	"fmt"
	"strings"
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

// rejectDuplicateKeys fails closed on two definitions in one list sharing a
// key. It is the key: selector's guard: a key shared by two rules cannot pick
// one of them.
func rejectDuplicateKeys(defs []Definition) error {
	byKey := make(map[string][]Definition, len(defs))
	for _, d := range defs {
		key := defKey(d)
		byKey[key] = append(byKey[key], d)
	}
	for key, group := range byKey {
		if len(group) > 1 {
			return duplicateKeyError(key, group)
		}
	}
	return nil
}

// rejectSharedSelectedKeys fails closed when a SELECTED rule's identity is
// shared in the loaded inventory. The state query is by
// datasource/group/name/file, so both siblings come back and stateRuleByKey
// would reduce whichever the backend lists first — not necessarily the one the
// selection matched. Narrowing the selection therefore does not make the rule
// observable; it must be excluded.
func rejectSharedSelectedKeys(all, selected []Definition) error {
	byKey := make(map[string][]Definition, len(all))
	for _, d := range all {
		key := defKey(d)
		byKey[key] = append(byKey[key], d)
	}
	for _, d := range selected {
		key := defKey(d)
		if group := byKey[key]; len(group) > 1 {
			return duplicateKeyError(key, group)
		}
	}
	return nil
}

// duplicateKeyError explains a key collision. A datasource collision is the
// interesting one: the rules are genuinely distinct (a backend may serve two
// same-name rules in one group/file) but cannot be told apart by the state API.
func duplicateKeyError(key string, defs []Definition) error {
	descs := make([]string, len(defs))
	for i, d := range defs {
		descs[i] = describeRule(d)
	}
	if defs[0].Kind == KindDatasourceManaged {
		return fmt.Errorf(
			"%d datasource-managed rules share the identity %s (%s) and cannot be told apart: the state query is by datasource/group/name/file, so this rule cannot be observed; exclude it from the selection",
			len(defs), key, strings.Join(descs, "; "))
	}
	return fmt.Errorf("%d rules share the key %s (%s)", len(defs), key, strings.Join(descs, "; "))
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
