package test_env

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/docker"
)

func TestNonDevGeth(t *testing.T) {
	network, err := docker.CreateNetwork()
	require.NoError(t, err)
	g := NewPrivateGethChain(&blockchain.SimulatedEVMNetwork, []string{network.Name})
	err = g.GetPrimaryNode().Start()
	require.NoError(t, err)
	err = g.GetPrimaryNode().ConnectToClient()
	require.NoError(t, err)
}
