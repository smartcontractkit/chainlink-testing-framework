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
	flag "github.com/spf13/pflag"
	"github.com/testcontainers/testcontainers-go/log"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/config"
	config_types "github.com/smartcontractkit/chainlink-testing-framework/lib/config/types"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/docker/test_env"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
)

var StartTestEnvCmd = &cobra.Command{
	Use:   "start-test-env",
	Short: "Start local test environment",
}

var startPrivateChain = &cobra.Command{
	Use:   "private-chain",
	Short: "Private chain with 1 node",
	RunE:  startPrivateEthChainE,
}

const (
	Flag_EthereumVersion      = "ethereum-version"
	Flag_ConsensusLayer       = "consensus-layer"
	Flag_ExecutionLayer       = "execution-layer"
	Flag_WaitForFinalization  = "wait-for-finalization"
	Flag_ChainID              = "chain-id"
	Flag_ExecutionClientImage = "execution-layer-image"
	Flag_ConsensucClientImage = "consensus-client-image"
	Flag_ValidatorImage       = "validator-image"
)

func init() {
	StartTestEnvCmd.AddCommand(startPrivateChain)

	StartTestEnvCmd.PersistentFlags().StringP(
		Flag_EthereumVersion,
		"v",
		"eth2",
		"ethereum version (eth1, eth2) (default: eth2)",
	)

	StartTestEnvCmd.PersistentFlags().StringP(
		Flag_ConsensusLayer,
		"l",
		"prysm",
		"consensus layer (prysm) (default: prysm)",
	)

	StartTestEnvCmd.PersistentFlags().StringP(
		Flag_ExecutionLayer,
		"e",
		"geth",
		"execution layer (geth, nethermind, besu or erigon) (default: geth)",
	)

	StartTestEnvCmd.PersistentFlags().BoolP(
		Flag_WaitForFinalization,
		"w",
		false,
		"wait for finalization of at least 1 epoch (might take up to 5 minutes) default: false",
	)

	StartTestEnvCmd.PersistentFlags().IntP(
		Flag_ChainID,
		"c",
		1337,
		"chain id",
	)

	StartTestEnvCmd.PersistentFlags().String(
		Flag_ExecutionClientImage,
		"",
		"custom Docker image for execution layer client",
	)

	StartTestEnvCmd.PersistentFlags().String(
		Flag_ConsensucClientImage,
		"",
		"custom Docker image for consensus layer client",
	)

	StartTestEnvCmd.PersistentFlags().String(
		Flag_ValidatorImage,
		"",
		"custom Docker image for validator",
	)

	// Set default log level for non-testcontainer code
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Discard testcontainers logs
	log.SetDefault(defaultlog.New(io.Discard, "", defaultlog.LstdFlags))
}

func startPrivateEthChainE(cmd *cobra.Command, args []string) error {
	log := logging.GetTestLogger(nil)
	flags := cmd.Flags()

	ethereumVersion, err := flags.GetString(Flag_EthereumVersion)
	if err != nil {
		return err
	}

	ethereumVersion = strings.ToLower(ethereumVersion)

	if ethereumVersion != "eth1" && ethereumVersion != "eth2" {
		return fmt.Errorf("invalid ethereum version: %s. use 'eth1' or 'eth2'", ethereumVersion)
	}

	consensusLayer, err := flags.GetString(Flag_ConsensusLayer)
	if err != nil {
		return err
	}

	consensusLayer = strings.ToLower(consensusLayer)

	if consensusLayer != "" && consensusLayer != "prysm" {
		return fmt.Errorf("invalid consensus layer: %s. use 'prysm'", consensusLayer)
	}

	if consensusLayer != "" && ethereumVersion == "eth1" {
		log.Warn().Msg("consensus layer was set, but it has no sense for a eth1. Ignoring it")
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

	consensusLayerToUse := config.ConsensusLayer(consensusLayer)
	if consensusLayer != "" && ethereumVersion == "eth1" {
		consensusLayerToUse = ""
	}

	customDockerImages, err := getCustomImages(flags)
	if err != nil {
		return err
	}

	builder := test_env.NewEthereumNetworkBuilder()
	builder = *builder.WithEthereumVersion(config_types.EthereumVersion(ethereumVersion)).
		WithConsensusLayer(consensusLayerToUse).
		WithExecutionLayer(config_types.ExecutionLayer(executionLayer)).
		WithEthereumChainConfig(config.EthereumChainConfig{
			ValidatorCount: 8,
			SlotsPerEpoch:  2,
			SecondsPerSlot: 6,
			ChainID:        chainId,
			HardForkEpochs: map[string]int{"Deneb": 500},
		})

	if waitForFinalization {
		builder = *builder.WithWaitingForFinalization()
	}

	if len(customDockerImages) > 0 {
		builder = *builder.WithCustomDockerImages(customDockerImages)
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

	err = cfg.Save()
	if err != nil {
		return err
	}

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

func getCustomImages(flags *flag.FlagSet) (map[config.ContainerType]string, error) {
	customImages := make(map[config.ContainerType]string)
	executionClientImage, err := flags.GetString(Flag_ExecutionClientImage)
	if err != nil {
		return nil, err
	}

	if executionClientImage != "" {
		customImages[config.ContainerType_ExecutionLayer] = executionClientImage
	}

	consensusClientImage, err := flags.GetString(Flag_ConsensucClientImage)
	if err != nil {
		return nil, err
	}

	if consensusClientImage != "" {
		customImages[config.ContainerType_ConsensusLayer] = consensusClientImage
	}

	validatorImage, err := flags.GetString(Flag_ValidatorImage)
	if err != nil {
		return nil, err
	}

	if validatorImage != "" {
		customImages[config.ContainerType_ConsensusValidator] = validatorImage
	}

	return customImages, nil
}
