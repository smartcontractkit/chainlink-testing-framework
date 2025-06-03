.PHONY: test_unit
test_unit:
	go test ./... -timeout 5m -json -cover -covermode=count -coverprofile=unit-test-coverage.out 2>&1 | tee /tmp/gotest.log | gotestloghelper -ci

.PHONY: test
test:
	gotestsum -- -cover ./...

.PHONY: test-package
test-package:
	gotestsum -- -cover ./$(PKG)

.PHONY: test-race
test-race:
	gotestsum -- -cover -count=1 -race

.PHONY: bench
bench:
	go test -run=^$$ -bench=. -benchmem ./...
