# BenchSpy - Custom Prometheus metrics

Similarly to what we have done with Loki, we can use custom metrics with Prometheus.

Most of the code is the same as in previous example. Differences start with the need to manually
create a `PrometheusQueryExecutor` with our custom queries:

```go
// no need to not pass name regexp pattern
// we provide them directly in custom queries
promConfig := benchspy.NewPrometheusConfig()

customPrometheus, err := benchspy.NewPrometheusQueryExecutor(
    map[string]string{
        // scalar value
        "95p_cpu_all_containers": "scalar(quantile(0.95, rate(container_cpu_usage_seconds_total{name=~\"node[^0]\"}[5m])) * 100)",
        // matrix value
        "cpu_rate_by_container": "rate(container_cpu_usage_seconds_total{name=~\"node[^0]\"}[1m])[30m:1m]",
    },
    *promConfig,
)
```

Then we pass them as custom query executor:
```go
baseLineReport, err := benchspy.NewStandardReport(
    "91ee9e3c903d52de12f3d0c1a07ac3c2a6d141fb",
    benchspy.WithQueryExecutors(customPrometheus),
    benchspy.WithGenerators(gen),
)
require.NoError(t, err, "failed to create baseline report")
```

> [!NOTE]
> Notice that when using custom Prometheus queries we don't need to pass the `PrometheusConfig`
> to `NewStandardReport()`, because we have already set it when creating `PrometheusQueryExecutor`.

Fetching of current and previous report remain unchanged, just like getting Prometheus metrics cast
to it's specific type:
```go
currentAsValues := benchspy.MustAllPrometheusResults(currentReport)
previousAsValues := benchspy.MustAllPrometheusResults(previousReport)

assert.Equal(t, len(currentAsValues), len(previousAsValues), "number of metrics in results should be the same")
```

But now comes another difference. All standard query results were instances of `model.Vector`. Our two custom queries
introduce two new types:
* `model.Matrix`
* `*model.Scalar`

And these differences are reflected in further casting that we do, before getting final metrics:
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

> [!NOTE]
> When casting to Prometheus' final types it's crucial to remember that two types have pointer receivers and the other two value receivers.
>
> Pointer receivers:
> * `*model.String`
> * `*model.Scalar`
>
> Value receivers:
> * `model.Vector`
> * `model.Matrix`