# Blockscout

You can use local [Blockscout](https://www.blockscout.com/) instance to debug EVM smart contracts.

```
ctf bs up
```
Your `Blockscout` instance is up on [localhost](http://localhost)

To remove it, we also clean up all Blockscout databases to prevent stale data when restarting your tests.
```
ctf bs down
```

<div class="warning">

Blockscout isn’t ideal for local, ephemeral environments, as it won’t re-index blocks and transactions on test reruns. The easiest approach is to set up Blockscout first, initialize the test environment, switch to the [cache](../components/caching.md) config, and run tests without restarting RPC nodes. 

Otherwise, use `ctf bs r` each time you restart your test with a fresh docker environment.
</div>

<div class="warning">

Blockscout integration is still WIP, for now Blockscout reads only one node that is on `:8545`, all our blockchain implementation expose this port by default.
</div>
