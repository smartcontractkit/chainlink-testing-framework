# Seth - Ethereum Client Library

[![Go Report Card](https://goreportcard.com/badge/github.com/smartcontractkit/chainlink-testing-framework/seth)](https://goreportcard.com/report/github.com/smartcontractkit/chainlink-testing-framework/seth)
[![Decoding tests](https://github.com/smartcontractkit/chainlink-testing-framework/actions/workflows/seth-test-decode.yml/badge.svg)](https://github.com/smartcontractkit/seth/actions/workflows/test_decode.yml)
[![Tracing tests](https://github.com/smartcontractkit/chainlink-testing-framework/actions/workflows/seth-test-trace.yml/badge.svg)](https://github.com/smartcontractkit/seth/actions/workflows/test_trace.yml)
[![Gas bumping tests](https://github.com/smartcontractkit/chainlink-testing-framework/actions/workflows/seth-test-bumping.yml/badge.svg)](https://github.com/smartcontractkit/seth/actions/workflows/test_cli.yml)
[![API tests](https://github.com/smartcontractkit/chainlink-testing-framework/actions/workflows/seth-test-api.yml/badge.svg)](https://github.com/smartcontractkit/seth/actions/workflows/test_api.yml)
[![CLI tests](https://github.com/smartcontractkit/chainlink-testing-framework/actions/workflows/seth-test-cli.yml/badge.svg)](https://github.com/smartcontractkit/seth/actions/workflows/test_cli.yml)
[![Integration tests (testnets)](https://github.com/smartcontractkit/chainlink-testing-framework/actions/workflows/seth-test-decode-testnet.yml/badge.svg)](https://github.com/smartcontractkit/seth/actions/workflows/test_decode_testnet.yml)

## Overview

Seth is a reliable and debug-friendly Ethereum client library that provides a thin, debuggable wrapper on top of `go-ethereum`. It's designed to make Ethereum development and testing easier by automatically decoding transaction inputs/outputs/logs and providing advanced debugging capabilities.

## Goals

- **Thin wrapper** on top of `go-ethereum` with minimal overhead
- **Debuggable** - comprehensive transaction decoding and tracing
- **Battle tested** - extensive test coverage including testnet integration
- **Automatic decoding** - decode all transaction inputs/outputs/logs for all ABIs
- **Simple synchronous API** - easy to use and understand
- **Resilient** - execute transactions even during gas spikes or RPC outages
- **Well tested** - comprehensive e2e test suite for testnet integration

## Key Features

### ✅ Implemented Features

- [x] **Decode named inputs** - Automatically decode function parameters
- [x] **Decode named outputs** - Decode function return values
- [x] **Decode anonymous outputs** - Handle unnamed return values
- [x] **Decode logs** - Parse event logs with full context
- [x] **Decode indexed logs** - Handle indexed event parameters
- [x] **Decode old string reverts** - Parse legacy revert messages
- [x] **Decode new typed reverts** - Handle modern revert types
- [x] **EIP-1559 support** - Full support for London hard fork features
- [x] **Multi-keys client support** - Use multiple private keys
- [x] **CLI to manipulate test keys** - Command-line key management
- [x] **Simple manual gas price estimation** - Built-in gas optimization
- [x] **Decode collided event hashes** - Handle event signature collisions
- [x] **Tracing support (4byte)** - Function signature tracing
- [x] **Tracing support (callTracer)** - Call trace analysis
- [x] **Tracing decoding** - Decode trace results
- [x] **Tracing tests** - Comprehensive tracing test coverage
- [x] **Saving deployed contracts mapping** - Track contract addresses
- [x] **Reading deployed contracts mappings** - Load contract mappings
- [x] **Automatic gas estimator** - Intelligent gas price estimation
- [x] **Block stats CLI** - Block analysis tools
- [x] **Pending nonce checking** - Prevent nonce conflicts
- [x] **DOT graph output** - Visualize transaction traces
- [x] **Gas bumping for slow transactions** - Automatic retry with gas increase

### 🚧 Planned Features

- [ ] **Fail over client logic** - Automatic RPC failover
- [ ] **Tracing support (prestate)** - Pre-state tracing

## Installation

### Using Go Modules
```bash
go get github.com/smartcontractkit/chainlink-testing-framework/seth
```

### Using Nix (Recommended for Development)
```bash
# Install nix
curl --proto '=https' --tlsv1.2 -sSf -L https://install.determinate.systems/nix | sh -s -- install

# Enter development shell
nix develop
```

## Quick Start

### Basic Usage

```go
package main

import (
    "context"
    "log"
    
    "github.com/smartcontractkit/chainlink-testing-framework/seth"
)

func main() {
    // Create a new Seth client
    client, err := seth.NewClient("http://localhost:8545")
    if err != nil {
        log.Fatal(err)
    }
    
    // Get the latest block number
    blockNumber, err := client.BlockNumber(context.Background())
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("Latest block: %d", blockNumber)
}
```

### With Configuration

```go
package main

import (
    "context"
    "log"
    
    "github.com/smartcontractkit/chainlink-testing-framework/seth"
)

func main() {
    // Create configuration
    cfg := seth.DefaultConfig("http://localhost:8545", []string{
        "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
    })
    
    // Create client with configuration
    client, err := seth.NewClientWithConfig(cfg)
    if err != nil {
        log.Fatal(err)
    }
    
    // Use the client
    balance, err := client.BalanceAt(context.Background(), client.Addresses[0], nil)
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("Balance: %s", balance.String())
}
```

## Configuration

### Simplified Configuration

For basic use cases, you can use the simplified configuration:

```go
cfg := seth.DefaultConfig("ws://localhost:8546", []string{
    "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
})
client, err := seth.NewClientWithConfig(cfg)
```

This uses reasonable defaults for:
- Gas price estimation
- Transaction timeout
- RPC timeout
- Gas limit estimation

### TOML Configuration

For more complex setups, use TOML configuration:

```toml
# seth.toml
[networks.anvil]
urls = ["http://localhost:8545"]
private_keys = ["ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"]
gas_price_estimation = true
gas_limit_estimation = true
gas_bump_percent = 10
gas_bump_wei = 1000000000
max_gas_price = 100000000000
min_gas_price = 1000000000
gas_bump_threshold = 3
rpc_timeout = "30s"
tx_timeout = "5m"
```

### Environment Variables

Seth supports various environment variables for configuration:

```bash
# Network configuration
SETH_NETWORK=anvil
SETH_URL=http://localhost:8545
SETH_PRIVATE_KEYS=ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80

# Gas configuration
SETH_GAS_PRICE_ESTIMATION=true
SETH_GAS_BUMP_PERCENT=10
SETH_MAX_GAS_PRICE=100000000000

# Timeout configuration
SETH_RPC_TIMEOUT=30s
SETH_TX_TIMEOUT=5m
```

## Transaction Decoding

### Automatic Decoding

Seth automatically decodes all transaction data when you have the corresponding ABI:

```go
// Deploy a contract
contractAddress, tx, _, err := counter.DeployCounter(client.NewTXOpts(), client)
if err != nil {
    log.Fatal(err)
}

// Decode the deployment transaction
decodedTx, err := client.Decode(tx, nil)
if err != nil {
    log.Fatal(err)
}

// Print decoded information
fmt.Printf("Contract deployed at: %s\n", contractAddress.Hex())
fmt.Printf("Gas used: %d\n", decodedTx.GasUsed)
fmt.Printf("Gas price: %s\n", decodedTx.GasPrice.String())
```

### Function Call Decoding

```go
// Call a function
tx, err := counter.Increment(client.NewTXOpts())
if err != nil {
    log.Fatal(err)
}

// Decode the transaction
decodedTx, err := client.Decode(tx, nil)
if err != nil {
    log.Fatal(err)
}

// Print function call details
fmt.Printf("Function: %s\n", decodedTx.FunctionName)
fmt.Printf("Inputs: %+v\n", decodedTx.Inputs)
```

### Event Log Decoding

```go
// Get logs for a specific event
logs, err := client.FilterLogs(context.Background(), query)
if err != nil {
    log.Fatal(err)
}

// Decode each log
for _, log := range logs {
    decodedLog, err := client.DecodeLog(log)
    if err != nil {
        continue
    }
    
    fmt.Printf("Event: %s\n", decodedLog.EventName)
    fmt.Printf("Data: %+v\n", decodedLog.Data)
}
```

## Transaction Tracing

### Call Tracing

```go
// Trace a transaction
trace, err := client.TraceTransaction(context.Background(), txHash)
if err != nil {
    log.Fatal(err)
}

// Print trace information
fmt.Printf("Trace type: %s\n", trace.Type)
fmt.Printf("Gas used: %d\n", trace.Result.GasUsed)
fmt.Printf("Calls: %d\n", len(trace.Result.Calls))
```

### Function Signature Tracing

```go
// Get function signature from 4byte database
signature, err := client.GetFunctionSignature("0xa9059cbb")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Function signature: %s\n", signature)
```

### DOT Graph Generation

```go
// Generate DOT graph for transaction trace
dotGraph, err := client.GenerateDOTGraph(txHash)
if err != nil {
    log.Fatal(err)
}

// Save to file
err = os.WriteFile("trace.dot", []byte(dotGraph), 0644)
if err != nil {
    log.Fatal(err)
}
```

## Gas Management

### Automatic Gas Estimation

```go
// Enable automatic gas estimation
cfg := seth.DefaultConfig("http://localhost:8545", []string{privateKey})
cfg.GasPriceEstimation = true
cfg.GasLimitEstimation = true

client, err := seth.NewClientWithConfig(cfg)
if err != nil {
    log.Fatal(err)
}

// Transactions will automatically use estimated gas
tx, err := counter.Increment(client.NewTXOpts())
if err != nil {
    log.Fatal(err)
}
```

### Manual Gas Price Estimation

```go
// Estimate gas price manually
gasPrice, err := client.EstimateGasPrice(context.Background())
if err != nil {
    log.Fatal(err)
}

// Use estimated gas price
opts := client.NewTXOpts()
opts.GasPrice = gasPrice

tx, err := counter.Increment(opts)
if err != nil {
    log.Fatal(err)
}
```

### Gas Bumping

```go
// Configure gas bumping
cfg := seth.DefaultConfig("http://localhost:8545", []string{privateKey})
cfg.GasBumpPercent = 10
cfg.GasBumpThreshold = 3

client, err := seth.NewClientWithConfig(cfg)
if err != nil {
    log.Fatal(err)
}

// If transaction is slow, it will be automatically bumped
tx, err := counter.Increment(client.NewTXOpts())
if err != nil {
    log.Fatal(err)
}
```

## Multi-Key Support

### Using Multiple Keys

```go
// Create client with multiple keys
cfg := seth.DefaultConfig("http://localhost:8545", []string{
    "key1...",
    "key2...",
    "key3...",
})

client, err := seth.NewClientWithConfig(cfg)
if err != nil {
    log.Fatal(err)
}

// Use specific key for transaction
opts := client.NewTXKeyOpts(1) // Use second key
tx, err := counter.Increment(opts)
if err != nil {
    log.Fatal(err)
}
```

### Key Management CLI

```bash
# List all keys
seth keys list

# Add a new key
seth keys add --private-key=0x...

# Remove a key
seth keys remove --index=0

# Export a key
seth keys export --index=0
```

## Contract Management

### Contract Mapping

```go
// Save contract mapping
err := client.SaveContractMapping("Counter", contractAddress, abi)
if err != nil {
    log.Fatal(err)
}

// Load contract mapping
abi, err := client.LoadContractMapping("Counter", contractAddress)
if err != nil {
    log.Fatal(err)
}
```

### ABI Finder

```go
// Find ABI by contract name
abi, err := client.FindABI("Counter")
if err != nil {
    log.Fatal(err)
}

// Find ABI by address
abi, err := client.FindABIByAddress(contractAddress)
if err != nil {
    log.Fatal(err)
}
```

## Block Analysis

### Block Stats

```go
// Get block statistics
stats, err := client.GetBlockStats(context.Background(), blockNumber)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Block %d:\n", stats.BlockNumber)
fmt.Printf("  Transactions: %d\n", stats.TransactionCount)
fmt.Printf("  Gas used: %d\n", stats.GasUsed)
fmt.Printf("  Gas limit: %d\n", stats.GasLimit)
fmt.Printf("  Base fee: %s\n", stats.BaseFee.String())
```

### Block Stats CLI

```bash
# Get stats for latest block
seth block stats

# Get stats for specific block
seth block stats --block=12345

# Get stats for range of blocks
seth block stats --from=1000 --to=1010
```

## Advanced Features

### Read-Only Mode

```go
// Create read-only client
client, err := seth.NewReadOnlyClient("http://localhost:8545")
if err != nil {
    log.Fatal(err)
}

// Only read operations are allowed
balance, err := client.BalanceAt(context.Background(), address, nil)
if err != nil {
    log.Fatal(err)
}
```

### RPC Traffic Logging

```go
// Enable RPC traffic logging
cfg := seth.DefaultConfig("http://localhost:8545", []string{privateKey})
cfg.LogRPC = true

client, err := seth.NewClientWithConfig(cfg)
if err != nil {
    log.Fatal(err)
}

// All RPC calls will be logged
```

### Bulk Transaction Tracing

```go
// Trace multiple transactions
txHashes := []common.Hash{hash1, hash2, hash3}
traces, err := client.TraceTransactions(context.Background(), txHashes)
if err != nil {
    log.Fatal(err)
}

for i, trace := range traces {
    fmt.Printf("Transaction %d: %s\n", i, trace.Type)
}
```

## Testing

### Running Tests

```bash
# Run all tests
make test

# Run specific test suite
make test_decode
make test_trace
make test_api
make test_cli

# Run tests on specific network
make network=Anvil test
make network=Geth test
```

### Testnet Integration

```bash
# Run tests on testnet
make network=Sepolia test
make network=Goerli test
```

## Examples

Check the [examples directory](https://github.com/smartcontractkit/chainlink-testing-framework/tree/main/seth/examples) for comprehensive examples:

- [Basic client usage](https://github.com/smartcontractkit/chainlink-testing-framework/tree/main/seth/examples/basic)
- [Contract deployment and interaction](https://github.com/smartcontractkit/chainlink-testing-framework/tree/main/seth/examples/contracts)
- [Transaction tracing](https://github.com/smartcontractkit/chainlink-testing-framework/tree/main/seth/examples/tracing)
- [Gas optimization](https://github.com/smartcontractkit/chainlink-testing-framework/tree/main/seth/examples/gas)

## CLI Reference

### Basic Commands

```bash
# Get help
seth --help

# Get version
seth version

# Get network info
seth network info
```

### Key Management

```bash
# List keys
seth keys list

# Add key
seth keys add --private-key=0x...

# Remove key
seth keys remove --index=0

# Export key
seth keys export --index=0
```

### Block Analysis

```bash
# Get block stats
seth block stats

# Get block info
seth block info --block=12345

# Get transaction info
seth tx info --hash=0x...
```

### Contract Management

```bash
# Save contract mapping
seth contract save --name=Counter --address=0x... --abi=path/to/abi.json

# Load contract mapping
seth contract load --name=Counter --address=0x...

# List contracts
seth contract list
```

## Best Practices

### 1. **Use Configuration Files**
- Store network configurations in TOML files
- Use environment variables for sensitive data
- Version control your configurations

### 2. **Handle Errors Gracefully**
```go
decodedTx, err := client.Decode(tx, nil)
if err != nil {
    // Log error but don't fail the test
    log.Printf("Failed to decode transaction: %v", err)
    return
}
```

### 3. **Use Appropriate Timeouts**
```go
cfg := seth.DefaultConfig("http://localhost:8545", []string{privateKey})
cfg.RPCTimeout = 30 * time.Second
cfg.TxTimeout = 5 * time.Minute
```

### 4. **Monitor Gas Prices**
```go
// Check gas price before sending transaction
gasPrice, err := client.SuggestGasPrice(context.Background())
if err != nil {
    log.Fatal(err)
}

if gasPrice.Cmp(big.NewInt(100000000000)) > 0 {
    log.Printf("High gas price: %s", gasPrice.String())
}
```

### 5. **Use Tracing for Debugging**
```go
// Enable tracing for complex transactions
trace, err := client.TraceTransaction(context.Background(), txHash)
if err != nil {
    log.Printf("Failed to trace transaction: %v", err)
} else {
    log.Printf("Transaction trace: %+v", trace)
}
```

## Troubleshooting

### Common Issues

#### Transaction Decoding Fails
- Ensure you have the correct ABI for the contract
- Check if the contract address is correct
- Verify the transaction hash is valid

#### Gas Estimation Issues
- Check if the RPC node supports gas estimation
- Verify the transaction parameters are valid
- Consider using manual gas price estimation

#### Tracing Not Available
- Ensure your RPC node supports tracing
- Check if tracing is enabled on the node
- Use a different RPC endpoint if needed

### Debug Mode

Enable debug logging:

```bash
export SETH_LOG_LEVEL=debug
go test -v
```

## API Reference

For detailed API documentation, see the [Go package documentation](https://pkg.go.dev/github.com/smartcontractkit/chainlink-testing-framework/seth).

## Contributing

We welcome contributions! Please see the [Contributing Guide](../../Contributing/Overview) for details on:

- Code of Conduct
- Development Setup
- Testing Guidelines
- Pull Request Process

## License

This project is licensed under the MIT License - see the [LICENSE](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/LICENSE) file for details. 