# Sui Blockchain Client

API is available on [localhost:9000](http://localhost:9000)

## Configuration

```toml
[blockchain_a]
  type = "sui"
  image = "mysten/sui-tools:mainnet" # if omitted default is mysten/sui-tools:devnet
  contracts_dir = "$your_dir"
```

## Usage

```golang
package examples

import (
	"context"
	"fmt"
	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/signer"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/stretchr/testify/require"
	"testing"
)

type CfgSui struct {
	BlockchainA *blockchain.Input `toml:"blockchain_a" validate:"required"`
}

func TestSuiSmoke(t *testing.T) {
	in, err := framework.Load[CfgSui](t)
	require.NoError(t, err)

	bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
	require.NoError(t, err)

	// network is already funded, here are the keys
	_ = bc.NetworkSpecificData.SuiAccount.Mnemonic
	_ = bc.NetworkSpecificData.SuiAccount.PublicBase64Key
	_ = bc.NetworkSpecificData.SuiAccount.SuiAddress

	// execute any additional commands, to deploy contracts or set up
	_, err = framework.ExecContainer(bc.ContainerName, []string{"ls", "-lah"})
	require.NoError(t, err)

	t.Run("test something", func(t *testing.T) {
		// use internal URL to connect Chainlink nodes
		_ = bc.Nodes[0].InternalHTTPUrl
		// use host URL to interact
		_ = bc.Nodes[0].ExternalHTTPUrl

		cli := sui.NewSuiClient(bc.Nodes[0].ExternalHTTPUrl)

		signerAccount, err := signer.NewSignertWithMnemonic(bc.NetworkSpecificData.SuiAccount.Mnemonic)
		require.NoError(t, err)
		rsp, err := cli.SuiXGetAllBalance(context.Background(), models.SuiXGetAllBalanceRequest{
			Owner: signerAccount.Address,
		})
		require.NoError(t, err)
		fmt.Printf("My funds: %v\n", rsp)
	})
}
```

## Test Private Keys

Since Sui doesn't have official local development chain we are using real node and generating mnemonic at start then funding that account through internal faucet, see
```
	// network is already funded, here are the keys
	_ = bc.NetworkSpecificData.SuiAccount.Mnemonic
	_ = bc.NetworkSpecificData.SuiAccount.PublicBase64Key
	_ = bc.NetworkSpecificData.SuiAccount.SuiAddress
```