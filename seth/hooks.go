package seth

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
)

// Hooks contain optional hooks that can be used to modify the behavior of the client.
// Hooks are called before and after certain operations.
// If a hook returns an error, the operation is aborted and the error is returned to the caller.
// If a hook returns nil, the operation continues as normal.
type Hooks struct {
	ContractDeployment ContractDeploymentHooks
	TxDecoding         TxDecodingHooks
}

type PreContractDeploymentHookFn func(auth *bind.TransactOpts, name string, abi abi.ABI, bytecode []byte, params ...interface{}) error
type PostContractDeploymentHookFn func(client *Client, tx *types.Transaction) error

type ContractDeploymentHooks struct {
	Pre  PreContractDeploymentHookFn
	Post PostContractDeploymentHookFn
}

type PreTxDecodingHookFn func(client *Client) error
type PostTxDecodingHookFn func(client *Client, decodedTx *DecodedTransaction, decodedErr error) error

type TxDecodingHooks struct {
	Pre  PreTxDecodingHookFn
	Post PostTxDecodingHookFn
}
