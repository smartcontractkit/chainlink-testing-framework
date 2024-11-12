# End-to-End Testing Project Maturity Model

## Developing
Run the tests locally
```
CTF_CONFIGS=smoke.toml go test -v -run TestSmoke
CTF_CONFIGS=performance_baseline.toml go test -v -run TestPerformanceBaseline
CTF_CONFIGS=chaos.toml go test -v -run TestChaos
CTF_CONFIGS=upgrade.toml go test -v -run TestUpgrade
```