package blockchain

import (
	"context"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"

	ethcontracts "github.com/smartcontractkit/chainlink-testing-framework/lib/contracts/ethereum"
)

const optimismGasOracleAddress string = "0x420000000000000000000000000000000000000F"

// OptimismMultinodeClient represents a multi-node, EVM compatible client for the Optimism network
type OptimismMultinodeClient struct {
	*EthereumMultinodeClient
}

// OptimismClient represents a single node, EVM compatible client for the Optimism network
type OptimismClient struct {
	*EthereumClient
}

func (o *OptimismClient) EstimateGas(callData ethereum.CallMsg) (GasEstimations, error) {
	// Optimism is unique in its usage of an L1 data fee on top of regular gas costs. Need to call their oracle
	// https://community.optimism.io/docs/developers/build/transaction-fees/#the-l1-data-fee
	gasOracle, err := ethcontracts.NewOptimismGas(common.HexToAddress(optimismGasOracleAddress), o.Client)
	if err != nil {
		return GasEstimations{}, err
	}
	opts := &bind.CallOpts{
		From:    common.HexToAddress(o.GetDefaultWallet().Address()),
		Context: context.Background(),
	}
	l1Fee, err := gasOracle.GetL1Fee(opts, types.DynamicFeeTx{}.Data)
	if err != nil {
		return GasEstimations{}, err
	}
	gasEstimations, err := o.EthereumClient.EstimateGas(callData)
	if err != nil {
		return GasEstimations{}, err
	}
	initialEstimate := gasEstimations.TotalGasCost
	gasEstimations.TotalGasCost.Add(initialEstimate, l1Fee)
	log.Debug().
		Uint64("New Total Cost", gasEstimations.TotalGasCost.Uint64()).
		Uint64("Initial Estimate", initialEstimate.Uint64()).
		Uint64("L1 Fee", l1Fee.Uint64()).
		Msg("Adding Optimism L1 Fee")
	return gasEstimations, err
}
