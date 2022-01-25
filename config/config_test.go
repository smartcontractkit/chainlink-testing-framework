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

	testEnvVals := map[string]string{
		"test_string_val": "someString",
		"test_int_val":    "420",
	}
	require.Equal(t, testEnvVals, cfg.ChainlinkEnvValues)
}

func TestNetworkConfig(t *testing.T) {
	t.Parallel()

	cfg, err := config.LoadNetworksConfig("./test_networks_config.yml")
	require.NoError(t, err)

	require.Equal(t, "huxtable", cfg.SelectedNetworks[0])
}

func TestChartCreation(t *testing.T) {
	t.Parallel()

	emptyConfig := config.FrameworkConfig{}
	emptyChartString := `{}`
	chart, err := emptyConfig.CreateChartOverrrides()
	require.NoError(t, err)
	require.JSONEq(t, emptyChartString, chart, "Expected an empty config to produce an empty object for chart overrides")

	gethOnlyConfig := config.FrameworkConfig{
		GethImage:   "testGethImage",
		GethVersion: "testGethVersion",
		GethArgs: []interface{}{
			"some",
			"args",
			15,
			"--address",
			"0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266",
		},
	}
	gethOnlyChartString := `{
		"geth":{
			"values":{
				"geth":{
					"image":{
						"image":"testGethImage",
						"version":"testGethVersion"
					}
				},
				"args": [
					"some",
					"args",
					15,
					"--address",
					"0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266"
				]
			}
		}
	}`
	chart, err = gethOnlyConfig.CreateChartOverrrides()
	require.NoError(t, err)
	require.JSONEq(t, gethOnlyChartString, chart, "Expected a config with only geth overrides")

	chainlinkOnlyConfig := config.FrameworkConfig{
		ChainlinkImage:   "testChainlinkImage",
		ChainlinkVersion: "testChainlinkVersion",
	}
	chainlinkOnlyChartString := `{
		"chainlink":{
			"values":{
				"chainlink":{
					"image":{
						"image":"testChainlinkImage",
						"version":"testChainlinkVersion"
					}
				}
			}
		}
	}`
	chart, err = chainlinkOnlyConfig.CreateChartOverrrides()
	require.NoError(t, err)
	require.JSONEq(t, chainlinkOnlyChartString, chart, "Expected a config with only chainlink image and version")

	chainlinkOnlyConfig.ChainlinkEnvValues = map[string]string{
		"test_string_val": "someString",
		"test_int_val":    "420",
	}
	chainlinkOnlyChartString = `{
		"chainlink":{
			"values":{
				"chainlink":{
					"image":{
						"image":"testChainlinkImage",
						"version":"testChainlinkVersion"
					}
				},
				"env": {
					"test_string_val": "someString",
					"test_int_val": "420"
				}
			}
		}
	}`
	chart, err = chainlinkOnlyConfig.CreateChartOverrrides()
	require.NoError(t, err)
	require.JSONEq(t, chainlinkOnlyChartString, chart, "Expected a config with chainlink image, version, and env vars")
}
