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

func TestBenchSpy_Standard_Generator_Metrics(t *testing.T) {
	gen, err := wasp.NewGenerator(&wasp.Config{
		T:           t,
		GenName:     "vu",
		CallTimeout: 100 * time.Millisecond,
		LoadType:    wasp.VU,
		Schedule:    wasp.Plain(10, 15*time.Second),
		VU: wasp.NewMockVU(&wasp.MockVirtualUserConfig{
			CallSleep: 50 * time.Millisecond,
		}),
	})
	require.NoError(t, err)

	gen.Run(true)

	fetchCtx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()

	baseLineReport, err := benchspy.NewStandardReport(
		"e7fc5826a572c09f8b93df3b9f674113372ce924",
		benchspy.WithStandardQueries(benchspy.StandardQueryExecutor_Generator),
		benchspy.WithGenerators(gen),
	)
	require.NoError(t, err, "failed to create baseline report")

	fetchErr := baseLineReport.FetchData(fetchCtx)
	require.NoError(t, fetchErr, "failed to fetch data for original report")

	path, storeErr := baseLineReport.Store()
	require.NoError(t, storeErr, "failed to store current report", path)

	newGen, err := wasp.NewGenerator(&wasp.Config{
		T:           t,
		GenName:     "vu",
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

	// currentReport is the report that we just created (baseLineReport)
	currentReport, previousReport, err := benchspy.FetchNewStandardReportAndLoadLatestPrevious(
		fetchCtx,
		"e7fc5826a572c09f8b93df3b9f674113372ce925",
		benchspy.WithStandardQueries(benchspy.StandardQueryExecutor_Generator),
		benchspy.WithGenerators(newGen),
	)
	require.NoError(t, err, "failed to fetch current report or load the previous one")

	// make sure that previous report is the same as the baseline report
	require.Equal(t, baseLineReport.CommitOrTag, previousReport.CommitOrTag, "current report should be the same as the original report")

	currentAsFloat64 := benchspy.MustAllGeneratorResults(currentReport)
	previousAsloat64 := benchspy.MustAllGeneratorResults(previousReport)

	var compareValues = func(
		metricName string,
		maxDiffPercentage float64,
	) {
		require.NotNil(t, currentAsFloat64[metricName], "%s results were missing from current report", metricName)
		require.NotNil(t, previousAsloat64[metricName], "%s results were missing from previous report", metricName)

		currentMetric := currentAsFloat64[metricName]
		previousMetric := previousAsloat64[metricName]

		var diffPrecentage float64
		if previousMetric != 0.0 {
			diffPrecentage = (currentMetric - previousMetric) / previousMetric * 100
		} else {
			diffPrecentage = currentMetric * 100.0
		}
		assert.LessOrEqual(t, math.Abs(diffPrecentage), maxDiffPercentage, "%s medians are more than 1% different", metricName, fmt.Sprintf("%.4f", diffPrecentage))
	}

	compareValues(string(benchspy.MedianLatency), 1.0)
	compareValues(string(benchspy.Percentile95Latency), 1.0)
	compareValues(string(benchspy.ErrorRate), 1.0)
}
