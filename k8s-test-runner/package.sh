#!/bin/bash

# Check if version argument is provided
if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <version>"
    exit 1
fi

# Capture the version from the script arguments
VERSION=$1

# Navigate to the cmd directory
cd cmd/

# Build the binary for Darwin ARM64 and append the version to its name
GOOS=darwin GOARCH=arm64 go build -o k8s-test-runner-v$VERSION-darwin-arm64

# Build the binary for Linux AMD64 and append the version to its name
GOOS=linux GOARCH=amd64 go build -o k8s-test-runner-v$VERSION-linux-amd64

# Navigate back to the root directory
cd ../

# Package the Dockerfile, Helm chart, and both binaries into a tarball, appending the version to the tarball's name
tar -czvf k8s-test-runner-v$VERSION.tar.gz Dockerfile.testbin chart cmd/k8s-test-runner-v$VERSION-darwin-arm64 cmd/k8s-test-runner-v$VERSION-linux-amd64
echo "Created k8s-test-runner-v$VERSION.tar.gz"
