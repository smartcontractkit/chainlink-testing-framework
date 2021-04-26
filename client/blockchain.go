package client

import (
	"fmt"
	"integrations-framework/config"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const EthereumHardhatID = "ethereum_hardhat"

// Generalized blockchain client for interaction with multiple different blockchains
type BlockchainClient interface {
	// Common blockchain interactions
	GetLatestBlock() (*Block, error)
	GetBlockByHash(blockHash string) (*Block, error)
	GetNativeBalance(addressHex string) (*big.Int, error)
	GetLinkBalance(addressHex string) (*big.Int, error)
	SendNativeTransaction(fromWallet BlockchainWallet, toHexAddress string, amount *big.Int) (string, error)
	SendLinkTransaction(fromWallet BlockchainWallet, toHexAddress string, amount *big.Int) (string, error)

	// Specific smart contract interactions
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

// BlockchainNetwork is the interface that when implemented, defines a new blockchain network that can be tested against
type BlockchainNetwork interface {
	ID() string
	URL() string
	ChainID() *big.Int
	Wallets() (BlockchainWallets, error)
	Config() *config.NetworkConfig
}

type BlockchainNetworkInit func(conf *config.Config) BlockchainNetwork

type Block struct {
	Hash   string
	Number uint64
}

// EthereumHardhat is the implementation of BlockchainNetwork for the local ETH dev server
type EthereumHardhat struct {
	networkConfig *config.NetworkConfig
}

// NewEthereumHardhat creates a way to interact with the ethereum hardhat blockchain
func NewEthereumHardhat(conf *config.Config) BlockchainNetwork {
	networkConf, _ := conf.GetNetworkConfig(EthereumHardhatID)
	return &EthereumHardhat{networkConf}
}

// ID returns the readable name of the hardhat network
func (e *EthereumHardhat) ID() string {
	return EthereumHardhatID
}

// URL returns the RPC URL used for connecting to hardhat
func (e *EthereumHardhat) URL() string {
	return e.networkConfig.URL
}

// ChainID returns the on-chain ID of the network being connected to, returning hardhat's default
func (e *EthereumHardhat) ChainID() *big.Int {
	return big.NewInt(e.networkConfig.ChainID)
}

// Config returns the blockchain network configuration
func (e *EthereumHardhat) Config() *config.NetworkConfig {
	return e.networkConfig
}

// Wallets returns all the viable wallets used for testing on chain, returning hardhat's default
func (e *EthereumHardhat) Wallets() (BlockchainWallets, error) {
	return newEthereumWallets(e.networkConfig.PrivateKeyStore)
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
func (w *Wallets) Default() BlockchainWallet {
	return w.wallets[w.defaultWallet]
}

// SetDefault changes the default wallet to be used for on-chain transactions
func (w *Wallets) SetDefault(i int) error {
	if err := walletSliceIndexInRange(w.wallets, i); err != nil {
		return err
	}
	w.defaultWallet = i
	return nil
}

// Wallet returns a wallet based on a given index in the slice
func (w *Wallets) Wallet(i int) (BlockchainWallet, error) {
	if err := walletSliceIndexInRange(w.wallets, i); err != nil {
		return nil, err
	}
	return w.wallets[i], nil
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

func newEthereumWallets(pkStore config.PrivateKeyStore) (BlockchainWallets, error) {
	// Check private keystore value, create wallets from such
	var processedWallets []BlockchainWallet
	keys, err := pkStore.Fetch()
	if err != nil {
		return nil, err
	}

	for _, key := range keys {
		wallet, err := NewEthereumWallet(strings.TrimSpace(key))
		if err != nil {
			return &Wallets{}, err
		}
		processedWallets = append(processedWallets, wallet)
	}

	return &Wallets{
		defaultWallet: 0,
		wallets:       processedWallets,
	}, nil
}

func walletSliceIndexInRange(wallets []BlockchainWallet, i int) error {
	if i > len(wallets)-1 {
		return fmt.Errorf("invalid index in list of wallets")
	}
	return nil
}
