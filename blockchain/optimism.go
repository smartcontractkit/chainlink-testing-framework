package blockchain

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/rs/zerolog/log"
)

// OptimismMultinodeClient represents a multi-node, EVM compatible client for the Optimism network
type OptimismMultinodeClient struct {
	*EthereumMultinodeClient
}

// OptimismClient represents a single node, EVM compatible client for the Optimism network
type OptimismClient struct {
	*EthereumClient
}

func (o *OptimismClient) ReturnFunds(fromKey *ecdsa.PrivateKey) error {
	log.Warn().Str("Key", fmt.Sprintf("%x", fromKey)).Msg("Unable to return funds from Optimism at this time. Do so manually")
	return nil
}
