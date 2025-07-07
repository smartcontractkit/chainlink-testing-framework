module github.com/smartcontractkit/chainlink-testing-framework/framework/examples

go 1.24.2

replace (
	github.com/smartcontractkit/chainlink-testing-framework/framework => ../../
	github.com/smartcontractkit/chainlink-testing-framework/framework/components/dockercompose => ../../components/dockercompose
	github.com/smartcontractkit/chainlink-testing-framework/havoc => ../../../havoc
	github.com/smartcontractkit/chainlink-testing-framework/wasp => ../../../wasp
)

require (
	github.com/block-vision/sui-go-sdk v1.0.6
	github.com/blocto/solana-go-sdk v1.30.0
	github.com/ethereum/go-ethereum v1.15.0
	github.com/go-resty/resty/v2 v2.16.3
	github.com/google/go-github/v72 v72.0.0
	github.com/smartcontractkit/chainlink-testing-framework/framework v0.8.9
	github.com/smartcontractkit/chainlink-testing-framework/framework/components/dockercompose v0.0.0-00010101000000-000000000000
	github.com/smartcontractkit/chainlink-testing-framework/havoc v1.50.2
	github.com/smartcontractkit/chainlink-testing-framework/seth v1.50.10
	github.com/smartcontractkit/chainlink-testing-framework/wasp v1.50.2
	github.com/smartcontractkit/chainlink/v2 v2.20.0
	github.com/stretchr/testify v1.10.0
	github.com/testcontainers/testcontainers-go v0.37.0
)

// replicating the replace directive on cosmos SDK
replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
