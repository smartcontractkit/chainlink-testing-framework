package internal

import (
	"io"
	defaultlog "log"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/testcontainers/testcontainers-go"

	"github.com/smartcontractkit/chainlink-testing-framework/docker/test_env"
)

var StartTestEnvCmd = &cobra.Command{
	Use:   "start-test-env",
	Short: "Start local test environment",
}

var startBesuPrysmChain = &cobra.Command{
	Use:   "besu-prysm-chain",
	Short: "Private Besu + Prysm chain with 1 node",
	RunE:  startBesuPrysmChainE,
}

func init() {
	StartTestEnvCmd.AddCommand(startBesuPrysmChain)

	// StartTestEnvCmd.PersistentFlags().StringP(
	// 	"config",
	// 	"c",
	// 	"",
	// 	"Path to test config (TOML)",
	// )
	// StartTestEnvCmd.MarkPersistentFlagRequired("config")

	// Set default log level for non-testcontainer code
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Discard testcontainers logs
	testcontainers.Logger = defaultlog.New(io.Discard, "", defaultlog.LstdFlags)
}

func startBesuPrysmChainE(cmd *cobra.Command, args []string) error {
	log.Info().Msg("Starting Besu + Prysm Chain..")

	builder := test_env.NewEthereumNetworkBuilder()
	cfg, err := builder.
		WithConsensusType(test_env.ConsensusType_PoS).
		WithCustomNetworkParticipants([]test_env.EthereumNetworkParticipant{
			{
				ConsensusLayer: test_env.ConsensusLayer_Prysm,
				ExecutionLayer: test_env.ExecutionLayer_Besu,
				Count:          1,
			},
		}).
		WithEthereumChainConfig(test_env.EthereumChainConfig{
			ValidatorCount: 8,
			SlotsPerEpoch:  2,
			SecondsPerSlot: 6,
		}).
		WithoutWaitingForFinalization().
		Build()
	if err != nil {
		return err
	}

	_, eth2, err := cfg.Start()

	if err != nil {
		return err
	}
	log.Info().Msg("Private chain is ready!\n")
	log.Info().Msgf("Public RPC WS URLs: %v", eth2.PublicWsUrls())
	log.Info().Msgf("Public RPC HTTP URLs: %v", eth2.PublicHttpUrls())

	handleExitSignal()

	return nil
}

func handleExitSignal() {
	// Create a channel to receive exit signals
	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, os.Interrupt, syscall.SIGTERM)

	log.Info().Msg("Press Ctrl+C to destroy the test environment")

	// Block until an exit signal is received
	<-exitChan
}
