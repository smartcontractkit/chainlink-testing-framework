package examples

import (
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/networktest"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
)

type Cfg struct {
	BlockchainA        *blockchain.Input `toml:"blockchain_a" validate:"required"`
	MockerDataProvider *fake.Input       `toml:"data_provider" validate:"required"`
	NodeSets           []*ns.Input       `toml:"nodesets" validate:"required"`
	JD                 *jd.Input         `toml:"jd" validate:"required"`
}

func TestSmoke(t *testing.T) {
	in, err := framework.Load[Cfg](t)
	require.NoError(t, err)

	bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
	require.NoError(t, err)
	_, err = fake.NewFakeDataProvider(in.MockerDataProvider)
	require.NoError(t, err)
	out, err := ns.NewSharedDBNodeSet(in.NodeSets[0], bc)
	require.NoError(t, err)
	_, err = jd.NewJD(in.JD)
	require.NoError(t, err)
	spew.Dump(in.NodeSets[0])
	_, err = networktest.NewNetworkTest(networktest.Input{Privileged: true, BlockInternet: true})
	require.NoError(t, err)
	dc, err := framework.NewDockerClient()
	require.NoError(t, err)
	sOut, err := dc.ExecContainer("networktest", []string{"ping", "-c", "2", "google.com"})
	require.NoError(t, err)
	fmt.Println(sOut)

	t.Run("test something", func(t *testing.T) {
		for _, n := range out.CLNodes {
			require.NotEmpty(t, n.Node.ExternalURL)
		}
	})
}
