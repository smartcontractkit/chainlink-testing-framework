package test_env

import (
	"context"
	"fmt"
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
		WithConsensusLayer(ConsensusLayer_Prysm).
		WithExecutionLayer(ExecutionLayer_Geth).
		WithEthereumChainConfig(EthereumChainConfig{
			SecondsPerSlot: 4,
			SlotsPerEpoch:  2,
		}).
		Build()
	require.NoError(t, err, "Builder validation failed")

	_, eth2, err := cfg.Start()
	require.NoError(t, err, "Couldn't start PoS network")

	ns := blockchain.SimulatedEVMNetwork
	ns.Name = "Simulated Geth + Prysm"
	ns.URLs = eth2.PublicWsUrls()
	c, err := blockchain.ConnectEVMClient(ns, l)
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
		WithConsensusLayer(ConsensusLayer_Prysm).
		WithExecutionLayer(ExecutionLayer_Geth).
		WithAddressesToFund([]string{addressToFund}).
		Build()
	require.NoError(t, err, "Builder validation failed")

	_, eth2, err := cfg.Start()
	require.NoError(t, err, "Couldn't start PoS network")

	ns := blockchain.SimulatedEVMNetwork
	ns.Name = "Simulated Geth + Prysm"
	ns.URLs = eth2.PublicWsUrls()
	c, err := blockchain.ConnectEVMClient(ns, l)
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
		WithConsensusLayer(ConsensusLayer_Prysm).
		WithExecutionLayer(ExecutionLayer_Geth).
		Build()
	require.NoError(t, err, "Builder validation failed")

	_, _, err = cfg.Start()
	require.NoError(t, err, "Couldn't start PoS network")

	newBuilder := NewEthereumNetworkBuilder()
	reusedCfg, err := newBuilder.
		WithExistingConfig(cfg).
		Build()
	require.NoError(t, err, "Builder validation failed")

	_, eth2, err := reusedCfg.Start()
	require.NoError(t, err, "Couldn't reuse PoS network")

	ns := blockchain.SimulatedEVMNetwork
	ns.Name = "Simulated Geth + Prysm"
	ns.URLs = eth2.PublicWsUrls()
	c, err := blockchain.ConnectEVMClient(ns, l)
	require.NoError(t, err, "Couldn't connect to the evm client")
	err = c.Close()
	require.NoError(t, err, "Couldn't close the client")
}
