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
	Nodes        int             `toml:"nodes" validate:"required"`
	OverrideMode string          `toml:"override_mode" validate:"required,oneof=all each"`
	NodeSpecs    []*clnode.Input `toml:"node_specs"`
	Out          *Output         `toml:"out"`
}

type Output struct {
	UseCache bool             `toml:"use_cache"`
	CLNodes  []*clnode.Output `toml:"cl_nodes"`
}

func oneNodeSharedDBConfiguration(in *Input, bcOut *blockchain.Output, fakeUrl string, overrideEach bool) (*Output, error) {
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
					Image:                   in.NodeSpecs[overrideIdx].Node.Image,
					Name:                    nodeName,
					PullImage:               in.NodeSpecs[overrideIdx].Node.PullImage,
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
		out, err = oneNodeSharedDBConfiguration(in, bcOut, fakeUrl, false)
		if err != nil {
			return nil, err
		}
	case "each":
		out, err = oneNodeSharedDBConfiguration(in, bcOut, fakeUrl, true)
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

func NewNodeSet(in *Input, bcOut *blockchain.Output, fakeUrl string) (*Output, error) {
	if in.Out != nil && in.Out.UseCache {
		return in.Out, nil
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
				DbInput:         in.NodeSpecs[i].DbInput,
				Node: &clnode.NodeInput{
					Image:                   in.NodeSpecs[i].Node.Image,
					Name:                    fmt.Sprintf("node%d", i),
					PullImage:               in.NodeSpecs[i].Node.PullImage,
					CapabilitiesBinaryPaths: in.NodeSpecs[i].Node.CapabilitiesBinaryPaths,
					CapabilityContainerDir:  in.NodeSpecs[i].Node.CapabilityContainerDir,
					TestConfigOverrides:     net,
					UserConfigOverrides:     in.NodeSpecs[i].Node.UserConfigOverrides,
					TestSecretsOverrides:    in.NodeSpecs[i].Node.TestSecretsOverrides,
					UserSecretsOverrides:    in.NodeSpecs[i].Node.UserSecretsOverrides,
				},
			}

			dbOut, err := postgres.NewPostgreSQL(in.NodeSpecs[i].DbInput)
			if err != nil {
				return err
			}
			o, err := clnode.NewNode(nodeSpec, dbOut)
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
