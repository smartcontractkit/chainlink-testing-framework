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

type CfgFake struct {
	BlockchainA *blockchain.Input `toml:"blockchain_a" validate:"required"`
	Fake        *fake.Input       `toml:"fake" validate:"required"`
	NodeSet     *ns.Input         `toml:"nodeset" validate:"required"`
}

func TestFakes(t *testing.T) {
	in, err := framework.Load[CfgFake](t)
	require.NoError(t, err)

	bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
	require.NoError(t, err)
	fakeOut, err := fake.NewFakeDataProvider(in.Fake)
	require.NoError(t, err)
	_, err = ns.NewSharedDBNodeSet(in.NodeSet, bc)
	require.NoError(t, err)

	t.Run("test fake on host machine", func(t *testing.T) {
		myFakeAPI := "/fake/api/one"
		// use fake.Func if you need full control over response
		err = fake.JSON(
			"GET",
			myFakeAPI,
			map[string]any{
				"data": "some_data",
			}, 200,
		)
		require.NoError(t, err)
		resp, err := resty.New().SetBaseURL(fakeOut.BaseURLHost).R().Get(myFakeAPI)
		require.NoError(t, err)
		require.Equal(t, 200, resp.StatusCode())

		// you can also access all recorded requests and responses
		data, err := fake.R.Get("GET", myFakeAPI)
		for _, rec := range data {
			fmt.Println(rec.Status)
			fmt.Println(rec.Method)
			fmt.Println(rec.Path)
			fmt.Println(rec.Headers)
			fmt.Println(rec.ReqBody)
			fmt.Println(rec.ResBody)
		}
	})
	t.Run("access fake from docker network", func(t *testing.T) {
		myFakeAPI := "/fake/api/two"
		err = fake.JSON(
			"GET",
			myFakeAPI,
			map[string]any{
				"data": "some_data",
			}, 200,
		)
		require.NoError(t, err)

		// use docker URL and path of your fake
		_ = fmt.Sprintf("%s%s", fakeOut.BaseURLDocker, myFakeAPI)
	})

	t.Run("verify that containers can access internally both locally and in CI", func(t *testing.T) {
		myFakeAPI := "/fake/api/internal"
		// use fake.Func if you need full control over response
		err = fake.JSON(
			"GET",
			myFakeAPI,
			map[string]any{
				"data": "some_data",
			}, 200,
		)
		require.NoError(t, err)
		err := components.NewDockerFakeTester(fmt.Sprintf("%s%s", fakeOut.BaseURLDocker, myFakeAPI))
		require.NoError(t, err)
	})
}
