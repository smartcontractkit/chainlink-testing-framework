# Framework Overview

The Chainlink Testing Framework (CTF) is designed to reduce the complexity of end-to-end testing, making complex system-level tests appear straightforward. It enables tests to run in any environment and serves as a single source of truth for system behavior as defined by requirements.

## Core Philosophy

### 1. **Straightforward and Sequential Test Composition**
Tests are readable and give you precise control over key aspects in a strict step-by-step order. No hidden magic or complex abstractions.

### 2. **Modular Configuration**
No arcane knowledge of framework settings is required. The configuration is simply a reflection of the components being used in the test. Components declare their own configuration—**what you see is what you get**.

### 3. **Component Isolation**
Components are decoupled via input/output structs, without exposing internal details. This allows for easy testing and replacement of individual components.

### 4. **Replaceability and Extensibility**
Since components are decoupled via outputs, any deployment component can be swapped with a real service without altering the test code.

### 5. **Quick Local Environments**
A common setup can be launched in just **15 seconds** 🚀 (when using caching).

### 6. **Integrated Observability Stack**
Get all the information you need to develop end-to-end tests: metrics, logs, traces, and profiles.

## Architecture

### Component-Based Design

CTF uses a component-based architecture where each component represents a service or system that can be deployed and tested:

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Blockchain    │    │  Chainlink Node │    │   Database      │
│   Component     │    │   Component     │    │   Component     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
                    ┌─────────────────┐
                    │   Test Logic    │
                    │   (Your Code)   │
                    └─────────────────┘
```

### Configuration-Driven Testing

Tests are driven by TOML configuration files that define what components to use and how to configure them:

```toml
[blockchain_a]
  type = "anvil"

[chainlink_node]
  image = "public.ecr.aws/chainlink/chainlink"
  tag = "2.7.0"
  pull_image = true

[postgres]
  image = "postgres:15"
  tag = "15"
  pull_image = true
```

### Input/Output Pattern

Each component follows a consistent input/output pattern:

```go
type Input struct {
    // Configuration fields
    Image     string `toml:"image" validate:"required"`
    Tag       string `toml:"tag" validate:"required"`
    PullImage bool   `toml:"pull_image"`
    
    // Output is embedded for caching
    Out *Output `toml:"out"`
}

type Output struct {
    UseCache bool   `toml:"use_cache"`
    URL      string `toml:"url"`
    // Other output fields
}
```

## Key Features

### 1. **Docker Integration**
- Automatic container management with [testcontainers-go](https://golang.testcontainers.org/)
- Built-in networking and service discovery
- Automatic cleanup and resource management

### 2. **Multi-Blockchain Support**
- **Ethereum**: Anvil, Geth, Besu
- **Solana**: Local validator
- **TON**: Local blockchain
- **Aptos**: Local blockchain
- **Sui**: Local blockchain
- **Tron**: Local blockchain

### 3. **Chainlink Integration**
- Chainlink node deployment and configuration
- Job management and monitoring
- Oracle contract integration
- Multi-node setups (DONs)

### 4. **Observability Stack**
- **Grafana**: Dashboards and visualization
- **Loki**: Log aggregation
- **Prometheus**: Metrics collection
- **Jaeger**: Distributed tracing
- **Pyroscope**: Performance profiling

### 5. **Caching System**
- Component state caching for faster test development
- Configurable cache invalidation
- Cross-test cache sharing

### 6. **CLI Tools**
- Environment management (`ctf obs up`, `ctf bs up`)
- Container cleanup (`ctf d rm`)
- Configuration validation
- Test execution helpers

## Component Types

### 1. **Blockchain Components**
- Local blockchain networks for testing
- Support for multiple blockchain types
- Automatic network configuration

### 2. **Chainlink Components**
- Chainlink node deployment
- Job management
- Oracle contract integration

### 3. **Infrastructure Components**
- Databases (PostgreSQL, etc.)
- Message queues
- External services
- Monitoring stacks

### 4. **Custom Components**
- User-defined components
- Integration with external services
- Specialized testing requirements

## Test Lifecycle

### 1. **Configuration Loading**
```go
in, err := framework.Load[Config](t)
require.NoError(t, err)
```

### 2. **Component Creation**
```go
bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
require.NoError(t, err)
```

### 3. **Test Execution**
```go
t.Run("test functionality", func(t *testing.T) {
    // Your test logic here
})
```

### 4. **Automatic Cleanup**
CTF automatically cleans up resources when tests complete.

## Environment Support

### Local Development
- Docker-based local environments
- Fast startup with caching
- Integrated observability

### CI/CD Integration
- GitHub Actions support
- Parallel test execution
- Resource optimization

### Kubernetes
- K8s test runner for distributed testing
- Chaos testing with Chaos Mesh
- Scalable test environments

## Performance Characteristics

### Startup Time
- **First run**: ~30-60 seconds (depending on Docker image pulls)
- **Cached run**: ~15 seconds
- **Component reuse**: ~5 seconds

### Resource Usage
- **Memory**: 2-8GB depending on components
- **CPU**: 2-4 cores recommended
- **Storage**: 10-20GB for Docker images and data

### Scalability
- **Parallel tests**: Full support with resource isolation
- **Large deployments**: Support for complex multi-component setups
- **Distributed testing**: Kubernetes integration for scale-out

## Best Practices

### 1. **Configuration Management**
- Use descriptive configuration names
- Validate required fields
- Document configuration options

### 2. **Test Organization**
- Use subtests for logical grouping
- Keep tests focused and readable
- Use meaningful test names

### 3. **Resource Management**
- Leverage caching for faster development
- Clean up resources appropriately
- Monitor resource usage

### 4. **Observability**
- Use the integrated observability stack
- Add custom metrics and logs
- Monitor test performance

## Comparison with Other Frameworks

| Feature | CTF | Traditional E2E | Unit Testing |
|---------|-----|-----------------|--------------|
| **Setup Complexity** | Low | High | Very Low |
| **Environment Management** | Automatic | Manual | None |
| **Component Isolation** | Built-in | Manual | Built-in |
| **Observability** | Integrated | External | Limited |
| **Multi-Service Testing** | Native | Complex | Not Applicable |
| **Blockchain Integration** | First-class | External | Not Applicable |

## Getting Started

1. **[Installation](../Getting-Started/Installation)** - Set up your development environment
2. **[First Test](../Getting-Started/First-Test)** - Write your first CTF test
3. **[Components](Components)** - Learn about available components
4. **[Configuration](Configuration)** - Understand configuration options
5. **[Observability](Observability)** - Set up monitoring and logging

## Next Steps

- Explore the **[Component System](Components)** to understand available components
- Learn about **[Configuration Management](Configuration)** for complex setups
- Set up **[Observability](Observability)** for better debugging
- Check out **[Advanced Patterns](Test-Patterns)** for complex scenarios 