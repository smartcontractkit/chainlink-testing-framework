package blockchain

import (
	"math/big"

	"github.com/ethereum/go-ethereum"
)


// BotanixMultinodeClient represents a multi-node, EVM compatible client for the Botanix network
type BotanixMultinodeClient struct {
	*EthereumMultinodeClient
}

// BotanixClient represents a single node, EVM compatible client for the Botanix network
type BotanixClient struct {
	*EthereumClient
}

func (b *BotanixClient) EstimateGas(callData ethereum.CallMsg) (GasEstimations, error) {
	gasEstimations, err := b.EthereumClient.EstimateGas(callData)
	if err != nil {
		return GasEstimations{}, err
	}
	gasEstimations.GasTipCap = big.NewInt(0)
	
	return gasEstimations, err
}
