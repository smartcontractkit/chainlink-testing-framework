package test_env

import (
	"fmt"
	"os"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/google/uuid"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/mirror"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/templates"
	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"
)

func NewBesuEth1(networks []string, chainConfg *EthereumChainConfig, opts ...EnvComponentOption) (*Besu, error) {
	parts := strings.Split(defaultBesuEth1Image, ":")
	g := &Besu{
		EnvComponent: EnvComponent{
			ContainerName:    fmt.Sprintf("%s-%s", "besu-eth1", uuid.NewString()[0:8]),
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

func (g *Besu) prepareEth1FilesAndDirs() error {
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
