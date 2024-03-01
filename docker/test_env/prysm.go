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
	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/docker"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/mirror"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/testcontext"
)

const (
	PRYSM_QUERY_RPC_PORT = "3500"
	PRYSM_NODE_RPC_PORT  = "4000"

	defaultPrysmBeaconChainImage = "gcr.io/prysmaticlabs/prysm/beacon-chain:v4.1.1"
	defaultPyrsmValidatorImage   = "gcr.io/prysmaticlabs/prysm/validator:v4.1.1"
)

type PrysmBeaconChain struct {
	EnvComponent
	InternalBeaconRpcProvider string
	InternalQueryRpcUrl       string
	ExternalBeaconRpcProvider string
	ExternalQueryRpcUrl       string
	generatedDataHostDir      string
	gethInternalExecutionURL  string
	chainConfig               *EthereumChainConfig
	l                         zerolog.Logger
	t                         *testing.T
}

func NewPrysmBeaconChain(networks []string, chainConfig *EthereumChainConfig, customConfigDataDir, gethExecutionURL string, opts ...EnvComponentOption) (*PrysmBeaconChain, error) {
	parts := strings.Split(defaultPrysmBeaconChainImage, ":")
	g := &PrysmBeaconChain{
		EnvComponent: EnvComponent{
			ContainerName:    fmt.Sprintf("%s-%s", "prysm-beacon-chain", uuid.NewString()[0:8]),
			Networks:         networks,
			ContainerImage:   parts[0],
			ContainerVersion: parts[1],
		},
		chainConfig:              chainConfig,
		generatedDataHostDir:     customConfigDataDir,
		gethInternalExecutionURL: gethExecutionURL,
		l:                        logging.GetTestLogger(nil),
	}
	g.SetDefaultHooks()
	for _, opt := range opts {
		opt(&g.EnvComponent)
	}
	// if the internal docker repo is set then add it to the version
	g.EnvComponent.ContainerImage = mirror.AddMirrorToImageIfSet(g.EnvComponent.ContainerImage)
	return g, nil
}

func (g *PrysmBeaconChain) WithTestInstance(t *testing.T) *PrysmBeaconChain {
	g.l = logging.GetTestLogger(t)
	g.t = t
	return g
}

func (g *PrysmBeaconChain) StartContainer() error {
	r, err := g.getContainerRequest(g.Networks)
	if err != nil {
		return err
	}

	l := logging.GetTestContainersGoTestLogger(g.t)
	ct, err := docker.StartContainerWithRetry(g.l, tc.GenericContainerRequest{
		ContainerRequest: *r,
		Reuse:            true,
		Started:          true,
		Logger:           l,
	})
	if err != nil {
		return fmt.Errorf("cannot start prysm beacon chain container: %w", err)
	}

	host, err := GetHost(testcontext.Get(g.t), ct)
	if err != nil {
		return err
	}
	queryPort, err := ct.MappedPort(testcontext.Get(g.t), NatPort(PRYSM_QUERY_RPC_PORT))
	if err != nil {
		return err
	}

	externalRcpPort, err := ct.MappedPort(testcontext.Get(g.t), NatPort(PRYSM_NODE_RPC_PORT))
	if err != nil {
		return err
	}

	g.Container = ct
	g.InternalBeaconRpcProvider = fmt.Sprintf("%s:%s", g.ContainerName, PRYSM_NODE_RPC_PORT)
	g.InternalQueryRpcUrl = fmt.Sprintf("%s:%s", g.ContainerName, PRYSM_QUERY_RPC_PORT)
	g.ExternalBeaconRpcProvider = FormatHttpUrl(host, externalRcpPort.Port())
	g.ExternalQueryRpcUrl = FormatHttpUrl(host, queryPort.Port())

	g.l.Info().Str("containerName", g.ContainerName).
		Msg("Started Prysm Beacon Chain container")

	return nil
}

func (g *PrysmBeaconChain) getContainerRequest(networks []string) (*tc.ContainerRequest, error) {
	return &tc.ContainerRequest{
		Name:          g.ContainerName,
		Image:         g.GetImageWithVersion(),
		ImagePlatform: "linux/amd64",
		Networks:      networks,
		WaitingFor: tcwait.ForAll(
			tcwait.ForLog("Starting beacon node").
				WithStartupTimeout(g.chainConfig.GetDefaultWaitDuration()).
				WithPollInterval(2 * time.Second),
		),
		Cmd: []string{
			"--accept-terms-of-use",
			"--datadir=/consensus-data",
			fmt.Sprintf("--chain-config-file=%s/config.yaml", GENERATED_DATA_DIR_INSIDE_CONTAINER),
			fmt.Sprintf("--genesis-state=%s/genesis.ssz", GENERATED_DATA_DIR_INSIDE_CONTAINER),
			fmt.Sprintf("--execution-endpoint=%s", g.gethInternalExecutionURL),
			"--rpc-host=0.0.0.0",
			"--grpc-gateway-host=0.0.0.0",
			"--grpc-gateway-corsdomain=*",
			"--suggested-fee-recipient=0x8943545177806ED17B9F23F0a21ee5948eCaa776",
			"--subscribe-all-subnets=true",
			fmt.Sprintf("--jwt-secret=%s", JWT_SECRET_FILE_LOCATION_INSIDE_CONTAINER),
			// mine, modify when running multi-node
			"--minimum-peers-per-subnet=0",
			"--min-sync-peers=0",
			"--interop-eth1data-votes",
		},
		ExposedPorts: []string{NatPortFormat(PRYSM_NODE_RPC_PORT), NatPortFormat(PRYSM_QUERY_RPC_PORT)},
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

type PrysmValidator struct {
	EnvComponent
	chainConfig               *EthereumChainConfig
	internalBeaconRpcProvider string
	valKeysDir                string
	generatedDataHostDir      string
	l                         zerolog.Logger
	t                         *testing.T
}

func NewPrysmValidator(networks []string, chainConfig *EthereumChainConfig, generatedDataHostDir, valKeysDir, internalBeaconRpcProvider string, opts ...EnvComponentOption) (*PrysmValidator, error) {
	parts := strings.Split(defaultPyrsmValidatorImage, ":")
	g := &PrysmValidator{
		EnvComponent: EnvComponent{
			ContainerName:    fmt.Sprintf("%s-%s", "prysm-validator", uuid.NewString()[0:8]),
			Networks:         networks,
			ContainerImage:   parts[0],
			ContainerVersion: parts[1],
		},
		chainConfig:               chainConfig,
		generatedDataHostDir:      generatedDataHostDir,
		valKeysDir:                valKeysDir,
		internalBeaconRpcProvider: internalBeaconRpcProvider,
		l:                         logging.GetTestLogger(nil),
	}
	g.SetDefaultHooks()
	for _, opt := range opts {
		opt(&g.EnvComponent)
	}
	// if the internal docker repo is set then add it to the version
	g.EnvComponent.ContainerImage = mirror.AddMirrorToImageIfSet(g.EnvComponent.ContainerImage)
	return g, nil
}

func (g *PrysmValidator) WithTestInstance(t *testing.T) *PrysmValidator {
	g.l = logging.GetTestLogger(t)
	g.t = t
	return g
}

func (g *PrysmValidator) StartContainer() error {
	r, err := g.getContainerRequest(g.Networks)
	if err != nil {
		return err
	}

	l := logging.GetTestContainersGoTestLogger(g.t)
	ct, err := docker.StartContainerWithRetry(g.l, tc.GenericContainerRequest{
		ContainerRequest: *r,
		Reuse:            true,
		Started:          true,
		Logger:           l,
	})
	if err != nil {
		return fmt.Errorf("cannot start prysm validator container: %w", err)
	}

	g.Container = ct

	g.l.Info().Str("containerName", g.ContainerName).
		Msg("Started Prysm Validator container")

	return nil
}

func (g *PrysmValidator) getContainerRequest(networks []string) (*tc.ContainerRequest, error) {
	return &tc.ContainerRequest{
		Name:          g.ContainerName,
		Image:         g.GetImageWithVersion(),
		Networks:      networks,
		ImagePlatform: "linux/x86_64",
		WaitingFor: tcwait.ForAll(
			tcwait.ForLog("Beacon chain started").
				WithStartupTimeout(g.chainConfig.GetDefaultWaitDuration()).
				WithPollInterval(2 * time.Second),
		),
		Cmd: []string{
			"--accept-terms-of-use",
			fmt.Sprintf("--chain-config-file=%s/config.yaml", GENERATED_DATA_DIR_INSIDE_CONTAINER),
			fmt.Sprintf("--beacon-rpc-provider=%s", g.internalBeaconRpcProvider),
			"--datadir=/consensus-data",
			"--suggested-fee-recipient=0x8943545177806ED17B9F23F0a21ee5948eCaa776",
			fmt.Sprintf("--wallet-dir=%s/prysm", NODE_0_DIR_INSIDE_CONTAINER),
			fmt.Sprintf("--wallet-password-file=%s", VALIDATOR_WALLET_PASSWORD_FILE_INSIDE_CONTAINER),
		},
		HostConfigModifier: func(hostConfig *container.HostConfig) {
			hostConfig.Mounts = append(hostConfig.Mounts, mount.Mount{
				Type:     mount.TypeBind,
				Source:   g.valKeysDir,
				Target:   GENERATED_VALIDATOR_KEYS_DIR_INSIDE_CONTAINER,
				ReadOnly: false,
			}, mount.Mount{
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
