package contracts

import (
	"context"
	"fmt"
	"math/big"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/smartcontractkit/integrations-framework/actions"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/contracts"
	"github.com/smartcontractkit/integrations-framework/environment"
)

var _ = Describe("OCR Feed", func() {
	var (
		suiteSetup     *actions.DefaultSuiteSetup
		chainlinkNodes []client.Chainlink
		adapter        environment.ExternalAdapter
		defaultWallet  client.BlockchainWallet
	)

	BeforeEach(func() {
		By("Deploying the environment", func() {
			var err error
			suiteSetup, err = actions.DefaultLocalSetup(
				environment.NewChainlinkCluster(5),
				client.NewNetworkFromConfig,
			)
			Expect(err).ShouldNot(HaveOccurred())
			adapter, err = environment.GetExternalAdapter(suiteSetup.Env)
			Expect(err).ShouldNot(HaveOccurred())
			chainlinkNodes, err = environment.GetChainlinkClients(suiteSetup.Env)
			Expect(err).ShouldNot(HaveOccurred())
			defaultWallet = suiteSetup.Wallets.Default()
			suiteSetup.Client.ParallelTransactions(true)
		})
	})

	It("Deploys an OCR feed", func() {
		var ocrInstance contracts.OffchainAggregator

		By("Funding nodes and deploying OCR contract", func() {
			err := actions.FundChainlinkNodes(
				chainlinkNodes,
				suiteSetup.Client,
				defaultWallet,
				big.NewFloat(2),
				big.NewFloat(2),
			)
			Expect(err).ShouldNot(HaveOccurred())

			// Deploy and config OCR contract
			deployer, err := contracts.NewContractDeployer(suiteSetup.Client)
			Expect(err).ShouldNot(HaveOccurred())

			ocrInstance, err = deployer.DeployOffChainAggregator(defaultWallet, contracts.DefaultOffChainAggregatorOptions())
			Expect(err).ShouldNot(HaveOccurred())
			err = ocrInstance.SetConfig(
				defaultWallet,
				chainlinkNodes,
				contracts.DefaultOffChainAggregatorConfig(len(chainlinkNodes)),
			)
			Expect(err).ShouldNot(HaveOccurred())
			err = ocrInstance.Fund(defaultWallet, nil, big.NewFloat(2))
			Expect(err).ShouldNot(HaveOccurred())
			err = suiteSetup.Client.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred())
		})

		By("Sending OCR jobs to chainlink nodes", func() {
			// Initialize bootstrap node
			bootstrapNode := chainlinkNodes[0]
			bootstrapP2PIds, err := bootstrapNode.ReadP2PKeys()
			Expect(err).ShouldNot(HaveOccurred())
			bootstrapP2PId := bootstrapP2PIds.Data[0].Attributes.PeerID
			bootstrapSpec := &client.OCRBootstrapJobSpec{
				ContractAddress: ocrInstance.Address(),
				P2PPeerID:       bootstrapP2PId,
				IsBootstrapPeer: true,
			}
			_, err = bootstrapNode.CreateJob(bootstrapSpec)
			Expect(err).ShouldNot(HaveOccurred())

			// Send OCR job to other nodes
			for index := 1; index < len(chainlinkNodes); index++ {
				nodeP2PIds, err := chainlinkNodes[index].ReadP2PKeys()
				Expect(err).ShouldNot(HaveOccurred())
				nodeP2PId := nodeP2PIds.Data[0].Attributes.PeerID
				nodeTransmitterAddress, err := chainlinkNodes[index].PrimaryEthAddress()
				Expect(err).ShouldNot(HaveOccurred())
				nodeOCRKeys, err := chainlinkNodes[index].ReadOCRKeys()
				Expect(err).ShouldNot(HaveOccurred())
				nodeOCRKeyId := nodeOCRKeys.Data[0].ID

				ocrSpec := &client.OCRTaskJobSpec{
					ContractAddress:    ocrInstance.Address(),
					P2PPeerID:          nodeP2PId,
					P2PBootstrapPeers:  []client.Chainlink{bootstrapNode},
					KeyBundleID:        nodeOCRKeyId,
					TransmitterAddress: nodeTransmitterAddress,
					ObservationSource:  client.ObservationSourceSpec(fmt.Sprintf("%s/variable", adapter.ClusterURL())),
				}
				_, err = chainlinkNodes[index].CreateJob(ocrSpec)
				Expect(err).ShouldNot(HaveOccurred())
			}
		})

		By("Checking OCR rounds", func() {
			roundTimeout := time.Minute * 2
			err := adapter.SetVariable(5)
			Expect(err).ShouldNot(HaveOccurred())
			err = ocrInstance.RequestNewRound(defaultWallet)
			Expect(err).ShouldNot(HaveOccurred())
			err = suiteSetup.Client.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred())

			// Wait for the first round
			ocrRound := contracts.NewOffchainAggregatorRoundConfirmer(ocrInstance, big.NewInt(1), roundTimeout)
			suiteSetup.Client.AddHeaderEventSubscription(ocrInstance.Address(), ocrRound)
			err = suiteSetup.Client.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred())

			// Check answer is as expected
			answer, err := ocrInstance.GetLatestAnswer(context.Background())
			Expect(err).ShouldNot(HaveOccurred())
			Expect(answer.Int64()).Should(Equal(int64(5)))

			// Change adapter answer
			err = adapter.SetVariable(10)
			Expect(err).ShouldNot(HaveOccurred())

			// Wait for the second round
			ocrRound = contracts.NewOffchainAggregatorRoundConfirmer(ocrInstance, big.NewInt(2), roundTimeout)
			suiteSetup.Client.AddHeaderEventSubscription(ocrInstance.Address(), ocrRound)
			err = suiteSetup.Client.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred())

			// Check answer is as expected
			answer, err = ocrInstance.GetLatestAnswer(context.Background())
			Expect(err).ShouldNot(HaveOccurred())
			Expect(answer.Int64()).Should(Equal(int64(10)))
		})
	})

	AfterEach(func() {
		By("Tearing down the environment", suiteSetup.TearDown())
	})
})
