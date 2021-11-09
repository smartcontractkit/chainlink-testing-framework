---
layout: default
title: Suite Setup
nav_order: 1
parent: Usage
---

The `SuiteSetup` interface provides a lot of convenient ways to setup and interact with your tests, and you should
probably be using it for most, if not all of your test scenarios.

```go
suiteSetup, err = actions.SingleNetworkSetup( // Indicating we only want to deal with a single blockchain network
  environment.NewChainlinkCluster(0),         // We're launching 0 chainlink nodes in this example
  hooks.EVMNetworkFromConfigHook,             // Default, gets network settings
  hooks.EthereumDeployerHook,                 // Default, creates a contract deployer for the network
  hooks.EthereumClientHook,                   // Default, establishes client connection to the network
  utils.ProjectRoot,                          // Default, the path of our config file.
)
Expect(err).ShouldNot(HaveOccurred())         // Make sure no errors happened
```

The `SuiteSetup` does a lot of work for us, allowing us to specify how many chainlink nodes we want to deploy,
and where to look for our config file (if we have one). There will be quite a few calls creating `SuiteSetup`s in our
examples.
