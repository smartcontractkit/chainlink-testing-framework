package performance

//revive:disable:dot-imports
import (
	"math/big"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/smartcontractkit/helmenv/environment"
	"github.com/smartcontractkit/helmenv/tools"
	"github.com/smartcontractkit/integrations-framework/actions"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/contracts"
	"github.com/smartcontractkit/integrations-framework/testsetups"
	"github.com/smartcontractkit/integrations-framework/utils"
)

var _ = Describe("Keeper performance suite @performance-keeper", func() {
	var (
		err                 error
		networks            *client.Networks
		contractDeployer    contracts.ContractDeployer
		linkToken           contracts.LinkToken
		chainlinkNodes      []client.Chainlink
		env                 *environment.Environment
		keeperBlockTimeTest *testsetups.KeeperBlockTimeTest
	)

	BeforeEach(func() {
		By("Deploying the environment", func() {
			env, err = environment.DeployOrLoadEnvironment(
				environment.NewChainlinkConfig(environment.ChainlinkReplicas(6, nil), ""),
				tools.ChartsRoot,
			)
			Expect(err).ShouldNot(HaveOccurred(), "Environment deployment shouldn't fail")
			err = env.ConnectAll()
			Expect(err).ShouldNot(HaveOccurred(), "Connecting to all nodes shouldn't fail")
		})

		By("Connecting to launched resources", func() {
			networkRegistry := client.NewDefaultNetworkRegistry()
			networks, err = networkRegistry.GetNetworks(env)
			Expect(err).ShouldNot(HaveOccurred(), "Connecting to blockchain nodes shouldn't fail")
			contractDeployer, err = contracts.NewContractDeployer(networks.Default)
			Expect(err).ShouldNot(HaveOccurred(), "Deploying contracts shouldn't fail")
			chainlinkNodes, err = client.ConnectChainlinkNodes(env)
			Expect(err).ShouldNot(HaveOccurred(), "Connecting to chainlink nodes shouldn't fail")
			networks.Default.ParallelTransactions(true)
		})

		By("Funding Chainlink nodes", func() {
			txCost, err := networks.Default.EstimateCostForChainlinkOperations(10)
			Expect(err).ShouldNot(HaveOccurred(), "Estimating cost for Chainlink Operations shouldn't fail")
			err = actions.FundChainlinkNodes(chainlinkNodes, networks.Default, txCost)
			Expect(err).ShouldNot(HaveOccurred(), "Funding Chainlink nodes shouldn't fail")
			// Edge case where simulated networks need some funds at the 0x0 address in order for keeper reads to work
			if networks.Default.GetNetworkType() == "eth_simulated" {
				err = actions.FundAddresses(networks.Default, big.NewFloat(1), "0x0")
				Expect(err).ShouldNot(HaveOccurred())
			}
		})

		By("Setup the Keeper test", func() {
			linkToken, err = contractDeployer.DeployLinkTokenContract()
			Expect(err).ShouldNot(HaveOccurred(), "Deploying Link Token Contract shouldn't fail")
			keeperBlockTimeTest = testsetups.NewKeeperBlockTimeTest(
				testsetups.KeeperBlockTimeTestInputs{
					NumberOfContracts: 20,
					BlockRange:        1000,
					BlockInterval:     50,
					ContractDeployer:  contractDeployer,
					ChainlinkNodes:    chainlinkNodes,
					Networks:          networks,
					LinkTokenContract: linkToken,
				},
			)
			keeperBlockTimeTest.Setup()
		})
	})

	Describe("Watching the keeper contracts to ensure they reply in time", func() {
		It("Watches for Upkeep counts", func() {
			keeperBlockTimeTest.Run()
		})
	})

	AfterEach(func() {
		By("Printing gas stats", func() {
			networks.Default.GasStats().PrintStats()
		})
		By("Tearing down the environment", func() {
			err = actions.TeardownSuite(env, networks, utils.ProjectRoot, &keeperBlockTimeTest.TestReporter)
			Expect(err).ShouldNot(HaveOccurred(), "Environment teardown shouldn't fail")
		})
	})
})
