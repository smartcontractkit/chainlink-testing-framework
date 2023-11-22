package test_env

import (
	"time"

	tc "github.com/testcontainers/testcontainers-go"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
)

const (
	ETH2_EXECUTION_PORT                = "8551"
	CONTAINER_ETH2_CONSENSUS_DIRECTORY = "/consensus"
	CONTAINER_ETH2_EXECUTION_DIRECTORY = "/execution"
	beaconConfigFile                   = "/consensus/config.yml"
	eth2GenesisFile                    = "/consensus/genesis.ssz"
	eth1GenesisFile                    = "/execution/genesis.json"
	jwtSecretFileLocation              = "/execution/jwtsecret" // #nosec G101
	VALIDATOR_BIPC39_MNEMONIC          = "giant issue aisle success illegal bike spike question tent bar rely arctic volcano long crawl hungry vocal artwork sniff fantasy very lucky have athlete"
)

type BeaconChainConfig struct {
	SecondsPerSlot   int
	SlotsPerEpoch    int
	GenesisDelay     int
	ValidatorCount   int
	ChainID          int
	GenesisTimestamp int
}

var DefaultBeaconChainConfig = func() BeaconChainConfig {
	config := BeaconChainConfig{
		SecondsPerSlot: 12,
		SlotsPerEpoch:  6,
		GenesisDelay:   30,
		ValidatorCount: 8,
		ChainID:        1337,
	}
	config.GenerateGenesisTimestamp()
	return config
}()

func (b *BeaconChainConfig) GetValidatorBasedGenesisDelay() int {
	return b.ValidatorCount * 5
}

func (b *BeaconChainConfig) GenerateGenesisTimestamp() {
	b.GenesisTimestamp = int(time.Now().Unix()) + b.GetValidatorBasedGenesisDelay()
}

type ExecutionClient interface {
	GetContainerName() string
	StartContainer() (blockchain.EVMNetwork, error)
	GetContainer() *tc.Container
	GetInternalExecutionURL() string
	GetExternalExecutionURL() string
	GetInternalHttpUrl() string
	GetInternalWsUrl() string
	GetExternalHttpUrl() string
	GetExternalWsUrl() string
	WaitUntilChainIsReady(waitTime time.Duration) error
}

// func buildGenesisJson(addressesToFund []string) (string, error) {
// 	for i := range addressesToFund {
// 		if has0xPrefix(addressesToFund[i]) {
// 			addressesToFund[i] = addressesToFund[i][2:]
// 		}
// 	}

// 	data := struct {
// 		AddressesToFund []string
// 	}{
// 		AddressesToFund: addressesToFund,
// 	}

// 	t, err := template.New("genesis-json").Funcs(funcMap).Parse(Eth1GenesisJSON)
// 	if err != nil {
// 		fmt.Println("Error parsing template:", err)
// 		os.Exit(1)
// 	}

// 	var buf bytes.Buffer
// 	err = t.Execute(&buf, data)

// 	return buf.String(), err
// }

// var funcMap = template.FuncMap{
// 	// The name "inc" is what the function will be called in the template text.
// 	"decrement": func(i int) int {
// 		return i - 1
// 	},
// }

// func has0xPrefix(str string) bool {
// 	return len(str) >= 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X')
// }
