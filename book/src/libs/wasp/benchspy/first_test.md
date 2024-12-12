# BenchSpy - Your first test

Let's start with a simplest case, which doesn't require you to have any of the observability stack, but only `WASP` and the application you are testing.
`BenchSpy` comes with some built-in `QueryExecutors` each of which additionaly has predefined metrics that you can use. One of these executors is the
`GeneratorQueryExecutor` that fetches metrics directly from `WASP` generators.

Our first test will follow the following logic:
* Run a simple load test
* Generate the performance report and store it
* Run the load again
* Generate a new report and compare it to the previous one

We will use some very simplified assertions, used only for the sake of example, and expect the performance to remain unchanged.

Let's start by defining and running a generator that will use a mocked service:
```go
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
```

Now that we have load data, let's generate a baseline performance report and store it in the local storage:
```go
fetchCtx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
defer cancelFn()

baseLineReport, err := benchspy.NewStandardReport(
    // random hash, this should be commit or hash of the Application Under Test (AUT)
    "e7fc5826a572c09f8b93df3b9f674113372ce924",
    // use built-in queries for an executor that fetches data directly from the WASP generator
    benchspy.WithStandardQueryExecutorType(benchspy.StandardQueryExecutor_Generator),
    // WASP generators
    benchspy.WithGenerators(gen),
)
require.NoError(t, err, "failed to create original report")

fetchErr := baseLineReport.FetchData(fetchCtx)
require.NoError(t, fetchErr, "failed to fetch data for original report")

path, storeErr := baseLineReport.Store()
require.NoError(t, storeErr, "failed to store current report", path)
```

> [!NOTE]
> There's quite a lot to unpack here and you are enouraged to read more about build-in `QueryExecutors` and
> standard metrics each comes with [here](./built_in_query_executors.md) and about the `StandardReport` [here](./standard_report.md).
>
> For now, it's enough for you to know that standard metrics that `StandardQueryExecutor_Generator` comes with are following:
> * median latency
> * p95 latency (95th percentile)
> * error rate

With baseline report ready let's run the load test again, but this time let's use a wrapper function
that will automatically load the previous report, generate a new one and make sure that they are actually comparable.
```go
// define a new generator using the same config values
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

// run the load
newGen.Run(true)

fetchCtx, cancelFn = context.WithTimeout(context.Background(), 60*time.Second)
defer cancelFn()

// currentReport is the report that we just created (baseLineReport)
currentReport, previousReport, err := benchspy.FetchNewStandardReportAndLoadLatestPrevious(
    fetchCtx,
    "e7fc5826a572c09f8b93df3b9f674113372ce925",
    benchspy.WithStandardQueryExecutorType(benchspy.StandardQueryExecutor_Generator),
    benchspy.WithGenerators(newGen),
)
require.NoError(t, err, "failed to fetch current report or load the previous one")
```

> [!NOTE]
> In real-world case, once you have the first report generated you should only need to use
> `benchspy.FetchNewStandardReportAndLoadLatestPrevious` function.

Okay, so we have two reports now, that's great, but how do we make sure that application's performance is as expected?
You'll find out in the [next chapter](./first_test_comparison.md).