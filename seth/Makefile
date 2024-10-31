.PHONY: build
build:
	solc --abi --overwrite -o contracts/abi contracts/NetworkDebugContract.sol
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

.PHONY: AnvilSync
AnvilSync:
	anvil

.PHONY: Anvil
Anvil:
	anvil > /dev/null 2>&1 &

.PHONY: Geth
Geth:
	rm -rf geth_data/geth
	geth init --datadir geth_data/ geth_data/clique_genesis.json
	geth --graphql --http --http.api admin,debug,web3,eth,txpool,personal,miner,net --http.corsdomain "*" --ws --ws.api admin,debug,web3,eth,txpool,personal,miner,net --ws.origins "*" --mine --miner.etherbase 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266 --unlock f39Fd6e51aad88F6F4ce6aB8827279cffFb92266 --allow-insecure-unlock --datadir ./geth_data --password geth_data/password.txt --nodiscover --vmdebug --networkid 1337 > /dev/null 2>&1 &

.PHONY: GethSync
GethSync:
	rm -rf geth_data/geth
	geth init --datadir geth_data/ geth_data/clique_genesis.json
	geth --graphql --http --http.api admin,debug,web3,eth,txpool,personal,miner,net --http.corsdomain "*" --ws --ws.api admin,debug,web3,eth,txpool,personal,miner,net --ws.origins "*" --mine --miner.etherbase 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266 --unlock f39Fd6e51aad88F6F4ce6aB8827279cffFb92266 --allow-insecure-unlock --datadir ./geth_data --password geth_data/password.txt --nodiscover --vmdebug --networkid 1337

.PHONY: test
test:
	SETH_NETWORK=$(network) SETH_ROOT_PRIVATE_KEY=$(root_private_key) go test -v -count 1 -race `go list ./... | grep -v examples` -run TestSmoke

.PHONY: kill_node
kill_node:
	pkill -f geth || pkill -f anvil || true

# this one is without -race flag, because zerolog is not thread safe and fails the run
.PHONY: test_api
test_api:
	SETH_NETWORK=$(network) SETH_ROOT_PRIVATE_KEY=$(root_private_key) go test -v -count 1 `go list ./... | grep -v examples` -run TestAPI

.PHONY: test_trace
test_trace:
	SETH_NETWORK=$(network) SETH_ROOT_PRIVATE_KEY=$(root_private_key) go test -v -count 1 -race `go list ./... | grep -v examples` -run TestTrace

.PHONY: test_cli
test_cli:
	SETH_NETWORK=$(network) SETH_ROOT_PRIVATE_KEY=$(root_private_key) go test -v -count 1 -race `go list ./... | grep -v examples` -run TestCLI

.PHONY: test_others
test_others:
	SETH_CONFIG_PATH="seth.toml" SETH_NETWORK=$(network) SETH_ROOT_PRIVATE_KEY=$(root_private_key) go test -v -count 1 -race `go list ./... | grep -v examples` -run "TestContractMap|TestGasEstimator|TestRPCHealthCheck|TestUtil|TestContract|TestConfig"

# this one is without -race flag, because zerolog is not thread safe and fails the run
.PHONY: test_gas_bumping
test_gas_bumping:
	SETH_NETWORK=$(network) SETH_ROOT_PRIVATE_KEY=$(root_private_key) go test -v -count 1 `go list ./... | grep -v examples` -run "TestGasBumping"

.PHONY: test+cover
test_cover:
	SETH_NETWORK=$(network) SETH_ROOT_PRIVATE_KEY=$(root_private_key) go test -v -coverprofile cover.out -count 1 `go list ./... | grep -v examples` -run "TestAPI|TestSmoke|TestContract|TestGasEstimator"
	go tool cover -html cover.out

.PHONY: lint
lint:
	golangci-lint --color=always run -v -c .golangci.yml
