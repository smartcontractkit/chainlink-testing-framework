package gate

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func selectedUIDs(defs []Definition) []string {
	out := make([]string, len(defs))
	for i, d := range defs {
		out[i] = d.UID
	}
	return out
}

func TestSelectByLabels_IncludeIsAnExactAND(t *testing.T) {
	defs := rulerDefs(t)

	selected, err := SelectByLabels(defs, []LabelMatcher{{Key: "team", Value: "example-team"}}, nil)
	require.NoError(t, err)
	require.Equal(t, []string{"rule0000006a", "rule0000006b", "rule0000009", "rule0000010"}, selectedUIDs(selected))

	selected, err = SelectByLabels(defs, []LabelMatcher{
		{Key: "team", Value: "example-team"},
		{Key: "severity", Value: "warning"},
	}, nil)
	require.NoError(t, err)
	require.Equal(t, []string{"rule0000009", "rule0000010"}, selectedUIDs(selected))
}

// A rule that does not carry an include label never matches, even when the
// label is missing rather than different.
func TestSelectByLabels_MissingIncludeLabelDoesNotMatch(t *testing.T) {
	defs := rulerDefs(t)

	selected, err := SelectByLabels(defs, []LabelMatcher{{Key: "zone", Value: "zone-a"}}, nil)
	require.NoError(t, err)
	require.Equal(t, []string{"rule0000006a", "rule0000006b"}, selectedUIDs(selected))

	_, err = SelectByLabels(defs, []LabelMatcher{{Key: "zone", Value: "zone-c"}}, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "no rule matches the include labels")
}

// Exclude drops rules that carry the pair; a missing label is never a reason
// to drop (unlike a Prometheus != matcher).
func TestSelectByLabels_ExcludeDropsCarriersOnly(t *testing.T) {
	defs := rulerDefs(t)

	selected, err := SelectByLabels(defs,
		[]LabelMatcher{{Key: "env", Value: "production"}},
		[]LabelMatcher{{Key: "severity", Value: "critical"}})
	require.NoError(t, err)
	// rule0000007 has env=production and no severity: the missing label must
	// keep it, not exclude it. Selection preserves the definitions' order.
	require.Equal(t, []string{"rule0000009", "rule0000010", "rule0000007"}, selectedUIDs(selected))
}

func TestSelectByLabels_AllExcludedIsAnError(t *testing.T) {
	defs := rulerDefs(t)

	_, err := SelectByLabels(defs,
		[]LabelMatcher{{Key: "team", Value: "example-team"}},
		[]LabelMatcher{{Key: "team", Value: "example-team"}})
	require.Error(t, err)
	require.Contains(t, err.Error(), "dropped by the exclude labels")
}

func TestSelectByLabels_RefusesMatchedUnsupportedKinds(t *testing.T) {
	dsDefs, err := ParseDefinitions(readFixture(t, "ruler_datasource_managed.json"))
	require.NoError(t, err)
	_, err = SelectByLabels(dsDefs, []LabelMatcher{{Key: "severity", Value: "warning"}}, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "datasource-managed")

	recDefs := []Definition{{UID: "r1", Title: "recorded", Kind: KindRecording, Labels: map[string]string{"env": "production"}}}
	_, err = SelectByLabels(recDefs, []LabelMatcher{{Key: "env", Value: "production"}}, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "recording rule")
}

// A caller-supplied empty matcher value must still require the key: a missing
// label is absent, not equal to the empty string.
func TestSelectByLabels_EmptyValueStillRequiresTheKey(t *testing.T) {
	defs := []Definition{
		{UID: "absent", Title: "absent", Kind: KindGrafanaManaged, Labels: map[string]string{}},
		{UID: "empty", Title: "empty", Kind: KindGrafanaManaged, Labels: map[string]string{"env": ""}},
	}
	selected, err := SelectByLabels(defs, []LabelMatcher{{Key: "env", Value: ""}}, nil)
	require.NoError(t, err)
	require.Equal(t, []string{"empty"}, selectedUIDs(selected))
}

// An empty include list matches every rule; the exclude list is then the only
// filter. Validation refuses this mode (it would watch the whole fleet), so
// this pins the pure function's behaviour only.
func TestSelectByLabels_EmptyIncludeIsNotAnError(t *testing.T) {
	defs := rulerDefs(t)
	selected, err := SelectByLabels(defs, nil, []LabelMatcher{{Key: "env", Value: "production"}})
	require.NoError(t, err)
	require.Equal(t, []string{"rule0000002"}, selectedUIDs(selected))
}

// Two distinct datasource rules can share one identity (same datasource,
// group, name and file, differing only by labels/query). A selection that
// matches both cannot observe them distinctly and must fail closed; one that
// matches a single rule is fine.
func TestResolveAlertSet_DuplicateIdentityOnlyFailsWhenBothSelected(t *testing.T) {
	a := Definition{
		Key: "ds:k", Title: "Same", Group: "G", Kind: KindDatasourceManaged,
		DatasourceUID: "vm", DatasourceName: "VM",
		Labels: map[string]string{"product": "ccip", "severity": "critical"},
	}
	b := a
	b.Labels = map[string]string{"product": "ccip", "severity": "warning"}
	defs := []Definition{a, b}

	_, _, err := resolveAlertSet(defs, nil, []LabelMatcher{{Key: "product", Value: "ccip"}}, nil, nil, "")
	require.Error(t, err)
	require.Contains(t, err.Error(), "share the identity")
	require.Contains(t, err.Error(), "narrow the selection")

	selected, _, err := resolveAlertSet(defs, nil, []LabelMatcher{{Key: "severity", Value: "critical"}}, nil, nil, "")
	require.NoError(t, err)
	require.Len(t, selected, 1)
}

// --exclude-alerts subtracts from an enumerated set: the names resolve like
// --alerts, so uid: forms work, and the result keeps the input order.
func TestResolveAlertSet_SubtractsExcludedNames(t *testing.T) {
	defs := rulerDefs(t)
	selected, _, err := resolveAlertSet(defs,
		[]string{"uid:rule0000009", "uid:rule0000010"}, nil, nil,
		[]string{"uid:rule0000010"}, "")
	require.NoError(t, err)
	require.Equal(t, []string{"rule0000009"}, selectedUIDs(selected))
}

// --exclude-alerts combines with label selection too: the label match is
// computed first, then the named rules are subtracted from it.
func TestResolveAlertSet_SubtractsExcludedNamesFromLabelSelection(t *testing.T) {
	defs := rulerDefs(t)
	selected, _, err := resolveAlertSet(defs, nil,
		[]LabelMatcher{{Key: "team", Value: "example-team"}}, nil,
		[]string{"uid:rule0000009", "uid:rule0000010"}, "")
	require.NoError(t, err)
	require.Equal(t, []string{"rule0000006a", "rule0000006b"}, selectedUIDs(selected))
}

// A typo in the excluded list is an error, never a silent no-op: the operator
// asked to remove something and the gate cannot prove it did.
func TestResolveAlertSet_UnknownExcludedNameFails(t *testing.T) {
	defs := rulerDefs(t)
	_, _, err := resolveAlertSet(defs, []string{"uid:rule0000009"}, nil, nil,
		[]string{"Does Not Exist"}, "")
	require.Error(t, err)
	require.Contains(t, err.Error(), "--exclude-alerts")
}

// Excluding the whole selection would leave nothing to observe; that is exit 2,
// not a vacuous pass.
func TestResolveAlertSet_ExcludingEverythingFails(t *testing.T) {
	defs := rulerDefs(t)
	_, _, err := resolveAlertSet(defs, []string{"uid:rule0000009"}, nil, nil,
		[]string{"uid:rule0000009"}, "")
	require.Error(t, err)
	require.Contains(t, err.Error(), "drops every selected rule")
}

// A datasource rule has no uid, so the listing must show its key instead of an
// empty pair of parentheses.
func TestPrintLabelSelection_DatasourceShowsKey(t *testing.T) {
	defs := []Definition{{
		Key: ruleKey("vm", "G", "A", "f", ""), Title: "A", Group: "G",
		Kind: KindDatasourceManaged, DatasourceUID: "vm", DatasourceName: "VM",
	}}
	var b strings.Builder
	printLabelSelection(&b, defs, []LabelMatcher{{Key: "k", Value: "v"}}, nil, 0)
	out := b.String()
	require.Contains(t, out, "  - A (ds:")
	require.NotContains(t, out, "()")
}

func TestPrintLabelSelection(t *testing.T) {
	defs := rulerDefs(t)
	selected, err := SelectByLabels(defs, []LabelMatcher{{Key: "severity", Value: "warning"}}, nil)
	require.NoError(t, err)

	var b strings.Builder
	printLabelSelection(&b, selected,
		[]LabelMatcher{{Key: "severity", Value: "warning"}},
		[]LabelMatcher{{Key: "env", Value: "stage"}}, 1)

	out := b.String()
	require.Contains(t, out, "alerts matching --include-labels severity=warning and --exclude-labels env=stage minus --exclude-alerts:\n")
	require.Contains(t, out, "  - Example Failure Ratio Above 10 Percent (rule0000009)\n")
	require.Contains(t, out, "  - Example Failure Ratio Above 10 Percent Weekly (rule0000010)\n")
}
