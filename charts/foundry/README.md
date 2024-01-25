## Anvil node for testing
See example [values.yaml](values.yaml) and [Dockerfile](Dockerfile) if you need to publish a custom image

Runs anvil node with an optional `--fork-url` and mines new block every `1s` by default

Change the URL in [values.yaml](values.yaml) and deploy

```
anvil:
  host: "0.0.0.0"
  port: "8545"
  blockTime: 1
  forkURL: "https://goerli.infura.io/v3/..."
```