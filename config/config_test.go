package config_test

import (
	"errors"
	"os"
	"testing"

	"github.com/smartcontractkit/integrations-framework/config"
	"github.com/smartcontractkit/integrations-framework/utils"
	"github.com/stretchr/testify/require"
)

func TestFrameworkConfig(t *testing.T) {
	t.Parallel()

	cfg, err := config.LoadFrameworkConfig("./test_framework_config.yml")
	require.NoError(t, err)

	require.Equal(t, "testChainlinkImage", cfg.ChainlinkImage)
	require.Equal(t, "testChainlinkVersion", cfg.ChainlinkVersion)

	testEnvVals := map[string]interface{}{
		"test_string_val": "someString",
		"test_int_val":    420,
	}
	require.Equal(t, testEnvVals, cfg.ChainlinkEnvValues)
}

func TestNetworkConfig(t *testing.T) {
	t.Parallel()

	cfg, err := config.LoadNetworksConfig("./test_networks_config.yml")
	require.NoError(t, err)

	require.Equal(t, "huxtable", cfg.SelectedNetworks[0])
}

func TestChainlinkValues(t *testing.T) {
	t.Parallel()

	config.ProjectFrameworkSettings = &config.FrameworkConfig{}
	loadedVals := config.ChainlinkVals()

	require.Equal(t, map[string]interface{}{}, loadedVals)

	config.ProjectFrameworkSettings = &config.FrameworkConfig{
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

	config.ProjectFrameworkSettings = &config.FrameworkConfig{
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

func TestRemoteRunnerConfig(t *testing.T) {
	t.Parallel()

	// Check if the config file already exists, if so, delete it
	if _, err := os.Stat(utils.RemoteRunnerConfigLocation); err == nil {
		err := os.Remove(utils.RemoteRunnerConfigLocation)
		require.NoError(t, err)
	} else if !errors.Is(err, os.ErrNotExist) {
		require.NoError(t, err)
	}
	_, err := config.ReadWriteRemoteRunnerConfig()
	require.Error(t, err, "Wrote an example config file at %s. Please fill in values and log back in", utils.RemoteRunnerConfigLocation)
	require.FileExists(t, utils.RemoteRunnerConfigLocation)

	remoteConfig, err := config.ReadWriteRemoteRunnerConfig()
	require.NoError(t, err)
	require.Equal(t, "@soak-ocr", remoteConfig.TestRegex)
	require.Equal(t, "abcdefg", remoteConfig.SlackAPIKey)
}
