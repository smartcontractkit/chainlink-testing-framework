BIN_DIR = bin
export GOPATH ?= $(shell go env GOPATH)
export GO111MODULE ?= on

LINUX=LINUX
OSX=OSX
WINDOWS=WIN32
OSFLAG :=
ifeq ($(OS),Windows_NT)
	OSFLAG = $(WINDOWS)
else
	UNAME_S := $(shell uname -s)
	ifeq ($(UNAME_S),Linux)
		OSFLAG = $(LINUX)
	endif
	ifeq ($(UNAME_S),Darwin)
		OSFLAG = $(OSX)
	endif
endif

lint:
	${BIN_DIR}/golangci-lint --color=always run ./... -v

golangci:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ${BIN_DIR} v1.42.0

go_mod:
	go mod download

install_tools:
ifeq ($(OSFLAG),$(LINUX))
	# used for linux and ci
	go install github.com/onsi/ginkgo/v2/ginkgo@v$(shell cat ./.tool-versions | grep ginkgo | sed -En "s/ginkgo.(.*)/\1/p")
endif
ifeq ($(OSFLAG),$(WINDOWS))
	echo "If you are running windows and know how to install what is needed, please contribute by adding it here!"
	exit 1
endif
ifeq ($(OSFLAG),$(OSX))
	brew install asdf
	asdf plugin-add nodejs https://github.com/asdf-vm/asdf-nodejs.git || true
	asdf plugin-add golang https://github.com/kennyp/asdf-golang.git || true
	asdf plugin-add ginkgo https://github.com/jimmidyson/asdf-ginkgo.git || true
	asdf install
endif

install: go_mod golangci install_tools

install_ci: go_mod install_tools

compile_contracts:
	python3 ./utils/compile_contracts.py

test_unit:
	ginkgo -r --junit-report=tests-unit-report.xml --keep-going --trace --randomize-all --randomize-suites --progress -cover -covermode=count -coverprofile=unit-test-coverage.out -nodes=10 ./client ./config ./gauntlet ./testreporters

test_soak:
	go test -v -count=1 ./suite/soak/soak_runner_test.go

test_smoke:
	ginkgo -v -r --junit-report=tests-smoke-report.xml --keep-going --trace --randomize-all --randomize-suites --progress $(args) ./suite/smoke 

test_performance:
	ginkgo -v -r -timeout=200h --junit-report=tests-performance-report.xml --keep-going --trace --randomize-all --randomize-suites --progress $(args) ./suite/performance 

test_chaos:
	ginkgo -r --junit-report=tests-chaos-report.xml --keep-going --trace --randomize-all --randomize-suites --progress $(args) ./suite/chaos 
