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

func TestErigonEth1(t *testing.T) {
	t.Skip("thorax/erigon is not available anymore, find the new image if needed and enable the test")
	l := logging.GetTestLogger(t)

	builder := NewEthereumNetworkBuilder()
	cfg, err := builder.
		//nolint:staticcheck //ignore SA1019
		WithEthereumVersion(config_types.EthereumVersion_Eth1_Legacy).
		WithExecutionLayer(config_types.ExecutionLayer_Erigon).
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

func TestErigonEth2(t *testing.T) {
	t.Skip("thorax/erigon is not available anymore, find the new image if needed and enable the test")
	l := logging.GetTestLogger(t)

	builder := NewEthereumNetworkBuilder()
	cfg, err := builder.
		WithEthereumVersion(config_types.EthereumVersion_Eth2).
		WithExecutionLayer(config_types.ExecutionLayer_Erigon).
		Build()
	require.NoError(t, err, "Builder validation failed")

	net, _, err := cfg.Start()
	require.NoError(t, err, "Couldn't start PoS network")

	c, err := blockchain.ConnectEVMClient(net, l)
	require.NoError(t, err, "Couldn't connect to the evm client")

	address := common.HexToAddress("0x90F8bf6A479f320ead074411a4B0e7944Ea8c9C1")
	err = sendAndCompareBalances(testcontext.Get(t), c, address)
	require.NoError(t, err, fmt.Sprintf("balance wasn't correctly updated for %s network", net.Name))

	err = c.Close()
	require.NoError(t, err, "Couldn't close the client")
}

func TestErigonEth2_Deneb(t *testing.T) {
	t.Skip("thorax/erigon is not available anymore, find the new image if needed and enable the test")
	l := logging.GetTestLogger(t)

	builder := NewEthereumNetworkBuilder()
	cfg, err := builder.
		//nolint:staticcheck //ignore SA1019
		WithCustomDockerImages(map[config.ContainerType]string{config.ContainerType_ExecutionLayer: "thorax/erigon:v2.59.0"}).
		WithConsensusLayer(config.ConsensusLayer_Prysm).
		WithExecutionLayer(config_types.ExecutionLayer_Erigon).
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

func TestErigonEth2_Shanghai(t *testing.T) {
	t.Skip("thorax/erigon is not available anymore, find the new image if needed and enable the test")
	l := logging.GetTestLogger(t)

	chainConfig := config.MustGetDefaultChainConfig()
	chainConfig.HardForkEpochs = map[string]int{"Deneb": 500}

	builder := NewEthereumNetworkBuilder()
	cfg, err := builder.
		WithCustomDockerImages(map[config.ContainerType]string{config.ContainerType_ExecutionLayer: "thorax/erigon:v2.58.0"}).
		WithExecutionLayer(config_types.ExecutionLayer_Erigon).
		WithEthereumChainConfig(chainConfig).
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
}
