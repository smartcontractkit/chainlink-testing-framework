# Kubernetes - Environments

As mentioned earlier, `CTFv1` creates `k8s` environments programmatically from existing building blocks. These include:

- `anvil`
- `blockscout` (cdk8s)
- `chainlink node`
- `geth`
- `goc` (cdk8s)
- `grafana`
- `influxdb`
- `kafka`
- `mock-adapter`
- `mockserver`
- `reorg controller`
- `schema registry`
- `solana validator`
- `starknet validator`
- `wiremock`

Unless noted otherwise, all components are based on `Helm` charts.

> [!NOTE]
> The process of creating new environments or modifying existing ones is explained in detail [here](../k8s/TUTORIAL.md). This document focuses on a practical example of creating a new `k8s` test environment with a basic setup.
>
> **It is highly recommended to read that tutorial before proceeding.**

---

## Example: Basic Testing Environment

We will create a simple testing environment consisting of:
- 6 `Chainlink nodes`
- 1 blockchain node (`go-ethereum`, aka `geth`)

---

### Step 1: Create Chainlink Node TOML Config

In real-world scenarios, you should dynamically generate or load Chainlink node configurations to suit your needs. For simplicity, we will use a hardcoded configuration:

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

This configuration enables the log poller and OCRv2 features while connecting to an EVM chain with `ChainID` `1337`. It uses the following RPC URLs:
- WebSocket: `ws://geth:8546`
- HTTP: `http://geth:8544`

These URLs correspond to the default ports for `geth` and match the `go-ethereum` service name in the `k8s` cluster.

---

### Step 2: Define the Chainlink Deployment

To define the Chainlink deployment, we configure the image, version, and other parameters such as replicas and database settings. Here's the detailed implementation:

```go
chainlinkImageCfg := &ctf_config.ChainlinkImageConfig{
    Image:   ptr.Ptr("public.ecr.aws/chainlink/chainlink"),
    Version: ptr.Ptr("2.19.0"),
}

var overrideFn = func(_ interface{}, target interface{}) {
    ctf_config.MustConfigOverrideChainlinkVersion(chainlinkImageCfg, target)
}

cd := chainlink.NewWithOverride(0, map[string]any{
    "replicas": 6,          // Number of Chainlink nodes
    "toml":     tomlConfig, // TOML configuration defined earlier
    "db": map[string]any{
        "stateful": true,   // Use stateful databases for tests
    },
}, chainlinkImageCfg, overrideFn)
```

**Key Details:**
- **Image and Version:** These are hardcoded here for simplicity but should ideally be configurable for different environments.
- **Replicas:** We specify 6 Chainlink nodes to simulate a multi-node setup.
- **Database Configuration:** The database is stateful to allow for persistence during soak tests.
- **Override Function:** This ensures that the specified image and version are applied to all Chainlink node deployments.

---

### Step 3: Label Resources

To track costs effectively, add required `chain.link` labels to all `k8s` resources:

```go
productName := "data-feedsv1.0"
nsLabels, err := environment.GetRequiredChainLinkNamespaceLabels(productName, "soak")
if err != nil {
    t.Fatal("Error creating namespace labels", err)
}

workloadPodLabels, err := environment.GetRequiredChainLinkWorkloadAndPodLabels(productName, "soak")
if err != nil {
    t.Fatal("Error creating workload and pod labels", err)
}
```

Set the following environment variables:
- `CHAINLINK_ENV_USER`: Name of the person running the test.
- `CHAINLINK_USER_TEAM`: Name of the team the test is for.

---

### Step 4: Create Environment Config

```go
baseEnvironmentConfig := &environment.Config{
    TTL:                time.Hour * 2,
    NamespacePrefix:    "my-namespace-prefix",
    Test:               t,
    PreventPodEviction: true,
    Labels:             nsLabels,
    WorkloadLabels:     workloadPodLabels,
    PodLabels:          workloadPodLabels,
}
```

**Key Fields:**
- **`TTL`**: Time-to-live for the namespace (auto-removal after this time).
- **`NamespacePrefix`**: Ensures unique namespace names.
- **`PreventPodEviction`**: Prevents pods from being evicted or restarted.

---

### Step 5: Define Blockchain Network

To set up the blockchain network, we use predefined properties for a simulated EVM network. Here's the detailed implementation:

```go
nodeNetwork := blockchain.SimulatedEVMNetwork

ethProps := &ethereum.Props{
    NetworkName: nodeNetwork.Name,         // Name of the network
    Simulated:   nodeNetwork.Simulated,   // Indicates that the network is simulated
    WsURLs:      nodeNetwork.URLs,        // WebSocket URLs for the network
    HttpURLs:    nodeNetwork.HTTPURLs,    // HTTP URLs for the network
}
```

**Details:**
- **Simulated Network:** Represents a private, ephemeral blockchain used for testing.
- **Dynamic Selection:** In real scenarios, use helper functions to dynamically select networks (public, private, or simulated) based on test requirements.
- **Custom URLs:** The `ethereum` chart requires explicit settings for the network name and URLs.

---

### Step 6: Build the Environment

```go
testEnv := environment.New(baseEnvironmentConfig).
    AddHelm(ethereum.New(ethProps)).    // Blockchain node
    AddHelm(cd)                         // Chainlink nodes

err = testEnv.Run()
if err != nil {
    t.Fatal("Error running environment", err)
}
```

---

### Step 7: Create Blockchain Client

```go
if !testEnv.Cfg.InsideK8s {
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

**Details:**
- **Local vs. Cluster Environment**: When running tests outside the k8s cluster, the service URLs (`ws://geth:8546`, `http://geth:8544`) are not directly accessible. Port forwarding ensures local access to these services.
- **Automatic Port Forwarding**: The `Environment` object manages forwarding for key services, including Geth in simulated mode, making these forwarded URLs available in the `URLs` map.
- **Dynamic Rewriting**: URLs are dynamically rewritten to switch between in-cluster and local connectivity.

---

### Step 8: Deploy LINK Contract

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

**Details:**
- **Deploy Contract:** Deploys the LINK token contract to the simulated blockchain.
- **Verify Deployment:** Ensures the total supply is greater than zero as a sanity check.

---

### Next Steps

Learn how to run long-duration tests using a `remote runner` in the [next chapter](./remote_runner.md).

> [!NOTE]
> This example can be found [here](https://github.com/smartcontractkit/chainlink-testing-framework/tree/main/lib/k8s/examples/link/link_test.go).