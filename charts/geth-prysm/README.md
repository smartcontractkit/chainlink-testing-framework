# Deployment initalisation flow
1. Generate validator keys
2. Generate eth1 and eth2 genesis
3. Start Geth
4. Generate common passwords, etc and save on shared volume
5. Wait for Geth to start
6. Start Prysm beacon chain
7. Start Prysm validator
8. Wait for first block to be produced (that's when `chain-ready` pod becomes ready)

# Default ports 
Note: These are ports that k8s services are exposed on, not localhost ports as local port forwarding has to be setup manually, but in order to access the RPC you should consider forwarding HTTP and WS RPCs ports.

## Geth
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
  # storage class to use for persistent volume that will be used to share data betwen containers
  class: hostpath
  # size of persistent volume
  size: 2Gi
```

# Requirements
1. `kubectl` installed
2. `helm` installed
3. Access to a remote k8s cluster or a local one (Docker Desktop will do)

# Usage
1. Connect with kubectl to the cluster you want to deploy to
2. Set the context/namespace you want to use (if the namespace doesn't exist you might need to create it manually)
3. Make sure that there's 1 node with with label `eth2=true` (this is used to schedule beacon chain and validator pods affinity to make sure they are deployed on the same node and have access to the same persistent volume). You can check it by running `kubectl get nodes --selector=eth2=true`. If there's no such node (which will especially be true on your local cluster, when running for the first time), run `kubectl get nodes --show-labels` to see all nodes and then pick one and run `kubectl label nodes <node-name> eth2=true` to add the label to it. It's best if you *don't do that* on remote clusters without previous consultation with the cluster owners. 
Once you have one labeled node you can proceed with chart installation.
3. Run `./install.sh`
This command by default uses `values.yaml` file, which is meant for local cluster use (because of the storage class it uses). If you want to deploy to SDLC cluster you should execute `./install.sh sdlc` and it will take values from `values-sdlc.yaml`, which uses `longhorn` storage class that is available in SDLC cluster and supports `ReadWriteMany` access mode (which is crucial, because multiple pods are using the same persistent volume).

That script will run lints, prepare a package and then install it.

Then you should wait for `chain-ready` container to become ready, as that will mean that chain started to produce blocks. You can check it's logs to see current latest unfinalized block.

# Limitations
* No support for restarting of geth app (it will try to initialize the chain from scratch every time and that will fail, becuase it will try to generate genesis.json based on previous chain state)
* Untested scalability
* I wasn't able to add a working readiness/liveliness probe for beacon chain