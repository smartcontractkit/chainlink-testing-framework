package test_env

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	tc "github.com/testcontainers/testcontainers-go"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/config"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/docker"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror"
)

const defaultEth2ValToolsImage = "protolambda/eth2-val-tools:latest"

type ValKeysGenerator struct {
	EnvComponent
	chainConfig        *config.EthereumChainConfig
	l                  zerolog.Logger
	valKeysHostDataDir string
	addressesToFund    []string
	t                  *testing.T
}

// NewValKeysGeneretor initializes a ValKeysGenerator for managing validator key generation in a specified Ethereum environment.
// It sets up the necessary container configuration and options, returning the generator instance or an error if initialization fails.
func NewValKeysGeneretor(chainConfig *config.EthereumChainConfig, valKeysHostDataDir string, opts ...EnvComponentOption) (*ValKeysGenerator, error) {
	parts := strings.Split(defaultEth2ValToolsImage, ":")
	g := &ValKeysGenerator{
		EnvComponent: EnvComponent{
			ContainerName:    fmt.Sprintf("%s-%s", "val-keys-generator", uuid.NewString()[0:8]),
			ContainerImage:   parts[0],
			ContainerVersion: parts[1],
			StartupTimeout:   1 * time.Minute,
		},
		chainConfig:        chainConfig,
		valKeysHostDataDir: valKeysHostDataDir,
		l:                  log.Logger,
		addressesToFund:    []string{},
	}
	g.SetDefaultHooks()
	for _, opt := range opts {
		opt(&g.EnvComponent)
	}

	// if the internal docker repo is set then add it to the version
	g.EnvComponent.ContainerImage = mirror.AddMirrorToImageIfSet(g.EnvComponent.ContainerImage)
	return g, nil
}

// WithTestInstance sets up a test logger and associates it with the ValKeysGenerator instance.
// This is useful for capturing log output during testing, ensuring that logs are directed to the test context.
func (g *ValKeysGenerator) WithTestInstance(t *testing.T) *ValKeysGenerator {
	g.l = logging.GetTestLogger(t)
	g.t = t
	return g
}

// StartContainer initializes and starts the validation keys generation container.
// It ensures the container is ready for use, logging the start event, and returns any errors encountered during the process.
func (g *ValKeysGenerator) StartContainer() error {
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
		return fmt.Errorf("cannot start val keys generation container: %w", err)
	}

	g.l.Info().Str("containerName", g.ContainerName).
		Msg("Started Val Keys Generation container")

	return nil
}

func (g *ValKeysGenerator) getContainerRequest(networks []string) (*tc.ContainerRequest, error) {
	return &tc.ContainerRequest{
		Name:          g.ContainerName,
		Image:         g.GetImageWithVersion(),
		ImagePlatform: "linux/x86_64",
		Networks:      networks,
		WaitingFor: NewExitCodeStrategy().WithExitCode(0).
			WithPollInterval(1 * time.Second).WithTimeout(g.StartupTimeout),
		Cmd: []string{"keystores",
			"--insecure",
			fmt.Sprintf("--prysm-pass=%s", WALLET_PASSWORD),
			fmt.Sprintf("--out-loc=%s", NODE_0_DIR_INSIDE_CONTAINER),
			fmt.Sprintf("--source-mnemonic=%s", VALIDATOR_BIP39_MNEMONIC),
			//if we ever have more than 1 node these indexes should be updated, so that we don't generate the same keys
			//e.g. if network has 2 nodes each with 10 validators, then the next source-min should be 10, and max should be 20
			"--source-min=0",
			fmt.Sprintf("--source-max=%d", g.chainConfig.ValidatorCount),
		},
		HostConfigModifier: func(hostConfig *container.HostConfig) {
			hostConfig.Mounts = append(hostConfig.Mounts, mount.Mount{
				Type:     mount.TypeBind,
				Source:   g.valKeysHostDataDir,
				Target:   GENERATED_VALIDATOR_KEYS_DIR_INSIDE_CONTAINER,
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
