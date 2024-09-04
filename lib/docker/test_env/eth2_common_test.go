package test_env

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/config"
)

func TestEthEnvReadDefaultChainConfig(t *testing.T) {
	t.Parallel()
	c := &config.EthereumChainConfig{}
	err := c.Default()
	require.NoError(t, err, "Couldn't read default config")

	require.Equal(t, 12, c.SecondsPerSlot, "SecondsPerSlot should be 12")
	require.Equal(t, 6, c.SlotsPerEpoch, "SlotsPerEpoch should be 6")
	require.Equal(t, 15, c.GenesisDelay, "Genesis delay should be 20")
	require.Equal(t, 8, c.ValidatorCount, "Validator count should be 8")
	require.Equal(t, 1337, c.ChainID, "Chain ID should be 1337")
	require.Len(t, c.AddressesToFund, 1, "Should have 1 address to fund")
	require.Equal(t, "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266", c.AddressesToFund[0], "Should have address 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266 to fund")
	require.Len(t, c.HardForkEpochs, 0, "Should have 0 hard fork epochs")
	require.Nil(t, c.HardForkEpochs, "Should have no hard fork epochs")
}
