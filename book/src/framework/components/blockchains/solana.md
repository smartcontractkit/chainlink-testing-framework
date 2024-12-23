# Solana Blockchain Client

Since `Solana` doesn't have official image for `arm64` we built it, images we use are:
```
amd64 solanalabs/solana:v1.18.26 - used in CI
arm64 f4hrenh9it/solana:latest - used locally
```

## Configuration
```toml
[blockchain_a]
  type = "solana"
  # public key for mint
  public_key = "9n1pyVGGo6V4mpiSDMVay5As9NurEkY283wwRk1Kto2C"
  # contracts directory, programs
  contracts_dir = "."
  # optional, in case you need some custom image
  # image = "solanalabs/solana:v1.18.26"
```

## Usage
```golang
package examples

import (
	"context"
	"fmt"
	"github.com/blocto/solana-go-sdk/client"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/stretchr/testify/require"
	"testing"
)

type CfgSolana struct {
	BlockchainA *blockchain.Input `toml:"blockchain_a" validate:"required"`
}

func TestSolanaSmoke(t *testing.T) {
	in, err := framework.Load[CfgSolana](t)
	require.NoError(t, err)

	bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
	require.NoError(t, err)

	t.Run("test something", func(t *testing.T) {
		// use internal URL to connect chainlink nodes
		_ = bc.Nodes[0].DockerInternalHTTPUrl
		// use host URL to deploy contracts
		c := client.NewClient(bc.Nodes[0].HostHTTPUrl)
		latestSlot, err := c.GetSlotWithConfig(context.Background(), client.GetSlotConfig{Commitment: "processed"})
		require.NoError(t, err)
		fmt.Printf("Latest slot: %v\n", latestSlot)
	})
}
```

## Test Private Keys

```
Public: 9n1pyVGGo6V4mpiSDMVay5As9NurEkY283wwRk1Kto2C
Private: [11,2,35,236,230,251,215,68,220,208,166,157,229,181,164,26,150,230,218,229,41,20,235,80,183,97,20,117,191,159,228,243,130,101,145,43,51,163,139,142,11,174,113,54,206,213,188,127,131,147,154,31,176,81,181,147,78,226,25,216,193,243,136,149]
```

