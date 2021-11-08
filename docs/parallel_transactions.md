---
layout: default
title: Parallel Transactions
nav_order: 6
has_children: false
---

# Parallel Transactions

By default, when the framework makes a blockchain transaction
(deploying or interacting with a contract, funding nodes or contracts, etc...) the code blocks until the transaction is
confirmed, then moves on to the next line of code. There are some scenarios where this isn't the most efficient
(e.g. deploying multiple contracts at once, not caring which is deployed first) and the behavior can be changed.

```go
suiteSetup, err = actions.SingleNetworkSetup(
  environment.NewChainlinkCluster(3),
  actions.EVMNetworkFromConfigHook,   // Default, uses the first network defined in our config file
  actions.EthereumDeployerHook,       // Default, creates a contract deployer for our network
  actions.EthereumClientHook,         // Default, establishes client connection to our network
  tools.ProjectRoot,
)
Expect(err).ShouldNot(HaveOccurred())

networkInfo = suiteSetup.DefaultNetwork()
networkInfo.Client.ParallelTransactions(true) // Make transactions with this network run in parallel

// Deploy a flux aggregator contract
fluxInstance, err = networkInfo.Deployer.DeployFluxAggregatorContract(networkInfo.Wallets.Default(), contracts.DefaultFluxAggregatorOptions())
Expect(err).ShouldNot(HaveOccurred())
err = fluxInstance.Fund(networkInfo.Wallets.Default(), nil, big.NewFloat(1)) // Fund the contract
Expect(err).ShouldNot(HaveOccurred())
err = fluxInstance.UpdateAvailableFunds(context.Background(), networkInfo.Wallets.Default()) // Interact with the contract
Expect(err).ShouldNot(HaveOccurred())

err = networkInfo.Client.WaitForEvents() // Wait for all transactions to actually finish, then continue
Expect(err).ShouldNot(HaveOccurred())
```

This is a bit of a niche use case, and probably only advisable in simulated networks, but can be helpful in eeking out
some extra performance in your tests.