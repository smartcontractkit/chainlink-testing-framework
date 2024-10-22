package test_env

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/config/types"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/docker/ethereum"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/config"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
)

var twentySecDuration = time.Duration(20 * time.Second)
var stopAllContainers = func(t *testing.T, containers []config.EthereumNetworkContainer) {
	for _, c := range containers {
		err := (*c.Container).Stop(context.Background(), &twentySecDuration)
		require.NoError(t, err, "Couldn't stop container")
	}

}

// 1.16 -> works
// ...
// 1.25.4 -> works
func TestEthEnvNethermindCompatibility(t *testing.T) {
	t.Skip("Execute manually when needed")
	l := logging.GetTestLogger(t)

	nethermindTcs := []string{}
	for i := 16; i < 26; i++ {
		nethermindTcs = append(nethermindTcs, fmt.Sprintf("nethermind/nethermind:1.%d.0", i))
	}

	nethermindTcs = append(nethermindTcs, ethereum.DefaultNethermindEth2Image)
	latest, err := FetchLatestEthereumClientDockerImageVersionIfNeed(fmt.Sprintf("nethermind/nethermind:%s", AUTOMATIC_STABLE_LATEST_TAG))
	require.NoError(t, err, "Couldn't fetch the latest Nethermind version")

	nethermindTcs = append(nethermindTcs, latest)
	nethermindTcs = UniqueStringSlice(nethermindTcs)

	for _, tc := range nethermindTcs {
		t.Run(fmt.Sprintf("nethermind-%s", tc), func(t *testing.T) {
			builder := NewEthereumNetworkBuilder()
			cfg, err := builder.
				WithExecutionLayer(types.ExecutionLayer_Nethermind).
				WithCustomDockerImages(map[config.ContainerType]string{
					config.ContainerType_ExecutionLayer: tc,
				}).
				Build()
			require.NoError(t, err, "Builder validation failed")

			net, _, err := cfg.Start()
			require.NoError(t, err, "Couldn't start Nethermind-based network")

			c, err := blockchain.ConnectEVMClient(net, l)
			require.NoError(t, err, "Couldn't connect to the evm client")

			address := common.HexToAddress("0x90F8bf6A479f320ead074411a4B0e7944Ea8c9C1")
			err = sendAndCompareBalances(testcontext.Get(t), c, address)
			require.NoError(t, err, fmt.Sprintf("balance wasn't correctly updated for %s network", net.Name))

			err = c.Close()
			require.NoError(t, err, "Couldn't close the client")

			stopAllContainers(t, cfg.Containers)
		})
	}
}

// 1.9 -> works
// ...
// 1.13 -> works
func TestEthEnvGethCompatibility(t *testing.T) {
	t.Skip("Execute manually when needed")
	l := logging.GetTestLogger(t)

	gethTcs := []string{}
	for i := 9; i < 12; i++ {
		gethTcs = append(gethTcs, fmt.Sprintf("ethereum/client-go:v1.%d.0", i))
	}

	gethTcs = append(gethTcs, ethereum.DefaultGethEth2Image)
	latest, err := FetchLatestEthereumClientDockerImageVersionIfNeed(fmt.Sprintf("ethereum/client-go:%s", AUTOMATIC_STABLE_LATEST_TAG))
	require.NoError(t, err, "Couldn't fetch the latest Go Ethereum version")

	gethTcs = append(gethTcs, latest)
	gethTcs = UniqueStringSlice(gethTcs)

	for _, tc := range gethTcs {
		t.Run(fmt.Sprintf("geth-%s", tc), func(t *testing.T) {
			builder := NewEthereumNetworkBuilder()
			cfg, err := builder.
				WithCustomDockerImages(map[config.ContainerType]string{
					config.ContainerType_ExecutionLayer: tc,
				}).
				Build()
			require.NoError(t, err, "Builder validation failed")

			net, _, err := cfg.Start()
			require.NoError(t, err, "Couldn't start Geth-based network")

			c, err := blockchain.ConnectEVMClient(net, l)
			require.NoError(t, err, "Couldn't connect to the evm client")

			address := common.HexToAddress("0x90F8bf6A479f320ead074411a4B0e7944Ea8c9C1")
			err = sendAndCompareBalances(testcontext.Get(t), c, address)
			require.NoError(t, err, fmt.Sprintf("balance wasn't correctly updated for %s network", net.Name))

			err = c.Close()
			require.NoError(t, err, "Couldn't close the client")

			stopAllContainers(t, cfg.Containers)
		})
	}
}

// v2.40 -> works
// ...
// v2.58.1 -> works
func TestEthEnvErigonCompatibility(t *testing.T) {
	t.Skip("Execute manually when needed")
	l := logging.GetTestLogger(t)

	erigonTcs := []string{}
	for i := 40; i < 59; i++ {
		erigonTcs = append(erigonTcs, fmt.Sprintf("thorax/erigon:v2.%d.0", i))
	}

	erigonTcs = append(erigonTcs, ethereum.DefaultErigonEth2Image)
	latest, err := FetchLatestEthereumClientDockerImageVersionIfNeed(fmt.Sprintf("thorax/erigon:%s", AUTOMATIC_STABLE_LATEST_TAG))
	require.NoError(t, err, "Couldn't fetch the latest Erigon version")

	erigonTcs = append(erigonTcs, latest)
	erigonTcs = UniqueStringSlice(erigonTcs)

	for _, tc := range erigonTcs {
		t.Run(fmt.Sprintf("erigon-%s", tc), func(t *testing.T) {
			builder := NewEthereumNetworkBuilder()
			cfg, err := builder.
				WithCustomDockerImages(map[config.ContainerType]string{
					config.ContainerType_ExecutionLayer: tc,
				}).
				Build()
			require.NoError(t, err, "Builder validation failed")

			net, _, err := cfg.Start()
			require.NoError(t, err, "Couldn't start Erigon-based network")

			c, err := blockchain.ConnectEVMClient(net, l)
			require.NoError(t, err, "Couldn't connect to the evm client")

			address := common.HexToAddress("0x90F8bf6A479f320ead074411a4B0e7944Ea8c9C1")
			err = sendAndCompareBalances(testcontext.Get(t), c, address)
			require.NoError(t, err, fmt.Sprintf("balance wasn't correctly updated for %s network", net.Name))

			err = c.Close()
			require.NoError(t, err, "Couldn't close the client")

			stopAllContainers(t, cfg.Containers)
		})
	}
}

// 22.1 -> works
// ...
// 24.1.2 -> works
func TestEthEnvBesuCompatibility(t *testing.T) {
	t.Skip("Execute manually when needed")
	l := logging.GetTestLogger(t)

	besuTcs := []string{}
	for i := 22; i < 25; i++ {
		besuTcs = append(besuTcs, fmt.Sprintf("hyperledger/besu:%d.1.0", i))
	}

	besuTcs = append(besuTcs, ethereum.DefaultBesuEth2Image)
	latest, err := FetchLatestEthereumClientDockerImageVersionIfNeed(fmt.Sprintf("hyperledger/besu:%s", AUTOMATIC_STABLE_LATEST_TAG))
	require.NoError(t, err, "Couldn't fetch the latest Erigon version")

	besuTcs = append(besuTcs, latest)
	besuTcs = UniqueStringSlice(besuTcs)

	for _, tc := range besuTcs {
		t.Run(fmt.Sprintf("besu-%s", tc), func(t *testing.T) {
			builder := NewEthereumNetworkBuilder()
			cfg, err := builder.
				WithCustomDockerImages(map[config.ContainerType]string{
					config.ContainerType_ExecutionLayer: tc,
				}).
				Build()
			require.NoError(t, err, "Builder validation failed")

			net, _, err := cfg.Start()
			require.NoError(t, err, "Couldn't start Besu-based network")

			c, err := blockchain.ConnectEVMClient(net, l)
			require.NoError(t, err, "Couldn't connect to the evm client")

			address := common.HexToAddress("0x90F8bf6A479f320ead074411a4B0e7944Ea8c9C1")
			err = sendAndCompareBalances(testcontext.Get(t), c, address)
			require.NoError(t, err, fmt.Sprintf("balance wasn't correctly updated for %s network", net.Name))

			err = c.Close()
			require.NoError(t, err, "Couldn't close the client")

			stopAllContainers(t, cfg.Containers)
		})
	}
}
