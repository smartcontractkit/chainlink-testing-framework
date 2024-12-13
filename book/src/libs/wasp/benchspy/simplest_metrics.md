# BenchSpy - Simplest metrics

As mentioned earlier, `BenchSpy` doesn't come with any comparison logic. It's up to the user to decide how to comapre metrics,
since there are not only various ways of doing it, but also various data formats returned by queries.

For example, if your query returns a time series, when comparing them you could:
* compare each time series data point individually
* compare averages, medians or min/max of each time series

Each of approaches as it's pros and cons and `BenchSpy` makes no judegment here. In our example we will use a very simplified
approach, which **should not be treated** as the golden standard. Our `QueryExecutor` returns a single data point for each metric,
so there's no dilemma, but with `Loki` and `Prometheus` things get a bit more complicated.

But first... since each of built-in `QueryExecutors` returns a different data type and we are using `interface{}` type to reflect that,
we will use convenience functions to cast them to more usable types:
```go
currentAsFloat64 := benchspy.MustAllDirectResults(currentReport)
previousAsloat64 := benchspy.MustAllDirectResults(previousReport)
```

> [!NOTE]
> All of the standard metrics for `GeneratorQueryExecutor` have `float64` type.

Now, let's define a simple function that compares two floats and makes sure that the difference between them is smaller than 1%.
```go
var compareValues = func(
    metricName string,
    maxDiffPercentage float64,
) {
    require.NotNil(t, currentAsFloat64[metricName], "%s results were missing from current report", metricName)
    require.NotNil(t, previousAsloat64[metricName], "%s results were missing from previous report", metricName)

    currentMetric := currentAsFloat64[metricName]
    previousMetric := previousAsloat64[metricName]

    var diffPrecentage float64
    if previousMetric != 0.0 {
        diffPrecentage = (currentMetric - previousMetric) / previousMetric * 100
    } else {
        diffPrecentage = currentMetric * 100.0
    }
    assert.LessOrEqual(t, math.Abs(diffPrecentage), maxDiffPercentage, "%s medians are more than 1% different", metricName, fmt.Sprintf("%.4f", diffPrecentage))
}

compareValues(string(benchspy.MedianLatency), 1.0)
compareValues(string(benchspy.Percentile95Latency), 1.0)
compareValues(string(benchspy.ErrorRate), 1.0)
```

And that's it! You have written your first test that uses `WASP` to generate the load and `BenchSpy` to make sure that neither median latency nor 95th latency percentile
nor error rate has changed significantly between the runs. You did that without even needing a Loki instance, but what if you wanted to leverage the power
of `LogQL`? We will look at that in the [next chapter](./using_loki.md).

You can find the full example [here](...).