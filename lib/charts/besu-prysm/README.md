# VERY IMPORTANT
**There needs to be least one node labeled with `eth2=true` for deployment to work! (see below how to do it)**

# Deployment initalisation flow
Each `StatefulSet` has the same `initContainers` that are responsible for:
1. Generating validator keys
2. Generating eth1 and eth2 genesis
3. Generating common passwords, etc

It's crucial that the chart is installed either via `install.sh` script or if installing it manually with identical values of `currentUnixTimestamp` for current package and `eth2-common` package, which can be achieved by running Helm install with:
```
now=$(date +%s)
...
--set "genesis.values.currentUnixTimestamp"="$now" --set "eth2-common.genesis.values.currentUnixTimestamp"="$now"
```

That's because we now generate genesis independently for each of the components, but they all need to have the same genesis time. Thanks to that we don't need to use a persistent volume with RWX access mode, which is not supported by most storage classes.

# Default ports
Note: These are ports that k8s services are exposed on, not localhost ports as local port forwarding has to be setup manually, but in order to access the RPC you should consider forwarding HTTP and WS RPCs ports.

## Besu
* `8544` - HTTP RPC
* `8545` - WS RPC
* `8551` - Execution RPC (used by consensus clients)

## Prysm
* `3500` - HTTP Query RPC
* `4000` - HTTP RPC (used by validators)
* `8080` - HTTP Metrics (useful endpoint is `/healtz`)

## Validator
None

# Configuration options
Description of only some selected, important options:
``` yaml
eth2-common:
  general:
    # network id that will be used for beacon chain and validator
    networkId: 1337
  genesis:
    values:
      # current timestamp in seconds that will be used to generate genesis time
      currentUnixTimestamp: 1600000000
general:
  # network id that will be used for execution client
  networkId: 1337

shared:
  genesis:
    values:
      # how many seconds should each slot last for validators to submit attestations
      secondsPerSlot: 12
      # how many slots should each epoch have (lower => shorter epoch => faster finality)
      slotsPerEpoch: 4
      # how many seconds in the future should the genesis time be set (this has to be after beacon chain starts )
      delaySeconds: 20
      # how many validators should the network have
      validatorCount: 8
      # array of adddresses that should be prefunded with ETH
      preminedAddresses:
        - "f39fd6e51aad88f6f4ce6ab8827279cfffb92266"
prysm:
  shared:
    # fee recipient for block validation
    feeRecipent: "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266"
    # how many validators should the network have
    validators: 8
    # how many seconds should initContainers or beacon chain and validator wait for Geth to start
    gethInitTimeoutSeconds: 600
storage:
  # size of persistent volume
  size: 2Gi
```

**Important**: remember to set `networkId` both for `eth2-common` and `general` sections, as otherwise you will end up in an inconsistent state. You need to override `networkId` for the `eth2-common` dependency.

# Requirements
1. `kubectl` installed
2. `helm` installed
3. Access to a remote k8s cluster or a local one (Docker Desktop will do)

# Usage
1. Connect with kubectl to the cluster you want to deploy to
2. Set the context/namespace you want to use (if the namespace doesn't exist you might need to create it manually)
3. Run `./install.sh`
This command uses `values.yaml` file while generating one value on the fly: `currentUnixTimestamp`. This is super important as all the components need to have the same genesis time.

That script will run lints, prepare a package and then install it.

Then you should wait for `chain-ready` container to become ready, as that will mean that chain started to produce blocks. You can check it's logs to see current latest unfinalized block.

It's recommended to remove the installation with `./uninstall.sh` script, as it will remove all persistent volume claims from the namespace (something that `helm uninstall` doesn't do).

# Limitations
* No support for restarting of beacon chain pod (validator gets slashed)
* Untested scalability
* no working readiness/liveliness probe for beacon chain
