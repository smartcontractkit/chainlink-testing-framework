package gate

import (
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
