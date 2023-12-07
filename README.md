<div align="center">

# Chainlink Testing Framework

[![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/smartcontractkit/chainlink-testing-framework)](https://github.com/smartcontractkit/chainlink-testing-framework/tags)
[![Go Report Card](https://goreportcard.com/badge/github.com/smartcontractkit/chainlink-testing-framework)](https://goreportcard.com/report/github.com/smartcontractkit/chainlink-testing-framework)
[![Go Reference](https://pkg.go.dev/badge/github.com/smartcontractkit/chainlink-testing-framework.svg)](https://pkg.go.dev/github.com/smartcontractkit/chainlink-testing-framework)
[![Go Version](https://img.shields.io/github/go-mod/go-version/smartcontractkit/chainlink-testing-framework)](https://go.dev/)
![Tests](https://github.com/smartcontractkit/chainlink-testing-framework/actions/workflows/test.yaml/badge.svg)
![Lint](https://github.com/smartcontractkit/chainlink-testing-framework/actions/workflows/lint.yaml/badge.svg)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

</div>

The Chainlink Testing Framework is a blockchain development framework written in Go. Its primary purpose is to help chainlink developers create extensive integration, e2e, performance, and chaos tests to ensure the stability of the chainlink project. It can also be helpful to those who just want to use chainlink oracles in their projects to help test their contracts, or even for those that aren't using chainlink.

If you're looking to implement a new chain integration for the testing framework, head over to the [blockchain](./blockchain/) directory for more info.

## k8s package
We have a k8s package we are using in tests, it provides:
- [cdk8s](https://cdk8s.io/) based wrappers
- High-level k8s API
- Automatic port forwarding

You can also use this package to spin up standalone environments.

### Local k8s cluster
Read [here](./k8s/KUBERNETES.md) about how to spin up a local cluster

#### Install
Set up deps, you need to have `node 14.x.x`, [helm](https://helm.sh/docs/intro/install/) and [yarn](https://classic.yarnpkg.com/lang/en/docs/install/#mac-stable)

Then use
```shell
make install_deps
```

### Running tests in k8s
To read how to run a test in k8s, read [here](./k8s/REMOTE_RUN.md)

### Usage
Create an env in a separate file and run it
```
export CHAINLINK_IMAGE="public.ecr.aws/chainlink/chainlink"
export CHAINLINK_TAG="1.4.0-root"
export CHAINLINK_ENV_USER="Satoshi"
go run k8s/examples/simple/env.go
```
For more features follow [tutorial](./k8s/TUTORIAL.md)

### Development
#### Running standalone example environment
```shell
go run k8s/examples/simple/env.go
```
If you have another env of that type, you can connect by overriding environment name
```
ENV_NAMESPACE="..."  go run k8s/examples/chainlink/env.go
```

Add more presets [here](./k8s/presets)

Add more programmatic examples [here](./k8s/examples/)

If you have [chaosmesh]() installed in your cluster you can pull and generated CRD in go like that
```
make chaosmesh
```

If you need to check your system tests coverage, use [that](./k8s/TUTORIAL.md#coverage)

# Chainlink Charts

This repository contains helm charts used by the chainlink organization mostly in QA.

## Chart Repository

You can add the published chart repository by pointing helm to the `gh-pages` branch with a personal access token (PAT) that has at least read-only access to the repository.

```sh
helm repo add chainlink-qa https://raw.githubusercontent.com/smartcontractkit/qa-charts/gh-pages/
helm search repo chainlink
```

## Releasing Charts

The following cases will trigger a chart release once a PR is merged into the `main` branch.
Modified packages or new packages get added and pushed to the `gh-pages` branch of the [qa-charts](https://github.com/smartcontractkit/qa-charts) repository.

- An existing chart is version bumped
- A new chart is added

Removed charts do not trigger a re-publish, the packages have to be removed and the index file regenerated in the `gh-pages` branch of the [qa-charts](https://github.com/smartcontractkit/qa-charts) repository.

Note: The qa-charts repository is scheduled to look for changes to the charts once every hour. This can be expedited by going to that repo and running the cd action via github UI.

# Simulated EVM chains

We have extended support for execution layer clients in simulated networks. Following ones are supported:
* `Geth`
* `Nethermind`
* `Besu`
* `Erigon`

When it comes to consensus layer we currently support only `Prysm`.

## Command line

You can start a simulated network with a single command:
```
go run docker/test_env/cmd/main.go start-test-env private-chain
```

By default it will start a network with 1 node running `Geth` and `Prysm`. It will use default chain id of `1337` and won't wait for the chain to finalize at least one epoch. Once the chain is started it will save the network configuration in a `JSON` file, which then you can use in your tests to connect to that chain (and thus save time it takes to start a new chain each time you run your test).

Following cmd line flags are available:
```
  -c, --chain-id int             chain id (default 1337)
  -l, --consensus-layer string   consensus layer (prysm) (default "prysm")
  -t, --consensus-type string    consensus type (pow or pos) (default "pos")
  -e, --execution-layer string   execution layer (geth, nethermind, besu or erigon) (default "geth")
  -w, --wait-for-finalization    wait for finalization of at least 1 epoch (might take up to 5 mintues)
```

To connect to that environment in your tests use the following code:
```
	builder := NewEthereumNetworkBuilder()
	cfg, err := builder.
		WihtExistingConfigFromEnvVar().
		Build()

    if err != nil {
        return err
    }

	net, rpc, err := cfg.Start()
    if err != nil {
        return err
    }
```
Builder will read the location of chain configuration from env var named `PRIVATE_ETHEREUM_NETWORK_CONFIG_PATH` (it will be printed in the console once the chain starts).

`net` is an instance of `blockchain.EVMNetwork`, which contains characteristics of the network and can be used to connect to it using an EVM client. `rpc` variable contains arrays of public and private RPC endpoints, where "private" means URL that's accessible from the same Docker network as the chain is running in.
# Using LogStream

LogStream is a package that allows to connect to a Docker container and then flush logs to configured targets. Currently 3 targets are supported:
* `file` - saves logs to a file in `./logs` folder
* `loki` - sends logs to Loki
* `in-memory` - stores logs in memory

It can be configured to use multiple targets at once. If no target is specified, it becomes a no-op.

Targets can be set in two ways:
* using `LOGSTREAM_LOG_TARGETS` environment variable, e.g. `Loki,in-MemOry` (case insensitive)
* using programmatic functional option `WithLogTarget()`

Functional option has higher priority than environment variable.

When you connect a contaier LogStream will create a new consumer and start a detached goroutine that listens to logs emitted by that container and which reconnects and re-requests logs if listening fails for whatever reason. Retry limit and timeout can both be configured using functional options. In most cases one container should have one consumer, but it's possible to have multiple consumers for one container.

LogStream stores all logs in gob temporary file. To actually send/save them, you need to flush them. When you do it, LogStream will decode the file and send logs to configured targets. If log handling results in an error it won't be retried and processing of logs for given consumer will stop (if you think we should add a retry mechanism please let us know).

*Important:* Flushing and accepting logs is blocking operation. That's because they both share the same cursor to temporary file and otherwise it's position would be racey and could result in mixed up logs.

When using `in-memory` or `file` target no other environment variables are required. When using `loki` target, following environment variables are required:
* `LOKI_TENTANT_ID` - tenant ID
* `LOKI_URL` - Loki URL to which logs will be pushed
* `LOKI_BASIC_AUTH` -- only needed when running in CI and using public endpoint

Also, do remember that different `LOKI_URL` should be used when running in CI and everywhere else. In CI it should be a public endpoint, while in local environment it should be a private one.

You can print log location for each target using this function: `(m *LogStream) PrintLogTargetsLocations()`. For `file` target it will print relative folder path, for `loki` it will print URL of a Grafana Dashboard scoped to current execution and container ids. For `in-memory` target it's no-op.

It is recommended to shutdown LogStream at the end of your tests. Here's an example:
```go

t.Cleanup(func() {
    l.Warn().Msg("Shutting down Log Stream")

    if t.Failed() || os.Getenv("TEST_LOG_COLLECT") == "true" {
        // we can't do much if this fails, so we just log the error
        _ = logStream.FlushLogsToTargets()
        logStream.PrintLogTargetsLocations()
        logStream.SaveLogLocationInTestSummary()
    }

    // we can't do much if this fails, so we just log the error
    _ = logStream.Shutdown(testcontext.Get(b.t))
    })
```

## Grouping test execution

When running tests in CI you're probably interested in grouping logs by test execution, so that you can easily find the logs in Loki. To do that your job should set `RUN_ID` environment variable. In GHA it's recommended to set it to workflow id. If that variable is not set, then a run id will be automatically generated and saved in `.run.id` file, so that it can be shared by tests that are part of the same execution, but are running in different processes.

## Test Summary

In order to facilitate displaying information in GH's step summary `testsummary` package was added. It exposes a single function `AddEntry(testName, key string, value interface{}) `. When you call it, it either creates a test summary JSON file or appends to it. The result is is a map of keys with values.

Example:
```JSON
{
   "file":[
      {
         "test_name":"TestOCRv2Basic",
         "value":"./logs/TestOCRv2Basic-2023-12-01T18-00-59-TestOCRv2Basic-38ac1e52-d0a6-48"
      }
   ],
   "loki":[
      {
         "test_name":"TestOCRv2Basic",
         "value":"https://grafana.ops.prod.cldev.sh/d/ddf75041-1e39-42af-aa46-361fe4c36e9e/ci-e2e-tests-logs?orgId=1\u0026var-run_id=TestOCRv2Basic-38ac1e52-d0a6-48\u0026var-container_id=cl-node-a179ca7d\u0026var-container_id=cl-node-76798f87\u0026var-container_id=cl-node-9ff7c3ae\u0026var-container_id=cl-node-43409b09\u0026var-container_id=cl-node-3b6810bd\u0026var-container_id=cl-node-69fed256\u0026from=1701449851165\u0026to=1701450124925"
      }
   ]
}
```

In GHA after tests have ended we can use tools like `jq` to extract the information we need and display it in step summary.
