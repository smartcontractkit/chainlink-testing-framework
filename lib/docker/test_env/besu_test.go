package test_env

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/config"
	config_types "github.com/smartcontractkit/chainlink-testing-framework/lib/config/types"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
)

func TestBesuEth1(t *testing.T) {
	l := logging.GetTestLogger(t)

	builder := NewEthereumNetworkBuilder()
	cfg, err := builder.
		//nolint:staticcheck //ignore SA1019
		WithEthereumVersion(config_types.EthereumVersion_Eth1_Legacy).
		WithExecutionLayer(config_types.ExecutionLayer_Besu).
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

func TestBesuEth2(t *testing.T) {
	l := logging.GetTestLogger(t)

	builder := NewEthereumNetworkBuilder()
	cfg, err := builder.
		WithEthereumVersion(config_types.EthereumVersion_Eth2).
		WithExecutionLayer(config_types.ExecutionLayer_Besu).
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

func TestBesuEth2_Deneb(t *testing.T) {
	l := logging.GetTestLogger(t)

	builder := NewEthereumNetworkBuilder()
	cfg, err := builder.
		WithCustomDockerImages(map[config.ContainerType]string{config.ContainerType_ExecutionLayer: "hyperledger/besu:24.1.0"}).
		WithConsensusLayer(config.ConsensusLayer_Prysm).
		WithExecutionLayer(config_types.ExecutionLayer_Besu).
		Build()
	require.NoError(t, err, "Builder validation failed")

	_, eth2, err := cfg.Start()
	require.NoError(t, err, "Couldn't start PoS network")

	nonEip1559Network := blockchain.SimulatedEVMNetwork
	nonEip1559Network.Name = "Simulated Besu + Prysm (non-EIP 1559)"
	nonEip1559Network.GasEstimationBuffer = 10_000_000_000
	nonEip1559Network.URLs = eth2.PublicWsUrls()
	clientOne, err := blockchain.ConnectEVMClient(nonEip1559Network, l)
	require.NoError(t, err, "Couldn't connect to the evm client")

	t.Cleanup(func() {
		err = clientOne.Close()
		require.NoError(t, err, "Couldn't close the client")
	})

	address := common.HexToAddress("0x90F8bf6A479f320ead074411a4B0e7944Ea8c9C1")
	err = sendAndCompareBalances(testcontext.Get(t), clientOne, address)
	require.NoError(t, err, fmt.Sprintf("balance wasn't correctly updated for %s network", nonEip1559Network.Name))

	eip1559Network := blockchain.SimulatedEVMNetwork
	eip1559Network.Name = "Simulated Besu + Prysm (EIP 1559)"
	eip1559Network.SupportsEIP1559 = true
	eip1559Network.URLs = eth2.PublicWsUrls()
	_, err = blockchain.ConnectEVMClient(eip1559Network, l)
	require.Error(t, err, "Could not connect to Besu")
	require.Contains(t, err.Error(), "Method not found", "Besu should not work EIP-1559 yet")
}

func TestBesuEth2_Shanghai(t *testing.T) {
	l := logging.GetTestLogger(t)

	chainConfig := config.MustGetDefaultChainConfig()
	chainConfig.HardForkEpochs = map[string]int{"Deneb": 500}

	builder := NewEthereumNetworkBuilder()
	cfg, err := builder.
		WithCustomDockerImages(map[config.ContainerType]string{config.ContainerType_ExecutionLayer: "hyperledger/besu:23.10"}).
		WithExecutionLayer(config_types.ExecutionLayer_Besu).
		WithEthereumChainConfig(chainConfig).
		Build()
	require.NoError(t, err, "Builder validation failed")

	_, eth2, err := cfg.Start()
	require.NoError(t, err, "Couldn't start PoS network")

	nonEip1559Network := blockchain.SimulatedEVMNetwork
	nonEip1559Network.Name = "Simulated Besu + Prysm (non-EIP 1559)"
	nonEip1559Network.GasEstimationBuffer = 10_000_000_000
	nonEip1559Network.URLs = eth2.PublicWsUrls()
	clientOne, err := blockchain.ConnectEVMClient(nonEip1559Network, l)
	require.NoError(t, err, "Couldn't connect to the evm client")

	t.Cleanup(func() {
		err = clientOne.Close()
		require.NoError(t, err, "Couldn't close the client")
	})

	address := common.HexToAddress("0x90F8bf6A479f320ead074411a4B0e7944Ea8c9C1")
	err = sendAndCompareBalances(testcontext.Get(t), clientOne, address)
	require.NoError(t, err, fmt.Sprintf("balance wasn't correctly updated for %s network", nonEip1559Network.Name))
}
