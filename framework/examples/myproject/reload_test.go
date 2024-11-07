package examples

import (
	"fmt"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/chaos"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type CfgReload struct {
	BlockchainA        *blockchain.Input `toml:"blockchain_a" validate:"required"`
	MockerDataProvider *fake.Input       `toml:"data_provider" validate:"required"`
	NodeSet            *ns.Input         `toml:"nodeset" validate:"required"`
}

func TestReload(t *testing.T) {
	in, err := framework.Load[CfgReload](t)
	require.NoError(t, err)

	bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
	require.NoError(t, err)
	dp, err := fake.NewFakeDataProvider(in.MockerDataProvider)
	require.NoError(t, err)

	// deploy first time
	out, err := ns.NewSharedDBNodeSet(in.NodeSet, bc, dp.BaseURLDocker)
	require.NoError(t, err)

	c, err := clclient.NewCLDefaultClients(out.CLNodes, framework.L)
	require.NoError(t, err)
	_, _, err = c[0].CreateJobRaw(`
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
		"""`)
	require.NoError(t, err)

	// deploy second time
	_, err = chaos.ExecPumba("rm --volumes=false re2:node.*|postgresql.*", 1*time.Second)
	require.NoError(t, err)
	ns.UpdateNodeConfigs(in.NodeSet, `
[Log]
level = 'info'
`)
	out, err = ns.NewSharedDBNodeSet(in.NodeSet, bc, dp.BaseURLDocker)
	require.NoError(t, err)
	jobs, _, err := c[0].ReadJobs()
	require.NoError(t, err)
	fmt.Println(jobs)

	t.Run("test something", func(t *testing.T) {
		for _, n := range out.CLNodes {
			require.NotEmpty(t, n.Node.HostURL)
			require.NotEmpty(t, n.Node.HostP2PURL)
		}
	})
}
