package blockchain

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-testing-framework/utils"
)

// ScrollMultinodeClient represents a multi-node, EVM compatible client for the Scroll network
type ScrollMultinodeClient struct {
	*EthereumMultinodeClient
}

// ScrollClient represents a single node, EVM compatible client for the Scroll network
type ScrollClient struct {
	*EthereumClient
}

// Fund sends some ETH to an address using the default wallet
func (p *ScrollClient) Fund(toAddress string, amount *big.Float) error {
	privateKey, err := crypto.HexToECDSA(p.DefaultWallet.PrivateKey())
	to := common.HexToAddress(toAddress)
	if err != nil {
		return fmt.Errorf("invalid private key: %v", err)
	}
	suggestedGasPrice, err := p.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}

	// Bump gas price
	gasPriceBuffer := big.NewInt(0).SetUint64(p.NetworkConfig.GasEstimationBuffer)
	suggestedGasPrice.Add(suggestedGasPrice, gasPriceBuffer)

	nonce, err := p.GetNonce(context.Background(), common.HexToAddress(p.DefaultWallet.Address()))
	if err != nil {
		return err
	}
	fmt.Println(ethereum.CallMsg{})
	estimatedGas, err := p.Client.EstimateGas(context.Background(), ethereum.CallMsg{})

	if err != nil {
		return err
	}

	tx, err := types.SignNewTx(privateKey, types.LatestSignerForChainID(p.GetChainID()), &types.LegacyTx{
		Nonce:    nonce,
		To:       &to,
		Value:    utils.EtherToWei(amount),
		GasPrice: suggestedGasPrice,
		Gas:      estimatedGas,
	})
	if err != nil {
		return err
	}

	log.Info().
		Str("Token", "ETH").
		Str("From", p.DefaultWallet.Address()).
		Str("To", toAddress).
		Str("Amount", amount.String()).
		Uint64("Estimated Gas Cost", new(big.Int).Mul(suggestedGasPrice, new(big.Int).SetUint64(estimatedGas)).Uint64()).
		Msg("Funding Address")
	if err := p.SendTransaction(context.Background(), tx); err != nil {
		return err
	}

	return p.ProcessTransaction(tx)
}

// DeployContract acts as a general contract deployment tool to an EVM chain
func (p *ScrollClient) DeployContract(
	contractName string,
	deployer ContractDeployer,
) (*common.Address, *types.Transaction, interface{}, error) {
	suggestedGasPrice, err := p.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, nil, nil, err
	}

	// Bump gas price
	gasPriceBuffer := big.NewInt(0).SetUint64(p.NetworkConfig.GasEstimationBuffer)
	suggestedGasPrice.Add(suggestedGasPrice, gasPriceBuffer)

	opts, err := p.TransactionOpts(p.DefaultWallet)
	if err != nil {
		return nil, nil, nil, err
	}
	opts.GasPrice = suggestedGasPrice

	contractAddress, transaction, contractInstance, err := deployer(opts, p.Client)
	if err != nil {
		return nil, nil, nil, err
	}

	if err = p.ProcessTransaction(transaction); err != nil {
		return nil, nil, nil, err
	}

	log.Info().
		Str("Contract Address", contractAddress.Hex()).
		Str("Contract Name", contractName).
		Str("From", p.DefaultWallet.Address()).
		Str("Total Gas Cost (ETH)", utils.WeiToEther(transaction.Cost()).String()).
		Str("Network Name", p.NetworkConfig.Name).
		Msg("Deployed contract")
	return &contractAddress, transaction, contractInstance, err
}

// Fund sends some ETH to an address using the default wallet
func (p *ScrollClient) ReturnFunds(fromPrivateKey *ecdsa.PrivateKey) error {
	to := common.HexToAddress(p.DefaultWallet.Address())

	// Scroll uses legacy transactions and gas estimations
	suggestedGasPrice, err := p.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}
	fromAddress, err := utils.PrivateKeyToAddress(fromPrivateKey)
	if err != nil {
		return err
	}

	balance, err := p.Client.BalanceAt(context.Background(), fromAddress, nil)
	if err != nil {
		return err
	}
	balance.Sub(balance, big.NewInt(1).Mul(suggestedGasPrice, big.NewInt(21000)))

	nonce, err := p.GetNonce(context.Background(), fromAddress)
	if err != nil {
		return err
	}
	estimatedGas, err := p.Client.EstimateGas(context.Background(), ethereum.CallMsg{})
	if err != nil {
		return err
	}

	tx, err := types.SignNewTx(fromPrivateKey, types.LatestSignerForChainID(p.GetChainID()), &types.LegacyTx{
		Nonce:    nonce,
		To:       &to,
		Value:    balance,
		GasPrice: suggestedGasPrice,
		Gas:      estimatedGas,
	})
	if err != nil {
		return err
	}

	log.Info().
		Str("Token", "ETH").
		Str("From", fromAddress.Hex()).
		Str("Amount", balance.String()).
		Msg("Returning Funds to Default Wallet")
	if err := p.SendTransaction(context.Background(), tx); err != nil {
		return err
	}

	return p.ProcessTransaction(tx)
}
