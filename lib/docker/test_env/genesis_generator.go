package test_env

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/config"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/docker"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/docker/ethereum"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror"
)

var generatorForkToImageMap = map[ethereum.Fork]string{
	ethereum.EthereumFork_Shanghai: ethereum.GenesisGeneratorShanghaiImage,
	ethereum.EthereumFork_Deneb:    ethereum.GenesisGeneratorDenebImage,
}

var generatorForkToDataDirMap = map[ethereum.Fork]string{
	ethereum.EthereumFork_Shanghai: "/data/custom_config_data",
	ethereum.EthereumFork_Deneb:    "/data/metadata",
}

type EthGenesisGenerator struct {
	EnvComponent
	chainConfig               config.EthereumChainConfig
	l                         zerolog.Logger
	generatedDataHostDir      string
	generatedDataContainerDir string
	t                         *testing.T
}

// NewEthGenesisGenerator initializes a new EthGenesisGenerator instance for generating Ethereum genesis data.
// It requires the Ethereum chain configuration, a directory for generated data, and the last fork used.
// This function is essential for setting up the environment to create and manage Ethereum genesis files.
func NewEthGenesisGenerator(chainConfig config.EthereumChainConfig, generatedDataHostDir string, lastFork ethereum.Fork, opts ...EnvComponentOption) (*EthGenesisGenerator, error) {
	genesisGeneratorImage, ok := generatorForkToImageMap[lastFork]
	if !ok {
		return nil, fmt.Errorf("unknown fork: %s", lastFork)
	}

	generatedDataContainerDir, ok := generatorForkToDataDirMap[lastFork]
	if !ok {
		return nil, fmt.Errorf("unknown fork: %s", lastFork)
	}

	parts := strings.Split(genesisGeneratorImage, ":")
	g := &EthGenesisGenerator{
		EnvComponent: EnvComponent{
			ContainerName:    fmt.Sprintf("%s-%s", "eth-genesis-generator", uuid.NewString()[0:8]),
			ContainerImage:   parts[0],
			ContainerVersion: parts[1],
			StartupTimeout:   30 * time.Second,
		},
		chainConfig:               chainConfig,
		generatedDataHostDir:      generatedDataHostDir,
		generatedDataContainerDir: generatedDataContainerDir,
		l:                         log.Logger,
	}
	g.SetDefaultHooks()
	for _, opt := range opts {
		opt(&g.EnvComponent)
	}
	// if the internal docker repo is set then add it to the version
	g.EnvComponent.ContainerImage = mirror.AddMirrorToImageIfSet(g.EnvComponent.ContainerImage)
	return g, nil
}

// WithTestInstance sets up the EthGenesisGenerator for testing by assigning a logger and test context.
// This allows for better logging and error tracking during test execution.
func (g *EthGenesisGenerator) WithTestInstance(t *testing.T) *EthGenesisGenerator {
	g.l = logging.GetTestLogger(t)
	g.t = t
	return g
}

// StartContainer initializes and starts the Ethereum genesis generation container.
// It ensures the container is ready for use, logging the process and any errors encountered.
func (g *EthGenesisGenerator) StartContainer() error {
	r, err := g.getContainerRequest(g.Networks)
	if err != nil {
		return err
	}

	l := logging.GetTestContainersGoTestLogger(g.t)
	_, err = docker.StartContainerWithRetry(g.l, tc.GenericContainerRequest{
		ContainerRequest: *r,
		Reuse:            true,
		Started:          true,
		Logger:           l,
	})
	if err != nil {
		return fmt.Errorf("cannot start eth genesis generation container: %w", err)
	}

	g.l.Info().Str("containerName", g.ContainerName).
		Msg("Started Eth Genesis container")

	return nil
}

func (g *EthGenesisGenerator) getContainerRequest(networks []string) (*tc.ContainerRequest, error) {
	valuesEnv, err := os.CreateTemp("", "values.env")
	if err != nil {
		return nil, err
	}

	bc, err := generateEnvValues(&g.chainConfig)
	if err != nil {
		return nil, err
	}
	_, err = valuesEnv.WriteString(bc)
	if err != nil {
		return nil, err
	}

	elGenesisFile, err := os.CreateTemp("", "genesis-config.yaml")
	if err != nil {
		return nil, err
	}
	_, err = elGenesisFile.WriteString(elGenesisConfig)
	if err != nil {
		return nil, err
	}

	clGenesisFile, err := os.CreateTemp("", "config.yaml")
	if err != nil {
		return nil, err
	}
	_, err = clGenesisFile.WriteString(clGenesisConfig)
	if err != nil {
		return nil, err
	}

	mnemonicsFile, err := os.CreateTemp("", "mnemonics.yaml")
	if err != nil {
		return nil, err
	}
	_, err = mnemonicsFile.WriteString(mnemonics)
	if err != nil {
		return nil, err
	}

	return &tc.ContainerRequest{
		Name:          g.ContainerName,
		Image:         g.GetImageWithVersion(),
		ImagePlatform: "linux/x86_64",
		Networks:      networks,
		WaitingFor: tcwait.ForAll(
			tcwait.ForLog("+ terminalTotalDifficulty=0"),
			tcwait.ForLog(fmt.Sprintf("+ sed -i 's/TERMINAL_TOTAL_DIFFICULTY:.*/TERMINAL_TOTAL_DIFFICULTY: 0/' %s/config.yaml", g.generatedDataContainerDir)).
				WithPollInterval(1*time.Second),
		).WithStartupTimeoutDefault(g.StartupTimeout),
		Cmd: []string{"all"},
		Files: []tc.ContainerFile{
			{
				HostFilePath:      valuesEnv.Name(),
				ContainerFilePath: "/config/values.env",
				FileMode:          0644,
			},
			{
				HostFilePath:      elGenesisFile.Name(),
				ContainerFilePath: "/config/el/genesis-config.yaml",
				FileMode:          0644,
			},
			{
				HostFilePath:      clGenesisFile.Name(),
				ContainerFilePath: "/config/cl/config.yaml",
				FileMode:          0644,
			},
			{
				HostFilePath:      mnemonicsFile.Name(),
				ContainerFilePath: "/config/cl/mnemonics.yaml",
				FileMode:          0644,
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

// GetGeneratedDataContainerDir returns the directory path for the generated data container.
// This is useful for accessing the location where the genesis data is stored after generation.
func (g *EthGenesisGenerator) GetGeneratedDataContainerDir() string {
	return g.generatedDataContainerDir
}
