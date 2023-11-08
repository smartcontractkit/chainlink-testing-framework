package test_env

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

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
	ExecutionLayer_Geth ExecutionLayer = "geth"
)

type ConsensusLayer string

const (
	ConsensusLayer_Prysm ConsensusLayer = "prysm"
)

type EthereumNetworkBuilder struct {
	t                 *testing.T
	l                 zerolog.Logger
	dockerNetworks    []string
	consensusType     *ConsensusType
	consensusLayer    *ConsensusLayer
	consensusNodes    int
	executionLayer    *ExecutionLayer
	executionNodes    int
	beaconChainConfig *BeaconChainConfig
	validated         bool
	existingConfig    *EthereumNetworkConfig
}

type Eth2Components struct {
	Geth        *Geth2
	BeaconChain *PrysmBeaconChain
	Validator   *PrysmValidator
}

func NewEthereumNetworkBuilder(t *testing.T) EthereumNetworkBuilder {
	return EthereumNetworkBuilder{
		t:              t,
		l:              logging.GetTestLogger(t),
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

func (b *EthereumNetworkBuilder) UsingDockerNetworks(networks []string) *EthereumNetworkBuilder {
	b.dockerNetworks = networks
	return b
}

func (b *EthereumNetworkBuilder) WithExistingConfig(config EthereumNetworkConfig) *EthereumNetworkBuilder {
	b.existingConfig = &config
	return b
}

func (b *EthereumNetworkBuilder) Build() error {
	b.importExistingConfig()
	err := b.validate()
	if err != nil {
		return err
	}
	b.validated = true
	return nil
}

func (b *EthereumNetworkBuilder) Start() (blockchain.EVMNetwork, RpcProvider, EthereumNetworkConfig, error) {
	if !b.validated {
		return blockchain.EVMNetwork{}, RpcProvider{}, EthereumNetworkConfig{}, errors.New("builder must be build and validated before starting. Execute builder.Build()")
	}

	switch *b.consensusType {
	case ConsensusType_PoS:
		return b.startPos()
	case ConsensusType_PoW:
		return b.startPow()
	default:
		return blockchain.EVMNetwork{}, RpcProvider{}, EthereumNetworkConfig{}, fmt.Errorf("unknown consensus type: %s", *b.consensusType)
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

	return nil
}

func (b *EthereumNetworkBuilder) startPos() (blockchain.EVMNetwork, RpcProvider, EthereumNetworkConfig, error) {
	if *b.consensusLayer != ConsensusLayer_Prysm {
		return blockchain.EVMNetwork{}, RpcProvider{}, EthereumNetworkConfig{}, fmt.Errorf("unsupported consensus layer: %s", *b.consensusLayer)
	}
	networkNames, err := b.getOrCreateDockerNetworks()
	if err != nil {
		return blockchain.EVMNetwork{}, RpcProvider{}, EthereumNetworkConfig{}, err
	}

	var hostExecutionDir, hostConsensusDir string

	// create host directories and run genesis containers only if we are NOT recreating existing containers
	if b.existingConfig == nil {
		hostExecutionDir, hostConsensusDir, err = createHostDirectories()
		if err != nil {
			return blockchain.EVMNetwork{}, RpcProvider{}, EthereumNetworkConfig{}, err
		}

		var beaconChainConfig BeaconChainConfig
		if b.beaconChainConfig != nil {
			beaconChainConfig = *b.beaconChainConfig
		} else {
			beaconChainConfig = DefaultBeaconChainConfig
		}

		bg := NewEth2Genesis(networkNames, beaconChainConfig, hostExecutionDir, hostConsensusDir).
			WithTestLogger(b.t)
		err = bg.StartContainer()
		if err != nil {
			return blockchain.EVMNetwork{}, RpcProvider{}, EthereumNetworkConfig{}, err
		}

		gg := NewEth1Genesis(networkNames, hostExecutionDir).WithTestLogger(b.t)
		err = gg.StartContainer()
		if err != nil {
			return blockchain.EVMNetwork{}, RpcProvider{}, EthereumNetworkConfig{}, err
		}
	}

	geth2 := NewGeth2(networkNames, hostExecutionDir, ConsensusLayer_Prysm, b.setExistingContainerName(ContainerType_Geth2)).WithTestLogger(b.t)
	net, err := geth2.StartContainer()
	if err != nil {
		return blockchain.EVMNetwork{}, RpcProvider{}, EthereumNetworkConfig{}, err
	}

	beacon := NewPrysmBeaconChain(networkNames, hostExecutionDir, hostConsensusDir, geth2.InternalExecutionURL, b.setExistingContainerName(ContainerType_PrysmBeacon)).WithTestLogger(b.t)
	err = beacon.StartContainer()
	if err != nil {
		return blockchain.EVMNetwork{}, RpcProvider{}, EthereumNetworkConfig{}, err
	}

	validator := NewPrysmValidator(networkNames, hostConsensusDir, beacon.InternalBeaconRpcProvider, b.setExistingContainerName(ContainerType_PrysmVal)).WithTestLogger(b.t)
	err = validator.StartContainer()
	if err != nil {
		return blockchain.EVMNetwork{}, RpcProvider{}, EthereumNetworkConfig{}, err
	}

	waitForFirstBlock := tcwait.NewLogStrategy("Chain head was updated").WithPollInterval(1 * time.Second).WithStartupTimeout(60 * time.Second)
	err = waitForFirstBlock.WaitUntilReady(context.Background(), geth2.Container)
	if err != nil {
		return blockchain.EVMNetwork{}, RpcProvider{}, EthereumNetworkConfig{}, err
	}

	rpcProvider := RpcProvider{
		privateHttpUrls: []string{geth2.InternalHttpUrl},
		privatelWsUrls:  []string{geth2.InternalWsUrl},
		publiclHttpUrls: []string{geth2.ExternalHttpUrl},
		publicsUrls:     []string{geth2.ExternalWsUrl},
	}

	ethNc := EthereumNetworkConfig{
		DockerNetworkNames: networkNames,
		ExecutionDir:       hostExecutionDir,
		ConsensusDir:       hostConsensusDir,
		Containers: EthereumNetworkContainers{
			{
				ContainerName: geth2.ContainerName,
				ContainerType: ContainerType_Geth2,
				Container:     &geth2.Container,
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
		},
	}

	b.saveUniversalDataInConfig(&ethNc)

	return net, rpcProvider, ethNc, nil
}

func (b *EthereumNetworkBuilder) startPow() (blockchain.EVMNetwork, RpcProvider, EthereumNetworkConfig, error) {
	if *b.executionLayer != ExecutionLayer_Geth {
		return blockchain.EVMNetwork{}, RpcProvider{}, EthereumNetworkConfig{}, fmt.Errorf("unsupported execution layer: %s", *b.executionLayer)
	}
	networkNames, err := b.getOrCreateDockerNetworks()
	if err != nil {
		return blockchain.EVMNetwork{}, RpcProvider{}, EthereumNetworkConfig{}, err
	}

	geth := NewGeth(networkNames, b.setExistingContainerName(ContainerType_Geth)).WithTestLogger(b.t)
	net, docker, err := geth.StartContainer()
	if err != nil {
		return blockchain.EVMNetwork{}, RpcProvider{}, EthereumNetworkConfig{}, err
	}

	ethNc := EthereumNetworkConfig{
		DockerNetworkNames: networkNames,
		Containers: EthereumNetworkContainers{
			{
				ContainerName: geth.ContainerName,
				ContainerType: ContainerType_Geth,
				Container:     &geth.Container,
			},
		},
	}

	b.saveUniversalDataInConfig(&ethNc)

	return net, RpcProvider{
		privateHttpUrls: []string{docker.HttpUrl},
		privatelWsUrls:  []string{docker.WsUrl},
		publiclHttpUrls: []string{geth.ExternalHttpUrl},
		publicsUrls:     []string{geth.ExternalWsUrl},
	}, ethNc, nil
}

func (b *EthereumNetworkBuilder) saveUniversalDataInConfig(c *EthereumNetworkConfig) {
	if b.consensusLayer != nil {
		c.ConsensusLayer = *b.consensusLayer
	} else {
		c.ConsensusLayer = ""
	}
	c.ConsensusNodes = b.consensusNodes
	c.ConsensusType = *b.consensusType

	c.ExecutionLayer = *b.executionLayer
	c.ExecutionNodes = b.executionNodes
}

func (b *EthereumNetworkBuilder) getOrCreateDockerNetworks() ([]string, error) {
	var networkNames []string

	if b.existingConfig != nil {
		return b.existingConfig.DockerNetworkNames, nil
	}

	if len(b.dockerNetworks) == 0 {
		network, err := docker.CreateNetwork(b.l)
		if err != nil {
			return networkNames, err
		}
		networkNames = []string{network.Name}
	} else {
		networkNames = b.dockerNetworks
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

func (b *EthereumNetworkConfig) Describe() string {
	return fmt.Sprintf("consensus type: %s, consensus layer: %s, execution layer: %s, consensus nodes: %d, execution nodes: %d",
		b.ConsensusType, b.ConsensusLayer, b.ExecutionLayer, b.ConsensusNodes, b.ExecutionNodes)
}

type ContainerType string

const (
	ContainerType_Geth        ContainerType = "geth"
	ContainerType_Geth2       ContainerType = "geth2"
	ContainerType_PrysmBeacon ContainerType = "prysm-beacon"
	ContainerType_PrysmVal    ContainerType = "prysm-validator"
)

type EthereumNetworkContainer struct {
	ContainerName string        `json:"container_name"`
	ContainerType ContainerType `json:"container_type"`
	Container     *tc.Container `json:"-"`
}

type EthereumNetworkConfig struct {
	ConsensusType      ConsensusType             `json:"consensus_type"`
	ConsensusLayer     ConsensusLayer            `json:"consensus_layer"`
	ConsensusNodes     int                       `json:"consensus_nodes"`
	ExecutionLayer     ExecutionLayer            `json:"execution_layer"`
	ExecutionNodes     int                       `json:"execution_nodes"`
	DockerNetworkNames []string                  `json:"docker_network_names"`
	ExecutionDir       string                    `json:"execution_dir"`
	ConsensusDir       string                    `json:"consensus_dir"`
	Containers         EthereumNetworkContainers `json:"containers"`
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

func (b *EthereumNetworkBuilder) setExistingContainerName(ct ContainerType) EnvComponentOption {
	if b.existingConfig == nil {
		return func(c *EnvComponent) {}
	}

	// in that way we can support restarting of multiple nodes out of the box
	for _, container := range b.existingConfig.Containers {
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
