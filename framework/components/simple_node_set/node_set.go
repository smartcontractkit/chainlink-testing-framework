package simple_node_set

import (
	"fmt"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/clnode"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/postgres"
	"golang.org/x/sync/errgroup"
	"slices"
	"strings"
	"sync"
)

const (
	DefaultHTTPPortStaticRangeStart = 10000
	DefaultP2PStaticRangeStart      = 12000
)

// Input is a node set configuration input
type Input struct {
	Nodes              int             `toml:"nodes" validate:"required"`
	HTTPPortRangeStart int             `toml:"http_port_range_start"`
	P2PPortRangeStart  int             `toml:"p2p_port_range_start"`
	OverrideMode       string          `toml:"override_mode" validate:"required,oneof=all each"`
	NodeSpecs          []*clnode.Input `toml:"node_specs" validate:"required"`
	Out                *Output         `toml:"out"`
}

// Output is a node set configuration output, used for caching or external components
type Output struct {
	UseCache bool             `toml:"use_cache"`
	CLNodes  []*clnode.Output `toml:"cl_nodes"`
}

// NewSharedDBNodeSet create a new node set with a shared database instance
// all the nodes have their own isolated database
func NewSharedDBNodeSet(in *Input, bcOut *blockchain.Output, fakeUrl string) (*Output, error) {
	if in.Out != nil && in.Out.UseCache {
		return in.Out, nil
	}
	var (
		out *Output
		err error
	)
	defer func() {
		printURLs(out)
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
	return out, nil
}

func printURLs(out *Output) {
	httpURLs, p2pURLs := make([]string, 0), make([]string, 0)
	for _, n := range out.CLNodes {
		httpURLs = append(httpURLs, n.Node.HostURL)
		p2pURLs = append(p2pURLs, n.Node.HostP2PURL)
	}
	framework.L.Info().Any("UI", httpURLs).Send()
	framework.L.Debug().Any("P2P", p2pURLs).Send()
}

func sharedDBSetup(in *Input, bcOut *blockchain.Output, fakeUrl string, overrideEach bool) (*Output, error) {
	in.NodeSpecs[0].DbInput.Databases = in.Nodes
	dbOut, err := postgres.NewPostgreSQL(in.NodeSpecs[0].DbInput)
	if err != nil {
		return nil, err
	}
	nodeOuts := make([]*clnode.Output, 0)

	// to make it easier for chaos testing we use static ports
	// there is no need to check them in advance since testcontainers-go returns a nice error
	var (
		httpPortRangeStart = DefaultHTTPPortStaticRangeStart
		p2pPortRangeStart  = DefaultP2PStaticRangeStart
	)
	if in.HTTPPortRangeStart != 0 {
		httpPortRangeStart = in.HTTPPortRangeStart
	}
	if in.P2PPortRangeStart != 0 {
		p2pPortRangeStart = in.P2PPortRangeStart
	}

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
					HTTPPort:                httpPortRangeStart + i,
					P2PPort:                 p2pPortRangeStart + i,
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
	sortNodeOutsByHostPort(nodeOuts)
	return &Output{
		UseCache: true,
		CLNodes:  nodeOuts,
	}, nil
}

func sortNodeOutsByHostPort(nodes []*clnode.Output) {
	slices.SortFunc[[]*clnode.Output, *clnode.Output](nodes, func(a, b *clnode.Output) int {
		aa := strings.Split(a.Node.HostURL, ":")
		bb := strings.Split(b.Node.HostURL, ":")
		if aa[len(aa)-1] < bb[len(bb)-1] {
			return -1
		} else {
			return 1
		}
	})
}
