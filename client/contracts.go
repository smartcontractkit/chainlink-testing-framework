package client

import (
	"context"
	"integrations-framework/contracts/ethereum"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

type Storage interface {
	Get(ctxt context.Context) (*big.Int, error)
	Set(context.Context, *big.Int) error
}

type EthereumStorage struct {
	client       *EthereumClient
	store        *ethereum.Store
	callerWallet BlockchainWallet
}

// Creates a new instance of the storage contract for ethereum chains
func NewEthereumStorage(client *EthereumClient, store *ethereum.Store, callerWallet BlockchainWallet) Storage {
	return &EthereumStorage{
		client:       client,
		store:        store,
		callerWallet: callerWallet,
	}
}

// Set sets a value in the storage contract
func (e *EthereumStorage) Set(ctxt context.Context, value *big.Int) error {
	opts, err := e.client.getTransactionOpts(e.callerWallet, big.NewInt(0))
	if err != nil {
		return err
	}

	transaction, err := e.store.Set(opts, value)
	if err != nil {
		return err
	}

	return e.client.waitForTransaction(transaction.Hash())
}

// Get retrieves a set value from the storage contract
func (e *EthereumStorage) Get(ctxt context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(e.callerWallet.Address()),
		Pending: true,
		Context: ctxt,
	}
	return e.store.Get(opts)

}
