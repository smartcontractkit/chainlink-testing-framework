package environment

import (
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/config"
	"github.com/smartcontractkit/integrations-framework/tools"
	"github.com/smartcontractkit/integrations-framework/types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Environment functionality @unit", func() {
	var (
		conf *config.Config
		env  Environment
	)

	BeforeEach(func() {
		var err error
		conf, err = config.NewConfig(tools.ProjectRoot)
		Expect(err).ShouldNot(HaveOccurred())
	})

	DescribeTable("single network environments", func(
		initFunc types.NewNetworkHook,
		envInitFunc K8sEnvSpecInit,
		nodeCount int,
	) {
		// Setup
		Skip("Not ready to be run in github")

		networkConfig, err := initFunc(conf)
		Expect(err).ShouldNot(HaveOccurred())

		env, err = NewK8sEnvironment(conf, networkConfig)
		Expect(err).ShouldNot(HaveOccurred())
		err = env.DeploySpecs(envInitFunc)
		Expect(err).ShouldNot(HaveOccurred())
		defer env.TearDown()

		clients, err := GetChainlinkClients(env)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(len(clients)).Should(Equal(nodeCount))

		for _, client := range clients {
			key, err := client.PrimaryEthAddress()
			Expect(err).ShouldNot(HaveOccurred())
			log.Info().Str("ETH Address", key).Msg("Got address")
		}
		Expect(err).ShouldNot(HaveOccurred())
	},
		Entry("1 node cluster", client.DefaultNetworkFromConfig, NewChainlinkCluster(1), 1),
		Entry("3 node cluster", client.DefaultNetworkFromConfig, NewChainlinkCluster(3), 3),
		Entry("mixed version cluster", client.DefaultNetworkFromConfig, NewMixedVersionChainlinkCluster(3, 2), 3),
	)

	AfterEach(func() {
		By("Tearing down the environment", func() {
			env.TearDown()
		})
	})
})
