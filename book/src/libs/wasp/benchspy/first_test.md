# BenchSpy - Your First Test

Let's start with the simplest case, which doesn't require any part of the observability stack—only `WASP` and the application you are testing.
`BenchSpy` comes with built-in `QueryExecutors`, each of which also has predefined metrics that you can use. One of these executors is the `DirectQueryExecutor`, which fetches metrics directly from `WASP` generators,
which means you can run it with Loki.

> [!NOTE]
> Not sure whether to use `Loki` or `Direct` query executors? [Read this!](./loki_dillema.md)

## Test Overview

Our first test will follow this logic:
- Run a simple load test.
- Generate a performance report and store it.
- Run the load test again.
- Generate a new report and compare it to the previous one.

We'll use very simplified assertions for this example and expect the performance to remain unchanged.

### Step 1: Define and Run a Generator

Let's start by defining and running a generator that uses a mocked service:

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

### Step 2: Generate a Baseline Performance Report

With load data available, let's generate a baseline performance report and store it in local storage:

```go
baseLineReport, err := benchspy.NewStandardReport(
    // random hash, this should be the commit or hash of the Application Under Test (AUT)
    "v1.0.0",
    // use built-in queries for an executor that fetches data directly from the WASP generator
    benchspy.WithStandardQueries(benchspy.StandardQueryExecutor_Direct),
    // WASP generators
    benchspy.WithGenerators(gen),
)
require.NoError(t, err, "failed to create baseline report")

fetchCtx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
defer cancelFn()

fetchErr := baseLineReport.FetchData(fetchCtx)
require.NoError(t, fetchErr, "failed to fetch data for baseline report")

path, storeErr := baseLineReport.Store()
require.NoError(t, storeErr, "failed to store baseline report", path)
```

> [!NOTE]
> There's a lot to unpack here, and you're encouraged to read more about the built-in `QueryExecutors` and the standard metrics they provide as well as about the `StandardReport` [here](./reports/standard_report.md).
>
> For now, it's enough to know that the standard metrics provided by `StandardQueryExecutor_Direct` include:
> - Median latency
> - P95 latency (95th percentile)
> - Max latency
> - Error rate

### Step 3: Run the Test Again and Compare Reports

With the baseline report ready, let's run the load test again. This time, we'll use a wrapper function to automatically load the previous report, generate a new one, and ensure they are comparable.

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
    // commit or tag of the new application version
    "v2.0.0",
    benchspy.WithStandardQueries(benchspy.StandardQueryExecutor_Direct),
    benchspy.WithGenerators(newGen),
)
require.NoError(t, err, "failed to fetch current report or load the previous one")
```

> [!NOTE]
> In a real-world case, once you've generated the first report, you should only need to use the `benchspy.FetchNewStandardReportAndLoadLatestPrevious` function.

### What's Next?

Now that we have two reports, how do we ensure that the application's performance meets expectations?
Find out in the [next chapter](./simplest_metrics.md).