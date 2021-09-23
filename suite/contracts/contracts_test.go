package contracts

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/smartcontractkit/integrations-framework/actions"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/contracts"
	"github.com/smartcontractkit/integrations-framework/environment"
	"github.com/smartcontractkit/integrations-framework/tools"
	"math/big"
)

var _ = Describe("Basic Contract Interactions @contract", func() {
	var s *actions.DefaultSuiteSetup
	var defaultWallet client.BlockchainWallet

	BeforeEach(func() {
		By("Deploying the environment", func() {
			var err error
			s, err = actions.DefaultLocalSetup(
				"basic-chainlink",
				environment.NewChainlinkCluster(0),
				client.NewNetworkFromConfig,
				tools.ProjectRoot,
			)
			Expect(err).ShouldNot(HaveOccurred())
			defaultWallet = s.Wallets.Default()
		})
	})

	It("can deploy all contracts", func() {
		By("basic interaction with a storage contract", func() {
			storeInstance, err := s.Deployer.DeployStorageContract(defaultWallet)
			Expect(err).ShouldNot(HaveOccurred())
			testVal := big.NewInt(5)
			err = storeInstance.Set(testVal)
			Expect(err).ShouldNot(HaveOccurred())
			val, err := storeInstance.Get(context.Background())
			Expect(err).ShouldNot(HaveOccurred())
			Expect(val).To(Equal(testVal))
		})

		By("deploying the flux monitor contract", func() {
			rac, err := s.Deployer.DeployReadAccessController(s.Wallets.Default())
			Expect(err).ShouldNot(HaveOccurred())
			flags, err := s.Deployer.DeployFlags(s.Wallets.Default(), rac.Address())
			Expect(err).ShouldNot(HaveOccurred())
			_, err = s.Deployer.DeployDeviationFlaggingValidator(s.Wallets.Default(), flags.Address(), big.NewInt(0))
			Expect(err).ShouldNot(HaveOccurred())
			fluxOptions := contracts.DefaultFluxAggregatorOptions()
			_, err = s.Deployer.DeployFluxAggregatorContract(defaultWallet, fluxOptions)
			Expect(err).ShouldNot(HaveOccurred())
		})

		By("deploying the ocr contract", func() {
			ocrOptions := contracts.DefaultOffChainAggregatorOptions()
			_, err := s.Deployer.DeployOffChainAggregator(defaultWallet, ocrOptions)
			Expect(err).ShouldNot(HaveOccurred())
		})
		By("deploying keeper contracts", func() {
			ef, err := s.Deployer.DeployMockETHLINKFeed(s.Wallets.Default(), big.NewInt(2e18))
			Expect(err).ShouldNot(HaveOccurred())
			gf, err := s.Deployer.DeployMockGasFeed(s.Wallets.Default(), big.NewInt(2e11))
			Expect(err).ShouldNot(HaveOccurred())
			_, err = s.Deployer.DeployKeeperRegistry(
				s.Wallets.Default(),
				&contracts.KeeperRegistryOpts{
					LinkAddr:             s.Link.Address(),
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
			err = s.Client.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred())
		})
		By("deploying vrf contract", func() {
			bhs, err := s.Deployer.DeployBlockhashStore(s.Wallets.Default())
			Expect(err).ShouldNot(HaveOccurred())
			coordinator, err := s.Deployer.DeployVRFCoordinator(s.Wallets.Default(), s.Link.Address(), bhs.Address())
			Expect(err).ShouldNot(HaveOccurred())
			_, err = s.Deployer.DeployVRFConsumer(s.Wallets.Default(), s.Link.Address(), coordinator.Address())
			Expect(err).ShouldNot(HaveOccurred())
			_, err = s.Deployer.DeployVRFContract(s.Wallets.Default())
			Expect(err).ShouldNot(HaveOccurred())
			err = s.Client.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred())
		})
		By("deploying direct request contract", func() {
			_, err := s.Deployer.DeployOracle(s.Wallets.Default(), s.Link.Address())
			Expect(err).ShouldNot(HaveOccurred())
			_, err = s.Deployer.DeployAPIConsumer(s.Wallets.Default(), s.Link.Address())
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	AfterEach(func() {
		By("Printing gas stats", func() {
			s.Client.GasStats().PrintStats()
		})
		By("Tearing down the environment", s.TearDown())
	})
})
