package client

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"math/big"
	"testing"
	"time"
)

type AnvilContainer struct {
	testcontainers.Container
	URL string
}

func startAnvil() (*AnvilContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:        "ghcr.io/foundry-rs/foundry",
		ExposedPorts: []string{"8545/tcp"},
		WaitingFor:   wait.ForListeningPort("8545").WithStartupTimeout(10 * time.Second),
		Entrypoint: []string{
			"anvil",
			"--host",
			"0.0.0.0",
			"--no-mine",
		},
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

func sendTestTransactions(t *testing.T, client *ethclient.Client, interval time.Duration) chan struct{} {
	firstAnvilPrivateKey := "5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a"
	secondAnvilAddress := "0x70997970C51812dc3A010C7d01b50e0d17dc79C8"
	privateKey, err := crypto.HexToECDSA(firstAnvilPrivateKey)
	require.NoError(t, err)
	publicKey := privateKey.Public()
	publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	stopCh := make(chan struct{})
	go func() {
		for {
			time.Sleep(interval)
			select {
			case <-stopCh:
				return
			default:
				nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
				require.NoError(t, err)
				value := big.NewInt(10)
				gasLimit := uint64(21000)
				toAddress := common.HexToAddress(secondAnvilAddress)
				chainID, err := client.NetworkID(context.Background())
				require.NoError(t, err)
				rawTx := &types.DynamicFeeTx{
					ChainID:   chainID,
					Nonce:     nonce,
					GasTipCap: big.NewInt(1000000000),
					GasFeeCap: big.NewInt(1000000000),
					Gas:       gasLimit,
					To:        &toAddress,
					Value:     value,
				}
				signedTx, err := types.SignNewTx(privateKey, types.LatestSignerForChainID(chainID), rawTx)
				require.NoError(t, err)
				err = client.SendTransaction(context.Background(), signedTx)
				require.NoError(t, err)
				t.Log("sent test transaction")
			}
		}
	}()
	return stopCh
}

func TestAnvilMiner(t *testing.T) {
	t.Run("mine periodically", func(t *testing.T) {
		period := 500 * time.Millisecond
		iterations := 10
		ac, err := startAnvil()
		require.NoError(t, err)
		client, err := ethclient.Dial(ac.URL)
		require.NoError(t, err)
		pm := NewAnvilMiner(ac.URL)
		pm.MinePeriodically(500 * time.Millisecond)
		time.Sleep(period * time.Duration(iterations))
		pm.Stop()
		bn, err := client.BlockNumber(context.Background())
		require.NoError(t, err)
		require.Equal(t, uint64(iterations), bn)
	})
	t.Run("mine batch", func(t *testing.T) {
		sendTransactionEvery := 500 * time.Millisecond
		txnInBlock := int64(10)
		iterations := 10
		ac, err := startAnvil()
		require.NoError(t, err)
		client, err := ethclient.Dial(ac.URL)
		require.NoError(t, err)
		stopTxns := sendTestTransactions(t, client, sendTransactionEvery)
		pm := NewAnvilMiner(ac.URL)
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
