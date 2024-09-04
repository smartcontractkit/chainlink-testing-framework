package test_env

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"strings"
	"time"

	config_types "github.com/smartcontractkit/chainlink-testing-framework/lib/config/types"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/google/uuid"
	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/config"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/docker/ethereum"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror"
)

// NewGethEth2 starts a new Geth Eth2 node running in Docker
func NewGethEth2(networks []string, chainConfig *config.EthereumChainConfig, generatedDataHostDir, generatedDataContainerDir string, consensusLayer config.ConsensusLayer, opts ...EnvComponentOption) (*Geth, error) {
	parts := strings.Split(ethereum.DefaultGethEth2Image, ":")
	g := &Geth{
		EnvComponent: EnvComponent{
			ContainerName:    fmt.Sprintf("%s-%s", "geth-eth2", uuid.NewString()[0:8]),
			Networks:         networks,
			ContainerImage:   parts[0],
			ContainerVersion: parts[1],
			StartupTimeout:   2 * time.Minute,
		},
		chainConfig:          chainConfig,
		posContainerSettings: posContainerSettings{generatedDataHostDir: generatedDataHostDir, generatedDataContainerDir: generatedDataContainerDir},
		consensusLayer:       consensusLayer,
		l:                    logging.GetTestLogger(nil),
		ethereumVersion:      config_types.EthereumVersion_Eth2,
	}
	g.SetDefaultHooks()
	for _, opt := range opts {
		opt(&g.EnvComponent)
	}

	if !g.WasRecreated {
		// set the container name again after applying functional options as version might have changed
		g.EnvComponent.ContainerName = fmt.Sprintf("%s-%s-%s", "geth-eth2", strings.Replace(g.ContainerVersion, ".", "_", -1), uuid.NewString()[0:8])
	}
	// if the internal docker repo is set then add it to the version
	g.EnvComponent.ContainerImage = mirror.AddMirrorToImageIfSet(g.EnvComponent.ContainerImage)
	return g, nil
}

func (g *Geth) getEth2ContainerRequest() (*tc.ContainerRequest, error) {
	initFile, err := os.CreateTemp("", "init.sh")
	if err != nil {
		return nil, err
	}

	initScriptContent, err := g.buildEth2dInitScript()
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
		ExposedPorts:  []string{NatPortFormat(DEFAULT_EVM_NODE_HTTP_PORT), NatPortFormat(DEFAULT_EVM_NODE_WS_PORT), NatPortFormat(ETH2_EXECUTION_PORT)},
		WaitingFor: tcwait.ForAll(
			tcwait.ForLog("WebSocket enabled").
				WithPollInterval(1 * time.Second)).
			WithStartupTimeoutDefault(g.StartupTimeout),
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
				Target:   g.generatedDataContainerDir,
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

func (g *Geth) buildEth2dInitScript() (string, error) {
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
		--graphql --graphql.corsdomain=* --unlock=0x123463a4b065722e99115d6c222f267d9cabb524 --verbosity={{.Verbosity}}`

	verbosity, err := g.logLevelToVerbosity()
	if err != nil {
		return "", err
	}

	data := struct {
		HttpPort             string
		WsPort               string
		ChainID              int
		GeneratedDataDir     string
		JwtFileLocation      string
		PasswordFileLocation string
		KeystoreDirLocation  string
		ExecutionDir         string
		Verbosity            int
	}{
		HttpPort:             DEFAULT_EVM_NODE_HTTP_PORT,
		WsPort:               DEFAULT_EVM_NODE_WS_PORT,
		ChainID:              g.chainConfig.ChainID,
		GeneratedDataDir:     g.generatedDataContainerDir,
		JwtFileLocation:      getJWTSecretFileLocationInsideContainer(g.generatedDataContainerDir),
		PasswordFileLocation: getAccountPasswordFileInsideContainer(g.generatedDataContainerDir),
		KeystoreDirLocation:  getKeystoreDirLocationInsideContainer(g.generatedDataContainerDir),
		ExecutionDir:         "/execution-data",
		Verbosity:            verbosity,
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
