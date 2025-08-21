package examples

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/clnode"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	testToken "github.com/smartcontractkit/chainlink-testing-framework/framework/examples/example_components/gethwrappers"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/examples/example_components/onchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/rpc"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
)

type CfgForkChainsOffChain struct {
	ContractsSrc  *onchain.Input    `toml:"contracts_src" validate:"required"`
	BlockchainSrc *blockchain.Input `toml:"blockchain_src" validate:"required"`
	// off-chain components
	NodeSets []*ns.Input `toml:"nodesets" validate:"required"`
}

func TestOffChainAndFork(t *testing.T) {
	in, err := framework.Load[CfgForkChainsOffChain](t)
	require.NoError(t, err)

	// spin up 2 anvils
	bcSrc, err := blockchain.NewBlockchainNetwork(in.BlockchainSrc)
	require.NoError(t, err)

	// create configs for 2 EVM networks
	srcNetworkCfg, err := clnode.NewNetworkCfg(&clnode.EVMNetworkConfig{
		MinIncomingConfirmations: 1,
		MinContractPayment:       "0.00001 link",
		ChainID:                  bcSrc.ChainID,
		EVMNodes: []*clnode.EVMNode{
			{
				SendOnly: false,
				Order:    100,
			},
		},
	}, bcSrc)
	in.NodeSets[0].NodeSpecs[0].Node.TestConfigOverrides = srcNetworkCfg

	_, err = ns.NewSharedDBNodeSet(in.NodeSets[0], bcSrc)
	require.NoError(t, err)

	scSrc, err := seth.NewClientBuilder().
		WithRpcUrl(bcSrc.Nodes[0].ExternalWSUrl).
		WithGasPriceEstimations(true, 0, seth.Priority_Fast).
		WithTracing(seth.TracingLevel_All, []string{seth.TraceOutput_Console}).
		WithPrivateKeys([]string{blockchain.DefaultAnvilPrivateKey}).
		Build()
	require.NoError(t, err)

	in.ContractsSrc.URL = bcSrc.Nodes[0].ExternalWSUrl
	contractsSrc, err := onchain.NewProductOnChainDeployment(scSrc, in.ContractsSrc)
	require.NoError(t, err)

	t.Run("test some contracts with fork", func(t *testing.T) {
		i, err := testToken.NewBurnMintERC677(contractsSrc.Addresses[0], scSrc.Client)
		require.NoError(t, err)
		balance, err := i.BalanceOf(scSrc.NewCallOpts(), contractsSrc.Addresses[0])
		require.NoError(t, err)
		fmt.Println(balance)

		// Use anvil methods, see https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/framework/rpc/rpc.go
		_ = rpc.New(bcSrc.Nodes[0].ExternalHTTPUrl, nil)
	})
}
