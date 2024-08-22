.PHONY: test
test:
	go test -v -count 1 `go list ./... | grep -v examples` -run TestSmoke

.PHONY: test_race
test_race:
	go test -v -race -count 1 `go list ./... | grep -v examples` -run TestSmoke

.PHONY: test+cover
test_cover:
	go test -v -coverprofile cover.out -count 1 `go list ./... | grep -v examples` -run "TestSmoke|TestAPI"
	go tool cover -html cover.out

.PHONY: install
install:
	go install cmd/havoc.go

.PHONE: build
build:
	go build cmd/havoc.go

.PHONY: lint
lint:
	golangci-lint --color=always run -v
