package networks

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/logging"
)

func TestMain(m *testing.M) {
	logging.Init()
	os.Exit(m.Run())
}

func TestMustGetSelectedNetworksFromEnv_Missing_SELECTED_NETWORKS(t *testing.T) {
	defer os.Clearenv()

	require.Panics(t, func() {
		MustGetSelectedNetworksFromEnv()
	})
}

func TestMustGetSelectedNetworksFromEnv_Missing_HTTP_URLS(t *testing.T) {
	defer os.Clearenv()

	networkKey := "ARBITRUM_GOERLI"

	os.Setenv("SELECTED_NETWORKS", networkKey)
	os.Setenv(fmt.Sprintf("%s_URLS", networkKey), "xxxx")
	os.Setenv(fmt.Sprintf("%s_KEYS", networkKey), "xxxx")

	require.PanicsWithError(t, fmt.Sprintf("set %s_HTTP_URLS env var", networkKey), func() {
		MustGetSelectedNetworksFromEnv()
	})
}

func TestMustGetSelectedNetworksFromEnv_Missing_KEYS(t *testing.T) {
	defer os.Clearenv()

	networkKey := "ARBITRUM_GOERLI"

	os.Setenv("SELECTED_NETWORKS", networkKey)
	os.Setenv(fmt.Sprintf("%s_URLS", networkKey), "xxxx")
	os.Setenv(fmt.Sprintf("%s_HTTP_URLS", networkKey), "xxxx")

	require.PanicsWithError(t, fmt.Sprintf("set %s_KEYS env var", networkKey), func() {
		MustGetSelectedNetworksFromEnv()
	})
}

func TestMustGetSelectedNetworksFromEnv_Missing_URLS(t *testing.T) {
	defer os.Clearenv()

	networkKey := "ARBITRUM_GOERLI"

	os.Setenv("SELECTED_NETWORKS", networkKey)
	os.Setenv(fmt.Sprintf("%s_HTTP_URLS", networkKey), "xxxx")
	os.Setenv(fmt.Sprintf("%s_KEYS", networkKey), "xxxx")

	require.PanicsWithError(t, fmt.Sprintf("set %s_URLS env var", networkKey), func() {
		MustGetSelectedNetworksFromEnv()
	})
}

func TestMustGetSelectedNetworksFromEnv_MultipleNetworks(t *testing.T) {
	defer os.Clearenv()

	networkKey := "ARBITRUM_GOERLI,OPTIMISM_GOERLI"
	os.Setenv("SELECTED_NETWORKS", networkKey)

	for _, network := range strings.Split(networkKey, ",") {
		os.Setenv(fmt.Sprintf("%s_URLS", network), "wss://devnet-1.mt/ABC/rpc/")
		os.Setenv(fmt.Sprintf("%s_HTTP_URLS", network), "https://devnet-1.mt/ABC/rpc/")
		os.Setenv(fmt.Sprintf("%s_KEYS", network), "1810868fc221b9f50b5b3e0186d8a5f343f892e51ce12a9e818f936ec0b651ed")
	}

	networks := MustGetSelectedNetworksFromEnv()
	require.Len(t, networks, 2)
	require.Equal(t, "Arbitrum Goerli", networks[0].Name)
	require.Equal(t, "Optimism Goerli", networks[1].Name)
}
