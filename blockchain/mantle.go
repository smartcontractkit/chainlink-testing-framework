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

var multiplier *big.Int = big.NewInt(1000)

// MantleMultinodeClient represents a multi-node, EVM compatible client for the Mantle network
type MantleMultinodeClient struct {
	*EthereumMultinodeClient
}

// MantleClient represents a single node, EVM compatible client for the Mantle network
type MantleClient struct {
	*EthereumClient
}

func (b *MantleClient) EstimateGas(callData ethereum.CallMsg) (GasEstimations, error) {
	gasEstimations, err := b.EthereumClient.EstimateGas(callData)
	if err != nil {
		return GasEstimations{}, err
	}
	gasEstimations.GasPrice.Mul(gasEstimations.GasPrice, multiplier)

	return gasEstimations, err
}

// DeployContract acts as a general contract deployment tool to an ethereum chain
func (b *MantleClient) DeployContract(
	contractName string,
	deployer ContractDeployer,
) (*common.Address, *types.Transaction, interface{}, error) {
	opts, err := b.TransactionOpts(b.DefaultWallet)
	if err != nil {
		return nil, nil, nil, err
	}
	fmt.Printf("%+v\n", opts)
	opts.GasPrice, err = b.EstimateGasPrice()
	fmt.Printf("Gas Price: %s\n", opts.GasPrice.String())
	if err != nil {
		return nil, nil, nil, err
	}
	opts.GasPrice.Mul(opts.GasPrice, multiplier)
	fmt.Printf("Gas Price: %s\n", opts.GasPrice.String())
	fmt.Printf("Gas Limit: %d\n", opts.GasLimit)
	// opts.GasPrice.Add(opts.GasPrice, big.NewInt(1))
	// fmt.Printf("Gas Price: %s\n", opts.GasPrice.String())
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
func (b *MantleClient) TransactionOpts(from *EthereumWallet) (*bind.TransactOpts, error) {
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
	opts.GasLimit = 500_000
	opts.GasPrice, err = b.EstimateGasPrice()
	if err != nil {
		return nil, err
	}
	opts.GasPrice.Mul(opts.GasPrice, multiplier)
	return opts, nil
}
