# Chainlink Node Set Environment Test

Let's use some external capability binaries in our tests and extend the [previous one](nodeset_environment.md).

We'll use a private repository example, so you should be authorized with [gh]()
```
gh auth login
gh auth setup-git
```

Download an example capability binary
```
export export GOPRIVATE=github.com/smartcontractkit/capabilities
go get github.com/smartcontractkit/capabilities/kvstore && go install github.com/smartcontractkit/capabilities/kvstore 
```

Create a configuration file `smoke.toml`
```toml
[blockchain_a]
  chain_id = "31337"
  image = "f4hrenh9it/foundry:latest"
  port = "8545"
  type = "anvil"

[nodeset]
  nodes = 5
  override_mode = "all"
  
  [nodeset.db]
    image = "postgres:15.6"

  [[nodeset.node_specs]]

    [nodeset.node_specs.node]
      # path to your capability binaries
      capabilities = ["./kvstore"]
      # default capabilities directory
      # capabilities_container_dir = "/home/capabilities"
      image = "public.ecr.aws/chainlink/chainlink:v2.17.0"
```

Run it
```bash
CTF_CONFIGS=smoke.toml go test -v -run TestNodeSet
```

Now you can configure your capability using `clclient.CreateJobRaw($raw_toml)`.

Capabilities are uploaded to `/home/capabilities` by default.

Summary:
- We deployed a node set with some capabilities


