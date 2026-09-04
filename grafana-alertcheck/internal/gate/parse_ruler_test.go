package gate

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestParseDefinitions_RulerRules(t *testing.T) {
	defs, err := ParseDefinitions(readFixture(t, "ruler_rules.json"))
	require.NoError(t, err)

	byUID := map[string]Definition{}
	for _, d := range defs {
		require.Equalf(t, KindGrafanaManaged, d.Kind, "rule %q: Kind", d.UID)
		byUID[d.UID] = d
	}

	// The real 2-way duplicate title: same folder, same group, same title,
	// distinct UIDs — only uid: can tell them apart.
	a, ok := byUID["rule0000006a"]
	require.True(t, ok, "missing rule0000006a")
	b, ok := byUID["rule0000006b"]
	require.True(t, ok, "missing rule0000006b")
	require.Equal(t, a.Title, b.Title, "duplicate-title pair should share Title")
	require.Equal(t, a.Folder, b.Folder, "duplicate-title pair should share Folder")
	require.Equal(t, a.Group, b.Group, "duplicate-title pair should share Group")
	require.NotEqual(t, a.UID, b.UID, "duplicate-title pair should have distinct UIDs")

	// The 3 real paused rules.
	pausedUIDs := []string{"rule0000002", "rule0000007", "rule0000008"}
	for _, uid := range pausedUIDs {
		d, ok := byUID[uid]
		require.True(t, ok, "missing paused rule %q", uid)
		require.Truef(t, d.IsPaused, "rule %q: IsPaused = false, want true", uid)
	}

	// for:1d and the derived for:1w rule.
	dayRule, ok := byUID["rule0000009"]
	require.Truef(t, ok, "missing rule0000009")
	require.Equal(t, 24*time.Hour, dayRule.For)
	weekRule, ok := byUID["rule0000010"]
	require.Truef(t, ok, "missing rule0000010")
	require.Equal(t, 7*24*time.Hour, weekRule.For)

	// Identity shared with testdata/state_paused.json.
	shared := byUID["rule0000002"]
	require.Equal(t, "folder0000002", shared.FolderUID)
}

func TestParseDefinitions_DatasourceManaged(t *testing.T) {
	defs, err := ParseDefinitions(readFixture(t, "ruler_datasource_managed.json"))
	require.NoError(t, err)
	require.Len(t, defs, 1)
	require.Equal(t, KindDatasourceManaged, defs[0].Kind)
	require.Equal(t, 5*time.Minute, defs[0].For)
	// A datasource-managed rule has no uid in this shape; its only identity
	// is the Prometheus "alert" name — a synthetic UID would be invented
	// shape, and an empty Title would make Resolve's refusal-by-name
	// unreachable.
	require.Equal(t, "ExampleTargetDown", defs[0].Title)
	require.Empty(t, defs[0].UID, "this shape has no uid")
}

func TestParseDefinitions_Recording(t *testing.T) {
	defs, err := ParseDefinitions(readFixture(t, "ruler_recording.json"))
	require.NoError(t, err)
	require.Len(t, defs, 1)
	d := defs[0]
	require.Equal(t, KindRecording, d.Kind)
	require.Equal(t, "rule0000011", d.UID)
	// The fixture deliberately omits no_data_state/exec_err_state/is_paused/
	// intervalSeconds/namespace_uid — alerting-only concepts a recording
	// rule may not carry. Requiring them would brick ParseDefinitions for
	// every named rule in the same response over one recording rule
	// elsewhere in the fleet; they must come back as zero values, not errors.
	require.Empty(t, d.NoDataState)
	require.Empty(t, d.ExecErrState)
	require.False(t, d.IsPaused)
	require.Zero(t, d.IntervalSeconds)
	require.Empty(t, d.FolderUID)
}
