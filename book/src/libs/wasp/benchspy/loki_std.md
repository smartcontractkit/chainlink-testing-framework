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
    "v1.0.0",
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
allCurrentAsStringSlice := benchspy.MustAllLokiResults(currentReport)
allPreviousAsStringSlice := benchspy.MustAllLokiResults(previousReport)

require.NotEmpty(t, allCurrentAsStringSlice, "current report is empty")
require.NotEmpty(t, allPreviousAsStringSlice, "previous report is empty")

currentAsStringSlice := allCurrentAsStringSlice[gen.Cfg.GenName]
previousAsStringSlice := allPreviousAsStringSlice[gen.Cfg.GenName]
```

An explanation is needed here: this function separates metrics for each generator, hence it returns a `map[string]map[string][]string`. Let's break it down:
- outer map's key is generator name
- inner map's key is metric name and the value is a series of measurements
In our case there's only a single generator, but in a complex test there might be a few.

## Step 4: Compare Metrics

Now, let’s compare metrics. Since we have `[]string`, we’ll first convert it to `[]float64`, calculate the median, and ensure the difference between the averages is less than 1%. Again, this is just an example—you should decide the best way to validate your metrics. Here we are explicitly aggregating them using an average to get a single number representation of each metric, but for your case a median or percentile or yet some other aggregate might be more appropriate.

```go
var compareAverages = func(t *testing.T, metricName string, currentAsStringSlice, previousAsStringSlice map[string][]string, maxPrecentageDiff float64) {
	require.NotEmpty(t, currentAsStringSlice[metricName], "%s results were missing from current report", metricName)
	require.NotEmpty(t, previousAsStringSlice[metricName], "%s results were missing from previous report", metricName)

	currentFloatSlice, err := benchspy.StringSliceToFloat64Slice(currentAsStringSlice[metricName])
	require.NoError(t, err, "failed to convert %s results to float64 slice", metricName)
	currentMedian, err := stats.Mean(currentFloatSlice)
	require.NoError(t, err, "failed to calculate median for %s results", metricName)

	previousFloatSlice, err := benchspy.StringSliceToFloat64Slice(previousAsStringSlice[metricName])
	require.NoError(t, err, "failed to convert %s results to float64 slice", metricName)
	previousMedian, err := stats.Mean(previousFloatSlice)
	require.NoError(t, err, "failed to calculate median for %s results", metricName)

	var diffPrecentage float64
	if previousMedian != 0.0 && currentMedian != 0.0 {
		diffPrecentage = (currentMedian - previousMedian) / previousMedian * 100
	} else if previousMedian == 0.0 && currentMedian == 0.0 {
		diffPrecentage = 0.0
	} else {
		diffPrecentage = 100.0
	}
	assert.LessOrEqual(t, math.Abs(diffPrecentage), maxPrecentageDiff, "%s medians are more than 1% different", metricName, fmt.Sprintf("%.4f", diffPrecentage))
}

compareAverages(
    t,
    string(benchspy.MedianLatency),
    currentAsStringSlice,
    previousAsStringSlice,
    1.0,
)
compareAverages(t, string(benchspy.Percentile95Latency), currentAsStringSlice, previousAsStringSlice, 1.0)
compareAverages(t, string(benchspy.MaxLatency), currentAsStringSlice, previousAsStringSlice, 1.0)
compareAverages(t, string(benchspy.ErrorRate), currentAsStringSlice, previousAsStringSlice, 1.0)
```

> [!WARNING]
> Standard Loki metrics are all calculated using a 10 seconds moving window, which results in smoothing of values due to aggregation.
> To learn what that means in details, please refer to [To Loki or Not to Loki](./loki_dillema.md) chapter.
>
> Also, due to the HTTP API endpoint used, namely the `query_range`, all query results **are always returned as a slice**. Execution of **instant queries**
> that return a single data point is currently **not supported**.

## What’s Next?

In this example, we used standard metrics, which are the same as in the first test. Now, [let’s explore how to use your custom LogQL queries](./loki_custom.md).

> [!NOTE]
> You can find the full example [here](https://github.com/smartcontractkit/chainlink-testing-framework/tree/main/wasp/examples/benchspy/loki_query_executor/loki_query_executor_test.go).