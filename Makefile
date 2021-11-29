BIN_DIR = bin
export GOPATH ?= $(shell go env GOPATH)
export GO111MODULE ?= on

lint:
	${BIN_DIR}/golangci-lint --color=always run ./... -v

golangci:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ${BIN_DIR} v1.42.0

go_mod:
	go mod download

install: go_mod golangci

test_unit:
	ginkgo -r --junit-report=tests-unit-report.xml --keep-going --trace --randomize-all --randomize-suites --progress -cover -covermode=count -coverprofile=unit-test-coverage.out -nodes=10 ./client ./config ./environment 

test_smoke:
	ginkgo -v -r --junit-report=tests-smoke-report.xml --keep-going --trace --randomize-all --randomize-suites --progress -nodes-15 ./suite/smoke

test_performance:
	NETWORKS=ethereum_geth_performance,ethereum_geth_performance ginkgo -r --junit-report=tests-performance-report.xml --keep-going --trace --randomize-all --randomize-suites --progress $(args) ./suite/performance 

test_chaos:
	NETWORKS=ethereum_geth_performance,ethereum_geth_performance ginkgo -r --junit-report=tests-chaos-report.xml --keep-going --trace --randomize-all --randomize-suites --progress $(args) ./suite/chaos 
