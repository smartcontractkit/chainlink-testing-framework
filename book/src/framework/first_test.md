# Writing your first test

The Chainlink Testing Framework (CTF) is a modular, data-driven tool that lets you explicitly define and configure various Chainlink components.

Let's spin up a simple component.

Create your configuration in `smoke.toml`
```toml
[blockchain_a]
  # choose between "anvil", "geth" or "besu"
  type = "anvil"
```

Create your test in `smoke_test.go`
```golang
package mymodule_test

import (
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/stretchr/testify/require"
	"testing"
)

type Config struct {
	BlockchainA *blockchain.Input `toml:"blockchain_a" validate:"required"`
}

func TestMe(t *testing.T) {
	in, err := framework.Load[Config](t)
	require.NoError(t, err)

	bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
	require.NoError(t, err)

	t.Run("test something", func(t *testing.T) {
		require.NotEmpty(t, bc.Nodes[0].ExternalHTTPUrl)
	})
}
```

Run the test
```bash
CTF_CONFIGS=smoke.toml go test -v -run TestMe
```

Remove containers (read more about cleanup [here](components/cleanup.md))
```
ctf d rm
```

Summary:
- We defined configuration for `BlockchainNetwork`
- We've used one CTF component in test and checked if it's working

Now let's connect the [Chainlink](./connecting_chainlink_node.md) node!

Learn more about [EVM Client](./components/blockchains/evm.md) component.