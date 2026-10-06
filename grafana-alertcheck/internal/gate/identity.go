package gate

import "encoding/json"

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
