package test_env

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/docker"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
)

func TestEth2(t *testing.T) {
	l := logging.GetTestLogger(t)
	network, err := docker.CreateNetwork(l)
	require.NoError(t, err)
	bg := NewBeaconChainGenesis([]string{network.Name}).
		WithTestLogger(t)
	err = bg.StartContainer()
	require.NoError(t, err)

	n, cmp, err := StartEth2(t, ConsensusLayer_Prysm)
	require.NoError(t, err, "Couldn't start eth2")

	_ = n

	ns := blockchain.SimulatedEVMNetwork
	ns.URLs = []string{cmp.Geth.ExternalWsUrl}
	c, err := blockchain.ConnectEVMClient(ns, l)
	require.NoError(t, err, "Couldn't connect to the evm client")
	err = c.Close()
	require.NoError(t, err, "Couldn't close the client")
}
