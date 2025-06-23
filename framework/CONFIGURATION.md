# Configuration Guide

Configuration is at the heart of the Chainlink Testing Framework (CTF). All components, networks, and services are configured using TOML files, which are loaded into Go structs for type-safe, validated test setup.

## Configuration Basics

- **TOML files** define the structure and parameters for all components in your test environment.
- **Go structs** mirror the TOML structure and provide validation.
- **Environment variables** can override or supplement TOML values for dynamic configuration.

## Example: Minimal Configuration

```toml
[blockchain_a]
  type = "anvil"
  image = "ghcr.io/foundry-rs/foundry"
  tag = "latest"
  pull_image = true

[chainlink_node]
  image = "public.ecr.aws/chainlink/chainlink"
  tag = "2.7.0"
  pull_image = true
```

## Loading Configuration in Go

```go
type Config struct {
    BlockchainA   *blockchain.Input `toml:"blockchain_a" validate:"required"`
    ChainlinkNode *clnode.Input     `toml:"chainlink_node" validate:"required"`
}

func TestMyTest(t *testing.T) {
    in, err := framework.Load[Config](t)
    require.NoError(t, err)
    // Use in.BlockchainA, in.ChainlinkNode, etc.
}
```

## TOML Structure

- Each top-level section represents a component or service.
- Fields map directly to struct fields in Go.
- Nested tables and arrays are supported for complex configurations.

### Example: Multi-Node, Multi-Chain

```toml
[blockchain_a]
  type = "anvil"
  image = "ghcr.io/foundry-rs/foundry"
  tag = "latest"
  pull_image = true

[blockchain_b]
  type = "geth"
  image = "ethereum/client-go"
  tag = "v1.12.0"
  pull_image = true

[[chainlink_nodes]]
  image = "public.ecr.aws/chainlink/chainlink"
  tag = "2.7.0"
  pull_image = true

[[chainlink_nodes]]
  image = "public.ecr.aws/chainlink/chainlink"
  tag = "2.7.0"
  pull_image = true
```

## Field Types

- **Strings**: `image = "..."`, `tag = "..."`
- **Booleans**: `pull_image = true`
- **Integers**: `port = 8545`
- **Arrays**: `urls = ["http://localhost:8545", "http://localhost:8546"]`
- **Tables**: Nested configuration for complex components

## Validation

- Use struct tags like `validate:"required"` to enforce required fields.
- CTF will fail fast if required fields are missing or invalid.

## Environment Variables

You can override or supplement TOML configuration with environment variables. This is useful for secrets, dynamic values, or CI/CD integration.

### Common Environment Variables

- `CTF_CONFIGS` – Path to the TOML config file(s)
- `CTF_LOG_LEVEL` – Set log verbosity (`info`, `debug`, `warn`, `error`)
- `CTF_DISABLE_CACHE` – Disable component caching (`true`/`false`)
- `CTF_NETWORK_NAME` – Set Docker network name
- `CTF_OBSERVABILITY` – Enable/disable observability stack
- `CTF_GRAFANA_URL` – Override Grafana URL
- `CTF_LOKI_URL` – Override Loki URL
- `CTF_PROMETHEUS_URL` – Override Prometheus URL

### Example Usage

```bash
CTF_CONFIGS=smoke.toml CTF_LOG_LEVEL=debug go test -v -run TestSmoke
```

## Advanced Configuration

### 1. **Component Caching**
- Use `out.use_cache = true` in TOML or `UseCache` in Go to enable caching
- See [Caching Guide](Caching) for details

### 2. **Secrets Management**
- Store secrets in environment variables or external files
- Reference them in TOML using `${ENV_VAR}` syntax if supported by your loader

### 3. **Multiple Config Files**
- You can specify multiple config files with `CTF_CONFIGS="a.toml,b.toml"`
- Later files override earlier ones for the same fields

### 4. **Dynamic Configuration in CI/CD**
- Use environment variables to inject dynamic values (e.g., image tags, URLs)
- Example:
  ```bash
  CTF_CONFIGS=ci.toml CTF_IMAGE_TAG=$GITHUB_SHA go test -v
  ```

## Best Practices

- **Keep configs minimal**: Only specify what you need for the test
- **Use arrays/tables** for repeated components (e.g., multiple nodes)
- **Validate** all required fields in your Go structs
- **Document** your config files for team clarity
- **Use environment variables** for secrets and dynamic values
- **Version control** your config files (except secrets)

## Troubleshooting

- **Missing fields**: Check for typos or missing required fields in TOML
- **Validation errors**: Ensure all `validate:"required"` fields are present
- **Environment overrides not working**: Confirm variable names and export status
- **Multiple configs**: Order matters; last file wins for duplicate fields

## Further Reading
- [Component System](Components)
- [Caching Guide](Caching)
- [Observability Guide](Observability)
- [Test Patterns](Test-Patterns)

## Test Configuration

In our framework, components control the configuration. Each component defines an Input that you embed into your test configuration. We automatically ensure consistency between TOML and struct definitions during validation.

An example of your component configuration:
```
// Input is a blockchain network configuration params declared by the component
type Input struct {
	Type                     string   `toml:"type" validate:"required,oneof=anvil geth" envconfig:"net_type"`
	Image                    string   `toml:"image" validate:"required"`
	Tag                      string   `toml:"tag" validate:"required"`
	Port                     string   `toml:"port" validate:"required"`
	ChainID                  string   `toml:"chain_id" validate:"required"`
	DockerCmdParamsOverrides []string `toml:"docker_cmd_params"`
	Out                      *Output  `toml:"out"`
}
```

How you use it in tests:
```
type Config struct {
	BlockchainA        *blockchain.Input `toml:"blockchain_a" validate:"required"`
}

func TestDON(t *testing.T) {
	in, err := framework.Load[Config](t)
	require.NoError(t, err)

	// deploy component
	bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
	require.NoError(t, err)
```

In `TOML`:
```
[blockchain_a]
chain_id = "1337"
image = "ghcr.io/foundry-rs/foundry:stable"
port = "8500"
type = "anvil"
docker_cmd_params = ["-b", "1"]
```

### Best practices for configuration and validation
- Avoid stateful types (e.g., loggers, clients) in your config.
- All `input` fields should include validate: "required", ensuring consistency between TOML and struct definitions.
- Add extra validation rules for URLs or "one-of" variants. Learn more here: go-playground/validator.

### Overriding configuration
To override any configuration, we merge multiple files into a single struct.

You can specify multiple file paths using `CTF_CONFIGS=path1,path2,path3`.

The framework will apply these configurations from right to left.