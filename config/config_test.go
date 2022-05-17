package config_test

import (
	"os"
	"testing"

	"github.com/smartcontractkit/chainlink-testing-framework/config"
	"github.com/stretchr/testify/require"
)

func setEnvVars(t *testing.T) {
	err := os.Setenv("FRAMEWORK_CONFIG_FILE", "./test_framework_config.yml")
	require.NoError(t, err)
	err = os.Setenv("NETWORKS_CONFIG_FILE", "./test_networks_config.yml")
	require.NoError(t, err)
}

func TestFrameworkConfig(t *testing.T) {
	setEnvVars(t)
	err := config.LoadFromEnv()
	require.NoError(t, err)

	require.Equal(t, "testChainlinkImage", config.ProjectConfig.FrameworkConfig.ChainlinkImage)
	require.Equal(t, "testChainlinkVersion", config.ProjectConfig.FrameworkConfig.ChainlinkVersion)

	testEnvVals := map[string]interface{}{
		"test_string_val": "someString",
		"test_int_val":    420,
	}
	require.Equal(t, testEnvVals, config.ProjectConfig.FrameworkConfig.ChainlinkEnvValues)
}

func TestNetworkConfig(t *testing.T) {
	setEnvVars(t)
	err := config.LoadFromEnv()
	require.NoError(t, err)
	require.Equal(t, "huxtable", config.ProjectConfig.NetworksConfig.SelectedNetworks[0])
}

func TestChainlinkValues(t *testing.T) {
	t.Parallel()

	config.ProjectConfig.FrameworkConfig = &config.FrameworkConfig{}
	loadedVals := config.ChainlinkVals()

	require.Equal(t, map[string]interface{}{}, loadedVals)

	config.ProjectConfig.FrameworkConfig = &config.FrameworkConfig{
		ChainlinkImage:   "image",
		ChainlinkVersion: "version",
	}
	loadedVals = config.ChainlinkVals()

	require.Equal(t, map[string]interface{}{
		"chainlink": map[string]interface{}{
			"image": map[string]interface{}{
				"image":   "image",
				"version": "version",
			},
		},
	}, loadedVals)

	config.ProjectConfig.FrameworkConfig = &config.FrameworkConfig{
		ChainlinkImage:   "image",
		ChainlinkVersion: "version",
		ChainlinkEnvValues: map[string]interface{}{
			"env": "value",
		},
	}
	loadedVals = config.ChainlinkVals()

	require.Equal(t, map[string]interface{}{
		"chainlink": map[string]interface{}{
			"image": map[string]interface{}{
				"image":   "image",
				"version": "version",
			},
		},
		"env": map[string]interface{}{
			"env": "value",
		},
	}, loadedVals)
}
