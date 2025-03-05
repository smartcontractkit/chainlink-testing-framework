package main

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"testing"
	"time"

	"github.com/montanaflynn/stats"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
	"github.com/smartcontractkit/chainlink-testing-framework/wasp/benchspy"
)

// this test requires CTFv2 observability stack to be running
func TestBenchSpy_Standard_Direct_And_Loki_Metrics(t *testing.T) {
	label := "benchspy-direct-loki"

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
		benchspy.WithStandardQueries(benchspy.StandardQueryExecutor_Direct, benchspy.StandardQueryExecutor_Loki),
		benchspy.WithGenerators(gen),
	)
	require.NoError(t, err, "failed to create original report")

	fetchCtx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()

	fetchErr := baseLineReport.FetchData(fetchCtx)
	require.NoError(t, fetchErr, "failed to fetch current report")

	currentAsLokiSlices := benchspy.MustAllLokiResults(baseLineReport)[gen.Cfg.GenName]
	currentAsDirectFloats := benchspy.MustAllDirectResults(baseLineReport)[gen.Cfg.GenName]

	require.NotEmpty(t, currentAsLokiSlices[string(benchspy.MedianLatency)], "%s results were missing for loki", string(benchspy.MedianLatency))
	require.NotEmpty(t, currentAsDirectFloats[string(benchspy.MedianLatency)], "%s results were missing for direct", string(benchspy.MedianLatency))

	var compareValues = func(t *testing.T, metricName string, lokiFloat, directFloat, maxDiffPercentage float64) {
		var diffPercentage float64
		if lokiFloat != 0.0 && directFloat != 0.0 {
			diffPercentage = (directFloat - lokiFloat) / lokiFloat * 100
		} else if lokiFloat == 0.0 && directFloat == 0.0 {
			diffPercentage = 0.0
		} else {
			diffPercentage = 100.0
		}
		assert.LessOrEqual(t, math.Abs(diffPercentage), maxDiffPercentage, "%s are more than 1% different", metricName, fmt.Sprintf("%.4f", diffPercentage))
	}

	lokiFloatSlice, err := benchspy.StringSliceToFloat64Slice(currentAsLokiSlices[string(benchspy.MedianLatency)])
	require.NoError(t, err, "failed to convert %s results to float64 slice", string(benchspy.MedianLatency))
	lokiMedian, err := stats.Median(lokiFloatSlice)
	require.NoError(t, err, "failed to calculate median for loki %s results", string(benchspy.MedianLatency))

	compareValues(t, string(benchspy.MedianLatency), lokiMedian, currentAsDirectFloats[string(benchspy.MedianLatency)], 1.0)

	lokip95, err := stats.Percentile(lokiFloatSlice, 95)
	require.NoError(t, err, "failed to calculate 95th percentile for loki %s results", string(benchspy.Percentile95Latency))

	// here the max diff is 1.5% because of higher impact of data aggregation in loki
	compareValues(t, string(benchspy.Percentile95Latency), lokip95, currentAsDirectFloats[string(benchspy.Percentile95Latency)], 1.5)

	lokiErrorRate := 0
	for _, v := range currentAsLokiSlices[string(benchspy.ErrorRate)] {
		asInt, err := strconv.Atoi(v)
		require.NoError(t, err)
		lokiErrorRate += int(asInt)
	}

	lokiErrorRate = lokiErrorRate / len(currentAsLokiSlices[string(benchspy.ErrorRate)])
	compareValues(t, string(benchspy.ErrorRate), float64(lokiErrorRate), currentAsDirectFloats[string(benchspy.ErrorRate)], 1.0)
}
