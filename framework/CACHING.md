# Caching Guide

Component caching is a powerful feature of the Chainlink Testing Framework (CTF) that enables faster test development by reusing previously deployed components and environments.

## What is Component Caching?

- **Component caching** allows you to skip the setup and deployment of components that have not changed since the last test run.
- Cached components are stored on disk and can be reused across multiple test runs, saving time and resources.
- Caching is especially useful for heavy components (e.g., blockchain networks, Chainlink nodes) that are expensive to start from scratch.

## How Caching Works

- Each component's configuration and output are hashed to create a unique cache key.
- If a component with the same configuration has already been deployed and cached, CTF will reuse the cached output instead of redeploying.
- If the configuration changes, the cache is invalidated and the component is redeployed.

## Enabling Caching

- Caching is enabled by default for most components.
- To enable caching for your custom component, include a `UseCache` field in the output struct and set it to `true`.

### Example

```go
type Output struct {
    UseCache    bool   `toml:"use_cache"`
    InternalURL string `toml:"internal_url"`
    ExternalURL string `toml:"external_url"`
}

func NewComponent(input *Input) (*Output, error) {
    if input.Out != nil && input.Out.UseCache {
        return input.Out, nil // Use cached output
    }
    // Deploy logic here
    return &Output{
        UseCache:    true,
        InternalURL: "http://container:8545",
        ExternalURL: "http://localhost:8545",
    }, nil
}
```

### In TOML

```toml
[blockchain_a]
  type = "anvil"
  use_cache = true
```

## Disabling Caching

- To force a fresh deployment, set `use_cache = false` in your TOML or `UseCache: false` in your Go struct.
- You can also disable caching globally with the environment variable:

```bash
export CTF_DISABLE_CACHE=true
```

## Cache Location

- By default, caches are stored in a `.ctf_cache` directory in your project root.
- You can configure the cache directory via environment variable:

```bash
export CTF_CACHE_DIR=/path/to/cache
```

## When to Use Caching

- **During development**: Rapidly iterate on tests without waiting for full environment setup
- **For stable components**: Reuse networks, databases, or nodes that don't change between tests
- **In CI/CD**: Use with care; ensure cache invalidation on config changes

## When Not to Use Caching

- **When testing upgrades or migrations**: Always start from a clean state
- **When debugging setup issues**: Disable cache to ensure fresh deployments
- **For tests that require full isolation**: Use fresh environments for each run

## Best Practices

- **Enable caching** for heavy, stable components during local development
- **Disable caching** for critical tests, upgrades, or when debugging
- **Document** which components are cacheable in your test README
- **Clean cache** periodically to avoid stale state:

```bash
rm -rf .ctf_cache
```

## Troubleshooting

- **Component not updating**: Check if `use_cache` is set to `true` and try disabling it
- **Stale state**: Clean the cache directory and rerun tests
- **Cache not being used**: Ensure `UseCache` is set in the output struct and TOML

## Further Reading
- [Component System](Components)
- [Configuration Guide](Configuration)
- [Observability Guide](Observability)
- [Test Patterns](Test-Patterns)

## Component caching

We use component caching to accelerate test development and enforce idempotent test actions.

Each component is designed to be pure with defined inputs/outputs.

You can use an environment variable to skip deployment steps and use cached outputs if your infrastructure is already running (locally or remotely):

```
export CTF_CONFIGS=smoke-cache.toml
```

### Using remote components

Because components are decoupled through outputs, you can use a cached config and switch outputs to any deployed infrastructure, such as staging. This allows you to reuse the same testing logic for behavior validation.

Example:
```
[blockchain_a.out]
use_cache = true
chain_id = '1337'

[[blockchain_a.out.nodes]]
ws_url = 'ws://127.0.0.1:33447'
http_url = 'http://127.0.0.1:33447'
internal_ws_url = 'ws://anvil-3716a:8900'
internal_http_url = 'http://anvil-3716a:8900'
```
Set flag `use_cache = true` on any component output and run your test again