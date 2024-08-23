## Running multi-key load test with Seth and WASP

To effectively simulate transaction workloads from multiple keys, you can utilize a "rotating wallet." Refer to the [example](client_wasp_test.go) code provided for guidance.

There are 2 modes: Ephemeral and a static private keys mode.

### Ephemeral mode

We generate 60 ephemeral keys and run the test, set `ephemeral_addresses_number` in `seth.toml`

This mode **should never be used on testnets or mainnets** in order not to lose funds. Please use it to test with simulated networks, like private `Geth` or `Anvil`

```toml
ephemeral_addresses_number = 60
```

Then start the Geth and run the test

```
nix develop
make GethSync

// another terminal, from examples_wasp dir
export SETH_LOG_LEVEL=debug
export SETH_CONFIG_PATH=seth.toml
export SETH_NETWORK=Geth
export SETH_ROOT_PRIVATE_KEY=ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
export LOKI_TENANT_ID=promtail
export LOKI_URL=...

go test -v -run TestWithWasp
```

See both [generator](client_wasp_test.go) and [test](client_wasp_test.go) implementation example

Check your results [here](https://grafana.ops.prod.cldev.sh/d/WaspDebug/waspdebug?orgId=1&from=now-5m&to=now)

If you see `key sync timeout`, just increase `ephemeral_addresses_number` to have more load

You can also change default `key_sync` values

```toml
[nonce_manager]
# 20 req/s limit for key syncing
key_sync_rate_limit_per_sec = 20
# key synchronization timeout, if it's more than N sec you'll see an error, raise amount of keys or increase the timeout
key_sync_timeout = "30s"
# key sync retry delay, each N seconds we'll updage each key nonce
key_sync_retry_delay = "1s"
# total number of retries until we throw an error
key_sync_retries = 30
```

### Static private keys mode

In that mode you should pass static keys that you already have as part of `Network` configuration. It's strongly recommended to do that programmatically, not via config file, since accidentally committing private keys to the repository will compromise the funds.
It would be better to read the TOML configuration first:

```go
cfg, err := seth.ReadCfg()
if err != nil {
    log.Fatal(err)
}
```

Then read the private keys in a safe manner. For example from a secure vault or environment variables:

```go
var privateKeys []string
var err error
privateKeys, err = some_utils.ReadPrivateKeysFromEnv()
if err != nil {
log.Fatal(err)
}
```

and then add them to the `Network` you plan to use. Let's assume it's called `Sepolia`:

```go
for i, network := range cfg.Networks {
    if network.Name == "Sepolia" {
        cfg.Networks[i].PrivateKeys = privateKeys
    }
}
```

Or if you aren't using `[[Networks]]` in your TOML config and have just a single `Network`:

```go
cfg.Network.PrivateKeys = privateKeys
```

Or... you can use the convenience function `AppendPksToNetwork()` to have them added to both the `Network` and `Networks` slice:

```go
added := cfg.AppendPksToNetwork(privateKeys, "Sepolia")
if !added {
    log.Fatal("Network Sepolia not found in the config")
}
```

Finally, proceed to create a new Seth instance:

```go
seth, err := seth.NewClientWithConfig(cfg)
if err != nil {
    log.Fatal(err)
}
```
