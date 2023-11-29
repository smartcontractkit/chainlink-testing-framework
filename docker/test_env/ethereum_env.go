package test_env

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	tc "github.com/testcontainers/testcontainers-go"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/docker"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
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
}

type EthereumNetworkParticipant struct {
	ConsensusLayer ConsensusLayer `json:"consensus_layer"` //nil means PoW
	ExecutionLayer ExecutionLayer `json:"execution_layer"`
	Count          int            `json:"count"`
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

	if b.existingConfig != nil {
		n.isRecreated = true
		n.Containers = b.existingConfig.Containers
	}

	n.WaitForFinalization = b.waitForFinalization
	n.ethereumChainConfig = b.ethereumChainConfig
	n.t = b.t

	return n
}

func (b *EthereumNetworkBuilder) Build() (EthereumNetwork, error) {
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
	b.dockerNetworks = b.existingConfig.DockerNetworkNames
	b.ethereumChainConfig = b.existingConfig.ethereumChainConfig

	return true
}

func (b *EthereumNetworkBuilder) validate() error {
	if b.consensusType == "" {
		return errors.New("consensus type is required")
	}

	if b.executionLayer == "" {
		return errors.New("execution layer is required")
	}

	if b.consensusType == ConsensusType_PoS && b.consensusLayer == nil {
		return errors.New("consensus layer is required for PoS")
	}

	if b.consensusType == ConsensusType_PoW && b.consensusLayer != nil {
		return errors.New("consensus layer is not allowed for PoW")
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
	ConsensusType        ConsensusType             `json:"consensus_type"`
	ConsensusLayer       *ConsensusLayer           `json:"consensus_layer"`
	ExecutionLayer       ExecutionLayer            `json:"execution_layer"`
	DockerNetworkNames   []string                  `json:"docker_network_names"`
	Containers           EthereumNetworkContainers `json:"containers"`
	WaitForFinalization  bool                      `json:"wait_for_finalization"`
	GeneratedDataHostDir string                    `json:"generated_data_host_dir"`
	ValKeysDir           string                    `json:"val_keys_dir"`
	isRecreated          bool
	ethereumChainConfig  *EthereumChainConfig
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
	var networkNames []string

	if *en.ConsensusLayer != ConsensusLayer_Prysm {
		return blockchain.EVMNetwork{}, RpcProvider{}, fmt.Errorf("unsupported consensus layer: %s. Use 'prysm'", *en.ConsensusLayer)
	}

	singleNetwork, err := en.getOrCreateDockerNetworks()
	if err != nil {
		return blockchain.EVMNetwork{}, RpcProvider{}, err
	}
	networkNames = append(networkNames, singleNetwork...)

	var generatedDataHostDir, valKeysDir string

	// create host directories and run genesis containers only if we are NOT recreating existing containers
	if !en.isRecreated {
		generatedDataHostDir, valKeysDir, err = createHostDirectories()

		en.GeneratedDataHostDir = generatedDataHostDir
		en.ValKeysDir = valKeysDir

		if err != nil {
			return blockchain.EVMNetwork{}, RpcProvider{}, err
		}

		valKeysGeneretor := NewValKeysGeneretor(en.ethereumChainConfig, valKeysDir).WithTestInstance(en.t)
		err = valKeysGeneretor.StartContainer()
		if err != nil {
			return blockchain.EVMNetwork{}, RpcProvider{}, err
		}

		genesis := NewEthGenesisGenerator(*en.ethereumChainConfig, generatedDataHostDir).WithTestInstance(en.t)
		err = genesis.StartContainer()
		if err != nil {
			return blockchain.EVMNetwork{}, RpcProvider{}, err
		}

		initHelper := NewInitHelper(*en.ethereumChainConfig, generatedDataHostDir).WithTestInstance(en.t)
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
	switch en.ExecutionLayer {
	case ExecutionLayer_Geth:
		client = NewGeth2(singleNetwork, en.ethereumChainConfig, generatedDataHostDir, ConsensusLayer_Prysm, en.setExistingContainerName(ContainerType_Geth2)).WithTestInstance(en.t)
	case ExecutionLayer_Nethermind:
		client = NewNethermind(singleNetwork, generatedDataHostDir, ConsensusLayer_Prysm, en.setExistingContainerName(ContainerType_Nethermind)).WithTestInstance(en.t)
	case ExecutionLayer_Erigon:
		client = NewErigon(singleNetwork, en.ethereumChainConfig, generatedDataHostDir, ConsensusLayer_Prysm, en.setExistingContainerName(ContainerType_Erigon)).WithTestInstance(en.t)
	case ExecutionLayer_Besu:
		client = NewBesu(singleNetwork, en.ethereumChainConfig, generatedDataHostDir, ConsensusLayer_Prysm, en.setExistingContainerName(ContainerType_Besu)).WithTestInstance(en.t)
	default:
		return blockchain.EVMNetwork{}, RpcProvider{}, fmt.Errorf("unsupported execution layer: %s", en.ExecutionLayer)
	}

	net, err = client.StartContainer()
	if err != nil {
		return blockchain.EVMNetwork{}, RpcProvider{}, err
	}

	beacon := NewPrysmBeaconChain(singleNetwork, en.ethereumChainConfig, generatedDataHostDir, client.GetInternalExecutionURL(), en.setExistingContainerName(ContainerType_PrysmBeacon)).WithTestInstance(en.t)
	err = beacon.StartContainer()
	if err != nil {
		return blockchain.EVMNetwork{}, RpcProvider{}, err
	}

	validator := NewPrysmValidator(singleNetwork, en.ethereumChainConfig, generatedDataHostDir, valKeysDir, beacon.InternalBeaconRpcProvider, en.setExistingContainerName(ContainerType_PrysmVal)).WithTestInstance(en.t)
	err = validator.StartContainer()
	if err != nil {
		return blockchain.EVMNetwork{}, RpcProvider{}, err
	}

	err = client.WaitUntilChainIsReady(en.ethereumChainConfig.GetDefaultWaitDuration())
	if err != nil {
		return blockchain.EVMNetwork{}, RpcProvider{}, err
	}

	en.DockerNetworkNames = networkNames
	net.ChainID = int64(en.ethereumChainConfig.ChainID)
	net.FinalityTag = true
	if en.ExecutionLayer == ExecutionLayer_Besu {
		// Besu doesn't support "eth_maxPriorityFeePerGas" https://github.com/hyperledger/besu/issues/5658
		// And if gas is too low, then transaction doesn't get to prioritized pool and is not a candidate for inclusion in the next block
		net.GasEstimationBuffer = 10_000_000_000
	} else {
		net.SupportsEIP1559 = true
	}

	logger := logging.GetTestLogger(en.t)
	if en.WaitForFinalization {
		evmClient, err := blockchain.NewEVMClientFromNetwork(net, logger)
		if err != nil {
			return blockchain.EVMNetwork{}, RpcProvider{}, err
		}

		err = waitForChainToFinaliseAnEpoch(logger, evmClient)
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
	var networkNames []string
	rpcProvider := RpcProvider{
		privateHttpUrls: []string{},
		privatelWsUrls:  []string{},
		publiclHttpUrls: []string{},
		publicsUrls:     []string{},
	}

	if en.ExecutionLayer != ExecutionLayer_Geth {
		return blockchain.EVMNetwork{}, RpcProvider{}, fmt.Errorf("unsupported execution layer: %s", en.ExecutionLayer)
	}
	singleNetwork, err := en.getOrCreateDockerNetworks()
	if err != nil {
		return blockchain.EVMNetwork{}, RpcProvider{}, err
	}

	geth := NewGeth(singleNetwork, *en.ethereumChainConfig, en.setExistingContainerName(ContainerType_Geth)).WithTestInstance(en.t)
	network, docker, err := geth.StartContainer()
	if err != nil {
		return blockchain.EVMNetwork{}, RpcProvider{}, err
	}

	net = network
	networkNames = append(networkNames, singleNetwork...)

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

	en.DockerNetworkNames = networkNames

	return net, rpcProvider, nil
}

func (en *EthereumNetwork) getOrCreateDockerNetworks() ([]string, error) {
	var networkNames []string

	if len(en.DockerNetworkNames) == 0 {
		network, err := docker.CreateNetwork(logging.GetTestLogger(en.t))
		if err != nil {
			return networkNames, err
		}
		networkNames = []string{network.Name}
	} else {
		networkNames = en.DockerNetworkNames
	}

	return networkNames, nil
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

func waitForChainToFinaliseAnEpoch(lggr zerolog.Logger, evmClient blockchain.EVMClient) error {
	lggr.Info().Msg("Waiting for chain to finalize an epoch")

	timeout := 180 * time.Second
	pollInterval := 15 * time.Second
	endTime := time.Now().Add(timeout)

	chainStarted := false
	for {
		finalized, err := evmClient.GetLatestFinalizedBlockHeader(context.Background())
		if err != nil {
			if strings.Contains(err.Error(), "finalized block not found") {
				lggr.Err(err).Msgf("error getting finalized block number for %s", evmClient.GetNetworkName())
			} else {
				lggr.Warn().Msgf("no epoch finalized yet for chain %s", evmClient.GetNetworkName())
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
