# End-to-End Testing Project Maturity Model

[![Smoke](https://img.shields.io/badge/Level_1-Smoke-blue?branch=maturity-model&job=TestSmoke)](https://github.com/smartcontractkit/chainlink-testing-framework/actions?query=workflow%3Aframework-golden-tests+branch%3Amaturity-mode)
[![Upgrade](https://img.shields.io/badge/Level_2-Upgrade-blue?branch=maturity-model&job=TestSmoke)](https://github.com/smartcontractkit/chainlink-testing-framework/actions?query=workflow%3Aframework-golden-tests+branch%3Amaturity-mode)
[![Performance Baseline](https://img.shields.io/badge/Level_3-Performance_baseline-blue?branch=maturity-model&job=TestSmoke)](https://github.com/smartcontractkit/chainlink-testing-framework/actions?query=workflow%3Aframework-golden-tests+branch%3Amaturity-mode)
[![Resiliency](https://img.shields.io/badge/Level_4-Resiliency-blue?branch=maturity-model&job=TestSmoke)](https://github.com/smartcontractkit/chainlink-testing-framework/actions?query=workflow%3Aframework-golden-tests+branch%3Amaturity-mode)
[![Scalability](https://img.shields.io/badge/Level_5-Scalability-blue?branch=maturity-model&job=TestSmoke)](https://github.com/smartcontractkit/chainlink-testing-framework/actions?query=workflow%3Aframework-golden-tests+branch%3Amaturity-mode)

## Level 1
The team maintains a system-level smoke test where all components are deployed using `docker`.

All on-chain changes are done through [chainlink-deployments](https://github.com/smartcontractkit/chainlink-deployments).

The test is readable, and the README clearly explains its purpose.

The test is reliable and stable when run with a `-count 30`.

If your project includes multiple use cases and functionality suitable for end-to-end testing, you can add additional tests at this level.

## Level 2
The team has an "upgrade" test to verify product compatibility with older versions.

While the number of compatible versions is team-determined, identifying incompatibilities at the system level early is a valuable, mature practice.

This test deploys specific platform and plugin versions, performs an end-to-end smoke test, and then upgrades (or migrates) the plugin(s) or platform on the same database to ensure that users remain unaffected by the upgrade.

## Level 3
The team has a baseline performance testing suite.

At this level, the focus is not on improving performance but on detecting any performance degradation, supported by a reliable CI pipeline.

This pipeline runs as needed—nightly or before releases—enabling early detection of performance issues across all critical on-chain and off-chain functionality.

## Level 4
The team incorporates chaos engineering practices to test the system’s failure modes.

This stage builds on [Level 3](#level-3), as it not only verifies that the system is reliable and can recover from reasonable failures but also ensures an understanding of how these failures impact performance and user experience.

## Level 5
The team has complete ownership of their persistent staging environment.

They can perform upgrades, data migrations, and run advanced load tests to validate the scalability of their applications.

## Developing
Run the tests locally
```
CTF_CONFIGS=smoke.toml go test -v -run TestSmoke
CTF_CONFIGS=upgrade.toml go test -v -run TestUpgrade
CTF_CONFIGS=performance_baseline.toml go test -v -run TestPerformanceBaseline
CTF_CONFIGS=chaos.toml go test -v -run TestChaos
CTF_CONFIGS=scalability.toml go test -v -run TestScalability
```