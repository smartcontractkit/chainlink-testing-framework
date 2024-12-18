package client

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/rand"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

// RPCClient is an RPC client for various node simulators
// API Reference https://book.getfoundry.sh/reference/anvil/
type RPCClient struct {
	client *resty.Client
	URL    string
}

// NewRPCClient creates Anvil client
func NewRPCClient(url string, headers http.Header) *RPCClient {
	isDebug := os.Getenv("RESTY_DEBUG") == "true"
	h := make(map[string]string)
	for k, v := range headers {
		h[k] = v[0]
	}
	// TODO: use proper certificated in CRIB
	//nolint
	return &RPCClient{
		URL: url,
		client: resty.New().
			SetDebug(isDebug).
			SetHeaders(h).
			SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}),
	}
}

// AnvilMine calls "evm_mine", mines one or more blocks, see the reference on RPCClient
// API Reference https://book.getfoundry.sh/reference/anvil/
func (m *RPCClient) AnvilMine(params []interface{}) error {
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

// AnvilSetAutoMine calls "evm_setAutomine", turns automatic mining on, see the reference on RPCClient
// API Reference https://book.getfoundry.sh/reference/anvil/
func (m *RPCClient) AnvilSetAutoMine(flag bool) error {
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

// AnvilTxPoolStatus calls "txpool_status", returns txpool status, see the reference on RPCClient
// API Reference https://book.getfoundry.sh/reference/anvil/
func (m *RPCClient) AnvilTxPoolStatus(params []interface{}) (*TxStatusResponse, error) {
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

// AnvilSetMinGasPrice sets min gas price (pre-EIP-1559 anvil is required)
// API Reference https://book.getfoundry.sh/reference/anvil/
func (m *RPCClient) AnvilSetMinGasPrice(params []interface{}) error {
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

// AnvilSetNextBlockBaseFeePerGas sets next block base fee per gas value
// API Reference https://book.getfoundry.sh/reference/anvil/
func (m *RPCClient) AnvilSetNextBlockBaseFeePerGas(params []interface{}) error {
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

// AnvilSetBlockGasLimit sets next block gas limit
// API Reference https://book.getfoundry.sh/reference/anvil/
func (m *RPCClient) AnvilSetBlockGasLimit(params []interface{}) error {
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

// AnvilDropTransaction removes transaction from tx pool
// API Reference https://book.getfoundry.sh/reference/anvil/
func (m *RPCClient) AnvilDropTransaction(params []interface{}) error {
	rInt, err := rand.Int()
	if err != nil {
		return err
	}
	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "anvil_dropTransaction",
		"params":  params,
		"id":      rInt,
	}
	if _, err := m.client.R().SetBody(payload).Post(m.URL); err != nil {
		return errors.Wrap(err, "anvil_dropTransaction")
	}
	return nil
}

// AnvilSetStorageAt sets storage at address
// API Reference https://book.getfoundry.sh/reference/anvil/
func (m *RPCClient) AnvilSetStorageAt(params []interface{}) error {
	rInt, err := rand.Int()
	if err != nil {
		return err
	}
	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "anvil_setStorageAt",
		"params":  params,
		"id":      rInt,
	}
	if _, err := m.client.R().SetBody(payload).Post(m.URL); err != nil {
		return errors.Wrap(err, "anvil_setStorageAt")
	}
	return nil
}

type CurrentBlockResponse struct {
	Result string `json:"result"`
}

// Call "eth_blockNumber" to get the current block number
func (m *RPCClient) BlockNumber() (int64, error) {
	rInt, err := rand.Int()
	if err != nil {
		return -1, err
	}
	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_blockNumber",
		"params":  []interface{}{},
		"id":      rInt,
	}
	resp, err := m.client.R().SetBody(payload).Post(m.URL)
	if err != nil {
		return -1, errors.Wrap(err, "eth_blockNumber")
	}
	var blockNumberResp *CurrentBlockResponse
	if err := json.Unmarshal(resp.Body(), &blockNumberResp); err != nil {
		return -1, err
	}
	bn, err := strconv.ParseInt(blockNumberResp.Result[2:], 16, 64)
	if err != nil {
		return -1, err
	}
	return bn, nil
}

// GethSetHead sets the Ethereum node's head to a specified block in the past.
// This function is useful for testing and debugging by allowing developers to
// manipulate the blockchain state to a previous block. It returns an error
// if the operation fails.
func (m *RPCClient) GethSetHead(blocksBack int) error {
	decimalLastBlock, err := m.BlockNumber()
	if err != nil {
		return err
	}
	moveToBlock := decimalLastBlock - int64(blocksBack)
	moveToBlockHex := strconv.FormatInt(moveToBlock, 16)

	rInt, err := rand.Int()
	if err != nil {
		return err
	}
	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "debug_setHead",
		"params":  []interface{}{fmt.Sprintf("0x%s", moveToBlockHex)},
		"id":      rInt,
	}
	if _, err := m.client.R().SetBody(payload).Post(m.URL); err != nil {
		return errors.Wrap(err, "debug_setHead")
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

type GethContainer struct {
	testcontainers.Container
	URL string
}

// StartAnvil initializes and starts an Anvil container for Ethereum development.
// It returns an AnvilContainer instance containing the container and its accessible URL.
// This function is useful for developers needing a local Ethereum node for testing and development.
func StartAnvil(params []string) (*AnvilContainer, error) {
	entryPoint := []string{"anvil", "--host", "0.0.0.0"}
	entryPoint = append(entryPoint, params...)
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
