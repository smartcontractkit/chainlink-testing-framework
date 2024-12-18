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

	"github.com/smartcontractkit/chainlink-testing-framework/lib/config"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/docker"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/docker/ethereum"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
)

const (
	PRYSM_QUERY_RPC_PORT = "3500"
	PRYSM_NODE_RPC_PORT  = "4000"
)

var beaconForkToImageMap = map[ethereum.Fork]string{
	ethereum.EthereumFork_Shanghai: "gcr.io/prysmaticlabs/prysm/beacon-chain:v4.1.1",
	ethereum.EthereumFork_Deneb:    "gcr.io/prysmaticlabs/prysm/beacon-chain:v5.0.4",
}

var validatorForkToImageMap = map[ethereum.Fork]string{
	ethereum.EthereumFork_Shanghai: "gcr.io/prysmaticlabs/prysm/validator:v4.1.1",
	ethereum.EthereumFork_Deneb:    "gcr.io/prysmaticlabs/prysm/validator:v5.0.4",
}

type PrysmBeaconChain struct {
	EnvComponent
	InternalBeaconRpcProvider string
	InternalQueryRpcUrl       string
	ExternalBeaconRpcProvider string
	ExternalQueryRpcUrl       string
	gethInternalExecutionURL  string
	chainConfig               *config.EthereumChainConfig
	l                         zerolog.Logger
	t                         *testing.T
	posContainerSettings
}

// NewPrysmBeaconChain initializes a new Prysm beacon chain instance with the specified network configurations and parameters.
// It is used to set up a beacon chain for Ethereum 2.0, enabling users to run and manage a consensus layer node.
func NewPrysmBeaconChain(networks []string, chainConfig *config.EthereumChainConfig, generatedDataHostDir, generatedDataContainerDir, gethExecutionURL string, baseEthereumFork ethereum.Fork, opts ...EnvComponentOption) (*PrysmBeaconChain, error) {
	prysmBeaconChainImage, ok := beaconForkToImageMap[baseEthereumFork]
	if !ok {
		return nil, fmt.Errorf("unknown fork: %s", baseEthereumFork)
	}
	parts := strings.Split(prysmBeaconChainImage, ":")
	g := &PrysmBeaconChain{
		EnvComponent: EnvComponent{
			ContainerName:    fmt.Sprintf("%s-%s", "prysm-beacon-chain", uuid.NewString()[0:8]),
			Networks:         networks,
			ContainerImage:   parts[0],
			ContainerVersion: parts[1],
			StartupTimeout:   2 * time.Minute,
		},
		chainConfig:              chainConfig,
		posContainerSettings:     posContainerSettings{generatedDataHostDir: generatedDataHostDir, generatedDataContainerDir: generatedDataContainerDir},
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

// WithTestInstance sets up the PrysmBeaconChain for testing by assigning a test logger and the testing context.
// This allows for better logging and error tracking during test execution.
func (g *PrysmBeaconChain) WithTestInstance(t *testing.T) *PrysmBeaconChain {
	g.l = logging.GetTestLogger(t)
	g.t = t
	return g
}

// StartContainer initializes and starts the Prysm Beacon Chain container.
// It sets up the necessary RPC endpoints and logs the container's status.
// This function is essential for deploying a Prysm-based Ethereum 2.0 beacon chain in a Docker environment.
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
	timeout := g.chainConfig.DefaultWaitDuration()
	if g.StartupTimeout < timeout {
		timeout = g.StartupTimeout
	}

	return &tc.ContainerRequest{
		Name:          g.ContainerName,
		Image:         g.GetImageWithVersion(),
		ImagePlatform: "linux/amd64",
		Networks:      networks,
		WaitingFor: tcwait.ForAll(
			tcwait.ForLog("Starting beacon node").
				WithStartupTimeout(timeout).
				WithPollInterval(2 * time.Second),
		),
		Cmd: []string{
			"--accept-terms-of-use",
			"--datadir=/consensus-data",
			fmt.Sprintf("--chain-config-file=%s/config.yaml", g.generatedDataContainerDir),
			fmt.Sprintf("--genesis-state=%s/genesis.ssz", g.generatedDataContainerDir),
			fmt.Sprintf("--execution-endpoint=%s", g.gethInternalExecutionURL),
			"--rpc-host=0.0.0.0",
			"--grpc-gateway-host=0.0.0.0",
			"--grpc-gateway-corsdomain=*",
			"--suggested-fee-recipient=0x8943545177806ED17B9F23F0a21ee5948eCaa776",
			"--subscribe-all-subnets=true",
			fmt.Sprintf("--jwt-secret=%s", getJWTSecretFileLocationInsideContainer(g.generatedDataContainerDir)),
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

type PrysmValidator struct {
	EnvComponent
	chainConfig               *config.EthereumChainConfig
	internalBeaconRpcProvider string
	valKeysDir                string
	l                         zerolog.Logger
	t                         *testing.T
	posContainerSettings
}

// NewPrysmValidator initializes a new Prysm validator instance with the specified network configurations and settings.
// It is used to set up a validator for Ethereum's consensus layer, ensuring proper integration with the blockchain environment.
func NewPrysmValidator(networks []string, chainConfig *config.EthereumChainConfig, generatedDataHostDir, generatedDataContainerDir, valKeysDir, internalBeaconRpcProvider string, baseEthereumFork ethereum.Fork, opts ...EnvComponentOption) (*PrysmValidator, error) {
	pyrsmValidatorImage, ok := validatorForkToImageMap[baseEthereumFork]
	if !ok {
		return nil, fmt.Errorf("unknown fork: %s", baseEthereumFork)
	}
	parts := strings.Split(pyrsmValidatorImage, ":")
	g := &PrysmValidator{
		EnvComponent: EnvComponent{
			ContainerName:    fmt.Sprintf("%s-%s", "prysm-validator", uuid.NewString()[0:8]),
			Networks:         networks,
			ContainerImage:   parts[0],
			ContainerVersion: parts[1],
		},
		chainConfig:               chainConfig,
		posContainerSettings:      posContainerSettings{generatedDataHostDir: generatedDataHostDir, generatedDataContainerDir: generatedDataContainerDir},
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

// WithTestInstance sets up the PrysmValidator with a test logger and the provided testing context.
// This allows for easier testing and debugging of the validator's behavior during unit tests.
func (g *PrysmValidator) WithTestInstance(t *testing.T) *PrysmValidator {
	g.l = logging.GetTestLogger(t)
	g.t = t
	return g
}

// StartContainer initializes and starts the Prysm validator container.
// It handles the setup and logging, ensuring the container is ready for use in Ethereum network operations.
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
				WithStartupTimeout(g.chainConfig.DefaultWaitDuration()).
				WithPollInterval(2 * time.Second),
		),
		Cmd: []string{
			"--accept-terms-of-use",
			fmt.Sprintf("--chain-config-file=%s/config.yaml", g.generatedDataContainerDir),
			fmt.Sprintf("--beacon-rpc-provider=%s", g.internalBeaconRpcProvider),
			"--datadir=/consensus-data",
			"--suggested-fee-recipient=0x8943545177806ED17B9F23F0a21ee5948eCaa776",
			fmt.Sprintf("--wallet-dir=%s/prysm", NODE_0_DIR_INSIDE_CONTAINER),
			fmt.Sprintf("--wallet-password-file=%s", getValidatorWalletPasswordFileInsideContainer(g.generatedDataContainerDir)),
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
