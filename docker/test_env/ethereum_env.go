package test_env

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
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
)

type ConsensusLayer string

const (
	ConsensusLayer_Prysm ConsensusLayer = "prysm"
)

type EthereumNetworkBuilder struct {
	t                   *testing.T
	dockerNetworks      []string
	consensusType       *ConsensusType
	ehtereumChainConfig *EthereumChainConfig
	existingConfig      *EthereumNetwork
	addressesToFund     []string
	participants        []EthereumNetworkParticipant
}

type EthereumNetworkParticipant struct {
	ConsensusLayer ConsensusLayer `json:"consensus_layer"` //nil means PoW
	ExecutionLayer ExecutionLayer `json:"execution_layer"`
	Count          int            `json:"count"`
}

func NewEthereumNetworkBuilder() EthereumNetworkBuilder {
	return EthereumNetworkBuilder{
		dockerNetworks: []string{},
		participants:   []EthereumNetworkParticipant{},
	}
}

func (b *EthereumNetworkBuilder) WithConsensusType(consensusType ConsensusType) *EthereumNetworkBuilder {
	b.consensusType = &consensusType
	return b
}

func (b *EthereumNetworkBuilder) WithDefaultNetworkParticipants(consensusType ConsensusType) *EthereumNetworkBuilder {
	b.consensusType = &consensusType
	switch consensusType {
	case ConsensusType_PoS:
		consensusLayer := ConsensusLayer_Prysm

		b.participants = []EthereumNetworkParticipant{
			{
				ConsensusLayer: consensusLayer,
				ExecutionLayer: ExecutionLayer_Geth,
				Count:          1,
			},
		}
	case ConsensusType_PoW:
		b.participants = []EthereumNetworkParticipant{
			{
				ConsensusLayer: "",
				ExecutionLayer: ExecutionLayer_Geth,
				Count:          1,
			},
		}
	default:
		panic(fmt.Sprintf("unknown consensus type: %s", consensusType))
	}

	return b
}

func (b *EthereumNetworkBuilder) WithCustomNetworkParticipants(participants []EthereumNetworkParticipant) *EthereumNetworkBuilder {
	if len(participants) != 1 {
		panic("only one participant is currently supported")
	}

	b.participants = participants
	return b
}

func (b *EthereumNetworkBuilder) WithEthereumChainConfig(config EthereumChainConfig) *EthereumNetworkBuilder {
	b.ehtereumChainConfig = &config
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

func (b *EthereumNetworkBuilder) buildNetworkConfig() EthereumNetwork {
	n := EthereumNetwork{
		ConsensusType: *b.consensusType,
		Participants:  b.participants,
	}

	if b.existingConfig != nil {
		n.isRecreated = true
		n.Containers = b.existingConfig.Containers
	}

	n.ehtereumChainConfig = b.ehtereumChainConfig
	n.t = b.t

	return n
}

func (b *EthereumNetworkBuilder) Build() (EthereumNetwork, error) {
	b.importExistingConfig()
	if b.ehtereumChainConfig == nil {
		defaultConfig := GetDefaultChainConfig()
		b.ehtereumChainConfig = &defaultConfig
	} else {
		b.ehtereumChainConfig.fillInMissingValuesWithDefault()
	}

	b.ehtereumChainConfig.GenerateGenesisTimestamp()
	err := b.validate()
	if err != nil {
		return EthereumNetwork{}, err
	}

	return b.buildNetworkConfig(), nil
}

func (b *EthereumNetworkBuilder) importExistingConfig() {
	if b.existingConfig == nil {
		return
	}

	b.participants = b.existingConfig.Participants
	b.consensusType = &b.existingConfig.ConsensusType
	b.dockerNetworks = b.existingConfig.DockerNetworkNames
	b.ehtereumChainConfig = b.existingConfig.ehtereumChainConfig
}

func (b *EthereumNetworkBuilder) validate() error {
	if b.consensusType == nil {
		return errors.New("consensus type is required")
	}

	if len(b.participants) != 1 {
		return errors.New("only one participant is currently supported")
	}

	for _, p := range b.participants {
		if p.ExecutionLayer == "" {
			return errors.New("execution layer is required")
		}

		if *b.consensusType == ConsensusType_PoS && p.ConsensusLayer == "" {
			return errors.New("consensus layer is required for PoS")
		}

		if *b.consensusType == ConsensusType_PoW && p.ConsensusLayer != "" {
			return errors.New("consensus layer is not allowed for PoW")
		}
	}

	//TODO when we support multiple participants, we need to validate that all of them are either PoW or PoS

	for _, addr := range b.addressesToFund {
		if !common.IsHexAddress(addr) {
			return fmt.Errorf("address %s is not a valid hex address", addr)
		}
	}

	err := b.ehtereumChainConfig.Validate(logging.GetTestLogger(b.t))
	if err != nil {
		return err
	}

	return nil
}

type EthereumNetwork struct {
	ConsensusType       ConsensusType                `json:"consensus_type"`
	DockerNetworkNames  []string                     `json:"docker_network_names"`
	Containers          EthereumNetworkContainers    `json:"containers"`
	Participants        []EthereumNetworkParticipant `json:"participants"`
	isRecreated         bool
	ehtereumChainConfig *EthereumChainConfig
	t                   *testing.T
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

	for _, p := range en.Participants {
		if p.ConsensusLayer != ConsensusLayer_Prysm {
			return blockchain.EVMNetwork{}, RpcProvider{}, fmt.Errorf("unsupported consensus layer: %s", p.ConsensusLayer)
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

			if err != nil {
				return blockchain.EVMNetwork{}, RpcProvider{}, err
			}

			valKeysGeneretor := NewValKeysGeneretor(en.ehtereumChainConfig, valKeysDir).WithTestInstance(en.t)
			err = valKeysGeneretor.StartContainer()
			if err != nil {
				return blockchain.EVMNetwork{}, RpcProvider{}, err
			}

			genesis := NewEthGenesisGenerator(*en.ehtereumChainConfig, generatedDataHostDir).WithTestInstance(en.t)
			err = genesis.StartContainer()
			if err != nil {
				return blockchain.EVMNetwork{}, RpcProvider{}, err
			}

			initHelper := NewInitHelper(*en.ehtereumChainConfig, generatedDataHostDir).WithTestInstance(en.t)
			err = initHelper.StartContainer()
			if err != nil {
				return blockchain.EVMNetwork{}, RpcProvider{}, err
			}
		} else {
			//TODO set to actual values, even if they do not matter for containers that are already running
			generatedDataHostDir = ""
			valKeysDir = ""
		}

		var client ExecutionClient
		switch p.ExecutionLayer {
		case ExecutionLayer_Geth:
			client = NewGeth2(singleNetwork, en.ehtereumChainConfig, generatedDataHostDir, ConsensusLayer_Prysm, en.setExistingContainerName(ContainerType_Geth2)).WithTestInstance(en.t)
		case ExecutionLayer_Nethermind:
			client = NewNethermind(singleNetwork, generatedDataHostDir, ConsensusLayer_Prysm, en.setExistingContainerName(ContainerType_Nethermind)).WithTestInstance(en.t)
		default:
			return blockchain.EVMNetwork{}, RpcProvider{}, fmt.Errorf("unsupported execution layer: %s", p.ExecutionLayer)
		}

		net, err = client.StartContainer()
		if err != nil {
			return blockchain.EVMNetwork{}, RpcProvider{}, err
		}

		beacon := NewPrysmBeaconChain(singleNetwork, en.ehtereumChainConfig, generatedDataHostDir, client.GetInternalExecutionURL(), en.setExistingContainerName(ContainerType_PrysmBeacon)).WithTestInstance(en.t)
		err = beacon.StartContainer()
		if err != nil {
			return blockchain.EVMNetwork{}, RpcProvider{}, err
		}

		validator := NewPrysmValidator(singleNetwork, en.ehtereumChainConfig, generatedDataHostDir, valKeysDir, beacon.InternalBeaconRpcProvider, en.setExistingContainerName(ContainerType_PrysmVal)).WithTestInstance(en.t)
		err = validator.StartContainer()
		if err != nil {
			return blockchain.EVMNetwork{}, RpcProvider{}, err
		}

		err = client.WaitUntilChainIsReady(en.ehtereumChainConfig.GetDefaultWaitDuration())
		if err != nil {
			return blockchain.EVMNetwork{}, RpcProvider{}, err
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
	}

	en.DockerNetworkNames = networkNames
	//TODO when we support multiple participants, we need to modify net so that it contains all the RPC URLs, not just the last one

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

	for _, p := range en.Participants {
		if p.ExecutionLayer != ExecutionLayer_Geth {
			return blockchain.EVMNetwork{}, RpcProvider{}, fmt.Errorf("unsupported execution layer: %s", p.ExecutionLayer)
		}
		singleNetwork, err := en.getOrCreateDockerNetworks()
		if err != nil {
			return blockchain.EVMNetwork{}, RpcProvider{}, err
		}

		geth := NewGeth(singleNetwork, *en.ehtereumChainConfig, en.setExistingContainerName(ContainerType_Geth)).WithTestLogger(en.t)
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
	}

	en.DockerNetworkNames = networkNames
	//TODO when we support multiple participants, we need to modify net so that it contains all the RPC URLs, not just the last one

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
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("consensus type: %s", en.ConsensusType))
	sb.WriteString("\nparticipants: {")
	for _, p := range en.Participants {
		sb.WriteString(fmt.Sprintf("\n\tconsensus layer: %s, execution layer: %s, count: %d", p.ConsensusLayer, p.ExecutionLayer, p.Count))
	}
	sb.WriteString("\n}")
	return sb.String()
}

func (en *EthereumNetwork) setExistingContainerName(ct ContainerType) EnvComponentOption {
	if !en.isRecreated {
		return func(c *EnvComponent) {}
	}

	// in that way we can support restarting of multiple nodes out of the box
	for _, container := range en.Containers {
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

// maybe only store ports here and depending on where the test is executed return different URLs?
// maybe 3 different constructors for each "perspective"?
// also it could expose 2 iterators:
// 1. that iterates until it has something to return
// 2. that iterates in a loop and always returns something
// why? because then client could decide not to care about how many RPCs there are and just be fine
// with any, even if all calls return the same RPC... and if there were more, then each node could
// use a different one, but the code for calling the provider would be the same
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

var restartedContainers = make(EthereumNetworkContainers, 0)

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
