package networks

import (
	"fmt"
	"testing"

	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/config"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
)

func getTestBaseToml() string {
	return `[OCR2]
Enabled = true

[P2P]
[P2P.V2]
AnnounceAddresses = ["0.0.0.0:6690"]
ListenAddresses = ["0.0.0.0:6690"]`
}

func TestAddNetworksConfigCaseInsensitive(t *testing.T) {
	networkTOML := `
	selected_networks = ["sIMulated"]
	`

	l := logging.GetTestLogger(t)
	networkCfg := config.NetworkConfig{}
	err := config.BytesToAnyTomlStruct(l, "test", "", &networkCfg, []byte(networkTOML))
	require.NoError(t, err, "error reading network config")

	s := AddNetworksConfig(getTestBaseToml(), &config.PyroscopeConfig{}, MustGetSelectedNetworkConfig(&networkCfg)[0])
	require.Contains(t, s, "[[EVM.Nodes]]")
	require.NotContains(t, s, "[Pyroscope]")
}

func TestAddNetworksConfigNoPyroscope(t *testing.T) {
	networkTOML := `
	selected_networks = ["SIMULATED"]
	`

	l := logging.GetTestLogger(t)
	networkCfg := config.NetworkConfig{}
	err := config.BytesToAnyTomlStruct(l, "test", "", &networkCfg, []byte(networkTOML))
	require.NoError(t, err, "error reading network config")

	s := AddNetworksConfig(getTestBaseToml(), &config.PyroscopeConfig{}, MustGetSelectedNetworkConfig(&networkCfg)[0])
	require.Contains(t, s, "[[EVM.Nodes]]")
	require.NotContains(t, s, "[Pyroscope]")
}

func TestAddNetworksConfigWithPyroscopeEnabled(t *testing.T) {
	networkTOML := `
	selected_networks = ["SIMULATED"]
	`
	peryscopeTOML := `
	enabled = true
	server_url = "pyroServer"
	environment = "pyroEnv"
	`

	l := logging.GetTestLogger(t)
	networkCfg := config.NetworkConfig{}
	err := config.BytesToAnyTomlStruct(l, "test", "", &networkCfg, []byte(networkTOML))
	require.NoError(t, err, "error reading network config")

	pyroCfg, err := readPyroscopeConfig(peryscopeTOML)
	require.NoError(t, err, "error reading pyroscope config")

	s := AddNetworksConfig(getTestBaseToml(), &pyroCfg, MustGetSelectedNetworkConfig(&networkCfg)[0])
	require.Contains(t, s, "[[EVM.Nodes]]")
	require.Contains(t, s, "[Pyroscope]")
	require.Contains(t, s, "pyroServer")
	require.Contains(t, s, "pyroEnv")
}

func TestAddNetworksConfigWithPyroscopeDisabled(t *testing.T) {
	networkTOML := `
	selected_networks = ["SIMULATED"]
	`
	peryscopeTOML := `
	enabled = false
	server_url = "pyroServer"
	environment = "pyroEnv"
	`

	l := logging.GetTestLogger(t)
	networkCfg := config.NetworkConfig{}
	err := config.BytesToAnyTomlStruct(l, "test", "", &networkCfg, []byte(networkTOML))
	require.NoError(t, err, "error reading network config")

	pyroCfg, err := readPyroscopeConfig(peryscopeTOML)
	require.NoError(t, err, "error reading pyroscope config")

	s := AddNetworksConfig(getTestBaseToml(), &pyroCfg, MustGetSelectedNetworkConfig(&networkCfg)[0])
	require.Contains(t, s, "[[EVM.Nodes]]")
	require.NotContains(t, s, "[Pyroscope]")
}

func TestAddSecretTomlConfig(t *testing.T) {
	s := AddSecretTomlConfig("url", "name", "pass")
	require.Contains(t, s, fmt.Sprintf("URL = '%s'", "url"))
	require.Contains(t, s, fmt.Sprintf("Username = '%s'", "name"))
	require.Contains(t, s, fmt.Sprintf("Password = '%s'", "pass"))
}

func readPyroscopeConfig(configDecoded string) (config.PyroscopeConfig, error) {
	var cfg config.PyroscopeConfig
	err := toml.Unmarshal([]byte(configDecoded), &cfg)
	if err != nil {
		return config.PyroscopeConfig{}, fmt.Errorf("error unmarshalling pyroscope config: %w", err)
	}

	return cfg, nil
}
