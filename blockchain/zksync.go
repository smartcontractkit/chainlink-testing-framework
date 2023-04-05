package blockchain

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/rs/zerolog/log"
	"github.com/zksync-sdk/zksync2-go"
	"math/big"
	"strconv"
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

	nonce, err := z.GetNonce(context.Background(), common.HexToAddress(from.Address()))
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
		txConfirmer = NewInstantConfirmer(z, tx.Hash(), nil, nil)
	} else {
		txConfirmer = NewTransactionConfirmer(z, tx, z.GetNetworkConfig().MinimumConfirmations)
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

func (z *ZKSyncClient) Fund(
	toAddress string,
	amount *big.Float,
) error {
	es, err := zksync2.NewEthSignerFromRawPrivateKey(common.Hex2Bytes(z.DefaultWallet.PrivateKey()), zksync2.ZkSyncChainIdMainnet)
	if err != nil {
		return err
	}
	log.Info().Str("ZKSync", fmt.Sprintf("Using L1 RPC url: %s", z.GetNetworkConfig().HTTPURLs[1])).Msg("")
	log.Info().Str("ZKSync", fmt.Sprintf("Using L2 RPC url: %s", z.GetNetworkConfig().HTTPURLs[0])).Msg("")
	zp, err := zksync2.NewDefaultProvider(z.GetNetworkConfig().HTTPURLs[0])
	if err != nil {
		return err
	}
	w, err := zksync2.NewWallet(es, zp)
	if err != nil {
		return err
	}

	ethRpc, err := rpc.Dial(z.GetNetworkConfig().HTTPURLs[1])
	ep, err := w.CreateEthereumProvider(ethRpc)
	transferAmt, _ := amount.Int64()
	log.Info().
		Str("From", z.DefaultWallet.Address()).
		Str("To", toAddress).
		Str("Amount", strconv.FormatInt(transferAmt, 10)).
		Msg("Transferring ETH")
	tx, err := ep.Deposit(
		zksync2.CreateETH(),
		big.NewInt(transferAmt),
		common.HexToAddress(toAddress),
		nil,
	)
	if err != nil {
		return err
	}
	log.Info().Str("ZKSync", fmt.Sprintf("TXHash %s", tx.Hash())).Msg("Executing ZKSync transaction")
	return nil
}
