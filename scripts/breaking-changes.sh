#!/usr/bin/env bash
cd tools/breakingchanges/
go run cmd/main.go -path ../.. || exit_code=$?
cd -
exit ${exit_code:-0}