package contracts

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type FluxAggregator interface {
	Description(context.Context) (string, error)
}

type LinkToken interface {
	Name(context.Context) (string, error)
}

type OffchainAggregator interface {
	Link(ctxt context.Context) (common.Address, error)
}

type Storage interface {
	Get(context.Context) (*big.Int, error)
	Set(context.Context, *big.Int) error
}

type VRF interface {
	ProofLength(context.Context) (*big.Int, error)
}
