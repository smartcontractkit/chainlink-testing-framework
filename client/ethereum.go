package client

import (
	"context"
	"integrations-framework/contracts"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
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

// SendTransaction sends a specified amount of WEI from a selected wallet to an address
func (e *EthereumClient) SendTransaction(fromWallet BlockchainWallet, toAddress string, weiAmount int64) (string, error) {
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
		types.NewTransaction(nonce, common.HexToAddress(toAddress), big.NewInt(weiAmount),
			e.Network.Config().TransactionLimit, gasPrice, nil)

	signedTransaction, err := types.SignTx(unsignedTransaction, types.NewEIP2930Signer(e.Network.ChainID()), privateKey)
	if err != nil {
		return "", err
	}

	err = e.Client.SendTransaction(context.Background(), signedTransaction)
	return signedTransaction.Hash().Hex(), err
}
