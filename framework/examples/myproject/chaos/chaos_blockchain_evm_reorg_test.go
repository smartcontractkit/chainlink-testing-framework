package chaos

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/clnode"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/rpc"
)

type CfgReorgTwoChains struct {
	BlockchainA        *blockchain.Input `toml:"blockchain_a" validate:"required"`
	BlockchainB        *blockchain.Input `toml:"blockchain_b" validate:"required"`
	MockerDataProvider *fake.Input       `toml:"data_provider" validate:"required"`
	NodeSets           []*ns.Input       `toml:"nodesets" validate:"required"`
}

func TestBlockchainReorgChaos(t *testing.T) {
	in, err := framework.Load[CfgReorgTwoChains](t)
	require.NoError(t, err)

	// Can replace deployments with CRIB here

	bcA, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
	require.NoError(t, err)
	bcB, err := blockchain.NewBlockchainNetwork(in.BlockchainB)
	require.NoError(t, err)
	// create network configs for 2 EVM networks
	srcNetworkCfg, err := clnode.NewNetworkCfg(&clnode.EVMNetworkConfig{
		MinIncomingConfirmations: 1,
		MinContractPayment:       "0.00001 link",
		ChainID:                  bcA.ChainID,
		EVMNodes: []*clnode.EVMNode{
			{
				SendOnly: false,
				Order:    100,
			},
		},
	}, bcA)
	dstNetworkConfig, err := clnode.NewNetworkCfg(&clnode.EVMNetworkConfig{
		MinIncomingConfirmations: 1,
		MinContractPayment:       "0.00001 link",
		ChainID:                  bcA.ChainID,
		EVMNodes: []*clnode.EVMNode{
			{
				SendOnly: false,
				Order:    100,
			},
		},
	}, bcB)
	// override the configuration to connect with 2 networks
	in.NodeSets[0].NodeSpecs[0].Node.TestConfigOverrides = srcNetworkCfg + dstNetworkConfig
	// create a node set
	nodesOut, err := ns.NewSharedDBNodeSet(in.NodeSets[0], bcA)
	require.NoError(t, err)

	c, err := clclient.New(nodesOut.CLNodes)
	require.NoError(t, err)

	testCases := []struct {
		name       string
		wait       time.Duration
		chainURL   string
		reorgDepth int
		validate   func(c []*clclient.ChainlinkClient) error
	}{
		{
			name:       "Reorg src with depth: 1",
			wait:       30 * time.Second,
			chainURL:   bcA.Nodes[0].ExternalHTTPUrl,
			reorgDepth: 1,
			validate: func(c []*clclient.ChainlinkClient) error {
				// add clients and validate
				return nil
			},
		},
		{
			name:       "Reorg dst with depth: 1",
			wait:       30 * time.Second,
			chainURL:   bcB.Nodes[0].ExternalHTTPUrl,
			reorgDepth: 1,
			validate: func(c []*clclient.ChainlinkClient) error {
				return nil
			},
		},
		{
			name:       "Reorg src with depth: 5",
			wait:       30 * time.Second,
			chainURL:   bcA.Nodes[0].ExternalHTTPUrl,
			reorgDepth: 5,
			validate: func(c []*clclient.ChainlinkClient) error {
				return nil
			},
		},
		{
			name:       "Reorg dst with depth: 5",
			wait:       30 * time.Second,
			chainURL:   bcB.Nodes[0].ExternalHTTPUrl,
			reorgDepth: 5,
			validate: func(c []*clclient.ChainlinkClient) error {
				return nil
			},
		},
	}

	// Start WASP load test here, apply average load profile that you expect in production!
	// Configure timeouts and validate all the test cases until the test ends

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Log(tc.name)
			r := rpc.New(tc.chainURL, nil)
			err := r.GethSetHead(tc.reorgDepth)
			require.NoError(t, err)
			time.Sleep(tc.wait)
			err = tc.validate(c)
			require.NoError(t, err)
		})
	}
}
