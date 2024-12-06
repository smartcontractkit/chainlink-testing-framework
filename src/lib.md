<div align="center">

# Framework v1 (Deprecated)

[![Lib tag](https://img.shields.io/github/v/tag/smartcontractkit/chainlink-testing-framework?filter=%2Alib%2A)](https://github.com/smartcontractkit/chainlink-testing-framework/tags)
[![Go Report Card](https://goreportcard.com/badge/github.com/smartcontractkit/chainlink-testing-framework)](https://goreportcard.com/report/github.com/smartcontractkit/chainlink-testing-framework)
[![Go Reference](https://pkg.go.dev/badge/github.com/smartcontractkit/chainlink-testing-framework.svg)](https://pkg.go.dev/github.com/smartcontractkit/chainlink-testing-framework)
[![Go Version](https://img.shields.io/github/go-mod/go-version/smartcontractkit/chainlink-testing-framework?filename=./lib/go.mod)](https://go.dev/)
![Tests](https://github.com/smartcontractkit/chainlink-testing-framework/actions/workflows/test.yaml/badge.svg)
![Lint](https://github.com/smartcontractkit/chainlink-testing-framework/actions/workflows/lint.yaml/badge.svg)

</div>

**DEPRECATED: This is v1 version and it is not actively maintained**

The purpose of this framework is to:
- Interact with different blockchains
- Configure CL jobs
- Deploy using `docker`
- Deploy using `k8s`

If you're looking to implement a new chain integration for the testing framework, head over to the [blockchain](lib/blockchain.md) directory for more info.

## k8s package

We have a k8s package we are using in tests, it provides:

- [cdk8s](https://cdk8s.io/) based wrappers
- High-level k8s API
- Automatic port forwarding

You can also use this package to spin up standalone environments.

### Local k8s cluster

Read [here](lib/k8s/KUBERNETES.md) about how to spin up a local cluster

#### Install

Set up deps, you need to have `node 14.x.x`, [helm](https://helm.sh/docs/intro/install/) and [yarn](https://classic.yarnpkg.com/lang/en/docs/install/#mac-stable)

Then use

```shell
make install_deps
```

##### Optional Nix

We have setup a nix shell which will produce a reliable environment that will behave the same locally and in ci. To use it instead of the above you will need to [install nix](https://nixos.org/download/)

To start the nix shell run:

```shell
make nix_shell
```

If you install [direnv](https://github.com/direnv/direnv/blob/master/docs/installation.md) you will be able to have your environment start the nix shell as soon as you cd into it once you have allowed the directory via:

```shell
direnv allow
```

### Running tests in k8s

To read how to run a test in k8s, read [here](lib/k8s/REMOTE_RUN.md)

### Usage

#### With env vars (deprecated)

Create an env in a separate file and run it

```sh
export CHAINLINK_IMAGE="public.ecr.aws/chainlink/chainlink"
export CHAINLINK_TAG="1.4.0-root"
export CHAINLINK_ENV_USER="Satoshi"
go run k8s/examples/simple/env.go
```

For more features follow [tutorial](lib/k8s/TUTORIAL.md)

#### With TOML config

It should be noted that using env vars for configuring CL nodes in k8s is deprecated. TOML config should be used instead:

```toml
[ChainlinkImage]
image="public.ecr.aws/chainlink/chainlink"
version="v2.8.0"
```

Check the example here: [env.go](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/lib/k8s/examples/simple_toml/env_toml_config.go)

### Development

#### Running standalone example environment

```shell
go run k8s/examples/simple/env.go
```

If you have another env of that type, you can connect by overriding environment name

```sh
ENV_NAMESPACE="..."  go run k8s/examples/chainlink/env.go
```

Add more presets [here](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/lib/k8s/presets)

Add more programmatic examples [here](../../lib/k8s/examples/)

If you have [chaosmesh](https://chaos-mesh.org/) installed in your cluster you can pull and generated CRD in go like that

```sh
make chaosmesh
```

If you need to check your system tests coverage, use [that](../../lib/k8s/TUTORIAL.md#coverage)

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

- `Geth`
- `Nethermind`
- `Besu`
- `Erigon`

When it comes to consensus layer we currently support only `Prysm`.

The easiest way to start a simulated network is to use a builder. It allows to configure the network in a fluent way and then start it. For example:

```go
builder := NewEthereumNetworkBuilder()
cfg, err: = builder.
    WithEthereumVersion(EthereumVersion_Eth2).
    WithExecutionLayer(ExecutionLayer_Geth).
    Build()
```

Since we support both `eth1` (aka pre-Merge) and `eth2` (aka post-Merge) client versions, you need to specify which one you want to use. You can do that by calling `WithEthereumVersion` method. There's no default provided. The only exception is when you use custom docker images (instead of default ones), because then we can determine which version it is based on the image version.

If you want your test to execute as fast as possible go for `eth1` since it's either using a fake PoW or PoA consensus and is much faster than `eth2` which uses PoS consensus (where there is a minimum viable length of slot/block, which is 4 seconds; for `eth1` it's 1 second). If you want to test the latest features, changes or forks in the Ethereum network and have your tests running on a network which is as close as possible to Ethereum Mainnet, go for `eth2`.

Every component has some default Docker image it uses, but builder has a method that allows to pass custom one:

```go
builder := NewEthereumNetworkBuilder()
cfg, err: = builder.
    WithEthereumVersion(EthereumVersion_Eth2).
    WithConsensusLayer(ConsensusLayer_Prysm).
    WithExecutionLayer(ExecutionLayer_Geth).
    WithCustomDockerImages(map[ContainerType]string{
        ContainerType_Geth: "my-custom-geth-pos-image:my-version"}).
    Build()
```

When using a custom image you can even further simplify the builder by calling only `WithCustomDockerImages` method. Based on the image name and version we will determine which execution layer client it is and whether it's `eth1` or `eth2` client:

```go
builder := NewEthereumNetworkBuilder()
cfg, err: = builder.
    WithCustomDockerImages(map[ContainerType]string{
        ContainerType_Geth: "ethereum/client-go:v1.13.10"}).
    Build()
```

In the case above we would launch a `Geth` client with `eth2` network and `Prysm` consensus layer.

You can also configure epochs at which hardforks will happen. Currently only `Deneb` is supported. Epoch must be >= 1. Example:

```go
builder := NewEthereumNetworkBuilder()
cfg, err: = builder.
    WithConsensusType(ConsensusType_PoS).
    WithConsensusLayer(ConsensusLayer_Prysm).
    WithExecutionLayer(ExecutionLayer_Geth).
    WithEthereumChainConfig(EthereumChainConfig{
        HardForkEpochs: map[string]int{"Deneb": 1},
    }).
    Build()
```

## Command line

You can start a simulated network with a single command:

```sh
go run docker/test_env/cmd/main.go start-test-env private-chain
```

By default it will start a network with 1 node running `Geth` and `Prysm`. It will use default chain id of `1337` and won't wait for the chain to finalize at least one epoch. Once the chain is started it will save the network configuration in a `JSON` file, which then you can use in your tests to connect to that chain (and thus save time it takes to start a new chain each time you run your test).

Following cmd line flags are available:

```sh
  -c, --chain-id int             chain id (default 1337)
  -l, --consensus-layer string   consensus layer (prysm) (default "prysm")
  -t, --consensus-type string    consensus type (pow or pos) (default "pos")
  -e, --execution-layer string   execution layer (geth, nethermind, besu or erigon) (default "geth")
  -w, --wait-for-finalization    wait for finalization of at least 1 epoch (might take up to 5 minutes)
      --consensus-client-image string   custom Docker image for consensus layer client
      --execution-layer-image string    custom Docker image for execution layer client
      --validator-image string          custom Docker image for validator
```

To connect to that environment in your tests use the following code:

```go
	builder := NewEthereumNetworkBuilder()
	cfg, err := builder.
		WithExistingConfigFromEnvVar().
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

## Logs
By default, we will save logs of all Docker containers running on the host machine, when the test ends (regardless whether it failed or succeeded). They will be available in the `./logs/<test-name><date>` directory. Same goes for dumping the databases of PostgresDBs
used by the Chainlink nodes. These will be saves in the `./db_dumps/<test-name><date>` directory.

## Loki and Grafana

If you need to pass Loki or Grafana configuration to your tests you can do that by providing the following config:

```toml
[Logging.Loki]
tenant_id="promtail"
url="https://change.me"
basic_auth_secret="my-secret-auth"
bearer_token_secret="bearer-token"
```

Also, do remember that different URL should be used when running in CI and everywhere else. In CI it should be a public endpoint, while in local environment it should be a private one.

If your test has a Grafana dashboard you should provide the following config:

```toml
[Logging.Grafana]
base_url="https://your-grafana-url"
dashboard_url="/my-dashboard"
dashboard_uid="my-dashboard-uid" # optional
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

# TOML Config

Basic and universal building blocks for TOML-based config are provided by `config` package. For more information do read [this](lib/config/config.md).

# ECR Mirror

An ecr mirror can be used to push images used often in order to bypass rate limit issues from dockerhub. The list of image mirrors can be found in the [matrix here](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/.github/workflows/update-internal-mirrors.yaml). This currently works with images with version numbers in dockerhub. Support for gcr is coming in the future. The images must also have a version number so putting `latest` will not work. We have a separate list for one offs we want that can be added to [here](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/.github/actions/update-internal-mirrors/scripts/mirror.json) that does work with gcr and latest images. Note however for `latest` it will only pull it once and will not update it in our mirror if the latest on the public repository has changed, in this case it is preferable to update it manually when you know that you need the new latest and the update will not break your tests.

For images in the mirrors you can use the INTERNAL_DOCKER_REPO environment variable when running tests and it will use that mirror in place of the public repository.

We have two ways to add new images to the ecr. The first two requirements are that you create the ecr repository with the same name as the one in dockerhub out in aws and then add that ecr to the infra permissions (ask TT if you don't know how to do this).

1. If it does not have version numbers or is gcr then you can add it [here](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/.github/actions/update-internal-mirrors/scripts/mirror.json)
2. You can add to the [mirror matrix](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/.github/actions/update-internal-mirrors/scripts/update_mirrors.sh) the new image name and an expression to get the latest versions added when the workflow runs. You can check the postgres one used in there for an example but basically the expression should filter out only the latest image or 2 for that particular version when calling the dockerhub endpoint, example curl call `curl -s "https://hub.docker.com/v2/repositories/${image_name}/tags/?page_size=100" | jq -r '.results[].name' | grep -E ${image_expression}` where image_name could be `library/postgres` and image_expression could be `'^[0-9]+\.[0-9]+$'`. Adding your ecr to this matrix should make sure we always have the latest versions for that expression.

## Debugging HTTP and RPC calls

```bash
export SETH_LOG_LEVEL=info
export RESTY_DEBUG=true
```

## Loki Client

The `LokiClient` allows you to easily query Loki logs from your tests. It supports basic authentication, custom queries, and can be configured for (Resty) debug mode.

### Debugging Resty and Loki Client

```bash
export LOKI_CLIENT_LOG_LEVEL=info
export RESTY_DEBUG=true
```

### Example usage:

```go
auth := LokiBasicAuth{
    Username: os.Getenv("LOKI_LOGIN"),
    Password: os.Getenv("LOKI_PASSWORD"),
}

queryParams := LokiQueryParams{
		Query:     `{namespace="test"} |= "test"`,
		StartTime: time.Now().AddDate(0, 0, -1),
		EndTime:   time.Now(),
		Limit:     100,
	}

lokiClient := NewLokiClient("https://loki.api.url", "my-tenant", auth, queryParams)
logEntries, err := lokiClient.QueryLogs(context.Background())
```
