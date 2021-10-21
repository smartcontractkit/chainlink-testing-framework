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

type FluxAggregatorSetOraclesOptions struct {
	AddList            []common.Address // oracle addresses to add
	RemoveList         []common.Address // oracle addresses to remove
	AdminList          []common.Address // oracle addresses to become admin
	MinSubmissions     uint32           // min amount of submissions in round
	MaxSubmissions     uint32           // max amount of submissions in round
	RestartDelayRounds uint32           // rounds to wait after oracles has changed
}

type SubmissionEvent struct {
	Contract    common.Address
	Submission  *big.Int
	Round       uint32
	BlockNumber uint64
	Oracle      common.Address
}

type FluxAggregator interface {
	Address() string
	Fund(fromWallet client.BlockchainWallet, ethAmount, linkAmount *big.Float) error
	LatestRoundID(ctx context.Context, blockNumber *big.Int) (*big.Int, error)
	LatestRoundData(ctx context.Context) (RoundData, error)
	GetContractData(ctxt context.Context) (*FluxAggregatorData, error)
	UpdateAvailableFunds(ctx context.Context, fromWallet client.BlockchainWallet) error
	PaymentAmount(ctx context.Context) (*big.Int, error)
	RequestNewRound(ctx context.Context, fromWallet client.BlockchainWallet) error
	WithdrawPayment(ctx context.Context, caller client.BlockchainWallet, from common.Address, to common.Address, amount *big.Int) error
	WithdrawablePayment(ctx context.Context, addr common.Address) (*big.Int, error)
	GetOracles(ctx context.Context) ([]string, error)
	SetOracles(client.BlockchainWallet, FluxAggregatorSetOraclesOptions) error
	Description(ctxt context.Context) (string, error)
	SetRequesterPermissions(ctx context.Context, fromWallet client.BlockchainWallet, addr common.Address, authorized bool, roundsDelay uint32) error
	WatchSubmissionReceived(ctx context.Context, eventChan chan<- *SubmissionEvent) error
}

type LinkToken interface {
	Address() string
	Approve(fromWallet client.BlockchainWallet, to string, amount *big.Int) error
	Transfer(fromWallet client.BlockchainWallet, to string, amount *big.Int) error
	BalanceOf(ctx context.Context, addr common.Address) (*big.Int, error)
	TransferAndCall(fromWallet client.BlockchainWallet, to string, amount *big.Int, data []byte) error
	Fund(fromWallet client.BlockchainWallet, ethAmount *big.Float) error
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
	Fund(fromWallet client.BlockchainWallet, nativeAmount, linkAmount *big.Float) error
	GetContractData(ctxt context.Context) (*OffchainAggregatorData, error)
	SetConfig(
		fromWallet client.BlockchainWallet,
		chainlinkNodes []client.Chainlink,
		ocrConfig OffChainAggregatorConfig,
	) error
	SetPayees(client.BlockchainWallet, []common.Address, []common.Address) error
	RequestNewRound(fromWallet client.BlockchainWallet) error
	GetLatestAnswer(ctxt context.Context) (*big.Int, error)
	GetLatestRound(ctxt context.Context) (*RoundData, error)
}

type Oracle interface {
	Address() string
	Fund(fromWallet client.BlockchainWallet, ethAmount, linkAmount *big.Float) error
	SetFulfillmentPermission(fromWallet client.BlockchainWallet, address string, allowed bool) error
}

type APIConsumer interface {
	Address() string
	RoundID(ctx context.Context) (*big.Int, error)
	Fund(fromWallet client.BlockchainWallet, ethAmount, linkAmount *big.Float) error
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
	WatchPerfEvents(ctx context.Context, eventChan chan<- *PerfEvent) error
}

type Storage interface {
	Get(ctxt context.Context) (*big.Int, error)
	Set(*big.Int) error
}

type VRF interface {
	Fund(fromWallet client.BlockchainWallet, ethAmount, linkAmount *big.Float) error
	ProofLength(context.Context) (*big.Int, error)
}

// JobByInstance helper struct to match job + instance ID
type JobByInstance struct {
	ID       string
	Instance string
}

type MockETHLINKFeed interface {
	Address() string
}

type MockGasFeed interface {
	Address() string
}

type UpkeepRegistrar interface {
	Address() string
	SetRegistrarConfig(
		fromWallet client.BlockchainWallet,
		autoRegister bool,
		windowSizeBlocks uint32,
		allowedPerWindow uint16,
		registryAddr string,
		minLinkJuels *big.Int,
	) error
	EncodeRegisterRequest(
		name string,
		email []byte,
		upkeepAddr string,
		gasLimit uint32,
		adminAddr string,
		checkData []byte,
		amount *big.Int,
		source uint8,
	) ([]byte, error)
	Fund(fromWallet client.BlockchainWallet, ethAmount, linkAmount *big.Float) error
}

type KeeperRegistry interface {
	Address() string
	Fund(fromWallet client.BlockchainWallet, ethAmount, linkAmount *big.Float) error
	SetRegistrar(fromWallet client.BlockchainWallet, registrarAddr string) error
	AddUpkeepFunds(fromWallet client.BlockchainWallet, id *big.Int, amount *big.Int) error
	GetUpkeepInfo(ctx context.Context, id *big.Int) (*UpkeepInfo, error)
	GetKeeperInfo(ctx context.Context, keeperAddr string) (*KeeperInfo, error)
	SetKeepers(fromWallet client.BlockchainWallet, keepers []string, payees []string) error
	GetKeeperList(ctx context.Context) ([]string, error)
	RegisterUpkeep(fromWallet client.BlockchainWallet, target string, gasLimit uint32, admin string, checkData []byte) error
}

type KeeperConsumer interface {
	Address() string
	Fund(fromWallet client.BlockchainWallet, ethAmount, linkAmount *big.Float) error
	Counter(ctx context.Context) (*big.Int, error)
}

// KeeperRegistryOpts opts to deploy keeper registry
type KeeperRegistryOpts struct {
	LinkAddr             string
	ETHFeedAddr          string
	GasFeedAddr          string
	PaymentPremiumPPB    uint32
	BlockCountPerTurn    *big.Int
	CheckGasLimit        uint32
	StalenessSeconds     *big.Int
	GasCeilingMultiplier uint16
	FallbackGasPrice     *big.Int
	FallbackLinkPrice    *big.Int
}

// KeeperInfo keeper status and balance info
type KeeperInfo struct {
	Payee   string
	Active  bool
	Balance *big.Int
}

// UpkeepInfo keeper target info
type UpkeepInfo struct {
	Target              string
	ExecuteGas          uint32
	CheckData           []byte
	Balance             *big.Int
	LastKeeper          string
	Admin               string
	MaxValidBlocknumber uint64
}

type BlockHashStore interface {
	Address() string
}

type VRFCoordinator interface {
	RegisterProvingKey(
		fromWallet client.BlockchainWallet,
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
	RequestRandomness(fromWallet client.BlockchainWallet, hash [32]byte, fee *big.Int) error
	CurrentRoundID(ctx context.Context) (*big.Int, error)
	RandomnessOutput(ctx context.Context) (*big.Int, error)
	WatchPerfEvents(ctx context.Context, eventChan chan<- *PerfEvent) error
	Fund(fromWallet client.BlockchainWallet, ethAmount, linkAmount *big.Float) error
}

type RoundData struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}

// ReadAccessController is read/write access controller, just named by interface
type ReadAccessController interface {
	Address() string
	AddAccess(fromWallet client.BlockchainWallet, addr string) error
	DisableAccessCheck(fromWallet client.BlockchainWallet) error
}

// Flags flags contract interface
type Flags interface {
	Address() string
	GetFlag(ctx context.Context, addr string) (bool, error)
}

// DeviationFlaggingValidator contract used as an external validator,
// fox ex. in flux monitor rounds validation
type DeviationFlaggingValidator interface {
	Address() string
}

// PerfEvent is used to get some metrics for contracts,
// it contrains roundID for Keeper/OCR/Flux tests and request id for VRF/Runlog
type PerfEvent struct {
	Contract       DeviationFlaggingValidator
	Round          *big.Int
	RequestID      [32]byte
	BlockTimestamp *big.Int
}

// OCRv2 contracts

// OCRv2AccessController access controller
type OCRv2AccessController interface {
	Address() string
	AddAccess(fromWallet client.BlockchainWallet, addr string) error
	RemoveAccess(fromWallet client.BlockchainWallet, addr string) error
	HasAccess(to string) (bool, error)
}

type OCRv2 interface {
	Address() string
	SetConfig(fromWallet client.BlockchainWallet) error
	TransferOwnership(fromWallet client.BlockchainWallet, to string) error
	SetBilling(fromWallet client.BlockchainWallet, observationPayment uint32, recommendedGasPrice uint32) error
	GetLatestConfigDetails() (map[string]interface{}, error)
	GetRoundData(roundID uint32) (map[string]interface{}, error)
	GetOwedPayment(transmitterAddr string) (map[string]interface{}, error)
}
