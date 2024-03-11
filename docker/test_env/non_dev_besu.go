package test_env

import (
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/mirror"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/templates"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/testcontext"
)

const (
	defaultNonDevBesuImage = "hyperledger/besu:24.1"
)

type PrivateBesuChain struct {
	PrimaryNode    *NonDevBesuNode
	Nodes          []*NonDevBesuNode
	NetworkConfig  *blockchain.EVMNetwork
	DockerNetworks []string
}

func NewPrivateBesuChain(networkCfg *blockchain.EVMNetwork, dockerNetworks []string) PrivateChain {
	evmChain := &PrivateBesuChain{
		NetworkConfig:  networkCfg,
		DockerNetworks: dockerNetworks,
	}
	evmChain.PrimaryNode = NewNonDevBesuNode(dockerNetworks, networkCfg)
	evmChain.Nodes = []*NonDevBesuNode{evmChain.PrimaryNode}
	return evmChain
}

func (p *PrivateBesuChain) GetPrimaryNode() NonDevNode {
	return p.PrimaryNode
}

func (p *PrivateBesuChain) GetNodes() []NonDevNode {
	nodes := make([]NonDevNode, 0)
	for _, node := range p.Nodes {
		nodes = append(nodes, node)
	}
	return nodes
}

func (p *PrivateBesuChain) GetNetworkConfig() *blockchain.EVMNetwork {
	return p.NetworkConfig
}

func (p *PrivateBesuChain) GetDockerNetworks() []string {
	return p.DockerNetworks
}

type NonDevBesuNode struct {
	EnvComponent
	Config          gethTxNodeConfig
	ExternalHttpUrl string
	InternalHttpUrl string
	ExternalWsUrl   string
	InternalWsUrl   string
	EVMClient       blockchain.EVMClient
	EthClient       *ethclient.Client
	t               *testing.T
	l               zerolog.Logger
}

func NewNonDevBesuNode(networks []string, networkCfg *blockchain.EVMNetwork) *NonDevBesuNode {
	n := &NonDevBesuNode{
		Config: gethTxNodeConfig{
			chainId:    strconv.FormatInt(networkCfg.ChainID, 10),
			networkCfg: networkCfg,
		},
		EnvComponent: EnvComponent{
			ContainerName: fmt.Sprintf("%s-%s",
				strings.ReplaceAll(networkCfg.Name, " ", "_"), uuid.NewString()[0:3]),
			Networks: networks,
		},
	}
	n.SetDefaultHooks()

	return n
}

func (g *NonDevBesuNode) WithTestInstance(t *testing.T) NonDevNode {
	g.t = t
	g.l = logging.GetTestLogger(t)
	return g
}

func (g *NonDevBesuNode) GetInternalHttpUrl() string {
	return g.InternalHttpUrl
}

func (g *NonDevBesuNode) GetInternalWsUrl() string {
	return g.InternalWsUrl
}

func (g *NonDevBesuNode) GetEVMClient() blockchain.EVMClient {
	return g.EVMClient
}

func (g *NonDevBesuNode) createMountDirs() error {
	keystorePath, err := os.MkdirTemp("", "keystore")
	if err != nil {
		return err
	}
	g.Config.keystorePath = keystorePath

	// Create keystore and ethereum account
	ks := keystore.NewKeyStore(g.Config.keystorePath, keystore.StandardScryptN, keystore.StandardScryptP)
	account, err := ks.NewAccount("")
	if err != nil {
		return err
	}

	g.Config.accountAddr = account.Address.Hex()
	addr := strings.Replace(account.Address.Hex(), "0x", "", 1)
	FundingAddresses[addr] = ""
	signerBytes, err := hex.DecodeString(addr)
	if err != nil {
		fmt.Println("Error decoding signer address:", err)
		return err
	}

	zeroBytes := make([]byte, 32)                      // Create 32 zero bytes
	extradata := append(zeroBytes, signerBytes...)     // Concatenate zero bytes and signer address
	extradata = append(extradata, make([]byte, 65)...) // Concatenate 65 more zero bytes

	fmt.Printf("Encoded extradata: 0x%s\n", hex.EncodeToString(extradata))

	i := 1
	var accounts []string
	for addr, v := range FundingAddresses {
		if v == "" {
			continue
		}
		f, err := os.Create(fmt.Sprintf("%s/%s", g.Config.keystorePath, fmt.Sprintf("key%d", i)))
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
	err = os.WriteFile(g.Config.keystorePath+"/password.txt", []byte(""), 0600)
	if err != nil {
		return err
	}

	genesisJsonStr, err := templates.BuildBesuGenesisJsonForNonDevChain(g.Config.chainId,
		accounts,
		fmt.Sprintf("0x%s", hex.EncodeToString(extradata)))
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

	g.Config.genesisPath = f.Name()

	configDir, err := os.MkdirTemp("", "config")
	if err != nil {
		return err
	}
	g.Config.rootPath = configDir

	return nil
}

func (g *NonDevBesuNode) ConnectToClient() error {
	ct := g.Container
	if ct == nil {
		return fmt.Errorf("container not started")
	}
	host, err := GetHost(testcontext.Get(g.t), ct)
	if err != nil {
		return err
	}
	port := NatPort(TX_GETH_HTTP_PORT)
	httpPort, err := ct.MappedPort(testcontext.Get(g.t), port)
	if err != nil {
		return err
	}
	port = NatPort(TX_NON_DEV_GETH_WS_PORT)
	wsPort, err := ct.MappedPort(testcontext.Get(g.t), port)
	if err != nil {
		return err
	}
	g.ExternalHttpUrl = fmt.Sprintf("http://%s:%s", host, httpPort.Port())
	g.InternalHttpUrl = fmt.Sprintf("http://%s:%s", g.ContainerName, TX_GETH_HTTP_PORT)
	g.ExternalWsUrl = fmt.Sprintf("ws://%s:%s", host, wsPort.Port())
	g.InternalWsUrl = fmt.Sprintf("ws://%s:%s", g.ContainerName, TX_NON_DEV_GETH_WS_PORT)

	networkConfig := g.Config.networkCfg
	networkConfig.URLs = []string{g.ExternalWsUrl}
	networkConfig.HTTPURLs = []string{g.ExternalHttpUrl}

	ec, err := blockchain.NewEVMClientFromNetwork(*networkConfig, g.l)
	if err != nil {
		return err
	}
	at, err := ec.BalanceAt(testcontext.Get(g.t), common.HexToAddress("0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266"))
	if err != nil {
		return err
	}
	fmt.Printf("balance: %s\n", at.String())
	g.EVMClient = ec
	// to make sure all the pending txs are done
	err = ec.WaitForEvents()
	if err != nil {
		return err
	}
	switch val := ec.(type) {
	case *blockchain.EthereumMultinodeClient:
		ethClient, ok := val.Clients[0].(*blockchain.EthereumClient)
		if !ok {
			return fmt.Errorf("could not get blockchain.EthereumClient from %+v", val)
		}
		g.EthClient = ethClient.Client
	default:
		return fmt.Errorf("%+v not supported for geth", val)
	}
	return err
}

func (g *NonDevBesuNode) Start() error {
	err := g.createMountDirs()
	if err != nil {
		return err
	}
	l := logging.GetTestContainersGoTestLogger(g.t)

	// Besu Bootnode setup: BEGIN
	// Generate public key for besu bootnode
	crbn, err := g.getBesuBootNodeContainerRequest()
	if err != nil {
		return err
	}
	bootNode, err := tc.GenericContainer(testcontext.Get(g.t),
		tc.GenericContainerRequest{
			ContainerRequest: crbn,
			Started:          true,
			Reuse:            true,
			Logger:           l,
		})
	if err != nil {
		return err
	}

	err = g.exportBesuBootNodeAddress(bootNode)
	if err != nil {
		return err
	}
	// Besu Bootnode setup: END

	host, err := GetHost(testcontext.Get(g.t), bootNode)
	if err != nil {
		return err
	}
	r, err := bootNode.CopyFileFromContainer(testcontext.Get(g.t), "/opt/besu/nodedata/bootnodes")
	if err != nil {
		return err
	}
	defer r.Close()
	b, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	bootnodePubKey := strings.TrimPrefix(strings.TrimSpace(string(b)), "0x")
	g.Config.bootNodeURL = fmt.Sprintf("enode://%s@%s:%s", bootnodePubKey, host, BOOTNODE_PORT)

	fmt.Printf("Besu Bootnode URL: %s\n", g.Config.bootNodeURL)

	cr, err := g.getBesuContainerRequest()
	if err != nil {
		return err
	}
	var ct tc.Container
	ct, err = tc.GenericContainer(testcontext.Get(g.t),
		tc.GenericContainerRequest{
			ContainerRequest: cr,
			Started:          true,
			Reuse:            true,
		})
	if err != nil {
		return err
	}
	g.Container = ct
	return nil
}

func (g *NonDevBesuNode) getBesuBootNodeContainerRequest() (tc.ContainerRequest, error) {
	besuImage := mirror.AddMirrorToImageIfSet(defaultNonDevBesuImage)
	return tc.ContainerRequest{
		Name:         g.ContainerName + "-bootnode",
		Image:        besuImage,
		Networks:     g.Networks,
		ExposedPorts: []string{"30301/udp"},
		WaitingFor: tcwait.ForLog("PeerDiscoveryAgent | P2P peer discovery agent started and listening on").
			WithStartupTimeout(999 * time.Second).
			WithPollInterval(1 * time.Second),
		Cmd: []string{
			"--genesis-file",
			"/opt/besu/nodedata/genesis.json",
			"--data-path",
			"/opt/besu/nodedata",
			"--logging=INFO",
			"--p2p-port=30301",
			"--bootnodes",
		},
		Files: []tc.ContainerFile{
			{
				HostFilePath:      g.Config.genesisPath,
				ContainerFilePath: "/opt/besu/nodedata/genesis.json",
				FileMode:          0644,
			},
		},
		HostConfigModifier: func(hostConfig *container.HostConfig) {
			hostConfig.Mounts = append(hostConfig.Mounts, mount.Mount{
				Type:     mount.TypeBind,
				Source:   g.Config.rootPath,
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

func (g *NonDevBesuNode) exportBesuBootNodeAddress(bootNode tc.Container) (err error) {
	resCode, _, err := bootNode.Exec(testcontext.Get(g.t), []string{
		"besu",
		"--genesis-file", "/opt/besu/nodedata/genesis.json",
		"--data-path", "/opt/besu/nodedata",
		"public-key", "export",
		"--to=/opt/besu/nodedata/bootnodes",
	})
	fmt.Printf("Export besu bootnode address, process exitcode: %d\n", resCode)
	if err != nil {
		return err
	}
	return nil
}

func (g *NonDevBesuNode) getBesuContainerRequest() (tc.ContainerRequest, error) {
	besuImage := mirror.AddMirrorToImageIfSet(defaultNonDevBesuImage)
	return tc.ContainerRequest{
		Name:  g.ContainerName,
		Image: besuImage,
		ExposedPorts: []string{
			NatPortFormat(TX_GETH_HTTP_PORT),
			NatPortFormat(TX_NON_DEV_GETH_WS_PORT),
			"30303/tcp", "30303/udp"},
		Networks: g.Networks,
		WaitingFor: tcwait.ForAll(
			tcwait.ForLog("WebSocketService | Websocket service started"),
			NewWebSocketStrategy(NatPort(TX_NON_DEV_GETH_WS_PORT), g.l),
			NewHTTPStrategy("/", NatPort(TX_GETH_HTTP_PORT)).WithStatusCode(201),
		),
		Entrypoint: []string{
			"besu",
			"--genesis-file", "/opt/besu/nodedata/genesis.json",
			"--host-allowlist", "*",
			// "--sync-mode", "X_SNAP", // Requires at least 5 peers in X_SNAP mode
			fmt.Sprintf("--bootnodes=%s", g.Config.bootNodeURL),
			"--rpc-http-enabled",
			"--rpc-http-cors-origins", "*",
			"--rpc-http-api", "ADMIN,DEBUG,WEB3,ETH,TXPOOL,CLIQUE,MINER,NET",
			"--rpc-http-host", "0.0.0.0",
			fmt.Sprintf("--rpc-http-port=%s", TX_GETH_HTTP_PORT),
			"--rpc-ws-enabled",
			"--rpc-ws-api", "ADMIN,DEBUG,WEB3,ETH,TXPOOL,CLIQUE,MINER,NET",
			"--rpc-ws-host", "0.0.0.0",
			fmt.Sprintf("--rpc-ws-port=%s", TX_NON_DEV_GETH_WS_PORT),
			"--miner-enabled=true",
			"--miner-coinbase", g.Config.accountAddr,
			fmt.Sprintf("--network-id=%s", g.Config.chainId),
			"--logging=DEBUG",
		},
		Files: []tc.ContainerFile{
			{
				HostFilePath:      g.Config.genesisPath,
				ContainerFilePath: "/opt/besu/nodedata/genesis.json",
				FileMode:          0644,
			},
		},
		HostConfigModifier: func(hostConfig *container.HostConfig) {
			hostConfig.Mounts = append(hostConfig.Mounts, mount.Mount{
				Type:     mount.TypeBind,
				Source:   g.Config.keystorePath,
				Target:   "/opt/besu/nodedata/keystore/",
				ReadOnly: false,
			}, mount.Mount{
				Type:     mount.TypeBind,
				Source:   g.Config.rootPath,
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
