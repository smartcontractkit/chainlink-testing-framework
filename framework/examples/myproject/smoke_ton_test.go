package examples

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/wallet"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
)

// getConnectionPoolFromLiteserverURL parses a liteserver:// URL and creates a connection pool
func getConnectionPoolFromLiteserverURL(ctx context.Context, liteserverURL string) (*liteclient.ConnectionPool, error) {
	// Parse the liteserver URL, expected format: liteserver://publickey@host:port
	if !strings.HasPrefix(liteserverURL, "liteserver://") {
		return nil, fmt.Errorf("invalid liteserver URL format: expected liteserver:// prefix")
	}

	// remove the liteserver:// prefix
	urlPart := strings.TrimPrefix(liteserverURL, "liteserver://")

	// split by @ to separate publickey and host:port
	parts := strings.Split(urlPart, "@")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid liteserver URL format: expected publickey@host:port")
	}

	publicKey := parts[0]
	hostPort := parts[1]

	connectionPool := liteclient.NewConnectionPool()

	// mirror the exact logic from AddConnectionsFromConfig
	timeout := 3 * time.Second
	if dl, ok := ctx.Deadline(); ok {
		timeout = time.Until(dl)
	}

	// create personal context for the connection attempt
	connCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err := connectionPool.AddConnection(connCtx, hostPort, publicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to add liteserver connection: %w", err)
	}

	return connectionPool, nil
}

type CfgTon struct {
	BlockchainA *blockchain.Input `toml:"blockchain_a" validate:"required"`
}

func TestTonSmoke(t *testing.T) {
	in, err := framework.Load[CfgTon](t)
	require.NoError(t, err)

	bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
	require.NoError(t, err)
	// we can also explicitly terminate the container after the test
	defer bc.Container.Terminate(t.Context())

	var client ton.APIClientWrapped

	t.Run("setup:connect", func(t *testing.T) {
		// bc.Nodes[0].ExternalHTTPUrl now contains: "liteserver://publickey@host:port"
		liteserverURL := bc.Nodes[0].ExternalHTTPUrl

		// Create connection pool from liteserver URL
		connectionPool, err := getConnectionPoolFromLiteserverURL(t.Context(), liteserverURL)
		require.NoError(t, err, "Failed to create connection pool from liteserver URL")

		// Create API client
		client = ton.NewAPIClient(connectionPool).WithRetry()

		t.Run("setup:faucet", func(t *testing.T) {
			// network is already funded
			rawHlWallet, err := wallet.FromSeed(client, strings.Fields(blockchain.DefaultTonHlWalletMnemonic), wallet.HighloadV2Verified)
			require.NoError(t, err, "failed to create highload wallet")
			mcFunderWallet, err := wallet.FromPrivateKeyWithOptions(client, rawHlWallet.PrivateKey(), wallet.HighloadV2Verified, wallet.WithWorkchain(-1))
			require.NoError(t, err, "failed to create highload wallet")
			funder, err := mcFunderWallet.GetSubwallet(uint32(42))
			require.NoError(t, err, "failed to get highload subwallet")

			// double check funder address
			require.Equal(t, funder.Address().StringRaw(), blockchain.DefaultTonHlWalletAddress, "funder address mismatch")

			// check funder balance
			master, err := client.GetMasterchainInfo(t.Context())
			require.NoError(t, err, "failed to get masterchain info for funder balance check")
			funderBalance, err := funder.GetBalance(t.Context(), master)
			require.NoError(t, err, "failed to get funder balance")
			require.Equal(t, funderBalance.Nano().String(), "1000000000000000", "funder balance mismatch")
		})
	})
}
