lint:
	golangci-lint --color=always run ./... --fix -v

build:
	go build -o ../../gotidy .

test_unit:
	go test -timeout 5m -json -cover -covermode=count -coverprofile=unit-test-coverage.out ./... 2>&1 | tee /tmp/gotest.log | gotestloghelper -ci
