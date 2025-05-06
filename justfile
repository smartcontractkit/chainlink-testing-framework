go_test_args := env('GO_TEST_ARGS', '')

# Print all the commands
default:
    @just --list

# Install pre-commit hooks
install:
    pre-commit install

# Install gotestloghelper for enhanced go test logs output
install-loghelper:
    go install github.com/smartcontractkit/chainlink-testing-framework/tools/gotestloghelper@latest

# Lint a module, example: just lint wasp
lint dir_path:
    cd {{dir_path}} && golangci-lint --color=always run -v -c {{invocation_directory()}}/.golangci.yaml

# Lint all the modules
lint-all:
    @just lint framework
    @just lint parrot
    @just lint wasp
    @just lint seth
    @just lint havoc
    @just lint k8s-test-runner
    @just lint lib
    @just lint tools/workflowresultparser
    @just lint tools/asciitable
    @just lint tools/ghlatestreleasechecker
    @just lint tools/testlistgenerator
    @just lint tools/ecrimagefetcher
    @just lint tools/gotestloghelper

# Run all the tests
test-all:
    @just test k8s-test-runner ./... &
    @just test parrot ./... &
    @just test tools/workflowresultparser ./... &
    @just test tools/asciitable ./... &
    @just test tools/ghlatestreleasechecker ./... &
    @just test tools/ecrimagefetcher ./... &
    @just test tools/testlistgenerator ./... &
    @just test tools/gotestloghelper ./... &
    @just test tools/citool ./... &
    wait
    @just test framework TestComponent
    @just test wasp TestSmoke
    @just test wasp TestBenchSpy

goimports-all:
    @just _default_goimports k8s-test-runner
    @just _default_goimports lib
    @just _default_goimports parrot
    @just _default_goimports tools/workflowresultparser
    @just _default_goimports tools/asciitable
    @just _default_goimports tools/ghlatestreleasechecker
    @just _default_goimports tools/ecrimagefetcher
    @just _default_goimports tools/testlistgenerator
    @just _default_goimports tools/gotestloghelper
    @just _default_goimports tools/citool
    @just _default_goimports framework
    @just _default_goimports wasp

_default_goimports dir:
    goimports -local github.com/smartcontractkit/chainlink-testing-framework -w {{dir}}

# Default test command (cacheable), set GO_TEST_ARGS="-count 1" to disable cache
_default_cached_test dir test_regex:
    cd {{dir}} && go test {{go_test_args}} -v -race `go list ./... | grep -v examples` -run {{test_regex}}

# Default test + coverage command (no-cache)
_default_cover_test dir test_regex:
    cd {{dir}} && go test {{go_test_args}} -v -race -cover -coverprofile=cover.out `go list ./... | grep -v examples` -run {{test_regex}}

# Run tests for a package, example: just test wasp TestSmoke, example: just test tools/citool ./...
test dir_path test_regex:
    just _default_cached_test {{dir_path}} {{test_regex}}

# Run tests for a package, example: just test wasp TestSmoke, example: just test tools/citool ./...
test-cover dir_path test_regex:
    just _default_cover_test {{dir_path}} {{test_regex}}

# Open module coverage
cover dir_path:
    cd {{dir_path}} && go tool cover -html cover.out

# WASP: Upload WASP dashboard
wasp-dashboard:
    cd wasp && go run dashboard/cmd/main.go

# Seth: build contracts and generate Go bindings
seth-build-contracts:
    @just seth-solc
    cd seth && solc --abi --overwrite -o contracts/abi contracts/NetworkDebugContract.sol
    solc --bin --overwrite -o contracts/bin contracts/NetworkDebugContract.sol
    abigen --bin=contracts/bin/NetworkDebugContract.bin --abi=contracts/abi/NetworkDebugContract.abi --pkg=network_debug_contract --out=contracts/bind/NetworkDebugContract/NetworkDebugContract.go
    solc --abi --overwrite -o contracts/abi contracts/NetworkDebugSubContract.sol
    solc --bin --overwrite -o contracts/bin contracts/NetworkDebugSubContract.sol
    abigen --bin=contracts/bin/NetworkDebugSubContract.bin --abi=contracts/abi/NetworkDebugSubContract.abi --pkg=network_debug_sub_contract --out=contracts/bind/NetworkDebugSubContract/NetworkDebugSubContract.go
    solc --abi --overwrite -o contracts/abi contracts/TestContractOne.sol
    solc --bin --overwrite -o contracts/bin contracts/TestContractOne.sol
    abigen --bin=contracts/bin/TestContractOne.bin --abi=contracts/abi/TestContractOne.abi --pkg=unique_event_one --out=contracts/bind/TestContractOne/TestContractOne.go
    solc --abi --overwrite -o contracts/abi contracts/TestContractTwo.sol
    solc --bin --overwrite -o contracts/bin contracts/TestContractTwo.sol
    abigen --bin=contracts/bin/TestContractTwo.bin --abi=contracts/abi/TestContractTwo.abi --pkg=unique_event_two --out=contracts/bind/TestContractTwo/TestContractTwo.go

# Seth: use solc version
seth-solc:
    cd seth && solc-select install 0.8.19
    solc-select use 0.8.19

# Seth: run anvil network
seth-Anvil:
    cd seth && anvil > /dev/null 2>&1 &

# Seth: run Geth network
seth-Geth:
    cd seth && rm -rf geth_data/geth && \
    geth init --datadir geth_data/ geth_data/clique_genesis.json && \
    geth --graphql --http --http.api admin,debug,web3,eth,txpool,personal,miner,net --http.corsdomain "*" --ws --ws.api admin,debug,web3,eth,txpool,personal,miner,net --ws.origins "*" --mine --miner.etherbase 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266 --unlock f39Fd6e51aad88F6F4ce6aB8827279cffFb92266 --allow-insecure-unlock --datadir ./geth_data --password geth_data/password.txt --nodiscover --vmdebug --networkid 1337 > /dev/null 2>&1 &

# Seth: run Seth tests, example: just seth-test Anvil http://localhost:8545 "TestAPI"
seth-test network url test_regex extra_flags:
    @just seth-{{network}}
    cd seth && SETH_URL={{url}} SETH_NETWORK={{network}} SETH_ROOT_PRIVATE_KEY=ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 go test -v "{{extra_flags}}" `go list ./... | grep -v examples` -run "{{test_regex}}" || pkill -f {{network}}

# Run pre-commit hooks, build, lint, tidy, check typos
pre-commit:
    pre-commit run --hook-stage pre-commit --show-diff-on-failure --color=always

# Generate Go modules graph
modgraph:
    go install github.com/jmank88/gomods@v0.1.5
    go install github.com/jmank88/modgraph@v0.1.0
    ./modgraph > go.md
    git diff --minimal --exit-code

# Run "go mod tidy" for all the packages
gomods-tidy:
    go install github.com/jmank88/gomods@v0.1.5
    gomods tidy

# CI: run go mod tidy for all packages and check there are no changes and go.sums are consistent everywhere
gomods-tidy-ci:
    go install github.com/jmank88/gomods@v0.1.5 && gomods tidy
    git add --all
    git diff --minimal --cached --exit-code

# Run GitHub Actions lint for CI workflows
actionlint:
    go install github.com/rhysd/actionlint/cmd/actionlint@latest
    actionlint

# Serve MDBook locally
book:
    cd book && mdbook serve -p 9999