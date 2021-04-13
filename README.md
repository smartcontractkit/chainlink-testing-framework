# Chainlink Integration Framework
![Tests](https://github.com/smartcontractkit/integrations-framework/actions/workflows/test.yaml/badge.svg)
![Lint](https://github.com/smartcontractkit/integrations-framework/actions/workflows/lint.yaml/badge.svg)

A framework for interacting with chainlink nodes, environments, and other blockchain systems. The framework is primarilly intended to facillitate testing chainlink features and stability.

### How to Test
`npx hardhat node` to start a local hardhat instance of ethereum

then

`go test ./...`

### // TODO
* Streamline the test running process
* Add more chainlink node checks
* Enable connecting chainlink node interfaces to actual running nodes in an environment
* Enable interaction with outside blockchains

Check out our [clubhouse board](https://app.clubhouse.io/chainlinklabs/project/5690/qa-team?vc_group_by=day) for a look into our progress.