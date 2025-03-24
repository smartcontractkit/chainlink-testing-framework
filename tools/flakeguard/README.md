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
- `create-tickets` automates the creation of Jira tickets for flaky tests detected by Flakeguard

Run with `--help` to see all flags for the commands.

### JSON Output

Both `find` and `run` commands support JSON output `--json`, making it easy to integrate Flakeguard with CI/CD pipelines and reporting tools.

### Creating JIRA Tickets
The `create-tickets` command allows you to automate the creation of JIRA tickets for flaky tests. It reads test results from a CSV file (typically exported from a Splunk view) and creates tickets in JIRA.

```
go run main.go create-tickets --jira-project=<JIRA_PROJECT_KEY> --flaky-test-json-db-path=<PATH_TO_FLAKY_TEST_DB_JSON> --assignee-mapping=<PATH_TO_JIRA_ASSIGNEE_MAPPING_JSON> --csv-path=<PATH_TO_CSV_FILE> [--skip-existing] [--dry-run]
```

Example:
```
go run main.go create-tickets --jira-project=DX --flaky-test-json-db-path=.flaky_test_db.json --assignee-mapping=.jira_assignee_mapping.json --skip-existing --csv-path '1742825894_77903.csv'
```

**Options:**

- `--jira-project`: The JIRA project key where tickets should be created (e.g., `DX`).
- `--test-db-path`: The path to a JSON database (`.json`) that stores information about existing flaky test tickets.
- `--assignee-mapping`: The path to a JSON file (`.json`) that maps test packages to JIRA assignees.
- `--csv-path`: The path to the CSV file containing the flaky test results.
- `--skip-existing`: (Optional) Skips creating tickets for tests that already have corresponding JIRA tickets in the database or JIRA.
- `--dry-run`: (Optional) Performs a dry run without actually creating JIRA tickets.

**Environment Variables:**

- `JIRA_DOMAIN`: The domain of your JIRA instance.
- `JIRA_EMAIL`: The email address used to authenticate with JIRA.
- `JIRA_API_KEY`: The API key used to authenticate with JIRA.

### Managing Tickets
The new tickets command lets you interactively manage your flaky test tickets stored in your local JSON database. With this command you can:

- Mark tests as skipped: Post a comment to the associated JIRA ticket and update the local DB.
- Unskip tests: Remove the skipped status and post an unskip comment.
- Navigate tickets: Use a TUI to browse through tickets.

```
go run main.go tickets --test-db-path=.flaky_test_db.json --jira-comment
```

**Options:**

- `--test-db-path`: The path to a JSON database (.json) that stores information about your flaky test tickets.
- `--jira-comment`: If set to true, posts a comment to the corresponding JIRA ticket when marking a test as skipped or unskipped.
`--dry-run`: (Optional) Runs the command in dry-run mode without posting comments to JIRA.
- `--hide-skipped`: (Optional) If set, tickets already marked as skipped are hidden from the interface.

**Environment Variables:**

- `JIRA_DOMAIN`: The domain of your JIRA instance.
- `JIRA_EMAIL`: The email address used to authenticate with JIRA.
- `JIRA_API_KEY`: The API key used to authenticate with JIRA.

### Example Run

You can find example usage and see outputs with:

```sh
make example             # Run an example flow of running tests, aggregating results, and reporting them to GitHub
make example_flaky_panic # Run example flow with flaky and panicking tests
ls example_results       # See results of each run and aggregation
```
