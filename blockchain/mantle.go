package blockchain

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/conversions"
)

// BttcMultinodeClient represents a multi-node, EVM compatible client for the Kava network
type MantleGoerliMultinodeClient struct {
	*EthereumMultinodeClient
}

// BttcClient represents a single node, EVM compatible client for the Kava network
type MantleGoerliClient struct {
	*EthereumClient
}

func (b *MantleGoerliClient) EstimateGas(callData ethereum.CallMsg) (GasEstimations, error) {
	gasEstimations, err := b.EthereumClient.EstimateGas(callData)
	if err != nil {
		return GasEstimations{}, err
	}
	multiplier := big.NewInt(100000)
	// gasEstimations.GasUnits = 1500000
	gasEstimations.GasPrice.Mul(gasEstimations.GasPrice, multiplier)
	return gasEstimations, err
}

// DeployContract acts as a general contract deployment tool to an ethereum chain
func (b *MantleGoerliClient) DeployContract(
	contractName string,
	deployer ContractDeployer,
) (*common.Address, *types.Transaction, interface{}, error) {
	multiplier := big.NewInt(100000)
	opts, err := b.TransactionOpts(b.DefaultWallet)
	if err != nil {
		return nil, nil, nil, err
	}
	opts.GasPrice, err = b.EstimateGasPrice()
	if err != nil {
		return nil, nil, nil, err
	}
	opts.GasPrice.Mul(opts.GasPrice, multiplier)

	contractAddress, transaction, contractInstance, err := deployer(opts, b.Client)
	if err != nil {
		if strings.Contains(err.Error(), "nonce") {
			err = fmt.Errorf("using nonce %d err: %w", opts.Nonce.Uint64(), err)
		}
		return nil, nil, nil, err
	}

	if err = b.ProcessTransaction(transaction); err != nil {
		return nil, nil, nil, err
	}

	b.l.Info().
		Str("Contract Address", contractAddress.Hex()).
		Str("Contract Name", contractName).
		Str("From", b.DefaultWallet.Address()).
		Str("Total Gas Cost", conversions.WeiToEther(transaction.Cost()).String()).
		Str("Network Name", b.NetworkConfig.Name).
		Msg("Deployed contract")
	return &contractAddress, transaction, contractInstance, err
}

// TransactionOpts returns the base Tx options for 'transactions' that interact with a smart contract. Since most
// contract interactions in this framework are designed to happen through abigen calls, it's intentionally quite bare.
func (b *MantleGoerliClient) TransactionOpts(from *EthereumWallet) (*bind.TransactOpts, error) {
	privateKey, err := crypto.HexToECDSA(from.PrivateKey())
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %w", err)
	}
	opts, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(b.NetworkConfig.ChainID))
	if err != nil {
		return nil, err
	}
	opts.From = common.HexToAddress(from.Address())
	opts.Context = context.Background()

	nonce, err := b.GetNonce(context.Background(), common.HexToAddress(from.Address()))
	if err != nil {
		return nil, err
	}
	opts.Nonce = big.NewInt(int64(nonce))

	if b.NetworkConfig.MinimumConfirmations <= 0 { // Wait for your turn to send on an L2 chain
		<-b.NonceSettings.registerInstantTransaction(from.Address(), nonce)
	}
	// if the gas limit is less than the default gas limit, use the default
	if b.NetworkConfig.DefaultGasLimit > opts.GasLimit {
		opts.GasLimit = b.NetworkConfig.DefaultGasLimit
	}
	multiplier := big.NewInt(100000)
	opts.GasPrice, err = b.EstimateGasPrice()
	if err != nil {
		return nil, err
	}
	opts.GasPrice.Mul(opts.GasPrice, multiplier)
	return opts, nil
}
