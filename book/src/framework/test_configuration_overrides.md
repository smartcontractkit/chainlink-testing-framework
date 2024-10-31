# Overriding Test Configuration

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
