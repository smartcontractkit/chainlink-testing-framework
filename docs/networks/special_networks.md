---
layout: default
title: Special Networks
nav_order: 3
parent: Networks
---

# Performance and Reorg Networks

For running performance and chaos tests, use `ethereum_geth_performance`. This creates our usual simulated geth network,
but with different settings tuned to allow for longer waits before transactions timeout.

For running tests with reorganization scenarios, `ethereum_geth_reorg`
