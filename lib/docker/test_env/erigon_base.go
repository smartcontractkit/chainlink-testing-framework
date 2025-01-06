package test_env

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/rs/zerolog"
	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/config"
	config_types "github.com/smartcontractkit/chainlink-testing-framework/lib/config/types"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/docker"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	docker_utils "github.com/smartcontractkit/chainlink-testing-framework/lib/utils/docker"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
)

type Erigon struct {
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
// This is useful for running tests that require a specific logging setup and context.
func (g *Erigon) WithTestInstance(t *testing.T) ExecutionClient {
	g.l = logging.GetTestLogger(t)
	g.t = t
	return g
}

// StartContainer initializes and starts an Erigon container for Ethereum execution.
// It configures network settings based on the Ethereum version and returns the
// blockchain network configuration along with any errors encountered during the process.
func (g *Erigon) StartContainer() (blockchain.EVMNetwork, error) {
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
		return blockchain.EVMNetwork{}, fmt.Errorf("cannot start erigon container: %w", err)
	}

	host, err := GetHost(testcontext.Get(g.t), ct)
	if err != nil {
		return blockchain.EVMNetwork{}, err
	}
	httpPort, err := ct.MappedPort(testcontext.Get(g.t), NatPort(DEFAULT_EVM_NODE_HTTP_PORT))
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
	g.ExternalWsUrl = FormatWsUrl(host, httpPort.Port())
	g.InternalWsUrl = FormatWsUrl(g.ContainerName, DEFAULT_EVM_NODE_HTTP_PORT)

	networkConfig := blockchain.SimulatedEVMNetwork
	if g.GetEthereumVersion() == config_types.EthereumVersion_Eth1 {
		networkConfig.Name = fmt.Sprintf("Private Eth-1-PoW [erigon %s]", g.ContainerVersion)
	} else {
		networkConfig.Name = fmt.Sprintf("Private Eth-2-PoS [erigon %s] + %s", g.ContainerVersion, g.consensusLayer)
	}
	networkConfig.URLs = []string{g.ExternalWsUrl}
	networkConfig.HTTPURLs = []string{g.ExternalHttpUrl}
	networkConfig.SimulationType = "Erigon"

	g.l.Info().Str("containerName", g.ContainerName).
		Msg("Started Erigon container")

	return networkConfig, nil
}

// GetInternalExecutionURL returns the internal execution URL for the Erigon client.
// It is used to retrieve the execution layer's endpoint, essential for connecting to the Ethereum network.
// If the Ethereum version is Eth1, it panics as Eth1 nodes do not support execution URLs.
func (g *Erigon) GetInternalExecutionURL() string {
	if g.GetEthereumVersion() == config_types.EthereumVersion_Eth1 {
		panic("eth1 node doesn't have an execution URL")
	}
	return g.InternalExecutionURL
}

// GetExternalExecutionURL returns the external execution URL for the Erigon instance.
// It panics if the Ethereum version is Eth1, as Eth1 nodes do not support execution URLs.
func (g *Erigon) GetExternalExecutionURL() string {
	if g.GetEthereumVersion() == config_types.EthereumVersion_Eth1 {
		panic("eth1 node doesn't have an execution URL")
	}
	return g.ExternalExecutionURL
}

// GetInternalHttpUrl returns the internal HTTP URL of the Erigon client.
// This URL is used to connect to the Erigon execution layer for internal communications.
func (g *Erigon) GetInternalHttpUrl() string {
	return g.InternalHttpUrl
}

// GetInternalWsUrl returns the internal WebSocket URL for the Erigon client.
// This URL is used to establish a WebSocket connection for real-time communication with the Erigon node.
func (g *Erigon) GetInternalWsUrl() string {
	return g.InternalWsUrl
}

// GetExternalHttpUrl returns the external HTTP URL for the Erigon client.
// This URL is used to interact with the Erigon execution layer over HTTP.
func (g *Erigon) GetExternalHttpUrl() string {
	return g.ExternalHttpUrl
}

// GetExternalWsUrl returns the external WebSocket URL for the Erigon client.
// This URL is essential for connecting to the Erigon node for real-time data and event subscriptions.
func (g *Erigon) GetExternalWsUrl() string {
	return g.ExternalWsUrl
}

// GetContainerName returns the name of the container associated with the Erigon instance.
// This function is useful for identifying and managing the container in a Docker environment.
func (g *Erigon) GetContainerName() string {
	return g.ContainerName
}

// GetContainer returns a pointer to the Container associated with the Erigon instance.
// This function is useful for accessing the container's properties and methods in a structured manner.
func (g *Erigon) GetContainer() *tc.Container {
	return &g.Container
}

// GetEthereumVersion returns the current Ethereum version of the Erigon instance.
// This information is essential for determining the appropriate execution URLs and consensus mechanisms.
func (g *Erigon) GetEthereumVersion() config_types.EthereumVersion {
	return g.ethereumVersion
}

// WaitUntilChainIsReady blocks until the Ethereum chain is ready for use, waiting for the first block to be built.
// It returns an error if the chain does not become ready within the specified wait time.
func (g *Erigon) WaitUntilChainIsReady(ctx context.Context, waitTime time.Duration) error {
	if g.GetEthereumVersion() == config_types.EthereumVersion_Eth1 {
		return nil
	}
	waitForFirstBlock := tcwait.NewLogStrategy("Built block").WithPollInterval(1 * time.Second).WithStartupTimeout(waitTime)
	return waitForFirstBlock.WaitUntilReady(ctx, *g.GetContainer())
}

// GetConsensusMechanism returns the consensus mechanism used by the Erigon instance.
// It identifies whether the Ethereum version is Eth1 (Proof of Work) or a later version (Proof of Stake).
func (g *Erigon) GetConsensusMechanism() ConsensusMechanism {
	if g.GetEthereumVersion() == config_types.EthereumVersion_Eth1 {
		return ConsensusMechanism_PoW
	}
	return ConsensusMechanism_PoS
}

func (g *Erigon) getExtraExecutionFlags() (string, error) {
	version, err := docker_utils.GetSemverFromImage(g.GetImageWithVersion())
	if err != nil {
		return "", err
	}

	extraExecutionFlags := ""

	// Erigon v2.47.0 and above have a new flag for disabling tx fee cap
	txFeeCapConstraint, err := semver.NewConstraint(">= 2.47.0")
	if err != nil {
		return "", err
	}

	if txFeeCapConstraint.Check(version) {
		extraExecutionFlags = " --rpc.txfeecap=0"
	}

	// Erigon v2.54.0 and above have a new flag for allowing unprotected txs
	allowUnprotectedTxsConstraint, err := semver.NewConstraint(">= 2.54.0")
	if err != nil {
		return "", err
	}

	if allowUnprotectedTxsConstraint.Check(version) {
		extraExecutionFlags += " --rpc.allow-unprotected-txs"
	}

	// Erigon v2.42.0 and above have a new flag for setting the db size limit
	dbSizeLimitConstraint, err := semver.NewConstraint(">= 2.42.0")
	if err != nil {
		return "", err
	}

	if dbSizeLimitConstraint.Check(version) {
		extraExecutionFlags += " --db.size.limit=8GB"
	}

	return extraExecutionFlags, nil
}
