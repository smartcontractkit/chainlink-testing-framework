# Anvil
[Anvil](https://book.getfoundry.sh/anvil/) is a Foundry local EVM blockchain simulator

## Configuration
```toml
[blockchain_a]
  # Blockchain node type, can be "anvil" or "geth"
  type = "anvil"
  # Chain ID
  chain_id = "31337"
  # Anvil command line params, ex.: docker_cmd_params = ['--block-time=1', '...']
  docker_cmd_params = []
  # Docker image and tag
  image = "f4hrenh9it/foundry:latest"
  # External port to expose
  port = "8545"
  # Pulls the image every time if set to 'true', used like that in CI. Can be set to 'false' to speed up local runs
  pull_image = false

  # Outputs are the results of deploying a component that can be used by another component
  [blockchain_a.out]
    chain_id = "31337"
    # If 'use_cache' equals 'true' we skip component setup when we run the test and return the outputs
    use_cache = true

    [[blockchain_a.out.nodes]]
      # URLs to access the node(s) inside docker network, used by other components
      docker_internal_http_url = "http://anvil-14411:8545"
      docker_internal_ws_url = "ws://anvil-14411:8545"
      # URLs to access the node(s) on your host machine or in CI
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

	// deploy anvil blockchain simulator
	bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
	require.NoError(t, err)
}
```

