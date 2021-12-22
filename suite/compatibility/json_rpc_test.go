package compatibility

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/smartcontractkit/helmenv/environment"
	"github.com/smartcontractkit/helmenv/tools"
)

var _ = Describe("Json RPC compatibility @json_rpc", func() {
	var (
		err error
		e   *environment.Environment
	)

	BeforeEach(func() {
		By("Deploying the environment", func() {
			e, err = environment.DeployOrLoadEnvironment(
				environment.NewChainlinkConfig(nil),
				tools.ChartsRoot,
			)
			Expect(err).ShouldNot(HaveOccurred())
			err = e.ConnectAll()
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
})
