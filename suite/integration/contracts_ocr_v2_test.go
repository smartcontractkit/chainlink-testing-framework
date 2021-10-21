package integration

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/smartcontractkit/integrations-framework/actions"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/contracts"
	"github.com/smartcontractkit/integrations-framework/environment"
	"github.com/smartcontractkit/integrations-framework/tools"
	"math/big"
)

var _ = Describe("OCR v2 @v2ocr", func() {
	var (
		suiteSetup  actions.SuiteSetup
		networkInfo actions.NetworkInfo
		//nodes         []client.Chainlink
		//nodeAddresses []common.Address
		wallets client.BlockchainWallets
		link    contracts.LinkToken
		err     error
	)

	BeforeEach(func() {
		By("Deploying the environment", func() {
			suiteSetup, err = actions.SingleNetworkSetup(
				environment.NewChainlinkHeadlessCluster(0),
				client.NewNetworkFromConfigWithDefault(client.NetworkTypeTerraLocal),
				tools.ProjectRoot,
			)
			Expect(err).ShouldNot(HaveOccurred())
			//nodes, err = environment.GetChainlinkClients(suiteSetup.Environment())
			//Expect(err).ShouldNot(HaveOccurred())
			//nodeAddresses, err = actions.ChainlinkNodeAddresses(nodes)
			//Expect(err).ShouldNot(HaveOccurred())
			//log.Debug().Interface("Addresses", nodeAddresses).Msg("Addresses")

			//err = suiteSetup.Environment().DeploySpecs(environment.NewRelays(1))
			//Expect(err).ShouldNot(HaveOccurred())
			networkInfo = suiteSetup.DefaultNetwork()
			wallets = networkInfo.Wallets
			link = networkInfo.Link
		})

		By("Setting up OCR contracts", func() {
			bac, err := networkInfo.Deployer.DeployOCRv2AccessController(wallets.Default())
			Expect(err).ShouldNot(HaveOccurred())
			rac, err := networkInfo.Deployer.DeployOCRv2AccessController(wallets.Default())
			Expect(err).ShouldNot(HaveOccurred())
			ocr2, err := networkInfo.Deployer.DeployOCRv2(wallets.Default(), bac.Address(), rac.Address(), link.Address())
			Expect(err).ShouldNot(HaveOccurred())
			err = link.Transfer(wallets.Default(), ocr2.Address(), big.NewInt(100))
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Describe("with OCRv2 job", func() {
		It("performs two rounds and has withdrawable payments for oracles", func() {

		})
	})

	AfterEach(func() {
		By("Tearing down the environment", suiteSetup.TearDown())
	})
})
