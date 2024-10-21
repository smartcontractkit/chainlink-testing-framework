package blockchain

import (
	"fmt"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

// Input is a blockchain network configuration params
type Input struct {
	Type                     string   `toml:"type" validate:"required,oneof=anvil geth"`
	Image                    string   `toml:"image" validate:"required"`
	Tag                      string   `toml:"tag" validate:"required"`
	Port                     string   `toml:"port" validate:"required"`
	ChainID                  string   `toml:"chain_id" validate:"required"`
	DockerCmdParamsOverrides []string `toml:"docker_cmd_params"`
	Out                      *Output  `toml:"out"`
}

// Output is a blockchain network output, ChainID and one or more nodes that forms the network
type Output struct {
	ChainID string  `toml:"chain_id" validate:"required"`
	Nodes   []*Node `toml:"nodes" validate:"required"`
}

// Node represents blockchain node output, URLs required for connection locally and inside docker network
type Node struct {
	WSUrl                 string `toml:"ws_url"`
	HTTPUrl               string `toml:"http_url"`
	DockerInternalWSUrl   string `toml:"docker_internal_ws_url" validate:"required"`
	DockerInternalHTTPUrl string `toml:"docker_internal_http_url" validate:"required"`
}

// NewBlockchainNetwork this is an abstraction that can spin up various blockchain network simulators
// - Anvil
// - Geth
func NewBlockchainNetwork(input *Input) (*Output, error) {
	if input.Out != nil && framework.UseCache() {
		return input.Out, nil
	}
	var out *Output
	var err error
	switch input.Type {
	case "anvil":
		out, err = deployAnvil(input)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("blockchain type is not supported or empty")
	}
	input.Out = out
	return out, nil
}
