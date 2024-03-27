#!/bin/bash

set -e

SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &>/dev/null && pwd)

cd "$SCRIPT_DIR"/../ || exit 1

# Check if any arguments are provided
if [ $# -eq 0 ]; then
    echo "Usage: $0 <tool path>, example: tools/gotesthelper"
    exit 1
fi

# The first argument is the tool
TOOL="$1"

# Function to check if the image exists in ECR
push_from_package_version() {
    local tool="$1"
    local version package tagexists

    version="v$(cmd < ./"${tool}"/package.json | jq -r '.version')"
    package="${tool}/${version}"
    tagexists=$(git tag -l "${package}")
    if [ -z "${tagexists}" ]; then
        git tag "${package}"
        git push "${package}"
    else
        echo "Tag ${package} already exists."
    fi
}

push_from_package_version "${TOOL}"
