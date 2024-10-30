# Anvil
[Anvil](https://book.getfoundry.sh/anvil/) is a Foundry local EVM blockchain simulator

Use `docker_cmd_params = ['--block-time=1', '...']` to provide more params

## Configuration
```toml
[blockchain_a]
  chain_id = "31337"
  docker_cmd_params = []
  image = "ghcr.io/gakonst/foundry:latest"
  port = "8545"
  pull_image = false
  type = "anvil"

  [blockchain_a.out]
    chain_id = "31337"
    use_cache = true

    [[blockchain_a.out.nodes]]
      docker_internal_http_url = "http://anvil-14411:8545"
      docker_internal_ws_url = "ws://anvil-14411:8545"
      http_url = "http://127.0.0.1:33955"
      ws_url = "ws://127.0.0.1:33955"
```

## Usage
```golang
package my_test

import (
	"os"
	"testing"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/stretchr/testify/require"
)

type Config struct {
	BlockchainA        *blockchain.Input `toml:"blockchain_a" validate:"required"`
}

func TestDON(t *testing.T) {
	in, err := framework.Load[Config](t)
	require.NoError(t, err)
	pkey := os.Getenv("PRIVATE_KEY")

	// deploy anvil blockchain simulator
	bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
	require.NoError(t, err)
}
```

