// Package client handles connections between chainlink nodes and different blockchain networks
package client

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"os"
	"reflect"
	"strings"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/integrations-framework/config"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// Commonly used variables
const (
	BlockchainTypeEVM          = "evm"
	BlockchainTypeEVMMultinode = "evm_multi"
	NetworkGethPerformance     = "ethereum_geth_performance"
)

// BlockchainClient is the interface that wraps a given client implementation for a blockchain, to allow for switching
// of network types within the test suite
type BlockchainClient interface {
	Get() interface{}
	GetNetworkName() string
	GetID() int
	SetID(id int)
	SetDefaultClient(clientID int) error
	GetClients() []BlockchainClient
	SuggestGasPrice(ctx context.Context) (*big.Int, error)
	HeaderHashByNumber(ctx context.Context, bn *big.Int) (string, error)
	BlockNumber(ctx context.Context) (uint64, error)
	HeaderTimestampByNumber(ctx context.Context, bn *big.Int) (uint64, error)
	CalculateTxGas(gasUsedValue *big.Int) (*big.Float, error)
	Fund(fromWallet BlockchainWallet, toAddress string, nativeAmount, linkAmount *big.Float) error
	GasStats() *GasStats
	ParallelTransactions(enabled bool)
	Close() error

	AddHeaderEventSubscription(key string, subscriber HeaderEventSubscription)
	DeleteHeaderEventSubscription(key string)
	WaitForEvents() error
}

// NewBlockchainClient returns an instantiated network client implementation based on the network configuration given
func NewBlockchainClient(network BlockchainNetwork) (BlockchainClient, error) {
	switch network.Type() {
	case BlockchainTypeEVM:
		return NewEthereumClient(network)
	case BlockchainTypeEVMMultinode:
		return NewEthereumClients(network)
	}
	return nil, errors.New("invalid blockchain network ID, not found")
}

// BlockchainNetwork is the interface that when implemented, defines a new blockchain network that can be tested against
type BlockchainNetwork interface {
	GasUsedEstimations
	ID() string
	ClusterURL() string
	LocalURL() string
	URLs() []string
	Type() string
	SetClusterURL(string)
	SetLocalURL(string)
	SetURLs(urls []string)
	ChainID() *big.Int
	RemotePort() uint16
	Wallets() (BlockchainWallets, error)
	Config() *config.NetworkConfig
}

// BlockchainNetworkInit is a helper function to obtain the network listed in the config file
type BlockchainNetworkInit func(conf *config.Config) (BlockchainNetwork, error)

// MultiNetworkInit is a helper function to create multiple blockchain networks at once
type MultiNetworkInit func(conf *config.Config) ([]BlockchainNetwork, error)

// EthereumNetwork is the implementation of BlockchainNetwork for the local ETH dev server
type EthereumNetwork struct {
	networkID     string
	networkConfig *config.NetworkConfig
}

// NewEthereumNetwork creates a way to interact with any specified EVM blockchain
func newEthereumNetwork(ID string, networkConfig config.NetworkConfig) (BlockchainNetwork, error) {
	return &EthereumNetwork{
		networkID:     ID,
		networkConfig: &networkConfig,
	}, nil
}

// DefaultNetworkFromConfig prepares settings for a connection the default blockchain specified in the config file
func DefaultNetworkFromConfig(conf *config.Config) (BlockchainNetwork, error) {
	if len(conf.Networks) <= 0 {
		return nil, fmt.Errorf("No default network(s) provided in config")
	}
	return NewNetworkFromConfig(conf, conf.Networks[0])
}

// DefaultNetworksFromConfig prepares settings for multiple connections to the default blockchains specified in the config file
func DefaultNetworksFromConfig(conf *config.Config) ([]BlockchainNetwork, error) {
	if len(conf.Networks) <= 0 {
		return nil, fmt.Errorf("No default networks provided in config")
	} else if len(conf.Networks) == 1 {
		return nil, fmt.Errorf("Only one network provided in config: '%s'", conf.Networks[0])
	}

	networks := []BlockchainNetwork{}
	for _, networkID := range conf.Networks {
		network, err := NewNetworkFromConfig(conf, networkID)
		if err != nil {
			return nil, err
		}
		networks = append(networks, network)
	}
	return networks, nil
}

// NewNetworkFromConfig creates a new blockchain network based on the ID
func NewNetworkFromConfig(conf *config.Config, networkID string) (BlockchainNetwork, error) {
	networkConfig, err := conf.GetNetworkConfig(networkID)
	if err != nil {
		return nil, err
	}
	switch networkConfig.Type {
	case BlockchainTypeEVM, BlockchainTypeEVMMultinode:
		return newEthereumNetwork(networkID, networkConfig)
	}
	return nil, fmt.Errorf(
		"network %s uses an unspported network type of: %s",
		networkID,
		networkConfig.Type,
	)
}

// NewNetworkFromConfigWithDefault will return a new network with config but with a customisable default in-case a test
// needs a different default from the config file default. For example, performance tests usually need a "performance"
// network with higher gas limits than the vanilla Geth network.
func NewNetworkFromConfigWithDefault(networkID string) BlockchainNetworkInit {
	return func(conf *config.Config) (BlockchainNetwork, error) {
		if conf == nil {
			return nil, errors.New("nil config was provided")
		}
		// Get the "yaml" tag on the config struct to check whether the ENV override has been specified and if so
		// use the override rather than the default given
		ct := reflect.TypeOf(*conf)
		field := ct.Field(0)
		networkKey := field.Tag.Get("yaml")
		if len(os.Getenv(strings.ToUpper(networkKey))) == 0 {
			conf.Networks = []string{networkID}
		}
		return DefaultNetworkFromConfig(conf)
	}
}

// ID returns the readable name of the EVM network
func (e *EthereumNetwork) ID() string {
	return e.networkID
}

// Type returns the readable type of the EVM network
func (e *EthereumNetwork) Type() string {
	return e.networkConfig.Type
}

// ClusterURL returns the RPC URL used for connecting to the network within the K8s cluster
func (e *EthereumNetwork) ClusterURL() string {
	return e.networkConfig.ClusterURL
}

// LocalURL returns the RPC URL used for connecting to the network from outside the K8s cluster
func (e *EthereumNetwork) LocalURL() string {
	return e.networkConfig.LocalURL
}

// URLs returns the RPC URLs used for connecting to the network nodes
func (e *EthereumNetwork) URLs() []string {
	return e.networkConfig.URLS
}

// SetURLs sets all nodes URLs
func (e *EthereumNetwork) SetURLs(urls []string) {
	e.networkConfig.URLS = urls
}

// SetClusterURL sets the RPC URL used to connect to the chain from within the K8s cluster
func (e *EthereumNetwork) SetClusterURL(newURL string) {
	e.networkConfig.ClusterURL = newURL
}

// SetLocalURL sets the RPC URL used to connect to the chain from outside the K8s cluster
func (e *EthereumNetwork) SetLocalURL(newURL string) {
	e.networkConfig.LocalURL = newURL
}

// ChainID returns the on-chain ID of the network being connected to
func (e *EthereumNetwork) ChainID() *big.Int {
	return big.NewInt(e.networkConfig.ChainID)
}

// Config returns the blockchain network configuration
func (e *EthereumNetwork) Config() *config.NetworkConfig {
	return e.networkConfig
}

// RemotePort returns the remote RPC port of the network
func (e *EthereumNetwork) RemotePort() uint16 {
	return e.networkConfig.RPCPort
}

// Wallets returns all the viable wallets used for testing on chain
func (e *EthereumNetwork) Wallets() (BlockchainWallets, error) {
	return newEthereumWallets(e.networkConfig.PrivateKeyStore)
}

// FluxMonitorSubmissionGasUsed Flux Monitor one submission gasUsed value
func (e *EthereumNetwork) FluxMonitorSubmissionGasUsed() (*big.Int, error) {
	if strings.HasPrefix(e.networkConfig.Name, "ethereum-geth") {
		return big.NewInt(400000), nil
	}
	return nil, fmt.Errorf("gas used estimation unavailable for the network name '%s'", e.networkConfig.Name)
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

// NodeBlock block with a node ID which mined it
type NodeBlock struct {
	NodeID int
	*types.Block
}

// HeaderEventSubscription is an interface for allowing callbacks when the client receives a new header
type HeaderEventSubscription interface {
	ReceiveBlock(header NodeBlock) error
	Wait() error
}

// GasUsedEstimations contains some known gas values for contracts for every network
type GasUsedEstimations interface {
	FluxMonitorSubmissionGasUsed() (*big.Int, error)
}
