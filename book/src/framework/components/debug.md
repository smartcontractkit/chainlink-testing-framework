# Debugging Tests

All container logs are saved in a directory named `logs`, which will appear in the same directory where you ran the test after the test is completed.

For verifying your on-chain components, refer to the [Blockscout documentation](https://docs.blockscout.com/devs/verification/foundry-verification) on smart contract verification. This guide provides detailed instructions on uploading your ABI and verifying your contracts.

Use `CTF_LOG_LEVEL=trace|debug|info|warn` to debug the framework.

Use `RESTY_DEBUG=true` to debug any API calls.

Use `SETH_LOG_LEVEL=trace|debug|info|warn` to debug [Seth](../../libs/seth.md).

## Using Delve (TBD)

You can use [Delve]() inside your containers to debug aplications.

Build them with `go build -gcflags="all=-N -l" -o myapp` and use an example `Dockerfile`:
```
FROM golang:1.20

# Install Delve
RUN go install github.com/go-delve/delve/cmd/dlv@latest

# Set working directory
WORKDIR /app

# Copy the application binary and source code (if needed for debugging)
COPY myapp /app/myapp
COPY . /app

# Expose the port for Delve
EXPOSE 40000

# Start Delve in headless mode for remote debugging
ENTRYPOINT ["dlv", "exec", "./myapp", "--headless", "--listen=:40000", "--api-version=2", "--accept-multiclient"]

```

Adding `Delve` to all our components is WIP right now.

To expose `Delve` port follow this [guide](state.md).



