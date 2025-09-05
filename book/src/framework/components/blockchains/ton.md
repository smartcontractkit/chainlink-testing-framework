# TON Blockchain Client

TON (The Open Network) support in the framework utilizes MyLocalTon Docker environment to provide a minimal local TON blockchain for testing purposes.

## Configuration

```toml
[blockchain_a]
  type = "ton"
  image = "ghcr.io/neodix42/mylocalton-docker:latest"
  port = "8000"
```

## Genesis Container Parameters

The genesis container supports additional environment variables that can be configured through the `custom_env` field. These parameters allow you to customize the blockchain behavior:

```toml
[blockchain_a]
  type = "ton"
  image = "ghcr.io/neodix42/mylocalton-docker:latest"
  port = "8000"
  
  [blockchain_a.custom_env]
  VERSION_CAPABILITIES = "11"
```

The custom_env parameters will override the default genesis container environment variables, allowing you to customize blockchain configuration as needed.
More info on parameters can be found here <https://github.com/neodix42/mylocalton-docker/wiki/Genesis-setup-parameters>.

## Network Configuration

The framework provides seamless access to the TON network configuration by embedding the config URL directly in the node URLs. The `ExternalHTTPUrl` and `InternalHTTPUrl` include the full path to `localhost.global.config.json`, which can be used directly with `liteclient.GetConfigFromUrl()` without additional URL formatting.

## Default Ports

The TON implementation exposes essential services:

* TON Simple HTTP Server: Port 8000
* TON Lite Server: Port derived from base port + 100

> Note: `tonutils-go` library is used for TON blockchain interactions, which requires a TON Lite Server connection. The framework embeds the config URL directly in the node URLs for convenient access to the global configuration file needed by `tonutils-go`.

## Usage

```go
package examples

import (
  "strings"
  "testing"

  "github.com/stretchr/testify/require"
  "github.com/xssnick/tonutils-go/liteclient"
  "github.com/xssnick/tonutils-go/ton"
  "github.com/xssnick/tonutils-go/ton/wallet"

  "github.com/smartcontractkit/chainlink-testing-framework/framework"
  "github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
)

type CfgTon struct {
  BlockchainA *blockchain.Input `toml:"blockchain_a" validate:"required"`
}

func TestTonSmoke(t *testing.T) {
  in, err := framework.Load[CfgTon](t)
  require.NoError(t, err)

  bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
  require.NoError(t, err)
  // we can also explicitly terminate the container after the test
  defer bc.Container.Terminate(t.Context())

  var client ton.APIClientWrapped

  t.Run("setup:connect", func(t *testing.T) {
    // Create a connection pool
    connectionPool := liteclient.NewConnectionPool()
    
    // Get the network configuration directly from the embedded config URL
    // The ExternalHTTPUrl already includes the full path to localhost.global.config.json
    cfg, cferr := liteclient.GetConfigFromUrl(t.Context(), bc.Nodes[0].ExternalHTTPUrl)
    require.NoError(t, cferr, "Failed to get config from URL")
    
    // Add connections from the config
    caerr := connectionPool.AddConnectionsFromConfig(t.Context(), cfg)
    require.NoError(t, caerr, "Failed to add connections from config")
    
    // Create an API client with retry functionality
    client = ton.NewAPIClient(connectionPool).WithRetry()

    t.Run("setup:faucet", func(t *testing.T) {
      // Create a wallet from the pre-funded high-load wallet seed
      rawHlWallet, err := wallet.FromSeed(client, strings.Fields(blockchain.DefaultTonHlWalletMnemonic), wallet.HighloadV2Verified)
      require.NoError(t, err, "failed to create highload wallet")
   
      // Create a workchain -1 (masterchain) wallet
      mcFunderWallet, err := wallet.FromPrivateKeyWithOptions(client, rawHlWallet.PrivateKey(), wallet.HighloadV2Verified, wallet.WithWorkchain(-1))
      require.NoError(t, err, "failed to create highload wallet")
   
      // Get subwallet with ID 42
      funder, err := mcFunderWallet.GetSubwallet(uint32(42))
      require.NoError(t, err, "failed to get highload subwallet")

      // Verify the funder address matches the expected default
      require.Equal(t, funder.Address().StringRaw(), blockchain.DefaultTonHlWalletAddress, "funder address mismatch")

      // Check the funder balance
      master, err := client.GetMasterchainInfo(t.Context())
      require.NoError(t, err, "failed to get masterchain info for funder balance check")
      funderBalance, err := funder.GetBalance(t.Context(), master)
      require.NoError(t, err, "failed to get funder balance")
      require.Equal(t, funderBalance.Nano().String(), "1000000000000000", "funder balance mismatch")
    })
  })
}
```

## Test Private Keys

The framework includes a pre-funded high-load wallet for testing purposes. This wallet type can send up to 254 messages per 1 external message, making it efficient for test scenarios.

Default High-Load Wallet:

```shell
Address: -1:5ee77ced0b7ae6ef88ab3f4350d8872c64667ffbe76073455215d3cdfab3294b
Mnemonic: twenty unfair stay entry during please water april fabric morning length lumber style tomorrow melody similar forum width ride render void rather custom coin
```

## Available Pre-funded Wallets

MyLocalTon Docker environment comes with several pre-funded wallets that can be used for testing:

1. Genesis Wallet (V3R2, WalletId: 42)
2. Validator Wallets (1-5) (V3R2, WalletId: 42)
3. Faucet Wallet (V3R2, WalletId: 42, Balance: 1 million TON)
4. Faucet Highload Wallet (Highload V2, QueryId: 0, Balance: 1 million TON)
5. Basechain Faucet Wallet (V3R2, WalletId: 42, Balance: 1 million TON)
6. Basechain Faucet Highload Wallet (Highload V2, QueryId: 0, Balance: 1 million TON)

For the complete list of addresses and mnemonics, refer to the [MyLocalTon Docker documentation](https://github.com/neodix42/mylocalton-docker).
