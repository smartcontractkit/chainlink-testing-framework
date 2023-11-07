package eth2

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/docker"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	tcwait "github.com/testcontainers/testcontainers-go/wait"
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
	t              *testing.T
	l              zerolog.Logger
	consensusType  *ConsensusType
	consensusLayer *ConsensusLayer
	consensusNodes int
	executionLayer *ExecutionLayer
	executionNodes int
}

type Eth2Components struct {
	Geth        *Geth2
	BeaconChain *BeaconChain
	Validator   *Validator
}

func NewEthereumNetworkBuilder(t *testing.T) *EthereumNetworkBuilder {
	return &EthereumNetworkBuilder{
		t: t,
		l: logging.GetTestLogger(t),
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

func (b *EthereumNetworkBuilder) Start() (blockchain.EVMNetwork, Eth2Components, error) {
	err := b.validate()
	if err != nil {
		return blockchain.EVMNetwork{}, Eth2Components{}, err
	}

	switch *b.consensusType {
	case ConsensusType_PoS:
		return b.startPos()
	case ConsensusType_PoW:
		return blockchain.EVMNetwork{}, Eth2Components{}, errors.New("PoW is not yet supported")
	default:
		return blockchain.EVMNetwork{}, Eth2Components{}, errors.New("unknown consensus type")
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

	return nil
}

func (b *EthereumNetworkBuilder) startPos() (blockchain.EVMNetwork, Eth2Components, error) {
	network, err := docker.CreateNetwork(b.l)
	if err != nil {
		return blockchain.EVMNetwork{}, Eth2Components{}, err
	}

	bg := NewEth2Genesis([]string{network.Name}).
		WithTestLogger(b.t)
	err = bg.StartContainer()
	if err != nil {
		return blockchain.EVMNetwork{}, Eth2Components{}, err
	}

	gg := NewEth1Genesis([]string{network.Name}, bg.ExecutionDir).WithTestLogger(b.t)
	err = gg.StartContainer()
	if err != nil {
		return blockchain.EVMNetwork{}, Eth2Components{}, err
	}

	geth := NewGeth2([]string{network.Name}, bg.ExecutionDir).WithTestLogger(b.t)
	net, _, err := geth.StartContainer()
	if err != nil {
		return blockchain.EVMNetwork{}, Eth2Components{}, err
	}

	beacon := NewBeaconChain([]string{network.Name}, bg.ExecutionDir, bg.ConsensusDir, geth.ExecutionURL).WithTestLogger(b.t)
	err = beacon.StartContainer()
	if err != nil {
		return blockchain.EVMNetwork{}, Eth2Components{}, err
	}

	validator := NewValidator([]string{network.Name}, bg.ConsensusDir, beacon.InternalRpcURL).WithTestLogger(b.t)
	err = validator.StartContainer()
	if err != nil {
		return blockchain.EVMNetwork{}, Eth2Components{}, err
	}

	waitForFirstBlock := tcwait.NewLogStrategy("Chain head was updated").WithPollInterval(1 * time.Second).WithStartupTimeout(60 * time.Second)
	err = waitForFirstBlock.WaitUntilReady(context.Background(), geth.Container)
	if err != nil {
		return blockchain.EVMNetwork{}, Eth2Components{}, err
	}

	return net, Eth2Components{
		Geth:        geth,
		BeaconChain: beacon,
		Validator:   validator,
	}, nil
}
