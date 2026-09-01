package gate

import (
	"fmt"
	"strings"
	"testing"
)

func rulerDefs(t *testing.T) []Definition {
	t.Helper()
	defs, err := ParseDefinitions(readFixture(t, "ruler_rules.json"))
	if err != nil {
		t.Fatalf("ParseDefinitions: unexpected error: %v", err)
	}
	return defs
}

func TestResolve_SingleMatch(t *testing.T) {
	defs := rulerDefs(t)
	resolved, notes, err := Resolve(defs, []string{"example_workflow_paused_rule"}, "")
	if err != nil {
		t.Fatalf("Resolve: unexpected error: %v", err)
	}
	if len(notes) != 0 {
		t.Errorf("notes = %v, want none", notes)
	}
	if len(resolved) != 1 || resolved[0].UID != "rule0000007" {
		t.Fatalf("resolved = %+v, want [rule0000007]", resolved)
	}
}

func TestResolve_UIDForm(t *testing.T) {
	defs := rulerDefs(t)
	resolved, _, err := Resolve(defs, []string{"uid:rule0000006a"}, "")
	if err != nil {
		t.Fatalf("Resolve: unexpected error: %v", err)
	}
	if len(resolved) != 1 || resolved[0].UID != "rule0000006a" {
		t.Fatalf("resolved = %+v, want [rule0000006a]", resolved)
	}
}

func TestResolve_FolderGroupTitleForm(t *testing.T) {
	defs := rulerDefs(t)
	resolved, _, err := Resolve(defs, []string{"Example-Zone-A/Gateway/Example No Gateways Available"}, "")
	if err == nil {
		t.Fatalf("Resolve: want ambiguous error (real 2-way collision), got resolved=%+v", resolved)
	}
	if !strings.Contains(err.Error(), "matches 2 rules") {
		t.Fatalf("Resolve: error = %q, want it to report 2 matches", err)
	}
	if !strings.Contains(err.Error(), "uid:rule0000006a") || !strings.Contains(err.Error(), "uid:rule0000006b") {
		t.Fatalf("Resolve: error = %q, want both candidate uids listed", err)
	}
}

func TestResolve_TrueCollisionResolvesByUID(t *testing.T) {
	defs := rulerDefs(t)
	resolved, _, err := Resolve(defs, []string{"uid:rule0000006a", "uid:rule0000006b"}, "")
	if err != nil {
		t.Fatalf("Resolve: unexpected error: %v", err)
	}
	if len(resolved) != 2 {
		t.Fatalf("resolved = %+v, want 2 distinct rules", resolved)
	}
}

func TestResolve_NoMatch(t *testing.T) {
	defs := rulerDefs(t)
	_, _, err := Resolve(defs, []string{"Does Not Exist"}, "")
	if err == nil {
		t.Fatal("Resolve: want error for unknown name")
	}
	if !strings.Contains(err.Error(), "no rule matched") || !strings.Contains(err.Error(), "list") {
		t.Errorf("Resolve: error = %q, want it to name 'no rule matched' and point at 'list'", err)
	}
}

func TestResolve_NoMatchSubstringSuggestion(t *testing.T) {
	defs := rulerDefs(t)
	_, _, err := Resolve(defs, []string{"paused rule"}, "")
	if err == nil {
		t.Fatal("Resolve: want error for unknown name")
	}
	if !strings.Contains(err.Error(), "did you mean") || !strings.Contains(err.Error(), "Example Paused Rule") {
		t.Errorf("Resolve: error = %q, want a case-insensitive substring suggestion", err)
	}
}

func TestResolve_RefusesDatasourceManaged(t *testing.T) {
	defs, err := ParseDefinitions(readFixture(t, "ruler_datasource_managed.json"))
	if err != nil {
		t.Fatalf("ParseDefinitions: unexpected error: %v", err)
	}
	_, _, err = Resolve(defs, []string{"ExampleTargetDown"}, "")
	if err == nil {
		t.Fatal("Resolve: want refusal for a datasource-managed rule")
	}
	if !strings.Contains(err.Error(), "datasource-managed") {
		t.Errorf("Resolve: error = %q, want it to name the datasource-managed kind", err)
	}
}

func TestResolve_RefusesRecording(t *testing.T) {
	defs, err := ParseDefinitions(readFixture(t, "ruler_recording.json"))
	if err != nil {
		t.Fatalf("ParseDefinitions: unexpected error: %v", err)
	}
	_, _, err = Resolve(defs, []string{"uid:rule0000011"}, "")
	if err == nil {
		t.Fatal("Resolve: want refusal for a recording rule")
	}
	if !strings.Contains(err.Error(), "recording rule") {
		t.Errorf("Resolve: error = %q, want it to name the recording kind", err)
	}
}

func TestResolve_RejectsEmptySegments(t *testing.T) {
	defs := rulerDefs(t)
	cases := []string{"/Title", "Folder/", "a//b", "//", "/"}
	for _, name := range cases {
		t.Run(name, func(t *testing.T) {
			_, _, err := Resolve(defs, []string{name}, "")
			if err == nil {
				t.Fatalf("Resolve(%q): want error for an empty /-separated segment", name)
			}
			if !strings.Contains(err.Error(), "empty") {
				t.Errorf("Resolve(%q): error = %q, want it to name the empty segment", name, err)
			}
		})
	}
}

func TestResolve_UIDEmptySuffix(t *testing.T) {
	// ruler_datasource_managed.json's only rule has UID == "" (P1.3: this
	// shape has no uid at all). "uid:" with an empty suffix must not match it
	// — that would report the misleading "datasource-managed rule, not
	// supported" for what is really a typo'd/empty uid.
	defs, err := ParseDefinitions(readFixture(t, "ruler_datasource_managed.json"))
	if err != nil {
		t.Fatalf("ParseDefinitions: unexpected error: %v", err)
	}
	_, _, err = Resolve(defs, []string{"uid:"}, "")
	if err == nil {
		t.Fatal("Resolve: want error for an empty uid: suffix")
	}
	if !strings.Contains(err.Error(), "no rule has this uid") {
		t.Errorf("Resolve: error = %q, want it to say no rule has this uid", err)
	}
	if strings.Contains(err.Error(), "datasource-managed") {
		t.Errorf("Resolve: error = %q, must not misreport this as a datasource-managed refusal", err)
	}
}

func TestResolve_UnsupportedKindsExcludedFromNoMatchSurfaces(t *testing.T) {
	dsDefs, err := ParseDefinitions(readFixture(t, "ruler_datasource_managed.json"))
	if err != nil {
		t.Fatalf("ParseDefinitions(datasource_managed): unexpected error: %v", err)
	}
	recDefs, err := ParseDefinitions(readFixture(t, "ruler_recording.json"))
	if err != nil {
		t.Fatalf("ParseDefinitions(recording): unexpected error: %v", err)
	}
	supported := rulerDefs(t)
	combined := append(append(append([]Definition{}, supported...), dsDefs...), recDefs...)

	_, _, err = Resolve(combined, []string{"Example"}, "")
	if err == nil {
		t.Fatal("Resolve: want a no-match error for a name matching no title exactly")
	}

	wantCount := fmt.Sprintf("(%d rules available", len(supported))
	if !strings.Contains(err.Error(), wantCount) {
		t.Errorf("Resolve: error = %q, want the available count scoped to the %d supported rules, not the %d combined", err, len(supported), len(combined))
	}
	if strings.Contains(err.Error(), "ExampleTargetDown") {
		t.Errorf("Resolve: error = %q, must not suggest the datasource-managed rule", err)
	}
	if strings.Contains(err.Error(), "example:recorded_metric:rate5m") {
		t.Errorf("Resolve: error = %q, must not suggest the recording rule", err)
	}
	if !strings.Contains(err.Error(), "Example Paused Rule") {
		t.Errorf("Resolve: error = %q, want it to still suggest a matching supported rule", err)
	}
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
	if err != nil {
		t.Fatalf("Resolve: unexpected error: %v", err)
	}
	if len(resolved) != 1 || resolved[0].UID != "supported-1" {
		t.Fatalf("resolved = %+v, want the supported rule alone, no ambiguity", resolved)
	}
}

func TestResolve_CollapseByUIDGivesNoteNotError(t *testing.T) {
	defs := rulerDefs(t)
	// The bare title and its Folder/Group/Title spelling both name the same
	// rule (rule0000007) — a duplicate-name copy mistake, not an error
	// (§17.3).
	resolved, notes, err := Resolve(defs, []string{
		"example_workflow_paused_rule",
		"ExampleObservability/Example Auth Production/example_workflow_paused_rule",
	}, "")
	if err != nil {
		t.Fatalf("Resolve: unexpected error: %v", err)
	}
	if len(resolved) != 1 || resolved[0].UID != "rule0000007" {
		t.Fatalf("resolved = %+v, want exactly one rule0000007", resolved)
	}
	if len(notes) != 1 {
		t.Fatalf("notes = %v, want exactly one collapse note", notes)
	}
}

func TestResolve_MinObservedCountIsPostCollapse(t *testing.T) {
	defs := rulerDefs(t)
	names := []string{
		"example_workflow_paused_rule",
		"ExampleObservability/Example Auth Production/example_workflow_paused_rule", // duplicate of the same rule
		"Example Paused Rule",
	}
	resolved, notes, err := Resolve(defs, names, "")
	if err != nil {
		t.Fatalf("Resolve: unexpected error: %v", err)
	}
	// §17.3: the default MinObserved must come from len(resolved) (2 distinct
	// rules) — never len(names) (3 input lines), which would be unsatisfiable.
	if len(resolved) != 2 {
		t.Fatalf("resolved = %+v, want 2 distinct rules after collapse", resolved)
	}
	if len(notes) != 1 {
		t.Fatalf("notes = %v, want exactly one collapse note", notes)
	}
}

func TestResolve_EmptyAndBlankLinesDiscarded(t *testing.T) {
	defs := rulerDefs(t)
	resolved, _, err := Resolve(defs, []string{"", "  ", "example_workflow_paused_rule", "   \t  "}, "")
	if err != nil {
		t.Fatalf("Resolve: unexpected error: %v", err)
	}
	if len(resolved) != 1 || resolved[0].UID != "rule0000007" {
		t.Fatalf("resolved = %+v, want [rule0000007]", resolved)
	}
}

func TestResolve_FolderScopesBareTitle(t *testing.T) {
	defs := rulerDefs(t)
	// Bare title, scoped to the wrong folder — must not match.
	_, _, err := Resolve(defs, []string{"example_workflow_paused_rule"}, "Example-Zone-A")
	if err == nil {
		t.Fatal("Resolve: want no-match when folder scope excludes the only candidate")
	}

	// Scoped to the right folder — must match.
	resolved, _, err := Resolve(defs, []string{"example_workflow_paused_rule"}, "ExampleObservability")
	if err != nil {
		t.Fatalf("Resolve: unexpected error: %v", err)
	}
	if len(resolved) != 1 || resolved[0].UID != "rule0000007" {
		t.Fatalf("resolved = %+v, want [rule0000007]", resolved)
	}
}
