package wasp_test

import (
	"context"
	"fmt"
	"math"
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

	lokiQueryExecutor := benchspy.NewLokiQueryExecutor(
		map[string]string{
			"vu_over_time": fmt.Sprintf("max_over_time({branch=~\"%s\", commit=~\"%s\", go_test_name=~\"%s\", test_data_type=~\"stats\", gen_name=~\"%s\"} | json | unwrap current_instances [10s]) by (node_id, go_test_name, gen_name)", label, label, t.Name(), gen.Cfg.GenName),
		},
		gen.Cfg.LokiConfig)

	currentReport.QueryExecutors = append(currentReport.QueryExecutors, lokiQueryExecutor)

	gen.Run(true)

	fetchCtx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()

	fetchErr := currentReport.FetchData(fetchCtx)
	require.NoError(t, fetchErr, "failed to fetch current report")

	// path, storeErr := currentReport.Store()
	// require.NoError(t, storeErr, "failed to store current report", path)

	// this is only needed, because we are using a non-standard directory
	// otherwise, the Load method would be able to find the file
	previousReport := benchspy.StandardReport{
		LocalStorage: benchspy.LocalStorage{
			Directory: "test_performance_reports",
		},
	}
	loadErr := previousReport.Load(t.Name(), "e7fc5826a572c09f8b93df3b9f674113372ce924")
	require.NoError(t, loadErr, "failed to load previous report")

	isComparableErrs := previousReport.IsComparable(&currentReport)
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
	previousAverage := previousSum / float64(len(previousReport.QueryExecutors[0].Results()["vu_over_time"]))

	require.Equal(t, currentAverage, previousAverage, "vu_over_time averages are not the same")
}

func TestBenchSpyWithTwoLokiQueries(t *testing.T) {
	label := "benchspy2"

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
		BasicData: benchspy.MustNewBasicData("e7fc5826a572c09f8b93df3b9f674113372ce924", gen),
		ResourceReporter: benchspy.ResourceReporter{
			ExecutionEnvironment: benchspy.ExecutionEnvironment_Docker,
		},
	}

	lokiQueryExecutor := benchspy.NewLokiQueryExecutor(
		map[string]string{
			"vu_over_time":        fmt.Sprintf("max_over_time({branch=~\"%s\", commit=~\"%s\", go_test_name=~\"%s\", test_data_type=~\"stats\", gen_name=~\"%s\"} | json | unwrap current_instances [10s]) by (node_id, go_test_name, gen_name)", label, label, t.Name(), gen.Cfg.GenName),
			"responses_over_time": fmt.Sprintf("sum(count_over_time({branch=~\"%s\", commit=~\"%s\", go_test_name=~\"%s\", test_data_type=~\"responses\", gen_name=~\"%s\"} [1s])) by (node_id, go_test_name, gen_name)", label, label, t.Name(), gen.Cfg.GenName),
		},
		gen.Cfg.LokiConfig)

	currentReport.QueryExecutors = append(currentReport.QueryExecutors, lokiQueryExecutor)

	gen.Run(true)

	fetchCtx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()

	fetchErr := currentReport.FetchData(fetchCtx)
	require.NoError(t, fetchErr, "failed to fetch current report")

	// path, storeErr := currentReport.Store()
	// require.NoError(t, storeErr, "failed to store current report", path)

	// this is only needed, because we are using a non-standard directory
	// otherwise, the Load method would be able to find the file
	previousReport := benchspy.StandardReport{
		LocalStorage: benchspy.LocalStorage{
			Directory: "test_performance_reports",
		},
	}
	loadErr := previousReport.Load(t.Name(), "e7fc5826a572c09f8b93df3b9f674113372ce924")
	require.NoError(t, loadErr, "failed to load previous report")

	isComparableErrs := previousReport.IsComparable(&currentReport)
	require.Empty(t, isComparableErrs, "reports were not comparable", isComparableErrs)
	// vu over time
	require.NotEmpty(t, currentReport.QueryExecutors[0].Results()["vu_over_time"], "vu_over_time results were missing from current report")
	require.NotEmpty(t, previousReport.QueryExecutors[0].Results()["vu_over_time"], "vu_over_time results were missing from current report")
	require.Equal(t, len(currentReport.QueryExecutors[0].Results()["vu_over_time"]), len(previousReport.QueryExecutors[0].Results()["vu_over_time"]), "vu_over_time results are not the same length")

	// compare each vu_over_time entry individually
	for i := range currentReport.QueryExecutors[0].Results()["vu_over_time"] {
		require.Equal(t, currentReport.QueryExecutors[0].Results()["vu_over_time"][i], previousReport.QueryExecutors[0].Results()["vu_over_time"][i], "vu_over_time results are not the same for given index")
	}

	//compare vu_over_time averages
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
	previousAverage := previousSum / float64(len(previousReport.QueryExecutors[0].Results()["vu_over_time"]))

	require.Equal(t, currentAverage, previousAverage, "vu_over_time averages are not the same")

	// responses over time
	require.NotEmpty(t, currentReport.QueryExecutors[0].Results()["responses_over_time"], "responses_over_time results were missing from current report")
	require.NotEmpty(t, previousReport.QueryExecutors[0].Results()["responses_over_time"], "responses_over_time results were missing from current report")
	require.Equal(t, len(currentReport.QueryExecutors[0].Results()["responses_over_time"]), len(previousReport.QueryExecutors[0].Results()["responses_over_time"]), "responses_over_time results are not the same length")

	//compare responses_over_time averages
	var currentRespSum float64
	for _, value := range currentReport.QueryExecutors[0].Results()["responses_over_time"] {
		asFloat, err := strconv.ParseFloat(value, 64)
		require.NoError(t, err, "failed to parse float")
		currentRespSum += asFloat
	}
	currentRespAverage := currentRespSum / float64(len(currentReport.QueryExecutors[0].Results()["responses_over_time"]))

	var previousRespSum float64
	for _, value := range previousReport.QueryExecutors[0].Results()["responses_over_time"] {
		asFloat, err := strconv.ParseFloat(value, 64)
		require.NoError(t, err, "failed to parse float")
		previousRespSum += asFloat
	}
	previousRespAverage := previousRespSum / float64(len(previousReport.QueryExecutors[0].Results()["responses_over_time"]))

	diffPrecentage := (currentRespAverage - previousRespAverage) / previousRespAverage * 100
	require.LessOrEqual(t, math.Abs(diffPrecentage), 1.0, "responses_over_time averages are more than 1% different", fmt.Sprintf("%.4f", diffPrecentage))
}

func TestBenchSpyWithStandardLokiMetrics(t *testing.T) {
	label := "benchspy-std"

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

	gen.Run(true)

	currentReport, err := benchspy.NewStandardReport("e7fc5826a572c09f8b93df3b9f674113372ce925", benchspy.ExecutionEnvironment_Docker, gen)
	require.NoError(t, err)

	fetchCtx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()

	fetchErr := currentReport.FetchData(fetchCtx)
	require.NoError(t, fetchErr, "failed to fetch current report")

	// path, storeErr := currentReport.Store()
	// require.NoError(t, storeErr, "failed to store current report", path)

	// this is only needed, because we are using a non-standard directory
	// otherwise, the Load method would be able to find the file
	previousReport := benchspy.StandardReport{
		LocalStorage: benchspy.LocalStorage{
			Directory: "test_performance_reports",
		},
	}
	loadErr := previousReport.Load(t.Name(), "e7fc5826a572c09f8b93df3b9f674113372ce924")
	require.NoError(t, loadErr, "failed to load previous report")

	isComparableErrs := previousReport.IsComparable(currentReport)
	require.Empty(t, isComparableErrs, "reports were not comparable", isComparableErrs)

	var compareAverages = func(metricName benchspy.StandardMetric) {
		require.NotEmpty(t, currentReport.QueryExecutors[0].Results()[string(metricName)], "%s results were missing from current report", string(metricName))
		require.NotEmpty(t, previousReport.QueryExecutors[0].Results()[string(metricName)], "%s results were missing from previous report", string(metricName))
		require.Equal(t, len(currentReport.QueryExecutors[0].Results()[string(metricName)]), len(previousReport.QueryExecutors[0].Results()[string(metricName)]), "%s results are not the same length", string(metricName))

		var currentAvgSum float64
		for _, value := range currentReport.QueryExecutors[0].Results()[string(metricName)] {
			asFloat, err := strconv.ParseFloat(value, 64)
			require.NoError(t, err, "failed to parse float")
			currentAvgSum += asFloat
		}
		currentAvgAverage := currentAvgSum / float64(len(currentReport.QueryExecutors[0].Results()[string(metricName)]))

		var previousAvgSum float64
		for _, value := range previousReport.QueryExecutors[0].Results()[string(metricName)] {
			asFloat, err := strconv.ParseFloat(value, 64)
			require.NoError(t, err, "failed to parse float")
			previousAvgSum += asFloat
		}
		previousAvgAverage := previousAvgSum / float64(len(previousReport.QueryExecutors[0].Results()[string(metricName)]))

		var diffPrecentage float64
		if previousAvgAverage != 0 {
			diffPrecentage = (currentAvgAverage - previousAvgAverage) / previousAvgAverage * 100
		} else {
			diffPrecentage = currentAvgAverage * 100
		}
		require.LessOrEqual(t, math.Abs(diffPrecentage), 1.0, "%s averages are more than 1% different", string(metricName), fmt.Sprintf("%.4f", diffPrecentage))
	}

	compareAverages(benchspy.AverageLatency)
	compareAverages(benchspy.Percentile95Latency)
	compareAverages(benchspy.ErrorRate)
}
