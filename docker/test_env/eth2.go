package test_env

import (
	"context"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/docker"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	tcwait "github.com/testcontainers/testcontainers-go/wait"
)

type ConsensusType = string

const (
	ConsensusType_PoS ConsensusType = "pos"
	ConsensusType_PoW ConsensusType = "pow"
)

type ConsensusLayer = string

const (
	ConsensusLayer_Prysm ConsensusLayer = "prysm"
)

type Eth2Components struct {
	Geth        *Geth2
	BeaconChain *BeaconChain
	Validator   *Validator
}

func StartEth2(t *testing.T, c ConsensusLayer) (blockchain.EVMNetwork, Eth2Components, error) {
	l := logging.GetTestLogger(t)
	network, err := docker.CreateNetwork(l)
	if err != nil {
		return blockchain.EVMNetwork{}, Eth2Components{}, err
	}

	bg := NewBeaconChainGenesis([]string{network.Name}).
		WithTestLogger(t)
	err = bg.StartContainer()
	if err != nil {
		return blockchain.EVMNetwork{}, Eth2Components{}, err
	}

	gg := NewGethGenesis([]string{network.Name}, bg.ExecutionDir).WithTestLogger(t)
	err = gg.StartContainer()
	if err != nil {
		return blockchain.EVMNetwork{}, Eth2Components{}, err
	}

	geth := NewGeth2([]string{network.Name}, bg.ExecutionDir).WithTestLogger(t)
	net, _, err := geth.StartContainer()
	if err != nil {
		return blockchain.EVMNetwork{}, Eth2Components{}, err
	}

	beacon := NewBeaconChain([]string{network.Name}, bg.ExecutionDir, bg.ConsensusDir, geth.ExecutionURL).WithTestLogger(t)
	err = beacon.StartContainer()
	if err != nil {
		return blockchain.EVMNetwork{}, Eth2Components{}, err
	}

	validator := NewValidator([]string{network.Name}, bg.ConsensusDir, beacon.InternalRpcURL).WithTestLogger(t)
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
