package test_env

import (
	"fmt"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/google/uuid"
	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/config"
	config_types "github.com/smartcontractkit/chainlink-testing-framework/lib/config/types"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/docker/ethereum"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror"
	docker_utils "github.com/smartcontractkit/chainlink-testing-framework/lib/utils/docker"
)

// NewBesuEth2 starts a new Besu Eth2 node running in Docker
func NewBesuEth2(networks []string, chainConfig *config.EthereumChainConfig, generatedDataHostDir, generatedDataContainerDir string, consensusLayer config.ConsensusLayer, opts ...EnvComponentOption) (*Besu, error) {
	parts := strings.Split(ethereum.DefaultBesuEth2Image, ":")
	g := &Besu{
		EnvComponent: EnvComponent{
			ContainerName:    fmt.Sprintf("%s-%s", "besu-eth2", uuid.NewString()[0:8]),
			Networks:         networks,
			ContainerImage:   parts[0],
			ContainerVersion: parts[1],
			StartupTimeout:   2 * time.Minute,
			l:                logging.GetTestLogger(nil),
		},
		chainConfig:          chainConfig,
		posContainerSettings: posContainerSettings{generatedDataHostDir: generatedDataHostDir, generatedDataContainerDir: generatedDataContainerDir},
		consensusLayer:       consensusLayer,
		ethereumVersion:      config_types.EthereumVersion_Eth2,
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
	cmd := []string{
		"--data-path=/opt/besu/execution-data",
		fmt.Sprintf("--genesis-file=%s/besu.json", g.generatedDataContainerDir),
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
		fmt.Sprintf("--engine-jwt-secret=%s", getJWTSecretFileLocationInsideContainer(g.generatedDataContainerDir)),
		"--engine-host-allowlist=*",
		fmt.Sprintf("--engine-rpc-port=%s", ETH2_EXECUTION_PORT),
		"--sync-mode=FULL",
		"--data-storage-format=BONSAI",
		fmt.Sprintf("--logging=%s", strings.ToUpper(g.LogLevel)),
		"--rpc-tx-feecap=0",
	}

	version, err := docker_utils.GetSemverFromImage(g.GetImageWithVersion())
	if err != nil {
		return nil, err
	}

	kgzConstraint, err := semver.NewConstraint(">=23.1 <24.0")
	if err != nil {
		return nil, fmt.Errorf("failed to parse constraint: %s", ">=23.1 && <23.7")
	}

	if kgzConstraint.Check(version) {
		cmd = append(cmd, "--kzg-trusted-setup", fmt.Sprintf("%s/trusted_setup.txt", g.generatedDataContainerDir))
	}

	bonsaiConstraint, err := semver.NewConstraint(">=24.6")
	if err != nil {
		return nil, fmt.Errorf("failed to parse constraint: %s", ">=24.6")
	}

	if bonsaiConstraint.Check(version) {
		// it crashes with sync-mode=FULL, and when we use a different sync mode then consensus client fails to propose correct blocks
		cmd = append(cmd, "--bonsai-limit-trie-logs-enabled=false")
	}

	return &tc.ContainerRequest{
		Name:     g.ContainerName,
		Image:    g.GetImageWithVersion(),
		Networks: g.Networks,
		// ImagePlatform: "linux/x86_64", //don't even try this on Apple Silicon, the node won't start due to JVM error
		ExposedPorts: []string{NatPortFormat(DEFAULT_EVM_NODE_HTTP_PORT), NatPortFormat(DEFAULT_EVM_NODE_WS_PORT), NatPortFormat(ETH2_EXECUTION_PORT)},
		WaitingFor: tcwait.ForAll(
			tcwait.ForLog("Ethereum main loop is up").
				WithPollInterval(1 * time.Second)).
			WithStartupTimeoutDefault(g.StartupTimeout),
		User: "0:0", //otherwise in CI we get "permission denied" error, when trying to access data from mounted volume
		Cmd:  cmd,
		Env: map[string]string{
			"JAVA_OPTS": "-agentlib:jdwp=transport=dt_socket,server=y,suspend=n",
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
