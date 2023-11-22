package test_env

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/docker"
)

type ValKeysGeneretor struct {
	EnvComponent
	beaconChainConfig  BeaconChainConfig
	l                  zerolog.Logger
	valKeysHostDataDir string
	addressesToFund    []string
}

func NewValKeysGeneretor(beaconChainConfig BeaconChainConfig, valKeysHostDataDir string, opts ...EnvComponentOption) *ValKeysGeneretor {
	g := &ValKeysGeneretor{
		EnvComponent: EnvComponent{
			ContainerName: fmt.Sprintf("%s-%s", "val-keys-generator", uuid.NewString()[0:8]),
		},
		beaconChainConfig:  beaconChainConfig,
		valKeysHostDataDir: valKeysHostDataDir,
		l:                  log.Logger,
		addressesToFund:    []string{},
	}
	for _, opt := range opts {
		opt(&g.EnvComponent)
	}
	return g
}

func (g *ValKeysGeneretor) WithLogger(l zerolog.Logger) *ValKeysGeneretor {
	g.l = l
	return g
}

func (g *ValKeysGeneretor) WithFundedAccounts(addresses []string) *ValKeysGeneretor {
	g.addressesToFund = addresses
	return g
}

func (g *ValKeysGeneretor) StartContainer() error {
	r, err := g.getContainerRequest(g.Networks)
	if err != nil {
		return err
	}

	_, err = docker.StartContainerWithRetry(g.l, tc.GenericContainerRequest{
		ContainerRequest: *r,
		Reuse:            true,
		Started:          true,
		Logger:           &g.l,
	})
	if err != nil {
		return errors.Wrapf(err, "cannot start val keys generation container")
	}

	g.l.Info().Str("containerName", g.ContainerName).
		Msg("Started Val Keys Generation container")

	return nil
}

func (g *ValKeysGeneretor) getContainerRequest(networks []string) (*tc.ContainerRequest, error) {
	return &tc.ContainerRequest{
		Name:          g.ContainerName,
		Image:         "protolambda/eth2-val-tools:latest",
		ImagePlatform: "linux/x86_64",
		Networks:      networks,
		//TODO add new strategy: exit with code
		WaitingFor: tcwait.ForExit().
			WithPollInterval(1 * time.Second).WithExitTimeout(10 * time.Second),
		Cmd: []string{"keystores",
			"--insecure",
			"--prysm-pass=password",
			"--out-loc=/keys/node-0",
			fmt.Sprintf("--source-mnemonic=%s", VALIDATOR_BIPC39_MNEMONIC),
			//if we ever have more than 1 node these indexes should be updated, so that we don't generate the same keys
			//e.g. if network has 2 nodes each with 10 validators, then the next source-min should be 10, and max should be 20
			"--source-min=0",
			fmt.Sprintf("--source-max=%d", g.beaconChainConfig.ValidatorCount),
		},
		Mounts: tc.ContainerMounts{
			tc.ContainerMount{
				Source: tc.GenericBindMountSource{
					HostPath: g.valKeysHostDataDir,
				},
				Target: "/keys",
			},
		},
	}, nil
}
