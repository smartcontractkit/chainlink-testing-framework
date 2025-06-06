package examples

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/wallet"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
)

type CfgTonParallel struct {
	BlockchainA *blockchain.Input `toml:"blockchain_a" validate:"required"`
	BlockchainB *blockchain.Input `toml:"blockchain_b" validate:"required"`
}

type DeploymentResult struct {
	Name       string
	Blockchain *blockchain.Output
	Client     ton.APIClientWrapped
	Error      error
	StartTime  time.Time
	EndTime    time.Time
}

func TestTonParallel(t *testing.T) {
	in, err := framework.Load[CfgTonParallel](t)
	require.NoError(t, err)

	configs := map[string]*blockchain.Input{
		"blockchain_a": in.BlockchainA,
		"blockchain_b": in.BlockchainB,
	}

	var mu sync.Mutex
	results := make(map[string]DeploymentResult)
	overallStartTime := time.Now()

	// Deploy all chains in parallel
	for name, config := range configs {
		t.Run("deploy_"+name, func(t *testing.T) {
			t.Parallel()

			result := deployTonInstance(name, config)

			mu.Lock()
			results[name] = result
			mu.Unlock()

			if result.Error != nil {
				t.Errorf("Failed to deploy %s: %v", result.Name, result.Error)
				return
			}

			deploymentDuration := result.EndTime.Sub(result.StartTime)
			t.Logf("✅ Successfully deployed %s in %v", result.Name, deploymentDuration)

			// Basic connectivity test
			validateConnectivity(t, result)

			// Wallet test
			validateWallet(t, result)
		})
	}

	overallDuration := time.Since(overallStartTime)
	t.Logf("🎉 All %d blockchain deployments completed in %v", len(configs), overallDuration)

	// Validate port isolation
	var successfulResults []DeploymentResult
	for _, result := range results {
		if result.Error == nil {
			successfulResults = append(successfulResults, result)
		}
	}

	if len(successfulResults) == len(configs) {
		validatePortIsolation(t, successfulResults)
	}
}

func deployTonInstance(name string, config *blockchain.Input) DeploymentResult {
	result := DeploymentResult{
		Name:      name,
		StartTime: time.Now(),
	}

	bc, err := blockchain.NewBlockchainNetwork(config)
	if err != nil {
		result.Error = fmt.Errorf("failed to create blockchain network: %w", err)
		result.EndTime = time.Now()
		return result
	}

	result.Blockchain = bc

	connectionPool := liteclient.NewConnectionPool()
	cfg, err := liteclient.GetConfigFromUrl(context.Background(),
		fmt.Sprintf("http://%s/localhost.global.config.json", bc.Nodes[0].ExternalHTTPUrl))
	if err != nil {
		result.Error = fmt.Errorf("failed to get config from URL: %w", err)
		result.EndTime = time.Now()
		return result
	}

	err = connectionPool.AddConnectionsFromConfig(context.Background(), cfg)
	if err != nil {
		result.Error = fmt.Errorf("failed to add connections from config: %w", err)
		result.EndTime = time.Now()
		return result
	}

	result.Client = ton.NewAPIClient(connectionPool).WithRetry()
	result.EndTime = time.Now()
	return result
}

// validateConnectivity tests basic blockchain connectivity
func validateConnectivity(t *testing.T, result DeploymentResult) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	master, err := result.Client.GetMasterchainInfo(ctx)
	require.NoError(t, err, "Should be able to get masterchain info for %s", result.Name)
	require.NotNil(t, master, "Masterchain info should not be nil for %s", result.Name)

	t.Logf("✓ %s: Connected, seqno: %d", result.Name, master.SeqNo)
}

func validateWallet(t *testing.T, result DeploymentResult) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	rawHlWallet, err := wallet.FromSeed(result.Client,
		strings.Fields(blockchain.DefaultTonHlWalletMnemonic), wallet.HighloadV2Verified)
	require.NoError(t, err, "Failed to create wallet for %s", result.Name)

	mcFunderWallet, err := wallet.FromPrivateKeyWithOptions(result.Client,
		rawHlWallet.PrivateKey(), wallet.HighloadV2Verified, wallet.WithWorkchain(-1))
	require.NoError(t, err, "Failed to create funder wallet for %s", result.Name)

	funder, err := mcFunderWallet.GetSubwallet(uint32(42))
	require.NoError(t, err, "Failed to get subwallet for %s", result.Name)

	require.Equal(t, funder.Address().StringRaw(), blockchain.DefaultTonHlWalletAddress,
		"Funder address mismatch for %s", result.Name)

	master, err := result.Client.GetMasterchainInfo(ctx)
	require.NoError(t, err, "Failed to get masterchain info for %s", result.Name)

	funderBalance, err := funder.GetBalance(ctx, master)
	require.NoError(t, err, "Failed to get funder balance for %s", result.Name)
	require.Equal(t, funderBalance.Nano().String(), "1000000000000000",
		"Funder balance mismatch for %s", result.Name)

	t.Logf("✓ %s: Wallet OK, balance: %s TON", result.Name, funderBalance.String())
}

func validatePortIsolation(t *testing.T, results []DeploymentResult) {
	portUsage := make(map[string][]string)

	for _, result := range results {
		httpUrl := result.Blockchain.Nodes[0].ExternalHTTPUrl
		parts := strings.Split(httpUrl, ":")
		if len(parts) == 2 {
			port := parts[1]
			portUsage[port] = append(portUsage[port], result.Name)
		}
	}

	conflicts := 0
	for port, chains := range portUsage {
		if len(chains) > 1 {
			t.Errorf("❌ Port conflict! Port %s used by: %v", port, chains)
			conflicts++
		}
	}

	if conflicts == 0 {
		t.Logf("🎯 Perfect port isolation across %d chains!", len(results))
	}
}
