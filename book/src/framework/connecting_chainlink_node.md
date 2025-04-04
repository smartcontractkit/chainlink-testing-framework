# Connecting Chainlink Node

Now let's have an example of Chainlink node connected to some local blockchain.

Create your configuration in `smoke.toml`
```toml
[blockchain_a]
  type = "anvil"
  docker_cmd_params = ["-b", "1"]

[cl_node]

  [cl_node.db]
    image = "postgres:12.0"

  [cl_node.node]
    image = "public.ecr.aws/chainlink/chainlink:v2.17.0"
```

Create your test in `smoke_test.go`
```golang

package capabilities_test

import (
	"fmt"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/clnode"
	"github.com/stretchr/testify/require"
	"testing"
)

type Config struct {
	BlockchainA *blockchain.Input `toml:"blockchain_a" validate:"required"`
	CLNode      *clnode.Input     `toml:"cl_node" validate:"required"`
}

func TestNode(t *testing.T) {
	in, err := framework.Load[Config](t)
	require.NoError(t, err)

	bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
	require.NoError(t, err)

	networkCfg, err := clnode.NewNetworkCfgOneNetworkAllNodes(bc)
	require.NoError(t, err)
	in.CLNode.Node.TestConfigOverrides = networkCfg

	output, err := clnode.NewNodeWithDB(in.CLNode)
	require.NoError(t, err)

	t.Run("test something", func(t *testing.T) {
		fmt.Printf("node url: %s\n", output.Node.ExternalURL)
		require.NotEmpty(t, output.Node.ExternalURL)
	})
}


```

Select your configuration by setting `CTF_CONFIGS=smoke.toml` and run it
```bash
go test -v -run TestNode
```

Check `node url: ...` in logs, open it and login using default credentials:
```
notreal@fakeemail.ch
fj293fbBnlQ!f9vNs
```

Summary:
- We defined configuration for `BlockchainNetwork` and `NodeWithDB` (Chainlink + PostgreSQL)
- We connected them together by creating common network config in `NewNetworkCfgOneNetworkAllNodes`
- We explored the Chainlink node UI

Let's proceed with another example of [using node sets](./nodeset_environment)