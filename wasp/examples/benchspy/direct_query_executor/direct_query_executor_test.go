package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
	"github.com/smartcontractkit/chainlink-testing-framework/wasp/benchspy"
)

// this test can be run without external dependencies
func TestBenchSpy_Standard_Direct_Metrics(t *testing.T) {
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
		"v1",
		benchspy.WithStandardQueries(benchspy.StandardQueryExecutor_Direct),
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
		"v2",
		benchspy.WithStandardQueries(benchspy.StandardQueryExecutor_Direct),
		benchspy.WithGenerators(newGen),
	)
	require.NoError(t, err, "failed to fetch current report or load the previous one")

	// make sure that previous report is the same as the baseline report
	require.Equal(t, baseLineReport.CommitOrTag, previousReport.CommitOrTag, "current report should be the same as the original report")

	hasErrors, errors := benchspy.CompareDirectWithThresholds(1.0, 1.0, 1.0, 1.0, currentReport, previousReport)
	require.False(t, hasErrors, fmt.Sprintf("errors found: %v", errors))
}

func TestBenchSpy_Standard_Direct_Metrics_Two_Generators(t *testing.T) {
	gen1, err := wasp.NewGenerator(&wasp.Config{
		T:           t,
		GenName:     "vu1",
		CallTimeout: 200 * time.Millisecond,
		LoadType:    wasp.VU,
		Schedule:    wasp.Plain(10, 15*time.Second),
		VU: wasp.NewMockVU(&wasp.MockVirtualUserConfig{
			CallSleep: 50 * time.Millisecond,
		}),
	})
	require.NoError(t, err)

	gen2, err := wasp.NewGenerator(&wasp.Config{
		T:           t,
		GenName:     "vu2",
		CallTimeout: 200 * time.Millisecond,
		LoadType:    wasp.VU,
		Schedule:    wasp.Plain(10, 20*time.Second),
		VU: wasp.NewMockVU(&wasp.MockVirtualUserConfig{
			CallSleep: 60 * time.Millisecond,
		}),
	})
	require.NoError(t, err)

	gen1.Run(false)
	gen2.Run(true)

	fetchCtx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()

	baseLineReport, err := benchspy.NewStandardReport(
		"v1",
		benchspy.WithStandardQueries(benchspy.StandardQueryExecutor_Direct),
		benchspy.WithGenerators(gen1, gen2),
	)
	require.NoError(t, err, "failed to create baseline report")

	fetchErr := baseLineReport.FetchData(fetchCtx)
	require.NoError(t, fetchErr, "failed to fetch data for original report")

	path, storeErr := baseLineReport.Store()
	require.NoError(t, storeErr, "failed to store current report", path)

	newGen1, err := wasp.NewGenerator(&wasp.Config{
		T:           t,
		GenName:     "vu1",
		CallTimeout: 200 * time.Millisecond,
		LoadType:    wasp.VU,
		Schedule:    wasp.Plain(10, 15*time.Second),
		VU: wasp.NewMockVU(&wasp.MockVirtualUserConfig{
			CallSleep: 50 * time.Millisecond,
		}),
	})
	require.NoError(t, err)

	newGen2, err := wasp.NewGenerator(&wasp.Config{
		T:           t,
		GenName:     "vu2",
		CallTimeout: 200 * time.Millisecond,
		LoadType:    wasp.VU,
		Schedule:    wasp.Plain(10, 20*time.Second),
		VU: wasp.NewMockVU(&wasp.MockVirtualUserConfig{
			CallSleep: 60 * time.Millisecond,
		}),
	})
	require.NoError(t, err)

	newGen1.Run(false)
	newGen2.Run(true)

	fetchCtx, cancelFn = context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()

	// currentReport is the report that we just created (baseLineReport)
	currentReport, previousReport, err := benchspy.FetchNewStandardReportAndLoadLatestPrevious(
		fetchCtx,
		"v2",
		benchspy.WithStandardQueries(benchspy.StandardQueryExecutor_Direct),
		benchspy.WithGenerators(newGen1, newGen2),
	)
	require.NoError(t, err, "failed to fetch current report or load the previous one")

	// make sure that previous report is the same as the baseline report
	require.Equal(t, baseLineReport.CommitOrTag, previousReport.CommitOrTag, "current report should be the same as the original report")

	hasErrors, errors := benchspy.CompareDirectWithThresholds(10.0, 10.0, 10.0, 10.0, currentReport, previousReport)
	require.False(t, hasErrors, fmt.Sprintf("errors found: %v", errors))
}
