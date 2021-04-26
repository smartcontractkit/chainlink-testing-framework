package client

import (
	"context"
	"errors"
	"integrations-framework/contracts"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
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

// GetLatestBlock retrieves the latest valid block from the EVM based chain
func (e *EthereumClient) GetLatestBlock() (*Block, error) {
	latestHeader, err := e.Client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return &Block{}, err
	}

	latestBlock, err := e.Client.BlockByNumber(context.Background(), latestHeader.Number)
	if err != nil {
		return &Block{}, err
	}

	transactions, err := getTransactions(latestBlock)
	if err != nil {
		return &Block{}, err
	}

	return &Block{
		Hash:         latestBlock.Hash().Hex(),
		Number:       latestBlock.Number().Uint64(),
		Transactions: transactions,
	}, nil
}

// GetBlockByHash retrieves a valid block from the EVM based chain, based on the provided hash
func (e *EthereumClient) GetBlockByHash(hash string) (*Block, error) {
	block, err := e.Client.BlockByHash(context.Background(), common.HexToHash(hash))
	if err != nil {
		return &Block{}, err
	}

	transactions, err := getTransactions(block)
	if err != nil {
		return &Block{}, err
	}

	return &Block{
		Hash:         block.Hash().Hex(),
		Number:       block.NumberU64(),
		Transactions: transactions,
	}, nil
}

// SendNativeTransaction sends a specified amount of WEI from a selected wallet to an address
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

// Marshalls ethereum specific transactions in a block into our generic transactions type
func getTransactions(block *types.Block) (Transactions, error) {
	transactions := make(Transactions)
	for _, tx := range block.Transactions() {
		message, err := tx.AsMessage(types.EIP155Signer{})
		if err != nil {
			return nil, err
		}

		transactions[tx.Hash().Hex()] = &Transaction{
			From:            message.From().Hex(),
			To:              tx.To().Hex(),
			NativeAmount:    tx.Value(),
			LinkTokenAmount: nil, // TODO: This is tricky, looking into it further
		}

	}

	return transactions, nil
}
