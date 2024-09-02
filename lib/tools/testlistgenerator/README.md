# Test List Generator

This Go script builds a JSON file containing tests to be run for a given product and Ethereum implementation. It reads command line arguments to construct and append entries to a JSON file.

## Usage

```bash
go run main.go -t <test_name> -o <output_file_name> -p <product> -r <test_regex> -f <file> -e <eth_implementation> -d <docker_images> -n <node_label> [-c <chain_ids>] [-w <networks>]
```

### Example

```bash
go run main.go -t "emv-test" -o "test_list.json" -p "ocr" -r "TestOCR.*" -f "./smoke/ocr_test.go" -e "besu" -d "hyperledger/besu:21.0.0,hyperledger/besu:22.0.0" -n "ubuntu-latest"
```

## Output

The script generates or updates a JSON file with entries structured as follows:

```json
[
  {
    "name": "emv-test-01",
    "os": "ubuntu-latest",
    "product": "ocr",
    "eth_implementation": "besu",
    "docker_image": "hyperledger/besu:21.0.0",
    "run": "-run 'TestOCR.*' ./smoke/ocr_test.go"
  },
  {
    "name": "emv-test-02",
    "os": "ubuntu-latest",
    "product": "ocr",
    "eth_implementation": "besu",
    "docker_image": "hyperledger/besu:22.0.0",
    "run": "-run 'TestOCR.*' ./smoke/ocr_test.go"
  }
]
```

If the script is run with optional `chain_id` flag the output is slightly different:

```bash
go run main.go -t "emv-test" -o "test_list.json" -p "ocr" -r "TestOCR.*" -f "./smoke/ocr_test.go" -e "besu" -d "hyperledger/besu:21.0.0,hyperledger/besu:22.0.0" -n "ubuntu-latest" -c 1337,2337 -w "mainnet,ropsten"
```

Output:

```json
[
  {
    "name": "emv-test-01",
    "os": "ubuntu-latest",
    "product": "ocr",
    "eth_implementation": "besu",
    "docker_image": "1337=hyperledger/besu:21.0.0,2337=hyperledger/besu:21.0.0",
    "run": "-run 'TestOCR.*' ./smoke/ocr_test.go",
    "networks": "mainnet,ropsten"
  },
  {
    "name": "emv-test-02",
    "os": "ubuntu-latest",
    "product": "ocr",
    "eth_implementation": "besu",
    "docker_image": "1337=hyperledger/besu:22.0.0,2337=hyperledger/besu:22.0.0",
    "run": "-run 'TestOCR.*' ./smoke/ocr_test.go",
    "networks": "mainnet,ropsten"
  }
]
```

## Command Line Arguments

- `-t <test_name>`: A prefix for the test name.
- `-o <output_file_name>`: The name of the JSON file where the test entries will be stored.
- `-p <product>`: The name of the product for which the tests are being generated.
- `-r <test_regex>`: The regular expression to match test names.
- `-f <file>`: The file path where the tests are defined.
- `-e <eth_implementation>`: The name of the Ethereum implementation.
- `-d <docker_images>`: A comma-separated list of Docker images to be used.
- `-n <node_label>`: The node label for the test environment.
- `-c <chain_ids>`: (Optional) A comma-separated list of chain IDs to associate with each Docker image.
- `-w <networks>`: (Optional) A comma-separated list of networks.

## Error Handling

The script will panic and display error messages in the following scenarios:

- Insufficient command line arguments.
- Empty parameters for test_name, output_file_name, product, test_regex, file, eth_implementation, docker_images, or node_label.
- Invalid Docker image format (should include a version tag).
- Invalid regular expression for test_regex.
- Invalid or non-integer chain IDs.
- Errors in file operations (opening, reading, writing).

## Detailed Steps

1. **Argument Parsing**: The script uses Cobra to parse command line arguments.
2. **Validation**: It validates the input parameters, ensuring none are empty, and the regular expression compiles.
3. **File Operations**:
   - If the output file exists, it reads and unmarshals the content.
   - If it doesn't exist, it creates a new file.
4. **Appending Entries**: For each Docker image, it creates new `OutputEntry` objects and appends them to the output.
5. **JSON Marshaling**: It marshals the updated output to JSON and writes it back to the file.
6. **Completion Message**: Prints a success message indicating the number of tests added.

## Notes

- Ensure the Docker image names include version tags, e.g., `hyperledger/besu:21.0.0`.
- The script appends new entries; it does not overwrite existing entries in the output file unless the output file is specified with new test entries.
- Optional parameters like `chain_ids` and `networks` allow for additional customization of the test entries. If `chain_ids` are provided, they will be included in the Docker image field. If `networks` are provided, they will be included as an additional field in the output JSON.

This update reflects the recent changes to include flags for `chain_ids` and `networks`, and ensures the documentation matches the new argument names and structure.
