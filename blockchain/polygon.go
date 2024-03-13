package blockchain

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// PolygonMultinodeClient represents a multi-node, EVM compatible client for the Klaytn network
type PolygonMultinodeClient struct {
	*EthereumMultinodeClient
}

// PolygonClient represents a single node, EVM compatible client for the Polygon network
type PolygonClient struct {
	*EthereumClient
}

// TransactionOpts returns the base Tx options for 'transactions' that interact with a smart contract. Since most
// contract interactions in this framework are designed to happen through abigen calls, it's intentionally quite bare.
func (e *PolygonClient) TransactionOpts(from *EthereumWallet) (*bind.TransactOpts, error) {
	privateKey, err := crypto.HexToECDSA(from.PrivateKey())
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %w", err)
	}
	opts, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(e.NetworkConfig.ChainID))
	if err != nil {
		return nil, err
	}
	opts.From = common.HexToAddress(from.Address())
	opts.Context = context.Background()

	nonce, err := e.GetNonce(context.Background(), common.HexToAddress(from.Address()))
	if err != nil {
		return nil, err
	}
	opts.Nonce = big.NewInt(int64(nonce))

	if e.NetworkConfig.MinimumConfirmations <= 0 { // Wait for your turn to send on an L2 chain
		<-e.NonceSettings.registerInstantTransaction(from.Address(), nonce)
	}
	// if the gas limit is less than the default gas limit, use the default
	if e.NetworkConfig.DefaultGasLimit > opts.GasLimit {
		opts.GasLimit = e.NetworkConfig.DefaultGasLimit
	}
	opts.GasFeeCap = big.NewInt(102052882926)
	opts.GasTipCap = big.NewInt(30000000000)
	if !e.NetworkConfig.SupportsEIP1559 {
		opts.GasPrice, err = e.EstimateGasPrice()
		if err != nil {
			return nil, err
		}
	}
	return opts, nil
}
