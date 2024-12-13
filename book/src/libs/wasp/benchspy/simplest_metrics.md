# BenchSpy - Simplest Metrics

As mentioned earlier, `BenchSpy` doesn't include any built-in comparison logic. It's up to you to decide how to compare metrics, as there are various ways to approach it and different data formats returned by queries.

For example, if your query returns a time series, you could:
- Compare each data point in the time series individually.
- Compare aggregates like averages, medians, or min/max values of the time series.

Each of these approaches has its pros and cons, and `BenchSpy` doesn't make any judgments here. In this example, we'll use a very simplified approach, which **should not be treated** as a gold standard. In our case, the `QueryExecutor` returns a single data point for each metric, eliminating the complexity. However, with `Loki` and `Prometheus`, things can get more complicated.

## Working with Built-in `QueryExecutors`

Since each built-in `QueryExecutor` returns a different data type, and we use the `interface{}` type to reflect this, convenience functions help cast these results into more usable types:

```go
currentAsFloat64 := benchspy.MustAllDirectResults(currentReport)
previousAsFloat64 := benchspy.MustAllDirectResults(previousReport)
```

> [!NOTE]
> All standard metrics for the `DirectQueryExecutor` have the `float64` type.

## Defining a Comparison Function

Next, let's define a simple function to compare two floats and ensure the difference between them is smaller than 1%:

```go
var compareValues = func(
    metricName string,
    maxDiffPercentage float64,
) {
    require.NotNil(t, currentAsFloat64[metricName], "%s results were missing from current report", metricName)
    require.NotNil(t, previousAsFloat64[metricName], "%s results were missing from previous report", metricName)

    currentMetric := currentAsFloat64[metricName]
    previousMetric := previousAsFloat64[metricName]

    var diffPercentage float64
    if previousMetric != 0.0 {
        diffPercentage = (currentMetric - previousMetric) / previousMetric * 100
    } else {
        diffPercentage = 100.0
    }
    assert.LessOrEqual(t, math.Abs(diffPercentage), maxDiffPercentage, "%s medians are more than 1% different", metricName, fmt.Sprintf("%.4f", diffPercentage))
}

compareValues(string(benchspy.MedianLatency), 1.0)
compareValues(string(benchspy.Percentile95Latency), 1.0)
compareValues(string(benchspy.ErrorRate), 1.0)
```

## Wrapping Up

And that's it! You've written your first test that uses `WASP` to generate load and `BenchSpy` to ensure that the median latency, 95th percentile latency, and error rate haven't changed significantly between runs. You accomplished this without even needing a Loki instance. But what if you wanted to leverage the power of `LogQL`? We'll explore that in the [next chapter](./loki_std.md).

> [!NOTE]
> You can find the full example [here](https://github.com/smartcontractkit/chainlink-testing-framework/tree/main/wasp/examples/benchspy/direct_query_executor/direct_query_executor_test.go).