# TRON Blockchain Client

## Configuration
```toml
[blockchain_a]
  type = "tron"
  # image = "tronbox/tre:1.0.3" is default image
```
Default port is `9090`

## Usage
```golang
package examples

import (
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/stretchr/testify/require"
	"testing"
)

type CfgTron struct {
	BlockchainA *blockchain.Input `toml:"blockchain_a" validate:"required"`
}

func TestTRONSmoke(t *testing.T) {
	in, err := framework.Load[CfgTron](t)
	require.NoError(t, err)

	bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
	require.NoError(t, err)

	// all private keys are funded
	_ = blockchain.TRONAccounts.PrivateKeys[0]

	t.Run("test something", func(t *testing.T) {
		// use internal URL to connect Chainlink nodes
		_ = bc.Nodes[0].InternalHTTPUrl
		// use host URL to interact
		_ = bc.Nodes[0].ExternalHTTPUrl

		// use bc.Nodes[0].ExternalHTTPUrl + "/wallet" to access full node
		// use bc.Nodes[0].ExternalHTTPUrl + "/walletsolidity" to access Solidity node
	})
}
```

## More info

Follow the [guide](https://developers.tron.network/reference/tronbox-quickstart) if you want to work with `TRONBox` environment via JS

## Golang HTTP Client

TRON doesn't have any library to interact with it in `Golang` but we maintain our internal fork [here](https://github.com/smartcontractkit/chainlink-internal-integrations/tree/69e35041cdea0bc38ddf642aa93fd3cc3fb5d0d9/tron/relayer/gotron-sdk)

Check TRON [HTTP API](https://tronprotocol.github.io/documentation-en/api/http/)

Full node is on `:9090/wallet`
```
curl -X POST http://127.0.0.1:9090/wallet/createtransaction -d '{                                                           
    "owner_address": "TRGhNNfnmgLegT4zHNjEqDSADjgmnHvubJ",
    "to_address": "TJCnKsPa7y5okkXvQAidZBzqx3QyQ6sxMW",
    "amount": 1000000,
    "visible": true
}'
```

Solidity node is on `:9090/walletsolidity`
```
curl -X POST  http://127.0.0.1:9090/walletsolidity/getaccount -d '{"address": "41E552F6487585C2B58BC2C9BB4492BC1F17132CD0"}'
```
