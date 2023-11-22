package test_env

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/docker"
)

const (
	PRYSM_QUERY_RPC_PORT = "3500"
	PRYSM_NODE_RPC_PORT  = "4000"
	PRYSM_IMAGE_TAG      = "v4.1.1-debug"
)

type PrysmBeaconChain struct {
	EnvComponent
	InternalBeaconRpcProvider string
	InternalQueryRpcUrl       string
	ExternalBeaconRpcProvider string
	ExternalQueryRpcUrl       string
	hostExecutionDir          string
	hostConsensusDir          string
	customConfigDataDir       string
	gethInternalExecutionURL  string
	beaconChainConfig         BeaconChainConfig
	l                         zerolog.Logger
}

type PrysmValidator struct {
	EnvComponent
	internalBeaconRpcProvider string
	hostConsensusDir          string
	valKeysDir                string
	customConfigDataDir       string
	beaconChainConfig         BeaconChainConfig
	l                         zerolog.Logger
}

func NewPrysmBeaconChain(networks []string, beaconChainConfig BeaconChainConfig, executionDir, consensusDir, customConfigDataDir, gethExecutionURL string, opts ...EnvComponentOption) *PrysmBeaconChain {
	g := &PrysmBeaconChain{
		EnvComponent: EnvComponent{
			ContainerName: fmt.Sprintf("%s-%s", "prysm-beacon-chain", uuid.NewString()[0:8]),
			Networks:      networks,
		},
		beaconChainConfig:        beaconChainConfig,
		hostExecutionDir:         executionDir,
		hostConsensusDir:         consensusDir,
		customConfigDataDir:      customConfigDataDir,
		gethInternalExecutionURL: gethExecutionURL,
		l:                        log.Logger,
	}
	for _, opt := range opts {
		opt(&g.EnvComponent)
	}
	return g
}

func (g *PrysmBeaconChain) WithLogger(l zerolog.Logger) *PrysmBeaconChain {
	g.l = l
	return g
}

func (g *PrysmBeaconChain) StartContainer() error {
	r, err := g.getContainerRequest(g.Networks)
	if err != nil {
		return err
	}

	ct, err := docker.StartContainerWithRetry(g.l, tc.GenericContainerRequest{
		ContainerRequest: *r,
		Reuse:            true,
		Started:          true,
		Logger:           &g.l,
	})
	if err != nil {
		return errors.Wrapf(err, "cannot start prysm beacon chain container")
	}

	host, err := GetHost(context.Background(), ct)
	if err != nil {
		return err
	}
	queryPort, err := ct.MappedPort(context.Background(), NatPort(PRYSM_QUERY_RPC_PORT))
	if err != nil {
		return err
	}

	externalRcpPort, err := ct.MappedPort(context.Background(), NatPort(PRYSM_NODE_RPC_PORT))
	if err != nil {
		return err
	}

	_ = externalRcpPort

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
	jwtSecretFile, err := os.CreateTemp("/tmp", "jwtsecret")
	if err != nil {
		return nil, err
	}
	_, err = jwtSecretFile.WriteString("0xfad2709d0bb03bf0e8ba3c99bea194575d3e98863133d1af638ed056d1d59345")
	if err != nil {
		return nil, err
	}

	waitDuration := time.Duration(g.beaconChainConfig.GenesisDelay+g.beaconChainConfig.GetValidatorBasedGenesisDelay()) * 2

	return &tc.ContainerRequest{
		Name: g.ContainerName,
		// AlwaysPullImage: true,
		Image:         fmt.Sprintf("gcr.io/prysmaticlabs/prysm/beacon-chain:%s", PRYSM_IMAGE_TAG),
		ImagePlatform: "linux/amd64",
		Networks:      networks,
		WaitingFor: tcwait.ForAll(
			tcwait.ForLog("Starting beacon node").
				WithStartupTimeout(waitDuration * time.Second).
				WithPollInterval(2 * time.Second),
		),
		Cmd: []string{
			"--accept-terms-of-use",
			"--datadir=/consensus-data",
			"--chain-config-file=/data/custom_config_data/config.yaml",
			"--genesis-state=/data/custom_config_data/genesis.ssz",
			fmt.Sprintf("--execution-endpoint=%s", g.gethInternalExecutionURL),
			"--rpc-host=0.0.0.0",
			"--grpc-gateway-host=0.0.0.0",
			"--grpc-gateway-corsdomain=*",
			"--suggested-fee-recipient=0x8943545177806ED17B9F23F0a21ee5948eCaa776",
			"--subscribe-all-subnets=true",
			"--jwt-secret=/data/jwtsecret",
			// mine
			"--minimum-peers-per-subnet=0",
			"--min-sync-peers=0",
			// unused
			// "--chain-id=1337",
			// "--bootstrap-node=",
			// "--contract-deployment-block=0",
			// "--enable-debug-rpc-endpoints",
			// "--verbosity=debug",
			// "--interop-eth1data-votes",
		},
		ExposedPorts: []string{NatPortFormat(PRYSM_NODE_RPC_PORT), NatPortFormat(PRYSM_QUERY_RPC_PORT)},
		Files: []tc.ContainerFile{
			{
				HostFilePath:      jwtSecretFile.Name(),
				ContainerFilePath: "/data/jwtsecret",
				FileMode:          0644,
			},
		},
		Mounts: tc.ContainerMounts{
			tc.ContainerMount{
				Source: tc.GenericBindMountSource{
					HostPath: g.hostExecutionDir,
				},
				Target: CONTAINER_ETH2_EXECUTION_DIRECTORY,
			},
			tc.ContainerMount{
				Source: tc.GenericBindMountSource{
					HostPath: g.hostConsensusDir,
				},
				Target: CONTAINER_ETH2_CONSENSUS_DIRECTORY,
			},
			tc.ContainerMount{
				Source: tc.GenericBindMountSource{
					HostPath: g.customConfigDataDir,
				},
				Target: "/data/custom_config_data",
			},
		},
	}, nil
}

func NewPrysmValidator(networks []string, beaconChainConfig BeaconChainConfig, consensusDir, customConfigDataDir, valKeysDir, internalBeaconRpcProvider string, opts ...EnvComponentOption) *PrysmValidator {
	g := &PrysmValidator{
		EnvComponent: EnvComponent{
			ContainerName: fmt.Sprintf("%s-%s", "prysm-validator", uuid.NewString()[0:8]),
			Networks:      networks,
		},
		beaconChainConfig:         beaconChainConfig,
		hostConsensusDir:          consensusDir,
		customConfigDataDir:       customConfigDataDir,
		valKeysDir:                valKeysDir,
		internalBeaconRpcProvider: internalBeaconRpcProvider,
		l:                         log.Logger,
	}
	for _, opt := range opts {
		opt(&g.EnvComponent)
	}
	return g
}

func (g *PrysmValidator) WithLogger(l zerolog.Logger) *PrysmValidator {
	g.l = l
	return g
}

func (g *PrysmValidator) StartContainer() error {
	r, err := g.getContainerRequest(g.Networks)
	if err != nil {
		return err
	}

	ct, err := docker.StartContainerWithRetry(g.l, tc.GenericContainerRequest{
		ContainerRequest: *r,
		Reuse:            true,
		Started:          true,
		Logger:           &g.l,
	})
	if err != nil {
		return errors.Wrapf(err, "cannot start prysm validator container")
	}

	g.Container = ct

	g.l.Info().Str("containerName", g.ContainerName).
		Msg("Started Prysm Validator container")

	return nil
}

func (g *PrysmValidator) getContainerRequest(networks []string) (*tc.ContainerRequest, error) {
	passwordFile, err := os.CreateTemp("", "password.txt")
	if err != nil {
		return nil, err
	}
	_, err = passwordFile.WriteString("password")
	if err != nil {
		return nil, err
	}

	waitDuration := time.Duration(g.beaconChainConfig.GenesisDelay+g.beaconChainConfig.GetValidatorBasedGenesisDelay()) * 2

	return &tc.ContainerRequest{
		Name: g.ContainerName,
		// AlwaysPullImage: true,
		Image:    fmt.Sprintf("gcr.io/prysmaticlabs/prysm/validator:%s", PRYSM_IMAGE_TAG),
		Networks: networks,
		WaitingFor: tcwait.ForAll(
			tcwait.ForLog("Beacon chain started").
				WithStartupTimeout(waitDuration * time.Second).
				WithPollInterval(2 * time.Second),
		),
		Cmd: []string{
			"--accept-terms-of-use",
			"--chain-config-file=/data/custom_config_data/config.yaml",
			fmt.Sprintf("--beacon-rpc-provider=%s", g.internalBeaconRpcProvider),
			"--datadir=/consensus-data",
			"--suggested-fee-recipient=0x8943545177806ED17B9F23F0a21ee5948eCaa776",
			"--wallet-dir=/keys/node-0/prysm",
			"--wallet-password-file=/keys/password.txt",
			// "--interop-num-validators=8",
			// "--interop-start-index=0",
			// "--verbosity=debug",
		},
		Files: []tc.ContainerFile{
			{
				HostFilePath:      passwordFile.Name(),
				ContainerFilePath: "/keys/password.txt",
				FileMode:          0644,
			},
		},
		Mounts: tc.ContainerMounts{
			tc.ContainerMount{
				Source: tc.GenericBindMountSource{
					HostPath: g.valKeysDir,
				},
				Target: "/keys",
			},
			tc.ContainerMount{
				Source: tc.GenericBindMountSource{
					HostPath: g.customConfigDataDir,
				},
				Target: "/data/custom_config_data",
			},
		},
	}, nil
}
