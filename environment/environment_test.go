package environment

import (
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Environment functionality", func() {
	var conf *config.Config

	BeforeEach(func() {
		var err error
		conf, err = config.NewWithPath(config.LocalConfig, "../config")
		Expect(err).ShouldNot(HaveOccurred())
	})

	DescribeTable("basic environment", func(
		envName string,
		initFunc client.BlockchainNetworkInit,
	) {
		// Setup
		networkConfig, err := initFunc(conf)
		Expect(err).ShouldNot(HaveOccurred())

		env, err := NewBasicEnvironment(envName, 1, networkConfig)
		Expect(err).ShouldNot(HaveOccurred())

		Expect(len(env.GetChainlinkNodes())).ShouldNot(Equal(0))

		mainNode := env.GetChainlinkNodes()[0]
		keys, err := mainNode.ReadETHKeys()
		Expect(err).ShouldNot(HaveOccurred())
		log.Info().Str("ETH Address", keys.Data[0].Attributes.Address).Msg("Got address")

		err = env.TearDown()
		Expect(err).ShouldNot(HaveOccurred())
	},
		Entry("on Ethereum Hardhat", "basic-hardhat", client.NewHardhatNetwork),
	)
})
