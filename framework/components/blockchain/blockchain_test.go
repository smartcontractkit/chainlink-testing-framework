package blockchain

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChains(t *testing.T) {
	testCases := []struct {
		name    string
		input   *Input
		chainId int64
	}{
		{
			name: "Anvil",
			input: &Input{
				Type:  "anvil",
				Image: "f4hrenh9it/foundry",
				Port:  "8555",
			},
			chainId: 31337,
		},
		{
			name: "AnvilZksync",
			input: &Input{
				Type: "anvil-zksync",
				Port: "8011",
			},
			chainId: 260,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testChain(t, tc.chainId, tc.input)
		})
	}
}

func testChain(t *testing.T, chainId int64, input *Input) {
	input.ChainID = strconv.FormatInt(chainId, 10)
	output, err := NewBlockchainNetwork(input)
	require.NoError(t, err)
	rpcUrl := output.Nodes[0].HostHTTPUrl

	reqBody := `{"jsonrpc": "2.0", "method": "eth_chainId", "params": [], "id": 1}`
	resp, err := http.Post(rpcUrl, "application/json", strings.NewReader(reqBody))
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

	actualChainId, err := strconv.ParseInt(strings.TrimPrefix(respJSON.Result, "0x"), 16, 64)
	require.NoError(t, err)

	require.Equal(t, chainId, actualChainId)
}
