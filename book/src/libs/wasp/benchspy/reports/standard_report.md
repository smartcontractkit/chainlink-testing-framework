# BenchSpy - Standard Report

`StandardReport` comes with built-in support for three types of data sources:
* `Direct`
* `Loki`
* `Prometheus`

Each of them allows you to both use pre-defined metrics or use your own.

## Pre-defined (standard) metrics

### Direct and Loki
Both query executors focus on the characteristics of the load generated with WASP.
The datasets they work on are almost identical, because the former allows you to query load-specific
data before its sent to Loki. The latter offers you richer querying options (via `LogQL`) and access
to actual load profile (as opposed to the configured one).

Both query executors have following predefined metrics:
* median latency
* 95th percentile latency
* error rate

Latency is understood as the round time from making a request to receiving a response
from the Application Under Test.

Error rate is the ratio of failed responses to the total number of responses. This include
both requests that timed out or returned an error from `Gun` or `Vu` implementation.

### Prometehus
On the other hand, these standard metrics focus on resource consumption by the application you are testing,
instead on the load generation.

They include the following:
* median CPU usage
* 95th percentil of CPU usage
* median memory usage
* 95th percentil of memory usage

In both cases queries focus on `total` consumption, which consists of the sum of what the underlaying system and
you appplication uses.

### How to use
As mentioned in the examples, to use predefined metrics you should use the `NewStandardReport` method:
```go
report, err := benchspy.NewStandardReport(
    "91ee9e3c903d52de12f3d0c1a07ac3c2a6d141fb",
    // Query executor types for which standard metrics should be generated
    benchspy.WithStandardQueries(benchspy.StandardQueryExecutor_Prometheus, benchspy.StandardQueryExecutor_Loki),
    // Prometheus configuration is required if using standard Prometheus metrics
    benchspy.WithPrometheusConfig(benchspy.NewPrometheusConfig("node[^0]")),
    // WASP generators
    benchspy.WithGenerators(gen),
)
require.NoError(t, err, "failed to create the report")
```

## Custom metrics
### WASP Generator
Since `WASP` stores AUT's responses in each generator you can create custom metrics that leverage them. Here's an example
of adding a function that returns the number of responses that timed out:
```go
var generator *wasp.Generator

var timeouts = func(responses *wasp.SliceBuffer[wasp.Response]) (float64, error) {
    if len(responses.Data) == 0 {
        return 0, nil
    }

    timeoutCount := 0.0
    inTimeCount := 0.0
    for _, response := range responses.Data {
        if response.Timeout {
            timeoutCount = timeoutCount + 1
        } else {
            inTimeCount = inTimeCount + 1
        }
    }

    return timeoutCount / (timeoutCount + inTimeCount), nil
}

directExectutor, err := NewDirectQueryExecutor(generator, map[string]DirectQueryFn{
    "timeout_ratio": timeouts,
})
require.NoError(t, err, "failed to create Direct Query Executor")
```

### Loki
Using custom `LogQL` queries is even simpler as all you need to do is create a new instance of
`NewLokiQueryExecutor` with a map of desired queries.
```go
var generator *wasp.Generator

lokiQueryExecutor := benchspy.NewLokiQueryExecutor(
    map[string]string{
        "responses_over_time": fmt.Sprintf("sum(count_over_time({my_label=~\"%s\", test_data_type=~\"responses\", gen_name=~\"%s\"} [1s])) by (node_id, go_test_name, gen_name)", label, gen.Cfg.GenName),
    },
    generator.Cfg.LokiConfig,
)
```
> [!NOTE]
> In order to effectively write `LogQL` queries for WASP you need to be familar with how to label
> your generators and what `test_data_types` WASP uses.

### Prometheus
Adding custom `PromQL` queries is equally straight-forward:
```go
promConfig := benchspy.NewPrometheusConfig()

prometheusExecutor, err := benchspy.NewPrometheusQueryExecutor(
    map[string]string{
        "cpu_rate_by_container": "rate(container_cpu_usage_seconds_total{name=~\"chainlink.*\"}[5m])[30m:1m]",
    },
    *promConfig,
)
require.NoError(t, err)
```

### How to use with StandardReport
Using custom queries with a `StandardReport` is rather simple. Instead of passing `StandardQueryExecutorType` with the
functional option `WithStandardQueries` you should pass the `QueryExecutors` created above with `WithQueryExecutors` option:
```go
report, err := benchspy.NewStandardReport(
    "2d1fa3532656c51991c0212afce5f80d2914e34e",
    benchspy.WithQueryExecutors(directExectutor, lokiQueryExecutor, prometheusExecutor),
    benchspy.WithGenerators(gen),
)
require.NoError(t, err, "failed to create baseline report")