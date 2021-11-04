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

test_smoke:
	ginkgo -r -keepGoing --trace --randomizeAllSpecs --randomizeSuites --progress $(args) ./suite/smoke 

test_performance:
	ginkgo -r -keepGoing --trace --randomizeAllSpecs --randomizeSuites --progress $(args) ./suite/performance 

test_chaos:
	ginkgo -r -keepGoing --trace --randomizeAllSpecs --randomizeSuites --progress $(args) ./suite/chaos 
