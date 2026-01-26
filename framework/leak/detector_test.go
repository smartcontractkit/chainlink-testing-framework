package leak_test

import (
	"fmt"
	"testing"
	"time"

	f "github.com/smartcontractkit/chainlink-testing-framework/framework"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/leak"
)

func mustTime(start string) time.Time {
	s, err := time.Parse(time.RFC3339, start)
	if err != nil {
		panic("can't convert time from RFC3339")
	}
	return s
}

func TestSmokeMeasure(t *testing.T) {
	qc := leak.NewFakeQueryClient()
	lc := leak.NewResourceLeakChecker(leak.WithQueryClient(qc))
	testCases := []struct {
		name           string
		startTime      time.Time
		endTime        time.Time
		startResponse  *f.PrometheusQueryResponse
		endResponse    *f.PrometheusQueryResponse
		warmUpDuration time.Duration
		expectedDiff   float64
		errorContains  string
	}{
		{
			name:          "diff is correct and > 0",
			startTime:     mustTime("2026-01-12T21:53:00Z"),
			endTime:       mustTime("2026-01-13T10:11:00Z"),
			startResponse: leak.PromSingleValueResponse("10"),
			endResponse:   leak.PromSingleValueResponse("20"),
			expectedDiff:  100,
		},
		{
			name:          "diff is correct and < 0",
			startTime:     mustTime("2026-01-12T21:53:00Z"),
			endTime:       mustTime("2026-01-13T10:11:00Z"),
			startResponse: leak.PromSingleValueResponse("20"),
			endResponse:   leak.PromSingleValueResponse("0"),
			expectedDiff:  -100,
		},
		{
			name:          "start > end time",
			startTime:     mustTime("2026-01-13T10:11:00Z"),
			endTime:       mustTime("2026-01-12T21:53:00Z"),
			startResponse: leak.PromSingleValueResponse("10"),
			endResponse:   leak.PromSingleValueResponse("20"),
			errorContains: "start time is greated than end time",
		},
		{
			name:           "works with warm up duration",
			startTime:      mustTime("2026-01-01T10:00:00Z"),
			endTime:        mustTime("2026-01-01T11:00:00Z"),
			startResponse:  leak.PromSingleValueResponse("10"),
			endResponse:    leak.PromSingleValueResponse("15"),
			warmUpDuration: 29 * time.Minute,
			expectedDiff:   50,
		},
		{
			name:           "warm up is too long",
			startTime:      mustTime("2026-01-01T10:00:00Z"),
			endTime:        mustTime("2026-01-01T11:00:00Z"),
			startResponse:  leak.PromSingleValueResponse("10"),
			endResponse:    leak.PromSingleValueResponse("20"),
			warmUpDuration: 31 * time.Minute,
			errorContains:  "warm up duration can't be more than 50 percent",
		},
		{
			name:      "no results for start time",
			startTime: mustTime("2026-01-12T21:53:00Z"),
			endTime:   mustTime("2026-01-13T10:11:00Z"),
			startResponse: &f.PrometheusQueryResponse{
				Data: &f.PromQueryResponseData{
					Result: nil,
				},
			},
			endResponse:   leak.PromSingleValueResponse("20"),
			errorContains: "no results for start timestamp",
		},
		{
			name:          "no results for end time",
			startTime:     mustTime("2026-01-12T21:53:00Z"),
			endTime:       mustTime("2026-01-13T10:11:00Z"),
			startResponse: leak.PromSingleValueResponse("10"),
			endResponse: &f.PrometheusQueryResponse{
				Data: &f.PromQueryResponseData{
					Result: nil,
				},
			},
			errorContains: "no results for end timestamp",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			qc.SetResponses(tc.startResponse, tc.endResponse)
			diff, err := lc.MeasureDelta(&leak.CheckConfig{
				// Prometheus returns good errors when query is invalid
				// so we do not test it since there is no additional validation
				Query:          ``,
				Start:          tc.startTime,
				End:            tc.endTime,
				WarmUpDuration: tc.warmUpDuration,
			})
			if tc.errorContains != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errorContains)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.expectedDiff, diff.Delta)
		})
	}
}

func TestRealCLNodesLeakDetectionLocalDevenv(t *testing.T) {
	t.Skip(`this test requires a real load run, see docs here https://github.com/smartcontractkit/chainlink/tree/develop/devenv, spin up the env and run "cl test load"`)

	cnd, err := leak.NewCLNodesLeakDetector(leak.NewResourceLeakChecker())
	require.NoError(t, err)
	errs := cnd.Check(&leak.CLNodesCheck{
		NumNodes: 4,
		// set timestamps for the run you are analyzing
		Start:           mustTime("2026-01-19T17:23:14Z"),
		End:             mustTime("2026-01-19T18:00:51Z"),
		CPUThreshold:    100.0,
		MemoryThreshold: 20.0,
	})
	require.NoError(t, errs)
	fmt.Println(errs)
}

func TestRealPrometheusLowLevelAPI(t *testing.T) {
	t.Skip(`this test requires a real load run, see docs here https://github.com/smartcontractkit/chainlink/tree/develop/devenv, spin up the env and run "cl test load"`)

	// demonstrates how to use low-level API for custom queries with CL nodes example
	donNodes := 4
	resourceLeaks := make([]*leak.Measurement, 0)

	lc := leak.NewResourceLeakChecker()
	for i := range donNodes {
		diff, err := lc.MeasureDelta(&leak.CheckConfig{
			Query: fmt.Sprintf(`quantile_over_time(0.5, container_memory_rss{name="don-node%d"}[1h]) / 1024 / 1024`, i),
			// set timestamps for the run you are analyzing
			Start:          mustTime("2026-01-12T21:53:00Z"),
			End:            mustTime("2026-01-13T10:11:00Z"),
			WarmUpDuration: 1 * time.Hour,
		})
		require.NoError(t, err)
		resourceLeaks = append(resourceLeaks, diff)
	}
	require.Len(t, resourceLeaks, 4)

	fmt.Println(resourceLeaks)
	for _, ml := range resourceLeaks {
		require.GreaterOrEqual(t, ml, 0.5)
	}
}
