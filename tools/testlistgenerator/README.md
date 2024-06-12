# Test List Generator

This Go script builds a JSON file containing tests to be run for a given product and Ethereum implementation. It reads command line arguments to construct and append entries to a JSON file.

## Usage

```bash
go run main.go <output_file_name> <product> <test_regex> <file> <eth_implementation> <docker_images>
```

### Example

```bash
go run main.go 'test_list.json' 'ocr' 'TestOCR.*' './smoke/ocr_test.go' 'besu' 'hyperledger/besu:21.0.0,hyperledger/besu:22.0.0'
```

## Output

The script generates or updates a JSON file with entries structured as follows:

```json
{
  "tests": [
    {
      "product": "ocr",
      "test_regex": "TestOCR.*",
      "file": "./smoke/ocr_test.go",
      "eth_implementation": "besu",
      "docker_image": "hyperledger/besu:21.0.0"
    },
    {
      "product": "ocr",
      "test_regex": "TestOCR.*",
      "file": "./smoke/ocr_test.go",
      "eth_implementation": "besu",
      "docker_image": "hyperledger/besu:22.0.0"
    }
  ]
}
```

## Command Line Arguments

- `<output_file_name>`: The name of the JSON file where the test entries will be stored.
- `<product>`: The name of the product for which the tests are being generated.
- `<test_regex>`: The regular expression to match test names.
- `<file>`: The file path where the tests are defined.
- `<eth_implementation>`: The name of the Ethereum implementation.
- `<docker_images>`: A comma-separated list of Docker images to be used.

## Error Handling

The script will panic and display error messages in the following scenarios:

- Insufficient command line arguments.
- Empty parameters for output_file_name, product, test_regex, file, eth_implementation, or docker_images.
- Invalid Docker image format (should include a version tag).
- Invalid regular expression for test_regex.
- Errors in file operations (opening, reading, writing).

## Detailed Steps

1. **Argument Parsing**: The script expects at least 7 command line arguments. It splits the `<docker_images>` argument into a slice.
2. **Validation**: It validates the input parameters, ensuring none are empty and the regular expression compiles.
3. **File Operations**:
   - If the output file exists, it reads and unmarshals the content.
   - If it doesn't exist, it creates a new file.
4. **Appending Entries**: For each Docker image, it creates a new `OutputEntry` and appends it to the output.
5. **JSON Marshaling**: It marshals the updated output to JSON and writes it back to the file.
6. **Completion Message**: Prints a success message indicating the number of tests added.

## Notes

- Ensure the Docker image names include version tags, e.g., `hyperledger/besu:21.0.0`.
- The script appends new entries; it does not overwrite existing entries in the output file.
