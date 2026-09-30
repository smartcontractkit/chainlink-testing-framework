package gate

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

// LabelMatcher is one exact-match label requirement.
type LabelMatcher struct {
	Key, Value string
}

// SelectByLabels resolves a label selection: a rule is selected when it
// carries every include pair, and dropped when it carries any exclude pair.
// A missing label never triggers an exclude (unlike a Prometheus != matcher);
// it only fails an include. Zero matches and a matched unsupported kind are
// errors, never a smaller watch set.
func SelectByLabels(defs []Definition, include, exclude []LabelMatcher) ([]Definition, error) {
	matchedInclude := 0
	selected := make([]Definition, 0, len(defs))
	for _, d := range defs {
		if !matchesAll(d.Labels, include) {
			continue
		}
		matchedInclude++
		if d.Kind != KindGrafanaManaged {
			return nil, fmt.Errorf("label selection matches %q, a %s, which is not supported", d.Title, kindName(d.Kind))
		}
		if matchesAny(d.Labels, exclude) {
			continue
		}
		selected = append(selected, d)
	}
	if matchedInclude == 0 {
		return nil, fmt.Errorf("no rule matches the include labels %s (%d rules visible)",
			formatLabelMatchers(include), len(defs))
	}
	if len(selected) == 0 {
		return nil, fmt.Errorf("all %d rule(s) matching the include labels %s are dropped by the exclude labels %s",
			matchedInclude, formatLabelMatchers(include), formatLabelMatchers(exclude))
	}
	return selected, nil
}

// resolveAlertSet picks the alert set: names or labels, minus --exclude-alerts.
func resolveAlertSet(defs []Definition, names []string, include, exclude []LabelMatcher,
	excludeAlerts []string, folder string) ([]Definition, []string, error) {

	var (
		selected []Definition
		notes    []string
		err      error
	)
	if len(include) > 0 {
		selected, err = SelectByLabels(defs, include, exclude)
	} else {
		selected, notes, err = Resolve(defs, names, folder)
	}
	if err != nil {
		return nil, nil, err
	}
	selected, err = subtractExcluded(defs, selected, excludeAlerts, folder)
	if err != nil {
		return nil, nil, err
	}
	return selected, notes, nil
}

// subtractExcluded drops every --exclude-alerts name from the selected set.
// Names resolve like --alerts, so a typo is an error; a name outside the
// selection removes nothing.
func subtractExcluded(all, selected []Definition, excludeAlerts []string, folder string) ([]Definition, error) {
	if len(excludeAlerts) == 0 {
		return selected, nil
	}
	excluded, _, err := Resolve(all, excludeAlerts, folder)
	if err != nil {
		return nil, fmt.Errorf("--exclude-alerts: %w", err)
	}
	drop := make(map[string]bool, len(excluded))
	for _, d := range excluded {
		drop[d.UID] = true
	}
	out := make([]Definition, 0, len(selected))
	for _, d := range selected {
		if !drop[d.UID] {
			out = append(out, d)
		}
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("--exclude-alerts drops every selected rule; refusing to run with an empty alert set")
	}
	return out, nil
}

// labelSelectionHeader names the include/exclude pairs and any --exclude-alerts.
func labelSelectionHeader(include, exclude []LabelMatcher, excludeCount int) string {
	header := fmt.Sprintf("alerts matching --include-labels %s", formatLabelMatchers(include))
	if len(exclude) > 0 {
		header += fmt.Sprintf(" and --exclude-labels %s", formatLabelMatchers(exclude))
	}
	if excludeCount > 0 {
		header += " minus --exclude-alerts"
	}
	return header
}

// printLabelSelection lists the matched rules, one per line.
func printLabelSelection(w io.Writer, resolved []Definition, include, exclude []LabelMatcher, excludeCount int) {
	fmt.Fprintf(w, "%s:\n", labelSelectionHeader(include, exclude, excludeCount))
	for _, d := range resolved {
		fmt.Fprintf(w, "  - %s (%s)\n", d.Title, d.UID)
	}
}

func matchesAll(labels map[string]string, matchers []LabelMatcher) bool {
	for _, m := range matchers {
		// The key check is load-bearing: without it an empty matcher value
		// would match every rule missing the key.
		if v, ok := labels[m.Key]; !ok || v != m.Value {
			return false
		}
	}
	return true
}

// validateSelection enforces the alert-set rules shared by watch and check:
// names and labels are mutually exclusive, exclusions need inclusions, and a
// label selection cannot be scoped by --folder. emptySelection phrases the
// "neither mode named anything" error, which differs per command.
func validateSelection(cmd string, named int, include, exclude []LabelMatcher, folder, emptySelection string) error {
	switch {
	case named > 0 && (len(include) > 0 || len(exclude) > 0):
		return fmt.Errorf("%s: --alerts cannot be combined with label selection", cmd)
	case len(exclude) > 0 && len(include) == 0:
		return fmt.Errorf("%s: --exclude-labels requires --include-labels; refusing to watch every rule", cmd)
	case named == 0 && len(include) == 0:
		return errors.New(emptySelection)
	case (len(include) > 0 || len(exclude) > 0) && folder != "":
		return fmt.Errorf("%s: --folder scopes alert names and cannot be combined with label selection", cmd)
	}
	return nil
}

func matchesAny(labels map[string]string, matchers []LabelMatcher) bool {
	for _, m := range matchers {
		if v, ok := labels[m.Key]; ok && v == m.Value {
			return true
		}
	}
	return false
}

func formatLabelMatchers(ms []LabelMatcher) string {
	parts := make([]string, len(ms))
	for i, m := range ms {
		parts[i] = m.Key + "=" + m.Value
	}
	return strings.Join(parts, ",")
}

func kindName(k RuleKind) string {
	switch k {
	case KindDatasourceManaged:
		return "datasource-managed rule"
	case KindRecording:
		return "recording rule"
	default:
		return "grafana-managed rule"
	}
}
