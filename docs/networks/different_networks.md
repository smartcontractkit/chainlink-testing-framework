---
layout: default
title: Different Networks
nav_order: 1
parent: Networks
---

# Using Different Networks

See the `network_configs` in our [config.yml](../../config.yml) file to see our list of supported networks. Most of them
are different EVM chains, whether test nets, or simulated chains. Non-EVM chains are trickier, but we're actively working
on integrating as many as we can.

If you want to run your tests on a different network, you can supply it as an environment variable, `NETWORKS` or change
the config file. Supply a comma-seperated list of networks you wish to run on.

```sh
NETWORKS=ethereum_kovan,ethereum_geth
```

or change the networks in your config file

```yaml
networks: # Selected network(s) for test execution
    - "ethereum_kovan"
    - "ethereum_geth"
```

## Non-EVM Networks

You'll notice that our network support doesn't have much besides typical EVM based networks. We're actively working
to expand to other network types however, and hope to have more to say on that soon.

## Simulated vs Live Networks

We've found it easiest and most effective to only use simulated networks when doing our testing, for a host of reasons
(speed, wallet management, even testnet currencies can be tricky to get large amounts of, etc...) and generally we'd
recommend you do the same. So, as of writing this, support for live test nets is a bit shaky. We'd obviously like to
improve on this, but it's not high on our priority list at the moment. However, you're welcome to lend us a hand on this
front by opening a PR!
