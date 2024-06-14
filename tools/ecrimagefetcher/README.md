# AWS ECR Image Fetcher

This Go script fetches the latest image tags from an AWS ECR repository based on specified criteria such as a grep string and optional semantic version constraints. It sorts and filters the tags and returns a specified number of the latest ones (according to semantic versioning).

## Usage

```bash
go run main.go <repository_name> <grep_string> <count> [<semver_constraints>]
```

### Example

```bash
go run main.go 'my-repo' '^v[0-9]+\.[0-9]+\.[0-9]+$' 5 '>=1.0.0, <2.0.0'
```

## Command Line Arguments

- `<repository_name>`: The name of the ECR repository.
- `<grep_string>`: A regex string to filter the tags.
- `<count>`: The number of latest tags to return.
- `[<semver_constraints>]`: Optional semantic version constraints to filter tags further.

## Output

- The script prints the specified number of latest image tags from the repository that match the given criteria, formatted as `repository_name:tag`, e.g.`hyperledger/besu:24.3,hyperledger/besu:24.2`.
- If the number of matching tags is less than the requested count, a warning is printed to stderr, and the available tags are returned.

## Error Handling

The script will panic and display error messages in the following scenarios:

- Insufficient command line arguments.
- Empty or invalid `repository_name`, `grep_string`, or `count`.
- Invalid regex for `grep_string`.
- Invalid integer value for `count`.
- Invalid semantic version constraints.
- Errors in fetching or parsing the image details from AWS ECR.

## Detailed Steps

1. **Argument Parsing and Validation**:
   - The script checks that at least 4 command line arguments are provided.
   - Validates the format and values of `repository_name`, `grep_string`, and `count`.
   - Optionally validates semantic version constraints if provided.
2. **Fetch Image Details**:
   - Constructs and executes an AWS CLI command to describe images in the specified ECR repository.
   - Parses the JSON output to extract image tags.
3. **Filter and Sort Tags**:
   - Filters tags based on the `grep_string` regex.
   - Optionally filters tags based on semantic version constraints.
   - Sorts the tags in descending order.
4. **Return Latest Tags**:
   - Constructs the output of the latest tags formatted as `repository_name:tag`.
   - Prints the result to stdout.

## Notes

- Ensure AWS CLI is configured and you have the necessary permissions to access the ECR repository.
- The environment variable `AWS_REGION` must be set to the appropriate AWS region.
- Semantic version constraints are optional and can be used to further filter the tags.
