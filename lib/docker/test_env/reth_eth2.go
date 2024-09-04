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

// NewRethEth2 starts a new Reth Eth2 node running in Docker
func NewRethEth2(networks []string, chainConfig *config.EthereumChainConfig, generatedDataHostDir, generatedDataContainerDir string, consensusLayer config.ConsensusLayer, opts ...EnvComponentOption) (*Reth, error) {
	parts := strings.Split(ethereum.DefaultRethEth2Image, ":")
	g := &Reth{
		EnvComponent: EnvComponent{
			ContainerName:    fmt.Sprintf("%s-%s", "reth-eth2", uuid.NewString()[0:8]),
			Networks:         networks,
			ContainerImage:   parts[0],
			ContainerVersion: parts[1],
			LogLevel:         "debug",
			StartupTimeout:   120 * time.Second,
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
		g.EnvComponent.ContainerName = fmt.Sprintf("%s-%s-%s", "reth-eth2", strings.Replace(g.ContainerVersion, ".", "_", -1), uuid.NewString()[0:8])
	}
	// if the internal docker repo is set then add it to the version
	g.EnvComponent.ContainerImage = mirror.AddMirrorToImageIfSet(g.EnvComponent.ContainerImage)
	return g, nil
}

func (g *Reth) getEth2ContainerRequest() (*tc.ContainerRequest, error) {
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
		ExposedPorts:  []string{NatPortFormat(DEFAULT_EVM_NODE_HTTP_PORT), NatPortFormat("8546"), NatPortFormat(DEFAULT_EVM_NODE_WS_PORT), NatPortFormat(ETH2_EXECUTION_PORT)},
		WaitingFor: tcwait.ForAll(
			tcwait.ForLog("Starting consensus engine").
				WithStartupTimeout(g.StartupTimeout).
				WithPollInterval(1 * time.Second),
		),
		User: "0:0",
		Entrypoint: []string{
			"sh",
			"/root/init.sh",
		},
		Files: []tc.ContainerFile{
			{
				HostFilePath:      initFile.Name(),
				ContainerFilePath: "/root/init.sh",
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

func (g *Reth) buildPosInitScript() (string, error) {
	initTemplate := `#!/bin/bash
	echo "Copied genesis file to {{.ExecutionDir}}"
	mkdir -p {{.ExecutionDir}}
	cp {{.GeneratedDataDir}}/genesis.json {{.ExecutionDir}}/genesis.json
	echo "Running reth init"
	reth init --datadir={{.ExecutionDir}} --chain={{.ExecutionDir}}/genesis.json -vvvv
	echo "Execution dir: {{.ExecutionDir}}"
	exit_code=$?
	if [ $exit_code -ne 0 ]; then
		echo "Reth init failed with exit code $exit_code"
		exit 1
	fi

	echo "Starting Reth..."
	command="reth node -d --ipcdisable --http --http.api=All --http.addr=0.0.0.0 --http.corsdomain='*' --http.port={{.HttpPort}} --ws --ws.addr=0.0.0.0 --ws.origins='*' --ws.port={{.WsPort}} --ws.api=All --authrpc.addr=0.0.0.0 --authrpc.jwtsecret={{.JwtFileLocation}} --chain={{.ExecutionDir}}/genesis.json --datadir={{.ExecutionDir}}"

	if [ "{{.LogLevel}}" = "error" ]; then
		command="$command -v"
	elif [ "{{.LogLevel}}" = "warning" ]; then
		command="$command -vv"
	elif [ "{{.LogLevel}}" = "info" ]; then
		command="$command -vvv"
	elif [ "{{.LogLevel}}" = "debug" ]; then
		command="$command -vvvv"
	elif [ "{{.LogLevel}}" = "trace" ]; then
		command="$command -vvvvv"
	fi

	echo "Running command: $command"
	eval $command`

	data := struct {
		HttpPort            string
		WsPort              string
		ChainID             int
		GeneratedDataDir    string
		JwtFileLocation     string
		ExecutionDir        string
		ExtraExecutionFlags string
		LogLevel            string
	}{
		HttpPort:         DEFAULT_EVM_NODE_HTTP_PORT,
		WsPort:           DEFAULT_EVM_NODE_WS_PORT,
		ChainID:          g.chainConfig.ChainID,
		GeneratedDataDir: g.generatedDataContainerDir,
		JwtFileLocation:  getJWTSecretFileLocationInsideContainer(g.generatedDataContainerDir),
		ExecutionDir:     "/root/.local",
		LogLevel:         g.LogLevel,
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
