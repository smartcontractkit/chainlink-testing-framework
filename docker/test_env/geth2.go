package test_env

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"os"
	"testing"
	"time"

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

type Geth2 struct {
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
	image                string
}

func NewGeth2(networks []string, chainConfg *EthereumChainConfig, generatedDataHostDir string, consensusLayer ConsensusLayer, opts ...EnvComponentOption) (*Geth2, error) {
	// currently it uses v1.13.5
	dockerImage, err := mirror.GetImage("ethereum/client-go:v1.13")
	if err != nil {
		return nil, err
	}

	g := &Geth2{
		EnvComponent: EnvComponent{
			ContainerName: fmt.Sprintf("%s-%s", "geth2", uuid.NewString()[0:8]),
			Networks:      networks,
		},
		chainConfg:           chainConfg,
		generatedDataHostDir: generatedDataHostDir,
		consensusLayer:       consensusLayer,
		l:                    logging.GetTestLogger(nil),
		image:                dockerImage,
	}
	g.SetDefaultHooks()
	for _, opt := range opts {
		opt(&g.EnvComponent)
	}
	return g, nil
}

func (g *Geth2) WithImage(imageWithTag string) *Geth2 {
	g.image = imageWithTag
	return g
}

func (g *Geth2) WithTestInstance(t *testing.T) ExecutionClient {
	g.l = logging.GetTestLogger(t)
	g.t = t
	return g
}

func (g *Geth2) GetImage() string {
	return g.image
}

func (g *Geth2) StartContainer() (blockchain.EVMNetwork, error) {
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
	networkConfig.Name = fmt.Sprintf("Simulated Eth2 (geth + %s)", g.consensusLayer)
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
		Image:         g.image,
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
		Mounts: tc.ContainerMounts{
			tc.ContainerMount{
				Source: tc.GenericBindMountSource{
					HostPath: g.generatedDataHostDir,
				},
				Target: tc.ContainerMountTarget(GENERATED_DATA_DIR_INSIDE_CONTAINER),
			},
		},
		LifecycleHooks: []tc.ContainerLifecycleHooks{
			{
				PostStarts: g.PostStartsHooks,
				PostStops:  g.PostStopsHooks,
			},
		},
	}, nil
}

func (g *Geth2) WaitUntilChainIsReady(ctx context.Context, waitTime time.Duration) error {
	waitForFirstBlock := tcwait.NewLogStrategy("Chain head was updated").WithPollInterval(1 * time.Second).WithStartupTimeout(waitTime)
	return waitForFirstBlock.WaitUntilReady(ctx, *g.GetContainer())
}

func (g *Geth2) buildInitScript() (string, error) {
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

func (g *Geth2) GetContainerType() ContainerType {
	return ContainerType_Geth2
}
