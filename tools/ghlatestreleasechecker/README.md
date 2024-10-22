# GitHub Latest Release Checker

This Go script checks if the latest release of a given GitHub repository is within a specified number of days from today. It fetches the latest release information from the GitHub API and compares the release date with the current date.

## Usage

```bash
go run main.go <repository_name> <days>
```

### Example

```bash
go run main.go 'owner/repo' 30
```

## Command Line Arguments

- `<repository_name>`: The GitHub repository in the format `owner/repo`.
- `<days>`: The number of days to check if the latest release is within.

## Output

- The script prints the tag name of the latest release if it was published within the specified number of days.
- If the release is older than the specified number of days, it prints `none`.

## Error Handling

The script will panic and display error messages in the following scenarios:

- Insufficient command line arguments.
- Invalid repository name format (must be `owner/repo`).
- Non-integer value for the `days` argument.
- Errors in fetching or parsing the latest release information from the GitHub API.
- Unexpected status codes from the GitHub API response.

## Detailed Steps

1. **Argument Parsing and Validation**:
   - The script checks that at least 3 command line arguments are provided.
   - Validates the format of the repository name.
   - Ensures the `days` argument is an integer.
2. **Fetch Latest Release**:
   - Constructs the API URL for the latest release of the specified repository.
   - Makes an HTTP GET request to the GitHub API.
   - Parses the JSON response to extract the latest release information.
3. **Check Release Date**:
   - Compares the release date with the current date to determine if it is within the specified number of days.
4. **Output Result**:
   - Prints the tag name of the latest release if it is recent.
   - Prints `none` if the release is older than the specified number of days.

## Notes

- Ensure the repository name is in the format `owner/repo`.
- The GitHub API has rate limits; be mindful of making too many requests in a short period.
