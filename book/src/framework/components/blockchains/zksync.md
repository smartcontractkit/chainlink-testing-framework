# ZKSync clients

We support [Anvil ZKSync](https://foundry-book.zksync.io/anvil-zksync/) memory blockchain in a Docker image.
It is a fork of [Anvil](https://book.getfoundry.sh/anvil/) with support for ZK Sync transactions and ZK VM.

Components are managed as [EVM](./evm) components.

> The component will create a temporary Dockerfile that pulls the Anvil ZKSync executables from https://raw.githubusercontent.com/matter-labs/foundry-zksync/main/install-foundry-zksync
> as per ZK Sync documentation.

## Configuration

Use `type: 'anvil-zksync'` to use Anvil ZKSync.

The configurable arguments are two:

- `chain_id`, defaults to `"260"`
- `port`, defaults to `"8011"`

## Test Private Keys

Testing keys are exported in the `blockchain` go module under `AnvilZKSyncRichAccountPks`.
