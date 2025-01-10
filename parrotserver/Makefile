.PHONY: lint
lint:
	golangci-lint --color=always run ./... --fix -v

.PHONY: test_unit
test_unit:
	go test ./...

.PHONY: bench
bench:
	go test -bench=. -run=^$$ ./...