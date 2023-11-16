package networks

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/k8s/config"
)

func getTestBaseToml() string {
	return `[OCR2]
Enabled = true

[P2P]
[P2P.V2]
AnnounceAddresses = ["0.0.0.0:6690"]
ListenAddresses = ["0.0.0.0:6690"]`
}

func TestAddNetworksConfigNoPyroscope(t *testing.T) {
	t.Setenv("SELECTED_NETWORKS", "SIMULATED")
	s := AddNetworksConfig(getTestBaseToml(), MustGetSelectedNetworksFromEnv()[0])
	require.Contains(t, s, "[[EVM.Nodes]]")
	require.NotContains(t, s, "[Pyroscope]")
}

func TestAddNetworksConfigWithPyroscope(t *testing.T) {
	t.Setenv("SELECTED_NETWORKS", "SIMULATED")
	t.Setenv(config.EnvVarPyroscopeServer, "pyroServer")
	t.Setenv(config.EnvVarPyroscopeEnvironment, "pyroEnv")
	s := AddNetworksConfig(getTestBaseToml(), MustGetSelectedNetworksFromEnv()[0])
	require.Contains(t, s, "[[EVM.Nodes]]")
	require.Contains(t, s, "[Pyroscope]")
	require.Contains(t, s, "pyroServer")
	require.Contains(t, s, "pyroEnv")
}

func TestAddSecretTomlConfig(t *testing.T) {
	s := AddSecretTomlConfig("url", "name", "pass")
	require.Contains(t, s, fmt.Sprintf("URL = '%s'", "url"))
	require.Contains(t, s, fmt.Sprintf("Username = '%s'", "name"))
	require.Contains(t, s, fmt.Sprintf("Password = '%s'", "pass"))
}
