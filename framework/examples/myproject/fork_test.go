package examples

import (
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/rpc"
	"github.com/stretchr/testify/require"
	"testing"
)

type CfgFork struct {
	BlockchainA *blockchain.Input `toml:"blockchain_a" validate:"required"`
}

func TestFork(t *testing.T) {
	in, err := framework.Load[CfgFork](t)
	require.NoError(t, err)

	bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
	require.NoError(t, err)

	t.Run("test some contracts with fork", func(t *testing.T) {
		ac := rpc.New(bc.Nodes[0].HostHTTPUrl, nil)
		err = ac.AnvilAutoImpersonate(true)
		require.NoError(t, err)
	})
}
