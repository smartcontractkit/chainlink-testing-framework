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

const defaultBesuPosImage = "hyperledger/besu:24.1"
const defaultBesuPoWImage = "hyperledger/besu:22.1"

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

func NewBesuPos(networks []string, chainConfg *EthereumChainConfig, generatedDataHostDir string, consensusLayer ConsensusLayer, opts ...EnvComponentOption) (*Besu, error) {
	parts := strings.Split(defaultBesuPosImage, ":")
	g := &Besu{
		EnvComponent: EnvComponent{
			ContainerName:    fmt.Sprintf("%s-%s", "besu-pos", uuid.NewString()[0:8]),
			Networks:         networks,
			ContainerImage:   parts[0],
			ContainerVersion: parts[1],
		},
		chainConfg:     chainConfg,
		posSettings:    posSettings{generatedDataHostDir: generatedDataHostDir},
		consensusLayer: consensusLayer,
		l:              logging.GetTestLogger(nil),
	}
	g.SetDefaultHooks()
	for _, opt := range opts {
		opt(&g.EnvComponent)
	}

	// if the internal docker repo is set then add it to the version
	g.EnvComponent.ContainerImage = mirror.AddMirrorToImageIfSet(g.EnvComponent.ContainerImage)

	return g, nil
}

func NewBesuPow(networks []string, chainConfg *EthereumChainConfig, opts ...EnvComponentOption) (*Besu, error) {
	parts := strings.Split(defaultBesuPoWImage, ":")
	g := &Besu{
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

func (g *Besu) WithTestInstance(t *testing.T) ExecutionClient {
	g.l = logging.GetTestLogger(t)
	g.t = t
	return g
}

func (g *Besu) StartContainer() (blockchain.EVMNetwork, error) {
	var r *tc.ContainerRequest
	var err error

	if g.GetEthereumVersion() == EthereumVersion_Eth1 {
		r, err = g.getPowContainerRequest()

	} else {
		r, err = g.getPosContainerRequest()
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
	if g.consensusLayer != "" {
		return EthereumVersion_Eth2
	}

	return EthereumVersion_Eth1
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

func (g *Besu) getPosContainerRequest() (*tc.ContainerRequest, error) {
	return &tc.ContainerRequest{
		Name:     g.ContainerName,
		Image:    g.GetImageWithVersion(),
		Networks: g.Networks,
		// ImagePlatform: "linux/x86_64", //don't even try this on Apple Silicon, the node won't start due to JVM error
		ExposedPorts: []string{NatPortFormat(DEFAULT_EVM_NODE_HTTP_PORT), NatPortFormat(DEFAULT_EVM_NODE_WS_PORT), NatPortFormat(ETH2_EXECUTION_PORT)},
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
			fmt.Sprintf("--rpc-http-port=%s", DEFAULT_EVM_NODE_HTTP_PORT),
			"--rpc-http-api=ADMIN,CLIQUE,ETH,NET,DEBUG,TXPOOL,ENGINE,TRACE,WEB3",
			"--rpc-http-cors-origins=*",
			"--rpc-ws-enabled=true",
			"--rpc-ws-host=0.0.0.0",
			fmt.Sprintf("--rpc-ws-port=%s", DEFAULT_EVM_NODE_WS_PORT),
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

func (g *Besu) getPowContainerRequest() (*tc.ContainerRequest, error) {
	err := g.prepareFilesAndDirs()
	if err != nil {
		return nil, err
	}

	return &tc.ContainerRequest{
		Name:  g.ContainerName,
		Image: g.GetImageWithVersion(),
		ExposedPorts: []string{
			NatPortFormat(DEFAULT_EVM_NODE_HTTP_PORT),
			NatPortFormat(DEFAULT_EVM_NODE_WS_PORT)},
		Networks: g.Networks,
		WaitingFor: tcwait.ForAll(
			tcwait.ForLog("WebSocketService | Websocket service started"),
			NewWebSocketStrategy(NatPort(DEFAULT_EVM_NODE_WS_PORT), g.l),
			NewHTTPStrategy("/", NatPort(DEFAULT_EVM_NODE_HTTP_PORT)).WithStatusCode(201),
		),
		Entrypoint: []string{
			"besu",
			"--genesis-file", "/opt/besu/nodedata/genesis.json",
			"--host-allowlist", "*",
			"--rpc-http-enabled",
			"--rpc-http-cors-origins", "*",
			"--rpc-http-api", "ADMIN,DEBUG,WEB3,ETH,TXPOOL,CLIQUE,MINER,NET",
			"--rpc-http-host", "0.0.0.0",
			fmt.Sprintf("--rpc-http-port=%s", DEFAULT_EVM_NODE_HTTP_PORT),
			"--rpc-ws-enabled",
			"--rpc-ws-api", "ADMIN,DEBUG,WEB3,ETH,TXPOOL,CLIQUE,MINER,NET",
			"--rpc-ws-host", "0.0.0.0",
			fmt.Sprintf("--rpc-ws-port=%s", DEFAULT_EVM_NODE_WS_PORT),
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

func (g *Besu) prepareFilesAndDirs() error {
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

	extraAddresses := []string{}
	for _, addr := range g.chainConfg.AddressesToFund {
		extraAddresses = append(extraAddresses, strings.Replace(addr, "0x", "", 1))
	}

	accounts = append(accounts, extraAddresses...)
	accounts, err = deduplicateAddresses(g.l, accounts)
	if err != nil {
		return err
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
