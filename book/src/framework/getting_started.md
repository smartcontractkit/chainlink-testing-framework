# Getting started

## Prerequisites
- `Docker` ([OrbStack](https://orbstack.dev/) or [Docker Desktop](https://www.docker.com/products/docker-desktop/))
- [Golang](https://go.dev/doc/install)

## Test setup

To start writing tests create a directory for your project with `go.mod` and pull the framework
```
go get github.com/smartcontractkit/chainlink-testing-framework/framework
```

Then download the CLI (runs from directory where you have `go.mod`)
```
go get github.com/smartcontractkit/chainlink-testing-framework/framework/cmd && \
go install github.com/smartcontractkit/chainlink-testing-framework/framework/cmd && \
mv ~/go/bin/cmd ~/go/bin/ctf
```
More CLI [docs](./cli.md)

Create an `.envrc` file and put common parameters there (you can use [direnv](https://direnv.net/) to sync them more easily)
```
export TESTCONTAINERS_RYUK_DISABLED=true # do not remove containers while we develop locally
export CTF_CONFIGS=smoke.toml            # our configuration file
```

Now you are ready to write your [first test](./first_test.md)

## Tools setup (Optional)

This setup is optional, and it explains how to create a local observability stack for on-chain and off-chain components.

Spin up your local obserability stack (Grafana LGTM)
```
ctf obs up
```

Spin up your `Blockscout` stack
```
ctf bs up
```
