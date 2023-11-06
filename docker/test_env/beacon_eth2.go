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
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/docker"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"
)

//TODO expose simple API that will start all these containers and return only data we care about
//TODO fund our addresses

const (
	ETH2_CONSENSUS_DIRECTORY = "/consensus"
	ETH2_EXECUTION_DIRECTORY = "/execution"
	GO_CLIENT_IMAGE          = "ethereum/client-go:latest" //TODO: fix version
	BEACON_RPC_PORT          = "4000"
	GETH_EXECUTION_PORT      = "8511"
)

type BeaconChainGenesis struct {
	EnvComponent
	ExecutionDir string
	ConsensusDir string
	l            zerolog.Logger
	t            *testing.T
}

type GethGenesis struct {
	EnvComponent
	ExecutionDir string
	l            zerolog.Logger
	t            *testing.T
}

type BeaconChain struct {
	EnvComponent
	InternalRpcURL   string
	ExecutionDir     string
	ConsensusDir     string
	GethExecutionURL string
	l                zerolog.Logger
	t                *testing.T
}

type Validator struct {
	EnvComponent
	InternalBeaconRpcProvider string
	ConsensusDir              string
	l                         zerolog.Logger
	t                         *testing.T
}

type Geth2 struct {
	EnvComponent
	ExternalHttpUrl string
	InternalHttpUrl string
	ExternalWsUrl   string
	InternalWsUrl   string
	ExecutionURL    string
	ExecutionDir    string
	l               zerolog.Logger
	t               *testing.T
}

func NewBeaconChainGenesis(networks []string, opts ...EnvComponentOption) *BeaconChainGenesis {
	g := &BeaconChainGenesis{
		EnvComponent: EnvComponent{
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

func NewGethGenesis(networks []string, executionDir string, opts ...EnvComponentOption) *GethGenesis {
	g := &GethGenesis{
		EnvComponent: EnvComponent{
			ContainerName: fmt.Sprintf("%s-%s", "geth-genesis", uuid.NewString()[0:8]),
			Networks:      networks,
		},
		ExecutionDir: executionDir,
		l:            log.Logger,
	}
	for _, opt := range opts {
		opt(&g.EnvComponent)
	}
	return g
}

func (g *GethGenesis) WithTestLogger(t *testing.T) *GethGenesis {
	g.l = logging.GetTestLogger(t)
	g.t = t
	return g
}

func (g *GethGenesis) StartContainer() error {
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
		return errors.Wrapf(err, "cannot start geth container")
	}

	g.l.Info().Str("containerName", g.ContainerName).
		Msg("Started Geth Genesis container")

	return nil
}

func (g *GethGenesis) getContainerRequest(networks []string) (*tc.ContainerRequest, error) {
	return &tc.ContainerRequest{
		Name:            g.ContainerName,
		AlwaysPullImage: true,
		Image:           GO_CLIENT_IMAGE,
		Networks:        networks,
		WaitingFor: tcwait.ForAll(
			tcwait.ForLog("Successfully wrote genesis state").
				WithStartupTimeout(120 * time.Second).
				WithPollInterval(1 * time.Second),
		),
		Cmd: []string{"--datadir=/execution",
			"init",
			"/execution/genesis.json",
		},
		Mounts: tc.ContainerMounts{
			tc.ContainerMount{
				Source: tc.GenericBindMountSource{
					HostPath: g.ExecutionDir,
				},
				Target: ETH2_EXECUTION_DIRECTORY,
			},
		},
	}, nil
}

func NewBeaconChain(networks []string, executionDir, consensusDir, gethExecutionURL string, opts ...EnvComponentOption) *BeaconChain {
	g := &BeaconChain{
		EnvComponent: EnvComponent{
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
	_, err = GetHost(context.Background(), ct)
	if err != nil {
		return err
	}

	_, err = ct.MappedPort(context.Background(), NatPort("3500"))
	if err != nil {
		return err
	}

	_, err = ct.MappedPort(context.Background(), NatPort("8080"))
	if err != nil {
		return err
	}

	_, err = ct.MappedPort(context.Background(), NatPort("6060"))
	if err != nil {
		return err
	}

	_, err = ct.MappedPort(context.Background(), NatPort("9090"))
	if err != nil {
		return err
	}

	externalRcpPort, err := ct.MappedPort(context.Background(), NatPort("4000"))
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
		Image:           "gcr.io/prysmaticlabs/prysm/beacon-chain:v4.0.8",
		ImagePlatform:   "linux/amd64",
		Networks:        networks,
		WaitingFor: tcwait.ForAll(
			tcwait.ForLog("Received state initialized event"),
			tcwait.ForLog("Node started p2p server").
				// tcwait.ForLog("Chain genesis time reached").
				WithStartupTimeout(120*time.Second).
				WithPollInterval(1*time.Second),
		),
		//write a bash script to execute that command?
		//I have no idea why it's failing with level=error msg="flag provided but not defined: -interop-eth1data-votesgeth" prefix=main
		//if according to output it should be there:
		//           --interop-eth1data-votes                                Enable mocking of eth1 data votes for proposers to package into blocks (default: false)
		Cmd: []string{
			"--datadir=/consensus/beacondata",
			"--min-sync-peers=0",
			"--genesis-state=/consensus/genesis.ssz",
			"--bootstrap-node=",
			//TODO check if genesis file is there
			"--chain-config-file=/consensus/config.yml",
			"--contract-deployment-block=0",
			"--chain-id=32382", //TODO change me
			"--rpc-host=0.0.0.0",
			"--grpc-gateway-host=0.0.0.0",
			fmt.Sprintf("--execution-endpoint=%s", g.GethExecutionURL),
			"--accept-terms-of-use",
			"--jwt-secret=/execution/jwtsecret",
			"--suggested-fee-recipient=0x123463a4b065722e99115d6c222f267d9cabb524",
			"--minimum-peers-per-subnet=0",
			"--enable-debug-rpc-endpoints",
			"--interop-eth1data-votesgeth",
		},
		ExposedPorts: []string{NatPortFormat(BEACON_RPC_PORT), NatPortFormat("3500"), NatPortFormat("8080"), NatPortFormat("6060"), NatPortFormat("9090")},
		// Files: []tc.ContainerFile{
		// 	{
		// 		HostFilePath:      jwtSecret.Name(),
		// 		ContainerFilePath: "/execution/jwtsecret",
		// 		FileMode:          0644,
		// 	},
		// },
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

func NewGeth2(networks []string, executionDir string, opts ...EnvComponentOption) *Geth2 {
	g := &Geth2{
		EnvComponent: EnvComponent{
			ContainerName: fmt.Sprintf("%s-%s", "geth2", uuid.NewString()[0:8]),
			Networks:      networks,
		},
		ExecutionDir: executionDir,
		l:            log.Logger,
	}
	for _, opt := range opts {
		opt(&g.EnvComponent)
	}
	return g
}

func (g *Geth2) WithTestLogger(t *testing.T) *Geth2 {
	g.l = logging.GetTestLogger(t)
	g.t = t
	return g
}

func (g *Geth2) StartContainer() (blockchain.EVMNetwork, InternalDockerUrls, error) {
	r, err := g.getContainerRequest(g.Networks)
	if err != nil {
		return blockchain.EVMNetwork{}, InternalDockerUrls{}, err
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
		return blockchain.EVMNetwork{}, InternalDockerUrls{}, errors.Wrapf(err, "cannot start geth container")
	}

	host, err := GetHost(context.Background(), ct)
	if err != nil {
		return blockchain.EVMNetwork{}, InternalDockerUrls{}, err
	}
	// host, err := ct.ContainerIP(context.Background())
	// host, err := ct.Host(context.Background())
	if err != nil {
		return blockchain.EVMNetwork{}, InternalDockerUrls{}, err
	}
	httpPort, err := ct.MappedPort(context.Background(), NatPort(TX_GETH_HTTP_PORT))
	if err != nil {
		return blockchain.EVMNetwork{}, InternalDockerUrls{}, err
	}
	wsPort, err := ct.MappedPort(context.Background(), NatPort(TX_GETH_WS_PORT))
	if err != nil {
		return blockchain.EVMNetwork{}, InternalDockerUrls{}, err
	}
	_, err = ct.MappedPort(context.Background(), NatPort("8551"))
	if err != nil {
		return blockchain.EVMNetwork{}, InternalDockerUrls{}, err
	}

	g.Container = ct
	g.ExternalHttpUrl = fmt.Sprintf("http://%s:%s", host, httpPort.Port())
	g.InternalHttpUrl = fmt.Sprintf("http://%s:%s", g.ContainerName, TX_GETH_HTTP_PORT)
	g.ExternalWsUrl = fmt.Sprintf("ws://%s:%s", host, wsPort.Port())
	g.InternalWsUrl = fmt.Sprintf("ws://%s:%s", g.ContainerName, TX_GETH_WS_PORT)
	g.ExecutionURL = fmt.Sprintf("http://%s:%s", g.ContainerName, "8551")

	networkConfig := blockchain.SimulatedEVMNetwork
	networkConfig.Name = "geth"
	networkConfig.URLs = []string{g.ExternalWsUrl}
	networkConfig.HTTPURLs = []string{g.ExternalHttpUrl}

	internalDockerUrls := InternalDockerUrls{
		HttpUrl: g.InternalHttpUrl,
		WsUrl:   g.InternalWsUrl,
	}

	g.l.Info().Str("containerName", g.ContainerName).
		Msg("Started Geth2 container")

	return networkConfig, internalDockerUrls, nil
}

func (g *Geth2) getContainerRequest(networks []string) (*tc.ContainerRequest, error) {
	err := os.WriteFile(g.ExecutionDir+"/password.txt", []byte(""), 0600)
	if err != nil {
		return nil, err
	}

	key1File, err := os.CreateTemp(g.ExecutionDir+"/keystore", "UTC--2022-08-19T17-38-31.257380510Z--123463a4b065722e99115d6c222f267d9cabb524")
	if err != nil {
		return nil, err
	}
	_, err = key1File.WriteString(`{"address":"123463a4b065722e99115d6c222f267d9cabb524","crypto":{"cipher":"aes-128-ctr","ciphertext":"93b90389b855889b9f91c89fd15b9bd2ae95b06fe8e2314009fc88859fc6fde9","cipherparams":{"iv":"9dc2eff7967505f0e6a40264d1511742"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"c07503bb1b66083c37527cd8f06f8c7c1443d4c724767f625743bd47ae6179a4"},"mac":"6d359be5d6c432d5bbb859484009a4bf1bd71b76e89420c380bd0593ce25a817"},"id":"622df904-0bb1-4236-b254-f1b8dfdff1ec","version":3}`)
	if err != nil {
		return nil, err
	}

	jwtSecret, err := os.CreateTemp(g.ExecutionDir, "jwtsecret")
	if err != nil {
		return nil, err
	}
	_, err = jwtSecret.WriteString("0xfad2709d0bb03bf0e8ba3c99bea194575d3e98863133d1af638ed056d1d59345")
	if err != nil {
		return nil, err
	}
	secretKey, err := os.CreateTemp(g.ExecutionDir, "sk.json")
	if err != nil {
		return nil, err
	}
	_, err = secretKey.WriteString("2e0834786285daccd064ca17f1654f67b4aef298acbb82cef9ec422fb4975622")
	if err != nil {
		return nil, err
	}

	return &tc.ContainerRequest{
		Name:            g.ContainerName,
		AlwaysPullImage: true,
		Image:           GO_CLIENT_IMAGE,
		Networks:        networks,
		ExposedPorts:    []string{NatPortFormat(TX_GETH_HTTP_PORT), NatPortFormat(TX_GETH_WS_PORT), NatPortFormat("8551")},
		// ExposedPorts: []string{NatPortFormat(TX_GETH_HTTP_PORT), NatPortFormat(TX_GETH_WS_PORT), "8551"},
		WaitingFor: tcwait.ForAll(
			// NewHTTPStrategy("/", NatPort("8551")),
			NewHTTPStrategy("/", NatPort(TX_GETH_HTTP_PORT)),
			tcwait.ForLog("WebSocket enabled"),
			tcwait.ForLog("Started P2P networking").
				WithStartupTimeout(120*time.Second).
				WithPollInterval(1*time.Second),
			NewWebSocketStrategy(NatPort(TX_GETH_WS_PORT), g.l),
		),
		Cmd: []string{"--http",
			"--http.api=eth,net,web3",
			"--http.addr=0.0.0.0",
			"--http.corsdomain=*",
			fmt.Sprintf("--http.port=%s", TX_GETH_HTTP_PORT),
			"--ws",
			"--ws.api=eth,net,web3",
			"--ws.addr=0.0.0.0",
			"--ws.origins=*",
			fmt.Sprintf("--ws.port=%s", TX_GETH_WS_PORT),
			"--authrpc.vhosts=*",
			"--authrpc.addr=0.0.0.0",
			"--authrpc.jwtsecret=/execution/jwtsecret",
			"--datadir=/execution",
			"--allow-insecure-unlock",
			"--unlock=0x123463a4b065722e99115d6c222f267d9cabb524", //TODO update me with ours?
			"--password=/execution/password.txt",
			"--nodiscover",
			"--syncmode=full",
		},
		Files: []tc.ContainerFile{
			{
				HostFilePath:      jwtSecret.Name(),
				ContainerFilePath: "/execution/jwtsecret",
				FileMode:          0644,
			},
			{
				HostFilePath:      secretKey.Name(),
				ContainerFilePath: "/execution/sk.json",
				FileMode:          0644,
			},
		},
		Mounts: tc.ContainerMounts{
			tc.ContainerMount{
				Source: tc.GenericBindMountSource{
					HostPath: g.ExecutionDir,
				},
				Target: ETH2_EXECUTION_DIRECTORY,
			},
		},
	}, nil
}

func NewValidator(networks []string, consensusDir, internalBeaconRpcProvider string, opts ...EnvComponentOption) *Validator {
	g := &Validator{
		EnvComponent: EnvComponent{
			ContainerName: fmt.Sprintf("%s-%s", "prysm-validator", uuid.NewString()[0:8]),
			Networks:      networks,
		},
		ConsensusDir:              consensusDir,
		InternalBeaconRpcProvider: internalBeaconRpcProvider,
		l:                         log.Logger,
	}
	for _, opt := range opts {
		opt(&g.EnvComponent)
	}
	return g
}

func (g *Validator) WithTestLogger(t *testing.T) *Validator {
	g.l = logging.GetTestLogger(t)
	g.t = t
	return g
}

func (g *Validator) StartContainer() error {
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

func (g *Validator) getContainerRequest(networks []string) (*tc.ContainerRequest, error) {
	return &tc.ContainerRequest{
		Name:            g.ContainerName,
		AlwaysPullImage: true,
		Image:           "gcr.io/prysmaticlabs/prysm/validator:v4.0.8",
		Networks:        networks,
		WaitingFor: tcwait.ForAll(
			tcwait.ForLog("Beacon chain started").
				WithStartupTimeout(120 * time.Second).
				WithPollInterval(1 * time.Second),
		),
		Cmd: []string{fmt.Sprintf("--beacon-rpc-provider==%s", g.InternalBeaconRpcProvider),
			"--datadir=/consensus/validatordata",
			"--accept-terms-of-use",
			"--interop-num-validators=64",
			"--interop-start-index=0",
			"--chain-config-file=/consensus/config.yml",
		},
		Mounts: tc.ContainerMounts{
			tc.ContainerMount{
				Source: tc.GenericBindMountSource{
					HostPath: g.ConsensusDir,
				},
				Target: ETH2_CONSENSUS_DIRECTORY,
			},
		},
	}, nil
}

var beaconConfigYAML = `
CONFIG_NAME: interop
PRESET_BASE: interop

# Genesis
GENESIS_FORK_VERSION: 0x20000089

# Altair
ALTAIR_FORK_EPOCH: 0
ALTAIR_FORK_VERSION: 0x20000090

# Merge
BELLATRIX_FORK_EPOCH: 0
BELLATRIX_FORK_VERSION: 0x20000091
TERMINAL_TOTAL_DIFFICULTY: 0

# Capella
CAPELLA_FORK_EPOCH: 0
CAPELLA_FORK_VERSION: 0x20000092
MAX_WITHDRAWALS_PER_PAYLOAD: 16

DENEB_FORK_VERSION: 0x20000093

# Time parameters
SECONDS_PER_SLOT: 12
SLOTS_PER_EPOCH: 6

# Deposit contract
DEPOSIT_CONTRACT_ADDRESS: 0x4242424242424242424242424242424242424242
`

// TODO change chainID and founded addresses
// TODO what about shanghaiTime?
var genesisJSON = `
{
	"config": {
		"chainId": 32382,
		"homesteadBlock": 0,
		"daoForkSupport": true,
		"eip150Block": 0,
		"eip155Block": 0,
		"eip158Block": 0,
		"byzantiumBlock": 0,
		"constantinopleBlock": 0,
		"petersburgBlock": 0,
		"istanbulBlock": 0,
		"muirGlacierBlock": 0,
		"berlinBlock": 0,
		"londonBlock": 0,
		"arrowGlacierBlock": 0,
		"grayGlacierBlock": 0,
		"shanghaiTime": 1699271757,
		"terminalTotalDifficulty": 0,
		"terminalTotalDifficultyPassed": true
	},
	"nonce": "0x0",
	"timestamp": "0x6548d44d",
	"extraData": "0x0000000000000000000000000000000000000000000000000000000000000000123463a4b065722e99115d6c222f267d9cabb5240000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
	"gasLimit": "0x1c9c380",
	"difficulty": "0x1",
	"mixHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
	"coinbase": "0x0000000000000000000000000000000000000000",
	"alloc": {
		"123463a4b065722e99115d6c222f267d9cabb524": {
			"balance": "0x43c33c1937564800000"
		},
		"14dc79964da2c08b23698b3d3cc7ca32193d9955": {
			"balance": "0x21e19e0c9bab2400000"
		},
		"15d34aaf54267db7d7c367839aaf71a00a2c6a65": {
			"balance": "0x21e19e0c9bab2400000"
		},
		"1cbd3b2770909d4e10f157cabc84c7264073c9ec": {
			"balance": "0x21e19e0c9bab2400000"
		},
		"23618e81e3f5cdf7f54c3d65f7fbc0abf5b21e8f": {
			"balance": "0x21e19e0c9bab2400000"
		},
		"2546bcd3c84621e976d8185a91a922ae77ecec30": {
			"balance": "0x21e19e0c9bab2400000"
		},
		"3c44cdddb6a900fa2b585dd299e03d12fa4293bc": {
			"balance": "0x21e19e0c9bab2400000"
		},
		"4242424242424242424242424242424242424242": {
			"code": "0x60806040526004361061003f5760003560e01c806301ffc9a71461004457806322895118146100b6578063621fd130146101e3578063c5f2892f14610273575b600080fd5b34801561005057600080fd5b5061009c6004803603602081101561006757600080fd5b8101908080357bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916906020019092919050505061029e565b604051808215151515815260200191505060405180910390f35b6101e1600480360360808110156100cc57600080fd5b81019080803590602001906401000000008111156100e957600080fd5b8201836020820111156100fb57600080fd5b8035906020019184600183028401116401000000008311171561011d57600080fd5b90919293919293908035906020019064010000000081111561013e57600080fd5b82018360208201111561015057600080fd5b8035906020019184600183028401116401000000008311171561017257600080fd5b90919293919293908035906020019064010000000081111561019357600080fd5b8201836020820111156101a557600080fd5b803590602001918460018302840111640100000000831117156101c757600080fd5b909192939192939080359060200190929190505050610370565b005b3480156101ef57600080fd5b506101f8610fd0565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101561023857808201518184015260208101905061021d565b50505050905090810190601f1680156102655780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561027f57600080fd5b50610288610fe2565b6040518082815260200191505060405180910390f35b60007f01ffc9a7000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916148061036957507f85640907000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916145b9050919050565b603087879050146103cc576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260268152602001806116ec6026913960400191505060405180910390fd5b60208585905014610428576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260368152602001806116836036913960400191505060405180910390fd5b60608383905014610484576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602981526020018061175f6029913960400191505060405180910390fd5b670de0b6b3a76400003410156104e5576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260268152602001806117396026913960400191505060405180910390fd5b6000633b9aca0034816104f457fe5b061461054b576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260338152602001806116b96033913960400191505060405180910390fd5b6000633b9aca00348161055a57fe5b04905067ffffffffffffffff80168111156105c0576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260278152602001806117126027913960400191505060405180910390fd5b60606105cb82611314565b90507f649bbc62d0e31342afea4e5cd82d4049e7e1ee912fc0889aa790803be39038c589898989858a8a610600602054611314565b60405180806020018060200180602001806020018060200186810386528e8e82818152602001925080828437600081840152601f19601f82011690508083019250505086810385528c8c82818152602001925080828437600081840152601f19601f82011690508083019250505086810384528a818151815260200191508051906020019080838360005b838110156106a657808201518184015260208101905061068b565b50505050905090810190601f1680156106d35780820380516001836020036101000a031916815260200191505b508681038352898982818152602001925080828437600081840152601f19601f820116905080830192505050868103825287818151815260200191508051906020019080838360005b8381101561073757808201518184015260208101905061071c565b50505050905090810190601f1680156107645780820380516001836020036101000a031916815260200191505b509d505050505050505050505050505060405180910390a1600060028a8a600060801b6040516020018084848082843780830192505050826fffffffffffffffffffffffffffffffff19166fffffffffffffffffffffffffffffffff1916815260100193505050506040516020818303038152906040526040518082805190602001908083835b6020831061080e57805182526020820191506020810190506020830392506107eb565b6001836020036101000a038019825116818451168082178552505050505050905001915050602060405180830381855afa158015610850573d6000803e3d6000fd5b5050506040513d602081101561086557600080fd5b8101908080519060200190929190505050905060006002808888600090604092610891939291906115da565b6040516020018083838082843780830192505050925050506040516020818303038152906040526040518082805190602001908083835b602083106108eb57805182526020820191506020810190506020830392506108c8565b6001836020036101000a038019825116818451168082178552505050505050905001915050602060405180830381855afa15801561092d573d6000803e3d6000fd5b5050506040513d602081101561094257600080fd5b8101908080519060200190929190505050600289896040908092610968939291906115da565b6000801b604051602001808484808284378083019250505082815260200193505050506040516020818303038152906040526040518082805190602001908083835b602083106109cd57805182526020820191506020810190506020830392506109aa565b6001836020036101000a038019825116818451168082178552505050505050905001915050602060405180830381855afa158015610a0f573d6000803e3d6000fd5b5050506040513d6020811015610a2457600080fd5b810190808051906020019092919050505060405160200180838152602001828152602001925050506040516020818303038152906040526040518082805190602001908083835b60208310610a8e5780518252602082019150602081019050602083039250610a6b565b6001836020036101000a038019825116818451168082178552505050505050905001915050602060405180830381855afa158015610ad0573d6000803e3d6000fd5b5050506040513d6020811015610ae557600080fd5b810190808051906020019092919050505090506000600280848c8c604051602001808481526020018383808284378083019250505093505050506040516020818303038152906040526040518082805190602001908083835b60208310610b615780518252602082019150602081019050602083039250610b3e565b6001836020036101000a038019825116818451168082178552505050505050905001915050602060405180830381855afa158015610ba3573d6000803e3d6000fd5b5050506040513d6020811015610bb857600080fd5b8101908080519060200190929190505050600286600060401b866040516020018084805190602001908083835b60208310610c085780518252602082019150602081019050602083039250610be5565b6001836020036101000a0380198251168184511680821785525050505050509050018367ffffffffffffffff191667ffffffffffffffff1916815260180182815260200193505050506040516020818303038152906040526040518082805190602001908083835b60208310610c935780518252602082019150602081019050602083039250610c70565b6001836020036101000a038019825116818451168082178552505050505050905001915050602060405180830381855afa158015610cd5573d6000803e3d6000fd5b5050506040513d6020811015610cea57600080fd5b810190808051906020019092919050505060405160200180838152602001828152602001925050506040516020818303038152906040526040518082805190602001908083835b60208310610d545780518252602082019150602081019050602083039250610d31565b6001836020036101000a038019825116818451168082178552505050505050905001915050602060405180830381855afa158015610d96573d6000803e3d6000fd5b5050506040513d6020811015610dab57600080fd5b81019080805190602001909291905050509050858114610e16576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252605481526020018061162f6054913960600191505060405180910390fd5b6001602060020a0360205410610e77576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602181526020018061160e6021913960400191505060405180910390fd5b60016020600082825401925050819055506000602054905060008090505b6020811015610fb75760018083161415610ec8578260008260208110610eb757fe5b018190555050505050505050610fc7565b600260008260208110610ed757fe5b01548460405160200180838152602001828152602001925050506040516020818303038152906040526040518082805190602001908083835b60208310610f335780518252602082019150602081019050602083039250610f10565b6001836020036101000a038019825116818451168082178552505050505050905001915050602060405180830381855afa158015610f75573d6000803e3d6000fd5b5050506040513d6020811015610f8a57600080fd5b8101908080519060200190929190505050925060028281610fa757fe5b0491508080600101915050610e95565b506000610fc057fe5b5050505050505b50505050505050565b6060610fdd602054611314565b905090565b6000806000602054905060008090505b60208110156111d057600180831614156110e05760026000826020811061101557fe5b01548460405160200180838152602001828152602001925050506040516020818303038152906040526040518082805190602001908083835b60208310611071578051825260208201915060208101905060208303925061104e565b6001836020036101000a038019825116818451168082178552505050505050905001915050602060405180830381855afa1580156110b3573d6000803e3d6000fd5b5050506040513d60208110156110c857600080fd5b810190808051906020019092919050505092506111b6565b600283602183602081106110f057fe5b015460405160200180838152602001828152602001925050506040516020818303038152906040526040518082805190602001908083835b6020831061114b5780518252602082019150602081019050602083039250611128565b6001836020036101000a038019825116818451168082178552505050505050905001915050602060405180830381855afa15801561118d573d6000803e3d6000fd5b5050506040513d60208110156111a257600080fd5b810190808051906020019092919050505092505b600282816111c057fe5b0491508080600101915050610ff2565b506002826111df602054611314565b600060401b6040516020018084815260200183805190602001908083835b6020831061122057805182526020820191506020810190506020830392506111fd565b6001836020036101000a0380198251168184511680821785525050505050509050018267ffffffffffffffff191667ffffffffffffffff1916815260180193505050506040516020818303038152906040526040518082805190602001908083835b602083106112a55780518252602082019150602081019050602083039250611282565b6001836020036101000a038019825116818451168082178552505050505050905001915050602060405180830381855afa1580156112e7573d6000803e3d6000fd5b5050506040513d60208110156112fc57600080fd5b81019080805190602001909291905050509250505090565b6060600867ffffffffffffffff8111801561132e57600080fd5b506040519080825280601f01601f1916602001820160405280156113615781602001600182028036833780820191505090505b50905060008260c01b90508060076008811061137957fe5b1a60f81b8260008151811061138a57fe5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a905350806006600881106113c657fe5b1a60f81b826001815181106113d757fe5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053508060056008811061141357fe5b1a60f81b8260028151811061142457fe5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053508060046008811061146057fe5b1a60f81b8260038151811061147157fe5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a905350806003600881106114ad57fe5b1a60f81b826004815181106114be57fe5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a905350806002600881106114fa57fe5b1a60f81b8260058151811061150b57fe5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053508060016008811061154757fe5b1a60f81b8260068151811061155857fe5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053508060006008811061159457fe5b1a60f81b826007815181106115a557fe5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a90535050919050565b600080858511156115ea57600080fd5b838611156115f757600080fd5b600185028301915084860390509450949250505056fe4465706f736974436f6e74726163743a206d65726b6c6520747265652066756c6c4465706f736974436f6e74726163743a207265636f6e7374727563746564204465706f7369744461746120646f6573206e6f74206d6174636820737570706c696564206465706f7369745f646174615f726f6f744465706f736974436f6e74726163743a20696e76616c6964207769746864726177616c5f63726564656e7469616c73206c656e6774684465706f736974436f6e74726163743a206465706f7369742076616c7565206e6f74206d756c7469706c65206f6620677765694465706f736974436f6e74726163743a20696e76616c6964207075626b6579206c656e6774684465706f736974436f6e74726163743a206465706f7369742076616c756520746f6f20686967684465706f736974436f6e74726163743a206465706f7369742076616c756520746f6f206c6f774465706f736974436f6e74726163743a20696e76616c6964207369676e6174757265206c656e677468a2646970667358221220230afd4b6e3551329e50f1239e08fa3ab7907b77403c4f237d9adf679e8e43cf64736f6c634300060b0033",
			"balance": "0x0"
		},
		"4e59b44847b379578588920ca78fbf26c0b4956c": {
			"code": "0x7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe03601600081602082378035828234f58015156039578182fd5b8082525050506014600cf3",
			"balance": "0x0"
		},
		"5678e9e827b3be0e3d4b910126a64a697a148267": {
			"balance": "0x43c33c1937564800000"
		},
		"70997970c51812dc3a010c7d01b50e0d17dc79c8": {
			"balance": "0x21e19e0c9bab2400000"
		},
		"71be63f3384f5fb98995898a86b02fb2426c5788": {
			"balance": "0x21e19e0c9bab2400000"
		},
		"8626f6940e2eb28930efb4cef49b2d1f2c9c1199": {
			"balance": "0x21e19e0c9bab2400000"
		},
		"90f79bf6eb2c4f870365e785982e1f101e93b906": {
			"balance": "0x21e19e0c9bab2400000"
		},
		"976ea74026e726554db657fa54763abd0c3a0aa9": {
			"balance": "0x21e19e0c9bab2400000"
		},
		"9965507d1a55bcc2695c58ba16fb37d819b0a4dc": {
			"balance": "0x21e19e0c9bab2400000"
		},
		"a0ee7a142d267c1f36714e4a8f75612f20a79720": {
			"balance": "0x21e19e0c9bab2400000"
		},
		"bcd4042de499d14e55001ccbb24a551f3b954096": {
			"balance": "0x21e19e0c9bab2400000"
		},
		"bda5747bfd65f08deb54cb465eb87d40e51b197e": {
			"balance": "0x21e19e0c9bab2400000"
		},
		"cd3b766ccdd6ae721141f452c550ca635964ce71": {
			"balance": "0x21e19e0c9bab2400000"
		},
		"dd2fd4581271e230360230f9337d5c0430bf44c0": {
			"balance": "0x21e19e0c9bab2400000"
		},
		"df3e18d64bc6a983f673ab319ccae4f1a57c7097": {
			"balance": "0x21e19e0c9bab2400000"
		},
		"f39fd6e51aad88f6f4ce6ab8827279cfffb92266": {
			"balance": "0x21e19e0c9bab2400000"
		},
		"fabb0ac9d68b0b445fb7357272ff202c5651694a": {
			"balance": "0x21e19e0c9bab2400000"
		}
	},
	"number": "0x0",
	"gasUsed": "0x0",
	"parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
	"baseFeePerGas": null,
	"excessBlobGas": null,
	"blobGasUsed": null
}
`
