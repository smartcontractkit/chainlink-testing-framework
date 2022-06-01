// Package blockchain handles connections to various blockchains
package blockchain

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/helmenv/environment"
	"gopkg.in/yaml.v2"

	"github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink-testing-framework/config"
)

// NewEVMClientFn enables connection to a new EVM client
type NewEVMClientFn func(networkSettings *config.EVMNetwork, testEnvironment *environment.Environment) (EVMClient, error)

// EVMClient is the interface that wraps a given client implementation for a blockchain, to allow for switching
// of network types within the test suite
// EVMClient can be connected to a single or multiple nodes,
type EVMClient interface {
	// Getters
	Get() interface{}
	GetNetworkName() string
	GetChainID() *big.Int
	GetClients() []EVMClient
	GetDefaultWallet() *EthereumWallet
	GetWallets() []*EthereumWallet
	GetNetworkConfig() *config.EVMNetwork
	GetEVMNodeAttributes() *client.EVMNodeAttributes
	GetEVMChainAttributes() *client.EVMChainAttributes

	// Setters
	SetID(id int)
	SetDefaultWallet(num int) error
	SetWallets([]*EthereumWallet)
	LoadWallets(ns interface{}) error
	SwitchNode(node int) error
	SetEVMChainAttributes(attrs *client.EVMChainAttributes)

	// On-chain Operations
	BalanceAt(ctx context.Context, address common.Address) (*big.Int, error)
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
	evmClients []EVMClient
	Default    EVMClient
}

// Teardown all clients
func (b *Networks) Teardown() error {
	for _, c := range b.evmClients {
		if err := c.Close(); err != nil {
			return err
		}
	}
	return nil
}

// SetDefault chooses default client
func (b *Networks) SetDefault(index int) error {
	if index > len(b.evmClients) {
		return fmt.Errorf("index of %d is out of bounds", index)
	}
	b.Default = b.evmClients[index]
	return nil
}

// Get gets blockchain network (client) by name
func (b *Networks) Get(index int) (EVMClient, error) {
	if index > len(b.evmClients) {
		return nil, fmt.Errorf("index of %d is out of bounds", index)
	}
	return b.evmClients[index], nil
}

// AllNetworks returns all the network clients
func (b *Networks) EVMNetworks() []EVMClient {
	return b.evmClients
}

// ConnectNetworks goes through the selected networks and builds clients and connections for them
func ConnectNetworks(env *environment.Environment) (*Networks, error) {
	networksConfig := config.ProjectConfig.NetworksConfig
	var evmClients []EVMClient
	for _, networkName := range networksConfig.SelectedNetworks {
		// If more network types are added, can add to this if-else chain to detect and properly connect to them
		if isEVMNetwork(networksConfig, networkName) {
			evmClient, err := connectEVMNetwork(env, networksConfig, networkName)
			if err != nil {
				return nil, err
			}
			evmClients = append(evmClients, evmClient)
		} else {
			return nil, fmt.Errorf("the network '%s' was not found defined anywhere in your networks.yaml file", networkName)
		}
	}
	var defaultClient EVMClient
	if len(evmClients) >= 1 {
		defaultClient = evmClients[0]
	}
	return &Networks{
		evmClients: evmClients,
		Default:    defaultClient,
	}, nil
}

// checks if the network name is registered as an EVM network
func isEVMNetwork(networksConfig *config.NetworksConfig, networkName string) bool {
	if networksConfig.EVMNetworkSettings == nil {
		log.Warn().Msg("No EVM network configs defined")
		return false
	}
	_, exists := networksConfig.EVMNetworkSettings.GetNetworkSettings(networkName)
	return exists
}

// connectEVMNetwork attempts to connect to the EVM network
func connectEVMNetwork(
	env *environment.Environment,
	networksConfig *config.NetworksConfig,
	networkName string,
) (EVMClient, error) {
	networkSettings, _ := networksConfig.EVMNetworkSettings.GetNetworkSettings(networkName)
	if networkSettings.ChainID == 0 {
		return nil, fmt.Errorf("Chain ID not set for network '%s'", networkName)
	}
	network := RegisteredEVMNetworks.GetNetwork(networkSettings.ChainID)
	connectedNetwork, err := network.ConnectClientFunc(networkSettings, env)
	if err != nil {
		return nil, err
	}
	connectedNetwork.SetEVMChainAttributes(network.ChainAttributes)
	return connectedNetwork, nil
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
