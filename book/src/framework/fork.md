# Fork Testing

We verify our on-chain and off-chain changes using forks of various networks.

Go to example project [dir](https://github.com/smartcontractkit/chainlink-testing-framework/tree/main/framework/examples/myproject) to try the examples yourself.

## On-chain Only
In this [example](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/framework/examples/myproject/fork_test.go), we:

- Create two `anvil` networks, each targeting the desired network (change [URLs](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/framework/examples/myproject/fork.toml) and `anvil` settings as required, see [full anvil](https://book.getfoundry.sh/reference/anvil/) reference).
- Connect two clients to the respective networks.
- Deploy two test contracts.
- Interact with the deployed contracts.
- Demonstrate interactions using the `anvil` RPC client (more client methods examples are [here](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/framework/rpc/rpc_test.go))

Run it 
```
CTF_CONFIGS=fork.toml go test -v -run TestFork
```

## On-chain + Off-chain

The chain setup remains the same as in the previous example, but now we have 5 `Chainlink` nodes [connected with 2 networks](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/framework/examples/myproject/fork_plus_offchain_test.go).

Run it
```
CTF_CONFIGS=fork_plus_offchain.toml go test -v -run TestOffChainAndFork
```

<div class="warning">

Be mindful of RPC rate limits, as your provider may enforce restrictions. Use `docker_cmd_params` field to configure appropriate rate limiting and retries with the following parameters:
```
--compute-units-per-second <CUPS>
--fork-retry-backoff <BACKOFF>
--retries <retries>
--timeout <timeout>
```
If the network imposes limits, the container will panic, triggering messages indicating that the container health check has failed.

</div>

