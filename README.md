# Chainlink Integration Framework

[![Go Report Card](https://goreportcard.com/badge/github.com/smartcontractkit/integrations-framework)](https://goreportcard.com/report/github.com/smartcontractkit/integrations-framework)
[![Go Reference](https://pkg.go.dev/badge/github.com/smartcontractkit/integrations-framework.svg)](https://pkg.go.dev/github.com/smartcontractkit/integrations-framework)
![Tests](https://github.com/smartcontractkit/integrations-framework/actions/workflows/test.yaml/badge.svg)
![Lint](https://github.com/smartcontractkit/integrations-framework/actions/workflows/lint.yaml/badge.svg)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

The Chainlnk Integration Framework is a blockchain development framework written in Go. Its primary purpose is to help
chainlink developers create extensive integration, e2e, performance, and chaos tests to ensure the stability of the
chainlink project. It can also be helpful to those who just want to use chainlink oracles in their projects to help
test their contracts, or even for those that aren't using chainlink.

See the [docs](https://smartcontractkit.github.io/integrations-framework/) or our
[go reference](https://pkg.go.dev/github.com/smartcontractkit/integrations-framework) page for more detailed info and
examples. If you just want a quick overview, keep reading.

## WIP

As of now, this framework is still very much a work in progress, and will have frequent changes, many of which
will probably be breaking.

## Setup

In order to use this framework, you must have a connection to an actively running Kubernetes cluster. If you don't have
one handy, check out [minikube](https://minikube.sigs.k8s.io/docs/start/) which should work fine for smaller tests,
but will likely need to be allocated more power, or you'll need to use a more powerful cluster in general to run tests
that require lots of services, like OCR.

## Usage

Here's a simple example on deploying and interacting with a basic storage contract using this framework and
[Ginkgo](https://github.com/onsi/ginkgo), a BDD testing framework we've come to really enjoy. You can use another testing
framework, including Go's default testing if you prefer otherwise.

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

## Config Values

You'll notice in the `SuiteSetup` that we provide a path to the config file, `tools.ProjectRoot`. This links to our default
config file, `config.yml`. For most cases, this will work out just fine, as you can pass in ENV variables to override those
config values. Below are some common ones we find ourselves using regularly.

| ENV Var                 | Description                                                 | Default                            |
|-------------------------|-------------------------------------------------------------|------------------------------------|
|`NETWORKS`               | Comma seperated list of blockchain networks to run tests on | ethereum_geth,ethereum_geth        |
|`APPS_CHAINLINK_IMAGE`   | Image location for a valid docker image of a chainlink node | public.ecr.aws/chainlink/chainlink |
|`APPS_CHAINLINK_VERSION` | Version to be used for the above mentioned image            | 0.10.14                            |
|`NETWORK_CONFIGS_<NETWORK_NAME>_PRIVATE_KEYS` | Comma seperated list of private keys for the network to use | Varies        |

If you want to provide your own config file instead, you can point `SuiteSetup` to the directory that the config file lives in.

```go
suiteSetup, err = actions.SingleNetworkSetup( 
  environment.NewChainlinkCluster(0), 
  client.DefaultNetworkFromConfig,    
  "../", // Look for a config.yml file in the parent directory of this test file.                 
)
```
