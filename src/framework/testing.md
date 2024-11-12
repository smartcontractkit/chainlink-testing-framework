# Testing Maturity Model

[Here](https://github.com/smartcontractkit/chainlink-testing-framework/actions/runs/11739154666/job/32703095118?pr=1311) are our "golden" templates for end-to-end tests, covering every test type:

- `Smoke`
- `PerformanceBaseline`
- `Chaos`
- `Upgrade`

These tests act as a maturity model and are implemented across all our products.

Refer to this README to understand the rationale behind our testing approach and to explore the stages of maturity in end-to-end testing.

## Developing
Run the tests locally
```
CTF_CONFIGS=smoke.toml go test -v -run TestSmoke
CTF_CONFIGS=performance_baseline.toml go test -v -run TestPerformanceBaseline
CTF_CONFIGS=chaos.toml go test -v -run TestChaos
CTF_CONFIGS=upgrade.toml go test -v -run TestUpgrade
```
