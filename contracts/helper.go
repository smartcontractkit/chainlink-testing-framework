package contracts

import (
	"fmt"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/config"
)

type DefaultSuiteSetup struct {
	Config   *config.Config
	Client   client.BlockchainClient
	Wallets  client.BlockchainWallets
	Deployer ContractDeployer
	Link     LinkToken
}

var (
	observationSourceTmpl = `fetch    [type=http method=POST url="%s" requestData="{}"];
			parse    [type=jsonparse path="data,result"];    
			fetch -> parse;`
)

func ObservationSourceSpec(url string) string {
	return fmt.Sprintf(observationSourceTmpl, url)
}

func DefaultLocalSetup(initFunc client.BlockchainNetworkInit) (*DefaultSuiteSetup, error) {
	conf, err := config.NewWithPath(config.LocalConfig, "../config")
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
	contractDeployer, err := NewContractDeployer(blockchainClient)
	if err != nil {
		return nil, err
	}
	link, err := contractDeployer.DeployLinkTokenContract(wallets.Default())
	if err != nil {
		return nil, err
	}
	return &DefaultSuiteSetup{
		Config:   conf,
		Client:   blockchainClient,
		Wallets:  wallets,
		Deployer: contractDeployer,
		Link:     link,
	}, nil
}
