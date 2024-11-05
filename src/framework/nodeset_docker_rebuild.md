# Local Docker Image Builds

In addition to [this common setup](nodeset_environment.md) you can also provide your local image path and quickly rebuild it automatically before starting the test.

Create a configuration file `smoke.toml`
```toml
[blockchain_a]
  chain_id = "31337"
  image = "f4hrenh9it/foundry:latest"
  port = "8545"
  type = "anvil"

[data_provider]
  port = 9111

[nodeset]
  nodes = 5
  override_mode = "all"

  [[nodeset.node_specs]]

    [nodeset.node_specs.db]
      image = "postgres:15.6"
      pull_image = true

    [nodeset.node_specs.node]
      docker_file = "../../core/chainlink.Dockerfile"
      docker_ctx = "../.."
      pull_image = true
```

These paths will work for `e2e/capabilities` in our main [repository](https://github.com/smartcontractkit/chainlink/tree/ctf-v2-tests/e2e/capabilities)

Summary:
- We learned how we can quickly re-build local docker image for CL node
