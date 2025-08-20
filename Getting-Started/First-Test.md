# Writing Your First Test

This guide will walk you through creating your first test with the Chainlink Testing Framework (CTF). We'll start with a simple blockchain test and gradually build up to more complex scenarios.

## Prerequisites

Before starting, make sure you have:
- [Installed CTF](Installation)
- Docker running
- A basic understanding of Go testing

## Basic Test Structure

CTF tests follow a simple pattern:
1. **Configuration** - Define what components you need
2. **Loading** - Load configuration into your test
3. **Component Creation** - Create and start components
4. **Testing** - Write your test logic
5. **Cleanup** - Clean up resources (automatic with CTF)

## Step 1: Create Configuration

Create a `smoke.toml` file in your project root:

```toml
[blockchain_a]
  type = "anvil"
```

This configuration tells CTF to create an Anvil blockchain instance. Anvil is a fast Ethereum development node that's perfect for testing.

## Step 2: Write Your Test

Create `smoke_test.go`:

```go
package mymodule_test

import (
    "github.com/smartcontractkit/chainlink-testing-framework/framework"
    "github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
    "github.com/stretchr/testify/require"
    "testing"
)

// Config defines the structure of your configuration file
type Config struct {
    BlockchainA *blockchain.Input `toml:"blockchain_a" validate:"required"`
}

func TestSmoke(t *testing.T) {
    // Load configuration from TOML file
    in, err := framework.Load[Config](t)
    require.NoError(t, err)

    // Create blockchain network
    bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
    require.NoError(t, err)

    // Your test logic goes here
    t.Run("blockchain is running", func(t *testing.T) {
        require.NotEmpty(t, bc.Nodes[0].ExternalHTTPUrl)
        require.NotEmpty(t, bc.Nodes[0].InternalHTTPUrl)
    })

    t.Run("blockchain is accessible", func(t *testing.T) {
        // Test that we can connect to the blockchain
        client := bc.GetDefaultClient()
        require.NotNil(t, client)
        
        // Get the latest block number
        blockNumber, err := client.BlockNumber(context.Background())
        require.NoError(t, err)
        require.GreaterOrEqual(t, blockNumber, uint64(0))
    })
}
```

## Step 3: Run Your Test

```bash
# Run the test
CTF_CONFIGS=smoke.toml go test -v -run TestSmoke

# Clean up containers when done
ctf d rm
```

## Understanding the Test

Let's break down what's happening:

### Configuration Loading
```go
in, err := framework.Load[Config](t)
```
This loads your `smoke.toml` file and validates it against the `Config` struct. The `validate:"required"` tag ensures the blockchain configuration is present.

### Component Creation
```go
bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
```
This creates a blockchain network based on your configuration. CTF will:
- Start a Docker container with Anvil
- Configure networking
- Wait for the service to be ready
- Return a network object with connection details

### Test Logic
```go
t.Run("blockchain is running", func(t *testing.T) {
    require.NotEmpty(t, bc.Nodes[0].ExternalHTTPUrl)
})
```
This tests that the blockchain is actually running and accessible. The `ExternalHTTPUrl` is the URL you can use to connect to the blockchain from your host machine.

## Adding More Components

Let's expand the test to include a Chainlink node:

### Updated Configuration
```toml
[blockchain_a]
  type = "anvil"

[chainlink_node]
  image = "public.ecr.aws/chainlink/chainlink"
  tag = "2.7.0"
  pull_image = true
  node_config = """
  [Log]
  Level = "info"
  
  [WebServer]
  TLSListenPort = 0
  Port = 6688
  
  [Database]
  URL = "postgresql://postgres:password@postgres:5432/chainlink?sslmode=disable"
  """
```

### Updated Test
```go
package mymodule_test

import (
    "context"
    "github.com/smartcontractkit/chainlink-testing-framework/framework"
    "github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
    "github.com/smartcontractkit/chainlink-testing-framework/framework/components/clnode"
    "github.com/stretchr/testify/require"
    "testing"
)

type Config struct {
    BlockchainA *blockchain.Input `toml:"blockchain_a" validate:"required"`
    ChainlinkNode *clnode.Input   `toml:"chainlink_node" validate:"required"`
}

func TestSmokeWithChainlink(t *testing.T) {
    in, err := framework.Load[Config](t)
    require.NoError(t, err)

    // Create blockchain network
    bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
    require.NoError(t, err)

    // Create Chainlink node
    cl, err := clnode.NewChainlinkNode(in.ChainlinkNode)
    require.NoError(t, err)

    t.Run("blockchain is running", func(t *testing.T) {
        require.NotEmpty(t, bc.Nodes[0].ExternalHTTPUrl)
    })

    t.Run("chainlink node is running", func(t *testing.T) {
        require.NotEmpty(t, cl.ExternalURL)
        
        // Test that the Chainlink API is accessible
        resp, err := http.Get(cl.ExternalURL + "/health")
        require.NoError(t, err)
        require.Equal(t, http.StatusOK, resp.StatusCode)
    })

    t.Run("chainlink can connect to blockchain", func(t *testing.T) {
        // This would typically involve setting up a job
        // and verifying it can interact with the blockchain
        require.NotEmpty(t, bc.Nodes[0].InternalHTTPUrl)
    })
}
```

## Test Patterns

### Using Subtests
CTF encourages the use of subtests (`t.Run`) to organize your test logic:

```go
func TestComprehensive(t *testing.T) {
    // Setup
    in, err := framework.Load[Config](t)
    require.NoError(t, err)

    bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
    require.NoError(t, err)

    // Test different aspects
    t.Run("connectivity", func(t *testing.T) {
        // Test network connectivity
    })

    t.Run("functionality", func(t *testing.T) {
        // Test core functionality
    })

    t.Run("performance", func(t *testing.T) {
        // Test performance characteristics
    })
}
```

### Parallel Testing
You can run tests in parallel for faster execution:

```go
func TestParallel(t *testing.T) {
    t.Parallel() // Enable parallel execution
    
    in, err := framework.Load[Config](t)
    require.NoError(t, err)
    // ... rest of test
}
```

### Test Cleanup
CTF automatically handles cleanup of Docker containers, but you can add custom cleanup:

```go
func TestWithCleanup(t *testing.T) {
    in, err := framework.Load[Config](t)
    require.NoError(t, err)

    bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
    require.NoError(t, err)

    // Add custom cleanup
    t.Cleanup(func() {
        // Custom cleanup logic
        log.Println("Cleaning up custom resources")
    })

    // Your test logic
}
```

## Running Tests

### Basic Test Execution
```bash
# Run all tests
go test -v

# Run specific test
go test -v -run TestSmoke

# Run tests with specific config
CTF_CONFIGS=smoke.toml go test -v -run TestSmoke

# Run tests multiple times
go test -v -run TestSmoke -count 5
```

### Test Flags
```bash
# Run tests with timeout
go test -v -timeout 10m -run TestSmoke

# Run tests with race detection
go test -v -race -run TestSmoke

# Run tests with coverage
go test -v -cover -run TestSmoke
```

### Environment Variables
```bash
# Set log level
CTF_LOG_LEVEL=debug go test -v -run TestSmoke

# Disable caching
CTF_DISABLE_CACHE=true go test -v -run TestSmoke

# Use specific Docker network
CTF_NETWORK_NAME=my-network go test -v -run TestSmoke
```

## Next Steps

Now that you've written your first test, you can:

1. **Add more components** - Try adding databases, external services, or multiple blockchain networks
2. **Learn about configuration** - Explore the [Configuration Guide](Configuration) for more options
3. **Set up observability** - Add monitoring and logging with the [Environment Setup Guide](Environment-Setup)
4. **Explore advanced patterns** - Check out the [Framework Test Patterns](../Framework/Test-Patterns)
5. **Look at examples** - Browse the [Examples section](../Examples/Basic-Examples) for more complex scenarios

## Common Issues

### Test Fails to Start
- Check that Docker is running
- Verify your configuration file syntax
- Ensure all required fields are present

### Components Not Ready
- CTF automatically waits for components to be ready
- Check logs for component startup issues
- Verify network connectivity

### Port Conflicts
- CTF automatically assigns ports to avoid conflicts
- If you get port binding errors, check what's using the ports
- Use `ctf d rm` to clean up existing containers

### Memory Issues
- Some components (like blockchain nodes) can be memory-intensive
- Consider using OrbStack instead of Docker Desktop
- Increase Docker memory limits if needed 