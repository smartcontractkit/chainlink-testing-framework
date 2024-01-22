BIN_DIR = bin
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
	golangci-lint --color=always run ./... --fix -v

go_mod:
	go mod tidy
	go mod download

.PHONY: install_gotestfmt
install_gotestfmt:
	go install github.com/gotesttools/gotestfmt/v2/cmd/gotestfmt@latest
	set -euo pipefail

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
	asdf plugin-add yarn || true
	asdf plugin-add k3d || true
	asdf plugin-add helm || true
	asdf plugin-add kubectl || true
	asdf plugin-add python || true
	asdf plugin add pre-commit || true
	asdf install
	mkdir /tmp/k3dvolume/ || true
	yarn global add cdk8s-cli@$(CDK8S_CLI_VERSION)
	curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash
	helm repo add chainlink-qa https://raw.githubusercontent.com/smartcontractkit/qa-charts/gh-pages/
	helm repo add grafana https://grafana.github.io/helm-charts
	helm repo update
	pre-commit install
endif

install: go_mod install_tools

install_ci: go_mod install_tools

docker_prune:
	docker system prune -a -f
	docker volume prune -f

compile_contracts:
	python3 ./utils/compile_contracts.py

test_unit: install_gotestfmt
	go test -json -cover -covermode=count -coverprofile=unit-test-coverage.out ./client ./gauntlet ./testreporters ./k8s/config ./utils/osutil 2 2>&1 | tee /tmp/gotest.log | gotestfmt

test_docker: install_gotestfmt
	go test -json -cover -covermode=count -coverprofile=unit-test-coverage.out ./docker/test_env ./logstream 2>&1 | tee /tmp/gotest.log | gotestfmt	


#######################
# K8s Helpers
#######################
.PHONY: create_cluster
create_cluster:
	k3d cluster create local --config ./k8s/k3d.yaml

.PHONY: start_cluster
start_cluster:
	k3d cluster start local

.PHONY: stop_cluster
stop_cluster:
	k3d cluster stop local

.PHONY: stop_cluster
delete_cluster:
	k3d cluster delete local

.PHONY: install_monitoring
install_monitoring:
	helm repo add grafana https://grafana.github.io/helm-charts
	helm repo update
	kubectl create namespace monitoring || true
	helm upgrade --wait --namespace monitoring --install loki grafana/loki-stack  --set grafana.enabled=true,prometheus.enabled=true,prometheus.alertmanager.persistentVolume.enabled=false,prometheus.server.persistentVolume.enabled=false,loki.persistence.enabled=false --values k8s/grafana/values.yml
	kubectl port-forward --namespace monitoring service/loki-grafana 3000:80

.PHONY: uninstall_monitoring
uninstall_monitoring:
	helm uninstall --namespace monitoring loki

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

k8s_test_e2e_ci:
	go test ./k8s/e2e/local-runner -count 1 -v -test.parallel=14 -test.timeout=1h -json 2>&1 | tee /tmp/gotest.log | gotestfmt

k8s_test_e2e_ci_remote_runner:
	go test ./k8s/e2e/remote-runner -count 1 -v -test.parallel=20 -test.timeout=1h -json 2>&1 | tee /tmp/remoterunnergotest.log | gotestfmt

.PHONY: examples
examples:
	go run k8s/cmd/test.go

.PHONY: chaosmesh
chaosmesh: ## there is currently a bug on JS side to import all CRDs from one yaml file, also a bug with stdin, so using cluster directly trough file
	kubectl get crd networkchaos.chaos-mesh.org -o json > tmp.json && cdk8s import -o k8s/imports/k8s/networkchaos tmp.json
	kubectl get crd stresschaos.chaos-mesh.org -o json > tmp.json && cdk8s import -o k8s/imports/k8s/stresschaos tmp.json
	kubectl get crd timechaos.chaos-mesh.org -o json > tmp.json && cdk8s import -o k8s/imports/k8s/timechaos tmp.json
	kubectl get crd podchaos.chaos-mesh.org -o json > tmp.json && cdk8s import -o k8s/imports/k8s/podchaos tmp.json
	kubectl get crd podiochaos.chaos-mesh.org -o json > tmp.json && cdk8s import -o k8s/imports/k8s/podiochaos tmp.json
	kubectl get crd httpchaos.chaos-mesh.org -o json > tmp.json && cdk8s import -o k8s/imports/k8s/httpchaos tmp.json
	kubectl get crd iochaos.chaos-mesh.org -o json > tmp.json && cdk8s import -o k8s/imports/k8s/iochaos tmp.json
	kubectl get crd podnetworkchaos.chaos-mesh.org -o json > tmp.json && cdk8s import -o k8s/imports/k8s/podnetworkchaos tmp.json
	rm -rf tmp.json
