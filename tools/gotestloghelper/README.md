# gotestloghelper CLI

gotestloghelper CLI is a command-line interface tool designed to enhance the output of Go test runs. It provides features such as colorized output for better readability, showing only errors, and removing log prefixed added by testing.T.Log, especially useful in Continuous Integration (CI) environments.

## Installation

To install gotestloghelper CLI, you need to have Go installed on your machine. With Go installed, run the following command:

```sh
go install github.com/smartcontractkit/chainlink-testing-framework/tools/gotestloghelper@latest
```

## Usage

After installation, you can run gotestloghelper CLI using the following syntax:

```sh
go test ./... -json | gotestloghelper [flags]
```

## Available Flags

    -tlogprefix: Set to true to remove the Go test log prefix. Default: false
    -json: Set to true to enable parsing the input from a go test -json output. Default: false
    -onlyerrors: Set to true to only print tests that failed. Note: Not compatible without -json. Default: false
    -color: Set to true to enable color output. Default: false
    -ci: Set to true to enable CI mode, which prints out logs with groupings when combined with -json. Default: false
    -singlepackage: Set to true if the Go test output is from a single package only. This prints tests out as they finish instead of waiting for the package to finish. Default: false

## Examples

To run gotestloghelper CLI with color output:

```sh
go test ./... -json | gotestloghelper -json -color
```

To filter only errors from JSON-formatted test output:

```sh
go test -json ./... | gotestloghelper -json -onlyerrors
```

## Additional Notes

- Interrupting the CLI (Ctrl+C) will cancel the current operation and print "Cancelling... interrupt again to exit". A second interrupt will exit the CLI immediately.
- Ensure that your Go test commands are compatible with the flags you use with gotestloghelper CLI, for example always use `-json` flags together.
