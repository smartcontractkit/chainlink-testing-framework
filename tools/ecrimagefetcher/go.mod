module github.com/smartcontractkit/chainlink-testing-framework/tools/ecrimagefetcher

go 1.22.5

require (
	github.com/Masterminds/semver/v3 v3.2.1
	github.com/itchyny/gojq v0.12.16
	github.com/stretchr/testify v1.9.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/itchyny/timefmt-go v0.1.6 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

retract [v1.999.0-test-release, v1.999.999-test-release]
