package test_env

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	tc "github.com/testcontainers/testcontainers-go"

	"github.com/smartcontractkit/chainlink-testing-framework/docker"
)

type AfterGenesisHelper struct {
	EnvComponent
	beaconChainConfig   EthereumChainConfig
	l                   zerolog.Logger
	customConfigDataDir string
	addressesToFund     []string
}

func NewInitHelper(beaconChainConfig EthereumChainConfig, customConfigDataDir string, opts ...EnvComponentOption) *AfterGenesisHelper {
	g := &AfterGenesisHelper{
		EnvComponent: EnvComponent{
			ContainerName: fmt.Sprintf("%s-%s", "after-genesis-helper", uuid.NewString()[0:8]),
		},
		beaconChainConfig:   beaconChainConfig,
		customConfigDataDir: customConfigDataDir,
		l:                   log.Logger,
		addressesToFund:     []string{},
	}
	for _, opt := range opts {
		opt(&g.EnvComponent)
	}
	return g
}

func (g *AfterGenesisHelper) WithLogger(l zerolog.Logger) *AfterGenesisHelper {
	g.l = l
	return g
}

func (g *AfterGenesisHelper) StartContainer() error {
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
		return errors.Wrapf(err, "cannot start init helper container")
	}

	g.l.Info().Str("containerName", g.ContainerName).
		Msg("Started Init Helper container")

	return nil
}

func (g *AfterGenesisHelper) getContainerRequest(networks []string) (*tc.ContainerRequest, error) {
	initScriptFile, err := os.CreateTemp("", "init.sh")
	if err != nil {
		return nil, err
	}

	initScript, err := g.buildInitScript()
	if err != nil {
		return nil, err
	}

	_, err = initScriptFile.WriteString(initScript)
	if err != nil {
		return nil, err
	}

	return &tc.ContainerRequest{
		Name:          g.ContainerName,
		Image:         "protolambda/eth2-val-tools:latest",
		ImagePlatform: "linux/x86_64",
		Networks:      networks,
		WaitingFor: NewExitCodeStrategy().WithExitCode(0).
			WithPollInterval(1 * time.Second).WithTimeout(10 * time.Second),
		Entrypoint: []string{"sh", "/init.sh"},
		Files: []tc.ContainerFile{
			{
				HostFilePath:      initScriptFile.Name(),
				ContainerFilePath: "/init.sh",
				FileMode:          0744,
			},
		},
		Mounts: tc.ContainerMounts{
			tc.ContainerMount{
				Source: tc.GenericBindMountSource{
					HostPath: g.customConfigDataDir,
				},
				Target: tc.ContainerMountTarget(GENERATED_DATA_DIR_INSIDE_CONTAINER),
			},
		},
	}, nil
}

func (g *AfterGenesisHelper) buildInitScript() (string, error) {
	initTemplate := `#!/bin/bash
echo "{{.WalletPassword}}" > {{.WalletPasswordFileLocation}}
echo "" > {{.AccountPasswordFileLocation}}
`

	data := struct {
		WalletPassword              string
		WalletPasswordFileLocation  string
		AccountPasswordFileLocation string
	}{
		WalletPassword:              WALLET_PASSWORD,
		WalletPasswordFileLocation:  VALIDATOR_WALLET_PASSWORD_FILE_INSIDE_CONTAINER,
		AccountPasswordFileLocation: EL_ACCOUNT_PASSWORD_FILE_INSIDE_CONTAINER,
	}

	t, err := template.New("init").Parse(initTemplate)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		os.Exit(1)
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, data)

	return buf.String(), err
}
