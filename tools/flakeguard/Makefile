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
	./run_example.sh 3 Panic,Timeout

.PHONY: example_flaky_panic
example_flaky_panic:
	./run_example.sh 3 TestPanic
