package gate

import "encoding/json"

// dsKeyPrefix marks a datasource-managed key. The key, not the prefix, is the
// identity; the prefix only lets a caller that has lost the Definition (an
// absent-rule poll, a log read) still tell the two source kinds apart.
const dsKeyPrefix = "ds:"

// ruleKey is the one map key for a rule across both source kinds. Grafana-managed
// rules keep their uid. Datasource-managed rules have no uid, so they get a
// JSON-encoded tuple; JSON keeps group/name separators from colliding. This is a
// key, not an identity the API gave us — Definition.UID stays empty for ds rules.
func ruleKey(dsUID, group, name, uid string) string {
	if uid != "" {
		return uid
	}
	b, _ := json.Marshal([3]string{dsUID, group, name})
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
