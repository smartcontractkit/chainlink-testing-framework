# 🚀 Getting started

## Prerequisites
- `Docker` [OrbStack](https://orbstack.dev/) or [Docker Desktop](https://www.docker.com/products/docker-desktop/), we recommend OrbStack (faster, smaller memory footprint)
- [Golang](https://go.dev/doc/install)

## Test setup

To start writing tests create a directory for your project with `go.mod` and pull the framework
```
go get github.com/smartcontractkit/chainlink-testing-framework/framework
```

Then download the CLI (runs from the directory where you have `go.mod`)
```
go get github.com/smartcontractkit/chainlink-testing-framework/framework/cmd && \
go install github.com/smartcontractkit/chainlink-testing-framework/framework/cmd && \
mv ~/go/bin/cmd ~/go/bin/ctf
```
Or download a binary release [here](https://github.com/smartcontractkit/chainlink-testing-framework/releases/tag/framework%2Fv0.1.7) and rename it to `ctf`

More CLI [docs](./cli.md)

Create an `.envrc` file and do `source .envrc`
```
export TESTCONTAINERS_RYUK_DISABLED=true # do not remove containers while we develop locally
```

Now you are ready to write your [first test](./first_test.md)

## Tools setup (Optional)

This setup is optional, and it explains how to create a local observability stack for on-chain and off-chain components.

Spin up your local obserability stack (Grafana LGTM)
```
ctf obs up
```
More [docs](observability/observability_stack.md)

Spin up your `Blockscout` stack
```
ctf bs up
```
More [docs](observability/blockscout.md)
