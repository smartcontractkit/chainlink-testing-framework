// Package client handles connections between chainlink nodes and different blockchain networks
package client

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/smartcontractkit/integrations-framework/config"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const BlockchainTypeEVM = "evm"

// BlockchainClient is the interface that wraps a given client implementation for a blockchain, to allow for switching
// of network types within the test suite
type BlockchainClient interface {
	Get() interface{}
	Fund(fromWallet BlockchainWallet, toAddress string, nativeAmount, linkAmount *big.Int) error
}

// NewBlockchainClient returns an instantiated network client implementation based on the network configuration given
func NewBlockchainClient(network BlockchainNetwork) (BlockchainClient, error) {
	switch network.Type() {
	case BlockchainTypeEVM:
		return NewEthereumClient(network)
	}
	return nil, errors.New("invalid blockchain network ID, not found")
}

// BlockchainNetwork is the interface that when implemented, defines a new blockchain network that can be tested against
type BlockchainNetwork interface {
	ID() string
	URL() string
	Type() string
	SetURL(string)
	ChainID() *big.Int
	Wallets() (BlockchainWallets, error)
	Config() *config.NetworkConfig
}

// BlockchainNetworkInit is a helper function to obtain different blockchain networks
type BlockchainNetworkInit func(conf *config.Config) (BlockchainNetwork, error)

// EthereumNetwork is the implementation of BlockchainNetwork for the local ETH dev server
type EthereumNetwork struct {
	networkID     string
	networkConfig *config.NetworkConfig
}

// NewEthereumNetwork creates a way to interact with any specified EVM blockchain
func newEthereumNetwork(ID string, networkConfig *config.NetworkConfig) (BlockchainNetwork, error) {
	return &EthereumNetwork{
		networkID:     ID,
		networkConfig: networkConfig,
	}, nil
}

// NewNetworkFromConfig prepares settings for a connection to a hardhat blockchain
func NewNetworkFromConfig(conf *config.Config) (BlockchainNetwork, error) {
	networkConfig, err := conf.GetNetworkConfig(conf.Network)
	if err != nil {
		return nil, err
	}
	switch networkConfig.Type {
	case BlockchainTypeEVM:
		return newEthereumNetwork(conf.Network, networkConfig)
	}
	return nil, fmt.Errorf(
		"network %s uses an unspported network type of: %s",
		conf.Network,
		networkConfig.Type,
	)
}

// ID returns the readable name of the EVM network
func (e *EthereumNetwork) ID() string {
	return e.networkID
}

// Type returns the readable type of the EVM network
func (e *EthereumNetwork) Type() string {
	return BlockchainTypeEVM
}

// URL returns the RPC URL used for connecting to the network
func (e *EthereumNetwork) URL() string {
	return e.networkConfig.URL
}

// SetURL sets the RPC URL, useful for when blockchain URLs might be dynamic
func (e *EthereumNetwork) SetURL(newURL string) {
	e.networkConfig.URL = newURL
}

// ChainID returns the on-chain ID of the network being connected to
func (e *EthereumNetwork) ChainID() *big.Int {
	return big.NewInt(e.networkConfig.ChainID)
}

// Config returns the blockchain network configuration
func (e *EthereumNetwork) Config() *config.NetworkConfig {
	return e.networkConfig
}

// Wallets returns all the viable wallets used for testing on chain
func (e *EthereumNetwork) Wallets() (BlockchainWallets, error) {
	return newEthereumWallets(e.networkConfig.PrivateKeyStore)
}

// BlockchainWallets is an interface that when implemented is a representation of a slice of wallets for
// a specific network
type BlockchainWallets interface {
	Default() BlockchainWallet
	All() []BlockchainWallet
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

// All returns the raw representation of Wallets
func (w *Wallets) All() []BlockchainWallet {
	return w.wallets
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
