# BenchSpy

BenchSpy (short for benchmark spy) is a WASP-coupled tool that allows for easy comparison of various performance metrics.

It's main characteristics are:
* three built-in data sources:
    * `Loki`
    * `Prometheus`
    * `Direct`
* standard/pre-defined metrics for each data source
* ease of extensibility with custom metrics
* ability to load latest performance report based on Git history
* 88% unit test coverage

It doesn't come with any comparation logic, other than making sure that performance reports are comparable (e.g. they mesure the same metrics in the same way),
leaving total freedom to the user.