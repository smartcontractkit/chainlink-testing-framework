package client

import (
	"math/big"
	"net/http"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

const HardhatConnectionString string = "http://localhost:8545"
const HardhatHexAddress string = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
const HardHatChainID int64 = 31337

// Tests ethereum contract deployment on a simulated blockchain
func TestEthereumClient_DeployStorageContract(t *testing.T) {
	checkIsHardhatUp(t)
	address := common.HexToAddress(HardhatHexAddress)

	simulatedClient := NewEthereumClient(HardhatConnectionString, big.NewInt(HardHatChainID), address)

	simulatedClient.DeployStorageContract()
}

// Tests if a local hardhat network is running, fails the test if it does not
func checkIsHardhatUp(t *testing.T) {
	_, isUp := http.Get(HardhatConnectionString)
	if isUp != nil {
		t.Fatal("Hardhat network has not been started!")
	}
}
