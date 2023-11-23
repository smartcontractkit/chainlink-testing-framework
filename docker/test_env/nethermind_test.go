package test_env

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
)

func TestEth2WithPrysmAndNethermind(t *testing.T) {
	l := logging.GetTestLogger(t)

	builder := NewEthereumNetworkBuilder()
	cfg, err := builder.
		WithConsensusType(ConsensusType_PoS).
		WithCustomNetworkParticipants([]EthereumNetworkParticipant{
			{
				ConsensusLayer: ConsensusLayer_Prysm,
				ExecutionLayer: ExecutionLayer_Nethermind,
				Count:          1,
			},
		}).
		Build()
	require.NoError(t, err, "Builder validation failed")

	_, eth2, err := cfg.Start()
	require.NoError(t, err, "Couldn't start PoS network")

	ns := blockchain.SimulatedEVMNetwork
	ns.Name = "Simulated Nethermind + Prysm"
	ns.URLs = eth2.PublicWsUrls()
	c, err := blockchain.ConnectEVMClient(ns, l)
	require.NoError(t, err, "Couldn't connect to the evm client")
	err = c.Close()
	require.NoError(t, err, "Couldn't close the client")
}
