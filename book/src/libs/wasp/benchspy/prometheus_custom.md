# BenchSpy - Custom Prometheus Metrics

Similar to what we did with Loki, we can use custom metrics with Prometheus.

Most of the code remains the same as in the previous example. However, the differences begin with the need to manually create a `PrometheusQueryExecutor` with our custom queries:

```go
// No need to pass the name regex pattern, as we provide it directly in the queries
// Remeber that you are free to use any other matching values or labels or none at all
promConfig := benchspy.NewPrometheusConfig()

customPrometheus, err := benchspy.NewPrometheusQueryExecutor(
    map[string]string{
        // Scalar value
        "95p_cpu_all_containers": "scalar(quantile(0.95, rate(container_cpu_usage_seconds_total{name=~\"node[^0]\"}[5m])) * 100)",
        // Matrix value
        "cpu_rate_by_container": "rate(container_cpu_usage_seconds_total{name=~\"node[^0]\"}[1m])[30m:1m]",
    },
    promConfig,
)
```

## Passing Custom Queries to the Report

Next, pass the custom queries as a query executor:

```go
baseLineReport, err := benchspy.NewStandardReport(
    "91ee9e3c903d52de12f3d0c1a07ac3c2a6d141fb",
    // notice the different functional option used to pass Prometheus executor with custom queries
    benchspy.WithQueryExecutors(customPrometheus),
    benchspy.WithGenerators(gen),
    // notice that no Prometehus config is passed here
)
require.NoError(t, err, "failed to create baseline report")
```

> [!NOTE]
> When using custom Prometheus queries, you don’t need to pass the `PrometheusConfig` to `NewStandardReport()`, as the URL already been set during the creation of the `PrometheusQueryExecutor`.

## Fetching and Casting Metrics

Fetching the current and previous reports remains unchanged, as does casting Prometheus metrics to their specific types:

```go
currentAsValues := benchspy.MustAllPrometheusResults(currentReport)
previousAsValues := benchspy.MustAllPrometheusResults(previousReport)

assert.Equal(t, len(currentAsValues), len(previousAsValues), "number of metrics in results should be the same")
```

## Handling Different Data Types

Here’s where things differ. While all standard query results are instances of `model.Vector`, the two custom queries introduce new types:
- `model.Matrix`
- `*model.Scalar`

These differences are reflected in the further casting process before accessing the final metrics:

```go
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
```

> [!WARNING]
> When casting to Prometheus' final types, it’s crucial to remember the distinction between pointer and value receivers:
>
> **Pointer receivers**:
> - `*model.String`
> - `*model.Scalar`
>
> **Value receivers**:
> - `model.Vector`
> - `model.Matrix`

And that's it! You know all you need to know to unlock the full power of `BenchSpy`!

> [!NOTE]
> You can find the full example [here](https://github.com/smartcontractkit/chainlink-testing-framework/tree/main/wasp/examples/benchspy/prometheus_query_executor/prometheus_query_executor_test.go).