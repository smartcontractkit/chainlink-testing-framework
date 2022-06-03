package contracts

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// AbigenLog is an interface for abigen generated log topics
type AbigenLog interface {
	Topic() common.Hash
}

type KeeperRegistryVersion int32

const (
	RegistryVersion_1_0 KeeperRegistryVersion = iota
	RegistryVersion_1_1
	RegistryVersion_1_2
)

type RoundData struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}

type SubmissionEvent struct {
	Contract    common.Address
	Submission  *big.Int
	Round       uint32
	BlockNumber uint64
	Oracle      common.Address
}

// PerfEvent is used to get some metrics for contracts,
// it contrains roundID for Keeper/OCR/Flux tests and request id for VRF/Runlog
type PerfEvent struct {
	Contract       DeviationFlaggingValidator
	Round          *big.Int
	RequestID      [32]byte
	BlockTimestamp *big.Int
}

// DeviationFlaggingValidator contract used as an external validator,
// fox ex. in flux monitor rounds validation
type DeviationFlaggingValidator interface {
	Address() string
}

type Oracle interface {
	Address() string
	Fund(ethAmount *big.Float) error
	SetFulfillmentPermission(address string, allowed bool) error
}

type APIConsumer interface {
	Address() string
	RoundID(ctx context.Context) (*big.Int, error)
	Fund(ethAmount *big.Float) error
	Data(ctx context.Context) (*big.Int, error)
	CreateRequestTo(
		oracleAddr string,
		jobID [32]byte,
		payment *big.Int,
		url string,
		path string,
		times *big.Int,
	) error
	WatchPerfEvents(ctx context.Context, eventChan chan<- *PerfEvent) error
}

type Storage interface {
	Get(ctxt context.Context) (*big.Int, error)
	Set(*big.Int) error
}

type MockETHLINKFeed interface {
	Address() string
	LatestRoundData() (*big.Int, error)
}

type MockGasFeed interface {
	Address() string
}

type BlockHashStore interface {
	Address() string
}

// ReadAccessController is read/write access controller, just named by interface
type ReadAccessController interface {
	Address() string
	AddAccess(addr string) error
	DisableAccessCheck() error
}

// Flags flags contract interface
type Flags interface {
	Address() string
	GetFlag(ctx context.Context, addr string) (bool, error)
}
