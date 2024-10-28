package test_env

import (
	"fmt"
	"testing"

	config_types "github.com/smartcontractkit/chainlink-testing-framework/lib/config/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/config"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/docker"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
)

func TestGethLegacy(t *testing.T) {
	l := logging.GetTestLogger(t)
	network, err := docker.CreateNetwork(l)
	require.NoError(t, err)
	defaultChainCfg := config.MustGetDefaultChainConfig()
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

func TestGethEth1(t *testing.T) {
	l := logging.GetTestLogger(t)

	builder := NewEthereumNetworkBuilder()
	cfg, err := builder.
		//nolint:staticcheck //ignore SA1019
		WithEthereumVersion(config_types.EthereumVersion_Eth1_Legacy).
		WithEthereumChainConfig(config.EthereumChainConfig{
			ChainID: 2337,
		}).
		WithExecutionLayer(config_types.ExecutionLayer_Geth).
		Build()
	require.NoError(t, err, "Builder validation failed")

	net, _, err := cfg.Start()
	require.NoError(t, err, "Couldn't start PoW network")

	c, err := blockchain.ConnectEVMClient(net, l)
	require.NoError(t, err, "Couldn't connect to the evm client")
	err = c.Close()
	require.NoError(t, err, "Couldn't close the client")
}

func TestGethEth2(t *testing.T) {
	l := logging.GetTestLogger(t)

	builder := NewEthereumNetworkBuilder()
	cfg, err := builder.
		WithEthereumVersion(config_types.EthereumVersion_Eth2).
		WithExecutionLayer(config_types.ExecutionLayer_Geth).
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
}

func TestGethEth2_Deneb(t *testing.T) {
	l := logging.GetTestLogger(t)

	builder := NewEthereumNetworkBuilder()
	cfg, err := builder.
		WithCustomDockerImages(map[config.ContainerType]string{config.ContainerType_ExecutionLayer: "ethereum/client-go:v1.13.12"}).
		WithConsensusLayer(config.ConsensusLayer_Prysm).
		WithExecutionLayer(config_types.ExecutionLayer_Geth).
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

func TestGethEth2_Shanghai(t *testing.T) {
	l := logging.GetTestLogger(t)

	chainConfig := config.MustGetDefaultChainConfig()
	chainConfig.HardForkEpochs = map[string]int{"Deneb": 500}

	builder := NewEthereumNetworkBuilder()
	cfg, err := builder.
		WithCustomDockerImages(map[config.ContainerType]string{config.ContainerType_ExecutionLayer: "ethereum/client-go:v1.13.11"}).
		WithExecutionLayer(config_types.ExecutionLayer_Geth).
		WithEthereumChainConfig(chainConfig).
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
}
