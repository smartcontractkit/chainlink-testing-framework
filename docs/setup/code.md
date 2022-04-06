---
layout: default
title: Test Setup Code
nav_order: 3
parent: Setup
---

# Test Setup Code

Now that we've got our config and Kubernetes sorted, we can write a bit of code that will deploy an environment for our test to run. To deploy our simulated geth, mock-server, and Chainlink instances, we rely on another chainlink library, [helmenv](https://github.com/smartcontractkit/helmenv/). This library handles deploying everything our test needs to the Kubernetes cluster.

```go
// We use the helmenv library to make and handle deployed resources
import "github.com/smartcontractkit/helmenv/environment"

// Deploy a testing environment, and receive it as the `env` variable. This is used to connect to resources.
env, err := environment.DeployOrLoadEnvironment( 
  // Define what sort of environment you would like to deploy. More on this below
  environment.NewChainlinkConfig(
    environment.ChainlinkReplicas(3, config.ChainlinkVals()),
    "namespace-prefix",
    config.GethNetworks()...,
  ),
  // Path to the helm charts you want to use (tools.ChartsRoot will work fine for 99% of cases)
  tools.ChartsRoot,
)
// Omitting error checking for brevity

// Connect to all the deployed resources to use later in the test
err = env.ConnectAll()
// Connect to all networks specified in the networks.yaml file
networkRegistry := client.NewDefaultNetworkRegistry()
// Retrieve these networks
networks, err := networkRegistry.GetNetworks(env)
// Get the default network (the first one in your listed selected_networks)
defaultNetwork := networks.Default
```

Most of the setup code will be the same for all your tests. Here's a more detailed explanation as to what some of the deployment code is doing to launch a few common test resources.

```go
env, err := environment.DeployOrLoadEnvironment( // Use helmenv to deploy a new environment
  environment.NewChainlinkConfig( // Indicate you want a standard EVM Chainlink testing environment
    environment.ChainlinkReplicas(3, config.ChainlinkVals()), // How many Chainlink nodes to launch, and what values to provide them
    "namespace-prefix", // The prefix of the namespace that will be created in Kubernetes
    config.GethNetworks()..., // All the settings of the simulated Geth networks that will be launched
  ),
  tools.ChartsRoot, // Default 
)
```

These common resources consist of

* A simulated Geth instance
* A basic mock server that serves as a mock adapter for Chainlink nodes
* A specified number of chainlink nodes

## Test Tear Down

When your test is done, you'll want to have a way to tear down the test environment you launched. You'll also want to be able to see the logs from your test environment. Below is a typical test flow.

```go
// Launch our environment
env, err := environment.DeployOrLoadEnvironment( 
  environment.NewChainlinkConfig(environment.ChainlinkReplicas(1, nil), "chainlink-test-setup"),
  tools.ChartsRoot,
)

// Put test logic here

// Tear down the test environment
// Prints some handy stats on Gas usage for the test, if you'd like to see that info.
networks.Default.GasStats().PrintStats()
// Tears down the test environment, according to options you selected in the `framework.yaml` config file
err = actions.TeardownSuite(
  env,                // The test environment object
  networks,           // The list of networks obtained from `networks, err := networkRegistry.GetNetworks(env)`
  utils.ProjectRoot,  // The folder location you'd like logs (on test failure) to be dumped to
  nil,                // An optional test reporter for more custom test statistics (we'll get to that later)
)
```
