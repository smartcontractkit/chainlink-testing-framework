package types

import (
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/config"
	"github.com/smartcontractkit/integrations-framework/contracts"
)

// NewDeployerHook external deployer function
type NewDeployerHook func(c client.BlockchainClient) (contracts.ContractDeployer, error)

// NewClientHook external client function
type NewClientHook func(network client.BlockchainNetwork) (client.BlockchainClient, error)

// NewNetworkHook is a helper function to obtain the network listed in the config file
type NewNetworkHook func(conf *config.Config) (client.BlockchainNetwork, error)

// NewMultinetworkHook is a helper function to create multiple blockchain networks at once
type NewMultinetworkHook func(conf *config.Config) ([]client.BlockchainNetwork, error)
