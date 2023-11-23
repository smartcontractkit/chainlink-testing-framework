package test_env

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	tc "github.com/testcontainers/testcontainers-go"

	"github.com/smartcontractkit/chainlink-testing-framework/docker"
)

type ValKeysGeneretor struct {
	EnvComponent
	chainConfig        *EthereumChainConfig
	l                  zerolog.Logger
	valKeysHostDataDir string
	addressesToFund    []string
}

func NewValKeysGeneretor(chainConfig *EthereumChainConfig, valKeysHostDataDir string, opts ...EnvComponentOption) *ValKeysGeneretor {
	g := &ValKeysGeneretor{
		EnvComponent: EnvComponent{
			ContainerName: fmt.Sprintf("%s-%s", "val-keys-generator", uuid.NewString()[0:8]),
		},
		chainConfig:        chainConfig,
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
	// walletPasswordFile, err := os.CreateTemp("", "password.txt")
	// if err != nil {
	// 	return nil, err
	// }
	// _, err = walletPasswordFile.WriteString(WALLET_PASSWORD)
	// if err != nil {
	// 	return nil, err
	// }

	// accountPasswordFile, err := os.CreateTemp("", "password.txt")
	// if err != nil {
	// 	return nil, err
	// }

	return &tc.ContainerRequest{
		Name:          g.ContainerName,
		Image:         "protolambda/eth2-val-tools:latest",
		ImagePlatform: "linux/x86_64",
		Networks:      networks,
		WaitingFor: NewExitCodeStrategy().WithExitCode(0).
			WithPollInterval(1 * time.Second).WithTimeout(10 * time.Second),
		// Entrypoint: []string{"sh", "/init.sh"},
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
		// Files: []tc.ContainerFile{
		// 	{
		// 		HostFilePath:      initScriptFile.Name(),
		// 		ContainerFilePath: "/init.sh",
		// 		FileMode:          0744,
		// 	},
		// },
		// Files: []tc.ContainerFile{
		// 	{
		// 		HostFilePath:      walletPasswordFile.Name(),
		// 		ContainerFilePath: WALLET_PASSWORD_FILE_INSIDE_CONTAINER,
		// 		FileMode:          0644,
		// 	},
		// 	{
		// 		HostFilePath:      accountPasswordFile.Name(),
		// 		ContainerFilePath: DEFAULT_EL_ACCOUNT_PASSWORD_FILE_INSIDE_CONTAINER,
		// 		FileMode:          0644,
		// 	},
		// },
		Mounts: tc.ContainerMounts{
			tc.ContainerMount{
				Source: tc.GenericBindMountSource{
					HostPath: g.valKeysHostDataDir,
				},
				Target: tc.ContainerMountTarget(GENERATED_VALIDATOR_KEYS_DIR_INSIDE_CONTAINER),
			},
		},
	}, nil
}
