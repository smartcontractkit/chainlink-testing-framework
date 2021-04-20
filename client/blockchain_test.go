package client

import (
	"integrations-framework/config"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Tests retrieving wallet values from env variables for ethereum wallets
func TestEthereumWallet_EnvVariable(t *testing.T) {
	conf, err := config.NewConfig(config.EnvironmentConfig)
	require.Nil(t, err)
	testCases := []struct {
		name          string
		envConfigName string
		network       BlockchainNetwork
		privateKey    string
		address       string
	}{
		{"Ethereum Hardhat", "ethereum_hardhat", NewEthereumHardhat(conf.Networks["etherum_hardhat"]),
			"cfff63a9931f8e948f8475795dd015065be59e5cecffeb7c2e2bfa48981d9d24",
			"0x5ff19251b8f8702d485f127d96232301023119f1"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// Set up
			err := os.Setenv("networks"+testCase.envConfigName+".private_keys", testCase.privateKey)
			require.Nil(t, err)

			wallets, err := testCase.network.Wallets()
			require.Nil(t, err)
			assert.Equal(t, testCase.privateKey, wallets.Default().PrivateKey())
			assert.Equal(t, testCase.address, strings.ToLower(wallets.Default().Address()))

			// Tear down
			err = os.Unsetenv(testCase.name)
			require.Nil(t, err)
		})
	}
}

// Tests retrieving wallet values from env variables for ethereum wallets
func TestEthereumWallet_ConfigFile(t *testing.T) {
	conf, err := config.NewConfig(config.FileConfig)
	require.Nil(t, err)
	testCases := []struct {
		name       string
		network    BlockchainNetwork
		privateKey string
		address    string
	}{
		{"Ethereum Hardhat", NewEthereumHardhat(conf.Networks["ethereum_hardhat"]),
			"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
			"0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			wallets, err := testCase.network.Wallets()
			require.Nil(t, err)
			assert.Equal(t, testCase.privateKey, wallets.Default().PrivateKey())
			assert.Equal(t, testCase.address, strings.ToLower(wallets.Default().Address()))
		})
	}
}

// Tests ethereum contract deployment on a simulated blockchain
func TestEthereumClient_DeployStorageContract(t *testing.T) {
	conf, err := config.NewConfig(config.FileConfig)
	require.Nil(t, err)
	testCases := []struct {
		name    string
		network BlockchainNetwork
	}{
		{"Ethereum Hardhat", NewEthereumHardhat(conf.Networks["etherum_hardhat"])},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			client, err := NewBlockchainClient(testCase.network)
			require.Nil(t, err)

			wallets, err := testCase.network.Wallets()
			require.Nil(t, err)

			err = client.DeployStorageContract(wallets.Default())
			assert.Nil(t, err)
		})
	}
}
