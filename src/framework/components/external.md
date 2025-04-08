# External Environment

### Using remote components

Because components are decoupled through outputs, you can use a cached config and switch outputs to any deployed infrastructure, such as staging. 

This allows you to reuse the same testing logic for behavior validation.

For example, to integrate with remote `k8s` environment you can use `CTF_CONFIGS=smoke_external.toml` and override all the outputs of components to connect to your remote env.

```toml
[blockchain_a]

  [blockchain_a.out]
    chain_id = "1337"
    use_cache = true

    [[blockchain_a.out.nodes]]
      # set up your RPC URLs
      http_url = "http://127.0.0.1:8545"
      ws_url = "ws://127.0.0.1:8545"

[[nodesets]]

  [[nodesets.node_specs]]
  ...

  [nodesets.out]
    use_cache = true

    [[nodesets.out.cl_nodes]]
      use_cache = true

      [nodesets.out.cl_nodes.node]
        # set up your user/password for API authorization
        api_auth_user = 'notreal@fakeemail.ch'
        api_auth_password = 'fj293fbBnlQ!f9vNs'
        # set up each node URLs
        p2p_url = "http://127.0.0.1:12000"
        url = "http://127.0.0.1:10000"

      [nodesets.out.cl_nodes.postgresql]
        # set up a database URL so tests can connect to your database if needed
        url = "postgresql://chainlink:thispasswordislongenough@127.0.0.1:13000/db_0?sslmode=disable"
      
      # more nodes in this array, configuration is the same ...
```
