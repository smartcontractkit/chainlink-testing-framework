package exampleProduct

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/clnode"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/dp"
	"github.com/stretchr/testify/require"
	"testing"
)

type Config struct {
	BlockchainA  *blockchain.Input `toml:"blockchain_a" validate:"required"`
	BlockchainB  *blockchain.Input `toml:"blockchain_b" validate:"required"`
	DataProvider *dp.Input         `toml:"data_provider" validate:"required"`
	CLNodeOne    *clnode.Input     `toml:"clnode_1" validate:"required"`
	CLNodeTwo    *clnode.Input     `toml:"clnode_2" validate:"required"`
}

func TestMultiNodeMultiNetwork(t *testing.T) {
	in, err := framework.Load[Config](t)
	require.NoError(t, err)

	bcNodes1, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
	require.NoError(t, err)
	bcNodes2, err := blockchain.NewBlockchainNetwork(in.BlockchainB)
	require.NoError(t, err)

	dpout, err := dp.NewMockedDataProvider(in.DataProvider)
	require.NoError(t, err)
	in.CLNodeOne.DataProviderURL = dpout.Urls[0]
	in.CLNodeTwo.DataProviderURL = dpout.Urls[0]

	networkA, err := clnode.NewNetworkCfg(&clnode.NetworkConfig{
		MinIncomingConfirmations: 1,
		MinContractPayment:       "0.0001 link",
		EVMNodes: []*clnode.EVMNode{
			{
				SendOnly: false,
				Order:    100,
			},
		},
	}, bcNodes1)
	require.NoError(t, err)
	networkB, err := clnode.NewNetworkCfg(&clnode.NetworkConfig{
		MinIncomingConfirmations: 1,
		MinContractPayment:       "0.0001 link",
		EVMNodes: []*clnode.EVMNode{
			{
				SendOnly: false,
				Order:    100,
			},
		},
	}, bcNodes2)
	require.NoError(t, err)

	in.CLNodeOne.Node.TestConfigOverrides = networkA + networkB
	in.CLNodeTwo.Node.TestConfigOverrides = networkB

	_, err = clnode.NewNode(in.CLNodeOne)
	require.NoError(t, err)

	_, err = clnode.NewNode(in.CLNodeTwo)
	require.NoError(t, err)

	// connect
	// seth, ctfCLClient, wasp, havoc

	// test / assert

	t.Run("test feature A1", func(t *testing.T) {
		client := resty.New()
		_, err := client.R().
			Get("http://localhost:8080/mock1")
		require.NoError(t, err)
	})
	t.Run("test feature A2", func(t *testing.T) {
		fmt.Println("Complex testing in progress...")
		fmt.Println("Complex testing in progress...")
		fmt.Println("Complex testing in progress... Done!")
	})
}
