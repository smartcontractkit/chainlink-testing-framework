package eth2

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
	"github.com/smartcontractkit/chainlink-testing-framework/docker"
	te "github.com/smartcontractkit/chainlink-testing-framework/docker/test_env"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"
)

const (
	ETH2_CONSENSUS_DIRECTORY = "/consensus"
	ETH2_EXECUTION_DIRECTORY = "/execution"
	PRYSM_RPC_PORT           = "4000"
	PRYSM_IMAGE_TAG          = "v4.1.1"
)

type BeaconChainGenesis struct {
	te.EnvComponent
	ExecutionDir string
	ConsensusDir string
	l            zerolog.Logger
	t            *testing.T
}

type BeaconChain struct {
	te.EnvComponent
	InternalRpcURL   string
	ExecutionDir     string
	ConsensusDir     string
	GethExecutionURL string
	l                zerolog.Logger
	t                *testing.T
}

type Validator struct {
	te.EnvComponent
	InternalBeaconRpcProvider string
	ConsensusDir              string
	l                         zerolog.Logger
	t                         *testing.T
}

func NewEth2Genesis(networks []string, opts ...te.EnvComponentOption) *BeaconChainGenesis {
	g := &BeaconChainGenesis{
		EnvComponent: te.EnvComponent{
			ContainerName: fmt.Sprintf("%s-%s", "beacon-chain-genesis", uuid.NewString()[0:8]),
			Networks:      networks,
		},
		l: log.Logger,
	}
	for _, opt := range opts {
		opt(&g.EnvComponent)
	}
	return g
}

func (g *BeaconChainGenesis) WithTestLogger(t *testing.T) *BeaconChainGenesis {
	g.l = logging.GetTestLogger(t)
	g.t = t
	return g
}

func (g *BeaconChainGenesis) StartContainer() error {
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
		return errors.Wrapf(err, "cannot start beacon chain genesis container")
	}

	g.l.Info().Str("containerName", g.ContainerName).
		Msg("Started Beacon Chain Genesis container")

	return nil
}

func (g *BeaconChainGenesis) getContainerRequest(networks []string) (*tc.ContainerRequest, error) {
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

	_, err = configFile.WriteString(beaconConfigYAML)
	if err != nil {
		return nil, err
	}

	genesisFile, err := os.CreateTemp("", "genesis_json")
	if err != nil {
		return nil, err
	}
	_, err = genesisFile.WriteString(genesisJSON)
	if err != nil {
		return nil, err
	}

	g.ExecutionDir = executionDir
	g.ConsensusDir = consensusDir

	return &tc.ContainerRequest{
		Name:            g.ContainerName,
		AlwaysPullImage: true,
		Image:           "gcr.io/prysmaticlabs/prysm/cmd/prysmctl:local-devnet",
		Networks:        networks,
		WaitingFor: tcwait.ForAll(
			tcwait.ForLog("Done writing genesis state to"),
			tcwait.ForLog("Command completed").
				WithStartupTimeout(120*time.Second).
				WithPollInterval(1*time.Second),
		),
		Cmd: []string{"testnet",
			"generate-genesis",
			"--fork=capella",
			"--num-validators=64",
			"--genesis-time-delay=15", //TODO: replace also here
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
				Target: ETH2_EXECUTION_DIRECTORY,
			},
			tc.ContainerMount{
				Source: tc.GenericBindMountSource{
					HostPath: consensusDir,
				},
				Target: ETH2_CONSENSUS_DIRECTORY,
			},
		},
	}, nil
}

func NewBeaconChain(networks []string, executionDir, consensusDir, gethExecutionURL string, opts ...te.EnvComponentOption) *BeaconChain {
	g := &BeaconChain{
		EnvComponent: te.EnvComponent{
			ContainerName: fmt.Sprintf("%s-%s", "prysm-beacon-chain", uuid.NewString()[0:8]),
			Networks:      networks,
		},
		ExecutionDir:     executionDir,
		ConsensusDir:     consensusDir,
		GethExecutionURL: gethExecutionURL,
		l:                log.Logger,
	}
	for _, opt := range opts {
		opt(&g.EnvComponent)
	}
	return g
}

func (g *BeaconChain) WithTestLogger(t *testing.T) *BeaconChain {
	g.l = logging.GetTestLogger(t)
	g.t = t
	return g
}

func (g *BeaconChain) StartContainer() error {
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
		return errors.Wrapf(err, "cannot start beacon chain container")
	}

	//TODO is this even needed?
	_, err = te.GetHost(context.Background(), ct)
	if err != nil {
		return err
	}

	_, err = ct.MappedPort(context.Background(), te.NatPort("3500"))
	if err != nil {
		return err
	}

	_, err = ct.MappedPort(context.Background(), te.NatPort("8080"))
	if err != nil {
		return err
	}

	_, err = ct.MappedPort(context.Background(), te.NatPort("6060"))
	if err != nil {
		return err
	}

	_, err = ct.MappedPort(context.Background(), te.NatPort("9090"))
	if err != nil {
		return err
	}

	externalRcpPort, err := ct.MappedPort(context.Background(), te.NatPort("4000"))
	if err != nil {
		return err
	}

	_ = externalRcpPort

	g.Container = ct
	g.InternalRpcURL = fmt.Sprintf("%s:%s", g.ContainerName, "4000")

	g.l.Info().Str("containerName", g.ContainerName).
		Msg("Started Prysm Beacon Chain container")

	return nil
}

func (g *BeaconChain) getContainerRequest(networks []string) (*tc.ContainerRequest, error) {
	// jwtSecret, err := os.CreateTemp(g.ExecutionDir, "jwtsecret")
	// if err != nil {
	// 	return nil, err
	// }
	// _, err = jwtSecret.WriteString("0xfad2709d0bb03bf0e8ba3c99bea194575d3e98863133d1af638ed056d1d59345")
	// if err != nil {
	// 	return nil, err
	// }

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
			fmt.Sprintf("--execution-endpoint=%s", g.GethExecutionURL),
			"--accept-terms-of-use",
			"--jwt-secret=/execution/jwtsecret",
			"--suggested-fee-recipient=0x123463a4b065722e99115d6c222f267d9cabb524",
			"--minimum-peers-per-subnet=0",
			"--enable-debug-rpc-endpoints",
			// "--interop-eth1data-votesgeth", //no idea why this flag results in error when passed here
		},
		ExposedPorts: []string{te.NatPortFormat(PRYSM_RPC_PORT), te.NatPortFormat("3500"), te.NatPortFormat("8080"), te.NatPortFormat("6060"), te.NatPortFormat("9090")},
		Mounts: tc.ContainerMounts{
			tc.ContainerMount{
				Source: tc.GenericBindMountSource{
					HostPath: g.ExecutionDir,
				},
				Target: ETH2_EXECUTION_DIRECTORY,
			},
			tc.ContainerMount{
				Source: tc.GenericBindMountSource{
					HostPath: g.ConsensusDir,
				},
				Target: ETH2_CONSENSUS_DIRECTORY,
			},
		},
	}, nil
}
