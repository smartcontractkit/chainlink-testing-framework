# EVM Blockchain Clients

We support 3 EVM clients at the moment: [Geth](https://geth.ethereum.org/docs/fundamentals/command-line-options), [Anvil](https://book.getfoundry.sh/anvil/) and [Besu](https://besu.hyperledger.org/)

## Configuration
```toml
[blockchain_a]
  # Blockchain node type, can be "anvil", "geth" or "besu
  type = "anvil"
  # Chain ID
  chain_id = "1337"
  # Anvil command line params, ex.: docker_cmd_params = ['--block-time=1', '...']
  docker_cmd_params = []
  # Docker image and tag
  image = "ghcr.io/foundry-rs/foundry:stable"
  # External port to expose (HTTP API)
  port = "8545"
  # External port to expose (WS API)
  port_ws = "8546"
  # Pulls the image every time if set to 'true', used like that in CI. Can be set to 'false' to speed up local runs
  pull_image = false

  # Outputs are the results of deploying a component that can be used by another component
  [blockchain_a.out]
    # If 'use_cache' equals 'true' we skip component setup when we run the test and return the outputs
    use_cache = true
    # Chain ID
    chain_id = "1337"
    # Chain family, "evm", "solana", "cosmos", "op", "arb"
    family = "evm"

    [[blockchain_a.out.nodes]]
      # URLs to access the node(s) inside docker network, used by other components
      internal_http_url = "http://anvil-14411:8545"
      internal_ws_url = "ws://anvil-14411:8545"
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

## Test Private Keys

For `Geth` we use
```
	DefaultGethPrivateKey = `ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80`
	DefaultGethPublicKey = `0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266`
```

For `Anvil` we use
```
	DefaultAnvilPrivateKey = `ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80`
	AnvilPrivateKey1       = `0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d`
	AnvilPrivateKey2       = `0x5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a`
	AnvilPrivateKey3       = `0x7c852118294e51e653712a81e05800f419141751be58f605c371e15141b007a6`
	AnvilPrivateKey4       = `0x47e179ec197488593b187f80a00eb0da91f1b9d0b13f8733639f19c30a34926a`
	
	DefaultAnvilPublicKey  = `0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266`
	AnvilPublicKey1        = `0x70997970C51812dc3A010C7d01b50e0d17dc79C8`
	AnvilPublicKey2        = `0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC`
	AnvilPublicKey3        = `0x90F79bf6EB2c4f870365E785982E1f101E93b906`
	AnvilPublicKey4        = `0x15d34AAf54267DB7D7c367839AAf71A00a2C6A65`
```

For `Besu` keys are
```
	DefaultBesuPrivateKey1 = "8f2a55949038a9610f50fb23b5883af3b4ecb3c3bb792cbcefbd1542c692be63"
	DefaultBesuPrivateKey2 = "c87509a1c067bbde78beb793e6fa76530b6382a4c0241e5e4a9ec0a0f44dc0d3"
	DefaultBesuPrivateKey3 = "ae6ae8e5ccbfb04590405997ee2d52d2b330726137b875053c36d94e974d162f"
	
	DefaultBesuPublicKey1 = "0xfe3b557e8fb62b89f4916b721be55ceb828dbd73"
	DefaultBesuPublicKey2 = "0x627306090abaB3A6e1400e9345bC60c78a8BEf57"
	DefaultBesuPublicKey3 = "0xf17f52151EbEF6C7334FAD080c5704D77216b732"
```

More docs for `Besu` can be found [here](https://besu.hyperledger.org/private-networks/reference/accounts-for-testing)

