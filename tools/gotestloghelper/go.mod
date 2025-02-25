module github.com/smartcontractkit/chainlink-testing-framework/tools/gotestloghelper

go 1.24.0

require (
	github.com/smartcontractkit/chainlink-testing-framework/lib v1.50.20-0.20250106135623-15722ca32b64
	github.com/stretchr/testify v1.9.0
)

require (
	dario.cat/mergo v1.0.1
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

retract [v1.999.0-test-release, v1.999.999-test-release]
