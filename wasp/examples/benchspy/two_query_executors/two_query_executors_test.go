package main

import (
	"context"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/prometheus/common/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
	"github.com/smartcontractkit/chainlink-testing-framework/wasp/benchspy"
)

func TestBenchSpy_Standard_Prometheus_And_Loki_Metrics(t *testing.T) {
	// this test requires CTFv2 node_set with observability stack to be running

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
		"91ee9e3c903d52de12f3d0c1a07ac3c2a6d141fb",
		benchspy.WithQueryExecutorType(benchspy.StandardQueryExecutor_Prometheus, benchspy.StandardQueryExecutor_Loki),
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
		"91ee9e3c903d52de12f3d0c1a07ac3c2a6d141fc",
		benchspy.WithQueryExecutorType(benchspy.StandardQueryExecutor_Prometheus, benchspy.StandardQueryExecutor_Loki),
		benchspy.WithPrometheusConfig(benchspy.NewPrometheusConfig("node[^0]")),
		benchspy.WithGenerators(newGen),
	)
	require.NoError(t, err, "failed to fetch current report or load the previous one")
	require.Equal(t, baseLineReport.CommitOrTag, previousReport.CommitOrTag, "current report should be the same as the original report")

	currentAsLokiSlices := benchspy.MustAllLokiResults(currentReport)
	previousAsLokiSlices := benchspy.MustAllLokiResults(previousReport)

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

var compareMedian = func(t *testing.T, metricName string, currentAsStringSlice, previousAsStringSlice map[string][]string) {
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
