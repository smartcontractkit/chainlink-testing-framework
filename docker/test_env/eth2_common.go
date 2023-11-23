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
	addressesToFund  []string
}

var DefaultBeaconChainConfig = func() EthereumChainConfig {
	config := EthereumChainConfig{
		SecondsPerSlot:  12,
		SlotsPerEpoch:   6,
		GenesisDelay:    15,
		ValidatorCount:  8,
		ChainID:         1337,
		addressesToFund: []string{"0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"},
	}
	config.GenerateGenesisTimestamp()
	return config
}()

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
