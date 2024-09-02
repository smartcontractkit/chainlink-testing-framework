package test_env

import (
	"fmt"
	"testing"

	config_types "github.com/smartcontractkit/chainlink-testing-framework/lib/config/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/config"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
)

func TestRethEth1(t *testing.T) {
	builder := NewEthereumNetworkBuilder()
	_, err := builder.
		WithEthereumVersion(config_types.EthereumVersion_Eth1).
		WithExecutionLayer(config_types.ExecutionLayer_Reth).
		Build()
	require.Error(t, err, "Builder validation failed")
	require.Contains(t, err.Error(), config.Eth1NotSupportedByRethMsg)
}

func TestRethEth2(t *testing.T) {
	l := logging.GetTestLogger(t)

	builder := NewEthereumNetworkBuilder()
	cfg, err := builder.
		WithEthereumVersion(config_types.EthereumVersion_Eth2).
		WithExecutionLayer(config_types.ExecutionLayer_Reth).
		Build()
	require.NoError(t, err, "Builder validation failed")

	_, eth2, err := cfg.Start()
	require.NoError(t, err, "Couldn't start PoS network")

	nonEip1559Network := blockchain.SimulatedEVMNetwork
	nonEip1559Network.Name = "Simulated Reth + Prysm (non-EIP 1559)"
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
	eip1559Network.Name = "Simulated Reth + Prysm (EIP 1559)"
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
