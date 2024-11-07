# Chainlink Node Set Environment Test

Let's create a full-fledged set of Chainlink nodes connected to some blockchain.

Create a configuration file `smoke.toml`
```toml
[blockchain_a]
  chain_id = "31337"
  image = "f4hrenh9it/foundry:latest"
  port = "8545"
  type = "anvil"

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
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/stretchr/testify/require"
	"testing"
)

type Config struct {
	BlockchainA        *blockchain.Input `toml:"blockchain_a" validate:"required"`
	MockerDataProvider *fake.Input       `toml:"data_provider" validate:"required"`
	NodeSet            *ns.Input         `toml:"nodeset" validate:"required"`
}

func TestNodeSet(t *testing.T) {
	in, err := framework.Load[Config](t)
	require.NoError(t, err)

	bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
	require.NoError(t, err)
	dp, err := fake.NewFakeDataProvider(in.MockerDataProvider)
	require.NoError(t, err)
	out, err := ns.NewSharedDBNodeSet(in.NodeSet, bc, dp.BaseURLDocker)
	require.NoError(t, err)

	t.Run("test something", func(t *testing.T) {
		for _, n := range out.CLNodes {
			require.NotEmpty(t, n.Node.HostURL)
			require.NotEmpty(t, n.Node.HostP2PURL)
		}
	})
}
```

Run it
```bash
CTF_CONFIGS=smoke.toml go test -v -run TestNodeSet
```

Check the logs to access the UI
```bash
12:41AM INF UI=["http://127.0.0.1:10000","http://127.0.0.1:10001", ...]
```

Use credentials to authorize:
```
notreal@fakeemail.ch
fj293fbBnlQ!f9vNs
```

Summary:
- We deployed fully-fledged set of Chainlink nodes connected to some blockchain and faked external data provider
- We explored the Chainlink node UI


