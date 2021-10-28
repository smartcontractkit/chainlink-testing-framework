---
layout: default
title: Different Networks
nav_order: 1
parent: Networks
---

# Using Different Networks

See the `network_configs` in our [config.yml](../config.yml) file to see our list of supported networks. Most of them
are different EVM chains, whether test nets, or simulated chains. Non-EVM chains are trickier, but we're actively working
on integrating as many as we can.

If you want to run your tests on a different network, you can supply it as an environment variable, `NETWORKS` or change
the config file. Supply a comma-seperated list of networks you wish to run on.

```sh
NETWORKS=ethereum_kovan,ethereum_kovan
```
