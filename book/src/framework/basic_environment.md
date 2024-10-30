# Chainlink Cluster (NodeSet) Environment Test

Create a configuration file `smoke.toml`
```toml
funds_eth = 30.0

[blockchain_a]
  chain_id = "31337"
  image = "f4hrenh9it/foundry:latest"
  port = "8545"
  type = "anvil"

[contracts]

[data_provider]
  port = 9111

[nodeset]
  nodes = 5
  override_mode = "all"

  [[nodeset.node_specs]]

    [nodeset.node_specs.db]
      image = "postgres:15.6"
      pull_image = true

    [nodeset.node_specs.node]
      image = "public.ecr.aws/chainlink/chainlink:v2.17.0"
      pull_image = true

```

Create a file `smoke_test.go`
```golang
package yourpackage_test

import (
	"fmt"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink/e2e/capabilities/components/onchain"
	"github.com/stretchr/testify/require"
	"testing"
)

type Config struct {
	FundingETH         float64           `toml:"funds_eth"`
	BlockchainA        *blockchain.Input `toml:"blockchain_a" validate:"required"`
	Contracts          *onchain.Input    `toml:"contracts" validate:"required"`
	MockerDataProvider *fake.Input       `toml:"data_provider" validate:"required"`
	NodeSet            *ns.Input         `toml:"nodeset" validate:"required"`
}

func TestNodeSet(t *testing.T) {
	in, err := framework.Load[Config](t)
	require.NoError(t, err)

	// deploy docker test environment
	bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
	require.NoError(t, err)
	dp, err := fake.NewFakeDataProvider(in.MockerDataProvider)
	require.NoError(t, err)
	out, err := ns.NewSharedDBNodeSet(in.NodeSet, bc, dp.BaseURLDocker)
	require.NoError(t, err)
}
```

Run it
```bash
go test -v -run TestNodeSet
```

