# Initalisation flow
1. Generate eth2 genesis
2. Generate eth1 genesis
3. Start Geth
4. Wait for Geth to start
5. Start Prysm beacon chain
6. Start Prysm validator
7. Wait for first block to be produced (that's when `chain-ready` pod becomes ready)

# Limitations
* No support for restarting of geth app (it will try to initialize the chain from scratch every time and that will fail, becuase it will try to generate genesis.json based on previous chain state)
* Untested scalability

# Default ports
## Geth
* `8544` - HTTP RPC
* `8545` - WS RPC
* `8551` - Execution RPC (used by consensus clients)
* `30303` - P2P

## Prysm
* `3500` - HTTP Query RPC
* `4000` - HTTP RPC (used by validators)
* `8080` - HTTP Metrics (useful endpoint is `/healtz`)

## Validator
None

# Configuration options
Description of only some selected, important options:
``` yaml
prysm: 
  shared: 
    # fee recipient for block validation
    feeRecipent: 0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266
    # how many validators should the network have
    validators: 8
    # how many seconds should initContainers or beacon chain and validator wait for Geth to start
    gethInitTimeoutSeconds: 600 
  genesis:
    values:
      # how many seconds should each slot last for validators to submit attestations
      secondsPerSlot: 12
      # how many slots should each epoch have (lower => shorter epoch => faster finality)
      slotsPerEpoch: 4
      # how many seconds in the future should the genesis time be set (this has to be after beacon chain starts )
      delaySeconds: 20        
storage:
  # storage class to use for persistent volume that will be used to share data betwen containers
  class: hostpath
  # size of persistent volume
  size: 2Gi
```

# Usage
1. Connect with kubectl to the cluster you want to deploy to
2. Set the context/namespace you want to use
3. Run `./reinstall-eth2.sh`

That script will use `Helm` to stop any existing deployment, remove PV and PVC, run lint, prepare a package and install a it.

Then you should wait for `chain-ready` container to become ready, as that will mean that chain started to produce blocks. You can check it's logs to see current latest unfinalized block.