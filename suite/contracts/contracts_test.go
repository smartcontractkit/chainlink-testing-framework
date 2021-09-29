package contracts

import (
	"context"
	"math/big"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/smartcontractkit/integrations-framework/actions"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/contracts"
	"github.com/smartcontractkit/integrations-framework/environment"
	"github.com/smartcontractkit/integrations-framework/tools"
)

var _ = Describe("Basic Contract Interactions @contract", func() {
	var suiteSetup *actions.DefaultSuiteSetup
	var defaultWallet client.BlockchainWallet

	BeforeEach(func() {
		By("Deploying the environment", func() {
			var err error
			suiteSetup, err = actions.DefaultLocalSetup(
				environment.NewChainlinkCluster(0),
				client.NewNetworkFromConfig,
				tools.ProjectRoot,
			)
			Expect(err).ShouldNot(HaveOccurred())
			defaultWallet = suiteSetup.Wallets.Default()
		})
	})

	It("can deploy all contracts", func() {
		By("basic interaction with a storage contract", func() {
			storeInstance, err := suiteSetup.Deployer.DeployStorageContract(defaultWallet)
			Expect(err).ShouldNot(HaveOccurred())
			testVal := big.NewInt(5)
			err = storeInstance.Set(testVal)
			Expect(err).ShouldNot(HaveOccurred())
			val, err := storeInstance.Get(context.Background())
			Expect(err).ShouldNot(HaveOccurred())
			Expect(val).To(Equal(testVal))
		})

		By("deploying the flux monitor contract", func() {
			rac, err := suiteSetup.Deployer.DeployReadAccessController(suiteSetup.Wallets.Default())
			Expect(err).ShouldNot(HaveOccurred())
			flags, err := suiteSetup.Deployer.DeployFlags(suiteSetup.Wallets.Default(), rac.Address())
			Expect(err).ShouldNot(HaveOccurred())
			_, err = suiteSetup.Deployer.DeployDeviationFlaggingValidator(suiteSetup.Wallets.Default(), flags.Address(), big.NewInt(0))
			Expect(err).ShouldNot(HaveOccurred())
			fluxOptions := contracts.DefaultFluxAggregatorOptions()
			_, err = suiteSetup.Deployer.DeployFluxAggregatorContract(defaultWallet, fluxOptions)
			Expect(err).ShouldNot(HaveOccurred())
		})

		By("deploying the ocr contract", func() {
			ocrOptions := contracts.DefaultOffChainAggregatorOptions()
			_, err := suiteSetup.Deployer.DeployOffChainAggregator(defaultWallet, ocrOptions)
			Expect(err).ShouldNot(HaveOccurred())
		})
		By("deploying keeper contracts", func() {
			ef, err := suiteSetup.Deployer.DeployMockETHLINKFeed(suiteSetup.Wallets.Default(), big.NewInt(2e18))
			Expect(err).ShouldNot(HaveOccurred())
			gf, err := suiteSetup.Deployer.DeployMockGasFeed(suiteSetup.Wallets.Default(), big.NewInt(2e11))
			Expect(err).ShouldNot(HaveOccurred())
			_, err = suiteSetup.Deployer.DeployKeeperRegistry(
				suiteSetup.Wallets.Default(),
				&contracts.KeeperRegistryOpts{
					LinkAddr:             suiteSetup.Link.Address(),
					ETHFeedAddr:          ef.Address(),
					GasFeedAddr:          gf.Address(),
					PaymentPremiumPPB:    uint32(200000000),
					BlockCountPerTurn:    big.NewInt(3),
					CheckGasLimit:        uint32(2500000),
					StalenessSeconds:     big.NewInt(90000),
					GasCeilingMultiplier: uint16(1),
					FallbackGasPrice:     big.NewInt(2e11),
					FallbackLinkPrice:    big.NewInt(2e18),
				},
			)
			Expect(err).ShouldNot(HaveOccurred())
			err = suiteSetup.Client.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred())
		})
		By("deploying vrf contract", func() {
			bhs, err := suiteSetup.Deployer.DeployBlockhashStore(suiteSetup.Wallets.Default())
			Expect(err).ShouldNot(HaveOccurred())
			coordinator, err := suiteSetup.Deployer.DeployVRFCoordinator(suiteSetup.Wallets.Default(), suiteSetup.Link.Address(), bhs.Address())
			Expect(err).ShouldNot(HaveOccurred())
			_, err = suiteSetup.Deployer.DeployVRFConsumer(suiteSetup.Wallets.Default(), suiteSetup.Link.Address(), coordinator.Address())
			Expect(err).ShouldNot(HaveOccurred())
			_, err = suiteSetup.Deployer.DeployVRFContract(suiteSetup.Wallets.Default())
			Expect(err).ShouldNot(HaveOccurred())
			err = suiteSetup.Client.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred())
		})
		By("deploying direct request contract", func() {
			_, err := suiteSetup.Deployer.DeployOracle(suiteSetup.Wallets.Default(), suiteSetup.Link.Address())
			Expect(err).ShouldNot(HaveOccurred())
			_, err = suiteSetup.Deployer.DeployAPIConsumer(suiteSetup.Wallets.Default(), suiteSetup.Link.Address())
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	AfterEach(func() {
		By("Printing gas stats", func() {
			suiteSetup.Client.GasStats().PrintStats()
		})
		By("Tearing down the environment", suiteSetup.TearDown())
	})
})
