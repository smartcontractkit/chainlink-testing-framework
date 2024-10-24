package simple_don

import (
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/clnode"
)

type Input struct {
	Nodes int `toml:"nodes" validate:"required"`
	*clnode.Input
	Out *Output `toml:"out"`
}

type Output struct {
	UseCache bool             `toml:"use_cache"`
	Nodes    []*clnode.Output `toml:"node"`
}

func NewSimpleDON(in *Input, bcOut *blockchain.Output, fakeUrl string) (*Output, error) {
	if in.Out.UseCache {
		return in.Out, nil
	}
	nodeOuts := make([]*clnode.Output, 0)
	for i := 0; i < in.Nodes; i++ {
		net, err := clnode.NewNetworkCfgOneNetworkAllNodes(bcOut)
		if err != nil {
			return nil, err
		}
		in.Input.Node.TestConfigOverrides = net
		in.Input.DataProviderURL = fakeUrl
		in.Input.Out = nil
		o, err := clnode.NewNode(in.Input)
		if err != nil {
			return nil, err
		}
		nodeOuts = append(nodeOuts, o)
	}
	out := &Output{
		UseCache: true,
		Nodes:    nodeOuts,
	}
	in.Out = out
	return out, nil
}
