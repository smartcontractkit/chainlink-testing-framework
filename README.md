# Chainlink Integration Framework

[![Go Report Card](https://goreportcard.com/badge/github.com/smartcontractkit/integrations-framework)](https://goreportcard.com/report/github.com/smartcontractkit/integrations-framework)
![Tests](https://github.com/smartcontractkit/integrations-framework/actions/workflows/test.yaml/badge.svg)
![Lint](https://github.com/smartcontractkit/integrations-framework/actions/workflows/lint.yaml/badge.svg)

A framework for interacting with chainlink nodes, environments, and other blockchain systems.
The framework is primarilly intended to facillitate testing chainlink features and stability.

## WIP

This framework is still very much a work in progress, and will have frequent changes, many of which will probably be
breaking.

## How to Test

1. Have a K8s cluster running that you can connect to locally. [Minikube](https://minikube.sigs.k8s.io/docs/)
   works well for local testing.
2. Run `go test ./...`

## Example Usage

You can see our tests for some basic usage examples. The most complete can be found in `contracts/contracts_test.go`
