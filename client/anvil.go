package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"time"

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

// SetMinGasPrice sets min gas price (pre-EIP-1559 anvil is required)
func (m *AnvilClient) SetMinGasPrice(params []interface{}) error {
	rInt, err := rand.Int()
	if err != nil {
		return err
	}
	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "anvil_setMinGasPrice",
		"params":  params,
		"id":      rInt,
	}
	if _, err := m.client.R().SetBody(payload).Post(m.URL); err != nil {
		return errors.Wrap(err, "anvil_setMinGasPrice")
	}
	return nil
}

// SetNextBlockBaseFeePerGas sets next block base fee per gas value
func (m *AnvilClient) SetNextBlockBaseFeePerGas(params []interface{}) error {
	rInt, err := rand.Int()
	if err != nil {
		return err
	}
	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "anvil_setNextBlockBaseFeePerGas",
		"params":  params,
		"id":      rInt,
	}
	if _, err := m.client.R().SetBody(payload).Post(m.URL); err != nil {
		return errors.Wrap(err, "anvil_setNextBlockBaseFeePerGas")
	}
	return nil
}

// SetBlockGasLimit sets next block gas limit
func (m *AnvilClient) SetBlockGasLimit(params []interface{}) error {
	rInt, err := rand.Int()
	if err != nil {
		return err
	}
	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "evm_setBlockGasLimit",
		"params":  params,
		"id":      rInt,
	}
	if _, err := m.client.R().SetBody(payload).Post(m.URL); err != nil {
		return errors.Wrap(err, "evm_setBlockGasLimit")
	}
	return nil
}

// TxStatusResponse common RPC response body
type TxStatusResponse struct {
	Result struct {
		Pending string `json:"pending"`
	} `json:"result"`
}

type AnvilContainer struct {
	testcontainers.Container
	URL string
}

func StartAnvil(params []string) (*AnvilContainer, error) {
	entryPoint := []string{"anvil", "--host", "0.0.0.0"}
	for _, p := range params {
		entryPoint = append(entryPoint, p)
	}
	req := testcontainers.ContainerRequest{
		Image:        "ghcr.io/foundry-rs/foundry",
		ExposedPorts: []string{"8545/tcp"},
		WaitingFor:   wait.ForListeningPort("8545").WithStartupTimeout(10 * time.Second),
		Entrypoint:   entryPoint,
	}
	container, err := testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}
	time.Sleep(1 * time.Second)
	mappedPort, err := container.MappedPort(context.Background(), "8545")
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("http://localhost:%s", mappedPort.Port())
	return &AnvilContainer{Container: container, URL: url}, nil
}
