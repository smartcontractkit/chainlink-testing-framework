package errors_test

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/smartcontractkit/integrations-framework/client"
)

var _ = Describe("Geth Errors @geth_errors", func() {
	var (
		err         error
		sendError   *client.SendError
		account1    *client.EthereumWallet
		account2    *client.EthereumWallet
		nonce1      uint64
		nonce2      uint64
		acc1Bal     = big.NewInt(0)
		acc2Bal     = big.NewInt(0)
		OneEth      = big.NewInt(1e18)
		ec          *client.EthereumClient
		networkInfo = client.NewNetworkConfig()
		toAddress = "0x62d7da380541bad6c50a90e932eb098e0fb26cf5"
	)

	BeforeEach(func() {
		By("Load network config and client network", func() {
			networkInfo = client.NewNetworkConfig()
			ec, err = client.NewEthereumClient(&networkInfo)
			Expect(err).ShouldNot(HaveOccurred(), "Connecting to network shouldn't fail")
			account1 = ec.DefaultWallet
			account2 = ec.Wallets[1]
		})
		By("Get balance for account 1", func() {
			bal1, err := ec.Client.BalanceAt(context.Background(), common.HexToAddress(account1.Address()), nil)
			Expect(err).ShouldNot(HaveOccurred(), "Unable to get balance for account 1")
			acc1Bal.Div(bal1, OneEth)
		})
		By("Get balance for account 2", func() {
			bal2, err := ec.Client.BalanceAt(context.Background(), common.HexToAddress(account2.Address()), nil)
			Expect(err).ShouldNot(HaveOccurred(), "Unable to get balance for account 1")
			acc2Bal.Div(bal2, OneEth)
		})
	})
	Describe("Error Scenario Designs", func() {
		It("insufficient funds for transfer", func() {
			amount := big.NewFloat(11)
			amountFloat, _ := amount.Float64()
			if float64(acc2Bal.Int64()) > amountFloat {
				// It("Transfer account 1 balance to account 2", func() {
				_, err = ec.SendTransaction(
					account2, common.HexToAddress(account1.Address()),
					big.NewFloat(float64(acc2Bal.Int64())))
				Expect(err).ShouldNot(HaveOccurred(), "Unable to reduce balance 2 balance to 0")
				// })
			}
			to := common.HexToAddress(account1.Address())
			_, err = ec.SendTransaction(
				account2, to,
				amount.Mul(amount, client.OneEth))
			sendError = client.NewSendError(err)
			Expect(sendError.IsInsufficientEth()).To(BeTrue(), fmt.Sprintf("Did not handle error - %v", err.Error()))
		})

		It("nonce too low", func() {
			amount := big.NewFloat(1)
			nonce1, err = ec.GetNonce(context.Background(), common.HexToAddress(account1.Address()))
			Expect(err).ShouldNot(HaveOccurred(), "Unable to get nonce")
			_, err = ec.SendTransactionWithConfig(
				account1, common.HexToAddress(account2.Address()),
				amount.Mul(amount, client.OneEth),
				&client.TxConfig{
					Nonce: nonce1 - 1,
				},
			)
			sendError = client.NewSendError(err)
			Expect(sendError.IsNonceTooLowError()).To(BeTrue(), fmt.Sprintf("Unable to handle %v", err.Error()))
		})
		// 	// max number that can cause an overflow 9,223,372,036,854,775,807. might not be possibl to reach that
		It("known transaction|already known", func() {
			amount := big.NewFloat(1)
			nonce1, err = ec.GetNonce(context.Background(), common.HexToAddress(account1.Address()))
			Expect(err).ShouldNot(HaveOccurred(), "Unable to get nonce 1")
			nonce2, err = ec.GetNonce(context.Background(), common.HexToAddress(account1.Address()))
			Expect(err).ShouldNot(HaveOccurred(), "Unable to get nonce 2")
			_, err = ec.SendTransactionWithNonce(
				account1, common.HexToAddress(account2.Address()),
				amount.Mul(amount, client.OneEth),
				nonce2,
			)
			sendError = client.NewSendError(err)
			Expect(sendError.IsTransactionAlreadyInMempool()).To(BeTrue(), fmt.Sprintf("Unable to handle %v", err.Error()))
		})
		//To reproduce "invalid sender error" ---
		//Sender function in core/types/transaction_signing.go compares the NetworkID with state NetworkId
		//changing ChainId from the network gives this error
		It("invalid sender", func() {
			amount := big.NewFloat(11)
			_, err = ec.SendTransactionWithConfig(
				account2, common.HexToAddress(account1.Address()),
				amount.Mul(amount, client.OneEth),
				&client.TxConfig{
					ChainID: 1224, // this produces error
				})
			sendError = client.NewSendError(err)
			Expect(sendError.Fatal()).To(BeTrue(), fmt.Sprintf("Unable to handle %v", err.Error()))
		})
		It("tx fee exceeds the configured cap", func() {
			amount := big.NewFloat(11)
			ecdsaKey, _ := crypto.HexToECDSA(account1.PrivateKey())
			value, _ := amount.Mul(amount, client.OneEth).Int(nil)
			suggestedGasPrice, _ := ec.Client.SuggestGasPrice(context.Background())
			gasTip, _ := ec.Client.SuggestGasTipCap(context.Background())
			nonce1, _ = ec.GetNonce(context.Background(), common.HexToAddress(account1.Address()))
			to := common.HexToAddress(account2.Address())
			tx, _ := types.SignNewTx(ecdsaKey, types.NewEIP2930Signer(big.NewInt(networkInfo.ChainID)),
				&types.LegacyTx{
					To:       &to,
					Value:    value,
					Data:     nil,
					Gas:      gasTip.Mul(gasTip, big.NewInt(3)).Uint64(), // 21000,
					GasPrice: suggestedGasPrice,
					Nonce:    nonce1,
				})
			_, err = ec.SendTransactionWithConfig(
				account2, common.HexToAddress(account1.Address()),
				amount.Mul(amount, client.OneEth),
				&client.TxConfig{
					Tx: tx,
				})
			sendError = client.NewSendError(err)
			Expect(sendError.IsTooExpensive()).To(BeTrue(), fmt.Sprintf("Unable to handle %v", err.Error()))
		})
		It("intrinsic gas too low", func() {
			amount := big.NewFloat(11)
			ecdsaKey, _ := crypto.HexToECDSA(account1.PrivateKey())
			value, _ := amount.Mul(amount, client.OneEth).Int(nil)
			suggestedGasPrice, _ := ec.Client.SuggestGasPrice(context.Background())
			gasTip, _ := ec.Client.SuggestGasTipCap(context.Background())
			gas := uint64(10)
			Expect(gas < gasTip.Uint64()).To(BeTrue(), fmt.Sprintf("Gas - %v is not less than max gas tip - %v", gas, gasTip.Uint64()))
			nonce, _ := ec.GetNonce(context.Background(), common.HexToAddress(account1.Address()))
			to := common.HexToAddress(account2.Address())
			tx, _ := types.SignNewTx(ecdsaKey, types.NewEIP2930Signer(big.NewInt(networkInfo.ChainID)),
				&types.LegacyTx{
					To:       &to,
					Value:    value,
					Data:     nil,
					Gas:      gas,
					GasPrice: suggestedGasPrice,
					Nonce:    nonce,
				})
			_, err = ec.SendTransactionWithConfig(
				account2, common.HexToAddress(account1.Address()),
				amount.Mul(amount, client.OneEth),
				&client.TxConfig{
					Tx: tx,
				})
			sendError = client.NewSendError(err)
			Expect(sendError.Fatal()).To(BeTrue(), fmt.Sprintf("Unable to handle %v", err.Error()))
		})
		It("exceeds block gas limit", func() {
			amount := big.NewFloat(11)
			ecdsaKey, _ := crypto.HexToECDSA(account1.PrivateKey())
			value, _ := amount.Mul(amount, client.OneEth).Int(nil)
			suggestedGasPrice, _ := ec.Client.SuggestGasPrice(context.Background())
			gasTip, _ := ec.Client.SuggestGasTipCap(context.Background())
			//Gas close to GasTipCap
			gas := gasTip.Div(gasTip, big.NewInt(10)).Uint64()
			nonce, _ := ec.GetNonce(context.Background(), common.HexToAddress(account1.Address()))
			to := common.HexToAddress(account2.Address())
			tx, _ := types.SignNewTx(ecdsaKey, types.NewEIP2930Signer(big.NewInt(networkInfo.ChainID)),
				&types.LegacyTx{
					To:       &to,
					Value:    value,
					Data:     nil,
					Gas:      gas,
					GasPrice: suggestedGasPrice,
					Nonce:    nonce,
				})
			_, err = ec.SendTransactionWithConfig(
				account2, common.HexToAddress(account1.Address()),
				amount.Mul(amount, client.OneEth),
				&client.TxConfig{
					Tx: tx,
				})
			sendError = client.NewSendError(err)
			Expect(sendError.Fatal()).To(BeTrue(), fmt.Sprintf("Unable to handle %v", err.Error()))
		})
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
		It("exceeds block gas limit", func() {

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

			block, err := ethClient.BlockByNumber(context.Background(),nil)
			Expect(err).ShouldNot(HaveOccurred(),"cannot get block")


			blockGasLimit := block.GasLimit()

			tx := client.LegacyTxConfig{
				PrivateKey: privateKey,
				Nonce: nonce,
				ChainID: chainID,
				GasLimit: blockGasLimit * 2,
				GasPrice: gasPrice,
			}

			_, err = tx.NewTransaction(ethClient, toAddress, big.NewInt(1000000000000))

			sendError := client.NewSendError(err)

			Expect(sendError.Fatal()).To(Equal(true))
		})
	})

})
