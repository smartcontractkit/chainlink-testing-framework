module github.com/smartcontractkit/chainlink-testing-framework/tools/breakingchanges

go 1.24.0

retract [v1.999.0-test-release, v1.999.999-test-release]

require (
	github.com/Masterminds/semver/v3 v3.3.1
	golang.org/x/mod v0.22.0
)
