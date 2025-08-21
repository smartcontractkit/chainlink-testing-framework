# Chainlink Node Set Environment Test

Let's create a full-fledged set of Chainlink nodes connected to some blockchain.

Create a configuration file `smoke.toml`
```toml
[blockchain_a]
  docker_cmd_params = ["-b", "1"]
  type = "anvil"

[data_provider]
  port = 9111

[[nodesets]]
  name = "don"
  nodes = 5
  override_mode = "all"

  [nodesets.db]
    image = "postgres:12.0"

  [[nodesets.node_specs]]

    [nodesets.node_specs.node]
      image = "public.ecr.aws/chainlink/chainlink:v2.17.0"
```

Create a file `smoke_test.go`
```golang
package examples

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
)

type Cfg struct {
	BlockchainA        *blockchain.Input `toml:"blockchain_a" validate:"required"`
	MockedDataProvider *fake.Input       `toml:"data_provider" validate:"required"`
	NodeSets           []*ns.Input       `toml:"nodesets" validate:"required"`
}

func TestNodeSet(t *testing.T) {
	in, err := framework.Load[Cfg](t)
	require.NoError(t, err)

	bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
	require.NoError(t, err)
	_, err = fake.NewFakeDataProvider(in.MockedDataProvider)
	require.NoError(t, err)
	out, err := ns.NewSharedDBNodeSet(in.NodeSets[0], bc)
	require.NoError(t, err)

	t.Run("test something", func(t *testing.T) {
		for _, n := range out.CLNodes {
			require.NotEmpty(t, n.Node.ExternalURL)
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


