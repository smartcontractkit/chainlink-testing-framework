# External Environment Integration

Developer environment exposes format which allows any external infra to be integrated.

## NodeSets + Blockchains + Fake

This is a basic example of data required to run developer environment commands on external infrastructure. Full list of fields can be found [here](https://smartcontractkit.github.io/chainlink-testing-framework/framework/developer_environment/toml.html)

Example `env-stage-1.toml`:
```toml
[fakes]
  [fakes.out]
    # URL to fakes server which represents some 3rd party which are mocked on external infrastructure
    base_url_host = 'https://chainlink-$product_name-fake:9111'

[[blockchains]]
  # One or more blockchains
  [blockchains.out]
    type = 'anvil'
    family = 'evm'
    chain_id = '1337'
    # One or more blockchain nodes
    [[blockchains.out.nodes]]
      ws_url = 'wss://blockchain-node-1:8545'
      http_url = 'https://blockchain-node-1:8545'

# One or more DON clusters
[[nodesets]]
  name = 'my-external-don-1'
  nodes = 2
  # Nodeset config output
  [nodesets.out]
    # First CL node connection data
    [[nodesets.out.cl_nodes]]
      [nodesets.out.cl_nodes.node]
        api_auth_user = 'some_user'
        api_auth_password = 'some_password'
        url = 'https://chainlink-node-1:6688'
      [nodesets.out.cl_nodes.postgresql]
        url = 'postgresql://chainlink:thispasswordislongenough@chainlink-node-1-db:13000/db_0?sslmode=disable'
    # Second node, etc..
    [[nodesets.out.cl_nodes]]
      [nodesets.out.cl_nodes.node]
        api_auth_user = 'some_user'
        api_auth_password = 'some_password'
        url = 'https://chainlink-node-2:6688'
      [nodesets.out.cl_nodes.postgresql]
        url = 'postgresql://chainlink:thispasswordislongenough@chainlink-node-2-db:13000/db_1?sslmode=disable'

# Second NodeSet, etc..
[[nodesets]]
  name = 'my-external-don-2'
  nodes = 2
  # All the other fields as in previous example above
```

Deploy the infrastructure first, then use `up env-stage-1.toml` to orchestrate staging environment.
