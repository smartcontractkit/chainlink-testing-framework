.PHONY: test_unit
test_unit:
	go test ./... -timeout 5m -json -cover -covermode=count -coverprofile=unit-test-coverage.out 2>&1 | tee /tmp/gotest.log | gotestloghelper -ci

.PHONY: test
test:
	go install github.com/gotesttools/gotestfmt/v2/cmd/gotestfmt@latest
	set -euo pipefail
	go test ./... -json -cover -coverprofile unit-test-coverage.out -v 2>&1 | tee /tmp/gotest.log | gotestfmt

.PHONY: test-package
test-package:
	go install github.com/gotesttools/gotestfmt/v2/cmd/gotestfmt@latest
	set -euo pipefail
	go test -json -cover -coverprofile unit-test-coverage.out -v $(PKG) 2>&1 | tee /tmp/gotest.log | gotestfmt

.PHONY: test-race
test-race:
	go test ./... -count=1 -race

.PHONY: bench
bench:
	go test -run=^$$ -bench=. -benchmem ./...
