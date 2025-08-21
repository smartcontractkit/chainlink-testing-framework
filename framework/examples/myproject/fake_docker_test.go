package examples

import (
	"fmt"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	components "github.com/smartcontractkit/chainlink-testing-framework/framework/examples/example_components"
)

type CfgDockerFake struct {
	BlockchainA *blockchain.Input `toml:"blockchain_a" validate:"required"`
	Fake        *fake.Input       `toml:"fake" validate:"required"`
	NodeSets    []*ns.Input       `toml:"nodesets" validate:"required"`
}

func TestDockerFakes(t *testing.T) {
	in, err := framework.Load[CfgDockerFake](t)
	require.NoError(t, err)

	bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
	require.NoError(t, err)
	fakeOut, err := fake.NewDockerFakeDataProvider(in.Fake)
	require.NoError(t, err)
	_, err = ns.NewSharedDBNodeSet(in.NodeSets[0], bc)
	require.NoError(t, err)

	t.Run("test fake on host machine", func(t *testing.T) {
		myFakeAPI := "/static-fake"
		resp, err := resty.New().SetBaseURL(fakeOut.BaseURLHost).R().Get(myFakeAPI)
		require.NoError(t, err)
		require.Equal(t, 200, resp.StatusCode())
	})

	t.Run("test fake inside Docker network", func(t *testing.T) {
		myFakeAPI := "/static-fake"
		err := components.NewDockerFakeTester(fmt.Sprintf("%s%s", fakeOut.BaseURLDocker, myFakeAPI))
		require.NoError(t, err)
	})
}
