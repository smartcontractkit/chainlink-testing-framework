package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Tests ethereum contract deployment on a simulated blockchain
func TestEthereumClient_DeployStorageContract(t *testing.T) {
	testCases := []struct {
		name    string
		network BlockchainNetwork
	}{
		{"Ethereum Hardhat", &EthereumHardhat{}},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			client, err := NewBlockchainClient(testCase.network)
			require.Nil(t, err)

			err = client.DeployStorageContract(testCase.network.Wallets().Default())
			assert.Nil(t, err)
		})
	}
}
