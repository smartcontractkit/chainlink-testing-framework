.PHONY: test_unit
test_unit:
	go list ./... | grep -v 'example_test_package' | xargs go test -timeout 5m -json -cover -covermode=count -coverprofile=unit-test-coverage.out 2>&1 | tee /tmp/gotest.log | gotestloghelper -ci

.PHONY: test
test:
	go install github.com/gotesttools/gotestfmt/v2/cmd/gotestfmt@latest
	set -euo pipefail
	go list ./... | grep -v 'example_test_package' | xargs go test -json -cover -coverprofile unit-test-coverage.out -v 2>&1 | tee /tmp/gotest.log | gotestfmt

.PHONY: test-package
test-package:
	go install github.com/gotesttools/gotestfmt/v2/cmd/gotestfmt@latest
	set -euo pipefail
	go test -json -cover -coverprofile unit-test-coverage.out -v $(PKG) 2>&1 | tee /tmp/gotest.log | gotestfmt

.PHONY: test-race
test-race:
	go list ./... | grep -v 'example_test_package' | xargs go test -count=1 -race

.PHONY: bench
bench:
	go test -run=^$$ -bench=. -benchmem ./...

.PHONY: example
example:
	rm -rf example_results
	mkdir -p example_results
	- go run . run --project-path=./runner --test-packages=./example_test_package --run-count=5 --skip-tests=Panic,Timeout --max-pass-ratio=1 --race=false --output-json=example_results/example_run_1.json
	- go run . run --project-path=./runner --test-packages=./example_test_package --run-count=5 --skip-tests=Panic,Timeout --max-pass-ratio=1 --race=false --output-json=example_results/example_run_2.json
	- go run . run --project-path=./runner --test-packages=./example_test_package --run-count=5 --skip-tests=Panic,Timeout --max-pass-ratio=1 --race=false --output-json=example_results/example_run_3.json
	go run . aggregate-results \
      --results-path ./example_results \
      --output-path ./example_results \
      --repo-url "https://github.com/smartcontractkit/chainlink-testing-framework" \
      --base-sha "abc" \
      --head-sha "xyz" \
      --github-workflow-name "example"

.PHONY: example_flaky_panic
example_flaky_panic:
	rm -rf example_results
	mkdir -p example_results
	- go run . run --project-path=./runner --test-packages=./example_test_package --run-count=5 --skip-tests=TestPanic --max-pass-ratio=1 --race=false --output-json=example_results/example_run_1.json
	- go run . run --project-path=./runner --test-packages=./example_test_package --run-count=5 --skip-tests=TestPanic --max-pass-ratio=1 --race=false --output-json=example_results/example_run_2.json
	- go run . run --project-path=./runner --test-packages=./example_test_package --run-count=5 --skip-tests=TestPanic --max-pass-ratio=1 --race=false --output-json=example_results/example_run_3.json
	go run . aggregate-results \
	  --results-path ./example_results \
	  --output-path ./example_results \
	  --repo-url "https://github.com/smartcontractkit/chainlink-testing-framework" \
      --base-sha "abc" \
      --head-sha "xyz" \
      --github-workflow-name "example"

.PHONY: example_timeout
example_timeout:
	rm -rf example_results
	mkdir -p example_results
	- go run . run --project-path=./runner --test-packages=./example_test_package --run-count=5 --select-tests=TestTimeout --timeout=1s --max-pass-ratio=1 --race=false --output-json=example_results/example_run_1.json
	- go run . run --project-path=./runner --test-packages=./example_test_package --run-count=5 --select-tests=TestTimeout --timeout=1s --max-pass-ratio=1 --race=false --output-json=example_results/example_run_2.json
	- go run . run --project-path=./runner --test-packages=./example_test_package --run-count=5 --select-tests=TestTimeout --timeout=1s --max-pass-ratio=1 --race=false --output-json=example_results/example_run_3.json
	go run . aggregate-results \
	  --results-path ./example_results \
	  --output-path ./example_results \
	  --repo-url "https://github.com/smartcontractkit/chainlink-testing-framework" \
      --base-sha "abc" \
      --head-sha "xyz" \
      --github-workflow-name "example"