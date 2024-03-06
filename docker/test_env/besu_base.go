package test_env

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/rs/zerolog"

	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/docker"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/testcontext"
)

const (
	defaultBesuEth1Image = "hyperledger/besu:22.1.0"
	defaultBesuEth2Image = "hyperledger/besu:24.1.0"
)

type posSettings struct {
	generatedDataHostDir string
}

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
	chainConfg           *EthereumChainConfig
	consensusLayer       ConsensusLayer
	l                    zerolog.Logger
	t                    *testing.T
	posSettings
	powSettings
}

func (g *Besu) WithTestInstance(t *testing.T) ExecutionClient {
	g.l = logging.GetTestLogger(t)
	g.t = t
	return g
}

func (g *Besu) StartContainer() (blockchain.EVMNetwork, error) {
	var r *tc.ContainerRequest
	var err error

	if g.GetEthereumVersion() == EthereumVersion_Eth1 {
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
		Reuse:            true,
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

	if g.GetEthereumVersion() == EthereumVersion_Eth2 {
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

	if g.GetEthereumVersion() == EthereumVersion_Eth1 {
		networkConfig.Name = "Simulated Eth-1-PoA (besu)"
	} else {
		networkConfig.Name = fmt.Sprintf("Simulated Eth-2-PoS (besu + %s)", g.consensusLayer)
	}

	g.l.Info().Str("containerName", g.ContainerName).
		Msg("Started Besu container")

	return networkConfig, nil
}

func (g *Besu) GetInternalExecutionURL() string {
	if g.GetEthereumVersion() == EthereumVersion_Eth1 {
		panic("eth1 node doesn't have an execution URL")
	}
	return g.InternalExecutionURL
}

func (g *Besu) GetExternalExecutionURL() string {
	if g.GetEthereumVersion() == EthereumVersion_Eth1 {
		panic("eth1 node doesn't have an execution URL")
	}
	return g.ExternalExecutionURL
}

func (g *Besu) GetInternalHttpUrl() string {
	return g.InternalHttpUrl
}

func (g *Besu) GetInternalWsUrl() string {
	return g.InternalWsUrl
}

func (g *Besu) GetExternalHttpUrl() string {
	return g.ExternalHttpUrl
}

func (g *Besu) GetExternalWsUrl() string {
	return g.ExternalWsUrl
}

func (g *Besu) GetContainerName() string {
	return g.ContainerName
}

func (g *Besu) GetContainer() *tc.Container {
	return &g.Container
}

func (g *Besu) GetEthereumVersion() EthereumVersion {
	if g.consensusLayer == "" {
		return EthereumVersion_Eth1
	}

	return EthereumVersion_Eth2
}

func (g *Besu) WaitUntilChainIsReady(ctx context.Context, waitTime time.Duration) error {
	if g.GetEthereumVersion() == EthereumVersion_Eth1 {
		return nil
	}
	waitForFirstBlock := tcwait.NewLogStrategy("Imported #1").WithPollInterval(1 * time.Second).WithStartupTimeout(waitTime)
	return waitForFirstBlock.WaitUntilReady(ctx, *g.GetContainer())
}

func (g *Besu) GetContainerType() ContainerType {
	return ContainerType_Besu
}

func (g *Besu) GethConsensusMechanism() ConsensusMechanism {
	if g.GetEthereumVersion() == EthereumVersion_Eth1 {
		return ConsensusMechanism_PoA
	}
	return ConsensusMechanism_PoS
}
