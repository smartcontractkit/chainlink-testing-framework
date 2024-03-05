package test_env

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"html/template"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/google/uuid"
	"github.com/rs/zerolog"

	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/docker"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/mirror"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/templates"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/testcontext"
)

const defaultGethPosImage = "ethereum/client-go:v1.13.10"
const defaultGethPoaImage = "ethereum/client-go:v1.12.0"

type Geth struct {
	EnvComponent
	ExternalHttpUrl      string
	InternalHttpUrl      string
	ExternalWsUrl        string
	InternalWsUrl        string
	InternalExecutionURL string
	ExternalExecutionURL string
	generatedDataHostDir string
	chainConfig          *EthereumChainConfig
	consensusLayer       ConsensusLayer
	l                    zerolog.Logger
	t                    *testing.T
}

func NewGethPos(networks []string, chainConfg *EthereumChainConfig, generatedDataHostDir string, consensusLayer ConsensusLayer, opts ...EnvComponentOption) (*Geth, error) {
	parts := strings.Split(defaultGethPosImage, ":")
	g := &Geth{
		EnvComponent: EnvComponent{
			ContainerName:    fmt.Sprintf("%s-%s", "geth-pos", uuid.NewString()[0:8]),
			Networks:         networks,
			ContainerImage:   parts[0],
			ContainerVersion: parts[1],
		},
		chainConfig:          chainConfg,
		generatedDataHostDir: generatedDataHostDir,
		consensusLayer:       consensusLayer,
		l:                    logging.GetTestLogger(nil),
	}
	g.SetDefaultHooks()
	for _, opt := range opts {
		opt(&g.EnvComponent)
	}
	// if the internal docker repo is set then add it to the version
	g.EnvComponent.ContainerImage = mirror.AddMirrorToImageIfSet(g.EnvComponent.ContainerImage)
	return g, nil
}

func NewGethPoa(networks []string, chainConfig *EthereumChainConfig, opts ...EnvComponentOption) *Geth {
	parts := strings.Split(defaultGethPoaImage, ":")
	g := &Geth{
		EnvComponent: EnvComponent{
			ContainerName:    fmt.Sprintf("%s-%s", "geth-poa", uuid.NewString()[0:8]),
			Networks:         networks,
			ContainerImage:   parts[0],
			ContainerVersion: parts[1],
		},
		chainConfig: chainConfig,
		l:           logging.GetTestLogger(nil),
	}
	g.SetDefaultHooks()
	for _, opt := range opts {
		opt(&g.EnvComponent)
	}
	// if the internal docker repo is set then add it to the version
	g.EnvComponent.ContainerImage = mirror.AddMirrorToImageIfSet(g.EnvComponent.ContainerImage)
	return g
}

func (g *Geth) WithTestInstance(t *testing.T) ExecutionClient {
	g.l = logging.GetTestLogger(t)
	g.t = t
	return g
}

func (g *Geth) StartContainer() (blockchain.EVMNetwork, error) {
	var r *tc.ContainerRequest
	var err error
	if g.GetEthereumVersion() == EthereumVersion_Eth1 {
		r, err = g.getPowGethContainerRequest()
	} else {
		r, err = g.getPosContainerRequest(g.Networks)
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
		return blockchain.EVMNetwork{}, fmt.Errorf("cannot start geth container: %w", err)
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
	if g.GetEthereumVersion() == EthereumVersion_Eth1 {
		networkConfig.Name = "Simulated Eth-1-PoA (geth)"
	} else {
		networkConfig.Name = fmt.Sprintf("Simulated Eth-2-PoS (geth + %s)", g.consensusLayer)
	}
	networkConfig.URLs = []string{g.ExternalWsUrl}
	networkConfig.HTTPURLs = []string{g.ExternalHttpUrl}

	g.l.Info().Str("containerName", g.ContainerName).
		Msg("Started Geth PoS container")

	return networkConfig, nil
}

func (g *Geth) GetInternalExecutionURL() string {
	if g.GetEthereumVersion() == EthereumVersion_Eth1 {
		panic("eth1 node doesn't have an execution URL")
	}
	return g.InternalExecutionURL
}

func (g *Geth) GetExternalExecutionURL() string {
	if g.GetEthereumVersion() == EthereumVersion_Eth1 {
		panic("eth1 node doesn't have an execution URL")
	}
	return g.ExternalExecutionURL
}

func (g *Geth) GetInternalHttpUrl() string {
	return g.InternalHttpUrl
}

func (g *Geth) GetInternalWsUrl() string {
	return g.InternalWsUrl
}

func (g *Geth) GetExternalHttpUrl() string {
	return g.ExternalHttpUrl
}

func (g *Geth) GetExternalWsUrl() string {
	return g.ExternalWsUrl
}

func (g *Geth) GetContainerName() string {
	return g.ContainerName
}

func (g *Geth) GetContainer() *tc.Container {
	return &g.Container
}

func (g *Geth) GetEthereumVersion() EthereumVersion {
	if g.consensusLayer != "" {
		return EthereumVersion_Eth2
	}

	return EthereumVersion_Eth1
}

func (g *Geth) WaitUntilChainIsReady(ctx context.Context, waitTime time.Duration) error {
	if g.GetEthereumVersion() == EthereumVersion_Eth1 {
		return nil
	}
	waitForFirstBlock := tcwait.NewLogStrategy("Chain head was updated").WithPollInterval(1 * time.Second).WithStartupTimeout(waitTime)
	return waitForFirstBlock.WaitUntilReady(ctx, *g.GetContainer())
}

func (g *Geth) getPosContainerRequest(networks []string) (*tc.ContainerRequest, error) {
	initFile, err := os.CreateTemp("", "init.sh")
	if err != nil {
		return nil, err
	}

	initScriptContent, err := g.buildInitScript()
	if err != nil {
		return nil, err
	}

	_, err = initFile.WriteString(initScriptContent)
	if err != nil {
		return nil, err
	}

	return &tc.ContainerRequest{
		Name:          g.ContainerName,
		Image:         g.GetImageWithVersion(),
		Networks:      networks,
		ImagePlatform: "linux/x86_64",
		ExposedPorts:  []string{NatPortFormat(DEFAULT_EVM_NODE_HTTP_PORT), NatPortFormat(DEFAULT_EVM_NODE_WS_PORT), NatPortFormat(ETH2_EXECUTION_PORT)},
		WaitingFor: tcwait.ForAll(
			tcwait.ForLog("WebSocket enabled").
				WithStartupTimeout(120 * time.Second).
				WithPollInterval(1 * time.Second),
		),
		Entrypoint: []string{
			"sh",
			"/init.sh",
		},
		Files: []tc.ContainerFile{
			{
				HostFilePath:      initFile.Name(),
				ContainerFilePath: "/init.sh",
				FileMode:          0744,
			},
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

func (g *Geth) getPowGethContainerRequest() (*tc.ContainerRequest, error) {
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
		ExposedPorts:    []string{NatPortFormat(DEFAULT_EVM_NODE_HTTP_PORT), NatPortFormat(DEFAULT_EVM_NODE_WS_PORT)},
		Networks:        g.Networks,
		WaitingFor: tcwait.ForAll(
			NewHTTPStrategy("/", NatPort(DEFAULT_EVM_NODE_HTTP_PORT)),
			tcwait.ForLog(websocketMsg),
			tcwait.ForLog("Started P2P networking").
				WithStartupTimeout(120*time.Second).
				WithPollInterval(1*time.Second),
			NewWebSocketStrategy(NatPort(DEFAULT_EVM_NODE_WS_PORT), g.l),
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

func (g *Geth) getEntryPointAndKeystoreLocation(minerAddress string) ([]string, error) {
	version, err := GetComparableVersionFromDockerImage(g.GetImageWithVersion())
	if err != nil {
		return nil, err
	}

	var enabledApis = "admin,debug,web3,eth,txpool,personal,clique,miner,net"

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
			enabledApis,
			"--rpccorsdomain",
			"*",
			"--rpcvhosts",
			"*",
			"--rpcaddr",
			"0.0.0.0",
			fmt.Sprintf("--rpcport=%s", DEFAULT_EVM_NODE_HTTP_PORT),
			"--ws",
			"--wsorigins",
			"*",
			"--wsaddr",
			"0.0.0.0",
			"--wsapi",
			enabledApis,
			fmt.Sprintf("--wsport=%s", DEFAULT_EVM_NODE_WS_PORT),
		)
	}

	if version >= 110 {
		entrypoint = append(entrypoint,
			"--http",
			"--http.vhosts",
			"*",
			"--http.api",
			enabledApis,
			"--http.corsdomain",
			"*",
			"--http.addr",
			"0.0.0.0",
			fmt.Sprintf("--http.port=%s", DEFAULT_EVM_NODE_HTTP_PORT),
			"--ws",
			"--ws.origins",
			"*",
			"--ws.addr",
			"0.0.0.0",
			"--ws.api",
			enabledApis,
			fmt.Sprintf("--ws.port=%s", DEFAULT_EVM_NODE_WS_PORT),
			"--rpc.allow-unprotected-txs",
			"--rpc.txfeecap",
			"0",
		)
	}

	return entrypoint, nil
}

func (g *Geth) getWebsocketEnabledMessage() (string, error) {
	version, err := GetComparableVersionFromDockerImage(g.GetImageWithVersion())
	if err != nil {
		return "", err
	}

	if version < 110 {
		return "WebSocket endpoint opened", nil
	}

	return "WebSocket enabled", nil
}

func (g *Geth) buildInitScript() (string, error) {
	initTemplate := `#!/bin/bash
	mkdir -p {{.ExecutionDir}}

	# copy general keystore to execution directory, because Geth doesn't allow to specify keystore location
	echo "Copying keystore to {{.ExecutionDir}}/keystore"
	cp -R {{.KeystoreDirLocation}} {{.ExecutionDir}}/keystore

	echo "Creating sk.json file"
	echo "2e0834786285daccd064ca17f1654f67b4aef298acbb82cef9ec422fb4975622" > {{.ExecutionDir}}/sk.json

	echo "Running geth init"
	# geth init --state.scheme=path --datadir={{.ExecutionDir}} {{.GeneratedDataDir}}/genesis.json
	geth init --datadir={{.ExecutionDir}} {{.GeneratedDataDir}}/genesis.json
	exit_code=$?
	if [ $exit_code -ne 0 ]; then
		echo "Geth init failed with exit code $exit_code"
		exit 1
	fi

	echo "Starting Geth..."
	geth --http --http.api=eth,net,web3,debug --http.addr=0.0.0.0 --http.corsdomain=* \
		--http.vhosts=* --http.port={{.HttpPort}} --ws --ws.api=admin,debug,web3,eth,txpool,net \
		--ws.addr=0.0.0.0 --ws.origins=* --ws.port={{.WsPort}} --authrpc.vhosts=* \
		--authrpc.addr=0.0.0.0 --authrpc.jwtsecret={{.JwtFileLocation}} --datadir={{.ExecutionDir}} \
		--rpc.allow-unprotected-txs --rpc.txfeecap=0 --allow-insecure-unlock \
		--password={{.PasswordFileLocation}} --nodiscover --syncmode=full --networkid={{.ChainID}} \
		--graphql --graphql.corsdomain=* --unlock=0x123463a4b065722e99115d6c222f267d9cabb524`

	data := struct {
		HttpPort             string
		WsPort               string
		ChainID              int
		GeneratedDataDir     string
		JwtFileLocation      string
		PasswordFileLocation string
		KeystoreDirLocation  string
		ExecutionDir         string
	}{
		HttpPort:             DEFAULT_EVM_NODE_HTTP_PORT,
		WsPort:               DEFAULT_EVM_NODE_WS_PORT,
		ChainID:              g.chainConfig.ChainID,
		GeneratedDataDir:     GENERATED_DATA_DIR_INSIDE_CONTAINER,
		JwtFileLocation:      JWT_SECRET_FILE_LOCATION_INSIDE_CONTAINER,
		PasswordFileLocation: ACCOUNT_PASSWORD_FILE_INSIDE_CONTAINER,
		KeystoreDirLocation:  KEYSTORE_DIR_LOCATION_INSIDE_CONTAINER,
		ExecutionDir:         "/execution-data",
	}

	t, err := template.New("init").Parse(initTemplate)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		os.Exit(1)
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, data)

	return buf.String(), err

}

func (g *Geth) GetContainerType() ContainerType {
	return ContainerType_Geth
}
