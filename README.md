# Chainlink Integration Framework

[![Go Report Card](https://goreportcard.com/badge/github.com/smartcontractkit/integrations-framework)](https://goreportcard.com/report/github.com/smartcontractkit/integrations-framework)
![Tests](https://github.com/smartcontractkit/integrations-framework/actions/workflows/test.yaml/badge.svg)
![Lint](https://github.com/smartcontractkit/integrations-framework/actions/workflows/lint.yaml/badge.svg)

A framework for interacting with chainlink nodes, environments, and other blockchain systems.
The framework is primarilly intended to facillitate testing chainlink features and stability.

## WIP

This framework is still very much a work in progress, and will have frequent changes, many of which will probably be
breaking.

## Execution Environment

Ephemeral environments are automatically deployed with Kubernetes. To run tests, you either need a deployed cluster 
in an environment, or a local installation.

### Locally

When running tests locally, it's advised to use minikube. To spin up a cluster, use:

```
minikube start
```

### Remotely

To run against a remote Kubernetes cluster, ensure your current context is the cluster you want to run against as the 
framework always uses current context.

## Test Execution

This framework advises the use of [Ginkgo](https://github.com/onsi/ginkgo) for test execution, but tests still can be 
ran with the go CLI.

### Ginkgo

Run:

```
ginkgo -r
```

### Go

Run:

```
go test ./..
```

### Volume tests
```
NETWORK="ethereum_geth_volume" make test_performance
```
