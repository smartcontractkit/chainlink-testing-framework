package gate

import (
	"testing"
	"time"
)

func TestParseDefinitions_RulerRules(t *testing.T) {
	defs, err := ParseDefinitions(readFixture(t, "ruler_rules.json"))
	if err != nil {
		t.Fatalf("ParseDefinitions: unexpected error: %v", err)
	}

	byUID := map[string]Definition{}
	for _, d := range defs {
		if d.Kind != KindGrafanaManaged {
			t.Errorf("rule %q: Kind = %v, want KindGrafanaManaged", d.UID, d.Kind)
		}
		byUID[d.UID] = d
	}

	// The real 2-way duplicate title: same folder, same group, same title,
	// distinct UIDs (§17, §22.2).
	a, ok := byUID["rule0000006a"]
	if !ok {
		t.Fatalf("missing rule0000006a")
	}
	b, ok := byUID["rule0000006b"]
	if !ok {
		t.Fatalf("missing rule0000006b")
	}
	if a.Title != b.Title || a.Folder != b.Folder || a.Group != b.Group {
		t.Errorf("duplicate-title pair should share Title/Folder/Group: a=%+v b=%+v", a, b)
	}
	if a.UID == b.UID {
		t.Errorf("duplicate-title pair should have distinct UIDs")
	}

	// The 3 real paused rules.
	pausedUIDs := []string{"rule0000002", "rule0000007", "rule0000008"}
	for _, uid := range pausedUIDs {
		d, ok := byUID[uid]
		if !ok {
			t.Fatalf("missing paused rule %q", uid)
		}
		if !d.IsPaused {
			t.Errorf("rule %q: IsPaused = false, want true", uid)
		}
	}

	// for:1d and the derived for:1w rule.
	dayRule, ok := byUID["rule0000009"]
	if !ok || dayRule.For != 24*time.Hour {
		t.Fatalf("rule0000009: For = %v, want 24h (ok=%v)", dayRule.For, ok)
	}
	weekRule, ok := byUID["rule0000010"]
	if !ok || weekRule.For != 7*24*time.Hour {
		t.Fatalf("rule0000010: For = %v, want 168h (ok=%v)", weekRule.For, ok)
	}

	// Identity shared with testdata/state_paused.json.
	shared := byUID["rule0000002"]
	if shared.FolderUID != "folder0000002" {
		t.Errorf("rule0000002: FolderUID = %q, want folder0000002", shared.FolderUID)
	}
}

func TestParseDefinitions_DatasourceManaged(t *testing.T) {
	defs, err := ParseDefinitions(readFixture(t, "ruler_datasource_managed.json"))
	if err != nil {
		t.Fatalf("ParseDefinitions: unexpected error: %v", err)
	}
	if len(defs) != 1 {
		t.Fatalf("got %d definitions, want 1", len(defs))
	}
	if defs[0].Kind != KindDatasourceManaged {
		t.Errorf("Kind = %v, want KindDatasourceManaged", defs[0].Kind)
	}
	if defs[0].For != 5*time.Minute {
		t.Errorf("For = %v, want 5m", defs[0].For)
	}
	// A datasource-managed rule has no uid in this shape; its only identity
	// is the Prometheus "alert" name — a synthetic UID would be invented
	// shape, and an empty Title would make P3's refusal-by-name unreachable.
	if defs[0].Title != "ExampleTargetDown" {
		t.Errorf("Title = %q, want ExampleTargetDown", defs[0].Title)
	}
	if defs[0].UID != "" {
		t.Errorf("UID = %q, want empty (this shape has no uid)", defs[0].UID)
	}
}

func TestParseDefinitions_Recording(t *testing.T) {
	defs, err := ParseDefinitions(readFixture(t, "ruler_recording.json"))
	if err != nil {
		t.Fatalf("ParseDefinitions: unexpected error: %v", err)
	}
	if len(defs) != 1 {
		t.Fatalf("got %d definitions, want 1", len(defs))
	}
	d := defs[0]
	if d.Kind != KindRecording {
		t.Errorf("Kind = %v, want KindRecording", d.Kind)
	}
	if d.UID != "rule0000011" {
		t.Errorf("UID = %q, want rule0000011", d.UID)
	}
	// The fixture deliberately omits no_data_state/exec_err_state/is_paused/
	// intervalSeconds/namespace_uid — alerting-only concepts a recording
	// rule may not carry. Requiring them would brick ParseDefinitions for
	// every named rule in the same response over one recording rule
	// elsewhere in the fleet; they must come back as zero values, not errors.
	if d.NoDataState != "" || d.ExecErrState != "" || d.IsPaused || d.IntervalSeconds != 0 || d.FolderUID != "" {
		t.Errorf("expected zero-valued alert-only fields for a recording rule, got %+v", d)
	}
}
