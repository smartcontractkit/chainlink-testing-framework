package client

import (
	"context"
	"integrations-framework/config"

	"github.com/ethereum/go-ethereum/crypto"
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
		suppliedPrivateKey, err := crypto.HexToECDSA(privateKeyString)
		Expect(err).ShouldNot(HaveOccurred())
		networkConfig := initFunc(conf)
		wallets, err := networkConfig.Wallets()
		privateKeyToCheckAgainst := wallets.Default().PrivateKey()
		Expect(err).ShouldNot(HaveOccurred())
		Expect(suppliedPrivateKey.Equal(privateKeyToCheckAgainst)).To(BeTrue())
		Expect(address).To(Equal(wallets.Default().Address()))
	},
		Entry("on Ethereum Hardhat", NewEthereumHardhat,
			"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
			"0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"),
	)

	DescribeTable("deploy and interact with the storage contract", func(
		initFunc BlockchainNetworkInit,
		contractVersion string,
	) {
		// Deploy contract
		networkConfig := initFunc(conf)
		client, err := NewBlockchainClient(networkConfig)
		Expect(err).ShouldNot(HaveOccurred())
		wallets, err := networkConfig.Wallets()
		Expect(err).ShouldNot(HaveOccurred())
		storageInstance, err := client.DeployStorageContract(wallets.Default(), wallets.Default())
		Expect(err).ShouldNot(HaveOccurred())

		// Interact with contract
		vers, err := storageInstance.Version(context.Background())
		Expect(err).ShouldNot(HaveOccurred())
		Expect(vers).To(Equal(contractVersion))
	},
		Entry("on Ethereum Hardhat", NewEthereumHardhat, "1.0"),
	)

})
