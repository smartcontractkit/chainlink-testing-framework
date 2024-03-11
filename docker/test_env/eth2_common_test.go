//go:build eth_env_tests
// +build eth_env_tests

package test_env

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDefaultChainConfig(t *testing.T) {
	t.Parallel()
	config := &EthereumChainConfig{}
	err := config.Default()
	require.NoError(t, err, "Couldn't read default config")

	require.Equal(t, 12, config.SecondsPerSlot, "SecondsPerSlot should be 12")
	require.Equal(t, 6, config.SlotsPerEpoch, "SlotsPerEpoch should be 6")
	require.Equal(t, 15, config.GenesisDelay, "Genesis delay should be 20")
	require.Equal(t, 8, config.ValidatorCount, "Validator count should be 8")
	require.Equal(t, 1337, config.ChainID, "Chain ID should be 1337")
	require.Len(t, config.AddressesToFund, 1, "Should have 1 address to fund")
	require.Equal(t, "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266", config.AddressesToFund[0], "Should have address 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266 to fund")
	require.Len(t, config.HardForkEpochs, 1, "Should have 1 hard fork epoch")
	require.Equal(t, map[string]int{"Deneb": 500}, config.HardForkEpochs, "Should have correct hard fork epochs")
}
