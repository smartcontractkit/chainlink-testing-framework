# Configuration

We generate config specification for all our components automatically.

Here is full specification for what can be specificed in `env.toml` for devenv.

`.out` fields are useful in case you want to use remote component, or integrate in other language than `Go`, otherwise can be ignored.

## Fakes

Fake HTTP service with recording capabilities.

```toml
{{#include fake.toml}}
```

## Blockchains

Various types of blockchain nodes.

```toml
{{#include blockchains.toml}}
```

## Node Set

Chainlink node set forming a Decentralized Oracle Network.

```toml
{{#include nodesets.toml}}
```
