package test_env

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/docker"
)

func TestNonDevGeth(t *testing.T) {
	//_ = os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")
	network, err := docker.CreateNetwork()
	require.NoError(t, err)
	g := NewPrivateGethChain(&blockchain.SimulatedEVMNetwork, []string{network.Name})
	err = g.PrimaryNode.Start()
	require.NoError(t, err)
	err = g.PrimaryNode.ConnectToClient()
	require.NoError(t, err)
}
