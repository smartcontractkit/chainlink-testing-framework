# gotestloghelper CLI

gotestloghelper CLI is a command-line interface tool designed to enhance the output of Go test runs. It provides features such as colorized output for better readability, hiding passing logs or passing tests altogether, and removing log prefix added by testing.T.Log, especially useful in Continuous Integration (CI) environments.

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
    -hidepassingtests: Set to true to hide passing tests, only compatible when used with -json.
    -hidepassinglogs: Set to true to hide passing test logs but not the tests themselves, only compatible when used with -json. Default: false
    -color: Set to true to enable color output. Default: false
    -ci: Set to true to enable CI mode, which prints out logs with groupings when combined with -json. Default: false
    -singlepackage: Set to true if the Go test output is from a single package only. This prints tests out as they finish instead of waiting for the package to finish. Default: false
    -errorattoplength: If the error message doesn't appear before this many lines, it will be printed at the top of the test output as well. Set to 0 to disable. Only works with -ci. Default: 100

    Deprecated:
    -onlyerrors: Now this sets -hidepassingtests to true. Set to true to only print tests that failed. Note: Only compatible with -json. Default: false

## Examples

To run gotestloghelper CLI with color output:

```sh
go test ./... -json | gotestloghelper -json -color
```

To filter only errors from JSON-formatted test output:

```sh
go test -json ./... | gotestloghelper -json -hidepassingtests
```

To filter passing test logs but still show all tests and failing test logs from JSON-formatted test output:

```sh
go test -json ./... | gotestloghelper -json -hidepassinglogs
```

## Additional Notes

- Interrupting the CLI (Ctrl+C) will cancel the current operation and print "Cancelling... interrupt again to exit". A second interrupt will exit the CLI immediately.
- Ensure that your Go test commands are compatible with the flags you use with gotestloghelper CLI, for example always use `-json` flags together.
