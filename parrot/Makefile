# Default test log level (can be overridden)
PARROT_TEST_LOG_LEVEL ?= ""

# Pass TEST_LOG_LEVEL as a flag to go test
TEST_ARGS ?= -testLogLevel=$(PARROT_TEST_LOG_LEVEL)

.PHONY: lint
lint:
	golangci-lint --color=always run ./... --fix -v

.PHONY: test
test:
	go install github.com/gotesttools/gotestfmt/v2/cmd/gotestfmt@latest
	set -euo pipefail
	go test $(TEST_ARGS) -json -cover -coverprofile cover.out -v ./... 2>&1 | tee /tmp/gotest.log | gotestfmt

.PHONY: test_race
test_race:
	go install github.com/gotesttools/gotestfmt/v2/cmd/gotestfmt@latest
	set -euo pipefail
	go test $(TEST_ARGS) -json -cover -count=1 -race -coverprofile cover.out -v ./... 2>&1 | tee /tmp/gotest.log | gotestfmt

.PHONY: test_unit
test_unit:
	go test $(TEST_ARGS) -coverprofile cover.out ./...

.PHONY: bench
bench:
	go test $(TEST_ARGS) -bench=. -run=^$$ ./...

.PHONY: build
build:
	go build -o ./parrot ./cmd