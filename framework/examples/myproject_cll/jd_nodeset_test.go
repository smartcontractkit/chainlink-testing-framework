package examples

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/rpc"
)

type CfgJDNodeSet struct {
	Blockchains        []*blockchain.Input `toml:"blockchains" validate:"required"`
	NodeSets           []*ns.Input         `toml:"nodesets" validate:"required"`
	MockedDataProvider *fake.Input         `toml:"data_provider" validate:"required"`
	JD                 *jd.Input           `toml:"jd" validate:"required"`
}

func TestJDNodeSet(t *testing.T) {
	in, err := framework.Load[CfgJDNodeSet](t)
	require.NoError(t, err)

	blockchainInfo, err := blockchain.NewBlockchainNetwork(in.Blockchains[0])
	require.NoError(t, err)

	// by default, nodeset is connected to the first Anvil but if you need more chains you can do it like that
	// 	srcNetworkCfg, err := clnode.NewNetworkCfg(&clnode.EVMNetworkConfig{
	//		MinIncomingConfirmations: 1,
	//		MinContractPayment:       "0.00001 link",
	//		ChainID:                  bcSrc.ChainID,
	//		EVMNodes: []*clnode.EVMNode{
	//			{
	//				SendOnly: false,
	//				Order:    100,
	//			},
	//		},
	//	}, bcSrc)
	//	in.NodeSets[0].NodeSpecs[0].Node.TestConfigOverrides = srcNetworkCfg

	jobDistributorInfo, err := jd.NewJD(in.JD)
	require.NoError(t, err)

	// connect JD with NodeSet

	nodeSetInfo, err := ns.NewSharedDBNodeSet(in.NodeSets[0], blockchainInfo)
	require.NoError(t, err)
	_, err = fake.NewFakeDataProvider(in.MockedDataProvider)
	require.NoError(t, err)

	// set up your contracts here with 0s blocks, control the mining speed later

	miner := rpc.NewRemoteAnvilMiner(blockchainInfo.Nodes[0].ExternalHTTPUrl, nil)
	miner.MinePeriodically(1 * time.Second)
	clClients, err := clclient.New(nodeSetInfo.CLNodes)
	require.NoError(t, err)

	t.Run("test #1", func(t *testing.T) {
		_ = clClients
		// create some jobs
		//_, _, err = c[0].CreateJobRaw()
		//require.NoError(t, err)
		_ = nodeSetInfo
		_ = blockchainInfo
		_ = jobDistributorInfo
	})

	t.Run("test #2", func(t *testing.T) {
		_ = clClients
		// create some jobs
		//_, _, err = c[0].CreateJobRaw()
		//require.NoError(t, err)
		_ = nodeSetInfo
		_ = blockchainInfo
		_ = jobDistributorInfo
	})
}
