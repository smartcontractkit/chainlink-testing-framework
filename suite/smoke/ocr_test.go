package smoke

import (
	"context"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/smartcontractkit/helmenv/environment"
	"github.com/smartcontractkit/helmenv/tools"
	"github.com/smartcontractkit/integrations-framework/actions"
)

var _ = Describe("OCR Feed @ocr", func() {
	var (
		err      error
		env      *environment.Environment
		ocrSetup *actions.OCRSetup
	)
	BeforeEach(func() {
		By("Deploying the environment", func() {
			env, err = environment.DeployOrLoadEnvironment(
				environment.NewChainlinkConfig(environment.ChainlinkReplicas(6, nil)),
				tools.ChartsRoot,
			)
			Expect(err).ShouldNot(HaveOccurred())
			err = env.ConnectAll()
			Expect(err).ShouldNot(HaveOccurred())
		})
		By("Getting the clients", func() {
			ocrSetup, err = actions.NewOCRSetup(env, []string{"chainlink"})
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
			err = ocrSetup.SetAdapterResults([]int{5, 5, 5, 5, 5})
			Expect(err).ShouldNot(HaveOccurred())

			err = ocrSetup.StartNewRound(1)

			answer, err := ocrSetup.OCRInstances[0].GetLatestAnswer(context.Background())
			Expect(err).ShouldNot(HaveOccurred())
			Expect(answer.Int64()).Should(Equal(int64(5)), "latest answer from OCR is not as expected")

			err = ocrSetup.SetAdapterResults([]int{10, 10, 10, 10, 10})

			err = ocrSetup.StartNewRound(2)

			answer, err = ocrSetup.OCRInstances[0].GetLatestAnswer(context.Background())
			Expect(err).ShouldNot(HaveOccurred())
			Expect(answer.Int64()).Should(Equal(int64(10)), "latest answer from OCR is not as expected")
		})
	})

	AfterEach(func() {
		By("Printing gas stats", func() {
			ocrSetup.Networks.Default.GasStats().PrintStats()
		})
		By("Tearing down the environment", func() {
			err = actions.TeardownSuite(env, ocrSetup.Networks)
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
})
