package actions

import (
	"errors"
	"fmt"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/config"
	"github.com/smartcontractkit/integrations-framework/contracts"
)

// Those hooks are just for pure compatibility to run some tests in integration-framework repo,
// remove this when we have new "hooks" implementation in the main repo

// EthereumDeployerHook deployer hook
func EthereumDeployerHook(bcClient client.BlockchainClient) (contracts.ContractDeployer, error) {
	return contracts.NewEthereumContractDeployer(bcClient.Get().(*client.EthereumClient)), nil
}

// EthereumClientHook client hook
func EthereumClientHook(network client.BlockchainNetwork) (client.BlockchainClient, error) {
	return client.NewEthereumClient(network)
}

// EVMNetworkFromConfigHook evm network from config hook
func EVMNetworkFromConfigHook(config *config.Config) (client.BlockchainNetwork, error) {
	firstNetwork := config.Networks[0]
	return client.NewEthereumNetwork(firstNetwork, config.NetworkConfigs[firstNetwork])
}

// NetworksFromConfigHook networks from config hook
func NetworksFromConfigHook(config *config.Config) ([]client.BlockchainNetwork, error) {
	networks := make([]client.BlockchainNetwork, 0)
	if len(config.NetworkConfigs) < 2 {
		return nil, errors.New("at least 2 evm networks are required")
	}
	for _, networkName := range config.Networks {
		if _, ok := config.NetworkConfigs[networkName]; !ok {
			return nil, fmt.Errorf("'%s' is not a supported network name. Check the network configs in you config file", networkName)
		}
		net, err := client.NewEthereumNetwork(networkName, config.NetworkConfigs[networkName])
		if err != nil {
			return nil, err
		}
		networks = append(networks, net)
	}
	return networks, nil
}

// EthereumPerfNetworkHook perf network func
func EthereumPerfNetworkHook(config *config.Config) (client.BlockchainNetwork, error) {
	return client.NewEthereumNetwork("ethereum_geth_performance", config.NetworkConfigs["ethereum_geth_performance"])
}
