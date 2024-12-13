# BenchSpy - Standard Loki Metrics

> [!WARNING]
> This example assumes you have access to Loki and Grafana instances. If you don't, learn how to launch them using CTFv2's [observability stack](../../../framework/observability/observability_stack.md).

In this example, our Loki workflow will differ from the previous one in just a few details:
- The generator will include a Loki configuration.
- The standard query executor type will be `benchspy.StandardQueryExecutor_Loki`.
- All results will be cast to `[]string`.
- We'll calculate medians for all metrics.

Ready?

## Step 1: Define a New Load Generator

Let's start by defining a new load generator:

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

## Step 2: Run the Generator and Save the Baseline Report

```go
gen.Run(true)

baseLineReport, err := benchspy.NewStandardReport(
    "c2cf545d733eef8bad51d685fcb302e277d7ca14",
    // notice the different standard query executor type
    benchspy.WithStandardQueries(benchspy.StandardQueryExecutor_Loki),
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

## Step 3: Skip to Metrics Comparison

Since the next steps are very similar to those in the first test, we’ll skip them and go straight to metrics comparison.

By default, the `LokiQueryExecutor` returns results as the `[]string` data type. Let’s use dedicated convenience functions to cast them from `interface{}` to string slices:

```go
currentAsStringSlice := benchspy.MustAllLokiResults(currentReport)
previousAsStringSlice := benchspy.MustAllLokiResults(previousReport)
```

## Step 4: Compare Metrics

Now, let’s compare metrics. Since we have `[]string`, we’ll first convert it to `[]float64`, calculate the median, and ensure the difference between the medians is less than 1%. Again, this is just an example—you should decide the best way to validate your metrics.

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

    var diffPercentage float64
	if previousMedian != 0.0 && currentMedian != 0.0 {
		diffPrecentage = (currentMedian - previousMedian) / previousMedian * 100
	} else if previousMedian == 0.0 && currentMedian == 0.0 {
		diffPrecentage = 0.0
	} else {
		diffPrecentage = 100.0
	}
    assert.LessOrEqual(t, math.Abs(diffPercentage), 1.0, "%s medians are more than 1% different", metricName, fmt.Sprintf("%.4f", diffPercentage))
}

compareMedian(string(benchspy.MedianLatency))
compareMedian(string(benchspy.Percentile95Latency))
compareMedian(string(benchspy.ErrorRate))
```

> [!WARNING]
> Standard Loki metrics are all calculated using a 10 seconds moving window, which results in smoothing of values due to aggregation.
> To learn what that means in details, please refer to [To Loki or Not to Loki](./loki_dillema.md) chapter.

## What’s Next?

In this example, we used standard metrics, which are the same as in the first test. Now, [let’s explore how to use your custom LogQL queries](./loki_custom.md).

> [!NOTE]
> You can find the full example [here](https://github.com/smartcontractkit/chainlink-testing-framework/tree/main/wasp/examples/benchspy/loki_query_executor/loki_query_executor_test.go).