package test_env

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/go-connections/nat"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/docker"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/mirror"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/templates"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/testcontext"
)

const (
	// RootFundingAddr is the static key that hardhat is using
	// https://hardhat.org/hardhat-runner/docs/getting-started
	// if you need more keys, keep them compatible, so we can swap Geth to Ganache/Hardhat in the future
	RootFundingAddr   = `0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266`
	RootFundingWallet = `{"address":"f39fd6e51aad88f6f4ce6ab8827279cfffb92266","crypto":{"cipher":"aes-128-ctr","ciphertext":"c36afd6e60b82d6844530bd6ab44dbc3b85a53e826c3a7f6fc6a75ce38c1e4c6","cipherparams":{"iv":"f69d2bb8cd0cb6274535656553b61806"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"80d5f5e38ba175b6b89acfc8ea62a6f163970504af301292377ff7baafedab53"},"mac":"f2ecec2c4d05aacc10eba5235354c2fcc3776824f81ec6de98022f704efbf065"},"id":"e5c124e9-e280-4b10-a27b-d7f3e516b408","version":3}`

	TX_GETH_HTTP_PORT = "8544"
	TX_GETH_WS_PORT   = "8545"

	defaultGethPoaImage = "ethereum/client-go:v1.12.0"
)

type InternalDockerUrls struct {
	HttpUrl string
	WsUrl   string
}

type GethPoa struct {
	EnvComponent
	ExternalHttpUrl string
	InternalHttpUrl string
	ExternalWsUrl   string
	InternalWsUrl   string
	chainConfig     *EthereumChainConfig
	l               zerolog.Logger
	t               *testing.T
}

func NewGethPoa(networks []string, chainConfig *EthereumChainConfig, opts ...EnvComponentOption) *GethPoa {
	parts := strings.Split(defaultGethPoaImage, ":")
	g := &GethPoa{
		EnvComponent: EnvComponent{
			ContainerName:    fmt.Sprintf("%s-%s", "geth-pow", uuid.NewString()[0:8]),
			Networks:         networks,
			ContainerImage:   parts[0],
			ContainerVersion: parts[1],
		},
		chainConfig: chainConfig,
		l:           log.Logger,
	}
	g.SetDefaultHooks()
	for _, opt := range opts {
		opt(&g.EnvComponent)
	}
	// if the internal docker repo is set then add it to the version
	g.EnvComponent.ContainerImage = mirror.AddMirrorToImageIfSet(g.EnvComponent.ContainerImage)
	return g
}

func (g *GethPoa) WithLogger(l zerolog.Logger) *GethPoa {
	g.l = l
	return g
}

func (g *GethPoa) WithTestInstance(t *testing.T) ExecutionClient {
	g.l = logging.GetTestLogger(t)
	g.t = t
	return g
}

func (g *GethPoa) StartContainer() (blockchain.EVMNetwork, error) {
	r, err := g.getGethContainerRequest()
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
		return blockchain.EVMNetwork{}, fmt.Errorf("cannot start geth-poa container: %w", err)
	}
	host, err := GetHost(testcontext.Get(g.t), ct)
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
	networkConfig.Name = "Simulated Ethereum-PoA (geth)"
	networkConfig.URLs = []string{g.ExternalWsUrl}
	networkConfig.HTTPURLs = []string{g.ExternalHttpUrl}

	g.l.Info().Str("containerName", g.ContainerName).
		Str("internalHttpUrl", g.InternalHttpUrl).
		Str("externalHttpUrl", g.ExternalHttpUrl).
		Str("externalWsUrl", g.ExternalWsUrl).
		Str("internalWsUrl", g.InternalWsUrl).
		Msg("Started Geth container")

	return networkConfig, nil
}

func (g *GethPoa) GetContainer() *tc.Container {
	return &g.Container
}

func (g *GethPoa) GetContainerName() string {
	return g.ContainerName
}

func (g *GethPoa) GetContainerType() ContainerType {
	return ContainerType_Geth
}

func (g *GethPoa) GetInternalExecutionURL() string {
	panic("not supported")
}

func (g *GethPoa) GetExternalExecutionURL() string {
	panic("not supported")
}

func (g *GethPoa) GetInternalHttpUrl() string {
	return g.InternalHttpUrl
}

func (g *GethPoa) GetInternalWsUrl() string {
	return g.InternalWsUrl
}

func (g *GethPoa) GetExternalHttpUrl() string {
	return g.ExternalHttpUrl
}

func (g *GethPoa) GetExternalWsUrl() string {
	return g.ExternalWsUrl
}

func (g *GethPoa) getGethContainerRequest() (*tc.ContainerRequest, error) {
	initScriptFile, err := os.CreateTemp("", "init_script")
	if err != nil {
		return nil, err
	}
	_, err = initScriptFile.WriteString(templates.InitGethScript)
	if err != nil {
		return nil, err
	}
	keystoreDir, err := os.MkdirTemp("", "keystore")
	if err != nil {
		return nil, err
	}

	generatedData, err := generateKeystoreAndExtraData(keystoreDir)
	if err != nil {
		return nil, err
	}

	genesisJsonStr, err := templates.GenesisJsonTemplate{
		ChainId:     fmt.Sprintf("%d", g.chainConfig.ChainID),
		AccountAddr: []string{generatedData.minerAccount.Address.Hex(), RootFundingAddr},
		Consensus:   templates.GethGenesisConsensus_Clique,
		ExtraData:   fmt.Sprintf("0x%s", hex.EncodeToString(generatedData.extraData)),
	}.String()
	if err != nil {
		return nil, err
	}
	genesisFile, err := os.CreateTemp("", "genesis_json")
	if err != nil {
		return nil, err
	}
	_, err = genesisFile.WriteString(genesisJsonStr)
	if err != nil {
		return nil, err
	}
	key1File, err := os.CreateTemp(keystoreDir, "key1")
	if err != nil {
		return nil, err
	}
	_, err = key1File.WriteString(RootFundingWallet)
	if err != nil {
		return nil, err
	}
	configDir, err := os.MkdirTemp("", "config")
	if err != nil {
		return nil, err
	}
	err = os.WriteFile(configDir+"/password.txt", []byte(""), 0600)
	if err != nil {
		return nil, err
	}

	entrypoint, err := g.getEntryPointAndKeystoreLocation(generatedData.minerAccount.Address.Hex())
	if err != nil {
		return nil, err
	}

	websocketMsg, err := g.getWebsocketEnabledMessage()
	if err != nil {
		return nil, err
	}

	return &tc.ContainerRequest{
		Name:            g.ContainerName,
		AlwaysPullImage: true,
		Image:           g.GetImageWithVersion(),
		ExposedPorts:    []string{NatPortFormat(TX_GETH_HTTP_PORT), NatPortFormat(TX_GETH_WS_PORT)},
		Networks:        g.Networks,
		WaitingFor: tcwait.ForAll(
			NewHTTPStrategy("/", NatPort(TX_GETH_HTTP_PORT)),
			tcwait.ForLog(websocketMsg),
			tcwait.ForLog("Started P2P networking").
				WithStartupTimeout(120*time.Second).
				WithPollInterval(1*time.Second),
			NewWebSocketStrategy(NatPort(TX_GETH_WS_PORT), g.l),
		),
		Entrypoint: entrypoint,
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
		HostConfigModifier: func(hostConfig *container.HostConfig) {
			hostConfig.Mounts = append(hostConfig.Mounts, mount.Mount{
				Type:     mount.TypeBind,
				Source:   keystoreDir,
				Target:   "/root/.ethereum/devchain/keystore/",
				ReadOnly: false,
			}, mount.Mount{
				Type:     mount.TypeBind,
				Source:   configDir,
				Target:   "/root/config/",
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

func (g *GethPoa) WaitUntilChainIsReady(_ context.Context, _ time.Duration) error {
	return nil
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
		host, err = GetHost(ctx, target.(tc.Container))
		if err != nil {
			w.l.Error().Msg("Failed to get the target host")
			return err
		}
		wsPort, err := target.MappedPort(ctx, w.Port)
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

func (g *GethPoa) getEntryPointAndKeystoreLocation(minerAddress string) ([]string, error) {
	version, err := GetComparableVersionFromDockerImage(g.GetImageWithVersion())
	if err != nil {
		return nil, err
	}

	entrypoint := []string{
		"sh", "./root/init.sh",
		"--password",
		"/root/config/password.txt",
		"--ipcdisable",
		"--graphql",
		"-graphql.corsdomain",
		"*",
		"--allow-insecure-unlock",
		"--vmdebug",
		fmt.Sprintf("--networkid=%d", g.chainConfig.ChainID),
		"--datadir",
		"/root/.ethereum/devchain",
		"--mine",
		"--miner.etherbase",
		minerAddress,
		"--unlock",
		minerAddress,
	}

	if version < 110 {
		entrypoint = append(entrypoint,
			"--rpc",
			"--rpcapi",
			"admin,debug,web3,eth,txpool,personal,clique,miner,net",
			"--rpccorsdomain",
			"*",
			"--rpcvhosts",
			"*",
			"--rpcaddr",
			"0.0.0.0",
			fmt.Sprintf("--rpcport=%s", TX_GETH_HTTP_PORT),
			"--ws",
			"--wsorigins",
			"*",
			"--wsaddr",
			"0.0.0.0",
			"--wsapi",
			"admin,debug,web3,eth,txpool,personal,clique,miner,net",
			fmt.Sprintf("--wsport=%s", TX_GETH_WS_PORT),
		)
	}

	if version >= 110 {
		entrypoint = append(entrypoint,
			"--http",
			"--http.vhosts",
			"*",
			"--http.api",
			"admin,debug,web3,eth,txpool,personal,clique,miner,net",
			"--http.corsdomain",
			"*",
			"--http.addr",
			"0.0.0.0",
			fmt.Sprintf("--http.port=%s", TX_GETH_HTTP_PORT),
			"--ws",
			"--ws.origins",
			"*",
			"--ws.addr",
			"0.0.0.0",
			"--ws.api",
			"admin,debug,web3,eth,txpool,personal,clique,miner,net",
			fmt.Sprintf("--ws.port=%s", TX_GETH_WS_PORT),
			"--rpc.allow-unprotected-txs",
			"--rpc.txfeecap",
			"0",
		)
	}

	return entrypoint, nil
}

func (g *GethPoa) getWebsocketEnabledMessage() (string, error) {
	version, err := GetComparableVersionFromDockerImage(g.GetImageWithVersion())
	if err != nil {
		return "", err
	}

	if version < 110 {
		return "WebSocket endpoint opened", nil
	}

	return "WebSocket enabled", nil
}

type HTTPStrategy struct {
	Path               string
	Port               nat.Port
	RetryDelay         time.Duration
	ExpectedStatusCode int
	timeout            time.Duration
}

func NewHTTPStrategy(path string, port nat.Port) *HTTPStrategy {
	return &HTTPStrategy{
		Path:               path,
		Port:               port,
		RetryDelay:         10 * time.Second,
		ExpectedStatusCode: 200,
		timeout:            2 * time.Minute,
	}
}

func (w *HTTPStrategy) WithTimeout(timeout time.Duration) *HTTPStrategy {
	w.timeout = timeout
	return w
}

func (w *HTTPStrategy) WithStatusCode(statusCode int) *HTTPStrategy {
	w.ExpectedStatusCode = statusCode
	return w
}

// WaitUntilReady implements Strategy.WaitUntilReady
func (w *HTTPStrategy) WaitUntilReady(ctx context.Context, target tcwait.StrategyTarget) (err error) {

	ctx, cancel := context.WithTimeout(ctx, w.timeout)
	defer cancel()

	host, err := GetHost(ctx, target.(tc.Container))
	if err != nil {
		return
	}

	var mappedPort nat.Port
	mappedPort, err = target.MappedPort(ctx, w.Port)
	if err != nil {
		return err
	}

	tripper := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	client := http.Client{Transport: tripper, Timeout: time.Second}
	address := net.JoinHostPort(host, strconv.Itoa(mappedPort.Int()))

	endpoint := url.URL{
		Scheme: "http",
		Host:   address,
		Path:   w.Path,
	}

	var body []byte
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(100 * time.Millisecond):
			state, err := target.State(ctx)
			if err != nil {
				return err
			}
			if !state.Running {
				return fmt.Errorf("container is not running %s", state.Status)
			}
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), bytes.NewReader(body))
			if err != nil {
				return err
			}
			resp, err := client.Do(req)
			if err != nil {
				continue
			}
			if resp.StatusCode != w.ExpectedStatusCode {
				_ = resp.Body.Close()
				continue
			}
			if err := resp.Body.Close(); err != nil {
				continue
			}
			return nil
		}
	}
}
