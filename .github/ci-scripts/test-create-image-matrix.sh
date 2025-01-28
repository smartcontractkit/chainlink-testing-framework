#!/bin/bash
set -x

export CHANGED_FILES="README.md
images/solana-validator/Dockerfile
havoc/chaos_listener.go
havoc/.gitignore
havoc/utils.go"

export GITHUB_OUTPUT="/dev/stdout"

./.github/ci-scripts/create-image-matrix.sh
