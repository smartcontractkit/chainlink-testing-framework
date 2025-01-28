# End-to-End Testing Project Maturity Model

[![Smoke](https://img.shields.io/badge/Level_1-Smoke-blue?branch=maturity-model&job=TestSmoke)](https://github.com/smartcontractkit/chainlink-testing-framework/actions/workflows/framework-golden-tests.yml)
[![Upgrade](https://img.shields.io/badge/Level_2-Upgrade-blue?branch=maturity-model&job=TestSmoke)](https://github.com/smartcontractkit/chainlink-testing-framework/actions/workflows/framework-golden-tests.yml)
[![Performance Baseline](https://img.shields.io/badge/Level_3-Performance_baseline-blue?branch=maturity-model&job=TestSmoke)](https://github.com/smartcontractkit/chainlink-testing-framework/actions/workflows/framework-golden-tests.yml)
[![Resiliency](https://img.shields.io/badge/Level_4-Resiliency-blue?branch=maturity-model&job=TestSmoke)](https://github.com/smartcontractkit/chainlink-testing-framework/actions/workflows/framework-golden-tests.yml)
[![Scalability](https://img.shields.io/badge/Level_5-Scalability-blue?branch=maturity-model&job=TestSmoke)](https://github.com/smartcontractkit/chainlink-testing-framework/actions/workflows/framework-golden-tests.yml)

## Level 0

The team creates and maintains a high-level test plan outlining the components involved and test cases in any format.

If the team decides on minimal or no manual testing and the project is trivial, they can consolidate all test cases into `go test` cases, outline the required implementations and commit templates up front.

The team identifies potential integration points with third-party software, blockchain, and external services and document any testing limitations.

If the system is complex and the team wants to test parts of it using mocks, a schema should be provided to outline which components will be separated and include examples of the mocks.

If new components are required, the team implements them following this guide: [Developing Components](https://smartcontractkit.github.io/chainlink-testing-framework/developing/developing_components.html).


## Level 1
The team maintains a system-level smoke test where all components are deployed using `docker`.

All on-chain changes are done through [chainlink-deployments](https://github.com/smartcontractkit/chainlink-deployments).

The test should not depend on `chainlink` core types.

The test is readable, and the README clearly explains what is tested.

The test is stable when run with a `-count 10`.

If the team has more than two distinct NodeSets (DONs), system-level tests in ephemeral environments should be limited to two NodeSets to allow for chaos testing later.

If your project includes multiple use cases and functionality suitable for end-to-end testing, you can add additional tests at this level.

## Level 2
The team has an "upgrade" test to verify product compatibility with older versions.

While the number of compatible versions is team-determined, identifying incompatibilities at the system level early is a valuable, mature practice.

This test deploys specific platform and plugin versions, performs an end-to-end smoke test, and then upgrades (or migrates) the plugin(s) or platform on the same database to ensure that users remain unaffected by the upgrade or, in case of breaking changes, migration process is tested.

## Level 3
The team has a baseline performance testing suite.

At this level, the focus is not on improving performance but on detecting any performance degradation early, supported by a reliable CI pipeline.

This pipeline runs as needed—nightly or before releases—enabling early detection of performance issues across all critical on-chain and off-chain functionality.

This stage combines performance testing with observability enhancements. The team should have fundamental system-level performance tests in place, along with dashboards to monitor product-specific metrics.

## Level 4
The team incorporates chaos engineering practices to test the system’s failure modes.

This stage builds on [Level 3](#level-3), as it not only verifies that the system is reliable and can recover from reasonable failures but also ensures an understanding of how these failures impact performance and user experience.

Recommended test cases:
- Component failure (reboot)
- Slow network
- RPC node outage
- Stress test (limiting container resources)

## Level 5
The team has complete ownership of their persistent staging environment.

The team can use [components caching](https://smartcontractkit.github.io/chainlink-testing-framework/framework/components/caching.html) to reuse testing logic on staging and substitute staging environment URLs instead of an ephemeral environment.

They can perform upgrades, data migrations, and run advanced load tests to validate the scalability of their applications.

## Explanation

It’s essential not to skip levels, as they enable us to manage complexity gradually and maintain focus on our product.

`Level 0` is crucial for understanding what and how we’ll test, planning test levels, and distinguishing between what can be tested in low-level integrations versus system-level testing.

`Level 1` centers on testing basic functionality end-to-end, analyzing configurability, and using the simplest infrastructure setup—Docker on a local machine.

`Levels 2-4` help the team establish a pipeline for migrations and prepare the product for real-world scenarios, such as load and chaos testing. Even if the initial product iteration isn’t perfect, it's critical to start measuring and analyzing performance degradation and failures early.

`Level 5` is the final stage before production. By this point, the product should be mature enough for stable staging, allowing the team to catch issues related to functionality, performance, failures, and incompatibilities before release. Staging is typically shared across the team; skipping levels and deploying early prototypes here can make debugging harder, risk disrupting other developers, and hinder quick, on-demand testing with low infrastructure complexity.

## Developing
Run the tests locally
```
CTF_CONFIGS=smoke.toml go test -v -run TestSmoke
CTF_CONFIGS=upgrade.toml go test -v -run TestUpgrade
CTF_CONFIGS=performance_baseline.toml go test -v -run TestPerformanceBaseline
CTF_CONFIGS=chaos.toml go test -v -run TestChaos
CTF_CONFIGS=scalability.toml go test -v -run TestScalability
```