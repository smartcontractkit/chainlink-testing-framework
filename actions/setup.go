package actions

import (
	"strings"
	"time"

	"github.com/onsi/ginkgo"
	"github.com/rs/zerolog/log"

	"github.com/avast/retry-go"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/config"
	"github.com/smartcontractkit/integrations-framework/contracts"
	"github.com/smartcontractkit/integrations-framework/environment"
)

const (
	KeepEnvironmentsNever  = "never"
	KeepEnvironmentsOnFail = "onfail"
	KeepEnvironmentsAlways = "always"
)

type DefaultSuiteSetup struct {
	Config   *config.Config
	Client   client.BlockchainClient
	Wallets  client.BlockchainWallets
	Deployer contracts.ContractDeployer
	Link     contracts.LinkToken
	Env      environment.Environment
}

// DefaultLocalSetup setup minimum required components for test
func DefaultLocalSetup(
	envInitFunc environment.K8sEnvSpecInit,
	initFunc client.BlockchainNetworkInit,
) (*DefaultSuiteSetup, error) {
	conf, err := config.NewConfig(config.LocalConfig)
	if err != nil {
		return nil, err
	}
	networkConfig, err := initFunc(conf)
	if err != nil {
		return nil, err
	}
	env, err := environment.NewK8sEnvironment(envInitFunc, conf, networkConfig)
	if err != nil {
		return nil, err
	}
	blockchainClient, err := environment.NewBlockchainClient(env, networkConfig)
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
		Env:      env,
	}, nil
}

func (s *DefaultSuiteSetup) TearDown() func() {
	return func() {
		switch strings.ToLower(s.Config.KeepEnvironments) {
		case KeepEnvironmentsNever:
			s.Env.TearDown()
		case KeepEnvironmentsOnFail:
			if !ginkgo.CurrentGinkgoTestDescription().Failed {
				s.Env.TearDown()
			} else {
				log.Info().Str("Namespace", s.Env.ID()).Msg("Kept environment due to test failure")
			}
		case KeepEnvironmentsAlways:
			log.Info().Str("Namespace", s.Env.ID()).Msg("Kept environment")
			return
		default:
			s.Env.TearDown()
		}
	}
}
