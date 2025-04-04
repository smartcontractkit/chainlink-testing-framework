package blockchain_test

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
)

func TestChains(t *testing.T) {
	testCases := []struct {
		name    string
		input   *blockchain.Input
		chainId int64
	}{
		{
			name: "Anvil",
			input: &blockchain.Input{
				Type:    "anvil",
				Image:   "f4hrenh9it/foundry",
				Port:    "8547",
				ChainID: "31337",
			},
		},
		{
			name: "AnvilZksync",
			input: &blockchain.Input{
				Type:    "anvil-zksync",
				Port:    "8011",
				ChainID: "260",
			},
		},
		{
			name: "Besu",
			input: &blockchain.Input{
				Type:    "besu",
				Port:    "8111",
				WSPort:  "8112",
				ChainID: "1337",
			},
		},
		{
			name: "Geth",
			input: &blockchain.Input{
				Type:    "geth",
				Port:    "8211",
				ChainID: "1337",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testChain(t, tc.input)
		})
	}
}

func testChain(t *testing.T, input *blockchain.Input) {
	chainId, err := strconv.ParseInt(input.ChainID, 10, 64)
	require.NoError(t, err)

	output, err := blockchain.NewBlockchainNetwork(input)
	require.NoError(t, err)

	rpcUrl := output.Nodes[0].ExternalHTTPUrl
	t.Logf("Testing RPC: %s", rpcUrl)
	reqBody := `{"jsonrpc": "2.0", "method": "eth_chainId", "params": [], "id": 1}`
	resp, err := http.Post(rpcUrl, "application/json", strings.NewReader(reqBody)) // nolint:gosec
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	responseData, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	t.Logf("JSON RPC Response: %s", responseData)
	var respJSON struct {
		Result string `json:"result"`
	}
	err = json.Unmarshal(responseData, &respJSON)
	require.NoError(t, err)
	result := respJSON.Result

	actualChainId, err := strconv.ParseInt(strings.TrimPrefix(result, "0x"), 16, 64)
	require.NoError(t, err)

	require.Equal(t, chainId, actualChainId)
}
