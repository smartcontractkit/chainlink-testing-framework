# Quick Contracts Deployment

You can control the mining pace to accelerate contract deployment. Start anvil with the following configuration:

```toml
[blockchain_a]
  type = "anvil"
```
Set the `miner` speed,
```golang
	// start periodic mining so nodes can receive heads (async)
	miner := rpc.NewRemoteAnvilMiner(bcSrc.Nodes[0].ExternalHTTPUrl, nil)
	miner.MinePeriodically(5 * time.Second)
```

Then use this [template](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/framework/examples/myproject/quick_deploy_test.go)
