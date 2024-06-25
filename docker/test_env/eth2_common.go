package test_env

import (
	"context"
	_ "embed" // leave me alone you linter
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
)

var (
	ETH2_EXECUTION_PORT                             = "8551"
	WALLET_PASSWORD                                 = "password"
	VALIDATOR_WALLET_PASSWORD_FILE_INSIDE_CONTAINER = fmt.Sprintf("%s/wallet_password.txt", GENERATED_DATA_DIR_INSIDE_CONTAINER)
	ACCOUNT_PASSWORD_FILE_INSIDE_CONTAINER          = fmt.Sprintf("%s/account_password.txt", GENERATED_DATA_DIR_INSIDE_CONTAINER)
	ACCOUNT_KEYSTORE_FILE_INSIDE_CONTAINER          = fmt.Sprintf("%s/account_key", KEYSTORE_DIR_LOCATION_INSIDE_CONTAINER)
	KEYSTORE_DIR_LOCATION_INSIDE_CONTAINER          = fmt.Sprintf("%s/keystore", GENERATED_DATA_DIR_INSIDE_CONTAINER)
	GENERATED_VALIDATOR_KEYS_DIR_INSIDE_CONTAINER   = "/keys"
	NODE_0_DIR_INSIDE_CONTAINER                     = fmt.Sprintf("%s/node-0", GENERATED_VALIDATOR_KEYS_DIR_INSIDE_CONTAINER)
	GENERATED_DATA_DIR_INSIDE_CONTAINER             = "/data/metadata"
	JWT_SECRET_FILE_LOCATION_INSIDE_CONTAINER       = fmt.Sprintf("%s/jwtsecret", GENERATED_DATA_DIR_INSIDE_CONTAINER) // #nosec G101
	VALIDATOR_BIP39_MNEMONIC                        = "giant issue aisle success illegal bike spike question tent bar rely arctic volcano long crawl hungry vocal artwork sniff fantasy very lucky have athlete"
)

func waitForChainToFinaliseAnEpoch(lggr zerolog.Logger, evmClient blockchain.EVMClient, timeout time.Duration) error {
	lggr.Info().Msg("Waiting for chain to finalize an epoch")

	timeoutC := time.After(timeout)
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeoutC:
			return fmt.Errorf("chain %s failed to finalize an epoch", evmClient.GetNetworkName())
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10*time.Second))
			finalized, err := evmClient.GetLatestFinalizedBlockHeader(ctx)
			if err != nil {
				if strings.Contains(err.Error(), "finalized block not found") {
					lggr.Err(err).Msgf("error getting finalized block number for %s", evmClient.GetNetworkName())
				} else {
					lggr.Warn().Msgf("no epoch finalized yet for chain %s", evmClient.GetNetworkName())
				}
			}
			cancel()

			if finalized != nil && finalized.Number.Int64() > 0 {
				lggr.Info().Msgf("Chain '%s' finalized an epoch", evmClient.GetNetworkName())
				return nil
			}
		}
	}
}
