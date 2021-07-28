// Package contracts handles deployment, management, and interactions of smart contracts on various chains
package contracts

import (
	"context"
	"math/big"
	"time"

	"github.com/smartcontractkit/integrations-framework/client"

	"github.com/ethereum/go-ethereum/common"
	ocrConfigHelper "github.com/smartcontractkit/libocr/offchainreporting/confighelper"
)

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

type SetOraclesOptions struct {
	AddList            []common.Address // oracle addresses to add
	RemoveList         []common.Address // oracle addresses to remove
	AdminList          []common.Address // oracle addresses to become admin
	MinSubmissions     uint32           // min amount of submissions in round
	MaxSubmissions     uint32           // max amount of submissions in round
	RestartDelayRounds uint32           // rounds to wait after oracles has changed
}

type FluxAggregator interface {
	Address() string
	Fund(fromWallet client.BlockchainWallet, ethAmount *big.Int, linkAmount *big.Int) error
	AwaitNextRoundFinalized(ctx context.Context) error
	LatestRound(ctx context.Context) (*big.Int, error)
	GetContractData(ctxt context.Context) (*FluxAggregatorData, error)
	UpdateAvailableFunds(ctx context.Context, fromWallet client.BlockchainWallet) error
	PaymentAmount(ctx context.Context) (*big.Int, error)
	RequestNewRound(ctx context.Context, fromWallet client.BlockchainWallet) error
	WithdrawPayment(ctx context.Context, caller client.BlockchainWallet, from common.Address, to common.Address, amount *big.Int) error
	WithdrawablePayment(ctx context.Context, addr common.Address) (*big.Int, error)
	GetOracles(ctx context.Context) ([]string, error)
	SetOracles(client.BlockchainWallet, SetOraclesOptions) error
	Description(ctxt context.Context) (string, error)
	SetRequesterPermissions(ctx context.Context, fromWallet client.BlockchainWallet, addr common.Address, authorized bool, roundsDelay uint32) error
}

type LinkToken interface {
	Address() string
	BalanceOf(ctx context.Context, addr common.Address) (*big.Int, error)
	Fund(fromWallet client.BlockchainWallet, ethAmount *big.Int) error
	Name(context.Context) (string, error)
}

type OffchainOptions struct {
	MaximumGasPrice           uint32         // The highest gas price for which transmitter will be compensated
	ReasonableGasPrice        uint32         // The transmitter will receive reward for gas prices under this value
	MicroLinkPerEth           uint32         // The reimbursement per ETH of gas cost, in 1e-6LINK units
	LinkGweiPerObservation    uint32         // The reward to the oracle for contributing an observation to a successfully transmitted report, in 1e-9LINK units
	LinkGweiPerTransmission   uint32         // The reward to the transmitter of a successful report, in 1e-9LINK units
	MinimumAnswer             *big.Int       // The lowest answer the median of a report is allowed to be
	MaximumAnswer             *big.Int       // The highest answer the median of a report is allowed to be
	BillingAccessController   common.Address // The access controller for billing admin functions
	RequesterAccessController common.Address // The access controller for requesting new rounds
	Decimals                  uint8          // Answers are stored in fixed-point format, with this many digits of precision
	Description               string         // A short description of what is being reported
}

// https://uploads-ssl.webflow.com/5f6b7190899f41fb70882d08/603651a1101106649eef6a53_chainlink-ocr-protocol-paper-02-24-20.pdf
type OffChainAggregatorConfig struct {
	DeltaProgress    time.Duration // The duration in which a leader must achieve progress or be replaced
	DeltaResend      time.Duration // The interval at which nodes resend NEWEPOCH messages
	DeltaRound       time.Duration // The duration after which a new round is started
	DeltaGrace       time.Duration // The duration of the grace period during which delayed oracles can still submit observations
	DeltaC           time.Duration // Limits how often updates are transmitted to the contract as long as the median isn’t changing by more then AlphaPPB
	AlphaPPB         uint64        // Allows larger changes of the median to be reported immediately, bypassing DeltaC
	DeltaStage       time.Duration // Used to stagger stages of the transmission protocol. Multiple Ethereum blocks must be mineable in this period
	RMax             uint8         // The maximum number of rounds in an epoch
	S                []int         // Transmission Schedule
	F                int           // The allowed number of "bad" oracles
	N                int           // The number of oracles
	OracleIdentities []ocrConfigHelper.OracleIdentityExtra
}

type OffchainAggregatorData struct {
	LatestRoundData RoundData // Data about the latest round
}

type OffchainAggregator interface {
	Address() string
	Fund(client.BlockchainWallet, *big.Int, *big.Int) error
	GetContractData(ctxt context.Context) (*OffchainAggregatorData, error)
	SetConfig(
		fromWallet client.BlockchainWallet,
		chainlinkNodes []client.Chainlink,
		ocrConfig OffChainAggregatorConfig,
	) error
	SetPayees(client.BlockchainWallet, []common.Address, []common.Address) error
	RequestNewRound(fromWallet client.BlockchainWallet) error
	Link(ctxt context.Context) (common.Address, error)
	GetLatestAnswer(ctxt context.Context) (*big.Int, error)
	GetLatestRound(ctxt context.Context) (*RoundData, error)
}

type Oracle interface {
	Address() string
	Fund(fromWallet client.BlockchainWallet, ethAmount *big.Int, linkAmount *big.Int) error
	SetFulfillmentPermission(fromWallet client.BlockchainWallet, address string, allowed bool) error
}

type APIConsumer interface {
	Address() string
	Fund(fromWallet client.BlockchainWallet, ethAmount *big.Int, linkAmount *big.Int) error
	Data(ctx context.Context) (*big.Int, error)
	CreateRequestTo(
		fromWallet client.BlockchainWallet,
		oracleAddr string,
		jobID [32]byte,
		payment *big.Int,
		url string,
		path string,
		times *big.Int,
	) error
}

type Storage interface {
	Get(ctxt context.Context) (*big.Int, error)
	Set(*big.Int) error
}

type VRF interface {
	Fund(client.BlockchainWallet, *big.Int, *big.Int) error
	ProofLength(context.Context) (*big.Int, error)
}

type RoundData struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}
