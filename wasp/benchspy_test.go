package wasp_test

import (
	"context"
	"fmt"
	"math"

	"strconv"
	"testing"
	"time"

	// "github.com/prometheus/common/model"
	// "github.com/smartcontractkit/chainlink-testing-framework/wasp"
	"github.com/prometheus/common/model"
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
	gen.Run(true)

	basicData, err := benchspy.NewBasicData("e7fc5826a572c09f8b93df3b9f674113372ce925", gen)
	require.NoError(t, err)

	currentReport := benchspy.StandardReport{
		BasicData: *basicData,
	}

	lokiQueryExecutor := benchspy.NewLokiQueryExecutor(
		map[string]string{
			"vu_over_time": fmt.Sprintf("max_over_time({branch=~\"%s\", commit=~\"%s\", go_test_name=~\"%s\", test_data_type=~\"stats\", gen_name=~\"%s\"} | json | unwrap current_instances [10s]) by (node_id, go_test_name, gen_name)", label, label, t.Name(), gen.Cfg.GenName),
		},
		gen.Cfg.LokiConfig)

	currentReport.QueryExecutors = append(currentReport.QueryExecutors, lokiQueryExecutor)

	fetchCtx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()

	fetchErr := currentReport.FetchData(fetchCtx)
	require.NoError(t, fetchErr, "failed to fetch current report")

	path, storeErr := currentReport.Store()
	require.NoError(t, storeErr, "failed to store current report", path)

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

	currentAsStringSlice, castErr := benchspy.ResultsAs([]string{}, currentReport.QueryExecutors, benchspy.StandardQueryExecutor_Loki)
	require.NoError(t, castErr, "failed to cast results to string slice")
	require.NotEmpty(t, currentAsStringSlice, "results were empty")

	previousAsStringSlice, castErr := benchspy.ResultsAs([]string{}, previousReport.QueryExecutors, benchspy.StandardQueryExecutor_Loki)
	require.NoError(t, castErr, "failed to cast results to string slice")
	require.NotEmpty(t, previousAsStringSlice, "results were empty")

	require.NotEmpty(t, currentAsStringSlice["vu_over_time"], "vu_over_time results were missing from current report")
	require.NotEmpty(t, previousAsStringSlice["vu_over_time"], "vu_over_time results were missing from current report")
	require.Equal(t, len(currentAsStringSlice["vu_over_time"]), len(previousAsStringSlice["vu_over_time"]), "vu_over_time results are not the same length")

	// compare each result entry individually
	for i := range currentAsStringSlice["vu_over_time"] {
		require.Equal(t, currentAsStringSlice["vu_over_time"][i], previousAsStringSlice["vu_over_time"][i], "vu_over_time results are not the same for given index")
	}

	//compare averages
	var currentSum float64
	for _, value := range currentAsStringSlice["vu_over_time"] {
		asFloat, err := strconv.ParseFloat(value, 64)
		require.NoError(t, err, "failed to parse float")
		currentSum += asFloat
	}
	currentAverage := currentSum / float64(len(currentAsStringSlice["vu_over_time"]))

	var previousSum float64
	for _, value := range previousAsStringSlice["vu_over_time"] {
		asFloat, err := strconv.ParseFloat(value, 64)
		require.NoError(t, err, "failed to parse float")
		previousSum += asFloat
	}
	previousAverage := previousSum / float64(len(previousAsStringSlice["vu_over_time"]))

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

	gen.Run(true)

	currentReport := benchspy.StandardReport{
		BasicData: benchspy.MustNewBasicData("e7fc5826a572c09f8b93df3b9f674113372ce925", gen),
	}

	lokiQueryExecutor := benchspy.NewLokiQueryExecutor(
		map[string]string{
			"vu_over_time":        fmt.Sprintf("max_over_time({branch=~\"%s\", commit=~\"%s\", go_test_name=~\"%s\", test_data_type=~\"stats\", gen_name=~\"%s\"} | json | unwrap current_instances [10s]) by (node_id, go_test_name, gen_name)", label, label, t.Name(), gen.Cfg.GenName),
			"responses_over_time": fmt.Sprintf("sum(count_over_time({branch=~\"%s\", commit=~\"%s\", go_test_name=~\"%s\", test_data_type=~\"responses\", gen_name=~\"%s\"} [1s])) by (node_id, go_test_name, gen_name)", label, label, t.Name(), gen.Cfg.GenName),
		},
		gen.Cfg.LokiConfig)

	currentReport.QueryExecutors = append(currentReport.QueryExecutors, lokiQueryExecutor)

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

	currentAsStringSlice, castErr := benchspy.ResultsAs([]string{}, currentReport.QueryExecutors, benchspy.StandardQueryExecutor_Loki)
	require.NoError(t, castErr, "failed to cast results to string slice")
	require.NotEmpty(t, currentAsStringSlice, "results were empty")

	previousAsStringSlice, castErr := benchspy.ResultsAs([]string{}, previousReport.QueryExecutors, benchspy.StandardQueryExecutor_Loki)
	require.NoError(t, castErr, "failed to cast results to string slice")
	require.NotEmpty(t, previousAsStringSlice, "results were empty")

	// vu over time
	require.NotEmpty(t, currentReport.QueryExecutors[0].Results()["vu_over_time"], "vu_over_time results were missing from current report")
	require.NotEmpty(t, previousAsStringSlice["vu_over_time"], "vu_over_time results were missing from current report")
	require.Equal(t, len(currentAsStringSlice["vu_over_time"]), len(previousAsStringSlice["vu_over_time"]), "vu_over_time results are not the same length")

	// compare each vu_over_time entry individually
	for i := range currentAsStringSlice["vu_over_time"] {
		require.Equal(t, currentAsStringSlice["vu_over_time"][i], previousAsStringSlice["vu_over_time"][i], "vu_over_time results are not the same for given index")
	}

	//compare vu_over_time averages
	var currentSum float64
	for _, value := range currentAsStringSlice["vu_over_time"] {
		asFloat, err := strconv.ParseFloat(value, 64)
		require.NoError(t, err, "failed to parse float")
		currentSum += asFloat
	}
	currentAverage := currentSum / float64(len(currentAsStringSlice["vu_over_time"]))

	var previousSum float64
	for _, value := range previousAsStringSlice["vu_over_time"] {
		asFloat, err := strconv.ParseFloat(value, 64)
		require.NoError(t, err, "failed to parse float")
		previousSum += asFloat
	}
	previousAverage := previousSum / float64(len(previousAsStringSlice["vu_over_time"]))

	require.Equal(t, currentAverage, previousAverage, "vu_over_time averages are not the same")

	// responses over time
	require.NotEmpty(t, currentAsStringSlice["responses_over_time"], "responses_over_time results were missing from current report")
	require.NotEmpty(t, previousReport.QueryExecutors[0].Results()["responses_over_time"], "responses_over_time results were missing from current report")
	require.Equal(t, len(currentAsStringSlice["responses_over_time"]), len(previousAsStringSlice["responses_over_time"]), "responses_over_time results are not the same length")

	//compare responses_over_time averages
	var currentRespSum float64
	for _, value := range currentAsStringSlice["responses_over_time"] {
		asFloat, err := strconv.ParseFloat(value, 64)
		require.NoError(t, err, "failed to parse float")
		currentRespSum += asFloat
	}
	currentRespAverage := currentRespSum / float64(len(currentAsStringSlice["responses_over_time"]))

	var previousRespSum float64
	for _, value := range currentAsStringSlice["responses_over_time"] {
		asFloat, err := strconv.ParseFloat(value, 64)
		require.NoError(t, err, "failed to parse float")
		previousRespSum += asFloat
	}
	previousRespAverage := previousRespSum / float64(len(currentAsStringSlice["responses_over_time"]))

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

	currentReport, err := benchspy.NewStandardReport("e7fc5826a572c09f8b93df3b9f674113372ce925", benchspy.WithStandardQueryExecutorType(benchspy.StandardQueryExecutor_Loki), benchspy.WithGenerators(gen))
	require.NoError(t, err)

	fetchCtx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()

	fetchErr := currentReport.FetchData(fetchCtx)
	require.NoError(t, fetchErr, "failed to fetch current report")

	path, storeErr := currentReport.Store()
	require.NoError(t, storeErr, "failed to store current report", path)

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

	currentAsStringSlice, castErr := benchspy.ResultsAs([]string{}, currentReport.QueryExecutors, benchspy.StandardQueryExecutor_Loki)
	require.NoError(t, castErr, "failed to cast results to string slice")
	require.NotEmpty(t, currentAsStringSlice, "results were empty")

	previousAsStringSlice, castErr := benchspy.ResultsAs([]string{}, previousReport.QueryExecutors, benchspy.StandardQueryExecutor_Loki)
	require.NoError(t, castErr, "failed to cast results to string slice")
	require.NotEmpty(t, previousAsStringSlice, "results were empty")

	var compareMedian = func(metricName benchspy.StandardLoadMetric) {
		require.NotEmpty(t, currentAsStringSlice[string(metricName)], "%s results were missing from current report", string(metricName))
		require.NotEmpty(t, previousAsStringSlice[string(metricName)], "%s results were missing from previous report", string(metricName))
		require.Equal(t, len(currentAsStringSlice[string(metricName)]), len(previousAsStringSlice[string(metricName)]), "%s results are not the same length", string(metricName))

		currentFloatSlice, err := benchspy.StringSliceToFloat64Slice(currentAsStringSlice[string(metricName)])
		require.NoError(t, err, "failed to convert %s results to float64 slice", string(metricName))
		currentMedian := benchspy.CalculatePercentile(currentFloatSlice, 0.5)

		previousFloatSlice, err := benchspy.StringSliceToFloat64Slice(previousAsStringSlice[string(metricName)])
		require.NoError(t, err, "failed to convert %s results to float64 slice", string(metricName))
		previousMedian := benchspy.CalculatePercentile(previousFloatSlice, 0.5)

		var diffPrecentage float64
		if previousMedian != 0 {
			diffPrecentage = (currentMedian - previousMedian) / previousMedian * 100
		} else {
			diffPrecentage = currentMedian * 100
		}
		require.LessOrEqual(t, math.Abs(diffPrecentage), 1.0, "%s medians are more than 1% different", string(metricName), fmt.Sprintf("%.4f", diffPrecentage))
	}

	compareMedian(benchspy.MedianLatency)
	compareMedian(benchspy.Percentile95Latency)
	compareMedian(benchspy.ErrorRate)
}

func TestBenchSpyWithStandardGeneratorMetrics(t *testing.T) {
	gen, err := wasp.NewGenerator(&wasp.Config{
		T: t,
		// notice lack of Loki config
		GenName:     "vu",
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

	currentReport, err := benchspy.NewStandardReport("e7fc5826a572c09f8b93df3b9f674113372ce925", benchspy.WithStandardQueryExecutorType(benchspy.StandardQueryExecutor_Generator), benchspy.WithGenerators(gen))
	require.NoError(t, err)

	// context is not really needed, since we are using a generator, but it's required by the FetchData method
	fetchErr := currentReport.FetchData(context.Background())
	require.NoError(t, fetchErr, "failed to fetch current report")

	path, storeErr := currentReport.Store()
	require.NoError(t, storeErr, "failed to store current report", path)

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

	currentAsString, castErr := benchspy.ResultsAs("", currentReport.QueryExecutors, benchspy.StandardQueryExecutor_Generator)
	require.NoError(t, castErr, "failed to cast results to string slice")
	require.NotEmpty(t, currentAsString, "results were empty")

	previousAsString, castErr := benchspy.ResultsAs("", previousReport.QueryExecutors, benchspy.StandardQueryExecutor_Generator)
	require.NoError(t, castErr, "failed to cast results to string slice")
	require.NotEmpty(t, previousAsString, "results were empty")

	var compareValues = func(metricName benchspy.StandardLoadMetric) {
		require.NotEmpty(t, currentAsString[string(metricName)], "%s results were missing from current report", string(metricName))
		require.NotEmpty(t, previousAsString[string(metricName)], "%s results were missing from previous report", string(metricName))

		currentFloat, err := strconv.ParseFloat(currentAsString[string(metricName)], 64)
		require.NoError(t, err, "failed to convert %s results to float64 slice", string(metricName))

		previousFloat, err := strconv.ParseFloat(previousAsString[string(metricName)], 64)
		require.NoError(t, err, "failed to convert %s results to float64 slice", string(metricName))

		var diffPrecentage float64
		if previousFloat != 0 {
			diffPrecentage = (currentFloat - previousFloat) / previousFloat * 100
		} else {
			diffPrecentage = currentFloat * 100
		}
		require.LessOrEqual(t, math.Abs(diffPrecentage), 1.0, "%s medians are more than 1% different", string(metricName), fmt.Sprintf("%.4f", diffPrecentage))
	}

	compareValues(benchspy.MedianLatency)
	compareValues(benchspy.Percentile95Latency)
	compareValues(benchspy.ErrorRate)
}

func TestBenchSpy_Prometheus(t *testing.T) {
	t.Skip("skipping test, since it requires a running CTFv2 node_set and observability stack [start it manually and run the test]")

	before := time.Now().Add(-5 * time.Minute)
	// exclude bootstrap node
	prometheusNodeReporter, err := benchspy.NewStandardPrometheusQueryExecutor("http://localhost:9090", before, time.Now(), `node[^0]`)
	require.NoError(t, err)

	fetchErr := prometheusNodeReporter.Execute(context.Background())
	require.NoError(t, fetchErr, "failed to fetch prometheus node resources")

	resourcesAsValue := prometheusNodeReporter.MustResultsAsValue()
	medianCpuUsagePerNode := resourcesAsValue[string(benchspy.MedianCPUUsage)]
	require.Equal(t, medianCpuUsagePerNode.Type(), model.ValVector, "median cpu usage per node should be a vector")

	medianCpuUsagePerNodeVector := medianCpuUsagePerNode.(model.Vector)
	require.NotEmpty(t, medianCpuUsagePerNodeVector, "median cpu usage per node vector should not be empty")
	require.Equal(t, 4, len(medianCpuUsagePerNodeVector), "median cpu usage per node vector should have 4 elements")

	for _, sample := range medianCpuUsagePerNodeVector {
		require.NotZero(t, sample.Value, "median cpu usage per node should not be zero")
	}
}
