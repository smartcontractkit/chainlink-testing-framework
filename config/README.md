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