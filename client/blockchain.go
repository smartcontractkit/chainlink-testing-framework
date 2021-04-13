package client

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
	"os"
)

const NetworkEthereumHardhat = "Ethereum Hardhat"

// Generalized blockchain client for interaction with multiple different blockchains
type BlockchainClient interface {
	DeployStorageContract(wallet BlockchainWallet) error
}

// NewBlockchainClient returns an implementation of a BlockchainClient based on the given network
func NewBlockchainClient(network BlockchainNetwork) (BlockchainClient, error) {
	switch network.(type) {
	case *EthereumHardhat:
		return NewEthereumClient(network)
	}
	return nil, fmt.Errorf("invalid blockchain network was given")
}

// BlockchainNetwork is the interface that when implemented, defines a new blockchain network that can
// be tested against
type BlockchainNetwork interface {
	Name() string
	URL() string
	ChainID() *big.Int
}

// EthereumHardhat is the implementation of BlockchainNetwork for the local ETH dev server
type EthereumHardhat struct{}

// Name returns the readable name of the hardhat network
func (e *EthereumHardhat) Name() string {
	return NetworkEthereumHardhat
}

// URL returns the RPC URL used for connecting to hardhat
func (e *EthereumHardhat) URL() string {
	return ethereumURL(e)
}

// ChainID returns the on-chain ID of the network being connected to, returning hardhats default
func (e *EthereumHardhat) ChainID() *big.Int {
	return big.NewInt(31337)
}

// BlockchainWallets is an interface that when implemented is a representation of a slice of wallets for
// a specific network
type BlockchainWallets interface {
	Default() BlockchainWallet
	SetDefault(i int) error
	Wallet(i int) (BlockchainWallet, error)
}

// Wallets is the default implementation of BlockchainWallets that holds a slice of wallets with the default
type Wallets struct {
	defaultWallet int
	wallets       []BlockchainWallet
}

// Default returns the default wallet to be used for a transaction on-chain
func (e *Wallets) Default() BlockchainWallet {
	return e.wallets[e.defaultWallet]
}

// SetDefault changes the default wallet to be used for on-chain transactions
func (e *Wallets) SetDefault(i int) error {
	if err := walletSliceIndexInRange(e.wallets, i); err != nil {
		return err
	}
	e.defaultWallet = i
	return nil
}

// Wallet returns a wallet based on a given index in the slice
func (e *Wallets) Wallet(i int) (BlockchainWallet, error) {
	if err := walletSliceIndexInRange(e.wallets, i); err != nil {
		return nil, err
	}
	return e.wallets[i], nil
}

// BlockchainWallet when implemented is the interface to allow multiple wallet implementations for each
// BlockchainNetwork that is supported
type BlockchainWallet interface {
	PrivateKey() string
	Address() string
}

// EthereumWallet is the implementation to allow testing with ETH based wallets
type EthereumWallet struct {
	privateKey string
	address    common.Address
}

// NewEthereumWallet returns the instantiated ETH wallet based on a given private key
func NewEthereumWallet(pk string) (*EthereumWallet, error) {
	privateKey, err := crypto.HexToECDSA(pk)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %v", err)
	}
	return &EthereumWallet{
		privateKey: pk,
		address:    crypto.PubkeyToAddress(privateKey.PublicKey),
	}, nil
}

// PrivateKey returns the private key for a given Ethereum wallet
func (e *EthereumWallet) PrivateKey() string {
	return e.privateKey
}

// Address returns the ETH address for a given wallet
func (e *EthereumWallet) Address() string {
	return e.address.String()
}

// DefaultHardhatWallets returns the instantiated BlockchainWallets containing the default set of Hardhat wallets
func DefaultHardhatWallets() BlockchainWallets {
	w0, _ := NewEthereumWallet("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	w1, _ := NewEthereumWallet("59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d")
	w2, _ := NewEthereumWallet("5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a")
	return &Wallets{
		defaultWallet: 0,
		wallets: []BlockchainWallet{w0, w1, w2},
	}
}

func ethereumURL(network BlockchainNetwork) string {
	env := getNetworkURLEnv(network.ChainID())
	if len(env) > 0 {
		return env
	}
	return "http://localhost:8545"
}

func getNetworkURLEnv(chainID *big.Int) string {
	return os.Getenv(fmt.Sprintf("NETWORK_%d_URL", chainID))
}

func walletSliceIndexInRange(wallets []BlockchainWallet, i int) error {
	if i > len(wallets)-1 {
		return fmt.Errorf("invalid index in list of wallets")
	}
	return nil
}
