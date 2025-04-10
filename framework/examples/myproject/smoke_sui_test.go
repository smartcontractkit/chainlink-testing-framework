package examples

import (
	"context"
	"fmt"
	"testing"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/signer"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
)

type CfgSui struct {
	BlockchainA        *blockchain.Input `toml:"blockchain_a" validate:"required"`
	MockerDataProvider *fake.Input       `toml:"data_provider" validate:"required"`
	NodeSets           []*ns.Input       `toml:"nodesets" validate:"required"`
}

func TestSuiSmoke(t *testing.T) {
	in, err := framework.Load[CfgSui](t)
	require.NoError(t, err)

	bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
	require.NoError(t, err)

	// network is already funded, here are the keys
	_ = bc.NetworkSpecificData.SuiAccount.Mnemonic
	_ = bc.NetworkSpecificData.SuiAccount.PublicBase64Key
	_ = bc.NetworkSpecificData.SuiAccount.SuiAddress

	// execute any additional commands, to deploy contracts or set up
	dc, err := framework.NewDockerClient()
	require.NoError(t, err)
	_, err = dc.ExecContainer(bc.ContainerName, []string{"ls", "-lah"})
	require.NoError(t, err)

	_, err = fake.NewFakeDataProvider(in.MockerDataProvider)
	require.NoError(t, err)

	fmt.Printf("Sui host HTTP URL: %s", bc.Nodes[0].ExternalHTTPUrl)
	fmt.Printf("Sui internal (docker) HTTP URL: %s", bc.Nodes[0].InternalHTTPUrl)
	for _, n := range in.NodeSets[0].NodeSpecs {
		// configure each CL node for Sui, just an example
		n.Node.TestConfigOverrides = `
											[Log]
											level = 'info'
`
	}
	out, err := ns.NewSharedDBNodeSet(in.NodeSets[0], nil)
	require.NoError(t, err)

	t.Run("test something", func(t *testing.T) {
		cli := sui.NewSuiClient(bc.Nodes[0].ExternalHTTPUrl)

		signerAccount, err := signer.NewSignertWithMnemonic(bc.NetworkSpecificData.SuiAccount.Mnemonic)
		require.NoError(t, err)
		rsp, err := cli.SuiXGetAllBalance(context.Background(), models.SuiXGetAllBalanceRequest{
			Owner: signerAccount.Address,
		})
		require.NoError(t, err)
		fmt.Printf("My funds: %v\n", rsp)
		clClients, err := clclient.New(out.CLNodes)
		require.NoError(t, err)
		// create jobs, etc
		for _, c := range clClients {
			_ = c
			// create jobs
			//_, _, err := c.CreateJobRaw(`...`)
			//require.NoError(t, err)
		}
	})
}
