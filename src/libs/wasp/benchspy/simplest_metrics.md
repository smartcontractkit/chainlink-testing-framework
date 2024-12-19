# BenchSpy - Simplest Metrics

As mentioned earlier, `BenchSpy` doesn't include any built-in comparison logic. It's up to you to decide how to compare metrics, as there are various ways to approach it and different data formats returned by queries.

For example, if your query returns a time series, you could:
- Compare each data point in the time series individually.
- Compare aggregates like averages, medians, or min/max values of the time series.

## Working with Built-in `QueryExecutors`
Each built-in `QueryExecutor` returns a different data type, and we use the `interface{}` type to reflect this. Since `Direct` executor always returns `float64` we have added a convenience function
that checks whether any of the standard metrics has **degraded** more than the threshold. If the performance has improved, no error will be returned.

```go
hasErrors, errors := benchspy.CompareDirectWithThresholds(
    // maximum differences in percentages for:
    1.0, // median latency
    1.0, // p95 latency
    1.0, // max latency
    1.0, // error rate
    currentReport,
    previousReport,
)
require.False(t, hasErrors, fmt.Sprintf("errors found: %v", errors))
```

If there are errors they will be returned as `map[string][]errors`, where key is the name of a generator.

> [!NOTE]
> Both `Direct` and `Loki` query executors support following standard performance metrics out of the box:
> - `median_latency`
> - `p95_latency`
> - `max_latency`
> - `error_rate`

The function also prints a table with the differences between two reports, regardless whether they were meaningful:
```bash
Generator: vu1
==============
+-------------------------+---------+---------+---------+
|         METRIC          |   V1    |   V2    | DIFF %  |
+-------------------------+---------+---------+---------+
| median_latency          | 50.1300 | 50.1179 | -0.0242 |
+-------------------------+---------+---------+---------+
| 95th_percentile_latency | 50.7387 | 50.7622 | 0.0463  |
+-------------------------+---------+---------+---------+
| max_latency             | 55.7195 | 51.7248 | -7.1692 |
+-------------------------+---------+---------+---------+
| error_rate              | 0.0000  | 0.0000  | 0.0000  |
+-------------------------+---------+---------+---------+
```

## Wrapping Up

And that's it! You've written your first test that uses `WASP` to generate load and `BenchSpy` to ensure that the median latency, 95th percentile latency, max latency and error rate haven't changed significantly between runs. You accomplished this without even needing a Loki instance. But what if you wanted to leverage the power of `LogQL`? We'll explore that in the [next chapter](./loki_std.md).

> [!NOTE]
> You can find the full example [here](https://github.com/smartcontractkit/chainlink-testing-framework/tree/main/wasp/examples/benchspy/direct_query_executor/direct_query_executor_test.go).