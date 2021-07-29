package contracts

import (
	"context"
	"math/big"

	"github.com/smartcontractkit/integrations-framework/config"
	"github.com/smartcontractkit/integrations-framework/contracts"
	"github.com/smartcontractkit/integrations-framework/environment"

	"github.com/smartcontractkit/integrations-framework/client"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Contracts", func() {
	var conf *config.Config

	BeforeEach(func() {
		var err error
		conf, err = config.NewWithPath(config.LocalConfig, "../../config")
		Expect(err).ShouldNot(HaveOccurred())
	})

	DescribeTable("deploy and interact with the storage contract", func(
		initFunc client.BlockchainNetworkInit,
		value *big.Int,
	) {
		network, err := initFunc(conf)
		Expect(err).ShouldNot(HaveOccurred())
		testEnv, err := environment.NewK8sEnvironment(environment.NewChainlinkCluster("../../", 0), conf, network)
		Expect(err).ShouldNot(HaveOccurred())
		defer testEnv.TearDown()

		blockchain, err := environment.NewBlockchainClient(testEnv, network)
		Expect(err).ShouldNot(HaveOccurred())
		deployer, err := contracts.NewContractDeployer(blockchain)
		Expect(err).ShouldNot(HaveOccurred())
		wallets, err := network.Wallets()
		Expect(err).ShouldNot(HaveOccurred())
		storeInstance, err := deployer.DeployStorageContract(wallets.Default())
		Expect(err).ShouldNot(HaveOccurred())

		// Interact with contract
		err = storeInstance.Set(value)
		Expect(err).ShouldNot(HaveOccurred())
		val, err := storeInstance.Get(context.Background())
		Expect(err).ShouldNot(HaveOccurred())
		Expect(val).To(Equal(value))
	},
		Entry("on Ethereum Hardhat", client.NewHardhatNetwork, big.NewInt(5)),
	)

	DescribeTable("deploy and interact with the FluxAggregator contract", func(
		initFunc client.BlockchainNetworkInit,
		fluxOptions contracts.FluxAggregatorOptions,
	) {
		network, err := initFunc(conf)
		Expect(err).ShouldNot(HaveOccurred())
		testEnv, err := environment.NewK8sEnvironment(environment.NewChainlinkCluster("../../", 0), conf, network)
		Expect(err).ShouldNot(HaveOccurred())
		defer testEnv.TearDown()

		blockchain, err := environment.NewBlockchainClient(testEnv, network)
		Expect(err).ShouldNot(HaveOccurred())
		deployer, err := contracts.NewContractDeployer(blockchain)
		Expect(err).ShouldNot(HaveOccurred())
		wallets, err := network.Wallets()
		Expect(err).ShouldNot(HaveOccurred())

		// Deploy LINK contract
		linkInstance, err := deployer.DeployLinkTokenContract(wallets.Default())
		Expect(err).ShouldNot(HaveOccurred())
		name, err := linkInstance.Name(context.Background())
		Expect(err).ShouldNot(HaveOccurred())
		Expect(name).To(Equal("ChainLink Token"))

		// Deploy FluxMonitor contract
		fluxInstance, err := deployer.DeployFluxAggregatorContract(wallets.Default(), fluxOptions)
		Expect(err).ShouldNot(HaveOccurred())
		err = fluxInstance.Fund(wallets.Default(), big.NewInt(0), big.NewInt(50000000000))
		Expect(err).ShouldNot(HaveOccurred())

		// Interact with contract
		desc, err := fluxInstance.Description(context.Background())
		Expect(err).ShouldNot(HaveOccurred())
		Expect(desc).To(Equal(fluxOptions.Description))
	},
		Entry("on Ethereum Hardhat", client.NewHardhatNetwork, contracts.DefaultFluxAggregatorOptions()),
	)

	DescribeTable("deploy and interact with the OffChain Aggregator contract", func(
		initFunc client.BlockchainNetworkInit,
		ocrOptions contracts.OffchainOptions,
	) {
		network, err := initFunc(conf)
		Expect(err).ShouldNot(HaveOccurred())
		testEnv, err := environment.NewK8sEnvironment(environment.NewChainlinkCluster("../../", 0), conf, network)
		Expect(err).ShouldNot(HaveOccurred())
		defer testEnv.TearDown()

		blockchain, err := environment.NewBlockchainClient(testEnv, network)
		Expect(err).ShouldNot(HaveOccurred())
		deployer, err := contracts.NewContractDeployer(blockchain)
		Expect(err).ShouldNot(HaveOccurred())
		wallets, err := network.Wallets()
		Expect(err).ShouldNot(HaveOccurred())

		// Deploy LINK contract
		linkInstance, err := deployer.DeployLinkTokenContract(wallets.Default())
		Expect(err).ShouldNot(HaveOccurred())
		name, err := linkInstance.Name(context.Background())
		Expect(err).ShouldNot(HaveOccurred())
		Expect(name).To(Equal("ChainLink Token"))

		// Deploy Offchain contract
		offChainInstance, err := deployer.DeployOffChainAggregator(wallets.Default(), ocrOptions)
		Expect(err).ShouldNot(HaveOccurred())
		err = offChainInstance.Fund(wallets.Default(), nil, big.NewInt(50000000000))
		Expect(err).ShouldNot(HaveOccurred())
	},
		Entry("on Ethereum Hardhat", client.NewNetworkFromConfig, contracts.DefaultOffChainAggregatorOptions()),
	)
})
