# Parrot Server

A simple, high-performing mockserver that can dynamically build new routes with customized responses, parroting back whatever you tell it to.

## Features

* Simplistic and fast design
* Run within your Go code, through a small binary, or in a minimal Docker container
* Easily record all incoming requests to the server to programmatically react to 

## Use

See our runnable examples in [examples_test.go](./examples_test.go) to see how to use Parrot programmatically.

## Run

```sh
go run ./cmd
go run ./cmd -h # See all config options 
```

## Test

```sh
make test
make test PARROT_TEST_LOG_LEVEL=trace # Set log level for tests
make test_race # Test with -race flag enabled
make bench # Benchmark
```

## Build

```sh
make goreleaser # Uses goreleaser to build binaries and docker containers
```
