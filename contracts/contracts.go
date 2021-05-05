package contracts

import (
	"context"
	"integrations-framework/client"
	"integrations-framework/contracts/ethereum"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

// Storage represents generalized interactions with the stoarge contract
type Storage interface {
	Get(ctxt context.Context) (*big.Int, error)
	Set(context.Context, *big.Int) error
}

// EthereumStorage acts as a conduit for the ethereum version of the storage contract
type EthereumStorage struct {
	client       *client.EthereumClient
	store        *ethereum.Store
	callerWallet client.BlockchainWallet
}

// NewEthereumStorage creates a new instance of the storage contract for ethereum chains
func NewEthereumStorage(client *client.EthereumClient, store *ethereum.Store, callerWallet client.BlockchainWallet) Storage {
	return &EthereumStorage{
		client:       client,
		store:        store,
		callerWallet: callerWallet,
	}
}

// DeployStorageContract deploys a vanilla storage contract that is a value store
func DeployStorageContract(ethClient *client.EthereumClient, fromWallet client.BlockchainWallet) (Storage, error) {
	_, _, instance, err := ethClient.DeployContract(fromWallet, func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return ethereum.DeployStore(auth, backend)
	})
	if err != nil {
		return nil, err
	}
	return NewEthereumStorage(ethClient, instance.(*ethereum.Store), fromWallet), nil
}

// Set sets a value in the storage contract
func (e *EthereumStorage) Set(ctxt context.Context, value *big.Int) error {
	opts, err := e.client.GetTransactionOpts(e.callerWallet, big.NewInt(0))
	if err != nil {
		return err
	}

	transaction, err := e.store.Set(opts, value)
	if err != nil {
		return err
	}
	return e.client.WaitForTransaction(transaction.Hash())
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
