# Blockscout

You can use local [Blockscout](https://www.blockscout.com/) instance to debug EVM smart contracts.

Some images require `ghcr` auth, login first and pass the token:

```
gh auth token | docker login ghcr.io -u github-username --password-stdin
```

Start Blockscout

```
ctf bs up
```

Your `Blockscout` instance is up on [localhost](http://localhost)

To remove it, we also clean up all Blockscout databases to prevent stale data when restarting your tests.

```
ctf bs down
```

## Selecting Blockchain Node

By default, we connect to the first `anvil` node, but you can select the node explicitly

```
ctf bs -r http://host.docker.internal:8545 d
ctf bs -r http://host.docker.internal:8555 d
```

<div class="warning">

Blockscout isn’t ideal for local, ephemeral environments, as it won’t re-index blocks and transactions on test reruns. The easiest approach is to set up Blockscout first, initialize the test environment, switch to the [cache](../components/caching.md) config, and run tests without restarting RPC nodes.

Otherwise, use `ctf bs r` each time you restart your test with a fresh docker environment.

</div>
