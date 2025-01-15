# Parrot Server

A simple, high-performing mockserver that can dynamically build new routes with customized responses, parroting back whatever you tell it to.

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