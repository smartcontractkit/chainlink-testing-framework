BIN_DIR = bin
export GOPATH ?= $(shell go env GOPATH)
export GO111MODULE ?= on
CDK8S_CLI_VERSION=2.1.48

.PHONY: lint
lint:
	golangci-lint --color=always run -v

.PHONY: docker_prune
docker_prune:
	docker system prune -a -f
	docker volume prune -f

.PHONY: install_deps
install_deps:
	asdf plugin-add nodejs || true
	asdf plugin-add yarn || true
	asdf plugin-add golang || true
	asdf plugin-add k3d || true
	asdf plugin-add helm || true
	asdf plugin-add kubectl || true
	asdf plugin-add golangci-lint || true
	asdf install
	mkdir /tmp/k3dvolume/ || true
	yarn global add cdk8s-cli@$(CDK8S_CLI_VERSION)
	curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash
	helm repo add chainlink-qa https://raw.githubusercontent.com/smartcontractkit/qa-charts/gh-pages/
	helm repo add grafana https://grafana.github.io/helm-charts
	helm repo update

.PHONY: create_cluster
create_cluster:
	k3d cluster create local --config k3d.yaml

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
	helm upgrade --wait --namespace monitoring --install loki grafana/loki-stack  --set grafana.enabled=true,prometheus.enabled=true,prometheus.alertmanager.persistentVolume.enabled=false,prometheus.server.persistentVolume.enabled=false,loki.persistence.enabled=false --values grafana/values.yml
	kubectl port-forward --namespace monitoring service/loki-grafana 3000:80

.PHONY: uninstall_monitoring
uninstall_monitoring:
	helm uninstall --namespace monitoring loki

.PHONY: build_test_base_image
build_test_base_image:
	./scripts/buildBaseImage "$(tag)"

.PHONY: build_test_image
build_test_image:
	./scripts/buildTestImage "$(tag)" "$(base_tag)"

.PHONY: test
test:
	go test -race ./config -count 1 -v

.PHONY: test_e2e
test_e2e:
	go test ./e2e/local-runner -count 1 -test.parallel=12 -v $(args)

test_e2e_ci:
	go test ./e2e/local-runner -count 1 -v -test.parallel=14 -test.timeout=1h -json 2>&1 | tee /tmp/gotest.log | gotestfmt

test_e2e_ci_remote_runner:
	go test ./e2e/remote-runner -count 1 -v -test.parallel=16 -test.timeout=1h -json 2>&1 | tee /tmp/remoterunnergotest.log | gotestfmt

.PHONY: examples
examples:
	go run cmd/test.go

.PHONY: chaosmesh
chaosmesh: ## there is currently a bug on JS side to import all CRDs from one yaml file, also a bug with stdin, so using cluster directly trough file
	kubectl get crd networkchaos.chaos-mesh.org -o json > tmp.json && cdk8s import -o imports/k8s/networkchaos tmp.json
	kubectl get crd stresschaos.chaos-mesh.org -o json > tmp.json && cdk8s import -o imports/k8s/stresschaos tmp.json
	kubectl get crd timechaos.chaos-mesh.org -o json > tmp.json && cdk8s import -o imports/k8s/timechaos tmp.json
	kubectl get crd podchaos.chaos-mesh.org -o json > tmp.json && cdk8s import -o imports/k8s/podchaos tmp.json
	kubectl get crd podiochaos.chaos-mesh.org -o json > tmp.json && cdk8s import -o imports/k8s/podiochaos tmp.json
	kubectl get crd httpchaos.chaos-mesh.org -o json > tmp.json && cdk8s import -o imports/k8s/httpchaos tmp.json
	kubectl get crd iochaos.chaos-mesh.org -o json > tmp.json && cdk8s import -o imports/k8s/iochaos tmp.json
	kubectl get crd podnetworkchaos.chaos-mesh.org -o json > tmp.json && cdk8s import -o imports/k8s/podnetworkchaos tmp.json
	rm -rf tmp.json