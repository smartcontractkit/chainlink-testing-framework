---
layout: default
title: Multiple Networks
nav_order: 2
parent: Networks
---

# Multiple Networks

You're likely wondering why `NETWORKS` takes a list of blockchain networks, and that's becuase we support running multiple
blockchain networks at once! This is handy for running scenarios where you want to test out
[chainlink's cross-chain capabilities](https://chain.link/cross-chain).

```go
suiteSetup, err = actions.MultiNetworkSetup( // Indicate you want to launch multiple networks
  environment.NewChainlinkCluster(3),        // Launch 3 Chainlink nodes
  actions.NetworksFromConfigHook,            // Default, uses all the networks defined in our config file
  actions.EthereumDeployerHook,              // Default, creates a contract deployer for each network
  actions.EthereumClientHook,                // Default, establishes client connection to each network
  tools.ProjectRoot,                         // Default, the path of our config file.
)
Expect(err).ShouldNot(HaveOccurred())

firstNetwork, err := suiteSetup.Network(0) // Retrieve the network at the 0 index of our network list
Expect(err).ShouldNot(HaveOccurred())
secondNetwork, err := suiteSetup.Network(1) // Retrieve the second network
Expect(err).ShouldNot(HaveOccurred())
defaultNetwork := suiteSetup.DefaultNetwork() // Get the default network: the one at the 0 index, same as firstNetwork
```

This gives you 2 seperate networks to work with, of the type listed in the `NETWORKS` ENV variable. If it's two of the
same simulated network, like `ethereum_geth,ethereum_geth`, then 2 simulated networks are launched and interacted with
seperately. Avoid using a test net combination like `ethereum_kovan,ethereum_kovan`, which will probably lead to unexpected
behaviors.
