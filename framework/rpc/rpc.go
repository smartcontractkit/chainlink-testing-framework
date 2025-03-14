package rpc

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"math/big"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/core/types"

	f "github.com/smartcontractkit/chainlink-testing-framework/framework"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

// RPCClient is an RPC client for various node simulators
// API Reference https://book.getfoundry.sh/reference/anvil/
// TODO: fix in next PRs
// nolint
type RPCClient struct {
	client *resty.Client
	URL    string
}

// New creates new RPC client that can be used with Geth or Anvil
// this is a high level wrapper for common calls we use
func New(url string, headers http.Header) *RPCClient {
	isDebug := os.Getenv("RESTY_DEBUG") == "true"
	h := make(map[string]string)
	for k, v := range headers {
		h[k] = v[0]
	}
	// TODO: use proper certificates in CRIB
	//nolint
	return &RPCClient{
		URL: url,
		client: resty.New().
			SetDebug(isDebug).
			SetHeaders(h).
			SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}),
	}
}

// AnvilAutoImpersonate sets auto impersonification to true or false
func (m *RPCClient) AnvilAutoImpersonate(b bool) error {
	//nolint:gosec
	rInt := rand.Int()
	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "anvil_autoImpersonateAccount",
		"params":  []interface{}{b},
		"id":      rInt,
	}
	if _, err := m.client.R().SetBody(payload).Post(m.URL); err != nil {
		return errors.Wrap(err, "failed to call anvil_autoImpersonateAccount")
	}
	return nil
}

// AnvilMine calls "evm_mine", mines one or more blocks, see the reference on RPCClient
// API Reference https://book.getfoundry.sh/reference/anvil/
func (m *RPCClient) AnvilMine(params []interface{}) error {
	//nolint:gosec
	rInt := rand.Int()
	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "anvil_mine",
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
	//nolint:gosec
	rInt := rand.Int()
	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "evm_setAutomine",
		"params":  []interface{}{flag},
		"id":      rInt,
	}
	_, err := m.client.R().SetBody(payload).Post(m.URL)
	if err != nil {
		return errors.Wrap(err, "failed to call evm_setAutomine")
	}
	return nil
}

// AnvilTxPoolStatus calls "txpool_status", returns txpool status, see the reference on RPCClient
// API Reference https://book.getfoundry.sh/reference/anvil/
func (m *RPCClient) AnvilTxPoolStatus(params []interface{}) (*TxStatusResponse, error) {
	//nolint:gosec
	rInt := rand.Int()
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
func (m *RPCClient) AnvilSetMinGasPrice(gas uint64) error {
	hexGasPrice := fmt.Sprintf("0x%x", gas)
	//nolint:gosec
	rInt := rand.Int()
	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "anvil_setMinGasPrice",
		"params":  []interface{}{hexGasPrice},
		"id":      rInt,
	}
	if _, err := m.client.R().SetBody(payload).Post(m.URL); err != nil {
		return errors.Wrap(err, "anvil_setMinGasPrice")
	}
	return nil
}

// AnvilSetNextBlockBaseFeePerGas sets next block base fee per gas value
// API Reference https://book.getfoundry.sh/reference/anvil/
func (m *RPCClient) AnvilSetNextBlockBaseFeePerGas(gas *big.Int) error {
	//nolint:gosec
	rInt := rand.Int()
	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "anvil_setNextBlockBaseFeePerGas",
		"params":  []interface{}{gas.String()},
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
	//nolint:gosec
	rInt := rand.Int()
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
	//nolint:gosec
	rInt := rand.Int()
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
	//nolint:gosec
	rInt := rand.Int()
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
	//nolint:gosec
	rInt := rand.Int()
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

type BlockResponse struct {
	Jsonrpc string        `json:"jsonrpc"`
	ID      int           `json:"id"`
	Result  *types.Header `json:"result"`
}

// GetHeaderByNumber retrieves block details by block number
// this is purely debug method to verify the gas chaos is applied
func (m *RPCClient) GetHeaderByNumber(blockNumber int64) (*types.Header, error) {
	//nolint:gosec
	rInt := rand.Int()
	blockNumberHex := fmt.Sprintf("0x%x", blockNumber)
	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_getBlockByNumber",
		"params":  []interface{}{blockNumberHex, false},
		"id":      rInt,
	}
	resp, err := m.client.R().SetBody(payload).Post(m.URL)
	if err != nil {
		return nil, errors.Wrap(err, "eth_getBlockByNumber")
	}
	var blk *BlockResponse
	if err := json.Unmarshal(resp.Body(), &blk); err != nil {
		return nil, err
	}
	if blk.Result == nil {
		return nil, errors.New("block not found")
	}
	return blk.Result, nil
}

func (m *RPCClient) GethSetHead(blocksBack int) error {
	decimalLastBlock, err := m.BlockNumber()
	if err != nil {
		return err
	}
	moveToBlock := decimalLastBlock - int64(blocksBack)
	moveToBlockHex := strconv.FormatInt(moveToBlock, 16)

	//nolint:gosec
	rInt := rand.Int()
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

// PrintBlockBaseFee prints block base fee
// this is purely debug method to verify gas chaos is applied
func (m *RPCClient) PrintBlockBaseFee() error {
	bn, err := m.BlockNumber()
	if err != nil {
		return err
	}
	b, err := m.GetHeaderByNumber(bn)
	if err != nil {
		return err
	}
	f.L.Info().Uint64("BaseFee", b.BaseFee.Uint64()).Msg("Current block")
	return nil
}
