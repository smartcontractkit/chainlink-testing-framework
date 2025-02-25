package seth

import (
	"context"
	"strings"

	"github.com/avast/retry-go"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

type PreDeployHook func(auth *bind.TransactOpts, name string, abi abi.ABI, bytecode []byte, params ...interface{}) error
type PostDeployHook func(client *Client, tx *types.Transaction) error

type ContractDeploymentHooks struct {
	Pre  PreDeployHook
	Post PostDeployHook
}

var RetryingPostDeployHook = func(client *Client, tx *types.Transaction) error {
	err := retry.Do(
		func() error {
			ctx, cancel := context.WithTimeout(context.Background(), client.Cfg.Network.TxnTimeout.Duration())
			_, err := bind.WaitDeployed(ctx, client.Client, tx)
			cancel()

			// let's make sure that deployment transaction was successful, before retrying
			if err != nil && !errors.Is(err, context.DeadlineExceeded) {
				ctx, cancel := context.WithTimeout(context.Background(), client.Cfg.Network.TxnTimeout.Duration())
				receipt, mineErr := bind.WaitMined(ctx, client.Client, tx)
				if mineErr != nil {
					cancel()
					return mineErr
				}
				cancel()

				if receipt.Status == 0 {
					return errors.New("deployment transaction was reverted")
				}
			}

			return err
		}, retry.OnRetry(func(i uint, retryErr error) {
			switch {
			case errors.Is(retryErr, context.DeadlineExceeded):
				replacementTx, replacementErr := prepareReplacementTransaction(client, tx)
				if replacementErr != nil {
					L.Debug().Str("Current error", retryErr.Error()).Str("Replacement error", replacementErr.Error()).Uint("Attempt", i+1).Msg("Failed to prepare replacement transaction for contract deployment. Retrying with the original one")
					return
				}
				tx = replacementTx
			default:
				// do nothing, just wait again until it's mined
			}
			L.Debug().Str("Current error", retryErr.Error()).Uint("Attempt", i+1).Msg("Waiting for contract to be deployed")
		}),
		retry.DelayType(retry.FixedDelay),
		// if gas bump retries are set to 0, we still want to retry 10 times, because what we will be retrying will be other errors (no code at address, etc.)
		// downside is that if retries are enabled and their number is low other retry errors will be retried only that number of times
		// (we could have custom logic for different retry count per error, but that seemed like an overkill, so it wasn't implemented)
		retry.Attempts(func() uint {
			if client.Cfg.GasBumpRetries() != 0 {
				return client.Cfg.GasBumpRetries()
			}
			return 10
		}()),
		retry.RetryIf(func(err error) bool {
			return strings.Contains(strings.ToLower(err.Error()), "no contract code at given address") ||
				strings.Contains(strings.ToLower(err.Error()), "no contract code after deployment") ||
				(client.Cfg.GasBumpRetries() != 0 && errors.Is(err, context.DeadlineExceeded))
		}),
	)

	if err != nil {
		// pass this specific error, so that Decode knows that it's not the actual revert reason
		_, _ = client.Decode(tx, errors.New(ErrContractDeploymentFailed))

		return wrapErrInMessageWithASuggestion(client.rewriteDeploymentError(err))
	}

	return nil
}
