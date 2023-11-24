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
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/docker"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/testcontext"
)

const (
	//TODO use Tate's mirror?
	ERIGON_IMAGE_TAG = "v2.54.0"
)

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

func NewErigon(networks []string, chainConfg *EthereumChainConfig, generatedDataHostDir string, consensusLayer ConsensusLayer, opts ...EnvComponentOption) *Erigon {
	g := &Erigon{
		EnvComponent: EnvComponent{
			ContainerName: fmt.Sprintf("%s-%s", "erigon", uuid.NewString()[0:8]),
			Networks:      networks,
		},
		chainConfg:           chainConfg,
		generatedDataHostDir: generatedDataHostDir,
		consensusLayer:       consensusLayer,
		l:                    log.Logger,
	}
	for _, opt := range opts {
		opt(&g.EnvComponent)
	}
	return g
}

func (g *Erigon) WithTestInstance(t *testing.T) *Erigon {
	g.l = logging.GetTestLogger(t)
	g.t = t
	return g
}

func (g *Erigon) StartContainer() (blockchain.EVMNetwork, error) {
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
		return blockchain.EVMNetwork{}, errors.Wrapf(err, "cannot start erigon container")
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
	executionPort, err := ct.MappedPort(testcontext.Get(g.t), NatPort(ETH2_EXECUTION_PORT))
	if err != nil {
		return blockchain.EVMNetwork{}, err
	}

	g.Container = ct
	g.ExternalHttpUrl = FormatHttpUrl(host, httpPort.Port())
	g.InternalHttpUrl = FormatHttpUrl(g.ContainerName, TX_GETH_HTTP_PORT)
	g.ExternalWsUrl = FormatWsUrl(host, httpPort.Port())
	g.InternalWsUrl = FormatWsUrl(g.ContainerName, TX_GETH_HTTP_PORT)
	g.InternalExecutionURL = FormatHttpUrl(g.ContainerName, ETH2_EXECUTION_PORT)
	g.ExternalExecutionURL = FormatHttpUrl(host, executionPort.Port())

	networkConfig := blockchain.SimulatedEVMNetwork
	networkConfig.Name = fmt.Sprintf("Simulated Eth2 (Erigon %s)", g.consensusLayer)
	networkConfig.URLs = []string{g.ExternalWsUrl}
	networkConfig.HTTPURLs = []string{g.ExternalHttpUrl}

	g.l.Info().Str("containerName", g.ContainerName).
		Msg("Started Erigon container")

	return networkConfig, nil
}

func (g *Erigon) GetInternalExecutionURL() string {
	return g.InternalExecutionURL
}

func (g *Erigon) GetExternalExecutionURL() string {
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

func (g *Erigon) getContainerRequest(networks []string) (*tc.ContainerRequest, error) {
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
		Image:         fmt.Sprintf("thorax/erigon:%s", ERIGON_IMAGE_TAG),
		Networks:      networks,
		ImagePlatform: "linux/x86_64",
		ExposedPorts:  []string{NatPortFormat(TX_GETH_HTTP_PORT), NatPortFormat(ETH2_EXECUTION_PORT)},
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
		Mounts: tc.ContainerMounts{
			tc.ContainerMount{
				Source: tc.GenericBindMountSource{
					HostPath: g.generatedDataHostDir,
				},
				Target: tc.ContainerMountTarget(GENERATED_DATA_DIR_INSIDE_CONTAINER),
			},
		},
	}, nil
}

func (g Erigon) WaitUntilChainIsReady(waitTime time.Duration) error {
	waitForFirstBlock := tcwait.NewLogStrategy("Built block").WithPollInterval(1 * time.Second).WithStartupTimeout(waitTime)
	return waitForFirstBlock.WaitUntilReady(context.Background(), *g.GetContainer())
}

// TODO copy genesis file to /hpme/erigon?
func (g Erigon) buildInitScript() (string, error) {
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
		--datadir={{.ExecutionDir}} --rpc.allow-unprotected-txs --rpc.txfeecap=0 --allow-insecure-unlock \
		--nodiscover --networkid={{.ChainID}}`

	data := struct {
		HttpPort         string
		ChainID          int
		GeneratedDataDir string
		JwtFileLocation  string
		ExecutionDir     string
	}{
		HttpPort:         TX_GETH_HTTP_PORT,
		ChainID:          g.chainConfg.ChainID,
		GeneratedDataDir: GENERATED_DATA_DIR_INSIDE_CONTAINER,
		JwtFileLocation:  JWT_SECRET_FILE_LOCATION_INSIDE_CONTAINER,
		ExecutionDir:     "/home/erigon/execution-data",
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

func (g *Erigon) GetContainerType() ContainerType {
	return ContainerType_Erigon
}
