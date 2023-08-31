package test_env

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/docker"
)

func TestGeth(t *testing.T) {
	network, err := docker.CreateNetwork()
	require.NoError(t, err)
	g := NewGeth([]string{network.Name})
	_, _, err = g.StartContainer()
	require.NoError(t, err)
	ns := blockchain.SimulatedEVMNetwork
	ns.URLs = []string{g.ExternalWsUrl}
	_, err = blockchain.ConnectEVMClient(ns)
	require.NoError(t, err, "Couldn't connect to the evm client")
}
