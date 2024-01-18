package test_env

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/docker"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/mirror"
)

type EthGenesisGeneretor struct {
	EnvComponent
	chainConfig          EthereumChainConfig
	l                    zerolog.Logger
	generatedDataHostDir string
	t                    *testing.T
	image                string
}

func NewEthGenesisGenerator(chainConfig EthereumChainConfig, generatedDataHostDir string, opts ...EnvComponentOption) (*EthGenesisGeneretor, error) {
	// currently it uses 2.0.5
	dockerImage, err := mirror.GetImage("tofelb/ethereum-genesis-generator:2.0.5")
	if err != nil {
		return nil, err
	}

	g := &EthGenesisGeneretor{
		EnvComponent: EnvComponent{
			ContainerName: fmt.Sprintf("%s-%s", "eth-genesis-generator", uuid.NewString()[0:8]),
		},
		chainConfig:          chainConfig,
		generatedDataHostDir: generatedDataHostDir,
		l:                    log.Logger,
		image:                dockerImage,
	}
	g.SetDefaultHooks()
	for _, opt := range opts {
		opt(&g.EnvComponent)
	}
	return g, nil
}

func (g *EthGenesisGeneretor) WithTestInstance(t *testing.T) *EthGenesisGeneretor {
	g.l = logging.GetTestLogger(t)
	g.t = t
	return g
}

func (g *EthGenesisGeneretor) WithImage(imageWithTag string) *EthGenesisGeneretor {
	g.image = imageWithTag
	return g
}

func (g *EthGenesisGeneretor) StartContainer() error {
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

func (g *EthGenesisGeneretor) getContainerRequest(networks []string) (*tc.ContainerRequest, error) {
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
		Image:         g.image,
		ImagePlatform: "linux/x86_64",
		Networks:      networks,
		WaitingFor: tcwait.ForAll(
			tcwait.ForLog("+ terminalTotalDifficulty=0"),
			tcwait.ForLog("+ sed -i 's/TERMINAL_TOTAL_DIFFICULTY:.*/TERMINAL_TOTAL_DIFFICULTY: 0/' /data/custom_config_data/config.yaml").
				WithStartupTimeout(20*time.Second).
				WithPollInterval(1*time.Second),
		),
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
