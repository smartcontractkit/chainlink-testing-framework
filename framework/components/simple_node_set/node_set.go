package simple_node_set

import (
	"fmt"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/clnode"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/postgres"
	"golang.org/x/sync/errgroup"
	"strings"
	"sync"
)

type Input struct {
	Nodes        int             `toml:"nodes" validate:"required"`
	OverrideMode string          `toml:"override_mode" validate:"required,oneof=all each"`
	NodeSpecs    []*clnode.Input `toml:"node_specs"`
	Out          *Output         `toml:"out"`
}

type Output struct {
	UseCache bool             `toml:"use_cache"`
	CLNodes  []*clnode.Output `toml:"cl_nodes"`
}

func NewSharedDBNodeSet(in *Input, bcOut *blockchain.Output, fakeUrl string) (*Output, error) {
	if in.Out != nil && in.Out.UseCache {
		return in.Out, nil
	}
	var (
		out *Output
		err error
	)
	defer func() {
		in.Out = out
	}()
	if len(in.NodeSpecs) != in.Nodes && in.OverrideMode == "each" {
		return nil, fmt.Errorf("amount of 'nodes' must be equal to specs provided in override_mode='each'")
	}
	switch in.OverrideMode {
	case "all":
		out, err = sharedDBSetup(in, bcOut, fakeUrl, false)
		if err != nil {
			return nil, err
		}
	case "each":
		out, err = sharedDBSetup(in, bcOut, fakeUrl, true)
		if err != nil {
			return nil, err
		}
	}
	printOut(out)
	return out, nil
}

func printOut(out *Output) {
	for i, n := range out.CLNodes {
		framework.L.Info().Str(fmt.Sprintf("Node-%d", i), n.Node.HostURL).Msg("Chainlink node url")
	}
}

func sharedDBSetup(in *Input, bcOut *blockchain.Output, fakeUrl string, overrideEach bool) (*Output, error) {
	dbOut, err := postgres.NewPostgreSQL(in.NodeSpecs[0].DbInput)
	if err != nil {
		return nil, err
	}
	nodeOuts := make([]*clnode.Output, 0)
	eg := &errgroup.Group{}
	mu := &sync.Mutex{}
	for i := 0; i < in.Nodes; i++ {
		i := i
		var overrideIdx int
		var nodeName string
		if overrideEach {
			overrideIdx = i
		} else {
			overrideIdx = 0
		}
		if in.NodeSpecs[overrideIdx].Node.Name == "" {
			nodeName = fmt.Sprintf("node%d", i)
		}
		eg.Go(func() error {
			net, err := clnode.NewNetworkCfgOneNetworkAllNodes(bcOut)
			if err != nil {
				return err
			}

			nodeSpec := &clnode.Input{
				DataProviderURL: fakeUrl,
				DbInput:         in.NodeSpecs[overrideIdx].DbInput,
				Node: &clnode.NodeInput{
					HTTPPort:                in.NodeSpecs[overrideIdx].Node.HTTPPort + i,
					P2PPort:                 in.NodeSpecs[overrideIdx].Node.P2PPort + i,
					Image:                   in.NodeSpecs[overrideIdx].Node.Image,
					Name:                    nodeName,
					PullImage:               in.NodeSpecs[overrideIdx].Node.PullImage,
					DockerFilePath:          in.NodeSpecs[overrideIdx].Node.DockerFilePath,
					DockerContext:           in.NodeSpecs[overrideIdx].Node.DockerContext,
					DockerImageName:         in.NodeSpecs[overrideIdx].Node.DockerImageName,
					CapabilitiesBinaryPaths: in.NodeSpecs[overrideIdx].Node.CapabilitiesBinaryPaths,
					CapabilityContainerDir:  in.NodeSpecs[overrideIdx].Node.CapabilityContainerDir,
					TestConfigOverrides:     net,
					UserConfigOverrides:     in.NodeSpecs[overrideIdx].Node.UserConfigOverrides,
					TestSecretsOverrides:    in.NodeSpecs[overrideIdx].Node.TestSecretsOverrides,
					UserSecretsOverrides:    in.NodeSpecs[overrideIdx].Node.UserSecretsOverrides,
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
	return &Output{
		UseCache: true,
		CLNodes:  nodeOuts,
	}, nil
}
