.PHONY: test_unit
test_unit:
	go list ./... | grep -v 'example_test_package' | xargs go test -timeout 5m -json -cover -covermode=count -coverprofile=unit-test-coverage.out 2>&1 | tee /tmp/gotest.log | gotestloghelper -ci

.PHONY: test
test:
	go install github.com/gotesttools/gotestfmt/v2/cmd/gotestfmt@latest
	set -euo pipefail
	go list ./... | grep -v 'example_test_package' | xargs go test -json -cover -coverprofile unit-test-coverage.out -v 2>&1 | tee /tmp/gotest.log | gotestfmt

.PHONY: example
example:
	go run . run --project-path=./runner --test-packages=./example_test_package --run-count=5 --skip-tests=TestPanic --max-pass-ratio=1 --race=false --output-json=example_results.json