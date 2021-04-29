package client

import (
	"context"
	"integrations-framework/contracts/ethereum"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

type Storage interface {
	// FilterItemSet(uint64, uint64, context.Context)
	Items(context.Context, [32]byte) ([32]byte, error)
	// ParseItemSet()
	SetItem(context.Context, [32]byte, [32]byte) error
	Version(context.Context) (string, error)
	// WatchItemSet()
}

type EthereumStorage struct {
	client       *EthereumClient
	storage      *ethereum.Storage
	callerWallet BlockchainWallet
}

func NewEthereumStorage(client *EthereumClient, storage *ethereum.Storage, callerWallet BlockchainWallet) Storage {
	return &EthereumStorage{
		client:       client,
		storage:      storage,
		callerWallet: callerWallet,
	}
}

func (e *EthereumStorage) Items(ctxt context.Context, key [32]byte) ([32]byte, error) {
	opts := &bind.CallOpts{
		Pending: false,
		Context: ctxt,
	}
	return e.storage.Items(opts, key)
}

func (e *EthereumStorage) SetItem(ctxt context.Context, key, value [32]byte) error {
	opts, err := e.client.getTransactionOpts(e.callerWallet, big.NewInt(0))
	if err != nil {
		return err
	}

	transaction, err := e.storage.SetItem(opts, key, value)
	if err != nil {
		return err
	}

	return e.client.waitForTransaction(transaction.Hash())
}

func (e *EthereumStorage) Version(ctxt context.Context) (string, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(e.callerWallet.Address()),
		Pending: false,
		Context: ctxt,
	}
	return e.storage.Version(opts)

}
