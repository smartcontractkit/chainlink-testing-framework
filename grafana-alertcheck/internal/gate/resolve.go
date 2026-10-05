package gate

import (
	"fmt"
	"slices"
	"sort"
	"strings"
)

// Resolve turns the operator-supplied alert names into resolved Definitions.
// Order is load-bearing:
//
//  1. Trim each name.
//  2. Discard empty lines.
//  3. Resolve each name to a rule (this is what resolveOne does).
//  4. Collapse the result by key — two names hitting the same rule is a note,
//     never an error (almost always a copy mistake, and a message costs the
//     user less than a failure).
//
// The caller-visible consequence: len(resolved) is the count *after* the
// collapse, and MinObserved must default from that length, never from
// len(names) — using the input line count would make one rule named twice turn
// an achievable default into an unsatisfiable one.
func Resolve(defs []Definition, names []string, folder string) (resolved []Definition, notes []string, err error) {
	seenKey := map[string]string{} // key -> the first input name that resolved to it
	for _, raw := range names {
		name := strings.TrimSpace(raw)
		if name == "" {
			continue
		}

		def, rerr := resolveOne(defs, name, folder)
		if rerr != nil {
			return nil, nil, rerr
		}

		key := defKey(def)
		if firstName, ok := seenKey[key]; ok {
			notes = append(notes, fmt.Sprintf(
				"%q and %q both resolve to %s (%s); counted once", firstName, name, def.Title, ruleRefLabel(def)))
			continue
		}
		seenKey[key] = name
		resolved = append(resolved, def)
	}
	return resolved, notes, nil
}

// resolveOne resolves a single trimmed, non-empty name against defs: one match
// wins outright, zero is an error with suggestions, two or more is an error
// listing every candidate. folder scopes a bare Grafana title to one folder; it
// is ignored for the /-separated forms and for datasource-managed rules, which
// name their own group and datasource.
//
// A name matches both kinds: Grafana forms are Title | Folder/Title |
// Folder/Group/Title, datasource forms are Title | Group/Title |
// DatasourceName/Group/Title, and an exact key: is unambiguous across both.
// Unsupported kinds (recording, and a datasource-managed rule with no
// datasource) are refused, and only supported candidates count for ambiguity.
func resolveOne(defs []Definition, name, folder string) (Definition, error) {
	if key, ok := strings.CutPrefix(name, "key:"); ok {
		if key != "" {
			for _, d := range defs {
				if defKey(d) == key {
					return refuseUnsupportedKind(name, d)
				}
			}
		}
		return Definition{}, fmt.Errorf("no rule matched %q: no rule has this key (run 'grafana-alertcheck list' to see keys)", name)
	}

	if uid, ok := strings.CutPrefix(name, "uid:"); ok {
		if uid != "" {
			for _, d := range defs {
				if d.UID == uid {
					return refuseUnsupportedKind(name, d)
				}
			}
		}
		// uid == "" falls through to the same message as "not found": several
		// Definition kinds legitimately carry UID == "" (datasource-managed
		// rules have no uid at all), so matching on an empty suffix
		// would silently hit one of those and report a misleading
		// kind-specific refusal for what is really an empty/typo'd uid. This
		// deliberately does not go through noMatchError: that function's
		// substring suggestion would degenerate to an empty needle, which
		// strings.Contains matches against every title — printing the whole
		// fleet instead of a real suggestion.
		return Definition{}, fmt.Errorf("no rule matched %q: no rule has this uid (run 'grafana-alertcheck list' to see uids)", name)
	}

	parts, err := parseNameForm(name)
	if err != nil {
		return Definition{}, err
	}

	var supportedCandidates, unsupportedCandidates []Definition
	for _, d := range defs {
		if !matchesName(d, parts, folder) {
			continue
		}
		if d.Kind == KindRecording || (d.Kind == KindDatasourceManaged && d.DatasourceUID == "") {
			unsupportedCandidates = append(unsupportedCandidates, d)
		} else {
			supportedCandidates = append(supportedCandidates, d)
		}
	}

	switch {
	case len(supportedCandidates) == 1:
		return supportedCandidates[0], nil
	case len(supportedCandidates) > 1:
		return Definition{}, ambiguousError(name, supportedCandidates)
	case len(unsupportedCandidates) > 0:
		return refuseUnsupportedKind(name, unsupportedCandidates[0])
	default:
		return Definition{}, noMatchError(supportedDefs(defs), name, parts[len(parts)-1])
	}
}

// matchesName reports whether d matches the /-separated name form. A bare
// Grafana title is scoped by folder; datasource forms never use folder.
func matchesName(d Definition, parts []string, folder string) bool {
	switch len(parts) {
	case 1:
		if d.Kind == KindGrafanaManaged {
			return (folder == "" || d.Folder == folder) && d.Title == parts[0]
		}
		return d.Title == parts[0]
	case 2:
		if d.Kind == KindGrafanaManaged {
			return d.Folder == parts[0] && d.Title == parts[1]
		}
		return d.Group == parts[0] && d.Title == parts[1]
	case 3:
		if d.Kind == KindGrafanaManaged {
			return d.Folder == parts[0] && d.Group == parts[1] && d.Title == parts[2]
		}
		return d.DatasourceName == parts[0] && d.Group == parts[1] && d.Title == parts[2]
	}
	return false
}

// supportedDefs filters out the refused kinds. Only these participate in
// name-based matching, the no-match rule count, and substring suggestions.
func supportedDefs(defs []Definition) []Definition {
	out := make([]Definition, 0, len(defs))
	for _, d := range defs {
		if d.Kind == KindRecording || (d.Kind == KindDatasourceManaged && d.DatasourceUID == "") {
			continue
		}
		out = append(out, d)
	}
	return out
}

// parseNameForm splits name into 1..3 /-separated segments. Every segment must
// be non-empty: without this, "/Title" would parse as an empty first segment —
// silently dropping the filter and matching unscoped, a fail-open — and
// "Folder/" would parse as an empty title, feeding noMatchError's substring
// search an empty needle that matches every title.
func parseNameForm(name string) ([]string, error) {
	parts := strings.Split(name, "/")
	if slices.Contains(parts, "") {
		return nil, fmt.Errorf("no rule matched %q: empty /-separated segment (want Title, Group/Title, or Datasource/Group/Title)", name)
	}
	if len(parts) > 3 {
		return nil, fmt.Errorf("no rule matched %q: too many /-separated segments (want Title, Group/Title, or Datasource/Group/Title)", name)
	}
	return parts, nil
}

// refuseUnsupportedKind rejects the unsupported kinds with a clear, specific
// error — distinct from "no match" and from "ambiguous".
func refuseUnsupportedKind(name string, d Definition) (Definition, error) {
	switch {
	case d.Kind == KindRecording:
		return Definition{}, fmt.Errorf("%q resolves to %s, a recording rule, which is not supported", name, d.Title)
	case d.Kind == KindDatasourceManaged && d.DatasourceUID == "":
		return Definition{}, fmt.Errorf("%q resolves to %s, a datasource-managed rule whose datasource is unknown, which is not supported", name, d.Title)
	default:
		return d, nil
	}
}

// ruleRefLabel names a rule for a note: a uid for Grafana, a copyable key for a
// datasource-managed rule.
func ruleRefLabel(d Definition) string {
	if d.UID != "" {
		return "uid:" + d.UID
	}
	return "key:" + defKey(d)
}

// noMatchError reports a no-match with the count of supported rules and
// case-insensitive substring suggestions.
func noMatchError(defs []Definition, name, wantTitle string) error {
	msg := fmt.Sprintf("no rule matched %q (%d rules available; run 'grafana-alertcheck list' to see titles)",
		name, len(defs))

	needle := strings.ToLower(wantTitle)
	var subs []string
	for _, d := range defs {
		if strings.Contains(strings.ToLower(d.Title), needle) {
			subs = append(subs, suggestionLabel(d))
		}
	}
	if len(subs) > 0 {
		sort.Strings(subs)
		msg += fmt.Sprintf("; did you mean: %s", strings.Join(subs, ", "))
	}
	return fmt.Errorf("%s", msg)
}

// suggestionLabel is the copyable name form for a supported rule.
func suggestionLabel(d Definition) string {
	if d.Kind == KindDatasourceManaged {
		return fmt.Sprintf("%s/%s/%s", d.DatasourceName, d.Group, d.Title)
	}
	return fmt.Sprintf("%s/%s/%s", d.Folder, d.Group, d.Title)
}

// ambiguousError lists every candidate with its source, its group, and the full
// copyable name — including the key:/uid: form, which resolves unambiguously on
// the next attempt.
func ambiguousError(name string, candidates []Definition) error {
	sorted := append([]Definition(nil), candidates...)
	sort.Slice(sorted, func(i, j int) bool { return defKey(sorted[i]) < defKey(sorted[j]) })

	var b strings.Builder
	fmt.Fprintf(&b, "%q matches %d rules; use key: or the full name:", name, len(sorted))
	for _, d := range sorted {
		if d.Kind == KindDatasourceManaged {
			fmt.Fprintf(&b, "\n  %s/%s/%s (datasource %s, key:%s)", d.DatasourceName, d.Group, d.Title, d.DatasourceName, defKey(d))
			continue
		}
		fmt.Fprintf(&b, "\n  %s/%s/%s (uid:%s)", d.Folder, d.Group, d.Title, d.UID)
	}
	return fmt.Errorf("%s", b.String())
}
