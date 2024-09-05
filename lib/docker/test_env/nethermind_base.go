package test_env

import (
	"context"
	"fmt"
	"time"

	config_types "github.com/smartcontractkit/chainlink-testing-framework/lib/config/types"

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
	posContainerSettings
}

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

func (g *Nethermind) GetInternalExecutionURL() string {
	if g.GetEthereumVersion() == config_types.EthereumVersion_Eth1 {
		panic("eth1 node doesn't have an execution URL")
	}
	return g.InternalExecutionURL
}

func (g *Nethermind) GetExternalExecutionURL() string {
	if g.GetEthereumVersion() == config_types.EthereumVersion_Eth1 {
		panic("eth1 node doesn't have an execution URL")
	}
	return g.ExternalExecutionURL
}

func (g *Nethermind) GetInternalHttpUrl() string {
	return g.InternalHttpUrl
}

func (g *Nethermind) GetInternalWsUrl() string {
	return g.InternalWsUrl
}

func (g *Nethermind) GetExternalHttpUrl() string {
	return g.ExternalHttpUrl
}

func (g *Nethermind) GetExternalWsUrl() string {
	return g.ExternalWsUrl
}

func (g *Nethermind) GetContainerName() string {
	return g.ContainerName
}

func (g *Nethermind) GetContainer() *tc.Container {
	return &g.Container
}

func (g *Nethermind) GetEthereumVersion() config_types.EthereumVersion {
	return g.ethereumVersion
}

func (g *Nethermind) GetEnvComponent() *EnvComponent {
	return &g.EnvComponent
}

func (g *Nethermind) WaitUntilChainIsReady(ctx context.Context, waitTime time.Duration) error {
	if g.GetEthereumVersion() == config_types.EthereumVersion_Eth1 {
		return nil
	}
	waitForFirstBlock := tcwait.NewLogStrategy("Improved post-merge block").WithPollInterval(1 * time.Second).WithStartupTimeout(waitTime)
	return waitForFirstBlock.WaitUntilReady(ctx, *g.GetContainer())
}

func (g *Nethermind) GethConsensusMechanism() ConsensusMechanism {
	if g.GetEthereumVersion() == config_types.EthereumVersion_Eth1 {
		return ConsensusMechanism_PoA
	}
	return ConsensusMechanism_PoS
}
