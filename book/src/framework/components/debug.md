# Debugging Tests

All container logs are saved in a directory named `logs`, which will appear in the same directory where you ran the test after the test is completed.

For verifying your on-chain components, refer to the [Blockscout documentation](https://docs.blockscout.com/devs/verification/foundry-verification) on smart contract verification. This guide provides detailed instructions on uploading your ABI and verifying your contracts.

Use `CTF_LOG_LEVEL=trace|debug|info|warn` to debug the framework.

Use `RESTY_DEBUG=true` to debug any API calls.

Use `SETH_LOG_LEVEL=trace|debug|info|warn` to debug [Seth](../../libs/seth.md).

## Using Delve Debugger

If you are using `Chainlink` image with [Delve](https://github.com/go-delve/delve) available in path you can use ports `40000..400XX` to connect to any node.
