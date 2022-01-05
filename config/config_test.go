package config_test

import (
	"testing"

	"github.com/smartcontractkit/integrations-framework/config"
	"github.com/stretchr/testify/require"
)

func TestFrameworkConfig(t *testing.T) {
	t.Parallel()

	cfg, err := config.LoadFrameworkConfig("./test_framework_config.yml")
	require.NoError(t, err)

	require.Equal(t, "testChainlinkImage", cfg.ChainlinkImage)
	require.Equal(t, "testChainlinkVersion", cfg.ChainlinkVersion)
	require.Equal(t, "testGethImage", cfg.GethImage)
	require.Equal(t, "testGethVersion", cfg.GethVersion)
}

func TestNetworkConfig(t *testing.T) {
	t.Parallel()

	cfg, err := config.LoadNetworksConfig("./test_networks_config.yml")
	require.NoError(t, err)

	require.Equal(t, "huxtable", cfg.SelectedNetworks[0])
}
