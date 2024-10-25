# Mockserver

A simple, high-performing mockserver that can dynamically build new routes with customized responses.

## Use

Call the `/register` endpoint to define a route.

### Curl

```sh
curl -X POST http://localhost:8080/register -d '{
  "method": "GET",
  "path": "/hello",
  "response": "{\"message\": \"Hello, world!\"}",
  "status_code": 200,
  "content_type": "application/json"
}' -H "Content-Type: application/json"
```

### Go and [Resty](https://github.com/go-resty/resty)

```go
client := resty.New()

route := map[string]interface{}{
    "method":      "GET",
    "path":        "/hello",
    "response":    "{\"message\":\"Hello, world!\"}",
    "status_code": 200,
    "content_type": "application/json",
}

resp, _ := client.R().
    SetHeader("Content-Type", "application/json").
    SetBody(route).
    Post("http://localhost:8080/register")
```

You can now call your endpoint and receive the JSON response back.

```sh
curl -X GET http://localhost:8080/hello -H "Content-Type: application/json"
# {"message":"Hello, world!"}
```

## Configure

Config is through environment variables.

| **Environment Variable** | **Description**                                                | **Default Value** |
| ------------------------ | -------------------------------------------------------------- | ----------------- |
| `LOG_LEVEL`              | Controls the logging level (`debug`, `info`, `warn`, `error`). | `debug`           |
| `SAVE_FILE`              | Path to the file where routes are saved and loaded.            | `save.json`       |

## Run

```sh
go run .
```

## Test

```sh
go test -cover -race ./...
```

## Benchmark

```sh
LOG_LEVEL=disabled go test -bench=. -benchmem -run=^$
```

Benchmark run on an Apple M3 Max.

```sh
goos: darwin
goarch: arm64
pkg: github.com/smartcontractkit/chainlink-testing-framework/mockserver
BenchmarkRegisterRoute-14    	  604978	      1967 ns/op	    6263 B/op	      29 allocs/op
BenchmarkRouteResponse-14    	16561670	        70.62 ns/op	      80 B/op	       1 allocs/op
BenchmarkSaveRoutes-14       	    1245	    956784 ns/op	  636042 B/op	    2014 allocs/op
BenchmarkLoadRoutes-14       	    1020	   1185990 ns/op	  348919 B/op	    9020 allocs/op
```

## Contribute

[Pre-commit](https://pre-commit.com/) is recommended to run checks before committing or pushing.

```sh
pre-commit install
```
