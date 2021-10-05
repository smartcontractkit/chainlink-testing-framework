package performance

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/smartcontractkit/integrations-framework/actions"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/environment"
	"github.com/smartcontractkit/integrations-framework/tools"
	"math/big"
	"time"
)

var _ = Describe("VRF soak test @soak-vrf", func() {
	var (
		suiteSetup *actions.DefaultSuiteSetup
		nodes      []client.Chainlink
		adapter    environment.ExternalAdapter
		perfTest   Test
		err        error
	)

	BeforeEach(func() {
		By("Deploying the environment", func() {
			suiteSetup, err = actions.DefaultLocalSetup(
				"vrf-soak",
				// more than one node is useless for VRF, because nodes are not cooperating for randomness
				environment.NewChainlinkCluster(1),
				client.NewNetworkFromConfig,
				tools.ProjectRoot,
			)
			Expect(err).ShouldNot(HaveOccurred())
			adapter, err = environment.GetExternalAdapter(suiteSetup.Env)
			Expect(err).ShouldNot(HaveOccurred())
			nodes, err = environment.GetChainlinkClients(suiteSetup.Env)
			Expect(err).ShouldNot(HaveOccurred())
			suiteSetup.Client.ParallelTransactions(true)
		})

		By("Funding the Chainlink nodes", func() {
			err := actions.FundChainlinkNodes(
				nodes,
				suiteSetup.Client,
				suiteSetup.Wallets.Default(),
				big.NewFloat(10),
				big.NewFloat(10),
			)
			Expect(err).ShouldNot(HaveOccurred())
		})

		By("Setting up the VRF soak test", func() {
			perfTest = NewVRFTest(
				VRFTestOptions{
					TestOptions: TestOptions{
						NumberOfContracts: 30,
					},
					RoundTimeout: 60 * time.Second,
					TestDuration: 3 * time.Minute,
				},
				suiteSetup.Env,
				suiteSetup.Link,
				suiteSetup.Client,
				suiteSetup.Wallets,
				suiteSetup.Deployer,
				adapter,
			)
			err = perfTest.Setup()
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Describe("VRF Soak test", func() {
		Measure("Measure VRF rounds", func(_ Benchmarker) {
			err = perfTest.Run()
			Expect(err).ShouldNot(HaveOccurred())
		}, 1)
	})

	AfterEach(func() {
		By("Tearing down the environment", suiteSetup.TearDown())
	})
})
