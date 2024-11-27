# WASP - Getting Started

## Prerequisites
* [Golang](https://go.dev/doc/install)

## Setup
To start writing tests, create a directory for your project with a `go.mod` file, then add the WASP package:

```bash
go get github.com/smartcontractkit/chainlink-testing-framework/wasp
```

> [!WARNING]  
> To execute any of the tests from the next chapters, you need access to Loki and Grafana.  
> You can find instructions on setting them up locally [here](./how-to/start_local_observability_stack.md).

That was simple, wasn't it? Time to write your [first test](./first_test.md).