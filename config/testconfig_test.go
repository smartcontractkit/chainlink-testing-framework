package config

import (
	"os"
	"reflect"
	"testing"

	"github.com/smartcontractkit/chainlink-testing-framework/utils/ptr"
)

func TestReadConfigValuesFromEnvVars(t *testing.T) {
	// Define multiple test cases
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
				os.Setenv("TEST_CONFIG_GROUP1_WALLET_KEY_1", "walletValue1")
				os.Setenv("TEST_CONFIG_GROUP2_RPC_HTTP_URL_1", "httpUrl1")
				os.Setenv("TEST_CONFIG_GROUP3_RPC_WS_URL_1", "wsUrl1")
				os.Setenv("TEST_CONFIG_CHAINLINK_IMAGE", "imageValue")
				os.Setenv("TEST_CONFIG_PYROSCOPE_ENABLED", "true")
			},
			cleanupFunc: func() {
				os.Unsetenv("TEST_CONFIG_GROUP1_WALLET_KEY_1")
				os.Unsetenv("TEST_CONFIG_GROUP2_RPC_HTTP_URL_1")
				os.Unsetenv("TEST_CONFIG_GROUP3_RPC_WS_URL_1")
				os.Unsetenv("TEST_CONFIG_CHAINLINK_IMAGE")
				os.Unsetenv("TEST_CONFIG_PYROSCOPE_ENABLED")
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
				os.Setenv("TEST_CONFIG_GROUP1_WALLET_KEY_1", "")
				os.Setenv("TEST_CONFIG_GROUP2_RPC_HTTP_URL_1", "")
				os.Setenv("TEST_CONFIG_GROUP3_RPC_WS_URL_1", "")
				os.Setenv("TEST_CONFIG_CHAINLINK_IMAGE", "")
			},
			cleanupFunc: func() {
				os.Unsetenv("TEST_CONFIG_GROUP1_WALLET_KEY_1")
				os.Unsetenv("TEST_CONFIG_GROUP2_RPC_HTTP_URL_1")
				os.Unsetenv("TEST_CONFIG_GROUP3_RPC_WS_URL_1")
				os.Unsetenv("TEST_CONFIG_CHAINLINK_IMAGE")
			},
			expectedConfig: TestConfig{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupFunc()         // Setup the test environment
			defer tt.cleanupFunc() // Ensure cleanup after the test

			c := &TestConfig{}

			// Execute
			err := c.ReadConfigValuesFromEnvVars()

			// Verify error handling
			if err != tt.expectedError {
				t.Errorf("Expected error to be %v, got %v", tt.expectedError, err)
			}

			// Verify the config
			if !reflect.DeepEqual(c, &tt.expectedConfig) {
				t.Errorf("Expected config to be %+v, got %+v", tt.expectedConfig, c)
			}
		})
	}
}

func newString(s string) *string {
	return &s
}
