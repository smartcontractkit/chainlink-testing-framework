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

func TestAnvil(t *testing.T) {
	anvilOut, err := NewBlockchainNetwork(&Input{
		Type: "anvil",
		Port: "8555",
	})
	require.NoError(t, err)

	testChain(t, 1337, anvilOut)
}

func TestAnvilZksync(t *testing.T) {
	anvilOut, err := NewBlockchainNetwork(&Input{
		Type: "anvil-zksync",
	})
	require.NoError(t, err)

	testChain(t, 260, anvilOut)
}

func testChain(t *testing.T, chainId int64, output *Output) {
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
