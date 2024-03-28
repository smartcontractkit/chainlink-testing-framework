package client

import (
	"encoding/json"

	"github.com/smartcontractkit/chainlink-testing-framework/utils/rand"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

// AnvilClient is an RPC client for Anvil node
// https://book.getfoundry.sh/anvil/
// API Reference https://book.getfoundry.sh/reference/anvil/
type AnvilClient struct {
	client *resty.Client
	URL    string
}

// NewAnvilClient creates Anvil client
func NewAnvilClient(url string) *AnvilClient {
	return &AnvilClient{URL: url, client: resty.New()}
}

// Mine calls "evm_mine", mines one or more blocks, see the reference on AnvilClient
func (m *AnvilClient) Mine(params []interface{}) error {
	rInt, err := rand.Int()
	if err != nil {
		return err
	}
	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "evm_mine",
		"params":  params,
		"id":      rInt,
	}
	if _, err := m.client.R().SetBody(payload).Post(m.URL); err != nil {
		return errors.Wrap(err, "failed to call evm_mine")
	}
	return nil
}

// SetAutoMine calls "evm_setAutomine", turns automatic mining on, see the reference on AnvilClient
func (m *AnvilClient) SetAutoMine(flag bool) error {
	rInt, err := rand.Int()
	if err != nil {
		return err
	}
	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "evm_setAutomine",
		"params":  []interface{}{flag},
		"id":      rInt,
	}
	_, err = m.client.R().SetBody(payload).Post(m.URL)
	if err != nil {
		return errors.Wrap(err, "failed to call evm_setAutomine")
	}
	return nil
}

// TxPoolStatus calls "txpool_status", returns txpool status, see the reference on AnvilClient
func (m *AnvilClient) TxPoolStatus(params []interface{}) (*TxStatusResponse, error) {
	rInt, err := rand.Int()
	if err != nil {
		return nil, err
	}
	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "txpool_status",
		"params":  params,
		"id":      rInt,
	}
	resp, err := m.client.R().SetBody(payload).Post(m.URL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to call txpool_status")
	}
	var txPoolStatusResponse *TxStatusResponse
	if err := json.Unmarshal(resp.Body(), &txPoolStatusResponse); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal txpool_status")
	}
	return txPoolStatusResponse, nil
}

// TxStatusResponse common RPC response body
type TxStatusResponse struct {
	Result struct {
		Pending string `json:"pending"`
	} `json:"result"`
}
