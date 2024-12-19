# Kubernetes - enviromments

As already mentioned `CTFv1` creates `k8s` environments programmatically from existing building blocks that currently include:
* `anvil`
* `blockscout` (cdk8s)
* `chainlink node`
* `geth`
* `goc` (cdk8s)
* `grafana`
* `influxdb`
* `kafka`
* `mock-adapter`
* `mockserver`
* `reorg controller`
* `schema registry`
* `solana validator`
* `starknet validator`
* `wiremock`

Unless marked otherwise, they are all based on `Helm` charts.

> [!NOTE]
> Creation of new environment and modification of existing ones is explained in detail [here](../k8s/TUTORIAL.md), so we won't repeat it here,
> but instead focus on practical example of creating a new `k8s` test that creates its own basic environment.
>
> **It is highly recommended that you read it before continuing**.

We will focus on creating a basic testing environment compromised of:
* 6 `chainlink nodes`
* 1 blockchain node (`go-ethereum` aka `geth`)

Let's start!

# Step 1: Create Chainlink node TOML config
In real-world scenario you should dynamically create or load Chainlink node configuration to match your needs.
Here, for simplification, we will use a hardcoded config that will work for our case.
```go
func TestSimpleDONWithLinkContract(t *testing.T) {
	tomlConfig := `[Feature]
FeedsManager = true
LogPoller = true
UICSAKeys = true

[Database]
MaxIdleConns = 20
MaxOpenConns = 40
MigrateOnStartup = true

[Log]
Level = "debug"
JSONConsole = true

[Log.File]
MaxSize = "0b"

[WebServer]
AllowOrigins = "*"
HTTPWriteTimeout = "3m0s"
HTTPPort = 6688
SecureCookies = false
SessionTimeout = "999h0m0s"

[WebServer.RateLimit]
Authenticated = 2000
Unauthenticated = 1000

[WebServer.TLS]
HTTPSPort = 0

[OCR]
Enabled = true

[P2P]

[P2P.V2]
ListenAddresses = ["0.0.0.0:6690"]

[[EVM]]
ChainID = "1337"
AutoCreateKey = true
FinalityDepth = 1
MinContractPayment = "0"

[EVM.GasEstimator]
PriceMax = "200 gwei"
LimitDefault = 6000000
FeeCapDefault = "200 gwei"

[[EVM.Nodes]]
Name = "Simulated Geth-0"
WSURL = "ws://geth:8546"
HTTPURL = "http://geth:8544"`
```

This configuration uses log poller and OCRv2 and connects to a single EVM chain with id `1337` that can be reached through RPC node with following URLS:
* `ws://geth:8546`
* `http://geth:8544`

These are standard `geth` ports, and `geth` is the default name of our `go-ethereum` k8s service. We will connect to a "simulated" blockchain, which
is a private ephemeral/on-demand blockchain composed of a single node.

Now, let's build the chart that describes our chainlink `k8s` deployment:
```go
chainlinkImageCfg := &ctf_config.ChainlinkImageConfig{
    Image:   ptr.Ptr("public.ecr.aws/chainlink/chainlink"),
    Version: ptr.Ptr("2.19.0"),
}

var overrideFn = func(_ interface{}, target interface{}) {
    ctf_config.MustConfigOverrideChainlinkVersion(chainlinkImageCfg, target)
}

cd := chainlink.NewWithOverride(0, map[string]any{
    "replicas": 6,          // number of nodes
    "toml":     tomlConfig,
    "db": map[string]any{
        "stateful": true,   // stateful DB by default for soak tests
    },
}, chainlinkImageCfg, overrideFn)
```

Here, we use a hardcoded image and version for the Chainlink node, but in real test you would like to make it configurable. This setup
will launch 6 nodes, with stateful set dbs. It does look complex, but for various legacy reasons after removing support for some env vars
this is how setting image name and version looks like.

# Step 2: Label resources
For the purpose of better expenses tracking in the next step we will create necessary `chain.link` labels that every k8s resource needs to have. We will
use existing convenience functions:
```go
productName := "data-feedsv1.0"
nsLabels, err := environment.GetRequiredChainLinkNamespaceLabels(productName, "soak")
if err != nil {
    t.Fatal("Error creating required chain.link labels for namespace", err)
}

workloadPodLabels, err := environment.GetRequiredChainLinkWorkloadAndPodLabels(productName, "soak")
if err != nil {
    t.Fatal("Error creating required chain.link labels for workload and pod", err)
}
```

> [!NOTE]
> As explained [here](../k8s/labels.md) there are two environment variables that need to be set
> to satisfy labelling requirements:
> - `CHAINLINK_ENV_USER` - name of person running the test
> - `CHAINLINK_USER_TEAM` - name of the team, for which the test is run

# Step 3: Create environment config
This step is pretty straightforward:
```go
baseEnvironmentConfig := &environment.Config{
    TTL:                time.Hour * 2,
    NamespacePrefix:    "my-namespace-prefix",
    Test:               t,
    PreventPodEviction: true,
    Labels:             nsLabels,           // pass labels created in previous step
    WorkloadLabels:     workloadPodLabels,  // pass labels created in previous step
    PodLabels:          workloadPodLabels,  // pass labels created in previous step
}
```
Just three explanations are necessary here:
* `TTL` is the amount of time after which the namespace will by automatically removed
* `NamespacePrefix` is the preffix to which unique hash will be attached to ensure name uniqueness
* `PreventPodEviction` will prevent our pods from being evicted or restarted by `k8s`

# Step 4: Define blockchain network
For simplicity, we will use a hardcoded "simulated" EVM network, which should more accurately be called
an ephemeral private blockchain. In real case scenario you would use existing convenienice functions
for dynamically selecting the network, to which nodes should connect as it could be either a "simulated" one or an existing network (public or private).
In the latter case your code should skip adding the `ethereum` chart that represents `go-ethereum`-based blockchain node, as it
be connecting an already available service.

```go
nodeNetwork := blockchain.SimulatedEVMNetwork

ethProps := &ethereum.Props{
    NetworkName: nodeNetwork.Name,
    Simulated:   nodeNetwork.Simulated,
    WsURLs:      nodeNetwork.URLs,
    HttpURLs:    nodeNetwork.HTTPURLs,
}
```
There's no default network name or URLs set for `ethereum` chart, so you need to set these as a minimum.

# Step 5: Build the environment
Now that we have all the building blocks lets put them together and build the environment:
```go
testEnv := environment.New(baseEnvironmentConfig).
    AddHelm(ethereum.New(ethProps)).    // blockchain node
    AddHelm(cd)                         // chainlink node

err = testEnv.Run()
if err != nil {
    t.Fatal("Error running environment", err)
}
```

# Step 6: Create new blockchain client
With our environment created, let's create blockchain client, which will connect to our EVM node and later on deploy
a contract. We will use [Seth](../../libs/seth.md) for that purpose:
```go
// if test is running inside K8s, nothing to do, default network urls are correct
if !testEnv.Cfg.InsideK8s {
    // Test is running locally, use forwarded URL of Geth blockchain node
    wsURLs := testEnv.URLs[blockchain.SimulatedEVMNetwork.Name]
    httpURLs := testEnv.URLs[blockchain.SimulatedEVMNetwork.Name+"_http"]
    if len(wsURLs) == 0 || len(httpURLs) == 0 {
        t.Fatal("Forwarded Geth URLs should not be empty")
    }
    nodeNetwork.URLs = wsURLs
    nodeNetwork.HTTPURLs = httpURLs
}

sethClient, err := seth.NewClientBuilder().
    WithRpcUrl(nodeNetwork.URLs[0]).
    WithPrivateKeys([]string{nodeNetwork.PrivateKeys[0]}).
    Build()
if err != nil {
    t.Fatal("Error creating Seth client", err)
}
```
Notice the URL rewriting for our `nodeNetwork`. That's required, because by default, that network uses the name
of `geth` service in the `k8s` as it's URI. That works inside `k8s`, but not when your test is executing
on local environment, as is currently the case.

`Environment` is capable of forwarding `k8s` ports to local machine and does that for some of applications automatically.
`Geth` running in "simulated" mode is one of these and adds forwarded ports to the `URLs` map, so we can just grab them from it.

# Step 7: Deploy LINK contract
Finally, let's deploy a LINK contract and assert that it's total supply isn't 0:
```go
linkTokenAbi, err := link_token_interface.LinkTokenMetaData.GetAbi()
if err != nil {
    t.Fatal("Error getting LinkToken ABI", err)
}
linkDeploymentData, err := sethClient.DeployContract(sethClient.NewTXOpts(), "LinkToken", *linkTokenAbi, common.FromHex(link_token_interface.LinkTokenMetaData.Bin))
if err != nil {
    t.Fatal("Error deploying LinkToken contract", err)
}
linkToken, err := link_token_interface.NewLinkToken(linkDeploymentData.Address, sethClient.Client)
if err != nil {
    t.Fatal("Error creating LinkToken contract instance", err)
}

totalSupply, err := linkToken.TotalSupply(sethClient.NewCallOpts())
if err != nil {
    t.Fatal("Error getting total supply of LinkToken", err)
}
if totalSupply.Cmp(big.NewInt(0)) <= 0 {
    t.Fatal("Total supply of LinkToken should be greater than 0")
}
```
In a real world scenario that could be the end of the setup phase. Well, you should probably deploy a couple more contracts,
maybe the data feeds? And then, generate some load, ideally using [WASP](../../libs/wasp/overview.md).

Let's say that is what you really want. And that on top of that you would like your test to run for 2 days without having to
keep your local machine up and running, or having to deal with CI limitations (6h maximum action duration in Github Actions).

In the [next chapter](./remote_runner.md) you'll learn how to achieve that.

> [!NOTE]
> You can find this example [here](https://github.com/smartcontractkit/chainlink-testing-framework/tree/main/lib/k8s/examples/link/link_test.go).