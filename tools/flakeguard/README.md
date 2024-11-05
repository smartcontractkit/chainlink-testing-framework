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

Flakeguard offers two main commands: `find` and `run`.

### `find` Command

The `find` command identifies test packages affected by recent changes.

#### Flags for `find` Command

- `--project-path, -r` : Path to the Go project. Defaults to the current directory.
- `--base-ref` (required): Git reference (branch, tag, or commit) for comparing changes.
- `--verbose, -v` : Enables verbose mode.
- `--json` : Outputs results in JSON format.
- `--filter-empty-tests` : Filters out test packages with no actual test functions (can be slow on large projects).
- `--excludes` : List of paths to exclude.
- `--levels, -l` : Levels of recursion for dependency search, with `0` for unlimited. `2` by default.
- `--find-by-test-files-diff` : Enable mode to find affected test packages by changes in test files.
- `--find-by-affected-packages` : Enable mode to find affected test packages based on changes in any project package.
- `--only-show-changed-test-files` : Display only changed test files and exit.

### `run` Command

The `run` command executes tests multiple times to identify flaky tests.

#### Flags for `run` Command

- `--project-path, -r` : Path to the Go project. Defaults to the current directory.
- `--test-packages-json` : JSON-encoded string of test packages to run.
- `--test-packages` : List of test packages to run.
- `--run-count, -c` : Number of times to run the tests.
- `--race` : Enable race condition detection.
- `--output-json` : File path to output test results in JSON format.
- `--threshold` : Threshold (0-1) for determining flakiness (e.g., `0.8` to pass if 80% successful).
- `--skip-tests` : List of tests to skip.
- `--fail-fast` : Stop on the first failure if threshold is set to `1.0`.

### JSON Output

Both `find` and `run` commands support JSON output `--json`, making it easy to integrate Flakeguard with CI/CD pipelines and reporting tools.
