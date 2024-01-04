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

It's up to the user to provide a way to read the config from file and unmarshal it into the struct. All of the building blocks do and should implement the following interface:
```golang
type GenericConfig[T any] interface {
	Validate() error
	ApplyOverride(from T) error
}
```

`Validate()` should be used to ensure that the config is valid. `ApplyOverride()` should be used to apply overrides from another config. Some of the building blocks have also a `Default()` method that can be used to get default values.

Some more explanation is needed for the `NetworkConfig`:
```golang
type NetworkConfig struct {
	SelectedNetworks []string            `toml:"selected_networks"`
	RpcHttpUrls      map[string][]string `toml:"RpcHttpUrls"`
	RpcWsUrls        map[string][]string `toml:"RpcWsUrls"`
	WalletKeys       map[string][]string `toml:"WalletKeys"`
}

func (n *NetworkConfig) ApplySecrets() error {
    ...
}
```

It not only stores the configuration of selected networks and RPC endpoints and wallet keys, but via `ApplySecrets()` method provides a way to read from env var `BASE64_NETWORK_CONFIG` a base64-ed configuration of RPC endpoints and wallet keys. This could prove useful in the CI, where we could store as a secret a default configuration of stable endpoints, so that when we run a test job all that we have to provide is the network name and nothing more as it's pretty tedious, especially for on-demand jobs, to have to pass the whole RPC/wallet configuration every time you run it.

## Providing custom values in the CI

Up to this point when we wanted to modify some dynamic tests parameters in the CI we would simply set env vars. That approach won't work anymore. The way to go around it is to build a TOML file, base64 it, mask it and then set is as `BASE64_CONFIG_OVERRIDE` env var that will be read by tests. Here's an example of a working snippet of how that could look:
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
key="$PYROSCOPE_KEY"

[Logging]
test_log_collect=$test_log_collect
run_id="$RUN_ID"

[Logging.LogStream]
log_targets=$log_targets

[Logging.Loki]
tenant_id="$LOKI_TENANT_ID"
url="$LOKI_URL"
basic_auth="$LOKI_BASIC_AUTH"

[Logging.Grafana]
url="$GRAFANA_URL"
EOF

BASE64_CONFIG_OVERRIDE=$(cat config.toml | base64 -w 0)
echo ::add-mask::$BASE64_CONFIG_OVERRIDE
echo "BASE64_CONFIG_OVERRIDE=$BASE64_CONFIG_OVERRIDE" >> $GITHUB_ENV
```

These two lines in that very order are super important:
```bash
BASE64_CONFIG_OVERRIDE=$(cat config.toml | base64 -w 0)
echo ::add-mask::$BASE64_CONFIG_OVERRIDE
```

`::add-mask::` has to be called only after env var has been set to it's final value, otherwise it won't be recognized and masked properly and secrets will be exposed in the logs.