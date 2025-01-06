package seth_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	pkg_seth "github.com/smartcontractkit/chainlink-testing-framework/seth"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/seth"
)

// Helper to create EVM networks with different configurations
func createEVMNetwork(name string, simulated bool, urls []string, httpUrls []string, supportsEIP1559 bool, defaultGasLimit uint64) blockchain.EVMNetwork {
	return blockchain.EVMNetwork{
		Name:            name,
		Simulated:       simulated,
		PrivateKeys:     []string{"key1"},
		URLs:            urls,
		HTTPURLs:        httpUrls,
		SupportsEIP1559: supportsEIP1559,
		DefaultGasLimit: defaultGasLimit,
	}
}

// Helper to create Seth config with initial network names
func createSethConfig(networkNames ...string) pkg_seth.Config {
	var networks []*pkg_seth.Network
	for _, name := range networkNames {
		networks = append(networks, &pkg_seth.Network{Name: name})
	}
	return pkg_seth.Config{Networks: networks}
}

func TestMergeSethAndEvmNetworkConfigs(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name             string
		evmNetwork       blockchain.EVMNetwork
		sethConfig       pkg_seth.Config
		expectedError    bool
		expectedErrorMsg string
		expectedURLs     []string
		expectedKeys     []string
		expectedEIP1559  bool
		expectedGasLimit uint64
	}{
		// Anvil tests
		{
			name: "Anvil network with both HTTP and WS URLs",
			evmNetwork: createEVMNetwork(
				"Anvil", true, []string{"ws://localhost:8546"}, []string{"http://localhost:8545"}, true, 21000,
			),
			sethConfig:       createSethConfig("Anvil"),
			expectedURLs:     []string{"ws://localhost:8546"},
			expectedKeys:     []string{"key1"},
			expectedEIP1559:  true,
			expectedGasLimit: 21000,
		},
		{
			name: "Anvil network with only HTTP URLs",
			evmNetwork: createEVMNetwork(
				"Anvil", true, nil, []string{"http://localhost:8545"}, true, 21000,
			),
			sethConfig:       createSethConfig("Anvil"),
			expectedURLs:     []string{"http://localhost:8545"},
			expectedKeys:     []string{"key1"},
			expectedEIP1559:  true,
			expectedGasLimit: 21000,
		},
		{
			name: "Anvil network with only WS URLs",
			evmNetwork: createEVMNetwork(
				"Anvil", true, []string{"ws://localhost:8546"}, nil, true, 21000,
			),
			sethConfig:       createSethConfig("Anvil"),
			expectedURLs:     []string{"ws://localhost:8546"},
			expectedKeys:     []string{"key1"},
			expectedEIP1559:  true,
			expectedGasLimit: 21000,
		},
		// Geth tests
		{
			name: "Geth network with both HTTP and WS URLs",
			evmNetwork: createEVMNetwork(
				"Geth", true, []string{"ws://localhost:8546"}, []string{"http://localhost:8545"}, true, 21000,
			),
			sethConfig:       createSethConfig("Geth"),
			expectedURLs:     []string{"ws://localhost:8546"},
			expectedKeys:     []string{"key1"},
			expectedEIP1559:  true,
			expectedGasLimit: 21000,
		},
		{
			name: "Geth network with only HTTP URLs",
			evmNetwork: createEVMNetwork(
				"Geth", true, nil, []string{"http://localhost:8545"}, true, 21000,
			),
			sethConfig:       createSethConfig("Geth"),
			expectedURLs:     []string{"http://localhost:8545"},
			expectedKeys:     []string{"key1"},
			expectedEIP1559:  true,
			expectedGasLimit: 21000,
		},
		{
			name: "Geth network with only WS URLs",
			evmNetwork: createEVMNetwork(
				"Geth", true, []string{"ws://localhost:8546"}, nil, true, 21000,
			),
			sethConfig:       createSethConfig("Geth"),
			expectedURLs:     []string{"ws://localhost:8546"},
			expectedKeys:     []string{"key1"},
			expectedEIP1559:  true,
			expectedGasLimit: 21000,
		},
		// Live Network tests
		{
			name: "Live_Network with both HTTP and WS URLs",
			evmNetwork: createEVMNetwork(
				"Live_Network", true, []string{"ws://localhost:8546"}, []string{"http://localhost:8545"}, true, 21000,
			),
			sethConfig:       createSethConfig("Live_Network"),
			expectedURLs:     []string{"ws://localhost:8546"},
			expectedKeys:     []string{"key1"},
			expectedEIP1559:  false,
			expectedGasLimit: 0,
		},
		{
			name: "Live_Network with only HTTP URLs",
			evmNetwork: createEVMNetwork(
				"Live_Network", false, nil, []string{"http://live-network:8545"}, true, 21000,
			),
			sethConfig:       createSethConfig("Live_Network"),
			expectedURLs:     []string{"http://live-network:8545"},
			expectedKeys:     []string{"key1"},
			expectedEIP1559:  false,
			expectedGasLimit: 0,
		},
		{
			name: "Live_Network with only WS URLs",
			evmNetwork: createEVMNetwork(
				"Live_Network", true, []string{"ws://localhost:8546"}, nil, true, 21000,
			),
			sethConfig:       createSethConfig("Live_Network"),
			expectedURLs:     []string{"ws://localhost:8546"},
			expectedKeys:     []string{"key1"},
			expectedEIP1559:  false,
			expectedGasLimit: 0,
		},
		// Default fallback network tests
		{
			name: "New network falls back to default config with both HTTP and WS URLs",
			evmNetwork: createEVMNetwork(
				"New_Network", true, []string{"ws://localhost:8546"}, []string{"http://localhost:8545"}, true, 21000,
			),
			sethConfig:       createSethConfig("New_Network"),
			expectedURLs:     []string{"ws://localhost:8546"},
			expectedKeys:     []string{"key1"},
			expectedEIP1559:  false,
			expectedGasLimit: 0,
		},
		{
			name: "New network falls back to default config with only HTTP",
			evmNetwork: createEVMNetwork(
				"New_Network", true, nil, []string{"http://localhost:8545"}, true, 21000,
			),
			sethConfig:       createSethConfig("New_Network"),
			expectedURLs:     []string{"http://localhost:8545"},
			expectedKeys:     []string{"key1"},
			expectedEIP1559:  false,
			expectedGasLimit: 0,
		},
		{
			name: "New network falls back to default config with only WS URLs",
			evmNetwork: createEVMNetwork(
				"New_Network", true, []string{"ws://localhost:8546"}, nil, true, 21000,
			),
			sethConfig:       createSethConfig("New_Network"),
			expectedURLs:     []string{"ws://localhost:8546"},
			expectedKeys:     []string{"key1"},
			expectedEIP1559:  false,
			expectedGasLimit: 0,
		},
		// Unknown network tests
		{
			name: "Unknown network with no matching or default config",
			evmNetwork: createEVMNetwork(
				"UnknownNetwork", false, []string{"ws://unknown:8546"}, nil, false, 0,
			),
			sethConfig:       createSethConfig("Anvil", "Geth"),
			expectedError:    true,
			expectedErrorMsg: "Failed to build network config for chain ID",
		},
	}

	// Run test cases
	for _, tc := range testCases {
		tc := tc //nolint capture range variable
		t.Run(tc.name, func(t *testing.T) {
			mergedConfig, err := seth.MergeSethAndEvmNetworkConfigs(tc.evmNetwork, tc.sethConfig)

			// Check for expected error
			if tc.expectedError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErrorMsg)
				return
			}
			require.NoError(t, err)

			// Validate merged configuration details
			require.Equal(t, tc.expectedURLs, mergedConfig.Network.URLs, "URLs mismatch")
			require.Equal(t, tc.expectedKeys, mergedConfig.Network.PrivateKeys, "PrivateKeys mismatch")
			require.Equal(t, tc.expectedEIP1559, mergedConfig.Network.EIP1559DynamicFees, "EIP1559DynamicFees mismatch")
			require.Equal(t, tc.expectedGasLimit, mergedConfig.Network.GasLimit, "GasLimit mismatch")
		})
	}
}
