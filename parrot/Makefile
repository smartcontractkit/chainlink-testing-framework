.PHONY: lint
lint:
	golangci-lint --color=always run ./... --fix -v

.PHONY: test
test:
	go install github.com/gotesttools/gotestfmt/v2/cmd/gotestfmt@latest
	set -euo pipefail
	go test -json -cover -coverprofile cover.out -v ./... 2>&1 | tee /tmp/gotest.log | gotestfmt

.PHONY: test_race
test_race:
	go install github.com/gotesttools/gotestfmt/v2/cmd/gotestfmt@latest
	set -euo pipefail
	go test -json -cover -count=1 -race -coverprofile cover.out -v ./... 2>&1 | tee /tmp/gotest.log | gotestfmt

.PHONY: test_unit
test_unit:
	go test -coverprofile cover.out ./...

.PHONY: bench
bench:
	go test -bench=. -run=^$$ ./...