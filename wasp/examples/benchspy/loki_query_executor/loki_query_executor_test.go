package main

import (
	"context"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/montanaflynn/stats"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
	"github.com/smartcontractkit/chainlink-testing-framework/wasp/benchspy"
)

// this test requires CTFv2 observability stack to be running
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
		"v1",
		benchspy.WithStandardQueries(benchspy.StandardQueryExecutor_Loki),
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

	fetchCtx, cancelFn = context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()

	currentReport, previousReport, err := benchspy.FetchNewStandardReportAndLoadLatestPrevious(
		fetchCtx,
		"v2",
		benchspy.WithStandardQueries(benchspy.StandardQueryExecutor_Loki),
		benchspy.WithGenerators(newGen),
	)
	require.NoError(t, err, "failed to fetch current report or load the previous one")
	require.Equal(t, baseLineReport.CommitOrTag, previousReport.CommitOrTag, "current report should be the same as the original report")

	allCurrentAsStringSlice := benchspy.MustAllLokiResults(currentReport)
	allPreviousAsStringSlice := benchspy.MustAllLokiResults(previousReport)

	require.NotEmpty(t, allCurrentAsStringSlice, "current report is empty")
	require.NotEmpty(t, allPreviousAsStringSlice, "previous report is empty")

	currentAsStringSlice := allCurrentAsStringSlice[gen.Cfg.GenName]
	previousAsStringSlice := allPreviousAsStringSlice[gen.Cfg.GenName]

	compareAverages(t, string(benchspy.MedianLatency), currentAsStringSlice, previousAsStringSlice, 1.0)
	compareAverages(t, string(benchspy.Percentile95Latency), currentAsStringSlice, previousAsStringSlice, 1.0)
	compareAverages(t, string(benchspy.MaxLatency), currentAsStringSlice, previousAsStringSlice, 1.0)
	compareAverages(t, string(benchspy.ErrorRate), currentAsStringSlice, previousAsStringSlice, 1.0)
}

// this test requires CTFv2 observability stack to be running
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
		gen.Cfg.GenName,
		map[string]string{
			"vu_over_time":        fmt.Sprintf("max_over_time({branch=~\"%s\", commit=~\"%s\", go_test_name=~\"%s\", test_data_type=~\"stats\", gen_name=~\"%s\"} | json | unwrap current_instances [10s]) by (node_id, go_test_name, gen_name)", label, label, t.Name(), gen.Cfg.GenName),
			"responses_over_time": fmt.Sprintf("sum(count_over_time({branch=~\"%s\", commit=~\"%s\", go_test_name=~\"%s\", test_data_type=~\"responses\", gen_name=~\"%s\"} [1s])) by (node_id, go_test_name, gen_name)", label, label, t.Name(), gen.Cfg.GenName),
		},
		gen.Cfg.LokiConfig,
	)

	baseLineReport, err := benchspy.NewStandardReport(
		"v1",
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

	fetchCtx, cancelFn = context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()

	currentReport, previousReport, err := benchspy.FetchNewStandardReportAndLoadLatestPrevious(
		fetchCtx,
		"v2",
		benchspy.WithQueryExecutors(lokiQueryExecutor),
		benchspy.WithGenerators(newGen),
	)
	require.NoError(t, err, "failed to fetch current report or load the previous one")
	require.Equal(t, baseLineReport.CommitOrTag, previousReport.CommitOrTag, "current report should be the same as the original report")

	allCurrentAsStringSlice := benchspy.MustAllLokiResults(currentReport)
	allPreviousAsStringSlice := benchspy.MustAllLokiResults(previousReport)

	require.NotEmpty(t, allCurrentAsStringSlice, "current report is empty")
	require.NotEmpty(t, allPreviousAsStringSlice, "previous report is empty")

	currentAsStringSlice := allCurrentAsStringSlice[gen.Cfg.GenName]
	previousAsStringSlice := allPreviousAsStringSlice[gen.Cfg.GenName]

	compareAverages(t, "vu_over_time", currentAsStringSlice, previousAsStringSlice, 1.0)
	compareAverages(t, "responses_over_time", currentAsStringSlice, previousAsStringSlice, 1.0)
}

var compareAverages = func(t *testing.T, metricName string, currentAsStringSlice, previousAsStringSlice map[string][]string, maxPercentageDiff float64) {
	require.NotEmpty(t, currentAsStringSlice[metricName], "%s results were missing from current report", metricName)
	require.NotEmpty(t, previousAsStringSlice[metricName], "%s results were missing from previous report", metricName)

	currentFloatSlice, err := benchspy.StringSliceToFloat64Slice(currentAsStringSlice[metricName])
	require.NoError(t, err, "failed to convert %s results to float64 slice", metricName)
	currentMedian, err := stats.Mean(currentFloatSlice)
	require.NoError(t, err, "failed to calculate median for %s results", metricName)

	previousFloatSlice, err := benchspy.StringSliceToFloat64Slice(previousAsStringSlice[metricName])
	require.NoError(t, err, "failed to convert %s results to float64 slice", metricName)
	previousMedian, err := stats.Mean(previousFloatSlice)
	require.NoError(t, err, "failed to calculate median for %s results", metricName)

	var diffPercentage float64
	if previousMedian != 0.0 && currentMedian != 0.0 {
		diffPercentage = (currentMedian - previousMedian) / previousMedian * 100
	} else if previousMedian == 0.0 && currentMedian == 0.0 {
		diffPercentage = 0.0
	} else {
		diffPercentage = 100.0
	}
	assert.LessOrEqual(t, math.Abs(diffPercentage), maxPercentageDiff, "%s medians are more than 1% different", metricName, fmt.Sprintf("%.4f", diffPercentage))
}

func TestBenchSpy_Standard_Loki_Metrics_Two_Generators(t *testing.T) {
	label := "benchspy-std"

	p1 := wasp.NewProfile()
	p1.Add(wasp.NewGenerator(&wasp.Config{
		T:          t,
		LokiConfig: wasp.NewEnvLokiConfig(),
		GenName:    "vu1",
		Labels: map[string]string{
			"branch": label,
			"commit": label,
		},
		CallTimeout: 200 * time.Millisecond,
		LoadType:    wasp.VU,
		Schedule:    wasp.Plain(10, 12*time.Second),
		VU: wasp.NewMockVU(&wasp.MockVirtualUserConfig{
			CallSleep: 50 * time.Millisecond,
		}),
	}))

	p1.Add(wasp.NewGenerator(&wasp.Config{
		T:          t,
		LokiConfig: wasp.NewEnvLokiConfig(),
		GenName:    "vu2",
		Labels: map[string]string{
			"branch": label,
			"commit": label,
		},
		CallTimeout: 100 * time.Millisecond,
		LoadType:    wasp.VU,
		Schedule:    wasp.Plain(7, 15*time.Second),
		VU: wasp.NewMockVU(&wasp.MockVirtualUserConfig{
			CallSleep: 48 * time.Millisecond,
		}),
	}))

	_, err := p1.Run(true)
	require.NoError(t, err)

	baseLineReport, err := benchspy.NewStandardReport(
		"c2cf545d733eef8bad51d685fcb302e277d7ca14",
		benchspy.WithStandardQueries(benchspy.StandardQueryExecutor_Loki),
		benchspy.WithGenerators(p1.Generators[0], p1.Generators[1]),
	)
	require.NoError(t, err, "failed to create original report")

	fetchCtx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()

	fetchErr := baseLineReport.FetchData(fetchCtx)
	require.NoError(t, fetchErr, "failed to fetch data for original report")

	path, storeErr := baseLineReport.Store()
	require.NoError(t, storeErr, "failed to store current report", path)

	p2 := wasp.NewProfile()
	p2.Add(wasp.NewGenerator(&wasp.Config{
		T:          t,
		LokiConfig: wasp.NewEnvLokiConfig(),
		GenName:    "vu1",
		Labels: map[string]string{
			"branch": label,
			"commit": label,
		},
		CallTimeout: 200 * time.Millisecond,
		LoadType:    wasp.VU,
		Schedule:    wasp.Plain(10, 12*time.Second),
		VU: wasp.NewMockVU(&wasp.MockVirtualUserConfig{
			CallSleep: 50 * time.Millisecond,
		}),
	}))

	p2.Add(wasp.NewGenerator(&wasp.Config{
		T:          t,
		LokiConfig: wasp.NewEnvLokiConfig(),
		GenName:    "vu2",
		Labels: map[string]string{
			"branch": label,
			"commit": label,
		},
		CallTimeout: 100 * time.Millisecond,
		LoadType:    wasp.VU,
		Schedule:    wasp.Plain(7, 15*time.Second),
		VU: wasp.NewMockVU(&wasp.MockVirtualUserConfig{
			CallSleep: 48 * time.Millisecond,
		}),
	}))

	_, err = p2.Run(true)
	require.NoError(t, err)

	fetchCtx, cancelFn = context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()

	currentReport, previousReport, err := benchspy.FetchNewStandardReportAndLoadLatestPrevious(
		fetchCtx,
		"c2cf545d733eef8bad51d685fcb302e277d7ca15",
		benchspy.WithStandardQueries(benchspy.StandardQueryExecutor_Loki),
		benchspy.WithGenerators(p2.Generators[0], p2.Generators[1]),
	)
	require.NoError(t, err, "failed to fetch current report or load the previous one")
	require.Equal(t, baseLineReport.CommitOrTag, previousReport.CommitOrTag, "current report should be the same as the original report")

	allCurrentAsStringSlice := benchspy.MustAllLokiResults(currentReport)
	allPreviousAsStringSlice := benchspy.MustAllLokiResults(previousReport)

	require.Equal(t, 2, len(allCurrentAsStringSlice), "current report doesn't have 2 loki results")
	require.Equal(t, 2, len(allPreviousAsStringSlice), "previous report doesn't have 2 loki results")

	currentAsStringSlice_vu1 := allCurrentAsStringSlice[p1.Generators[0].Cfg.GenName]
	previousAsStringSlice_vu1 := allPreviousAsStringSlice[p2.Generators[0].Cfg.GenName]

	compareAverages(t, string(benchspy.MedianLatency), currentAsStringSlice_vu1, previousAsStringSlice_vu1, 10.0)
	compareAverages(t, string(benchspy.Percentile95Latency), currentAsStringSlice_vu1, previousAsStringSlice_vu1, 10.0)
	compareAverages(t, string(benchspy.MaxLatency), currentAsStringSlice_vu1, previousAsStringSlice_vu1, 10.0)
	compareAverages(t, string(benchspy.ErrorRate), currentAsStringSlice_vu1, previousAsStringSlice_vu1, 10.0)

	currentAsStringSlice_vu2 := allCurrentAsStringSlice[p1.Generators[1].Cfg.GenName]
	previousAsStringSlice_vu2 := allPreviousAsStringSlice[p2.Generators[1].Cfg.GenName]

	compareAverages(t, string(benchspy.MedianLatency), currentAsStringSlice_vu2, previousAsStringSlice_vu2, 10.0)
	compareAverages(t, string(benchspy.Percentile95Latency), currentAsStringSlice_vu2, previousAsStringSlice_vu2, 10.0)
	compareAverages(t, string(benchspy.MaxLatency), currentAsStringSlice_vu2, previousAsStringSlice_vu2, 10.0)
	compareAverages(t, string(benchspy.ErrorRate), currentAsStringSlice_vu2, previousAsStringSlice_vu2, 10.0)
}
