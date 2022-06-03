package contracts

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type FluxAggregator interface {
	Address() string
	Fund(ethAmount *big.Float) error
	LatestRoundID(ctx context.Context) (*big.Int, error)
	LatestRoundData(ctx context.Context) (RoundData, error)
	GetContractData(ctxt context.Context) (*FluxAggregatorData, error)
	UpdateAvailableFunds() error
	PaymentAmount(ctx context.Context) (*big.Int, error)
	RequestNewRound(ctx context.Context) error
	WithdrawPayment(ctx context.Context, from common.Address, to common.Address, amount *big.Int) error
	WithdrawablePayment(ctx context.Context, addr common.Address) (*big.Int, error)
	GetOracles(ctx context.Context) ([]string, error)
	SetOracles(opts FluxAggregatorSetOraclesOptions) error
	Description(ctxt context.Context) (string, error)
	SetRequesterPermissions(ctx context.Context, addr common.Address, authorized bool, roundsDelay uint32) error
	WatchSubmissionReceived(ctx context.Context, eventChan chan<- *SubmissionEvent) error
}

type FluxAggregatorOptions struct {
	PaymentAmount *big.Int       // The amount of LINK paid to each oracle per submission, in wei (units of 10⁻¹⁸ LINK)
	Timeout       uint32         // The number of seconds after the previous round that are allowed to lapse before allowing an oracle to skip an unfinished round
	Validator     common.Address // An optional contract address for validating external validation of answers
	MinSubValue   *big.Int       // An immutable check for a lower bound of what submission values are accepted from an oracle
	MaxSubValue   *big.Int       // An immutable check for an upper bound of what submission values are accepted from an oracle
	Decimals      uint8          // The number of decimals to offset the answer by
	Description   string         // A short description of what is being reported
}

type FluxAggregatorData struct {
	AllocatedFunds  *big.Int         // The amount of payment yet to be withdrawn by oracles
	AvailableFunds  *big.Int         // The amount of future funding available to oracles
	LatestRoundData RoundData        // Data about the latest round
	Oracles         []common.Address // Addresses of oracles on the contract
}

type FluxAggregatorSetOraclesOptions struct {
	AddList            []common.Address // oracle addresses to add
	RemoveList         []common.Address // oracle addresses to remove
	AdminList          []common.Address // oracle addresses to become admin
	MinSubmissions     uint32           // min amount of submissions in round
	MaxSubmissions     uint32           // max amount of submissions in round
	RestartDelayRounds uint32           // rounds to wait after oracles has changed
}
