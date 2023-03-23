package blockchain

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/rs/zerolog/log"
	"github.com/zksync-sdk/zksync2-go"
	"math/big"
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
	log.Info().Str("ZKSync", fmt.Sprintf("About to fund %s with %d", toAddress, transferAmt)).Msg("Executing ZKSync transaction")
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
