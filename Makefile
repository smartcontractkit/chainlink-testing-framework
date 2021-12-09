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
	NETWORKS=alfajores_celo ginkgo -r -keepGoing --trace --randomizeAllSpecs --randomizeSuites --progress -covermode=count -coverprofile=unit-test-coverage.out -nodes=10 ./client ./config ./environment 

test_smoke:
	NETWORKS=alfajores_celo ginkgo -v -r -keepGoing --trace --randomizeAllSpecs --randomizeSuites --progress $(args) ./suite/smoke 

test_performance:
	NETWORKS=alfajores_celo ginkgo -r -timeout=200h -keepGoing --trace --randomizeAllSpecs --randomizeSuites --progress $(args) ./suite/performance 

test_chaos:
	NETWORKS=alfajores_celo ginkgo -r -keepGoing --trace --randomizeAllSpecs --randomizeSuites --progress $(args) ./suite/chaos 
