package test_env

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/docker"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
)

func TestNonDevGeth(t *testing.T) {
	l := logging.GetTestLogger(t)
	network, err := docker.CreateNetwork(l)
	require.NoError(t, err)
	g := NewPrivateGethChain(&blockchain.SimulatedEVMNetwork, []string{network.Name}).
		WithTestLogger(t)
	err = g.PrimaryNode.Start()
	require.NoError(t, err)
	err = g.PrimaryNode.ConnectToClient()
	require.NoError(t, err)
}
