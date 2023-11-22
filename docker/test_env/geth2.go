package test_env

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/docker"
)

const (
	GO_CLIENT_IMAGE_TAG = "v1.13.4"
)

type GethGenesis struct {
	EnvComponent
	ExecutionDir string
	l            zerolog.Logger
}

type Geth2 struct {
	EnvComponent
	ExternalHttpUrl      string
	InternalHttpUrl      string
	ExternalWsUrl        string
	InternalWsUrl        string
	InternalExecutionURL string
	ExternalExecutionURL string
	ExecutionDir         string
	consensusLayer       ConsensusLayer
	l                    zerolog.Logger
}

func NewEth1Genesis(networks []string, executionDir string, opts ...EnvComponentOption) *GethGenesis {
	g := &GethGenesis{
		EnvComponent: EnvComponent{
			ContainerName: fmt.Sprintf("%s-%s", "geth-eth1-genesis", uuid.NewString()[0:8]),
			Networks:      networks,
		},
		ExecutionDir: executionDir,
		l:            log.Logger,
	}
	for _, opt := range opts {
		opt(&g.EnvComponent)
	}
	return g
}

func (g *GethGenesis) WithLogger(l zerolog.Logger) *GethGenesis {
	g.l = l
	return g
}

func (g *GethGenesis) StartContainer() error {
	r, err := g.getContainerRequest(g.Networks)
	if err != nil {
		return err
	}

	_, err = docker.StartContainerWithRetry(g.l, tc.GenericContainerRequest{
		ContainerRequest: *r,
		Reuse:            true,
		Started:          true,
		Logger:           &g.l,
	})
	if err != nil {
		return errors.Wrapf(err, "cannot start geth eth1 genesis container")
	}

	g.l.Info().Str("containerName", g.ContainerName).
		Msg("Started Geth Eth1 Genesis container")

	return nil
}

func (g *GethGenesis) getContainerRequest(networks []string) (*tc.ContainerRequest, error) {
	return &tc.ContainerRequest{
		Name:            g.ContainerName,
		AlwaysPullImage: true,
		Image:           fmt.Sprintf("ethereum/client-go:%s", GO_CLIENT_IMAGE_TAG),
		Networks:        networks,
		WaitingFor: tcwait.ForAll(
			tcwait.ForLog("Successfully wrote genesis state").
				WithStartupTimeout(120 * time.Second).
				WithPollInterval(1 * time.Second),
		),
		Cmd: []string{"--datadir=/execution",
			"init",
			eth1GenesisFile,
		},
		Mounts: tc.ContainerMounts{
			tc.ContainerMount{
				Source: tc.GenericBindMountSource{
					HostPath: g.ExecutionDir,
				},
				Target: CONTAINER_ETH2_EXECUTION_DIRECTORY,
			},
		},
	}, nil
}

func NewGeth2(networks []string, executionDir string, consensusLayer ConsensusLayer, opts ...EnvComponentOption) *Geth2 {
	g := &Geth2{
		EnvComponent: EnvComponent{
			ContainerName: fmt.Sprintf("%s-%s", "geth2", uuid.NewString()[0:8]),
			Networks:      networks,
		},
		ExecutionDir:   executionDir,
		consensusLayer: consensusLayer,
		l:              log.Logger,
	}
	for _, opt := range opts {
		opt(&g.EnvComponent)
	}
	return g
}

func (g *Geth2) WithLogger(l zerolog.Logger) *Geth2 {
	g.l = l
	return g
}

func (g *Geth2) StartContainer() (blockchain.EVMNetwork, error) {
	r, err := g.getContainerRequest(g.Networks)
	if err != nil {
		return blockchain.EVMNetwork{}, err
	}

	ct, err := docker.StartContainerWithRetry(g.l, tc.GenericContainerRequest{
		ContainerRequest: *r,
		Reuse:            true,
		Started:          true,
		Logger:           &g.l,
	})
	if err != nil {
		return blockchain.EVMNetwork{}, errors.Wrapf(err, "cannot start geth container")
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
	networkConfig.Name = fmt.Sprintf("geth-eth2-%s", g.consensusLayer)
	networkConfig.URLs = []string{g.ExternalWsUrl}
	networkConfig.HTTPURLs = []string{g.ExternalHttpUrl}

	g.l.Info().Str("containerName", g.ContainerName).
		Msg("Started Geth2 container")

	return networkConfig, nil
}

func (g *Geth2) GetInternalExecutionURL() string {
	return g.InternalExecutionURL
}

func (g *Geth2) GetExternalExecutionURL() string {
	return g.ExternalExecutionURL
}

func (g *Geth2) GetInternalHttpUrl() string {
	return g.InternalHttpUrl
}

func (g *Geth2) GetInternalWsUrl() string {
	return g.InternalWsUrl
}

func (g *Geth2) GetExternalHttpUrl() string {
	return g.ExternalHttpUrl
}

func (g *Geth2) GetExternalWsUrl() string {
	return g.ExternalWsUrl
}

func (g *Geth2) GetContainerName() string {
	return g.ContainerName
}

func (g *Geth2) GetContainer() *tc.Container {
	return &g.Container
}

func (g *Geth2) getContainerRequest(networks []string) (*tc.ContainerRequest, error) {
	passwordFile, err := os.CreateTemp("", "password.txt")
	if err != nil {
		return nil, err
	}

	key1File, err := os.CreateTemp("", "key1")
	if err != nil {
		return nil, err
	}
	_, err = key1File.WriteString(`{"address":"123463a4b065722e99115d6c222f267d9cabb524","crypto":{"cipher":"aes-128-ctr","ciphertext":"93b90389b855889b9f91c89fd15b9bd2ae95b06fe8e2314009fc88859fc6fde9","cipherparams":{"iv":"9dc2eff7967505f0e6a40264d1511742"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"c07503bb1b66083c37527cd8f06f8c7c1443d4c724767f625743bd47ae6179a4"},"mac":"6d359be5d6c432d5bbb859484009a4bf1bd71b76e89420c380bd0593ce25a817"},"id":"622df904-0bb1-4236-b254-f1b8dfdff1ec","version":3}`)
	if err != nil {
		return nil, err
	}

	jwtSecret, err := os.CreateTemp("", "jwtsecret")
	if err != nil {
		return nil, err
	}
	_, err = jwtSecret.WriteString("0xfad2709d0bb03bf0e8ba3c99bea194575d3e98863133d1af638ed056d1d59345")
	if err != nil {
		return nil, err
	}
	secretKey, err := os.CreateTemp("", "sk.json")
	if err != nil {
		return nil, err
	}
	_, err = secretKey.WriteString("2e0834786285daccd064ca17f1654f67b4aef298acbb82cef9ec422fb4975622")
	if err != nil {
		return nil, err
	}

	return &tc.ContainerRequest{
		Name:            g.ContainerName,
		AlwaysPullImage: true,
		Image:           fmt.Sprintf("ethereum/client-go:%s", GO_CLIENT_IMAGE_TAG),
		Networks:        networks,
		ExposedPorts:    []string{NatPortFormat(TX_GETH_HTTP_PORT), NatPortFormat(TX_GETH_WS_PORT), NatPortFormat(ETH2_EXECUTION_PORT)},
		WaitingFor: tcwait.ForAll(
			NewHTTPStrategy("/", NatPort(TX_GETH_HTTP_PORT)),
			tcwait.ForLog("WebSocket enabled"),
			tcwait.ForLog("Started P2P networking").
				WithStartupTimeout(120*time.Second).
				WithPollInterval(1*time.Second),
			NewWebSocketStrategy(NatPort(TX_GETH_WS_PORT), g.l),
		),
		Cmd: []string{"--http",
			"--http.api=eth,net,web3,debug",
			"--http.addr=0.0.0.0",
			"--http.corsdomain=*",
			"--http.vhosts=*",
			fmt.Sprintf("--http.port=%s", TX_GETH_HTTP_PORT),
			"--ws",
			"--ws.api=admin,debug,web3,eth,txpool,net",
			"--ws.addr=0.0.0.0",
			"--ws.origins=*",
			fmt.Sprintf("--ws.port=%s", TX_GETH_WS_PORT),
			"--authrpc.vhosts=*",
			"--authrpc.addr=0.0.0.0",
			"--authrpc.jwtsecret=" + jwtSecretFileLocation,
			"--datadir=/execution",
			"--rpc.allow-unprotected-txs",
			"--rpc.txfeecap=0",
			"--allow-insecure-unlock",
			"--unlock=0x123463a4b065722e99115d6c222f267d9cabb524",
			"--password=/execution/password.txt",
			"--nodiscover",
			"--syncmode=full",
			"--networkid=1337",
			"--graphql",
			"--graphql.corsdomain=*",
		},
		Files: []tc.ContainerFile{
			{
				HostFilePath:      passwordFile.Name(),
				ContainerFilePath: "/execution/password.txt",
				FileMode:          0644,
			},
			{
				HostFilePath:      jwtSecret.Name(),
				ContainerFilePath: jwtSecretFileLocation,
				FileMode:          0644,
			},
			{
				HostFilePath:      key1File.Name(),
				ContainerFilePath: "/execution/keystore/key1",
				FileMode:          0644,
			},
			{
				HostFilePath:      secretKey.Name(),
				ContainerFilePath: "/execution/sk.json",
				FileMode:          0644,
			},
		},
		Mounts: tc.ContainerMounts{
			tc.ContainerMount{
				Source: tc.GenericBindMountSource{
					HostPath: g.ExecutionDir,
				},
				Target: CONTAINER_ETH2_EXECUTION_DIRECTORY,
			},
		},
	}, nil
}

func (g Geth2) WaitUntilChainIsReady(waitTime time.Duration) error {
	waitForFirstBlock := tcwait.NewLogStrategy("Chain head was updated").WithPollInterval(1 * time.Second).WithStartupTimeout(waitTime)
	return waitForFirstBlock.WaitUntilReady(context.Background(), *g.GetContainer())
}
