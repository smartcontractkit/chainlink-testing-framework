package test_env

import (
	"fmt"
	"os"
	"strings"
	"time"

	config_types "github.com/smartcontractkit/chainlink-testing-framework/lib/config/types"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/google/uuid"
	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/config"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/docker/ethereum"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/templates"
)

// NewBesuEth1 starts a new Besu Eth1 node running in Docker
func NewBesuEth1(networks []string, chainConfig *config.EthereumChainConfig, opts ...EnvComponentOption) (*Besu, error) {
	parts := strings.Split(ethereum.DefaultBesuEth1Image, ":")
	g := &Besu{
		EnvComponent: EnvComponent{
			ContainerName:    fmt.Sprintf("%s-%s", "besu-eth1", uuid.NewString()[0:8]),
			Networks:         networks,
			ContainerImage:   parts[0],
			ContainerVersion: parts[1],
			StartupTimeout:   2 * time.Minute,
		},
		chainConfig:     chainConfig,
		l:               logging.GetTestLogger(nil),
		ethereumVersion: config_types.EthereumVersion_Eth1,
	}
	g.SetDefaultHooks()
	for _, opt := range opts {
		opt(&g.EnvComponent)
	}

	if !g.WasRecreated {
		// set the container name again after applying functional options as version might have changed
		g.EnvComponent.ContainerName = fmt.Sprintf("%s-%s-%s", "besu-eth1", strings.Replace(g.ContainerVersion, ".", "_", -1), uuid.NewString()[0:8])
	}

	// if the internal docker repo is set then add it to the version
	g.EnvComponent.ContainerImage = mirror.AddMirrorToImageIfSet(g.EnvComponent.ContainerImage)

	return g, nil
}

func (g *Besu) getEth1ContainerRequest() (*tc.ContainerRequest, error) {
	err := g.prepareEth1FilesAndDirs()
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
			NewHTTPStrategy("/", NatPort(DEFAULT_EVM_NODE_HTTP_PORT)).WithStatusCode(201)).
			WithStartupTimeoutDefault(g.StartupTimeout),
		Entrypoint: []string{
			"besu",
			"--genesis-file", "/opt/besu/genesis/genesis.json",
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
			fmt.Sprintf("--network-id=%d", g.chainConfig.ChainID),
			fmt.Sprintf("--logging=%s", strings.ToUpper(g.LogLevel)),
		},
		Files: []tc.ContainerFile{
			{
				HostFilePath:      g.genesisPath,
				ContainerFilePath: "/opt/besu/genesis/genesis.json",
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

func (g *Besu) prepareEth1FilesAndDirs() error {
	keystorePath, err := os.MkdirTemp("", "keystore")
	if err != nil {
		return err
	}
	g.keystorePath = keystorePath

	toFund := g.chainConfig.AddressesToFund
	toFund = append(toFund, RootFundingAddr)
	generatedData, err := generateKeystoreAndExtraData(keystorePath, toFund)
	if err != nil {
		return err
	}

	g.accountAddr = generatedData.minerAccount.Address.Hex()

	err = os.WriteFile(g.keystorePath+"/password.txt", []byte(""), 0600)
	if err != nil {
		return err
	}

	genesisJsonStr, err := templates.BuildBesuGenesisJsonForNonDevChain(fmt.Sprint(g.chainConfig.ChainID),
		generatedData.accountsToFund, "0x")
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
