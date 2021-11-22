// Package client handles connections between chainlink nodes and different blockchain networks
package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/helmenv/environment"
	"github.com/smartcontractkit/integrations-framework/utils"
	"gopkg.in/yaml.v2"
	"math/big"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/smartcontractkit/integrations-framework/config"
)

// Commonly used blockchain network types
const (
	ETHNetworkType         = "eth_multinode"
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
	LoadWallets(ns interface{}) error
	SetWallet(num int) error

	CalculateTXSCost(txs int64) (*big.Float, error)
	CalculateTxGas(gasUsedValue *big.Int) (*big.Float, error)

	Get() interface{}
	GetNetworkName() string
	SwitchNode(node int) error
	GetClients() []BlockchainClient
	SuggestGasPrice(ctx context.Context) (*big.Int, error)
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
	if len(b.clients) >= index {
		return fmt.Errorf("index of %d is out of bounds", index)
	}
	b.Default = b.clients[index]
	return nil
}

// Get gets blockchain network (client) by name
func (b *Networks) Get(index int) (BlockchainClient, error) {
	if len(b.clients) >= index {
		return nil, fmt.Errorf("index of %d is out of bounds", index)
	}
	return b.clients[index], nil
}

// NewMockServerClientFromEnv creates new mockserver from env
func NewMockServerClientFromEnv(e *environment.Environment) (*MockserverClient, error) {
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

// NetworkRegistry holds all the registered network types that can be initialised, allowing
// external libraries to register alternative network types to use
type NetworkRegistry struct {
	registeredNetworks map[string]registeredNetwork
}

type registeredNetwork struct {
	newBlockchainClientFn NewBlockchainClientFn
	blockchainClientURLFn BlockchainClientURLFn
}

// NewNetworkRegistry returns an instance of the network registry with the default supported networks registered
func NewNetworkRegistry() *NetworkRegistry {
	return &NetworkRegistry{
		registeredNetworks: map[string]registeredNetwork{
			ETHNetworkType: {
				newBlockchainClientFn: NewEthereumMultiNodeClient,
				blockchainClientURLFn: EthereumMultiNodeURLs,
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

// GetNetworks returns a networks object with all the BlockchainClient(s) initialised
func (n *NetworkRegistry) GetNetworks(env *environment.Environment) (*Networks, error) {
	nc, err := config.LoadNetworksConfig(filepath.Join(utils.ProjectRoot, "networks.yaml"))
	if err != nil {
		return nil, err
	}
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
	if len(clients) == 1 {
		for _, c := range clients {
			defaultClient = c
		}
	}
	return &Networks{
		clients: clients,
		Default: defaultClient,
	}, nil
}

// NewChainlinkClients creates new chainlink clients
func NewChainlinkClients(e *environment.Environment) ([]Chainlink, error) {
	var clients []Chainlink

	urls, err := e.Charts.Connections("chainlink").LocalURLsByPort("access", environment.HTTP)
	if err != nil {
		return nil, err
	}
	for _, chainlinkURL := range urls {
		c, err := NewChainlink(&ChainlinkConfig{
			URL:      chainlinkURL.String(),
			Email:    "notreal@fakeemail.ch",
			Password: "twochains",
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

// UnmarshalNetworkConfig is a generic function to unmarshall a yaml map into a given object
func UnmarshalNetworkConfig(config map[string]interface{}, obj interface{}) error {
	b, err := json.Marshal(config)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(b, obj); err != nil {
		return err
	}
	return nil
}
