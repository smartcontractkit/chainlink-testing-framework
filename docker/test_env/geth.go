package test_env

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/docker"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/templates"
)

const (
	// RootFundingAddr is the static key that hardhat is using
	// https://hardhat.org/hardhat-runner/docs/getting-started
	// if you need more keys, keep them compatible, so we can swap Geth to Ganache/Hardhat in the future
	RootFundingAddr   = `0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266`
	RootFundingWallet = `{"address":"f39fd6e51aad88f6f4ce6ab8827279cfffb92266","crypto":{"cipher":"aes-128-ctr","ciphertext":"c36afd6e60b82d6844530bd6ab44dbc3b85a53e826c3a7f6fc6a75ce38c1e4c6","cipherparams":{"iv":"f69d2bb8cd0cb6274535656553b61806"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"80d5f5e38ba175b6b89acfc8ea62a6f163970504af301292377ff7baafedab53"},"mac":"f2ecec2c4d05aacc10eba5235354c2fcc3776824f81ec6de98022f704efbf065"},"id":"e5c124e9-e280-4b10-a27b-d7f3e516b408","version":3}`

	TX_GETH_HTTP_PORT = "8544"
	TX_GETH_WS_PORT   = "8545"
)

type InternalDockerUrls struct {
	HttpUrl string
	WsUrl   string
}

type Geth struct {
	EnvComponent
	ExternalHttpUrl string
	InternalHttpUrl string
	ExternalWsUrl   string
	InternalWsUrl   string
	l               zerolog.Logger
	t               *testing.T
}

func NewGeth(networks []string, opts ...EnvComponentOption) *Geth {
	g := &Geth{
		EnvComponent: EnvComponent{
			ContainerName: fmt.Sprintf("%s-%s", "geth", uuid.NewString()[0:8]),
			Networks:      networks,
		},
		l: log.Logger,
	}
	for _, opt := range opts {
		opt(&g.EnvComponent)
	}
	return g
}

func (g *Geth) WithTestLogger(t *testing.T) *Geth {
	g.l = logging.GetTestLogger(t)
	g.t = t
	return g
}

func (g *Geth) StartContainer() (blockchain.EVMNetwork, InternalDockerUrls, error) {
	r, _, _, err := g.getGethContainerRequest(g.Networks)
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
	host, err := ct.Host(context.Background())
	if err != nil {
		return blockchain.EVMNetwork{}, InternalDockerUrls{}, err
	}
	httpPort, err := ct.MappedPort(context.Background(), NatPort(TX_GETH_HTTP_PORT))
	if err != nil {
		return blockchain.EVMNetwork{}, InternalDockerUrls{}, err
	}
	wsPort, err := getUniqueWsPort(ct, TX_GETH_HTTP_PORT, TX_GETH_WS_PORT, g.l)
	if err != nil {
		return blockchain.EVMNetwork{}, InternalDockerUrls{}, err
	}

	g.Container = ct
	g.ExternalHttpUrl = fmt.Sprintf("http://%s:%s", host, httpPort.Port())
	g.InternalHttpUrl = fmt.Sprintf("http://%s:%s", g.ContainerName, TX_GETH_HTTP_PORT)
	g.ExternalWsUrl = fmt.Sprintf("ws://%s:%s", host, wsPort.Port())
	g.InternalWsUrl = fmt.Sprintf("ws://%s:%s", g.ContainerName, TX_GETH_WS_PORT)

	networkConfig := blockchain.SimulatedEVMNetwork
	networkConfig.Name = "geth"
	networkConfig.URLs = []string{g.ExternalWsUrl}
	networkConfig.HTTPURLs = []string{g.ExternalHttpUrl}

	internalDockerUrls := InternalDockerUrls{
		HttpUrl: g.InternalHttpUrl,
		WsUrl:   g.InternalWsUrl,
	}

	g.l.Info().Str("containerName", g.ContainerName).
		Str("internalHttpUrl", g.InternalHttpUrl).
		Str("externalHttpUrl", g.ExternalHttpUrl).
		Str("externalWsUrl", g.ExternalWsUrl).
		Str("internalWsUrl", g.InternalWsUrl).
		Msg("Started Geth container")

	return networkConfig, internalDockerUrls, nil
}

func getUniqueWsPort(ct tc.Container, httpInternalPort, wsInternalPort string, l zerolog.Logger) (nat.Port, error) {
	p, err := ct.Ports(context.Background())
	if err != nil {
		return "", err
	}
	l.Info().
		Interface("ports", p).
		Msg("Ports mapped")
	httpPorts := p[NatPort(httpInternalPort)]
	wsPorts := p[NatPort(wsInternalPort)]
	wsPort := NatPort(wsPorts[0].HostPort)
	if len(wsPorts) > 1 {
		if httpPorts[0].HostPort == wsPorts[0].HostPort || httpPorts[1].HostPort == wsPorts[0].HostPort {
			wsPort = NatPort(wsPorts[1].HostPort)
		}
	}
	return wsPort, nil
}

func (g *Geth) getGethContainerRequest(networks []string) (*tc.ContainerRequest, *keystore.KeyStore, *accounts.Account, error) {
	chainId := "1337"
	blocktime := "1"

	initScriptFile, err := os.CreateTemp("", "init_script")
	if err != nil {
		return nil, nil, nil, err
	}
	_, err = initScriptFile.WriteString(templates.InitGethScript)
	if err != nil {
		return nil, nil, nil, err
	}
	keystoreDir, err := os.MkdirTemp("", "keystore")
	if err != nil {
		return nil, nil, nil, err
	}
	// Create keystore and ethereum account
	ks := keystore.NewKeyStore(keystoreDir, keystore.StandardScryptN, keystore.StandardScryptP)
	account, err := ks.NewAccount("")
	if err != nil {
		return nil, ks, &account, err
	}
	genesisJsonStr, err := templates.GenesisJsonTemplate{
		ChainId:     chainId,
		AccountAddr: account.Address.Hex(),
	}.String()
	if err != nil {
		return nil, ks, &account, err
	}
	genesisFile, err := os.CreateTemp("", "genesis_json")
	if err != nil {
		return nil, ks, &account, err
	}
	_, err = genesisFile.WriteString(genesisJsonStr)
	if err != nil {
		return nil, ks, &account, err
	}
	key1File, err := os.CreateTemp(keystoreDir, "key1")
	if err != nil {
		return nil, ks, &account, err
	}
	_, err = key1File.WriteString(RootFundingWallet)
	if err != nil {
		return nil, ks, &account, err
	}
	configDir, err := os.MkdirTemp("", "config")
	if err != nil {
		return nil, ks, &account, err
	}
	err = os.WriteFile(configDir+"/password.txt", []byte(""), 0600)
	if err != nil {
		return nil, ks, &account, err
	}

	return &tc.ContainerRequest{
		Name:            g.ContainerName,
		AlwaysPullImage: true,
		Image:           "ethereum/client-go:v1.12.0",
		ExposedPorts:    []string{NatPortFormat(TX_GETH_HTTP_PORT), NatPortFormat(TX_GETH_WS_PORT)},
		Networks:        networks,
		WaitingFor: tcwait.ForAll(
			tcwait.NewHTTPStrategy("/").
				WithPort(NatPort(TX_GETH_HTTP_PORT)),
			tcwait.ForLog("WebSocket enabled"),
			tcwait.ForLog("Started P2P networking").
				WithStartupTimeout(120*time.Second).
				WithPollInterval(1*time.Second),
			NewWebSocketStrategy(NatPort(TX_GETH_WS_PORT), g.l),
		),
		Entrypoint: []string{"sh", "./root/init.sh",
			"--dev",
			"--password", "/root/config/password.txt",
			"--datadir",
			"/root/.ethereum/devchain",
			"--unlock",
			RootFundingAddr,
			"--mine",
			"--miner.etherbase",
			RootFundingAddr,
			"--ipcdisable",
			"--http",
			"--http.vhosts",
			"*",
			"--http.addr",
			"0.0.0.0",
			fmt.Sprintf("--http.port=%s", TX_GETH_HTTP_PORT),
			"--ws",
			"--ws.origins",
			"*",
			"--ws.addr",
			"0.0.0.0",
			"--ws.api", "admin,debug,web3,eth,txpool,personal,clique,miner,net",
			fmt.Sprintf("--ws.port=%s", TX_GETH_WS_PORT),
			"--graphql",
			"-graphql.corsdomain",
			"*",
			"--allow-insecure-unlock",
			"--rpc.allow-unprotected-txs",
			"--http.api",
			"eth,web3,debug",
			"--http.corsdomain",
			"*",
			"--vmdebug",
			fmt.Sprintf("--networkid=%s", chainId),
			"--rpc.txfeecap",
			"0",
			"--dev.period",
			blocktime,
		},
		Files: []tc.ContainerFile{
			{
				HostFilePath:      initScriptFile.Name(),
				ContainerFilePath: "/root/init.sh",
				FileMode:          0644,
			},
			{
				HostFilePath:      genesisFile.Name(),
				ContainerFilePath: "/root/genesis.json",
				FileMode:          0644,
			},
		},
		Mounts: tc.ContainerMounts{
			tc.ContainerMount{
				Source: tc.GenericBindMountSource{
					HostPath: keystoreDir,
				},
				Target: "/root/.ethereum/devchain/keystore/",
			},
			tc.ContainerMount{
				Source: tc.GenericBindMountSource{
					HostPath: configDir,
				},
				Target: "/root/config/",
			},
		},
	}, ks, &account, nil
}

type WebSocketStrategy struct {
	Port       nat.Port
	RetryDelay time.Duration
	timeout    time.Duration
	l          zerolog.Logger
}

func NewWebSocketStrategy(port nat.Port, l zerolog.Logger) *WebSocketStrategy {
	return &WebSocketStrategy{
		Port:       port,
		RetryDelay: 10 * time.Second,
		timeout:    2 * time.Minute,
	}
}

func (w *WebSocketStrategy) WithTimeout(timeout time.Duration) *WebSocketStrategy {
	w.timeout = timeout
	return w
}

func (w *WebSocketStrategy) WaitUntilReady(ctx context.Context, target tcwait.StrategyTarget) (err error) {
	var client *rpc.Client
	var host string
	ctx, cancel := context.WithTimeout(ctx, w.timeout)
	defer cancel()
	i := 0
	for {
		host, err = target.Host(ctx)
		if err != nil {
			w.l.Error().Msg("Failed to get the target host")
			return err
		}
		wsPort, err := getUniqueWsPort(target.(tc.Container), TX_GETH_HTTP_PORT, w.Port.Port(), w.l)
		if err != nil {
			return err
		}

		url := fmt.Sprintf("ws://%s:%s", host, wsPort.Port())
		w.l.Info().Msgf("Attempting to dial %s", url)
		client, err = rpc.DialContext(ctx, url)
		if err == nil {
			client.Close()
			w.l.Info().Msg("WebSocket rpc port is ready")
			return nil
		}
		if client != nil {
			client.Close() // Close client if DialContext failed
			client = nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(w.RetryDelay):
			i++
			w.l.Info().Msgf("WebSocket attempt %d failed: %s. Retrying...", i, err)
		}
	}
}
