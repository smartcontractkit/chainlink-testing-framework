# ASCII Table

This Go script reads a JSON file containing parsed results, optionally extracts specific key values, and formats the data into a table. The output is saved to a specified file.

## Usage

```bash
go run main.go --jsonfile <path_to_json_file> [options]
```

### Example

```bash
go run main.go --jsonfile data.json --section "Some section" --namedKey "keyName" --outputFile results.txt
```

### Sample output

```bash
+-----------------+--------+
| Value           | Result |
+-----------------+--------+
| No section here | √      |
+-----------------+--------+
|       Some section       |
+-----------------+--------+
| With a section  | X      |
+-----------------+--------+
```

## Command Line Arguments

- `--jsonfile`: Path to the JSON input file (required).
- `--section`: Optional section name to categorize results in the output file.
- `--namedKey`: Optional key to extract specific data from the JSON file.
- `--outputFile`: Optional output file name to save the results (default: `output.txt`).
- `--firstColumn`: Header for the first column in the output table (default: `Value`).
- `--secondColumn`: Header for the second column in the output table (default: `Result`).

## Output

- The script formats the data from the JSON file into a table and saves it to the specified output file.
- If a `namedKey` is provided, only the corresponding data is included in the output.

## Error Handling

The script will panic and display error messages in the following scenarios:

- Missing required `--jsonfile` flag.
- Errors in reading the JSON file.
- Errors in parsing the JSON data.
- Errors in writing to the output file.

## Detailed Steps

1. **Argument Parsing and Validation**:
   - The script checks that the `--jsonfile` flag is provided.
   - Validates the presence of optional flags and sets default values if not provided.
2. **Read JSON File**:
   - Reads the content of the specified JSON file.
3. **Parse JSON Data**:
   - Parses the JSON data to extract results.
   - If `namedKey` is provided, extracts specific data from the JSON.
4. **Format Data**:
   - Constructs a table with headers for the first and second columns.
   - Adds the parsed results to the table.
5. **Write to Output File**:
   - Writes the formatted table to the specified output file.
   - Creates a new file if it doesn't exist or overwrites the existing file.

## Notes

- Ensure the JSON file path is correct and accessible.
- The JSON file should contain properly formatted data as expected by the script.
- Use the `--section` flag to categorize results in the output file for better organization.
