package wasp_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
	"github.com/smartcontractkit/chainlink-testing-framework/wasp/comparator"
	"github.com/stretchr/testify/require"
)

func TestLokiComparator(t *testing.T) {
	label := "performance_comparator_tool"

	gen, err := wasp.NewGenerator(&wasp.Config{
		T:          t,
		LokiConfig: wasp.NewEnvLokiConfig(),
		GenName:    "vu",
		Labels: map[string]string{
			"branch": label,
			"commit": label,
		},
		CallTimeout: 100 * time.Millisecond,
		LoadType:    wasp.VU,
		Schedule: wasp.CombineAndRepeat(
			2,
			wasp.Steps(10, 1, 10, 10*time.Second),
			wasp.Plain(30, 15*time.Second),
			wasp.Steps(20, -1, 10, 5*time.Second),
		),
		VU: wasp.NewMockVU(&wasp.MockVirtualUserConfig{
			CallSleep: 50 * time.Millisecond,
		}),
	})
	require.NoError(t, err)

	currentReport := comparator.BasicReport{
		GeneratorConfigs: map[string]*wasp.Config{
			gen.Cfg.GenName: gen.Cfg,
		},
		ExecutionEnvironment: comparator.ExecutionEnvironment_Docker,
		LokiConfig:           gen.Cfg.LokiConfig,
		LokiQueries: map[string]string{
			"vu_over_time": fmt.Sprintf("max_over_time({branch=~\"%s\", commit=~\"%s\", go_test_name=~\"%s\", test_data_type=~\"stats\", gen_name=~\"%s\"} | json | unwrap current_instances [10s]) by (node_id, go_test_name, gen_name)", label, label, t.Name(), gen.Cfg.GenName),
		},
		TestName:    "TestLokiComparator",
		TestStart:   time.Now(),
		CommitOrTag: "current",
	}

	gen.Run(true)
	currentReport.TestEnd = time.Now()

	fetchErr := currentReport.Fetch()
	require.NoError(t, fetchErr, "failed to fetch current report")

	path, storeErr := currentReport.Store()
	require.NoError(t, storeErr, "failed to store current report", path)

	previousRelease := "old-one"
	previousReport := comparator.BasicReport{
		CommitOrTag: previousRelease,
		TestName:    "TestLokiComparator",
	}
	loadErr := previousReport.Load()
	require.NoError(t, loadErr, "failed to load previous report", previousRelease)

	isComparable, isComparableErrs := previousReport.IsComparable(currentReport)
	require.True(t, isComparable, "reports are not comparable", isComparableErrs)
	require.Empty(t, isComparableErrs, "reports were not comparable", isComparableErrs)
	require.Equal(t, len(currentReport.Results["vu_over_time"]), len(previousReport.Results["vu_over_time"]), "vu_over_time results are not the same length")

	// compare each result individually
	for i := range currentReport.Results["vu_over_time"] {
		require.Equal(t, currentReport.Results["vu_over_time"][i], previousReport.Results["vu_over_time"][i], "vu_over_time results are not the same for given index")
	}

	//get previous release
	//previousRelease := "something"
	//previousReport := comparator.BasicReport{}
	//loadErr := previousReport.Load()
	//require.NoError(t, loadErr, "failed to load previous report", previousRelease)

	//
	//areTheSame, perfDiffs, compareErr := previousReport.Compare(&currentReport)
	//require.NoError(t, compareErr, "failed to compare reports")
	//require.True(t, areTheSame, "performance of both reports is not the same", perfDiffs)
}
