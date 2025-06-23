# Observability Guide

Observability is a core feature of the Chainlink Testing Framework (CTF). The framework provides an integrated stack for monitoring, logging, tracing, and profiling your tests and environments.

## Why Observability?

- **Debugging**: Quickly identify issues in your tests or deployed components
- **Performance Analysis**: Monitor resource usage, latency, and throughput
- **Reliability**: Detect failures, bottlenecks, and regressions
- **Transparency**: Gain insight into system behavior during tests

## Observability Stack

CTF integrates with the following tools:

- **Grafana**: Dashboards and visualization
- **Loki**: Log aggregation
- **Prometheus**: Metrics collection
- **Jaeger**: Distributed tracing
- **Pyroscope**: Performance profiling

## Quick Start

### 1. Spin Up the Observability Stack

Use the CTF CLI to start the observability stack locally:

```bash
ctf obs up
```

This will launch Grafana, Loki, Prometheus, and Pyroscope (if configured) in Docker containers.

### 2. Access Grafana

- Open [http://localhost:3000](http://localhost:3000) in your browser
- Default credentials: `admin` / `admin` (change after first login)
- Explore pre-built dashboards for blockchain, Chainlink nodes, and test metrics

### 3. Log and Metric Collection

- All logs from components are sent to Loki
- Metrics are scraped by Prometheus
- Custom metrics and logs can be pushed from your tests

## Configuration

You can configure observability endpoints in your TOML or via environment variables:

```toml
[observability]
  grafana_url = "http://localhost:3000"
  loki_url = "http://localhost:3100"
  prometheus_url = "http://localhost:9090"
  pyroscope_url = "http://localhost:4040"
```

Or via environment variables:

```bash
export CTF_GRAFANA_URL=http://localhost:3000
export CTF_LOKI_URL=http://localhost:3100
export CTF_PROMETHEUS_URL=http://localhost:9090
export CTF_PYROSCOPE_URL=http://localhost:4040
```

## Using Observability in Tests

### Logging

- Use the built-in logger (`framework/logging`) for structured logs
- All logs are automatically sent to Loki
- You can add custom log fields for better filtering in Grafana

### Metrics

- Expose custom metrics from your components or tests
- Use Prometheus client libraries to define and push metrics
- Metrics are visualized in Grafana dashboards

### Tracing

- Distributed tracing is enabled for supported components
- Use Jaeger UI to view traces ([http://localhost:16686](http://localhost:16686))
- Trace requests across services and components

### Profiling

- Pyroscope collects CPU and memory profiles
- Access Pyroscope UI at [http://localhost:4040](http://localhost:4040)
- Analyze performance bottlenecks

## Grafana Dashboards

- Pre-built dashboards for blockchain, Chainlink nodes, and test metrics
- Custom dashboards can be created for your use case
- Use dashboard annotations to mark test events (e.g., chaos experiments, load tests)

### Example: Annotating Dashboards

```go
import (
    "github.com/smartcontractkit/chainlink-testing-framework/wasp"
)

profile := wasp.NewProfile()
grafanaOpts := &wasp.GrafanaOpts{
    GrafanaURL: "http://localhost:3000",
    GrafanaToken: "your-token",
    AnnotateDashboardUIDs: []string{"dashboard-uid"},
}
profile.WithGrafana(grafanaOpts)
```

## Best Practices

- **Always enable observability** in CI and local runs
- **Tag logs and metrics** with test names, component names, and environment
- **Use dashboard annotations** for key events (start/end of tests, chaos, upgrades)
- **Monitor resource usage** to detect leaks or bottlenecks
- **Set up alerts** in Grafana for critical metrics (e.g., error rates, latency)

## Troubleshooting

- **No logs in Grafana**: Check Loki container status and log configuration
- **No metrics in dashboards**: Ensure Prometheus is scraping targets
- **No traces in Jaeger**: Verify tracing is enabled and endpoints are correct
- **Pyroscope not collecting**: Check agent configuration and endpoint

## Further Reading
- [Component System](Components)
- [Configuration Guide](Configuration)
- [Caching Guide](Caching)
- [Test Patterns](Test-Patterns)
- [WASP Load Testing](../Libraries/WASP)
- [Havoc Chaos Testing](../Libraries/Havoc) 