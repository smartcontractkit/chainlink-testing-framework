# BenchSpy - To Loki or Not to Loki?

You might be wondering whether to use the `Loki` or `Direct` query executor if all you need are basic latency metrics.

## Rule of Thumb

You should opt for the `Direct` query executor if all you need is a single number, such as the median latency or error rate, and you're not interested in:
- Comparing time series directly,
- Examining minimum or maximum values over time, or
- Performing advanced calculations on raw data,

## Why Choose `Direct`?

The `Direct` executor returns a single value for each standard metric using the same raw data that Loki would use. It accesses data stored in the `WASP` generator, which is later pushed to Loki.

This means you can:
- Run your load test without a Loki instance.
- Avoid calculating metrics like the median, 95th percentile latency, or error ratio yourself.

By using `Direct`, you save resources and simplify the process when advanced analysis isn't required.

> [!WARNING]
> Metrics calculated by the two query executors may differ slightly due to differences in their data processing and calculation methods:
> - **`Direct` QueryExecutor**: This method processes all individual data points from the raw dataset, ensuring that every value is taken into account for calculations like averages, percentiles, or other statistics. It provides the most granular and precise results but may also be more sensitive to outliers and noise in the data.
> - **`Loki` QueryExecutor**: This method aggregates data using a default window size of 10 seconds. Within each window, multiple raw data points are combined (e.g., through averaging, summing, or other aggregation functions), which reduces the granularity of the dataset. While this approach can improve performance and reduce noise, it also smooths the data, which may obscure outliers or small-scale variability.

> #### Why This Matters for Percentiles:
> Percentiles, such as the 95th percentile (p95), are particularly sensitive to the granularity of the input data:
> - In the **`Direct` QueryExecutor**, the p95 is calculated across all raw data points, capturing the true variability of the dataset, including any extreme values or spikes.
> - In the **`Loki` QueryExecutor**, the p95 is calculated over aggregated data (i.e. using the 10-second window). As a result, the raw values within each window are smoothed into a single representative value, potentially lowering or altering the calculated p95. For example, an outlier that would significantly affect the p95 in the `Direct` calculation might be averaged out in the `Loki` window, leading to a slightly lower percentile value.

> #### Direct caveats:
> - **buffer limitations:** `WASP` generator use a [StringBuffer](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/wasp/buffer.go) with fixed size to store the responses. Once full capacity is reached
> oldest entries are replaced with incoming ones. The size of the buffer can be set in generator's config. By default, it is limited to 50k entries to lower resource consumption and potential OOMs.
>
> - **sampling:** `WASP` generators support optional sampling of successful responses. It is disabled by deafult, but if you do enable it, then the calculations would no longer be done over a full dataset.

> #### Key Takeaway:
> The difference arises because `Direct` prioritizes precision by using raw data, while `Loki` prioritizes efficiency and scalability by using aggregated data. When interpreting results, it’s essential to consider how the smoothing effect of `Loki` might impact the representation of variability or extremes in the dataset. This is especially important for metrics like percentiles, where such details can significantly influence the outcome.