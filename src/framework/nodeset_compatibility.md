# Chainlink Node Set Compatibility Testing Environment

The difference between this and [basic node set configuration](nodeset_environment.md) is that here you can provide any custom configuration for CL nodes.

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
  override_mode = "each"

  [[nodeset.node_specs]]

    [nodeset.node_specs.db]
      image = "postgres:15.6"
      pull_image = true

    [nodeset.node_specs.node]
      image = "public.ecr.aws/chainlink/chainlink:v2.17.0"
      pull_image = true
      user_config_overrides = "      [Log]\n      level = 'info'\n      "
      user_secrets_overrides = ""

  [[nodeset.node_specs]]

    [nodeset.node_specs.db]
      image = "postgres:15.6"
      pull_image = true

    [nodeset.node_specs.node]
      image = "public.ecr.aws/chainlink/chainlink:v2.17.0"
      pull_image = true
      user_config_overrides = "      [Log]\n      level = 'info'\n      "
      user_secrets_overrides = ""

  [[nodeset.node_specs]]

    [nodeset.node_specs.db]
      image = "postgres:15.6"
      pull_image = true

    [nodeset.node_specs.node]
      image = "public.ecr.aws/chainlink/chainlink:v2.17.0"
      pull_image = true
      user_config_overrides = "      [Log]\n      level = 'info'\n      "
      user_secrets_overrides = ""

  [[nodeset.node_specs]]

    [nodeset.node_specs.db]
      image = "postgres:15.6"
      pull_image = true

    [nodeset.node_specs.node]
      image = "public.ecr.aws/chainlink/chainlink:v2.17.0"
      pull_image = true
      user_config_overrides = "      [Log]\n      level = 'info'\n      "
      user_secrets_overrides = ""

  [[nodeset.node_specs]]

    [nodeset.node_specs.db]
      image = "postgres:15.6"
      pull_image = true

    [nodeset.node_specs.node]
      image = "public.ecr.aws/chainlink/chainlink:v2.17.0"
      pull_image = true
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
