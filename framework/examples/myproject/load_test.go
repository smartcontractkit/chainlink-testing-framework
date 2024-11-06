package examples

import (
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/chaos"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
	"github.com/stretchr/testify/require"
	"os"
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
	out, err := ns.NewSharedDBNodeSet(in.NodeSet, bc, dp.BaseURLDocker)
	require.NoError(t, err)

	var lokiCfg *wasp.LokiConfig
	// temp fix, we can't reach shared Loki instance in CI
	if os.Getenv("CI") != "true" {
		lokiCfg = wasp.NewEnvLokiConfig()
	}

	c, err := clclient.NewCLDefaultClients(out.CLNodes, framework.L)
	require.NoError(t, err)

	t.Run("run the cluster and simulate slow network", func(t *testing.T) {
		p, err := wasp.NewProfile().
			Add(wasp.NewGenerator(&wasp.Config{
				T:        t,
				LoadType: wasp.RPS,
				Schedule: wasp.Combine(
					wasp.Steps(1, 1, 9, 30*time.Second),
					wasp.Plain(10, 30*time.Second),
					wasp.Steps(10, -1, 10, 30*time.Second),
				),
				Gun: NewCLNodeGun(c[0], "bridges"),
				Labels: map[string]string{
					"gen_name": "cl_node_api_call",
					"branch":   "example",
					"commit":   "example",
				},
				LokiConfig: lokiCfg,
			})).
			Run(false)
		require.NoError(t, err)
		_, err = chaos.ExecPumba("netem --tc-image=gaiadocker/iproute2 --duration=1m delay --time=300 re2:node.*")
		require.NoError(t, err)
		p.Wait()
	})
}
