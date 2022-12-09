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
	golangci-lint --color=always run ./... --fix -v

go_mod:
	go mod tidy
	go mod download

install_tools:
ifeq ($(OSFLAG),$(WINDOWS))
	echo "If you are running windows and know how to install what is needed, please contribute by adding it here!"
	exit 1
endif
ifeq ($(OSFLAG),$(OSX))
	brew install asdf
	asdf plugin-add nodejs https://github.com/asdf-vm/asdf-nodejs.git || true
	asdf plugin-add golang https://github.com/kennyp/asdf-golang.git || true
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
	go test -cover -covermode=count -coverprofile=unit-test-coverage.out ./client ./gauntlet ./testreporters
