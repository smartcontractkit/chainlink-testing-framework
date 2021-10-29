package environment

import (
	"github.com/smartcontractkit/integrations-framework/client"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/smartcontractkit/integrations-framework/config"
	"github.com/smartcontractkit/integrations-framework/tools"
)

var _ = Describe("Environment with Helm @helm_deploy", func() {
	var conf *config.Config

	Describe("Chart deployments", func() {
		var env Environment

		BeforeEach(func() {
			var err error
			conf, err = config.NewConfig(tools.ProjectRoot)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("Deploy Geth reorg chart", func() {
			Skip("Not ready to be run in github")
			conf.Networks = []string{"ethereum_geth_reorg"}
			networkConfig, err := client.DefaultNetworkFromConfig(conf)
			Expect(err).ShouldNot(HaveOccurred())
			env, err = NewK8sEnvironment(conf, networkConfig)
			Expect(err).ShouldNot(HaveOccurred())
			err = env.DeploySpecs(NewChainlinkCluster(1))
			Expect(err).ShouldNot(HaveOccurred())
			// check service details has EVM port
			sd, err := env.GetServiceDetails(networkConfig.RemotePort())
			Expect(err).ShouldNot(HaveOccurred())
			Expect(sd.RemoteURL).Should(ContainSubstring(strconv.Itoa(int(networkConfig.RemotePort()))))
		})

		AfterEach(func() {
			By("Tearing down the environment", func() {
				env.TearDown()
			})
		})
	})
})
