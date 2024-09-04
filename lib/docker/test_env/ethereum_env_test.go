package test_env

import (
	"context"
	"fmt"
	"os"
	"testing"

	config_types "github.com/smartcontractkit/chainlink-testing-framework/lib/config/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pelletier/go-toml/v2"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/config"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
)

func TestEthEnvCustomConfig(t *testing.T) {
	l := logging.GetTestLogger(t)

	builder := NewEthereumNetworkBuilder()
	cfg, err := builder.
		WithEthereumVersion(config_types.EthereumVersion_Eth2).
		WithConsensusLayer(config.ConsensusLayer_Prysm).
		WithExecutionLayer(config_types.ExecutionLayer_Geth).
		WithEthereumChainConfig(config.EthereumChainConfig{
			SecondsPerSlot: 6,
			SlotsPerEpoch:  2,
		}).
		Build()
	require.NoError(t, err, "Builder validation failed")

	net, _, err := cfg.Start()
	require.NoError(t, err, "Couldn't start PoS network")

	c, err := blockchain.ConnectEVMClient(net, l)
	require.NoError(t, err, "Couldn't connect to the evm client")
	err = c.Close()
	require.NoError(t, err, "Couldn't close the client")
}

func TestEthEnvExtraFunding(t *testing.T) {
	l := logging.GetTestLogger(t)

	addressToFund := "0x14dc79964da2c08b23698b3d3cc7ca32193d9955"

	builder := NewEthereumNetworkBuilder()
	cfg, err := builder.
		//nolint:staticcheck //ignore SA1019
		WithEthereumVersion(config_types.EthereumVersion_Eth2_Legacy).
		WithConsensusLayer(config.ConsensusLayer_Prysm).
		WithExecutionLayer(config_types.ExecutionLayer_Geth).
		WithEthereumChainConfig(config.EthereumChainConfig{
			AddressesToFund: []string{addressToFund},
		}).
		Build()
	require.NoError(t, err, "Builder validation failed")

	net, _, err := cfg.Start()
	require.NoError(t, err, "Couldn't start PoS network")

	c, err := blockchain.ConnectEVMClient(net, l)
	require.NoError(t, err, "Couldn't connect to the evm client")

	balance, err := c.BalanceAt(context.Background(), common.HexToAddress(addressToFund))
	require.NoError(t, err, "Couldn't get balance")
	require.Equal(t, "11515845246265065472", fmt.Sprintf("%d", balance.Uint64()), "Balance is not correct")

	err = c.Close()
	require.NoError(t, err, "Couldn't close the client")
}

func TestEthEnvWithPrysmAndGethReuseConfig(t *testing.T) {
	l := logging.GetTestLogger(t)

	builder := NewEthereumNetworkBuilder()
	cfg, err := builder.
		WithEthereumVersion(config_types.EthereumVersion_Eth2).
		WithConsensusLayer(config.ConsensusLayer_Prysm).
		WithExecutionLayer(config_types.ExecutionLayer_Geth).
		Build()
	require.NoError(t, err, "Builder validation failed")

	_, _, err = cfg.Start()
	require.NoError(t, err, "Couldn't start PoS network")

	newBuilder := NewEthereumNetworkBuilder()
	reusedCfg, err := newBuilder.
		WithExistingConfig(cfg.EthereumNetworkConfig).
		Build()
	require.NoError(t, err, "Builder validation failed")

	net, _, err := reusedCfg.Start()
	require.NoError(t, err, "Couldn't reuse PoS network")

	c, err := blockchain.ConnectEVMClient(net, l)
	require.NoError(t, err, "Couldn't connect to the evm client")
	err = c.Close()
	require.NoError(t, err, "Couldn't close the client")
}

func TestEthEnvWithPrysmAndGethReuseFromEnv(t *testing.T) {
	t.Skip("for demonstration purposes only")
	l := logging.GetTestLogger(t)

	os.Setenv(CONFIG_ENV_VAR_NAME, "change-me-to-the-path-of-your-config-file")
	builder := NewEthereumNetworkBuilder()
	cfg, err := builder.
		WithExistingConfigFromEnvVar().
		Build()
	require.NoError(t, err, "Builder validation failed")

	net, _, err := cfg.Start()
	require.NoError(t, err, "Couldn't start PoS network")

	c, err := blockchain.ConnectEVMClient(net, l)
	require.NoError(t, err, "Couldn't connect to the evm client")
	err = c.Close()
	require.NoError(t, err, "Couldn't close the client")
}

func TestEthEnvExecClientFromToml(t *testing.T) {
	t.Parallel()
	toml := `
	[EthereumNetwork]
	consensus_type="pos"
	consensus_layer="prysm"
	execution_layer="besu"
	wait_for_finalization=false

	[EthereumNetwork.EthereumChainConfig]
	seconds_per_slot=12
	slots_per_epoch=2
	genesis_delay=20
	validator_count=8
	chain_id=1234
	addresses_to_fund=["0x742d35Cc6634C0532925a3b844Bc454e4438f44e", "0x742d35Cc6634C0532925a3b844Bc454e4438f44f"]
	`

	tomlCfg, err := readEthereumNetworkConfig(toml)
	require.NoError(t, err, "Couldn't read config")

	tomlCfg.EthereumChainConfig.GenerateGenesisTimestamp()

	err = tomlCfg.Validate()
	require.NoError(t, err, "Couldn't validate TOML config")

	builder := NewEthereumNetworkBuilder()
	cfg, err := builder.
		WithExistingConfig(tomlCfg).
		Build()
	require.NoError(t, err, "Builder validation failed")
	require.Equal(t, config_types.ExecutionLayer_Besu, *cfg.ExecutionLayer, "Execution layer should be Besu")
	require.NotNil(t, cfg.ConsensusLayer, "Consensus layer should not be nil")
	require.Equal(t, config.ConsensusLayer_Prysm, *cfg.ConsensusLayer, "Consensus layer should be Prysm")
	require.Equal(t, config_types.EthereumVersion_Eth2, *cfg.EthereumVersion, "Ethereum Version should be Eth2")
	require.NotNil(t, cfg.WaitForFinalization, "Wait for finalization should not be nil")
	require.False(t, *cfg.WaitForFinalization, "Wait for finalization should be false")
	require.Equal(t, 2, len(cfg.EthereumChainConfig.AddressesToFund), "Should have 2 addresses to fund")
	require.Equal(t, 12, cfg.EthereumChainConfig.SecondsPerSlot, "Seconds per slot should be 12")
	require.Equal(t, 2, cfg.EthereumChainConfig.SlotsPerEpoch, "Slots per epoch should be 2")
	require.Equal(t, 20, cfg.EthereumChainConfig.GenesisDelay, "Genesis delay should be 20")
	require.Equal(t, 8, cfg.EthereumChainConfig.ValidatorCount, "Validator count should be 8")
	require.Equal(t, 1234, cfg.EthereumChainConfig.ChainID, "Chain ID should be 1234")
}

func TestEthEnvCustomDockerImagesFromToml(t *testing.T) {
	t.Parallel()
	toml := `
	[EthereumNetwork]
	consensus_type="pos"
	consensus_layer="prysm"
	execution_layer="geth"
	wait_for_finalization=false

	[EthereumNetwork.EthereumChainConfig]
	seconds_per_slot=12
	slots_per_epoch=2
	genesis_delay=20
	validator_count=8
	chain_id=1234
	addresses_to_fund=["0x742d35Cc6634C0532925a3b844Bc454e4438f44e", "0x742d35Cc6634C0532925a3b844Bc454e4438f44f"]

	[EthereumNetwork.CustomDockerImages]
	execution_layer="i-dont-exist:tag-me"
	`

	tomlCfg, err := readEthereumNetworkConfig(toml)
	require.NoError(t, err, "Couldn't read config")

	tomlCfg.EthereumChainConfig.GenerateGenesisTimestamp()

	err = tomlCfg.Validate()
	require.NotNil(t, err, "Should have failed to validate TOML config")
	require.Error(t, err, "Couldn't validate TOML config")
	require.Equal(t, fmt.Sprintf("failed to parse docker image and extract version: %s", "i-dont-exist:tag-me"), err.Error(), "Error message is not correct")
}

func TestEthNoLogLevelDefaultsToInfo(t *testing.T) {
	t.Parallel()
	tomlStr := `
	[EthereumNetwork]
	consensus_type="pos"
	consensus_layer="prysm"
	execution_layer="geth"
	wait_for_finalization=false

	[EthereumNetwork.EthereumChainConfig]
	seconds_per_slot=12
	slots_per_epoch=2
	genesis_delay=20
	validator_count=8
	chain_id=1234
	addresses_to_fund=["0x742d35Cc6634C0532925a3b844Bc454e4438f44e", "0x742d35Cc6634C0532925a3b844Bc454e4438f44f"]
	`

	tomlCfg, err := readEthereumNetworkConfig(tomlStr)
	require.NoError(t, err, "Couldn't read config")

	tomlCfg.EthereumChainConfig.GenerateGenesisTimestamp()

	err = tomlCfg.Validate()
	require.NoError(t, err, "Failed to validate correct TOML config")
	require.Equal(t, config.DefaultNodeLogLevel, *tomlCfg.NodeLogLevel, "Incorrect default log level")
}

func TestEthLogLevelTrace(t *testing.T) {
	t.Parallel()
	tomlStr := `
	[EthereumNetwork]
	consensus_type="pos"
	consensus_layer="prysm"
	execution_layer="geth"
	wait_for_finalization=false
	node_log_level="trace"

	[EthereumNetwork.EthereumChainConfig]
	seconds_per_slot=12
	slots_per_epoch=2
	genesis_delay=20
	validator_count=8
	chain_id=1234
	addresses_to_fund=["0x742d35Cc6634C0532925a3b844Bc454e4438f44e", "0x742d35Cc6634C0532925a3b844Bc454e4438f44f"]
	`

	tomlCfg, err := readEthereumNetworkConfig(tomlStr)
	require.NoError(t, err, "Couldn't read config")

	tomlCfg.EthereumChainConfig.GenerateGenesisTimestamp()

	err = tomlCfg.Validate()
	require.NoError(t, err, "Failed to validate correct TOML config")
	require.Equal(t, "trace", *tomlCfg.NodeLogLevel, "Incorrect default log level")
}

func TestEthInvalidLogLevel(t *testing.T) {
	t.Parallel()
	tomlStr := `
	[EthereumNetwork]
	consensus_type="pos"
	consensus_layer="prysm"
	execution_layer="geth"
	wait_for_finalization=false
	node_log_level="random"

	[EthereumNetwork.EthereumChainConfig]
	seconds_per_slot=12
	slots_per_epoch=2
	genesis_delay=20
	validator_count=8
	chain_id=1234
	addresses_to_fund=["0x742d35Cc6634C0532925a3b844Bc454e4438f44e", "0x742d35Cc6634C0532925a3b844Bc454e4438f44f"]
	`

	tomlCfg, err := readEthereumNetworkConfig(tomlStr)
	require.NoError(t, err, "Couldn't read config")

	tomlCfg.EthereumChainConfig.GenerateGenesisTimestamp()

	err = tomlCfg.Validate()
	require.Error(t, err, "Successfully validated incorrect TOML config")
}

type ethereumNetworkWrapper struct {
	EthereumNetwork *config.EthereumNetworkConfig `toml:"EthereumNetwork"`
}

func readEthereumNetworkConfig(configDecoded string) (config.EthereumNetworkConfig, error) {
	var net ethereumNetworkWrapper
	err := toml.Unmarshal([]byte(configDecoded), &net)
	if err != nil {
		return config.EthereumNetworkConfig{}, fmt.Errorf("error unmarshalling ethereum network config: %w", err)
	}

	return *net.EthereumNetwork, nil
}

func TestEthEnvCustomDockerNetworks(t *testing.T) {
	t.Parallel()
	networks := []string{"test-network"}

	builder := NewEthereumNetworkBuilder()
	cfg, err := builder.
		WithEthereumVersion(config_types.EthereumVersion_Eth2).
		WithConsensusLayer(config.ConsensusLayer_Prysm).
		WithExecutionLayer(config_types.ExecutionLayer_Geth).
		WithDockerNetworks(networks).
		Build()
	require.NoError(t, err, "Builder validation failed")
	require.Equal(t, networks, cfg.DockerNetworkNames, "Incorrect docker networks in config")
}

func TestEthEnvValidHardForks(t *testing.T) {
	t.Parallel()
	chainConfig := config.MustGetDefaultChainConfig()
	chainConfig.HardForkEpochs = map[string]int{"Deneb": 500}

	builder := NewEthereumNetworkBuilder()
	_, err := builder.
		WithCustomDockerImages(map[config.ContainerType]string{config.ContainerType_ExecutionLayer: "ethereum/client-go:v1.13.11"}).
		WithExecutionLayer(config_types.ExecutionLayer_Geth).
		WithEthereumChainConfig(chainConfig).
		Build()
	require.NoError(t, err, "Builder validation failed")
}

func TestEthEnvInvalidHardForks(t *testing.T) {
	t.Parallel()
	builder := NewEthereumNetworkBuilder()
	_, err := builder.
		WithEthereumVersion(config_types.EthereumVersion_Eth2).
		WithConsensusLayer(config.ConsensusLayer_Prysm).
		WithExecutionLayer(config_types.ExecutionLayer_Geth).
		WithEthereumChainConfig(config.EthereumChainConfig{
			HardForkEpochs: map[string]int{"Electra": 0},
		}).
		Build()
	require.NoError(t, err, "Builder validation failed")
}

func TestEthEnvAutoEthereumVersionEth1Minor(t *testing.T) {
	t.Parallel()
	builder := NewEthereumNetworkBuilder()
	cfg, err := builder.
		WithExecutionLayer(config_types.ExecutionLayer_Besu).
		WithCustomDockerImages(map[config.ContainerType]string{
			config.ContainerType_ExecutionLayer: "hyperledger/besu:20.1"}).
		Build()
	require.NoError(t, err, "Builder validation failed")
	require.Equal(t, config_types.EthereumVersion_Eth1, *cfg.EthereumVersion, "Ethereum Version should be Eth1")
	require.Nil(t, cfg.ConsensusLayer, "Consensus layer should be nil")
}

func TestEthEnvAutoEthereumVersionEth2Minor(t *testing.T) {
	builder := NewEthereumNetworkBuilder()
	cfg, err := builder.
		WithExecutionLayer(config_types.ExecutionLayer_Nethermind).
		WithCustomDockerImages(map[config.ContainerType]string{
			config.ContainerType_ExecutionLayer: "nethermind/nethermind:1.17.2"}).
		Build()
	require.NoError(t, err, "Builder validation failed")
	require.Equal(t, config_types.EthereumVersion_Eth2, *cfg.EthereumVersion, "Ethereum Version should be Eth2")
	require.Equal(t, config.ConsensusLayer_Prysm, *cfg.ConsensusLayer, "Consensus layer should be Prysm")
}

func TestEthEnvAutoEthereumVersionReleaseCandidate(t *testing.T) {
	t.Parallel()
	builder := NewEthereumNetworkBuilder()
	cfg, err := builder.
		WithExecutionLayer(config_types.ExecutionLayer_Nethermind).
		WithCustomDockerImages(map[config.ContainerType]string{
			config.ContainerType_ExecutionLayer: "nethermind/nethermind:1.17.0-RC2"}).
		Build()
	require.NoError(t, err, "Builder validation failed")
	require.Equal(t, config_types.EthereumVersion_Eth2, *cfg.EthereumVersion, "Ethereum Version should be Eth2")
	require.Equal(t, config.ConsensusLayer_Prysm, *cfg.ConsensusLayer, "Consensus layer should be Prysm")
}

func TestEthEnvAutoEthereumVersionWithLettersInVersion(t *testing.T) {
	t.Parallel()
	builder := NewEthereumNetworkBuilder()
	cfg, err := builder.
		WithConsensusLayer(config.ConsensusLayer_Prysm).
		WithExecutionLayer(config_types.ExecutionLayer_Geth).
		WithCustomDockerImages(map[config.ContainerType]string{
			config.ContainerType_ExecutionLayer: "ethereum/client-go:v1.13.10"}).
		Build()
	require.NoError(t, err, "Builder validation failed")
	require.Equal(t, config_types.EthereumVersion_Eth2, *cfg.EthereumVersion, "Ethereum Version should be Eth2")
	require.Equal(t, config_types.ExecutionLayer_Geth, *cfg.ExecutionLayer, "Execution layer should be Geth")
	require.Equal(t, config.ConsensusLayer_Prysm, *cfg.ConsensusLayer, "Consensus layer should be Prysm")
}

func TestEthEnvAutoEthereumVersionOnlyMajor(t *testing.T) {
	t.Parallel()
	builder := NewEthereumNetworkBuilder()
	cfg, err := builder.
		WithCustomDockerImages(map[config.ContainerType]string{
			config.ContainerType_ExecutionLayer: "hyperledger/besu:v24.1"}).
		Build()
	require.NoError(t, err, "Builder validation failed")
	require.Equal(t, config_types.EthereumVersion_Eth2, *cfg.EthereumVersion, "Ethereum Version should be Eth2")
	require.Equal(t, config_types.ExecutionLayer_Besu, *cfg.ExecutionLayer, "Execution layer should be Besu")
	require.Equal(t, config.ConsensusLayer_Prysm, *cfg.ConsensusLayer, "Consensus layer should be Prysm")
}

func TestEthEnvLatestVersionFromGithub(t *testing.T) {
	t.Parallel()
	builder := NewEthereumNetworkBuilder()
	cfg, err := builder.
		WithCustomDockerImages(map[config.ContainerType]string{
			config.ContainerType_ExecutionLayer: fmt.Sprintf("hyperledger/besu:%s", AUTOMATIC_STABLE_LATEST_TAG)}).
		Build()
	require.NoError(t, err, "Builder validation failed")
	require.Equal(t, config_types.EthereumVersion_Eth2, *cfg.EthereumVersion, "Ethereum Version should be Eth2")
	require.Equal(t, config_types.ExecutionLayer_Besu, *cfg.ExecutionLayer, "Execution layer should be Besu")
	require.Equal(t, config.ConsensusLayer_Prysm, *cfg.ConsensusLayer, "Consensus layer should be Prysm")
	require.NotContains(t, cfg.CustomDockerImages[config.ContainerType_ExecutionLayer], AUTOMATIC_STABLE_LATEST_TAG, "Automatic tag should be replaced")
}

func TestEthEnvMischmachedExecutionClient(t *testing.T) {
	t.Parallel()
	builder := NewEthereumNetworkBuilder()
	_, err := builder.
		WithExecutionLayer(config_types.ExecutionLayer_Nethermind).
		WithCustomDockerImages(map[config.ContainerType]string{
			config.ContainerType_ExecutionLayer: "hyperledger/besu:v1.22.0"}).
		Build()
	require.Error(t, err, "Builder validation succeeded")
	require.Equal(t, fmt.Sprintf(MsgMismatchedExecutionClient, config_types.ExecutionLayer_Besu, config_types.ExecutionLayer_Nethermind), err.Error(), "Error message is not correct")
}
