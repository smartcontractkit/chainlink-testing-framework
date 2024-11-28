package blockchain

import (
	"fmt"
)

// Input is a blockchain network configuration params
type Input struct {
	Type                     string   `toml:"type" validate:"required,oneof=anvil geth besu" envconfig:"net_type"`
	Image                    string   `toml:"image"`
	PullImage                bool     `toml:"pull_image"`
	Port                     string   `toml:"port"`
	WSPort                   string   `toml:"port_ws"`
	ChainID                  string   `toml:"chain_id"`
	DockerCmdParamsOverrides []string `toml:"docker_cmd_params"`
	Out                      *Output  `toml:"out"`
}

// Output is a blockchain network output, ChainID and one or more nodes that forms the network
type Output struct {
	UseCache      bool    `toml:"use_cache"`
	Family        string  `toml:"family"`
	ContainerName string  `toml:"container_name"`
	ChainID       string  `toml:"chain_id"`
	Nodes         []*Node `toml:"nodes"`
}

// Node represents blockchain node output, URLs required for connection locally and inside docker network
type Node struct {
	HostWSUrl             string `toml:"ws_url"`
	HostHTTPUrl           string `toml:"http_url"`
	DockerInternalWSUrl   string `toml:"docker_internal_ws_url"`
	DockerInternalHTTPUrl string `toml:"docker_internal_http_url"`
}

// NewBlockchainNetwork this is an abstraction that can spin up various blockchain network simulators
// - Anvil
// - Geth
func NewBlockchainNetwork(in *Input) (*Output, error) {
	if in.Out != nil && in.Out.UseCache {
		return in.Out, nil
	}
	var out *Output
	var err error
	switch in.Type {
	case "anvil":
		out, err = newAnvil(in)
	case "geth":
		out, err = newGeth(in)
	case "besu":
		out, err = newBesu(in)
	default:
		return nil, fmt.Errorf("blockchain type is not supported or empty, must be 'anvil' or 'geth'")
	}
	if err != nil {
		return nil, err
	}
	in.Out = out
	return out, nil
}
