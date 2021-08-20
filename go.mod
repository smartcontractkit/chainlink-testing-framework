module github.com/smartcontractkit/integrations-framework

go 1.16

require (
	github.com/avast/retry-go v3.0.0+incompatible
	github.com/ethereum/go-ethereum v1.10.4
	github.com/ghodss/yaml v1.0.0
	github.com/google/go-github/v38 v38.1.0
	github.com/hashicorp/go-multierror v1.1.1
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.13.0
	github.com/pkg/errors v0.9.1
	github.com/rs/zerolog v1.21.0
	github.com/satori/go.uuid v1.2.0
	github.com/smartcontractkit/integrations-framework/explorer v0.0.0
	github.com/smartcontractkit/libocr v0.0.0-20210803133922-ddddd3dce7e5
	github.com/spf13/cobra v1.1.1
	github.com/spf13/viper v1.8.0
	k8s.io/api v0.22.0
	k8s.io/apimachinery v0.22.0
	k8s.io/client-go v0.22.0
)

replace github.com/smartcontractkit/integrations-framework/explorer => ./explorer/
