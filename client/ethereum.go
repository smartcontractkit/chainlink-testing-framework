package client

import (
	"context"
	"encoding/hex"
	"errors"
	"integrations-framework/contracts"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
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
func (e *EthereumClient) SendRawTransaction(fromWallet BlockchainWallet, rawTxDataHex string) (string, error) {
	rawTxData, err := hex.DecodeString(rawTxDataHex)
	if err != nil {
		return "", err
	}

	// Marshal raw data into a transaction
	transaction := new(types.Transaction)
	rlp.DecodeBytes(rawTxData, &transaction)

	privateKey, _ := crypto.HexToECDSA(fromWallet.PrivateKey())

	signedTransaction, err := types.SignTx(transaction, types.NewEIP2930Signer(e.Network.ChainID()), privateKey)
	if err != nil {
		return "", err
	}

	err = e.Client.SendTransaction(context.Background(), signedTransaction)
	if err != nil {
		return "", err
	}

	e.waitForTransaction(signedTransaction.Hash())

	return signedTransaction.Hash().Hex(), err
}

// SendNativeTransaction sends a specified amount of WEI from a selected wallet to an address, and blocks until the
// transaction completes
func (e *EthereumClient) SendNativeTransaction(
	fromWallet BlockchainWallet, toHexAddress string, amount *big.Int) (string, error) {

	gasPrice, err := e.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", err
	}

	nonce, err := e.Client.PendingNonceAt(context.Background(), common.HexToAddress(fromWallet.Address()))
	if err != nil {
		return "", err
	}

	privateKey, _ := crypto.HexToECDSA(fromWallet.PrivateKey())

	unsignedTransaction :=
		types.NewTransaction(nonce, common.HexToAddress(toHexAddress), amount,
			e.Network.Config().TransactionLimit, gasPrice, nil)

	signedTransaction, err := types.SignTx(unsignedTransaction, types.NewEIP2930Signer(e.Network.ChainID()), privateKey)
	if err != nil {
		return "", err
	}

	err = e.Client.SendTransaction(context.Background(), signedTransaction)
	if err != nil {
		return "", err
	}

	e.waitForTransaction(signedTransaction.Hash())

	return signedTransaction.Hash().Hex(), err
}

// SendLinkTransaction sends a specified amount of LINK from a wallet to a public address
func (e *EthereumClient) SendLinkTransaction(
	fromWallet BlockchainWallet, toHexAddress string, amount *big.Int) (string, error) {

	linkTokenAddress := common.HexToAddress(e.Network.Config().LinkTokenAddress)
	toAddress := common.HexToAddress(toHexAddress)
	gasPrice, err := e.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", err
	}

	nonce, err := e.Client.PendingNonceAt(context.Background(), common.HexToAddress(fromWallet.Address()))
	if err != nil {
		return "", err
	}

	privateKey, _ := crypto.HexToECDSA(fromWallet.PrivateKey())

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

	signedTransaction, err := types.SignTx(unsignedTransaction, types.NewEIP2930Signer(e.Network.ChainID()), privateKey)
	if err != nil {
		return "", err
	}

	err = e.Client.SendTransaction(context.Background(), signedTransaction)
	if err != nil {
		return "", err
	}

	e.waitForTransaction(signedTransaction.Hash())

	return signedTransaction.Hash().Hex(), err
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
	gasPrice, err := e.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}

	nonce, err := e.Client.PendingNonceAt(context.Background(), common.HexToAddress(wallet.Address()))
	if err != nil {
		return err
	}

	privateKey, _ := crypto.HexToECDSA(wallet.PrivateKey())
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

// Keep checking until the transaction is no longer pending, or if there is an error
func (e *EthereumClient) waitForTransaction(txHash common.Hash) (bool, error) {
	_, isPending, err := e.Client.TransactionByHash(context.Background(), txHash)
	done := 0
	for isPending {
		if err != nil {
			break
		}
		time.Sleep(1 * time.Second)
		_, isPending, err = e.Client.TransactionByHash(context.Background(), txHash)
		done++
	}
	return isPending, err
}
