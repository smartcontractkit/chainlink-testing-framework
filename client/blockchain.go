// Package client handles connections between chainlink nodes and different blockchain networks
package client

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/url"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/helmenv/environment"
	"gopkg.in/yaml.v2"

	"github.com/smartcontractkit/integrations-framework/config"
)

// Commonly used blockchain network types
const (
	SimulatedEthNetwork    = "eth_simulated"
	LiveEthTestNetwork     = "eth_testnet"
	NetworkGethPerformance = "ethereum_geth_performance"
)

// NewBlockchainClientFn external client implementation function
// networkName must match a key in "networks" in networks.yaml config
// networkConfig is just an arbitrary config you provide in "networks" for your key
type NewBlockchainClientFn func(
	networkName string,
	networkConfig map[string]interface{},
	urls []*url.URL,
) (BlockchainClient, error)

// BlockchainClientURLFn are used to be able to return a list of URLs from the environment to connect
type BlockchainClientURLFn func(e *environment.Environment) ([]*url.URL, error)

// BlockchainClient is the interface that wraps a given client implementation for a blockchain, to allow for switching
// of network types within the test suite
// BlockchainClient can be connected to a single or multiple nodes,
type BlockchainClient interface {
	ContractsDeployed() bool
	LoadWallets(ns interface{}) error
	SetWallet(num int) error
	GetDefaultWallet() *EthereumWallet

	EstimateCostForChainlinkOperations(amountOfOperations int) (*big.Float, error)
	EstimateTransactionGasCost() (*big.Int, error)

	Get() interface{}
	GetNetworkName() string
	GetNetworkType() string
	GetChainID() int64
	SwitchNode(node int) error
	GetClients() []BlockchainClient
	HeaderHashByNumber(ctx context.Context, bn *big.Int) (string, error)
	BlockNumber(ctx context.Context) (uint64, error)
	HeaderTimestampByNumber(ctx context.Context, bn *big.Int) (uint64, error)
	Fund(toAddress string, amount *big.Float) error
	GasStats() *GasStats
	ParallelTransactions(enabled bool)
	Close() error

	AddHeaderEventSubscription(key string, subscriber HeaderEventSubscription)
	DeleteHeaderEventSubscription(key string)
	WaitForEvents() error
}

// Networks is a thin wrapper that just selects client connected to some network
// if there is only one client it is chosen as Default
// if there is multiple you just get clients you need in test
type Networks struct {
	clients []BlockchainClient
	Default BlockchainClient
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
func (b *Networks) Get(index int) (BlockchainClient, error) {
	if index > len(b.clients) {
		return nil, fmt.Errorf("index of %d is out of bounds", index)
	}
	return b.clients[index], nil
}

// AllNetworks returns all the network clients
func (b *Networks) AllNetworks() []BlockchainClient {
	return b.clients
}

// ConnectMockServer creates a connection to a deployed mockserver in the environment
func ConnectMockServer(e *environment.Environment) (*MockserverClient, error) {
	localURL, err := e.Charts.Connections("mockserver").LocalURLByPort("serviceport", environment.HTTP)
	if err != nil {
		return nil, err
	}
	remoteURL, err := e.Config.Charts.Connections("mockserver").RemoteURLByPort("serviceport", environment.HTTP)
	if err != nil {
		return nil, err
	}
	c := NewMockserverClient(&MockserverConfig{
		LocalURL:   localURL.String(),
		ClusterURL: remoteURL.String(),
	})
	return c, nil
}

// ConnectMockServerSoak creates a connection to a deployed mockserver, assuming runner is in a soak test runner
func ConnectMockServerSoak(e *environment.Environment) (*MockserverClient, error) {
	remoteURL, err := e.Config.Charts.Connections("mockserver").RemoteURLByPort("serviceport", environment.HTTP)
	if err != nil {
		return nil, err
	}
	c := NewMockserverClient(&MockserverConfig{
		LocalURL:   remoteURL.String(),
		ClusterURL: remoteURL.String(),
	})
	return c, nil
}

// NetworkRegistry holds all the registered network types that can be initialized, allowing
// external libraries to register alternative network types to use
type NetworkRegistry struct {
	registeredNetworks map[string]registeredNetwork
}

type registeredNetwork struct {
	newBlockchainClientFn NewBlockchainClientFn
	blockchainClientURLFn BlockchainClientURLFn
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
		},
	}
}

// RegisterNetwork registers a new type of network within the registry
func (n *NetworkRegistry) RegisterNetwork(networkType string, fn NewBlockchainClientFn, urlFn BlockchainClientURLFn) {
	n.registeredNetworks[networkType] = registeredNetwork{
		newBlockchainClientFn: fn,
		blockchainClientURLFn: urlFn,
	}
}

// GetNetworks returns a networks object with all the BlockchainClient(s) initialized
func (n *NetworkRegistry) GetNetworks(env *environment.Environment) (*Networks, error) {
	nc := config.ProjectNetworkSettings
	var clients []BlockchainClient
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
	var defaultClient BlockchainClient
	if len(clients) >= 1 {
		defaultClient = clients[0]
	}
	return &Networks{
		clients: clients,
		Default: defaultClient,
	}, nil
}

// ConnectChainlinkNodes creates new chainlink clients
func ConnectChainlinkNodes(e *environment.Environment) ([]Chainlink, error) {
	return ConnectChainlinkNodesByCharts(e, []string{"chainlink"})
}

func ConnectChainlinkDBs(e *environment.Environment) ([]*PostgresConnector, error) {
	return ConnectChainlinkDBByCharts(e, []string{"chainlink"})
}

// ConnectChainlinkDBByCharts creates new chainlink DBs clients by charts
func ConnectChainlinkDBByCharts(e *environment.Environment, charts []string) ([]*PostgresConnector, error) {
	var dbs []*PostgresConnector
	for _, chart := range charts {
		pgUrls, err := e.Charts.Connections(chart).LocalURLsByPort("postgres", environment.HTTP)
		if err != nil {
			return nil, err
		}
		for _, u := range pgUrls {
			c, err := NewPostgresConnector(&PostgresConfig{
				Host:     "localhost",
				Port:     u.Port(),
				User:     "postgres",
				Password: "node",
				DBName:   "chainlink",
			})
			dbs = append(dbs, c)
			if err != nil {
				return nil, err
			}
		}
	}
	return dbs, nil
}

// ConnectChainlinkNodesByCharts creates new chainlink clients by charts
func ConnectChainlinkNodesByCharts(e *environment.Environment, charts []string) ([]Chainlink, error) {
	var clients []Chainlink

	for _, chart := range charts {
		localURLs, err := e.Charts.Connections(chart).LocalURLsByPort("access", environment.HTTP)
		if err != nil {
			return nil, err
		}
		remoteURLs, err := e.Charts.Connections(chart).RemoteURLsByPort("access", environment.HTTP)
		if err != nil {
			return nil, err
		}
		for urlIndex, localURL := range localURLs {
			c, err := NewChainlink(&ChainlinkConfig{
				URL:      localURL.String(),
				Email:    "notreal@fakeemail.ch",
				Password: "twochains",
				RemoteIP: remoteURLs[urlIndex].Hostname(),
			}, http.DefaultClient)
			clients = append(clients, c)
			if err != nil {
				return nil, err
			}
		}
	}
	return clients, nil
}

// ConnectChainlinkNodesSoak assumes that the tests are being run from an internal soak test runner
func ConnectChainlinkNodesSoak(e *environment.Environment) ([]Chainlink, error) {
	var clients []Chainlink

	remoteURLs, err := e.Charts.Connections("chainlink").RemoteURLsByPort("access", environment.HTTP)
	if err != nil {
		return nil, err
	}
	for urlIndex, localURL := range remoteURLs {
		c, err := NewChainlink(&ChainlinkConfig{
			URL:      localURL.String(),
			Email:    "notreal@fakeemail.ch",
			Password: "twochains",
			RemoteIP: remoteURLs[urlIndex].Hostname(),
		}, http.DefaultClient)
		clients = append(clients, c)
		if err != nil {
			return nil, err
		}
	}
	return clients, nil
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
