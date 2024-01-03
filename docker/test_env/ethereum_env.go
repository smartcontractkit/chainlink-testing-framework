package test_env

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	tc "github.com/testcontainers/testcontainers-go"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/docker"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	utils "github.com/smartcontractkit/chainlink-testing-framework/utils/json"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/testcontext"
)

const (
	CONFIG_ENV_VAR_NAME      = "PRIVATE_ETHEREUM_NETWORK_CONFIG_PATH"
	EXEC_CLIENT_ENV_VAR_NAME = "ETH2_EL_CLIENT"
)

var (
	ErrMissingConsensusType     = errors.New("consensus type is required")
	ErrMissingExecutionLayer    = errors.New("execution layer is required")
	ErrMissingConsensusLayer    = errors.New("consensus layer is required for PoS")
	ErrConsensusLayerNotAllowed = errors.New("consensus layer is not allowed for PoW")
	ErrTestConfigNotSaved       = errors.New("could not save test env config")
)

type ConsensusType string

const (
	ConsensusType_PoS ConsensusType = "pos"
	ConsensusType_PoW ConsensusType = "pow"
)

type ExecutionLayer string

const (
	ExecutionLayer_Geth       ExecutionLayer = "geth"
	ExecutionLayer_Nethermind ExecutionLayer = "nethermind"
	ExecutionLayer_Erigon     ExecutionLayer = "erigon"
	ExecutionLayer_Besu       ExecutionLayer = "besu"
)

type ConsensusLayer string

const (
	ConsensusLayer_Prysm ConsensusLayer = "prysm"
)

type EthereumNetworkBuilder struct {
	t                   *testing.T
	dockerNetworks      []string
	consensusType       ConsensusType
	consensusLayer      *ConsensusLayer
	executionLayer      ExecutionLayer
	ethereumChainConfig *EthereumChainConfig
	existingConfig      *EthereumNetwork
	addressesToFund     []string
	waitForFinalization bool
	existingFromEnvVar  bool
}

func NewEthereumNetworkBuilder() EthereumNetworkBuilder {
	return EthereumNetworkBuilder{
		dockerNetworks:      []string{},
		waitForFinalization: false,
	}
}

func (b *EthereumNetworkBuilder) WithConsensusType(consensusType ConsensusType) *EthereumNetworkBuilder {
	b.consensusType = consensusType
	return b
}

func (b *EthereumNetworkBuilder) WithConsensusLayer(consensusLayer ConsensusLayer) *EthereumNetworkBuilder {
	b.consensusLayer = &consensusLayer
	return b
}

func (b *EthereumNetworkBuilder) WithExecutionLayer(executionLayer ExecutionLayer) *EthereumNetworkBuilder {
	b.executionLayer = executionLayer
	return b
}

func (b *EthereumNetworkBuilder) WithEthereumChainConfig(config EthereumChainConfig) *EthereumNetworkBuilder {
	b.ethereumChainConfig = &config
	return b
}

func (b *EthereumNetworkBuilder) WithDockerNetworks(networks []string) *EthereumNetworkBuilder {
	b.dockerNetworks = networks
	return b
}

func (b *EthereumNetworkBuilder) WithExistingConfig(config EthereumNetwork) *EthereumNetworkBuilder {
	b.existingConfig = &config
	return b
}

func (b *EthereumNetworkBuilder) WihtExistingConfigFromEnvVar() *EthereumNetworkBuilder {
	b.existingFromEnvVar = true
	return b
}

func (b *EthereumNetworkBuilder) WithTest(t *testing.T) *EthereumNetworkBuilder {
	b.t = t
	return b
}

func (b *EthereumNetworkBuilder) WithWaitingForFinalization() *EthereumNetworkBuilder {
	b.waitForFinalization = true
	return b
}

func (b *EthereumNetworkBuilder) buildNetworkConfig() EthereumNetwork {
	n := EthereumNetwork{
		ConsensusType:  b.consensusType,
		ExecutionLayer: b.executionLayer,
		ConsensusLayer: b.consensusLayer,
	}

	if b.existingConfig != nil && len(b.existingConfig.Containers) > 0 {
		n.isRecreated = true
		n.Containers = b.existingConfig.Containers
	}

	n.DockerNetworkNames = b.dockerNetworks
	n.WaitForFinalization = &b.waitForFinalization
	n.EthereumChainConfig = b.ethereumChainConfig
	n.t = b.t

	return n
}

func (b *EthereumNetworkBuilder) Build() (EthereumNetwork, error) {
	if b.existingFromEnvVar {
		path := os.Getenv(CONFIG_ENV_VAR_NAME)
		if path == "" {
			return EthereumNetwork{}, fmt.Errorf("environment variable %s is not set, but build from env var was requested", CONFIG_ENV_VAR_NAME)
		}

		config, err := NewPrivateChainEnvConfigFromFile(path)
		if err != nil {
			return EthereumNetwork{}, err
		}

		config.isRecreated = true

		return config, nil
	}

	if !b.importExistingConfig() {
		if b.ethereumChainConfig == nil {
			defaultConfig := GetDefaultChainConfig()
			b.ethereumChainConfig = &defaultConfig
		} else {
			b.ethereumChainConfig.fillInMissingValuesWithDefault()
		}

		b.ethereumChainConfig.GenerateGenesisTimestamp()
	}

	err := b.validate()
	if err != nil {
		return EthereumNetwork{}, err
	}

	return b.buildNetworkConfig(), nil
}

func (b *EthereumNetworkBuilder) importExistingConfig() bool {
	if b.existingConfig == nil {
		return false
	}

	b.consensusType = b.existingConfig.ConsensusType
	b.consensusLayer = b.existingConfig.ConsensusLayer
	b.executionLayer = b.existingConfig.ExecutionLayer
	if len(b.existingConfig.DockerNetworkNames) > 0 {
		b.dockerNetworks = b.existingConfig.DockerNetworkNames
	}
	b.ethereumChainConfig = b.existingConfig.EthereumChainConfig

	return true
}

func (b *EthereumNetworkBuilder) validate() error {
	if b.consensusType == "" {
		return ErrMissingConsensusType
	}

	if b.executionLayer == "" {
		return ErrMissingExecutionLayer
	}

	if b.consensusType == ConsensusType_PoS && b.consensusLayer == nil {
		return ErrMissingConsensusLayer
	}

	if b.consensusType == ConsensusType_PoW && b.consensusLayer != nil {
		return ErrConsensusLayerNotAllowed
	}

	for _, addr := range b.addressesToFund {
		if !common.IsHexAddress(addr) {
			return fmt.Errorf("address %s is not a valid hex address", addr)
		}
	}

	err := b.ethereumChainConfig.Validate(logging.GetTestLogger(b.t))
	if err != nil {
		return err
	}

	return nil
}

type EthereumNetwork struct {
	ConsensusType        ConsensusType             `json:"consensus_type" toml:"consensus_type"`
	ConsensusLayer       *ConsensusLayer           `json:"consensus_layer" toml:"consensus_layer"`
	ExecutionLayer       ExecutionLayer            `json:"execution_layer" toml:"execution_layer"`
	DockerNetworkNames   []string                  `json:"docker_network_names"`
	Containers           EthereumNetworkContainers `json:"containers"`
	WaitForFinalization  *bool                     `json:"wait_for_finalization" toml:"wait_for_finalization"`
	GeneratedDataHostDir string                    `json:"generated_data_host_dir"`
	ValKeysDir           string                    `json:"val_keys_dir"`
	EthereumChainConfig  *EthereumChainConfig      `json:"ethereum_chain_config" toml:"EthereumChainConfig"`
	isRecreated          bool
	t                    *testing.T
}

func (en *EthereumNetwork) Start() (blockchain.EVMNetwork, RpcProvider, error) {
	switch en.ConsensusType {
	case ConsensusType_PoS:
		return en.startPos()
	case ConsensusType_PoW:
		return en.startPow()
	default:
		return blockchain.EVMNetwork{}, RpcProvider{}, fmt.Errorf("unknown consensus type: %s", en.ConsensusType)
	}
}

func (en *EthereumNetwork) startPos() (blockchain.EVMNetwork, RpcProvider, error) {
	rpcProvider := RpcProvider{
		privateHttpUrls: []string{},
		privatelWsUrls:  []string{},
		publiclHttpUrls: []string{},
		publicsUrls:     []string{},
	}

	var net blockchain.EVMNetwork

	if *en.ConsensusLayer != ConsensusLayer_Prysm {
		return blockchain.EVMNetwork{}, RpcProvider{}, fmt.Errorf("unsupported consensus layer: %s. Use 'prysm'", *en.ConsensusLayer)
	}

	dockerNetworks, err := en.getOrCreateDockerNetworks()
	if err != nil {
		return blockchain.EVMNetwork{}, RpcProvider{}, err
	}
	var generatedDataHostDir, valKeysDir string

	// create host directories and run genesis containers only if we are NOT recreating existing containers
	if !en.isRecreated {
		generatedDataHostDir, valKeysDir, err = createHostDirectories()

		en.GeneratedDataHostDir = generatedDataHostDir
		en.ValKeysDir = valKeysDir

		if err != nil {
			return blockchain.EVMNetwork{}, RpcProvider{}, err
		}

		valKeysGeneretor, err := NewValKeysGeneretor(en.EthereumChainConfig, valKeysDir)
		if err != nil {
			return blockchain.EVMNetwork{}, RpcProvider{}, err
		}
		valKeysGeneretor.WithTestInstance(en.t)

		err = valKeysGeneretor.StartContainer()
		if err != nil {
			return blockchain.EVMNetwork{}, RpcProvider{}, err
		}

		genesis, err := NewEthGenesisGenerator(*en.EthereumChainConfig, generatedDataHostDir)
		if err != nil {
			return blockchain.EVMNetwork{}, RpcProvider{}, err
		}

		genesis.WithTestInstance(en.t)

		err = genesis.StartContainer()
		if err != nil {
			return blockchain.EVMNetwork{}, RpcProvider{}, err
		}

		initHelper := NewInitHelper(*en.EthereumChainConfig, generatedDataHostDir).WithTestInstance(en.t)
		err = initHelper.StartContainer()
		if err != nil {
			return blockchain.EVMNetwork{}, RpcProvider{}, err
		}
	} else {
		// we don't set actual values to not increase complexity, as they do not matter for containers that are already running
		generatedDataHostDir = ""
		valKeysDir = ""
	}

	var client ExecutionClient
	var clientErr error
	switch en.ExecutionLayer {
	case ExecutionLayer_Geth:
		client, clientErr = NewGeth2(dockerNetworks, en.EthereumChainConfig, generatedDataHostDir, ConsensusLayer_Prysm, en.setExistingContainerName(ContainerType_Geth2))
	case ExecutionLayer_Nethermind:
		client, clientErr = NewNethermind(dockerNetworks, generatedDataHostDir, ConsensusLayer_Prysm, en.setExistingContainerName(ContainerType_Nethermind))
	case ExecutionLayer_Erigon:
		client, clientErr = NewErigon(dockerNetworks, en.EthereumChainConfig, generatedDataHostDir, ConsensusLayer_Prysm, en.setExistingContainerName(ContainerType_Erigon))
	case ExecutionLayer_Besu:
		client, clientErr = NewBesu(dockerNetworks, en.EthereumChainConfig, generatedDataHostDir, ConsensusLayer_Prysm, en.setExistingContainerName(ContainerType_Besu))
	default:
		return blockchain.EVMNetwork{}, RpcProvider{}, fmt.Errorf("unsupported execution layer: %s", en.ExecutionLayer)
	}

	if clientErr != nil {
		return blockchain.EVMNetwork{}, RpcProvider{}, clientErr
	}

	client.WithTestInstance(en.t)

	net, err = client.StartContainer()
	if err != nil {
		return blockchain.EVMNetwork{}, RpcProvider{}, err
	}

	beacon, err := NewPrysmBeaconChain(dockerNetworks, en.EthereumChainConfig, generatedDataHostDir, client.GetInternalExecutionURL(), en.setExistingContainerName(ContainerType_PrysmBeacon))
	if err != nil {
		return blockchain.EVMNetwork{}, RpcProvider{}, err
	}

	beacon.WithTestInstance(en.t)
	err = beacon.StartContainer()
	if err != nil {
		return blockchain.EVMNetwork{}, RpcProvider{}, err
	}

	validator, err := NewPrysmValidator(dockerNetworks, en.EthereumChainConfig, generatedDataHostDir, valKeysDir, beacon.
		InternalBeaconRpcProvider, en.setExistingContainerName(ContainerType_PrysmVal))
	if err != nil {
		return blockchain.EVMNetwork{}, RpcProvider{}, err
	}

	validator.WithTestInstance(en.t)
	err = validator.StartContainer()
	if err != nil {
		return blockchain.EVMNetwork{}, RpcProvider{}, err
	}

	err = client.WaitUntilChainIsReady(testcontext.Get(en.t), en.EthereumChainConfig.GetDefaultWaitDuration())
	if err != nil {
		return blockchain.EVMNetwork{}, RpcProvider{}, err
	}

	en.DockerNetworkNames = dockerNetworks
	net.ChainID = int64(en.EthereumChainConfig.ChainID)
	// use a higher value than the default, because eth2 is slower than dev-mode eth1
	net.Timeout = blockchain.StrDuration{Duration: time.Duration(4 * time.Minute)}
	net.FinalityTag = true
	net.FinalityDepth = 0

	if en.ExecutionLayer == ExecutionLayer_Besu {
		// Besu doesn't support "eth_maxPriorityFeePerGas" https://github.com/hyperledger/besu/issues/5658
		// And if gas is too low, then transaction doesn't get to prioritized pool and is not a candidate for inclusion in the next block
		net.GasEstimationBuffer = 10_000_000_000
	} else {
		net.SupportsEIP1559 = true
	}

	logger := logging.GetTestLogger(en.t)
	if en.WaitForFinalization != nil && *en.WaitForFinalization {
		evmClient, err := blockchain.NewEVMClientFromNetwork(net, logger)
		if err != nil {
			return blockchain.EVMNetwork{}, RpcProvider{}, err
		}

		err = waitForChainToFinaliseAnEpoch(logger, evmClient, en.EthereumChainConfig.GetDefaultFinalizationWaitDuration())
		if err != nil {
			return blockchain.EVMNetwork{}, RpcProvider{}, err
		}
	} else {
		logger.Info().Msg("Not waiting for chain to finalize first epoch")
	}

	containers := EthereumNetworkContainers{
		{
			ContainerName: client.GetContainerName(),
			ContainerType: client.GetContainerType(),
			Container:     client.GetContainer(),
		},
		{
			ContainerName: beacon.ContainerName,
			ContainerType: ContainerType_PrysmBeacon,
			Container:     &beacon.Container,
		},
		{
			ContainerName: validator.ContainerName,
			ContainerType: ContainerType_PrysmVal,
			Container:     &validator.Container,
		},
	}

	en.Containers = append(en.Containers, containers...)

	rpcProvider.privateHttpUrls = append(rpcProvider.privateHttpUrls, client.GetInternalHttpUrl())
	rpcProvider.privatelWsUrls = append(rpcProvider.privatelWsUrls, client.GetInternalWsUrl())
	rpcProvider.publiclHttpUrls = append(rpcProvider.publiclHttpUrls, client.GetExternalHttpUrl())
	rpcProvider.publicsUrls = append(rpcProvider.publicsUrls, client.GetExternalWsUrl())

	return net, rpcProvider, nil
}

func (en *EthereumNetwork) startPow() (blockchain.EVMNetwork, RpcProvider, error) {
	var net blockchain.EVMNetwork
	rpcProvider := RpcProvider{
		privateHttpUrls: []string{},
		privatelWsUrls:  []string{},
		publiclHttpUrls: []string{},
		publicsUrls:     []string{},
	}

	if en.ExecutionLayer != ExecutionLayer_Geth {
		return blockchain.EVMNetwork{}, RpcProvider{}, fmt.Errorf("unsupported execution layer: %s", en.ExecutionLayer)
	}
	dockerNetworks, err := en.getOrCreateDockerNetworks()
	if err != nil {
		return blockchain.EVMNetwork{}, RpcProvider{}, err
	}

	geth := NewGeth(dockerNetworks, en.EthereumChainConfig, en.setExistingContainerName(ContainerType_Geth)).WithTestInstance(en.t)
	network, docker, err := geth.StartContainer()
	if err != nil {
		return blockchain.EVMNetwork{}, RpcProvider{}, err
	}

	net = network
	containers := EthereumNetworkContainers{
		{
			ContainerName: geth.ContainerName,
			ContainerType: ContainerType_Geth,
			Container:     &geth.Container,
		},
	}

	en.Containers = append(en.Containers, containers...)
	rpcProvider.privateHttpUrls = append(rpcProvider.privateHttpUrls, docker.HttpUrl)
	rpcProvider.privatelWsUrls = append(rpcProvider.privatelWsUrls, docker.WsUrl)
	rpcProvider.publiclHttpUrls = append(rpcProvider.publiclHttpUrls, geth.ExternalHttpUrl)
	rpcProvider.publicsUrls = append(rpcProvider.publicsUrls, geth.ExternalWsUrl)

	en.DockerNetworkNames = dockerNetworks

	return net, rpcProvider, nil
}

func (en *EthereumNetwork) getOrCreateDockerNetworks() ([]string, error) {
	if len(en.DockerNetworkNames) != 0 {
		return en.DockerNetworkNames, nil
	}

	network, err := docker.CreateNetwork(logging.GetTestLogger(en.t))
	if err != nil {
		return []string{}, err
	}

	return []string{network.Name}, nil
}

func (en *EthereumNetwork) Describe() string {
	cL := "prysm"
	if en.ConsensusLayer == nil {
		cL = "(none)"
	}
	return fmt.Sprintf("consensus type: %s, execution layer: %s, consensus layer: %s", en.ConsensusType, en.ExecutionLayer, cL)
}

func (en *EthereumNetwork) setExistingContainerName(ct ContainerType) EnvComponentOption {
	if !en.isRecreated {
		return func(c *EnvComponent) {}
	}

	for _, container := range en.Containers {
		if container.ContainerType == ct {
			return func(c *EnvComponent) {
				if container.ContainerName != "" {
					c.ContainerName = container.ContainerName
				}
			}
		}
	}

	return func(c *EnvComponent) {}
}

func (en *EthereumNetwork) Save() error {
	name := fmt.Sprintf("ethereum_network_%s", uuid.NewString()[0:8])
	confPath, err := utils.SaveStructAsJson(en, ".private_chains", name)
	if err != nil {
		return ErrTestConfigNotSaved
	}

	log := logging.GetTestLogger(en.t)
	log.Info().Msgf("Saved private Ethereum Network config. To reuse in e2e tests, set: %s=%s", CONFIG_ENV_VAR_NAME, confPath)

	return nil
}

func (en *EthereumNetwork) Validate() error {
	if en.ConsensusType == "" {
		return ErrMissingConsensusType
	}

	if en.ExecutionLayer == "" {
		return ErrMissingExecutionLayer
	}

	if en.ConsensusType == ConsensusType_PoS && en.ConsensusLayer == nil {
		return ErrMissingConsensusLayer
	}

	if en.ConsensusType == ConsensusType_PoW && en.ConsensusLayer != nil {
		return ErrConsensusLayerNotAllowed
	}

	if en.EthereumChainConfig == nil {
		return errors.New("ethereum chain config is required")
	}

	err := en.EthereumChainConfig.Validate(logging.GetTestLogger(nil))
	if err != nil {
		return err
	}

	return nil
}

func (en *EthereumNetwork) ApplyOverrides(from *EthereumNetwork) error {
	if from == nil {
		return nil
	}
	if from.ConsensusLayer != nil {
		en.ConsensusLayer = from.ConsensusLayer
	}
	if from.ExecutionLayer != "" {
		en.ExecutionLayer = from.ExecutionLayer
	}
	if from.ConsensusType != "" {
		en.ConsensusType = from.ConsensusType
	}
	if from.WaitForFinalization != nil {
		en.WaitForFinalization = from.WaitForFinalization
	}

	if from.EthereumChainConfig != nil {
		if en.EthereumChainConfig == nil {
			en.EthereumChainConfig = from.EthereumChainConfig
		} else {
			err := en.EthereumChainConfig.ApplyOverrides(from.EthereumChainConfig)
			if err != nil {
				return errors.Wrapf(err, "error applying overrides from network config file to config")
			}
		}
	}

	return nil
}

//go:embed tomls/default_ethereum_env.toml
var defaultEthEnvConfig []byte

func (en *EthereumNetwork) Default() error {
	wrapper := struct {
		EthereumNetwork *EthereumNetwork `toml:"PrivateEthereumNetwork"`
	}{}
	if err := toml.Unmarshal(defaultEthEnvConfig, &wrapper); err != nil {
		return errors.Wrapf(err, "error unmarshaling ethereum network config")
	}

	*en = *wrapper.EthereumNetwork

	if en.EthereumChainConfig != nil && en.EthereumChainConfig.genesisTimestamp == 0 {
		en.EthereumChainConfig.GenerateGenesisTimestamp()
	}

	return nil
}

// maybe only store ports here and depending on where the test is executed return different URLs?
// maybe 3 different constructors for each "perspective"? (docker, k8s with local runner, k8s with remote runner)
type RpcProvider struct {
	privateHttpUrls []string
	privatelWsUrls  []string
	publiclHttpUrls []string
	publicsUrls     []string
}

func (s *RpcProvider) PrivateHttpUrls() []string {
	return s.privateHttpUrls
}

func (s *RpcProvider) PrivateWsUrsl() []string {
	return s.privatelWsUrls
}

func (s *RpcProvider) PublicHttpUrls() []string {
	return s.publiclHttpUrls
}

func (s *RpcProvider) PublicWsUrls() []string {
	return s.publicsUrls
}

type ContainerType string

const (
	ContainerType_Geth        ContainerType = "geth"
	ContainerType_Geth2       ContainerType = "geth2"
	ContainerType_Erigon      ContainerType = "erigon"
	ContainerType_Besu        ContainerType = "besu"
	ContainerType_Nethermind  ContainerType = "nethermind"
	ContainerType_PrysmBeacon ContainerType = "prysm-beacon"
	ContainerType_PrysmVal    ContainerType = "prysm-validator"
)

type EthereumNetworkContainer struct {
	ContainerName string        `json:"container_name"`
	ContainerType ContainerType `json:"container_type"`
	Container     *tc.Container `json:"-"`
}

type EthereumNetworkContainers []EthereumNetworkContainer

func createHostDirectories() (string, string, error) {
	customConfigDataDir, err := os.MkdirTemp("", "custom_config_data")
	if err != nil {
		return "", "", err
	}

	valKeysDir, err := os.MkdirTemp("", "val_keys")
	if err != nil {
		return "", "", err
	}

	return customConfigDataDir, valKeysDir, nil
}

func waitForChainToFinaliseAnEpoch(lggr zerolog.Logger, evmClient blockchain.EVMClient, timeout time.Duration) error {
	lggr.Info().Msg("Waiting for chain to finalize an epoch")

	pollInterval := 15 * time.Second
	endTime := time.Now().Add(timeout)

	chainStarted := false
	for {
		finalized, err := evmClient.GetLatestFinalizedBlockHeader(context.Background())
		if err != nil {
			if strings.Contains(err.Error(), "finalized block not found") {
				lggr.Err(err).Msgf("error getting finalized block number for %s", evmClient.GetNetworkName())
			} else {
				timeLeft := time.Until(endTime).Seconds()
				lggr.Warn().Msgf("no epoch finalized yet for chain %s. Time left: %d sec", evmClient.GetNetworkName(), int(timeLeft))
			}
		}

		if finalized != nil && finalized.Number.Int64() > 0 || time.Now().After(endTime) {
			lggr.Info().Msgf("Chain '%s' finalized an epoch", evmClient.GetNetworkName())
			chainStarted = true
			break
		}

		time.Sleep(pollInterval)
	}

	if !chainStarted {
		return fmt.Errorf("chain %s failed to finalize an epoch", evmClient.GetNetworkName())
	}

	return nil
}

func NewPrivateChainEnvConfigFromFile(path string) (EthereumNetwork, error) {
	c := EthereumNetwork{}
	err := utils.OpenJsonFileAsStruct(path, &c)
	return c, err
}
