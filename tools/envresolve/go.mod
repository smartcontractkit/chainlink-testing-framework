module github.com/smartcontractkit/chainlink-testing-framework/tools/envresolve

go 1.24.0

require (
	github.com/smartcontractkit/chainlink-testing-framework/lib v1.50.20-0.20250106135623-15722ca32b64
	github.com/spf13/cobra v1.8.1
)

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
)

retract [v1.999.0-test-release, v1.999.999-test-release]
