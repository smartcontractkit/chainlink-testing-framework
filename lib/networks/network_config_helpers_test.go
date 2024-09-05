package networks

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/config"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"
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

	l := logging.GetTestLogger(t)
	networkCfg := config.NetworkConfig{}
	err := config.BytesToAnyTomlStruct(l, "test", "", &networkCfg, []byte(networkTOML))
	require.NoError(t, err, "error reading network config")

	pyroCfg := config.PyroscopeConfig{
		Enabled:     ptr.Ptr(true),
		ServerUrl:   ptr.Ptr("pyroServer"),
		Environment: ptr.Ptr("pyroEnv"),
	}

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

	l := logging.GetTestLogger(t)
	networkCfg := config.NetworkConfig{}
	err := config.BytesToAnyTomlStruct(l, "test", "", &networkCfg, []byte(networkTOML))
	require.NoError(t, err, "error reading network config")

	pyroCfg := config.PyroscopeConfig{
		Enabled:     ptr.Ptr(false),
		ServerUrl:   ptr.Ptr("pyroServer"),
		Environment: ptr.Ptr("pyroEnv"),
	}

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
