package blockchain

import (
	"context"
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
	TypeTon         = "ton"
	TypeCanton      = "canton"
	TypeStellar     = "stellar"
)

// Blockchain node family
const (
	FamilyEVM     = "evm"
	FamilySolana  = "solana"
	FamilyAptos   = "aptos"
	FamilySui     = "sui"
	FamilyTron    = "tron"
	FamilyTon     = "ton"
	FamilyCanton  = "canton"
	FamilyStellar = "stellar"
)

// Input is a blockchain network configuration params
type Input struct {
	// Common EVM fields
	Type          string `toml:"type" validate:"required,oneof=anvil geth besu solana aptos tron sui ton canton stellar" envconfig:"net_type" comment:"Type can be one of: anvil geth besu solana aptos tron sui ton canton stellar, this struct describes common configuration we are using across all blockchains"`
	Image         string `toml:"image" comment:"Blockchain node image in format: $registry:$image, ex.: ghcr.io/foundry-rs/foundry:stable"`
	PullImage     bool   `toml:"pull_image" comment:"Whether to pull image or not when creating Docker container"`
	Port          string `toml:"port" comment:"The port Docker container will expose"`
	ContainerName string `toml:"container_name" comment:"Docker container name"`
	// Not applicable to Solana, ws port for Solana is +1 of port
	WSPort                   string   `toml:"port_ws" comment:"WebSocket port container will expose"`
	ChainID                  string   `toml:"chain_id" comment:"Blockchain chain ID"`
	DockerCmdParamsOverrides []string `toml:"docker_cmd_params" comment:"Docker command parameters override, ex. for Anvil: [\"-b\", \"1\", \"--mixed-mining\"]"`
	Out                      *Output  `toml:"out" comment:"blockchain deployment output"`

	// Solana fields
	// publickey to mint when solana-test-validator starts
	PublicKey    string `toml:"public_key" comment:"Public key to mint when solana-test-validator starts"`
	ContractsDir string `toml:"contracts_dir" comment:"Solana's contracts directory"`
	// programs to deploy on solana-test-validator start
	// a map of program name to program id
	// there needs to be a matching .so file in contracts_dir
	SolanaPrograms     map[string]string             `toml:"solana_programs" comment:"Solana's programs, map of program name to program ID, there needs to be a matching .so file in contracts_dir"`
	ContainerResources *framework.ContainerResources `toml:"resources" comment:"Docker container resources"`
	CustomPorts        []string                      `toml:"custom_ports" comment:"Custom ports pairs in format $host_port_number:$docker_port_number"`

	// Sui specific: faucet port for funding accounts
	FaucetPort string `toml:"faucet_port" comment:"Sui blockchain network faucet port"`

	// Canton specific
	NumberOfCantonValidators int  `toml:"number_of_canton_validators" comment:"Number of Canton network validators"`
	EnableSplice             bool `toml:"enable_splice" comment:"Whether to enable Splice service for Canton network (default: false). Splice is only needed when interacting with CC (Cross-Chain) features."`

	// GAPv2 specific params
	HostNetworkMode  bool   `toml:"host_network_mode" comment:"GAPv2 specific paramter: host netowork mode, if 'true' will run environment in host network mode"`
	CertificatesPath string `toml:"certificates_path" comment:"GAPv2 specific parameter: path to default Ubuntu's certificates"`

	// Optional params
	ImagePlatform *string `toml:"image_platform" comment:"Docker image platform, default is 'linux/amd64'"`
	// Custom environment variables for the container
	CustomEnv map[string]string `toml:"custom_env" comment:"Docker container environment variables in TOML format key = value"`
}

// Output is a blockchain network output, ChainID and one or more nodes that forms the network
type Output struct {
	UseCache            bool                     `toml:"use_cache" comment:"Whether to respect caching or not, if cache = true component won't be deployed again"`
	Type                string                   `toml:"type" comment:"Type can be one of: anvil geth besu solana aptos tron sui ton canton stellar, this struct describes common configuration we are using across all blockchains"`
	Family              string                   `toml:"family" comment:"Blockchain family, can be one of: evm solana aptos sui tron ton canton stellar"`
	ContainerName       string                   `toml:"container_name" comment:"Blockchain Docker container name"`
	NetworkSpecificData *NetworkSpecificData     `toml:"network_specific_data" comment:"Blockchain network-specific data"`
	Container           testcontainers.Container `toml:"-"`
	ChainID             string                   `toml:"chain_id" comment:"Chain ID"`
	Nodes               []*Node                  `toml:"nodes" comment:"Blockchain nodes info"`
}

type NetworkSpecificData struct {
	SuiAccount      *SuiWalletInfo      `toml:"sui_account" comment:"Sui network account info"`
	CantonEndpoints *CantonEndpoints    `toml:"canton_endpoints" comment:"Canton network endpoints info"`
	StellarNetwork  *StellarNetworkInfo `toml:"stellar_network" comment:"Stellar network info"`
}

// Node represents blockchain node output, URLs required for connection locally and inside docker network
type Node struct {
	ExternalWSUrl   string `toml:"ws_url" comment:"External blockchain node WebSocket URL"`
	ExternalHTTPUrl string `toml:"http_url" comment:"External blockchain node HTTP URL"`
	InternalWSUrl   string `toml:"internal_ws_url" comment:"Internal blockchain node WebSocket URL"`
	InternalHTTPUrl string `toml:"internal_http_url" comment:"Internal blockchain node HTTP URL"`
}

func NewBlockchainNetwork(in *Input) (*Output, error) {
	// pass context to input if needed in the future
	return NewWithContext(context.Background(), in)
}

// NewBlockchainNetwork this is an abstraction that can spin up various blockchain network simulators
func NewWithContext(ctx context.Context, in *Input) (*Output, error) {
	var out *Output
	var err error
	switch in.Type {
	case TypeAnvil:
		out, err = newAnvil(ctx, in)
	case TypeGeth:
		out, err = newGeth(ctx, in)
	case TypeBesu:
		out, err = newBesu(ctx, in)
	case TypeSolana:
		out, err = newSolana(ctx, in)
	case TypeAptos:
		out, err = newAptos(ctx, in)
	case TypeSui:
		out, err = newSui(ctx, in)
	case TypeTron:
		out, err = newTron(ctx, in)
	case TypeAnvilZKSync:
		out, err = newAnvilZksync(ctx, in)
	case TypeTon:
		out, err = newTon(ctx, in)
	case TypeCanton:
		out, err = newCanton(ctx, in)
	case TypeStellar:
		out, err = newStellar(ctx, in)
	default:
		return nil, fmt.Errorf("blockchain type is not supported or empty, must be 'anvil' or 'geth'")
	}
	if err != nil {
		return nil, err
	}
	in.Out = out
	return out, nil
}

type ChainFamily string

func TypeToFamily(t string) (ChainFamily, error) {
	switch t {
	case TypeAnvil, TypeGeth, TypeBesu, TypeAnvilZKSync:
		return ChainFamily(FamilyEVM), nil
	case TypeSolana:
		return ChainFamily(FamilySolana), nil
	case TypeAptos:
		return ChainFamily(FamilyAptos), nil
	case TypeSui:
		return ChainFamily(FamilySui), nil
	case TypeTron:
		return ChainFamily(FamilyTron), nil
	case TypeTon:
		return ChainFamily(FamilyTon), nil
	case TypeCanton:
		return ChainFamily(FamilyCanton), nil
	case TypeStellar:
		return ChainFamily(FamilyStellar), nil
	default:
		return "", fmt.Errorf("blockchain type is not supported or empty: %s", t)
	}
}
