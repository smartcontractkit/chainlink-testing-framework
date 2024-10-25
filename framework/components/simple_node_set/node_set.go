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
		i := i
		eg.Go(func() error {
			net, err := clnode.NewNetworkCfgOneNetworkAllNodes(bcOut)
			if err != nil {
				return err
			}

			nodeSpec := &clnode.Input{
				DataProviderURL: fakeUrl,
				DbInput:         in.NodeSpec.DbInput,
				Node: &clnode.NodeInput{
					Image:                   in.NodeSpec.Node.Image,
					Tag:                     in.NodeSpec.Node.Tag,
					Name:                    fmt.Sprintf("node%d", i),
					PullImage:               in.NodeSpec.Node.PullImage,
					Port:                    in.NodeSpec.Node.Port,
					P2PPort:                 in.NodeSpec.Node.P2PPort,
					CapabilitiesBinaryPaths: in.NodeSpec.Node.CapabilitiesBinaryPaths,
					CapabilityContainerDir:  in.NodeSpec.Node.CapabilityContainerDir,
					TestConfigOverrides:     net,
					UserConfigOverrides:     in.NodeSpec.Node.UserConfigOverrides,
					TestSecretsOverrides:    in.NodeSpec.Node.TestSecretsOverrides,
					UserSecretsOverrides:    in.NodeSpec.Node.UserSecretsOverrides,
				},
			}

			dbURL := strings.Replace(dbOut.DockerInternalURL, "/chainlink?sslmode=disable", fmt.Sprintf("/db_%d?sslmode=disable", i), -1)
			dbSpec := &postgres.Output{
				Url:               dbOut.Url,
				DockerInternalURL: dbURL,
			}

			o, err := clnode.NewNode(nodeSpec, dbSpec)
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
