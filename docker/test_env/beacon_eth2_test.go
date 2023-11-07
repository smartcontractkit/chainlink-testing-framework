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

	gg := NewGethGenesis([]string{network.Name}, bg.ExecutionDir).WithTestLogger(t)
	err = gg.StartContainer()
	require.NoError(t, err)

	// geth2Name := fmt.Sprintf("%s-%s", "geth2", uuid.NewString()[0:8])

	// geth := NewGeth2([]string{network.Name}, bg.ExecutionDir, WithContainerName(geth2Name)).WithTestLogger(t)
	geth := NewGeth2([]string{network.Name}, bg.ExecutionDir).WithTestLogger(t)
	n, docker, err := geth.StartContainer()
	require.NoError(t, err)

	l.Error().Msgf("geth execution url: %s", geth.ExecutionURL)

	// beacon := NewBeaconChain([]string{network.Name}, bg.ExecutionDir, bg.ConsensusDir, fmt.Sprintf("http://%s:%s", geth2Name, GETH_EXECUTION_PORT)).WithTestLogger(t)
	beacon := NewBeaconChain([]string{network.Name}, bg.ExecutionDir, bg.ConsensusDir, geth.ExecutionURL).WithTestLogger(t)
	err = beacon.StartContainer()
	require.NoError(t, err)

	_ = n
	_ = docker

	// l.Error().Msgf("beacon rcp: %s", beacon.InternalRpcURL)

	validator := NewValidator([]string{network.Name}, bg.ConsensusDir, beacon.InternalRpcURL).WithTestLogger(t)
	err = validator.StartContainer()
	require.NoError(t, err)

	ns := blockchain.SimulatedEVMNetwork
	ns.URLs = []string{geth.ExternalWsUrl}
	c, err := blockchain.ConnectEVMClient(ns, l)
	require.NoError(t, err, "Couldn't connect to the evm client")
	err = c.Close()
	require.NoError(t, err, "Couldn't close the client")
}
