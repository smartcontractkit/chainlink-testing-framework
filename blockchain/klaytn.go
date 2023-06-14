package blockchain

import (
	"github.com/ethereum/go-ethereum"
)

// Handles specific issues with the Klaytn EVM chain: https://docs.klaytn.com/

// KlaytnMultinodeClient represents a multi-node, EVM compatible client for the Klaytn network
type KlaytnMultinodeClient struct {
	*EthereumMultinodeClient
}

// KlaytnClient represents a single node, EVM compatible client for the Klaytn network
type KlaytnClient struct {
	*EthereumClient
}

func (k *KlaytnClient) EstimateGas(callData ethereum.CallMsg) (GasEstimations, error) {
	// Klaytn is unique in its usage of a gas tip cap, enforcing it be the same
	// https://docs.klaytn.com/klaytn/design/transaction-fees#unit-price
	gasEstimations, err := k.EthereumClient.EstimateGas(callData)
	if err != nil {
		return GasEstimations{}, err
	}
	gasEstimations.GasTipCap = gasEstimations.GasPrice
	gasEstimations.GasFeeCap = gasEstimations.GasPrice
	return gasEstimations, err
}
