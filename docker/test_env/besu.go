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
	"github.com/smartcontractkit/chainlink-testing-framework/utils/testcontext"
)

const (
	//TODO use Tate's mirror?
	BESU_IMAGE_TAG = "23.10"
)

type Besu struct {
	EnvComponent
	ExternalHttpUrl      string
	InternalHttpUrl      string
	ExternalWsUrl        string
	InternalWsUrl        string
	InternalExecutionURL string
	ExternalExecutionURL string
	generatedDataHostDir string
	chainConfg           *EthereumChainConfig
	consensusLayer       ConsensusLayer
	l                    zerolog.Logger
	t                    *testing.T
}

func NewBesu(networks []string, chainConfg *EthereumChainConfig, generatedDataHostDir string, consensusLayer ConsensusLayer, opts ...EnvComponentOption) *Besu {
	g := &Besu{
		EnvComponent: EnvComponent{
			ContainerName: fmt.Sprintf("%s-%s", "besu", uuid.NewString()[0:8]),
			Networks:      networks,
		},
		chainConfg:           chainConfg,
		generatedDataHostDir: generatedDataHostDir,
		consensusLayer:       consensusLayer,
		l:                    log.Logger,
	}
	for _, opt := range opts {
		opt(&g.EnvComponent)
	}
	return g
}

func (g *Besu) WithTestInstance(t *testing.T) *Besu {
	g.l = logging.GetTestLogger(t)
	g.t = t
	return g
}

func (g *Besu) StartContainer() (blockchain.EVMNetwork, error) {
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
		return blockchain.EVMNetwork{}, errors.Wrapf(err, "cannot start Besu container")
	}

	host, err := GetHost(testcontext.Get(g.t), ct)
	if err != nil {
		return blockchain.EVMNetwork{}, err
	}
	if err != nil {
		return blockchain.EVMNetwork{}, err
	}
	httpPort, err := ct.MappedPort(testcontext.Get(g.t), NatPort(TX_GETH_HTTP_PORT))
	if err != nil {
		return blockchain.EVMNetwork{}, err
	}
	wsPort, err := ct.MappedPort(testcontext.Get(g.t), NatPort(TX_GETH_WS_PORT))
	if err != nil {
		return blockchain.EVMNetwork{}, err
	}
	executionPort, err := ct.MappedPort(testcontext.Get(g.t), NatPort(ETH2_EXECUTION_PORT))
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
	networkConfig.Name = fmt.Sprintf("Simulated Eth2 (Besu %s)", g.consensusLayer)
	networkConfig.URLs = []string{g.ExternalWsUrl}
	networkConfig.HTTPURLs = []string{g.ExternalHttpUrl}

	g.l.Info().Str("containerName", g.ContainerName).
		Msg("Started Besu container")

	return networkConfig, nil
}

func (g *Besu) GetInternalExecutionURL() string {
	return g.InternalExecutionURL
}

func (g *Besu) GetExternalExecutionURL() string {
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

func (g *Besu) getContainerRequest(networks []string) (*tc.ContainerRequest, error) {
	return &tc.ContainerRequest{
		Name:     g.ContainerName,
		Image:    fmt.Sprintf("hyperledger/besu:%s", BESU_IMAGE_TAG),
		Networks: networks,
		// ImagePlatform: "linux/x86_64", //don't even try this on Apple Silicon, the node won't start due to JVM error
		ExposedPorts: []string{NatPortFormat(TX_GETH_HTTP_PORT), NatPortFormat(TX_GETH_WS_PORT), NatPortFormat(ETH2_EXECUTION_PORT)},
		WaitingFor: tcwait.ForAll(
			tcwait.ForLog("Ethereum main loop is up").
				WithStartupTimeout(120 * time.Second).
				WithPollInterval(1 * time.Second),
		),
		User: "0:0", //otherwise in CI we get "permission denied" error, when trying to access data from mounted volume
		Cmd: []string{
			"--data-path=/opt/besu/execution-data",
			fmt.Sprintf("--genesis-file=%s/besu.json", GENERATED_DATA_DIR_INSIDE_CONTAINER),
			fmt.Sprintf("--network-id=%d", g.chainConfg.ChainID),
			"--host-allowlist=*",
			"--rpc-http-enabled=true",
			"--rpc-http-host=0.0.0.0",
			fmt.Sprintf("--rpc-http-port=%s", TX_GETH_HTTP_PORT),
			"--rpc-http-api=ADMIN,CLIQUE,ETH,NET,DEBUG,TXPOOL,ENGINE,TRACE,WEB3",
			"--rpc-http-cors-origins=*",
			"--rpc-ws-enabled=true",
			"--rpc-ws-host=0.0.0.0",
			fmt.Sprintf("--rpc-ws-port=%s", TX_GETH_WS_PORT),
			"--rpc-ws-api=ADMIN,CLIQUE,ETH,NET,DEBUG,TXPOOL,ENGINE,TRACE,WEB3",
			"--engine-rpc-enabled=true",
			fmt.Sprintf("--engine-jwt-secret=%s", JWT_SECRET_FILE_LOCATION_INSIDE_CONTAINER),
			"--engine-host-allowlist=*",
			fmt.Sprintf("--engine-rpc-port=%s", ETH2_EXECUTION_PORT),
			"--sync-mode=FULL",
			"--data-storage-format=BONSAI",
			// "--logging=DEBUG",
			"--rpc-tx-feecap=0",
		},
		Env: map[string]string{
			"JAVA_OPTS": "-agentlib:jdwp=transport=dt_socket,server=y,suspend=n",
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

func (g Besu) WaitUntilChainIsReady(waitTime time.Duration) error {
	waitForFirstBlock := tcwait.NewLogStrategy("Imported #1").WithPollInterval(1 * time.Second).WithStartupTimeout(waitTime)
	return waitForFirstBlock.WaitUntilReady(context.Background(), *g.GetContainer())
}

func (g *Besu) GetContainerType() ContainerType {
	return ContainerType_Besu
}
