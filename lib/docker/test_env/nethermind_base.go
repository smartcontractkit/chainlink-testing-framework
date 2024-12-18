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
)

type Nethermind struct {
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

// WithTestInstance sets up the execution client with a test logger and test context.
// This is useful for running tests with specific logging and context requirements.
func (g *Nethermind) WithTestInstance(t *testing.T) ExecutionClient {
	g.l = logging.GetTestLogger(t)
	g.t = t
	return g
}

// StartContainer initializes and starts a Nethermind container for Ethereum execution.
// It configures network settings based on the Ethereum version and returns the network configuration
// along with any errors encountered during the process.
func (g *Nethermind) StartContainer() (blockchain.EVMNetwork, error) {
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
		return blockchain.EVMNetwork{}, fmt.Errorf("cannot start nethermind container: %w", err)
	}

	host, err := GetHost(context.Background(), ct)
	if err != nil {
		return blockchain.EVMNetwork{}, err
	}
	httpPort, err := ct.MappedPort(context.Background(), NatPort(DEFAULT_EVM_NODE_HTTP_PORT))
	if err != nil {
		return blockchain.EVMNetwork{}, err
	}
	wsPort, err := ct.MappedPort(context.Background(), NatPort(DEFAULT_EVM_NODE_WS_PORT))
	if err != nil {
		return blockchain.EVMNetwork{}, err
	}

	if g.GetEthereumVersion() == config_types.EthereumVersion_Eth2 {
		executionPort, err := ct.MappedPort(context.Background(), NatPort(ETH2_EXECUTION_PORT))
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
		networkConfig.Name = fmt.Sprintf("Private Eth-1-PoA [nethermind %s", g.ContainerVersion)
		networkConfig.GasEstimationBuffer = 100_000_000_000
	} else {
		networkConfig.Name = fmt.Sprintf("Private Eth-2-PoS [nethermind %s] + %s", g.ContainerVersion, g.consensusLayer)
	}
	networkConfig.URLs = []string{g.ExternalWsUrl}
	networkConfig.HTTPURLs = []string{g.ExternalHttpUrl}
	networkConfig.SimulationType = "Nethermind"

	g.l.Info().Str("containerName", g.ContainerName).
		Msg("Started Nethermind container")

	return networkConfig, nil
}

// GetInternalExecutionURL returns the internal execution URL for the Ethereum client.
// It is used to retrieve the endpoint for executing transactions in a network that supports it.
// If the client is an Eth1 node, it will panic as Eth1 does not have an execution URL.
func (g *Nethermind) GetInternalExecutionURL() string {
	if g.GetEthereumVersion() == config_types.EthereumVersion_Eth1 {
		panic("eth1 node doesn't have an execution URL")
	}
	return g.InternalExecutionURL
}

// GetExternalExecutionURL returns the external execution URL for the Nethermind node.
// It panics if the node is running on Ethereum version 1, as it does not support execution URLs.
func (g *Nethermind) GetExternalExecutionURL() string {
	if g.GetEthereumVersion() == config_types.EthereumVersion_Eth1 {
		panic("eth1 node doesn't have an execution URL")
	}
	return g.ExternalExecutionURL
}

// GetInternalHttpUrl returns the internal HTTP URL of the Nethermind client.
// This URL is used to communicate with the execution layer in a secure manner.
func (g *Nethermind) GetInternalHttpUrl() string {
	return g.InternalHttpUrl
}

// GetInternalWsUrl returns the internal WebSocket URL for the Nethermind client.
// This URL is essential for establishing WebSocket connections to the execution layer.
func (g *Nethermind) GetInternalWsUrl() string {
	return g.InternalWsUrl
}

// GetExternalHttpUrl returns the external HTTP URL for the Nethermind client.
// This URL is used to interact with the Nethermind execution layer over HTTP.
func (g *Nethermind) GetExternalHttpUrl() string {
	return g.ExternalHttpUrl
}

// GetExternalWsUrl returns the external WebSocket URL for the Nethermind client.
// This URL is essential for connecting to the Nethermind execution layer from external services.
func (g *Nethermind) GetExternalWsUrl() string {
	return g.ExternalWsUrl
}

// GetContainerName returns the name of the container associated with the Nethermind instance.
// This function is useful for identifying and managing the container in a Docker environment.
func (g *Nethermind) GetContainerName() string {
	return g.ContainerName
}

// GetContainer returns a pointer to the Container associated with the Nethermind instance.
// This function is useful for accessing the container's properties and methods in a Dockerized environment.
func (g *Nethermind) GetContainer() *tc.Container {
	return &g.Container
}

// GetEthereumVersion returns the current Ethereum version of the Nethermind instance.
// This information is crucial for determining the appropriate execution URLs and consensus mechanisms.
func (g *Nethermind) GetEthereumVersion() config_types.EthereumVersion {
	return g.ethereumVersion
}

// WaitUntilChainIsReady blocks until the Ethereum chain is fully operational or the specified wait time elapses.
// It is useful for ensuring that the blockchain client is ready for transactions or queries before proceeding.
func (g *Nethermind) WaitUntilChainIsReady(ctx context.Context, waitTime time.Duration) error {
	if g.GetEthereumVersion() == config_types.EthereumVersion_Eth1 {
		return nil
	}
	waitForFirstBlock := tcwait.NewLogStrategy("Improved post-merge block").WithPollInterval(1 * time.Second).WithStartupTimeout(waitTime)
	return waitForFirstBlock.WaitUntilReady(ctx, *g.GetContainer())
}

// GetConsensusMechanism returns the consensus mechanism used by the Nethermind instance.
// It determines whether the Ethereum version is Eth1 or not, returning either Proof of Authority (PoA)
// or Proof of Stake (PoS) accordingly. This is useful for understanding the network's validation method.
func (g *Nethermind) GetConsensusMechanism() ConsensusMechanism {
	if g.GetEthereumVersion() == config_types.EthereumVersion_Eth1 {
		return ConsensusMechanism_PoA
	}
	return ConsensusMechanism_PoS
}
