# BenchSpy - Standard Prometheus Metrics

Now that we've seen how to query and assert load-related metrics, let's explore how to query and assert on resource usage by our Application Under Test (AUT).

If you're unsure why this is important, consider the following situation: the p95 latency of a new release matches the previous version, but memory consumption is 34% higher. Not ideal, right?

## Step 1: Prometheus Configuration

Since `WASP` has no built-in integration with `Prometheus`, we need to pass its configuration separately:

```go
promConfig := benchspy.NewPrometheusConfig("node[^0]")
```

This constructor loads the URL from the environment variable `PROMETHEUS_URL` and adds a single regex pattern to match containers **by name**. In this case, it excludes the bootstrap Chainlink node (named `node0` in the `CTFv2` stack).

> [!WARNING]
> This example assumes that you have both the observability stack and basic node set running.
> If you have the [CTF CLI](../../../framework/getting_started.md), you can start it by running: `ctf b ns`.

> [!NOTE]
> Matching containers **by name** should work both for most k8s and Docker setups using `CTFv2` observability stack.

## Step 2: Fetching and Storing a Baseline Report

As in previous examples, we'll use built-in Prometheus metrics and fetch and store a baseline report:

```go
baseLineReport, err := benchspy.NewStandardReport(
    "91ee9e3c903d52de12f3d0c1a07ac3c2a6d141fb",
    // notice the different standard query executor type
    benchspy.WithStandardQueries(benchspy.StandardQueryExecutor_Prometheus),
    benchspy.WithPrometheusConfig(promConfig),
    // Required to calculate test time range based on generator start/end times.
    benchspy.WithGenerators(gen),
)
require.NoError(t, err, "failed to create baseline report")

fetchCtx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
defer cancelFn()

fetchErr := baseLineReport.FetchData(fetchCtx)
require.NoError(t, fetchErr, "failed to fetch baseline report")

path, storeErr := baseLineReport.Store()
require.NoError(t, storeErr, "failed to store baseline report", path)
```

> [!NOTE]
> Standard metrics for Prometheus differ from those used by `Loki` or `Direct` query executors.
> Prometheus metrics focus on resource usage by the AUT, while `Loki`/`Direct` metrics measure load characteristics.
>
> Standard Prometheus metrics include:
> - `median_cpu_usage`
> - `median_mem_usage`
> - `max_cpu_usage`
> - `p95_cpu_usage`
> - `p95_mem_usage`
> - `max_mem_usage`
>
> These are calculated at the **container level**, based on total usage (user + system).

## Step 3: Handling Prometheus Result Types

Unlike Loki and Generator, Prometheus results can have various data types:
- `scalar`
- `string`
- `vector`
- `matrix`

This makes asserting results a bit more complex.

### Converting Results to `model.Value`

First, convert results to the `model.Value` interface using convenience functions:

```go
currentAsValues := benchspy.MustAllPrometheusResults(currentReport)
previousAsValues := benchspy.MustAllPrometheusResults(previousReport)
```

### Casting to Specific Types

Next, determine the data type returned by your query and cast it accordingly:

```go
// Fetch a single metric
currentMedianCPUUsage := currentAsValues[string(benchspy.MedianCPUUsage)]
previousMedianCPUUsage := previousAsValues[string(benchspy.MedianCPUUsage)]

assert.Equal(t, currentMedianCPUUsage.Type(), previousMedianCPUUsage.Type(), "types of metrics should be the same")

// In this case, we know the query returns a Vector
currentMedianCPUUsageVector := currentMedianCPUUsage.(model.Vector)
previousMedianCPUUsageVector := previousMedianCPUUsage.(model.Vector)
```

Since these metrics are not related to load generation, the convenience function a `map[string](model.Value)`, where key is resource metric name.

> [!WARNING]
> All standard Prometheus metrics bundled with `BenchSpy` return `model.Vector`.
> However, if you use custom queries, you must manually verify their return types.

## Skipping Assertions for Resource Usage

We skip the assertion part because, unless you're comparing resource usage under stable loads, significant differences between reports are likely.
For example:
- The first report might be generated right after the node set starts.
- The second report might be generated after the node set has been running for some time.

## What’s Next?

In the next chapter, we’ll [explore custom Prometheus queries](./prometheus_custom.md).

> [!NOTE]
> You can find the full example [here](https://github.com/smartcontractkit/chainlink-testing-framework/tree/main/wasp/examples/benchspy/prometheus_query_executor/prometheus_query_executor_test.go).