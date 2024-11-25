// nolint
package client

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

const (
	firstAnvilPrivateKey = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	secondAnvilAddress   = "0x70997970C51812dc3A010C7d01b50e0d17dc79C8"
)

func sendTestTransaction(t *testing.T, client *ethclient.Client, gasFeeCap *big.Int, gasTipCap *big.Int, wait bool) (*types.Transaction, error) {
	privateKey, err := crypto.HexToECDSA(firstAnvilPrivateKey)
	require.NoError(t, err)
	publicKey := privateKey.Public()
	publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return nil, err
	}
	value := big.NewInt(1)
	gasLimit := uint64(21000)
	toAddress := common.HexToAddress(secondAnvilAddress)
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return nil, err
	}
	rawTx := &types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasFeeCap: gasFeeCap,
		GasTipCap: gasTipCap,
		Gas:       gasLimit,
		To:        &toAddress,
		Value:     value,
	}
	signedTx, err := types.SignNewTx(privateKey, types.LatestSignerForChainID(chainID), rawTx)
	if err != nil {
		return nil, err
	}
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return nil, err
	}
	t.Log("sent test transaction")
	if wait {
		receipt, err := bind.WaitMined(context.Background(), client, signedTx)
		require.NoError(t, err)
		t.Logf("receipt: %v", receipt)
	}
	return signedTx, err
}

func sendTestTransactions(t *testing.T, client *ethclient.Client, interval time.Duration, gasFeeCap *big.Int, gasTipCap *big.Int, wait bool) (chan struct{}, chan error) {
	stopCh := make(chan struct{})
	errChan := make(chan error, 100)
	go func() {
		for {
			time.Sleep(interval)
			select {
			case <-stopCh:
				close(errChan)
				return
			default:
				_, err := sendTestTransaction(t, client, gasFeeCap, gasTipCap, wait)
				errChan <- err
			}
		}
	}()
	return stopCh, errChan
}

func printGasPrices(t *testing.T, client *ethclient.Client) {
	gasPrice, err := client.SuggestGasPrice(context.Background())
	require.NoError(t, err)
	t.Logf("gas price: %s", gasPrice.String())
	gasTipCap, err := client.SuggestGasTipCap(context.Background())
	require.NoError(t, err)
	t.Logf("gas tip cap: %s", gasTipCap.String())
}

func TestRPCAPI(t *testing.T) {
	t.Run("(geth) test reorg", func(t *testing.T) {
		t.Skip("manual test")
		// TODO: can't use our eth1/eth2 Geth to test here, need to decouple mockserver API
		url := "http://localhost:8545"
		client, err := ethclient.Dial(url)
		require.NoError(t, err)

		bnBefore, err := client.BlockNumber(context.Background())
		require.NoError(t, err)

		ac := NewRPCClient(url, nil)
		err = ac.GethSetHead(10)
		require.NoError(t, err)
		bnAfter, err := client.BlockNumber(context.Background())
		require.NoError(t, err)
		require.Greater(t, bnBefore, bnAfter)
	})

	t.Run("(anvil) test drop transaction", func(t *testing.T) {
		ac, err := StartAnvil([]string{"--balance", "1", "--block-time", "5"})
		require.NoError(t, err)
		client, err := ethclient.Dial(ac.URL)
		require.NoError(t, err)

		tx, err := sendTestTransaction(t, client, big.NewInt(1e9), big.NewInt(1e9), false)
		require.NoError(t, err)

		anvilClient := NewRPCClient(ac.URL, nil)
		err = anvilClient.AnvilDropTransaction([]interface{}{tx.Hash().String()})
		require.NoError(t, err)
		status, err := anvilClient.AnvilTxPoolStatus(nil)
		require.NoError(t, err)
		require.Equal(t, status.Result.Pending, "0x0")
		t.Logf("status: %v", status)
	})

	t.Run("(anvil) test set storage at address", func(t *testing.T) {
		ac, err := StartAnvil([]string{"--balance", "1", "--block-time", "5"})
		require.NoError(t, err)
		client, err := ethclient.Dial(ac.URL)
		require.NoError(t, err)

		randomAddress := common.HexToAddress("0x0d2026b3EE6eC71FC6746ADb6311F6d3Ba1C000B")
		storeValue := "0x0000000000000000000000000000000000000000000000000000000000000001"

		anvilClient := NewRPCClient(ac.URL, nil)
		err = anvilClient.AnvilSetStorageAt([]interface{}{randomAddress.Hex(), "0x0", storeValue})
		require.NoError(t, err, "unable to set storage at address")

		value, err := client.StorageAt(context.Background(), randomAddress, common.HexToHash("0x0"), nil)
		require.NoError(t, err)
		decodedStoreValue, err := hex.DecodeString(storeValue[2:])
		require.NoError(t, err, "unable to decode store value")
		require.Equal(t, decodedStoreValue, value)

		t.Logf("value: %v", value)
	})

	t.Run("(anvil) test we can shrink the block and control transaction inclusion", func(t *testing.T) {
		ac, err := StartAnvil([]string{"--balance", "1", "--block-time", "1"})
		require.NoError(t, err)
		client, err := ethclient.Dial(ac.URL)
		require.NoError(t, err)

		anvilClient := NewRPCClient(ac.URL, nil)
		err = anvilClient.AnvilSetBlockGasLimit([]interface{}{"1"})
		require.NoError(t, err)

		_, err = sendTestTransaction(t, client, big.NewInt(1e9), big.NewInt(1e9), true)
		require.Error(t, err)

		err = anvilClient.AnvilSetBlockGasLimit([]interface{}{"30000000"})
		require.NoError(t, err)

		_, err = sendTestTransaction(t, client, big.NewInt(1e9), big.NewInt(1e9), true)
		require.NoError(t, err)
	})

	t.Run("(anvil) test we can change next block base fee per gas and make tx pass or fail", func(t *testing.T) {
		ac, err := StartAnvil([]string{"--balance", "1", "--block-time", "1"})
		require.NoError(t, err)
		client, err := ethclient.Dial(ac.URL)
		require.NoError(t, err)
		printGasPrices(t, client)

		anvilClient := NewRPCClient(ac.URL, nil)
		err = anvilClient.AnvilSetNextBlockBaseFeePerGas([]interface{}{"10000000000"})
		require.NoError(t, err)
		printGasPrices(t, client)
		gasPrice, err := client.SuggestGasPrice(context.Background())
		require.NoError(t, err)
		require.Equal(t, int64(11000000000), gasPrice.Int64())

		_, err = sendTestTransaction(t, client, big.NewInt(1e9), big.NewInt(1e9), true)
		require.Error(t, err)

		err = anvilClient.AnvilSetNextBlockBaseFeePerGas([]interface{}{"1"})
		require.NoError(t, err)

		_, err = sendTestTransaction(t, client, big.NewInt(1e9), big.NewInt(1e9), true)
		require.NoError(t, err)
	})

	t.Run("(anvil) test we can mine sub-second blocks", func(t *testing.T) {
		period := 500 * time.Millisecond
		iterations := 10
		ac, err := StartAnvil([]string{"--no-mine"})
		require.NoError(t, err)
		client, err := ethclient.Dial(ac.URL)
		require.NoError(t, err)
		pm := NewRemoteAnvilMiner(ac.URL, nil)
		pm.MinePeriodically(500 * time.Millisecond)
		time.Sleep(period * time.Duration(iterations))
		pm.Stop()
		bn, err := client.BlockNumber(context.Background())
		require.NoError(t, err)
		require.GreaterOrEqual(t, uint64(iterations), bn-1)
		require.LessOrEqual(t, uint64(iterations), bn+1)
	})

	t.Run("(anvil) test we can mine blocks with strictly N+ transactions", func(t *testing.T) {
		sendTransactionEvery := 500 * time.Millisecond
		txnInBlock := int64(10)
		iterations := 10
		ac, err := StartAnvil([]string{"--no-mine"})
		require.NoError(t, err)
		client, err := ethclient.Dial(ac.URL)
		require.NoError(t, err)
		stopTxns, errCh := sendTestTransactions(t, client, sendTransactionEvery, big.NewInt(1e9), big.NewInt(1e9), false)
		pm := NewRemoteAnvilMiner(ac.URL, nil)
		pm.MineBatch(txnInBlock, 1*time.Second, 1*time.Minute)
		time.Sleep(sendTransactionEvery * time.Duration(iterations) * 2)
		pm.Stop()
		stopTxns <- struct{}{}
		for e := range errCh {
			require.NoError(t, e)
		}
		bn, err := client.BlockNumber(context.Background())
		require.NoError(t, err)
		for i := 1; i <= int(bn); i++ {
			block, err := client.BlockByNumber(context.Background(), big.NewInt(int64(i)))
			require.NoError(t, err)
			require.GreaterOrEqual(t, int64(block.Transactions().Len()), txnInBlock)
		}
	})
}
