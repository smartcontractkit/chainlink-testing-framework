# ghsecrets 

ghsecrets is a command-line tool designed to manage and set test secrets in GitHub via the GitHub CLI. 

## Installation

To install ghsecrets CLI, you need to have Go installed on your machine. With Go installed, run the following command:

```sh
go install github.com/smartcontractkit/chainlink-testing-framework/tools/ghsecrets@latest
```

Please install GitHub CLI to use this tool - https://cli.github.com/

## Usage

Set default test secrets from ~/.testsecrets file:
```sh
ghsecrets set
```

## FAQ

### Q: What should I do if I get "command not found: ghsecrets" after installation?

This error typically means that the directory where Go installs its binaries is not included in your system's PATH. The binaries are usually installed in $GOPATH/bin or $GOBIN. Here's how you can resolve this issue:

1. Add Go bin directory to PATH:

- First, find out where your Go bin directory is by running:
    ```sh
    echo $(go env GOPATH)/bin
    ```
    This command will print the path where Go binaries are installed, typically something like /home/username/go/bin

- Add the following line at the end of the file:
    ```sh
    export PATH="$PATH:<path-to-go-bin>"
    ```

-  Apply the changes by sourcing the file:   
    ```sh
    source ~/.bashrc  # Use the appropriate file like .zshrc if needed
    ```

2. Alternatively, run using the full path:

    If you prefer not to alter your PATH, or if you are troubleshooting temporarily, you can run the tool directly using its full path:
    ```sh
    $(go env GOPATH)/bin/ghsecrets set
    ```
    Ensure you have correctly installed ghsecrets by checking its presence in the Go bin directory ($GOPATH/bin or $GOBIN). If the problem persists, you might want to reinstall the tool and watch for any errors during the installation process.