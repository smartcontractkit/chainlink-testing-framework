package performance

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/smartcontractkit/integrations-framework/actions"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/contracts"
	"github.com/smartcontractkit/integrations-framework/environment"
	"github.com/smartcontractkit/integrations-framework/tools"
	"math/big"
	"time"
)

var _ = Describe("OCR soak test @soak-ocr", func() {
	var (
		s        *actions.DefaultSuiteSetup
		nodes    []client.Chainlink
		adapter  environment.ExternalAdapter
		perfTest Test
		err      error
	)

	BeforeEach(func() {
		By("Deploying the environment", func() {
			s, err = actions.DefaultLocalSetup(
				"ocr-soak",
				environment.NewChainlinkCluster(5),
				client.NewNetworkFromConfig,
				tools.ProjectRoot,
			)
			Expect(err).ShouldNot(HaveOccurred())
			adapter, err = environment.GetExternalAdapter(s.Env)
			Expect(err).ShouldNot(HaveOccurred())
			nodes, err = environment.GetChainlinkClients(s.Env)
			Expect(err).ShouldNot(HaveOccurred())
			s.Client.ParallelTransactions(true)
		})

		By("Funding the Chainlink nodes", func() {
			err := actions.FundChainlinkNodes(
				nodes,
				s.Client,
				s.Wallets.Default(),
				big.NewFloat(10),
				big.NewFloat(10),
			)
			Expect(err).ShouldNot(HaveOccurred())
		})

		By("Setting up the OCR soak test", func() {
			perfTest = NewOCRTest(
				OCRTestOptions{
					TestOptions: TestOptions{
						NumberOfContracts: 5,
					},
					RoundTimeout: 180 * time.Second,
					AdapterValue: 5,
					TestDuration: 10 * time.Minute,
				},
				contracts.DefaultOffChainAggregatorOptions(),
				s.Env,
				s.Client,
				s.Wallets,
				s.Deployer,
				adapter,
			)
			err = perfTest.Setup()
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Describe("OCR Soak test", func() {
		Measure("Measure OCR rounds", func(_ Benchmarker) {
			err = perfTest.Run()
			Expect(err).ShouldNot(HaveOccurred())
		}, 1)
	})

	AfterEach(func() {
		By("Tearing down the environment", s.TearDown())
	})
})
