package contracts

import (
	"context"
	"math/big"

	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
)

type KeeperRegistryVersion int32

const (
	RegistryVersion_1_0 KeeperRegistryVersion = iota
	RegistryVersion_1_1
	RegistryVersion_1_2
)

type UpkeepRegistrar interface {
	Address() string
	SetRegistrarConfig(
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
	Fund(ethAmount *big.Float) error
}

type KeeperRegistry interface {
	Address() string
	Fund(ethAmount *big.Float) error
	SetConfig(config KeeperRegistrySettings) error
	SetRegistrar(registrarAddr string) error
	AddUpkeepFunds(id *big.Int, amount *big.Int) error
	GetUpkeepInfo(ctx context.Context, id *big.Int) (*UpkeepInfo, error)
	GetKeeperInfo(ctx context.Context, keeperAddr string) (*KeeperInfo, error)
	SetKeepers(keepers []string, payees []string) error
	GetKeeperList(ctx context.Context) ([]string, error)
	RegisterUpkeep(target string, gasLimit uint32, admin string, checkData []byte) error
	CancelUpkeep(id *big.Int) error
	SetUpkeepGasLimit(id *big.Int, gas uint32) error
	ParseUpkeepIdFromRegisteredLog(log *types.Log) (*big.Int, error)
}

type KeeperConsumer interface {
	Address() string
	Fund(ethAmount *big.Float) error
	Counter(ctx context.Context) (*big.Int, error)
}

type UpkeepCounter interface {
	Address() string
	Fund(ethAmount *big.Float) error
	Counter(ctx context.Context) (*big.Int, error)
	SetSpread(testRange *big.Int, interval *big.Int) error
}

type UpkeepPerformCounterRestrictive interface {
	Address() string
	Fund(ethAmount *big.Float) error
	Counter(ctx context.Context) (*big.Int, error)
	SetSpread(testRange *big.Int, interval *big.Int) error
}

// KeeperConsumerPerformance is a keeper consumer contract that is more complicated than the typical consumer,
// it's intended to only be used for performance tests.
type KeeperConsumerPerformance interface {
	Address() string
	Fund(ethAmount *big.Float) error
	CheckEligible(ctx context.Context) (bool, error)
	GetUpkeepCount(ctx context.Context) (*big.Int, error)
	SetCheckGasToBurn(ctx context.Context, gas *big.Int) error
	SetPerformGasToBurn(ctx context.Context, gas *big.Int) error
}

// KeeperRegistryOpts opts to deploy keeper registry version
type KeeperRegistryOpts struct {
	RegistryVersion KeeperRegistryVersion
	LinkAddr        string
	ETHFeedAddr     string
	GasFeedAddr     string
	TranscoderAddr  string
	RegistrarAddr   string
	Settings        KeeperRegistrySettings
}

// KeeperRegistrySettings represents the settins to fine tune keeper registry
type KeeperRegistrySettings struct {
	PaymentPremiumPPB    uint32   // payment premium rate oracles receive on top of being reimbursed for gas, measured in parts per billion
	FlatFeeMicroLINK     uint32   // flat fee charged for each upkeep
	BlockCountPerTurn    *big.Int // number of blocks each oracle has during their turn to perform upkeep before it will be the next keeper's turn to submit
	CheckGasLimit        uint32   // gas limit when checking for upkeep
	StalenessSeconds     *big.Int // number of seconds that is allowed for feed data to be stale before switching to the fallback pricing
	GasCeilingMultiplier uint16   // multiplier to apply to the fast gas feed price when calculating the payment ceiling for keepers
	MinUpkeepSpend       *big.Int // minimum spend required by an upkeep before they can withdraw funds
	MaxPerformGas        uint32   // max gas allowed for an upkeep within perform
	FallbackGasPrice     *big.Int // gas price used if the gas price feed is stale
	FallbackLinkPrice    *big.Int // LINK price used if the LINK price feed is stale
}

// KeeperRegistrarSettings represents settings for registrar contract
type KeeperRegistrarSettings struct {
	AutoRegister     bool
	WindowSizeBlocks uint32
	AllowedPerWindow uint16
	RegistryAddr     string
	MinLinkJuels     *big.Int
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
