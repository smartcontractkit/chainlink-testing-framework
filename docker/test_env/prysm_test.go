package test_env

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
)

func TestEth2WithPrysmAndGethDefaultConfig(t *testing.T) {
	l := logging.GetTestLogger(t)

	builder := NewEthereumNetworkBuilder(t)
	err := builder.
		WithConsensusType(ConsensusType_PoS).
		WithConsensusLayer(ConsensusLayer_Prysm).
		WithExecutionLayer(ExecutionLayer_Geth).
		Build()
	require.NoError(t, err, "Builder validation failed")

	_, eth2, _, err := builder.Start()
	require.NoError(t, err, "Couldn't start PoS network")

	ns := blockchain.SimulatedEVMNetwork
	ns.URLs = eth2.PublicWsUrsl()
	c, err := blockchain.ConnectEVMClient(ns, l)
	require.NoError(t, err, "Couldn't connect to the evm client")
	err = c.Close()
	require.NoError(t, err, "Couldn't close the client")
}

func TestEth2WithPrysmAndGethCustomConfig(t *testing.T) {
	l := logging.GetTestLogger(t)

	builder := NewEthereumNetworkBuilder(t)
	err := builder.
		WithConsensusType(ConsensusType_PoS).
		WithConsensusLayer(ConsensusLayer_Prysm).
		WithExecutionLayer(ExecutionLayer_Geth).
		WithBeaconChainConfig(BeaconChainConfig{
			SecondsPerSlot: 4,
			SlotsPerEpoch:  2,
		}).
		Build()
	require.NoError(t, err, "Builder validation failed")

	_, eth2, _, err := builder.Start()
	require.NoError(t, err, "Couldn't start PoS network")

	ns := blockchain.SimulatedEVMNetwork
	ns.URLs = eth2.PublicWsUrsl()
	c, err := blockchain.ConnectEVMClient(ns, l)
	require.NoError(t, err, "Couldn't connect to the evm client")
	err = c.Close()
	require.NoError(t, err, "Couldn't close the client")
}

func TestEth2WithPrysmRestart(t *testing.T) {
	t.Skip("add support for restarting -- meaning that we shouldn't create any config files, keystores, etc anymore")
	l := logging.GetTestLogger(t)

	builder := NewEthereumNetworkBuilder(t)
	err := builder.
		WithConsensusType(ConsensusType_PoS).
		WithConsensusLayer(ConsensusLayer_Prysm).
		WithExecutionLayer(ExecutionLayer_Geth).
		Build()
	require.NoError(t, err, "Builder validation failed")

	_, eth2, cfg, err := builder.Start()
	require.NoError(t, err, "Couldn't start PoS network")

	ns := blockchain.SimulatedEVMNetwork
	ns.URLs = eth2.PublicWsUrsl()
	c, err := blockchain.ConnectEVMClient(ns, l)
	require.NoError(t, err, "Couldn't connect to the evm client")
	err = c.Close()
	require.NoError(t, err, "Couldn't close the client")

	ctx := context.Background()
	stopDuration := 5 * time.Second
	for _, c := range cfg.Containers {
		err := (*c.Container).Stop(ctx, &stopDuration)
		require.NoError(t, err, fmt.Sprintf("Couldn't stop %s container", c.ContainerName))
	}

	builder = NewEthereumNetworkBuilder(t)
	err = builder.
		WithExistingConfig(cfg).
		Build()
	require.NoError(t, err, "Builder validation failed")

	_, eth2, _, err = builder.Start()
	require.NoError(t, err, "Couldn't start PoS network")

	c, err = blockchain.ConnectEVMClient(ns, l)
	require.NoError(t, err, "Couldn't connect to the evm client")
	err = c.Close()
	require.NoError(t, err, "Couldn't close the client")
}
