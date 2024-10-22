package pkg

import "github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"

// Common labels for k8s envs
const (
	TTLLabelKey       = "cleanup.kyverno.io/ttl"
	NamespaceLabelKey = "namespace"
)

// Environment types, envs got selected by having a label of that type
const (
	EnvTypeEVM5             = "evm-5-minimal"
	EnvTypeEVM5RemoteRunner = "evm-5-remote-runner"
)

func PGIsReadyCheck() *[]*string {
	return &[]*string{
		ptr.Ptr("pg_isready"),
		ptr.Ptr("-U"),
		ptr.Ptr("postgres"),
	}
}
