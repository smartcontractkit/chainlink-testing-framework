# Advanced Examples

This page provides advanced and real-world test scenarios using the Chainlink Testing Framework (CTF). These examples demonstrate how to leverage CTF for complex, production-like workflows.

---

## 1. Multi-Chain Integration Test

Test a system that interacts with multiple blockchains (e.g., Ethereum and Solana) and verifies cross-chain data flow.

```toml
[ethereum]
  type = "geth"
  image = "ethereum/client-go"
  tag = "v1.12.0"
  pull_image = true

[solana]
  type = "solana"
  image = "solanalabs/solana"
  tag = "v1.16.0"
  pull_image = true

[chainlink_node]
  image = "public.ecr.aws/chainlink/chainlink"
  tag = "2.7.0"
  pull_image = true
```

```go
func TestMultiChainIntegration(t *testing.T) {
    type Config struct {
        Ethereum      *blockchain.Input `toml:"ethereum" validate:"required"`
        Solana        *blockchain.Input `toml:"solana" validate:"required"`
        ChainlinkNode *clnode.Input     `toml:"chainlink_node" validate:"required"`
    }
    in, err := framework.Load[Config](t)
    require.NoError(t, err)
    eth, err := blockchain.NewBlockchainNetwork(in.Ethereum)
    require.NoError(t, err)
    sol, err := blockchain.NewBlockchainNetwork(in.Solana)
    require.NoError(t, err)
    cl, err := clnode.NewChainlinkNode(in.ChainlinkNode)
    require.NoError(t, err)
    // ... deploy contracts, set up jobs, verify cross-chain data
}
```

---

## 2. Upgrade and Migration Test

Test that a Chainlink node or contract can be upgraded without breaking existing functionality.

```toml
[chainlink_node]
  image = "public.ecr.aws/chainlink/chainlink"
  tag = "2.6.0"
  pull_image = true
```

```go
func TestUpgrade(t *testing.T) {
    // Deploy with old version
    in, err := framework.Load[Config](t)
    require.NoError(t, err)
    cl, err := clnode.NewChainlinkNode(in.ChainlinkNode)
    require.NoError(t, err)
    // ... run smoke test

    // Upgrade node
    in.ChainlinkNode.Tag = "2.7.0"
    clUpgraded, err := clnode.NewChainlinkNode(in.ChainlinkNode)
    require.NoError(t, err)
    // ... verify data, jobs, and state are preserved
}
```

---

## 3. Chaos Engineering Test

Inject failures and network issues to test system resilience using Havoc.

```go
func TestChaosResilience(t *testing.T) {
    in, err := framework.Load[Config](t)
    require.NoError(t, err)
    cl, err := clnode.NewChainlinkNode(in.ChainlinkNode)
    require.NoError(t, err)
    // Start chaos experiment
    client, err := havoc.NewClient()
    require.NoError(t, err)
    chaos, err := havoc.NewChaos(client, createNetworkChaos())
    require.NoError(t, err)
    err = chaos.Create(context.Background())
    require.NoError(t, err)
    // ... run test logic while chaos is active
    chaos.Delete(context.Background())
}
```

---

## 4. Performance and Load Test

Use WASP to generate load and measure system performance under stress.

```go
func TestPerformance(t *testing.T) {
    profile := wasp.NewProfile()
    generator := wasp.NewGenerator(&wasp.Config{
        T: 60 * time.Second,
        RPS: 100,
        LoadType: wasp.RPS,
        Schedule: wasp.Plain(100, 60*time.Second),
    })
    generator.AddRequestFn(func(ctx context.Context) error {
        // Simulate contract call or API request
        return nil
    })
    profile.Add(generator)
    _, err := profile.Run(true)
    require.NoError(t, err)
}
```

---

## 5. End-to-End Oracle Test

Test the full workflow from contract deployment to job fulfillment and data reporting.

```go
func TestOracleE2E(t *testing.T) {
    in, err := framework.Load[Config](t)
    require.NoError(t, err)
    bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
    require.NoError(t, err)
    cl, err := clnode.NewChainlinkNode(in.ChainlinkNode)
    require.NoError(t, err)
    // Deploy oracle contract
    // Register job on Chainlink node
    // Send request and verify fulfillment
}
```

---

## 6. Real-World: Staging Environment Reuse

Reuse cached components and substitute staging URLs for persistent environment testing.

```toml
[blockchain_a]
  type = "anvil"
  use_cache = true
  external_url = "https://staging-eth.example.com"

[chainlink_node]
  use_cache = true
  external_url = "https://staging-cl.example.com"
```

```go
func TestStagingReuse(t *testing.T) {
    in, err := framework.Load[Config](t)
    require.NoError(t, err)
    // Use staging URLs for integration tests
}
```

---

## More Examples
- [WASP Load Testing](../Libraries/WASP)
- [Havoc Chaos Testing](../Libraries/Havoc)
- [Seth Ethereum Client](../Libraries/Seth)
- [Framework Examples Directory](https://github.com/smartcontractkit/chainlink-testing-framework/tree/main/framework/examples/myproject) 