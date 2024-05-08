package test_env

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/google/uuid"
	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/config"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/mirror"
)

// NewNethermindEth2 starts a new Nethermin Eth2 node running in Docker
func NewNethermindEth2(networks []string, chainConfig *config.EthereumChainConfig, generatedDataHostDir string, consensusLayer config.ConsensusLayer, opts ...EnvComponentOption) (*Nethermind, error) {
	parts := strings.Split(defaultNethermindEth2Image, ":")
	g := &Nethermind{
		EnvComponent: EnvComponent{
			ContainerName:    fmt.Sprintf("%s-%s", "nethermind-eth2", uuid.NewString()[0:8]),
			Networks:         networks,
			ContainerImage:   parts[0],
			ContainerVersion: parts[1],
		},
		generatedDataHostDir: generatedDataHostDir,
		chainConfig:          chainConfig,
		consensusLayer:       consensusLayer,
		l:                    logging.GetTestLogger(nil),
		ethereumVersion:      config.EthereumVersion_Eth2,
	}
	g.SetDefaultHooks()
	for _, opt := range opts {
		opt(&g.EnvComponent)
	}

	if !g.WasRecreated {
		// set the container name again after applying functional options as version might have changed
		g.EnvComponent.ContainerName = fmt.Sprintf("%s-%s-%s", "nethermind-eth2", strings.Replace(g.ContainerVersion, ".", "_", -1), uuid.NewString()[0:8])
	}
	// if the internal docker repo is set then add it to the version
	g.EnvComponent.ContainerImage = mirror.AddMirrorToImageIfSet(g.EnvComponent.ContainerImage)
	return g, nil
}

func (g *Nethermind) getEth2ContainerRequest() (*tc.ContainerRequest, error) {
	// create empty cfg file since if we don't pass any
	// default mainnet.cfg will be used
	noneCfg, err := os.CreateTemp("", "none.cfg")
	if err != nil {
		return nil, err
	}

	_, err = noneCfg.WriteString("{}")
	if err != nil {
		return nil, err
	}

	command := []string{
		"--datadir=/nethermind",
		"--config=/none.cfg",
		fmt.Sprintf("--Init.ChainSpecPath=%s/chainspec.json", GENERATED_DATA_DIR_INSIDE_CONTAINER),
		"--Init.DiscoveryEnabled=false",
		"--Init.WebSocketsEnabled=true",
		fmt.Sprintf("--JsonRpc.WebSocketsPort=%s", DEFAULT_EVM_NODE_WS_PORT),
		fmt.Sprintf("--Blocks.SecondsPerSlot=%d", g.chainConfig.SecondsPerSlot),
		"--JsonRpc.Enabled=true",
		"--JsonRpc.EnabledModules=net,eth,consensus,subscribe,web3,admin,txpool,debug,trace",
		"--JsonRpc.Host=0.0.0.0",
		fmt.Sprintf("--JsonRpc.Port=%s", DEFAULT_EVM_NODE_HTTP_PORT),
		"--JsonRpc.EngineHost=0.0.0.0",
		"--JsonRpc.EnginePort=" + ETH2_EXECUTION_PORT,
		fmt.Sprintf("--JsonRpc.JwtSecretFile=%s", JWT_SECRET_FILE_LOCATION_INSIDE_CONTAINER),
		fmt.Sprintf("--KeyStore.KeyStoreDirectory=%s", KEYSTORE_DIR_LOCATION_INSIDE_CONTAINER),
		"--KeyStore.BlockAuthorAccount=0x123463a4b065722e99115d6c222f267d9cabb524",
		"--KeyStore.UnlockAccounts=0x123463a4b065722e99115d6c222f267d9cabb524",
		fmt.Sprintf("--KeyStore.PasswordFiles=%s", ACCOUNT_PASSWORD_FILE_INSIDE_CONTAINER),
		"--Network.MaxActivePeers=0",
		"--Network.OnlyStaticPeers=true",
		"--HealthChecks.Enabled=true", // default slug /health
		fmt.Sprintf("--log=%s", strings.ToUpper(g.LogLevel)),
	}

	if g.LogLevel == "trace" {
		command = append(command, "--TraceStore.Enabled=true")
		command = append(command, "--Network.DiagTracerEnabled=true")
		command = append(command, "--TxPool.ReportMinutes=1")
	}

	return &tc.ContainerRequest{
		Name:            g.ContainerName,
		Image:           g.GetImageWithVersion(),
		Networks:        g.Networks,
		AlwaysPullImage: true,
		// ImagePlatform: "linux/x86_64",  //don't even try this on Apple Silicon, the node won't start due to .NET error
		ExposedPorts: []string{NatPortFormat(DEFAULT_EVM_NODE_HTTP_PORT), NatPortFormat(DEFAULT_EVM_NODE_WS_PORT), NatPortFormat(ETH2_EXECUTION_PORT)},
		WaitingFor: tcwait.ForAll(
			tcwait.ForLog("Nethermind initialization completed").
				WithStartupTimeout(120 * time.Second).
				WithPollInterval(1 * time.Second),
		),
		Cmd: command,
		Files: []tc.ContainerFile{
			{
				HostFilePath:      noneCfg.Name(),
				ContainerFilePath: "/none.cfg",
				FileMode:          0644,
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
