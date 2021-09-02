package environment

import (
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/config"
	"github.com/smartcontractkit/integrations-framework/tools"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Environment functionality @unit", func() {
	var conf *config.Config

	BeforeEach(func() {
		var err error
		conf, err = config.NewConfig(tools.ProjectRoot)
		Expect(err).ShouldNot(HaveOccurred())
	})

	DescribeTable("basic environment", func(
		initFunc client.BlockchainNetworkInit,
		chainlinkGroupsInit K8sChainlinkGroupsInit,
		nodeCount int,
	) {
		Skip("Not ready to be run in github")
		// Setup
		networkConfig, err := initFunc(conf)
		Expect(err).ShouldNot(HaveOccurred())

		env, err := NewChainlinkEnvironment(chainlinkGroupsInit, nodeCount, conf, networkConfig)
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
		Entry("1 node cluster", client.NewNetworkFromConfig, NewChainlinkNodesGroups, 1),
		Entry("3 node cluster", client.NewNetworkFromConfig, NewChainlinkNodesGroups, 3),
		Entry("mixed version cluster", client.NewNetworkFromConfig, NewMixedVersionChainlinkGroupInit(2), 3),
	)
})
