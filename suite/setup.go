package suite

import (
	"github.com/avast/retry-go"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/config"
	"github.com/smartcontractkit/integrations-framework/contracts"
	"math/big"
	"net/http"
	"time"
)

type DefaultSuiteSetup struct {
	Config   *config.Config
	Client   client.BlockchainClient
	Wallets  client.BlockchainWallets
	Deployer contracts.ContractDeployer
	Link     contracts.LinkToken
}

// DefaultLocalSetup setup minimum required components for test
func DefaultLocalSetup(initFunc client.BlockchainNetworkInit) (*DefaultSuiteSetup, error) {
	conf, err := config.NewWithPath(config.LocalConfig, "../../config")
	if err != nil {
		return nil, err
	}
	networkConfig, err := initFunc(conf)
	if err != nil {
		return nil, err
	}
	blockchainClient, err := client.NewBlockchainClient(networkConfig)
	if err != nil {
		return nil, err
	}
	wallets, err := networkConfig.Wallets()
	if err != nil {
		return nil, err
	}
	contractDeployer, err := contracts.NewContractDeployer(blockchainClient)
	if err != nil {
		return nil, err
	}
	link, err := contractDeployer.DeployLinkTokenContract(wallets.Default())
	if err != nil {
		return nil, err
	}
	// configure default retry
	retry.DefaultAttempts = conf.Retry.Attempts
	// linear waiting
	retry.DefaultDelayType = func(n uint, err error, config *retry.Config) time.Duration {
		return conf.Retry.LinearDelay
	}
	return &DefaultSuiteSetup{
		Config:   conf,
		Client:   blockchainClient,
		Wallets:  wallets,
		Deployer: contractDeployer,
		Link:     link,
	}, nil
}

// FundTemplateNodes funds Chainlink nodes with ETH/LINK
func FundTemplateNodes(blockchainClient client.BlockchainClient, wallets client.BlockchainWallets, nodes []client.Chainlink, ethAmount int64, linkAmount int64) error {
	for _, node := range nodes {
		nodeEthKeys, err := node.ReadETHKeys()
		if err != nil {
			return err
		}
		primaryEthKey := nodeEthKeys.Data[0]

		err = blockchainClient.Fund(
			wallets.Default(),
			primaryEthKey.Attributes.Address,
			big.NewInt(ethAmount), big.NewInt(linkAmount),
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// ConnectToTemplateNodes assumes that 5 template nodes are running locally, check out a quick setup for that here:
// https://github.com/smartcontractkit/chainlink-node-compose
func ConnectToTemplateNodes() ([]client.Chainlink, []common.Address, error) {
	urlBase := "http://localhost:"
	ports := []string{"6711", "6722", "6733", "6744", "6755"}
	// Checks if those nodes are actually up and healthy
	for _, port := range ports {
		_, err := http.Get(urlBase + port)
		if err != nil {
			log.Err(err).Str("URL", urlBase+port).Msg("Chainlink node unhealthy / not up. Make sure nodes are already up")
			return nil, nil, err
		}
		log.Info().Str("URL", urlBase+port).Msg("Chainlink Node Healthy")
	}

	var cls []client.Chainlink
	var clsAddresses []common.Address
	for _, port := range ports {
		c := &client.ChainlinkConfig{
			URL:      urlBase + port,
			Email:    "notreal@fakeemail.ch",
			Password: "twochains",
		}
		cl, err := client.NewChainlink(c, http.DefaultClient)
		if err != nil {
			return nil, nil, err
		}
		cls = append(cls, cl)
		nodeEthKeys, err := cl.ReadETHKeys()
		if err != nil {
			log.Err(err).Str("Node URL", urlBase+port).Msg("Issue establishing connection to node")
		}
		primaryEthKey := nodeEthKeys.Data[0]
		ethAddress := primaryEthKey.Attributes.Address
		clsAddresses = append(clsAddresses, common.HexToAddress(ethAddress))
	}

	return cls, clsAddresses, nil
}
