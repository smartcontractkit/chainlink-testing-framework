# "Golden" Templates

[Here](https://github.com/smartcontractkit/chainlink-testing-framework/actions/runs/11739154666/job/32703095118?pr=1311) are our "golden" templates for end-to-end tests, covering every test type:

- `Smoke`
- `Performance`
- `Chaos`
- `Upgrade`

These tests act as a maturity model and are implemented across all our products.

Run them locally:
```
CTF_CONFIGS=smoke.toml go test -v -run TestSmoke
CTF_CONFIGS=load.toml go test -v -run TestLoad
CTF_CONFIGS=chaos.toml go test -v -run TestChaos
CTF_CONFIGS=upgrade_some.toml go test -v -run TestUpgradeSome
```

Use this [workflow](https://github.com/smartcontractkit/chainlink-testing-framework/actions/runs/11739154666/workflow?pr=1311) as a starting point for developing a new end-to-end integration test.

Set the count to 5-10 during development, and once stable, set the timeout and proceed to merge.

If you need to structure a lot of different tests (not only end-to-end) follow [this](https://github.com/smartcontractkit/.github/tree/main/.github/workflows) guide.