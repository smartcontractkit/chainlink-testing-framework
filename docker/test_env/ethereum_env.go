package test_env

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"

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
)

type ConsensusLayer string

const (
	ConsensusLayer_Prysm ConsensusLayer = "prysm"
)

type EthereumNetworkBuilder struct {
	t                 *testing.T
	dockerNetworks    []string
	consensusType     *ConsensusType
	consensusLayer    *ConsensusLayer
	consensusNodes    int
	executionLayer    *ExecutionLayer
	executionNodes    int
	beaconChainConfig *BeaconChainConfig
	existingConfig    *EthereumNetwork
	addressesToFund   []string
}

type Eth2Components struct {
	Geth        *Geth2
	BeaconChain *PrysmBeaconChain
	Validator   *PrysmValidator
}

func NewEthereumNetworkBuilder() EthereumNetworkBuilder {
	return EthereumNetworkBuilder{
		dockerNetworks: []string{},
		executionNodes: 1,
		consensusNodes: 1,
	}
}

func (b *EthereumNetworkBuilder) WithConsensusType(consensusType ConsensusType) *EthereumNetworkBuilder {
	b.consensusType = &consensusType
	return b
}

func (b *EthereumNetworkBuilder) WithConsensusLayer(consensusLayer ConsensusLayer) *EthereumNetworkBuilder {
	b.consensusLayer = &consensusLayer
	return b
}

func (b *EthereumNetworkBuilder) WithExecutionLayer(executionLayer ExecutionLayer) *EthereumNetworkBuilder {
	b.executionLayer = &executionLayer
	return b
}

func (b *EthereumNetworkBuilder) WithConsensusNodes(consensusNodes int) *EthereumNetworkBuilder {
	b.consensusNodes = consensusNodes
	return b
}

func (b *EthereumNetworkBuilder) WithExecutionNodes(executionNodes int) *EthereumNetworkBuilder {
	b.executionNodes = executionNodes
	return b
}

func (b *EthereumNetworkBuilder) WithBeaconChainConfig(config BeaconChainConfig) *EthereumNetworkBuilder {
	b.beaconChainConfig = &config
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

func (b *EthereumNetworkBuilder) WithAddressesToFund(addresses []string) *EthereumNetworkBuilder {
	b.addressesToFund = addresses
	return b
}

func (b *EthereumNetworkBuilder) buildConfig() EthereumNetwork {
	n := EthereumNetwork{
		ConsensusType:  *b.consensusType,
		ConsensusNodes: b.consensusNodes,
		ExecutionLayer: *b.executionLayer,
		ExecutionNodes: b.executionNodes,
	}

	if b.consensusLayer != nil {
		consensusLayer := ConsensusLayer(*b.consensusLayer)
		n.ConsensusLayer = consensusLayer
	} else {
		n.ConsensusLayer = ""
	}

	if b.existingConfig != nil {
		n.isRecreated = true
		n.ExecutionDir = b.existingConfig.ExecutionDir
		n.ConsensusDir = b.existingConfig.ConsensusDir
		n.Containers = b.existingConfig.Containers
	} else {
		n.beaconChainConfig = b.beaconChainConfig
	}

	n.logger = logging.GetTestLogger(b.t)
	n.addressesToFund = b.addressesToFund

	return n
}

func (b *EthereumNetworkBuilder) Build() (EthereumNetwork, error) {
	b.importExistingConfig()
	err := b.validate()
	if err != nil {
		return EthereumNetwork{}, err
	}

	return b.buildConfig(), nil
}

func (b *EthereumNetwork) Start() (blockchain.EVMNetwork, RpcProvider, error) {
	switch b.ConsensusType {
	case ConsensusType_PoS:
		return b.startPos()
	case ConsensusType_PoW:
		return b.startPow()
	default:
		return blockchain.EVMNetwork{}, RpcProvider{}, fmt.Errorf("unknown consensus type: %s", b.ConsensusType)
	}
}

func (b *EthereumNetworkBuilder) importExistingConfig() {
	if b.existingConfig == nil {
		return
	}

	if b.existingConfig.ConsensusLayer != "" {
		consensusLayer := ConsensusLayer(b.existingConfig.ConsensusLayer)
		b.consensusLayer = &consensusLayer
	} else {
		b.consensusType = nil
	}
	b.consensusType = &b.existingConfig.ConsensusType
	b.consensusNodes = b.existingConfig.ConsensusNodes
	b.executionLayer = &b.existingConfig.ExecutionLayer
	b.executionNodes = b.existingConfig.ExecutionNodes
	b.dockerNetworks = b.existingConfig.DockerNetworkNames
}

func (b *EthereumNetworkBuilder) validate() error {
	if b.consensusType == nil {
		return errors.New("consensus type is required")
	}
	if b.executionLayer == nil {
		return errors.New("execution layer is required")
	}

	if *b.consensusType == ConsensusType_PoS && b.consensusLayer == nil {
		return errors.New("consensus layer is required for PoS")
	}

	if *b.consensusType == ConsensusType_PoW && b.consensusLayer != nil {
		return errors.New("consensus layer is not allowed for PoW")
	}

	if b.consensusNodes > 1 {
		return errors.New("only one consensus node is currently supported")
	}

	if b.executionNodes > 1 {
		return errors.New("only one execution node is currently supported")
	}

	if *b.consensusType == ConsensusType_PoW && b.beaconChainConfig != nil {
		return errors.New("beacon chain config is not allowed for PoW")
	}

	if *b.consensusType == ConsensusType_PoW {
		b.consensusNodes = 0
	}

	for _, addr := range b.addressesToFund {
		if !common.IsHexAddress(addr) {
			return fmt.Errorf("address %s is not a valid hex address", addr)
		}
	}

	return nil
}

func (b *EthereumNetwork) startPos() (blockchain.EVMNetwork, RpcProvider, error) {
	if b.ConsensusLayer != ConsensusLayer_Prysm {
		return blockchain.EVMNetwork{}, RpcProvider{}, fmt.Errorf("unsupported consensus layer: %s", b.ConsensusLayer)
	}
	networkNames, err := b.getOrCreateDockerNetworks()
	if err != nil {
		return blockchain.EVMNetwork{}, RpcProvider{}, err
	}

	var hostExecutionDir, hostConsensusDir string
	genesisTime := time.Now().Add(20 * time.Second).Unix()

	// create host directories and run genesis containers only if we are NOT recreating existing containers
	if !b.isRecreated {
		hostExecutionDir, hostConsensusDir, err = createHostDirectories()

		if err != nil {
			return blockchain.EVMNetwork{}, RpcProvider{}, err
		}

		b.ExecutionDir = hostExecutionDir
		b.ConsensusDir = hostConsensusDir

		var beaconChainConfig BeaconChainConfig
		if b.beaconChainConfig != nil {
			beaconChainConfig = *b.beaconChainConfig
		} else {
			beaconChainConfig = DefaultBeaconChainConfig
			beaconChainConfig.MinGenesisTime = int(genesisTime)
		}

		bg := NewEth2Genesis(networkNames, beaconChainConfig, hostExecutionDir, hostConsensusDir).
			WithLogger(b.logger).WithFundedAccounts(b.addressesToFund)
		err = bg.StartContainer()
		if err != nil {
			return blockchain.EVMNetwork{}, RpcProvider{}, err
		}

		//TODO make this nicer
		if b.ExecutionLayer == ExecutionLayer_Geth {
			gg := NewEth1Genesis(networkNames, hostExecutionDir).WithLogger(b.logger)
			err = gg.StartContainer()
			if err != nil {
				return blockchain.EVMNetwork{}, RpcProvider{}, err
			}
		}
	}

	var net blockchain.EVMNetwork
	var client ExecutionClient
	if b.ExecutionLayer == ExecutionLayer_Geth {
		client = NewGeth2(networkNames, hostExecutionDir, ConsensusLayer_Prysm, b.setExistingContainerName(ContainerType_Geth2)).WithLogger(b.logger)
	} else {
		client = NewNethermind(networkNames, hostExecutionDir, ConsensusLayer_Prysm, int(genesisTime), b.setExistingContainerName(ContainerType_Nethermind)).WithLogger(b.logger)
	}

	net, err = client.StartContainer()
	if err != nil {
		return blockchain.EVMNetwork{}, RpcProvider{}, err
	}

	beacon := NewPrysmBeaconChain(networkNames, hostExecutionDir, hostConsensusDir, client.GetInternalExecutionURL(), b.setExistingContainerName(ContainerType_PrysmBeacon)).WithLogger(b.logger)
	err = beacon.StartContainer()
	if err != nil {
		return blockchain.EVMNetwork{}, RpcProvider{}, err
	}

	validator := NewPrysmValidator(networkNames, hostConsensusDir, beacon.InternalBeaconRpcProvider, b.setExistingContainerName(ContainerType_PrysmVal)).WithLogger(b.logger)
	err = validator.StartContainer()
	if err != nil {
		return blockchain.EVMNetwork{}, RpcProvider{}, err
	}

	waitForFirstBlock := tcwait.NewLogStrategy("Chain head was updated").WithPollInterval(1 * time.Second).WithStartupTimeout(60 * time.Second)
	err = waitForFirstBlock.WaitUntilReady(context.Background(), *client.GetContainer())
	if err != nil {
		return blockchain.EVMNetwork{}, RpcProvider{}, err
	}

	rpcProvider := RpcProvider{
		privateHttpUrls: []string{client.GetInternalHttpUrl()},
		privatelWsUrls:  []string{client.GetInternalWsUrl()},
		publiclHttpUrls: []string{client.GetExternalHttpUrl()},
		publicsUrls:     []string{client.GetExternalWsUrl()},
	}

	b.DockerNetworkNames = networkNames
	b.Containers = EthereumNetworkContainers{
		{
			ContainerName: client.GetContainerName(),
			ContainerType: ContainerType_Geth2,
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

	return net, rpcProvider, nil
}

func (b *EthereumNetwork) startPow() (blockchain.EVMNetwork, RpcProvider, error) {
	if b.ExecutionLayer != ExecutionLayer_Geth {
		return blockchain.EVMNetwork{}, RpcProvider{}, fmt.Errorf("unsupported execution layer: %s", b.ExecutionLayer)
	}
	networkNames, err := b.getOrCreateDockerNetworks()
	if err != nil {
		return blockchain.EVMNetwork{}, RpcProvider{}, err
	}

	geth := NewGeth(networkNames, b.setExistingContainerName(ContainerType_Geth)).WithLogger(b.logger)
	net, docker, err := geth.StartContainer()
	if err != nil {
		return blockchain.EVMNetwork{}, RpcProvider{}, err
	}

	b.DockerNetworkNames = networkNames
	b.Containers = EthereumNetworkContainers{
		{
			ContainerName: geth.ContainerName,
			ContainerType: ContainerType_Geth,
			Container:     &geth.Container,
		},
	}

	return net, RpcProvider{
		privateHttpUrls: []string{docker.HttpUrl},
		privatelWsUrls:  []string{docker.WsUrl},
		publiclHttpUrls: []string{geth.ExternalHttpUrl},
		publicsUrls:     []string{geth.ExternalWsUrl},
	}, nil
}

func (b *EthereumNetwork) getOrCreateDockerNetworks() ([]string, error) {
	var networkNames []string

	if len(b.DockerNetworkNames) == 0 {
		network, err := docker.CreateNetwork(b.logger)
		if err != nil {
			return networkNames, err
		}
		networkNames = []string{network.Name}
	} else {
		networkNames = b.DockerNetworkNames
	}

	return networkNames, nil
}

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

func (s *RpcProvider) PublicWsUrsl() []string {
	return s.publicsUrls
}

func (b *EthereumNetwork) Describe() string {
	return fmt.Sprintf("consensus type: %s, consensus layer: %s, execution layer: %s, consensus nodes: %d, execution nodes: %d",
		b.ConsensusType, b.ConsensusLayer, b.ExecutionLayer, b.ConsensusNodes, b.ExecutionNodes)
}

type ContainerType string

const (
	ContainerType_Geth        ContainerType = "geth"
	ContainerType_Geth2       ContainerType = "geth2"
	ContainerType_Nethermind  ContainerType = "nethermind"
	ContainerType_PrysmBeacon ContainerType = "prysm-beacon"
	ContainerType_PrysmVal    ContainerType = "prysm-validator"
)

type EthereumNetworkContainer struct {
	ContainerName string        `json:"container_name"`
	ContainerType ContainerType `json:"container_type"`
	Container     *tc.Container `json:"-"`
}

type EthereumNetwork struct {
	ConsensusType      ConsensusType             `json:"consensus_type"`
	ConsensusLayer     ConsensusLayer            `json:"consensus_layer"`
	ConsensusNodes     int                       `json:"consensus_nodes"`
	ExecutionLayer     ExecutionLayer            `json:"execution_layer"`
	ExecutionNodes     int                       `json:"execution_nodes"`
	DockerNetworkNames []string                  `json:"docker_network_names"`
	ExecutionDir       string                    `json:"execution_dir"`
	ConsensusDir       string                    `json:"consensus_dir"`
	Containers         EthereumNetworkContainers `json:"containers"`
	logger             zerolog.Logger
	isRecreated        bool
	beaconChainConfig  *BeaconChainConfig
	addressesToFund    []string
}

type EthereumNetworkContainers []EthereumNetworkContainer

func (e *EthereumNetworkContainers) add(container EthereumNetworkContainer) {
	*e = append(*e, container)
}

func (e *EthereumNetworkContainers) wasAlreadyRestarted(containerName string) bool {
	for _, container := range *e {
		if container.ContainerName == containerName {
			return true
		}
	}
	return false
}

var restartedContainers = make(EthereumNetworkContainers, 0)

func (b *EthereumNetwork) setExistingContainerName(ct ContainerType) EnvComponentOption {
	if !b.isRecreated {
		return func(c *EnvComponent) {}
	}

	// in that way we can support restarting of multiple nodes out of the box
	for _, container := range b.Containers {
		if container.ContainerType == ct && !restartedContainers.wasAlreadyRestarted(container.ContainerName) {
			restartedContainers.add(container)
			return func(c *EnvComponent) {
				if container.ContainerName != "" {
					c.ContainerName = container.ContainerName
				}
			}
		}
	}

	return func(c *EnvComponent) {}
}

func createHostDirectories() (string, string, error) {
	executionDir, err := os.MkdirTemp("", "execution")
	if err != nil {
		return "", "", err
	}

	consensusDir, err := os.MkdirTemp("", "consensus")
	if err != nil {
		return "", "", err
	}

	return executionDir, consensusDir, nil
}
