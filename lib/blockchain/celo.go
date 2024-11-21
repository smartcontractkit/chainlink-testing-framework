package blockchain

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/conversions"
)

// Handles specific issues with the Celo EVM chain: https://docs.celo.org/

// CeloMultinodeClient represents a multi-node, EVM compatible client for the Celo network
type CeloMultinodeClient struct {
	*EthereumMultinodeClient
}

// CeloClient represents a single node, EVM compatible client for the Celo network
type CeloClient struct {
	*EthereumClient
}

// DeployContract uses legacy txs for Celo to bypass Geth checking Celo headers which do not have a required
// sha3Uncles field
func (e *CeloClient) DeployContract(
	contractName string,
	deployer ContractDeployer,
) (*common.Address, *types.Transaction, interface{}, error) {
	suggestedPrice, err := e.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, nil, nil, err
	}
	gasPriceBuffer := big.NewInt(0).SetUint64(e.NetworkConfig.GasEstimationBuffer)

	opts, err := e.TransactionOpts(e.DefaultWallet)
	if err != nil {
		return nil, nil, nil, err
	}
	opts.GasPrice = suggestedPrice.Add(gasPriceBuffer, suggestedPrice)

	contractAddress, transaction, contractInstance, err := deployer(opts, e.Client)
	if err != nil {
		if strings.Contains(err.Error(), "nonce") {
			err = fmt.Errorf("using nonce %d err: %w", opts.Nonce.Uint64(), err)
		}
		return nil, nil, nil, err
	}

	if err = e.ProcessTransaction(transaction); err != nil {
		return nil, nil, nil, err
	}

	log.Info().
		Str("Contract Address", contractAddress.Hex()).
		Str("Contract Name", contractName).
		Str("From", e.DefaultWallet.Address()).
		Str("Total Gas Cost (CELO)", conversions.WeiToEther(transaction.Cost()).String()).
		Str("Network Name", e.NetworkConfig.Name).
		Msg("Deployed contract")
	return &contractAddress, transaction, contractInstance, err
}

func (e *CeloClient) TransactionOpts(from *EthereumWallet) (*bind.TransactOpts, error) {
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
	if nonce > math.MaxInt64 {
		return nil, fmt.Errorf("nonce value %d exceeds int64 range", nonce)
	}
	opts.Nonce = big.NewInt(int64(nonce))

	gasPrice, err := e.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}
	gasPriceBuffer := big.NewInt(0).SetUint64(e.NetworkConfig.GasEstimationBuffer)
	opts.GasPrice = gasPrice.Add(gasPriceBuffer, gasPrice)

	if e.NetworkConfig.MinimumConfirmations <= 0 { // Wait for your turn to send on an L2 chain
		<-e.NonceSettings.registerInstantTransaction(from.Address(), nonce)
	}
	return opts, nil
}
