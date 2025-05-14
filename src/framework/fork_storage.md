# Fork Testing (Mutating Storage)

We provide API to use `anvil_setStorageAt` more easily so in case you can't edit EVM smart contracts code you can still mutate your contract values.

You need to build your contract layout first
```
forge build || forge inspect Counter storageLayout --json > layout.json
```

And then you can use `AnvilSetStorageAt` to override contract's storage data
```
		r := rpc.New(rpcURL, nil)
		err = r.AnvilSetStorageAt([]interface{}{contractAddr, slot, data})
```
See [examples](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/framework/evm_storage/layout_api_test.go) of how you can use API to encode/mutate different values.

Keep in mind that values <32 bytes are packed together, see `encodeCustomStructFunc` example and `offset` example to understnad how to change them properly.

```
cd framework/evm_storage
./setup.sh
go test -v -run TestLayoutAPI
./teardown.sh
```

Read more about Solidity storage layout [here](https://docs.soliditylang.org/en/latest/internals/layout_in_storage.html#)
