package contracts

import (
	"context"
	"math/big"
)

type VRF interface {
	Fund(ethAmount *big.Float) error
	ProofLength(context.Context) (*big.Int, error)
}

type VRFCoordinator interface {
	RegisterProvingKey(
		fee *big.Int,
		oracleAddr string,
		publicProvingKey [2]*big.Int,
		jobID [32]byte,
	) error
	HashOfKey(ctx context.Context, pubKey [2]*big.Int) ([32]byte, error)
	Address() string
}

type VRFConsumer interface {
	Address() string
	RequestRandomness(hash [32]byte, fee *big.Int) error
	CurrentRoundID(ctx context.Context) (*big.Int, error)
	RandomnessOutput(ctx context.Context) (*big.Int, error)
	WatchPerfEvents(ctx context.Context, eventChan chan<- *PerfEvent) error
	Fund(ethAmount *big.Float) error
}
