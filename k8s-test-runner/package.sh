#!/bin/bash

# Check if sufficient arguments are provided
if [ "$#" -ne 2 ]; then
    echo "Usage: $0 <version> <os(darwin|linux)>"
    exit 1
fi

# Capture the version and OS from the script arguments
VERSION=$1
OS=$2

# Validate the OS argument
if [ "$OS" != "darwin" ] && [ "$OS" != "linux" ]; then
    echo "Invalid OS. Please select either 'darwin' or 'linux'."
    exit 1
fi

# Navigate to the cmd directory
cd cmd/ || exit 1

# Initialize binary name
BINARY_NAME=""

# Build the binary based on the selected OS
if [ "$OS" == "darwin" ]; then
    GOOS=darwin GOARCH=arm64 go build -o "k8s-test-runner-$OS-arm64"
    BINARY_NAME="k8s-test-runner-$OS-arm64"
elif [ "$OS" == "linux" ]; then
    GOOS=linux GOARCH=amd64 go build -o "k8s-test-runner-$OS-amd64"
    BINARY_NAME="k8s-test-runner-$OS-amd64"
fi

# Navigate back to the root directory
cd ../

# Package the Dockerfile, Helm chart, and the selected binary into a tarball, appending the version and OS to the tarball's name
tar -czvf "k8s-test-runner-$OS-v$VERSION.tar.gz" Dockerfile.testbin chart -C cmd "$BINARY_NAME"
echo "Created k8s-test-runner-$OS-v$VERSION.tar.gz"
