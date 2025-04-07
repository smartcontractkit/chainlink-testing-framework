Flakeguard
==========

**Flakeguard** is a tool designed to help identify flaky tests within a Go project. Flaky tests are tests that intermittently fail without changes to the code, often due to race conditions or other non-deterministic behavior. Flakeguard assists by analyzing the impact of code changes on test packages and by running tests multiple times to determine stability.

In addition to detecting flaky tests, Flakeguard can also integrate with Jira to track known flaky tests. It maintains a local database of tests and their associated Jira tickets, allowing you to create and manage these tickets directly from the command line.

Features
--------

*   **Identify Impacted Tests:** Detects test packages that may be affected by changes in your Go project files.
*   **Run Tests for Flakiness:** Runs tests multiple times to determine their flakiness.
*   **Output Results in JSON:** Allows easy integration with CI pipelines and custom reporting.
*   **Supports Exclusion Lists:** Configurable to exclude specified packages or paths from the analysis.
*   **Recursive Dependency Analysis:** Detects all impacted packages through dependency levels.
*   **Jira Integration (Optional):** Create, review, and manage flaky test tickets in Jira.
*   **Local Database:** Store known flaky tests and their associated Jira tickets for easy reference.

Prerequisites
-------------

1.  **Go:** Version 1.21 or later recommended.
2.  **Jira API Access (optional):** Required only if you want to use the Jira-related commands (`create-tickets`, `review-tickets`, `sync-jira`):
    *   `JIRA_DOMAIN`: Your Jira instance domain (e.g. _your-company.atlassian.net_).
    *   `JIRA_EMAIL`: The email address associated with your Jira account used for API access.
    *   `JIRA_API_KEY`: Your Jira API token (generate one in your Jira account settings).
    *   `JIRA_PROJECT_KEY`: (Optional) The default Jira project key to use if `--jira-project` is not specified for `create-tickets`.

Installation
------------

To install the `flakeguard` CLI, ensure you have Go installed. Then run:

go install github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard@latest

You can also clone the repository and build from source:

git clone https://github.com/smartcontractkit/chainlink-testing-framework.git
cd chainlink-testing-framework/tools/flakeguard

go build -o flakeguard main.go


Usage (Flaky Test Detection)
----------------------------------

Flakeguard provides two primary commands for local detection of flaky tests without involving Jira: `find` and `run`.

### `flakeguard find`

The `find` command scans your Go project to determine which test packages may be impacted by recent file changes. It conducts a dependency analysis to see which test packages depend on the changed files, helping you focus your testing efforts on the most relevant areas.


### `flakeguard run`

After identifying packages of interest (via `flakeguard find` or otherwise), use the `run` command to execute tests multiple times to detect flakiness.


Configuration Files (for Jira Integration)
------------------------------------------

When using Flakeguard’s Jira integration and local database features, you may need these configuration files:

1.  ### Local Database (`flaky_tests_db.json `)
    
    *   **Purpose:** Stores mappings between test definitions (package + name) and their Jira tickets, along with the assignee ID.
    *   **Location:** By default, stored as `flaky_tests_db.json ` in your user home directory (~). This path is determined by `localdb.DefaultDBPath()`.
    *   **Override:** Use `--test-db-path` with any command to specify a different location.
    *   **Format:** Internal JSON structure managed by the tool (_localdb/localdb.go_).
2.  ### User Mapping (`user_mapping.json`)
    
    *   **Purpose:** Maps Jira User Account IDs to Pillar Names or team/group identifiers. Used by:
        *   `create-tickets`: To set the Pillar Name custom field when creating a new Jira ticket if an assignee is determined.
        *   `review-tickets`: To allow setting the Pillar Name (\[i\] action) based on the ticket's assignee stored in the local DB.
    *   **Format:** A JSON array of objects.
    
    *   **Flag:** Use `--user-mapping-path` (default: `user_mapping.json`) when needed.
3.  ### User Test Mapping (`user_test_mapping.json`)
    
    *   **Purpose:** Defines regex patterns to match test packages (or full test paths) to suggest an assignee (Jira User Account ID) for new tickets created by `create-tickets`. The first pattern match determines the suggested assignee.
    *   **Format:** A JSON array of objects. The `.*` pattern can be used as a fallback.
    
    *   **Flag:** Use `--user-test-mapping-path` (default: `user_test_mapping.json`) when creating tickets.


Usage (Jira Integration)
------------------------

If you only need to identify flaky tests locally (via `find` and `run`) and do not intend to create or manage Jira tickets, you can skip these commands.

### `flakeguard create-tickets`

Interactively process a CSV file of flaky tests, suggest assignees based on patterns, check for existing Jira tickets, and create new tickets if needed.

*   **Key Features:**
    *   Reads test details from CSV.
    *   Suggests assignee using `user_test_mapping.json`.
    *   Checks local DB and Jira for existing tickets.
    *   Sets Pillar Name from `user_mapping.json` when creating new tickets.
    *   Provides a text-based UI for review and confirmation.
    *   Updates the local DB with created/manually assigned ticket keys.
    *   Outputs tests that were not confirmed for ticket creation to a _\_remaining.csv_ file.
*   **Flags:**
    *   `--csv-path` (required): Path to the input CSV file of flaky tests.
    *   `--jira-project`: Jira project key for new tickets (required if `JIRA_PROJECT_KEY` is not set).
    *   `--jira-issue-type`: Jira issue type for new tickets (default: _Task_).
    *   `--jira-search-label`: Label used to search for existing tickets (default: _flaky\_test_).
    *   `--test-db-path`: Local database file path (default: `~/flaky_tests_db.json `).
    *   `--user-mapping-path`: Path to user mapping JSON (default: `user_mapping.json`).
    *   `--user-test-mapping-path`: Path to user-test mapping JSON (default: `user_test_mapping.json`).
    *   `--skip-existing`: If set, automatically skips tests that already have a known Jira key.
    *   `--dry-run`: Simulates actions without creating Jira tickets or modifying the DB.

**Example:**

```
 go run main.go create-tickets \                      
  --jira-project=DX \
  --test-db-path=flaky_test_db.json \
  --user-mapping-path=user_mapping.json \
  --user-test-mapping-path=user_test_mapping.json \
  --skip-existing \
  --dry-run=false \
  --csv-path '1744018947_31163.csv'
```

### `flakeguard review-tickets`

Interactively review tickets in the local database. Fetches current status and Pillar Name from Jira, and lets you set the Pillar Name in Jira.

*   **Key Features:**
    *   Reads test data from the local DB.
    *   Fetches current Status and Pillar Name from Jira.
    *   Displays Test Info, Assignee, Pillar Name, and Jira Status.
    *   Allows setting the Pillar Name in Jira using the **\[i\]** action, based on `user_mapping.json`.
*   **TUI Actions:**
    *   **\[i\]** Set Pillar Name in Jira if it’s empty (and if prerequisites are met).
    *   **\[p\]** View previous ticket.
    *   **\[n\]** View next ticket.
    *   **\[q\]** Quit the TUI.
*   **Flags:**
    *   `--test-db-path`: Local DB path (default: `~/flaky_tests_db.json `).
    *   `--user-mapping-path`: Path to user mapping JSON (default: `user_mapping.json`).
    *   `--user-test-mapping-path`: Path to user-test mapping JSON (default: `user_test_mapping.json`).
    *   `--missing-pillars`: Only show tickets missing a Pillar Name.
    *   `--dry-run`: Prevents the **\[i\]** action from actually updating Jira.

**Examples:**

```
go run main.go review-tickets --test-db-path ".flaky_test_db.json" --dry-run=false --user-mapping-path "user_mapping.json"
```

### `flakeguard sync-jira`

Scans Jira for all tickets matching a specific label and ensures they exist in the local database. Updates local assignees based on Jira's current assignee.

*   **Key Features:**
    *   Fetches all Jira tickets
    *   Adds missing tickets to the local DB (where the test package/name might be unknown).
    *   Updates the _AssigneeID_ in the local DB if it differs from the current Jira assignee.
*   **Flags:**
    *   `--test-db-path`: Local DB file path (default: `~/flaky_tests_db.json `).
    *   `--jira-search-label`: Label used to find relevant tickets (default: _flaky\_test_).
    *   `--dry-run`: Shows what would be updated without modifying the DB.

**Example:**

```
go run main.go sync-jira --test-db-path=.flaky_test_db.json
```