package test_env

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	tc "github.com/testcontainers/testcontainers-go"

	"github.com/smartcontractkit/chainlink-testing-framework/docker"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/mirror"
)

type ValKeysGeneretor struct {
	EnvComponent
	chainConfig        *EthereumChainConfig
	l                  zerolog.Logger
	valKeysHostDataDir string
	addressesToFund    []string
	t                  *testing.T
	image              string
}

func NewValKeysGeneretor(chainConfig *EthereumChainConfig, valKeysHostDataDir string, opts ...EnvComponentOption) (*ValKeysGeneretor, error) {
	// currently it uses latest (no fixed version available)
	dockerImage, err := mirror.GetImage("protolambda/eth2-val-tools:l")
	if err != nil {
		return nil, err
	}

	g := &ValKeysGeneretor{
		EnvComponent: EnvComponent{
			ContainerName: fmt.Sprintf("%s-%s", "val-keys-generator", uuid.NewString()[0:8]),
		},
		chainConfig:        chainConfig,
		valKeysHostDataDir: valKeysHostDataDir,
		l:                  log.Logger,
		addressesToFund:    []string{},
		image:              dockerImage,
	}
	g.SetDefaultHooks()
	for _, opt := range opts {
		opt(&g.EnvComponent)
	}

	return g, nil
}

func (g *ValKeysGeneretor) WithTestInstance(t *testing.T) *ValKeysGeneretor {
	g.l = logging.GetTestLogger(t)
	g.t = t
	return g
}

func (g *ValKeysGeneretor) StartContainer() error {
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

func (g *ValKeysGeneretor) getContainerRequest(networks []string) (*tc.ContainerRequest, error) {
	return &tc.ContainerRequest{
		Name:          g.ContainerName,
		Image:         g.image,
		ImagePlatform: "linux/x86_64",
		Networks:      networks,
		WaitingFor: NewExitCodeStrategy().WithExitCode(0).
			WithPollInterval(1 * time.Second).WithTimeout(10 * time.Second),
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
		Mounts: tc.ContainerMounts{
			tc.ContainerMount{
				Source: tc.GenericBindMountSource{
					HostPath: g.valKeysHostDataDir,
				},
				Target: tc.ContainerMountTarget(GENERATED_VALIDATOR_KEYS_DIR_INSIDE_CONTAINER),
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
