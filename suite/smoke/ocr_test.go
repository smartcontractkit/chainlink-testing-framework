package smoke

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/smartcontractkit/helmenv/environment"
	"github.com/smartcontractkit/helmenv/tools"
	"github.com/smartcontractkit/integrations-framework/actions"
	"github.com/smartcontractkit/integrations-framework/contracts"
	"math/big"
	"time"
)

var _ = FDescribe("OCR Feed @ocr", func() {
	var (
		err      error
		e        *environment.Environment
		ocrSetup *actions.OCRSetup
	)
	BeforeEach(func() {
		By("Deploying the environment", func() {
			e, err = environment.DeployOrLoadEnvironment(
				environment.NewChainlinkConfig(environment.ChainlinkReplicas(6, nil)),
				tools.ChartsRoot,
			)
			Expect(err).ShouldNot(HaveOccurred())
			err = e.ConnectAll()
			Expect(err).ShouldNot(HaveOccurred())
		})
		By("Getting the clients", func() {
			ocrSetup, err = actions.NewOCRSetup(e, []string{"chainlink"})
			Expect(err).ShouldNot(HaveOccurred())
		})
		By("Funding Chainlink nodes", func() {
			err = ocrSetup.FundNodes()
			Expect(err).ShouldNot(HaveOccurred())
		})
		By("Deploying OCR contracts", func() {
			err = ocrSetup.DeployOCRContracts(1)
			Expect(err).ShouldNot(HaveOccurred())
		})
		By("Creating OCR jobs", func() {
			err = ocrSetup.CreateOCRJobs()
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Describe("with OCR job", func() {
		It("performs two rounds", func() {
			ocrRoundTimeout := 2 * time.Minute
			ocrRound := contracts.NewOffchainAggregatorRoundConfirmer(ocrSetup.OCRInstances[0], big.NewInt(1), ocrRoundTimeout)
			ocrSetup.Networks.Default.AddHeaderEventSubscription(ocrSetup.OCRInstances[0].Address(), ocrRound)
			err = ocrSetup.SetAdapterResults([]int{5, 5, 5, 5, 5})
			Expect(err).ShouldNot(HaveOccurred())

			answer, err := ocrSetup.OCRInstances[0].GetLatestAnswer(context.Background())
			Expect(err).ShouldNot(HaveOccurred())
			Expect(answer.Int64()).Should(Equal(int64(5)), "latest answer from OCR is not as expected")

			err = ocrSetup.SetAdapterResults([]int{10, 10, 10, 10, 10})
			Expect(err).ShouldNot(HaveOccurred())

			err = ocrSetup.OCRInstances[0].RequestNewRound()
			Expect(err).ShouldNot(HaveOccurred())
			ocrRound2 := contracts.NewOffchainAggregatorRoundConfirmer(ocrSetup.OCRInstances[0], big.NewInt(2), ocrRoundTimeout)
			ocrSetup.Networks.Default.AddHeaderEventSubscription(ocrSetup.OCRInstances[0].Address(), ocrRound2)
			err = ocrSetup.Networks.Default.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred())

			answer2, err := ocrSetup.OCRInstances[0].GetLatestAnswer(context.Background())
			Expect(err).ShouldNot(HaveOccurred())
			Expect(answer2.Int64()).Should(Equal(int64(10)), "latest answer from OCR is not as expected")

			//err = ocrSetup.SetAdapterResults([]int{5, 5, 5, 5, 5})
			//Expect(err).ShouldNot(HaveOccurred())
			//
			//err = ocrSetup.StartNewRound(1)
			//
			//answer, err := ocrSetup.OCRInstances[0].GetLatestAnswer(context.Background())
			//Expect(err).ShouldNot(HaveOccurred())
			//Expect(answer.Int64()).Should(Equal(int64(5)), "latest answer from OCR is not as expected")
			//
			//err = ocrSetup.SetAdapterResults([]int{10, 10, 10, 10, 10})
			//
			//err = ocrSetup.StartNewRound(2)
			//
			//answer, err = ocrSetup.OCRInstances[0].GetLatestAnswer(context.Background())
			//Expect(err).ShouldNot(HaveOccurred())
			//Expect(answer.Int64()).Should(Equal(int64(10)), "latest answer from OCR is not as expected")
		})
	})

	AfterEach(func() {
		By("Printing gas stats", func() {
			ocrSetup.Networks.Default.GasStats().PrintStats()
		})
		By("Tearing down the environment", func() {
			err = actions.TeardownSuite(e, ocrSetup.Networks)
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
})
