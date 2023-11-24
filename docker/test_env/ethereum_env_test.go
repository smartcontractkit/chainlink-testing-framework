package test_env

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
)

func TestEth2CustomConfig(t *testing.T) {
	l := logging.GetTestLogger(t)

	builder := NewEthereumNetworkBuilder()
	cfg, err := builder.
		WithConsensusType(ConsensusType_PoS).
		WithCustomNetworkParticipants([]EthereumNetworkParticipant{
			{
				ConsensusLayer: ConsensusLayer_Prysm,
				ExecutionLayer: ExecutionLayer_Geth,
				Count:          1,
			},
		}).
		WithEthereumChainConfig(EthereumChainConfig{
			SecondsPerSlot: 6,
			SlotsPerEpoch:  2,
		}).
		WithoutWaitingForFinalization().
		Build()
	require.NoError(t, err, "Builder validation failed")

	net, _, err := cfg.Start()
	require.NoError(t, err, "Couldn't start PoS network")

	c, err := blockchain.ConnectEVMClient(net, l)
	require.NoError(t, err, "Couldn't connect to the evm client")
	err = c.Close()
	require.NoError(t, err, "Couldn't close the client")
}

func TestEth2ExtraFunding(t *testing.T) {
	l := logging.GetTestLogger(t)

	addressToFund := "0x14dc79964da2c08b23698b3d3cc7ca32193d9955"

	builder := NewEthereumNetworkBuilder()
	cfg, err := builder.
		WithConsensusType(ConsensusType_PoS).
		WithCustomNetworkParticipants([]EthereumNetworkParticipant{
			{
				ConsensusLayer: ConsensusLayer_Prysm,
				ExecutionLayer: ExecutionLayer_Geth,
				Count:          1,
			},
		}).
		WithEthereumChainConfig(EthereumChainConfig{
			AddressesToFund: []string{addressToFund},
		}).
		WithoutWaitingForFinalization().
		Build()
	require.NoError(t, err, "Builder validation failed")

	net, _, err := cfg.Start()
	require.NoError(t, err, "Couldn't start PoS network")

	c, err := blockchain.ConnectEVMClient(net, l)
	require.NoError(t, err, "Couldn't connect to the evm client")

	balance, err := c.BalanceAt(context.Background(), common.HexToAddress(addressToFund))
	require.NoError(t, err, "Couldn't get balance")
	require.Equal(t, "11515845246265065472", fmt.Sprintf("%d", balance.Uint64()), "Balance is not correct")

	err = c.Close()
	require.NoError(t, err, "Couldn't close the client")
}

func TestEth2WithPrysmAndGethReuseNetwork(t *testing.T) {
	l := logging.GetTestLogger(t)

	builder := NewEthereumNetworkBuilder()
	cfg, err := builder.
		WithConsensusType(ConsensusType_PoS).
		WithCustomNetworkParticipants([]EthereumNetworkParticipant{
			{
				ConsensusLayer: ConsensusLayer_Prysm,
				ExecutionLayer: ExecutionLayer_Geth,
				Count:          1,
			},
		}).
		WithoutWaitingForFinalization().
		Build()
	require.NoError(t, err, "Builder validation failed")

	_, _, err = cfg.Start()
	require.NoError(t, err, "Couldn't start PoS network")

	newBuilder := NewEthereumNetworkBuilder()
	reusedCfg, err := newBuilder.
		WithExistingConfig(cfg).
		Build()
	require.NoError(t, err, "Builder validation failed")

	net, _, err := reusedCfg.Start()
	require.NoError(t, err, "Couldn't reuse PoS network")

	c, err := blockchain.ConnectEVMClient(net, l)
	require.NoError(t, err, "Couldn't connect to the evm client")
	err = c.Close()
	require.NoError(t, err, "Couldn't close the client")
}

func TestEth2WithPrysmAndGethReuseFromEnv(t *testing.T) {
	t.Skip("for demonstration purposes only")
	l := logging.GetTestLogger(t)

	os.Setenv(CONFIG_ENV_VAR_NAME, "CHANGE_ME")
	builder := NewEthereumNetworkBuilder()
	cfg, err := builder.
		FromEnvVar().
		Build()
	require.NoError(t, err, "Builder validation failed")

	net, _, err := cfg.Start()
	require.NoError(t, err, "Couldn't start PoS network")

	c, err := blockchain.ConnectEVMClient(net, l)
	require.NoError(t, err, "Couldn't connect to the evm client")
	err = c.Close()
	require.NoError(t, err, "Couldn't close the client")
}
