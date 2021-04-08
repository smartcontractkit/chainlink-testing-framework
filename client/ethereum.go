package client

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"integrations-framework/contracts"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Etherum client that wraps the go-ethereum client and adds some helper methods
type EthereumClient struct {
	Client bind.ContractBackend
}

// Builds a new ethereum client based on a connection string
func NewEthereumClient(rpcConnectionString string) EthereumClient {
	cl, err := ethclient.Dial(rpcConnectionString)
	if err != nil {
		log.Fatal(err)
	}

	return EthereumClient{Client: cl}
}

// Builds a new ethereum client pointing to a locally simulated blockchain
func NewSimulatedEthereumClient(backend *backends.SimulatedBackend) EthereumClient {
	return EthereumClient{Client: backend}
}

// Creates a default contract (need to parameterize this)
func (clientWrapper EthereumClient) DeployStorageContract() (common.Address, *types.Transaction, *contracts.Storage) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
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

	chainID := big.NewInt(0)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatal(err)
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(99999999999) // in wei
	auth.GasLimit = uint64(1)            // in units
	auth.GasPrice = gasPrice

	addr, tx, instance, err := contracts.DeployStorage(auth, clientWrapper.Client, "1.0")
	if err != nil {
		log.Fatal(err)
	}

	return addr, tx, instance
}
