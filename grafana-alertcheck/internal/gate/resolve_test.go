package gate

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func rulerDefs(t *testing.T) []Definition {
	t.Helper()
	defs, err := ParseDefinitions(readFixture(t, "ruler_rules.json"))
	require.NoError(t, err)
	return defs
}

func TestResolve_SingleMatch(t *testing.T) {
	defs := rulerDefs(t)
	resolved, notes, err := Resolve(defs, []string{"example_workflow_paused_rule"}, "")
	require.NoError(t, err)
	require.Empty(t, notes)
	require.Len(t, resolved, 1)
	require.Equal(t, "rule0000007", resolved[0].UID)
}

func TestResolve_UIDForm(t *testing.T) {
	defs := rulerDefs(t)
	resolved, _, err := Resolve(defs, []string{"uid:rule0000006a"}, "")
	require.NoError(t, err)
	require.Len(t, resolved, 1)
	require.Equal(t, "rule0000006a", resolved[0].UID)
}

func TestResolve_FolderTitleForm(t *testing.T) {
	defs := rulerDefs(t)
	resolved, _, err := Resolve(defs, []string{"ExampleFeeds/TEMP - Example depeg alert"}, "")
	require.NoError(t, err)
	require.Len(t, resolved, 1)
	require.Equal(t, "rule0000008", resolved[0].UID)
}

func TestResolve_FolderGroupTitleForm(t *testing.T) {
	defs := rulerDefs(t)
	resolved, _, err := Resolve(defs, []string{"Example-Zone-A/Gateway/Example No Gateways Available"}, "")
	require.Error(t, err, "want ambiguous error (real 2-way collision), got resolved=%+v", resolved)
	require.Contains(t, err.Error(), "matches 2 rules")
	require.Contains(t, err.Error(), "uid:rule0000006a")
	require.Contains(t, err.Error(), "uid:rule0000006b")
}

func TestResolve_TrueCollisionResolvesByUID(t *testing.T) {
	defs := rulerDefs(t)
	resolved, _, err := Resolve(defs, []string{"uid:rule0000006a", "uid:rule0000006b"}, "")
	require.NoError(t, err)
	require.Len(t, resolved, 2)
}

func TestResolve_NoMatch(t *testing.T) {
	defs := rulerDefs(t)
	_, _, err := Resolve(defs, []string{"Does Not Exist"}, "")
	require.Error(t, err)
	require.Contains(t, err.Error(), "no rule matched")
	require.Contains(t, err.Error(), "list")
}

func TestResolve_RefusesDatasourceManaged(t *testing.T) {
	defs, err := ParseDefinitions(readFixture(t, "ruler_datasource_managed.json"))
	require.NoError(t, err)
	_, _, err = Resolve(defs, []string{"ExampleTargetDown"}, "")
	require.Error(t, err)
	require.Contains(t, err.Error(), "datasource-managed")
}

func TestResolve_RefusesRecording(t *testing.T) {
	defs, err := ParseDefinitions(readFixture(t, "ruler_recording.json"))
	require.NoError(t, err)
	// Both the uid form and the rule's real Folder/Group/Title path must match
	// and be refused specifically — a recording rule must not fall through the
	// datasource name forms to a generic no-match.
	for _, name := range []string{
		"uid:rule0000011",
		"ExampleMetrics/Recording Group/example:recorded_metric:rate5m",
		"ExampleMetrics/example:recorded_metric:rate5m",
	} {
		_, _, err := Resolve(defs, []string{name}, "")
		require.Errorf(t, err, "name %q", name)
		require.Containsf(t, err.Error(), "recording rule", "name %q", name)
	}
}

func TestResolve_RejectsEmptySegments(t *testing.T) {
	defs := rulerDefs(t)
	cases := []string{"/Title", "Folder/", "a//b", "//", "/"}
	for _, name := range cases {
		t.Run(name, func(t *testing.T) {
			_, _, err := Resolve(defs, []string{name}, "")
			require.Error(t, err)
			require.Contains(t, err.Error(), "empty")
		})
	}
}

func TestResolve_UIDEmptySuffix(t *testing.T) {
	// ruler_datasource_managed.json's only rule has UID == "" — that shape has
	// no uid at all. "uid:" with an empty suffix must not match it
	// — that would report the misleading "datasource-managed rule, not
	// supported" for what is really a typo'd/empty uid.
	defs, err := ParseDefinitions(readFixture(t, "ruler_datasource_managed.json"))
	require.NoError(t, err)
	_, _, err = Resolve(defs, []string{"uid:"}, "")
	require.Error(t, err)
	require.Contains(t, err.Error(), "no rule has this uid")
	require.NotContains(t, err.Error(), "datasource-managed")
}

func TestResolve_UnsupportedKindsExcludedFromNoMatchSurfaces(t *testing.T) {
	dsDefs, err := ParseDefinitions(readFixture(t, "ruler_datasource_managed.json"))
	require.NoError(t, err)
	recDefs, err := ParseDefinitions(readFixture(t, "ruler_recording.json"))
	require.NoError(t, err)
	supported := rulerDefs(t)
	combined := append(append(append([]Definition{}, supported...), dsDefs...), recDefs...)

	_, _, err = Resolve(combined, []string{"Example"}, "")
	require.Error(t, err, "want a no-match error for a name matching no title exactly")

	wantCount := fmt.Sprintf("(%d rules available", len(supported))
	require.Contains(t, err.Error(), wantCount)
}

func TestResolve_UnsupportedHomonymResolvesSupportedSilently(t *testing.T) {
	// Synthetic: a supported and an unsupported rule sharing an identical
	// Folder/Group/Title. Real Grafana data has no such case in the capture,
	// but the policy must not treat this as ambiguous — the unsupported rule
	// is invisible next to a same-named supported one.
	defs := []Definition{
		{UID: "supported-1", Folder: "F", Group: "G", Title: "Shared Title", Kind: KindGrafanaManaged},
		{UID: "", Folder: "F", Group: "G", Title: "Shared Title", Kind: KindDatasourceManaged},
	}
	resolved, _, err := Resolve(defs, []string{"F/G/Shared Title"}, "")
	require.NoError(t, err)
	require.Len(t, resolved, 1)
	require.Equal(t, "supported-1", resolved[0].UID)
}

func TestResolve_CollapseByUIDGivesNoteNotError(t *testing.T) {
	defs := rulerDefs(t)
	// The bare title and its Folder/Group/Title spelling both name the same
	// rule (rule0000007) — a duplicate-name copy mistake, not an error.
	resolved, notes, err := Resolve(defs, []string{
		"example_workflow_paused_rule",
		"ExampleObservability/Example Auth Production/example_workflow_paused_rule",
	}, "")
	require.NoError(t, err)
	require.Len(t, resolved, 1)
	require.Equal(t, "rule0000007", resolved[0].UID)
	require.Len(t, notes, 1)
}

// The same rule named twice with the identical string must collapse to one
// rule — distinct from the different-spellings case above.
func TestResolve_IdenticalDuplicateNameCollapsesWithNote(t *testing.T) {
	defs := rulerDefs(t)
	resolved, notes, err := Resolve(defs, []string{
		"example_workflow_paused_rule",
		"example_workflow_paused_rule",
	}, "")
	require.NoError(t, err)
	require.Len(t, resolved, 1)
	require.Equal(t, "rule0000007", resolved[0].UID)
	require.Len(t, notes, 1)
}

func TestResolve_MinObservedCountIsPostCollapse(t *testing.T) {
	defs := rulerDefs(t)
	names := []string{
		"example_workflow_paused_rule",
		"ExampleObservability/Example Auth Production/example_workflow_paused_rule", // duplicate of the same rule
		"Example Paused Rule",
	}
	resolved, notes, err := Resolve(defs, names, "")
	require.NoError(t, err)
	// The default MinObserved must come from len(resolved) (2 distinct
	// rules) — never len(names) (3 input lines), which would be unsatisfiable.
	require.Len(t, resolved, 2)
	require.Len(t, notes, 1)
}

func TestResolve_EmptyAndBlankLinesDiscarded(t *testing.T) {
	defs := rulerDefs(t)
	resolved, _, err := Resolve(defs, []string{"", "  ", "example_workflow_paused_rule", "   \t  "}, "")
	require.NoError(t, err)
	require.Len(t, resolved, 1)
	require.Equal(t, "rule0000007", resolved[0].UID)
}

func dsResolveDef(ds, dsName, group, title string) Definition {
	return Definition{
		Key: ruleKey(ds, group, title, "", ""), Title: title, Group: group,
		Kind: KindDatasourceManaged, DatasourceUID: ds, DatasourceName: dsName,
	}
}

func TestResolve_DatasourceFormsAndKey(t *testing.T) {
	defs := []Definition{dsResolveDef("vm", "VictoriaMetrics - Prod", "ExampleMetrics", "ExampleTargetDown")}

	for _, name := range []string{
		"ExampleTargetDown",
		"ExampleMetrics/ExampleTargetDown",
		"VictoriaMetrics - Prod/ExampleMetrics/ExampleTargetDown",
		"key:" + ruleKey("vm", "ExampleMetrics", "ExampleTargetDown", "", ""),
	} {
		resolved, _, err := Resolve(defs, []string{name}, "")
		require.NoErrorf(t, err, "name %q", name)
		require.Len(t, resolved, 1)
		require.Equal(t, "vm", resolved[0].DatasourceUID)
	}
}

// A datasource-managed rule's name can itself contain "/", so the whole input
// must be tried as an exact title before the segmented forms.
func TestResolve_DatasourceNameWithSlashes(t *testing.T) {
	name := "devex-cicd/prod/griddle-github: ContainersNotReady"
	defs := []Definition{{
		Key:   ruleKey("ds", "DevexCICDGriddleGitHubServiceAlerts", name, "f", ""),
		Title: name, Group: "DevexCICDGriddleGitHubServiceAlerts",
		Kind: KindDatasourceManaged, DatasourceUID: "ds", DatasourceName: "VM",
	}}
	resolved, _, err := Resolve(defs, []string{name}, "")
	require.NoError(t, err)
	require.Len(t, resolved, 1)
	require.Equal(t, name, resolved[0].Title)
}

// The same string can be one datasource rule's exact title AND a Grafana
// Folder/Title selector. That must be reported as ambiguous, not silently
// resolved to whichever interpretation is tried first.
func TestResolve_ExactTitleVsSegmentedIsAmbiguous(t *testing.T) {
	ds := Definition{
		Key:   ruleKey("vm", "G", "Platform/HighErrorRate", "f", ""),
		Title: "Platform/HighErrorRate", Group: "G",
		Kind: KindDatasourceManaged, DatasourceUID: "vm", DatasourceName: "VM",
	}
	grafana := Definition{Key: "u1", UID: "u1", Title: "HighErrorRate", Folder: "Platform", Kind: KindGrafanaManaged}
	_, _, err := Resolve([]Definition{grafana, ds}, []string{"Platform/HighErrorRate"}, "")
	require.Error(t, err)
	require.Contains(t, err.Error(), "matches 2 rules")
}

// A key shared by two distinct rules cannot select one of them.
func TestResolve_KeySharedByDistinctRulesIsAmbiguous(t *testing.T) {
	a := Definition{
		Key: "ds:k", Title: "Same", Group: "G", Kind: KindDatasourceManaged,
		DatasourceUID: "vm", DatasourceName: "VM",
	}
	b := a
	_, _, err := Resolve([]Definition{a, b}, []string{"key:ds:k"}, "")
	require.Error(t, err)
	require.Contains(t, err.Error(), "share the identity")
}

func TestResolve_DatasourceAmbiguityAcrossSources(t *testing.T) {
	defs := []Definition{
		dsResolveDef("vm-a", "A", "G", "Same"),
		dsResolveDef("vm-b", "B", "G", "Same"),
	}
	_, _, err := Resolve(defs, []string{"Same"}, "")
	require.Error(t, err)
	require.Contains(t, err.Error(), "matches 2 rules")
	require.Contains(t, err.Error(), "A")
	require.Contains(t, err.Error(), "B")
}

func TestResolve_FolderScopesBareTitle(t *testing.T) {
	defs := rulerDefs(t)
	// Bare title, scoped to the wrong folder — must not match.
	_, _, err := Resolve(defs, []string{"example_workflow_paused_rule"}, "Example-Zone-A")
	require.Error(t, err, "want no-match when folder scope excludes the only candidate")

	// Scoped to the right folder — must match.
	resolved, _, err := Resolve(defs, []string{"example_workflow_paused_rule"}, "ExampleObservability")
	require.NoError(t, err)
	require.Len(t, resolved, 1)
	require.Equal(t, "rule0000007", resolved[0].UID)
}
