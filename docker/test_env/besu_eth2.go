package test_env

import (
	"fmt"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/google/uuid"
	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/mirror"
)

// NewBesuEth2 starts a new Besu Eth2 node running in Docker
func NewBesuEth2(networks []string, chainConfig *EthereumChainConfig, generatedDataHostDir string, consensusLayer ConsensusLayer, opts ...EnvComponentOption) (*Besu, error) {
	parts := strings.Split(defaultBesuEth2Image, ":")
	g := &Besu{
		EnvComponent: EnvComponent{
			ContainerName:    fmt.Sprintf("%s-%s", "besu-eth2", uuid.NewString()[0:8]),
			Networks:         networks,
			ContainerImage:   parts[0],
			ContainerVersion: parts[1],
		},
		chainConfig:     chainConfig,
		posSettings:     posSettings{generatedDataHostDir: generatedDataHostDir},
		consensusLayer:  consensusLayer,
		l:               logging.GetTestLogger(nil),
		ethereumVersion: EthereumVersion_Eth2,
	}
	g.SetDefaultHooks()
	for _, opt := range opts {
		opt(&g.EnvComponent)
	}

	if !g.WasRecreated {
		// set the container name again after applying functional options as version might have changed
		g.EnvComponent.ContainerName = fmt.Sprintf("%s-%s-%s", "besu-eth2", strings.Replace(g.ContainerVersion, ".", "_", -1), uuid.NewString()[0:8])
	}
	// if the internal docker repo is set then add it to the version
	g.EnvComponent.ContainerImage = mirror.AddMirrorToImageIfSet(g.EnvComponent.ContainerImage)

	return g, nil
}

func (g *Besu) getEth2ContainerRequest() (*tc.ContainerRequest, error) {
	return &tc.ContainerRequest{
		Name:     g.ContainerName,
		Image:    g.GetImageWithVersion(),
		Networks: g.Networks,
		// ImagePlatform: "linux/x86_64", //don't even try this on Apple Silicon, the node won't start due to JVM error
		ExposedPorts: []string{NatPortFormat(DEFAULT_EVM_NODE_HTTP_PORT), NatPortFormat(DEFAULT_EVM_NODE_WS_PORT), NatPortFormat(ETH2_EXECUTION_PORT)},
		WaitingFor: tcwait.ForAll(
			tcwait.ForLog("Ethereum main loop is up").
				WithStartupTimeout(120 * time.Second).
				WithPollInterval(1 * time.Second),
		),
		User: "0:0", //otherwise in CI we get "permission denied" error, when trying to access data from mounted volume
		Cmd: []string{
			"--data-path=/opt/besu/execution-data",
			fmt.Sprintf("--genesis-file=%s/besu.json", GENERATED_DATA_DIR_INSIDE_CONTAINER),
			fmt.Sprintf("--network-id=%d", g.chainConfig.ChainID),
			"--host-allowlist=*",
			"--rpc-http-enabled=true",
			"--rpc-http-host=0.0.0.0",
			fmt.Sprintf("--rpc-http-port=%s", DEFAULT_EVM_NODE_HTTP_PORT),
			"--rpc-http-api=ADMIN,CLIQUE,ETH,NET,DEBUG,TXPOOL,ENGINE,TRACE,WEB3",
			"--rpc-http-cors-origins=*",
			"--rpc-ws-enabled=true",
			"--rpc-ws-host=0.0.0.0",
			fmt.Sprintf("--rpc-ws-port=%s", DEFAULT_EVM_NODE_WS_PORT),
			"--rpc-ws-api=ADMIN,CLIQUE,ETH,NET,DEBUG,TXPOOL,ENGINE,TRACE,WEB3",
			"--engine-rpc-enabled=true",
			fmt.Sprintf("--engine-jwt-secret=%s", JWT_SECRET_FILE_LOCATION_INSIDE_CONTAINER),
			"--engine-host-allowlist=*",
			fmt.Sprintf("--engine-rpc-port=%s", ETH2_EXECUTION_PORT),
			"--sync-mode=FULL",
			"--data-storage-format=BONSAI",
			// "--logging=DEBUG",
			"--rpc-tx-feecap=0",
		},
		Env: map[string]string{
			"JAVA_OPTS": "-agentlib:jdwp=transport=dt_socket,server=y,suspend=n",
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
