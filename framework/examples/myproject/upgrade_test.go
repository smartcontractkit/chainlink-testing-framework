package examples

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
)

const (
	testJob = `
		type            = "cron"
		schemaVersion   = 1
		schedule        = "CRON_TZ=UTC */10 * * * * *" # every 10 secs
		observationSource   = """
		   // data source 2
		   fetch         [type=http method=GET url="https://min-api.cryptocompare.com/data/pricemultifull?fsyms=ETH&tsyms=USD"];
		   parse       [type=jsonparse path="RAW,ETH,USD,PRICE"];
		   multiply    [type="multiply" input="$(parse)" times=100]
		   encode_tx   [type="ethabiencode"
		                abi="submit(uint256 value)"
		                data="{ \\"value\\": $(multiply) }"]
		   submit_tx   [type="ethtx" to="0x859AAa51961284C94d970B47E82b8771942F1980" data="$(encode_tx)"]
		
		   fetch -> parse -> multiply -> encode_tx -> submit_tx
		"""`
)

type CfgReload struct {
	BlockchainA        *blockchain.Input `toml:"blockchain_a" validate:"required"`
	MockerDataProvider *fake.Input       `toml:"data_provider" validate:"required"`
	NodeSets           []*ns.Input       `toml:"nodesets" validate:"required"`
}

func TestUpgrade(t *testing.T) {
	in, err := framework.Load[CfgReload](t)
	require.NoError(t, err)

	bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
	require.NoError(t, err)
	_, err = fake.NewFakeDataProvider(in.MockerDataProvider)
	require.NoError(t, err)

	// deploy first time
	out, err := ns.NewSharedDBNodeSet(in.NodeSets[0], bc)
	require.NoError(t, err)

	c, err := clclient.New(out.CLNodes)
	require.NoError(t, err)
	_, _, err = c[0].CreateJobRaw(testJob)
	require.NoError(t, err)

	in.NodeSets[0].NodeSpecs[0].Node.Image = "public.ecr.aws/chainlink/chainlink:v2.17.0"
	in.NodeSets[0].NodeSpecs[0].Node.UserConfigOverrides = `
											[Log]
											level = 'info'
	`
	in.NodeSets[0].NodeSpecs[4].Node.Image = "public.ecr.aws/chainlink/chainlink:v2.17.0"
	in.NodeSets[0].NodeSpecs[4].Node.UserConfigOverrides = `
											[Log]
											level = 'info'
	`

	out, err = ns.UpgradeNodeSet(t, in.NodeSets[0], bc, 3*time.Second)
	require.NoError(t, err)

	jobs, _, err := c[0].ReadJobs()
	require.NoError(t, err)
	fmt.Println(jobs)

	t.Run("test something", func(t *testing.T) {
		for _, n := range out.CLNodes {
			require.NotEmpty(t, n.Node.ExternalURL)
		}
	})
}
