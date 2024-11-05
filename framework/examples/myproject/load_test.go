package examples

import (
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
	BlockchainA        *blockchain.Input `toml:"blockchain_a" validate:"required"`
	MockerDataProvider *fake.Input       `toml:"data_provider" validate:"required"`
	NodeSet            *ns.Input         `toml:"nodeset" validate:"required"`
}

func TestLoad(t *testing.T) {
	in, err := framework.Load[CfgLoad](t)
	require.NoError(t, err)

	bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
	require.NoError(t, err)
	dp, err := fake.NewFakeDataProvider(in.MockerDataProvider)
	require.NoError(t, err)
	_, err = ns.NewSharedDBNodeSet(in.NodeSet, bc, dp.BaseURLDocker)
	require.NoError(t, err)

	t.Run("test something", func(t *testing.T) {
		srv := wasp.NewHTTPMockServer(nil)
		srv.Run()
		labels := map[string]string{
			"go_test_name": "generator_healthcheck",
			"gen_name":     "generator_healthcheck",
			"branch":       "generator_healthcheck",
			"commit":       "generator_healthcheck",
		}
		gen, err := wasp.NewGenerator(&wasp.Config{
			LoadType:   wasp.RPS,
			Schedule:   wasp.Plain(5, 60*time.Second),
			Gun:        NewExampleHTTPGun(srv.URL()),
			Labels:     labels,
			LokiConfig: wasp.NewEnvLokiConfig(),
		})
		require.NoError(t, err)
		gen.Run(true)
	})
}
