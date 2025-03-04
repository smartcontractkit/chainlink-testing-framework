package examples

import (
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/stretchr/testify/require"
)

type CfgAptos struct {
	BlockchainA *blockchain.Input `toml:"blockchain_a" validate:"required"`
}

func TestAptosSmoke(t *testing.T) {
	in, err := framework.Load[CfgAptos](t)
	require.NoError(t, err)

	bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
	require.NoError(t, err)

	// execute any additional commands, to deploy contracts or set up
	// network is already funded, here are the keys
	_ = blockchain.DefaultAptosAccount
	_ = blockchain.DefaultAptosPrivateKey

	_, err = framework.ExecContainer(bc.ContainerName, []string{"ls", "-lah"})
	require.NoError(t, err)

	t.Run("test something", func(t *testing.T) {
		// use internal URL to connect Chainlink nodes
		_ = bc.Nodes[0].DockerInternalHTTPUrl
		// use host URL to interact
		_ = bc.Nodes[0].HostHTTPUrl
		r := resty.New().SetBaseURL(bc.Nodes[0].HostHTTPUrl).EnableTrace()
		_, err := r.R().Get("/v1/transactions")
		require.NoError(t, err)
	})
}
