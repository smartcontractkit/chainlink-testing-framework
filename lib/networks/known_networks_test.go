// nolint
package networks

import (
	"encoding/base64"
	"os"
	"testing"

	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/config"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"
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
	httpOnlyNetwork := newNetwork
	httpOnlyNetwork.URLs = nil
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
		{
			name: "case insensitive network key to EVMNetworks",
			configOverrideTOML: `
[Network]
selected_networks = ["NEW_NETWORK"]

[Network.EVMNetworks.new_Network]
evm_name = "new_test_network"
evm_chain_id = 100009
evm_urls = ["ws://localhost:8546"]
evm_http_urls = ["http://localhost:8545"]
client_implementation = "Ethereum"
evm_keys = ["ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"]
evm_simulated = true
evm_chainlink_transaction_limit = 5000
evm_minimum_confirmations = 1
evm_gas_estimation_buffer = 10000
evm_supports_eip1559 = true
evm_default_gas_limit = 6000000
		`,
			expNetworks: []blockchain.EVMNetwork{newNetwork},
		},
		{
			name: "case insensitive network key of fork config",
			configOverrideTOML: `
[Network]
selected_networks = ["KROMA_SEPOLIA"]

[Network.AnvilConfigs.kroma_SEPOLIA]
url = "ws://localhost:8546"
block_number = 100
		`,
			expNetworks: []blockchain.EVMNetwork{KromaSepolia},
		},
		{
			name: "forked network for existing network",
			configOverrideTOML: `
[Network]
selected_networks = ["KROMA_SEPOLIA"]

[Network.AnvilConfigs.KROMA_SEPOLIA]
url = "ws://localhost:8546"
block_number = 100
		`,
			expNetworks: []blockchain.EVMNetwork{KromaSepolia},
		},
		{
			name: "existing network and new network together in one config",
			configOverrideTOML: `
[Network]
selected_networks = ["new_network","arbitrum_goerli", "optimism_goerli"]

[Network.EVMNetworks.new_network]
evm_name = "new_test_network"
evm_chain_id = 100009
evm_urls = ["ws://localhost:8546"]
evm_http_urls = ["http://localhost:8545"]
client_implementation = "Ethereum"
evm_keys = ["ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"]
evm_simulated = true
evm_chainlink_transaction_limit = 5000
evm_minimum_confirmations = 1
evm_gas_estimation_buffer = 10000
evm_supports_eip1559 = true
evm_default_gas_limit = 6000000
				`,
			setupFunc: func() {
				os.Setenv("E2E_TEST_NEW_NETWORK_RPC_HTTP_URL", "https://devnet-1.mt/ABC/rpc/")
				os.Setenv("E2E_TEST_NEW_NETWORK_RPC_WS_URL", "wss://devnet-1.mt/ABC/rpc/")
				os.Setenv("E2E_TEST_NEW_NETWORK_WALLET_KEY", "1810868fc221b9f50b5b3e0186d8a5f343f892e51ce12a9e818f936ec0b651ed")
				os.Setenv("E2E_TEST_ARBITRUM_GOERLI_RPC_HTTP_URL", "https://devnet-1.mt/ABC/rpc/")
				os.Setenv("E2E_TEST_ARBITRUM_GOERLI_RPC_WS_URL", "wss://devnet-1.mt/ABC/rpc/")
				os.Setenv("E2E_TEST_ARBITRUM_GOERLI_WALLET_KEY", "1810868fc221b9f50b5b3e0186d8a5f343f892e51ce12a9e818f936ec0b651ed")
				os.Setenv("E2E_TEST_OPTIMISM_GOERLI_RPC_HTTP_URL", "https://devnet-1.mt/ABC/rpc/")
				os.Setenv("E2E_TEST_OPTIMISM_GOERLI_RPC_WS_URL", "wss://devnet-1.mt/ABC/rpc/")
				os.Setenv("E2E_TEST_OPTIMISM_GOERLI_WALLET_KEY", "1810868fc221b9f50b5b3e0186d8a5f343f892e51ce12a9e818f936ec0b651ed")
			},
			cleanupFunc: func() {
				os.Unsetenv("E2E_TEST_NEW_NETWORK_RPC_HTTP_URL")
				os.Unsetenv("E2E_TEST_NEW_NETWORK_RPC_WS_URL")
				os.Unsetenv("E2E_TEST_NEW_NETWORK_WALLET_KEY")
				os.Unsetenv("E2E_TEST_ARBITRUM_GOERLI_RPC_HTTP_URL")
				os.Unsetenv("E2E_TEST_ARBITRUM_GOERLI_RPC_WS_URL")
				os.Unsetenv("E2E_TEST_ARBITRUM_GOERLI_WALLET_KEY")
				os.Unsetenv("E2E_TEST_OPTIMISM_GOERLI_RPC_HTTP_URL")
				os.Unsetenv("E2E_TEST_OPTIMISM_GOERLI_RPC_WS_URL")
				os.Unsetenv("E2E_TEST_OPTIMISM_GOERLI_WALLET_KEY")
			},
			expNetworks: []blockchain.EVMNetwork{
				newNetwork, ArbitrumGoerli, OptimismGoerli,
			},
		},
		{
			name: "new network with empty chain id",
			configOverrideTOML: `
[Network]
selected_networks = ["new_network"]

[Network.EVMNetworks.new_network]
evm_name = "new_test_network"
evm_urls = ["ws://localhost:8546"]
evm_http_urls = ["http://localhost:8545"]
evm_keys = ["ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"]
evm_simulated = true
evm_chainlink_transaction_limit = 5000
client_implementation = "Ethereum"
evm_minimum_confirmations = 1
evm_gas_estimation_buffer = 10000
evm_supports_eip1559 = true
evm_default_gas_limit = 6000000`,
			isNetworkConfigError: true,
		},
		{
			name: "new network with empty client implementation",
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
evm_supports_eip1559 = true
evm_default_gas_limit = 6000000`,
			isNetworkConfigError: true,
		},
		{
			name: "new network without rpc urls",
			configOverrideTOML: `
[Network]
selected_networks = ["new_network"]

[Network.EVMNetworks.new_network]
evm_name = "new_test_network"
evm_chain_id = 100009
evm_keys = ["ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"]
evm_simulated = true
evm_chainlink_transaction_limit = 5000
evm_minimum_confirmations = 1
evm_gas_estimation_buffer = 10000
client_implementation = "Ethereum"
evm_supports_eip1559 = true
evm_default_gas_limit = 6000000`,
			isNetworkConfigError: true,
		},
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
		{
			name: "missing wallet key for new network",
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
			},
			cleanupFunc: func() {
				os.Unsetenv("E2E_TEST_NEW_NETWORK_RPC_HTTP_URL")
				os.Unsetenv("E2E_TEST_NEW_NETWORK_RPC_WS_URL")
			},
			isNetworkConfigError: true,
		},
		{
			name: "missing rpc http urls for new network",
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
				os.Setenv("E2E_TEST_NEW_NETWORK_WALLET_KEY", newNetwork.PrivateKeys[0])
				os.Setenv("E2E_TEST_NEW_NETWORK_RPC_WS_URL", newNetwork.URLs[0])
			},
			cleanupFunc: func() {
				os.Unsetenv("E2E_TEST_NEW_NETWORK_WALLET_KEY")
				os.Unsetenv("E2E_TEST_NEW_NETWORK_RPC_WS_URL")
			},
			isNetworkConfigError: true,
		},
		{
			name: "missing rpc ws urls for new network",
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
				os.Setenv("E2E_TEST_NEW_NETWORK_WALLET_KEY", newNetwork.PrivateKeys[0])
				os.Setenv("E2E_TEST_NEW_NETWORK_RPC_HTTP_URL", newNetwork.HTTPURLs[0])
			},
			cleanupFunc: func() {
				os.Unsetenv("E2E_TEST_NEW_NETWORK_WALLET_KEY")
				os.Unsetenv("E2E_TEST_NEW_NETWORK_RPC_HTTP_URL")
			},
			expNetworks: []blockchain.EVMNetwork{httpOnlyNetwork},
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

func TestSetNetworks(t *testing.T) {
	// Helper to create a simple NetworkConfig for testing
	createNetworkConfig := func(selectedNetworks []string, evmNetworks map[string]*blockchain.EVMNetwork,
		rpcHttpUrls map[string][]string, rpcWsUrls map[string][]string, walletKeys map[string][]string,
		anvilConfigs map[string]*config.AnvilConfig) config.NetworkConfig {

		return config.NetworkConfig{
			SelectedNetworks: selectedNetworks,
			EVMNetworks:      evmNetworks,
			RpcHttpUrls:      rpcHttpUrls,
			RpcWsUrls:        rpcWsUrls,
			WalletKeys:       walletKeys,
			AnvilConfigs:     anvilConfigs,
		}
	}

	// Define test cases
	testCases := []struct {
		name             string
		networkCfg       config.NetworkConfig
		expectedError    bool
		expectedErrorMsg string
		expectedCount    int // Expected number of valid networks returned
	}{
		{
			name: "Basic network setup with both HTTP and WS URLs",
			networkCfg: createNetworkConfig(
				[]string{"ETHEREUM_MAINNET"},
				map[string]*blockchain.EVMNetwork{
					"ETHEREUM_MAINNET": {Name: "ETHEREUM_MAINNET", ChainID: 1},
				},
				map[string][]string{"ETHEREUM_MAINNET": {"http://localhost:8545"}},
				map[string][]string{"ETHEREUM_MAINNET": {"ws://localhost:8546"}},
				map[string][]string{"ETHEREUM_MAINNET": {"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"}},
				nil,
			),
			expectedError: false,
			expectedCount: 1,
		},
		{
			name: "Simulated network bypasses validations",
			networkCfg: createNetworkConfig(
				[]string{"SIMULATED"},
				map[string]*blockchain.EVMNetwork{
					"SIMULATED": {Name: "Simulated", Simulated: true, ChainID: 1337},
				},
				nil, nil, nil, nil,
			),
			expectedError: false,
			expectedCount: 1,
		},
		{
			name: "Forked network skips RPC and wallet validation",
			networkCfg: createNetworkConfig(
				[]string{"ETHEREUM_MAINNET"},
				nil,
				nil,
				nil,
				nil,
				map[string]*config.AnvilConfig{
					"ETHEREUM_MAINNET": {URL: ptr.Ptr("http://localhost:8545")},
				},
			),
			expectedError: false,
			expectedCount: 1,
		},
		{
			name: "Only HTTP URLs, valid setup",
			networkCfg: createNetworkConfig(
				[]string{"ETHEREUM_MAINNET"},
				nil,
				map[string][]string{"ETHEREUM_MAINNET": {"http://localhost:8545"}},
				nil,
				map[string][]string{"ETHEREUM_MAINNET": {"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"}},
				nil,
			),
			expectedError: false,
			expectedCount: 1,
		},
		{
			name: "Only WS URLs without HTTP, invalid setup",
			networkCfg: createNetworkConfig(
				[]string{"MAINNET"},
				nil,
				nil,
				map[string][]string{"MAINNET": {"ws://localhost:8546"}},
				map[string][]string{"MAINNET": {"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"}},
				nil,
			),
			expectedError:    true,
			expectedErrorMsg: "WS RPC endpoint for MAINNET network is set without an HTTP endpoint",
		},
		{
			name: "Multiple networks with valid configurations",
			networkCfg: createNetworkConfig(
				[]string{"ETHEREUM_MAINNET", "OPTIMISM_MAINNET"},
				map[string]*blockchain.EVMNetwork{
					"ETHEREUM_MAINNET": {Name: "ETHEREUM_MAINNET", ChainID: 1},
					"OPTIMISM_MAINNET": {Name: "OPTIMISM_MAINNET", ChainID: 10},
				},
				map[string][]string{
					"ETHEREUM_MAINNET": {"http://localhost:8545"},
					"OPTIMISM_MAINNET": {"http://localhost:8547"},
				},
				map[string][]string{
					"ETHEREUM_MAINNET": {"ws://localhost:8546"},
					"OPTIMISM_MAINNET": {"ws://localhost:8548"},
				},
				map[string][]string{
					"ETHEREUM_MAINNET": {"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"},
					"OPTIMISM_MAINNET": {"59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d"},
				},
				nil,
			),
			expectedError: false,
			expectedCount: 2,
		},
		{
			name: "Missing wallet keys, invalid setup",
			networkCfg: createNetworkConfig(
				[]string{"ETHEREUM_MAINNET"},
				nil,
				map[string][]string{"ETHEREUM_MAINNET": {"http://localhost:8545"}},
				nil,
				nil,
				nil,
			),
			expectedError:    true,
			expectedErrorMsg: "no wallet keys found in config for 'ETHEREUM_MAINNET' network",
		},
		{
			name: "Unknown network, invalid setup skips validation",
			networkCfg: createNetworkConfig(
				[]string{"UNKNOWN_NETWORK"},
				nil,
				nil,
				nil,
				nil,
				nil,
			),
			expectedError:    true,
			expectedErrorMsg: "at least one HTTP RPC endpoint for UNKNOWN_NETWORK network must be set",
		},
		{
			name: "Valid known network from MappedNetworks",
			networkCfg: createNetworkConfig(
				[]string{"SIMULATED_1"},
				nil,
				nil,
				nil,
				nil,
				nil,
			),
			expectedError: false,
			expectedCount: 1,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			nets, err := SetNetworks(tc.networkCfg)

			// Check for expected error
			if tc.expectedError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErrorMsg)
				return
			}
			require.NoError(t, err)

			// Check the expected count of networks returned
			require.Equal(t, tc.expectedCount, len(nets))
		})
	}
}
