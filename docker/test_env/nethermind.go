package test_env

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/docker"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
)

const (
	NETHERMIND_IMAGE_TAG = "1.22.0"
)

type Nethermind struct {
	EnvComponent
	ExternalHttpUrl      string
	InternalHttpUrl      string
	ExternalWsUrl        string
	InternalWsUrl        string
	InternalExecutionURL string
	ExternalExecutionURL string
	generatedDataHostDir string
	consensusLayer       ConsensusLayer
	l                    zerolog.Logger
	t                    *testing.T
}

func NewNethermind(networks []string, generatedDataHostDir string, consensusLayer ConsensusLayer, opts ...EnvComponentOption) *Nethermind {
	g := &Nethermind{
		EnvComponent: EnvComponent{
			ContainerName: fmt.Sprintf("%s-%s", "nethermind", uuid.NewString()[0:8]),
			Networks:      networks,
		},
		generatedDataHostDir: generatedDataHostDir,
		consensusLayer:       consensusLayer,
		l:                    log.Logger,
	}
	for _, opt := range opts {
		opt(&g.EnvComponent)
	}
	return g
}

func (g *Nethermind) WithTestInstance(t *testing.T) *Nethermind {
	g.l = logging.GetTestLogger(t)
	g.t = t
	return g
}

func (g *Nethermind) StartContainer() (blockchain.EVMNetwork, error) {
	r, err := g.getContainerRequest(g.Networks)
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
		return blockchain.EVMNetwork{}, errors.Wrapf(err, "cannot start nethermind container")
	}

	host, err := GetHost(context.Background(), ct)
	if err != nil {
		return blockchain.EVMNetwork{}, err
	}
	if err != nil {
		return blockchain.EVMNetwork{}, err
	}
	httpPort, err := ct.MappedPort(context.Background(), NatPort(TX_GETH_HTTP_PORT))
	if err != nil {
		return blockchain.EVMNetwork{}, err
	}
	wsPort, err := ct.MappedPort(context.Background(), NatPort(TX_GETH_WS_PORT))
	if err != nil {
		return blockchain.EVMNetwork{}, err
	}
	executionPort, err := ct.MappedPort(context.Background(), NatPort(ETH2_EXECUTION_PORT))
	if err != nil {
		return blockchain.EVMNetwork{}, err
	}

	g.Container = ct
	g.ExternalHttpUrl = FormatHttpUrl(host, httpPort.Port())
	g.InternalHttpUrl = FormatHttpUrl(g.ContainerName, TX_GETH_HTTP_PORT)
	g.ExternalWsUrl = FormatWsUrl(host, wsPort.Port())
	g.InternalWsUrl = FormatWsUrl(g.ContainerName, TX_GETH_WS_PORT)
	g.InternalExecutionURL = FormatHttpUrl(g.ContainerName, ETH2_EXECUTION_PORT)
	g.ExternalExecutionURL = FormatHttpUrl(host, executionPort.Port())

	networkConfig := blockchain.SimulatedEVMNetwork
	networkConfig.Name = fmt.Sprintf("Simulated Eth2 (nethermind + %s)", g.consensusLayer)
	networkConfig.URLs = []string{g.ExternalWsUrl}
	networkConfig.HTTPURLs = []string{g.ExternalHttpUrl}

	g.l.Info().Str("containerName", g.ContainerName).
		Msg("Started Nethermind container")

	return networkConfig, nil
}

func (g *Nethermind) GetInternalExecutionURL() string {
	return g.InternalExecutionURL
}

func (g *Nethermind) GetExternalExecutionURL() string {
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

func (g *Nethermind) getContainerRequest(networks []string) (*tc.ContainerRequest, error) {
	return &tc.ContainerRequest{
		Name:            g.ContainerName,
		Image:           fmt.Sprintf("nethermind/nethermind:%s", NETHERMIND_IMAGE_TAG),
		Networks:        networks,
		AlwaysPullImage: true,
		// ImagePlatform: "linux/x86_64", // this breaks everything, don't try it
		ExposedPorts: []string{NatPortFormat(TX_GETH_HTTP_PORT), NatPortFormat(TX_GETH_WS_PORT), NatPortFormat(ETH2_EXECUTION_PORT)},
		WaitingFor: tcwait.ForAll(
			tcwait.ForLog("Nethermind initialization completed").
				WithStartupTimeout(120 * time.Second).
				WithPollInterval(1 * time.Second),
		),
		Cmd: []string{
			"--datadir=/nethermind",
			"--config=none.cfg",
			fmt.Sprintf("--Init.ChainSpecPath=%s/chainspec.json", GENERATED_DATA_DIR_INSIDE_CONTAINER),
			"--Init.DiscoveryEnabled=false",
			"--Init.WebSocketsEnabled=true",
			fmt.Sprintf("--JsonRpc.WebSocketsPort=%s", TX_GETH_WS_PORT),
			"--JsonRpc.Enabled=true",
			"--JsonRpc.EnabledModules=net,eth,consensus,subscribe,web3,admin",
			"--JsonRpc.Host=0.0.0.0",
			fmt.Sprintf("--JsonRpc.Port=%s", TX_GETH_HTTP_PORT),
			"--JsonRpc.EngineHost=0.0.0.0",
			"--JsonRpc.EnginePort=" + ETH2_EXECUTION_PORT,
			fmt.Sprintf("--JsonRpc.JwtSecretFile=%s", JWT_SECRET_FILE_LOCATION_INSIDE_CONTAINER),
			fmt.Sprintf("--KeyStore.KeyStoreDirectory=%s", KEYSTORE_DIR_LOCATION_INSIDE_CONTAINER),
			"--KeyStore.BlockAuthorAccount=0x123463a4b065722e99115d6c222f267d9cabb524",
			"--KeyStore.UnlockAccounts=0x123463a4b065722e99115d6c222f267d9cabb524",
			fmt.Sprintf("--KeyStore.PasswordFiles=%s", ACCOUNT_PASSWORD_FILE_INSIDE_CONTAINER),
			"--Network.MaxActivePeers=0",
			"--Network.OnlyStaticPeers=true",
			"--HealthChecks.Enabled=true", // default slug /health
		},
		Mounts: tc.ContainerMounts{
			tc.ContainerMount{
				Source: tc.GenericBindMountSource{
					HostPath: g.generatedDataHostDir,
				},
				Target: tc.ContainerMountTarget(GENERATED_DATA_DIR_INSIDE_CONTAINER),
			},
		},
	}, nil
}

func (g Nethermind) WaitUntilChainIsReady(waitTime time.Duration) error {
	waitForFirstBlock := tcwait.NewLogStrategy("Improved post-merge block").WithPollInterval(1 * time.Second).WithStartupTimeout(waitTime)
	return waitForFirstBlock.WaitUntilReady(context.Background(), *g.GetContainer())
}

func (g *Nethermind) GetContainerType() ContainerType {
	return ContainerType_Nethermind
}
