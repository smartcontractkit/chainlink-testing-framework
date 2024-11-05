# NodeSet + External Blockchain Network

It is simple to replace any component by declaring the outputs!

Check our [previous](./nodeset_environment.md) test and add another config file on top.

Create a configuration file `smoke-fuji.toml`
```toml
[blockchain_a]

  [blockchain_a.out]
    chain_id = "43113"
    use_cache = true

    [[blockchain_a.out.nodes]]
      docker_internal_http_url = "https://ava-testnet.public.blastapi.io/ext/bc/C/rpc"
      docker_internal_ws_url = "wss://avalanche-fuji-c-chain-rpc.publicnode.com"
      http_url = "https://ava-testnet.public.blastapi.io/ext/bc/C/rpc"
      ws_url = "wss://avalanche-fuji-c-chain-rpc.publicnode.com"

```

Set both configs and replace your private key:
```bash
export PRIVATE_KEY=...
export CTF_CONFIGS=smoke.toml,smoke-fuji.toml
```

Run it
```bash
CTF_CONFIGS=smoke.toml go test -v -run TestNodeSet
```

Summary:
- We deployed fully-fledged set of Chainlink nodes connected to an external blockchain

To understand outputs of various components [search here](./components/overview.md)


