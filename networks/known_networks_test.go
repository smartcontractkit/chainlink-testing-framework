package networks

import (
	"encoding/base64"
	"fmt"
	"os"
	"testing"

	"github.com/pelletier/go-toml/v2"
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

	// err = networkCfg.Default()
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
	// err = networkCfg.Default()
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
		network := MappedNetworks["VALID_KEY"]
		err := NewEVMNetwork(&network, nil, nil, nil)
		require.NoError(t, err)
		require.Equal(t, MappedNetworks["VALID_KEY"].HTTPURLs, network.HTTPURLs)
		require.Equal(t, MappedNetworks["VALID_KEY"].URLs, network.URLs)
	})

	t.Run("overwriting default values", func(t *testing.T) {
		walletKeys := []string{"1810868fc221b9f50b5b3e0186d8a5f343f892e51ce12a9e818f936ec0b651ed"}
		httpUrls := []string{"http://newurl.com"}
		wsUrls := []string{"ws://newwsurl.com"}
		network := MappedNetworks["VALID_KEY"]
		err := NewEVMNetwork(&network, walletKeys, httpUrls, wsUrls)
		require.NoError(t, err)
		require.Equal(t, httpUrls, network.HTTPURLs)
		require.Equal(t, wsUrls, network.URLs)
		require.Equal(t, walletKeys, network.PrivateKeys)
	})
}

func TestVariousNetworkConfig(t *testing.T) {
	newNetwork := blockchain.EVMNetwork{
		Name:                      "new_test_network",
		ChainID:                   100009,
		Simulated:                 true,
		ChainlinkTransactionLimit: 5000,
		MinimumConfirmations:      1,
		GasEstimationBuffer:       10000,
		ClientImplementation:      blockchain.EthereumClientImplementation,
		HTTPURLs: []string{
			"http://localhost:8545",
		},
		URLs: []string{
			"ws://localhost:8546",
		},
		SupportsEIP1559: true,
		DefaultGasLimit: 6000000,
		PrivateKeys: []string{
			"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
		},
	}
	forkedNetwork := newNetwork
	forkedNetwork.HTTPURLs = nil
	forkedNetwork.URLs = nil
	forkedNetwork.PrivateKeys = nil
	t.Cleanup(func() {
		ArbitrumGoerli.URLs = []string{}
		ArbitrumGoerli.HTTPURLs = []string{}
		ArbitrumGoerli.PrivateKeys = []string{}
		OptimismGoerli.URLs = []string{}
		OptimismGoerli.HTTPURLs = []string{}
		OptimismGoerli.PrivateKeys = []string{}
	})
	ArbitrumGoerli.URLs = []string{"wss://devnet-1.mt/ABC/rpc/"}
	ArbitrumGoerli.HTTPURLs = []string{"https://devnet-1.mt/ABC/rpc/"}
	ArbitrumGoerli.PrivateKeys = []string{"1810868fc221b9f50b5b3e0186d8a5f343f892e51ce12a9e818f936ec0b651ed"}
	OptimismGoerli.URLs = []string{"wss://devnet-1.mt/ABC/rpc/"}
	OptimismGoerli.HTTPURLs = []string{"https://devnet-1.mt/ABC/rpc/"}
	OptimismGoerli.PrivateKeys = []string{"1810868fc221b9f50b5b3e0186d8a5f343f892e51ce12a9e818f936ec0b651ed"}

	testcases := []struct {
		name                 string
		configOverrideTOML   string
		isNetworkConfigError bool
		isEVMNetworkError    bool
		expNetworks          []blockchain.EVMNetwork
		setupFunc            func()
		cleanupFunc          func()
	}{
		// 		{
		// 			name: "case insensitive network key to EVMNetworks",
		// 			configOverrideTOML: `
		// selected_networks = ["NEW_NETWORK"]

		// [EVMNetworks.new_Network]
		// evm_name = "new_test_network"
		// evm_chain_id = 100009
		// evm_urls = ["ws://localhost:8546"]
		// evm_http_urls = ["http://localhost:8545"]
		// client_implementation = "Ethereum"
		// evm_keys = ["ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"]
		// evm_simulated = true
		// evm_chainlink_transaction_limit = 5000
		// evm_minimum_confirmations = 1
		// evm_gas_estimation_buffer = 10000
		// evm_supports_eip1559 = true
		// evm_default_gas_limit = 6000000
		// `,
		// 			expNetworks: []blockchain.EVMNetwork{newNetwork},
		// 		},
		// 		{
		// 			name: "case insensitive network key of fork config",
		// 			configOverrideTOML: `
		// selected_networks = ["KROMA_SEPOLIA"]

		// [AnvilConfigs.kroma_SEPOLIA]
		// url = "ws://localhost:8546"
		// block_number = 100
		// `,
		// 			expNetworks: []blockchain.EVMNetwork{KromaSepolia},
		// 		},
		// 		{
		// 			name: "override with new AnvilConfigs and new EVMNetworks",
		// 			configOverrideTOML: `
		// selected_networks = ["KROMA_SEPOLIA","NEW_NETWORK"]

		// [AnvilConfigs.KROMA_SEPOLIA]
		// url = "ws://localhost:8546"
		// block_number = 100

		// [EVMNetworks.new_network]
		// evm_name = "new_test_network"
		// evm_chain_id = 100009
		// evm_simulated = true
		// client_implementation = "Ethereum"
		// evm_chainlink_transaction_limit = 5000
		// evm_minimum_confirmations = 1
		// evm_gas_estimation_buffer = 10000
		// evm_supports_eip1559 = true
		// evm_default_gas_limit = 6000000

		// [AnvilConfigs.new_network]
		// url = "ws://localhost:8546"
		// block_number = 100
		// `,
		// 			expNetworks: []blockchain.EVMNetwork{KromaSepolia, forkedNetwork},
		// 		},
		// 		{
		// 			name: "forked network for existing network",
		// 			configOverrideTOML: `
		// selected_networks = ["KROMA_SEPOLIA"]

		// [AnvilConfigs.KROMA_SEPOLIA]
		// url = "ws://localhost:8546"
		// block_number = 100
		// `,
		// 			expNetworks: []blockchain.EVMNetwork{KromaSepolia},
		// 		},

		// 		{
		// 			name: "forked network for new network",
		// 			configOverrideTOML: `
		// selected_networks = ["new_network"]

		// [EVMNetworks.new_network]
		// evm_name = "new_test_network"
		// evm_chain_id = 100009
		// evm_simulated = true
		// client_implementation = "Ethereum"
		// evm_chainlink_transaction_limit = 5000
		// evm_minimum_confirmations = 1
		// evm_gas_estimation_buffer = 10000
		// evm_supports_eip1559 = true
		// evm_default_gas_limit = 6000000

		// [AnvilConfigs.new_network]
		// url = "ws://localhost:8546"
		// block_number = 100
		// `,
		// 			expNetworks: []blockchain.EVMNetwork{forkedNetwork},
		// 		},
		// 		{
		// 			name: "existing network and new network together in one config",
		// 			configOverrideTOML: `
		// selected_networks = ["new_network","arbitrum_goerli", "optimism_goerli"]

		// [EVMNetworks.new_network]
		// evm_name = "new_test_network"
		// evm_chain_id = 100009
		// evm_urls = ["ws://localhost:8546"]
		// evm_http_urls = ["http://localhost:8545"]
		// client_implementation = "Ethereum"
		// evm_keys = ["ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"]
		// evm_simulated = true
		// evm_chainlink_transaction_limit = 5000
		// evm_minimum_confirmations = 1
		// evm_gas_estimation_buffer = 10000
		// evm_supports_eip1559 = true
		// evm_default_gas_limit = 6000000

		// [RpcHttpUrls]
		// arbitrum_goerli = ["https://devnet-1.mt/ABC/rpc/"]
		// optimism_goerli = ["https://devnet-1.mt/ABC/rpc/"]

		// [RpcWsUrls]
		// arbitrum_goerli = ["wss://devnet-1.mt/ABC/rpc/"]
		// optimism_goerli = ["wss://devnet-1.mt/ABC/rpc/"]

		// [WalletKeys]
		// arbitrum_goerli = ["1810868fc221b9f50b5b3e0186d8a5f343f892e51ce12a9e818f936ec0b651ed"]
		// optimism_goerli = ["1810868fc221b9f50b5b3e0186d8a5f343f892e51ce12a9e818f936ec0b651ed"]
		// 		`,
		// 			expNetworks: []blockchain.EVMNetwork{
		// 				newNetwork, ArbitrumGoerli, OptimismGoerli,
		// 			},
		// 		},
		// 		{
		// 			name: "new network with empty chain id",
		// 			configOverrideTOML: `
		// selected_networks = ["new_network"]

		// [EVMNetworks.new_network]
		// evm_name = "new_test_network"
		// evm_urls = ["ws://localhost:8546"]
		// evm_http_urls = ["http://localhost:8545"]
		// evm_keys = ["ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"]
		// evm_simulated = true
		// evm_chainlink_transaction_limit = 5000
		// client_implementation = "Ethereum"
		// evm_minimum_confirmations = 1
		// evm_gas_estimation_buffer = 10000
		// evm_supports_eip1559 = true
		// evm_default_gas_limit = 6000000
		// 		`,
		// 			isNetworkConfigError: true,
		// 		},
		// 		{
		// 			name: "new network with empty client implementation",
		// 			configOverrideTOML: `
		// selected_networks = ["new_network"]

		// [EVMNetworks.new_network]
		// evm_name = "new_test_network"
		// evm_chain_id = 100009
		// evm_urls = ["ws://localhost:8546"]
		// evm_http_urls = ["http://localhost:8545"]
		// evm_keys = ["ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"]
		// evm_simulated = true
		// evm_chainlink_transaction_limit = 5000
		// evm_minimum_confirmations = 1
		// evm_gas_estimation_buffer = 10000
		// evm_supports_eip1559 = true
		// evm_default_gas_limit = 6000000
		// 		`,
		// 			isNetworkConfigError: true,
		// 		},
		// 		{
		// 			name: "new network without rpc urls",
		// 			configOverrideTOML: `
		// selected_networks = ["new_network"]

		// [EVMNetworks.new_network]
		// evm_name = "new_test_network"
		// evm_chain_id = 100009
		// evm_keys = ["ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"]
		// evm_simulated = true
		// evm_chainlink_transaction_limit = 5000
		// evm_minimum_confirmations = 1
		// evm_gas_estimation_buffer = 10000
		// client_implementation = "Ethereum"
		// evm_supports_eip1559 = true
		// evm_default_gas_limit = 6000000
		// `,
		// 			isNetworkConfigError: true,
		// 		},
		{
			name: "new network with rpc urls and wallet keys both in EVMNetworks and Rpc<Http/Ws>Urls and WalletKeys",
			configOverrideTOML: `
[Network]
selected_networks = ["new_network"]

[Network.EVMNetworks.new_network]
evm_name = "new_test_network"
evm_chain_id = 100009
evm_urls = ["ws://localhost:8546"]
evm_http_urls = ["http://localhost:8545"]
evm_keys = ["ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"]
evm_simulated = true
evm_chainlink_transaction_limit = 5000
evm_minimum_confirmations = 1
evm_gas_estimation_buffer = 10000
client_implementation = "Ethereum"
evm_supports_eip1559 = true
evm_default_gas_limit = 6000000`,
			setupFunc: func() {
				os.Setenv("E2E_TEST_NEW_NETWORK_RPC_HTTP_URL", "http://localhost:iamnotvalid")
				os.Setenv("E2E_TEST_NEW_NETWORK_RPC_WS_URL", "ws://localhost:iamnotvalid")
				os.Setenv("E2E_TEST_NEW_NETWORK_WALLET_KEY", "something random")
			},
			cleanupFunc: func() {
				os.Unsetenv("E2E_TEST_NEW_NETWORK_RPC_HTTP_URL")
				os.Unsetenv("E2E_TEST_NEW_NETWORK_RPC_WS_URL")
				os.Unsetenv("E2E_TEST_NEW_NETWORK_WALLET_KEY")
			},
			expNetworks: []blockchain.EVMNetwork{newNetwork},
		},
		{
			name: "new network with rpc urls and wallet keys in EVMNetworks",
			configOverrideTOML: `
[Network]
selected_networks = ["new_network"]

[Network.EVMNetworks.new_network]
evm_name = "new_test_network"
evm_chain_id = 100009
evm_urls = ["ws://localhost:8546"]
evm_http_urls = ["http://localhost:8545"]
evm_keys = ["ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"]
evm_simulated = true
evm_chainlink_transaction_limit = 5000
evm_minimum_confirmations = 1
evm_gas_estimation_buffer = 10000
client_implementation = "Ethereum"
evm_supports_eip1559 = true
evm_default_gas_limit = 6000000`,
			expNetworks: []blockchain.EVMNetwork{newNetwork},
		},
		{
			name: "new network with rpc urls in EVMNetworks and wallet keys in WalletKeys NetworkConfig",
			configOverrideTOML: `
[Network]
selected_networks = ["new_network"]

[Network.EVMNetworks.new_network]
evm_name = "new_test_network"
evm_chain_id = 100009
evm_urls = ["ws://localhost:8546"]
evm_http_urls = ["http://localhost:8545"]
evm_simulated = true
evm_chainlink_transaction_limit = 5000
evm_minimum_confirmations = 1
evm_gas_estimation_buffer = 10000
client_implementation = "Ethereum"
evm_supports_eip1559 = true
evm_default_gas_limit = 6000000
`,
			setupFunc: func() {
				os.Setenv("E2E_TEST_NEW_NETWORK_WALLET_KEY", "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
			},
			cleanupFunc: func() {
				os.Unsetenv("E2E_TEST_NEW_NETWORK_WALLET_KEY")
			},
			expNetworks: []blockchain.EVMNetwork{newNetwork},
		},
		{
			name: "new network with rpc urls and wallet keys in NetworkConfig",
			configOverrideTOML: `
[Network]
selected_networks = ["new_network"]

[Network.EVMNetworks.new_network]
evm_name = "new_test_network"
evm_chain_id = 100009
evm_simulated = true
evm_chainlink_transaction_limit = 5000
evm_minimum_confirmations = 1
evm_gas_estimation_buffer = 10000
client_implementation = "Ethereum"
evm_supports_eip1559 = true
evm_default_gas_limit = 6000000
`,
			setupFunc: func() {
				os.Setenv("E2E_TEST_NEW_NETWORK_RPC_HTTP_URL", newNetwork.HTTPURLs[0])
				os.Setenv("E2E_TEST_NEW_NETWORK_RPC_WS_URL", newNetwork.URLs[0])
				os.Setenv("E2E_TEST_NEW_NETWORK_WALLET_KEY", newNetwork.PrivateKeys[0])
			},
			cleanupFunc: func() {
				os.Unsetenv("E2E_TEST_NEW_NETWORK_RPC_HTTP_URL")
				os.Unsetenv("E2E_TEST_NEW_NETWORK_RPC_WS_URL")
				os.Unsetenv("E2E_TEST_NEW_NETWORK_WALLET_KEY")
			},
			expNetworks: []blockchain.EVMNetwork{newNetwork},
		},
	}
	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupFunc != nil {
				tc.setupFunc()
			}
			if tc.cleanupFunc != nil {
				defer tc.cleanupFunc()
			}

			// Read from test config override
			cfg := &config.TestConfig{}
			encoded := base64.StdEncoding.EncodeToString([]byte(tc.configOverrideTOML))
			decoded, err := base64.StdEncoding.DecodeString(encoded)
			require.NoError(t, err, "error decoding base64 config")
			err = toml.Unmarshal(decoded, &cfg)
			require.NoError(t, err, "error unmarshalling config")

			// Read from config env vars (test secrets)
			err = cfg.ReadFromEnvVar()
			require.NoError(t, err, "error reading from env var")

			cfg.Network.UpperCaseNetworkNames()
			cfg.Network.OverrideURLsAndKeysFromEVMNetwork()

			err = cfg.Network.Validate()
			if tc.isNetworkConfigError {
				require.Error(t, err, "expected network config error")
				return
			}
			require.NoError(t, err, "error validating network config")
			actualNets, err := SetNetworks(*cfg.Network)
			if tc.isEVMNetworkError {
				t.Log(err)
				require.Error(t, err, "expected evmNetwork set up error")
				return
			}
			require.NoError(t, err, "unexpected error")
			require.Equal(t, tc.expNetworks, actualNets, "unexpected networks")
		})
	}
}
