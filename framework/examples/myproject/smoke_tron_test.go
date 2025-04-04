package examples

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
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
