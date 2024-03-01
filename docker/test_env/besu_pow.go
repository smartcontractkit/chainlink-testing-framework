package test_env

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/docker"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/mirror"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/templates"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/testcontext"
)

const defaultBesuPoWImage = "hyperledger/besu:22.1"

type BesuPoW struct {
	EnvComponent
	ExternalHttpUrl      string
	InternalHttpUrl      string
	ExternalWsUrl        string
	InternalWsUrl        string
	InternalExecutionURL string
	ExternalExecutionURL string
	chainConfg           *EthereumChainConfig
	l                    zerolog.Logger
	t                    *testing.T
	genesisPath          string
	rootPath             string
	keystorePath         string
	accountAddr          string
}

func NewBesuPow(networks []string, chainConfg *EthereumChainConfig, opts ...EnvComponentOption) (*BesuPoW, error) {
	parts := strings.Split(defaultBesuPoWImage, ":")
	g := &BesuPoW{
		EnvComponent: EnvComponent{
			ContainerName:    fmt.Sprintf("%s-%s", "besu-pow", uuid.NewString()[0:8]),
			Networks:         networks,
			ContainerImage:   parts[0],
			ContainerVersion: parts[1],
		},
		chainConfg: chainConfg,
		l:          logging.GetTestLogger(nil),
	}
	g.SetDefaultHooks()
	for _, opt := range opts {
		opt(&g.EnvComponent)
	}

	// if the internal docker repo is set then add it to the version
	g.EnvComponent.ContainerImage = mirror.AddMirrorToImageIfSet(g.EnvComponent.ContainerImage)

	return g, nil
}

func (g *BesuPoW) WithTestInstance(t *testing.T) ExecutionClient {
	g.l = logging.GetTestLogger(t)
	g.t = t
	return g
}

func (g *BesuPoW) StartContainer() (blockchain.EVMNetwork, error) {
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
		return blockchain.EVMNetwork{}, fmt.Errorf("cannot start Besu container: %w", err)
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

	g.Container = ct
	g.ExternalHttpUrl = FormatHttpUrl(host, httpPort.Port())
	g.InternalHttpUrl = FormatHttpUrl(g.ContainerName, TX_GETH_HTTP_PORT)
	g.ExternalWsUrl = FormatWsUrl(host, wsPort.Port())
	g.InternalWsUrl = FormatWsUrl(g.ContainerName, TX_GETH_WS_PORT)

	networkConfig := blockchain.SimulatedEVMNetwork
	networkConfig.Name = "Simulated Ethereum-PoW (besu)"
	networkConfig.URLs = []string{g.ExternalWsUrl}
	networkConfig.HTTPURLs = []string{g.ExternalHttpUrl}
	networkConfig.GasEstimationBuffer = 10_000_000_000

	g.l.Info().Str("containerName", g.ContainerName).
		Msg("Started Besu container")

	return networkConfig, nil
}

func (g *BesuPoW) GetInternalExecutionURL() string {
	return g.InternalExecutionURL
}

func (g *BesuPoW) GetExternalExecutionURL() string {
	return g.ExternalExecutionURL
}

func (g *BesuPoW) GetInternalHttpUrl() string {
	return g.InternalHttpUrl
}

func (g *BesuPoW) GetInternalWsUrl() string {
	return g.InternalWsUrl
}

func (g *BesuPoW) GetExternalHttpUrl() string {
	return g.ExternalHttpUrl
}

func (g *BesuPoW) GetExternalWsUrl() string {
	return g.ExternalWsUrl
}

func (g *BesuPoW) GetContainerName() string {
	return g.ContainerName
}

func (g *BesuPoW) GetContainer() *tc.Container {
	return &g.Container
}

func (g *BesuPoW) createMountDirs() error {
	keystorePath, err := os.MkdirTemp("", "keystore")
	if err != nil {
		return err
	}
	g.keystorePath = keystorePath

	// Create keystore and ethereum account
	ks := keystore.NewKeyStore(g.keystorePath, keystore.StandardScryptN, keystore.StandardScryptP)
	account, err := ks.NewAccount("")
	if err != nil {
		return err
	}

	g.accountAddr = account.Address.Hex()
	addr := strings.Replace(account.Address.Hex(), "0x", "", 1)
	FundingAddresses[addr] = ""

	i := 1
	var accounts []string
	for addr, v := range FundingAddresses {
		if v == "" {
			continue
		}
		f, err := os.Create(fmt.Sprintf("%s/%s", g.keystorePath, fmt.Sprintf("key%d", i)))
		if err != nil {
			return err
		}
		_, err = f.WriteString(v)
		if err != nil {
			return err
		}
		i++
		accounts = append(accounts, addr)
	}
	err = os.WriteFile(g.keystorePath+"/password.txt", []byte(""), 0600)
	if err != nil {
		return err
	}

	genesisJsonStr, err := templates.BuildBesuGenesisJsonForNonDevChain(fmt.Sprint(g.chainConfg.ChainID),
		accounts, "0x")
	if err != nil {
		return err
	}
	f, err := os.CreateTemp("", "genesis_json")
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(genesisJsonStr)
	if err != nil {
		return err
	}

	g.genesisPath = f.Name()

	configDir, err := os.MkdirTemp("", "config")
	if err != nil {
		return err
	}
	g.rootPath = configDir

	return nil
}

func (g *BesuPoW) getContainerRequest() (*tc.ContainerRequest, error) {
	err := g.createMountDirs()
	if err != nil {
		return nil, err
	}

	return &tc.ContainerRequest{
		Name:  g.ContainerName,
		Image: g.GetImageWithVersion(),
		ExposedPorts: []string{
			NatPortFormat(TX_GETH_HTTP_PORT),
			NatPortFormat(TX_GETH_WS_PORT)},
		Networks: g.Networks,
		WaitingFor: tcwait.ForAll(
			tcwait.ForLog("WebSocketService | Websocket service started"),
			NewWebSocketStrategy(NatPort(TX_GETH_WS_PORT), g.l),
			NewHTTPStrategy("/", NatPort(TX_GETH_HTTP_PORT)).WithStatusCode(201),
		),
		Entrypoint: []string{
			"besu",
			"--genesis-file", "/opt/besu/nodedata/genesis.json",
			"--host-allowlist", "*",
			"--rpc-http-enabled",
			"--rpc-http-cors-origins", "*",
			"--rpc-http-api", "ADMIN,DEBUG,WEB3,ETH,TXPOOL,CLIQUE,MINER,NET",
			"--rpc-http-host", "0.0.0.0",
			fmt.Sprintf("--rpc-http-port=%s", TX_GETH_HTTP_PORT),
			"--rpc-ws-enabled",
			"--rpc-ws-api", "ADMIN,DEBUG,WEB3,ETH,TXPOOL,CLIQUE,MINER,NET",
			"--rpc-ws-host", "0.0.0.0",
			fmt.Sprintf("--rpc-ws-port=%s", TX_GETH_WS_PORT),
			"--miner-enabled=true",
			"--miner-coinbase", RootFundingAddr,
			fmt.Sprintf("--network-id=%d", g.chainConfg.ChainID),
			"--logging=DEBUG",
		},
		Files: []tc.ContainerFile{
			{
				HostFilePath:      g.genesisPath,
				ContainerFilePath: "/opt/besu/nodedata/genesis.json",
				FileMode:          0644,
			},
		},
		HostConfigModifier: func(hostConfig *container.HostConfig) {
			hostConfig.Mounts = append(hostConfig.Mounts, mount.Mount{
				Type:     mount.TypeBind,
				Source:   g.keystorePath,
				Target:   "/opt/besu/nodedata/keystore/",
				ReadOnly: false,
			}, mount.Mount{
				Type:     mount.TypeBind,
				Source:   g.rootPath,
				Target:   "/opt/besu/nodedata/",
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

func (g *BesuPoW) WaitUntilChainIsReady(ctx context.Context, waitTime time.Duration) error {
	waitForFirstBlock := tcwait.NewLogStrategy("Imported #1").WithPollInterval(1 * time.Second).WithStartupTimeout(waitTime)
	return waitForFirstBlock.WaitUntilReady(ctx, *g.GetContainer())
}

func (g *BesuPoW) GetContainerType() ContainerType {
	return ContainerType_Besu
}
