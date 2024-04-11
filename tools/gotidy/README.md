# gotidy CLI

gotidy CLI is a tool designed to check if you have your go project mod files tidy'd. It will also commit the tidy changes if specified to.

## Installation

To install gotidy CLI, you need to have Go installed on your machine. With Go installed, run the following command:

```sh
go install github.com/smartcontractkit/chainlink-testing-framework/tools/gotidy@latest
```

## Usage

After installation, you can run gotidy CLI using the following syntax:

```sh
go tidy -path="./" -commit=true
```

## Available Flags

    -path: Path to the go project to check for tidy. Default: .
    -commit: Commit the changes if there are any. Default: false
    -onlyerrors: Set to true to only print tests that failed. Note: Not compatible without -json. Default: false

## Examples

To run gotidy CLI to checky tidyness:

```sh
gotidy
```

To check a subproject and commit changes if any:

```sh
gotidy -path="./sub/project" -commit=true
```

## Additional Notes

- Interrupting the CLI (Ctrl+C) will cancel the current operation and print "Cancelling... interrupt again to exit". A second interrupt will exit the CLI immediately.
