package gate

import (
	"fmt"
	"slices"
	"sort"
	"strings"
)

// Resolve turns the operator-supplied alert names into resolved Definitions
// (§17). Order is load-bearing (§17.3):
//
//  1. Trim each name.
//  2. Discard empty lines.
//  3. Resolve each name to a UID (this is what resolveOne does).
//  4. Collapse the result by UID — two names hitting the same rule is a note,
//     never an error (almost always a copy mistake, and a message costs the
//     user less than a failure).
//
// The caller-visible consequence: len(resolved) is the count *after* the
// collapse. A later phase's MinObserved must default from that length, never
// from len(names) — using the input line count would make one rule named
// twice turn an achievable default into an unsatisfiable one (§17.3).
func Resolve(defs []Definition, names []string, folder string) (resolved []Definition, notes []string, err error) {
	seenUID := map[string]string{} // uid -> the first input name that resolved to it
	for _, raw := range names {
		name := strings.TrimSpace(raw)
		if name == "" {
			continue
		}

		def, rerr := resolveOne(defs, name, folder)
		if rerr != nil {
			return nil, nil, rerr
		}

		if firstName, ok := seenUID[def.UID]; ok {
			notes = append(notes, fmt.Sprintf(
				"%q and %q both resolve to %s (uid:%s); counted once", firstName, name, def.Title, def.UID))
			continue
		}
		seenUID[def.UID] = name
		resolved = append(resolved, def)
	}
	return resolved, notes, nil
}

// resolveOne resolves a single trimmed, non-empty name against defs (§17.1):
// one match wins outright, zero is an error with suggestions, two or more is
// an error listing every candidate. folder scopes a bare title (no "/" in the
// name) to one folder; it is ignored for the "Folder/Title" and
// "Folder/Group/Title" forms, which already name their own folder.
//
// Policy on unsupported kinds (datasource-managed, recording) — decided here
// because §17.1 only says to refuse them, not how they interact with the
// no-match/ambiguous surfaces: a name can still match an unsupported rule (so
// naming one by title still gets the specific, named refusal, not a bare "no
// match"), but only *supported* candidates count for ambiguity — an
// unsupported rule sharing a title with a supported one is resolved silently
// in the supported rule's favor rather than reported as ambiguous — and the
// "%d rules available" count and substring suggestions in a genuine no-match
// are scoped to supported rules only, so an unsupported rule never inflates
// or pollutes either. uid: is always exact regardless of kind (typically
// copy-pasted from `list`, which already shows Kind).
func resolveOne(defs []Definition, name, folder string) (Definition, error) {
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
		// rules have no uid at all, P1.3), so matching on an empty suffix
		// would silently hit one of those and report a misleading
		// kind-specific refusal for what is really an empty/typo'd uid. This
		// deliberately does not go through noMatchError: that function's
		// substring suggestion would degenerate to an empty needle, which
		// strings.Contains matches against every title — printing the whole
		// fleet instead of a real suggestion.
		return Definition{}, fmt.Errorf("no rule matched %q: no rule has this uid (run 'grafana-alertcheck list' to see uids)", name)
	}

	wantFolder, wantGroup, wantTitle, err := classifyForm(name, folder)
	if err != nil {
		return Definition{}, err
	}

	var supportedCandidates, unsupportedCandidates []Definition
	for _, d := range defs {
		if wantFolder != "" && d.Folder != wantFolder {
			continue
		}
		if wantGroup != "" && d.Group != wantGroup {
			continue
		}
		if d.Title != wantTitle {
			continue
		}
		if d.Kind == KindGrafanaManaged {
			supportedCandidates = append(supportedCandidates, d)
		} else {
			unsupportedCandidates = append(unsupportedCandidates, d)
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
		return Definition{}, noMatchError(supportedDefs(defs), name, wantTitle)
	}
}

// supportedDefs filters out the two kinds §17.1 refuses. Only these
// participate in name-based matching, the no-match rule count, and substring
// suggestions (see the policy note on resolveOne).
func supportedDefs(defs []Definition) []Definition {
	out := make([]Definition, 0, len(defs))
	for _, d := range defs {
		if d.Kind == KindGrafanaManaged {
			out = append(out, d)
		}
	}
	return out
}

// classifyForm splits name into the Title | Folder/Title | Folder/Group/Title
// forms (§17). A bare title is scoped by folder when the caller supplied one;
// the two- and three-segment forms already carry their own folder and ignore
// it.
//
// Every segment must be non-empty. Without this, "/Title" would parse as an
// empty wantFolder — silently dropping the folder filter and matching
// unscoped, a fail-open — and "Folder/" would parse as an empty wantTitle,
// which would then feed noMatchError's substring search an empty needle that
// matches every title.
func classifyForm(name, folder string) (wantFolder, wantGroup, wantTitle string, err error) {
	parts := strings.Split(name, "/")
	if slices.Contains(parts, "") {
		return "", "", "", fmt.Errorf("no rule matched %q: empty /-separated segment (want Title, Folder/Title, or Folder/Group/Title)", name)
	}
	switch len(parts) {
	case 1:
		return folder, "", parts[0], nil
	case 2:
		return parts[0], "", parts[1], nil
	case 3:
		return parts[0], parts[1], parts[2], nil
	default:
		return "", "", "", fmt.Errorf("no rule matched %q: too many /-separated segments (want Title, Folder/Title, or Folder/Group/Title)", name)
	}
}

// refuseUnsupportedKind rejects the two kinds §17.1 names explicitly with a
// clear, specific error — distinct from "no match" and from "ambiguous" — so
// an operator who names a recording or datasource-managed rule learns why,
// not just that nothing matched.
func refuseUnsupportedKind(name string, d Definition) (Definition, error) {
	switch d.Kind {
	case KindDatasourceManaged:
		return Definition{}, fmt.Errorf("%q resolves to %s, a datasource-managed rule, which is not supported", name, d.Title)
	case KindRecording:
		return Definition{}, fmt.Errorf("%q resolves to %s, a recording rule, which is not supported", name, d.Title)
	default:
		return d, nil
	}
}

// noMatchError reports a no-match with the count of rules the gate could see
// and, per Context decision 4, case-insensitive substring matches in place of
// the source plan's cut Levenshtein suggestions (§17.2).
func noMatchError(defs []Definition, name, wantTitle string) error {
	msg := fmt.Sprintf("no rule matched %q (%d rules available; run 'grafana-alertcheck list' to see titles)", name, len(defs))

	needle := strings.ToLower(wantTitle)
	var subs []string
	for _, d := range defs {
		if strings.Contains(strings.ToLower(d.Title), needle) {
			subs = append(subs, fmt.Sprintf("%s/%s/%s", d.Folder, d.Group, d.Title))
		}
	}
	if len(subs) > 0 {
		sort.Strings(subs)
		msg += fmt.Sprintf("; did you mean: %s", strings.Join(subs, ", "))
	}
	return fmt.Errorf("%s", msg)
}

// ambiguousError lists every candidate with its folder, its group, and the
// full copyable Folder/Group/Title (§17.1) — including the uid: form, which
// resolves unambiguously on the next attempt.
func ambiguousError(name string, candidates []Definition) error {
	sorted := append([]Definition(nil), candidates...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].UID < sorted[j].UID })

	var b strings.Builder
	fmt.Fprintf(&b, "%q matches %d rules; use uid: or the full Folder/Group/Title:", name, len(sorted))
	for _, d := range sorted {
		fmt.Fprintf(&b, "\n  %s/%s/%s (uid:%s)", d.Folder, d.Group, d.Title, d.UID)
	}
	return fmt.Errorf("%s", b.String())
}
