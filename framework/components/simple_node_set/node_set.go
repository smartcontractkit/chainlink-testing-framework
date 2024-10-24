package simple_node_set

import (
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/clnode"
)

type Input struct {
	Nodes    int           `toml:"nodes" validate:"required"`
	NodeSpec *clnode.Input `toml:"node_spec" validate:"required"`
	Out      *Output       `toml:"out"`
}

type Output struct {
	UseCache bool             `toml:"use_cache"`
	CLNodes  []*clnode.Output `toml:"cl_nodes"`
}

// NewNodeSet creates a simple set of CL nodes
func NewNodeSet(in *Input, bcOut *blockchain.Output, fakeUrl string) (*Output, error) {
	if in.Out.UseCache {
		return in.Out, nil
	}
	nodeOuts := make([]*clnode.Output, 0)
	for i := 0; i < in.Nodes; i++ {
		net, err := clnode.NewNetworkCfgOneNetworkAllNodes(bcOut)
		if err != nil {
			return nil, err
		}
		newIn := in.NodeSpec
		newIn.Node.TestConfigOverrides = net
		newIn.DataProviderURL = fakeUrl
		newIn.Out = nil
		o, err := clnode.NewNode(newIn)
		if err != nil {
			return nil, err
		}
		nodeOuts = append(nodeOuts, o)
	}
	out := &Output{
		UseCache: true,
		CLNodes:  nodeOuts,
	}
	in.Out = out
	return out, nil
}
