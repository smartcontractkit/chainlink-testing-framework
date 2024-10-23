## Component caching

We use component caching to accelerate test development and enforce idempotent test actions.

Each component is designed to be pure with defined inputs/outputs.

You can use an environment variable to skip deployment steps and use cached outputs if your infrastructure is already running (locally or remotely):

```
export CTF_CONFIGS=smoke-cache.toml
export CTF_USE_CACHED_OUTPUTS=true
```

### Using remote components

Because components are decoupled through outputs, you can use a cached config and switch outputs to any deployed infrastructure, such as staging. This allows you to reuse the same testing logic for behavior validation.

Example:
```
[blockchain_a.out]
chain_id = '31337'

[[blockchain_a.out.nodes]]
ws_url = 'ws://127.0.0.1:33447'
http_url = 'http://127.0.0.1:33447'
docker_internal_ws_url = 'ws://anvil-3716a:8900'
docker_internal_http_url = 'http://anvil-3716a:8900'
```