# Test Configuration

Since end-to-end system-level test configurations can become complex, we use a single, generic method to marshal any configuration.

Here’s an example of how you can extend your configuration.

```golang
type Cfg struct {
    // Usually you'll have basic components
	BlockchainA        *blockchain.Input `toml:"blockchain_a" validate:"required"`
	MockerDataProvider *fake.Input       `toml:"data_provider" validate:"required"`
	NodeSet            *ns.Input         `toml:"nodeset" validate:"required"`
    // And some custom test-related data
    MyCustomTestData   string            `toml:"my_custom_test_data" validate:"required"`
}

func TestSmoke(t *testing.T) {
    in, err := framework.Load[Cfg](t)
    require.NoError(t, err)
    in.MyCustomTestData // do something
    ...
}
```

We use [validator](https://github.com/go-playground/validator) to make sure anyone can properly configure your test.

All basic components configuration is described in our docs, but it’s recommended to comment on your configuration in TOML.

Additionally, use `validate:"required"` or `validate:"required,oneof=anvil geth"` for fields with strict value options on your custom configuration.
```
# My custom config does X
[MyCustomConfig]
# Can be a number
a = 1
# Can be Y or Z
b = "Z"
```


## Overriding Test Configuration

To override any test configuration, we merge multiple files into a single struct.

You can specify multiple file paths using `CTF_CONFIGS=path1,path2,path3`.

The framework will apply these configurations from right to left and marshal them to a single test config structure.

Use it to structure the variations of your test, ex.:
```
export CTF_CONFIGS=smoke-test-feature-a-simulated-network.toml
export CTF_CONFIGS=smoke-test-feature-a-simulated-network.toml,smoke-test-feature-a-testnet.toml

export CTF_CONFIGS=smoke-test-feature-a.toml
export CTF_CONFIGS=smoke-test-feature-a.toml,smoke-test-feature-b.toml

export CTF_CONFIGS=load-profile-api-service-1.toml
export CTF_CONFIGS=load-profile-api-service-1.toml,load-profile-api-service-2.toml
```

This helps reduce duplication in the configuration.

> [!NOTE]
> We designed overrides to be as simple as possible, as frameworks like [envconfig](https://github.com/kelseyhightower/envconfig) and [viper](https://github.com/spf13/viper) offer extensive flexibility but can lead to inconsistent configurations prone to drift.
> 
> This feature is meant to override test setup configurations, not test logic. Avoid using TOML to alter test logic.
> 
> Tests should remain straightforward, readable, and perform a single set of actions (potentially across different CI/CD environments). If variations in test logic are required, consider splitting them into separate tests.

> [!WARNING]  
> When override slices remember that you should replace the full slice, it won't be extended by default!

## Overriding Components Configuration

The same override logic applies across components, files, and configuration fields in code, configs are applied in order:

1. Implicit component defaults that are defined inside application
2. Component defaults defined in the framework or external component, ex.: [CLNode](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/framework/components/clnode/default.go)
3. `Test.*Override`
4. `user_.*_overrides`

Use `Test.*Override` in test `code` to override component configurations, and `user_.*_overrides` in `TOML` for the same purpose.