package test_env

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/testcontext"
)

func TestEth2WithPrysmAndBesu(t *testing.T) {
	l := logging.GetTestLogger(t)

	builder := NewEthereumNetworkBuilder()
	cfg, err := builder.
		WithConsensusType(ConsensusType_PoS).
		WithConsensusLayer(ConsensusLayer_Prysm).
		WithExecutionLayer(ExecutionLayer_Besu).
		Build()
	require.NoError(t, err, "Builder validation failed")

	_, eth2, err := cfg.Start()
	require.NoError(t, err, "Couldn't start PoS network")

	nonEip1559Network := blockchain.SimulatedEVMNetwork
	nonEip1559Network.Name = "Simulated Besu + Prysm (non-EIP 1559)"
	nonEip1559Network.GasEstimationBuffer = 10_000_000_000
	nonEip1559Network.URLs = eth2.PublicWsUrls()
	clientOne, err := blockchain.ConnectEVMClient(nonEip1559Network, l)
	require.NoError(t, err, "Couldn't connect to the evm client")

	t.Cleanup(func() {
		err = clientOne.Close()
		require.NoError(t, err, "Couldn't close the client")
	})

	address := common.HexToAddress("0x90F8bf6A479f320ead074411a4B0e7944Ea8c9C1")
	err = sendAndCompareBalances(testcontext.Get(t), clientOne, address)
	require.NoError(t, err, fmt.Sprintf("balance wasn't correctly updated when %s network", nonEip1559Network.Name))

	eip1559Network := blockchain.SimulatedEVMNetwork
	eip1559Network.Name = "Simulated Besu + Prysm (EIP 1559)"
	eip1559Network.SupportsEIP1559 = true
	eip1559Network.URLs = eth2.PublicWsUrls()
	_, err = blockchain.ConnectEVMClient(eip1559Network, l)
	require.Error(t, err, "Couldn't connect to the evm client")
	require.Contains(t, err.Error(), "Method not found", "Besu should not work EIP-1559 yet")
}

func sendAndCompareBalances(ctx context.Context, c blockchain.EVMClient, address common.Address) error {
	balanceBefore, err := c.BalanceAt(ctx, address)
	if err != nil {
		return err
	}

	toSendEth := big.NewFloat(1)
	gasEstimations, err := c.EstimateGas(ethereum.CallMsg{
		To: &address,
	})
	if err != nil {
		return err
	}
	err = c.Fund(address.Hex(), toSendEth, gasEstimations)
	if err != nil {
		return err
	}

	balanceAfter, err := c.BalanceAt(ctx, address)
	if err != nil {
		return err
	}

	exp := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	toSendEthInt := new(big.Int)
	_, _ = toSendEth.Int(toSendEthInt)
	sentInWei := new(big.Int).Mul(toSendEthInt, exp)

	expected := big.NewInt(0).Add(balanceBefore, sentInWei)

	if expected.Cmp(balanceAfter) != 0 {
		return errors.Errorf("Balance is incorrect. Expected %s, got %s", expected.String(), balanceAfter.String())
	}

	return nil
}
