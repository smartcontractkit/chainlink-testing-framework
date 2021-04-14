package client

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Tests retrieving wallet values from env variables for ethereum wallets
func TestEthereumWallet_EnvVariable(t *testing.T) {
	testCases := []struct {
		name       string // Name of the test case as well as the env variable to set
		network    BlockchainNetwork
		privateKey string
		address    string
	}{
		{"EthereumHardhat", &EthereumHardhat{},
			"cfff63a9931f8e948f8475795dd015065be59e5cecffeb7c2e2bfa48981d9d24",
			"0x5ff19251b8f8702d485f127d96232301023119f1"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			os.Setenv(testCase.name, testCase.privateKey)
			wallets := testCase.network.Wallets()
			assert.Equal(t, testCase.privateKey, wallets.Default().PrivateKey())
			assert.Equal(t, testCase.address, strings.ToLower(wallets.Default().Address()))
		})
	}
}

// Tests retrieving wallet values from env variables for ethereum wallets
func TestEthereumWallet_ConfigFile(t *testing.T) {
	testCases := []struct {
		name       string
		network    BlockchainNetwork
		privateKey string
		address    string
	}{
		{"EthereumHardhat", &EthereumHardhat{},
			"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
			"0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			wallets := testCase.network.Wallets()
			assert.Equal(t, testCase.privateKey, wallets.Default().PrivateKey())
			assert.Equal(t, testCase.address, strings.ToLower(wallets.Default().Address()))
		})
	}
}
