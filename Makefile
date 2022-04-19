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
	golangci-lint --color=always run ./... -v

go_mod:
	go mod tidy
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
	asdf plugin add k3d https://github.com/spencergilbert/asdf-k3d.git || true
	asdf plugin add act https://github.com/grimoh/asdf-act.git || true
	asdf plugin add golangci-lint https://github.com/hypnoglow/asdf-golangci-lint.git || true
	asdf plugin add actionlint || true
	asdf plugin add shellcheck || true
	asdf install
endif

install: go_mod install_tools

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

compile_soak:
	ginkgo build ./suite/soak/ -o ./soak.test

compile_smoke:
	ginkgo build ./suite/smoke/ -o ./smoke.test

compile_soak_all:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 ginkgo build ./suite/soak/ -o ./linux_amd64_soak.test
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 ginkgo build ./suite/soak/ -o ./linux_arm64_soak.test

	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 ginkgo build ./suite/soak/ -o ./darwin_amd64_soak.test
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 ginkgo build ./suite/soak/ -o ./darwin_arm64_soak.test

	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 ginkgo build ./suite/soak/ -o ./windows_amd64_soak.test

compile_smoke_all:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 ginkgo build ./suite/smoke/ -o ./linux_amd64_smoke.test
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 ginkgo build ./suite/smoke/ -o ./linux_arm64_smoke.test

	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 ginkgo build ./suite/smoke/ -o ./darwin_amd64_smoke.test
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 ginkgo build ./suite/smoke/ -o ./darwin_arm64_smoke.test

	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 ginkgo build ./suite/smoke/ -o ./windows_amd64_smoke.test