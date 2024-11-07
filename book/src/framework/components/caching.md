## Component caching

We use component caching to accelerate test development and enforce idempotent test actions development.

Each component is isolated by means of inputs and outputs.

If cached config has any outputs with `use_cache = true` it will be used instead of deploying a component again.

```
export CTF_CONFIGS=smoke-cache.toml
```

### Using remote components

Because components are decoupled through outputs, you can use a cached config and switch outputs to any deployed infrastructure, such as staging. This allows you to reuse the same testing logic for behavior validation.

Example:
```
[blockchain_a.out]
use_cache = true
chain_id = '31337'

[[blockchain_a.out.nodes]]
ws_url = 'ws://127.0.0.1:33447'
http_url = 'http://127.0.0.1:33447'
docker_internal_ws_url = 'ws://anvil-3716a:8900'
docker_internal_http_url = 'http://anvil-3716a:8900'
```
Set flag `use_cache = true` on any component output, change output fields as needed and run your test again.