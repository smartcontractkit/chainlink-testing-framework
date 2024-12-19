# BenchSpy

BenchSpy (short for Benchmark Spy) is a [WASP](../overview.md)-coupled tool designed for easy comparison of various performance metrics.

## Key Features
- **Three built-in data sources**:
  - `Loki`
  - `Prometheus`
  - `Direct`
- **Standard/pre-defined metrics** for each data source.
- **Ease of extensibility** with custom metrics.
- **Ability to load the latest performance report** based on Git history.

BenchSpy does not include any built-in comparison logic beyond ensuring that performance reports are comparable (e.g., they measure the same metrics in the same way), offering complete freedom to the user for interpretation and analysis.

## Why could you need it?
`BenchSpy` was created with two main goals in mind:
* **measuring application performance programmatically**,
* **finding performance-related changes or regression issues between different commits or releases**.