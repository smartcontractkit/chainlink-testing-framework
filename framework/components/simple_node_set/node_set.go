package simple_node_set

import (
	"context"
	"fmt"
	"os"
	"slices"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/clnode"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/postgres"
	"github.com/smartcontractkit/pods"
)

const (
	DefaultHTTPPortStaticRangeStart = 10000
	DefaultP2PStaticRangeStart      = 12000
)

// Input is a node set configuration input
type Input struct {
	Name               string          `toml:"name" validate:"required" comment:"Node set name, ex.:'don-1', Docker containers will be prefixed with this name so tests can distinguish one DON from another"`
	Nodes              int             `toml:"nodes" validate:"required" comment:"Number of nodes in node set"`
	HTTPPortRangeStart int             `toml:"http_port_range_start" comment:"HTTP ports range starting with port X and increasing by 1"`
	P2PPortRangeStart  int             `toml:"p2p_port_range_start" comment:"P2P ports range starting with port X and increasing by 1"`
	DlvPortRangeStart  int             `toml:"dlv_port_range_start" comment:"Delve debugger ports range starting with port X and increasing by 1"`
	OverrideMode       string          `toml:"override_mode" validate:"required,oneof=all each" comment:"Override mode, applicable only to 'localcre'. Changes how config overrides to TOML nodes apply"`
	DbInput            *postgres.Input `toml:"db" validate:"required" comment:"Shared node set data base input for PostgreSQL"`
	NodeSpecs          []*clnode.Input `toml:"node_specs" validate:"required" comment:"Chainlink node TOML configurations"`
	NoDNS              bool            `toml:"no_dns" comment:"Turn DNS on, helpful to isolate container from the internet"`
	Out                *Output         `toml:"out" comment:"Nodeset config output"`
}

// Output is a node set configuration output, used for caching or external components
type Output struct {
	// UseCache Whether to respect caching or not, if cache = true component won't be deployed again
	UseCache bool `toml:"use_cache" comment:"Whether to respect caching or not, if cache = true component won't be deployed again"`
	// DBOut Nodeset shared database output (PostgreSQL)
	DBOut *postgres.Output `toml:"db_out" comment:"Nodeset shared database output (PostgreSQL)"`
	// CLNodes Chainlink node config outputs
	CLNodes []*clnode.Output `toml:"cl_nodes" comment:"Chainlink node config outputs"`
}

// NewSharedDBNodeSet create a new node set with a shared database instance
// all the nodes have their own isolated database
func NewSharedDBNodeSet(in *Input, bcOut *blockchain.Output) (*Output, error) {
	return NewSharedDBNodeSetWithContext(context.Background(), in, bcOut)
}

// NewSharedDBNodeSetWithContext create a new node set with a shared database instance
// all the nodes have their own isolated database
func NewSharedDBNodeSetWithContext(ctx context.Context, in *Input, bcOut *blockchain.Output) (*Output, error) {
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
	out, err = sharedDBSetup(ctx, in, bcOut)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func NodeNamePrefix(nodeSetName string) string {
	return nodeSetName + "-" + "node"
}

func printURLs(out *Output) {
	if out == nil {
		return
	}
	httpURLs, _, pgURLs := make([]string, 0), make([]string, 0), make([]string, 0)
	for _, n := range out.CLNodes {
		httpURLs = append(httpURLs, n.Node.ExternalURL)
		pgURLs = append(pgURLs, n.PostgreSQL.Url)
	}
	framework.L.Info().Any("UI", httpURLs).Send()
	framework.L.Debug().Any("DB", pgURLs).Send()
}

func sharedDBSetup(ctx context.Context, in *Input, bcOut *blockchain.Output) (*Output, error) {
	in.DbInput.Name = fmt.Sprintf("%s-%s", in.Name, "ns-postgresql")
	in.DbInput.VolumeName = in.Name

	// create database for each node
	in.DbInput.Databases = in.Nodes
	dbOut, err := postgres.NewWithContext(ctx, in.DbInput)
	if err != nil {
		return nil, err
	}
	nodeOuts := make([]*clnode.Output, 0)

	envImage := os.Getenv("CTF_CHAINLINK_IMAGE")

	// to make it easier for chaos testing we use static ports
	// there is no need to check them in advance since testcontainers-go returns a nice error
	var (
		httpPortRangeStart = DefaultHTTPPortStaticRangeStart
		p2pPortRangeStart  = DefaultP2PStaticRangeStart
		dlvPortStart       = clnode.DefaultDebuggerPort
	)
	if in.HTTPPortRangeStart != 0 {
		httpPortRangeStart = in.HTTPPortRangeStart
	}
	if in.P2PPortRangeStart != 0 {
		p2pPortRangeStart = in.P2PPortRangeStart
	}
	if in.DlvPortRangeStart != 0 {
		dlvPortStart = in.DlvPortRangeStart
	}

	eg := &errgroup.Group{}
	mu := &sync.Mutex{}
	for i := 0; i < in.Nodes; i++ {
		overrideIdx := i
		if in.OverrideMode == "all" {
			if len(in.NodeSpecs[overrideIdx].Node.CustomPorts) > 0 {
				return nil, fmt.Errorf("custom_ports can be used only with override_mode = 'each'")
			}
		}

		eg.Go(func() error {
			var net string
			var err error
			if bcOut != nil {
				net, err = clnode.NewNetworkCfgOneNetworkAllNodes(bcOut)
				if err != nil {
					return err
				}
			}
			if in.NodeSpecs[overrideIdx].Node.TestConfigOverrides != "" {
				net = in.NodeSpecs[overrideIdx].Node.TestConfigOverrides
			}
			nodeWithNodeSetPrefixName := NodeNamePrefix(in.Name) + fmt.Sprint(i)

			nodeSpec := &clnode.Input{
				NoDNS:   in.NoDNS,
				DbInput: in.DbInput,
				Node: &clnode.NodeInput{
					HTTPPort:                httpPortRangeStart + i,
					P2PPort:                 p2pPortRangeStart + i,
					DebuggerPort:            dlvPortStart + i,
					CustomPorts:             in.NodeSpecs[overrideIdx].Node.CustomPorts,
					Image:                   in.NodeSpecs[overrideIdx].Node.Image,
					Name:                    nodeWithNodeSetPrefixName,
					PullImage:               in.NodeSpecs[overrideIdx].Node.PullImage,
					DockerFilePath:          in.NodeSpecs[overrideIdx].Node.DockerFilePath,
					DockerContext:           in.NodeSpecs[overrideIdx].Node.DockerContext,
					DockerBuildArgs:         in.NodeSpecs[overrideIdx].Node.DockerBuildArgs,
					CapabilitiesBinaryPaths: in.NodeSpecs[overrideIdx].Node.CapabilitiesBinaryPaths,
					CapabilityContainerDir:  in.NodeSpecs[overrideIdx].Node.CapabilityContainerDir,
					TestConfigOverrides:     net,
					UserConfigOverrides:     in.NodeSpecs[overrideIdx].Node.UserConfigOverrides,
					TestSecretsOverrides:    in.NodeSpecs[overrideIdx].Node.TestSecretsOverrides,
					UserSecretsOverrides:    in.NodeSpecs[overrideIdx].Node.UserSecretsOverrides,
					ContainerResources:      in.NodeSpecs[overrideIdx].Node.ContainerResources,
					EnvVars:                 in.NodeSpecs[overrideIdx].Node.EnvVars,
				},
			}

			if envImage != "" {
				nodeSpec.Node.Image = envImage
				// unset docker build context and file path to avoid conflicts, image provided via env var takes precedence
				nodeSpec.Node.DockerContext = ""
				nodeSpec.Node.DockerFilePath = ""
			}

			dbURLHost := strings.Replace(dbOut.Url, "/chainlink?sslmode=disable", fmt.Sprintf("/db_%d?sslmode=disable", i), -1)
			dbURL := strings.Replace(dbOut.InternalURL, "/chainlink?sslmode=disable", fmt.Sprintf("/db_%d?sslmode=disable", i), -1)
			dbSpec := &postgres.Output{
				Url:         dbURLHost,
				InternalURL: dbURL,
			}

			o, err := clnode.NewNodeWithContext(ctx, nodeSpec, dbSpec)
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
	// wait for all K8s services at once
	if os.Getenv(components.K8sNamespaceEnvVar) != "" {
		pods.WaitReady(3 * time.Minute)
	}
	return &Output{
		UseCache: true,
		DBOut:    dbOut,
		CLNodes:  nodeOuts,
	}, nil
}

func sortNodeOutsByHostPort(nodes []*clnode.Output) {
	slices.SortFunc(nodes, func(a, b *clnode.Output) int {
		aa := strings.Split(a.Node.ExternalURL, ":")
		bb := strings.Split(b.Node.ExternalURL, ":")
		if aa[len(aa)-1] < bb[len(bb)-1] {
			return -1
		}
		return 1
	})
}
