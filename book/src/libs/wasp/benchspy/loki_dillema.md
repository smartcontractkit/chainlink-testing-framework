# BenchSpy - To Loki or Not to Loki?

You might be wondering whether to use the `Loki` or `Direct` query executor if all you need are basic latency metrics.

## Rule of Thumb

If all you need is a single number, such as the median latency or error rate, and you're not interested in:
- Comparing time series directly,
- Examining minimum or maximum values, or
- Performing advanced calculations on raw data,

then you should opt for the `Direct` query executor.

## Why Choose `Direct`?

The `Direct` executor returns a single value for each standard metric using the same raw data that Loki would use. It accesses data stored in the `WASP` generator, which is later pushed to Loki.

This means you can:
- Run your load test without a Loki instance.
- Avoid calculating metrics like the median, 95th percentile latency, or error ratio yourself.

By using `Direct`, you save resources and simplify the process when advanced analysis isn't required.
