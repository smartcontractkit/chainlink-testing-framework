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
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
)

// PolygonZkEvmMultinodeClient represents a multi-node, EVM compatible client for the Polygon zkEVM network
type PolygonZkEvmMultinodeClient struct {
	*EthereumMultinodeClient
}

// PolygonZkEvmClient represents a single node, EVM compatible client for the Polygon zkEVM network
type PolygonZkEvmClient struct {
	*EthereumClient
}

func (p *PolygonZkEvmClient) EstimateGas(callData ethereum.CallMsg) (GasEstimations, error) {
	var gasEstimations GasEstimations
	// gas := big.NewInt(100000)
	gasEstimations.GasTipCap = nil
	gasEstimations.GasFeeCap = nil
	gasEstimations.TotalGasCost = nil
	gasEstimations.GasPrice = nil
	gasEstimations.TotalGasCost = nil
	return gasEstimations, nil
}

// TransactionOpts returns the base Tx options for 'transactions' that interact with a smart contract. Since most
// contract interactions in this framework are designed to happen through abigen calls, it's intentionally quite bare.
func (e *PolygonZkEvmClient) TransactionOpts(from *EthereumWallet) (*bind.TransactOpts, error) {
	privateKey, err := crypto.HexToECDSA(from.PrivateKey())
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %v", err)
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
	if !e.NetworkConfig.SupportsEIP1559 {
		gasEstimations, err := e.EstimateGas(ethereum.CallMsg{})
		if err != nil {
			return nil, err
		}
		opts.GasPrice = gasEstimations.GasPrice
	}
	return opts, nil
}

// DeployContract acts as a general contract deployment tool to an ethereum chain
func (e *PolygonZkEvmClient) DeployContract(
	contractName string,
	deployer ContractDeployer,
) (*common.Address, *types.Transaction, interface{}, error) {
	opts, err := e.TransactionOpts(e.DefaultWallet)
	if err != nil {
		return nil, nil, nil, err
	}
	if !e.NetworkConfig.SupportsEIP1559 {
		gasEstimations, err := e.EstimateGas(ethereum.CallMsg{})
		if err != nil {
			return nil, nil, nil, err
		}
		opts.GasPrice = gasEstimations.GasPrice
	}

	contractAddress, transaction, contractInstance, err := deployer(opts, e.Client)
	if err != nil {
		if strings.Contains(err.Error(), "nonce") {
			err = errors.Wrap(err, fmt.Sprintf("using nonce %d", opts.Nonce.Uint64()))
		}
		return nil, nil, nil, err
	}

	e.l.Debug().Msg("Processing Tx")
	if err = e.ProcessTransaction(transaction); err != nil {
		return nil, nil, nil, err
	}

	e.l.Info().
		Str("Contract Address", contractAddress.Hex()).
		Str("Contract Name", contractName).
		Str("From", e.DefaultWallet.Address()).
		Str("Total Gas Cost", utils.WeiToEther(transaction.Cost()).String()).
		Str("Network Name", e.NetworkConfig.Name).
		Msg("Deployed contract")
	return &contractAddress, transaction, contractInstance, err
}

// Fund sends some ETH to an address using the default wallet
func (e *PolygonZkEvmClient) Fund(
	toAddress string,
	amount *big.Float,
	gasEstimations GasEstimations,
) error {
	privateKey, err := crypto.HexToECDSA(e.DefaultWallet.PrivateKey())
	if err != nil {
		return fmt.Errorf("invalid private key: %v", err)
	}
	to := common.HexToAddress(toAddress)

	nonce, err := e.GetNonce(context.Background(), common.HexToAddress(e.DefaultWallet.Address()))
	if err != nil {
		return err
	}
	gasEstimations.GasPrice = big.NewInt(1000000000)
	gasEstimations.GasUnits = 21000

	tx, err := e.NewTx(privateKey, nonce, to, utils.EtherToWei(amount), gasEstimations)
	if err != nil {
		return err
	}

	e.l.Info().
		Str("Token", "ETH").
		Str("From", e.DefaultWallet.Address()).
		Str("To", toAddress).
		Str("Hash", tx.Hash().Hex()).
		Uint64("Nonce", tx.Nonce()).
		Str("Network Name", e.GetNetworkName()).
		Str("Amount", amount.String()).
		Uint64("Estimated Gas Cost", tx.Cost().Uint64()).
		Msg("Funding Address")
	if err := e.SendTransaction(context.Background(), tx); err != nil {
		if strings.Contains(err.Error(), "nonce") {
			err = errors.Wrap(err, fmt.Sprintf("using nonce %d", nonce))
		}
		return err
	}

	return e.ProcessTransaction(tx)
}
