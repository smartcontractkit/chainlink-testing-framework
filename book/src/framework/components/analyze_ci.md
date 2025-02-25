# Analyzing CI Runs

We offer a straightforward CLI tool designed to analyze CI runs, focusing on Jobs and Steps, to provide deeper insights into system-level tests.

## Examples
```
# GITHUB_TOKEN must have access to "actions" API
export GITHUB_TOKEN=...

# E2E tests from core, the last day
ctf ci -r "smartcontractkit/chainlink" -w "Integration Tests"

# Last 3 days runs for e2e framework tests
ctf ci -r "smartcontractkit/chainlink-testing-framework" -w "Framework Golden Tests Examples" -d 3
```
