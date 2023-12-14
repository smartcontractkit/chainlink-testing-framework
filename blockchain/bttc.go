package blockchain

import (
	"math/big"

	"github.com/ethereum/go-ethereum"
)

// BTTCMultinodeClient represents a multi-node, EVM compatible client for the BTTC network
type BTTCMultinodeClient struct {
	*EthereumMultinodeClient
}

// BTTCClient represents a single node, EVM compatible client for the BTTC network
type BTTCClient struct {
	*EthereumClient
}

func (k *BTTCClient) EstimateGas(callData ethereum.CallMsg) (GasEstimations, error) {
	gasEstimations, err := k.EthereumClient.EstimateGas(callData)
	if err != nil {
		return GasEstimations{}, err
	}
	multiplier := big.NewInt(1000)
	gasEstimations.GasUnits = 1500000
	gasEstimations.GasPrice.Mul(gasEstimations.GasPrice, multiplier)
	return gasEstimations, err
}
