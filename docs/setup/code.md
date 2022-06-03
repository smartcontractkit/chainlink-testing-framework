---
layout: default
title: Test Setup Code
nav_order: 3
parent: Setup
---

# Test Setup Code

Now that we've got our config and Kubernetes sorted, we can write a bit of code that will deploy an environment for our test to run. To deploy our simulated geth, mock-server, and Chainlink instances, we rely on another chainlink library, [chainlink-env](https://github.com/smartcontractkit/chainlink-env/). This library handles deploying everything our test needs to the Kubernetes cluster.

```go
// We use the chainlink-env library to make and handle deployed resources
import "github.com/smartcontractkit/chainlink-env/environment"

// Deploy a testing environment, and receive it as the `env` variable. This is used to connect to resources.
e = environment.New(nil)
err := e.
    AddHelm(mockservercfg.New(nil)).
    AddHelm(mockserver.New(nil)).
    AddHelm(geth.New(nil)).
    AddHelm(chainlink.New(nil)).
    Run()
Expect(err).ShouldNot(HaveOccurred(), "Environment deployment shouldn't fail")
// Connect to all networks specified in the networks.yaml file
networkRegistry := client.NewDefaultNetworkRegistry()
// Retrieve these networks
networks, err := networkRegistry.GetNetworks(env)
// Get the default network (the first one in your listed selected_networks)
defaultNetwork := networks.Default
```

Most of the setup code will be the same for all your tests. Here's a more detailed explanation as to what some of the deployment code is doing to launch a few common test resources.

```go
e := environment.New(&environment.Config{
    Labels: []string{fmt.Sprintf("envType=%s", pkg.EnvTypeEVM5)}, // set more additional labels
})
err := e.
    AddHelm(mockservercfg.New(nil)). // add more Helm charts, all charts got merged in a manifest and deployed with kubectl when you call Run()
    AddHelm(mockserver.New(nil)).
    Run()
Expect(err).ShouldNot(HaveOccurred(), "Environment deployment shouldn't fail")
// do some other stuff with deployed charts if you need to interact with deployed services
err = e.
    AddChart(blockscout.New(&blockscout.Props{})). // you can also add cdk8s charts if you like Go code
    AddHelm(geth.New(nil)).
    AddHelm(chainlink.New(nil)).
    Run()
// Connect to all networks specified in the networks.yaml file
networkRegistry := client.NewDefaultNetworkRegistry()
// Retrieve these networks
networks, err := networkRegistry.GetNetworks(env)
// Get the default network (the first one in your listed selected_networks)
defaultNetwork := networks.Default
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
