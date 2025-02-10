#!/bin/bash
# Used to run flakeguard on the example test package

run_flakeguard_example() {
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

  go run . aggregate-results \
		--results-path ./example_results \
		--output-path ./example_results \
		--repo-url "https://github.com/smartcontractkit/chainlink-testing-framework" \
		--branch-name "example-branch" \
		--base-sha "abc" \
		--head-sha "xyz" \
		--github-workflow-name "ExampleWorkflowName" \
		--github-workflow-run-url "https://github.com/example/repo/actions/runs/1" \
		--splunk-url "https://splunk.example.com" \
		--splunk-token "splunk-token" \
		--splunk-event "example-splunk-event"
  local EXIT_CODE=$?
  if [ $EXIT_CODE -eq 2 ]; then
    echo "ERROR: Flakeguard encountered an error while aggregating results"
    exit $EXIT_CODE
  fi

  GITHUB_TOKEN="EXAMPLE_GITHUB_TOKEN" go run . generate-report \
    --aggregated-results-path ./example_results/all-test-results.json \
    --output-path ./example_results \
    --github-repository "smartcontractkit/chainlink-testing-framework" \
    --github-run-id "1" \
    --failed-tests-artifact-name "failed-test-results-with-logs.json" \
    --base-branch "exampleBaseBranch" \
    --current-branch "exampleCurrentBranch" \
    --current-commit-sha "abc" \
    --repo-url "https://github.com/smartcontractkit/chainlink-testing-framework" \
    --action-run-id "1" \
    --max-pass-ratio "1.0"
  local EXIT_CODE=$?
  if [ $EXIT_CODE -eq 2 ]; then
    echo "ERROR: Flakeguard encountered an error while generating report for the aggregated results"
    exit $EXIT_CODE
  fi
}

# Run the commands
rm -rf example_results
mkdir -p example_results

run_flakeguard_example "$1" "$2"