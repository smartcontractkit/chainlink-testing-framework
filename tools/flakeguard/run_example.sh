#!/bin/bash

run_flakeguard() {
  local run_count=$1
  local skip_tests=$2

  echo "Running Flakeguard $run_count times"

  for ((i=1; i<=$run_count; i++)); do
    output_file="example_results/example_run_$i.json"
    go run . run \
      --project-path=./runner \
      --test-packages=./example_test_package \
      --run-count=5 \
      --skip-tests=$skip_tests \
      --max-pass-ratio=1 \
      --race=false \
      --output-json=$output_file
    local EXIT_CODE=$?
    if [ $EXIT_CODE -eq 2 ]; then
      echo "ERROR: Flakeguard encountered an error while running tests"
      exit $EXIT_CODE
    fi
  done
}

# Run the commands
rm -rf example_results
mkdir -p example_results

run_flakeguard "$1" "$2"