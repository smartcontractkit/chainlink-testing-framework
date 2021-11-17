package smoke

import (
	"context"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/helmenv/environment"
	"github.com/smartcontractkit/integrations-framework/actions"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/contracts"
	"github.com/smartcontractkit/integrations-framework/utils"
	"math/big"
	"path/filepath"
	"time"
)

var _ = XDescribe("OCR Feed @ocr", func() {
	var (
		err             error
		nets            *client.Networks
		cd              contracts.ContractDeployer
		lt              contracts.LinkToken
		ocr             contracts.OffchainAggregator
		ocrRoundTimeout = 2 * time.Minute
		cls             []client.Chainlink
		adapterPath     string
		mockserver      *client.MockserverClient
		e               *environment.Environment
	)
	BeforeEach(func() {
		By("Deploying the environment", func() {
			e, err = environment.NewEnvironmentFromPreset(filepath.Join(utils.PresetRoot, "chainlink-cluster-6"))
			Expect(err).ShouldNot(HaveOccurred())
			err = e.Connect()
			Expect(err).ShouldNot(HaveOccurred())
		})
		By("Getting the clients", func() {
			nets, err = client.NewNetworks(e, nil)
			Expect(err).ShouldNot(HaveOccurred())
			cd, err = contracts.NewContractDeployer(nets.Default)
			Expect(err).ShouldNot(HaveOccurred())
			cls, err = client.NewChainlinkClients(e)
			Expect(err).ShouldNot(HaveOccurred())
			mockserver, err = client.NewMockServerClientFromEnv(e)
			Expect(err).ShouldNot(HaveOccurred())
			nets.Default.ParallelTransactions(true)
		})
		By("Funding Chainlink nodes", func() {
			txCost, err := nets.Default.CalculateTXSCost(200)
			Expect(err).ShouldNot(HaveOccurred())
			err = actions.FundChainlinkNodes(cls, nets.Default, txCost)
			Expect(err).ShouldNot(HaveOccurred())
		})
		By("Deploying OCR contracts", func() {
			lt, err = cd.DeployLinkTokenContract()
			Expect(err).ShouldNot(HaveOccurred())
			ocr, err = cd.DeployOffChainAggregator(lt.Address(), contracts.DefaultOffChainAggregatorOptions())
			Expect(err).ShouldNot(HaveOccurred())
			err = ocr.SetConfig(
				cls[1:],
				contracts.DefaultOffChainAggregatorConfig(len(cls[1:])),
			)
			Expect(err).ShouldNot(HaveOccurred())
			err = lt.Transfer(ocr.Address(), big.NewInt(2e18))
			Expect(err).ShouldNot(HaveOccurred())
			err = nets.Default.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred())
		})
		By("Creating OCR jobs", func() {
			bootstrapNode := cls[0]
			bootstrapP2PIds, err := bootstrapNode.ReadP2PKeys()
			Expect(err).ShouldNot(HaveOccurred())
			bootstrapP2PId := bootstrapP2PIds.Data[0].Attributes.PeerID
			bootstrapSpec := &client.OCRBootstrapJobSpec{
				Name:            fmt.Sprintf("bootstrap-%s", uuid.NewV4().String()),
				ContractAddress: ocr.Address(),
				P2PPeerID:       bootstrapP2PId,
				IsBootstrapPeer: true,
			}
			_, err = bootstrapNode.CreateJob(bootstrapSpec)
			Expect(err).ShouldNot(HaveOccurred())

			uuid := uuid.NewV4().String()
			adapterPath = fmt.Sprintf("/variable_%s", uuid)
			adapterURL := fmt.Sprintf("%s%s", mockserver.Config.ClusterURL, adapterPath)

			for nodeIndex := 1; nodeIndex < len(cls); nodeIndex++ {
				nodeP2PIds, err := cls[nodeIndex].ReadP2PKeys()
				Expect(err).ShouldNot(HaveOccurred())
				nodeP2PId := nodeP2PIds.Data[0].Attributes.PeerID
				nodeTransmitterAddress, err := cls[nodeIndex].PrimaryEthAddress()
				Expect(err).ShouldNot(HaveOccurred())
				nodeOCRKeys, err := cls[nodeIndex].ReadOCRKeys()
				Expect(err).ShouldNot(HaveOccurred())
				nodeOCRKeyId := nodeOCRKeys.Data[0].ID

				bta := client.BridgeTypeAttributes{
					Name: fmt.Sprintf("variable_%s", uuid),
					URL:  adapterURL,
				}

				err = cls[nodeIndex].CreateBridge(&bta)
				Expect(err).ShouldNot(HaveOccurred())

				ocrSpec := &client.OCRTaskJobSpec{
					ContractAddress:    ocr.Address(),
					P2PPeerID:          nodeP2PId,
					P2PBootstrapPeers:  []client.Chainlink{bootstrapNode},
					KeyBundleID:        nodeOCRKeyId,
					TransmitterAddress: nodeTransmitterAddress,
					ObservationSource:  client.ObservationSourceSpecBridge(bta),
				}
				_, err = cls[nodeIndex].CreateJob(ocrSpec)
				Expect(err).ShouldNot(HaveOccurred())
			}
		})
	})

	Describe("with OCR job", func() {
		It("performs two rounds", func() {
			ocrRound := contracts.NewOffchainAggregatorRoundConfirmer(ocr, big.NewInt(1), ocrRoundTimeout)
			nets.Default.AddHeaderEventSubscription(ocr.Address(), ocrRound)
			err = mockserver.SetValuePath(adapterPath, 5)
			Expect(err).ShouldNot(HaveOccurred())
			err = ocr.RequestNewRound()
			Expect(err).ShouldNot(HaveOccurred())
			err = nets.Default.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred())

			answer, err := ocr.GetLatestAnswer(context.Background())
			Expect(err).ShouldNot(HaveOccurred())
			Expect(answer.Int64()).Should(Equal(int64(5)), "latest answer from OCR is not as expected")

			err = mockserver.SetValuePath(adapterPath, 10)
			Expect(err).ShouldNot(HaveOccurred())

			err = ocr.RequestNewRound()
			Expect(err).ShouldNot(HaveOccurred())
			ocrRound2 := contracts.NewOffchainAggregatorRoundConfirmer(ocr, big.NewInt(2), ocrRoundTimeout)
			nets.Default.AddHeaderEventSubscription(ocr.Address(), ocrRound2)
			err = nets.Default.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred())

			answer2, err := ocr.GetLatestAnswer(context.Background())
			Expect(err).ShouldNot(HaveOccurred())
			Expect(answer2.Int64()).Should(Equal(int64(10)), "latest answer from OCR is not as expected")
		})
	})

	AfterEach(func() {
		By("Printing gas stats", func() {
			nets.Default.GasStats().PrintStats()
		})
		By("Tearing down the environment", func() {
			err = actions.TeardownSuite(e, nets)
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
})
