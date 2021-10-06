package actions

import (
	"fmt"
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

// NetworkInfo helps delineate network information in a multi-network setup
type NetworkInfo struct {
	Client   client.BlockchainClient
	Wallets  client.BlockchainWallets
	Deployer contracts.ContractDeployer
	Link     contracts.LinkToken
	Network  client.BlockchainNetwork
}

// buildNetworkInfo initializes the network's blockchain client and gathers all test-relevant network information
func buildNetworkInfo(network client.BlockchainNetwork, env environment.Environment) (NetworkInfo, error) {
	// Initialize blockchain client
	var bcc client.BlockchainClient
	var err error
	switch network.Config().Type {
	case client.BlockchainTypeEVMMultinode:
		bcc, err = environment.NewBlockchainClients(env, network)
	case client.BlockchainTypeEVM:
		bcc, err = environment.NewBlockchainClient(env, network)
	}
	if err != nil {
		return NetworkInfo{}, err
	}

	// Initialize wallets
	wallets, err := network.Wallets()
	if err != nil {
		return NetworkInfo{}, err
	}
	contractDeployer, err := contracts.NewContractDeployer(bcc)
	if err != nil {
		return NetworkInfo{}, err
	}
	link, err := contractDeployer.DeployLinkTokenContract(wallets.Default())
	if err != nil {
		return NetworkInfo{}, err
	}
	return NetworkInfo{
		Client:   bcc,
		Wallets:  wallets,
		Deployer: contractDeployer,
		Link:     link,
		Network:  network,
	}, nil
}

// SuiteSetup enables common use cases, and safe handling of different blockchain networks for test scenarios
type SuiteSetup interface {
	Config() *config.Config
	Environment() environment.Environment

	DefaultNetwork() NetworkInfo
	Network(networkID string) (NetworkInfo, error)
	Networks(networkID string) ([]NetworkInfo, error)

	TearDown() func()
}

// SingleNetworkSuiteSetup holds the data for a default setup
type SingleNetworkSuiteSetup struct {
	config  *config.Config
	env     environment.Environment
	network NetworkInfo
}

// SingleNetworkSetup setup minimum required components for test
func SingleNetworkSetup(
	initialDeployInitFunc environment.K8sEnvSpecInit,
	initFunc client.BlockchainNetworkInit,
	configPath string,
) (SuiteSetup, error) {
	conf, err := config.NewConfig(configPath)
	if err != nil {
		return nil, err
	}
	network, err := initFunc(conf)
	if err != nil {
		return nil, err
	}

	env, err := environment.NewK8sEnvironment(conf, network)
	if err != nil {
		return nil, err
	}
	err = env.DeploySpecs(initialDeployInitFunc)
	if err != nil {
		return nil, err
	}

	networkInfo, err := buildNetworkInfo(network, env)
	if err != nil {
		return nil, err
	}

	// configure default retry
	retry.DefaultAttempts = conf.Retry.Attempts
	retry.DefaultDelayType = func(n uint, err error, config *retry.Config) time.Duration {
		return conf.Retry.LinearDelay
	}

	return &SingleNetworkSuiteSetup{
		config:  conf,
		env:     env,
		network: networkInfo,
	}, nil
}

// Config retrieves the general config for the suite
func (s *SingleNetworkSuiteSetup) Config() *config.Config {
	return s.config
}

// Environment retrieves the general environment for the suite
func (s *SingleNetworkSuiteSetup) Environment() environment.Environment {
	return s.env
}

// DefaultNetwork returns the only network in a single network environment
func (s *SingleNetworkSuiteSetup) DefaultNetwork() NetworkInfo {
	return s.network
}

// Network returns the only network in a single network environment
func (s *SingleNetworkSuiteSetup) Network(networkID string) (NetworkInfo, error) {
	return s.network, nil
}

// Networks returns the only network in a single network environment
func (s *SingleNetworkSuiteSetup) Networks(networkID string) ([]NetworkInfo, error) {
	return []NetworkInfo{s.network}, nil
}

// TearDown checks for test failure, writes logs if there is one, then tears down the test environment, based on the
// keep_environments config value
func (s *SingleNetworkSuiteSetup) TearDown() func() {
	return teardown(*s.config, s.env, s.network.Client)
}

// multiNetworkSuiteSetup holds the data for a multiple network setup
type multiNetworkSuiteSetup struct {
	config   *config.Config
	env      environment.Environment
	networks []NetworkInfo
}

// MultiNetworkSetup enables testing across multiple networks
func MultiNetworkSetup(
	initialDeployInitFunc environment.K8sEnvSpecInit,
	multiNetworkInitialization client.MultiNetworkInit,
	configPath string,
) (SuiteSetup, error) {
	conf, err := config.NewConfig(configPath)
	if err != nil {
		return nil, err
	}
	networks, err := multiNetworkInitialization(conf)
	if err != nil {
		return nil, err
	}

	env, err := environment.NewK8sEnvironment(conf, networks...)
	if err != nil {
		return nil, err
	}

	err = env.DeploySpecs(initialDeployInitFunc)
	if err != nil {
		return nil, err
	}

	allNetworks := make([]NetworkInfo, len(networks))
	for index, network := range networks {
		networkInfo, err := buildNetworkInfo(network, env)
		if err != nil {
			return nil, err
		}
		allNetworks[index] = networkInfo
	}

	// configure default retry
	retry.DefaultAttempts = conf.Retry.Attempts
	retry.DefaultDelayType = func(n uint, err error, config *retry.Config) time.Duration {
		return conf.Retry.LinearDelay
	}
	return &multiNetworkSuiteSetup{
		config:   conf,
		env:      env,
		networks: allNetworks,
	}, nil
}

// Config retrieves the general config for the suite
func (s *multiNetworkSuiteSetup) Config() *config.Config {
	return s.config
}

// Environment retrieves the general environment for the suite
func (s *multiNetworkSuiteSetup) Environment() environment.Environment {
	return s.env
}

// DefaultNetwork returns the network information for the first / only network in the suite
func (s *multiNetworkSuiteSetup) DefaultNetwork() NetworkInfo {
	return s.networks[0]
}

// Network returns the network information for the network with the supplied ID. If there is more than 1 network with
// that ID, the first one encountered is returned.
func (s *multiNetworkSuiteSetup) Network(networkID string) (NetworkInfo, error) {
	networkIDs := make([]string, 0)
	for _, network := range s.networks {
		networkIDs = append(networkIDs, network.Client.GetName())
		if network.Client.GetName() == networkID {
			return network, nil
		}
	}
	return NetworkInfo{}, fmt.Errorf("Unable to find any networks with the ID '%s'. All found networks: %v", networkID, networkIDs)
}

// Networks returns the network information for all the networks with the supplied ID.
func (s *multiNetworkSuiteSetup) Networks(networkID string) ([]NetworkInfo, error) {
	networkIDs := make([]string, 0)
	networks := make([]NetworkInfo, 0)
	for _, network := range s.networks {
		networkIDs = append(networkIDs, network.Client.GetName())
		if network.Client.GetName() == networkID {
			networks = append(networks, network)
		}
	}
	if len(networks) == 0 {
		return nil, fmt.Errorf("Unable to find any networks with the ID '%s'. All found networks: %v", networkID, networkIDs)
	} else {
		return networks, nil
	}
}

// TearDown checks for test failure, writes logs if there is one, then tears down the test environment, based on the
// keep_environments config value
func (s *multiNetworkSuiteSetup) TearDown() func() {
	clients := make([]client.BlockchainClient, len(s.networks))
	for index, network := range s.networks {
		clients[index] = network.Client
	}
	return teardown(*s.config, s.env, clients...)
}

// TearDown checks for test failure, writes logs if there is one, then tears down the test environment, based on the
// keep_environments config value
func teardown(config config.Config, env environment.Environment, clients ...client.BlockchainClient) func() {
	if ginkgo.CurrentGinkgoTestDescription().Failed { // If a test fails, dump logs
		logsFolder := filepath.Join(config.ConfigFileLocation, "/logs/")
		if _, err := os.Stat(logsFolder); os.IsNotExist(err) {
			if err = os.Mkdir(logsFolder, 0755); err != nil {
				log.Err(err).Str("Log Folder", logsFolder).Msg("Error creating logs directory")
			}
		}
		testLogFolder := filepath.Join(logsFolder, strings.Replace(ginkgo.CurrentGinkgoTestDescription().TestText, " ", "-", -1)+
			"_"+env.ID()+"/")
		// Create specific test folder
		if _, err := os.Stat(testLogFolder); os.IsNotExist(err) {
			if err = os.Mkdir(testLogFolder, 0755); err != nil {
				log.Err(err).Str("Log Folder", testLogFolder).Msg("Error creating logs directory")
			}
		}

		env.WriteArtifacts(testLogFolder)
		log.Info().Str("Log Folder", testLogFolder).Msg("Wrote environment logs")
	}
	return func() {
		for _, client := range clients {
			if err := client.Close(); err != nil {
				log.Err(err).
					Str("Network", client.GetName()).
					Msgf("Error while closing the Blockchain client")
			}
		}

		switch strings.ToLower(config.KeepEnvironments) {
		case KeepEnvironmentsNever:
			env.TearDown()
		case KeepEnvironmentsOnFail:
			if !ginkgo.CurrentGinkgoTestDescription().Failed {
				env.TearDown()
			} else {
				log.Info().Str("Namespace", env.ID()).Msg("Kept environment due to test failure")
			}
		case KeepEnvironmentsAlways:
			log.Info().Str("Namespace", env.ID()).Msg("Kept environment")
			return
		default:
			env.TearDown()
		}
	}
}
