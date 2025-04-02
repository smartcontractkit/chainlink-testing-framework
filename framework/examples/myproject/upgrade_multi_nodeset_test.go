package examples

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
)

type CfgUpgradeMulti struct {
	BlockchainA        *blockchain.Input `toml:"blockchain_a" validate:"required"`
	MockerDataProvider *fake.Input       `toml:"data_provider" validate:"required"`
	NodeSetA           *ns.Input         `toml:"nodeset_a" validate:"required"`
	NodeSetB           *ns.Input         `toml:"nodeset_b" validate:"required"`
}

func TestMultiUpgrade(t *testing.T) {
	in, err := framework.Load[CfgUpgradeMulti](t)
	require.NoError(t, err)

	bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
	require.NoError(t, err)
	_, err = fake.NewFakeDataProvider(in.MockerDataProvider)
	require.NoError(t, err)

	// deploy first time
	ns1, err := ns.NewSharedDBNodeSet(in.NodeSetA, bc)
	require.NoError(t, err)
	ns2, err := ns.NewSharedDBNodeSet(in.NodeSetB, bc)
	require.NoError(t, err)

	// reboot both node sets and upgrade with new configs

	in.NodeSetA.NodeSpecs[0].Node.Image = "public.ecr.aws/chainlink/chainlink:v2.17.0"
	in.NodeSetA.NodeSpecs[0].Node.UserConfigOverrides = `
											[Log]
											level = 'info'
	`

	ns1, err = ns.UpgradeNodeSet(t, in.NodeSetA, bc, 3*time.Second)
	require.NoError(t, err)

	in.NodeSetB.NodeSpecs[0].Node.Image = "public.ecr.aws/chainlink/chainlink:v2.17.0"
	in.NodeSetB.NodeSpecs[0].Node.UserConfigOverrides = `
											[Log]
											level = 'info'
	`

	ns2, err = ns.UpgradeNodeSet(t, in.NodeSetB, bc, 3*time.Second)
	require.NoError(t, err)

	t.Run("test something", func(t *testing.T) {
		for _, n := range ns1.CLNodes {
			require.NotEmpty(t, n.Node.ExternalURL)
		}
		for _, n := range ns2.CLNodes {
			require.NotEmpty(t, n.Node.ExternalURL)
		}
	})
}
