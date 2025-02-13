# Default test log level is none
PARROT_TEST_LOG_LEVEL ?= ""
TEST_ARGS ?= -testLogLevel=$(PARROT_TEST_LOG_LEVEL)

TEST_TIMEOUT ?= 30s

.PHONY: lint
lint:
	golangci-lint --color=always run ./... --fix -v

.PHONY: test
test:
	go install github.com/gotesttools/gotestfmt/v2/cmd/gotestfmt@latest
	set -euo pipefail
	go test $(TEST_ARGS) -json -timeout $(TEST_TIMEOUT) -cover -coverprofile cover.out -v ./... 2>&1 | tee /tmp/gotest.log | gotestfmt

.PHONY: test_race
test_race:
	go install github.com/gotesttools/gotestfmt/v2/cmd/gotestfmt@latest
	set -euo pipefail
	go test $(TEST_ARGS) -json -timeout $(TEST_TIMEOUT) -cover -count=1 -race -coverprofile cover.out -v ./... 2>&1 | tee /tmp/gotest.log | gotestfmt


.PHONY: test_unit
test_unit:
	go test $(TEST_ARGS) -timeout $(TEST_TIMEOUT) -coverprofile cover.out ./...

.PHONY: bench
bench:
	go test $(TEST_ARGS) -bench=. -run=^$$ ./...

.PHONY: fuzz_tests
fuzz:
	go test -list="Fuzz" ./...

.PHONY: build
build:
	cd .. && goreleaser release --snapshot --clean -f ./parrot/.goreleaser.yaml
	echo "Build done, check in parrot/dist for binaries"
