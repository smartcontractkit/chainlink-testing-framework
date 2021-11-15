// Package client handles connections between chainlink nodes and different blockchain networks
package client

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/helmenv/environment"
	"github.com/smartcontractkit/integrations-framework/utils"
	"gopkg.in/yaml.v2"
	"math/big"
	"path/filepath"

	"github.com/smartcontractkit/integrations-framework/config"
)

// Commonly used blockchain network types
const (
	ETHNetworkType         = "eth_multinode"
	NetworkGethPerformance = "ethereum_geth_performance"
)

// ExternalClientImplFunc external client implementation function
// networkName must match a key in "networks" in networks.yaml config
// networkConfig is just an arbitrary config you provide in "networks" for your key
type ExternalClientImplFunc func(networkName string, networkConfig map[string]interface{}, e *environment.Environment) (BlockchainClient, error)

// BlockchainClient is the interface that wraps a given client implementation for a blockchain, to allow for switching
// of network types within the test suite
// BlockchainClient can be connected to a single or multiple nodes,
type BlockchainClient interface {
	LoadWallets(ns interface{}) error
	Get() interface{}
	GetNetworkName() string
	GetID() int
	SetID(id int)
	SwitchNode(node int) error
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

// Networks is a thin wrapper that just selects client connected to some network
// if there is only one client it is chosen as Default
// if there is multiple you just get clients you need in test
type Networks struct {
	clients map[string]BlockchainClient
	Default BlockchainClient
}

func (b *Networks) SetDefault(name string) error {
	if _, ok := b.clients[name]; !ok {
		return fmt.Errorf("no client found for network %s", name)
	}
	return nil
}

func (b *Networks) Get(name string) (BlockchainClient, error) {
	c, ok := b.clients[name]
	if !ok {
		return nil, fmt.Errorf("no client found for network %s", name)
	}
	return c, nil
}

// NewNetworks trying to create blockchain networks using "selected_networks" from networks config
// each network may have only one client instance, some clients can connect to multiple nodes but should implement
// BlockchainClient interface
// If there is no implementation found it will be created from ExternalClientImplFunc
// If network is deployed in the current environment URL got overridden from env
// If network is external (deployed somewhere else) we just use url from config
func NewNetworks(env *environment.Environment, extClients map[string]ExternalClientImplFunc) (*Networks, error) {
	nc, err := config.LoadNetworksConfig(filepath.Join(utils.ProjectRoot, "networks.yaml"))
	if err != nil {
		return nil, err
	}
	clients := map[string]BlockchainClient{}
	for _, networkName := range nc.SelectedNetworks {
		ns := nc.NetworkSettings[networkName].(map[string]interface{})
		networkType := ns["type"].(string)
		switch networkType {
		case ETHNetworkType:
			d, err := yaml.Marshal(ns)
			if err != nil {
				return nil, err
			}
			var cfg *config.ETHNetwork
			if err := yaml.Unmarshal(d, &cfg); err != nil {
				return nil, err
			}
			if !cfg.External {
				if _, ok := env.Config.NetworksURLs[networkName]; !ok {
					return nil, fmt.Errorf("network %s is not found in environment URLs", networkName)
				}
				cfg.URLs = env.Config.NetworksURLs[networkName]
			}
			cfg.ID = networkName
			ec, err := NewEthereumMultiNodeClient(cfg)
			if err != nil {
				return nil, err
			}
			clients[networkName] = ec
		default:
			log.Info().
				Str("Name", networkName).
				Msg("Creating client using external implementation")
			f, ok := extClients[networkName]
			if !ok {
				return nil, fmt.Errorf("client implementation for network %s is not provided, client func is nil", networkName)
			}
			ec, err := f(networkName, ns, env)
			if err != nil {
				return nil, err
			}
			clients[networkName] = ec
		}
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
