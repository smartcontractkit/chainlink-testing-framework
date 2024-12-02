# Verifying Contracts

Check out our [example](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/framework/examples/myproject/verify_test.go) of programmatically verifying contracts using `Blockscout` and `Foundry`. You'll need to provide:

- The path to your Foundry directory
- The path to the contract
- The contract name

```golang
		err := blockchain.VerifyContract(blockchainComponentOutput, c.Addresses[0].String(),
			"example_components/onchain",
			"src/Counter.sol",
			"Counter",
		)
		require.NoError(t, err)
```