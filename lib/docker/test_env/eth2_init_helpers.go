package test_env

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
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
)

type AfterGenesisHelper struct {
	EnvComponent
	chainConfig     config.EthereumChainConfig
	addressesToFund []string
	l               zerolog.Logger
	t               *testing.T
	posContainerSettings
}

// NewInitHelper initializes a new AfterGenesisHelper instance with the provided chain configuration and directory paths.
// It sets up necessary environment components and options, facilitating the management of post-genesis operations in a blockchain network.
func NewInitHelper(chainConfig config.EthereumChainConfig, generatedDataHostDir, generatedDataContainerDir string, opts ...EnvComponentOption) *AfterGenesisHelper {
	g := &AfterGenesisHelper{
		EnvComponent: EnvComponent{
			ContainerName:  fmt.Sprintf("%s-%s", "after-genesis-helper", uuid.NewString()[0:8]),
			StartupTimeout: 20 * time.Second,
		},
		chainConfig:          chainConfig,
		posContainerSettings: posContainerSettings{generatedDataHostDir: generatedDataHostDir, generatedDataContainerDir: generatedDataContainerDir},
		l:                    log.Logger,
		addressesToFund:      []string{},
	}
	g.SetDefaultHooks()
	for _, opt := range opts {
		opt(&g.EnvComponent)
	}
	return g
}

// WithTestInstance sets up the AfterGenesisHelper for testing by assigning a logger and test context.
// This allows for better logging and error tracking during test execution.
func (g *AfterGenesisHelper) WithTestInstance(t *testing.T) *AfterGenesisHelper {
	g.l = logging.GetTestLogger(t)
	g.t = t
	return g
}

// StartContainer initializes and starts the After Genesis Helper container.
// It handles container requests and logs the startup process, ensuring the container is ready for use.
func (g *AfterGenesisHelper) StartContainer() error {
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
		return fmt.Errorf("cannot start after genesis helper container: %w", err)
	}

	g.l.Info().Str("containerName", g.ContainerName).
		Msg("Started After Genesis Helper container")

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
		Image:         defaultEth2ValToolsImage,
		ImagePlatform: "linux/x86_64",
		Networks:      networks,
		WaitingFor: NewExitCodeStrategy().WithExitCode(0).
			WithPollInterval(1 * time.Second).WithTimeout(g.StartupTimeout),
		Entrypoint: []string{"sh", "/init.sh"},
		Files: []tc.ContainerFile{
			{
				HostFilePath:      initScriptFile.Name(),
				ContainerFilePath: "/init.sh",
				FileMode:          0744,
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

func (g *AfterGenesisHelper) buildInitScript() (string, error) {
	initTemplate := `#!/bin/bash
echo "Saving wallet password to {{.WalletPasswordFileLocation}}"
echo "{{.WalletPassword}}" > {{.WalletPasswordFileLocation}}
echo "Saving execution client keystore file to {{.AccountKeystoreFileLocation}}"
mkdir -p {{.KeystoreDirLocation}}
echo '{"address":"123463a4b065722e99115d6c222f267d9cabb524","crypto":{"cipher":"aes-128-ctr","ciphertext":"93b90389b855889b9f91c89fd15b9bd2ae95b06fe8e2314009fc88859fc6fde9","cipherparams":{"iv":"9dc2eff7967505f0e6a40264d1511742"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"c07503bb1b66083c37527cd8f06f8c7c1443d4c724767f625743bd47ae6179a4"},"mac":"6d359be5d6c432d5bbb859484009a4bf1bd71b76e89420c380bd0593ce25a817"},"id":"622df904-0bb1-4236-b254-f1b8dfdff1ec","version":3}' > {{.AccountKeystoreFileLocation}}
echo "Saving execution client account password to {{.AccountPasswordFileLocation}}"
echo "" > {{.AccountPasswordFileLocation}}
echo "Saving jwt secret to {{.JwtFileLocation}}"
echo "0xfad2709d0bb03bf0e8ba3c99bea194575d3e98863133d1af638ed056d1d59345" > {{.JwtFileLocation}}
echo "All done!"
echo
echo "------------------------------------------------------------------"
formatted_genesis_date=$(date -d "@{{.GenesisTimestamp}}" '+%Y-%m-%d %H:%M:%S')
echo "Chain genesis timestamp: $formatted_genesis_date UTC"

current_timestamp=$(date +%s)
time_diff=$(({{.GenesisTimestamp}} - current_timestamp))
echo "More or less $time_diff seconds from now"
echo "------------------------------------------------------------------"
`

	data := struct {
		WalletPassword              string
		WalletPasswordFileLocation  string
		AccountPasswordFileLocation string
		JwtFileLocation             string
		AccountKeystoreFileLocation string
		KeystoreDirLocation         string
		GenesisTimestamp            int
	}{
		WalletPassword:              WALLET_PASSWORD,
		WalletPasswordFileLocation:  getValidatorWalletPasswordFileInsideContainer(g.generatedDataContainerDir),
		AccountPasswordFileLocation: getAccountPasswordFileInsideContainer(g.generatedDataContainerDir),
		JwtFileLocation:             getJWTSecretFileLocationInsideContainer(g.generatedDataContainerDir),
		AccountKeystoreFileLocation: getAccountKeystoreFileInsideContainer(g.generatedDataContainerDir),
		KeystoreDirLocation:         getKeystoreDirLocationInsideContainer(g.generatedDataContainerDir),
		GenesisTimestamp:            g.chainConfig.GenesisTimestamp,
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
