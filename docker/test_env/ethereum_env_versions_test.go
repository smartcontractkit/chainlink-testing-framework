package test_env

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/testcontext"
	"github.com/stretchr/testify/require"
)

var twentySecDuration = time.Duration(20 * time.Second)
var stopAllContainers = func(t *testing.T, containers []EthereumNetworkContainer) {
	for _, c := range containers {
		err := (*c.Container).Stop(context.Background(), &twentySecDuration)
		require.NoError(t, err, "Couldn't stop container")
	}

}

// 1.16 -> works
// ...
// 1.25 -> works
func TestNethermindCompatiblity(t *testing.T) {
	t.Skip("Execute manually when needed")
	l := logging.GetTestLogger(t)

	nethermindTcs := []string{}
	for i := 16; i < 26; i++ {
		nethermindTcs = append(nethermindTcs, fmt.Sprintf("nethermind/nethermind:1.%d.0", i))
	}

	nethermindTcs = append(nethermindTcs, defaultNethermindPosImage)

	for _, tc := range nethermindTcs {
		t.Run(fmt.Sprintf("nethermind-%s", tc), func(t *testing.T) {
			builder := NewEthereumNetworkBuilder()
			cfg, err := builder.
				WithExecutionLayer(ExecutionLayer_Nethermind).
				WithCustomDockerImages(map[ContainerType]string{
					ContainerType_Nethermind: tc,
				}).
				Build()
			require.NoError(t, err, "Builder validation failed")

			net, _, err := cfg.Start()
			require.NoError(t, err, "Couldn't start Nethermind-based network")

			c, err := blockchain.ConnectEVMClient(net, l)
			require.NoError(t, err, "Couldn't connect to the evm client")

			address := common.HexToAddress("0x90F8bf6A479f320ead074411a4B0e7944Ea8c9C1")
			err = sendAndCompareBalances(testcontext.Get(t), c, address)
			require.NoError(t, err, fmt.Sprintf("balance wasn't correctly updated when %s network", net.Name))

			err = c.Close()
			require.NoError(t, err, "Couldn't close the client")

			stopAllContainers(t, cfg.Containers)
		})
	}
}

// 1.9 -> works
// ...
// 1.13 -> works
func TestGethCompatiblity(t *testing.T) {
	t.Skip("Execute manually when needed")
	l := logging.GetTestLogger(t)

	gethTcs := []string{}
	for i := 9; i < 12; i++ {
		gethTcs = append(gethTcs, fmt.Sprintf("ethereum/client-go:v1.%d.0", i))
	}

	gethTcs = append(gethTcs, defaultGethPosImage)

	for _, tc := range gethTcs {
		t.Run(fmt.Sprintf("geth-%s", tc), func(t *testing.T) {
			builder := NewEthereumNetworkBuilder()
			cfg, err := builder.
				WithCustomDockerImages(map[ContainerType]string{
					ContainerType_Geth: tc,
				}).
				Build()
			require.NoError(t, err, "Builder validation failed")

			net, _, err := cfg.Start()
			require.NoError(t, err, "Couldn't start Geth-based network")

			c, err := blockchain.ConnectEVMClient(net, l)
			require.NoError(t, err, "Couldn't connect to the evm client")

			address := common.HexToAddress("0x90F8bf6A479f320ead074411a4B0e7944Ea8c9C1")
			err = sendAndCompareBalances(testcontext.Get(t), c, address)
			require.NoError(t, err, fmt.Sprintf("balance wasn't correctly updated when %s network", net.Name))

			err = c.Close()
			require.NoError(t, err, "Couldn't close the client")

			stopAllContainers(t, cfg.Containers)
		})
	}
}

// v2.40 -> works
// ...
// v2.58 -> works
func TestErigonCompatiblity(t *testing.T) {
	t.Skip("Execute manually when needed")
	l := logging.GetTestLogger(t)

	erigonTcs := []string{}
	for i := 40; i < 59; i++ {
		erigonTcs = append(erigonTcs, fmt.Sprintf("thorax/erigon:v2.%d.0", i))
	}

	erigonTcs = append(erigonTcs, defaultErigonPosImage)

	for _, tc := range erigonTcs {
		t.Run(fmt.Sprintf("erigon-%s", tc), func(t *testing.T) {
			builder := NewEthereumNetworkBuilder()
			cfg, err := builder.
				WithCustomDockerImages(map[ContainerType]string{
					ContainerType_Erigon: tc,
				}).
				Build()
			require.NoError(t, err, "Builder validation failed")

			net, _, err := cfg.Start()
			require.NoError(t, err, "Couldn't start Erigon-based network")

			c, err := blockchain.ConnectEVMClient(net, l)
			require.NoError(t, err, "Couldn't connect to the evm client")

			address := common.HexToAddress("0x90F8bf6A479f320ead074411a4B0e7944Ea8c9C1")
			err = sendAndCompareBalances(testcontext.Get(t), c, address)
			require.NoError(t, err, fmt.Sprintf("balance wasn't correctly updated when %s network", net.Name))

			err = c.Close()
			require.NoError(t, err, "Couldn't close the client")

			stopAllContainers(t, cfg.Containers)
		})
	}
}

// 22.1 -> works
// ...
// 24.1 -> works
func TestBesuCompatiblity(t *testing.T) {
	t.Skip("Execute manually when needed")
	l := logging.GetTestLogger(t)

	besuTcs := []string{}
	for i := 22; i < 25; i++ {
		besuTcs = append(besuTcs, fmt.Sprintf("hyperledger/besu:%d.1.0", i))
	}

	besuTcs = append(besuTcs, defaultBesuPosImage)

	for _, tc := range besuTcs {
		t.Run(fmt.Sprintf("besu-%s", tc), func(t *testing.T) {
			builder := NewEthereumNetworkBuilder()
			cfg, err := builder.
				WithCustomDockerImages(map[ContainerType]string{
					ContainerType_Besu: tc,
				}).
				Build()
			require.NoError(t, err, "Builder validation failed")

			net, _, err := cfg.Start()
			require.NoError(t, err, "Couldn't start Besu-based network")

			c, err := blockchain.ConnectEVMClient(net, l)
			require.NoError(t, err, "Couldn't connect to the evm client")

			address := common.HexToAddress("0x90F8bf6A479f320ead074411a4B0e7944Ea8c9C1")
			err = sendAndCompareBalances(testcontext.Get(t), c, address)
			require.NoError(t, err, fmt.Sprintf("balance wasn't correctly updated when %s network", net.Name))

			err = c.Close()
			require.NoError(t, err, "Couldn't close the client")

			stopAllContainers(t, cfg.Containers)
		})
	}
}
