package client

import (
	"context"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

const (
	firstAnvilPrivateKey = "5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a"
	secondAnvilAddress   = "0x70997970C51812dc3A010C7d01b50e0d17dc79C8"
)

var ()

func sendTestTransaction(t *testing.T, client *ethclient.Client, gasFeeCap *big.Int, gasTipCap *big.Int, wait bool) error {
	privateKey, err := crypto.HexToECDSA(firstAnvilPrivateKey)
	require.NoError(t, err)
	publicKey := privateKey.Public()
	publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return err
	}
	value := big.NewInt(1)
	gasLimit := uint64(21000)
	toAddress := common.HexToAddress(secondAnvilAddress)
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return err
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
		return err
	}
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return err
	}
	t.Log("sent test transaction")
	if wait {
		receipt, err := bind.WaitMined(context.Background(), client, signedTx)
		require.NoError(t, err)
		t.Logf("receipt: %v", receipt)
	}
	return err
}

func sendTestTransactions(t *testing.T, client *ethclient.Client, interval time.Duration, gasFeeCap *big.Int, gasTipCap *big.Int, wait bool) (chan struct{}, chan error) {
	stopCh := make(chan struct{})
	errChan := make(chan error)
	go func() {
		for {
			time.Sleep(interval)
			select {
			case <-stopCh:
				return
			default:
				sendTestTransaction(t, client, gasFeeCap, gasTipCap, wait)
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

func TestAnvilAPIs(t *testing.T) {
	t.Run("test we can shrink the block and control transaction inclusion", func(t *testing.T) {
		ac, err := StartAnvil([]string{"--balance", "1", "--block-time", "1"})
		require.NoError(t, err)
		client, err := ethclient.Dial(ac.URL)
		require.NoError(t, err)

		anvilClient := NewAnvilClient(ac.URL)
		err = anvilClient.SetBlockGasLimit([]interface{}{"1"})
		require.NoError(t, err)

		err = sendTestTransaction(t, client, big.NewInt(1e9), big.NewInt(1e9), true)
		require.Error(t, err)

		err = anvilClient.SetBlockGasLimit([]interface{}{"30000000"})
		require.NoError(t, err)

		err = sendTestTransaction(t, client, big.NewInt(1e9), big.NewInt(1e9), true)
		require.NoError(t, err)
	})

	t.Run("test we can change next block base fee per gas and make tx pass or fail", func(t *testing.T) {
		ac, err := StartAnvil([]string{"--balance", "1", "--block-time", "1"})
		require.NoError(t, err)
		client, err := ethclient.Dial(ac.URL)
		require.NoError(t, err)
		printGasPrices(t, client)

		anvilClient := NewAnvilClient(ac.URL)
		err = anvilClient.SetNextBlockBaseFeePerGas([]interface{}{"10000000000"})
		require.NoError(t, err)
		printGasPrices(t, client)
		gasPrice, err := client.SuggestGasPrice(context.Background())
		require.NoError(t, err)
		require.Equal(t, int64(11000000000), gasPrice.Int64())

		err = sendTestTransaction(t, client, big.NewInt(1e9), big.NewInt(1e9), true)
		require.Error(t, err)

		err = anvilClient.SetNextBlockBaseFeePerGas([]interface{}{"1"})
		require.NoError(t, err)

		err = sendTestTransaction(t, client, big.NewInt(1e9), big.NewInt(1e9), true)
		require.NoError(t, err)
	})

	t.Run("test we can mine sub-second blocks", func(t *testing.T) {
		period := 500 * time.Millisecond
		iterations := 10
		ac, err := StartAnvil([]string{"--no-mine"})
		require.NoError(t, err)
		client, err := ethclient.Dial(ac.URL)
		require.NoError(t, err)
		pm := NewRemoteAnvilMiner(ac.URL)
		pm.MinePeriodically(500 * time.Millisecond)
		time.Sleep(period * time.Duration(iterations))
		pm.Stop()
		bn, err := client.BlockNumber(context.Background())
		require.NoError(t, err)
		require.Equal(t, uint64(iterations), bn)
	})

	t.Run("test we can mine blocks with strictly N+ transactions", func(t *testing.T) {
		sendTransactionEvery := 500 * time.Millisecond
		txnInBlock := int64(10)
		iterations := 10
		ac, err := StartAnvil([]string{"--no-mine"})
		require.NoError(t, err)
		client, err := ethclient.Dial(ac.URL)
		require.NoError(t, err)
		stopTxns, _ := sendTestTransactions(t, client, sendTransactionEvery, big.NewInt(1e9), big.NewInt(1e9), false)
		pm := NewRemoteAnvilMiner(ac.URL)
		pm.MineBatch(txnInBlock, 1*time.Second, 1*time.Minute)
		time.Sleep(sendTransactionEvery * time.Duration(iterations) * 2)
		pm.Stop()
		stopTxns <- struct{}{}
		bn, err := client.BlockNumber(context.Background())
		require.NoError(t, err)
		for i := 1; i <= int(bn); i++ {
			block, err := client.BlockByNumber(context.Background(), big.NewInt(int64(i)))
			require.NoError(t, err)
			require.GreaterOrEqual(t, int64(block.Transactions().Len()), txnInBlock)
		}
	})
}
