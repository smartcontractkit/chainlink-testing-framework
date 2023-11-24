package internal

import (
	"fmt"
	"io"
	defaultlog "log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/testcontainers/testcontainers-go"

	"github.com/smartcontractkit/chainlink-testing-framework/docker/test_env"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
)

var StartTestEnvCmd = &cobra.Command{
	Use:   "start-test-env",
	Short: "Start local test environment",
}

var startPrivateChain = &cobra.Command{
	Use:   "private-chain",
	Short: "Private chain with 1 node",
	RunE:  startBesuPrysmChainE,
}

const (
	Flag_ConsensusType       = "consensus-type"
	Flag_ConsensusLayer      = "consensus-layer"
	Flag_ExecutionLayer      = "execution-layer"
	Flag_WaitForFinalization = "wait-for-finalization"
	Flag_ChainID             = "chain-id"
)

func init() {
	StartTestEnvCmd.AddCommand(startPrivateChain)

	StartTestEnvCmd.PersistentFlags().StringP(
		Flag_ConsensusType,
		"t",
		"pos",
		"consensus type (pow or pos)",
	)

	StartTestEnvCmd.PersistentFlags().StringP(
		Flag_ConsensusLayer,
		"l",
		"prysm",
		"consensus layer (prysm)",
	)

	StartTestEnvCmd.PersistentFlags().StringP(
		Flag_ExecutionLayer,
		"e",
		"geth",
		"execution layer (geth, nethermind, besu or erigon)",
	)

	StartTestEnvCmd.PersistentFlags().BoolP(
		Flag_WaitForFinalization,
		"w",
		false,
		"wait for finalization of at least 1 epoch (might take up to 5 mintues)",
	)

	StartTestEnvCmd.PersistentFlags().IntP(
		Flag_ChainID,
		"c",
		1337,
		"chain id",
	)

	// Set default log level for non-testcontainer code
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Discard testcontainers logs
	testcontainers.Logger = defaultlog.New(io.Discard, "", defaultlog.LstdFlags)
}

func startBesuPrysmChainE(cmd *cobra.Command, args []string) error {
	log := logging.GetTestLogger(nil)
	flags := cmd.Flags()

	consensusType, err := flags.GetString(Flag_ConsensusType)
	if err != nil {
		return err
	}

	consensusType = strings.ToLower(consensusType)

	if consensusType != "pos" && consensusType != "pow" {
		return fmt.Errorf("invalid consensus type: %s. use 'pow' or 'pos'", consensusType)
	}

	consensusLayer, err := flags.GetString(Flag_ConsensusLayer)
	if err != nil {
		return err
	}

	consensusLayer = strings.ToLower(consensusLayer)

	if consensusLayer != "" && consensusLayer != "prysm" {
		return fmt.Errorf("invalid consensus layer: %s. use 'prysm'", consensusLayer)
	}

	if consensusLayer != "" && consensusType == "pow" {
		log.Warn().Msg("consensus layer was set, but it has no sense for a PoW conensus. Ignoring it")
	}

	executionLayer, err := flags.GetString(Flag_ExecutionLayer)
	if err != nil {
		return err
	}

	executionLayer = strings.ToLower(executionLayer)
	switch executionLayer {
	case "geth", "nethermind", "besu", "erigon":
	default:
		return fmt.Errorf("invalid execution layer: %s. use 'geth', 'nethermind', 'besu' or 'erigon'", executionLayer)
	}

	waitForFinalization, err := flags.GetBool(Flag_WaitForFinalization)
	if err != nil {
		return err
	}

	chainId, err := flags.GetInt(Flag_ChainID)
	if err != nil {
		return err
	}

	consensusLayerToUse := test_env.ConsensusLayer(consensusLayer)
	if consensusLayer != "" && consensusType == "pow" {
		consensusLayerToUse = ""
	}

	builder := test_env.NewEthereumNetworkBuilder()
	builder = *builder.WithConsensusType(test_env.ConsensusType(consensusType)).
		WithCustomNetworkParticipants([]test_env.EthereumNetworkParticipant{
			{
				ConsensusLayer: consensusLayerToUse,
				ExecutionLayer: test_env.ExecutionLayer(executionLayer),
				Count:          1,
			},
		}).
		WithEthereumChainConfig(test_env.EthereumChainConfig{
			ValidatorCount: 8,
			SlotsPerEpoch:  2,
			SecondsPerSlot: 6,
			ChainID:        chainId,
		})

	if !waitForFinalization {
		builder = *builder.WithoutWaitingForFinalization()
	}

	cfg, err := builder.
		Build()

	log.Info().Str("chain", cfg.Describe()).Msg("Starting private chain")

	if err != nil {
		return err
	}

	_, eth2, err := cfg.Start()

	if err != nil {
		return err
	}
	log.Info().Msg("---------- Private chain is ready ----------")
	log.Info().Msgf("Public RPC WS URLs: %v", eth2.PublicWsUrls())
	log.Info().Msgf("Public RPC HTTP URLs: %v", eth2.PublicHttpUrls())

	handleExitSignal()

	return nil
}

func handleExitSignal() {
	// Create a channel to receive exit signals
	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, os.Interrupt, syscall.SIGTERM)

	log := logging.GetTestLogger(nil)
	log.Info().Msg("Press Ctrl+C to destroy the test environment")

	// Block until an exit signal is received
	<-exitChan
}
