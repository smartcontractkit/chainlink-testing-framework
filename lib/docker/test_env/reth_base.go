package test_env

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	config_types "github.com/smartcontractkit/chainlink-testing-framework/lib/config/types"

	"github.com/rs/zerolog"

	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/config"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/docker"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
)

type Reth struct {
	EnvComponent
	ExternalHttpUrl      string
	InternalHttpUrl      string
	ExternalWsUrl        string
	InternalWsUrl        string
	InternalExecutionURL string
	ExternalExecutionURL string
	chainConfig          *config.EthereumChainConfig
	consensusLayer       config.ConsensusLayer
	ethereumVersion      config_types.EthereumVersion
	l                    zerolog.Logger
	t                    *testing.T
	posContainerSettings
}

// WithTestInstance sets up the execution client with a test logger and the provided testing context.
// This is useful for integrating testing frameworks with the execution client, enabling better logging and error tracking during tests.
func (g *Reth) WithTestInstance(t *testing.T) ExecutionClient {
	g.l = logging.GetTestLogger(t)
	g.t = t
	return g
}

// StartContainer initializes and starts a Reth execution client container.
// It returns the configured EVM network details or an error if the process fails.
// This function is essential for setting up a local Ethereum environment for testing or development.
func (g *Reth) StartContainer() (blockchain.EVMNetwork, error) {
	if g.GetEthereumVersion() == config_types.EthereumVersion_Eth1 {
		return blockchain.EVMNetwork{}, errors.New(config.Eth1NotSupportedByRethMsg)
	}

	r, err := g.getEth2ContainerRequest()

	if err != nil {
		return blockchain.EVMNetwork{}, err
	}

	l := logging.GetTestContainersGoTestLogger(g.t)
	ct, err := docker.StartContainerWithRetry(g.l, tc.GenericContainerRequest{
		ContainerRequest: *r,
		Reuse:            g.WasRecreated,
		Started:          true,
		Logger:           l,
	})
	if err != nil {
		return blockchain.EVMNetwork{}, fmt.Errorf("cannot start reth container: %w", err)
	}

	host, err := GetHost(testcontext.Get(g.t), ct)
	if err != nil {
		return blockchain.EVMNetwork{}, err
	}
	httpPort, err := ct.MappedPort(testcontext.Get(g.t), NatPort(DEFAULT_EVM_NODE_HTTP_PORT))
	if err != nil {
		return blockchain.EVMNetwork{}, err
	}

	wsPort, err := ct.MappedPort(testcontext.Get(g.t), NatPort(DEFAULT_EVM_NODE_WS_PORT))
	if err != nil {
		return blockchain.EVMNetwork{}, err
	}

	if g.GetEthereumVersion() == config_types.EthereumVersion_Eth2 {
		executionPort, err := ct.MappedPort(testcontext.Get(g.t), NatPort(ETH2_EXECUTION_PORT))
		if err != nil {
			return blockchain.EVMNetwork{}, err
		}
		g.InternalExecutionURL = FormatHttpUrl(g.ContainerName, ETH2_EXECUTION_PORT)
		g.ExternalExecutionURL = FormatHttpUrl(host, executionPort.Port())
	}

	g.Container = ct
	g.ExternalHttpUrl = FormatHttpUrl(host, httpPort.Port())
	g.InternalHttpUrl = FormatHttpUrl(g.ContainerName, DEFAULT_EVM_NODE_HTTP_PORT)
	g.ExternalWsUrl = FormatWsUrl(host, wsPort.Port())
	g.InternalWsUrl = FormatWsUrl(g.ContainerName, DEFAULT_EVM_NODE_WS_PORT)

	networkConfig := blockchain.SimulatedEVMNetwork
	if g.GetEthereumVersion() == config_types.EthereumVersion_Eth1 {
		networkConfig.Name = fmt.Sprintf("Private Eth-1-PoW [reth %s]", g.ContainerVersion)
	} else {
		networkConfig.Name = fmt.Sprintf("Private Eth-2-PoS [reth %s] + %s", g.ContainerVersion, g.consensusLayer)
	}
	networkConfig.URLs = []string{g.ExternalWsUrl}
	networkConfig.HTTPURLs = []string{g.ExternalHttpUrl}
	networkConfig.SimulationType = "Reth"

	g.l.Info().Str("containerName", g.ContainerName).
		Msg("Started Reth container")

	return networkConfig, nil
}

// GetInternalExecutionURL returns the internal execution URL for the Ethereum client.
// It is used to retrieve the endpoint for interacting with the execution layer,
// ensuring compatibility with Ethereum 2.0 clients.
func (g *Reth) GetInternalExecutionURL() string {
	if g.GetEthereumVersion() == config_types.EthereumVersion_Eth1 {
		panic("eth1 node doesn't have an execution URL")
	}
	return g.InternalExecutionURL
}

// GetExternalExecutionURL returns the external execution URL for the Reth instance.
// It panics if the Ethereum version is Eth1, as Eth1 nodes do not support execution URLs.
func (g *Reth) GetExternalExecutionURL() string {
	if g.GetEthereumVersion() == config_types.EthereumVersion_Eth1 {
		panic("eth1 node doesn't have an execution URL")
	}
	return g.ExternalExecutionURL
}

// GetInternalHttpUrl returns the internal HTTP URL of the execution client.
// This URL is essential for establishing communication with the client in a secure manner.
func (g *Reth) GetInternalHttpUrl() string {
	return g.InternalHttpUrl
}

// GetInternalWsUrl returns the internal WebSocket URL for the Reth execution client.
// This URL is essential for establishing WebSocket connections to the client for real-time data streaming.
func (g *Reth) GetInternalWsUrl() string {
	return g.InternalWsUrl
}

// GetExternalHttpUrl returns the external HTTP URL for the Reth execution client.
// This URL is useful for connecting to the client from external applications or services.
func (g *Reth) GetExternalHttpUrl() string {
	return g.ExternalHttpUrl
}

// GetExternalWsUrl returns the external WebSocket URL for the Reth execution client.
// This URL is essential for connecting to the client for real-time data and event subscriptions.
func (g *Reth) GetExternalWsUrl() string {
	return g.ExternalWsUrl
}

// GetContainerName returns the name of the container associated with the Reth instance.
// This function is useful for identifying and managing the container in a Docker environment.
func (g *Reth) GetContainerName() string {
	return g.ContainerName
}

// GetContainer returns a pointer to the container associated with the Reth instance.
// This function is useful for accessing the container's properties and methods in order to manage or interact with the execution environment.
func (g *Reth) GetContainer() *tc.Container {
	return &g.Container
}

// GetEthereumVersion returns the current Ethereum version of the Reth instance.
// This information is essential for determining compatibility and functionality
// with various Ethereum features and services.
func (g *Reth) GetEthereumVersion() config_types.EthereumVersion {
	return g.ethereumVersion
}

// WaitUntilChainIsReady blocks until the Ethereum chain is ready for use, waiting for the first block to be committed if necessary.
// This function is essential for ensuring that the blockchain environment is fully operational before proceeding with further operations.
func (g *Reth) WaitUntilChainIsReady(ctx context.Context, waitTime time.Duration) error {
	if g.GetEthereumVersion() == config_types.EthereumVersion_Eth1 {
		return nil
	}
	waitForFirstBlock := tcwait.NewLogStrategy("Canonical chain committed").WithPollInterval(1 * time.Second).WithStartupTimeout(waitTime)
	return waitForFirstBlock.WaitUntilReady(ctx, *g.GetContainer())
}

// GetConsensusMechanism returns the consensus mechanism used by the Ethereum network.
// It identifies whether the network is operating on Proof of Work (PoW) or Proof of Stake (PoS)
// based on the current Ethereum version, aiding in understanding network dynamics.
func (g *Reth) GetConsensusMechanism() ConsensusMechanism {
	if g.GetEthereumVersion() == config_types.EthereumVersion_Eth1 {
		return ConsensusMechanism_PoW
	}
	return ConsensusMechanism_PoS
}
