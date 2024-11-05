package examples

import (
	"fmt"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type CfgLoad struct {
	BlockchainA *blockchain.Input `toml:"blockchain_a" validate:"required"`
	//Contracts          *onchain.Input    `toml:"contracts" validate:"required"`
	MockerDataProvider *fake.Input `toml:"data_provider" validate:"required"`
	NodeSet            *ns.Input   `toml:"nodeset" validate:"required"`
}

func TestLoad(t *testing.T) {
	in, err := framework.Load[CfgLoad](t)
	require.NoError(t, err)
	//pkey := os.Getenv("PRIVATE_KEY")

	bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
	require.NoError(t, err)
	dp, err := fake.NewFakeDataProvider(in.MockerDataProvider)
	require.NoError(t, err)
	_, err = ns.NewSharedDBNodeSet(in.NodeSet, bc, dp.BaseURLDocker)
	require.NoError(t, err)

	//// deploy product contracts
	//in.Contracts.URL = bc.Nodes[0].HostWSUrl
	//contracts, err := onchain.NewProductOnChainDeployment(in.Contracts)
	//require.NoError(t, err)
	//
	//sc, err := seth.NewClientBuilder().
	//	WithRpcUrl(bc.Nodes[0].HostWSUrl).
	//	WithGasPriceEstimations(true, 0, seth.Priority_Fast).
	//	WithTracing(seth.TracingLevel_All, []string{seth.TraceOutput_Console}).
	//	WithPrivateKeys([]string{pkey}).
	//	Build()
	//require.NoError(t, err)

	t.Run("test something", func(t *testing.T) {
		labels := map[string]string{
			"go_test_name": "generator_healthcheck",
			"gen_name":     "generator_healthcheck",
			"branch":       "generator_healthcheck",
			"commit":       "generator_healthcheck",
		}
		gen, err := wasp.NewGenerator(&wasp.Config{
			LoadType: wasp.RPS,
			Schedule: wasp.Combine(
				wasp.Steps(1, 1, 9, 30*time.Second),
				wasp.Plain(10, 30*time.Second),
				wasp.Steps(10, -1, 10, 30*time.Second),
			),
			Gun:        NewExampleHTTPGun(fmt.Sprintf("%s/mock1", dp.BaseURLHost)),
			Labels:     labels,
			LokiConfig: wasp.NewEnvLokiConfig(),
		})
		require.NoError(t, err)
		gen.Run(true)
	})
}
