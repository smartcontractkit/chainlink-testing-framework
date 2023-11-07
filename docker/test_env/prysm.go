package test_env

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/docker"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
)

const (
	PRYSM_QUERY_RPC_PORT = "3500"
	PRYSM_NODE_RPC_PORT  = "4000"
	PRYSM_IMAGE_TAG      = "v4.1.1"
)

type PrysmGenesis struct {
	EnvComponent
	hostExecutionDir  string
	hostConsensusDir  string
	beaconChainConfig BeaconChainConfig
	l                 zerolog.Logger
	t                 *testing.T
}

type PrysmBeaconChain struct {
	EnvComponent
	InternalBeaconRpcProvider string
	InternalQueryRpcUrl       string
	ExternalBeaconRpcProvider string
	ExternalQueryRpcUrl       string
	hostExecutionDir          string
	hostConsensusDir          string
	gethInternalExecutionURL  string
	l                         zerolog.Logger
	t                         *testing.T
}

type PrysmValidator struct {
	EnvComponent
	internalBeaconRpcProvider string
	hostConsensusDir          string
	l                         zerolog.Logger
	t                         *testing.T
}

func NewEth2Genesis(networks []string, beaconChainConfig BeaconChainConfig, opts ...EnvComponentOption) *PrysmGenesis {
	g := &PrysmGenesis{
		EnvComponent: EnvComponent{
			ContainerName: fmt.Sprintf("%s-%s", "prysm-eth2-genesis", uuid.NewString()[0:8]),
			Networks:      networks,
		},
		beaconChainConfig: beaconChainConfig,
		l:                 log.Logger,
	}
	for _, opt := range opts {
		opt(&g.EnvComponent)
	}
	return g
}

func (g *PrysmGenesis) WithTestLogger(t *testing.T) *PrysmGenesis {
	g.l = logging.GetTestLogger(t)
	g.t = t
	return g
}

func (g *PrysmGenesis) StartContainer() error {
	r, err := g.getContainerRequest(g.Networks)
	if err != nil {
		return err
	}

	l := tc.Logger
	if g.t != nil {
		l = logging.CustomT{
			T: g.t,
			L: g.l,
		}
	}

	_, err = docker.StartContainerWithRetry(g.l, tc.GenericContainerRequest{
		ContainerRequest: *r,
		Reuse:            true,
		Started:          true,
		Logger:           l,
	})
	if err != nil {
		return errors.Wrapf(err, "cannot start prysm beacon chain genesis container")
	}

	g.l.Info().Str("containerName", g.ContainerName).
		Msg("Started Prysm Beacon Chain Genesis container")

	return nil
}

func (g *PrysmGenesis) getContainerRequest(networks []string) (*tc.ContainerRequest, error) {
	executionDir, err := os.MkdirTemp("", "execution")
	if err != nil {
		return nil, err
	}

	consensusDir, err := os.MkdirTemp("", "consensus")
	if err != nil {
		return nil, err
	}

	configFile, err := os.CreateTemp("", "config.yml")
	if err != nil {
		return nil, err
	}

	bc, err := GenerateBeaconChainConfig(&g.beaconChainConfig)
	if err != nil {
		return nil, err
	}
	_, err = configFile.WriteString(bc)
	if err != nil {
		return nil, err
	}

	genesisFile, err := os.CreateTemp("", "genesis_json")
	if err != nil {
		return nil, err
	}
	_, err = genesisFile.WriteString(Eth1GenesisJSON)
	if err != nil {
		return nil, err
	}

	g.hostExecutionDir = executionDir
	g.hostConsensusDir = consensusDir

	return &tc.ContainerRequest{
		Name:            g.ContainerName,
		AlwaysPullImage: true,
		Image:           "gcr.io/prysmaticlabs/prysm/cmd/prysmctl:local-devnet",
		Networks:        networks,
		WaitingFor: tcwait.ForAll(
			tcwait.ForLog("Done writing genesis state to"),
			tcwait.ForLog("Command completed").
				WithStartupTimeout(20*time.Second).
				WithPollInterval(1*time.Second),
		),
		Cmd: []string{"testnet",
			"generate-genesis",
			"--fork=capella",
			"--num-validators=64",
			"--genesis-time-delay=15",
			"--output-ssz=/consensus/genesis.ssz",
			"--chain-config-file=/consensus/config.yml",
			"--geth-genesis-json-in=/execution/genesis.json",
			"--geth-genesis-json-out=/execution/genesis.json",
		},
		Files: []tc.ContainerFile{
			{
				HostFilePath:      configFile.Name(),
				ContainerFilePath: "/consensus/config.yml",
				FileMode:          0644,
			},
			{
				HostFilePath:      genesisFile.Name(),
				ContainerFilePath: "/execution/genesis.json",
				FileMode:          0644,
			},
		},
		Mounts: tc.ContainerMounts{
			tc.ContainerMount{
				Source: tc.GenericBindMountSource{
					HostPath: executionDir,
				},
				Target: CONTAINER_ETH2_EXECUTION_DIRECTORY,
			},
			tc.ContainerMount{
				Source: tc.GenericBindMountSource{
					HostPath: consensusDir,
				},
				Target: CONTAINER_ETH2_CONSENSUS_DIRECTORY,
			},
		},
	}, nil
}

func NewPrysmBeaconChain(networks []string, executionDir, consensusDir, gethExecutionURL string, opts ...EnvComponentOption) *PrysmBeaconChain {
	g := &PrysmBeaconChain{
		EnvComponent: EnvComponent{
			ContainerName: fmt.Sprintf("%s-%s", "prysm-beacon-chain", uuid.NewString()[0:8]),
			Networks:      networks,
		},
		hostExecutionDir:         executionDir,
		hostConsensusDir:         consensusDir,
		gethInternalExecutionURL: gethExecutionURL,
		l:                        log.Logger,
	}
	for _, opt := range opts {
		opt(&g.EnvComponent)
	}
	return g
}

func (g *PrysmBeaconChain) WithTestLogger(t *testing.T) *PrysmBeaconChain {
	g.l = logging.GetTestLogger(t)
	g.t = t
	return g
}

func (g *PrysmBeaconChain) StartContainer() error {
	r, err := g.getContainerRequest(g.Networks)
	if err != nil {
		return err
	}

	l := tc.Logger
	if g.t != nil {
		l = logging.CustomT{
			T: g.t,
			L: g.l,
		}
	}
	ct, err := docker.StartContainerWithRetry(g.l, tc.GenericContainerRequest{
		ContainerRequest: *r,
		Reuse:            true,
		Started:          true,
		Logger:           l,
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
	g.ExternalBeaconRpcProvider = fmt.Sprintf("http://%s:%s", host, externalRcpPort.Port())
	g.ExternalQueryRpcUrl = fmt.Sprintf("http://%s:%s", host, queryPort.Port())

	g.l.Info().Str("containerName", g.ContainerName).
		Msg("Started Prysm Beacon Chain container")

	return nil
}

func (g *PrysmBeaconChain) getContainerRequest(networks []string) (*tc.ContainerRequest, error) {
	return &tc.ContainerRequest{
		Name:            g.ContainerName,
		AlwaysPullImage: true,
		Image:           fmt.Sprintf("gcr.io/prysmaticlabs/prysm/beacon-chain:%s", PRYSM_IMAGE_TAG),
		ImagePlatform:   "linux/amd64",
		Networks:        networks,
		WaitingFor: tcwait.ForAll(
			tcwait.ForLog("Received state initialized event"),
			tcwait.ForLog("Node started p2p server").
				WithStartupTimeout(120*time.Second).
				WithPollInterval(1*time.Second),
		),
		Cmd: []string{
			"--datadir=/consensus/beacondata",
			"--min-sync-peers=0",
			"--genesis-state=/consensus/genesis.ssz",
			"--bootstrap-node=",
			"--chain-config-file=/consensus/config.yml",
			"--contract-deployment-block=0",
			"--chain-id=1337",
			"--rpc-host=0.0.0.0",
			"--grpc-gateway-host=0.0.0.0",
			fmt.Sprintf("--execution-endpoint=%s", g.gethInternalExecutionURL),
			"--accept-terms-of-use",
			"--jwt-secret=/execution/jwtsecret",
			"--suggested-fee-recipient=0x123463a4b065722e99115d6c222f267d9cabb524",
			"--minimum-peers-per-subnet=0",
			"--enable-debug-rpc-endpoints",
			// "--interop-eth1data-votesgeth", //no idea why this flag results in error when passed here
		},
		ExposedPorts: []string{NatPortFormat(PRYSM_NODE_RPC_PORT), NatPortFormat(PRYSM_QUERY_RPC_PORT)},
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
		},
	}, nil
}

func NewPrysmValidator(networks []string, consensusDir, internalBeaconRpcProvider string, opts ...EnvComponentOption) *PrysmValidator {
	g := &PrysmValidator{
		EnvComponent: EnvComponent{
			ContainerName: fmt.Sprintf("%s-%s", "prysm-validator", uuid.NewString()[0:8]),
			Networks:      networks,
		},
		hostConsensusDir:          consensusDir,
		internalBeaconRpcProvider: internalBeaconRpcProvider,
		l:                         log.Logger,
	}
	for _, opt := range opts {
		opt(&g.EnvComponent)
	}
	return g
}

func (g *PrysmValidator) WithTestLogger(t *testing.T) *PrysmValidator {
	g.l = logging.GetTestLogger(t)
	g.t = t
	return g
}

func (g *PrysmValidator) StartContainer() error {
	r, err := g.getContainerRequest(g.Networks)
	if err != nil {
		return err
	}

	l := tc.Logger
	if g.t != nil {
		l = logging.CustomT{
			T: g.t,
			L: g.l,
		}
	}
	ct, err := docker.StartContainerWithRetry(g.l, tc.GenericContainerRequest{
		ContainerRequest: *r,
		Reuse:            true,
		Started:          true,
		Logger:           l,
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
	return &tc.ContainerRequest{
		Name:            g.ContainerName,
		AlwaysPullImage: true,
		Image:           fmt.Sprintf("gcr.io/prysmaticlabs/prysm/validator:%s", PRYSM_IMAGE_TAG),
		Networks:        networks,
		WaitingFor: tcwait.ForAll(
			tcwait.ForLog("Beacon chain started").
				WithStartupTimeout(120 * time.Second).
				WithPollInterval(1 * time.Second),
		),
		Cmd: []string{fmt.Sprintf("--beacon-rpc-provider=%s", g.internalBeaconRpcProvider),
			"--datadir=/consensus/validatordata",
			"--accept-terms-of-use",
			"--interop-num-validators=64",
			"--interop-start-index=0",
			"--chain-config-file=/consensus/config.yml",
		},
		Mounts: tc.ContainerMounts{
			tc.ContainerMount{
				Source: tc.GenericBindMountSource{
					HostPath: g.hostConsensusDir,
				},
				Target: CONTAINER_ETH2_CONSENSUS_DIRECTORY,
			},
		},
	}, nil
}
