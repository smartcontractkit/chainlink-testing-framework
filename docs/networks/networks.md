---
layout: default
title: Networks
nav_order: 3
has_children: true
permalink: docs/networks
---

# Networks

Most of the time we run our tests on `ethereum_geth`, a simulated ethereum chain that makes our tests fast and cheap.
But some tests you'll want to run on other simulated chains, or even live test nets. You can see an extensive amount of
pre-configured networks defined in our [config.yml](../../config.yml) file. You can add more by adding to that list,
adjusting config values as necessary.

## Config Values

| Config Value           | Type   | Description |
|------------------------|--------|-------------|
|  name                  | string | Human-readable name of the network
|  chain_id              | int    | Chain ID number
|  type                  | string | The type of chain, right now we only really have `evm` at the moment, but are expanding
|  secret_private_keys   | bool   | Whether to look for private keys in you Kubernetes deployment or not. Default to false
|  namespace_for_secret  | string | Which namespace to look for private keys in, probably `default`
|  private_keys          | list   | List of private keys to use as wallets for tests. Usually the defaults for simulated chains.
|  transaction_limit     | int    | A limit on how much gas to use in a transaction
|  transaction_timeout   | time   | How long to wait before timing out on a transaction.
|  minimum_confirmations | int    | How many block confirmations to wait before declaring a transaction as successful.
|  gas_estimation_buffer | int    | A buffer to add to gas estimations if you want to ensure transactions going through
|  block_gas_limit       | int    | The gas limit per-block for the network
