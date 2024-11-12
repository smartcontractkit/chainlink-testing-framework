# End-to-End Testing Project Maturity Model

[![Smoke](https://img.shields.io/badge/Level_1-TestSmoke?branch=maturity-model&job=TestSmoke)](https://github.com/smartcontractkit/chainlink-testing-framework/actions/workflows/framework-golden-tests.yml)

[![Smoke](https://github.com/smartcontractkit/chainlink-testing-framework/actions/workflows/framework-golden-tests.yml/badge.svg?branch=maturity-model&job=TestSmoke)](https://github.com/smartcontractkit/chainlink-testing-framework/actions?query=workflow%3Aframework-golden-tests+branch%3Amaturity-model)
[![Performance Baseline](https://github.com/smartcontractkit/chainlink-testing-framework/actions/workflows/framework-golden-tests.yml/badge.svg?branch=maturity-model&job=PerformanceBaseline)](https://github.com/smartcontractkit/chainlink-testing-framework/actions/workflows/framework-golden-tests.yml)
[![Chaos](https://github.com/smartcontractkit/chainlink-testing-framework/actions/workflows/framework-golden-tests.yml/badge.svg?branch=maturity-model&job=TestChaos)](https://github.com/smartcontractkit/chainlink-testing-framework/actions/workflows/framework-golden-tests.yml)
[![Upgrade](https://github.com/smartcontractkit/chainlink-testing-framework/actions/workflows/framework-golden-tests.yml/badge.svg?branch=maturity-model&job=TestUpgrade)](https://github.com/smartcontractkit/chainlink-testing-framework/actions/workflows/framework-golden-tests.yml)

## Developing
Run the tests locally
```
CTF_CONFIGS=smoke.toml go test -v -run TestSmoke
CTF_CONFIGS=performance_baseline.toml go test -v -run TestPerformanceBaseline
CTF_CONFIGS=chaos.toml go test -v -run TestChaos
CTF_CONFIGS=upgrade.toml go test -v -run TestUpgrade
```