BIN_DIR = bin
export GOPATH ?= $(shell go env GOPATH)
export GO111MODULE ?= on

.PHONY: lint
lint: ## run linter
	${BIN_DIR}/golangci-lint --color=always run ./... -v --timeout=3m

.PHONY: golangci
golangci: ## install golangci-linter
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ${BIN_DIR} v1.41.1

.PHONY: gomod
gomod: ## install go modules
	go mod download

install-deps: gomod golangci ## install necessary dependencies

.PHONY: test
test: ## run tests
	go test -v ./suite/contracts ./client -count 1 -p 1

.PHONY: test
test_performance: ## run performance tests
	go test -v ./suite/performance/... -count 1 -p 1 -timeout 100m

.PHONY: test_race
test_refill: ## runs refill suite
	go test -v ./suite/refill -count 1 -p 1

.PHONY: test_race
test_race: ## run tests with race
	go test -v ./... -race -count 1 -p 1

.PHONY: test_nightly
test_nightly: ## run nightly tests
	go test -v ./... -race -count 20 -p 1