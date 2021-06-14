package client

import (
	"integrations-framework/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {
	var conf *config.Config

	BeforeEach(func() {
		var err error
		conf, err = config.NewWithPath(config.LocalConfig, "../config")
		Expect(err).ShouldNot(HaveOccurred())
	})

	DescribeTable("create new wallet configurations", func(
		initFunc BlockchainNetworkInit,
	) {
		networkConfig, err := initFunc(conf)
		Expect(err).ShouldNot(HaveOccurred())
		wallets, err := networkConfig.Wallets()
		Expect(err).ShouldNot(HaveOccurred())
		rawWallets := wallets.All()
		for index := range rawWallets {
			_, err := wallets.Wallet(index)
			Expect(err).ShouldNot(HaveOccurred())
		}
	},
		Entry("on Ethereum Hardhat", NewHardhatNetwork),
		Entry("on Ethereum Kovan", NewKovanNetwork),
		Entry("on Ethereum Goerli", NewGoerliNetwork),
	)
})
