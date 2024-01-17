package networks

import (
	"encoding/base64"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/config"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
)

func TestMain(m *testing.M) {
	logging.Init()
	os.Exit(m.Run())
}

func TestMustGetSelectedNetworkConfig_MissingSelectedNetwork(t *testing.T) {
	require.Panics(t, func() {
		MustGetSelectedNetworkConfig(&config.NetworkConfig{})
	})
}

func TestMustGetSelectedNetworkConfig_Missing_RpcHttpUrls(t *testing.T) {
	networkName := "arbitrum_goerli"
	testTOML := `
	selected_networks = ["arbitrum_goerli"]
	
	[RpcWsUrls]
	arbitrum_goerli = ["wss://devnet-1.mt/ABC/rpc/"]

	[WalletKeys]
	arbitrum_goerli = ["1810868fc221b9f50b5b3e0186d8a5f343f892e51ce12a9e818f936ec0b651ed"]
	`

	l := logging.GetTestLogger(t)
	networkCfg := config.NetworkConfig{}
	err := config.BytesToAnyTomlStruct(l, "test", "", &networkCfg, []byte(testTOML))

	require.NoError(t, err, "error reading network config")

	require.PanicsWithError(t, fmt.Sprintf("no rpc http urls found in config for '%s' network", networkName), func() {
		MustGetSelectedNetworkConfig(&networkCfg)
	})
}

func TestMustGetSelectedNetworkConfig_Missing_RpcWsUrls(t *testing.T) {
	networkName := "arbitrum_goerli"
	testTOML := `
	selected_networks = ["arbitrum_goerli"]

	[RpcHttpUrls]
	arbitrum_goerli = ["https://devnet-1.mt/ABC/rpc/"]

	[WalletKeys]
	arbitrum_goerli = ["1810868fc221b9f50b5b3e0186d8a5f343f892e51ce12a9e818f936ec0b651ed"]
	`

	l := logging.GetTestLogger(t)
	networkCfg := config.NetworkConfig{}
	err := config.BytesToAnyTomlStruct(l, "test", "", &networkCfg, []byte(testTOML))
	require.NoError(t, err, "error reading network config")

	require.PanicsWithError(t, fmt.Sprintf("no rpc ws urls found in config for '%s' network", networkName), func() {
		MustGetSelectedNetworkConfig(&networkCfg)
	})
}

func TestMustGetSelectedNetworkConfig_Missing_WalletKeys(t *testing.T) {
	networkName := "arbitrum_goerli"
	testTOML := `
	selected_networks = ["arbitrum_goerli"]

	[RpcHttpUrls]
	arbitrum_goerli = ["https://devnet-1.mt/ABC/rpc/"]

	[RpcWsUrls]
	arbitrum_goerli = ["wss://devnet-1.mt/ABC/rpc/"]
	`

	l := logging.GetTestLogger(t)
	networkCfg := config.NetworkConfig{}
	err := config.BytesToAnyTomlStruct(l, "test", "", &networkCfg, []byte(testTOML))
	require.NoError(t, err, "error reading network config")

	require.PanicsWithError(t, fmt.Sprintf("no wallet keys found in config for '%s' network", networkName), func() {
		MustGetSelectedNetworkConfig(&networkCfg)
	})
}

func TestMustGetSelectedNetworkConfig_DefaultUrlsFromEnv(t *testing.T) {
	networkConfigTOML := `
	[RpcHttpUrls]
	arbitrum_goerli = ["https://devnet-1.mt/ABC/rpc/"]

	[RpcWsUrls]
	arbitrum_goerli = ["wss://devnet-1.mt/ABC/rpc/"]
	`
	encoded := base64.StdEncoding.EncodeToString([]byte(networkConfigTOML))
	err := os.Setenv("BASE64_NETWORK_CONFIG", encoded)
	require.NoError(t, err, "error setting env var")

	testTOML := `
	selected_networks = ["arbitrum_goerli"]

	[WalletKeys]
	arbitrum_goerli = ["1810868fc221b9f50b5b3e0186d8a5f343f892e51ce12a9e818f936ec0b651ed"]
	`

	l := logging.GetTestLogger(t)
	networkCfg := config.NetworkConfig{}
	err = config.BytesToAnyTomlStruct(l, "test", "", &networkCfg, []byte(testTOML))
	require.NoError(t, err, "error reading network config")

	networkCfg.UpperCaseNetworkNames()

	err = networkCfg.Default()
	require.NoError(t, err, "error reading default network config")

	err = networkCfg.Validate()
	require.NoError(t, err, "error validating network config")

	networks := MustGetSelectedNetworkConfig(&networkCfg)
	require.Len(t, networks, 1, "should have 1 network")
	require.Equal(t, "Arbitrum Goerli", networks[0].Name, "first network should be arbitrum")
	require.Equal(t, []string{"wss://devnet-1.mt/ABC/rpc/"}, networks[0].URLs, "should have default ws url for arbitrum")
	require.Equal(t, []string{"https://devnet-1.mt/ABC/rpc/"}, networks[0].HTTPURLs, "should have default http url for arbitrum")
	require.Equal(t, []string{"1810868fc221b9f50b5b3e0186d8a5f343f892e51ce12a9e818f936ec0b651ed"}, networks[0].PrivateKeys, "should have correct wallet key for arbitrum")
}

func TestMustGetSelectedNetworkConfig_MultipleNetworks(t *testing.T) {
	testTOML := `
	selected_networks = ["arbitrum_goerli", "optimism_goerli"]

	[RpcHttpUrls]
	arbitrum_goerli = ["https://devnet-1.mt/ABC/rpc/"]
	optimism_goerli = ["https://devnet-1.mt/ABC/rpc/"]

	[RpcWsUrls]
	arbitrum_goerli = ["wss://devnet-1.mt/ABC/rpc/"]
	optimism_goerli = ["wss://devnet-1.mt/ABC/rpc/"]

	[WalletKeys]
	arbitrum_goerli = ["1810868fc221b9f50b5b3e0186d8a5f343f892e51ce12a9e818f936ec0b651ed"]
	optimism_goerli = ["1810868fc221b9f50b5b3e0186d8a5f343f892e51ce12a9e818f936ec0b651ed"]
	`

	l := logging.GetTestLogger(t)
	networkCfg := config.NetworkConfig{}
	err := config.BytesToAnyTomlStruct(l, "test", "", &networkCfg, []byte(testTOML))
	require.NoError(t, err, "error reading network config")

	networks := MustGetSelectedNetworkConfig(&networkCfg)
	require.Len(t, networks, 2)
	require.Equal(t, "Arbitrum Goerli", networks[0].Name)
	require.Equal(t, "Optimism Goerli", networks[1].Name)
}

func TestMustGetSelectedNetworkConfig_DefaultUrlsFromSecret_OverrideOne(t *testing.T) {
	networkConfigTOML := `
	[RpcHttpUrls]
	arbitrum_goerli = ["https://devnet-1.mt/ABC/rpc/"]
	optimism_goerli = ["https://devnet-1.mt/ABC/rpc/"]

	[RpcWsUrls]
	arbitrum_goerli = ["wss://devnet-1.mt/ABC/rpc/"]
	optimism_goerli = ["wss://devnet-1.mt/ABC/rpc/"]
	`
	encoded := base64.StdEncoding.EncodeToString([]byte(networkConfigTOML))
	err := os.Setenv("BASE64_NETWORK_CONFIG", encoded)
	require.NoError(t, err, "error setting env var")

	testTOML := `
	selected_networks = ["arbitrum_goerli", "optimism_goerli"]

	[RpcHttpUrls]
	arbitrum_goerli = ["https://devnet-2.mt/ABC/rpc/"]

	[WalletKeys]
	arbitrum_goerli = ["1810868fc221b9f50b5b3e0186d8a5f343f892e51ce12a9e818f936ec0b651ed"]
	optimism_goerli = ["1810868fc221b9f50b5b3e0186d8a5f343f892e51ce12a9e818f936ec0b651ed"]
	`

	l := logging.GetTestLogger(t)
	networkCfg := config.NetworkConfig{}
	err = config.BytesToAnyTomlStruct(l, "test", "", &networkCfg, []byte(testTOML))
	require.NoError(t, err, "error reading network config")

	networkCfg.UpperCaseNetworkNames()
	err = networkCfg.Default()
	require.NoError(t, err, "error reading default network config")

	networks := MustGetSelectedNetworkConfig(&networkCfg)
	require.Len(t, networks, 2, "should have 2 networks")
	require.Equal(t, "Arbitrum Goerli", networks[0].Name, "first network should be arbitrum")
	require.Equal(t, []string{"wss://devnet-1.mt/ABC/rpc/"}, networks[0].URLs, "should have default ws url for arbitrum")
	require.Equal(t, []string{"https://devnet-2.mt/ABC/rpc/"}, networks[0].HTTPURLs, "should have overridden http url for arbitrum")
	require.Equal(t, []string{"1810868fc221b9f50b5b3e0186d8a5f343f892e51ce12a9e818f936ec0b651ed"}, networks[0].PrivateKeys, "should have correct wallet key for arbitrum")

	require.Equal(t, "Optimism Goerli", networks[1].Name, "first network should be optimism")
	require.Equal(t, []string{"wss://devnet-1.mt/ABC/rpc/"}, networks[1].URLs, "should have default ws url for optimism")
	require.Equal(t, []string{"https://devnet-1.mt/ABC/rpc/"}, networks[1].HTTPURLs, "should have default http url for optimism")
	require.Equal(t, []string{"1810868fc221b9f50b5b3e0186d8a5f343f892e51ce12a9e818f936ec0b651ed"}, networks[1].PrivateKeys, "should have correct wallet key for optimism")
}

func TestNewEVMNetwork(t *testing.T) {
	// Set up a mock mapping and revert it after test
	originalMappedNetworks := MappedNetworks
	MappedNetworks = map[string]blockchain.EVMNetwork{
		"VALID_KEY": {
			HTTPURLs: []string{"default_http"},
			URLs:     []string{"default_ws"},
		},
	}
	defer func() {
		MappedNetworks = originalMappedNetworks
	}()

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
