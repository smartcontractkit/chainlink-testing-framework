package contracts

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
)

type LinkToken interface {
	Address() string
	Approve(to string, amount *big.Int) error
	Transfer(to string, amount *big.Int) error
	BalanceOf(ctx context.Context, addr string) (*big.Int, error)
	TransferAndCall(to string, amount *big.Int, data []byte) (*types.Transaction, error)
	Name(context.Context) (string, error)
}
