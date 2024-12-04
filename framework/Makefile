export GOPATH ?= $(shell go env GOPATH)
export GO111MODULE ?= on
CDK8S_CLI_VERSION=2.1.48

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
	pre-commit run go-lint --all-files --show-diff-on-failure --color=always

.PHONY: tidy
tidy:
	pre-commit run go-mod-tidy --hook-stage pre-push --all-files --show-diff-on-failure --color=always

.PHONY: go_mod
go_mod:
	go mod download

test_component: go_mod
	go test -timeout 5m -count 1 -run TestComponentDocker ./...