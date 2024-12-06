# Local Docker Image Builds

In addition to [this common setup](nodeset_environment.md) you can also provide your local image path and quickly rebuild it automatically before starting the test.

Create a configuration file `smoke.toml`
```toml
[blockchain_a]
  type = "anvil"
  docker_cmd_params = ["-b", "1"]

[nodeset]
  nodes = 5
  override_mode = "all"
  
  [nodeset.db]
    image = "postgres:12.0"

  [[nodeset.node_specs]]

    [nodeset.node_specs.node]
      # Dockerfile path is relative to "docker_ctx"
      docker_file = "core/chainlink.Dockerfile"
      docker_ctx = "../.."
```

These paths will work for `e2e/capabilities` in our main [repository](https://github.com/smartcontractkit/chainlink/tree/ctf-v2-tests/e2e/capabilities)

Also check how you can add rebuild to your [components](../developing/developing_components.md#building-local-images).

Summary:
- We learned how we can quickly re-build local docker image for CL node
