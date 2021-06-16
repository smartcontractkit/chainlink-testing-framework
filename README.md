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

1. Start a local hardhat network. You can easily do so by using our
 [docker container](https://hub.docker.com/r/smartcontract/hardhat-network). You could also deploy
 [your own local version](https://hardhat.org/hardhat-network/), if you are so inclined.
2. Start few local chainlink nodes, utilizing our `docker-compose` setup
   [here](https://github.com/smartcontractkit/chainlink-node-compose)
3. Run `go test ./...`

## Example Usage

You can see our tests for some basic usage examples. The most complete can be found in `contracts/contracts_test.go`
