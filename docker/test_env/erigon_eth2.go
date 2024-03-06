package test_env

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/google/uuid"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/mirror"
	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"
)

func NewErigonEth2(networks []string, chainConfg *EthereumChainConfig, generatedDataHostDir string, consensusLayer ConsensusLayer, opts ...EnvComponentOption) (*Erigon, error) {
	parts := strings.Split(defaultErigonEth2Image, ":")
	g := &Erigon{
		EnvComponent: EnvComponent{
			ContainerName:    fmt.Sprintf("%s-%s", "erigon-eth2", uuid.NewString()[0:8]),
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

	g.EnvComponent.ContainerName = fmt.Sprintf("%s-%s-%s", "erigon-eth2", strings.Replace(g.ContainerVersion, ".", "_", -1), uuid.NewString()[0:8])
	// if the internal docker repo is set then add it to the version
	g.EnvComponent.ContainerImage = mirror.AddMirrorToImageIfSet(g.EnvComponent.ContainerImage)
	return g, nil
}

func (g *Erigon) getEth2ContainerRequest() (*tc.ContainerRequest, error) {
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
