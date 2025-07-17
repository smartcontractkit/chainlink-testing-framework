# WASP - Load Testing Library

<div align="center">

[![Go Report Card](https://goreportcard.com/badge/github.com/smartcontractkit/wasp)](https://goreportcard.com/report/github.com/smartcontractkit/wasp)
[![Component Tests](https://github.com/smartcontractkit/chainlink-testing-framework/actions/workflows/wasp-test.yml/badge.svg)](https://github.com/smartcontractkit/chainlink-testing-framework/actions/workflows/wasp-test.yml)
[![E2E tests](https://github.com/smartcontractkit/chainlink-testing-framework/actions/workflows/wasp-test-e2e.yml/badge.svg)](https://github.com/smartcontractkit/chainlink-testing-framework/actions/workflows/wasp-test-e2e.yml)
<a href='https://github.com/jpoles1/gopherbadger' target='_blank'>![gopherbadger-tag-do-not-edit](https://img.shields.io/badge/Go%20Coverage-80%25-brightgreen.svg?longCache=true&style=flat)</a>

**Scalable protocol-agnostic load testing library for Go**

</div>

## Overview

WASP is a powerful load testing library designed for Go applications that need to test the performance and scalability of their systems. It's particularly well-suited for blockchain and Chainlink applications, but can be used for any protocol or service.

## Goals

- **Easy to reuse** any custom client Go code
- **Easy to grasp** - simple API with predictable behavior
- **Slim codebase** (500-1k lines of code)
- **No test harness or CLI** - easy to integrate and run with plain `go test`
- **Predictable performance footprint** - consistent resource usage
- **Easy to create synthetic or user-based scenarios**
- **Scalable in Kubernetes** without complicated configuration
- **Non-opinionated reporting** - push any data to Loki

## Key Features

### 1. **Protocol Agnostic**
- Works with any protocol (HTTP, gRPC, WebSocket, custom protocols)
- No built-in assumptions about the target system
- Easy to integrate existing client code

### 2. **Simple API**
```go
// Basic usage
profile := wasp.NewProfile()
profile.Add(wasp.NewGenerator(config))
profile.Run(true)
```

### 3. **Flexible Configuration**
- Support for various load patterns (constant, ramp-up, step, etc.)
- Configurable RPS (requests per second) and VU (virtual users)
- Time-based and iteration-based test duration

### 4. **Observability Integration**
- Built-in Grafana dashboard support
- Loki integration for log aggregation
- Prometheus metrics export
- Custom metric collection

### 5. **Kubernetes Ready**
- Designed for distributed load testing
- No complex UI dependencies
- Easy deployment and scaling

## Installation

### Using Go Modules
```bash
go get github.com/smartcontractkit/chainlink-testing-framework/wasp
```

### Using Nix (Recommended for Development)
```bash
# Install nix
curl --proto '=https' --tlsv1.2 -sSf -L https://install.determinate.systems/nix | sh -s -- install

# Enter development shell
nix develop
```

## Quick Start

### Basic Load Test

```go
package main

import (
    "context"
    "log"
    "time"
    
    "github.com/smartcontractkit/chainlink-testing-framework/wasp"
)

func main() {
    // Create a new load test profile
    profile := wasp.NewProfile()
    
    // Define your load generator
    generator := wasp.NewGenerator(&wasp.Config{
        T: 10 * time.Second, // Test duration
        RPS: 100,            // Requests per second
        LoadType: wasp.RPS,
        Schedule: wasp.Plain(100, 10*time.Second),
    })
    
    // Add your custom request function
    generator.AddRequestFn(func(ctx context.Context) error {
        // Your custom request logic here
        // Example: HTTP request, blockchain transaction, etc.
        return nil
    })
    
    // Add generator to profile
    profile.Add(generator)
    
    // Run the load test
    _, err := profile.Run(true)
    if err != nil {
        log.Fatal(err)
    }
}
```

### HTTP Load Test

```go
package main

import (
    "context"
    "log"
    "net/http"
    "time"
    
    "github.com/smartcontractkit/chainlink-testing-framework/wasp"
)

func main() {
    profile := wasp.NewProfile()
    
    // Create HTTP client
    client := &http.Client{
        Timeout: 30 * time.Second,
    }
    
    generator := wasp.NewGenerator(&wasp.Config{
        T: 30 * time.Second,
        RPS: 50,
        LoadType: wasp.RPS,
        Schedule: wasp.Plain(50, 30*time.Second),
    })
    
    generator.AddRequestFn(func(ctx context.Context) error {
        req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/api/health", nil)
        if err != nil {
            return err
        }
        
        resp, err := client.Do(req)
        if err != nil {
            return err
        }
        defer resp.Body.Close()
        
        return nil
    })
    
    profile.Add(generator)
    
    _, err := profile.Run(true)
    if err != nil {
        log.Fatal(err)
    }
}
```

### Blockchain Load Test

```go
package main

import (
    "context"
    "log"
    "time"
    
    "github.com/ethereum/go-ethereum/ethclient"
    "github.com/smartcontractkit/chainlink-testing-framework/wasp"
)

func main() {
    profile := wasp.NewProfile()
    
    // Connect to blockchain
    client, err := ethclient.Dial("http://localhost:8545")
    if err != nil {
        log.Fatal(err)
    }
    
    generator := wasp.NewGenerator(&wasp.Config{
        T: 60 * time.Second,
        RPS: 10, // Lower RPS for blockchain
        LoadType: wasp.RPS,
        Schedule: wasp.Plain(10, 60*time.Second),
    })
    
    generator.AddRequestFn(func(ctx context.Context) error {
        // Get latest block number
        blockNumber, err := client.BlockNumber(ctx)
        if err != nil {
            return err
        }
        
        // Your blockchain interaction logic here
        log.Printf("Block number: %d", blockNumber)
        return nil
    })
    
    profile.Add(generator)
    
    _, err = profile.Run(true)
    if err != nil {
        log.Fatal(err)
    }
}
```

## Configuration Options

### Load Types

#### RPS (Requests Per Second)
```go
config := &wasp.Config{
    T: 30 * time.Second,
    RPS: 100,
    LoadType: wasp.RPS,
    Schedule: wasp.Plain(100, 30*time.Second),
}
```

#### VU (Virtual Users)
```go
config := &wasp.Config{
    T: 30 * time.Second,
    VU: 50,
    LoadType: wasp.VU,
    Schedule: wasp.Plain(50, 30*time.Second),
}
```

### Load Schedules

#### Plain (Constant Load)
```go
Schedule: wasp.Plain(100, 30*time.Second) // 100 RPS for 30 seconds
```

#### Ramp Up
```go
Schedule: wasp.RampUp(0, 100, 10*time.Second) // Ramp from 0 to 100 RPS over 10 seconds
```

#### Step
```go
Schedule: wasp.Step(10, 50, 5*time.Second) // Step from 10 to 50 RPS every 5 seconds
```

#### Custom Schedule
```go
Schedule: wasp.Custom([]wasp.Step{
    {Duration: 10 * time.Second, RPS: 10},
    {Duration: 20 * time.Second, RPS: 50},
    {Duration: 10 * time.Second, RPS: 100},
})
```

## Observability Integration

### Grafana Dashboard

WASP provides built-in Grafana dashboard support:

```go
profile := wasp.NewProfile()

// Configure Grafana options
grafanaOpts := &wasp.GrafanaOpts{
    GrafanaURL: "http://localhost:3000",
    GrafanaToken: "your-token",
    AnnotateDashboardUIDs: []string{"wasp-dashboard"},
    CheckDashboardAlertsAfterRun: []string{"wasp-dashboard"},
}

// Add Grafana integration
profile.WithGrafana(grafanaOpts)

// Run with dashboard annotations
_, err := profile.Run(true)
```

### Loki Integration

Send logs to Loki for aggregation:

```go
// Configure Loki
lokiConfig := &wasp.LokiConfig{
    URL: "http://localhost:3100/loki/api/v1/push",
    TenantID: "test-tenant",
    BatchWait: 5 * time.Second,
    BatchSize: 500 * 1024,
    Timeout: 20 * time.Second,
}

// Add Loki to generator
generator.WithLoki(lokiConfig)
```

### Custom Metrics

Collect custom metrics during your load test:

```go
generator.AddRequestFn(func(ctx context.Context) error {
    start := time.Now()
    
    // Your request logic here
    
    duration := time.Since(start)
    
    // Record custom metric
    generator.RecordMetric("custom_duration", duration.Seconds())
    
    return nil
})
```

## Advanced Usage

### Multiple Generators

```go
profile := wasp.NewProfile()

// Add different types of load
profile.Add(wasp.NewGenerator(&wasp.Config{
    T: 30 * time.Second,
    RPS: 50,
    LoadType: wasp.RPS,
    Schedule: wasp.Plain(50, 30*time.Second),
}))

profile.Add(wasp.NewGenerator(&wasp.Config{
    T: 30 * time.Second,
    VU: 10,
    LoadType: wasp.VU,
    Schedule: wasp.Plain(10, 30*time.Second),
}))

_, err := profile.Run(true)
```

### Conditional Requests

```go
generator.AddRequestFn(func(ctx context.Context) error {
    // Add some randomness to your requests
    if rand.Float64() < 0.1 {
        // 10% of requests are writes
        return performWrite(ctx)
    } else {
        // 90% of requests are reads
        return performRead(ctx)
    }
})
```

### Request Context

```go
generator.AddRequestFn(func(ctx context.Context) error {
    // Check if context is cancelled
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
    }
    
    // Your request logic here
    return nil
})
```

## Kubernetes Deployment

### Basic Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: wasp-load-test
spec:
  replicas: 3
  selector:
    matchLabels:
      app: wasp-load-test
  template:
    metadata:
      labels:
        app: wasp-load-test
    spec:
      containers:
      - name: wasp
        image: your-wasp-image:latest
        env:
        - name: TARGET_URL
          value: "http://your-service:8080"
        - name: TEST_DURATION
          value: "300s"
        - name: RPS
          value: "100"
```

### Distributed Load Testing

```go
// Configure for distributed testing
config := &wasp.Config{
    T: 300 * time.Second,
    RPS: 100,
    LoadType: wasp.RPS,
    Schedule: wasp.Plain(100, 300*time.Second),
    // Add distributed testing options
    Distributed: true,
    WorkerID: os.Getenv("WORKER_ID"),
    CoordinatorURL: os.Getenv("COORDINATOR_URL"),
}
```

## Best Practices

### 1. **Start Small**
- Begin with low RPS and short duration
- Gradually increase load to find system limits
- Monitor system resources during tests

### 2. **Use Realistic Scenarios**
- Model real user behavior
- Include think time between requests
- Vary request types and parameters

### 3. **Monitor System Health**
- Use the integrated observability stack
- Monitor target system metrics
- Set up alerts for critical thresholds

### 4. **Test in Production-Like Environment**
- Use similar hardware and configuration
- Include all dependencies and services
- Test with realistic data volumes

### 5. **Document Your Tests**
- Document test scenarios and parameters
- Record baseline performance metrics
- Track performance changes over time

## Troubleshooting

### Common Issues

#### High Memory Usage
```go
// Reduce batch sizes
lokiConfig := &wasp.LokiConfig{
    BatchSize: 100 * 1024, // Smaller batches
    BatchWait: 1 * time.Second,
}
```

#### Network Timeouts
```go
// Increase timeouts
client := &http.Client{
    Timeout: 60 * time.Second,
    Transport: &http.Transport{
        IdleConnTimeout: 30 * time.Second,
    },
}
```

#### Loki Connection Issues
```bash
# Check Loki connectivity
curl -X GET "http://localhost:3100/ready"

# Check WASP logs
WASP_LOG_LEVEL=trace go test -v
```

### Debug Mode

Enable debug logging:

```bash
export WASP_LOG_LEVEL=debug
go test -v
```

## Examples

Check the [examples directory](https://github.com/smartcontractkit/chainlink-testing-framework/tree/main/wasp/examples) for more comprehensive examples:

- [Basic HTTP load testing](https://github.com/smartcontractkit/chainlink-testing-framework/tree/main/wasp/examples/http)
- [Blockchain transaction testing](https://github.com/smartcontractkit/chainlink-testing-framework/tree/main/wasp/examples/blockchain)
- [Chainlink oracle testing](https://github.com/smartcontractkit/chainlink-testing-framework/tree/main/wasp/examples/chainlink)
- [Kubernetes deployment](https://github.com/smartcontractkit/chainlink-testing-framework/tree/main/wasp/examples/k8s)

## API Reference

For detailed API documentation, see the [Go package documentation](https://pkg.go.dev/github.com/smartcontractkit/chainlink-testing-framework/wasp).

## Contributing

We welcome contributions! Please see the [Contributing Guide](../../Contributing/Overview) for details on:

- Code of Conduct
- Development Setup
- Testing Guidelines
- Pull Request Process

## License

This project is licensed under the MIT License - see the [LICENSE](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/LICENSE) file for details. 