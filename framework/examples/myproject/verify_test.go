package examples

import (
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/examples/example_components/onchain"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	"github.com/stretchr/testify/require"
	"testing"
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

	// connect 2 clients
	scSrc, err := seth.NewClientBuilder().
		WithRpcUrl(bc.Nodes[0].HostWSUrl).
		WithGasPriceEstimations(true, 0, seth.Priority_Fast).
		WithTracing(seth.TracingLevel_All, []string{seth.TraceOutput_Console}).
		WithPrivateKeys([]string{blockchain.DefaultAnvilPrivateKey}).
		Build()
	require.NoError(t, err)
	in.ContractsSrc.URL = bc.Nodes[0].HostWSUrl
	c, err := onchain.NewProductOnChainDeployment(scSrc, in.ContractsSrc)
	require.NoError(t, err)

	t.Run("test something", func(t *testing.T) {
		err := framework.VerifyContract(
			bc.Nodes[0].HostHTTPUrl,
			c.Addresses[0].String(),
			"/Users/fahrenheit/GolandProjects/chainlink-testing-framework/framework/examples/myproject/example_components/onchain/src/Counter.sol",
			"Counter",
			"http://localhost",
		)
		require.NoError(t, err)
	})
}
