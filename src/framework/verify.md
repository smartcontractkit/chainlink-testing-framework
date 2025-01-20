# Verifying Contracts

You need to install [Foundry](https://book.getfoundry.sh/getting-started/installation) first, `forge` should be available in your `$PATH`.

Check out our [example](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/framework/examples/myproject/verify_test.go) of programmatically verifying contracts using `Blockscout` and `Foundry`. You'll need to provide:

- The path to your Foundry directory
- The path to the contract
- The contract name

```golang
		err := blockchain.VerifyContract(blockchainComponentOutput, c.Addresses[0].String(),
			"example_components/onchain",
			"src/Counter.sol",
			"Counter",
            "0.8.20" // check your file compiler version on the first line
		)
		require.NoError(t, err)
```
