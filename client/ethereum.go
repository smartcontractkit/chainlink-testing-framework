package client

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"integrations-framework/contracts"
	"integrations-framework/contracts/ethereum"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
)

// EthereumClient wraps the client and the BlockChain network to interact with an EVM based Blockchain
type EthereumClient struct {
	Client  *ethclient.Client
	Network BlockchainNetwork
}

// NewEthereumClient returns an instantiated instance of the Ethereum client that has connected to the server
func NewEthereumClient(network BlockchainNetwork) (*EthereumClient, error) {
	cl, err := ethclient.Dial(network.URL())
	if err != nil {
		return nil, err
	}

	return &EthereumClient{
		Client:  cl,
		Network: network,
	}, nil
}

// SendTransaction sends a specified amount of WEI from a selected wallet to an address, and blocks until the
// transaction completes
func (e *EthereumClient) SendTransaction(
	fromWallet BlockchainWallet, toHexAddress string, amount int64) (string, error) {

	gasPrice, nonce, pk, err := e.getEthTransactionBasics(fromWallet)
	if err != nil {
		return "", err
	}
	privateKey, err := crypto.HexToECDSA(pk)
	if err != nil {
		return "", fmt.Errorf("invalid private key: %v", err)
	}

	unsignedTransaction :=
		types.NewTransaction(nonce.Uint64(), common.HexToAddress(toHexAddress), big.NewInt(amount),
			e.Network.Config().TransactionLimit, gasPrice, nil)

	txHash, err := e.signAndSendTransaction(unsignedTransaction, privateKey)
	if err != nil {
		return "", err
	}

	err = e.waitForTransaction(txHash)
	return txHash.Hex(), err
}

// DeployStorageContract deploys a vanilla storage contract that is a kv store
func (e *EthereumClient) DeployStorageContract(fromWallet, fundingWallet BlockchainWallet) (contracts.Storage, error) {
	opts, err := e.getTransactionOpts(fromWallet, big.NewInt(0))
	if err != nil {
		return nil, err
	}

	// Deploy contract
	contractAddress, transaction, storageInstance, err := ethereum.DeployStore(opts, e.Client)
	if err != nil {
		return nil, err
	}

	log.Info().Str("Contract address", contractAddress.Hex()).Msg("Deployed storage contract")
	err = e.waitForTransaction(transaction.Hash())
	if err != nil {
		return nil, err
	}

	return contracts.NewEthereumStorage(e, storageInstance, fromWallet), err
}

// Returns the suggested gas price, nonce, private key, and any errors encountered
func (e *EthereumClient) getEthTransactionBasics(wallet BlockchainWallet) (*big.Int, *big.Int, string, error) {
	gasPrice, err := e.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, nil, "", err
	}

	nonce, err := e.Client.PendingNonceAt(context.Background(), common.HexToAddress(wallet.Address()))
	if err != nil {
		return nil, nil, "", err
	}

	return gasPrice, new(big.Int).SetUint64(nonce), wallet.PrivateKey(), err
}

// Helper function to sign and send any ethereum transaction
func (e *EthereumClient) signAndSendTransaction(
	unsignedTransaction *types.Transaction, privateKey *ecdsa.PrivateKey) (common.Hash, error) {

	signedTransaction, err := types.SignTx(unsignedTransaction, types.NewEIP2930Signer(e.Network.ChainID()), privateKey)
	if err != nil {
		return signedTransaction.Hash(), err
	}

	err = e.Client.SendTransaction(context.Background(), signedTransaction)
	if err != nil {
		return signedTransaction.Hash(), err
	}
	log.Info().Str("TX Hash", signedTransaction.Hash().Hex()).Msg("Sending transaction")

	return signedTransaction.Hash(), err
}

// Helper function that waits for a specified transaction to clear
func (e *EthereumClient) waitForTransaction(transactionHash common.Hash) error {
	headerChannel := make(chan *types.Header)
	subscription, err := e.Client.SubscribeNewHead(context.Background(), headerChannel)
	defer subscription.Unsubscribe()
	if err != nil {
		return err
	}

	// Hardhat is a specific case due to instant block mining
	if e.Network.ID() == EthereumHardhatID {
		for {
			_, isPending, err := e.Client.TransactionByHash(context.Background(), transactionHash)
			if err != nil {
				return err
			}
			if !isPending {
				return err
			}
		}
	}

	// Wait for new block to show in subscription
	for {
		select {
		case err := <-subscription.Err():
			return err
		case header := <-headerChannel:
			// Get latest block
			block, err := e.Client.BlockByHash(context.Background(), header.Hash())
			if err != nil {
				return err
			}
			log.Info().Str("Block Hash", block.Hash().Hex()).Msg("New block mined")
			// Look through it for our transaction
			_, isPending, err := e.Client.TransactionByHash(context.Background(), transactionHash)
			if err != nil {
				return err
			}
			if !isPending {
				return err
			}
		}
	}
}

// Builds the default TransactOpts object used for various eth transaction types
func (e *EthereumClient) getTransactionOpts(fromWallet BlockchainWallet, value *big.Int) (*bind.TransactOpts, error) {
	gasPrice, nonce, pk, err := e.getEthTransactionBasics(fromWallet)
	if err != nil {
		return nil, err
	}
	privateKey, err := crypto.HexToECDSA(pk)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %v", err)
	}

	opts, err := bind.NewKeyedTransactorWithChainID(privateKey, e.Network.ChainID())
	if err != nil {
		return nil, err
	}

	opts.Nonce = nonce
	opts.Value = value                                  // in wei
	opts.GasLimit = e.Network.Config().TransactionLimit // in units
	opts.GasPrice = gasPrice

	return opts, err
}
