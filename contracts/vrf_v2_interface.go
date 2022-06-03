package contracts

import (
	"context"
	"math/big"

	"github.com/smartcontractkit/chainlink-testing-framework/contracts/ethereum"
)

type VRFCoordinatorV2 interface {
	SetConfig(
		minimumRequestConfirmations uint16,
		maxGasLimit uint32,
		stalenessSeconds uint32,
		gasAfterPaymentCalculation uint32,
		fallbackWeiPerUnitLink *big.Int,
		feeConfig ethereum.VRFCoordinatorV2FeeConfig,
	) error
	RegisterProvingKey(
		oracleAddr string,
		publicProvingKey [2]*big.Int,
	) error
	HashOfKey(ctx context.Context, pubKey [2]*big.Int) ([32]byte, error)
	Address() string
}

type VRFConsumerV2 interface {
	Address() string
	CurrentSubscription() (uint64, error)
	CreateFundedSubscription(funds *big.Int) error
	TopUpSubscriptionFunds(funds *big.Int) error
	RequestRandomness(hash [32]byte, subID uint64, confs uint16, gasLimit uint32, numWords uint32) error
	RandomnessOutput(ctx context.Context, arg0 *big.Int) (*big.Int, error)
	GetAllRandomWords(ctx context.Context, num int) ([]*big.Int, error)
	GasAvailable() (*big.Int, error)
	Fund(ethAmount *big.Float) error
}
