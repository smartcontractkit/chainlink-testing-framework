module github.com/smartcontractkit/chainlink-testing-framework/tools/citool

go 1.24.0

require (
	github.com/smartcontractkit/chainlink-testing-framework/lib v1.50.20-0.20250106135623-15722ca32b64
	github.com/spf13/cobra v1.8.1
	gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
)

retract [v1.999.0-test-release, v1.999.999-test-release]
