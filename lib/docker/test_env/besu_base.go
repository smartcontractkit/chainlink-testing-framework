package test_env

import (
	"context"
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

type powSettings struct {
	genesisPath  string
	rootPath     string
	keystorePath string
	accountAddr  string
}

type Besu struct {
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
	powSettings
}

// WithTestInstance sets up the execution client for testing by assigning a test logger and the testing context.
// This allows for better logging and error tracking during test execution.
func (g *Besu) WithTestInstance(t *testing.T) ExecutionClient {
	g.l = logging.GetTestLogger(t)
	g.t = t
	return g
}

// StartContainer initializes and starts a Besu container for Ethereum execution.
// It configures network settings based on the Ethereum version and returns the
// network configuration along with any errors encountered during the process.
func (g *Besu) StartContainer() (blockchain.EVMNetwork, error) {
	var r *tc.ContainerRequest
	var err error

	if g.GetEthereumVersion() == config_types.EthereumVersion_Eth1 {
		r, err = g.getEth1ContainerRequest()

	} else {
		r, err = g.getEth2ContainerRequest()
	}
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
		return blockchain.EVMNetwork{}, fmt.Errorf("cannot start Besu container: %w", err)
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
	networkConfig.URLs = []string{g.ExternalWsUrl}
	networkConfig.HTTPURLs = []string{g.ExternalHttpUrl}
	networkConfig.GasEstimationBuffer = 10_000_000_000
	networkConfig.SimulationType = "Besu"

	if g.GetEthereumVersion() == config_types.EthereumVersion_Eth1 {
		networkConfig.Name = fmt.Sprintf("Private Eth-1-PoA [besu %s]", g.ContainerVersion)
	} else {
		networkConfig.Name = fmt.Sprintf("Private Eth-2-PoS [besu %s] + %s", g.ContainerVersion, g.consensusLayer)
	}

	g.l.Info().Str("containerName", g.ContainerName).
		Msg("Started Besu container")

	return networkConfig, nil
}

// GetInternalExecutionURL returns the internal execution URL for the Besu client.
// It is used to retrieve the endpoint for executing transactions in Ethereum 2.0 networks,
// ensuring compatibility with the Ethereum version in use.
func (g *Besu) GetInternalExecutionURL() string {
	if g.GetEthereumVersion() == config_types.EthereumVersion_Eth1 {
		panic("eth1 node doesn't have an execution URL")
	}
	return g.InternalExecutionURL
}

// GetExternalExecutionURL returns the external execution URL for the Besu instance.
// It panics if the Ethereum version is Eth1, as Eth1 nodes do not support execution URLs.
func (g *Besu) GetExternalExecutionURL() string {
	if g.GetEthereumVersion() == config_types.EthereumVersion_Eth1 {
		panic("eth1 node doesn't have an execution URL")
	}
	return g.ExternalExecutionURL
}

// GetInternalHttpUrl returns the internal HTTP URL of the Besu client.
// This URL is essential for establishing communication with the Besu node in a private network setup.
func (g *Besu) GetInternalHttpUrl() string {
	return g.InternalHttpUrl
}

// GetInternalWsUrl returns the internal WebSocket URL for the Besu client.
// This URL is essential for establishing WebSocket connections to the Besu node for real-time data and event subscriptions.
func (g *Besu) GetInternalWsUrl() string {
	return g.InternalWsUrl
}

// GetExternalHttpUrl returns the external HTTP URL of the Besu client.
// This URL is used to interact with the Besu node from external applications or services.
func (g *Besu) GetExternalHttpUrl() string {
	return g.ExternalHttpUrl
}

// GetExternalWsUrl returns the external WebSocket URL for the Besu client.
// This URL is essential for connecting to the Besu node from external services or clients.
func (g *Besu) GetExternalWsUrl() string {
	return g.ExternalWsUrl
}

// GetContainerName returns the name of the container associated with the Besu instance.
// This function is useful for identifying and managing the container in a Docker environment.
func (g *Besu) GetContainerName() string {
	return g.ContainerName
}

// GetContainer returns a pointer to the container associated with the Besu instance.
// This function is useful for accessing the container's properties and methods in other operations.
func (g *Besu) GetContainer() *tc.Container {
	return &g.Container
}

// GetEthereumVersion returns the current Ethereum version of the Besu instance.
// This information is crucial for determining the appropriate container configuration and consensus mechanism.
func (g *Besu) GetEthereumVersion() config_types.EthereumVersion {
	return g.ethereumVersion
}

// WaitUntilChainIsReady blocks until the Ethereum chain is ready for operations.
// It is useful for ensuring that the execution client has fully synchronized with the network before proceeding with further actions.
func (g *Besu) WaitUntilChainIsReady(ctx context.Context, waitTime time.Duration) error {
	if g.GetEthereumVersion() == config_types.EthereumVersion_Eth1 {
		return nil
	}
	waitForFirstBlock := tcwait.NewLogStrategy("Imported #1").WithPollInterval(1 * time.Second).WithStartupTimeout(waitTime)
	return waitForFirstBlock.WaitUntilReady(ctx, *g.GetContainer())
}

// GetConsensusMechanism returns the consensus mechanism used by the Besu instance.
// It identifies whether the Ethereum version is Eth1 or not, returning either Proof of Authority (PoA)
// or Proof of Stake (PoS) accordingly. This is useful for understanding the network's validation method.
func (g *Besu) GetConsensusMechanism() ConsensusMechanism {
	if g.GetEthereumVersion() == config_types.EthereumVersion_Eth1 {
		return ConsensusMechanism_PoA
	}
	return ConsensusMechanism_PoS
}
