package client

import (
	"integrations-framework/contracts"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Wallet struct {
}

// Generalized blockchain client for interaction with multiple different blockchains
type BlockchainClient interface {
	SetDefaultWallet(wallet Wallet) // For future reference
	DeployStorageContract() (common.Address, *types.Transaction, *contracts.Storage)
}
