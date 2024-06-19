# Workflow Result Parser

This Go script fetches jobs from a specified GitHub Actions workflow run, filters them based on a regex pattern, and outputs the results. The output can be saved to a specified file and optionally stored under a named key.

## Usage

```bash
go run main.go --githubToken <token> --githubRepo <repository> --workflowRunID <workflow_run_id> --jobNameRegex <regex_pattern> [options]
```

### Example

```bash
go run main.go --githubToken ghp_exampleToken --githubRepo owner/repo --workflowRunID 123456789 --jobNameRegex "Test-.*" --outputFile results.json
```

### Sample output

```json
{
  "results": [
    {
      "conclusion": ":white_check_mark:",
      "cap": "1",
      "html_url": "http://example.com/job1"
    },
    {
      "conclusion": ":x:",
      "cap": "2",
      "html_url": "http://example.com/job2"
    }
  ]
}
```

## Command Line Arguments

- `--githubToken`: GitHub token for authentication (required).
- `--githubRepo`: GitHub repository in the format `owner/repo` (required).
- `--workflowRunID`: ID of the GitHub Actions workflow run (required).
- `--jobNameRegex`: Regex pattern to match job names (required).
- `--namedKey`: Optional named key under which results will be stored.
- `--outputFile`: Optional output file name to save the results.

## Output

- The script fetches jobs from the specified GitHub Actions workflow run.
- Filters the jobs based on the provided regex pattern.
- Formats and saves the results to the specified output file.
- Prints the results to the console if no output file is specified.

## Error Handling

The script will panic and display error messages in the following scenarios:

- Missing required flags: `--githubToken`, `--githubRepo`, `--workflowRunID`, `--jobNameRegex`.
- Errors in making HTTP requests to the GitHub API.
- Errors in reading or writing files.
- Errors in parsing JSON data.

## Detailed Steps

1. **Argument Parsing and Validation**:
   - The script checks that all required flags are provided.
   - Validates optional flags and sets default values if not provided.
2. **Fetch GitHub Jobs**:
   - Constructs the API URL for fetching jobs from the specified workflow run.
   - Makes an HTTP GET request to the GitHub API with the provided token.
   - Handles pagination to fetch all jobs if necessary.
3. **Parse and Filter Jobs**:
   - Parses the JSON response from the GitHub API.
   - Filters the jobs based on the provided regex pattern.
4. **Process and Save Results**:
   - Formats the filtered results.
   - Saves the results to the specified output file.
   - Prints the results to the console if no output file is specified.

## Notes

- Ensure the GitHub token has the necessary permissions to access the repository and workflow runs.
- Use the `--namedKey` flag to categorize results in the output file for better organization.
- The GitHub API has rate limits; be mindful of making too many requests in a short period.
