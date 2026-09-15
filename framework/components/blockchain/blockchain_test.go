package blockchain_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

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
				Image:   "ghcr.io/foundry-rs/foundry:stable",
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

// TestStellar covers the Stellar quickstart chain the same way TestChains covers the
// EVM chains: it starts the container through NewBlockchainNetwork and asserts the
// factory's contract — the Soroban RPC is healthy, the network passphrase matches the
// standalone network, and the Friendbot URL the factory exposes actually funds an account.
// Stellar has no EVM-style chain id / eth_chainId, so it gets its own helper rather than
// being folded into the EVM-centric testChain.
func TestStellar(t *testing.T) {
	input := &blockchain.Input{
		Type: "stellar",
		// Distinct host port so it never clashes with the EVM cases (8547/8011/8111/8211)
		// or the examples smoke_stellar.toml (8100) when both happen to run on the same host.
		Port: "8310",
	}
	testStellar(t, input)
}

func testStellar(t *testing.T, input *blockchain.Input) {
	t.Helper()

	output, err := blockchain.NewBlockchainNetwork(input)
	require.NoError(t, err)

	netInfo := output.NetworkSpecificData.StellarNetwork
	require.NotNil(t, netInfo, "Stellar network info should be present")

	rpcURL := output.Nodes[0].ExternalHTTPUrl
	t.Logf("Testing Stellar RPC: %s", rpcURL)
	t.Logf("Friendbot URL: %s", netInfo.FriendbotURL)

	// getHealth must report healthy — this is the readiness contract newStellar waits on.
	health, err := callStellarRPC[struct {
		Status string `json:"status"`
	}](rpcURL, "getHealth", nil)
	require.NoError(t, err)
	require.Equal(t, "healthy", health.Status, "Stellar RPC should be healthy")

	// getNetwork passphrase must match the standalone network the factory starts.
	network, err := callStellarRPC[struct {
		Passphrase      string `json:"passphrase"`
		ProtocolVersion int    `json:"protocolVersion"`
	}](rpcURL, "getNetwork", nil)
	require.NoError(t, err)
	require.Equal(t, blockchain.DefaultStellarNetworkPassphrase, network.Passphrase,
		"network passphrase should match the standalone network")
	t.Logf("Protocol version: %d", network.ProtocolVersion)

	// The Friendbot URL the factory exposes must actually fund an account. Friendbot can
	// still be warming up right after RPC readiness (it returns 502/503 until ready), so
	// retry until it accepts the funding request. 200 = funded, 400 = already funded on a
	// cached/reused network — both are success.
	addr := "GAAZI4TCR3TY5OJHCTJC2A4QSY6CJWJH5IAJTGKIN2ER7LBNVKOCCWN7"
	require.Eventually(t, func() bool {
		resp, gerr := http.Get(fmt.Sprintf("%s?addr=%s", netInfo.FriendbotURL, addr)) //nolint:gosec
		if gerr != nil {
			return false
		}
		defer resp.Body.Close()
		return resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusBadRequest
	}, 2*time.Minute, 5*time.Second, "friendbot never became ready at %s", netInfo.FriendbotURL)
}

// callStellarRPC is a minimal JSON-RPC 2.0 POST helper for the Soroban RPC endpoint
// (the framework intentionally has no Stellar SDK dependency). It returns the parsed
// `result` object or an error if the RPC returned a JSON-RPC error.
func callStellarRPC[T any](rpcURL, method string, params any) (*T, error) {
	reqBody := map[string]any{"jsonrpc": "2.0", "id": 1, "method": method}
	if params != nil {
		reqBody["params"] = params
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	resp, err := http.Post(rpcURL, "application/json", strings.NewReader(string(body))) //nolint:gosec
	if err != nil {
		return nil, fmt.Errorf("rpc call: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read rpc response: %w", err)
	}

	var envelope struct {
		Result json.RawMessage `json:"result,omitempty"`
		Error  *struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error,omitempty"`
	}
	if err := json.Unmarshal(respBody, &envelope); err != nil {
		return nil, fmt.Errorf("unmarshal rpc response: %w (body: %s)", err, respBody)
	}
	if envelope.Error != nil {
		return nil, fmt.Errorf("rpc error %d: %s", envelope.Error.Code, envelope.Error.Message)
	}

	var result T
	if err := json.Unmarshal(envelope.Result, &result); err != nil {
		return nil, fmt.Errorf("unmarshal rpc result: %w", err)
	}
	return &result, nil
}
