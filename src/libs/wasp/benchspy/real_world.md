# BenchSpy - Real-World Example

Now that we've covered all possible usages, you might wonder how to write a test that compares performance between different releases of your application. Here’s a practical example to guide you through the process.

## Typical Workflow

1. Write a performance test.
2. At the end of the test, generate a performance report, store it, and commit it to Git.
3. Modify the test to fetch both the latest report and create a new one.
4. Write assertions to validate your performance metrics.

---

## Writing the Performance Test

We'll use a simple mock for the application under test. This mock waits for `50 ms` before returning a 200 response code.

```go
generator, err := wasp.NewGenerator(&wasp.Config{
    T:           t,
    GenName:     "vu",
    CallTimeout: 100 * time.Millisecond,
    LoadType:    wasp.VU,
    Schedule:    wasp.Plain(10, 15*time.Second),
    VU: wasp.NewMockVU(&wasp.MockVirtualUserConfig{
        CallSleep: 50 * time.Millisecond,
    }),
})
require.NoError(t, err)
generator.Run(true)
```

---

## Generating the First Report

Here, we'll generate a performance report for version `v1.0.0` using the `Direct` query executor. The report will be saved to a custom directory named `test_reports`. This report will later be used to compare the performance of new versions.

```go
fetchCtx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
defer cancelFn()

baseLineReport, err := benchspy.NewStandardReport(
    "v1.0.0",
    benchspy.WithStandardQueries(benchspy.StandardQueryExecutor_Direct),
    benchspy.WithReportDirectory("test_reports"),
    benchspy.WithGenerators(gen),
)
require.NoError(t, err, "failed to create baseline report")

fetchErr := baseLineReport.FetchData(fetchCtx)
require.NoError(t, fetchErr, "failed to fetch data for original report")

path, storeErr := baseLineReport.Store()
require.NoError(t, storeErr, "failed to store current report", path)
```

---

## Modifying Report Generation

With the baseline report for `v1.0.0` stored, we'll modify the test to support future releases. The code from the previous step will change as follows:

```go
currentVersion := os.Getenv("CURRENT_VERSION")
require.NotEmpty(t, currentVersion, "No current version provided")

fetchCtx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
defer cancelFn()

currentReport, previousReport, err := benchspy.FetchNewStandardReportAndLoadLatestPrevious(
    fetchCtx,
    currentVersion,
    benchspy.WithStandardQueries(benchspy.StandardQueryExecutor_Direct),
    benchspy.WithReportDirectory("test_reports"),
    benchspy.WithGenerators(generator),
)
require.NoError(t, err, "failed to fetch current report or load the previous one")
```

This function fetches the current report (for version passed as environment variable `CURRENT_VERSION`) while loading the latest stored report from the `test_reports` directory.

---

## Adding Assertions

Let’s assume you want to ensure that none of the performance metrics degrade by more than **1%** between releases (and that error rate has not changed at all). Here's how you can write assertions using a convenient function for the `Direct` query executor:

```go
hasFailed, error := benchspy.CompareDirectWithThresholds(
    1.0, // Max 1% worse median latency
    1.0, // Max 1% worse p95 latency
    1.0, // Max 1% worse maximum latency
    0.0, // No increase in error rate
    currentReport, previousReport)
require.False(t, hasError, fmt.Sprintf("issues found: %v", error))
```

Error returned by this function is a concatenation of all threshold violations found for each standard metric and generator.

---

## Conclusion

You’re now ready to use `BenchSpy` to ensure that your application’s performance does not degrade below your specified thresholds!

> [!NOTE]
> [Here](https://github.com/smartcontractkit/chainlink-testing-framework/tree/main/wasp/examples/benchspy/direct_query_executor/direct_query_real_case.go) you can find an example test where performance has degraded significantly,
> because mock's latency has been increased from `50ms` to `60ms`.
>
> **This test passes because it is designed to expect performance degradation. Of course, in a real application, your goal should be to prevent such regressions.** 😊