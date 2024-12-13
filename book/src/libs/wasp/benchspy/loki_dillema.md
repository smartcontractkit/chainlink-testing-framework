# BenchSpy - To Loki or not to Loki?

You might be asking yourself whether you should use `Loki` or `Direct` query executor if all you
need are basic latency metrics.

As a rule of thumb, if all you need is a single number that describes the median latency or error rate
and you are not interested in directly comparing time series, minimum or maximum values or any kinds
of more advanced calculation on raw data, then you should go with the `Direct`.

Why?

Because it returns a single value for each of standard metrics using the same raw data that Loki would use
(it accesses the data stored in the `WASP`'s generator that would later be pushed to Loki).
This way you can run your load test without a Loki instance and save yourself the need of calculating the
median and 95th percentile latency or the error ratio.