package examples

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
)

type CfgStellar struct {
	BlockchainA *blockchain.Input `toml:"blockchain_a" validate:"required"`
}

type StellarRPCResponse struct {
	JSONRPC string           `json:"jsonrpc"`
	ID      int              `json:"id"`
	Result  json.RawMessage  `json:"result,omitempty"`
	Error   *StellarRPCError `json:"error,omitempty"`
}

type StellarRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type StellarHealthResult struct {
	Status string `json:"status"`
}

type StellarNetworkResult struct {
	FriendbotURL    string `json:"friendbotUrl,omitempty"`
	Passphrase      string `json:"passphrase"`
	ProtocolVersion int    `json:"protocolVersion"`
}

func TestStellarSmoke(t *testing.T) {
	in, err := framework.Load[CfgStellar](t)
	require.NoError(t, err)

	bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
	require.NoError(t, err)

	networkInfo := bc.NetworkSpecificData.StellarNetwork
	require.NotNil(t, networkInfo, "Stellar network info should be present")

	t.Logf("Stellar RPC URL: %s", bc.Nodes[0].ExternalHTTPUrl)
	t.Logf("Stellar Internal URL: %s", bc.Nodes[0].InternalHTTPUrl)
	t.Logf("Network Passphrase: %s", networkInfo.NetworkPassphrase)
	t.Logf("Friendbot URL: %s", networkInfo.FriendbotURL)

	t.Run("verify RPC health", func(t *testing.T) {
		result, err := callStellarRPC[StellarHealthResult](bc.Nodes[0].ExternalHTTPUrl, "getHealth", nil)
		require.NoError(t, err)
		require.Equal(t, "healthy", result.Status, "Stellar RPC should be healthy")
		t.Logf("Health status: %s", result.Status)
	})

	t.Run("verify network info", func(t *testing.T) {
		result, err := callStellarRPC[StellarNetworkResult](bc.Nodes[0].ExternalHTTPUrl, "getNetwork", nil)
		require.NoError(t, err)
		require.Equal(t, blockchain.DefaultStellarNetworkPassphrase, result.Passphrase, "Network passphrase should match")
		t.Logf("Network passphrase: %s", result.Passphrase)
		t.Logf("Protocol version: %d", result.ProtocolVersion)
	})

	t.Run("fund account via Friendbot", func(t *testing.T) {
		testAddress := "GAAZI4TCR3TY5OJHCTJC2A4QSY6CJWJH5IAJTGKIN2ER7LBNVKOCCWN7"

		friendbotURL := fmt.Sprintf("%s?addr=%s", networkInfo.FriendbotURL, testAddress)
		resp, err := http.Get(friendbotURL)
		require.NoError(t, err)
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		t.Logf("Friendbot response status: %d", resp.StatusCode)
		t.Logf("Friendbot response: %s", string(body))

		switch resp.StatusCode {
		case http.StatusOK:
			t.Log("Account funded successfully")
		case http.StatusBadRequest:
			t.Log("Account already funded (expected on retry)")
		case http.StatusBadGateway, http.StatusServiceUnavailable:
			t.Log("Friendbot still initializing - this is expected shortly after startup")
		default:
			t.Errorf("Unexpected Friendbot response: %d", resp.StatusCode)
		}
	})
}

func callStellarRPC[T any](rpcURL, method string, params any) (*T, error) {
	reqBody := map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  method,
	}
	if params != nil {
		reqBody["params"] = params
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := http.Post(rpcURL, "application/json", strings.NewReader(string(jsonBody)))
	if err != nil {
		return nil, fmt.Errorf("failed to call RPC: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var rpcResp StellarRPCResponse
	if err := json.Unmarshal(body, &rpcResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if rpcResp.Error != nil {
		return nil, fmt.Errorf("RPC error %d: %s", rpcResp.Error.Code, rpcResp.Error.Message)
	}

	var result T
	if err := json.Unmarshal(rpcResp.Result, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal result: %w", err)
	}

	return &result, nil
}
