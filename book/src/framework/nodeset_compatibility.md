# Chainlink Node Set Compatibility Testing Environment

The difference between this and [basic node set configuration](nodeset_environment.md) is that here you can provide any custom configuration for CL nodes.

Create a configuration file `smoke.toml`
```toml
[blockchain_a]
  type = "anvil"
  docker_cmd_params = ["-b", "1"]

[[nodesets]]
  name = "don"
  nodes = 5
  override_mode = "each"

  [nodesets.db]
    image = "postgres:12.0"

  [[nodesets.node_specs]]

    [nodesets.node_specs.node]
      image = "public.ecr.aws/chainlink/chainlink:v2.17.0"
      user_config_overrides = "      [Log]\n      level = 'info'\n      "
      user_secrets_overrides = ""

  [[nodesets.node_specs]]

    [nodesets.node_specs.node]
      image = "public.ecr.aws/chainlink/chainlink:v2.17.0"
      user_config_overrides = "      [Log]\n      level = 'info'\n      "
      user_secrets_overrides = ""

  [[nodesets.node_specs]]

    [nodesets.node_specs.node]
      image = "public.ecr.aws/chainlink/chainlink:v2.17.0"
      user_config_overrides = "      [Log]\n      level = 'info'\n      "
      user_secrets_overrides = ""

  [[nodesets.node_specs]]

    [nodesets.node_specs.node]
      image = "public.ecr.aws/chainlink/chainlink:v2.17.0"
      user_config_overrides = "      [Log]\n      level = 'info'\n      "
      user_secrets_overrides = ""

  [[nodesets.node_specs]]

    [nodesets.node_specs.node]
      image = "public.ecr.aws/chainlink/chainlink:v2.17.0"
      user_config_overrides = "      [Log]\n      level = 'info'\n      "
      user_secrets_overrides = ""
```

You can reuse `smoke_test.go` from previous [setup](nodeset_environment.md)

Run it
```bash
CTF_CONFIGS=smoke.toml go test -v -run TestNodeSet
```

Summary:
- We deployed fully-fledged set of Chainlink nodes connected to some blockchain and faked external data provider
- We understood how we can test different versions of Chainlink nodes for compatibility and override configs
