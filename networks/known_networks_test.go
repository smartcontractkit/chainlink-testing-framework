package networks

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
)

func TestMain(m *testing.M) {
	logging.Init()
	os.Exit(m.Run())
}

func TestMustGetSelectedNetworksFromEnv_Missing_SELECTED_NETWORKS(t *testing.T) {
	require.Panics(t, func() {
		MustGetSelectedNetworksFromEnv()
	})
}

func TestMustGetSelectedNetworksFromEnv_Missing_HTTP_URLS(t *testing.T) {
	networkKey := "ARBITRUM_GOERLI"

	t.Setenv("SELECTED_NETWORKS", networkKey)
	t.Setenv(fmt.Sprintf("%s_URLS", networkKey), "xxxx")
	t.Setenv(fmt.Sprintf("%s_KEYS", networkKey), "xxxx")

	require.PanicsWithError(t, fmt.Sprintf("set %s_HTTP_URLS env var", networkKey), func() {
		MustGetSelectedNetworksFromEnv()
	})
}

func TestMustGetSelectedNetworksFromEnv_Missing_KEYS(t *testing.T) {
	networkKey := "ARBITRUM_GOERLI"

	t.Setenv("SELECTED_NETWORKS", networkKey)
	t.Setenv(fmt.Sprintf("%s_URLS", networkKey), "xxxx")
	t.Setenv(fmt.Sprintf("%s_HTTP_URLS", networkKey), "xxxx")

	require.PanicsWithError(t, fmt.Sprintf("set %s_KEYS env var", networkKey), func() {
		MustGetSelectedNetworksFromEnv()
	})
}

func TestMustGetSelectedNetworksFromEnv_Missing_URLS(t *testing.T) {
	networkKey := "ARBITRUM_GOERLI"

	t.Setenv("SELECTED_NETWORKS", networkKey)
	t.Setenv(fmt.Sprintf("%s_HTTP_URLS", networkKey), "xxxx")
	t.Setenv(fmt.Sprintf("%s_KEYS", networkKey), "xxxx")

	require.PanicsWithError(t, fmt.Sprintf("set %s_URLS env var", networkKey), func() {
		MustGetSelectedNetworksFromEnv()
	})
}

func TestMustGetSelectedNetworksFromEnv_MultipleNetworks(t *testing.T) {
	networkKey := "ARBITRUM_GOERLI,OPTIMISM_GOERLI"
	t.Setenv("SELECTED_NETWORKS", networkKey)

	for _, network := range strings.Split(networkKey, ",") {
		t.Setenv(fmt.Sprintf("%s_URLS", network), "wss://devnet-1.mt/ABC/rpc/")
		t.Setenv(fmt.Sprintf("%s_HTTP_URLS", network), "https://devnet-1.mt/ABC/rpc/")
		t.Setenv(fmt.Sprintf("%s_KEYS", network), "1810868fc221b9f50b5b3e0186d8a5f343f892e51ce12a9e818f936ec0b651ed")
	}

	networks := MustGetSelectedNetworksFromEnv()
	require.Len(t, networks, 2)
	require.Equal(t, "Arbitrum Goerli", networks[0].Name)
	require.Equal(t, "Optimism Goerli", networks[1].Name)
}

func TestNewEVMNetwork(t *testing.T) {
	// Set up a mock mapping for this test
	MappedNetworks = map[string]blockchain.EVMNetwork{
		"VALID_KEY": {
			HTTPURLs: []string{"default_http"},
			URLs:     []string{"default_ws"},
		},
	}

	t.Run("valid networkKey", func(t *testing.T) {
		network, err := NewEVMNetwork("VALID_KEY", nil, nil, nil)
		require.NoError(t, err)
		require.Equal(t, MappedNetworks["VALID_KEY"].HTTPURLs, network.HTTPURLs)
		require.Equal(t, MappedNetworks["VALID_KEY"].URLs, network.URLs)
	})

	t.Run("invalid networkKey", func(t *testing.T) {
		_, err := NewEVMNetwork("INVALID_KEY", nil, nil, nil)
		require.Error(t, err)
		require.Contains(t, err.Error(), "network key: 'INVALID_KEY' is invalid")
	})

	t.Run("overwriting default values", func(t *testing.T) {
		walletKeys := []string{"1810868fc221b9f50b5b3e0186d8a5f343f892e51ce12a9e818f936ec0b651ed"}
		httpUrls := []string{"http://newurl.com"}
		wsUrls := []string{"ws://newwsurl.com"}

		network, err := NewEVMNetwork("VALID_KEY", walletKeys, httpUrls, wsUrls)
		require.NoError(t, err)
		require.Equal(t, httpUrls, network.HTTPURLs)
		require.Equal(t, wsUrls, network.URLs)
		require.Equal(t, walletKeys, network.PrivateKeys)
	})
}
