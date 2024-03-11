# TOML Config

These basic building blocks can be used to create a TOML config file. For example:

```golang
import (
    ctf_config "github.com/smartcontractkit/chainlink-testing-framework/config"
	ctf_test_env "github.com/smartcontractkit/chainlink-testing-framework/docker/test_env"
)

type TestConfig struct {
	ChainlinkImage         *ctf_config.ChainlinkImageConfig `toml:"ChainlinkImage"`
	ChainlinkUpgradeImage  *ctf_config.ChainlinkImageConfig `toml:"ChainlinkUpgradeImage"`
	Logging                *ctf_config.LoggingConfig        `toml:"Logging"`
	Network                *ctf_config.NetworkConfig        `toml:"Network"`
	Pyroscope              *ctf_config.PyroscopeConfig      `toml:"Pyroscope"`
	PrivateEthereumNetwork *ctf_test_env.EthereumNetwork    `toml:"PrivateEthereumNetwork"`
}
```

It's up to the user to provide a way to read the config from file and unmarshal it into the struct. You can check [testconfig.go](../config/examples/testconfig.go) to see one way it chould be done..

`Validate()` should be used to ensure that the config is valid. Some of the building blocks have also a `Default()` method that can be used to get default values.

Also you might find `BytesToAnyTomlStruct(logger zerolog.Logger, filename, configurationName string, target any, content []byte) error` utility method useful for unmarshalling TOMLs read from env var or files into a struct

## Secrets in TOML config

For all values regarded as secrets, their keys should end with the `_secret` suffix. For example, use `basic_auth_secret="basic-auth"` instead of `basic_auth="basic-auth"`.

## Working example

For a full working example making use of all the building blocks see [testconfig.go](../config/examples/testconfig.go). It provides methods for reading TOML, applying overrides and validating non-empty config blocks. It supports 4 levels of overrides, in order of precedence:

- `BASE64_CONFIG_OVERRIDE` env var
- `overrides.toml`
- `[product_name].toml`
- `default.toml`

All you need to do now to get the config is execute `func GetConfig(configurationName string, product string) (TestConfig, error)`. It will first look for folder with file `.root_dir` and from there it will look for config files in all subfolders, so that you can place the config files in whatever folder(s) work for you. It assumes that all configuration versions for a single product are kept in `[product_name].toml` under different configuration names (that can represent anything you want: a single test, a test type, a test group, etc).

Overrides of config files are done in a super-simple way. We try to unmarshall consecutive files into the same struct. Since it's all pointer based only not-nil keys are overwritten.

## IMPORTANT!

It is **required** to add `overrides.toml` to `.gitignore` in your project, so that you don't accidentally commit it as it might contain secrets.

## Network config (and default RPC endpoints)

Some more explanation is needed for the `NetworkConfig`:

```golang
type NetworkConfig struct {
	// list of networks that should be used for testing
	SelectedNetworks []string            `toml:"selected_networks"`
	// map of network name to EVMNetworks where key is network name and value is a pointer to EVMNetwork
	// if not set, it will try to find the network from defined networks in MappedNetworks under known_networks.go
	// it doesn't matter if you use `arbitrum_sepolia` or `ARBITRUM_SEPOLIA` or even `arbitrum_SEPOLIA` as key
	// as all keys will be uppercased when loading the Default config
	EVMNetworks map[string]*blockchain.EVMNetwork `toml:"EVMNetworks,omitempty"`
	// map of network name to ForkConfigs where key is network name and value is a pointer to ForkConfig
	// only used if network fork is needed, if provided, the network will be forked with the given config
	// networkname is fetched first from the EVMNetworks and
	// if not defined with EVMNetworks, it will try to find the network from defined networks in MappedNetworks under known_networks.go
    ForkConfigs map[string]*ForkConfig `toml:"ForkConfigs,omitempty"`
	// map of network name to RPC endpoints where key is network name and value is a list of RPC HTTP endpoints
	RpcHttpUrls      map[string][]string `toml:"RpcHttpUrls"`
	// map of network name to RPC endpoints where key is network name and value is a list of RPC WS endpoints
	RpcWsUrls        map[string][]string `toml:"RpcWsUrls"`
	// map of network name to wallet keys where key is network name and value is a list of private keys (aka funding keys)
	WalletKeys       map[string][]string `toml:"WalletKeys"`
}

func (n *NetworkConfig) Default() error {
    ...
}
```

Sample TOML config:

```toml
selected_networks = ["arbitrum_goerli", "optimism_goerli", "new_network"]

[EVMNetworks.new_network]
evm_name = "new_test_network"
evm_chain_id = 100009
evm_simulated = true
evm_chainlink_transaction_limit = 5000
evm_minimum_confirmations = 1
evm_gas_estimation_buffer = 10000
client_implementation = "Ethereum"
evm_supports_eip1559 = true
evm_default_gas_limit = 6000000

[ForkConfigs.new_network]
url = "ws://localhost:8546"
block_number = 100

[RpcHttpUrls]
arbitrum_goerli = ["https://devnet-2.mt/ABC/rpc/"]
new_network = ["http://localhost:8545"]

[RpcWsUrls]
arbitrum_goerli = ["wss://devnet-2.mt/ABC/ws/"]
new_network = ["ws://localhost:8546"]

[WalletKeys]
arbitrum_goerli = ["1810868fc221b9f50b5b3e0186d8a5f343f892e51ce12a9e818f936ec0b651ed"]
optimism_goerli = ["1810868fc221b9f50b5b3e0186d8a5f343f892e51ce12a9e818f936ec0b651ed"]
new_network = ["ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"]
```

Whenver you are adding a new EVMNetwork to the config, you can either

- provide the rpcs and wallet keys in Rpc<Http/Ws>Urls and WalletKeys. Like in the example above, you can see that `new_network` is added to the `selected_networks` and `EVMNetworks` and then the rpcs and wallet keys are provided in `RpcHttpUrls`, `RpcWsUrls` and `WalletKeys` respectively.
- provide the rpcs and wallet keys in the `EVMNetworks` itself. Like in the example below, you can see that `new_network` is added to the `selected_networks` and `EVMNetworks` and then the rpcs and wallet keys are provided in `EVMNetworks` itself.

```toml

selected_networks = ["new_network"]

[EVMNetworks.new_network]
evm_name = "new_test_network"
evm_chain_id = 100009
evm_urls = ["ws://localhost:8546"]
evm_http_urls = ["http://localhost:8545"]
evm_keys = ["ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"]
evm_simulated = true
evm_chainlink_transaction_limit = 5000
evm_minimum_confirmations = 1
evm_gas_estimation_buffer = 10000
client_implementation = "Ethereum"
evm_supports_eip1559 = true
evm_default_gas_limit = 6000000
```

If your config struct looks like that:

```golang

type TestConfig struct {
	Network *ctf_config.NetworkConfig `toml:"Network"`
}
```

then your TOML file should look like that:

```toml
[Network]
selected_networks = ["arbitrum_goerli","new_network"]

[Network.EVMNetworks.new_network]
evm_name = "new_test_network"
evm_chain_id = 100009
evm_simulated = true
evm_chainlink_transaction_limit = 5000
evm_minimum_confirmations = 1
evm_gas_estimation_buffer = 10000
client_implementation = "Ethereum"
evm_supports_eip1559 = true
evm_default_gas_limit = 6000000

[Network.RpcHttpUrls]
arbitrum_goerli = ["https://devnet-2.mt/ABC/rpc/"]
new_network = ["http://localhost:8545"]

[Network.RpcWsUrls]
arbitrum_goerli = ["ws://devnet-2.mt/ABC/rpc/"]
new_network = ["ws://localhost:8546"]

[Network.WalletKeys]
arbitrum_goerli = ["1810868fc221b9f50b5b3e0186d8a5f343f892e51ce12a9e818f936ec0b651ed"]
new_network = ["ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"]
```

It not only stores the configuration of selected networks and RPC endpoints and wallet keys, but via `Default()` method provides a way to read from env var `BASE64_NETWORK_CONFIG` a base64-ed configuration of RPC endpoints and wallet keys. This could prove useful in the CI, where we could store as a secret a default configuration of stable endpoints, so that when we run a test job all that we have to provide is the network name and nothing more as it's pretty tedious, especially for on-demand jobs, to have to pass the whole RPC/wallet configuration every time you run it.

If in your product config you want to support case-insensitive network names and map keys remember to run `NetworkConfig.UpperCaseNetworkNames()` on your config before using it.

## Providing custom values in the CI

Up to this point when we wanted to modify some dynamic tests parameters in the CI we would simply set env vars. That approach won't work anymore. The way to go around it is to build a TOML file, `base64` it, mask it and then set is as `BASE64_CONFIG_OVERRIDE` env var that will be read by tests. Here's an example of a working snippet of how that could look:

```bash
convert_to_toml_array() {
	local IFS=','
	local input_array=($1)
	local toml_array_format="["

	for element in "${input_array[@]}"; do
		toml_array_format+="\"$element\","
	done

	toml_array_format="${toml_array_format%,}]"
	echo "$toml_array_format"
}

selected_networks=$(convert_to_toml_array "$SELECTED_NETWORKS")
log_targets=$(convert_to_toml_array "$LOGSTREAM_LOG_TARGETS")

if [ -n "$PYROSCOPE_SERVER" ]; then
	pyroscope_enabled=true
else
	pyroscope_enabled=false
fi

if [ -n "$ETH2_EL_CLIENT" ]; then
	execution_layer="$ETH2_EL_CLIENT"
else
	execution_layer="geth"
fi

if [ -n "$TEST_LOG_COLLECT" ]; then
	test_log_collect=true
else
	test_log_collect=false
fi

cat << EOF > config.toml
[Network]
selected_networks=$selected_networks

[ChainlinkImage]
image="$CHAINLINK_IMAGE"
version="$CHAINLINK_VERSION"

[Pyroscope]
enabled=$pyroscope_enabled
server_url="$PYROSCOPE_SERVER"
environment="$PYROSCOPE_ENVIRONMENT"
key_secret="$PYROSCOPE_KEY"

[Logging]
test_log_collect=$test_log_collect
run_id="$RUN_ID"

[Logging.LogStream]
log_targets=$log_targets

[Logging.Loki]
tenant_id="$LOKI_TENANT_ID"
url="$LOKI_URL"
basic_auth_secret="$LOKI_BASIC_AUTH"
bearer_token_secret="$LOKI_BEARER_TOKEN"

[Logging.Grafana]
url="$GRAFANA_URL"
EOF

BASE64_CONFIG_OVERRIDE=$(cat config.toml | base64 -w 0)
echo ::add-mask::$BASE64_CONFIG_OVERRIDE
echo "BASE64_CONFIG_OVERRIDE=$BASE64_CONFIG_OVERRIDE" >> $GITHUB_ENV
```

**These two lines in that very order are super important**

```bash
BASE64_CONFIG_OVERRIDE=$(cat config.toml | base64 -w 0)
echo ::add-mask::$BASE64_CONFIG_OVERRIDE
```

`::add-mask::` has to be called only after env var has been set to it's final value, otherwise it won't be recognized and masked properly and secrets will be exposed in the logs.

## Providing custom values for local execution

For local execution it's best to put custom variables in `overrides.toml` file.

## Providing custom values in k8s

It's easy. All you need to do is:

- Create TOML file with these values
- Base64 it: `cat your.toml | base64`
- Set the base64 result as `BASE64_CONFIG_OVERRIDE` environment variable.

Both `BASE64_CONFIG_OVERRIDE` and `BASE64_NETWORK_CONFIG` will be automatically forwarded to k8s (as long as they are set and available to the test process), when creating the environment programmatically via `environment.New()`.

Quick example:

```bash
BASE64_CONFIG_OVERRIDE=$(cat your.toml | base64) go test your-test-that-runs-in-k8s ./file/with/your/test
```

# Not moved to TOML

Not moved to TOML:

- `SLACK_API_KEY`
- `SLACK_USER`
- `SLACK_CHANNEL`
- `TEST_LOG_LEVEL`
- `CHAINLINK_ENV_USER`
- `DETACH_RUNNER`
- `ENV_JOB_IMAGE`
- most of k8s-specific env variables were left untouched
