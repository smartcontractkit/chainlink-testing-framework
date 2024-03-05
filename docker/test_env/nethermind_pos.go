package test_env

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/google/uuid"
	"github.com/rs/zerolog"

	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/docker"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/mirror"
)

const defaultNethermindPosImage = "nethermind/nethermind:1.25.1"

type NethermindPos struct {
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

func NewNethermindPos(networks []string, generatedDataHostDir string, consensusLayer ConsensusLayer, opts ...EnvComponentOption) (*NethermindPos, error) {
	parts := strings.Split(defaultNethermindPosImage, ":")
	g := &NethermindPos{
		EnvComponent: EnvComponent{
			ContainerName:    fmt.Sprintf("%s-%s", "nethermind-pos", uuid.NewString()[0:8]),
			Networks:         networks,
			ContainerImage:   parts[0],
			ContainerVersion: parts[1],
		},
		generatedDataHostDir: generatedDataHostDir,
		consensusLayer:       consensusLayer,
		l:                    logging.GetTestLogger(nil),
	}
	for _, opt := range opts {
		opt(&g.EnvComponent)
	}
	// if the internal docker repo is set then add it to the version
	g.EnvComponent.ContainerImage = mirror.AddMirrorToImageIfSet(g.EnvComponent.ContainerImage)
	return g, nil
}

func (g *NethermindPos) WithTestInstance(t *testing.T) ExecutionClient {
	g.l = logging.GetTestLogger(t)
	g.t = t
	return g
}

func (g *NethermindPos) StartContainer() (blockchain.EVMNetwork, error) {
	r, err := g.getContainerRequest()
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
		return blockchain.EVMNetwork{}, fmt.Errorf("cannot start nethermind container: %w", err)
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
	networkConfig.Name = fmt.Sprintf("Simulated Ethereum-PoS (nethermind + %s)", g.consensusLayer)
	networkConfig.URLs = []string{g.ExternalWsUrl}
	networkConfig.HTTPURLs = []string{g.ExternalHttpUrl}

	g.l.Info().Str("containerName", g.ContainerName).
		Msg("Started Nethermind container")

	return networkConfig, nil
}

func (g *NethermindPos) GetInternalExecutionURL() string {
	return g.InternalExecutionURL
}

func (g *NethermindPos) GetExternalExecutionURL() string {
	return g.ExternalExecutionURL
}

func (g *NethermindPos) GetInternalHttpUrl() string {
	return g.InternalHttpUrl
}

func (g *NethermindPos) GetInternalWsUrl() string {
	return g.InternalWsUrl
}

func (g *NethermindPos) GetExternalHttpUrl() string {
	return g.ExternalHttpUrl
}

func (g *NethermindPos) GetExternalWsUrl() string {
	return g.ExternalWsUrl
}

func (g *NethermindPos) GetContainerName() string {
	return g.ContainerName
}

func (g *NethermindPos) GetContainer() *tc.Container {
	return &g.Container
}

func (g *NethermindPos) getContainerRequest() (*tc.ContainerRequest, error) {
	// create empty cfg file since if we don't pass any
	// default mainnet.cfg will be used
	noneCfg, err := os.CreateTemp("", "none.cfg")
	if err != nil {
		return nil, err
	}

	_, err = noneCfg.WriteString("{}")
	if err != nil {
		return nil, err
	}

	return &tc.ContainerRequest{
		Name:            g.ContainerName,
		Image:           g.GetImageWithVersion(),
		Networks:        g.Networks,
		AlwaysPullImage: true,
		// ImagePlatform: "linux/x86_64",  //don't even try this on Apple Silicon, the node won't start due to .NET error
		ExposedPorts: []string{NatPortFormat(TX_GETH_HTTP_PORT), NatPortFormat(TX_GETH_WS_PORT), NatPortFormat(ETH2_EXECUTION_PORT)},
		WaitingFor: tcwait.ForAll(
			tcwait.ForLog("Nethermind initialization completed").
				WithStartupTimeout(120 * time.Second).
				WithPollInterval(1 * time.Second),
		),
		Cmd: []string{
			"--datadir=/nethermind",
			"--config=/none.cfg",
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
		Files: []tc.ContainerFile{
			{
				HostFilePath:      noneCfg.Name(),
				ContainerFilePath: "/none.cfg",
				FileMode:          0644,
			},
		},
		HostConfigModifier: func(hostConfig *container.HostConfig) {
			hostConfig.Mounts = append(hostConfig.Mounts, mount.Mount{
				Type:     mount.TypeBind,
				Source:   g.generatedDataHostDir,
				Target:   GENERATED_DATA_DIR_INSIDE_CONTAINER,
				ReadOnly: false,
			})
		},
		LifecycleHooks: []tc.ContainerLifecycleHooks{
			{
				PostStarts: g.PostStartsHooks,
				PostStops:  g.PostStopsHooks,
			},
		},
	}, nil
}

func (g *NethermindPos) WaitUntilChainIsReady(ctx context.Context, waitTime time.Duration) error {
	waitForFirstBlock := tcwait.NewLogStrategy("Improved post-merge block").WithPollInterval(1 * time.Second).WithStartupTimeout(waitTime)
	return waitForFirstBlock.WaitUntilReady(ctx, *g.GetContainer())
}

func (g *NethermindPos) GetContainerType() ContainerType {
	return ContainerType_Nethermind
}
