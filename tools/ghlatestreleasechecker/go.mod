module github.com/smartcontractkit/chainlink-testing-framework/tools/ghlatestreleasechecker

go 1.24.0

require github.com/stretchr/testify v1.9.0

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

retract [v1.999.0-test-release, v1.999.999-test-release]
