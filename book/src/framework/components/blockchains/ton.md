# TON Blockchain Client

TON (The Open Network) support in the framework utilizes MyLocalTon Docker Compose environment to provide a local TON blockchain for testing purposes.

## Configuration

```toml
[blockchain_a]
  type = "ton"
  # By default uses MyLocalTon Docker Compose file
  image = "https://raw.githubusercontent.com/neodix42/mylocalton-docker/main/docker-compose.yaml"
  # Optional: Specify only core services needed for testing (useful in CI environments)
  ton_core_services = [
    "genesis",
    "tonhttpapi",
    "event-cache",
    "index-postgres", 
    "index-worker", 
    "index-api"
  ]
```

## Default Ports

The TON implementation exposes several services:

- TON Lite Server: Port 40004
- TON HTTP API: Port 8081
- TON Simple HTTP Server: Port 8000
- TON Explorer: Port 8080

> **Note**: By default, only the lite client service is exposed externally. Other services may need additional configuration to be accessible outside the Docker network.

## Validator Configuration

By default, the MyLocalTon environment starts with only one validator enabled. If multiple validators are needed (up to 6 are supported), the Docker Compose file must be provided with modified version with corresponding service definition in toml file before starting the environment.

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

	var client ton.APIClientWrapped

	t.Run("setup:connect", func(t *testing.T) {
		// Create a connection pool
		connectionPool := liteclient.NewConnectionPool()
		
		// Get the network configuration from the global config URL
		cfg, cferr := liteclient.GetConfigFromUrl(t.Context(), fmt.Sprintf("http://%s/localhost.global.config.json", bc.Nodes[0].ExternalHTTPUrl))
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
```
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
