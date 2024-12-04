#!/bin/bash

go test -v -race -coverprofile cover.out -count 1 `go list ./... | grep -v examples | grep benchspy` -run TestBenchSpy
coverage=$(go tool cover -func=cover.out | grep total | awk '{print substr($3, 1, length($3)-1)}')
if [ -z "$coverage" ]; then
    echo "Error: Could not determine test coverage";
    exit 1
fi

if [[ $(echo "$coverage < 85" | bc -l) -eq 1 ]]; then
    echo "Test coverage $coverage% is below minimum 85%"
    exit 1
fi
echo "Test coverage: $coverage%"