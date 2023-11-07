package test_env

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/rs/zerolog"
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
	BeaconChainConfig *BeaconChainConfig
	validated         bool
}

type Eth2Components struct {
	Geth        *Geth2
	BeaconChain *PrysmBeaconChain
	Validator   *PrysmValidator
}

func NewEthereumNetworkBuilder(t *testing.T) *EthereumNetworkBuilder {
	return &EthereumNetworkBuilder{
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
	b.BeaconChainConfig = &config
	return b
}

func (b *EthereumNetworkBuilder) UsingDockerNetworks(networks []string) *EthereumNetworkBuilder {
	b.dockerNetworks = networks
	return b
}

func (b *EthereumNetworkBuilder) Build() error {
	err := b.validate()
	if err != nil {
		return err
	}
	b.validated = true
	return nil
}

func (b *EthereumNetworkBuilder) Start() (blockchain.EVMNetwork, RpcProvider, error) {
	if !b.validated {
		return blockchain.EVMNetwork{}, RpcProvider{}, errors.New("builder must be build and validated before starting. Execute builder.Build()")
	}

	switch *b.consensusType {
	case ConsensusType_PoS:
		return b.startPos()
	case ConsensusType_PoW:
		return b.startPow()
	default:
		return blockchain.EVMNetwork{}, RpcProvider{}, fmt.Errorf("unknown consensus type: %s", *b.consensusType)
	}
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

	if *b.consensusType == ConsensusType_PoW {
		b.consensusNodes = 0
	}

	return nil
}

func (b *EthereumNetworkBuilder) startPos() (blockchain.EVMNetwork, RpcProvider, error) {
	if *b.consensusLayer != ConsensusLayer_Prysm {
		return blockchain.EVMNetwork{}, RpcProvider{}, fmt.Errorf("unsupported consensus layer: %s", *b.consensusLayer)
	}
	networkNames, err := b.getOrCreateDockerNetworks()
	if err != nil {
		return blockchain.EVMNetwork{}, RpcProvider{}, err
	}

	var beaconChainConfig BeaconChainConfig
	if b.BeaconChainConfig != nil {
		beaconChainConfig = *b.BeaconChainConfig
	} else {
		beaconChainConfig = DefaultBeaconChainConfig
	}

	bg := NewEth2Genesis(networkNames, beaconChainConfig).
		WithTestLogger(b.t)
	err = bg.StartContainer()
	if err != nil {
		return blockchain.EVMNetwork{}, RpcProvider{}, err
	}

	gg := NewEth1Genesis(networkNames, bg.hostExecutionDir).WithTestLogger(b.t)
	err = gg.StartContainer()
	if err != nil {
		return blockchain.EVMNetwork{}, RpcProvider{}, err
	}

	geth2 := NewGeth2(networkNames, bg.hostExecutionDir).WithTestLogger(b.t)
	net, _, err := geth2.StartContainer()
	if err != nil {
		return blockchain.EVMNetwork{}, RpcProvider{}, err
	}

	beacon := NewPrysmBeaconChain(networkNames, bg.hostExecutionDir, bg.hostConsensusDir, geth2.InternalExecutionURL).WithTestLogger(b.t)
	err = beacon.StartContainer()
	if err != nil {
		return blockchain.EVMNetwork{}, RpcProvider{}, err
	}

	validator := NewPrysmValidator(networkNames, bg.hostConsensusDir, beacon.InternalBeaconRpcProvider).WithTestLogger(b.t)
	err = validator.StartContainer()
	if err != nil {
		return blockchain.EVMNetwork{}, RpcProvider{}, err
	}

	waitForFirstBlock := tcwait.NewLogStrategy("Chain head was updated").WithPollInterval(1 * time.Second).WithStartupTimeout(60 * time.Second)
	err = waitForFirstBlock.WaitUntilReady(context.Background(), geth2.Container)
	if err != nil {
		return blockchain.EVMNetwork{}, RpcProvider{}, err
	}

	return net, RpcProvider{
		privateHttpUrls: []string{geth2.InternalHttpUrl},
		privatelWsUrls:  []string{geth2.InternalWsUrl},
		publiclHttpUrls: []string{geth2.ExternalHttpUrl},
		publicsUrls:     []string{geth2.ExternalWsUrl},
	}, nil
}

func (b *EthereumNetworkBuilder) startPow() (blockchain.EVMNetwork, RpcProvider, error) {
	if *b.executionLayer != ExecutionLayer_Geth {
		return blockchain.EVMNetwork{}, RpcProvider{}, fmt.Errorf("unsupported execution layer: %s", *b.executionLayer)
	}
	networkNames, err := b.getOrCreateDockerNetworks()
	if err != nil {
		return blockchain.EVMNetwork{}, RpcProvider{}, err
	}

	geth := NewGeth(networkNames).WithTestLogger(b.t)
	net, docker, err := geth.StartContainer()
	if err != nil {
		return blockchain.EVMNetwork{}, RpcProvider{}, err
	}

	return net, RpcProvider{
		privateHttpUrls: []string{docker.HttpUrl},
		privatelWsUrls:  []string{docker.WsUrl},
		publiclHttpUrls: []string{geth.ExternalHttpUrl},
		publicsUrls:     []string{geth.ExternalWsUrl},
	}, nil
}

func (b *EthereumNetworkBuilder) getOrCreateDockerNetworks() ([]string, error) {
	var networkNames []string

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

func (b *EthereumNetworkBuilder) Describe() string {
	return fmt.Sprintf("consensus type: %s, consensus layer: %s, execution layer: %s, consensus nodes: %d, execution nodes: %d",
		*b.consensusType, *b.consensusLayer, *b.executionLayer, b.consensusNodes, b.executionNodes)
}
