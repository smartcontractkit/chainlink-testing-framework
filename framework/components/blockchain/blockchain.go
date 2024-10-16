package blockchain

import (
	"fmt"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

type Input struct {
	Type                     string   `toml:"type" validate:"required,oneof=anvil geth"`
	Port                     string   `toml:"port" validate:"required"`
	ChainID                  string   `toml:"chain_id" validate:"required"`
	DockerCmdParamsOverrides []string `toml:"docker_cmd_params"`
	Out                      *Output  `toml:"out"`
}

type Output struct {
	ChainID string  `toml:"chain_id" validate:"required"`
	Nodes   []*Node `toml:"nodes" validate:"required"`
}

type Node struct {
	WSUrl   string `toml:"ws_url" validate:"required"`
	HTTPUrl string `toml:"http_url" validate:"required"`
}

func NewBlockchainNetwork(input *Input) (*Output, error) {
	if input.Out != nil && framework.NoCache() {
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
