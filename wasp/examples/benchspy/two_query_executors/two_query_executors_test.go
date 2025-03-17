package main

import (
	"context"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/montanaflynn/stats"
	"github.com/prometheus/common/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
	"github.com/smartcontractkit/chainlink-testing-framework/wasp/benchspy"
)

// this test requires CTFv2 node_set with observability stack to be running
func TestBenchSpy_Standard_Prometheus_And_Loki_Metrics(t *testing.T) {
	label := "benchspy-two-query-executors"

	gen, err := wasp.NewGenerator(&wasp.Config{
		T:           t,
		GenName:     "vu",
		CallTimeout: 100 * time.Millisecond,
		LoadType:    wasp.VU,
		Schedule:    wasp.Plain(1, 10*time.Second),
		VU: wasp.NewMockVU(&wasp.MockVirtualUserConfig{
			CallSleep: 50 * time.Millisecond,
		}),
		Labels: map[string]string{
			"branch": label,
			"commit": label,
		},
		LokiConfig: wasp.NewEnvLokiConfig(),
	})
	require.NoError(t, err)

	gen.Run(true)

	baseLineReport, err := benchspy.NewStandardReport(
		"v1.0.0",
		benchspy.WithStandardQueries(benchspy.StandardQueryExecutor_Prometheus, benchspy.StandardQueryExecutor_Loki),
		benchspy.WithPrometheusConfig(benchspy.NewPrometheusConfig("node[^0]")),
		benchspy.WithGenerators(gen),
	)
	require.NoError(t, err, "failed to create original report")

	fetchCtx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()

	fetchErr := baseLineReport.FetchData(fetchCtx)
	require.NoError(t, fetchErr, "failed to fetch current report")

	path, storeErr := baseLineReport.Store()
	require.NoError(t, storeErr, "failed to store current report", path)

	newGen, err := wasp.NewGenerator(&wasp.Config{
		T:           t,
		GenName:     "vu",
		CallTimeout: 100 * time.Millisecond,
		LoadType:    wasp.VU,
		Schedule:    wasp.Plain(1, 10*time.Second),
		VU: wasp.NewMockVU(&wasp.MockVirtualUserConfig{
			CallSleep: 50 * time.Millisecond,
		}),
		Labels: map[string]string{
			"branch": label,
			"commit": label,
		},
		LokiConfig: wasp.NewEnvLokiConfig(),
	})
	require.NoError(t, err)

	newGen.Run(true)

	fetchCtx, cancelFn = context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()

	currentReport, previousReport, err := benchspy.FetchNewStandardReportAndLoadLatestPrevious(
		fetchCtx,
		"v1.1.0",
		benchspy.WithStandardQueries(benchspy.StandardQueryExecutor_Prometheus, benchspy.StandardQueryExecutor_Loki),
		benchspy.WithPrometheusConfig(benchspy.NewPrometheusConfig("node[^0]")),
		benchspy.WithGenerators(newGen),
	)
	require.NoError(t, err, "failed to fetch current report or load the previous one")
	require.Equal(t, baseLineReport.CommitOrTag, previousReport.CommitOrTag, "current report should be the same as the original report")

	allCurrentAsStringSlice := benchspy.MustAllLokiResults(currentReport)
	allPreviousAsStringSlice := benchspy.MustAllLokiResults(previousReport)

	require.NotEmpty(t, allCurrentAsStringSlice, "current report is empty")
	require.NotEmpty(t, allPreviousAsStringSlice, "previous report is empty")

	currentAsLokiSlices := allCurrentAsStringSlice[gen.Cfg.GenName]
	previousAsLokiSlices := allPreviousAsStringSlice[gen.Cfg.GenName]

	compareMedian(t, string(benchspy.MedianLatency), currentAsLokiSlices, previousAsLokiSlices)

	currentPromValues := benchspy.MustAllPrometheusResults(currentReport)
	previousPromValues := benchspy.MustAllPrometheusResults(previousReport)

	assert.Equal(t, len(currentPromValues), len(previousPromValues), "number of metrics in results should be the same")

	currentMedianCPUUsage := currentPromValues[string(benchspy.MedianCPUUsage)]
	previousMedianCPUUsage := previousPromValues[string(benchspy.MedianCPUUsage)]

	assert.Equal(t, currentMedianCPUUsage.Type(), previousMedianCPUUsage.Type(), "types of metrics should be the same")

	currentMedianCPUUsageVector := currentMedianCPUUsage.(model.Vector)
	previousMedianCPUUsageVector := previousMedianCPUUsage.(model.Vector)

	assert.Equal(t, len(currentMedianCPUUsageVector), len(previousMedianCPUUsageVector), "number of samples in vectors should be the same")

	// here we could compare actual values, but most likely they will be very different and the test will fail
}

// this test requires CTFv2 observability stack to be running
func TestBenchSpy_Two_Loki_Executors(t *testing.T) {
	label := "benchspy-two-loki-executors"

	gen, err := wasp.NewGenerator(&wasp.Config{
		T:           t,
		GenName:     "vu",
		CallTimeout: 100 * time.Millisecond,
		LoadType:    wasp.VU,
		Schedule:    wasp.Plain(1, 10*time.Second),
		VU: wasp.NewMockVU(&wasp.MockVirtualUserConfig{
			CallSleep: 50 * time.Millisecond,
		}),
		Labels: map[string]string{
			"branch": label,
			"commit": label,
		},
		LokiConfig: wasp.NewEnvLokiConfig(),
	})
	require.NoError(t, err)

	gen.Run(true)

	firstLokiQueryExecutor := benchspy.NewLokiQueryExecutor(
		gen.Cfg.GenName,
		map[string]string{
			"vu_over_time": fmt.Sprintf("max_over_time({branch=~\"%s\", commit=~\"%s\", go_test_name=~\"%s\", test_data_type=~\"stats\", gen_name=~\"%s\"} | json | unwrap current_instances [10s]) by (node_id, go_test_name, gen_name)", label, label, t.Name(), gen.Cfg.GenName),
		},
		gen.Cfg.LokiConfig,
	)

	secondLokiQueryExecutor := benchspy.NewLokiQueryExecutor(
		gen.Cfg.GenName,
		map[string]string{
			"responses_over_time": fmt.Sprintf("sum(count_over_time({branch=~\"%s\", commit=~\"%s\", go_test_name=~\"%s\", test_data_type=~\"responses\", gen_name=~\"%s\"} [1s])) by (node_id, go_test_name, gen_name)", label, label, t.Name(), gen.Cfg.GenName),
		},
		gen.Cfg.LokiConfig,
	)

	baseLineReport, err := benchspy.NewStandardReport(
		"v1.0.0",
		benchspy.WithQueryExecutors(firstLokiQueryExecutor, secondLokiQueryExecutor),
		benchspy.WithGenerators(gen),
	)
	require.NoError(t, err, "failed to create original report")

	fetchCtx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()

	fetchErr := baseLineReport.FetchData(fetchCtx)
	require.NoError(t, fetchErr, "failed to fetch current report")

	path, storeErr := baseLineReport.Store()
	require.NoError(t, storeErr, "failed to store current report", path)

	newGen, err := wasp.NewGenerator(&wasp.Config{
		T:           t,
		GenName:     "vu",
		CallTimeout: 100 * time.Millisecond,
		LoadType:    wasp.VU,
		Schedule:    wasp.Plain(1, 10*time.Second),
		VU: wasp.NewMockVU(&wasp.MockVirtualUserConfig{
			CallSleep: 50 * time.Millisecond,
		}),
		Labels: map[string]string{
			"branch": label,
			"commit": label,
		},
		LokiConfig: wasp.NewEnvLokiConfig(),
	})
	require.NoError(t, err)

	newGen.Run(true)

	fetchCtx, cancelFn = context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()

	currentReport, previousReport, err := benchspy.FetchNewStandardReportAndLoadLatestPrevious(
		fetchCtx,
		"v1.1.0",
		benchspy.WithQueryExecutors(firstLokiQueryExecutor, secondLokiQueryExecutor),
		benchspy.WithGenerators(newGen),
	)
	require.NoError(t, err, "failed to fetch current report or load the previous one")
	require.Equal(t, baseLineReport.CommitOrTag, previousReport.CommitOrTag, "current report should be the same as the original report")

	allCurrentAsStringSlice := benchspy.MustAllLokiResults(currentReport)
	allPreviousAsStringSlice := benchspy.MustAllLokiResults(previousReport)

	require.NotEmpty(t, allCurrentAsStringSlice, "current report is empty")
	require.NotEmpty(t, allPreviousAsStringSlice, "previous report is empty")

	currentAsLokiSlices := allCurrentAsStringSlice[gen.Cfg.GenName]
	previousAsLokiSlices := allPreviousAsStringSlice[gen.Cfg.GenName]

	compareMedian(t, "vu_over_time", currentAsLokiSlices, previousAsLokiSlices)
	compareMedian(t, "responses_over_time", currentAsLokiSlices, previousAsLokiSlices)

	// here we could compare actual values, but most likely they will be very different and the test will fail
}

var compareMedian = func(t *testing.T, metricName string, currentAsStringSlice, previousAsStringSlice map[string][]string) {
	require.NotEmpty(t, currentAsStringSlice[metricName], "%s results were missing from current report", metricName)
	require.NotEmpty(t, previousAsStringSlice[metricName], "%s results were missing from previous report", metricName)

	currentFloatSlice, err := benchspy.StringSliceToFloat64Slice(currentAsStringSlice[metricName])
	require.NoError(t, err, "failed to convert %s results to float64 slice", metricName)
	currentMedian, err := stats.Median(currentFloatSlice)
	require.NoError(t, err, "failed to calculate median for %s results", metricName)

	previousFloatSlice, err := benchspy.StringSliceToFloat64Slice(previousAsStringSlice[metricName])
	require.NoError(t, err, "failed to convert %s results to float64 slice", metricName)
	previousMedian, err := stats.Median(previousFloatSlice)
	require.NoError(t, err, "failed to calculate median for %s results", metricName)

	var diffPercentage float64
	if previousMedian != 0.0 && currentMedian != 0.0 {
		diffPercentage = (currentMedian - previousMedian) / previousMedian * 100
	} else if previousMedian == 0.0 && currentMedian == 0.0 {
		diffPercentage = 0.0
	} else {
		diffPercentage = 100.0
	}
	assert.LessOrEqual(t, math.Abs(diffPercentage), 1.0, "%s medians are more than 1% different", metricName, fmt.Sprintf("%.4f", diffPercentage))
}
