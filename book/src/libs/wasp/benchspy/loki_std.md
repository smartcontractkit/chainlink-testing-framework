# BenchSpy - Standard Loki metrics

> [!NOTE]
> This example assumes you have access to Loki and Grafana instances. If you don't
> find out how to launch them using CTFv2's [observability stack](../../../framework/observability/observability_stack.md).

Our Loki example, will vary from the previous one in just a couple of details:
* generator will have Loki config
* standard query executor type will be `benchspy.StandardQueryExecutor_Loki`
* we will cast all results to `[]string`
* and calculate medians for all metrics

Ready?

Let's define new load generation first:
```go
label := "benchspy-std"

gen, err := wasp.NewGenerator(&wasp.Config{
    T:          t,
    // read Loki config from environment
    LokiConfig: wasp.NewEnvLokiConfig(),
    GenName:    "vu",
    // set unique labels
    Labels: map[string]string{
        "branch": label,
        "commit": label,
    },
    CallTimeout: 100 * time.Millisecond,
    LoadType:    wasp.VU,
    Schedule:    wasp.Plain(10, 15*time.Second),
    VU: wasp.NewMockVU(&wasp.MockVirtualUserConfig{
        CallSleep: 50 * time.Millisecond,
    }),
})
require.NoError(t, err)
```

Now let's run the generator and save baseline report:
```go
gen.Run(true)

fetchCtx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
defer cancelFn()

baseLineReport, err := benchspy.NewStandardReport(
    "c2cf545d733eef8bad51d685fcb302e277d7ca14",
    // notice the different standard executor type
    benchspy.WithStandardQueries(benchspy.StandardQueryExecutor_Loki),
    benchspy.WithGenerators(gen),
)
require.NoError(t, err, "failed to create original report")

fetchErr := baseLineReport.FetchData(fetchCtx)
require.NoError(t, fetchErr, "failed to fetch data for original report")

path, storeErr := baseLineReport.Store()
require.NoError(t, storeErr, "failed to store current report", path)
```

Since next steps are very similar to the ones used in the first test we will skip them and jump straight
to metrics comparison.

By default, `LokiQueryExecutor` returns `[]string` data type, so let's use dedicated convenience functions
to cast them from `interface{}` to string slice:
```go
currentAsStringSlice := benchspy.MustAllLokiResults(currentReport)
previousAsStringSlice := benchspy.MustAllLokiResults(previousReport)
```

And finally, time to compare metrics. Since we have a `[]string` we will first convert it to `[]float64` and
then calculate a median and assume it hasn't changed by more than 1%. Again, remember that this is just an illustration.
You should decide yourself what's the best way to assert the metrics.

```go
var compareMedian = func(metricName string) {
    require.NotEmpty(t, currentAsStringSlice[metricName], "%s results were missing from current report", metricName)
    require.NotEmpty(t, previousAsStringSlice[metricName], "%s results were missing from previous report", metricName)

    currentFloatSlice, err := benchspy.StringSliceToFloat64Slice(currentAsStringSlice[metricName])
    require.NoError(t, err, "failed to convert %s results to float64 slice", metricName)
    currentMedian := benchspy.CalculatePercentile(currentFloatSlice, 0.5)

    previousFloatSlice, err := benchspy.StringSliceToFloat64Slice(previousAsStringSlice[metricName])
    require.NoError(t, err, "failed to convert %s results to float64 slice", metricName)
    previousMedian := benchspy.CalculatePercentile(previousFloatSlice, 0.5)

    var diffPrecentage float64
    if previousMedian != 0 {
        diffPrecentage = (currentMedian - previousMedian) / previousMedian * 100
    } else {
        diffPrecentage = currentMedian * 100
    }
    assert.LessOrEqual(t, math.Abs(diffPrecentage), 1.0, "%s medians are more than 1% different", metricName, fmt.Sprintf("%.4f", diffPrecentage))
}

compareMedian(string(benchspy.MedianLatency))
compareMedian(string(benchspy.Percentile95Latency))
compareMedian(string(benchspy.ErrorRate))
```

We have used standard metrics, which are the same as in the first test, now let's see how you can use your custom LogQl queries.

> [!NOTE]
> Don't know whether to use `Loki` or `Direct` query executors? [Read this!](./loki_dillema.md)

You can find the full example [here](...).