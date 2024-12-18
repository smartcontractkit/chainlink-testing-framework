# BenchSpy - Real world example

Now that we have seen all possible usages, you might wonder how you should write a test that compares performance between different
releases of your application.

Usually steps to follow would look like this:
1. Write the performance test.
2. At the end of the test fetch the report, store it and commit to git.
3. Modify the previous point, so that it fetches both latest report and creates a new one.
4. Write your assertions for metrics.

# Writing the performance test
We will use a simple mock for the application under test. All that it does is wait for `50 ms` before
returning a 200 response code.

```go
generator, err := wasp.NewGenerator(&wasp.Config{
    T:           t,
    GenName:     "vu",
    CallTimeout: 100 * time.Millisecond,
    LoadType:    wasp.VU,
    Schedule:    wasp.Plain(10, 15*time.Second),
    VU: wasp.NewMockVU(&wasp.MockVirtualUserConfig{
        // notice lower latency
        CallSleep: 50 * time.Millisecond,
    }),
})
require.NoError(t, err)

generator.Run(true)
```

# Generating first report
Here we generate a new performance report for `v1.0.0`. We will use `Direct` query executor and save the report to a custom directory
called `test_reports`. We will use this report later to compare the performance of new versions.

```go
fetchCtx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
defer cancelFn()

baseLineReport, err := benchspy.NewStandardReport(
    "v1.0.0",
    benchspy.WithStandardQueries(benchspy.StandardQueryExecutor_Direct),
    benchspy.WithReportDirectory("test_reports"),
    benchspy.WithGenerators(gen),
)
require.NoError(t, err, "failed to create baseline report")

fetchErr := baseLineReport.FetchData(fetchCtx)
require.NoError(t, fetchErr, "failed to fetch data for original report")

path, storeErr := baseLineReport.Store()
require.NoError(t, storeErr, "failed to store current report", path)
```

# Modifying report generation
Now that we have a baseline report stored for `v1.0.0` lets modify the test, so that we can use it with future releases of our application.
That means that the code from previous step has to change to:
```go
fetchCtx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
defer cancelFn()

currentReport, previousReport, err := benchspy.FetchNewStandardReportAndLoadLatestPrevious(
    fetchCtx,
    "v1.1.0",
    benchspy.WithStandardQueries(benchspy.StandardQueryExecutor_Direct),
    benchspy.WithReportDirectory("test_reports"),
    benchspy.WithGenerators(generator),
)
require.NoError(t, err, "failed to fetch current report or load the previous one")
```

As you remember this function will load latest report from `test_reports` directory and fetch a current one, in this case for `v1.1.0`.

# Adding assertions
Let's assume we don't want any of performance metrics to get more than **1% worse** between releases and use a convenience function
for `Direct` query executor:
```go
hasErrors, errors := benchspy.CompareDirectWithThresholds(
    1.0, // max 1% worse median latency
    1.0, // max 1% worse p95 latency
    1.0, // max 1% worse maximum latency
    0.0, // no change in error rate
    currentReport, previousReport)
require.False(t, hasErrors, fmt.Sprintf("errors found: %v", errors))
```

Done, you're ready to use `BenchSpy` to make sure that the performance of your application didn't degrade below your chosen thresholds!

> [!NOTE]
> You can find a test example, where the performance has degraded significantly [here](https://github.com/smartcontractkit/chainlink-testing-framework/tree/main/wasp/examples/benchspy/direct_query_executor/direct_query_real_case.go).
>
> This test passes, because we expect the performance to be worse. This is, of course, the opposite what you should do in case of a real application :-)