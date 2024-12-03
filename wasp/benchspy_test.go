package wasp_test

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
	"github.com/smartcontractkit/chainlink-testing-framework/wasp/benchspy"
	"github.com/stretchr/testify/require"
)

func TestBenchSpyWithLokiQuery(t *testing.T) {
	label := "benchspy"

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

	currentReport := benchspy.StandardReport{
		BasicData: benchspy.BasicData{
			GeneratorConfigs: map[string]*wasp.Config{
				gen.Cfg.GenName: gen.Cfg,
			},
			TestName:    t.Name(),
			TestStart:   time.Now(),
			CommitOrTag: "e7fc5826a572c09f8b93df3b9f674113372ce925",
		},
		ResourceReporter: benchspy.ResourceReporter{
			ExecutionEnvironment: benchspy.ExecutionEnvironment_Docker,
		},
	}

	lokiQueryExecutor := benchspy.NewLokiQuery(
		map[string]string{
			"vu_over_time": fmt.Sprintf("max_over_time({branch=~\"%s\", commit=~\"%s\", go_test_name=~\"%s\", test_data_type=~\"stats\", gen_name=~\"%s\"} | json | unwrap current_instances [10s]) by (node_id, go_test_name, gen_name)", label, label, t.Name(), gen.Cfg.GenName),
		},
		gen.Cfg.LokiConfig)

	currentReport.QueryExecutors = append(currentReport.QueryExecutors, lokiQueryExecutor)

	gen.Run(true)
	currentReport.TestEnd = time.Now()

	fetchErr := currentReport.Fetch()
	require.NoError(t, fetchErr, "failed to fetch current report")

	// path, storeErr := currentReport.Store()
	// require.NoError(t, storeErr, "failed to store current report", path)

	previousReport := benchspy.StandardReport{
		BasicData: benchspy.BasicData{
			TestName:    t.Name(),
			CommitOrTag: "e7fc5826a572c09f8b93df3b9f674113372ce924",
		},
		LocalReportStorage: benchspy.LocalReportStorage{
			Directory: "test_performance_reports",
		},
	}
	loadErr := previousReport.Load()
	require.NoError(t, loadErr, "failed to load previous report")

	isComparableErrs := previousReport.IsComparable(currentReport)
	require.Empty(t, isComparableErrs, "reports were not comparable", isComparableErrs)
	require.NotEmpty(t, currentReport.QueryExecutors[0].Results()["vu_over_time"], "vu_over_time results were missing from current report")
	require.NotEmpty(t, previousReport.QueryExecutors[0].Results()["vu_over_time"], "vu_over_time results were missing from current report")
	require.Equal(t, len(currentReport.QueryExecutors[0].Results()["vu_over_time"]), len(previousReport.QueryExecutors[0].Results()["vu_over_time"]), "vu_over_time results are not the same length")

	// compare each result entry individually
	for i := range currentReport.QueryExecutors[0].Results()["vu_over_time"] {
		require.Equal(t, currentReport.QueryExecutors[0].Results()["vu_over_time"][i], previousReport.QueryExecutors[0].Results()["vu_over_time"][i], "vu_over_time results are not the same for given index")
	}

	//compare averages
	var currentSum float64
	for _, value := range currentReport.QueryExecutors[0].Results()["vu_over_time"] {
		asFloat, err := strconv.ParseFloat(value, 64)
		require.NoError(t, err, "failed to parse float")
		currentSum += asFloat
	}
	currentAverage := currentSum / float64(len(currentReport.QueryExecutors[0].Results()["vu_over_time"]))

	var previousSum float64
	for _, value := range previousReport.QueryExecutors[0].Results()["vu_over_time"] {
		asFloat, err := strconv.ParseFloat(value, 64)
		require.NoError(t, err, "failed to parse float")
		previousSum += asFloat
	}
	previousAverage := currentSum / float64(len(previousReport.QueryExecutors[0].Results()["vu_over_time"]))

	require.Equal(t, currentAverage, previousAverage, "vu_over_time averages are not the same")
}
