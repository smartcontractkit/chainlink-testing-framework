package examples

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/examples/generators"
	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
)

type CfgScalability struct {
	BlockchainA        *blockchain.Input `toml:"blockchain_a" validate:"required"`
	MockerDataProvider *fake.Input       `toml:"data_provider" validate:"required"`
	NodeSets           []*ns.Input       `toml:"nodesets" validate:"required"`
}

func TestScalability(t *testing.T) {
	in, err := framework.Load[CfgScalability](t)
	require.NoError(t, err)

	// TODO: these outputs of components come from CRIB staging environment
	bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
	require.NoError(t, err)
	_, err = fake.NewFakeDataProvider(in.MockerDataProvider)
	require.NoError(t, err)
	out, err := ns.NewSharedDBNodeSet(in.NodeSets[0], bc)
	require.NoError(t, err)

	var lokiCfg *wasp.LokiConfig
	// temp fix, we can't reach shared Loki instance in CI
	if os.Getenv("CI") != "true" {
		lokiCfg = wasp.NewEnvLokiConfig()
	}

	c, err := clclient.New(out.CLNodes)
	require.NoError(t, err)

	t.Run("scalability test for your product", func(t *testing.T) {
		_, err := wasp.NewProfile().
			Add(wasp.NewGenerator(&wasp.Config{
				T:        t,
				LoadType: wasp.RPS,
				Schedule: wasp.Combine(
					wasp.Steps(1, 10, 9, 30*time.Second),
					wasp.Plain(10, 30*time.Second),
					wasp.Steps(10, -1, 10, 30*time.Second),
				),
				Gun: generators.NewCLNodeGun(c[0], "bridges"),
				Labels: map[string]string{
					"gen_name": "cl_node_api_call",
					"branch":   "example",
					"commit":   "example",
				},
				LokiConfig: lokiCfg,
			})).
			Run(true)
		require.NoError(t, err)
	})
}
