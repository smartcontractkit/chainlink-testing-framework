package test_env

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/docker"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
)

func TestGeth(t *testing.T) {
	l := logging.GetTestLogger(t)
	network, err := docker.CreateNetwork(l)
	require.NoError(t, err)
	g := NewGeth([]string{network.Name}).
		WithTestLogger(t)
	_, _, err = g.StartContainer()
	require.NoError(t, err)
	ns := blockchain.SimulatedEVMNetwork
	ns.URLs = []string{g.ExternalWsUrl}
	c, err := blockchain.ConnectEVMClient(ns, l)
	require.NoError(t, err, "Couldn't connect to the evm client")
	err = c.Close()
	require.NoError(t, err, "Couldn't close the client")
}
