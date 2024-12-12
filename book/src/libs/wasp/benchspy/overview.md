# BenchSpy

BenchSpy (short for benchmark spy) is a WASP-coupled tool that allows for easy comparison of various performance metrics.
It supports three types of data sources:
* `Loki`
* `Prometheus`
* `WASP generators`

And can be easily extended to support additional ones.

Since it's main goal is comparison of performance between various releases or versions of applications (for example, to catch performance degradation)
it is `Git`-aware and is able to automatically find the latest relevant performance report.

It doesn't come with any comparation logic, other than making sure that performance reports are comparable (e.g. they mesure the same metrics in the same way),
leaving total freedom to the user.