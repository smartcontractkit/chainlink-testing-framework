# Flakeguard

**Flakeguard** is a tool designed to help identify flaky tests within a Go project. Flaky tests are tests that intermittently fail without changes to the code, often due to race conditions or other non-deterministic behavior. Flakeguard assists by analyzing the impact of code changes on test packages and by running tests multiple times to determine stability.

## Features

- **Identify Impacted Tests**: Detects test packages that may be affected by changes in your Go project files.
- **Run Tests for Flakiness**: Runs tests multiple times to determine their flakiness.
- **Output Results in JSON**: Allows easy integration with CI pipelines and custom reporting.
- **Supports Exclusion Lists**: Configurable to exclude specified packages or paths from the analysis.
- **Recursive Dependency Analysis**: Detects all impacted packages through dependency levels.

## Installation

To install `flakeguard` CLI, you need to have Go installed on your machine. With Go installed, run the following command:

```sh
go install github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard@latest
```

## Usage

Flakeguard offers two main commands:

- `find` identifies test packages affected by recent changes.
- `run` executes tests multiple times to identify flaky tests

Run with `--help` to see all flags for the commands.

### JSON Output

Both `find` and `run` commands support JSON output `--json`, making it easy to integrate Flakeguard with CI/CD pipelines and reporting tools.

### Example Run

You can find example usage and see outputs with:

```sh
make example             # Run an example flow of running tests, aggregating results, and reporting them to GitHub
make example_flaky_panic # Run example flow with flaky and panicking tests
make example_timeout     # Run example flow tests the timeout
ls example_results       # See results of each run and aggregation
```
