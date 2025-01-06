package test_env

import (
	"context"
	"fmt"
	"testing"
	"time"

	config_types "github.com/smartcontractkit/chainlink-testing-framework/lib/config/types"

	"github.com/Masterminds/semver/v3"
	"github.com/rs/zerolog"
	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/config"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/docker"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	docker_utils "github.com/smartcontractkit/chainlink-testing-framework/lib/utils/docker"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
)

type Geth struct {
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

// WithTestInstance sets up the execution client with a test logger and associates it with the provided testing context.
// This is useful for capturing logs during testing and ensuring that the client operates in a test environment.
func (g *Geth) WithTestInstance(t *testing.T) ExecutionClient {
	g.l = logging.GetTestLogger(t)
	g.t = t
	return g
}

// StartContainer initializes and starts a Geth container for Ethereum, returning the network configuration.
// It supports both Eth1 and Eth2 versions, setting up necessary URLs for communication with the blockchain.
func (g *Geth) StartContainer() (blockchain.EVMNetwork, error) {
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
		return blockchain.EVMNetwork{}, fmt.Errorf("cannot start geth container: %w", err)
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
		networkConfig.Name = fmt.Sprintf("Private Eth-1-PoA [geth %s]", g.ContainerVersion)
	} else {
		networkConfig.Name = fmt.Sprintf("Private Eth-2-PoS [geth %s] + %s", g.ContainerVersion, g.consensusLayer)
	}
	networkConfig.URLs = []string{g.ExternalWsUrl}
	networkConfig.HTTPURLs = []string{g.ExternalHttpUrl}
	networkConfig.SimulationType = "Geth"

	version, err := docker_utils.GetSemverFromImage(g.GetImageWithVersion())
	if err != nil {
		return blockchain.EVMNetwork{}, err
	}

	constraint, err := semver.NewConstraint(">=1.10 <1.11")
	if err != nil {
		return blockchain.EVMNetwork{}, fmt.Errorf("failed to parse constraint: %s", ">=1.10 <1.11")
	}

	if constraint.Check(version) {
		// Geth v1.10.x will not set it itself if it's set 0, like later versions do
		networkConfig.DefaultGasLimit = 9_000_000
	}

	g.l.Info().Str("containerName", g.ContainerName).
		Msg("Started Geth PoS container")

	return networkConfig, nil
}

// GetInternalExecutionURL returns the internal execution URL for the Ethereum client.
// It is used to retrieve the endpoint for executing transactions in Ethereum 2.0 environments.
// An error is raised if called on an Ethereum 1.0 node.
func (g *Geth) GetInternalExecutionURL() string {
	if g.GetEthereumVersion() == config_types.EthereumVersion_Eth1 {
		panic("eth1 node doesn't have an execution URL")
	}
	return g.InternalExecutionURL
}

// GetExternalExecutionURL returns the external execution URL for the Geth instance.
// It panics if the Ethereum version is Eth1, as Eth1 nodes do not support execution URLs.
func (g *Geth) GetExternalExecutionURL() string {
	if g.GetEthereumVersion() == config_types.EthereumVersion_Eth1 {
		panic("eth1 node doesn't have an execution URL")
	}
	return g.ExternalExecutionURL
}

// GetInternalHttpUrl returns the internal HTTP URL of the Geth client.
// This URL is essential for connecting to the execution layer in a private network setup.
func (g *Geth) GetInternalHttpUrl() string {
	return g.InternalHttpUrl
}

// GetInternalWsUrl returns the internal WebSocket URL for the Geth client.
// This URL is essential for establishing WebSocket connections to the Geth node, enabling real-time data streaming and event handling.
func (g *Geth) GetInternalWsUrl() string {
	return g.InternalWsUrl
}

// GetExternalHttpUrl returns the external HTTP URL for the Geth client.
// This URL is used to interact with the Ethereum network externally, enabling communication with other services or clients.
func (g *Geth) GetExternalHttpUrl() string {
	return g.ExternalHttpUrl
}

// GetExternalWsUrl returns the external WebSocket URL for the Geth client.
// This URL is essential for connecting to the Ethereum network via WebSocket.
func (g *Geth) GetExternalWsUrl() string {
	return g.ExternalWsUrl
}

// GetContainerName returns the name of the container associated with the Geth instance.
// This function is useful for identifying and managing the container in a Docker environment.
func (g *Geth) GetContainerName() string {
	return g.ContainerName
}

// GetContainer returns a pointer to the container associated with the Geth instance.
// This function is useful for accessing the container's properties and methods in a seamless manner.
func (g *Geth) GetContainer() *tc.Container {
	return &g.Container
}

// GetEthereumVersion returns the current Ethereum version of the Geth instance.
// This information is essential for determining the appropriate execution URL and consensus mechanism.
func (g *Geth) GetEthereumVersion() config_types.EthereumVersion {
	return g.ethereumVersion
}

// GetConsensusMechanism returns the consensus mechanism used by the Geth instance.
// It determines whether the Ethereum version is Eth1 or not, returning either Proof of Authority (PoA)
// or Proof of Stake (PoS) accordingly. This is useful for understanding the underlying protocol.
func (g *Geth) GetConsensusMechanism() ConsensusMechanism {
	if g.GetEthereumVersion() == config_types.EthereumVersion_Eth1 {
		return ConsensusMechanism_PoA
	}
	return ConsensusMechanism_PoS
}

// WaitUntilChainIsReady blocks until the Ethereum chain is ready for use, or returns an error if it times out.
// This function is essential for ensuring that the blockchain environment is fully operational before proceeding with further operations.
func (g *Geth) WaitUntilChainIsReady(ctx context.Context, waitTime time.Duration) error {
	if g.GetEthereumVersion() == config_types.EthereumVersion_Eth1 {
		return nil
	}
	waitForFirstBlock := tcwait.NewLogStrategy("Chain head was updated").WithPollInterval(1 * time.Second).WithStartupTimeout(waitTime)
	return waitForFirstBlock.WaitUntilReady(ctx, *g.GetContainer())
}

func (g *Geth) getEntryPointAndKeystoreLocation(minerAddress string) ([]string, error) {
	version, err := docker_utils.GetSemverFromImage(g.GetImageWithVersion())
	if err != nil {
		return nil, err
	}

	var enabledApis = "admin,debug,web3,eth,txpool,personal,clique,miner,net"
	var localhost = "0.0.0.0"

	entrypoint := []string{
		"sh", "./root/init.sh",
		"--password", "/root/config/password.txt",
		"--ipcdisable",
		"--graphql",
		"-graphql.corsdomain", "*",
		"--allow-insecure-unlock",
		"--vmdebug",
		fmt.Sprintf("--networkid=%d", g.chainConfig.ChainID),
		"--datadir", "/root/.ethereum/devchain",
		"--mine",
		"--miner.etherbase", minerAddress,
		"--unlock", minerAddress,
	}

	constraint, err := semver.NewConstraint("<1.10")
	if err != nil {
		return nil, fmt.Errorf("failed to parse constraint: %s", "<1.10")
	}

	if constraint.Check(version) {
		entrypoint = append(entrypoint,
			"--rpc",
			"--rpcapi", enabledApis,
			"--rpccorsdomain", "*",
			"--rpcvhosts", "*",
			"--rpcaddr", localhost,
			fmt.Sprintf("--rpcport=%s", DEFAULT_EVM_NODE_HTTP_PORT),
			"--ws",
			"--wsorigins", "*",
			"--wsaddr", localhost,
			"--wsapi", enabledApis,
			fmt.Sprintf("--wsport=%s", DEFAULT_EVM_NODE_WS_PORT),
		)
	} else {
		entrypoint = append(entrypoint,
			"--http",
			"--http.vhosts", "*",
			"--http.api", enabledApis,
			"--http.corsdomain", "*",
			"--http.addr", localhost,
			fmt.Sprintf("--http.port=%s", DEFAULT_EVM_NODE_HTTP_PORT),
			"--ws",
			"--ws.origins", "*",
			"--ws.addr", localhost,
			"--ws.api", enabledApis,
			fmt.Sprintf("--ws.port=%s", DEFAULT_EVM_NODE_WS_PORT),
			"--rpc.allow-unprotected-txs",
			"--rpc.txfeecap", "0",
		)
	}

	return entrypoint, nil
}

func (g *Geth) getWebsocketEnabledMessage() (string, error) {
	version, err := docker_utils.GetSemverFromImage(g.GetImageWithVersion())
	if err != nil {
		return "", err
	}

	constraint, err := semver.NewConstraint("<1.10")
	if err != nil {
		return "", fmt.Errorf("failed to parse constraint: %s", "<1.10")
	}

	if constraint.Check(version) {
		return "WebSocket endpoint opened", nil
	}

	return "WebSocket enabled", nil
}

func (g *Geth) logLevelToVerbosity() (int, error) {
	switch g.LogLevel {
	case "trace":
		return 5, nil
	case "debug":
		return 4, nil
	case "info":
		return 3, nil
	case "warn":
		return 2, nil
	case "error":
		return 1, nil
	case "silent":
		return 0, nil
	default:
		return -1, fmt.Errorf("unknown log level: %s", g.LogLevel)
	}
}
