---
layout: default
title: Config
nav_order: 2
parent: Setup
---

# Config

The framework draws on 2 config files, `framework.yaml` and `networks.yaml`. Whether you're working in the framework's main repo, or importing the framework as a library, it will look for these two config files in order to launch, configure, and connect all the testing resources.

## `framework.yaml`

This config handles how the framework should handle logging, deploying, and tearing down test environments. Check out the example below, or the most up-to-date version [here](https://github.com/smartcontractkit/integrations-framework/blob/main/framework.yaml).

Location of this file can be overridden by setting a file path as the `FRAMEWORK_CONFIG_FILE` environment variable.

```yaml
# Retains default and configurable name-values for the integration framework
#
# All configuration can be set at runtime with environment variables, for example:
# KEEP_ENVIRONMENTS=OnFail
keep_environments: Never # Options: Always, OnFail, Never
logging:
    # panic=5, fatal=4, error=3, warn=2, info=1, debug=0, trace=-1
    level: 0

# Specify the image and version of the Chainlink image you want to run tests against. Leave blank for default.
chainlink_image: public.ecr.aws/chainlink/chainlink
chainlink_version: 1.2.1
chainlink_env_values:

# Setting an environment file allows for persistent, not ephemeral environments on test execution
#
# For example, if an environment is created with helmenv CLI, then the YAML file outputted on creation can be
# referenced for use of that environment during all the tests
environment_file:
```

## `networks.yaml`

This file handles the settings for each network you want to connect to. This is a truncated version, see the full one [here](https://github.com/smartcontractkit/integrations-framework/blob/main/networks.yaml).

Location of this file can be overridden by setting a file path as the `NETWORKS_CONFIG_FILE` environment variable. **The first private key listed will be considered the default one to use for deploying contracts and funding addresses**.

```yaml
private_keys: &private_keys 
  private_keys: # Private keys that are used for simulated networks. These are publicly known keys for use only in simulated networks.
    - ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
    - 59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d
    - 5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a
    - 7c852118294e51e653712a81e05800f419141751be58f605c371e15141b007a6
    - 47e179ec197488593b187f80a00eb0da91f1b9d0b13f8733639f19c30a34926a
    - 8b3a350cf5c34c9194ca85829a2df0ec3153be0318b5e2d3348e872092edffba
    - 92db14e403b83dfe3df233f83dfa3a0d7096f21ca9b0d6d6b8d88b2b4ec1564e
    - 4bbbf85ce3377467afe5d46f804f221813b2bb87f24d81f60f1fcdbf7cbf4356
    - dbda1821b80551c9d65939329250298aa3472ba22feea921c0cf5d620ea67b97
    - 2a871d0798f97d79848a013d4936a73bf4cc922c825d33c1cf7073dff6d409c6
    - f214f2b2cd398c806f84e317254e0f0b801d0643303237d97a22a48e01628897
    - 701b615bbdfb9de65240bc28bd21bbc0d996645a3dd57e7b12bc2bdf6f192c82
    - a267530f49f8280200edf313ee7af6b827f2a8bce2897751d06a843f644967b1
    - 47c99abed3324a2707c28affff1267e45918ec8c3f20b8aa892e8b065d2942dd
    - c526ee95bf44d8fc405a158bb884d9d1238d99f0612e9f33d006bb0789009aaa
    - 8166f546bab6da521a8369cab06c5d2b9e46670292d85c875ee9ec20e84ffb61
    - ea6c44ac03bff858b476bba40716402b03e41b8e97e276d1baec7c37d42484a0
    - 689af8efa8c651a91ad287602527f3af2fe9f6501a7ac4b061667b5a93e037fd
    - de9be858da4a475276426320d5e9262ecfc3ba460bfac56360bfa6c4c28b4ee0
    - df57089febbacf7ba0bc227dafbffa9fc08a93fdc68e1e42411a14efcf23656e

selected_networks: # Selected network(s) for test execution
  - "geth"

networks:
  ... # See full file
```

Each network will have an identifier that can be listed in the `selected_networks` section.

```yaml
geth:                                 # Network identifier
  name: "Ethereum Geth dev"           # Human-readable name for network
  chain_id: 1337                      # ETH chain ID
  type: eth_simulated                 # eth_simulated or eth_testnet
  secret_private_keys: false          # Experimental feature for storing private keys as Kubernetes secrets
  namespace_for_secret: default       # Experimental feature for storing private keys as Kubernetes secrets
  private_keys:                       # List of private keys for this network, used for funding and deploying contracts
    - ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
  chainlink_transaction_limit: 500000 # Estimated gas limit that a single Chainlink tx might take (used for funding Chainlink nodes)
  transaction_timeout: 2m             # Duration of how long for the framework to wait for a confirmed transaction before timeout
  minimum_confirmations: 1            # How many blocks to wait for transaction to be confirmed
  gas_estimation_buffer: 10000        # How much gas to bump transaction an contract creations by (added to auto-estimations)
  block_gas_limit: 40000000           # How much gas each block of the network should be using
```

There are a couple values available for launching simulated Geth instances. If you choose them in your `selected_networks` list, they will launch with the following properties:

### `geth`

The default geth instance, small footprint with fast block times.

```yaml
resources:
  requests:
    cpu: .2
    memory: 1000Mi

config_args:
  "--dev.period": "1"
  "--miner.threads": "1"
  "--miner.gasprice": "10000000000"
  "--miner.gastarget": "80000000000"
```

### `geth_performance`

Used for performance tests, launching a powerful geth instance with large blocks and fast block times.

```yaml
resources:
  requests:
    cpu: 4
    memory: 4096Mi
  limits:
    cpu: 4
    memory: 4096Mi

config_args:
  "--dev.period": "1"
  "--miner.threads": "4"
  "--miner.gasprice": "10000000000"
  "--miner.gastarget": "30000000000"
```

### `geth_realistic`

Launches a powerful geth instance that tries to simulate Ethereum mainnet as close as possible.

```yaml
resources:
  requests:
    cpu: 4
    memory: 4096Mi
  limits:
    cpu: 4
    memory: 4096Mi

config_args:
  "--dev.period": "14"
  "--miner.threads": "4"
  "--miner.gasprice": "10000000000"
  "--miner.gastarget": "15000000000"
```
