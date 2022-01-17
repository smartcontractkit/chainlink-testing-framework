package client_test

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/stretchr/testify/assert"
)

func Test_Eth_Client_Errors(t *testing.T) {
	t.Parallel()
	var (
		err       error
		sendError *client.SendError
		hash      common.Hash
		nonce1    uint64
		nonce2    uint64
		ec        *client.EthereumClient
	)
	ZERO_ADDRESS := "0x0000000000000000000000000000000000000000000000000000000000000000"
	networkInfo := client.NewNetworkConfig()
	t.Run("insufficient funds for transfer", func(t *testing.T) {
		ec, err := client.NewEthereumClient(&networkInfo)
		if err != nil {
			t.Error(err)
		}
		account1 := ec.DefaultWallet
		account2 := ec.Wallets[1]
		OneEth := big.NewInt(1e18)
		amount := big.NewFloat(11)
		bal1, _ := ec.Client.BalanceAt(context.Background(), common.HexToAddress(account1.Address()), nil)
		bal2, _ := ec.Client.BalanceAt(context.Background(), common.HexToAddress(account2.Address()), nil)
		amountFloat, _ := amount.Float64()
		var acc1Bal, acc2Bal big.Int
		acc1Bal.Div(bal1, OneEth)
		acc2Bal.Div(bal2, OneEth)
		if float64(acc2Bal.Int64()) > amountFloat {
			hash, e := ec.SendTransaction(
				account2, common.HexToAddress(account1.Address()),
				big.NewFloat(float64(bal2.Int64())))
			if e != nil {
				t.Error(e)
			}
			fmt.Printf("Reducing account 2 balance - %v balance %v", hash.String(), bal2)
		}
		hash1, err := ec.SendTransaction(
			account2, common.HexToAddress(account1.Address()),
			amount.Mul(amount, client.OneEth))
		sendError := client.NewSendError(err)
		assert.Equal(t, sendError.IsInsufficientEth(), true)
		assert.Equal(t, hash1.String(), ZERO_ADDRESS)
	})

	//TODO find out how nonce can be less than of the nonce in geth state using multiple transaction
	t.Run("nonce too low", func(t *testing.T) {
		ec, err := client.NewEthereumClient(&networkInfo)
		if err != nil {
			t.Error(err)
		}
		account1 := ec.DefaultWallet
		account2 := ec.Wallets[1]
		amount := big.NewFloat(1)
		nonce1, err3 := ec.GetNonce(context.Background(), common.HexToAddress(account1.Address()))
		if err3 != nil {
			t.Error(err3)
		}
		hash, err = ec.SendTransactionWithConfig(
			account1, common.HexToAddress(account2.Address()),
			amount.Mul(amount, client.OneEth),
			&client.TxConfig{
				Nonce: nonce1 - 1,
			},
		)
		sendError := client.NewSendError(err)
		fmt.Printf("### Error String : %v\n", err.Error())
		assert.Equal(t, sendError.IsNonceTooLowError(), true)
		assert.Equal(t, hash.String(), ZERO_ADDRESS)
	})
	// max number that can cause an overflow 9,223,372,036,854,775,807. might not be possibl to reach that
	t.Run("known transaction|already known", func(t *testing.T) {
		ec, err = client.NewEthereumClient(&networkInfo)
		if err != nil {
			t.Error(err)
		}
		account1 := ec.DefaultWallet
		account2 := ec.Wallets[1]
		amount := big.NewFloat(1)
		if nonce1, err = ec.GetNonce(context.Background(), common.HexToAddress(account1.Address())); err != nil {
			t.Error(err)
		}
		if nonce2, err = ec.GetNonce(context.Background(), common.HexToAddress(account1.Address())); err != nil {
			t.Error(err)
		}
		hash, err = ec.SendTransactionWithNonce(
			account1, common.HexToAddress(account2.Address()),
			amount.Mul(amount, client.OneEth),
			nonce2,
		)
		sendError = client.NewSendError(err)
		fmt.Printf("### Error String : %v\n", err.Error())
		assert.NotEqual(t, nonce1, nonce2)
		assert.Equal(t, sendError.IsTransactionAlreadyInMempool(), true)
		assert.Equal(t, hash.String(), ZERO_ADDRESS)
	})

	t.Run("nonce has max value", func(t *testing.T) {

	})
	t.Run("gas limit reached", func(t *testing.T) {

	})
	//To reproduce "invalid sender error" ---
	//Sender function in core/types/transaction_signing.go compares the NetworkID with state NetworkId
	//changing ChainId from the the network gives this error
	t.Run("invalid sender", func(t *testing.T) {
		ec, err := client.NewEthereumClient(&networkInfo)
		if err != nil {
			t.Error(err)
		}
		account1 := ec.DefaultWallet
		account2 := ec.Wallets[1]
		amount := big.NewFloat(11)
		hash1, err := ec.SendTransactionWithConfig(
			account2, common.HexToAddress(account1.Address()),
			amount.Mul(amount, client.OneEth),
			&client.TxConfig{
				ChainID: 1224, // this produces error
			})
		sendError := client.NewSendError(err)
		fmt.Printf("### Error String : %v\n", err.Error())
		assert.Equal(t, sendError.Fatal(), true)
		assert.Equal(t, hash1.String(), ZERO_ADDRESS)
	})
	t.Run("transaction underpriced", func(t *testing.T) {
		ec, err := client.NewEthereumClient(&networkInfo)
		if err != nil {
			t.Error(err)
		}
		account1 := ec.DefaultWallet
		account2 := ec.Wallets[1]
		amount := big.NewFloat(11)
		ecdsaKey, _ := crypto.HexToECDSA(account1.PrivateKey())
		value, _ := amount.Mul(amount, client.OneEth).Int(nil)
		suggestedGasPrice, _ := ec.Client.SuggestGasPrice(context.Background())
		nonce, err := ec.GetNonce(context.Background(), common.HexToAddress(account1.Address()))
		fmt.Printf("Nonce Error %v:\n", err)
		fmt.Printf("Nonce : %v\n", nonce)
		to := common.HexToAddress(account2.Address())
		tx, _ := types.SignNewTx(ecdsaKey, types.NewEIP2930Signer(big.NewInt(networkInfo.ChainID)),
			&types.LegacyTx{
				To:       &to,
				Value:    value,
				Data:     nil,
				Gas:      21000,
				GasPrice: suggestedGasPrice,
				Nonce:    nonce,
			})
		fmt.Printf("Tx value %v\n", tx)
		// hash1, err := ec.SendTransactionWithConfig(
		// 	account2, common.HexToAddress(account1.Address()),
		// 	amount.Mul(amount, client.OneEth),
		// 	&client.TxConfig{
		// 		Tx: tx,
		// 	})
		// fmt.Printf("### Error String : %v\n", err)
		// tx2, _ := types.SignNewTx(ecdsaKey, types.NewEIP2930Signer(big.NewInt(networkInfo.ChainID)),
		// 	&types.LegacyTx{
		// 		To:       &to,
		// 		Value:    value,
		// 		Data:     nil,
		// 		Gas:      21000,
		// 		GasPrice: big.NewInt(1),
		// 		Nonce:    nonce,
		// 	})
		// hash2, err := ec.SendTransactionWithConfig(
		// 	account2, common.HexToAddress(account1.Address()),
		// 	amount.Mul(amount, client.OneEth),
		// 	&client.TxConfig{
		// 		Tx: tx2,
		// 	})
		// fmt.Print(hash2)
		// fmt.Print(nonce, nonce)
		// fmt.Printf("### Error String : %v\n", err)
		// sendError := client.NewSendError(err)
		// assert.Equal(t, sendError.Fatal(), true)
		// assert.Equal(t, hash1.String(), ZERO_ADDRESS)
	})

}

func Int64(float *big.Float) {
	panic("unimplemented")
}
