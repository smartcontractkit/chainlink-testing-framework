package blockchain

import (
	"fmt"

	"github.com/testcontainers/testcontainers-go"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

// Blockchain node type
const (
	TypeAnvil       = "anvil"
	TypeAnvilZKSync = "anvil-zksync"
	TypeGeth        = "geth"
	TypeBesu        = "besu"
	TypeSolana      = "solana"
	TypeAptos       = "aptos"
	TypeSui         = "sui"
	TypeTron        = "tron"
)

// Blockchain node family
const (
	FamilyEVM    = "evm"
	FamilySolana = "solana"
	FamilyAptos  = "aptos"
	FamilySui    = "sui"
	FamilyTron   = "tron"
)

// Input is a blockchain network configuration params
type Input struct {
	// Common EVM fields
	Type      string `toml:"type" validate:"required,oneof=anvil geth besu solana aptos tron sui" envconfig:"net_type"`
	Image     string `toml:"image"`
	PullImage bool   `toml:"pull_image"`
	Port      string `toml:"port"`
	// Not applicable to Solana, ws port for Solana is +1 of port
	WSPort                   string   `toml:"port_ws"`
	ChainID                  string   `toml:"chain_id"`
	DockerCmdParamsOverrides []string `toml:"docker_cmd_params"`
	Out                      *Output  `toml:"out"`

	// Solana fields
	// publickey to mint when solana-test-validator starts
	PublicKey    string `toml:"public_key"`
	ContractsDir string `toml:"contracts_dir"`
	// programs to deploy on solana-test-validator start
	// a map of program name to program id
	// there needs to be a matching .so file in contracts_dir
	SolanaPrograms     map[string]string             `toml:"solana_programs"`
	ContainerResources *framework.ContainerResources `toml:"resources"`
	CustomPorts        []string                      `toml:"custom_ports"`
}

// Output is a blockchain network output, ChainID and one or more nodes that forms the network
type Output struct {
	UseCache            bool                     `toml:"use_cache"`
	Type                string                   `toml:"type"`
	Family              string                   `toml:"family"`
	ContainerName       string                   `toml:"container_name"`
	NetworkSpecificData *NetworkSpecificData     `toml:"network_specific_data"`
	Container           testcontainers.Container `toml:"-"`
	ChainID             string                   `toml:"chain_id"`
	Nodes               []*Node                  `toml:"nodes"`
}

type NetworkSpecificData struct {
	SuiAccount *SuiWalletInfo
}

// Node represents blockchain node output, URLs required for connection locally and inside docker network
type Node struct {
	ExternalWSUrl   string `toml:"ws_url"`
	ExternalHTTPUrl string `toml:"http_url"`
	InternalWSUrl   string `toml:"internal_ws_url"`
	InternalHTTPUrl string `toml:"internal_http_url"`
}

// NewBlockchainNetwork this is an abstraction that can spin up various blockchain network simulators
func NewBlockchainNetwork(in *Input) (*Output, error) {
	var out *Output
	var err error
	switch in.Type {
	case TypeAnvil:
		out, err = newAnvil(in)
	case TypeGeth:
		out, err = newGeth(in)
	case TypeBesu:
		out, err = newBesu(in)
	case TypeSolana:
		out, err = newSolana(in)
	case TypeAptos:
		out, err = newAptos(in)
	case TypeSui:
		out, err = newSui(in)
	case TypeTron:
		out, err = newTron(in)
	case TypeAnvilZKSync:
		out, err = newAnvilZksync(in)
	default:
		return nil, fmt.Errorf("blockchain type is not supported or empty, must be 'anvil' or 'geth'")
	}
	if err != nil {
		return nil, err
	}
	in.Out = out
	return out, nil
}
