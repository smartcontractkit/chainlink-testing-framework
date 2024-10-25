package simple_node_set

import (
	"fmt"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/clnode"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/postgres"
	"golang.org/x/sync/errgroup"
	"strings"
	"sync"
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
		o, err := clnode.NewNodeWithDB(newIn)
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

func NewSharedDBNodeSet(in *Input, bcOut *blockchain.Output, fakeUrl string) (*Output, error) {
	if in.Out.UseCache {
		return in.Out, nil
	}
	dbOut, err := postgres.NewPostgreSQL(in.NodeSpec.DbInput)
	if err != nil {
		return nil, err
	}
	nodeOuts := make([]*clnode.Output, 0)
	eg := &errgroup.Group{}
	mu := &sync.Mutex{}
	for i := 0; i < in.Nodes; i++ {
		eg.Go(func() error {
			net, err := clnode.NewNetworkCfgOneNetworkAllNodes(bcOut)
			if err != nil {
				return err
			}
			in.NodeSpec.Node.TestConfigOverrides = net
			in.NodeSpec.DataProviderURL = fakeUrl
			in.NodeSpec.Out = nil

			dbURL := strings.Replace(dbOut.DockerInternalURL, "/chainlink?sslmode=disable", fmt.Sprintf("/db_%d?sslmode=disable", i), -1)
			o, err := clnode.NewNode(in.NodeSpec, &postgres.Output{
				Url:               dbOut.Url,
				DockerInternalURL: dbURL,
			})
			if err != nil {
				return err
			}
			mu.Lock()
			nodeOuts = append(nodeOuts, o)
			mu.Unlock()
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}
	out := &Output{
		UseCache: true,
		CLNodes:  nodeOuts,
	}
	in.Out = out
	return out, nil
}
