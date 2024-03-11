//go:build erigon_tests
// +build erigon_tests

package test_env

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/testcontext"
)

func TestEth1WithErigon(t *testing.T) {
	l := logging.GetTestLogger(t)

	builder := NewEthereumNetworkBuilder()
	cfg, err := builder.
		WithEthereumVersion(EthereumVersion_Eth1_Legacy).
		WithExecutionLayer(ExecutionLayer_Erigon).
		Build()
	require.NoError(t, err, "Builder validation failed")

	net, _, err := cfg.Start()
	require.NoError(t, err, "Couldn't start PoW network")

	c, err := blockchain.ConnectEVMClient(net, l)
	require.NoError(t, err, "Couldn't connect to the evm client")

	address := common.HexToAddress("0x90F8bf6A479f320ead074411a4B0e7944Ea8c9C1")
	err = sendAndCompareBalances(testcontext.Get(t), c, address)
	require.NoError(t, err, fmt.Sprintf("balance wasn't correctly updated for %s network", net.Name))

	err = c.Close()
	require.NoError(t, err, "Couldn't close the client")
}

func TestEth2WithPrysmAndErigon(t *testing.T) {
	l := logging.GetTestLogger(t)

	builder := NewEthereumNetworkBuilder()
	cfg, err := builder.
		WithEthereumVersion(EthereumVersion_Eth2_Legacy).
		WithConsensusLayer(ConsensusLayer_Prysm).
		WithExecutionLayer(ExecutionLayer_Erigon).
		Build()
	require.NoError(t, err, "Builder validation failed")

	_, eth2, err := cfg.Start()
	require.NoError(t, err, "Couldn't start PoS network")

	nonEip1559Network := blockchain.SimulatedEVMNetwork
	nonEip1559Network.Name = "Simulated Erigon + Prysm (non-EIP 1559)"
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
	eip1559Network.Name = "Simulated Erigon + Prysm (EIP 1559)"
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
