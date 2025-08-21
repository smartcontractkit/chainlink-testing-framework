package examples

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/examples/example_components/onchain"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
)

type VerifyCfg struct {
	ContractsSrc *onchain.Input    `toml:"contracts_src" validate:"required"`
	BlockchainA  *blockchain.Input `toml:"blockchain_a" validate:"required"`
}

func TestVerify(t *testing.T) {
	in, err := framework.Load[VerifyCfg](t)
	require.NoError(t, err)

	bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
	require.NoError(t, err)

	scSrc, err := seth.NewClientBuilder().
		WithRpcUrl(bc.Nodes[0].ExternalWSUrl).
		WithPrivateKeys([]string{blockchain.DefaultAnvilPrivateKey}).
		Build()
	require.NoError(t, err)
	in.ContractsSrc.URL = bc.Nodes[0].ExternalWSUrl
	c, err := onchain.NewCounterDeployment(scSrc, in.ContractsSrc)
	require.NoError(t, err)

	t.Run("verify contract and test with debug", func(t *testing.T) {
		// give Blockscout some time to index your transactions
		// there is no API for that
		time.Sleep(10 * time.Second)
		err := blockchain.VerifyContract(bc, c.Addresses[0].String(),
			"example_components/onchain",
			"src/Counter.sol",
			"Counter",
			"0.8.13",
		)
		require.NoError(t, err)
	})
}
