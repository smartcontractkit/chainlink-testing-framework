//go:build geth_tests
// +build geth_tests

package test_env

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/docker"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/testcontext"
)

func TestOldGeth(t *testing.T) {
	l := logging.GetTestLogger(t)
	network, err := docker.CreateNetwork(l)
	require.NoError(t, err)
	defaultChainCfg := GetDefaultChainConfig()
	g := NewGethEth1([]string{network.Name}, &defaultChainCfg).
		WithTestInstance(t)
	_, err = g.StartContainer()
	require.NoError(t, err)
	ns := blockchain.SimulatedEVMNetwork
	ns.URLs = []string{g.GetExternalWsUrl()}
	c, err := blockchain.ConnectEVMClient(ns, l)
	require.NoError(t, err, "Couldn't connect to the evm client")
	err = c.Close()
	require.NoError(t, err, "Couldn't close the client")
}

func TestEth1WithGeth(t *testing.T) {
	l := logging.GetTestLogger(t)

	builder := NewEthereumNetworkBuilder()
	cfg, err := builder.
		WithEthereumVersion(EthereumVersion_Eth1_Legacy).
		WithExecutionLayer(ExecutionLayer_Geth).
		Build()
	require.NoError(t, err, "Builder validation failed")

	_, eth2, err := cfg.Start()
	require.NoError(t, err, "Couldn't start PoW network")

	ns := blockchain.SimulatedEVMNetwork
	ns.URLs = eth2.PublicWsUrls()
	c, err := blockchain.ConnectEVMClient(ns, l)
	require.NoError(t, err, "Couldn't connect to the evm client")
	err = c.Close()
	require.NoError(t, err, "Couldn't close the client")
}

func TestEth2WithPrysmAndGeth(t *testing.T) {
	l := logging.GetTestLogger(t)

	builder := NewEthereumNetworkBuilder()
	cfg, err := builder.
		WithEthereumVersion(EthereumVersion_Eth2_Legacy).
		WithConsensusLayer(ConsensusLayer_Prysm).
		WithExecutionLayer(ExecutionLayer_Geth).
		Build()
	require.NoError(t, err, "Builder validation failed")

	_, eth2, err := cfg.Start()
	require.NoError(t, err, "Couldn't start PoS network")

	nonEip1559Network := blockchain.SimulatedEVMNetwork
	nonEip1559Network.Name = "Simulated Geth + Prysm (non-EIP 1559)"
	nonEip1559Network.URLs = eth2.PublicWsUrls()
	clientOne, err := blockchain.ConnectEVMClient(nonEip1559Network, l)
	require.NoError(t, err, "Couldn't connect to the evm client")

	t.Cleanup(func() {
		err = clientOne.Close()
		require.NoError(t, err, "Couldn't close the client")
	})

	ctx := testcontext.Get(t)
	address := common.HexToAddress("0x90F8bf6A479f320ead074411a4B0e7944Ea8c9C1")
	err = sendAndCompareBalances(ctx, clientOne, address)
	require.NoError(t, err, fmt.Sprintf("balance wasn't correctly updated for %s network", nonEip1559Network.Name))

	eip1559Network := blockchain.SimulatedEVMNetwork
	eip1559Network.Name = "Simulated Geth + Prysm (EIP 1559)"
	eip1559Network.SupportsEIP1559 = true
	eip1559Network.URLs = eth2.PublicWsUrls()
	clientTwo, err := blockchain.ConnectEVMClient(eip1559Network, l)
	require.NoError(t, err, "Couldn't connect to the evm client")

	t.Cleanup(func() {
		err = clientTwo.Close()
		require.NoError(t, err, "Couldn't close the client")
	})

	err = sendAndCompareBalances(ctx, clientTwo, address)
	require.NoError(t, err, fmt.Sprintf("balance wasn't correctly updated for %s network", eip1559Network.Name))
}
