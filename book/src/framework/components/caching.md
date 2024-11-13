## Component caching

We use component caching to accelerate test development and enforce idempotent test actions development.

Each component is isolated by means of inputs and outputs.

If cached config has any outputs with `use_cache = true` it will be used instead of deploying a component again.

```
export CTF_CONFIGS=smoke-cache.toml
```
