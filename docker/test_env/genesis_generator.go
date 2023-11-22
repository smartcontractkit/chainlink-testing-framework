package test_env

import (
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/docker"
	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"
)

type EthGenesisGeneretor struct {
	EnvComponent
	beaconChainConfig   BeaconChainConfig
	l                   zerolog.Logger
	customConfigDataDir string
	addressesToFund     []string
}

func NewEthGenesisGenerator(beaconChainConfig BeaconChainConfig, hostSharedDataDir string, opts ...EnvComponentOption) *EthGenesisGeneretor {
	g := &EthGenesisGeneretor{
		EnvComponent: EnvComponent{
			ContainerName: fmt.Sprintf("%s-%s", "eth-genesis-generator", uuid.NewString()[0:8]),
		},
		beaconChainConfig:   beaconChainConfig,
		customConfigDataDir: hostSharedDataDir,
		l:                   log.Logger,
		addressesToFund:     []string{},
	}
	for _, opt := range opts {
		opt(&g.EnvComponent)
	}
	return g
}

func (g *EthGenesisGeneretor) WithLogger(l zerolog.Logger) *EthGenesisGeneretor {
	g.l = l
	return g
}

func (g *EthGenesisGeneretor) WithFundedAccounts(addresses []string) *EthGenesisGeneretor {
	g.addressesToFund = addresses
	return g
}

func (g *EthGenesisGeneretor) StartContainer() error {
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
		return errors.Wrapf(err, "cannot start eth genesis generation container")
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

	bc, err := generateEnvValues(&g.beaconChainConfig)
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
		Image:         "tofelb/ethereum-genesis-generator:2.0.4-slots-per-epoch",
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
					HostPath: g.customConfigDataDir,
				},
				Target: "/data/custom_config_data",
			},
		},
	}, nil
}
