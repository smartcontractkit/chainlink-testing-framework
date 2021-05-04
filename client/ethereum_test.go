package client

import (
	"integrations-framework/config"
	"math/big"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Ethereum functionality", func() {
	var conf *config.Config

	BeforeEach(func() {
		var err error
		conf, err = config.NewConfigWithPath(config.LocalConfig, "../config")
		Expect(err).ShouldNot(HaveOccurred())
	})

	DescribeTable("eth transaction basics", func(
		initFunc BlockchainNetworkInit,
		networkID BlockchainNetworkID,
	) {
		// Setup
		networkConfig, err := initFunc(conf, networkID)
		Expect(err).ShouldNot(HaveOccurred())
		client, err := NewEthereumClient(networkConfig)
		Expect(err).ShouldNot(HaveOccurred())
		wallets, err := networkConfig.Wallets()
		Expect(err).ShouldNot(HaveOccurred())

		// Transaction Settings
		_, _, _, err = client.GetEthTransactionBasics(wallets.Default())
		Expect(err).ShouldNot(HaveOccurred())
		_, err = client.GetTransactionOpts(wallets.Default(), big.NewInt(0))
		Expect(err).ShouldNot(HaveOccurred())

		// Actual transaction
		toWallet, err := wallets.Wallet(1)
		Expect(err).ShouldNot(HaveOccurred())
		_, err = client.SendTransaction(wallets.Default(), toWallet.Address(), 0)
		Expect(err).ShouldNot(HaveOccurred())
	},
		Entry("on Ethereum Hardhat", NewEthereumNetwork, EthereumHardhatID),
	)
})
