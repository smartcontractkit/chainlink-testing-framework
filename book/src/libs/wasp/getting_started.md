# WASP - Getting Started

## Pre-requisites
* [Golang](https://go.dev/doc/install)

## Setup
To start writing tests create a directory for your project with `go.mod` and add WASP package
```bash
go get github.com/smartcontractkit/chainlink-testing-framework/wasp
```

> [!WARNING]
> To execute any of the following examples you will also need access to Loki and Grafana. You can find instructions on how to set them up locally [here](./local_loki_grafana_stack.md).

That was simple, wasn't it? Time to write your [first test](./first_test.md)