package contracts

import (
	"context"
	"math/big"

	"github.com/smartcontractkit/integrations-framework/actions"
	"github.com/smartcontractkit/integrations-framework/contracts"
	"github.com/smartcontractkit/integrations-framework/environment"

	"github.com/smartcontractkit/integrations-framework/client"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Basic Contract Interactions", func() {
	var suiteSetup *actions.DefaultSuiteSetup

	BeforeSuite(func() {
		By("Deploying the environment", func() {
			var err error
			suiteSetup, err = actions.DefaultLocalSetup(
				environment.NewChainlinkCluster("../../", 0),
				client.NewNetworkFromConfig,
			)
			Expect(err).ShouldNot(HaveOccurred())

			// Deploy LINK contract
			linkInstance, err := suiteSetup.Deployer.DeployLinkTokenContract(suiteSetup.Wallets.Default())
			Expect(err).ShouldNot(HaveOccurred())
			name, err := linkInstance.Name(context.Background())
			Expect(err).ShouldNot(HaveOccurred())
			Expect(name).To(Equal("ChainLink Token"))
		})
	})

	Describe("with the storage contract", func() {
		It("should properly store and return values", func() {
			// Deploy storage
			storeInstance, err := suiteSetup.Deployer.DeployStorageContract(suiteSetup.Wallets.Default())
			Expect(err).ShouldNot(HaveOccurred())

			testVal := big.NewInt(5)
			// Interact with contract
			err = storeInstance.Set(testVal)
			Expect(err).ShouldNot(HaveOccurred())
			val, err := storeInstance.Get(context.Background())
			Expect(err).ShouldNot(HaveOccurred())
			Expect(val).To(Equal(testVal))
		})
	})

	Describe("with the FluxAggregator contract", func() {
		It("should deploy and be able to interact", func() {
			// Deploy FluxMonitor contract
			fluxOptions := contracts.DefaultFluxAggregatorOptions()
			fluxInstance, err := suiteSetup.Deployer.DeployFluxAggregatorContract(suiteSetup.Wallets.Default(), fluxOptions)
			Expect(err).ShouldNot(HaveOccurred())
			err = fluxInstance.Fund(suiteSetup.Wallets.Default(), big.NewInt(0), big.NewInt(50000000000))
			Expect(err).ShouldNot(HaveOccurred())

			// Interact with contract
			desc, err := fluxInstance.Description(context.Background())
			Expect(err).ShouldNot(HaveOccurred())
			Expect(desc).To(Equal(fluxOptions.Description))
		})

	})

	Describe("the OffChain Aggregator contract", func() {
		It("should deploy and be able to interact", func() {
			// Deploy Offchain contract
			ocrOptions := contracts.DefaultOffChainAggregatorOptions()
			offChainInstance, err := suiteSetup.Deployer.DeployOffChainAggregator(suiteSetup.Wallets.Default(), ocrOptions)
			Expect(err).ShouldNot(HaveOccurred())
			err = offChainInstance.Fund(suiteSetup.Wallets.Default(), nil, big.NewInt(50000000000))
			Expect(err).ShouldNot(HaveOccurred())

			// Check a round
			ans, err := offChainInstance.GetLatestAnswer(context.Background())
			Expect(err).ShouldNot(HaveOccurred())
			Expect(ans).ShouldNot(Equal(nil))
		})
	})

	AfterSuite(func() {
		By("Tearing down the environment", func() {
			suiteSetup.Env.TearDown()
		})
	})
})
