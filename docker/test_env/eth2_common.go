package test_env

import (
	"fmt"
	"time"

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
	SecondsPerSlot   int
	SlotsPerEpoch    int
	GenesisDelay     int
	ValidatorCount   int
	ChainID          int
	genesisTimestamp int
	AddressesToFund  []string
}

var DefaultChainConfig = func() EthereumChainConfig {
	config := EthereumChainConfig{
		SecondsPerSlot:  12,
		SlotsPerEpoch:   6,
		GenesisDelay:    15,
		ValidatorCount:  8,
		ChainID:         1337,
		AddressesToFund: []string{"0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"},
	}
	config.GenerateGenesisTimestamp()
	return config
}()

func (c *EthereumChainConfig) Validate() error {
	if c.ValidatorCount < 4 {
		return fmt.Errorf("validator count must be >= 4")
	}
	if c.SecondsPerSlot < 4 {
		return fmt.Errorf("seconds per slot must be >= 4")
	}
	if c.SlotsPerEpoch < 2 {
		return fmt.Errorf("slots per epoch must be >= 2")
	}
	if c.GenesisDelay < 10 {
		return fmt.Errorf("genesis delay must be >= 10")
	}
	if c.ChainID < 1 {
		return fmt.Errorf("chain id must be >= 0")
	}
	return nil
}

func (c *EthereumChainConfig) fillInMissingValuesWithDefault() {
	if c.ValidatorCount == 0 {
		c.ValidatorCount = DefaultChainConfig.ValidatorCount
	}
	if c.SecondsPerSlot == 0 {
		c.SecondsPerSlot = DefaultChainConfig.SecondsPerSlot
	}
	if c.SlotsPerEpoch == 0 {
		c.SlotsPerEpoch = DefaultChainConfig.SlotsPerEpoch
	}
	if c.GenesisDelay == 0 {
		c.GenesisDelay = DefaultChainConfig.GenesisDelay
	}
	if c.ChainID == 0 {
		c.ChainID = DefaultChainConfig.ChainID
	}
	if len(c.AddressesToFund) == 0 {
		c.AddressesToFund = DefaultChainConfig.AddressesToFund
	} else {
		c.AddressesToFund = append(c.AddressesToFund, DefaultChainConfig.AddressesToFund...)
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
	WaitUntilChainIsReady(waitTime time.Duration) error
}
