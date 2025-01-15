package examples

import (
	"context"
	"fmt"
	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/signer"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/stretchr/testify/require"
	"testing"
)

type CfgSui struct {
	BlockchainA *blockchain.Input `toml:"blockchain_a" validate:"required"`
}

func TestSuiSmoke(t *testing.T) {
	in, err := framework.Load[CfgSui](t)
	require.NoError(t, err)

	bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
	require.NoError(t, err)

	// execute any additional commands, to deploy contracts or set up
	// network is already funded, here are the keys
	_ = bc.GeneratedData.Mnemonic

	_, err = framework.ExecContainer(bc.ContainerName, []string{"ls", "-lah"})
	require.NoError(t, err)

	t.Run("test something", func(t *testing.T) {
		// use internal URL to connect Chainlink nodes
		_ = bc.Nodes[0].DockerInternalHTTPUrl
		// use host URL to interact
		_ = bc.Nodes[0].HostHTTPUrl
		cli := sui.NewSuiClient("http://localhost:9000")

		signerAccount, err := signer.NewSignertWithMnemonic(bc.GeneratedData.Mnemonic)
		require.NoError(t, err)
		rsp, err := cli.SuiXGetAllBalance(context.Background(), models.SuiXGetAllBalanceRequest{
			Owner: signerAccount.Address,
		})
		require.NoError(t, err)
		fmt.Printf("My funds: %v\n", rsp)
	})
}
