package main

// TODO: fix errors assertions

//// this test can be run without external dependencies
//// it demontsrates a case, where performance degradation was found between two releases
//func TestBenchSpy_Standard_Direct_Metrics_RealCase(t *testing.T) {
//	// uncomment the code below and comment the rest
//	// to generate the v1.0.0 report
//
//	// generator, err := wasp.NewGenerator(&wasp.Config{
//	// 	T:           t,
//	// 	GenName:     "vu",
//	// 	CallTimeout: 100 * time.Millisecond,
//	// 	LoadType:    wasp.VU,
//	// 	Schedule:    wasp.Plain(10, 15*time.Second),
//	// 	VU: wasp.NewMockVU(&wasp.MockVirtualUserConfig{
//	// 		// notice lower latency
//	// 		CallSleep: 50 * time.Millisecond,
//	// 	}),
//	// })
//	// require.NoError(t, err)
//
//	// generator.Run(true)
//
//	// fetchCtx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
//	// defer cancelFn()
//
//	// baseLineReport, err := benchspy.NewStandardReport(
//	// 	"v1.0.0",
//	// 	benchspy.WithStandardQueries(benchspy.StandardQueryExecutor_Direct),
//	// 	benchspy.WithGenerators(generator),
//	// )
//	// require.NoError(t, err, "failed to create baseline report")
//
//	// fetchErr := baseLineReport.FetchData(fetchCtx)
//	// require.NoError(t, fetchErr, "failed to fetch data for original report")
//
//	// path, storeErr := baseLineReport.Store()
//	// require.NoError(t, storeErr, "failed to store current report", path)
//
//	// comment this part, when generating the v1.0.0 report
//	generator, err := wasp.NewGenerator(&wasp.Config{
//		T:           t,
//		GenName:     "vu",
//		CallTimeout: 100 * time.Millisecond,
//		LoadType:    wasp.VU,
//		Schedule:    wasp.Plain(10, 15*time.Second),
//		VU: wasp.NewMockVU(&wasp.MockVirtualUserConfig{
//			// increase latency by 10ms
//			CallSleep: 60 * time.Millisecond,
//		}),
//	})
//	require.NoError(t, err)
//	generator.Run(true)
//
//	currentVersion := os.Getenv("CURRENT_VERSION")
//	require.NotEmpty(t, currentVersion, "No current version provided")
//
//	fetchCtx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
//	defer cancelFn()
//
//	currentReport, previousReport, err := benchspy.FetchNewStandardReportAndLoadLatestPrevious(
//		fetchCtx,
//		currentVersion,
//		benchspy.WithStandardQueries(benchspy.StandardQueryExecutor_Direct),
//		benchspy.WithReportDirectory("test_reports"),
//		benchspy.WithGenerators(generator),
//	)
//	require.NoError(t, err, "failed to fetch current report or load the previous one")
//
//	hasErrors, errors := benchspy.CompareDirectWithThresholds(1.0, 1.0, 1.0, 1.0, currentReport, previousReport)
//	require.True(t, hasErrors, "Found no errors, but expected some")
//	require.Equal(t, 3, len(errors), "Expected 3 errors, got %d", len(errors))
//
//	expectedErrors := []benchspy.StandardLoadMetric{benchspy.MedianLatency, benchspy.Percentile95Latency, benchspy.MaxLatency}
//	var foundErrors []benchspy.StandardLoadMetric
//
//	for _, e := range errors[generator.Cfg.GenName] {
//		for _, expected := range expectedErrors {
//			if strings.Contains(e.Error(), string(expected)) {
//				foundErrors = append(foundErrors, expected)
//				break
//			}
//		}
//	}
//
//	require.EqualValues(t, expectedErrors, foundErrors, "Expected errors not found")
//}
