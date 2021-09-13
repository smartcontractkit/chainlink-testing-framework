package actions

import (
	"os"
	"path/filepath"
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

// Keep Environments options
const (
	KeepEnvironmentsNever  = "never"
	KeepEnvironmentsOnFail = "onfail"
	KeepEnvironmentsAlways = "always"
)

// DefaultSuiteSetup holds the data for a default setup
type DefaultSuiteSetup struct {
	Config   *config.Config
	Client   client.BlockchainClient
	Wallets  client.BlockchainWallets
	Deployer contracts.ContractDeployer
	Link     contracts.LinkToken
	Env      environment.Environment
	Network  client.BlockchainNetwork
}

// DefaultLocalSetup setup minimum required components for test
func DefaultLocalSetup(
	envInitFunc environment.K8sEnvSpecInit,
	initFunc client.BlockchainNetworkInit,
	configPath string,
) (*DefaultSuiteSetup, error) {
	conf, err := config.NewConfig(configPath)
	if err != nil {
		return nil, err
	}
	network, err := initFunc(conf)
	if err != nil {
		return nil, err
	}
	env, err := environment.NewK8sEnvironment(envInitFunc, conf, network)
	if err != nil {
		return nil, err
	}
	blockchainClient, err := environment.NewBlockchainClient(env, network)
	if err != nil {
		return nil, err
	}
	wallets, err := network.Wallets()
	if err != nil {
		return nil, err
	}
	contractDeployer, err := contracts.NewContractDeployer(blockchainClient)
	if err != nil {
		return nil, err
	}
	if err := contracts.AwaitMining(blockchainClient); err != nil {
		return nil, err
	}
	link, err := contractDeployer.DeployLinkTokenContract(wallets.Default())
	if err != nil {
		return nil, err
	}
	// configure default retry
	retry.DefaultAttempts = conf.Retry.Attempts
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
		Network:  network,
	}, nil
}

// TearDown checks for test failure, writes logs if there is one, then tears down the test environment, based on the
// keep_environments config value
func (s *DefaultSuiteSetup) TearDown() func() {
	if ginkgo.CurrentGinkgoTestDescription().Failed { // If a test fails, dump logs
		logsFolder := filepath.Join(s.Config.ConfigFileLocation, "/logs/")
		if _, err := os.Stat(logsFolder); os.IsNotExist(err) {
			if err = os.Mkdir(logsFolder, 0755); err != nil {
				log.Err(err).Str("Log Folder", logsFolder).Msg("Error creating logs directory")
			}
		}
		testLogFolder := filepath.Join(logsFolder, strings.Replace(ginkgo.CurrentGinkgoTestDescription().TestText, " ", "-", -1)+
			"_"+s.Env.ID()+"/")
		// Create specific test folder
		if _, err := os.Stat(testLogFolder); os.IsNotExist(err) {
			if err = os.Mkdir(testLogFolder, 0755); err != nil {
				log.Err(err).Str("Log Folder", testLogFolder).Msg("Error creating logs directory")
			}
		}

		s.Env.WriteArtifacts(testLogFolder)
		log.Info().Str("Log Folder", testLogFolder).Msg("Wrote environment logs")
	}
	return func() {
		if err := s.Client.Close(); err != nil {
			log.Error().
				Str("Network", s.Config.Network).
				Msgf("Error while closing the Blockchain client: %v", err)
		}

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
