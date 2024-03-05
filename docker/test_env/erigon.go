package test_env

import (
	"bytes"
	"context"
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

const defaultErigonPosImage = "thorax/erigon:v2.56.2"
const defaultErigonPowImage = "thorax/erigon:v2.40.0"

type Erigon struct {
	EnvComponent
	ExternalHttpUrl      string
	InternalHttpUrl      string
	ExternalWsUrl        string
	InternalWsUrl        string
	InternalExecutionURL string
	ExternalExecutionURL string
	generatedDataHostDir string
	chainConfg           *EthereumChainConfig
	consensusLayer       ConsensusLayer
	l                    zerolog.Logger
	t                    *testing.T
}

func NewErigonPos(networks []string, chainConfg *EthereumChainConfig, generatedDataHostDir string, consensusLayer ConsensusLayer, opts ...EnvComponentOption) (*Erigon, error) {
	parts := strings.Split(defaultErigonPosImage, ":")
	g := &Erigon{
		EnvComponent: EnvComponent{
			ContainerName:    fmt.Sprintf("%s-%s", "erigon-pos", uuid.NewString()[0:8]),
			Networks:         networks,
			ContainerImage:   parts[0],
			ContainerVersion: parts[1],
		},
		chainConfg:           chainConfg,
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

func NewErigonPow(networks []string, chainConfg *EthereumChainConfig, opts ...EnvComponentOption) (*Erigon, error) {
	parts := strings.Split(defaultErigonPowImage, ":")
	g := &Erigon{
		EnvComponent: EnvComponent{
			ContainerName:    fmt.Sprintf("%s-%s", "erigon-pow", uuid.NewString()[0:8]),
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

func (g *Erigon) WithTestInstance(t *testing.T) ExecutionClient {
	g.l = logging.GetTestLogger(t)
	g.t = t
	return g
}

func (g *Erigon) StartContainer() (blockchain.EVMNetwork, error) {
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
		return blockchain.EVMNetwork{}, fmt.Errorf("cannot start erigon container: %w", err)
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
	g.ExternalWsUrl = FormatWsUrl(host, httpPort.Port())
	g.InternalWsUrl = FormatWsUrl(g.ContainerName, DEFAULT_EVM_NODE_HTTP_PORT)

	networkConfig := blockchain.SimulatedEVMNetwork
	if g.GetEthereumVersion() == EthereumVersion_Eth1 {
		networkConfig.Name = "Simulated Eth-1-PoW (erigon)"
	} else {
		networkConfig.Name = fmt.Sprintf("Simulated Eth-2-PoS (erigon + %s)", g.consensusLayer)
	}
	networkConfig.URLs = []string{g.ExternalWsUrl}
	networkConfig.HTTPURLs = []string{g.ExternalHttpUrl}

	g.l.Info().Str("containerName", g.ContainerName).
		Msg("Started Erigon container")

	return networkConfig, nil
}

func (g *Erigon) GetInternalExecutionURL() string {
	if g.GetEthereumVersion() == EthereumVersion_Eth1 {
		panic("eth1 node doesn't have an execution URL")
	}
	return g.InternalExecutionURL
}

func (g *Erigon) GetExternalExecutionURL() string {
	if g.GetEthereumVersion() == EthereumVersion_Eth1 {
		panic("eth1 node doesn't have an execution URL")
	}
	return g.ExternalExecutionURL
}

func (g *Erigon) GetInternalHttpUrl() string {
	return g.InternalHttpUrl
}

func (g *Erigon) GetInternalWsUrl() string {
	return g.InternalWsUrl
}

func (g *Erigon) GetExternalHttpUrl() string {
	return g.ExternalHttpUrl
}

func (g *Erigon) GetExternalWsUrl() string {
	return g.ExternalWsUrl
}

func (g *Erigon) GetContainerName() string {
	return g.ContainerName
}

func (g *Erigon) GetContainer() *tc.Container {
	return &g.Container
}

func (g *Erigon) GetEthereumVersion() EthereumVersion {
	if g.consensusLayer != "" {
		return EthereumVersion_Eth2
	}

	return EthereumVersion_Eth1
}

func (g *Erigon) WaitUntilChainIsReady(ctx context.Context, waitTime time.Duration) error {
	if g.GetEthereumVersion() == EthereumVersion_Eth1 {
		return nil
	}
	waitForFirstBlock := tcwait.NewLogStrategy("Built block").WithPollInterval(1 * time.Second).WithStartupTimeout(waitTime)
	return waitForFirstBlock.WaitUntilReady(ctx, *g.GetContainer())
}

func (g *Erigon) GetContainerType() ContainerType {
	return ContainerType_Erigon
}

func (g *Erigon) getPosContainerRequest() (*tc.ContainerRequest, error) {
	initFile, err := os.CreateTemp("", "init.sh")
	if err != nil {
		return nil, err
	}

	initScriptContent, err := g.buildPosInitScript()
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
		Networks:      g.Networks,
		ImagePlatform: "linux/x86_64",
		ExposedPorts:  []string{NatPortFormat(DEFAULT_EVM_NODE_HTTP_PORT), NatPortFormat(ETH2_EXECUTION_PORT)},
		WaitingFor: tcwait.ForAll(
			tcwait.ForLog("Started P2P networking").
				WithStartupTimeout(120 * time.Second).
				WithPollInterval(1 * time.Second),
		),
		User: "0:0",
		Entrypoint: []string{
			"sh",
			"/home/erigon/init.sh",
		},
		Files: []tc.ContainerFile{
			{
				HostFilePath:      initFile.Name(),
				ContainerFilePath: "/home/erigon/init.sh",
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

func (g *Erigon) getPowContainerRequest() (*tc.ContainerRequest, error) {
	keystoreDir, err := os.MkdirTemp("", "keystore")
	if err != nil {
		return nil, err
	}

	generatedData, err := generateKeystoreAndExtraData(keystoreDir)
	if err != nil {
		return nil, err
	}

	genesisJsonStr, err := templates.GenesisJsonTemplate{
		ChainId:     fmt.Sprintf("%d", g.chainConfg.ChainID),
		AccountAddr: []string{RootFundingAddr, generatedData.minerAccount.Address.Hex()},
		Consensus:   templates.GethGenesisConsensus_Ethash,
	}.String()
	if err != nil {
		return nil, err
	}

	initFile, err := os.CreateTemp("", "init.sh")
	if err != nil {
		return nil, err
	}

	initScriptContent, err := g.buildPowInitScript(generatedData.minerAccount.Address.Hex())
	if err != nil {
		return nil, err
	}

	_, err = initFile.WriteString(initScriptContent)
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

	return &tc.ContainerRequest{
		Name:          g.ContainerName,
		Image:         g.GetImageWithVersion(),
		Networks:      g.Networks,
		ImagePlatform: "linux/x86_64",
		ExposedPorts:  []string{NatPortFormat(DEFAULT_EVM_NODE_HTTP_PORT)},
		WaitingFor: tcwait.ForAll(
			tcwait.ForLog("Started P2P networking").
				WithStartupTimeout(120*time.Second).
				WithPollInterval(1*time.Second),
			NewWebSocketStrategy(NatPort(DEFAULT_EVM_NODE_HTTP_PORT), g.l),
		),
		User: "0:0",
		Entrypoint: []string{
			"sh",
			"/home/erigon/init.sh",
		},
		Files: []tc.ContainerFile{
			{
				HostFilePath:      initFile.Name(),
				ContainerFilePath: "/root/init.sh",
				FileMode:          0644,
			},
			{
				HostFilePath:      genesisFile.Name(),
				ContainerFilePath: "/root/genesis.json",
				FileMode:          0644,
			},
			{
				HostFilePath:      initFile.Name(),
				ContainerFilePath: "/home/erigon/init.sh",
				FileMode:          0744,
			},
		},
		HostConfigModifier: func(hostConfig *container.HostConfig) {
			hostConfig.Mounts = append(hostConfig.Mounts, mount.Mount{
				Type:     mount.TypeBind,
				Source:   keystoreDir,
				Target:   "/root/.local/share/erigon/keystore/",
				ReadOnly: false,
			},
			)
		},
		LifecycleHooks: []tc.ContainerLifecycleHooks{
			{
				PostStarts: g.PostStartsHooks,
				PostStops:  g.PostStopsHooks,
			},
		},
	}, nil
}

func (g *Erigon) getExtraExecutionFlags() (string, error) {
	version, err := GetComparableVersionFromDockerImage(g.GetImageWithVersion())
	if err != nil {
		return "", err
	}

	extraExecutionFlags := ""
	if version > 247 {
		extraExecutionFlags = " --rpc.txfeecap=0"
	}

	if version > 254 {
		extraExecutionFlags += " --rpc.allow-unprotected-txs"
	}

	if version > 255 {
		extraExecutionFlags += " --db.size.limit=8TB"
	}

	return extraExecutionFlags, nil
}

func (g *Erigon) buildPosInitScript() (string, error) {
	extraExecutionFlags, err := g.getExtraExecutionFlags()
	if err != nil {
		return "", err
	}

	initTemplate := `#!/bin/bash
	echo "Copied genesis file to {{.ExecutionDir}}"
	mkdir -p {{.ExecutionDir}}
	cp {{.GeneratedDataDir}}/genesis.json {{.ExecutionDir}}/genesis.json
	echo "Running erigon init"
	erigon init --datadir={{.ExecutionDir}} {{.ExecutionDir}}/genesis.json
	exit_code=$?
	if [ $exit_code -ne 0 ]; then
		echo "Erigon init failed with exit code $exit_code"
		exit 1
	fi

	echo "Starting Erigon..."
	erigon --http --http.api=eth,erigon,engine,web3,net,debug,trace,txpool,admin --http.addr=0.0.0.0 --http.corsdomain=* \
		--http.vhosts=* --http.port={{.HttpPort}} --ws --authrpc.vhosts=* --authrpc.addr=0.0.0.0 --authrpc.jwtsecret={{.JwtFileLocation}} \
		--datadir={{.ExecutionDir}} {{.ExtraExecutionFlags}} --allow-insecure-unlock --nodiscover --networkid={{.ChainID}}`

	data := struct {
		HttpPort            string
		ChainID             int
		GeneratedDataDir    string
		JwtFileLocation     string
		ExecutionDir        string
		ExtraExecutionFlags string
	}{
		HttpPort:            DEFAULT_EVM_NODE_HTTP_PORT,
		ChainID:             g.chainConfg.ChainID,
		GeneratedDataDir:    GENERATED_DATA_DIR_INSIDE_CONTAINER,
		JwtFileLocation:     JWT_SECRET_FILE_LOCATION_INSIDE_CONTAINER,
		ExecutionDir:        "/home/erigon/execution-data",
		ExtraExecutionFlags: extraExecutionFlags,
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

func (g *Erigon) buildPowInitScript(minerAddr string) (string, error) {
	initTemplate := `#!/bin/bash
	echo "Running erigon init"
	erigon init /root/genesis.json
	exit_code=$?
	if [ $exit_code -ne 0 ]; then
		echo "Erigon init failed with exit code $exit_code"
		exit 1
	fi

	echo "Starting Erigon..."
	erigon --http --http.api=eth,erigon,engine,web3,net,debug,trace,txpool,admin --http.addr=0.0.0.0 --http.corsdomain=* --http.vhosts=* --http.port={{.HttpPort}} --ws \
	--allow-insecure-unlock  --nodiscover --networkid={{.ChainID}} --mine --miner.etherbase={{.MinerAddr}} --fakepow`

	data := struct {
		HttpPort  string
		ChainID   int
		MinerAddr string
	}{
		HttpPort:  DEFAULT_EVM_NODE_HTTP_PORT,
		ChainID:   g.chainConfg.ChainID,
		MinerAddr: minerAddr,
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
