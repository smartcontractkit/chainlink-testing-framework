package test_env

import (
	"context"
	"fmt"
	"os"
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
	GO_CLIENT_IMAGE_TAG      = "v1.13.4"
	GETH_ETH2_EXECUTION_PORT = "8551"
)

type GethGenesis struct {
	EnvComponent
	ExecutionDir string
	l            zerolog.Logger
	t            *testing.T
}

type Geth2 struct {
	EnvComponent
	ExternalHttpUrl string
	InternalHttpUrl string
	ExternalWsUrl   string
	InternalWsUrl   string
	ExecutionURL    string
	ExecutionDir    string
	l               zerolog.Logger
	t               *testing.T
}

func NewEth1Genesis(networks []string, executionDir string, opts ...EnvComponentOption) *GethGenesis {
	g := &GethGenesis{
		EnvComponent: EnvComponent{
			ContainerName: fmt.Sprintf("%s-%s", "geth-genesis", uuid.NewString()[0:8]),
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

func (g *GethGenesis) WithTestLogger(t *testing.T) *GethGenesis {
	g.l = logging.GetTestLogger(t)
	g.t = t
	return g
}

func (g *GethGenesis) StartContainer() error {
	r, err := g.getContainerRequest(g.Networks)
	if err != nil {
		return err
	}

	l := tc.Logger
	if g.t != nil {
		l = logging.CustomT{
			T: g.t,
			L: g.l,
		}
	}
	_, err = docker.StartContainerWithRetry(g.l, tc.GenericContainerRequest{
		ContainerRequest: *r,
		Reuse:            true,
		Started:          true,
		Logger:           l,
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
			"/execution/genesis.json",
		},
		Mounts: tc.ContainerMounts{
			tc.ContainerMount{
				Source: tc.GenericBindMountSource{
					HostPath: g.ExecutionDir,
				},
				Target: EXECUTION_DIRECTORY,
			},
		},
	}, nil
}

func NewGeth2(networks []string, executionDir string, opts ...EnvComponentOption) *Geth2 {
	g := &Geth2{
		EnvComponent: EnvComponent{
			ContainerName: fmt.Sprintf("%s-%s", "geth2", uuid.NewString()[0:8]),
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

func (g *Geth2) WithTestLogger(t *testing.T) *Geth2 {
	g.l = logging.GetTestLogger(t)
	g.t = t
	return g
}

func (g *Geth2) StartContainer() (blockchain.EVMNetwork, InternalDockerUrls, error) {
	r, err := g.getContainerRequest(g.Networks)
	if err != nil {
		return blockchain.EVMNetwork{}, InternalDockerUrls{}, err
	}

	l := tc.Logger
	if g.t != nil {
		l = logging.CustomT{
			T: g.t,
			L: g.l,
		}
	}
	ct, err := docker.StartContainerWithRetry(g.l, tc.GenericContainerRequest{
		ContainerRequest: *r,
		Reuse:            true,
		Started:          true,
		Logger:           l,
	})
	if err != nil {
		return blockchain.EVMNetwork{}, InternalDockerUrls{}, errors.Wrapf(err, "cannot start geth container")
	}

	host, err := GetHost(context.Background(), ct)
	if err != nil {
		return blockchain.EVMNetwork{}, InternalDockerUrls{}, err
	}
	if err != nil {
		return blockchain.EVMNetwork{}, InternalDockerUrls{}, err
	}
	httpPort, err := ct.MappedPort(context.Background(), NatPort(TX_GETH_HTTP_PORT))
	if err != nil {
		return blockchain.EVMNetwork{}, InternalDockerUrls{}, err
	}
	wsPort, err := ct.MappedPort(context.Background(), NatPort(TX_GETH_WS_PORT))
	if err != nil {
		return blockchain.EVMNetwork{}, InternalDockerUrls{}, err
	}
	_, err = ct.MappedPort(context.Background(), NatPort(GETH_ETH2_EXECUTION_PORT))
	if err != nil {
		return blockchain.EVMNetwork{}, InternalDockerUrls{}, err
	}

	g.Container = ct
	g.ExternalHttpUrl = fmt.Sprintf("http://%s:%s", host, httpPort.Port())
	g.InternalHttpUrl = fmt.Sprintf("http://%s:%s", g.ContainerName, TX_GETH_HTTP_PORT)
	g.ExternalWsUrl = fmt.Sprintf("ws://%s:%s", host, wsPort.Port())
	g.InternalWsUrl = fmt.Sprintf("ws://%s:%s", g.ContainerName, TX_GETH_WS_PORT)
	g.ExecutionURL = fmt.Sprintf("http://%s:%s", g.ContainerName, GETH_ETH2_EXECUTION_PORT)

	networkConfig := blockchain.SimulatedEVMNetwork
	networkConfig.Name = "geth"
	networkConfig.URLs = []string{g.ExternalWsUrl}
	networkConfig.HTTPURLs = []string{g.ExternalHttpUrl}

	internalDockerUrls := InternalDockerUrls{
		HttpUrl: g.InternalHttpUrl,
		WsUrl:   g.InternalWsUrl,
	}

	g.l.Info().Str("containerName", g.ContainerName).
		Msg("Started Geth2 container")

	return networkConfig, internalDockerUrls, nil
}

func (g *Geth2) getContainerRequest(networks []string) (*tc.ContainerRequest, error) {
	err := os.WriteFile(g.ExecutionDir+"/password.txt", []byte(""), 0600)
	if err != nil {
		return nil, err
	}

	key1File, err := os.CreateTemp(g.ExecutionDir+"/keystore", "UTC--2022-08-19T17-38-31.257380510Z--123463a4b065722e99115d6c222f267d9cabb524")
	if err != nil {
		return nil, err
	}
	_, err = key1File.WriteString(`{"address":"123463a4b065722e99115d6c222f267d9cabb524","crypto":{"cipher":"aes-128-ctr","ciphertext":"93b90389b855889b9f91c89fd15b9bd2ae95b06fe8e2314009fc88859fc6fde9","cipherparams":{"iv":"9dc2eff7967505f0e6a40264d1511742"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"c07503bb1b66083c37527cd8f06f8c7c1443d4c724767f625743bd47ae6179a4"},"mac":"6d359be5d6c432d5bbb859484009a4bf1bd71b76e89420c380bd0593ce25a817"},"id":"622df904-0bb1-4236-b254-f1b8dfdff1ec","version":3}`)
	if err != nil {
		return nil, err
	}

	jwtSecret, err := os.CreateTemp(g.ExecutionDir, "jwtsecret")
	if err != nil {
		return nil, err
	}
	_, err = jwtSecret.WriteString("0xfad2709d0bb03bf0e8ba3c99bea194575d3e98863133d1af638ed056d1d59345")
	if err != nil {
		return nil, err
	}
	secretKey, err := os.CreateTemp(g.ExecutionDir, "sk.json")
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
		ExposedPorts:    []string{NatPortFormat(TX_GETH_HTTP_PORT), NatPortFormat(TX_GETH_WS_PORT), NatPortFormat(GETH_ETH2_EXECUTION_PORT)},
		WaitingFor: tcwait.ForAll(
			NewHTTPStrategy("/", NatPort(TX_GETH_HTTP_PORT)),
			tcwait.ForLog("WebSocket enabled"),
			tcwait.ForLog("Started P2P networking").
				WithStartupTimeout(120*time.Second).
				WithPollInterval(1*time.Second),
			NewWebSocketStrategy(NatPort(TX_GETH_WS_PORT), g.l),
		),
		Cmd: []string{"--http",
			"--http.api=eth,net,web3",
			"--http.addr=0.0.0.0",
			"--http.corsdomain=*",
			fmt.Sprintf("--http.port=%s", TX_GETH_HTTP_PORT),
			"--ws",
			"--ws.api=eth,net,web3",
			"--ws.addr=0.0.0.0",
			"--ws.origins=*",
			fmt.Sprintf("--ws.port=%s", TX_GETH_WS_PORT),
			"--authrpc.vhosts=*",
			"--authrpc.addr=0.0.0.0",
			"--authrpc.jwtsecret=/execution/jwtsecret",
			"--datadir=/execution",
			"--allow-insecure-unlock",
			"--unlock=0x123463a4b065722e99115d6c222f267d9cabb524",
			"--password=/execution/password.txt",
			"--nodiscover",
			"--syncmode=full",
		},
		Files: []tc.ContainerFile{
			{
				HostFilePath:      jwtSecret.Name(),
				ContainerFilePath: "/execution/jwtsecret",
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
				Target: EXECUTION_DIRECTORY,
			},
		},
	}, nil
}
