package main

import (
	"context"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
	"github.com/smartcontractkit/chainlink-testing-framework/wasp/benchspy"
)

func TestBenchSpy_Standard_Loki_Metrics(t *testing.T) {
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
		Schedule:    wasp.Plain(10, 15*time.Second),
		VU: wasp.NewMockVU(&wasp.MockVirtualUserConfig{
			CallSleep: 50 * time.Millisecond,
		}),
	})
	require.NoError(t, err)

	gen.Run(true)

	baseLineReport, err := benchspy.NewStandardReport(
		"c2cf545d733eef8bad51d685fcb302e277d7ca14",
		benchspy.WithQueryExecutorType(benchspy.StandardQueryExecutor_Loki),
		benchspy.WithGenerators(gen),
	)
	require.NoError(t, err, "failed to create original report")

	fetchCtx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()

	fetchErr := baseLineReport.FetchData(fetchCtx)
	require.NoError(t, fetchErr, "failed to fetch data for original report")

	path, storeErr := baseLineReport.Store()
	require.NoError(t, storeErr, "failed to store current report", path)

	newGen, err := wasp.NewGenerator(&wasp.Config{
		T:          t,
		LokiConfig: wasp.NewEnvLokiConfig(),
		GenName:    "vu",
		Labels: map[string]string{
			"branch": label,
			"commit": label,
		},
		CallTimeout: 100 * time.Millisecond,
		LoadType:    wasp.VU,
		Schedule:    wasp.Plain(10, 15*time.Second),
		VU: wasp.NewMockVU(&wasp.MockVirtualUserConfig{
			CallSleep: 50 * time.Millisecond,
		}),
	})
	require.NoError(t, err)

	newGen.Run(true)

	currentReport, previousReport, err := benchspy.FetchNewStandardReportAndLoadLatestPrevious(
		fetchCtx,
		"c2cf545d733eef8bad51d685fcb302e277d7ca15",
		benchspy.WithQueryExecutorType(benchspy.StandardQueryExecutor_Loki),
		benchspy.WithGenerators(newGen),
	)
	require.NoError(t, err, "failed to fetch current report or load the previous one")
	require.Equal(t, baseLineReport.CommitOrTag, previousReport.CommitOrTag, "current report should be the same as the original report")

	currentAsStringSlice := benchspy.MustAllLokiResults(currentReport)
	previousAsStringSlice := benchspy.MustAllLokiResults(previousReport)

	var compareMedian = func(metricName string) {
		require.NotEmpty(t, currentAsStringSlice[metricName], "%s results were missing from current report", metricName)
		require.NotEmpty(t, previousAsStringSlice[metricName], "%s results were missing from previous report", metricName)

		currentFloatSlice, err := benchspy.StringSliceToFloat64Slice(currentAsStringSlice[metricName])
		require.NoError(t, err, "failed to convert %s results to float64 slice", metricName)
		currentMedian := benchspy.CalculatePercentile(currentFloatSlice, 0.5)

		previousFloatSlice, err := benchspy.StringSliceToFloat64Slice(previousAsStringSlice[metricName])
		require.NoError(t, err, "failed to convert %s results to float64 slice", metricName)
		previousMedian := benchspy.CalculatePercentile(previousFloatSlice, 0.5)

		var diffPrecentage float64
		if previousMedian != 0 {
			diffPrecentage = (currentMedian - previousMedian) / previousMedian * 100
		} else {
			diffPrecentage = currentMedian * 100
		}
		assert.LessOrEqual(t, math.Abs(diffPrecentage), 1.0, "%s medians are more than 1% different", metricName, fmt.Sprintf("%.4f", diffPrecentage))
	}

	compareMedian(string(benchspy.MedianLatency))
	compareMedian(string(benchspy.Percentile95Latency))
	compareMedian(string(benchspy.ErrorRate))
}

func TestBenchSpy_Custom_Loki_Metrics(t *testing.T) {
	label := "benchspy-custom"

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
		Schedule:    wasp.Plain(10, 15*time.Second),
		VU: wasp.NewMockVU(&wasp.MockVirtualUserConfig{
			CallSleep: 50 * time.Millisecond,
		}),
	})
	require.NoError(t, err)

	gen.Run(true)

	lokiQueryExecutor := benchspy.NewLokiQueryExecutor(
		map[string]string{
			"vu_over_time":        fmt.Sprintf("max_over_time({branch=~\"%s\", commit=~\"%s\", go_test_name=~\"%s\", test_data_type=~\"stats\", gen_name=~\"%s\"} | json | unwrap current_instances [10s]) by (node_id, go_test_name, gen_name)", label, label, t.Name(), gen.Cfg.GenName),
			"responses_over_time": fmt.Sprintf("sum(count_over_time({branch=~\"%s\", commit=~\"%s\", go_test_name=~\"%s\", test_data_type=~\"responses\", gen_name=~\"%s\"} [1s])) by (node_id, go_test_name, gen_name)", label, label, t.Name(), gen.Cfg.GenName),
		},
		gen.Cfg.LokiConfig,
	)

	baseLineReport, err := benchspy.NewStandardReport(
		"2d1fa3532656c51991c0212afce5f80d2914e34e",
		benchspy.WithQueryExecutors(lokiQueryExecutor),
		benchspy.WithGenerators(gen),
	)
	require.NoError(t, err, "failed to create baseline report")

	fetchCtx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()

	fetchErr := baseLineReport.FetchData(fetchCtx)
	require.NoError(t, fetchErr, "failed to fetch data for original report")

	path, storeErr := baseLineReport.Store()
	require.NoError(t, storeErr, "failed to store current report", path)

	newGen, err := wasp.NewGenerator(&wasp.Config{
		T:          t,
		LokiConfig: wasp.NewEnvLokiConfig(),
		GenName:    "vu",
		Labels: map[string]string{
			"branch": label,
			"commit": label,
		},
		CallTimeout: 100 * time.Millisecond,
		LoadType:    wasp.VU,
		Schedule:    wasp.Plain(10, 15*time.Second),
		VU: wasp.NewMockVU(&wasp.MockVirtualUserConfig{
			CallSleep: 50 * time.Millisecond,
		}),
	})
	require.NoError(t, err)

	newGen.Run(true)

	currentReport, previousReport, err := benchspy.FetchNewStandardReportAndLoadLatestPrevious(
		fetchCtx,
		"2d1fa3532656c51991c0212afce5f80d2914e34f",
		benchspy.WithQueryExecutors(lokiQueryExecutor),
		benchspy.WithGenerators(newGen),
	)
	require.NoError(t, err, "failed to fetch current report or load the previous one")
	require.Equal(t, baseLineReport.CommitOrTag, previousReport.CommitOrTag, "current report should be the same as the original report")

	currentAsStringSlice := benchspy.MustAllLokiResults(currentReport)
	previousAsStringSlice := benchspy.MustAllLokiResults(previousReport)

	var compareMedian = func(metricName string) {
		require.NotEmpty(t, currentAsStringSlice[metricName], "%s results were missing from current report", metricName)
		require.NotEmpty(t, previousAsStringSlice[metricName], "%s results were missing from previous report", metricName)

		currentFloatSlice, err := benchspy.StringSliceToFloat64Slice(currentAsStringSlice[metricName])
		require.NoError(t, err, "failed to convert %s results to float64 slice", metricName)
		currentMedian := benchspy.CalculatePercentile(currentFloatSlice, 0.5)

		previousFloatSlice, err := benchspy.StringSliceToFloat64Slice(previousAsStringSlice[metricName])
		require.NoError(t, err, "failed to convert %s results to float64 slice", metricName)
		previousMedian := benchspy.CalculatePercentile(previousFloatSlice, 0.5)

		var diffPrecentage float64
		if previousMedian != 0 {
			diffPrecentage = (currentMedian - previousMedian) / previousMedian * 100
		} else {
			diffPrecentage = currentMedian * 100
		}
		assert.LessOrEqual(t, math.Abs(diffPrecentage), 1.0, "%s medians are more than 1% different", metricName, fmt.Sprintf("%.4f", diffPrecentage))
	}

	compareMedian("vu_over_time")
	compareMedian("responses_over_time")
}
