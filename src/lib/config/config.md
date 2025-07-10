# TOML Config

<div class="warning">

This documentation is deprecated, we are using it in [Chainlink Integration Tests](https://github.com/smartcontractkit/chainlink/tree/develop/integration-tests)

If you want to test our new products use [v2](../framework/overview.md)
</div>

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

It's up to the user to provide a way to read the config from file and unmarshal it into the struct. You can check [testconfig.go](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/lib/config/examples/testconfig.go) to see one way it could be done.

`Validate()` should be used to ensure that the config is valid. Some of the building blocks have also a `Default()` method that can be used to get default values.

Also, you might find `BytesToAnyTomlStruct(logger zerolog.Logger, filename, configurationName string, target any, content []byte) error` utility method useful for unmarshalling TOMLs read from env var or files into a struct

## Test Secrets

Test secrets are not stored directly within the `TestConfig` TOML due to security reasons. Instead, they are passed into `TestConfig` via environment variables. Below is a list of all available secrets. Set only the secrets required for your specific tests, like so: `E2E_TEST_CHAINLINK_IMAGE=qa_ecr_image_url`.

### Default Secret Loading

By default, secrets are loaded from the `~/.testsecrets` dotenv file. Example of a local `~/.testsecrets` file:

```bash
E2E_TEST_CHAINLINK_IMAGE=qa_ecr_image_url
E2E_TEST_CHAINLINK_UPGRADE_IMAGE=qa_ecr_image_url
E2E_TEST_ARBITRUM_SEPOLIA_WALLET_KEY=wallet_key
```

### All E2E Test Secrets

| Secret                        | Env Var                                                             | Example                                                                                                                                                                                         |
| ----------------------------- | ------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Chainlink Image               | `E2E_TEST_CHAINLINK_IMAGE`                                          | `E2E_TEST_CHAINLINK_IMAGE=qa_ecr_image_url`                                                                                                                                                     |
| Chainlink Upgrade Image       | `E2E_TEST_CHAINLINK_UPGRADE_IMAGE`                                  | `E2E_TEST_CHAINLINK_UPGRADE_IMAGE=qa_ecr_image_url`                                                                                                                                             |
| Wallet Key per network        | `E2E_TEST_(.+)_WALLET_KEY` or `E2E_TEST_(.+)_WALLET_KEY_(\d+)$`     | `E2E_TEST_ARBITRUM_SEPOLIA_WALLET_KEY=wallet_key` or `E2E_TEST_ARBITRUM_SEPOLIA_WALLET_KEY_1=wallet_key_1`, `E2E_TEST_ARBITRUM_SEPOLIA_WALLET_KEY_2=wallet_key_2` for multiple keys per network |
| RPC HTTP URL per network      | `E2E_TEST_(.+)_RPC_HTTP_URL` or `E2E_TEST_(.+)_RPC_HTTP_URL_(\d+)$` | `E2E_TEST_ARBITRUM_SEPOLIA_RPC_HTTP_URL=url` or `E2E_TEST_ARBITRUM_SEPOLIA_RPC_HTTP_URL_1=url`, `E2E_TEST_ARBITRUM_SEPOLIA_RPC_HTTP_URL_2=url` for multiple http urls per network               |
| RPC WebSocket URL per network | `E2E_TEST_(.+)_RPC_WS_URL` or `E2E_TEST_(.+)_RPC_WS_URL_(\d+)$`     | `E2E_TEST_ARBITRUM_RPC_WS_URL=ws_url` or `E2E_TEST_ARBITRUM_RPC_WS_URL_1=ws_url_1`, `E2E_TEST_ARBITRUM_RPC_WS_URL_2=ws_url_2` for multiple ws urls per network                                  |
| Loki Tenant ID                | `E2E_TEST_LOKI_TENANT_ID`                                           | `E2E_TEST_LOKI_TENANT_ID=tenant_id`                                                                                                                                                             |
| Loki Endpoint                 | `E2E_TEST_LOKI_ENDPOINT`                                            | `E2E_TEST_LOKI_ENDPOINT=url`                                                                                                                                                                    |
| Loki Basic Auth               | `E2E_TEST_LOKI_BASIC_AUTH`                                          | `E2E_TEST_LOKI_BASIC_AUTH=token`                                                                                                                                                                |
| Loki Bearer Token             | `E2E_TEST_LOKI_BEARER_TOKEN`                                        | `E2E_TEST_LOKI_BEARER_TOKEN=token`                                                                                                                                                              |
| Grafana Bearer Token          | `E2E_TEST_GRAFANA_BEARER_TOKEN`                                     | `E2E_TEST_GRAFANA_BEARER_TOKEN=token`                                                                                                                                                           |
| Pyroscope Server URL          | `E2E_TEST_PYROSCOPE_SERVER_URL`                                     | `E2E_TEST_PYROSCOPE_SERVER_URL=url`                                                                                                                                                             |
| Pyroscope Key                 | `E2E_TEST_PYROSCOPE_KEY`                                            | `E2E_TEST_PYROSCOPE_KEY=key`                                                                                                                                                                    |

### Run GitHub Workflow with Your Test Secrets

By default, GitHub workflows execute with a set of predefined secrets. However, you can use custom secrets by specifying a unique identifier for your secrets when running the `gh workflow` command.

#### Steps to Use Custom Secrets

1. **Upload Local Secrets to GitHub Secrets Vault:**

    - **Install `ghsecrets` tool:**
      Install the `ghsecrets` tool to manage GitHub Secrets more efficiently.

      ```bash
      go install github.com/smartcontractkit/chainlink-testing-framework/tools/ghsecrets@latest
      ```

      If you use `asdf`, run `asdf reshim`

    - **Upload Secrets:**
      Run `ghsecrets set` from local core repo to upload the content of your `~/.testsecrets` file to the GitHub Secrets Vault and generate a unique identifier (referred to as `your_ghsecret_id`).

   ```bash
   cd path-to-chainlink-core-repo
   ```

   ```bash
   ghsecrets set
   ```

   For more details about `ghsecrets`, visit https://github.com/smartcontractkit/chainlink-testing-framework/tree/main/tools/ghsecrets#faq

2. **Execute the Workflow with Custom Secrets:**
    - To use the custom secrets in your GitHub Actions workflow, pass the `-f test_secrets_override_key={your_ghsecret_id}` flag when running the `gh workflow` command.
      ```bash
      gh workflow run <workflow_name> -f test_secrets_override_key={your_ghsecret_id}
      ```

#### Default Secrets Handling

If the `test_secrets_override_key` is not provided, the workflow will default to using the secrets preconfigured in the CI environment.

### Creating New Test Secrets in TestConfig

When adding a new secret to the `TestConfig`, such as a token or other sensitive information, the method `ReadConfigValuesFromEnvVars()` in `config/testconfig.go` must be extended to include the new secret. Ensure that the new environment variable starts with the `E2E_TEST_` prefix. This prefix is crucial for ensuring that the secret is correctly propagated to Kubernetes tests when using the Remote Runner.

Here’s a quick checklist for adding a new test secret:

- Add the secret to ~/.testsecrets with the `E2E_TEST_` prefix to ensure proper handling.
- Extend the `config/testconfig.go:ReadConfigValuesFromEnvVars()` method to load the secret in `TestConfig`
- Add the secrets to [All E2E Test Secrets](#all-e2e-test-secrets) table.

## Working example

For a full working example making use of all the building blocks see [testconfig.go](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/lib/config/examples/testconfig.go). It provides methods for reading TOML, applying overrides and validating non-empty config blocks. It supports 4 levels of overrides, in order of precedence:

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

Whenever you are adding a new EVMNetwork to the config, you can either

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

`BASE64_CONFIG_OVERRIDE` will be automatically forwarded to k8s (as long as it is set and available to the test process), when creating the environment programmatically via `environment.New()`.

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
