package main

import (
	"context"
	"testing"
	"time"

	"github.com/prometheus/common/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
	"github.com/smartcontractkit/chainlink-testing-framework/wasp/benchspy"
)

// this test requires CTFv2 node_set with observability stack to be running
func TestBenchSpy_Standard_Prometheus_Metrics(t *testing.T) {
	gen, err := wasp.NewGenerator(&wasp.Config{
		T:           t,
		GenName:     "vu",
		CallTimeout: 100 * time.Millisecond,
		LoadType:    wasp.VU,
		Schedule:    wasp.Plain(1, 10*time.Second),
		VU: wasp.NewMockVU(&wasp.MockVirtualUserConfig{
			CallSleep: 50 * time.Millisecond,
		}),
	})
	require.NoError(t, err)

	gen.Run(true)

	// exclude bootstrap node
	promConfig := benchspy.NewPrometheusConfig("node[^0]")

	baseLineReport, err := benchspy.NewStandardReport(
		"91ee9e3c903d52de12f3d0c1a07ac3c2a6d141fb",
		benchspy.WithStandardQueries(benchspy.StandardQueryExecutor_Prometheus),
		benchspy.WithPrometheusConfig(promConfig),
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
	})
	require.NoError(t, err)

	newGen.Run(true)

	fetchCtx, cancelFn = context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()

	currentReport, previousReport, err := benchspy.FetchNewStandardReportAndLoadLatestPrevious(
		fetchCtx,
		"2d1fa3532656c51991c0212afce5f80d2914e34f",
		benchspy.WithStandardQueries(benchspy.StandardQueryExecutor_Prometheus),
		benchspy.WithPrometheusConfig(promConfig),
		benchspy.WithGenerators(newGen),
	)
	require.NoError(t, err, "failed to fetch current report or load the previous one")
	require.Equal(t, baseLineReport.CommitOrTag, previousReport.CommitOrTag, "current report should be the same as the original report")

	currentAsValues := benchspy.MustAllPrometheusResults(currentReport)
	previousAsValues := benchspy.MustAllPrometheusResults(previousReport)

	assert.Equal(t, len(currentAsValues), len(previousAsValues), "number of metrics in results should be the same")
	assert.Equal(t, 6, len(currentAsValues), "there should be 6 metrics in the report")

	for _, metric := range benchspy.StandardResourceMetrics {
		assert.NotEmpty(t, currentAsValues[string(metric)], "current report should contain metric %s", metric)
		assert.NotEmpty(t, previousAsValues[string(metric)], "previous report should contain metric %s", metric)
	}

	currentMedianCPUUsage := currentAsValues[string(benchspy.MedianCPUUsage)]
	previousMedianCPUUsage := previousAsValues[string(benchspy.MedianCPUUsage)]

	assert.Equal(t, currentMedianCPUUsage.Type(), previousMedianCPUUsage.Type(), "types of metrics should be the same")

	currentMedianCPUUsageVector := currentMedianCPUUsage.(model.Vector)
	previousMedianCPUUsageVector := previousMedianCPUUsage.(model.Vector)

	assert.Equal(t, len(currentMedianCPUUsageVector), len(previousMedianCPUUsageVector), "number of samples in vectors should be the same")

	// here we could compare actual values, but most likely they will be very different and the test will fail
}

// this test requires CTFv2 node_set with observability stack to be running
func TestBenchSpy_Custom_Prometheus_Metrics(t *testing.T) {
	gen, err := wasp.NewGenerator(&wasp.Config{
		T:           t,
		GenName:     "vu",
		CallTimeout: 100 * time.Millisecond,
		LoadType:    wasp.VU,
		Schedule:    wasp.Plain(1, 10*time.Second),
		VU: wasp.NewMockVU(&wasp.MockVirtualUserConfig{
			CallSleep: 50 * time.Millisecond,
		}),
	})
	require.NoError(t, err)

	gen.Run(true)

	// no need to not pass name regexp pattern
	// we provide them directly in custom queries
	promConfig := benchspy.NewPrometheusConfig()

	customPrometheus, err := benchspy.NewPrometheusQueryExecutor(
		map[string]string{
			// scalar value
			"95p_cpu_all_containers": "scalar(quantile(0.95, rate(container_cpu_usage_seconds_total{name=~\"node[^0]\"}[5m])) * 100)",
			// matrix value
			"cpu_rate_by_container": "rate(container_cpu_usage_seconds_total{name=~\"node[^0]\"}[1m])[30m:1m]",
		},
		promConfig,
	)
	require.NoError(t, err)

	baseLineReport, err := benchspy.NewStandardReport(
		"91ee9e3c903d52de12f3d0c1a07ac3c2a6d141fb",
		benchspy.WithQueryExecutors(customPrometheus),
		benchspy.WithGenerators(gen),
	)
	require.NoError(t, err, "failed to create baseline report")

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
	})
	require.NoError(t, err)

	newGen.Run(true)

	fetchCtx, cancelFn = context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()

	currentReport, previousReport, err := benchspy.FetchNewStandardReportAndLoadLatestPrevious(
		fetchCtx,
		"2d1fa3532656c51991c0212afce5f80d2914e340",
		benchspy.WithQueryExecutors(customPrometheus),
		benchspy.WithGenerators(newGen),
	)
	require.NoError(t, err, "failed to fetch current report or load the previous one")
	require.Equal(t, baseLineReport.CommitOrTag, previousReport.CommitOrTag, "current report should be the same as the original report")

	currentAsValues := benchspy.MustAllPrometheusResults(currentReport)
	previousAsValues := benchspy.MustAllPrometheusResults(previousReport)

	assert.Equal(t, len(currentAsValues), len(previousAsValues), "number of metrics in results should be the same")

	current95CPUUsage := currentAsValues["95p_cpu_all_containers"]
	previous95CPUUsage := previousAsValues["95p_cpu_all_containers"]

	assert.Equal(t, current95CPUUsage.Type(), previous95CPUUsage.Type(), "types of metrics should be the same")
	assert.IsType(t, current95CPUUsage, &model.Scalar{}, "current metric should be a scalar")

	currentCPUByContainer := currentAsValues["cpu_rate_by_container"]
	previousCPUByContainer := previousAsValues["cpu_rate_by_container"]

	assert.Equal(t, currentCPUByContainer.Type(), previousCPUByContainer.Type(), "types of metrics should be the same")
	assert.IsType(t, currentCPUByContainer, model.Matrix{}, "current metric should be a scalar")

	current95CPUUsageAsMatrix := currentCPUByContainer.(model.Matrix)
	previous95CPUUsageAsMatrix := currentCPUByContainer.(model.Matrix)

	assert.Equal(t, len(current95CPUUsageAsMatrix), len(previous95CPUUsageAsMatrix), "number of samples in matrices should be the same")

	// here we could compare actual values, but most likely they will be very different and the test will fail
}
