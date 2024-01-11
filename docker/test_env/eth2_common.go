package test_env

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pelletier/go-toml/v2"
	"github.com/rs/zerolog"
	tc "github.com/testcontainers/testcontainers-go"

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
	GENERATED_DATA_DIR_INSIDE_CONTAINER             = "/data/custom_config_data"
	JWT_SECRET_FILE_LOCATION_INSIDE_CONTAINER       = fmt.Sprintf("%s/jwtsecret", GENERATED_DATA_DIR_INSIDE_CONTAINER) // #nosec G101
	VALIDATOR_BIP39_MNEMONIC                        = "giant issue aisle success illegal bike spike question tent bar rely arctic volcano long crawl hungry vocal artwork sniff fantasy very lucky have athlete"
)

type EthereumChainConfig struct {
	SecondsPerSlot   int      `json:"seconds_per_slot" toml:"seconds_per_slot"`
	SlotsPerEpoch    int      `json:"slots_per_epoch" toml:"slots_per_epoch"`
	GenesisDelay     int      `json:"genesis_delay" toml:"genesis_delay"`
	ValidatorCount   int      `json:"validator_count" toml:"validator_count"`
	ChainID          int      `json:"chain_id" toml:"chain_id"`
	genesisTimestamp int      // this is not serialized
	AddressesToFund  []string `json:"addresses_to_fund" toml:"addresses_to_fund"`
}

//go:embed tomls/default_ethereum_env.toml
var defaultEthereumChainConfig []byte

func (c *EthereumChainConfig) Default() error {
	wrapper := struct {
		EthereumNetwork *EthereumNetwork `toml:"PrivateEthereumNetwork"`
	}{}
	if err := toml.Unmarshal(defaultEthereumChainConfig, &wrapper); err != nil {
		return fmt.Errorf("error unmarshaling ethereum network config: %w", err)
	}

	if wrapper.EthereumNetwork == nil {
		return errors.New("[EthereumNetwork] was not present in default TOML file")
	}

	*c = *wrapper.EthereumNetwork.EthereumChainConfig

	if c.genesisTimestamp == 0 {
		c.GenerateGenesisTimestamp()
	}

	return nil
}

func GetDefaultChainConfig() EthereumChainConfig {
	config := EthereumChainConfig{}
	if err := config.Default(); err != nil {
		panic(err)
	}
	return config
}

func (c *EthereumChainConfig) Validate(logger zerolog.Logger) error {
	if c.ValidatorCount < 4 {
		return fmt.Errorf("validator count must be >= 4")
	}
	if c.SecondsPerSlot < 3 {
		return fmt.Errorf("seconds per slot must be >= 3")
	}
	if c.SlotsPerEpoch < 2 {
		return fmt.Errorf("slots per epoch must be >= 1")
	}
	if c.GenesisDelay < 10 {
		return fmt.Errorf("genesis delay must be >= 10")
	}
	if c.ChainID < 1 {
		return fmt.Errorf("chain id must be >= 0")
	}
	if c.genesisTimestamp == 0 {
		return fmt.Errorf("genesis timestamp must be generated by calling GenerateGenesisTimestamp()")
	}

	addressSet := make(map[string]struct{})
	deduplicated := make([]string, 0)

	for _, addr := range c.AddressesToFund {
		if !common.IsHexAddress(addr) {
			return fmt.Errorf("address %s is not a valid hex address", addr)
		}

		if _, exists := addressSet[addr]; exists {
			logger.Warn().Str("address", addr).Msg("duplicate address in addresses to fund, this should not happen, removing it so that genesis generation doesn't crash")
			continue
		}

		addressSet[addr] = struct{}{}
		deduplicated = append(deduplicated, addr)
	}

	c.AddressesToFund = deduplicated

	return nil
}

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

func (c *EthereumChainConfig) fillInMissingValuesWithDefault() {
	defaultConfig := GetDefaultChainConfig()
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
}

func (c *EthereumChainConfig) GetValidatorBasedGenesisDelay() int {
	return c.ValidatorCount * 5
}

func (c *EthereumChainConfig) GenerateGenesisTimestamp() {
	c.genesisTimestamp = int(time.Now().Unix()) + c.GetValidatorBasedGenesisDelay()
}

func (c *EthereumChainConfig) GetDefaultWaitDuration() time.Duration {
	return time.Duration((c.GenesisDelay+c.GetValidatorBasedGenesisDelay())*2) * time.Second
}

func (c *EthereumChainConfig) GetDefaultFinalizationWaitDuration() time.Duration {
	return time.Duration(5 * time.Minute)
}

type ExecutionClient interface {
	GetContainerName() string
	StartContainer() (blockchain.EVMNetwork, error)
	GetContainer() *tc.Container
	GetContainerType() ContainerType
	GetInternalExecutionURL() string
	GetExternalExecutionURL() string
	GetInternalHttpUrl() string
	GetInternalWsUrl() string
	GetExternalHttpUrl() string
	GetExternalWsUrl() string
	WaitUntilChainIsReady(ctx context.Context, waitTime time.Duration) error
	WithTestInstance(t *testing.T) ExecutionClient
}

type HasDockerImage[T any] interface {
	WithImage(imageWithTag string) T
	GetImage() string
}
