package config

import (
	_ "embed"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/pelletier/go-toml/v2"
	"github.com/rs/zerolog"
	tc "github.com/testcontainers/testcontainers-go"

	config_types "github.com/smartcontractkit/chainlink-testing-framework/lib/config/types"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/docker/ethereum"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/slice"
)

var (
	ErrMissingEthereumVersion = errors.New("ethereum version is required")
	ErrMissingExecutionLayer  = errors.New("execution layer is required")
	Eth1NotSupportedByRethMsg = "eth1 is not supported by Reth, please use eth2"
	DefaultNodeLogLevel       = "info"
)

type EthereumNetworkConfig struct {
	ConsensusType        *config_types.EthereumVersion `toml:"consensus_type"`
	EthereumVersion      *config_types.EthereumVersion `toml:"ethereum_version"`
	ConsensusLayer       *ConsensusLayer               `toml:"consensus_layer"`
	ExecutionLayer       *config_types.ExecutionLayer  `toml:"execution_layer"`
	DockerNetworkNames   []string                      `toml:"docker_network_names"`
	Containers           EthereumNetworkContainers     `toml:"containers"`
	WaitForFinalization  *bool                         `toml:"wait_for_finalization"`
	GeneratedDataHostDir *string                       `toml:"generated_data_host_dir"`
	ValKeysDir           *string                       `toml:"val_keys_dir"`
	EthereumChainConfig  *EthereumChainConfig          `toml:"EthereumChainConfig"`
	CustomDockerImages   map[ContainerType]string      `toml:"CustomDockerImages"`
	NodeLogLevel         *string                       `toml:"node_log_level,omitempty"`
}

func (en *EthereumNetworkConfig) Validate() error {
	l := logging.GetTestLogger(nil)

	// logically it doesn't belong here, but placing it here guarantees it will always run without changing API
	if en.EthereumVersion != nil && en.ConsensusType != nil {
		l.Warn().Msg("Both EthereumVersion and ConsensusType are set. ConsensusType as a _deprecated_ field will be ignored")
	}

	if en.EthereumVersion == nil && en.ConsensusType != nil {
		l.Debug().Msg("Using _deprecated_ ConsensusType as EthereumVersion")
		tempEthVersion := en.ConsensusType
		switch *tempEthVersion {
		//nolint:staticcheck //ignore SA1019
		case config_types.EthereumVersion_Eth1, config_types.EthereumVersion_Eth1_Legacy:
			*tempEthVersion = config_types.EthereumVersion_Eth1
		//nolint:staticcheck //ignore SA1019
		case config_types.EthereumVersion_Eth2, config_types.EthereumVersion_Eth2_Legacy:
			*tempEthVersion = config_types.EthereumVersion_Eth2
		default:
			return fmt.Errorf("unknown ethereum version (consensus type): %s", *en.ConsensusType)
		}

		en.EthereumVersion = tempEthVersion
	}

	if (en.EthereumVersion == nil || *en.EthereumVersion == "") && len(en.CustomDockerImages) == 0 {
		return ErrMissingEthereumVersion
	}

	if (en.ExecutionLayer == nil || *en.ExecutionLayer == "") && len(en.CustomDockerImages) == 0 {
		return ErrMissingExecutionLayer
	}

	//nolint:staticcheck //ignore SA1019
	if (en.EthereumVersion != nil && (*en.EthereumVersion == config_types.EthereumVersion_Eth2_Legacy || *en.EthereumVersion == config_types.EthereumVersion_Eth2)) && (en.ConsensusLayer == nil || *en.ConsensusLayer == "") {
		l.Warn().Msg("Consensus layer is not set, but is required for PoS. Defaulting to Prysm")
		en.ConsensusLayer = &ConsensusLayer_Prysm
	}

	//nolint:staticcheck //ignore SA1019
	if (en.EthereumVersion != nil && (*en.EthereumVersion == config_types.EthereumVersion_Eth1_Legacy || *en.EthereumVersion == config_types.EthereumVersion_Eth1)) && (en.ConsensusLayer != nil && *en.ConsensusLayer != "") {
		l.Warn().Msg("Consensus layer is set, but is not allowed for PoW. Ignoring")
		en.ConsensusLayer = nil
	}

	if en.NodeLogLevel == nil {
		en.NodeLogLevel = &DefaultNodeLogLevel
	}

	if *en.EthereumVersion == config_types.EthereumVersion_Eth1 && *en.ExecutionLayer == config_types.ExecutionLayer_Reth {
		msg := `%s

If you are using builder to create the network, please change the EthereumVersion to EthereumVersion_Eth2 by calling this method:
WithEthereumVersion(config.EthereumVersion_Eth2).

If you are using a TOML file, please change the EthereumVersion to "eth2" in the TOML file:
[PrivateEthereumNetwork]
ethereum_version="eth2"
`
		return fmt.Errorf(msg, Eth1NotSupportedByRethMsg)
	}

	switch strings.ToLower(*en.NodeLogLevel) {
	case "trace", "debug", "info", "warn", "error":
	default:
		return fmt.Errorf("invalid node log level: %s", *en.NodeLogLevel)
	}

	if en.EthereumChainConfig == nil {
		return errors.New("ethereum chain config is required")
	}

	return en.EthereumChainConfig.Validate(l, en.EthereumVersion, en.ExecutionLayer, en.CustomDockerImages)
}

func (en *EthereumNetworkConfig) ApplyOverrides(from *EthereumNetworkConfig) error {
	if from == nil {
		return nil
	}
	if from.ConsensusLayer != nil {
		en.ConsensusLayer = from.ConsensusLayer
	}
	if from.ExecutionLayer != nil {
		en.ExecutionLayer = from.ExecutionLayer
	}
	if from.EthereumVersion != nil {
		en.EthereumVersion = from.EthereumVersion
	}
	if from.WaitForFinalization != nil {
		en.WaitForFinalization = from.WaitForFinalization
	}

	if from.EthereumChainConfig != nil {
		if en.EthereumChainConfig == nil {
			en.EthereumChainConfig = from.EthereumChainConfig
		} else {
			err := en.EthereumChainConfig.ApplyOverrides(from.EthereumChainConfig)
			if err != nil {
				return fmt.Errorf("error applying overrides from network config file to config: %w", err)
			}
		}
	}

	return nil
}

func (en *EthereumNetworkConfig) Describe() string {
	cL := "prysm"
	if en.ConsensusLayer == nil {
		cL = "(none)"
	}
	return fmt.Sprintf("ethereum version: %s, execution layer: %s, consensus layer: %s", *en.EthereumVersion, *en.ExecutionLayer, cL)
}

type EthereumNetworkContainer struct {
	ContainerName string        `toml:"container_name"`
	ContainerType ContainerType `toml:"container_type"`
	Container     *tc.Container `toml:"-"`
}

// Deprecated: use EthereumVersion instead
type ConsensusType string

const (
	// Deprecated: use EthereumVersion_Eth2 instead
	ConsensusType_PoS ConsensusType = "pos"
	// Deprecated: use EthereumVersion_Eth1 instead
	ConsensusType_PoW ConsensusType = "pow"
)

type ConsensusLayer string

var ConsensusLayer_Prysm ConsensusLayer = "prysm"

type EthereumNetworkContainers []EthereumNetworkContainer

type ContainerType string

const (
	ContainerType_ExecutionLayer     ContainerType = "execution_layer"
	ContainerType_ConsensusLayer     ContainerType = "consensus_layer"
	ContainerType_ConsensusValidator ContainerType = "consensus_validator"
	ContainerType_GenesisGenerator   ContainerType = "genesis_generator"
	ContainerType_ValKeysGenerator   ContainerType = "val_keys_generator"
)

type EthereumChainConfig struct {
	SecondsPerSlot   int            `json:"seconds_per_slot" toml:"seconds_per_slot"`
	SlotsPerEpoch    int            `json:"slots_per_epoch" toml:"slots_per_epoch"`
	GenesisDelay     int            `json:"genesis_delay" toml:"genesis_delay"`
	ValidatorCount   int            `json:"validator_count" toml:"validator_count"`
	ChainID          int            `json:"chain_id" toml:"chain_id"`
	GenesisTimestamp int            // this is not serialized
	AddressesToFund  []string       `json:"addresses_to_fund" toml:"addresses_to_fund"`
	HardForkEpochs   map[string]int `json:"HardForkEpochs" toml:"HardForkEpochs"`
}

//go:embed tomls/default_ethereum_env.toml
var defaultEthereumChainConfig []byte

// Default sets the EthereumChainConfig to the default values
func (c *EthereumChainConfig) Default() error {
	wrapper := struct {
		EthereumNetwork *EthereumNetworkConfig `toml:"PrivateEthereumNetwork"`
	}{}
	if err := toml.Unmarshal(defaultEthereumChainConfig, &wrapper); err != nil {
		return fmt.Errorf("error unmarshalling ethereum network config: %w", err)
	}

	if wrapper.EthereumNetwork == nil {
		return errors.New("[EthereumNetwork] was not present in default TOML file")
	}

	*c = *wrapper.EthereumNetwork.EthereumChainConfig

	if c.GenesisTimestamp == 0 {
		c.GenerateGenesisTimestamp()
	}

	return nil
}

// MustGetDefaultChainConfig returns the default EthereumChainConfig or panics if it can't be loaded
func MustGetDefaultChainConfig() EthereumChainConfig {
	config := EthereumChainConfig{}
	if err := config.Default(); err != nil {
		panic(err)
	}
	return config
}

// Validate validates the EthereumChainConfig
func (c *EthereumChainConfig) Validate(l zerolog.Logger, ethereumVersion *config_types.EthereumVersion, executionLayer *config_types.ExecutionLayer, customDockerImages map[ContainerType]string) error {
	if c.ChainID < 1 {
		return fmt.Errorf("chain id must be >= 0")
	}

	// don't like it 100% but in cases where we load private ethereum network config from TOML it might be incomplete
	// until we pass it to ethereum network builder that will fill in defaults
	//nolint:staticcheck //ignore SA1019
	if ethereumVersion == nil || (*ethereumVersion == config_types.EthereumVersion_Eth1_Legacy || *ethereumVersion == config_types.EthereumVersion_Eth1) {
		return nil
	}

	if c.ValidatorCount < 4 {
		return fmt.Errorf("validator count must be >= 4")
	}
	if c.SecondsPerSlot < 3 {
		return fmt.Errorf("seconds per slot must be >= 3")
	}
	if c.SlotsPerEpoch < 2 {
		return fmt.Errorf("slots per epoch must be >= 2")
	}
	if c.GenesisDelay < 10 {
		return fmt.Errorf("genesis delay must be >= 10")
	}
	if c.GenesisTimestamp == 0 {
		return fmt.Errorf("genesis timestamp must be generated by calling GenerateGenesisTimestamp()")
	}

	if err := c.ValidateHardForks(l, ethereumVersion, executionLayer, customDockerImages); err != nil {
		return err
	}

	var err error
	var hadDuplicates bool
	// we need to deduplicate addresses to fund, because if present they will crash the genesis
	c.AddressesToFund, hadDuplicates, err = slice.ValidateAndDeduplicateAddresses(c.AddressesToFund)
	if err != nil {
		return err
	}
	if hadDuplicates {
		l.Warn().Msg("Duplicate addresses found in addresses_to_fund. Removed them. You might want to review your configuration.")
	}

	return nil
}

// ValidateHardForks validates hard forks based either on custom or default docker images for eth2 execution layer
func (c *EthereumChainConfig) ValidateHardForks(l zerolog.Logger, ethereumVersion *config_types.EthereumVersion, executionLayer *config_types.ExecutionLayer, customDockerImages map[ContainerType]string) error {
	//nolint:staticcheck //ignore SA1019
	if ethereumVersion == nil || (*ethereumVersion == config_types.EthereumVersion_Eth1_Legacy || *ethereumVersion == config_types.EthereumVersion_Eth1) {
		return nil
	}

	customImage := customDockerImages[ContainerType_ExecutionLayer]
	var baseEthereumFork ethereum.Fork
	var err error
	if customImage == "" {
		if executionLayer == nil {
			return ErrMissingExecutionLayer
		}
		var dockerImage string
		switch *executionLayer {
		case config_types.ExecutionLayer_Geth:
			dockerImage = ethereum.DefaultGethEth2Image
		case config_types.ExecutionLayer_Nethermind:
			dockerImage = ethereum.DefaultNethermindEth2Image
		case config_types.ExecutionLayer_Erigon:
			dockerImage = ethereum.DefaultErigonEth2Image
		case config_types.ExecutionLayer_Besu:
			dockerImage = ethereum.DefaultBesuEth2Image
		case config_types.ExecutionLayer_Reth:
			dockerImage = ethereum.DefaultRethEth2Image
		}
		baseEthereumFork, err = ethereum.LastSupportedForkForEthereumClient(dockerImage)
	} else {
		baseEthereumFork, err = ethereum.LastSupportedForkForEthereumClient(customImage)
	}

	if err != nil {
		return err
	}

	validFutureForks, err := baseEthereumFork.ValidFutureForks()
	if err != nil {
		return err
	}

	validForks := make(map[string]int)

	// latest Prysm Beacon Chain doesn't support any fork (Electra is coming in 2025)
	// but older versions do support Deneb
	for fork, epoch := range c.HardForkEpochs {
		isValid := false
		for _, validFork := range validFutureForks {
			if strings.EqualFold(fork, string(validFork)) {
				isValid = true
				validForks[fork] = epoch
				break
			}
		}

		if !isValid {
			l.Debug().Msgf("Fork %s is not supported. Removed it from configuration", fork)
		}
	}

	// at the same time for Shanghai-based forks we need to add Deneb to the list if it's not there, so that genesis is valid
	if _, ok := c.HardForkEpochs[string(ethereum.EthereumFork_Deneb)]; !ok && baseEthereumFork == ethereum.EthereumFork_Shanghai {
		l.Debug().Msg("Adding Deneb to fork setup, because it's required, but was missing from the configuration. It's scheduled for epoch 1000")
		validForks[string(ethereum.EthereumFork_Deneb)] = 1000
	}

	c.HardForkEpochs = validForks

	return nil
}

// ApplyOverrides applies overrides from another EthereumChainConfig
func (c *EthereumChainConfig) ApplyOverrides(from *EthereumChainConfig) error {
	if from == nil {
		return nil
	}
	if from.ValidatorCount != 0 {
		c.ValidatorCount = from.ValidatorCount
	}
	if from.SecondsPerSlot != 0 {
		c.SecondsPerSlot = from.SecondsPerSlot
	}
	if from.SlotsPerEpoch != 0 {
		c.SlotsPerEpoch = from.SlotsPerEpoch
	}
	if from.GenesisDelay != 0 {
		c.GenesisDelay = from.GenesisDelay
	}
	if from.ChainID != 0 {
		c.ChainID = from.ChainID
	}
	if len(from.AddressesToFund) != 0 {
		c.AddressesToFund = append([]string{}, from.AddressesToFund...)
	}
	return nil
}

// FillInMissingValuesWithDefault fills in missing/zero values with default values
func (c *EthereumChainConfig) FillInMissingValuesWithDefault() {
	defaultConfig := MustGetDefaultChainConfig()
	if c.ValidatorCount == 0 {
		c.ValidatorCount = defaultConfig.ValidatorCount
	}
	if c.SecondsPerSlot == 0 {
		c.SecondsPerSlot = defaultConfig.SecondsPerSlot
	}
	if c.SlotsPerEpoch == 0 {
		c.SlotsPerEpoch = defaultConfig.SlotsPerEpoch
	}
	if c.GenesisDelay == 0 {
		c.GenesisDelay = defaultConfig.GenesisDelay
	}
	if c.ChainID == 0 {
		c.ChainID = defaultConfig.ChainID
	}
	if len(c.AddressesToFund) == 0 {
		c.AddressesToFund = append([]string{}, defaultConfig.AddressesToFund...)
	} else {
		c.AddressesToFund = append(append([]string{}, c.AddressesToFund...), defaultConfig.AddressesToFund...)
	}

	if len(c.HardForkEpochs) == 0 {
		c.HardForkEpochs = defaultConfig.HardForkEpochs
	}
}

// ValidatorBasedGenesisDelay returns the delay in seconds based on the number of validators
func (c *EthereumChainConfig) ValidatorBasedGenesisDelay() int {
	return c.ValidatorCount * 5
}

func (c *EthereumChainConfig) GenerateGenesisTimestamp() {
	c.GenesisTimestamp = int(time.Now().Unix()) + c.ValidatorBasedGenesisDelay()
}

// DefaultWaitDuration returns the default wait duration for the network based on the genesis delay and the number of validators
func (c *EthereumChainConfig) DefaultWaitDuration() time.Duration {
	return time.Duration((c.GenesisDelay+c.ValidatorBasedGenesisDelay())*2) * time.Second
}

// DefaultFinalizationWaitDuration returns the default wait duration for finalization
func (c *EthereumChainConfig) DefaultFinalizationWaitDuration() time.Duration {
	return 5 * time.Minute
}
