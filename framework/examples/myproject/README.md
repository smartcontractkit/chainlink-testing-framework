# End-to-End Testing Project Maturity Model

[![Smoke](https://img.shields.io/badge/Level_1-Smoke-blue?branch=maturity-model&job=TestSmoke)](https://github.com/smartcontractkit/chainlink-testing-framework/actions?query=workflow%3Aframework-golden-tests+branch%3Amaturity-mode)
[![Upgrade](https://img.shields.io/badge/Level_2-Upgrade-blue?branch=maturity-model&job=TestSmoke)](https://github.com/smartcontractkit/chainlink-testing-framework/actions?query=workflow%3Aframework-golden-tests+branch%3Amaturity-mode)
[![Performance Baseline](https://img.shields.io/badge/Level_3-Performance_baseline-blue?branch=maturity-model&job=TestSmoke)](https://github.com/smartcontractkit/chainlink-testing-framework/actions?query=workflow%3Aframework-golden-tests+branch%3Amaturity-mode)
[![Resiliency](https://img.shields.io/badge/Level_4-Resiliency-blue?branch=maturity-model&job=TestSmoke)](https://github.com/smartcontractkit/chainlink-testing-framework/actions?query=workflow%3Aframework-golden-tests+branch%3Amaturity-mode)
[![Scalability](https://img.shields.io/badge/Level_5-Scalability-blue?branch=maturity-model&job=TestSmoke)](https://github.com/smartcontractkit/chainlink-testing-framework/actions?query=workflow%3Aframework-golden-tests+branch%3Amaturity-mode)

## Developing
Run the tests locally
```
CTF_CONFIGS=smoke.toml go test -v -run TestSmoke
CTF_CONFIGS=upgrade.toml go test -v -run TestUpgrade
CTF_CONFIGS=performance_baseline.toml go test -v -run TestPerformanceBaseline
CTF_CONFIGS=chaos.toml go test -v -run TestChaos
CTF_CONFIGS=scalability.toml go test -v -run TestScalability
```