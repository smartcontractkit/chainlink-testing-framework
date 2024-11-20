.PHONY: test_unit
test_unit:
	go list ./... | grep -v 'flaky_test_package' | xargs go test -timeout 5m -json -cover -covermode=count -coverprofile=unit-test-coverage.out 2>&1 | tee /tmp/gotest.log | gotestloghelper -ci

.PHONY: test
test:
	go install github.com/gotesttools/gotestfmt/v2/cmd/gotestfmt@latest
	set -euo pipefail
	go list ./... | grep -v 'flaky_test_package' | xargs go test -json -cover -coverprofile unit-test-coverage.out -v 2>&1 | tee /tmp/gotest.log | gotestfmt
