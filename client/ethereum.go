package client

import (
	"context"
	"crypto/ecdsa"
	"integrations-framework/contracts"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Etherum client that wraps the go-ethereum client and adds some helper methods
type EthereumClient struct {
	Client        *ethclient.Client
	EthChainID    *big.Int
	SourceAddress common.Address
}

// Builds a new ethereum client based on a connection string
// Need to handle rpc over websocket as well
func NewEthereumClient(rpcConnectionString string, chainID *big.Int) EthereumClient {
	cl, err := ethclient.Dial(rpcConnectionString)
	if err != nil {
		log.Fatal(err)
	}

	return EthereumClient{
		Client:     cl,
		EthChainID: chainID,
	}
}

// Creates a default contract (need to parameterize this)
func (clientWrapper EthereumClient) DeployStorageContract() (common.Address, *types.Transaction, *contracts.Storage) {
	// Needs paramaterization to work with hardhat and others
	privateKey, err := crypto.HexToECDSA("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := clientWrapper.Client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := clientWrapper.Client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, clientWrapper.EthChainID)
	if err != nil {
		log.Fatal(err)
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(3) // in wei
	// Needs parameterization
	auth.GasLimit = 9500000 // in units
	auth.GasPrice = gasPrice

	addr, tx, instance, err := contracts.DeployStorage(auth, clientWrapper.Client, "1.0")
	if err != nil {
		log.Fatal(err)
	}

	return addr, tx, instance
}
