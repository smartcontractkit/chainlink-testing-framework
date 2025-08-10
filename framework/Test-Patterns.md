# Test Patterns

This guide covers common and advanced test patterns for writing robust, maintainable, and efficient tests with the Chainlink Testing Framework (CTF).

## Subtests

Subtests help organize your test logic and make it easier to debug and maintain. Use `t.Run` to group related assertions and steps.

### Example
```go
func TestEndToEnd(t *testing.T) {
    in, err := framework.Load[Config](t)
    require.NoError(t, err)
    bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
    require.NoError(t, err)

    t.Run("blockchain is running", func(t *testing.T) {
        require.NotEmpty(t, bc.Nodes[0].ExternalHTTPUrl)
    })

    t.Run("chainlink node is running", func(t *testing.T) {
        cl, err := clnode.NewChainlinkNode(in.ChainlinkNode)
        require.NoError(t, err)
        require.NotEmpty(t, cl.ExternalURL)
    })
}
```

## Parallelism

CTF supports parallel test execution for faster feedback and better resource utilization. Use `t.Parallel()` to enable parallelism.

### Example
```go
func TestParallel(t *testing.T) {
    t.Parallel()
    in, err := framework.Load[Config](t)
    require.NoError(t, err)
    // ... rest of test
}
```

- Use parallelism for independent tests that do not share resources.
- Be careful with shared state (e.g., Docker networks, files).

## Cleanup

CTF automatically cleans up Docker containers and resources after tests. You can add custom cleanup logic using `t.Cleanup`.

### Example
```go
func TestWithCleanup(t *testing.T) {
    in, err := framework.Load[Config](t)
    require.NoError(t, err)
    bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
    require.NoError(t, err)

    t.Cleanup(func() {
        // Custom cleanup logic
        log.Println("Cleaning up custom resources")
    })
}
```

## Table-Driven Tests

Table-driven tests are useful for running the same logic with different inputs.

### Example
```go
testCases := []struct {
    name   string
    config string
}{
    {"anvil", "anvil.toml"},
    {"geth", "geth.toml"},
}

for _, tc := range testCases {
    t.Run(tc.name, func(t *testing.T) {
        t.Parallel()
        in, err := framework.LoadConfig(tc.config)
        require.NoError(t, err)
        // ... test logic
    })
}
```

## Advanced Patterns

### 1. **Retry Logic**
Use retry logic for flaky external dependencies.
```go
require.Eventually(t, func() bool {
    // Check condition
    return isReady()
}, 30*time.Second, 2*time.Second)
```

### 2. **Timeouts**
Set timeouts for long-running tests.
```bash
go test -v -timeout 10m
```

### 3. **Dynamic Configuration**
Use environment variables or flags to inject dynamic values.
```bash
CTF_CONFIGS=smoke.toml CTF_LOG_LEVEL=debug go test -v -run TestSmoke
```

### 4. **Observability in Tests**
Add custom logs and metrics for better debugging.
```go
log := logging.NewLogger()
log.Info().Msg("Starting test")
```

### 5. **Combining Load and Chaos**
Combine WASP load tests and Havoc chaos experiments for resilience testing.
```go
go func() {
    // Start chaos
    havoc.NewChaos(...).Create(context.Background())
}()
profile.Run(true) // WASP load test
```

## Best Practices

- **Organize tests** with subtests and table-driven patterns
- **Use parallelism** for independent tests
- **Clean up** resources with `t.Cleanup`
- **Add observability** (logs, metrics, traces) to all tests
- **Document** test cases and expected outcomes
- **Use caching** for faster development, but disable for critical tests
- **Monitor** resource usage and test performance

## Further Reading
- [Component System](Components)
- [Configuration Guide](Configuration)
- [Observability Guide](Observability)
- [Caching Guide](Caching)
- [WASP Load Testing](../Libraries/WASP)
- [Havoc Chaos Testing](../Libraries/Havoc) 