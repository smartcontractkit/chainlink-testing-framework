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
	ContractsDst  *onchain.Input    `toml:"contracts_dst" validate:"required"`
	BlockchainSrc *blockchain.Input `toml:"blockchain_src" validate:"required"`
	BlockchainDst *blockchain.Input `toml:"blockchain_dst" validate:"required"`
	// off-chain components
	NodeSet *ns.Input `toml:"nodeset" validate:"required"`
}

func TestOffChainAndFork(t *testing.T) {
	in, err := framework.Load[CfgForkChainsOffChain](t)
	require.NoError(t, err)

	// spin up 2 anvils
	bcSrc, err := blockchain.NewBlockchainNetwork(in.BlockchainSrc)
	require.NoError(t, err)

	bcDst, err := blockchain.NewBlockchainNetwork(in.BlockchainDst)
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
	dstNetworkConfig, err := clnode.NewNetworkCfg(&clnode.EVMNetworkConfig{
		MinIncomingConfirmations: 1,
		MinContractPayment:       "0.00001 link",
		ChainID:                  bcSrc.ChainID,
		EVMNodes: []*clnode.EVMNode{
			{
				SendOnly: false,
				Order:    100,
			},
		},
	}, bcDst)
	// override the configuration
	in.NodeSet.NodeSpecs[0].Node.TestConfigOverrides = srcNetworkCfg + dstNetworkConfig

	// create a node set
	_, err = ns.NewSharedDBNodeSet(in.NodeSet, bcSrc)
	require.NoError(t, err)

	// connect 2 clients
	scSrc, err := seth.NewClientBuilder().
		WithRpcUrl(bcSrc.Nodes[0].ExternalWSUrl).
		WithGasPriceEstimations(true, 0, seth.Priority_Fast).
		WithTracing(seth.TracingLevel_All, []string{seth.TraceOutput_Console}).
		WithPrivateKeys([]string{blockchain.DefaultAnvilPrivateKey}).
		Build()
	require.NoError(t, err)
	scDst, err := seth.NewClientBuilder().
		WithRpcUrl(bcDst.Nodes[0].ExternalWSUrl).
		WithGasPriceEstimations(true, 0, seth.Priority_Fast).
		WithTracing(seth.TracingLevel_All, []string{seth.TraceOutput_Console}).
		WithPrivateKeys([]string{blockchain.DefaultAnvilPrivateKey}).
		Build()
	require.NoError(t, err)

	// deploy 2 example product contracts
	// you should replace it with chainlink-deployments
	in.ContractsSrc.URL = bcSrc.Nodes[0].ExternalWSUrl
	contractsSrc, err := onchain.NewProductOnChainDeployment(scSrc, in.ContractsSrc)
	require.NoError(t, err)
	in.ContractsDst.URL = bcDst.Nodes[0].ExternalWSUrl
	contractsDst, err := onchain.NewProductOnChainDeployment(scDst, in.ContractsDst)
	require.NoError(t, err)

	t.Run("test some contracts with fork", func(t *testing.T) {
		// interact with a source chain
		i, err := testToken.NewBurnMintERC677(contractsSrc.Addresses[0], scSrc.Client)
		require.NoError(t, err)
		balance, err := i.BalanceOf(scSrc.NewCallOpts(), contractsSrc.Addresses[0])
		require.NoError(t, err)
		fmt.Println(balance)

		// interact with a destination chain
		i, err = testToken.NewBurnMintERC677(contractsDst.Addresses[0], scDst.Client)
		require.NoError(t, err)
		balance, err = i.BalanceOf(scDst.NewCallOpts(), contractsDst.Addresses[0])
		require.NoError(t, err)
		fmt.Println(balance)

		// Use anvil methods, see https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/framework/rpc/rpc.go
		_ = rpc.New(bcSrc.Nodes[0].ExternalHTTPUrl, nil)
		_ = rpc.New(bcDst.Nodes[0].ExternalHTTPUrl, nil)
	})
}
