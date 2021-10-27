# Chainlink Integration Framework

[![Go Report Card](https://goreportcard.com/badge/github.com/smartcontractkit/integrations-framework)](https://goreportcard.com/report/github.com/smartcontractkit/integrations-framework)
[![Go Reference](https://pkg.go.dev/badge/github.com/smartcontractkit/integrations-framework.svg)](https://pkg.go.dev/github.com/smartcontractkit/integrations-framework)

The Chainlnk Integration Framework is a blockchain development framework written in Go. Its primary purpose is to help
chainlink developers create extensive integration, e2e, performance, and chaos tests to ensure the stability of the
chainlink project. It can also be helpful to those who just want to use chainlink oracles in their projects to help
test their contracts, or even for those that aren't using chainlink.

## WIP

As of now, this framework is still very much a work in progress, and will have frequent changes, many of which
will probably be breaking.

{:toc}

## Setup

In order to use this framework, you must have a connection to an actively running 

## Usage

A simple example on deploying and interaction with a simple contract using this framework and [Ginkgo](https://github.com/onsi/ginkgo)

```go
var _ = Describe("Basic Contract Interactions", func() {
  var (
    suiteSetup    actions.SuiteSetup
    networkInfo   actions.NetworkInfo
    defaultWallet client.BlockchainWallet
  )

  BeforeEach(func() {
    By("Deploying the environment", func() {
      var err error
      suiteSetup, err = actions.SingleNetworkSetup(
        environment.NewChainlinkCluster(0),
        client.DefaultNetworkFromConfig,
        tools.ProjectRoot,
      )
      Expect(err).ShouldNot(HaveOccurred())
      networkInfo = suiteSetup.DefaultNetwork()
      defaultWallet = networkInfo.Wallets.Default()
    })
  })

  It("exercises basic contract usage", func() {
    By("deploying the storage contract", func() {
      // Deploy storage
      storeInstance, err := suiteSetup.Deployer.DeployStorageContract(defaultWallet)
      Expect(err).ShouldNot(HaveOccurred())

      testVal := big.NewInt(5)

      // Interact with contract
      err = storeInstance.Set(testVal)
      Expect(err).ShouldNot(HaveOccurred())
      val, err := storeInstance.Get(context.Background())
      Expect(err).ShouldNot(HaveOccurred())
      Expect(val).To(Equal(testVal))
     })
  })

  AfterEach(func() {
    By("Tearing down the environment", suiteSetup.TearDown())
  })
})
```

## Execution Environment

Ephemeral environments are automatically deployed with Kubernetes. To run tests, you either need a deployed cluster
in an environment, or a local installation.

### Locally

When running tests locally, it's advised to use minikube. To spin up a cluster, use:

```sh
minikube start
```

### Remotely

To run against a remote Kubernetes cluster, ensure your current context is the cluster you want to run against as the
framework always uses current context.

## Test Execution

This framework advises the use of [Ginkgo](https://github.com/onsi/ginkgo) for test execution, but tests still can be
ran with the go CLI.

### Ginkgo

Run:

```sh
ginkgo -r
```

### Go

Run:

```sh
go test ./..
```

### Volume tests

```sh
NETWORK="ethereum_geth_performance" make test_performance
```

### Build contracts

Example of generating go bindings for Ethereum contracts

```sh
ifcli build_contracts -c config.yml
```

### Create environment

Example of creating environment with one Chainlink node and Geth dev network

```sh
ifcli create_env -n ethereum_geth -t chainlink -c 1
```
