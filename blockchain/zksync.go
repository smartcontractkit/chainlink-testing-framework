package blockchain

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-testing-framework/utils/conversions"
)

// ZKSyncClient represents a single node, EVM compatible client for the ZKSync network
type ZKSyncClient struct {
	*EthereumClient
}

// ZKSyncMultiNodeClient represents a multi-node, EVM compatible client for the ZKSync network
type ZKSyncMultiNodeClient struct {
	*EthereumMultinodeClient
}

// TransactionOpts returns the base Tx options for 'transactions' that interact with a smart contract. Since most
// contract interactions in this framework are designed to happen through abigen calls, it's intentionally quite bare.
func (z *ZKSyncClient) TransactionOpts(from *EthereumWallet) (*bind.TransactOpts, error) {
	privateKey, err := crypto.HexToECDSA(from.PrivateKey())
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %v", err)
	}
	opts, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(z.NetworkConfig.ChainID))
	if err != nil {
		return nil, err
	}
	price, err := z.Client.SuggestGasPrice(opts.Context)
	opts.GasPrice = price
	if err != nil {
		return nil, err
	}
	opts.From = common.HexToAddress(from.Address())
	opts.Context = context.Background()

	nonce, err := z.GetNonce(context.Background(), common.HexToAddress(from.Address()), true)
	if err != nil {
		return nil, err
	}
	opts.Nonce = big.NewInt(int64(nonce))

	if z.NetworkConfig.MinimumConfirmations <= 0 { // Wait for your turn to send on an L2 chain
		<-z.NonceSettings.registerInstantTransaction(from.Address(), nonce)
	}

	return opts, nil
}

// ProcessTransaction will queue or wait on a transaction depending on whether parallel transactions are enabled
func (z *ZKSyncClient) ProcessTransaction(tx *types.Transaction) error {
	var txConfirmer HeaderEventSubscription
	if z.GetNetworkConfig().MinimumConfirmations <= 0 {
		fromAddr, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
		if err != nil {
			return err
		}
		z.NonceSettings.sentInstantTransaction(fromAddr.Hex()) // On an L2 chain, indicate the tx has been sent
		txConfirmer = NewInstantConfirmer(z, tx.Hash(), nil, nil, z.l)
	} else {
		txConfirmer = NewTransactionConfirmer(z, tx, z.GetNetworkConfig().MinimumConfirmations, z.l)
	}

	z.AddHeaderEventSubscription(tx.Hash().String(), txConfirmer)

	if !z.queueTransactions { // For sequential transactions
		log.Debug().Str("Hash", tx.Hash().String()).Msg("Waiting for TX to confirm before moving on")
		defer z.DeleteHeaderEventSubscription(tx.Hash().String())
		return txConfirmer.Wait()
	}
	return nil
}

// IsTxConfirmed checks if the transaction is confirmed on chain or not
// Temp changes until the TransactionByHash method is fixed
func (z *ZKSyncClient) IsTxConfirmed(txHash common.Hash) (bool, error) {
	isPending := false
	receipt, err := z.Client.TransactionReceipt(context.Background(), txHash)
	if err != nil {
		return !isPending, err
	}
	z.gasStats.AddClientTXData(TXGasData{
		TXHash:            txHash.String(),
		GasUsed:           receipt.GasUsed,
		CumulativeGasUsed: receipt.CumulativeGasUsed,
	})
	if receipt.Status == 0 { // 0 indicates failure, 1 indicates success
		if err != nil {
			log.Warn().Str("TX Hash", txHash.Hex()).
				Msg("Transaction failed and was reverted! Unable to retrieve reason!")
			return false, err
		}
		log.Warn().Str("TX Hash", txHash.Hex()).
			Msg("Transaction failed and was reverted!")
	}
	return !isPending, err
}

//
// func (z *ZKSyncClient) Fund(
//	toAddress string,
//	amount *big.Float,
//	gasEstimations GasEstimations,
// ) error {
//	// Connect to zkSync network
//	client, err := clients.Dial(z.GetNetworkConfig().HTTPURLs[0])
//	if err != nil {
//		return err
//	}
//	defer client.Close()
//
//	// Connect to Ethereum network
//	ethClient, err := ethclient.Dial(z.GetNetworkConfig().HTTPURLs[1])
//	if err != nil {
//		return err
//	}
//	defer ethClient.Close()
//
//	// Create wallet
//	w, err := accounts.NewWallet(common.Hex2Bytes(z.DefaultWallet.PrivateKey()), &client, ethClient)
//	if err != nil {
//		return err
//	}
//
//	opts := &accounts.TransactOpts{}
//
//	transferAmt, _ := amount.Int64()
//	log.Info().
//		Str("From", z.DefaultWallet.Address()).
//		Str("To", toAddress).
//		Str("Amount", strconv.FormatInt(transferAmt, 10)).
//		Msg("Transferring ETH")
//
//	// Create a transfer transaction
//	tx, err := w.Transfer(opts, accounts.TransferTransaction{
//		To:     common.HexToAddress(toAddress),
//		Amount: big.NewInt(transferAmt),
//		Token:  zkutils.EthAddress,
//	})
//	if err != nil {
//		return err
//	}
//	log.Info().Str("ZKSync", fmt.Sprintf("TXHash %s", tx.Hash())).Msg("Executing ZKSync transaction")
//	return z.ProcessTransaction(tx)
// }

// // ReturnFunds overrides the EthereumClient.ReturnFunds method.
// // This is needed to call the ZKSyncClient.ProcessTransaction method instead of the EthereumClient.ProcessTransaction method.
// func (z *ZKSyncClient) ReturnFunds(fromKey *ecdsa.PrivateKey) error {
//	var tx *types.Transaction
//	var err error
//	for attempt := 0; attempt < 20; attempt++ {
//		tx, err = attemptZKSyncReturn(z, fromKey, attempt)
//		if err == nil {
//			return z.ProcessTransaction(tx)
//		}
//		z.l.Debug().Err(err).Int("Attempt", attempt+1).Msg("Error returning funds from Chainlink node, trying again")
//		time.Sleep(time.Millisecond * 500)
//	}
//	return err
// }

// This is just a 1:1 copy of attemptReturn, which can't be reused as-is for ZKSync as it doesn't
// accept an interface.
func attemptZKSyncReturn(z *ZKSyncClient, fromKey *ecdsa.PrivateKey, _ int) (*types.Transaction, error) {
	to := common.HexToAddress(z.DefaultWallet.Address())
	fromAddress, err := conversions.PrivateKeyToAddress(fromKey)
	if err != nil {
		return nil, err
	}
	nonce, err := z.GetNonce(context.Background(), fromAddress, true)
	if err != nil {
		return nil, err
	}
	balance, err := z.BalanceAt(context.Background(), fromAddress)
	if err != nil {
		return nil, err
	}
	gasEstimations, err := z.EstimateGas(ethereum.CallMsg{
		From: fromAddress,
		To:   &to,
	})
	if err != nil {
		return nil, err
	}
	totalGasCost := gasEstimations.TotalGasCost
	balanceGasDelta := big.NewInt(0).Sub(balance, totalGasCost)

	if balanceGasDelta.Sign() < 0 { // Try with 0.5 gwei if we have no or negative margin. Might as well
		balanceGasDelta = big.NewInt(500_000_000)
	}

	tx, err := z.NewTx(fromKey, nonce, to, balanceGasDelta, gasEstimations)
	if err != nil {
		return nil, err
	}
	z.l.Info().
		Str("Amount", balance.String()).
		Str("From", fromAddress.Hex()).
		Str("To", to.Hex()).
		Str("Total Gas Cost", totalGasCost.String()).
		Msg("Returning Funds to Default Wallet")
	return tx, z.SendTransaction(context.Background(), tx)
}
