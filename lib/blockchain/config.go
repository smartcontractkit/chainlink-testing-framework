package blockchain

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
)

var (
	// SimulatedEVMNetwork ensures that the test will use a default simulated geth instance
	SimulatedEVMNetwork = EVMNetwork{
		Name:                 "Simulated Geth",
		ClientImplementation: EthereumClientImplementation,
		Simulated:            true,
		ChainID:              1337,
		URLs:                 []string{"ws://geth:8546"},
		HTTPURLs:             []string{"http://geth:8544"},
		PrivateKeys: []string{
			"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
			"59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d",
			"5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a",
			"8d351f5fc88484c65c15d44c7d3aa8779d12383373fb42d802e4576a50f765e5",
			"44fd8327d465031c71b20d7a5ba60bb01d33df8256fba406467bcb04e6f7262c",
			"809871f5c72d01a953f44f65d8b7bd0f3e39aee084d8cd0bc17ba3c386391814",
			"f29f5fda630ac9c0e39a8b05ec5b4b750a2e6ef098e612b177c6641bb5a675e1",
			"99b256477c424bb0102caab28c1792a210af906b901244fa67e2b704fac5a2bb",
			"bb74c3a9439ca83d09bcb4d3e5e65d8bc4977fc5b94be4db73772b22c3ff3d1a",
			"58845406a51d98fb2026887281b4e91b8843bbec5f16b89de06d5b9a62b231e8",
		},
		ChainlinkTransactionLimit: 500000,
		Timeout:                   StrDuration{2 * time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       10000,
	}
)

// EVMNetwork configures all the data the test needs to connect and operate on an EVM compatible network
type EVMNetwork struct {
	// Human-readable name of the network:
	Name string `toml:"evm_name" json:"evm_name"`
	// Chain ID for the blockchain
	ChainID int64 `toml:"evm_chain_id" json:"evm_chain_id"`
	// List of websocket URLs you want to connect to
	URLs []string `toml:"evm_urls" json:"evm_urls"`
	// List of websocket URLs you want to connect to
	HTTPURLs []string `toml:"evm_http_urls" json:"evm_http_urls"`
	// True if the network is simulated like a geth instance in dev mode. False if the network is a real test or mainnet
	Simulated bool `toml:"evm_simulated" json:"evm_simulated"`
	// Type of chain client node. Values: "none" | "geth" | "besu"
	SimulationType string `toml:"evm_simulation_type" json:"evm_simulation_type"`
	// List of private keys to fund the tests
	PrivateKeys []string `toml:"evm_keys" json:"evm_keys"`
	// Default gas limit to assume that Chainlink nodes will use. Used to try to estimate the funds that Chainlink
	// nodes require to run the tests.
	ChainlinkTransactionLimit uint64 `toml:"evm_chainlink_transaction_limit" json:"evm_chainlink_transaction_limit"`
	// How long to wait for on-chain operations before timing out an on-chain operation
	Timeout StrDuration `toml:"evm_transaction_timeout" json:"evm_transaction_timeout"`
	// How many block confirmations to wait to confirm on-chain events
	MinimumConfirmations int `toml:"evm_minimum_confirmations" json:"evm_minimum_confirmations"`
	// How much WEI to add to gas estimations for sending transactions
	GasEstimationBuffer uint64 `toml:"evm_gas_estimation_buffer" json:"evm_gas_estimation_buffer"`
	// ClientImplementation is the blockchain client to use when interacting with the test chain
	ClientImplementation ClientImplementation `toml:"client_implementation" json:"client_implementation"`
	// SupportsEIP1559 indicates if the client should try to use EIP1559 style gas and transactions
	SupportsEIP1559 bool `toml:"evm_supports_eip1559" json:"evm_supports_eip1559"`

	// Default gaslimit to use when sending transactions. If set this will override the transactionOptions gaslimit in case the
	// transactionOptions gaslimit is lesser than the defaultGasLimit.
	DefaultGasLimit uint64 `toml:"evm_default_gas_limit" json:"evm_default_gas_limit"`
	// Few chains use finality tags to mark blocks as finalized. This is used to determine if the chain uses finality tags.
	FinalityTag bool `toml:"evm_finality_tag" json:"evm_finality_tag"`
	// If the chain does not use finality tags, this is used to determine how many blocks to wait for before considering a block finalized.
	FinalityDepth uint64 `toml:"evm_finality_depth" json:"evm_finality_depth"`

	// TimeToReachFinality is the time it takes for a block to be considered final. This is used to determine how long to wait for a block to be considered final.
	TimeToReachFinality StrDuration `toml:"evm_time_to_reach_finality" json:"evm_time_to_reach_finality"`

	// Only used internally, do not set
	URL string `ignored:"true"`

	// Only used internally, do not set
	Headers http.Header `toml:"evm_headers" json:"evm_headers"`
}

// LoadNetworkFromEnvironment loads an EVM network from default environment variables. Helpful in soak tests
func LoadNetworkFromEnvironment() EVMNetwork {
	var network EVMNetwork
	if err := envconfig.Process("", &network); err != nil {
		log.Fatal().Err(err).Msg("Error loading network settings from environment variables")
	}
	log.Debug().Str("Name", network.Name).Int64("Chain ID", network.ChainID).Msg("Loaded Network")
	return network
}

// ToMap marshals the network's values to a generic map, useful for setting env vars on instances like the remote runner
// Map Structure
// "envconfig_key": stringValue
func (e *EVMNetwork) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"evm_name":                        e.Name,
		"evm_chain_id":                    fmt.Sprint(e.ChainID),
		"evm_urls":                        strings.Join(e.URLs, ","),
		"evm_http_urls":                   strings.Join(e.HTTPURLs, ","),
		"evm_simulated":                   fmt.Sprint(e.Simulated),
		"evm_simulation_type":             e.SimulationType,
		"evm_keys":                        strings.Join(e.PrivateKeys, ","),
		"evm_chainlink_transaction_limit": fmt.Sprint(e.ChainlinkTransactionLimit),
		"evm_transaction_timeout":         fmt.Sprint(e.Timeout),
		"evm_minimum_confirmations":       fmt.Sprint(e.MinimumConfirmations),
		"evm_gas_estimation_buffer":       fmt.Sprint(e.GasEstimationBuffer),
		"client_implementation":           fmt.Sprint(e.ClientImplementation),
	}
}

var (
	evmNetworkTOML = `[[EVM]]
ChainID = '%d'
MinContractPayment = '0'
%s`

	evmNodeTOML = `[[EVM.Nodes]]
Name = '%s'
WSURL = '%s'
HTTPURL = '%s'`
)

// MustChainlinkTOML marshals EVM network values into a TOML setting snippet. Will fail if error is encountered
// Can provide more detailed config for the network if non-default behaviors are desired.
func (e *EVMNetwork) MustChainlinkTOML(networkDetails string) string {
	if len(e.HTTPURLs) != len(e.URLs) || len(e.HTTPURLs) == 0 || len(e.URLs) == 0 {
		log.Fatal().
			Int("WS Count", len(e.URLs)).
			Int("HTTP Count", len(e.HTTPURLs)).
			Interface("WS URLs", e.URLs).
			Interface("HTTP URLs", e.HTTPURLs).
			Msg("Amount of HTTP and WS URLs should match, and not be empty")
		return ""
	}
	netString := fmt.Sprintf(evmNetworkTOML, e.ChainID, networkDetails)
	for index := range e.URLs {
		netString = fmt.Sprintf("%s\n\n%s", netString,
			fmt.Sprintf(evmNodeTOML, fmt.Sprintf("node-%d", index), e.URLs[index], e.HTTPURLs[index]))
	}

	return netString
}

// StrDuration is JSON/TOML friendly duration that can be parsed from "1h2m0s" Go format
type StrDuration struct {
	time.Duration
}

func (d *StrDuration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *StrDuration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case string:
		var err error
		d.Duration, err = time.ParseDuration(value)
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.New("invalid duration")
	}
}

// MarshalText implements the text.Marshaler interface (used by toml)
func (d StrDuration) MarshalText() ([]byte, error) {
	return []byte(d.Duration.String()), nil
}

// UnmarshalText implements the text.Unmarshaler interface (used by toml)
func (d *StrDuration) UnmarshalText(b []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(b))
	if err != nil {
		return err
	}
	return nil

}
