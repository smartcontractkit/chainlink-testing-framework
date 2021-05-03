package client

import (
	"context"
	"integrations-framework/config"
	"math/big"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {
	var conf *config.Config

	BeforeEach(func() {
		var err error
		conf, err = config.NewConfigWithPath(config.LocalConfig, "../config")
		Expect(err).ShouldNot(HaveOccurred())
	})

	DescribeTable("create new wallet configurations", func(
		initFunc BlockchainNetworkInit,
		privateKeyString string,
		address string,
	) {
		networkConfig := initFunc(conf)
		wallets, err := networkConfig.Wallets()
		Expect(err).ShouldNot(HaveOccurred())
		Expect(wallets.Default().PrivateKey()).To(Equal(privateKeyString))
		Expect(address).To(Equal(wallets.Default().Address()))
	},
		Entry("on Ethereum Hardhat", NewEthereumHardhat,
			"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
			"0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"),
	)

	DescribeTable("deploy and interact with the storage contract", func(
		initFunc BlockchainNetworkInit,
		value *big.Int,
	) {
		// Deploy contract
		networkConfig := initFunc(conf)
		client, err := NewBlockchainClient(networkConfig)
		Expect(err).ShouldNot(HaveOccurred())
		wallets, err := networkConfig.Wallets()
		Expect(err).ShouldNot(HaveOccurred())
		storeInstance, err := client.DeployStorageContract(wallets.Default(), wallets.Default())
		Expect(err).ShouldNot(HaveOccurred())

		// Interact with contract
		err = storeInstance.Set(context.Background(), big.NewInt(5))
		Expect(err).ShouldNot(HaveOccurred())
		val, err := storeInstance.Get(context.Background())
		Expect(err).ShouldNot(HaveOccurred())
		Expect(val).To(Equal(value))
	},
		Entry("on Ethereum Hardhat", NewEthereumHardhat, big.NewInt(5)),
	)

})
