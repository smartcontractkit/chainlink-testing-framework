// Package blockchain handles connections to various blockchains
package blockchain

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/url"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/helmenv/environment"
	"gopkg.in/yaml.v2"

	"github.com/smartcontractkit/integrations-framework/config"
)

// Commonly used blockchain network types
const (
	SimulatedEthNetwork   = "eth_simulated"
	LiveEthTestNetwork    = "eth_testnet"
	LiveKlaytnTestNetwork = "klaytn_testnet"
)

// NewBlockchainClientFn external client implementation function
// networkName must match a key in "networks" in networks.yaml config
// networkConfig is just an arbitrary config you provide in "networks" for your key
type NewEVMClientFn func(
	networkName string,
	networkConfig map[string]interface{},
	urls []*url.URL,
) (EVMClient, error)

// ClientURLFn are used to be able to return a list of URLs from the environment to connect
type ClientURLFn func(e *environment.Environment) ([]*url.URL, error)

// EVMClient is the interface that wraps a given client implementation for a blockchain, to allow for switching
// of network types within the test suite
// EVMClient can be connected to a single or multiple nodes,
type EVMClient interface {
	// Getters
	Get() interface{}
	GetNetworkName() string
	GetNetworkType() string
	GetChainID() *big.Int
	GetClients() []EVMClient
	GetDefaultWallet() *EthereumWallet
	GetWallets() []*EthereumWallet
	GetNetworkConfig() *config.ETHNetwork

	// Setters
	SetID(id int)
	SetDefaultWallet(num int) error
	SetWallets([]*EthereumWallet)
	LoadWallets(ns interface{}) error
	SwitchNode(node int) error

	// On-chain Operations
	HeaderHashByNumber(ctx context.Context, bn *big.Int) (string, error)
	HeaderTimestampByNumber(ctx context.Context, bn *big.Int) (uint64, error)
	LatestBlockNumber(ctx context.Context) (uint64, error)
	Fund(toAddress string, amount *big.Float) error
	DeployContract(
		contractName string,
		deployer ContractDeployer,
	) (*common.Address, *types.Transaction, interface{}, error)
	TransactionOpts(from *EthereumWallet) (*bind.TransactOpts, error)
	ProcessTransaction(tx *types.Transaction) error
	IsTxConfirmed(txHash common.Hash) (bool, error)
	ParallelTransactions(enabled bool)
	Close() error

	// Gas Operations
	EstimateCostForChainlinkOperations(amountOfOperations int) (*big.Float, error)
	EstimateTransactionGasCost() (*big.Int, error)
	GasStats() *GasStats

	// Event Subscriptions
	AddHeaderEventSubscription(key string, subscriber HeaderEventSubscription)
	DeleteHeaderEventSubscription(key string)
	WaitForEvents() error
}

// Networks is a thin wrapper that just selects client connected to some network
// if there is only one client it is chosen as Default
// if there is multiple you just get clients you need in test
type Networks struct {
	clients []EVMClient
	Default EVMClient
}

// Teardown all clients
func (b *Networks) Teardown() error {
	for _, c := range b.clients {
		if err := c.Close(); err != nil {
			return err
		}
	}
	return nil
}

// SetDefault chooses default client
func (b *Networks) SetDefault(index int) error {
	if index > len(b.clients) {
		return fmt.Errorf("index of %d is out of bounds", index)
	}
	b.Default = b.clients[index]
	return nil
}

// Get gets blockchain network (client) by name
func (b *Networks) Get(index int) (EVMClient, error) {
	if index > len(b.clients) {
		return nil, fmt.Errorf("index of %d is out of bounds", index)
	}
	return b.clients[index], nil
}

// AllNetworks returns all the network clients
func (b *Networks) AllNetworks() []EVMClient {
	return b.clients
}

// NetworkRegistry holds all the registered network types that can be initialized, allowing
// external libraries to register alternative network types to use
type NetworkRegistry struct {
	registeredNetworks map[string]registeredNetwork
}

type registeredNetwork struct {
	newBlockchainClientFn NewEVMClientFn
	blockchainClientURLFn ClientURLFn
}

// NewDefaultNetworkRegistry returns an instance of the network registry with the default supported networks registered
func NewDefaultNetworkRegistry() *NetworkRegistry {
	return &NetworkRegistry{
		registeredNetworks: map[string]registeredNetwork{
			SimulatedEthNetwork: {
				newBlockchainClientFn: NewEthereumMultiNodeClient,
				blockchainClientURLFn: SimulatedEthereumURLs,
			},
			LiveEthTestNetwork: {
				newBlockchainClientFn: NewEthereumMultiNodeClient,
				blockchainClientURLFn: LiveEthTestnetURLs,
			},
			LiveKlaytnTestNetwork: {
				newBlockchainClientFn: NewKlaytnMultiNodeClient,
				blockchainClientURLFn: LiveEthTestnetURLs,
			},
		},
	}
}

// NewSoakNetworkRegistry retrieves a network registry for use in soak tests
func NewSoakNetworkRegistry() *NetworkRegistry {
	return &NetworkRegistry{
		registeredNetworks: map[string]registeredNetwork{
			SimulatedEthNetwork: {
				newBlockchainClientFn: NewEthereumMultiNodeClient,
				blockchainClientURLFn: SimulatedSoakEthereumURLs,
			},
			LiveEthTestNetwork: {
				newBlockchainClientFn: NewEthereumMultiNodeClient,
				blockchainClientURLFn: LiveEthTestnetURLs,
			},
			LiveKlaytnTestNetwork: {
				newBlockchainClientFn: NewKlaytnMultiNodeClient,
				blockchainClientURLFn: LiveEthTestnetURLs,
			},
		},
	}
}

// RegisterNetwork registers a new type of network within the registry
func (n *NetworkRegistry) RegisterNetwork(networkType string, fn NewEVMClientFn, urlFn ClientURLFn) {
	n.registeredNetworks[networkType] = registeredNetwork{
		newBlockchainClientFn: fn,
		blockchainClientURLFn: urlFn,
	}
}

// GetNetworks returns a networks object with all the BlockchainClient(s) initialized
func (n *NetworkRegistry) GetNetworks(env *environment.Environment) (*Networks, error) {
	nc := config.ProjectNetworkSettings
	var clients []EVMClient
	for _, networkName := range nc.SelectedNetworks {
		networkSettings, ok := nc.NetworkSettings[networkName]
		if !ok {
			return nil, fmt.Errorf("network with the name of '%s' doesn't exist in the network config", networkName)
		}
		networkType, ok := networkSettings["type"]
		if !ok {
			return nil, fmt.Errorf("network config for '%s' doesn't define a 'type'", networkName)
		}
		initFn, ok := n.registeredNetworks[fmt.Sprint(networkType)]
		if !ok {
			return nil, fmt.Errorf("network '%s' of type '%s' hasn't been registered", networkName, networkType)
		}
		urls, err := initFn.blockchainClientURLFn(env)
		if err != nil {
			return nil, err
		}
		client, err := initFn.newBlockchainClientFn(networkName, networkSettings, urls)
		if err != nil {
			return nil, err
		}
		clients = append(clients, client)
	}
	var defaultClient EVMClient
	if len(clients) >= 1 {
		defaultClient = clients[0]
	}
	return &Networks{
		clients: clients,
		Default: defaultClient,
	}, nil
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

// UnmarshalNetworkConfig is a generic function to unmarshal a yaml map into a given object
func UnmarshalNetworkConfig(config map[string]interface{}, obj interface{}) error {
	b, err := json.Marshal(config)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(b, obj)
}

// ContractDeployer acts as a go-between function for general contract deployment
type ContractDeployer func(auth *bind.TransactOpts, backend bind.ContractBackend) (
	common.Address,
	*types.Transaction,
	interface{},
	error,
)
