# Chainlink Integration Framework

![Tests](https://github.com/smartcontractkit/integrations-framework/actions/workflows/test.yaml/badge.svg)
![Lint](https://github.com/smartcontractkit/integrations-framework/actions/workflows/lint.yaml/badge.svg)

A framework for interacting with chainlink nodes, environments, and other blockchain systems. The framework is primarilly intended to facillitate testing chainlink features and stability.

## How to Test

1. Start a local hardhat network. You can easily do so by using our [docker container](https://hub.docker.com/r/smartcontract/hardhat-network). You could also deploy [your own local version](https://hardhat.org/hardhat-network/), if you are so inclined.
2. Run `go test ./...`

## // TODO

* Streamline the test running process
* Add more chainlink node checks
* Enable connecting chainlink node interfaces to actual running nodes in an environment
* Enable interaction with outside blockchains
* Check out [hardhat deploy](https://hardhat.org/plugins/hardhat-deploy.html) to help setup test environments
* Look into logging frameworks like [zerolog](https://github.com/rs/zerolog) or [zap](https://github.com/uber-go/zap)

Check out our [clubhouse board](https://app.clubhouse.io/chainlinklabs/project/5690/qa-team?vc_group_by=day) for a look into our progress.
