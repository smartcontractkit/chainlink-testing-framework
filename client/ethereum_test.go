package client

import (
	"math/big"
	"net/http"
	"testing"
)

const HardhatConnectionString string = "http://localhost:8545"
const HardHatChainID int64 = 31337

// Tests ethereum contract deployment on a simulated blockchain
func TestEthereumClient_DeployStorageContract(t *testing.T) {
	checkIsHardhatUp(t)

	simulatedClient := NewEthereumClient(HardhatConnectionString, big.NewInt(HardHatChainID))

	simulatedClient.DeployStorageContract()
}

// Tests if a local hardhat network is running, fails the test if it does not
func checkIsHardhatUp(t *testing.T) {
	_, isUp := http.Get(HardhatConnectionString)
	if isUp != nil {
		t.Fatal("Hardhat network has not been started!")
	}
}
