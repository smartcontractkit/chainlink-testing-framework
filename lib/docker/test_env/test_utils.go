package test_env

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/blockchain"
)

func sendAndCompareBalances(ctx context.Context, c blockchain.EVMClient, address common.Address) error {
	balanceBefore, err := c.BalanceAt(ctx, address)
	if err != nil {
		return err
	}

	toSendEth := big.NewFloat(1)

	exp := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	toSendEthInt := new(big.Int)
	_, _ = toSendEth.Int(toSendEthInt)
	sentInWei := new(big.Int).Mul(toSendEthInt, exp)

	gasEstimations, err := c.EstimateGas(ethereum.CallMsg{
		From:  common.HexToAddress(c.GetDefaultWallet().Address()),
		To:    &address,
		Value: sentInWei,
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

	expected := big.NewInt(0).Add(balanceBefore, sentInWei)

	if expected.Cmp(balanceAfter) != 0 {
		return fmt.Errorf("balance is incorrect. Expected %s, got %s", expected.String(), balanceAfter.String())
	}

	return nil
}
