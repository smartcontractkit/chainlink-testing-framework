---
layout: default
title: Chainlink Nodes
nav_order: 5
has_children: false
---

# Chainlink Nodes

An important portion of our test environments are chainlink nodes, getting them to connect to our contracts and read
their states. You can launch a chosen amount of chainlink nodes when calling for a `SuiteSetup`. You can modify the
chainlink node image and version by setting `APPS_CHAINLINK_IMAGE` and `APPS_CHAINLINK_VERSION`.

```go
// Launch 3 chainlink nodes of the image supplied in your config
environment.NewChainlinkCluster(3) 
// Launch 6 chainlink nodes, mixing in your supplied version with the 2 latest released versions
environment.NewMixedVersionChainlinkCluster(6, 2) 
```

Once launched, you can interact with the nodes, funding them and creating jobs for the nodes to run. A good example of
the typical workflow here can be seen in [our flux aggregator test](../suite/smoke/contracts_flux_test.go). Below
is a small portion of it.

```go
var (
  suiteSetup    actions.SuiteSetup
  networkInfo   actions.NetworkInfo
  adapter       environment.ExternalAdapter
  nodes         []client.Chainlink
  nodeAddresses []common.Address
  fluxInstance  contracts.FluxAggregator
  err           error
)

By("Deploying the environment", func() {
  suiteSetup, err = actions.SingleNetworkSetup(
    environment.NewMixedVersionChainlinkCluster(6, 2), // Use a mix of different chainlink versions
    actions.EVMNetworkFromConfigHook,                  // Default, uses the first network defined in our config file
    actions.EthereumDeployerHook,                      // Default, creates a contract deployer for our network
    actions.EthereumClientHook,                        // Default, establishes client connection to our network
    tools.ProjectRoot,
  )
  Expect(err).ShouldNot(HaveOccurred())
  nodes, err = environment.GetChainlinkClients(suiteSetup.Environment()) // Get all our chainlink nodes
  Expect(err).ShouldNot(HaveOccurred())
  adapter, err = environment.GetExternalAdapter(suiteSetup.Environment())
  Expect(err).ShouldNot(HaveOccurred())
  networkInfo = suiteSetup.DefaultNetwork()

  networkInfo.Client.ParallelTransactions(true)
})

By("Deploying and funding contract", func() {
  fluxInstance, err = networkInfo.Deployer.DeployFluxAggregatorContract(networkInfo.Wallets.Default(), contracts.DefaultFluxAggregatorOptions())
  Expect(err).ShouldNot(HaveOccurred())
  err = fluxInstance.Fund(networkInfo.Wallets.Default(), nil, big.NewFloat(1))
  Expect(err).ShouldNot(HaveOccurred())
  err = fluxInstance.UpdateAvailableFunds(context.Background(), networkInfo.Wallets.Default())
  Expect(err).ShouldNot(HaveOccurred())
  err = networkInfo.Client.WaitForEvents()
  Expect(err).ShouldNot(HaveOccurred())
})

By("Funding Chainlink nodes", func() {
  nodeAddresses, err = actions.ChainlinkNodeAddresses(nodes)
  Expect(err).ShouldNot(HaveOccurred())
  ethAmount, err := networkInfo.Deployer.CalculateETHForTXs(networkInfo.Wallets.Default(), networkInfo.Network.Config(), 3)
  Expect(err).ShouldNot(HaveOccurred())
  err = actions.FundChainlinkNodes( // Fund all our chainlink nodes
    nodes,
    networkInfo.Client,
    networkInfo.Wallets.Default(),
    ethAmount,
    nil,
  )
  Expect(err).ShouldNot(HaveOccurred())
})
```
