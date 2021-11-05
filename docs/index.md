---
layout: default
title: Home
nav_order: 1
description: "A general blockchain integration testing framework geared towards Chainlink projects"
permalink: /
---

# Chainlink Integrations Framework

[![Go Report Card](https://goreportcard.com/badge/github.com/smartcontractkit/integrations-framework)](https://goreportcard.com/report/github.com/smartcontractkit/integrations-framework)
[![Go Reference](https://pkg.go.dev/badge/github.com/smartcontractkit/integrations-framework.svg)](https://pkg.go.dev/github.com/smartcontractkit/integrations-framework)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

The Chainlnk Integration Framework is a blockchain development and testing framework written in Go. Its primary purpose
is to help chainlink developers create extensive integration, e2e, performance, and chaos tests to ensure the stability
of the chainlink project. It can also be helpful to those who just want to use chainlink oracles in their projects to help
test their contracts, or even for those that aren't using chainlink.

## Setup

In order to use this framework, you must have a connection to an actively running
[Kubernetes cluster](https://kubernetes.io/docs/tasks/access-application-cluster/access-cluster/). If you don't have
one handy, check out [minikube](https://minikube.sigs.k8s.io/docs/start/) which should work fine for smaller tests,
but will likely need to be allocated more power, or you'll need to use a more powerful cluster in general to run tests
that require lots of services, like OCR.

### Why?

There's a lot of different components to bring up for each test, most of which involve:

* A simulated blockchain
* Some number of Chainlink nodes
* An equal number of postgres DBs to support the chainlink nodes
* At least one external adapter

Following the good testing practice of having clean, non-dependent test environments means we're creating a lot of these
components for each test, and tearing them down soon after. In order to organize these test environments, and after
finding `docker compose` to be woefully inadequate after a certain point, kubernetes was the obvious choice.

## Usage

We recommend using [Ginkgo](https://github.com/onsi/ginkgo), a BDD testing framework for Go that we've found handy
for organizing and running tests. You should be able to use any other testing framework you like, including Go's built-in
`testing` package, but the examples you find here will be in Ginkgo and its accompanying assertions library,
[gomega](https://onsi.github.io/gomega/)

```go
var _ = Describe("Basic Contract Interactions", func() {
  var ( // Create variables that we're going to be using across test steps
    suiteSetup    actions.SuiteSetup
    networkInfo   actions.NetworkInfo
    defaultWallet client.BlockchainWallet
  )

  It("Exercises basic smart contract usage", func() {
    By("Deploying the environment", func() {
      var err error
      // SuiteSetup creates an ephemeral environment for the test, launching a simulated blockchain, an external adapter
      // and as many chainlink nodes as you would like.
      suiteSetup, err = actions.SingleNetworkSetup( 
        environment.NewChainlinkCluster(0), // We're launching this test with 0 chainlnk nodes
        client.DefaultNetworkFromConfig,    // Using the first network defined in our config file
        tools.ProjectRoot,                  // The path of our config file.
      )
      Expect(err).ShouldNot(HaveOccurred())
      networkInfo = suiteSetup.DefaultNetwork()
      defaultWallet = networkInfo.Wallets.Default()
    })

    By("Deploying and using the storage contract", func() {
      // Deploy a storage contract, all it does is store a value, then regurgitate that value when called for
      storeInstance, err := suiteSetup.Deployer.DeployStorageContract(defaultWallet)
      Expect(err).ShouldNot(HaveOccurred())

      // Value we're going to store
      testVal := big.NewInt(5)

      // Set the contract value
      err = storeInstance.Set(testVal)
      Expect(err).ShouldNot(HaveOccurred())
      // Retrieve the value
      val, err := storeInstance.Get(context.Background())
      // Make sure no errors happened, and the value is what we expect
      Expect(err).ShouldNot(HaveOccurred())
      Expect(val).To(Equal(testVal))
     })
  })

  AfterEach(func() {
    // Tears down the environment, deleting everything that the SuiteSetup launched, and collecting logs if the test failed
    By("Tearing down the environment", suiteSetup.TearDown())
  })
})
```

The `SuiteSetup` does a lot of work for us, allowing us to specify how many chainlink nodes we want to deploy,
what network we want to run on, and where to look for our config file. There will be quite a few calls creating
`SuiteSetup`s in our examples.

## Environments

By default, the `TearDown()` method deletes the environment that was launched by the `SuiteSetup`. Sometimes that's not
desired though, like when debugging failing tests. For that, there's a handy ENV variable, `KEEP_ENVIRONMENTS`.

```sh
KEEP_ENVIRONMENTS = Never # Options: Always, OnFail, Never
```

## Logs

`TearDown()` also checks if the test has failed. If so, it builds a `logs/` directory, and dumps the logs and contents
of each piece of the environment that was launched.