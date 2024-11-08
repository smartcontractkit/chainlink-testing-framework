# "Golden" Templates

[Here](https://github.com/smartcontractkit/chainlink-testing-framework/actions/runs/11739154666/job/32703095118?pr=1311) are our "golden" templates for end-to-end tests, covering every test type:

- `Smoke`
- `Performance`
- `Chaos`
- `Upgrade`

These tests act as a maturity model and are implemented across all our products.

Use this [workflow](https://github.com/smartcontractkit/chainlink-testing-framework/actions/runs/11739154666/workflow?pr=1311) as a starting point for developing a new end-to-end integration test.

Set the count to 5-10 during development, and once stable, set the timeout and proceed to merge.