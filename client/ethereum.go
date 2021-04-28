package client

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"integrations-framework/contracts"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/sha3"
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

// SendRawTransaction uses a specified wallet and raw hex data to sign and send a raw transaction
func (e *EthereumClient) SendRawTransaction(options TransactionOptions) (string, error) {
	rawHex, err := options.Hex()
	if err != nil {
		return "", err
	}
	rawTxData, err := hex.DecodeString(rawHex)
	if err != nil {
		return "", err
	}

	// Marshal raw data into a transaction
	transaction := new(types.Transaction)
	err = rlp.DecodeBytes(rawTxData, &transaction)
	if err != nil {
		return "", err
	}

	err = e.Client.SendTransaction(context.Background(), transaction)
	if err != nil {
		return "", err
	}

	err = e.waitForTransaction(transaction.Hash())
	return transaction.Hash().Hex(), err
}

// SendNativeTransaction sends a specified amount of WEI from a selected wallet to an address, and blocks until the
// transaction completes
func (e *EthereumClient) SendNativeTransaction(
	fromWallet BlockchainWallet, toHexAddress string, amount *big.Int) (string, error) {

	gasPrice, nonce, privateKey, err := e.getEthTransactionBasics(fromWallet)
	if err != nil {
		return "", err
	}

	unsignedTransaction :=
		types.NewTransaction(nonce, common.HexToAddress(toHexAddress), amount,
			e.Network.Config().TransactionLimit, gasPrice, nil)

	txHash, err := e.signAndSendTransaction(unsignedTransaction, privateKey)

	return txHash.Hex(), err
}

// SendLinkTransaction sends a specified amount of LINK from a wallet to a public address
func (e *EthereumClient) SendLinkTransaction(
	fromWallet BlockchainWallet, toHexAddress string, amount *big.Int) (string, error) {

	linkTokenAddress := common.HexToAddress(e.Network.Config().LinkTokenAddress)
	toAddress := common.HexToAddress(toHexAddress)
	gasPrice, nonce, privateKey, err := e.getEthTransactionBasics(fromWallet)
	if err != nil {
		return "", err
	}

	// Prepare data to transfer LINK token
	transferFnSignature := []byte("transfer(address,uint256)")
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]
	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)

	// Marshall data
	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	unsignedTransaction := types.NewTransaction(nonce, linkTokenAddress, big.NewInt(0),
		e.Network.Config().TransactionLimit, gasPrice, data)

	txHash, err := e.signAndSendTransaction(unsignedTransaction, privateKey)

	return txHash.Hex(), err
}

// GetNativeBalance returns the balance of ETH a public address has in WEI
func (e *EthereumClient) GetNativeBalance(addressHex string) (*big.Int, error) {
	accountAddress := common.HexToAddress(addressHex)
	return e.Client.BalanceAt(context.Background(), accountAddress, nil)
}

// GetLinkBalance returns to balance of LINK a public address has
func (e *EthereumClient) GetLinkBalance(addressHex string) (*big.Int, error) {
	// TODO: Needs LINK token in hardhat
	return nil, errors.New("not implemented yet")
}

// DeployStorageContract deploys a vanilla storage contract that is a kv store
func (e *EthereumClient) DeployStorageContract(wallet BlockchainWallet) error {
	gasPrice, nonce, privateKey, err := e.getEthTransactionBasics(wallet)
	if err != nil {
		return err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, e.Network.ChainID())
	if err != nil {
		return err
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(3)                          // in wei
	auth.GasLimit = e.Network.Config().TransactionLimit // in units
	auth.GasPrice = gasPrice

	_, _, _, err = contracts.DeployStorage(auth, e.Client, "1.0")
	return err
}

// Returns the suggested gas price, nonce, private key, and any errors encountered
func (e *EthereumClient) getEthTransactionBasics(wallet BlockchainWallet) (*big.Int, uint64, *ecdsa.PrivateKey, error) {
	gasPrice, err := e.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, 0, nil, err
	}

	nonce, err := e.Client.PendingNonceAt(context.Background(), common.HexToAddress(wallet.Address()))
	if err != nil {
		return nil, 0, nil, err
	}

	return gasPrice, nonce, wallet.PrivateKey(), err
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
	log.Println("Sending transaction. Hash: ", signedTransaction.Hash().Hex())

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
			log.Println("New block mined. Hash: ", block.Hash().Hex())
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
