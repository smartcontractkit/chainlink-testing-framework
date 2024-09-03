package test_env

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/semver/v3"
	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/config"
	config_types "github.com/smartcontractkit/chainlink-testing-framework/config/types"
	"github.com/smartcontractkit/chainlink-testing-framework/docker"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	docker_utils "github.com/smartcontractkit/chainlink-testing-framework/utils/docker"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/testcontext"
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
	//l                    zerolog.Logger
	//t                    *testing.T
	posContainerSettings
}

//func (g *Erigon) WithTestInstance(t *testing.T) ExecutionClient {
//	g.l = logging.GetTestLogger(t)
//	g.t = t
//	return g
//}

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

func (g *Erigon) GetInternalExecutionURL() string {
	if g.GetEthereumVersion() == config_types.EthereumVersion_Eth1 {
		panic("eth1 node doesn't have an execution URL")
	}
	return g.InternalExecutionURL
}

func (g *Erigon) GetExternalExecutionURL() string {
	if g.GetEthereumVersion() == config_types.EthereumVersion_Eth1 {
		panic("eth1 node doesn't have an execution URL")
	}
	return g.ExternalExecutionURL
}

func (g *Erigon) GetInternalHttpUrl() string {
	return g.InternalHttpUrl
}

func (g *Erigon) GetInternalWsUrl() string {
	return g.InternalWsUrl
}

func (g *Erigon) GetExternalHttpUrl() string {
	return g.ExternalHttpUrl
}

func (g *Erigon) GetExternalWsUrl() string {
	return g.ExternalWsUrl
}

func (g *Erigon) GetContainerName() string {
	return g.ContainerName
}

func (g *Erigon) GetContainer() *tc.Container {
	return &g.Container
}

func (g *Erigon) GetEthereumVersion() config_types.EthereumVersion {
	return g.ethereumVersion
}

func (g *Erigon) GetEnvComponent() *EnvComponent {
	return &g.EnvComponent
}

func (g *Erigon) WaitUntilChainIsReady(ctx context.Context, waitTime time.Duration) error {
	if g.GetEthereumVersion() == config_types.EthereumVersion_Eth1 {
		return nil
	}
	waitForFirstBlock := tcwait.NewLogStrategy("Built block").WithPollInterval(1 * time.Second).WithStartupTimeout(waitTime)
	return waitForFirstBlock.WaitUntilReady(ctx, *g.GetContainer())
}

func (g *Erigon) GethConsensusMechanism() ConsensusMechanism {
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
