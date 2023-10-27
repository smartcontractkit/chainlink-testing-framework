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