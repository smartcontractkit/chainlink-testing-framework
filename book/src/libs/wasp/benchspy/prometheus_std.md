# BenchSpy - Standard Prometheus metrics

Now that we have seen how we can query and assert on the load side of things, let's check how we can
query and assert on resource usage by our Application Under Test (AUT).

If you don't know why that's important think about following situation: p95 latency of the new release of your application is the same as the previous
one, but memory consumption is 34% higher. Not ideal, right? And how could you know that it's even a thing?

To begin, since `WASP` has no built-in integration with `Prometheus` we need to pass its configuration separately:
```go
promConfig := benchspy.NewPrometheusConfig("node[^0]")
```

This constructor loads the url from environment variable `PROMETHEUS_URL` and adds a single regexp pattern that will be used to match container **by name**.
In this very case it will exclude the bootstrap Chainlink node.

> [!NOTE]
> This example assumes that you have both the observability stack and basic node set running.
> If you have the `CTF CLI` you can start it by running: `ctf b ns`

Just like in previous examples we will use built-in Prometheus metrics and fetch and store a baseline report:
```go
baseLineReport, err := benchspy.NewStandardReport(
    "91ee9e3c903d52de12f3d0c1a07ac3c2a6d141fb",
    benchspy.WithQueryExecutorType(benchspy.StandardQueryExecutor_Prometheus),
    benchspy.WithPrometheusConfig(promConfig),
    // needed even if we don't query Loki or Generator,
    // because we calculate test time range based on
    // generator start/end times
    benchspy.WithGenerators(gen),
)

fetchCtx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
defer cancelFn()

fetchErr := baseLineReport.FetchData(fetchCtx)
require.NoError(t, fetchErr, "failed to fetch current report")

path, storeErr := baseLineReport.Store()
require.NoError(t, storeErr, "failed to store current report", path)
```

> [!NOTE]
> Standard metrics for Prometheus are different from the ones used by Loki or Generator query executors.
> That's because with Prometheus we are interested in resource usage by the AUT and with Loki/Generator
> with load characteristics metrics.
>
> These standard Prometheus metrics include:
> * `median_cpu_usage`
> * `median_mem_usage`
> * `p95_cpu_usage`
> * `p95_mem_usage`
>
> And are calculated on the **container level** based on total usages (user + system).

Contrary to Loki and Generator, Prometheus uses a variety of data types for results, such as:
* `scalar`
* `string`
* `vector`
* `matrix`

Therefore, asserting on the results is a bit more complex.

First, we need to conver them to `mode.Value` interface using convenience functions:
```go
currentAsValues := benchspy.MustAllPrometheusResults(currentReport)
previousAsValues := benchspy.MustAllPrometheusResults(previousReport)
```

Then, once we know which data type our query returned we need to cast it to the specific type:
```go
// fetch a single metric
currentMedianCPUUsage := currentAsValues[string(benchspy.MedianCPUUsage)]
previousMedianCPUUsage := previousAsValues[string(benchspy.MedianCPUUsage)]

assert.Equal(t, currentMedianCPUUsage.Type(), previousMedianCPUUsage.Type(), "types of metrics should be the same")

// in this case we know that this query returns a Vector
currentMedianCPUUsageVector := currentMedianCPUUsage.(model.Vector)
previousMedianCPUUsageVector := previousMedianCPUUsage.(model.Vector)
```

> [!WARNING]
> All of standard Prometheus metrics that are bundled with `BenchSpy` return `model.Vector`,
> but if you decide to use your own queries you need to manually their return types.

We skip the assertion part, because unless you are comparing resource usage over periods with stable loads
it's almost certain that there would be big differences between two reports (e.g. when the first report is created
node set, presumably, has just started; when the second one is created it has already been up for some time).

In the next chapter we will look at custom Prometheus queries.

You will find the full example [here](...).