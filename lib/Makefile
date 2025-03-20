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

docker_prune:
	docker system prune -a -f
	docker volume prune -f

test_unit: go_mod
	go test -timeout 5m -json -cover -covermode=count -coverprofile=unit-test-coverage.out $(shell go list ./... | grep -v /k8s/e2e/ | grep -v /k8s/examples/ | grep -v /docker/test_env) 2>&1 | tee /tmp/gotest.log | gotestloghelper -ci

test_docker: go_mod
	go test -timeout 20m -json -failfast -parallel 3 -cover -covermode=atomic -coverprofile=unit-test-coverage.out ./docker/test_env 2>&1 | tee /tmp/gotest.log | gotestloghelper -ci


.PHONY: build_test_base_image
build_k8s_test_base_image:
	./k8s/scripts/buildBaseImage "$(tag)"

.PHONY: build_test_image
build_k8s_test_image:
	./k8s/scripts/buildTestImage "$(tag)" "$(base_tag)"

k8s_test:
	go test -race ./k8s/config -count 1 -v

k8s_test_e2e:
	go test ./k8s/e2e/local-runner -count 1 -test.parallel=12 -v $(args)

k8s_test_e2e_ci: go_mod
	CHAINLINK_USER_TEAM=TT go test ./k8s/e2e/local-runner -count 1 -test.parallel=14 -test.timeout=1h -json 2>&1 | tee /tmp/gotest.log | gotestloghelper -ci -singlepackage

k8s_test_e2e_ci_remote_runner: go_mod
	CHAINLINK_USER_TEAM=TT go test ./k8s/e2e/remote-runner -count 1 -test.parallel=20 -test.timeout=1h -json 2>&1 | tee /tmp/remoterunnergotest.log | gotestloghelper -ci -singlepackage

.PHONY: examples
examples:
	go run k8s/cmd/test.go

.PHONY: tools_build
gotestloghelper_build:
	cd ../tools/gotestloghelper && go build -o ../../gotestloghelper . && cd -

.PHONY: typos
typos:
	pre-commit run detect-typos --all-files --show-diff-on-failure --color=always

.PHONY: nix_shell
nix_shell:
	nix develop
