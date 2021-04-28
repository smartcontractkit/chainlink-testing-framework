package client

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// TransactionOptions contain all the variables needed to transact on a blockchain network.
type TransactionOptions interface {
	// Converts transaction options into a single hex string
	Hex() (string, error)
}

// EthereumTransactionOptions specify all the variables needed to transact on the ethereum blockchain
type EthereumTransactionOptions struct {
	// Should we replace these by just passing in a client so we can estimate these / grab them easier?
	// Or leave as is? More customization this way, but likely more errors
	ChainID  *big.Int
	Nonce    *big.Int
	GasLimit *big.Int
	GasPrice *big.Int
	From     BlockchainWallet
	To       string
	Value    *big.Int
	Data     []byte
}

func (options *EthereumTransactionOptions) Hex() (string, error) {
	// Marshall eth transaction options to transaction object
	toAddress := common.HexToAddress(options.To)
	unsignedTransaction := types.NewTransaction(options.Nonce.Uint64(), toAddress, options.Value,
		options.GasLimit.Uint64(), options.GasPrice, options.Data)

	signedTransaction, err := types.SignTx(unsignedTransaction, types.NewEIP2930Signer(options.ChainID),
		options.From.PrivateKey())

	return signedTransaction.Hash().Hex(), err
}
