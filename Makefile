BIN_DIR = bin
export GOPATH ?= $(shell go env GOPATH)
export GO111MODULE ?= on

.PHONY: lint
lint: ## run linter
	${BIN_DIR}/golangci-lint --color=always run ./... -v

.PHONY: golangci
golangci: ## install golangci-linter
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ${BIN_DIR} v1.41.1

.PHONY: gomod
gomod: ## install go modules
	go mod download

install-deps: gomod golangci ## install necessary dependencies

.PHONY: install_cli
install_cli: ## install framework CLI
	go install cli/ifcli.go

.PHONY: test
test: ## run tests
	go test -v ./suite/contracts ./client -count 1 -p 1

.PHONY: test
test_performance: ## run performance tests
	go test -v ./suite/performance/... -count 1 -p 1 -timeout 100m

.PHONY: test_refill
test_refill: ## runs refill suite
	go test -v ./suite/refill -count 1 -p 1

.PHONY: test_race
test_race: ## run tests with race
	go test -v ./... -race -count 1 -p 1

.PHONY: test_nightly
test_nightly: ## run nightly tests
	go test -v ./... -race -count 20 -p 1

.PHONY: test_coverage
test_coverage: ## run tests with coverage
	go test ./client ./config -v -covermode=count -coverprofile=coverage.out

.PHONY: test_unit
test_unit: ## run unit tests
	ginkgo -r --focus=@unit

.PHONY: test_cron
test_cron: ## run cron tests
	ginkgo -r --focus=@cron

.PHONY: test_flux
test_flux: ## run flux tests
	ginkgo -r --focus=@flux

.PHONY: test_keeper
test_keeper: ## run keeper tests
	ginkgo -r --focus=@keeper

.PHONY: test_ocr
test_ocr: ## run ocr tests
	ginkgo -r --focus=@ocr

.PHONY: test_runlog
test_runlog: ## run runlog tests
	ginkgo -r --focus=@runlog

.PHONY: test_contract
test_contract: ## run contract tests
	ginkgo -r --focus=@contract

.PHONY: test_vrf
test_vrf: ## run vrf tests
	ginkgo -r --focus=@vrf