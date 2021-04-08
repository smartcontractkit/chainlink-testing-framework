package client

import (
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/core"
)

// Tests ethereum contract deployment on a simulated blockchain
func TestEthereumClient_DeployStorageContract(t *testing.T) {
	sim := backends.NewSimulatedBackend(core.DefaultGenesisBlock().Alloc, 1000)
	defer sim.Close()
	client := NewSimulatedEthereumClient(sim)

	client.DeployStorageContract()

	sim.Commit()
}
