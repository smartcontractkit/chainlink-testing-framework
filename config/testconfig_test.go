package config

import (
	"errors"
	"os"
	"reflect"
	"testing"

	"github.com/smartcontractkit/chainlink-testing-framework/utils/ptr"
)

func TestReadConfigValuesFromEnvVars(t *testing.T) {
	tests := []struct {
		name           string
		setupFunc      func()
		cleanupFunc    func()
		expectedConfig TestConfig
		expectedError  error
	}{
		{
			name: "All configurations set correctly",
			setupFunc: func() {
				os.Setenv("E2E_TEST_GROUP1_WALLET_KEY_1", "walletValue1")
				os.Setenv("E2E_TEST_GROUP2_RPC_HTTP_URL_1", "httpUrl1")
				os.Setenv("E2E_TEST_GROUP3_RPC_WS_URL_1", "wsUrl1")
				os.Setenv("E2E_TEST_CHAINLINK_IMAGE", "imageValue")
				os.Setenv("E2E_TEST_PYROSCOPE_ENABLED", "true")
			},
			cleanupFunc: func() {
				os.Unsetenv("E2E_TEST_GROUP1_WALLET_KEY_1")
				os.Unsetenv("E2E_TEST_GROUP2_RPC_HTTP_URL_1")
				os.Unsetenv("E2E_TEST_GROUP3_RPC_WS_URL_1")
				os.Unsetenv("E2E_TEST_CHAINLINK_IMAGE")
				os.Unsetenv("E2E_TEST_PYROSCOPE_ENABLED")
			},
			expectedConfig: TestConfig{
				Network: &NetworkConfig{
					WalletKeys:  map[string][]string{"GROUP1": {"walletValue1"}},
					RpcHttpUrls: map[string][]string{"GROUP2": {"httpUrl1"}},
					RpcWsUrls:   map[string][]string{"GROUP3": {"wsUrl1"}},
				},
				Pyroscope:      &PyroscopeConfig{Enabled: ptr.Ptr[bool](true)},
				ChainlinkImage: &ChainlinkImageConfig{Image: newString("imageValue")},
			},
		},
		{
			name: "Environment variables are empty strings",
			setupFunc: func() {
				os.Setenv("E2E_TEST_GROUP1_WALLET_KEY_1", "")
				os.Setenv("E2E_TEST_GROUP2_RPC_HTTP_URL_1", "")
				os.Setenv("E2E_TEST_GROUP3_RPC_WS_URL_1", "")
				os.Setenv("E2E_TEST_CHAINLINK_IMAGE", "")
			},
			cleanupFunc: func() {
				os.Unsetenv("E2E_TEST_GROUP1_WALLET_KEY_1")
				os.Unsetenv("E2E_TEST_GROUP2_RPC_HTTP_URL_1")
				os.Unsetenv("E2E_TEST_GROUP3_RPC_WS_URL_1")
				os.Unsetenv("E2E_TEST_CHAINLINK_IMAGE")
			},
			expectedConfig: TestConfig{},
		},
		{
			name: "Environment variables set with mixed numeric suffixes",
			setupFunc: func() {
				os.Setenv("E2E_TEST_ARBITRUM_SEPOLIA_RPC_HTTP_URL", "url")
				os.Setenv("E2E_TEST_ARBITRUM_SEPOLIA_RPC_HTTP_URL_1", "url1")
				os.Setenv("E2E_TEST_ARBITRUM_SEPOLIA_RPC_WS_URL", "wsurl")
				os.Setenv("E2E_TEST_ARBITRUM_SEPOLIA_RPC_WS_URL_1", "wsurl1")
				os.Setenv("E2E_TEST_OPTIMISM_SEPOLIA_WALLET_KEY", "wallet")
				os.Setenv("E2E_TEST_OPTIMISM_SEPOLIA_WALLET_KEY_1", "wallet1")
				os.Setenv("E2E_TEST_OPTIMISM_SEPOLIA_WALLET_KEY_2", "wallet2")
			},
			cleanupFunc: func() {
				os.Unsetenv("E2E_TEST_ARBITRUM_SEPOLIA_RPC_HTTP_URL_1")
				os.Unsetenv("E2E_TEST_ARBITRUM_SEPOLIA_RPC_WS_URL")
				os.Unsetenv("E2E_TEST_ARBITRUM_SEPOLIA_RPC_WS_URL_1")
				os.Unsetenv("E2E_TEST_OPTIMISM_SEPOLIA_WALLET_KEY")
				os.Unsetenv("E2E_TEST_OPTIMISM_SEPOLIA_WALLET_KEY_1")
				os.Unsetenv("E2E_TEST_ARBITRUM_SEPOLIA_RPC_HTTP_URL_2")
				os.Unsetenv("E2E_TEST_OPTIMISM_SEPOLIA_WALLET_KEY_2")
			},
			expectedConfig: TestConfig{
				Network: &NetworkConfig{
					RpcHttpUrls: map[string][]string{"ARBITRUM_SEPOLIA": {"url", "url1"}},
					RpcWsUrls:   map[string][]string{"ARBITRUM_SEPOLIA": {"wsurl", "wsurl1"}},
					WalletKeys:  map[string][]string{"OPTIMISM_SEPOLIA": {"wallet", "wallet1", "wallet2"}},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			localTt := tt               // Create a local copy of tt
			localTt.setupFunc()         // Setup the test environment
			defer localTt.cleanupFunc() // Ensure cleanup after the test

			c := &TestConfig{}

			// Execute
			err := c.ReadFromEnvVar()

			// Verify error handling
			if !errors.Is(err, localTt.expectedError) {
				t.Errorf("Expected error to be %v, got %v", localTt.expectedError, err)
			}

			// Verify the config
			if !reflect.DeepEqual(c, &localTt.expectedConfig) {
				t.Errorf("Expected config to be %+v, got %+v", localTt.expectedConfig, c)
			}
		})
	}
}

func newString(s string) *string {
	return &s
}
