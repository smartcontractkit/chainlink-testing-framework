package client_test

import (
	"context"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/smartcontractkit/integrations-framework/client"
	"math/big"
)

var _ = FDescribe("ExtraErrorsClient", func() {
	networkInfo := client.NewNetworkConfig()

	It("transaction is underpriced", func() {
		ethClient, err := ethclient.Dial(networkInfo.URL)
		Expect(err).ShouldNot(HaveOccurred(),"cannot connect to client")

		privateKey , err := crypto.HexToECDSA(networkInfo.PrivateKeys[0])
		Expect(err).ShouldNot(HaveOccurred(),"cannot connect to client: ", err)

		publicKey := privateKey.Public()
		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
		Expect(ok).To(Equal(true),"error casting public key to ECDSA: ", err)

		accountAPubAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

		chainID, err := ethClient.NetworkID(context.Background())
		Expect(err).ShouldNot(HaveOccurred(),"could not determine Network ID: ", err)

		nonce, err := ethClient.PendingNonceAt(context.Background(), accountAPubAddress)
		Expect(err).ShouldNot(HaveOccurred(),"error retrieving account nonce: ", err)

		gasPrice, err := ethClient.SuggestGasPrice(context.Background())
		Expect(err).ShouldNot(HaveOccurred(),"error getting latest gas price: ", err)

		gasLimit := uint64(21000)

		tx := client.LegacyTxConfig{
			PrivateKey: privateKey,
			Nonce: nonce,
			ChainID: chainID,
			GasLimit: gasLimit,
			GasPrice: gasPrice,
		}

		toAddress := "0x62d7da380541bad6c50a90e932eb098e0fb26cf5"
		gasLimit = uint64(21000)
		errors := make(chan error)

		//Initiate first transaction in a non-blocking way
		go func() {
			_, err := tx.NewTransaction(ethClient, toAddress, big.NewInt(1000000000000))
			errors <- err
		}()

		tx2 := tx
		tx2.GasPrice = big.NewInt(10000)

		//Using the same nonce and lower gas price initiate a second transaction
		go func() {
			_, err := tx2.NewTransaction(ethClient, toAddress, big.NewInt(1000000000000))
			errors <- err
		}()

		tx3 := tx
		tx3.GasPrice = big.NewInt(1000)

		//Using the same nonce and lower gas price initiate a third transaction
		go func() {
			_, err := tx3.NewTransaction(ethClient, toAddress, big.NewInt(1000000000000))
			errors <- err
		}()

		tx4 := tx
		tx4.GasPrice = big.NewInt(100)

		//Using the same nonce and lower gas price initiate a fourth transaction
		go func() {
			_, err := tx4.NewTransaction(ethClient, toAddress, big.NewInt(1000000000000))
			errors <- err
		}()

		tx5 := tx
		tx5.GasPrice = big.NewInt(10)

		//Using the same nonce and lower gas price initiate a fifth transaction
		go func() {
			_, err := tx5.NewTransaction(ethClient, toAddress, big.NewInt(1000000000000))
			errors <- err
		}()

		for i := 0; i < 5; i++{
			txErr := <- errors
			if txErr != nil{
				sendError := client.NewSendError(txErr)
				Expect(sendError.IsReplacementUnderpriced()).To(Equal(true))
			}
		}
	})
})
