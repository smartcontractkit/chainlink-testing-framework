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
	"github.com/smartcontractkit/chainlink-testing-framework/utils/testcontext"
)

const defaultGethPosImage = "ethereum/client-go:v1.13.10"

type GethPos struct {
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

func NewGethPos(networks []string, chainConfg *EthereumChainConfig, generatedDataHostDir string, consensusLayer ConsensusLayer, opts ...EnvComponentOption) (*GethPos, error) {
	parts := strings.Split(defaultGethPosImage, ":")
	g := &GethPos{
		EnvComponent: EnvComponent{
			ContainerName:    fmt.Sprintf("%s-%s", "geth-pos", uuid.NewString()[0:8]),
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

func (g *GethPos) WithTestInstance(t *testing.T) ExecutionClient {
	g.l = logging.GetTestLogger(t)
	g.t = t
	return g
}

func (g *GethPos) StartContainer() (blockchain.EVMNetwork, error) {
	r, err := g.getContainerRequest(g.Networks)
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
	httpPort, err := ct.MappedPort(testcontext.Get(g.t), NatPort(TX_GETH_HTTP_PORT))
	if err != nil {
		return blockchain.EVMNetwork{}, err
	}
	wsPort, err := ct.MappedPort(testcontext.Get(g.t), NatPort(TX_GETH_WS_PORT))
	if err != nil {
		return blockchain.EVMNetwork{}, err
	}
	executionPort, err := ct.MappedPort(testcontext.Get(g.t), NatPort(ETH2_EXECUTION_PORT))
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
	networkConfig.Name = fmt.Sprintf("Simulated Ethereum-PoS (geth + %s)", g.consensusLayer)
	networkConfig.URLs = []string{g.ExternalWsUrl}
	networkConfig.HTTPURLs = []string{g.ExternalHttpUrl}

	g.l.Info().Str("containerName", g.ContainerName).
		Msg("Started Geth PoS container")

	return networkConfig, nil
}

func (g *GethPos) GetInternalExecutionURL() string {
	return g.InternalExecutionURL
}

func (g *GethPos) GetExternalExecutionURL() string {
	return g.ExternalExecutionURL
}

func (g *GethPos) GetInternalHttpUrl() string {
	return g.InternalHttpUrl
}

func (g *GethPos) GetInternalWsUrl() string {
	return g.InternalWsUrl
}

func (g *GethPos) GetExternalHttpUrl() string {
	return g.ExternalHttpUrl
}

func (g *GethPos) GetExternalWsUrl() string {
	return g.ExternalWsUrl
}

func (g *GethPos) GetContainerName() string {
	return g.ContainerName
}

func (g *GethPos) GetContainer() *tc.Container {
	return &g.Container
}

func (g *GethPos) getContainerRequest(networks []string) (*tc.ContainerRequest, error) {
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
		ExposedPorts:  []string{NatPortFormat(TX_GETH_HTTP_PORT), NatPortFormat(TX_GETH_WS_PORT), NatPortFormat(ETH2_EXECUTION_PORT)},
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

func (g *GethPos) WaitUntilChainIsReady(ctx context.Context, waitTime time.Duration) error {
	waitForFirstBlock := tcwait.NewLogStrategy("Chain head was updated").WithPollInterval(1 * time.Second).WithStartupTimeout(waitTime)
	return waitForFirstBlock.WaitUntilReady(ctx, *g.GetContainer())
}

func (g *GethPos) buildInitScript() (string, error) {
	initTemplate := `#!/bin/bash
	mkdir -p {{.ExecutionDir}}

	# copy general keystore to execution directory, because Geth doesn't allow to specify keystore location
	echo "Copying keystore to {{.ExecutionDir}}/keystore"
	cp -R {{.KeystoreDirLocation}} {{.ExecutionDir}}/keystore

	echo "Creating sk.json file"
	echo "2e0834786285daccd064ca17f1654f67b4aef298acbb82cef9ec422fb4975622" > {{.ExecutionDir}}/sk.json

	echo "Running geth init"
	geth init --state.scheme=path --datadir={{.ExecutionDir}} {{.GeneratedDataDir}}/genesis.json
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
		HttpPort:             TX_GETH_HTTP_PORT,
		WsPort:               TX_GETH_WS_PORT,
		ChainID:              g.chainConfg.ChainID,
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

func (g *GethPos) GetContainerType() ContainerType {
	return ContainerType_Geth
}
