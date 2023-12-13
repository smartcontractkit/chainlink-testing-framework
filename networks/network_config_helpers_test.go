package networks

import (
	"fmt"
	"testing"

	"github.com/smartcontractkit/chainlink-testing-framework/config"
	"github.com/stretchr/testify/require"
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
	testTOML := `
	[Network]
	selected_networks = ["sIMulated"]"]	
	`
	networkCfg := config.NetworkConfig{}
	err := networkCfg.ApplyDecoded(testTOML)
	require.NoError(t, err, "error reading network config")

	s := AddNetworksConfig(getTestBaseToml(), &config.PyroscopeConfig{}, MustGetSelectedNetworkConfig(&networkCfg)[0])
	require.Contains(t, s, "[[EVM.Nodes]]")
	require.NotContains(t, s, "[Pyroscope]")
}

func TestAddNetworksConfigNoPyroscope(t *testing.T) {
	testTOML := `
	[Network]
	selected_networks = ["SIMULATED"]	
	`
	networkCfg := config.NetworkConfig{}
	err := networkCfg.ApplyDecoded(testTOML)
	require.NoError(t, err, "error reading network config")

	s := AddNetworksConfig(getTestBaseToml(), &config.PyroscopeConfig{}, MustGetSelectedNetworkConfig(&networkCfg)[0])
	require.Contains(t, s, "[[EVM.Nodes]]")
	require.NotContains(t, s, "[Pyroscope]")
}

func TestAddNetworksConfigWithPyroscopeEnabled(t *testing.T) {
	testTOML := `
	[Network]
	selected_networks = ["SIMULATED"]	

	[Pyroscope]
	enabled = true
	server_url = "pyroServer"
	environment = "pyroEnv"
	`

	networkCfg := config.NetworkConfig{}
	err := networkCfg.ApplyDecoded(testTOML)
	require.NoError(t, err, "error reading network config")

	pyroCfg := config.PyroscopeConfig{}
	err = pyroCfg.ReadDecoded(testTOML)
	require.NoError(t, err, "error reading pyroscope config")

	s := AddNetworksConfig(getTestBaseToml(), &pyroCfg, MustGetSelectedNetworkConfig(&networkCfg)[0])
	require.Contains(t, s, "[[EVM.Nodes]]")
	require.Contains(t, s, "[Pyroscope]")
	require.Contains(t, s, "pyroServer")
	require.Contains(t, s, "pyroEnv")
}

func TestAddNetworksConfigWithPyroscopeDisabled(t *testing.T) {
	testTOML := `
	[Network]
	selected_networks = ["SIMULATED"]	

	[Pyroscope]
	enabled = false
	server_url = "pyroServer"
	environment = "pyroEnv"
	`

	networkCfg := config.NetworkConfig{}
	err := networkCfg.ApplyDecoded(testTOML)
	require.NoError(t, err, "error reading network config")

	pyroCfg := config.PyroscopeConfig{}
	err = pyroCfg.ReadDecoded(testTOML)
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
